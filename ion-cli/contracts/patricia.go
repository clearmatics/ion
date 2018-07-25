// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// PatriciaTrieABI is the input ABI used to generate the binding from.
const PatriciaTrieABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_value\",\"type\":\"bytes\"},{\"name\":\"_parentNodes\",\"type\":\"bytes\"},{\"name\":\"_path\",\"type\":\"bytes\"},{\"name\":\"_root\",\"type\":\"bytes32\"}],\"name\":\"verifyProof\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// PatriciaTrieBin is the compiled bytecode used for deploying new contracts.
const PatriciaTrieBin = `611278610030600b82828239805160001a6073146000811461002057610022565bfe5b5030600052607381538281f3007300000000000000000000000000000000000000003014608060405260043610610058576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680634f7142ad1461005d575b600080fd5b610151600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803560001916906020019092919050505061016b565b604051808215151515815260200191505060405180910390f35b60006101756111d5565b606060006060806000606060008060608061018f8f61060f565b9a5061019a8b61066a565b99508c9850600095506101ae8e6000610727565b9450600093505b89518410156105fb576101de8a858151811015156101cf57fe5b90602001906020020151610ad7565b9750876040518082805190602001908083835b60208310151561021657805182526020820191506020810190506020830392506101f1565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902060001916896000191614151561025c5760009b506105fc565b61027c8a8581518110151561026d57fe5b9060200190602002015161066a565b9650601187511415610446578451861415610397578f6040518082805190602001908083835b6020831015156102c757805182526020820191506020810190506020830392506102a2565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390206000191661031988601081518110151561030a57fe5b90602001906020020151610ad7565b6040518082805190602001908083835b60208310151561034e5780518252602082019150602081019050602083039250610329565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902060001916141561038e5760019b506105fc565b60009b506105fc565b84868151811015156103a557fe5b9060200101517f010000000000000000000000000000000000000000000000000000000000000090047f0100000000000000000000000000000000000000000000000000000000000000027f010000000000000000000000000000000000000000000000000000000000000090049250610439878461ffff1681518110151561042a57fe5b90602001906020020151610b3c565b98506001860195506105ee565b6002875114156105e45761047187600081518110151561046257fe5b90602001906020020151610b51565b915061047e826001610727565b51860195508451861415610592578f6040518082805190602001908083835b6020831015156104c2578051825260208201915060208101905060208303925061049d565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390206000191661051488600181518110151561050557fe5b90602001906020020151610b51565b6040518082805190602001908083835b6020831015156105495780518252602082019150602081019050602083039250610524565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390206000191614156105895760019b506105fc565b60009b506105fc565b600061059f836001610727565b5114156105af5760009b506105fc565b6105d08760018151811015156105c157fe5b90602001906020020151610b51565b90506105dd816000610bc3565b98506105ed565b60009b506105fc565b5b83806001019450506101b5565b5b5050505050505050505050949350505050565b6106176111d5565b6000808351915060008214156106455760408051908101604052806000815260200160008152509250610663565b60208401905060408051908101604052808281526020018381525092505b5050919050565b606060006106766111ef565b600061068185610ca1565b151561068c57600080fd5b61069585610cd3565b9250826040519080825280602002602001820160405280156106d157816020015b6106be611210565b8152602001906001900390816106b65790505b5093506106dd85610d4d565b91505b6106e982610d95565b1561071f576106f782610dbf565b848281518110151561070557fe5b9060200190602002018190525080806001019150506106e0565b505050919050565b60608060008061073561122a565b6060600060ff6040519080825280601f01601f19166020018201604052801561076d5781602001602082028038833980820191505090505b50955060009450600093505b88518410156109db576107e3898581518110151561079357fe5b9060200101517f010000000000000000000000000000000000000000000000000000000000000090047f010000000000000000000000000000000000000000000000000000000000000002610e1b565b92508780156107f25750600084145b1561091a5760017f01000000000000000000000000000000000000000000000000000000000000000283600060028110151561082a57fe5b60200201517effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614806108b4575060037f01000000000000000000000000000000000000000000000000000000000000000283600060028110151561088b57fe5b60200201517effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916145b15610915578260016002811015156108c857fe5b6020020151868660ff168151811015156108de57fe5b9060200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053506001850194505b6109ce565b82600060028110151561092957fe5b6020020151868660ff1681518110151561093f57fe5b9060200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535082600160028110151561097e57fe5b6020020151866001870160ff1681518110151561099757fe5b9060200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053506002850194505b8380600101945050610779565b8460ff166040519080825280601f01601f191660200182016040528015610a115781602001602082028038833980820191505090505b509150600090505b8460ff16811015610ac8578581815181101515610a3257fe5b9060200101517f010000000000000000000000000000000000000000000000000000000000000090047f0100000000000000000000000000000000000000000000000000000000000000028282815181101515610a8b57fe5b9060200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508080600101915050610a19565b81965050505050505092915050565b60606000826020015190506000811415610af057610b36565b806040519080825280601f01601f191660200182016040528015610b235781602001602082028038833980820191505090505b509150610b3583600001518383610f01565b5b50919050565b6000610b4782610f42565b6001029050919050565b6060600080610b5f84610f99565b1515610b6a57600080fd5b610b7384610fca565b8092508193505050806040519080825280601f01601f191660200182016040528015610bae5781602001602082028038833980820191505090505b509250610bbc828483610f01565b5050919050565b60008060008090505b6020811015610c96576008810260ff7f01000000000000000000000000000000000000000000000000000000000000000286838701815181101515610c0d57fe5b9060200101517f010000000000000000000000000000000000000000000000000000000000000090047f010000000000000000000000000000000000000000000000000000000000000002167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916600019169060020a9004821791508080600101915050610bcc565b819250505092915050565b600080600083602001511415610cba5760009150610ccd565b8260000151905060c0815160001a101591505b50919050565b600080600080600080610ce587610ca1565b1515610cf45760009550610d43565b86600001519350835160001a9450610d0b8761104b565b840192506001876020015185010391505b8183111515610d3f57610d2e836110d7565b830192508080600101915050610d1c565b8095505b5050505050919050565b610d556111ef565b6000610d6083610ca1565b1515610d6b57600080fd5b610d748361104b565b83600001510190508282600001819052508082602001818152505050919050565b6000610d9f6111d5565b826000015190508060200151816000015101836020015110915050919050565b610dc76111d5565b600080610dd384610d95565b15610e0f5783602001519150610de8826110d7565b90508183600001818152505080836020018181525050808201846020018181525050610e14565b600080fd5b5050919050565b610e2361122a565b600080610e31846004611171565b9150600f7f010000000000000000000000000000000000000000000000000000000000000002841690506040805190810160405280837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681525092505050919050565b60006020601f83010484602085015b828414610f2f5760208402808301518183015260018501945050610f10565b6000865160200187015250505050505050565b6000806000610f5084610f99565b1515610f5b57600080fd5b610f6484610fca565b80925081935050506020811180610f7b5750600081145b15610f8557600080fd5b806020036101000a82510492505050919050565b600080600083602001511415610fb25760009150610fc4565b8260000151905060c0815160001a1091505b50919050565b6000806000806000610fdb86610f99565b1515610fe657600080fd5b85600001519150815160001a925060808310156110095781945060019350611043565b60b88310156110275760018660200151039350600182019450611042565b60b78303905080600187602001510303935060018183010194505b5b505050915091565b6000806000808460200151141561106557600092506110d0565b83600001519050805160001a9150608082101561108557600092506110d0565b60b88210806110a1575060c082101580156110a0575060f882105b5b156110af57600192506110d0565b60c08210156110c657600160b783030192506110d0565b600160f783030192505b5050919050565b600080825160001a905060808110156110f3576001915061116b565b60b881101561110a5760016080820301915061116a565b60c08110156111345760b78103806020036101000a60018501510480820160010193505050611169565b60f881101561114b57600160c08203019150611168565b60f78103806020036101000a600185015104808201600101935050505b5b5b5b50919050565b60008160020a60ff16837f0100000000000000000000000000000000000000000000000000000000000000900460ff168115156111aa57fe5b047f010000000000000000000000000000000000000000000000000000000000000002905092915050565b604080519081016040528060008152602001600081525090565b606060405190810160405280611203611210565b8152602001600081525090565b604080519081016040528060008152602001600081525090565b60408051908101604052806002906020820280388339808201915050905050905600a165627a7a7230582027380a33dbfdf1a14b565c857047780931c383048b159e2e7117549db774476a0029`

// DeployPatriciaTrie deploys a new Ethereum contract, binding an instance of PatriciaTrie to it.
func DeployPatriciaTrie(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *PatriciaTrie, error) {
	parsed, err := abi.JSON(strings.NewReader(PatriciaTrieABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(PatriciaTrieBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PatriciaTrie{PatriciaTrieCaller: PatriciaTrieCaller{contract: contract}, PatriciaTrieTransactor: PatriciaTrieTransactor{contract: contract}, PatriciaTrieFilterer: PatriciaTrieFilterer{contract: contract}}, nil
}

// PatriciaTrie is an auto generated Go binding around an Ethereum contract.
type PatriciaTrie struct {
	PatriciaTrieCaller     // Read-only binding to the contract
	PatriciaTrieTransactor // Write-only binding to the contract
	PatriciaTrieFilterer   // Log filterer for contract events
}

// PatriciaTrieCaller is an auto generated read-only Go binding around an Ethereum contract.
type PatriciaTrieCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PatriciaTrieTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PatriciaTrieTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PatriciaTrieFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PatriciaTrieFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PatriciaTrieSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PatriciaTrieSession struct {
	Contract     *PatriciaTrie     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PatriciaTrieCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PatriciaTrieCallerSession struct {
	Contract *PatriciaTrieCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// PatriciaTrieTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PatriciaTrieTransactorSession struct {
	Contract     *PatriciaTrieTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// PatriciaTrieRaw is an auto generated low-level Go binding around an Ethereum contract.
type PatriciaTrieRaw struct {
	Contract *PatriciaTrie // Generic contract binding to access the raw methods on
}

// PatriciaTrieCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PatriciaTrieCallerRaw struct {
	Contract *PatriciaTrieCaller // Generic read-only contract binding to access the raw methods on
}

// PatriciaTrieTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PatriciaTrieTransactorRaw struct {
	Contract *PatriciaTrieTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPatriciaTrie creates a new instance of PatriciaTrie, bound to a specific deployed contract.
func NewPatriciaTrie(address common.Address, backend bind.ContractBackend) (*PatriciaTrie, error) {
	contract, err := bindPatriciaTrie(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PatriciaTrie{PatriciaTrieCaller: PatriciaTrieCaller{contract: contract}, PatriciaTrieTransactor: PatriciaTrieTransactor{contract: contract}, PatriciaTrieFilterer: PatriciaTrieFilterer{contract: contract}}, nil
}

// NewPatriciaTrieCaller creates a new read-only instance of PatriciaTrie, bound to a specific deployed contract.
func NewPatriciaTrieCaller(address common.Address, caller bind.ContractCaller) (*PatriciaTrieCaller, error) {
	contract, err := bindPatriciaTrie(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PatriciaTrieCaller{contract: contract}, nil
}

// NewPatriciaTrieTransactor creates a new write-only instance of PatriciaTrie, bound to a specific deployed contract.
func NewPatriciaTrieTransactor(address common.Address, transactor bind.ContractTransactor) (*PatriciaTrieTransactor, error) {
	contract, err := bindPatriciaTrie(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PatriciaTrieTransactor{contract: contract}, nil
}

// NewPatriciaTrieFilterer creates a new log filterer instance of PatriciaTrie, bound to a specific deployed contract.
func NewPatriciaTrieFilterer(address common.Address, filterer bind.ContractFilterer) (*PatriciaTrieFilterer, error) {
	contract, err := bindPatriciaTrie(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PatriciaTrieFilterer{contract: contract}, nil
}

// bindPatriciaTrie binds a generic wrapper to an already deployed contract.
func bindPatriciaTrie(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PatriciaTrieABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PatriciaTrie *PatriciaTrieRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _PatriciaTrie.Contract.PatriciaTrieCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PatriciaTrie *PatriciaTrieRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PatriciaTrie.Contract.PatriciaTrieTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PatriciaTrie *PatriciaTrieRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PatriciaTrie.Contract.PatriciaTrieTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PatriciaTrie *PatriciaTrieCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _PatriciaTrie.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PatriciaTrie *PatriciaTrieTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PatriciaTrie.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PatriciaTrie *PatriciaTrieTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PatriciaTrie.Contract.contract.Transact(opts, method, params...)
}

// VerifyProof is a free data retrieval call binding the contract method 0x4f7142ad.
//
// Solidity: function verifyProof(_value bytes, _parentNodes bytes, _path bytes, _root bytes32) constant returns(bool)
func (_PatriciaTrie *PatriciaTrieCaller) VerifyProof(opts *bind.CallOpts, _value []byte, _parentNodes []byte, _path []byte, _root [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _PatriciaTrie.contract.Call(opts, out, "verifyProof", _value, _parentNodes, _path, _root)
	return *ret0, err
}

// VerifyProof is a free data retrieval call binding the contract method 0x4f7142ad.
//
// Solidity: function verifyProof(_value bytes, _parentNodes bytes, _path bytes, _root bytes32) constant returns(bool)
func (_PatriciaTrie *PatriciaTrieSession) VerifyProof(_value []byte, _parentNodes []byte, _path []byte, _root [32]byte) (bool, error) {
	return _PatriciaTrie.Contract.VerifyProof(&_PatriciaTrie.CallOpts, _value, _parentNodes, _path, _root)
}

// VerifyProof is a free data retrieval call binding the contract method 0x4f7142ad.
//
// Solidity: function verifyProof(_value bytes, _parentNodes bytes, _path bytes, _root bytes32) constant returns(bool)
func (_PatriciaTrie *PatriciaTrieCallerSession) VerifyProof(_value []byte, _parentNodes []byte, _path []byte, _root [32]byte) (bool, error) {
	return _PatriciaTrie.Contract.VerifyProof(&_PatriciaTrie.CallOpts, _value, _parentNodes, _path, _root)
}
