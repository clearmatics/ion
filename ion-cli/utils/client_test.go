package utils_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/clearmatics/ion/ion-cli/utils"
)

const URL = "https://mainnet.infura.io"

// NOTE: This tests depend on an external network (not really good)

func TestClient(t *testing.T) {
	client := utils.Client(URL)
	client.Close()
}

func TestGetReceipts(t *testing.T) {
	expectedTotalReceipts := 92

	client := utils.Client(URL)
	defer client.Close()

	blockNumber := big.NewInt(6021002)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		t.Error(err)
	}

	receiptArr := utils.GetBlockTxReceipts(client, block)

	if len(receiptArr) != expectedTotalReceipts {
		t.Errorf("Got %d receipts and expected %d receipts!\n", len(receiptArr), expectedTotalReceipts)
	}
}

func TestBlockNumberByTransactionHash(t *testing.T) {
	client := utils.Client(URL)
	defer client.Close()

	blockNumber := big.NewInt(6021002)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		t.Fatal(err)
	}
	txArr := block.Transactions()
	tx := txArr[0]
	txHash := tx.Hash()

	// needs to use the ClientRPC because we make the request directly to the RPC in order to get the blocknumber
	clientRPC := utils.ClientRPC(URL)
	defer clientRPC.Close()

	bNumber, _, err := utils.BlockNumberByTransactionHash(context.Background(), clientRPC, txHash)
	if err != nil {
		t.Fatal(err)
	}

	var bNumberInt big.Int
	t.Log(bNumber)
	t.Log((*bNumber)[2:])
	bNumberInt.SetString((*bNumber)[2:], 16)
	t.Log(bNumberInt)

	if blockNumber.Cmp(&bNumberInt) != 0 {
		t.Errorf("Blocknumber retrieved by transaction hash is not right. It expected %s but got %s\n", blockNumber.String(), bNumberInt.String())
	}
}
