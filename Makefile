SOLC=solc --optimize
PYTHON=python
GANACHE=./node_modules/.bin/ganache-cli
TRUFFLE=./node_modules/.bin/truffle
CONTRACTS=Sodium Fluoride IonLock IonLink ERC223 Token HTLC
CONTRACTS_BIN=$(addprefix build/,$(addsuffix .bin,$(CONTRACTS)))
CONTRACTS_ABI=$(addprefix abi/,$(addsuffix .abi,$(CONTRACTS)))

PROTOCOLS=ion/proto/chain
PROTOCOLS_PY=$(addsuffix _pb2.py,$(PROTOCOLS))

PYLINT_IGNORE=C0330,invalid-name,line-too-long,missing-docstring,bad-whitespace,consider-using-ternary,wrong-import-position,wrong-import-order,trailing-whitespace

all: $(CONTRACTS_BIN) $(CONTRACTS_ABI) $(PROTOCOLS_PY) test truffle-test dist/ion lint README.pdf

README.pdf: README.md
	pandoc --toc --reference-links --number-sections --listings --template docs/eisvogel -f markdown -t latex -o $@ $<

build:
	mkdir -p build

pyflakes:
	$(PYTHON) -mpyflakes ion

pylint:
	$(PYTHON) -mpylint -d $(PYLINT_IGNORE) ion

lint: pyflakes pylint

dist:
	$(PYTHON) setup.py bdist_egg --exclude-source-files

dist/ion: dist
	$(PYTHON) -mPyInstaller ion.spec

dev-yarn:
	echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list
	curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add -
	apt-get update
	apt-get install yarn

dev-nodejs:
	curl -sL https://deb.nodesource.com/setup_8.x | bash -
	apt install nodejs
	ln -s /usr/bin/nodejs  /usr/bin/node

dev-python:
	$(PYTHON) -mpip install pylint pyflakes snakefood pycallgraph
	apt install protobuf-compiler

dev: dev-python dev-nodejs dev-yarn

docs/deps-modules.dot:
	pydepgraph -p ion > $@

docs/deps-files.dot:
	sfood -i -r ion | sfood-graph > $@

docker-build: dist/ion
	docker build --rm=true -t clearmatics/ion:latest -f Dockerfile.alpine-libc .

docker-run:
	docker run --rm=true -ti clearmatics/ion:latest

yarn:
	yarn

$(TRUFFLE): yarn
$(GANACHE): yarn

truffle-test: $(TRUFFLE)
	$(TRUFFLE) test

truffle-compile: $(TRUFFLE)
	$(TRUFFLE) compile

truffle-deploy:
	$(TRUFFLE) deploy

ion/proto/%_pb2.py: ion/proto/%.proto
	protoc -I. --python_out=. $<

requirements: requirements.txt
	$(PYTHON) -mpip install -r requirements.txt

abi:
	mkdir -p abi

abi/%.abi: build/%.abi abi
	cp $< $@

build/%.bin: contracts/%.sol
	$(SOLC) -o build --asm --bin --overwrite --abi $<

build/%.combined.bin: build/%.combined.sol
	$(SOLC) -o build --asm --bin --overwrite --abi $<

build/%.combined.sol: contracts/%.sol build
	cat $< | sed -e 's/\bimport\(\b.*\);/#include \1/g' | cpp -Icontracts | sed -e 's/^#.*$$//g' > $@

clean:
	rm -rf build chaindata dist
	find . -name '*.pyc' -exec rm '{}' ';'
	rm -rf *.pyc *.pdf *.egg-info

testrpc:
	yarn testrpc

test-genesis:
	rm -rf chaindata
	$(PYTHON) -mion.chain -g -r 10

test-client:
	$(PYTHON) -mion.client --inproc --test

test-merkle:
	$(PYTHON) -mion.merkle

test-payment:
	$(PYTHON) -mion.payment --random -b 0xed39af75a8367cad4689e3b4ffe7e189171eb33e32663c70cf503690dbc49d98 -v 1234 -j | $(PYTHON) -mion.payment -i /dev/stdin -b 0xed39af75a8367cad4689e3b4ffe7e189171eb33e32663c70cf503690dbc49d98 -m

test: test-genesis test-client test-merkle test-payment
