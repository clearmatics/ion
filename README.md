# Ion Interoperability Protocol

The Ion Interoperability Protocol provides mechanisms to perform atomic swaps and currency transfers across
multiple turing-complete blockchains. It is a sidechain and series of smart-contracts which 

## Getting Started

The development environment relies on Python 2.7, NodeJS 8.x and Yarn, the following command will install
all required dependencies.

```
make dev
```

## Tests

```
make testrpc &
make test
```

## Docker

This creates the docker image `clearmatics/ion:latest` 

```
make docker-build 
```
