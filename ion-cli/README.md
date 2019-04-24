# Ion CLI

The Command-Line Interface reference.

## Building the CLI

Clean the directory: `make clean`.

Run `make build` to fetch sources and compile the binary. This will error if attempted to be built outside of your `GOPATH`.

### Requirements

* [`solc`](https://github.com/ethereum/solidity/releases) ensure matching compiler version with target solidity code as 0.4.x will not be compatible with solc 0.5

## Usage
Run `./ion-cli` to launch the interface.

```
addAccount                  use:   addAccount [name] [path/to/keystore]
                            description: Add account to be used for transactions

Adds an account to be referenced by given name that can be used to sign transactions from. Takes a name
(string) that can be any chosen name, and a path to a keystore file containing an encrypted private key.

Accounts are not persisted between sessions.
```
```
addContractInstance         use:   addContractInstance [name] [path/to/solidity/contract]
                            description: Compiles a contract for use

Adds a compiled contract instance. This contract instance can be deployed to a connected chain or used
to make function calls to deployed contracts. Takes any chosen name (string) and a path to the Solidity
contract to be compiled.

Contracts are not persisted between sessions.
```
```
clear                       clear the screen
```
```
connectToClient             use:   connectToClient [rpc url]
                            description: Connects to an RPC client to be used

Connects to an Ethereum blockchain specified by a given url. This chain will be used for contract
deployment, calls, chain data etc. unless specified otherwise.

Only one client can be connected at a time.
```
```
deployContract              use:   deployContract [contract name] [account name] [gas limit]
                            description: Deploys specified contract instance to connected client

Deploys a specified compiled contract instance by name from the specified account name to the connected
client. Requires gas limit to be provided.
```
```
exit                        exit the program
```
```
getBlockByHash              use:   getBlockByHash [optional rpc url] [hash]
                            description: Returns block header specified by hash from connected client
                            or from specific endpoint

Retrieves the block header from the connected or specified endpoint at given block hash.
```
```
getBlockByHash_Clique       use:   getBlockByHash_Clique [optional rpc url] [hash]
                            description: Returns signed and unsigned RLP-encoded block headers by
                            block hash required for submission to Clique validation connected
                            client or specified endpoint

Retrieves RLP-encoded block headers that have been signed and unsigned by validators of the Clique
mechanism specified by block hash. This is used to be submitted to the Clique validation contract
scheme for interoperation with Clique proof-of-authority chains.
```
```
getBlockByNumber            use:   getBlockByNumber [optional rpc url] [integer]
                            description: Returns block header specified by height from connected
                            client or from specified endpoint

Retrieves the block header from the connected or specified endpoint at given block number.
```
```
getBlockByNumber_Clique     use:   getBlockByNumber_Clique [optional rpc url] [integer]
                            description: Returns signed and unsigned RLP-encoded block headers by
                            block number required for submission to Clique validation from connected
                            client or specified endpoint

Retrieves RLP-encoded block headers that have been signed and unsigned by validators of the Clique
mechanism specified by block number. This is used to be submitted to the Clique validation contract
scheme for interoperation with Clique proof-of-authority chains.
```
```
getProof                    use:   getProof [optional rpc url] [Transaction Hash]
                            description: Returns a merkle patricia proof of a specific transaction
                            and its receipt in a block

Retrieves a merkle patricia trie proof of a transaction in a block. It checks against the Transaction
Trie and the Receipt Trie in a block, reconstructs the trie and generates a proof that the specified
transaction exists in the trie and returns the data to verify the proof.
```
```
getTransactionByHash        use:   getTransactionByHash [optional rpc url] [hash]
                            description: Returns transaction specified by hash from connected
                            client or specified endpoint

Retrieves the transaction data of a specified transaction by hash.
```
```
help                        display help
```
```
linkAndDeployContract       use:   linkAndDeployContract [contract name] [account name] [gas limit]
                            description: Deploys specified contract instance to connected client

Links any required libraries to a contract instance before deploying to the connect client.
```
```
listAccounts                use:   listAccounts
                            description: List all added accounts

Lists all added accounts.
```
```
listContracts               use:   listContracts
                            description: List compiled contract instances

Lists all added contract instances.
```
```
transactionMessage          use:   transactionMessage [contract name] [function name] [from account name] [deployed contract address] [amount] [gasLimit]
                            description: Calls a contract function as a transaction.

Calls a contract instance function that has been deployed to the connected client as a transaction. This
will mutate the world state of the destination chain if the transaction is accepted into a block.
Requires gas.

Will be prompted to provide function input parameters.
```
