package ionflow

import (
	"bytes"
	"context"
	"math/big"
	"regexp"
	"strings"
	"testing"

	contract "github.com/clearmatics/ion/ion-cli/contracts"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestRawTransactionSimulated(t *testing.T) {
	ctx := context.Background()
	initialBalance := big.NewInt(1000000000)

	userAKey, _ := crypto.GenerateKey()
	userAAddr := crypto.PubkeyToAddress(userAKey.PublicKey)

	alloc := make(core.GenesisAlloc)
	alloc[userAAddr] = core.GenesisAccount{
		Balance: initialBalance,
	}

	blockchain := backends.NewSimulatedBackend(alloc)

	userBKey, _ := crypto.GenerateKey()
	userBAddr := crypto.PubkeyToAddress(userBKey.PublicKey)

	// create transaction
	// NewTransaction(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte)
	from := userAAddr                                  // useless variable, just to make clear who the signer is
	nonce, err := blockchain.PendingNonceAt(ctx, from) // uint64(0)
	if err != nil {
		t.Fatal("Error getting pending nonce ", err)
	}
	to := userBAddr
	amount := big.NewInt(10000)                      // random amount
	gasLimit := uint64(30000)                        // random magic number (we could estimate)
	gasPrice, err := blockchain.SuggestGasPrice(ctx) //new(big.Int)
	if err != nil {
		t.Fatal("Error suggesting gas price ", err)
	}
	// data := []byte{}
	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, nil)

	// to better understand the different signers
	// worth looking into https://github.com/ethereum/go-ethereum/blob/cbfb40b0aab093e1b612f3b16834894b2cc67882/core/types/transaction_signing.go#L42-L53
	// signer need chainID
	//chainID := big.NewInt(18)
	//signer := types.NewEIP155Signer(chainID)
	// simulated backend used homestead signer... but the latest might be EIP155?
	signer := types.HomesteadSigner{}
	signedTx, err := types.SignTx(tx, signer, userAKey)
	if err != nil {
		t.Fatal("Error signing tx: ", err)
	}

	blockchain.SendTransaction(ctx, signedTx)
	blockchain.Commit()

	txReceipt, _ := blockchain.TransactionReceipt(ctx, signedTx.Hash())

	gasUsed := big.NewInt(int64(txReceipt.GasUsed))
	balA, err := blockchain.BalanceAt(ctx, userAAddr, nil)
	if err != nil {
		t.Fatal("Error retrieving balance of user A")
	}
	balB, err := blockchain.BalanceAt(ctx, userBAddr, nil)
	if err != nil {
		t.Fatal("Error retrieving balance of user B")
	}
	sum := new(big.Int)
	sum.Add(balA, balB)
	sum.Add(sum, gasUsed)

	// assert
	if sum.Cmp(initialBalance) != 0 {
		t.Fatal("FAILED: bad sum of balances and gas cost of transfer")
	}
}

// inspired by https://medium.com/@akshay_111meher/creating-offline-raw-transactions-with-go-ethereum-8d6cc8174c5d
func TestDeployRawContract(t *testing.T) {
	ctx := context.Background()
	initialBalance := big.NewInt(1000000000)

	userAKey, _ := crypto.GenerateKey()
	userAAddr := crypto.PubkeyToAddress(userAKey.PublicKey)

	// start simulated blockchain
	alloc := make(core.GenesisAlloc)
	alloc[userAAddr] = core.GenesisAccount{
		Balance: initialBalance,
	}
	blockchain := backends.NewSimulatedBackend(alloc)

	// ---------------------------------------------
	// DEPLOY PATRICIA LIB ADDRESS
	// ---------------------------------------------
	// generate payload bytes (we are using PatriciaTrie in our example)
	contractBinStr := contract.PatriciaTrieBin
	contractABIStr := contract.PatriciaTrieABI
	bytecode := common.Hex2Bytes(contractBinStr)
	abiPatriciaTrie, err := abi.JSON(strings.NewReader(contractABIStr))
	if err != nil {
		t.Fatal("ERROR reading PatriciaTrie ABI ", err)
	}
	// packedABI, err := abi.Pack("",arg1,arg2,arg3) // if there is construxtor
	packedABI, err := abiPatriciaTrie.Pack("")
	if err != nil {
		t.Fatal("ERROR packing ABI ", err)
	}
	payloadBytecode := append(bytecode, packedABI...)

	// transaction parameters
	from := userAAddr                                  // useless variable, just to make clear who the signer is
	nonce, err := blockchain.PendingNonceAt(ctx, from) // uint64(0)
	if err != nil {
		t.Fatal("Error getting pending nonce ", err)
	}
	amount := big.NewInt(0)                          // random amount
	gasLimit := uint64(3000000)                      // random magic number (we could estimate)
	gasPrice, err := blockchain.SuggestGasPrice(ctx) //new(big.Int)
	if err != nil {
		t.Fatal("Error suggesting gas price ", err)
	}

	// create contract transaction NewContractCreation is the same has NewTransaction with `to` == nil
	// tx := types.NewTransaction(nonce, nil, amount, gasLimit, gasPrice, payloadBytecode)
	tx := types.NewContractCreation(nonce, amount, gasLimit, gasPrice, payloadBytecode)

	// sign transaction
	signer := types.HomesteadSigner{}
	signedTx, err := types.SignTx(tx, signer, userAKey)
	if err != nil {
		t.Fatal("Error signing tx: ", err)
	}

	blockchain.SendTransaction(ctx, signedTx)
	blockchain.Commit()

	txReceipt, err := blockchain.TransactionReceipt(ctx, signedTx.Hash())
	if err != nil {
		t.Fatal("ERROR getting tx receipt", err)
	}

	// ---------------------------------------------
	// DEPLOY ION CONTRACT WITH PATRICIA LIB ADDRESS
	// ---------------------------------------------
	// after PatriciaTrie is deployed we deploy the Ion contract
	patriciaTrieLibAddr := txReceipt.ContractAddress

	// generate payload bytes (we are using  Ion in our example)
	// we added the key word Ion to make the variables different form the previous
	contractIonBinStr := contract.IonBin
	contractIonABIStr := contract.IonABI

	// add library address to Ion bytecode
	var re = regexp.MustCompile(`__.*PatriciaTrie.*__`)
	contractIonBinStrWithLibAddr := re.ReplaceAllString(contractIonBinStr, patriciaTrieLibAddr.Hex()[2:])

	bytecodeIon := common.Hex2Bytes(contractIonBinStrWithLibAddr)
	abiIon, err := abi.JSON(strings.NewReader(contractIonABIStr))
	if err != nil {
		t.Fatal("ERROR reading PatriciaTrie ABI ", err)
	}

	constructorArg1Ion := crypto.Keccak256Hash([]byte("test argument")) // Ion argument
	packedABIIon, err := abiIon.Pack("", constructorArg1Ion)            // the Ion constructor argument is added here!
	if err != nil {
		t.Fatal("ERROR packing ABI ", err)
	}
	payloadBytecodeIon := append(bytecodeIon, packedABIIon...)

	// transaction parameters
	fromIon := userAAddr                                     // useless variable, just to make clear who the signer is
	nonceIon, err := blockchain.PendingNonceAt(ctx, fromIon) // uint64(0)
	if err != nil {
		t.Fatal("Error getting pending nonce ", err)
	}
	amountIon := big.NewInt(0)
	gasLimitIon := uint64(3000000)                      // random magic number (we could estimate)
	gasPriceIon, err := blockchain.SuggestGasPrice(ctx) //new(big.Int)
	if err != nil {
		t.Fatal("Error suggesting gas price ", err)
	}

	// create transaction
	txIon := types.NewContractCreation(nonceIon, amountIon, gasLimitIon, gasPriceIon, payloadBytecodeIon)

	// sign transaction
	signedTxIon, err := types.SignTx(txIon, signer, userAKey)
	if err != nil {
		t.Fatal("Error signing tx: ", err)
	}

	blockchain.SendTransaction(ctx, signedTxIon)
	blockchain.Commit()

	// test to see if chain id in contract matches the one sent
	txReceiptIon, err := blockchain.TransactionReceipt(ctx, signedTxIon.Hash())
	if err != nil {
		t.Fatal("ERROR getting tx receipt", err)
	}
	ionAddr := txReceiptIon.ContractAddress

	// contract Call function (to actually run a write function we would need to add transaction)
	methodName := "chainId"
	out := new([32]byte)
	input, err := abiIon.Pack(methodName)
	if err != nil {
		t.Fatal("ERROR packing the method name for the contract call", err)
	}
	msg := ethereum.CallMsg{From: userAAddr, To: &ionAddr, Data: input}
	output, err := blockchain.CallContract(ctx, msg, nil)
	if err != nil {
		t.Fatal("ERROR calling the Ion Contract", err)
	}
	err = abiIon.Unpack(out, methodName, output)
	if err != nil {
		t.Fatal("ERROR upacking the call", err)
	}

	if !bytes.Equal((*out)[:], constructorArg1Ion.Bytes()) {
		t.Fatalf("ERROR bytes stored in contract differ from bytes expected\n\tExpected:\t% 0x\n\tResult:\t\t% 0x\n", constructorArg1Ion.Bytes(), *out)
	}
}
