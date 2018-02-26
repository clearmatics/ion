# Ion Interoperability Protocol

The Ion Interoperability Protocol provides mechanisms to perform atomic swaps and currency transfers across
multiple turing-complete blockchains. It is a sidechain and series of smart-contracts which 

## Getting Started

Building Ion requires many components such as Python, NodeJS, Yarn, Solidity compiler, Docker etc.
The `make dev` command will install all necessary dependencies required to develop and build Ion.

All Ion functionality is available via the built-in shell,

### Python

```commandline
python -mion
```

### Docker

The small Alpine Linux container is under 30mb and contains a self-contained PyInstaller executable.

```commandline
make docker-build
docker run -ti --rm clearmatics/ion:latest 
```

### PyInstaller

The PyInstaller `ion` executable is a self-contained Python bundle which includes all necessary dependencies,
this makes Ion easier to distribute.  

```commandline
make dist/ion
./dist/ion
```

## Tests

```
make testrpc &
make test
```
