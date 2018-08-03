package ionflow

import (
	"bytes"
	"context"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

func TestCompileAndDeployIon(t *testing.T) {
	// ---------------------------------------------
	// START BLOCKCHAIN SIMULATOR
	// ---------------------------------------------

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

	// create a chain id
	chainID := crypto.Keccak256Hash([]byte("test argument")) // Ion argument

	// start compile and deploy ion
	contractChan := CompileAndDeployIon(ctx, blockchain, userAKey, chainID)

	// commit first block after sent transaction for deployment of patricia trie lib
	blockchain.Commit()
	// patriciaTrieContractInstance := <-contractChan
	<-contractChan

	// commit after transaction for deployment of Ion (with reference to Patricia Trie lib) as been sent
	blockchain.Commit()
	ionContractInstance := <-contractChan

	// call contract variable
	methodName := "chainId"
	out := new([32]byte)
	CallContract(ctx, blockchain, ionContractInstance.Contract, userAAddr, ionContractInstance.Address, methodName, out)

	if !bytes.Equal((*out)[:], chainID.Bytes()) {
		t.Fatal("ERROR chainID result from contract call, and sent to contract constructor differ")
	}
}

func TestRegisterChain(t *testing.T) {
	// check comments on TestCompileAndDeploy()
	ctx := context.Background()
	initialBalance := big.NewInt(1000000000)
	userAKey, _ := crypto.GenerateKey()
	userAAddr := crypto.PubkeyToAddress(userAKey.PublicKey)
	alloc := make(core.GenesisAlloc)
	alloc[userAAddr] = core.GenesisAccount{
		Balance: initialBalance,
	}
	blockchain := backends.NewSimulatedBackend(alloc)
	var chainID [32]byte
	copy(chainID[:], crypto.Keccak256Hash([]byte("DEPLOYEDCHAINID")).Bytes())
	contractChan := CompileAndDeployIon(ctx, blockchain, userAKey, chainID)
	blockchain.Commit()
	<-contractChan
	blockchain.Commit()
	ionContractInstance := <-contractChan

	// deploy validation contract
	contractChan = CompileAndDeployValidation(ctx, blockchain, userAKey, chainID)
	blockchain.Commit()
	validationContractInstance := <-contractChan

	var chainIDA, addressArray [32]byte
	copy(addressArray[:], validationContractInstance.Address.Bytes())
	copy(chainIDA[:], crypto.Keccak256Hash([]byte("TESTCHAINID")).Bytes())
	txRegisterChain := TransactionContract(
		ctx,
		blockchain,
		userAKey,
		ionContractInstance.Contract,
		ionContractInstance.Address,
		nil,
		uint64(3000000),
		"RegisterChain",
		chainIDA,
		addressArray,
	)
	blockchain.Commit()

	registerChainReceipt, err := bind.WaitMined(ctx, blockchain, txRegisterChain)
	if err != nil {
		log.Fatal("ERROR while waiting for contract deployment")
	}
	if registerChainReceipt.Status == 0 {
		t.Fatalf("ERROR transaction of RegisterChain failed!: %#v\n", registerChainReceipt)
	}

	methodName := "chains"
	var isChainRegistered bool
	CallContract(
		ctx,
		blockchain,
		ionContractInstance.Contract,
		userAAddr,
		ionContractInstance.Address,
		methodName,
		&isChainRegistered,
		chainIDA,
	)

	if !isChainRegistered {
		t.Log("ERROR expecting value of chains(validation.address) to be true, but it was ", isChainRegistered)
	}
}

func TestVerifyTx(t *testing.T) {
	ctx := context.Background()

	testValidators := [7]common.Hash{
		common.HexToHash("0x42eb768f2244c8811c63729a21a3569731535f06"),
		common.HexToHash("0x6635f83421bf059cd8111f180f0727128685bae4"),
		common.HexToHash("0x7ffc57839b00206d1ad20c69a1981b489f772031"),
		common.HexToHash("0xb279182d99e65703f0076e4812653aab85fca0f0"),
		common.HexToHash("0xd6ae8250b8348c94847280928c79fb3b63ca453e"),
		common.HexToHash("0xda35dee8eddeaa556e4c26268463e26fb91ff74f"),
		common.HexToHash("0xfc18cbc391de84dbd87db83b20935d3e89f5dd91"),
	}

	deployedChainID := common.HexToHash("0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075")
	testChainID := common.HexToHash("0x22b55e8a4f7c03e1689da845dd463b09299cb3a574e64c68eafc4e99077a7254")

	urlEventChain := "https://rinkeby.infura.io"
	txHashWithEvent := common.HexToHash("0xafc3ab60059ed38e71c7f6bea036822abe16b2c02fcf770a4f4b5fffcbfe6e7e")

	userKey, _ := crypto.GenerateKey()
	userAddr := crypto.PubkeyToAddress(userKey.PublicKey)
	userIntialBalance := big.NewInt(1000000000)

	// ---------------------------------------------
	// GET BLOCK WITH EVENT FROM RINKEBY CHAIN
	// ---------------------------------------------

	clientRPC := ClientRPC(urlEventChain)
	defer clientRPC.Close()

	blockNumberStr, txTrigger, err := BlockNumberByTransactionHash(ctx, clientRPC, txHashWithEvent)
	if err != nil {
		t.Fatal("ERROR couldn't find block by tx hash: ", err)
	}

	var blockNumber big.Int
	blockNumber.SetString((*blockNumberStr)[2:], 16)

	client := ethclient.NewClient(clientRPC)
	eventTxBlockNumber := blockNumber
	block, err := client.BlockByNumber(ctx, &eventTxBlockNumber)
	if err != nil {
		t.Fatal("ERROR retriving block: ", err)
	}

	// ---------------------------------------------
	// START BLOCKCHAIN SIMULATOR
	// ---------------------------------------------
	alloc := make(core.GenesisAlloc)
	alloc[userAddr] = core.GenesisAccount{Balance: userIntialBalance}
	blockchain := backends.NewSimulatedBackend(alloc)

	// ---------------------------------------------
	// COMPILE AND DEPLOY ION
	// ---------------------------------------------
	contractChan := CompileAndDeployIon(ctx, blockchain, userKey, deployedChainID)
	blockchain.Commit()
	<-contractChan // PatriciaTrie libraryContractInstance
	blockchain.Commit()
	ionContractInstance := <-contractChan

	// ---------------------------------------------
	// COMPILE AND DEPLOY VALIDATION
	// ---------------------------------------------
	contractChan = CompileAndDeployValidation(ctx, blockchain, userKey, deployedChainID)
	blockchain.Commit()
	validationContractInstance := <-contractChan

	// ---------------------------------------------
	// REGISTER CHAIN ON ION
	// ---------------------------------------------
	var validationContractAddr [20]byte
	copy(validationContractAddr[:], validationContractInstance.Address.Bytes())
	txRegisterChainIon := TransactionContract(
		ctx,
		blockchain,
		userKey,
		ionContractInstance.Contract,
		ionContractInstance.Address,
		nil,
		uint64(3000000),
		"RegisterChain",
		testChainID,
		validationContractAddr,
	)
	blockchain.Commit()
	registerChainIonReceipt, err := bind.WaitMined(ctx, blockchain, txRegisterChainIon)
	if err != nil || registerChainIonReceipt.Status == 0 {
		t.Fatal("ERROR while waiting for contract deployment")
	}

	// ---------------------------------------------
	// REGISTER CHAIN ON VALIDATION
	// ---------------------------------------------
	var genesisHash [32]byte
	copy(genesisHash[:], block.ParentHash().Bytes())
	txRegisterChainValidation := TransactionContract(
		ctx,
		blockchain,
		userKey,
		validationContractInstance.Contract,
		validationContractInstance.Address,
		nil,
		uint64(3000000),
		"RegisterChain",
		testChainID,
		testValidators,
		genesisHash,
	)
	blockchain.Commit()
	registerChainValidationReceipt, err := bind.WaitMined(ctx, blockchain, txRegisterChainValidation)
	if err != nil || registerChainValidationReceipt.Status == 0 {
		t.Fatal("ERROR while waiting for contract deployment")
	}

	// ---------------------------------------------
	// SUBMIT BLOCK ON VALIDATION
	// ---------------------------------------------
	blockHeader := block.Header()
	extraData := blockHeader.Extra
	unsignedExtraData := extraData[:len(extraData)-(64+1)] // 64 bytes + 1 vanity byte
	signedBlockHeaderRLP, _ := rlp.EncodeToBytes(blockHeader)
	blockHeader.Extra = unsignedExtraData
	unsignedBlockHeaderRLP, _ := rlp.EncodeToBytes(blockHeader)
	txSubmitBlockValidation := TransactionContract(
		ctx,
		blockchain,
		userKey,
		validationContractInstance.Contract,
		validationContractInstance.Address,
		nil,
		uint64(3000000),
		"SubmitBlock",
		testChainID,
		unsignedBlockHeaderRLP,
		signedBlockHeaderRLP,
	)

	blockchain.Commit()
	submitBlockValidationReceipt, err := bind.WaitMined(ctx, blockchain, txSubmitBlockValidation)
	if err != nil || submitBlockValidationReceipt.Status == 0 {
		t.Fatal("ERROR while waiting for contract deployment")
	}

	// ---------------------------------------------
	// CHECK ROOTS PROOF ON ION
	// ---------------------------------------------
	blockHash := block.Hash()
	txTrie := TxTrie(block.Transactions())
	blockReceipts := GetBlockTxReceipts(client, block)
	receiptTrie := ReceiptTrie(blockReceipts)

	txKey := []byte{0x01}
	txProofArr := Proof(txTrie, txKey)
	receiptKey := []byte{0x01}
	receiptProofArr := Proof(receiptTrie, receiptKey)

	checkRootsProofIon := TransactionContract(
		ctx,
		blockchain,
		userKey,
		ionContractInstance.Contract,
		ionContractInstance.Address,
		nil,
		uint64(3000000),
		"CheckRootsProof",
		testChainID,
		blockHash,
		txProofArr,
		receiptProofArr,
	)

	blockchain.Commit()
	chackRootsProofIonReceipt, err := bind.WaitMined(ctx, blockchain, checkRootsProofIon)
	if err != nil || chackRootsProofIonReceipt.Status == 0 {
		t.Fatal("ERROR while waiting for contract deployment", err)
	}

	// ---------------------------------------------
	// COMPILE AND DEPLOY TRIGGER VERIFIER AND CONSUMER FUNCTION
	// ---------------------------------------------
	contractChan = CompileAndDeployTriggerVerifierAndConsumerFunction(
		ctx,
		blockchain,
		userKey,
		ionContractInstance.Address,
	)
	blockchain.Commit()
	<-contractChan // triggerEventVerifierContractInstance := <-contractChan
	blockchain.Commit()
	consumerFunctionContractInstance := <-contractChan

	// ---------------------------------------------
	// FIXME
	// VERIFY FUNCTION EXECUITION
	// let tx = await functionContract.verifyAndExecute(TESTCHAINID, TESTBLOCK.hash, TRIG_DEPLOYED_RINKEBY_ADDR, TEST_PATH, TEST_TX_VALUE, TEST_TX_NODES, TEST_RECEIPT_VALUE, TEST_RECEIPT_NODES, TRIG_CALLED_BY);
	// ---------------------------------------------
	txTriggerPath := []byte{0x13} // SHOULD SOMEHOW BE DYNAMIC!
	txTriggerRLP, _ := rlp.EncodeToBytes(txTrigger)
	txTriggerProofArr := Proof(txTrie, txTriggerPath[:])
	receiptTrigger, _ := rlp.EncodeToBytes(blockReceipts[0x13])
	receiptTriggerProofArr := Proof(receiptTrie, txTriggerPath[:])

	// get tx sender TODO!!!
	signer := types.HomesteadSigner{} // blockchain simulater signer is this one
	triggerCalledBy, _ := types.Sender(signer, txTrigger)

	txVerifyAndExecuteFunction := TransactionContract(
		ctx,
		blockchain,
		userKey,
		consumerFunctionContractInstance.Contract,
		consumerFunctionContractInstance.Address,
		nil,
		uint64(3000000),
		"verifyAndExecute",
		testChainID,
		blockHash,
		txTrigger.To(),         // TRIG_DEPLOYED_RINKEBY_ADDR,
		txTriggerPath,          // TEST_PATH,
		txTriggerRLP,           // TEST_TX_VALUE,
		txTriggerProofArr,      // TEST_TX_NODES,
		receiptTrigger,         // TEST_RECEIPT_VALUE,
		receiptTriggerProofArr, // TEST_RECEIPT_NODES,
		triggerCalledBy,        // TRIG_CALLED_BY,
	)

	blockchain.Commit()
	verifyAndExecuteFunctionReceipt, err := bind.WaitMined(ctx, blockchain, txVerifyAndExecuteFunction)
	if err != nil || verifyAndExecuteFunctionReceipt.Status == 0 {
		t.Logf("\n\n%#v\n\n%#v\n", txTrigger, verifyAndExecuteFunctionReceipt)
		t.Fatal("ERROR while waiting for contract deployment", err)
	}

	// TODO check logs to confirm executed
	//t.Logf("verifyAndExecuteFunctionReceipt: %#v\n", verifyAndExecuteFunctionReceipt.Logs)

	//blockchain.Commit()
	//<-contractChan // PatriciaTrie libraryContractInstance

	// ====================================================================================================
	// TODO
	// let tx = await functionContract.verifyAndExecute(TESTCHAINID, TESTBLOCK.hash, TRIG_DEPLOYED_RINKEBY_ADDR, TEST_PATH, TEST_TX_VALUE, TEST_TX_NODES, TEST_RECEIPT_VALUE, TEST_RECEIPT_NODES, TRIG_CALLED_BY);
}
