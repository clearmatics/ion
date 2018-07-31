// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func LinkDeployIon(auth *bind.TransactOpts, backend bind.ContractBackend, _id [32]byte, linkAddr common.Address, linkString string) (common.Address, *types.Transaction, *Ion, error) {
	// Convert address to string and replace library reference in Bin
	linkAddrStr := hex.EncodeToString(linkAddr.Bytes())
	NewIonBin := strings.Replace(IonBin, linkString, linkAddrStr, 1)

	parsed, err := abi.JSON(strings.NewReader(IonABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(NewIonBin), backend, _id)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Ion{IonCaller: IonCaller{contract: contract}, IonTransactor: IonTransactor{contract: contract}, IonFilterer: IonFilterer{contract: contract}}, nil
}
