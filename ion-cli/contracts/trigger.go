// Copyright (c) 2018 Clearmatics Technologies Ltd
package contract

import (
	"context"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/core/types"
)

// Registers chain with Validation contract specified
func Fire(
	ctx context.Context,
	backend bind.ContractBackend,
	userKey *ecdsa.PrivateKey,
	contract *compiler.Contract,
	toAddr common.Address,
) (tx *types.Transaction) {
	tx = TransactionContract(
		ctx,
		backend,
		userKey,
		contract,
		toAddr,
		nil,
		uint64(3000000),
		"fire",
	)

	return
}
