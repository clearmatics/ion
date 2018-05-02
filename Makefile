SOLC=solc --optimize
PYTHON=python
GANACHE=./node_modules/.bin/ganache-cli
TRUFFLE=./node_modules/.bin/truffle
CONTRACTS=IonLock IonLink ERC223 Token HTLC
CONTRACTS_BIN=$(addprefix build/,$(addsuffix .bin,$(CONTRACTS)))
CONTRACTS_ABI=$(addprefix abi/,$(addsuffix .abi,$(CONTRACTS)))

PROTOCOLS=ion/proto/chain
PROTOCOLS_PY=$(addsuffix _pb2.py,$(PROTOCOLS))

PYLINT_IGNORE=C0330,invalid-name,line-too-long,missing-docstring,bad-whitespace,consider-using-ternary,wrong-import-position,wrong-import-order,trailing-whitespace

all: $(CONTRACTS_BIN) $(CONTRACTS_ABI) $(PROTOCOLS_PY) test truffle-test pylint

build:
	mkdir -p build
	npm install
	$(PYTHON) -mpip install -r requirements.txt

.PHONY: docs/deps-modules.dot
docs/deps-modules.dot:
	pydepgraph -p ion > $@

.PHONY: docs/deps-files.dot
docs/deps-files.dot:
	sfood -i -r ion | sfood-graph > $@

docker-build: dist/ion
	docker build --rm=true -t clearmatics/ion:latest -f Dockerfile.alpine-glibc .

docker-run:
	docker run --rm=true -ti clearmatics/ion:latest shell

lint:
	$(PYTHON) -mpylint ion/

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
	npm run testrpca

test-js:
	npm run test

test-unit:
	$(PYTHON) -m unittest discover test/

test: test-unit test-js
