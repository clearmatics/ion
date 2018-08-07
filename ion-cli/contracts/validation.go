package contract

import (
	"context"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/core/types"
)

func RegisterChain(ctx context.Context, backend bind.ContractBackend, userKey *ecdsa.PrivateKey, contract *compiler.Contract, to common.Address, chainId common.Hash, ionAddr common.Address, validators []common.Address, registerHash common.Hash) (tx *types.Transaction) {
	tx = TransactionContract(
		ctx,
		backend,
		userKey,
		contract,
		to,
		nil,
		uint64(3000000),
		"RegisterChain",
		chainId,
		ionAddr,
		validators,
		registerHash,
	)

	return
}
