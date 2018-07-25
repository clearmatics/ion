// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"encoding/hex"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// IonABI is the input ABI used to generate the binding from.
const IonABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"chains\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"CheckReceiptProof\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_id\",\"type\":\"bytes32\"},{\"name\":\"_blockHash\",\"type\":\"bytes32\"},{\"name\":\"_rlpBlockHeader\",\"type\":\"bytes\"}],\"name\":\"SubmitBlock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_blockHash\",\"type\":\"bytes32\"}],\"name\":\"getBlockHeader\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[3]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"m_blockheaders\",\"outputs\":[{\"name\":\"prevBlockHash\",\"type\":\"bytes32\"},{\"name\":\"txRootHash\",\"type\":\"bytes32\"},{\"name\":\"receiptRootHash\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"chainId\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_id\",\"type\":\"bytes32\"}],\"name\":\"RegisterChain\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_id\",\"type\":\"bytes32\"},{\"name\":\"_blockHash\",\"type\":\"bytes32\"},{\"name\":\"_value\",\"type\":\"bytes\"},{\"name\":\"_parentNodes\",\"type\":\"bytes\"},{\"name\":\"_path\",\"type\":\"bytes\"}],\"name\":\"CheckTxProof\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"CheckRootsProof\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"m_blockhashes\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_id\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"chainId\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"name\":\"VerifiedTxProof\",\"type\":\"event\"}]"

// IonBin is the compiled bytecode used for deploying new contracts.
const IonBin = `608060405234801561001057600080fd5b50604051602080611487833981018060405281019080805190602001909291905050508060008160001916905550506114398061004e6000396000f3006080604052600436106100a4576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063550325b5146100a957806359349832146100f25780635a0235e4146101095780636b4f9b9d1461018e578063927a32e3146101fb5780639a8a0592146102665780639e43d86b14610299578063affd8be9146102ca578063e318df54146103f3578063fecc37c31461040a575b600080fd5b3480156100b557600080fd5b506100d460048036038101908080359060200190929190505050610461565b60405180826000191660001916815260200191505060405180910390f35b3480156100fe57600080fd5b50610107610484565b005b34801561011557600080fd5b5061018c60048036038101908080356000191690602001909291908035600019169060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610486565b005b34801561019a57600080fd5b506101bd60048036038101908080356000191690602001909291905050506106be565b6040518082600360200280838360005b838110156101e85780820151818401526020810190506101cd565b5050505090500191505060405180910390f35b34801561020757600080fd5b5061022a600480360381019080803560001916906020019092919050505061072d565b60405180846000191660001916815260200183600019166000191681526020018260001916600019168152602001935050505060405180910390f35b34801561027257600080fd5b5061027b610757565b60405180826000191660001916815260200191505060405180910390f35b3480156102a557600080fd5b506102c8600480360381019080803560001916906020019092919050505061075d565b005b3480156102d657600080fd5b506103d960048036038101908080356000191690602001909291908035600019169060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506108f4565b604051808215151515815260200191505060405180910390f35b3480156103ff57600080fd5b50610408610d83565b005b34801561041657600080fd5b50610443600480360381019080803560001916906020019092919080359060200190929190505050610d85565b60405180826000191660001916815260200191505060405180910390f35b60018181548110151561047057fe5b906000526020600020016000915090505481565b565b60606000808560008060009150600090505b6001805490508110156104e4576001818154811015156104b457fe5b906000526020600020015460001916836000191614156104d757600191506104e4565b8080600101915050610498565b811515610559576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f436861696e206973206e6f74207265676973746572656400000000000000000081525060200191505060405180910390fd5b61056a61056588610db5565b610e10565b9550866040518082805190602001908083835b6020831015156105a2578051825260208201915060208101905060208303925061057d565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390209450876000191685600019161415156105e757600080fd5b600360008960001916600019168152602001908152602001600020935061062f61062887600081518110151561061957fe5b90602001906020020151610ecd565b6001610f32565b84600001816000191690555061066661065f87600481518110151561065057fe5b90602001906020020151610ecd565b6001610f32565b84600101816000191690555061069d61069687600581518110151561068757fe5b90602001906020020151610ecd565b6001610f32565b8460020181600019169055506106b38989611010565b505050505050505050565b6106c6611395565b6000600360008460001916600019168152602001908152602001600020905060606040519081016040528082600001546000191660001916815260200182600101546000191660001916815260200182600201546000191660001916815250915050919050565b60036020528060005260406000206000915090508060000154908060010154908060020154905083565b60005481565b6000805460001916826000191614151515610806576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a8152602001807f43616e6e6f7420616464207468697320636861696e20696420746f206368616981526020017f6e2072656769737465720000000000000000000000000000000000000000000081525060400191505060405180910390fd5b600090505b6001805490508110156108be57816000191660018281548110151561082c57fe5b906000526020600020015460001916141515156108b1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f436861696e20616c72656164792065786973747300000000000000000000000081525060200191505060405180910390fd5b808060010191505061080b565b60018290806001815401808255809150509060018203906000526020600020016000909192909190915090600019169055505050565b6000808660008060009150600090505b6001805490508110156109505760018181548110151561092057fe5b906000526020600020015460001916836000191614156109435760019150610950565b8080600101915050610904565b8115156109c5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f436861696e206973206e6f74207265676973746572656400000000000000000081525060200191505060405180910390fd5b8989600060606000809250600260008660001916600019168152602001908152602001600020805480602002602001604051908101604052809291908181526020018280548015610a3957602002820191906000526020600020905b81546000191681526020019060010190808311610a21575b50505050509150600090505b8151811015610a8b578181815181101515610a5c57fe5b906020019060200201516000191684600019161415610a7e5760019250610a8b565b8080600101915050610a45565b821515610b00576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f426c6f636b20646f6573206e6f7420657869737420666f7220636861696e000081525060200191505060405180910390fd5b600360008f60001916600019168152602001908152602001600020985073__./contracts/libraries/PatriciaTrie.s__634f7142ad8e8e8e8d600101546040518563ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808060200180602001806020018560001916600019168152602001848103845288818151815260200191508051906020019080838360005b83811015610bc0578082015181840152602081019050610ba5565b50505050905090810190601f168015610bed5780820380516001836020036101000a031916815260200191505b50848103835287818151815260200191508051906020019080838360005b83811015610c26578082015181840152602081019050610c0b565b50505050905090810190601f168015610c535780820380516001836020036101000a031916815260200191505b50848103825286818151815260200191508051906020019080838360005b83811015610c8c578082015181840152602081019050610c71565b50505050905090810190601f168015610cb95780820380516001836020036101000a031916815260200191505b5097505050505050505060206040518083038186803b158015610cdb57600080fd5b505af4158015610cef573d6000803e3d6000fd5b505050506040513d6020811015610d0557600080fd5b81019080805190602001909291905050501515610d1e57fe5b7f202dc9388a5d69cb591d889833cb0b5cd344fd68afaffc1aaffa5bfe8e79c6a18f8f60405180836000191660001916815260200182600019166000191681526020019250505060405180910390a16001995050505050505050505095945050505050565b565b600260205281600052604060002081815481101515610da057fe5b90600052602060002001600091509150505481565b610dbd6113b8565b600080835191506000821415610deb5760408051908101604052806000815260200160008152509250610e09565b60208401905060408051908101604052808281526020018381525092505b5050919050565b60606000610e1c6113d2565b6000610e27856110b4565b1515610e3257600080fd5b610e3b856110e6565b925082604051908082528060200260200182016040528015610e7757816020015b610e646113f3565b815260200190600190039081610e5c5790505b509350610e8385611160565b91505b610e8f826111a8565b15610ec557610e9d826111d2565b8482815181101515610eab57fe5b906020019060200201819052508080600101915050610e86565b505050919050565b60606000826020015190506000811415610ee657610f2c565b806040519080825280601f01601f191660200182016040528015610f195781602001602082028038833980820191505090505b509150610f2b8360000151838361122e565b5b50919050565b60008060008090505b6020811015611005576008810260ff7f01000000000000000000000000000000000000000000000000000000000000000286838701815181101515610f7c57fe5b9060200101517f010000000000000000000000000000000000000000000000000000000000000090047f010000000000000000000000000000000000000000000000000000000000000002167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916600019169060020a9004821791508080600101915050610f3b565b819250505092915050565b6000806002600085600019166000191681526020019081526020016000209150600090505b818054905081101561107d578260001916828281548110151561105457fe5b9060005260206000200154600019161415151561107057600080fd5b8080600101915050611035565b8183908060018154018082558091505090600182039060005260206000200160009091929091909150906000191690555050505050565b6000806000836020015114156110cd57600091506110e0565b8260000151905060c0815160001a101591505b50919050565b6000806000806000806110f8876110b4565b15156111075760009550611156565b86600001519350835160001a945061111e8761126f565b840192506001876020015185010391505b818311151561115257611141836112fb565b83019250808060010191505061112f565b8095505b5050505050919050565b6111686113d2565b6000611173836110b4565b151561117e57600080fd5b6111878361126f565b83600001510190508282600001819052508082602001818152505050919050565b60006111b26113b8565b826000015190508060200151816000015101836020015110915050919050565b6111da6113b8565b6000806111e6846111a8565b1561122257836020015191506111fb826112fb565b90508183600001818152505080836020018181525050808201846020018181525050611227565b600080fd5b5050919050565b60006020601f83010484602085015b82841461125c576020840280830151818301526001850194505061123d565b6000865160200187015250505050505050565b6000806000808460200151141561128957600092506112f4565b83600001519050805160001a915060808210156112a957600092506112f4565b60b88210806112c5575060c082101580156112c4575060f882105b5b156112d357600192506112f4565b60c08210156112ea57600160b783030192506112f4565b600160f783030192505b5050919050565b600080825160001a90506080811015611317576001915061138f565b60b881101561132e5760016080820301915061138e565b60c08110156113585760b78103806020036101000a6001850151048082016001019350505061138d565b60f881101561136f57600160c0820301915061138c565b60f78103806020036101000a600185015104808201600101935050505b5b5b5b50919050565b606060405190810160405280600390602082028038833980820191505090505090565b604080519081016040528060008152602001600081525090565b6060604051908101604052806113e66113f3565b8152602001600081525090565b6040805190810160405280600081526020016000815250905600a165627a7a7230582092fe2217ce4fc7e2c1df216da8d39cd9ccf25c15d607088915e2dede2778065d0029`

// DeployIon deploys a new Ethereum contract, binding an instance of Ion to it.
func DeployIon(auth *bind.TransactOpts, backend bind.ContractBackend, _id [32]byte) (common.Address, *types.Transaction, *Ion, error) {
	parsed, err := abi.JSON(strings.NewReader(IonABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(IonBin), backend, _id)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Ion{IonCaller: IonCaller{contract: contract}, IonTransactor: IonTransactor{contract: contract}, IonFilterer: IonFilterer{contract: contract}}, nil
}

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

// Deploy and Link Ion to libraries
// func LinkDeployIon(auth *bind.TransactOpts, backend bind.ContractBackend, _id [32]byte, linkAddr common.Address) (common.Address, *types.Transaction, *Ion, error) {
// 	parsed, err := abi.JSON(strings.NewReader(IonABI))
// 	if err != nil {
// 		return common.Address{}, nil, nil, err
// 	}
// 	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(IonBin), backend, _id)
// 	if err != nil {
// 		return common.Address{}, nil, nil, err
// 	}
// 	return address, tx, &Ion{IonCaller: IonCaller{contract: contract}, IonTransactor: IonTransactor{contract: contract}, IonFilterer: IonFilterer{contract: contract}}, nil
// }

// Ion is an auto generated Go binding around an Ethereum contract.
type Ion struct {
	IonCaller     // Read-only binding to the contract
	IonTransactor // Write-only binding to the contract
	IonFilterer   // Log filterer for contract events
}

// IonCaller is an auto generated read-only Go binding around an Ethereum contract.
type IonCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IonTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IonTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IonFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IonFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IonSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IonSession struct {
	Contract     *Ion              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IonCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IonCallerSession struct {
	Contract *IonCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// IonTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IonTransactorSession struct {
	Contract     *IonTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IonRaw is an auto generated low-level Go binding around an Ethereum contract.
type IonRaw struct {
	Contract *Ion // Generic contract binding to access the raw methods on
}

// IonCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IonCallerRaw struct {
	Contract *IonCaller // Generic read-only contract binding to access the raw methods on
}

// IonTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IonTransactorRaw struct {
	Contract *IonTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIon creates a new instance of Ion, bound to a specific deployed contract.
func NewIon(address common.Address, backend bind.ContractBackend) (*Ion, error) {
	contract, err := bindIon(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ion{IonCaller: IonCaller{contract: contract}, IonTransactor: IonTransactor{contract: contract}, IonFilterer: IonFilterer{contract: contract}}, nil
}

// NewIonCaller creates a new read-only instance of Ion, bound to a specific deployed contract.
func NewIonCaller(address common.Address, caller bind.ContractCaller) (*IonCaller, error) {
	contract, err := bindIon(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IonCaller{contract: contract}, nil
}

// NewIonTransactor creates a new write-only instance of Ion, bound to a specific deployed contract.
func NewIonTransactor(address common.Address, transactor bind.ContractTransactor) (*IonTransactor, error) {
	contract, err := bindIon(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IonTransactor{contract: contract}, nil
}

// NewIonFilterer creates a new log filterer instance of Ion, bound to a specific deployed contract.
func NewIonFilterer(address common.Address, filterer bind.ContractFilterer) (*IonFilterer, error) {
	contract, err := bindIon(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IonFilterer{contract: contract}, nil
}

// bindIon binds a generic wrapper to an already deployed contract.
func bindIon(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IonABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ion *IonRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Ion.Contract.IonCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ion *IonRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ion.Contract.IonTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ion *IonRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ion.Contract.IonTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ion *IonCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Ion.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ion *IonTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ion.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ion *IonTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ion.Contract.contract.Transact(opts, method, params...)
}

// CheckReceiptProof is a free data retrieval call binding the contract method 0x59349832.
//
// Solidity: function CheckReceiptProof() constant returns()
func (_Ion *IonCaller) CheckReceiptProof(opts *bind.CallOpts) error {
	var ()
	out := &[]interface{}{}
	err := _Ion.contract.Call(opts, out, "CheckReceiptProof")
	return err
}

// CheckReceiptProof is a free data retrieval call binding the contract method 0x59349832.
//
// Solidity: function CheckReceiptProof() constant returns()
func (_Ion *IonSession) CheckReceiptProof() error {
	return _Ion.Contract.CheckReceiptProof(&_Ion.CallOpts)
}

// CheckReceiptProof is a free data retrieval call binding the contract method 0x59349832.
//
// Solidity: function CheckReceiptProof() constant returns()
func (_Ion *IonCallerSession) CheckReceiptProof() error {
	return _Ion.Contract.CheckReceiptProof(&_Ion.CallOpts)
}

// CheckRootsProof is a free data retrieval call binding the contract method 0xe318df54.
//
// Solidity: function CheckRootsProof() constant returns()
func (_Ion *IonCaller) CheckRootsProof(opts *bind.CallOpts) error {
	var ()
	out := &[]interface{}{}
	err := _Ion.contract.Call(opts, out, "CheckRootsProof")
	return err
}

// CheckRootsProof is a free data retrieval call binding the contract method 0xe318df54.
//
// Solidity: function CheckRootsProof() constant returns()
func (_Ion *IonSession) CheckRootsProof() error {
	return _Ion.Contract.CheckRootsProof(&_Ion.CallOpts)
}

// CheckRootsProof is a free data retrieval call binding the contract method 0xe318df54.
//
// Solidity: function CheckRootsProof() constant returns()
func (_Ion *IonCallerSession) CheckRootsProof() error {
	return _Ion.Contract.CheckRootsProof(&_Ion.CallOpts)
}

// ChainId is a free data retrieval call binding the contract method 0x9a8a0592.
//
// Solidity: function chainId() constant returns(bytes32)
func (_Ion *IonCaller) ChainId(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Ion.contract.Call(opts, out, "chainId")
	return *ret0, err
}

// ChainId is a free data retrieval call binding the contract method 0x9a8a0592.
//
// Solidity: function chainId() constant returns(bytes32)
func (_Ion *IonSession) ChainId() ([32]byte, error) {
	return _Ion.Contract.ChainId(&_Ion.CallOpts)
}

// ChainId is a free data retrieval call binding the contract method 0x9a8a0592.
//
// Solidity: function chainId() constant returns(bytes32)
func (_Ion *IonCallerSession) ChainId() ([32]byte, error) {
	return _Ion.Contract.ChainId(&_Ion.CallOpts)
}

// Chains is a free data retrieval call binding the contract method 0x550325b5.
//
// Solidity: function chains( uint256) constant returns(bytes32)
func (_Ion *IonCaller) Chains(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Ion.contract.Call(opts, out, "chains", arg0)
	return *ret0, err
}

// Chains is a free data retrieval call binding the contract method 0x550325b5.
//
// Solidity: function chains( uint256) constant returns(bytes32)
func (_Ion *IonSession) Chains(arg0 *big.Int) ([32]byte, error) {
	return _Ion.Contract.Chains(&_Ion.CallOpts, arg0)
}

// Chains is a free data retrieval call binding the contract method 0x550325b5.
//
// Solidity: function chains( uint256) constant returns(bytes32)
func (_Ion *IonCallerSession) Chains(arg0 *big.Int) ([32]byte, error) {
	return _Ion.Contract.Chains(&_Ion.CallOpts, arg0)
}

// GetBlockHeader is a free data retrieval call binding the contract method 0x6b4f9b9d.
//
// Solidity: function getBlockHeader(_blockHash bytes32) constant returns(bytes32[3])
func (_Ion *IonCaller) GetBlockHeader(opts *bind.CallOpts, _blockHash [32]byte) ([3][32]byte, error) {
	var (
		ret0 = new([3][32]byte)
	)
	out := ret0
	err := _Ion.contract.Call(opts, out, "getBlockHeader", _blockHash)
	return *ret0, err
}

// GetBlockHeader is a free data retrieval call binding the contract method 0x6b4f9b9d.
//
// Solidity: function getBlockHeader(_blockHash bytes32) constant returns(bytes32[3])
func (_Ion *IonSession) GetBlockHeader(_blockHash [32]byte) ([3][32]byte, error) {
	return _Ion.Contract.GetBlockHeader(&_Ion.CallOpts, _blockHash)
}

// GetBlockHeader is a free data retrieval call binding the contract method 0x6b4f9b9d.
//
// Solidity: function getBlockHeader(_blockHash bytes32) constant returns(bytes32[3])
func (_Ion *IonCallerSession) GetBlockHeader(_blockHash [32]byte) ([3][32]byte, error) {
	return _Ion.Contract.GetBlockHeader(&_Ion.CallOpts, _blockHash)
}

// MBlockhashes is a free data retrieval call binding the contract method 0xfecc37c3.
//
// Solidity: function m_blockhashes( bytes32,  uint256) constant returns(bytes32)
func (_Ion *IonCaller) MBlockhashes(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Ion.contract.Call(opts, out, "m_blockhashes", arg0, arg1)
	return *ret0, err
}

// MBlockhashes is a free data retrieval call binding the contract method 0xfecc37c3.
//
// Solidity: function m_blockhashes( bytes32,  uint256) constant returns(bytes32)
func (_Ion *IonSession) MBlockhashes(arg0 [32]byte, arg1 *big.Int) ([32]byte, error) {
	return _Ion.Contract.MBlockhashes(&_Ion.CallOpts, arg0, arg1)
}

// MBlockhashes is a free data retrieval call binding the contract method 0xfecc37c3.
//
// Solidity: function m_blockhashes( bytes32,  uint256) constant returns(bytes32)
func (_Ion *IonCallerSession) MBlockhashes(arg0 [32]byte, arg1 *big.Int) ([32]byte, error) {
	return _Ion.Contract.MBlockhashes(&_Ion.CallOpts, arg0, arg1)
}

// MBlockheaders is a free data retrieval call binding the contract method 0x927a32e3.
//
// Solidity: function m_blockheaders( bytes32) constant returns(prevBlockHash bytes32, txRootHash bytes32, receiptRootHash bytes32)
func (_Ion *IonCaller) MBlockheaders(opts *bind.CallOpts, arg0 [32]byte) (struct {
	PrevBlockHash   [32]byte
	TxRootHash      [32]byte
	ReceiptRootHash [32]byte
}, error) {
	ret := new(struct {
		PrevBlockHash   [32]byte
		TxRootHash      [32]byte
		ReceiptRootHash [32]byte
	})
	out := ret
	err := _Ion.contract.Call(opts, out, "m_blockheaders", arg0)
	return *ret, err
}

// MBlockheaders is a free data retrieval call binding the contract method 0x927a32e3.
//
// Solidity: function m_blockheaders( bytes32) constant returns(prevBlockHash bytes32, txRootHash bytes32, receiptRootHash bytes32)
func (_Ion *IonSession) MBlockheaders(arg0 [32]byte) (struct {
	PrevBlockHash   [32]byte
	TxRootHash      [32]byte
	ReceiptRootHash [32]byte
}, error) {
	return _Ion.Contract.MBlockheaders(&_Ion.CallOpts, arg0)
}

// MBlockheaders is a free data retrieval call binding the contract method 0x927a32e3.
//
// Solidity: function m_blockheaders( bytes32) constant returns(prevBlockHash bytes32, txRootHash bytes32, receiptRootHash bytes32)
func (_Ion *IonCallerSession) MBlockheaders(arg0 [32]byte) (struct {
	PrevBlockHash   [32]byte
	TxRootHash      [32]byte
	ReceiptRootHash [32]byte
}, error) {
	return _Ion.Contract.MBlockheaders(&_Ion.CallOpts, arg0)
}

// CheckTxProof is a paid mutator transaction binding the contract method 0xaffd8be9.
//
// Solidity: function CheckTxProof(_id bytes32, _blockHash bytes32, _value bytes, _parentNodes bytes, _path bytes) returns(bool)
func (_Ion *IonTransactor) CheckTxProof(opts *bind.TransactOpts, _id [32]byte, _blockHash [32]byte, _value []byte, _parentNodes []byte, _path []byte) (*types.Transaction, error) {
	return _Ion.contract.Transact(opts, "CheckTxProof", _id, _blockHash, _value, _parentNodes, _path)
}

// CheckTxProof is a paid mutator transaction binding the contract method 0xaffd8be9.
//
// Solidity: function CheckTxProof(_id bytes32, _blockHash bytes32, _value bytes, _parentNodes bytes, _path bytes) returns(bool)
func (_Ion *IonSession) CheckTxProof(_id [32]byte, _blockHash [32]byte, _value []byte, _parentNodes []byte, _path []byte) (*types.Transaction, error) {
	return _Ion.Contract.CheckTxProof(&_Ion.TransactOpts, _id, _blockHash, _value, _parentNodes, _path)
}

// CheckTxProof is a paid mutator transaction binding the contract method 0xaffd8be9.
//
// Solidity: function CheckTxProof(_id bytes32, _blockHash bytes32, _value bytes, _parentNodes bytes, _path bytes) returns(bool)
func (_Ion *IonTransactorSession) CheckTxProof(_id [32]byte, _blockHash [32]byte, _value []byte, _parentNodes []byte, _path []byte) (*types.Transaction, error) {
	return _Ion.Contract.CheckTxProof(&_Ion.TransactOpts, _id, _blockHash, _value, _parentNodes, _path)
}

// RegisterChain is a paid mutator transaction binding the contract method 0x9e43d86b.
//
// Solidity: function RegisterChain(_id bytes32) returns()
func (_Ion *IonTransactor) RegisterChain(opts *bind.TransactOpts, _id [32]byte) (*types.Transaction, error) {
	return _Ion.contract.Transact(opts, "RegisterChain", _id)
}

// RegisterChain is a paid mutator transaction binding the contract method 0x9e43d86b.
//
// Solidity: function RegisterChain(_id bytes32) returns()
func (_Ion *IonSession) RegisterChain(_id [32]byte) (*types.Transaction, error) {
	return _Ion.Contract.RegisterChain(&_Ion.TransactOpts, _id)
}

// RegisterChain is a paid mutator transaction binding the contract method 0x9e43d86b.
//
// Solidity: function RegisterChain(_id bytes32) returns()
func (_Ion *IonTransactorSession) RegisterChain(_id [32]byte) (*types.Transaction, error) {
	return _Ion.Contract.RegisterChain(&_Ion.TransactOpts, _id)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0x5a0235e4.
//
// Solidity: function SubmitBlock(_id bytes32, _blockHash bytes32, _rlpBlockHeader bytes) returns()
func (_Ion *IonTransactor) SubmitBlock(opts *bind.TransactOpts, _id [32]byte, _blockHash [32]byte, _rlpBlockHeader []byte) (*types.Transaction, error) {
	return _Ion.contract.Transact(opts, "SubmitBlock", _id, _blockHash, _rlpBlockHeader)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0x5a0235e4.
//
// Solidity: function SubmitBlock(_id bytes32, _blockHash bytes32, _rlpBlockHeader bytes) returns()
func (_Ion *IonSession) SubmitBlock(_id [32]byte, _blockHash [32]byte, _rlpBlockHeader []byte) (*types.Transaction, error) {
	return _Ion.Contract.SubmitBlock(&_Ion.TransactOpts, _id, _blockHash, _rlpBlockHeader)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0x5a0235e4.
//
// Solidity: function SubmitBlock(_id bytes32, _blockHash bytes32, _rlpBlockHeader bytes) returns()
func (_Ion *IonTransactorSession) SubmitBlock(_id [32]byte, _blockHash [32]byte, _rlpBlockHeader []byte) (*types.Transaction, error) {
	return _Ion.Contract.SubmitBlock(&_Ion.TransactOpts, _id, _blockHash, _rlpBlockHeader)
}

// IonVerifiedTxProofIterator is returned from FilterVerifiedTxProof and is used to iterate over the raw logs and unpacked data for VerifiedTxProof events raised by the Ion contract.
type IonVerifiedTxProofIterator struct {
	Event *IonVerifiedTxProof // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IonVerifiedTxProofIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IonVerifiedTxProof)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IonVerifiedTxProof)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IonVerifiedTxProofIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IonVerifiedTxProofIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IonVerifiedTxProof represents a VerifiedTxProof event raised by the Ion contract.
type IonVerifiedTxProof struct {
	ChainId   [32]byte
	BlockHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterVerifiedTxProof is a free log retrieval operation binding the contract event 0x202dc9388a5d69cb591d889833cb0b5cd344fd68afaffc1aaffa5bfe8e79c6a1.
//
// Solidity: event VerifiedTxProof(chainId bytes32, blockHash bytes32)
func (_Ion *IonFilterer) FilterVerifiedTxProof(opts *bind.FilterOpts) (*IonVerifiedTxProofIterator, error) {

	logs, sub, err := _Ion.contract.FilterLogs(opts, "VerifiedTxProof")
	if err != nil {
		return nil, err
	}
	return &IonVerifiedTxProofIterator{contract: _Ion.contract, event: "VerifiedTxProof", logs: logs, sub: sub}, nil
}

// WatchVerifiedTxProof is a free log subscription operation binding the contract event 0x202dc9388a5d69cb591d889833cb0b5cd344fd68afaffc1aaffa5bfe8e79c6a1.
//
// Solidity: event VerifiedTxProof(chainId bytes32, blockHash bytes32)
func (_Ion *IonFilterer) WatchVerifiedTxProof(opts *bind.WatchOpts, sink chan<- *IonVerifiedTxProof) (event.Subscription, error) {

	logs, sub, err := _Ion.contract.WatchLogs(opts, "VerifiedTxProof")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IonVerifiedTxProof)
				if err := _Ion.contract.UnpackLog(event, "VerifiedTxProof", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
