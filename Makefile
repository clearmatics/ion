SOLC=solc --optimize
PYTHON=python
GANACHE=./node_modules/.bin/ganache-cli
TRUFFLE=./node_modules/.bin/truffle

build:
	mkdir -p build
	npm install
	$(PYTHON) -mpip install -r requirements.txt

docker-build: dist/ion
	docker build --rm=true -t clearmatics/ion:latest -f Dockerfile.alpine-glibc .

docker-run:
	docker run --rm=true -ti clearmatics/ion:latest shell

python-lint:
	$(PYTHON) -mpylint ion/

solidity-lint:
	npm run lint

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
