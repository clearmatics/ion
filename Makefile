ROOT_DIR := $(shell dirname $(realpath $(MAKEFILE_LIST)))

SOLC=$(ROOT_DIR)/node_modules/.bin/solcjs
PYTHON=python3
NPM=npm
GANACHE=$(ROOT_DIR)/node_modules/.bin/ganache-cli
TRUFFLE=$(ROOT_DIR)/node_modules/.bin/truffle

DOCKER_TAG_NAME=clearmatics/ion:latest

UTIL_IMPORTS=$(ROOT_DIR)/utils/extract-imports.sh

CONTRACTS=IonLock IonLink ERC223 Token HTLC
CONTRACTS_BIN=$(addprefix build/,$(addsuffix .bin,$(CONTRACTS)))
CONTRACTS_ABI=$(addprefix abi/,$(addsuffix .abi,$(CONTRACTS)))

PYLINT_IGNORE=C0330,invalid-name,line-too-long,missing-docstring,bad-whitespace,consider-using-ternary,wrong-import-position,wrong-import-order,trailing-whitespace


all: check-prereqs contracts python-pyflakes test python-pylint

check-prereqs:
	@if [ ! -f "$(SOLC)" ]; then \
		echo -e "Dependencies not found!\nInstall prerequisites first! See README.md"; \
		false; \
	fi

clean:
	rm -rf build chaindata dist
	find . -name '*.pyc' -exec rm '{}' ';'
	rm -rf *.pyc *.pdf *.egg-info *.pid *.log


#######################################################################
#
# Packaging and distribution

docker-build: dist/ion
	docker build --rm=true -t $(DOCKER_TAG_NAME) -f Dockerfile.alpine-glibc .

docker-run:
	docker run --rm=true -ti $(DOCKER_TAG_NAME) --help

bdist:
	$(PYTHON) setup.py bdist_egg --exclude-source-files
	$(PYTHON) setup.py bdist_wheel --universal

dist:
	mkdir -p $@

dist/ion: dist
	$(PYTHON) -mPyInstaller ion.spec


#######################################################################
#
# Linting and anti-retardery measures

python-pyflakes:
	$(PYTHON) -mpyflakes ion

python-pylint:
	$(PYTHON) -mpylint -d $(PYLINT_IGNORE) ion

python-lint: python-pyflakes python-pylint

solidity-lint:
	$(NPM) run lint


#######################################################################
#
# Install dependencies / requirements etc. for Python and NodeJS
#

nodejs-requirements:
	$(NPM) install

# Useful shortcut for development, install packages to user path by default
python-pip-user:
	mkdir -p $(HOME)/.pip/
	echo -e "[global]\nuser = true\n" > $(HOME)/.pip/pip.conf

python-requirements: requirements.txt
	$(PYTHON) -mpip install -r $<

python-dev-requirements: requirements-dev.txt
	$(PYTHON) -mpip install -r $<

requirements-dev: nodejs-requirements python-dev-requirements

requirements: python-requirements

fedora-dev:
	# use `nvm` to manage nodejs versions, rather than relying on system node
	curl https://raw.githubusercontent.com/creationix/nvm/master/install.sh | bash
	nvm install --lts


#######################################################################
#
# Builds Solidity contracts and ABI files
#

contracts: $(CONTRACTS_BIN) $(CONTRACTS_ABI)

abi:
	mkdir -p abi

abi/%.abi: build/%.abi abi contracts/%.sol
	cp $< $@

build:
	mkdir -p build

build/%.abi: build/%.bin

build/%.bin: contracts/%.sol build
	$(eval contract_name := $(shell echo $(shell basename $<) | cut -f 1 -d .))
	cd $(shell dirname $<) && $(SOLC) --optimize -o ../build --asm --bin --overwrite --abi $(shell basename $<) $(shell $(UTIL_IMPORTS) $<)
	cp build/$(contract_name)_sol_$(contract_name).bin build/$(contract_name).bin
	cp build/$(contract_name)_sol_$(contract_name).abi build/$(contract_name).abi

build/%.combined.bin: build/%.combined.sol
	$(SOLC) -o build --asm --bin --overwrite --abi $<

build/%.combined.sol: contracts/%.sol build
	cat $< | sed -e 's/\bimport\(\b.*\);/#include \1/g' | cpp -Icontracts | sed -e 's/^#.*$$//g' > $@


#######################################################################
#
# Testing and unit test harnesses
#

# runs an instance of testrpc in background, then waits for it to be ready
travis-testrpc-start: travis-testrpc-stop
	$(NPM) run testrpca > .testrpc.log & echo $$! > .testrpc.pid
	while true; do echo -n . ; curl http://localhost:8545 &> /dev/null && break || sleep 1; done

# Stops previ
travis-testrpc-stop:
	if [ -f .testrpc.pid ]; then kill `cat .testrpc.pid` || true; rm -f .testrpc.pid; fi

travis: travis-testrpc-start truffle-deploy-a contracts test


testrpc:
	$(NPM) run testrpca

testrpc-b:
	$(NPM) run testrpcb

test-js:
	$(NPM) run test

test-unit:
	$(PYTHON) -m unittest discover test/

test-coordserver:
	$(PYTHON) -mion htlc coordinator --contract 0xd833215cbcc3f914bd1c9ece3ee7bf8b14f841bb

test-coordclient:
	PYTHONPATH=. $(PYTHON) ./test/test_coordclient.py

test: test-unit test-js


#######################################################################
#
# Truffle utils
#

truffle-deploy:
	$(TRUFFLE) deploy

truffle-deploy-a:
	$(TRUFFLE) deploy --network testrpca --reset

truffle-deploy-b:
	$(TRUFFLE) deployb --network testrpcb --reset

truffle-console:
	$(TRUFFLE) console