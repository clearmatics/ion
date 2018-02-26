# Ion Interoperability Protocol

The Ion Interoperability Protocol provides mechanisms to perform atomic swaps and currency transfers across
multiple turing-complete blockchains. It is a sidechain and series of smart-contracts which 

## Getting Started

Building Ion requires many components such as Python, NodeJS, Yarn, Solidity compiler, Docker etc.
The `make dev` command will install all necessary dependencies required to develop and build Ion.

All Ion functionality is available via the built-in shell, or via the `ion` command:

```commandline
$ python -mion shell
> onchain Token transfer --help
```

or 

```
$ python -mion onchain Token transfer --help
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

# Commands

 * `shell` - REPL environment
 * `etheventrelay` - Relay Ethereum events as merkle roots 
 * `plasma payment` - Create Plasma payment
 * `plasma chain` - Perform operations on the Plasma chain
 * `rpc server` - Run a Plasma RPC server which syncs to a smart contract
 * `rpc client` - JSON-RPC client for Plasma RPC server
 
## Smart-Contract integration

All `.abi` files in the `abi` directory are available via the `onchain` command,
for example the `abi/Token.abi` file is available as a command within the Ion shell:

```commandline
> onchain Token transfer --help
Usage: __main__.py  onchain Token transfer [OPTIONS] to value data

  transfer(address,uint256,bytes)

Options:
  -w      Wait until mined
  -c      Commit transaction
  --help  Show this message and exit.
```
