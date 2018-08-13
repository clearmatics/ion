// Copyright (c) 2018 Clearmatics Technologies Ltd
package contract

import (
	"bytes"
	"context"
	"math/big"
	"testing"

	"github.com/clearmatics/ion/ion-cli/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

// TODO
// avoid having to get data from Rinkeby, make it deploy a trigger function into a PoA chain

// TestVerifyTx test for the full flow of Ion
func Test_VerifyTx(t *testing.T) {
	ctx := context.Background()

	// ---------------------------------------------
	// HARD CODED DATA
	// ---------------------------------------------
	testValidators := []common.Address{
		common.HexToAddress("0x42eb768f2244c8811c63729a21a3569731535f06"),
		common.HexToAddress("0x6635f83421bf059cd8111f180f0727128685bae4"),
		common.HexToAddress("0x7ffc57839b00206d1ad20c69a1981b489f772031"),
		common.HexToAddress("0xb279182d99e65703f0076e4812653aab85fca0f0"),
		common.HexToAddress("0xd6ae8250b8348c94847280928c79fb3b63ca453e"),
		common.HexToAddress("0xda35dee8eddeaa556e4c26268463e26fb91ff74f"),
		common.HexToAddress("0xfc18cbc391de84dbd87db83b20935d3e89f5dd91"),
	}

	deployedChainID := common.HexToHash("0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075")
	testChainID := common.HexToHash("0x22b55e8a4f7c03e1689da845dd463b09299cb3a574e64c68eafc4e99077a7254")

	urlEventChain := "https://rinkeby.infura.io"
	txHashWithEvent := common.HexToHash("0xafc3ab60059ed38e71c7f6bea036822abe16b2c02fcf770a4f4b5fffcbfe6e7e")

	signer := types.HomesteadSigner{} // blockchain simulater signer is this one
	userKey, _ := crypto.GenerateKey()
	userAddr := crypto.PubkeyToAddress(userKey.PublicKey)
	userIntialBalance := big.NewInt(1000000000)

	// ---------------------------------------------
	// GET BLOCK WITH EVENT FROM RINKEBY CHAIN
	// ---------------------------------------------

	clientRPC := utils.ClientRPC(urlEventChain)
	defer clientRPC.Close()

	blockNumberStr, txTrigger, err := utils.BlockNumberByTransactionHash(ctx, clientRPC, txHashWithEvent)
	if err != nil {
		t.Fatal("ERROR couldn't find block by tx hash: ", err)
	}

	var blockNumber big.Int
	blockNumber.SetString((*blockNumberStr)[2:], 16)

	client := ethclient.NewClient(clientRPC)
	eventTxBlockNumber := blockNumber
	block, err := client.BlockByNumber(ctx, &eventTxBlockNumber)
	if err != nil {
		t.Fatal("ERROR retrieving block: ", err)
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
	// REGISTER CHAIN ON VALIDATION
	// ---------------------------------------------
	var ionContractAddr [20]byte
	copy(ionContractAddr[:], ionContractInstance.Address.Bytes())
	var genesisHash [32]byte
	copy(genesisHash[:], block.ParentHash().Bytes())
	txRegisterChainValidation := RegisterChain(
		ctx,
		blockchain,
		userKey,
		validationContractInstance.Contract,
		validationContractInstance.Address,
		testChainID,
		ionContractAddr,
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
	txSubmitBlockValidation := SubmitBlock(
		ctx,
		blockchain,
		userKey,
		validationContractInstance.Contract,
		validationContractInstance.Address,
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
	blockTransactions := block.Transactions()
	txTrie := utils.TxTrie(blockTransactions)
	blockReceipts := utils.GetBlockTxReceipts(client, block)
	receiptTrie := utils.ReceiptTrie(blockReceipts)

	txKey := []byte{0x01}
	txProofArr := utils.Proof(txTrie, txKey)
	receiptKey := []byte{0x01}
	receiptProofArr := utils.Proof(receiptTrie, receiptKey)

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
	// VERIFY FUNCTION EXECUITION
	// ---------------------------------------------
	triggerCalledBy, _ := types.Sender(signer, txTrigger)

	// Generate the proof
	txPath, txValue, txNodes, receiptValue, receiptNodes := utils.GenerateProof(
		ctx,
		clientRPC,
		txHashWithEvent,
	)

	txVerifyAndExecuteFunction := VerifyExecute(
		ctx,
		blockchain,
		userKey,
		consumerFunctionContractInstance.Contract,
		consumerFunctionContractInstance.Address,
		testChainID,
		blockHash,
		*txTrigger.To(), // TRIG_DEPLOYED_RINKEBY_ADDR,
		txPath,          // TEST_PATH,
		txValue,         // TEST_TX_VALUE,
		txNodes,         // TEST_TX_NODES,
		receiptValue,    // TEST_RECEIPT_VALUE,
		receiptNodes,    // TEST_RECEIPT_NODES,
		triggerCalledBy, // TRIG_CALLED_BY,
		nil,
	)

	blockchain.Commit()
	verifyAndExecuteFunctionReceipt, err := bind.WaitMined(ctx, blockchain, txVerifyAndExecuteFunction)
	if err != nil || verifyAndExecuteFunctionReceipt.Status == 0 {
		t.Logf("\n\n%#v\n\n%#v\n", txTrigger, verifyAndExecuteFunctionReceipt)
		t.Fatal("ERROR while waiting for contract deployment", err)
	}

	// confirm the Executed event was emited by Consumer Function
	eventSignatureHash := crypto.Keccak256Hash([]byte("Executed()")) // Ion argument

	foundExecuted := false
	for _, vlog := range verifyAndExecuteFunctionReceipt.Logs {
		if len(vlog.Topics) < 1 {
			continue
		}
		foundExecuted = bytes.Equal(vlog.Topics[0].Bytes(), eventSignatureHash.Bytes())
		if foundExecuted {
			break
		}
	}
	if !foundExecuted {
		t.Fatal("ERROR did not find Executed() event")
	}
}
