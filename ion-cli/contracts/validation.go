// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// ValidationABI is the input ABI used to generate the binding from.
const ValidationABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validators\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"GetBlock\",\"outputs\":[{\"name\":\"_blockHash\",\"type\":\"bytes32\"},{\"name\":\"_prevBlockHash\",\"type\":\"bytes32\"},{\"name\":\"_blockHeight\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"LatestBlock\",\"outputs\":[{\"name\":\"_latestBlock\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_header\",\"type\":\"bytes\"},{\"name\":\"prefixHeader\",\"type\":\"bytes\"},{\"name\":\"prefixExtraData\",\"type\":\"bytes\"}],\"name\":\"ValidateBlock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"m_blockheaders\",\"outputs\":[{\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"name\":\"prevBlockHash\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"GetValidators\",\"outputs\":[{\"name\":\"_validators\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_validators\",\"type\":\"address[]\"},{\"name\":\"genesisHash\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"broadcastSig\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"header\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"parentHash\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"rootHash\",\"type\":\"bytes\"}],\"name\":\"broadcastHashData\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"name\":\"broadcastHash\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"blockHash\",\"type\":\"bytes\"}],\"name\":\"broadcastHash2\",\"type\":\"event\"}]"

// ValidationBin is the compiled bytecode used for deploying new contracts.
const ValidationBin = `608060405234801561001057600080fd5b50604051610ca4380380610ca483398101806040528101908080518201929190602001805190602001909291905050506000336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600090505b825181101561018957600183828151811015156100a057fe5b9060200190602002015190806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050600160056000858481518110151561012157fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508080600101915050610087565b816003816000191690555081600460008081526020019081526020016000206000018160001916905550505050610adf806101c56000396000f300608060405260043610610078576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806335aa2e441461007d57806340ce4b6c146100ea5780639db7d9f714610149578063aae933e41461017c578063c4cf3ed614610271578063d405af3d146102c9575b600080fd5b34801561008957600080fd5b506100a860048036038101908080359060200190929190505050610335565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156100f657600080fd5b5061011560048036038101908080359060200190929190505050610373565b6040518084600019166000191681526020018360001916600019168152602001828152602001935050505060405180910390f35b34801561015557600080fd5b5061015e6103b6565b60405180826000191660001916815260200191505060405180910390f35b34801561018857600080fd5b5061026f600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506103c0565b005b34801561027d57600080fd5b5061029c60048036038101908080359060200190929190505050610857565b60405180836000191660001916815260200182600019166000191681526020019250505060405180910390f35b3480156102d557600080fd5b506102de61087b565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b83811015610321578082015181840152602081019050610306565b505050509050019250505060405180910390f35b60018181548110151561034457fe5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000806000600460008581526020019081526020016000206000015492506004600085815260200190815260200160002060010154915060025490509193909250565b6000600354905090565b600080600060608060608060606000806103db8d6004610909565b9950600354600019168a6000191614151561045e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f4e6f74206368696c64206f662070726576696f757320626c6f636b210000000081525060200191505060405180910390fd5b8c5198508c6040518082805190602001908083835b6020831015156104985780518252602082019150602081019050602083039250610473565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902097507fcd7ee33e1a630d6301d87631aab1d4ddce7e1942593cd2689aa989f76d67cf018860405180826000191660001916815260200191505060405180910390a1608d89036040519080825280601f01601f19166020018201604052801561053c5781602001602082028038833980820191505090505b509650601f6040519080825280601f01601f1916602001820160405280156105735781602001602082028038833980820191505090505b50955060416040519080825280601f01601f1916602001820160405280156105aa5781602001602082028038833980820191505090505b509450602a6040519080825280601f01601f1916602001820160405280156105e15781602001602082028038833980820191505090505b5093506105f0878e600061091d565b600260218801600260208f016004610bb8fa50610611868e608c8c0361091d565b600160208701600160208e016004610bb8fa50610632848e602a8c0361091d565b61063d87878661093f565b9250826040518082805190602001908083835b6020831015156106755780518252602082019150602081019050602083039250610650565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902091506106b2858e606b8c0361091d565b6106bc82866109eb565b905060011515600560008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161515141515610786576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f5369676e6572206e6f7420612076616c696461746f722100000000000000000081525060200191505060405180910390fd5b87600381600019169055506002600081548092919060010191905055508760046000600254815260200190815260200160002060000181600019169055508960046000600254815260200190815260200160002060010181600019169055507fba2fe28067a0918af64c5359b1579f887bf1479dd3163c7e5d456314168854a581604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a150505050505050505050505050565b60046020528060005260406000206000915090508060000154908060010154905082565b606060018054806020026020016040519081016040528092919081815260200182805480156108ff57602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190600101908083116108b5575b5050505050905090565b600060208201915081830151905092915050565b6000835190506020820191508060208501828486016004610bb8fa5050505050565b606060008060008060008060608a5196508951955060208701945088519350600260208888010101925060028487890101019150816040519080825280601f01601f1916602001820160405280156109a65781602001602082028038833980820191505090505b50905086602082018860208e016004610bb8fa50858582018760208d016004610bb8fa50838382018560208c016004610bb8fa50809750505050505050509392505050565b60008060008060418551141515610a0157600080fd5b6020850151925060408501519150606085015160001a9050601b8160ff161015610a2c57601b810190505b601b8160ff161480610a415750601c8160ff16145b1515610a4c57600080fd5b610a5886828585610a63565b935050505092915050565b60008060006040518781528660208201528560408201528460608201526020816080836001610bb8fa9250805191505060011515821515141515610aa657600080fd5b80925050509493505050505600a165627a7a72305820cf83fdc91c1477a7ff4a9251338903b37b04c10ccb8538af298dd16fece5a0ba0029`

// DeployValidation deploys a new Ethereum contract, binding an instance of Validation to it.
func DeployValidation(auth *bind.TransactOpts, backend bind.ContractBackend, _validators []common.Address, genesisHash [32]byte) (common.Address, *types.Transaction, *Validation, error) {
	parsed, err := abi.JSON(strings.NewReader(ValidationABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ValidationBin), backend, _validators, genesisHash)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Validation{ValidationCaller: ValidationCaller{contract: contract}, ValidationTransactor: ValidationTransactor{contract: contract}, ValidationFilterer: ValidationFilterer{contract: contract}}, nil
}

// Validation is an auto generated Go binding around an Ethereum contract.
type Validation struct {
	ValidationCaller     // Read-only binding to the contract
	ValidationTransactor // Write-only binding to the contract
	ValidationFilterer   // Log filterer for contract events
}

// ValidationCaller is an auto generated read-only Go binding around an Ethereum contract.
type ValidationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ValidationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ValidationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ValidationSession struct {
	Contract     *Validation       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ValidationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ValidationCallerSession struct {
	Contract *ValidationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ValidationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ValidationTransactorSession struct {
	Contract     *ValidationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ValidationRaw is an auto generated low-level Go binding around an Ethereum contract.
type ValidationRaw struct {
	Contract *Validation // Generic contract binding to access the raw methods on
}

// ValidationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ValidationCallerRaw struct {
	Contract *ValidationCaller // Generic read-only contract binding to access the raw methods on
}

// ValidationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ValidationTransactorRaw struct {
	Contract *ValidationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewValidation creates a new instance of Validation, bound to a specific deployed contract.
func NewValidation(address common.Address, backend bind.ContractBackend) (*Validation, error) {
	contract, err := bindValidation(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Validation{ValidationCaller: ValidationCaller{contract: contract}, ValidationTransactor: ValidationTransactor{contract: contract}, ValidationFilterer: ValidationFilterer{contract: contract}}, nil
}

// NewValidationCaller creates a new read-only instance of Validation, bound to a specific deployed contract.
func NewValidationCaller(address common.Address, caller bind.ContractCaller) (*ValidationCaller, error) {
	contract, err := bindValidation(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ValidationCaller{contract: contract}, nil
}

// NewValidationTransactor creates a new write-only instance of Validation, bound to a specific deployed contract.
func NewValidationTransactor(address common.Address, transactor bind.ContractTransactor) (*ValidationTransactor, error) {
	contract, err := bindValidation(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ValidationTransactor{contract: contract}, nil
}

// NewValidationFilterer creates a new log filterer instance of Validation, bound to a specific deployed contract.
func NewValidationFilterer(address common.Address, filterer bind.ContractFilterer) (*ValidationFilterer, error) {
	contract, err := bindValidation(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ValidationFilterer{contract: contract}, nil
}

// bindValidation binds a generic wrapper to an already deployed contract.
func bindValidation(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ValidationABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Validation *ValidationRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Validation.Contract.ValidationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Validation *ValidationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validation.Contract.ValidationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Validation *ValidationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Validation.Contract.ValidationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Validation *ValidationCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Validation.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Validation *ValidationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validation.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Validation *ValidationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Validation.Contract.contract.Transact(opts, method, params...)
}

// GetBlock is a free data retrieval call binding the contract method 0x40ce4b6c.
//
// Solidity: function GetBlock(blockNumber uint256) constant returns(_blockHash bytes32, _prevBlockHash bytes32, _blockHeight uint256)
func (_Validation *ValidationCaller) GetBlock(opts *bind.CallOpts, blockNumber *big.Int) (struct {
	BlockHash     [32]byte
	PrevBlockHash [32]byte
	BlockHeight   *big.Int
}, error) {
	ret := new(struct {
		BlockHash     [32]byte
		PrevBlockHash [32]byte
		BlockHeight   *big.Int
	})
	out := ret
	err := _Validation.contract.Call(opts, out, "GetBlock", blockNumber)
	return *ret, err
}

// GetBlock is a free data retrieval call binding the contract method 0x40ce4b6c.
//
// Solidity: function GetBlock(blockNumber uint256) constant returns(_blockHash bytes32, _prevBlockHash bytes32, _blockHeight uint256)
func (_Validation *ValidationSession) GetBlock(blockNumber *big.Int) (struct {
	BlockHash     [32]byte
	PrevBlockHash [32]byte
	BlockHeight   *big.Int
}, error) {
	return _Validation.Contract.GetBlock(&_Validation.CallOpts, blockNumber)
}

// GetBlock is a free data retrieval call binding the contract method 0x40ce4b6c.
//
// Solidity: function GetBlock(blockNumber uint256) constant returns(_blockHash bytes32, _prevBlockHash bytes32, _blockHeight uint256)
func (_Validation *ValidationCallerSession) GetBlock(blockNumber *big.Int) (struct {
	BlockHash     [32]byte
	PrevBlockHash [32]byte
	BlockHeight   *big.Int
}, error) {
	return _Validation.Contract.GetBlock(&_Validation.CallOpts, blockNumber)
}

// GetValidators is a free data retrieval call binding the contract method 0xd405af3d.
//
// Solidity: function GetValidators() constant returns(_validators address[])
func (_Validation *ValidationCaller) GetValidators(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _Validation.contract.Call(opts, out, "GetValidators")
	return *ret0, err
}

// GetValidators is a free data retrieval call binding the contract method 0xd405af3d.
//
// Solidity: function GetValidators() constant returns(_validators address[])
func (_Validation *ValidationSession) GetValidators() ([]common.Address, error) {
	return _Validation.Contract.GetValidators(&_Validation.CallOpts)
}

// GetValidators is a free data retrieval call binding the contract method 0xd405af3d.
//
// Solidity: function GetValidators() constant returns(_validators address[])
func (_Validation *ValidationCallerSession) GetValidators() ([]common.Address, error) {
	return _Validation.Contract.GetValidators(&_Validation.CallOpts)
}

// LatestBlock is a free data retrieval call binding the contract method 0x9db7d9f7.
//
// Solidity: function LatestBlock() constant returns(_latestBlock bytes32)
func (_Validation *ValidationCaller) LatestBlock(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Validation.contract.Call(opts, out, "LatestBlock")
	return *ret0, err
}

// LatestBlock is a free data retrieval call binding the contract method 0x9db7d9f7.
//
// Solidity: function LatestBlock() constant returns(_latestBlock bytes32)
func (_Validation *ValidationSession) LatestBlock() ([32]byte, error) {
	return _Validation.Contract.LatestBlock(&_Validation.CallOpts)
}

// LatestBlock is a free data retrieval call binding the contract method 0x9db7d9f7.
//
// Solidity: function LatestBlock() constant returns(_latestBlock bytes32)
func (_Validation *ValidationCallerSession) LatestBlock() ([32]byte, error) {
	return _Validation.Contract.LatestBlock(&_Validation.CallOpts)
}

// MBlockheaders is a free data retrieval call binding the contract method 0xc4cf3ed6.
//
// Solidity: function m_blockheaders( uint256) constant returns(blockHash bytes32, prevBlockHash bytes32)
func (_Validation *ValidationCaller) MBlockheaders(opts *bind.CallOpts, arg0 *big.Int) (struct {
	BlockHash     [32]byte
	PrevBlockHash [32]byte
}, error) {
	ret := new(struct {
		BlockHash     [32]byte
		PrevBlockHash [32]byte
	})
	out := ret
	err := _Validation.contract.Call(opts, out, "m_blockheaders", arg0)
	return *ret, err
}

// MBlockheaders is a free data retrieval call binding the contract method 0xc4cf3ed6.
//
// Solidity: function m_blockheaders( uint256) constant returns(blockHash bytes32, prevBlockHash bytes32)
func (_Validation *ValidationSession) MBlockheaders(arg0 *big.Int) (struct {
	BlockHash     [32]byte
	PrevBlockHash [32]byte
}, error) {
	return _Validation.Contract.MBlockheaders(&_Validation.CallOpts, arg0)
}

// MBlockheaders is a free data retrieval call binding the contract method 0xc4cf3ed6.
//
// Solidity: function m_blockheaders( uint256) constant returns(blockHash bytes32, prevBlockHash bytes32)
func (_Validation *ValidationCallerSession) MBlockheaders(arg0 *big.Int) (struct {
	BlockHash     [32]byte
	PrevBlockHash [32]byte
}, error) {
	return _Validation.Contract.MBlockheaders(&_Validation.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators( uint256) constant returns(address)
func (_Validation *ValidationCaller) Validators(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Validation.contract.Call(opts, out, "validators", arg0)
	return *ret0, err
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators( uint256) constant returns(address)
func (_Validation *ValidationSession) Validators(arg0 *big.Int) (common.Address, error) {
	return _Validation.Contract.Validators(&_Validation.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators( uint256) constant returns(address)
func (_Validation *ValidationCallerSession) Validators(arg0 *big.Int) (common.Address, error) {
	return _Validation.Contract.Validators(&_Validation.CallOpts, arg0)
}

// ValidateBlock is a paid mutator transaction binding the contract method 0xaae933e4.
//
// Solidity: function ValidateBlock(_header bytes, prefixHeader bytes, prefixExtraData bytes) returns()
func (_Validation *ValidationTransactor) ValidateBlock(opts *bind.TransactOpts, _header []byte, prefixHeader []byte, prefixExtraData []byte) (*types.Transaction, error) {
	return _Validation.contract.Transact(opts, "ValidateBlock", _header, prefixHeader, prefixExtraData)
}

// ValidateBlock is a paid mutator transaction binding the contract method 0xaae933e4.
//
// Solidity: function ValidateBlock(_header bytes, prefixHeader bytes, prefixExtraData bytes) returns()
func (_Validation *ValidationSession) ValidateBlock(_header []byte, prefixHeader []byte, prefixExtraData []byte) (*types.Transaction, error) {
	return _Validation.Contract.ValidateBlock(&_Validation.TransactOpts, _header, prefixHeader, prefixExtraData)
}

// ValidateBlock is a paid mutator transaction binding the contract method 0xaae933e4.
//
// Solidity: function ValidateBlock(_header bytes, prefixHeader bytes, prefixExtraData bytes) returns()
func (_Validation *ValidationTransactorSession) ValidateBlock(_header []byte, prefixHeader []byte, prefixExtraData []byte) (*types.Transaction, error) {
	return _Validation.Contract.ValidateBlock(&_Validation.TransactOpts, _header, prefixHeader, prefixExtraData)
}

// ValidationBroadcastHashIterator is returned from FilterBroadcastHash and is used to iterate over the raw logs and unpacked data for BroadcastHash events raised by the Validation contract.
type ValidationBroadcastHashIterator struct {
	Event *ValidationBroadcastHash // Event containing the contract specifics and raw log

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
func (it *ValidationBroadcastHashIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidationBroadcastHash)
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
		it.Event = new(ValidationBroadcastHash)
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
func (it *ValidationBroadcastHashIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidationBroadcastHashIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidationBroadcastHash represents a BroadcastHash event raised by the Validation contract.
type ValidationBroadcastHash struct {
	BlockHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBroadcastHash is a free log retrieval operation binding the contract event 0xcd7ee33e1a630d6301d87631aab1d4ddce7e1942593cd2689aa989f76d67cf01.
//
// Solidity: event broadcastHash(blockHash bytes32)
func (_Validation *ValidationFilterer) FilterBroadcastHash(opts *bind.FilterOpts) (*ValidationBroadcastHashIterator, error) {

	logs, sub, err := _Validation.contract.FilterLogs(opts, "broadcastHash")
	if err != nil {
		return nil, err
	}
	return &ValidationBroadcastHashIterator{contract: _Validation.contract, event: "broadcastHash", logs: logs, sub: sub}, nil
}

// WatchBroadcastHash is a free log subscription operation binding the contract event 0xcd7ee33e1a630d6301d87631aab1d4ddce7e1942593cd2689aa989f76d67cf01.
//
// Solidity: event broadcastHash(blockHash bytes32)
func (_Validation *ValidationFilterer) WatchBroadcastHash(opts *bind.WatchOpts, sink chan<- *ValidationBroadcastHash) (event.Subscription, error) {

	logs, sub, err := _Validation.contract.WatchLogs(opts, "broadcastHash")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidationBroadcastHash)
				if err := _Validation.contract.UnpackLog(event, "broadcastHash", log); err != nil {
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

// ValidationBroadcastHash2Iterator is returned from FilterBroadcastHash2 and is used to iterate over the raw logs and unpacked data for BroadcastHash2 events raised by the Validation contract.
type ValidationBroadcastHash2Iterator struct {
	Event *ValidationBroadcastHash2 // Event containing the contract specifics and raw log

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
func (it *ValidationBroadcastHash2Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidationBroadcastHash2)
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
		it.Event = new(ValidationBroadcastHash2)
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
func (it *ValidationBroadcastHash2Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidationBroadcastHash2Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidationBroadcastHash2 represents a BroadcastHash2 event raised by the Validation contract.
type ValidationBroadcastHash2 struct {
	BlockHash []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBroadcastHash2 is a free log retrieval operation binding the contract event 0xd5796acb84f214e8b630fe40f93fdbf75c0a53b553e23be97a99e61152f174d2.
//
// Solidity: event broadcastHash2(blockHash bytes)
func (_Validation *ValidationFilterer) FilterBroadcastHash2(opts *bind.FilterOpts) (*ValidationBroadcastHash2Iterator, error) {

	logs, sub, err := _Validation.contract.FilterLogs(opts, "broadcastHash2")
	if err != nil {
		return nil, err
	}
	return &ValidationBroadcastHash2Iterator{contract: _Validation.contract, event: "broadcastHash2", logs: logs, sub: sub}, nil
}

// WatchBroadcastHash2 is a free log subscription operation binding the contract event 0xd5796acb84f214e8b630fe40f93fdbf75c0a53b553e23be97a99e61152f174d2.
//
// Solidity: event broadcastHash2(blockHash bytes)
func (_Validation *ValidationFilterer) WatchBroadcastHash2(opts *bind.WatchOpts, sink chan<- *ValidationBroadcastHash2) (event.Subscription, error) {

	logs, sub, err := _Validation.contract.WatchLogs(opts, "broadcastHash2")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidationBroadcastHash2)
				if err := _Validation.contract.UnpackLog(event, "broadcastHash2", log); err != nil {
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

// ValidationBroadcastHashDataIterator is returned from FilterBroadcastHashData and is used to iterate over the raw logs and unpacked data for BroadcastHashData events raised by the Validation contract.
type ValidationBroadcastHashDataIterator struct {
	Event *ValidationBroadcastHashData // Event containing the contract specifics and raw log

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
func (it *ValidationBroadcastHashDataIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidationBroadcastHashData)
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
		it.Event = new(ValidationBroadcastHashData)
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
func (it *ValidationBroadcastHashDataIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidationBroadcastHashDataIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidationBroadcastHashData represents a BroadcastHashData event raised by the Validation contract.
type ValidationBroadcastHashData struct {
	Header     []byte
	ParentHash []byte
	RootHash   []byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBroadcastHashData is a free log retrieval operation binding the contract event 0x8511795469a13c04a2bc22c3f1309fc0bd918a0a25a3e7e222a0417b719274c7.
//
// Solidity: event broadcastHashData(header bytes, parentHash bytes, rootHash bytes)
func (_Validation *ValidationFilterer) FilterBroadcastHashData(opts *bind.FilterOpts) (*ValidationBroadcastHashDataIterator, error) {

	logs, sub, err := _Validation.contract.FilterLogs(opts, "broadcastHashData")
	if err != nil {
		return nil, err
	}
	return &ValidationBroadcastHashDataIterator{contract: _Validation.contract, event: "broadcastHashData", logs: logs, sub: sub}, nil
}

// WatchBroadcastHashData is a free log subscription operation binding the contract event 0x8511795469a13c04a2bc22c3f1309fc0bd918a0a25a3e7e222a0417b719274c7.
//
// Solidity: event broadcastHashData(header bytes, parentHash bytes, rootHash bytes)
func (_Validation *ValidationFilterer) WatchBroadcastHashData(opts *bind.WatchOpts, sink chan<- *ValidationBroadcastHashData) (event.Subscription, error) {

	logs, sub, err := _Validation.contract.WatchLogs(opts, "broadcastHashData")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidationBroadcastHashData)
				if err := _Validation.contract.UnpackLog(event, "broadcastHashData", log); err != nil {
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

// ValidationBroadcastSigIterator is returned from FilterBroadcastSig and is used to iterate over the raw logs and unpacked data for BroadcastSig events raised by the Validation contract.
type ValidationBroadcastSigIterator struct {
	Event *ValidationBroadcastSig // Event containing the contract specifics and raw log

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
func (it *ValidationBroadcastSigIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidationBroadcastSig)
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
		it.Event = new(ValidationBroadcastSig)
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
func (it *ValidationBroadcastSigIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidationBroadcastSigIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidationBroadcastSig represents a BroadcastSig event raised by the Validation contract.
type ValidationBroadcastSig struct {
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterBroadcastSig is a free log retrieval operation binding the contract event 0xba2fe28067a0918af64c5359b1579f887bf1479dd3163c7e5d456314168854a5.
//
// Solidity: event broadcastSig(owner address)
func (_Validation *ValidationFilterer) FilterBroadcastSig(opts *bind.FilterOpts) (*ValidationBroadcastSigIterator, error) {

	logs, sub, err := _Validation.contract.FilterLogs(opts, "broadcastSig")
	if err != nil {
		return nil, err
	}
	return &ValidationBroadcastSigIterator{contract: _Validation.contract, event: "broadcastSig", logs: logs, sub: sub}, nil
}

// WatchBroadcastSig is a free log subscription operation binding the contract event 0xba2fe28067a0918af64c5359b1579f887bf1479dd3163c7e5d456314168854a5.
//
// Solidity: event broadcastSig(owner address)
func (_Validation *ValidationFilterer) WatchBroadcastSig(opts *bind.WatchOpts, sink chan<- *ValidationBroadcastSig) (event.Subscription, error) {

	logs, sub, err := _Validation.contract.WatchLogs(opts, "broadcastSig")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidationBroadcastSig)
				if err := _Validation.contract.UnpackLog(event, "broadcastSig", log); err != nil {
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
