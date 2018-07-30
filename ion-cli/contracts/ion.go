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

// IonABI is the input ABI used to generate the binding from.
const IonABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"registeredChains\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"validation_addr\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"m_blockhashes\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_id\",\"type\":\"bytes32\"},{\"name\":\"_rlpBlockHeader\",\"type\":\"bytes\"},{\"name\":\"_rlpSignedBlockHeader\",\"type\":\"bytes\"}],\"name\":\"SubmitBlock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"m_validators\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_id\",\"type\":\"bytes32\"},{\"name\":\"validationAddr\",\"type\":\"address\"},{\"name\":\"_validators\",\"type\":\"address[]\"},{\"name\":\"_genesisHash\",\"type\":\"bytes32\"}],\"name\":\"RegisterChain\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"chainId\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_id\",\"type\":\"bytes32\"},{\"name\":\"_blockHash\",\"type\":\"bytes32\"},{\"name\":\"_value\",\"type\":\"bytes\"},{\"name\":\"_parentNodes\",\"type\":\"bytes\"},{\"name\":\"_path\",\"type\":\"bytes\"}],\"name\":\"CheckTxProof\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_id\",\"type\":\"bytes32\"},{\"name\":\"_blockHash\",\"type\":\"bytes32\"},{\"name\":\"_value\",\"type\":\"bytes\"},{\"name\":\"_parentNodes\",\"type\":\"bytes\"},{\"name\":\"_path\",\"type\":\"bytes\"}],\"name\":\"CheckReceiptProof\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"chains\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"m_blockheaders\",\"outputs\":[{\"name\":\"blockHeight\",\"type\":\"uint256\"},{\"name\":\"prevBlockHash\",\"type\":\"bytes32\"},{\"name\":\"txRootHash\",\"type\":\"bytes32\"},{\"name\":\"receiptRootHash\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"blockHash\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"blockHeight\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_id\",\"type\":\"bytes32\"},{\"name\":\"_blockHash\",\"type\":\"bytes32\"},{\"name\":\"_txNodes\",\"type\":\"bytes\"},{\"name\":\"_receiptNodes\",\"type\":\"bytes\"}],\"name\":\"CheckRootsProof\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_id\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"chainId\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"proofType\",\"type\":\"uint256\"}],\"name\":\"VerifiedProof\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"broadcastSignature\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"name\":\"broadcastHash\",\"type\":\"event\"}]"

// IonBin is the compiled bytecode used for deploying new contracts.
const IonBin = `608060405234801561001057600080fd5b5060405160208061221f833981018060405281019080805190602001909291905050508060018160001916905550506121d18061004e6000396000f3006080604052600436106100d0576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063200ab0d3146100d5578063261e50731461011e5780634b3550301461018f57806352824374146101e657806353fe62e6146102a35780637558a01a1461030c5780639a8a0592146103ae578063affd8be9146103e1578063bec205b91461050a578063c18de0ef14610633578063e79b136c1461067c578063f22a195e146106fc578063f44ff7121461072f578063f484c1f71461075a575b600080fd5b3480156100e157600080fd5b506101006004803603810190808035906020019092919050505061083d565b60405180826000191660001916815260200191505060405180910390f35b34801561012a57600080fd5b5061014d6004803603810190808035600019169060200190929190505050610860565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561019b57600080fd5b506101cc60048036038101908080356000191690602001909291908035600019169060200190929190505050610893565b604051808215151515815260200191505060405180910390f35b3480156101f257600080fd5b506102a16004803603810190808035600019169060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506108c2565b005b3480156102af57600080fd5b506102f26004803603810190808035600019169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610d71565b604051808215151515815260200191505060405180910390f35b34801561031857600080fd5b506103ac6004803603810190808035600019169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843782019150505050505091929192908035600019169060200190929190505050610da0565b005b3480156103ba57600080fd5b506103c3611124565b60405180826000191660001916815260200191505060405180910390f35b3480156103ed57600080fd5b506104f060048036038101908080356000191690602001909291908035600019169060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050919291929050505061112a565b604051808215151515815260200191505060405180910390f35b34801561051657600080fd5b5061061960048036038101908080356000191690602001909291908035600019169060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050919291929050505061152e565b604051808215151515815260200191505060405180910390f35b34801561063f57600080fd5b506106626004803603810190808035600019169060200190929190505050611932565b604051808215151515815260200191505060405180910390f35b34801561068857600080fd5b506106b960048036038101908080356000191690602001909291908035600019169060200190929190505050611952565b6040518085815260200184600019166000191681526020018360001916600019168152602001826000191660001916815260200194505050505060405180910390f35b34801561070857600080fd5b5061071161198f565b60405180826000191660001916815260200191505060405180910390f35b34801561073b57600080fd5b50610744611995565b6040518082815260200191505060405180910390f35b34801561076657600080fd5b5061082360048036038101908080356000191690602001909291908035600019169060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050919291929050505061199b565b604051808215151515815260200191505060405180910390f35b60028181548110151561084c57fe5b906000526020600020016000915090505481565b60056020528060005260406000206000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60066020528160005260406000206020528060005260406000206000915091509054906101000a900460ff1681565b60608060008060008060008960046000826000191660001916815260200190815260200160002060009054906101000a900460ff16151561096b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f436861696e206973206e6f74207265676973746572656400000000000000000081525060200191505060405180910390fd5b61097c6109778b611bda565b611c35565b975061098f61098a8a611bda565b611c35565b9650600560008c6000191660001916815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1695508594508473ffffffffffffffffffffffffffffffffffffffff166309fbd5368c8c8c6040518463ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004018084600019166000191681526020018060200180602001838103835285818151815260200191508051906020019080838360005b83811015610a6e578082015181840152602081019050610a53565b50505050905090810190601f168015610a9b5780820380516001836020036101000a031916815260200191505b50838103825284818151815260200191508051906020019080838360005b83811015610ad4578082015181840152602081019050610ab9565b50505050905090810190601f168015610b015780820380516001836020036101000a031916815260200191505b5095505050505050600060405180830381600087803b158015610b2357600080fd5b505af1158015610b37573d6000803e3d6000fd5b50505050610b66610b5f896000815181101515610b5057fe5b90602001906020020151611cf2565b6001611d57565b9350886040518082805190602001908083835b602083101515610b9e5780518252602082019150602081019050602083039250610b79565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902092506001600760008d60001916600019168152602001908152602001600020600086600019166000191681526020019081526020016000206000015401915081600760008d60001916600019168152602001908152602001600020600085600019166000191681526020019081526020016000206000018190555083600760008d60001916600019168152602001908152602001600020600085600019166000191681526020019081526020016000206001018160001916905550610cb1610caa896004815181101515610c9b57fe5b90602001906020020151611cf2565b6001611d57565b600760008d60001916600019168152602001908152602001600020600085600019166000191681526020019081526020016000206002018160001916905550610d1b610d14896005815181101515610d0557fe5b90602001906020020151611cf2565b6001611d57565b600760008d60001916600019168152602001908152602001600020600085600019166000191681526020019081526020016000206003018160001916905550610d648b84611d6b565b5050505050505050505050565b60086020528160005260406000206020528060005260406000206000915091509054906101000a900460ff1681565b600060015460001916856000191614151515610e4a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a8152602001807f43616e6e6f7420616464207468697320636861696e20696420746f206368616981526020017f6e2072656769737465720000000000000000000000000000000000000000000081525060400191505060405180910390fd5b60046000866000191660001916815260200190815260200160002060009054906101000a900460ff16151515610ee8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f436861696e20616c72656164792065786973747300000000000000000000000081525060200191505060405180910390fd5b600160046000876000191660001916815260200190815260200160002060006101000a81548160ff021916908315150217905550600285908060018154018082558091505090600182039060005260206000200160009091929091909150906000191690555060016006600087600019166000191681526020019081526020016000206000846000191660001916815260200190815260200160002060006101000a81548160ff021916908315150217905550600060076000876000191660001916815260200190815260200160002060008460001916600019168152602001908152602001600020600001819055508390508073ffffffffffffffffffffffffffffffffffffffff16630b5abdd08685856040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808460001916600019168152602001806020018360001916600019168152602001828103825284818151815260200191508051906020019060200280838360005b83811015611084578082015181840152602081019050611069565b50505050905001945050505050600060405180830381600087803b1580156110ab57600080fd5b505af11580156110bf573d6000803e3d6000fd5b505050508360056000876000191660001916815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050505050565b60015481565b6000808660046000826000191660001916815260200190815260200160002060009054906101000a900460ff1615156111cb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f436861696e206973206e6f74207265676973746572656400000000000000000081525060200191505060405180910390fd5b87876006600083600019166000191681526020019081526020016000206000826000191660001916815260200190815260200160002060009054906101000a900460ff161515611283576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f426c6f636b20646f6573206e6f7420657869737420666f7220636861696e000081525060200191505060405180910390fd5b600760008b6000191660001916815260200190815260200160002060008a60001916600019168152602001908152602001600020935073__./contracts/libraries/PatriciaTrie.s__634f7142ad89898988600201546040518563ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808060200180602001806020018560001916600019168152602001848103845288818151815260200191508051906020019080838360005b8381101561135c578082015181840152602081019050611341565b50505050905090810190601f1680156113895780820380516001836020036101000a031916815260200191505b50848103835287818151815260200191508051906020019080838360005b838110156113c25780820151818401526020810190506113a7565b50505050905090810190601f1680156113ef5780820380516001836020036101000a031916815260200191505b50848103825286818151815260200191508051906020019080838360005b8381101561142857808201518184015260208101905061140d565b50505050905090810190601f1680156114555780820380516001836020036101000a031916815260200191505b5097505050505050505060206040518083038186803b15801561147757600080fd5b505af415801561148b573d6000803e3d6000fd5b505050506040513d60208110156114a157600080fd5b810190808051906020019092919050505015156114ba57fe5b7ff0bc00f5b90f382e1bbca216713ca9e2e8e298f9d7717d30847905395f2870468a8a600060028111156114ea57fe5b6040518084600019166000191681526020018360001916600019168152602001828152602001935050505060405180910390a1600194505050505095945050505050565b6000808660046000826000191660001916815260200190815260200160002060009054906101000a900460ff1615156115cf576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f436861696e206973206e6f74207265676973746572656400000000000000000081525060200191505060405180910390fd5b87876006600083600019166000191681526020019081526020016000206000826000191660001916815260200190815260200160002060009054906101000a900460ff161515611687576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f426c6f636b20646f6573206e6f7420657869737420666f7220636861696e000081525060200191505060405180910390fd5b600760008b6000191660001916815260200190815260200160002060008a60001916600019168152602001908152602001600020935073__./contracts/libraries/PatriciaTrie.s__634f7142ad89898988600301546040518563ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808060200180602001806020018560001916600019168152602001848103845288818151815260200191508051906020019080838360005b83811015611760578082015181840152602081019050611745565b50505050905090810190601f16801561178d5780820380516001836020036101000a031916815260200191505b50848103835287818151815260200191508051906020019080838360005b838110156117c65780820151818401526020810190506117ab565b50505050905090810190601f1680156117f35780820380516001836020036101000a031916815260200191505b50848103825286818151815260200191508051906020019080838360005b8381101561182c578082015181840152602081019050611811565b50505050905090810190601f1680156118595780820380516001836020036101000a031916815260200191505b5097505050505050505060206040518083038186803b15801561187b57600080fd5b505af415801561188f573d6000803e3d6000fd5b505050506040513d60208110156118a557600080fd5b810190808051906020019092919050505015156118be57fe5b7ff0bc00f5b90f382e1bbca216713ca9e2e8e298f9d7717d30847905395f2870468a8a600160028111156118ee57fe5b6040518084600019166000191681526020018360001916600019168152602001828152602001935050505060405180910390a1600194505050505095945050505050565b60046020528060005260406000206000915054906101000a900460ff1681565b6007602052816000526040600020602052806000526040600020600091509150508060000154908060010154908060020154908060030154905084565b60005481565b60035481565b6000808560046000826000191660001916815260200190815260200160002060009054906101000a900460ff161515611a3c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f436861696e206973206e6f74207265676973746572656400000000000000000081525060200191505060405180910390fd5b86866006600083600019166000191681526020019081526020016000206000826000191660001916815260200190815260200160002060009054906101000a900460ff161515611af4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f426c6f636b20646f6573206e6f7420657869737420666f7220636861696e000081525060200191505060405180910390fd5b600760008a60001916600019168152602001908152602001600020600089600019166000191681526020019081526020016000209350611b3387611dbc565b60001916846002015460001916141515611b4957fe5b611b5286611dbc565b60001916846003015460001916141515611b6857fe5b7ff0bc00f5b90f382e1bbca216713ca9e2e8e298f9d7717d30847905395f2870468989600280811115611b9757fe5b6040518084600019166000191681526020018360001916600019168152602001828152602001935050505060405180910390a16001945050505050949350505050565b611be2612150565b600080835191506000821415611c105760408051908101604052806000815260200160008152509250611c2e565b60208401905060408051908101604052808281526020018381525092505b5050919050565b60606000611c4161216a565b6000611c4c85611e6f565b1515611c5757600080fd5b611c6085611ea1565b925082604051908082528060200260200182016040528015611c9c57816020015b611c8961218b565b815260200190600190039081611c815790505b509350611ca885611f1b565b91505b611cb482611f63565b15611cea57611cc282611f8d565b8482815181101515611cd057fe5b906020019060200201819052508080600101915050611cab565b505050919050565b60606000826020015190506000811415611d0b57611d51565b806040519080825280601f01601f191660200182016040528015611d3e5781602001602082028038833980820191505090505b509150611d5083600001518383611fe9565b5b50919050565b600060208201915081830151905092915050565b60016006600084600019166000191681526020019081526020016000206000836000191660001916815260200190815260200160002060006101000a81548160ff0219169083151502179055505050565b6000611dc6612150565b606080611dd285611bda565b9250611ddd83611c35565b9150611e00826000815181101515611df157fe5b90602001906020020151611cf2565b9050806040518082805190602001908083835b602083101515611e385780518252602082019150602081019050602083039250611e13565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390209350505050919050565b600080600083602001511415611e885760009150611e9b565b8260000151905060c0815160001a101591505b50919050565b600080600080600080611eb387611e6f565b1515611ec25760009550611f11565b86600001519350835160001a9450611ed98761202a565b840192506001876020015185010391505b8183111515611f0d57611efc836120b6565b830192508080600101915050611eea565b8095505b5050505050919050565b611f2361216a565b6000611f2e83611e6f565b1515611f3957600080fd5b611f428361202a565b83600001510190508282600001819052508082602001818152505050919050565b6000611f6d612150565b826000015190508060200151816000015101836020015110915050919050565b611f95612150565b600080611fa184611f63565b15611fdd5783602001519150611fb6826120b6565b90508183600001818152505080836020018181525050808201846020018181525050611fe2565b600080fd5b5050919050565b60006020601f83010484602085015b8284146120175760208402808301518183015260018501945050611ff8565b6000865160200187015250505050505050565b6000806000808460200151141561204457600092506120af565b83600001519050805160001a9150608082101561206457600092506120af565b60b8821080612080575060c0821015801561207f575060f882105b5b1561208e57600192506120af565b60c08210156120a557600160b783030192506120af565b600160f783030192505b5050919050565b600080825160001a905060808110156120d2576001915061214a565b60b88110156120e957600160808203019150612149565b60c08110156121135760b78103806020036101000a60018501510480820160010193505050612148565b60f881101561212a57600160c08203019150612147565b60f78103806020036101000a600185015104808201600101935050505b5b5b5b50919050565b604080519081016040528060008152602001600081525090565b60606040519081016040528061217e61218b565b8152602001600081525090565b6040805190810160405280600081526020016000815250905600a165627a7a7230582063bd23fd4aba60ad5fe5003097cfeb975ac723f7e78172be2b48321603b5cabd0029`

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

// BlockHash is a free data retrieval call binding the contract method 0xf22a195e.
//
// Solidity: function blockHash() constant returns(bytes32)
func (_Ion *IonCaller) BlockHash(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Ion.contract.Call(opts, out, "blockHash")
	return *ret0, err
}

// BlockHash is a free data retrieval call binding the contract method 0xf22a195e.
//
// Solidity: function blockHash() constant returns(bytes32)
func (_Ion *IonSession) BlockHash() ([32]byte, error) {
	return _Ion.Contract.BlockHash(&_Ion.CallOpts)
}

// BlockHash is a free data retrieval call binding the contract method 0xf22a195e.
//
// Solidity: function blockHash() constant returns(bytes32)
func (_Ion *IonCallerSession) BlockHash() ([32]byte, error) {
	return _Ion.Contract.BlockHash(&_Ion.CallOpts)
}

// BlockHeight is a free data retrieval call binding the contract method 0xf44ff712.
//
// Solidity: function blockHeight() constant returns(uint256)
func (_Ion *IonCaller) BlockHeight(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Ion.contract.Call(opts, out, "blockHeight")
	return *ret0, err
}

// BlockHeight is a free data retrieval call binding the contract method 0xf44ff712.
//
// Solidity: function blockHeight() constant returns(uint256)
func (_Ion *IonSession) BlockHeight() (*big.Int, error) {
	return _Ion.Contract.BlockHeight(&_Ion.CallOpts)
}

// BlockHeight is a free data retrieval call binding the contract method 0xf44ff712.
//
// Solidity: function blockHeight() constant returns(uint256)
func (_Ion *IonCallerSession) BlockHeight() (*big.Int, error) {
	return _Ion.Contract.BlockHeight(&_Ion.CallOpts)
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

// Chains is a free data retrieval call binding the contract method 0xc18de0ef.
//
// Solidity: function chains( bytes32) constant returns(bool)
func (_Ion *IonCaller) Chains(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Ion.contract.Call(opts, out, "chains", arg0)
	return *ret0, err
}

// Chains is a free data retrieval call binding the contract method 0xc18de0ef.
//
// Solidity: function chains( bytes32) constant returns(bool)
func (_Ion *IonSession) Chains(arg0 [32]byte) (bool, error) {
	return _Ion.Contract.Chains(&_Ion.CallOpts, arg0)
}

// Chains is a free data retrieval call binding the contract method 0xc18de0ef.
//
// Solidity: function chains( bytes32) constant returns(bool)
func (_Ion *IonCallerSession) Chains(arg0 [32]byte) (bool, error) {
	return _Ion.Contract.Chains(&_Ion.CallOpts, arg0)
}

// MBlockhashes is a free data retrieval call binding the contract method 0x4b355030.
//
// Solidity: function m_blockhashes( bytes32,  bytes32) constant returns(bool)
func (_Ion *IonCaller) MBlockhashes(opts *bind.CallOpts, arg0 [32]byte, arg1 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Ion.contract.Call(opts, out, "m_blockhashes", arg0, arg1)
	return *ret0, err
}

// MBlockhashes is a free data retrieval call binding the contract method 0x4b355030.
//
// Solidity: function m_blockhashes( bytes32,  bytes32) constant returns(bool)
func (_Ion *IonSession) MBlockhashes(arg0 [32]byte, arg1 [32]byte) (bool, error) {
	return _Ion.Contract.MBlockhashes(&_Ion.CallOpts, arg0, arg1)
}

// MBlockhashes is a free data retrieval call binding the contract method 0x4b355030.
//
// Solidity: function m_blockhashes( bytes32,  bytes32) constant returns(bool)
func (_Ion *IonCallerSession) MBlockhashes(arg0 [32]byte, arg1 [32]byte) (bool, error) {
	return _Ion.Contract.MBlockhashes(&_Ion.CallOpts, arg0, arg1)
}

// MBlockheaders is a free data retrieval call binding the contract method 0xe79b136c.
//
// Solidity: function m_blockheaders( bytes32,  bytes32) constant returns(blockHeight uint256, prevBlockHash bytes32, txRootHash bytes32, receiptRootHash bytes32)
func (_Ion *IonCaller) MBlockheaders(opts *bind.CallOpts, arg0 [32]byte, arg1 [32]byte) (struct {
	BlockHeight     *big.Int
	PrevBlockHash   [32]byte
	TxRootHash      [32]byte
	ReceiptRootHash [32]byte
}, error) {
	ret := new(struct {
		BlockHeight     *big.Int
		PrevBlockHash   [32]byte
		TxRootHash      [32]byte
		ReceiptRootHash [32]byte
	})
	out := ret
	err := _Ion.contract.Call(opts, out, "m_blockheaders", arg0, arg1)
	return *ret, err
}

// MBlockheaders is a free data retrieval call binding the contract method 0xe79b136c.
//
// Solidity: function m_blockheaders( bytes32,  bytes32) constant returns(blockHeight uint256, prevBlockHash bytes32, txRootHash bytes32, receiptRootHash bytes32)
func (_Ion *IonSession) MBlockheaders(arg0 [32]byte, arg1 [32]byte) (struct {
	BlockHeight     *big.Int
	PrevBlockHash   [32]byte
	TxRootHash      [32]byte
	ReceiptRootHash [32]byte
}, error) {
	return _Ion.Contract.MBlockheaders(&_Ion.CallOpts, arg0, arg1)
}

// MBlockheaders is a free data retrieval call binding the contract method 0xe79b136c.
//
// Solidity: function m_blockheaders( bytes32,  bytes32) constant returns(blockHeight uint256, prevBlockHash bytes32, txRootHash bytes32, receiptRootHash bytes32)
func (_Ion *IonCallerSession) MBlockheaders(arg0 [32]byte, arg1 [32]byte) (struct {
	BlockHeight     *big.Int
	PrevBlockHash   [32]byte
	TxRootHash      [32]byte
	ReceiptRootHash [32]byte
}, error) {
	return _Ion.Contract.MBlockheaders(&_Ion.CallOpts, arg0, arg1)
}

// MValidators is a free data retrieval call binding the contract method 0x53fe62e6.
//
// Solidity: function m_validators( bytes32,  address) constant returns(bool)
func (_Ion *IonCaller) MValidators(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Ion.contract.Call(opts, out, "m_validators", arg0, arg1)
	return *ret0, err
}

// MValidators is a free data retrieval call binding the contract method 0x53fe62e6.
//
// Solidity: function m_validators( bytes32,  address) constant returns(bool)
func (_Ion *IonSession) MValidators(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _Ion.Contract.MValidators(&_Ion.CallOpts, arg0, arg1)
}

// MValidators is a free data retrieval call binding the contract method 0x53fe62e6.
//
// Solidity: function m_validators( bytes32,  address) constant returns(bool)
func (_Ion *IonCallerSession) MValidators(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _Ion.Contract.MValidators(&_Ion.CallOpts, arg0, arg1)
}

// RegisteredChains is a free data retrieval call binding the contract method 0x200ab0d3.
//
// Solidity: function registeredChains( uint256) constant returns(bytes32)
func (_Ion *IonCaller) RegisteredChains(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Ion.contract.Call(opts, out, "registeredChains", arg0)
	return *ret0, err
}

// RegisteredChains is a free data retrieval call binding the contract method 0x200ab0d3.
//
// Solidity: function registeredChains( uint256) constant returns(bytes32)
func (_Ion *IonSession) RegisteredChains(arg0 *big.Int) ([32]byte, error) {
	return _Ion.Contract.RegisteredChains(&_Ion.CallOpts, arg0)
}

// RegisteredChains is a free data retrieval call binding the contract method 0x200ab0d3.
//
// Solidity: function registeredChains( uint256) constant returns(bytes32)
func (_Ion *IonCallerSession) RegisteredChains(arg0 *big.Int) ([32]byte, error) {
	return _Ion.Contract.RegisteredChains(&_Ion.CallOpts, arg0)
}

// ValidationAddr is a free data retrieval call binding the contract method 0x261e5073.
//
// Solidity: function validation_addr( bytes32) constant returns(address)
func (_Ion *IonCaller) ValidationAddr(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Ion.contract.Call(opts, out, "validation_addr", arg0)
	return *ret0, err
}

// ValidationAddr is a free data retrieval call binding the contract method 0x261e5073.
//
// Solidity: function validation_addr( bytes32) constant returns(address)
func (_Ion *IonSession) ValidationAddr(arg0 [32]byte) (common.Address, error) {
	return _Ion.Contract.ValidationAddr(&_Ion.CallOpts, arg0)
}

// ValidationAddr is a free data retrieval call binding the contract method 0x261e5073.
//
// Solidity: function validation_addr( bytes32) constant returns(address)
func (_Ion *IonCallerSession) ValidationAddr(arg0 [32]byte) (common.Address, error) {
	return _Ion.Contract.ValidationAddr(&_Ion.CallOpts, arg0)
}

// CheckReceiptProof is a paid mutator transaction binding the contract method 0xbec205b9.
//
// Solidity: function CheckReceiptProof(_id bytes32, _blockHash bytes32, _value bytes, _parentNodes bytes, _path bytes) returns(bool)
func (_Ion *IonTransactor) CheckReceiptProof(opts *bind.TransactOpts, _id [32]byte, _blockHash [32]byte, _value []byte, _parentNodes []byte, _path []byte) (*types.Transaction, error) {
	return _Ion.contract.Transact(opts, "CheckReceiptProof", _id, _blockHash, _value, _parentNodes, _path)
}

// CheckReceiptProof is a paid mutator transaction binding the contract method 0xbec205b9.
//
// Solidity: function CheckReceiptProof(_id bytes32, _blockHash bytes32, _value bytes, _parentNodes bytes, _path bytes) returns(bool)
func (_Ion *IonSession) CheckReceiptProof(_id [32]byte, _blockHash [32]byte, _value []byte, _parentNodes []byte, _path []byte) (*types.Transaction, error) {
	return _Ion.Contract.CheckReceiptProof(&_Ion.TransactOpts, _id, _blockHash, _value, _parentNodes, _path)
}

// CheckReceiptProof is a paid mutator transaction binding the contract method 0xbec205b9.
//
// Solidity: function CheckReceiptProof(_id bytes32, _blockHash bytes32, _value bytes, _parentNodes bytes, _path bytes) returns(bool)
func (_Ion *IonTransactorSession) CheckReceiptProof(_id [32]byte, _blockHash [32]byte, _value []byte, _parentNodes []byte, _path []byte) (*types.Transaction, error) {
	return _Ion.Contract.CheckReceiptProof(&_Ion.TransactOpts, _id, _blockHash, _value, _parentNodes, _path)
}

// CheckRootsProof is a paid mutator transaction binding the contract method 0xf484c1f7.
//
// Solidity: function CheckRootsProof(_id bytes32, _blockHash bytes32, _txNodes bytes, _receiptNodes bytes) returns(bool)
func (_Ion *IonTransactor) CheckRootsProof(opts *bind.TransactOpts, _id [32]byte, _blockHash [32]byte, _txNodes []byte, _receiptNodes []byte) (*types.Transaction, error) {
	return _Ion.contract.Transact(opts, "CheckRootsProof", _id, _blockHash, _txNodes, _receiptNodes)
}

// CheckRootsProof is a paid mutator transaction binding the contract method 0xf484c1f7.
//
// Solidity: function CheckRootsProof(_id bytes32, _blockHash bytes32, _txNodes bytes, _receiptNodes bytes) returns(bool)
func (_Ion *IonSession) CheckRootsProof(_id [32]byte, _blockHash [32]byte, _txNodes []byte, _receiptNodes []byte) (*types.Transaction, error) {
	return _Ion.Contract.CheckRootsProof(&_Ion.TransactOpts, _id, _blockHash, _txNodes, _receiptNodes)
}

// CheckRootsProof is a paid mutator transaction binding the contract method 0xf484c1f7.
//
// Solidity: function CheckRootsProof(_id bytes32, _blockHash bytes32, _txNodes bytes, _receiptNodes bytes) returns(bool)
func (_Ion *IonTransactorSession) CheckRootsProof(_id [32]byte, _blockHash [32]byte, _txNodes []byte, _receiptNodes []byte) (*types.Transaction, error) {
	return _Ion.Contract.CheckRootsProof(&_Ion.TransactOpts, _id, _blockHash, _txNodes, _receiptNodes)
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

// RegisterChain is a paid mutator transaction binding the contract method 0x7558a01a.
//
// Solidity: function RegisterChain(_id bytes32, validationAddr address, _validators address[], _genesisHash bytes32) returns()
func (_Ion *IonTransactor) RegisterChain(opts *bind.TransactOpts, _id [32]byte, validationAddr common.Address, _validators []common.Address, _genesisHash [32]byte) (*types.Transaction, error) {
	return _Ion.contract.Transact(opts, "RegisterChain", _id, validationAddr, _validators, _genesisHash)
}

// RegisterChain is a paid mutator transaction binding the contract method 0x7558a01a.
//
// Solidity: function RegisterChain(_id bytes32, validationAddr address, _validators address[], _genesisHash bytes32) returns()
func (_Ion *IonSession) RegisterChain(_id [32]byte, validationAddr common.Address, _validators []common.Address, _genesisHash [32]byte) (*types.Transaction, error) {
	return _Ion.Contract.RegisterChain(&_Ion.TransactOpts, _id, validationAddr, _validators, _genesisHash)
}

// RegisterChain is a paid mutator transaction binding the contract method 0x7558a01a.
//
// Solidity: function RegisterChain(_id bytes32, validationAddr address, _validators address[], _genesisHash bytes32) returns()
func (_Ion *IonTransactorSession) RegisterChain(_id [32]byte, validationAddr common.Address, _validators []common.Address, _genesisHash [32]byte) (*types.Transaction, error) {
	return _Ion.Contract.RegisterChain(&_Ion.TransactOpts, _id, validationAddr, _validators, _genesisHash)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0x52824374.
//
// Solidity: function SubmitBlock(_id bytes32, _rlpBlockHeader bytes, _rlpSignedBlockHeader bytes) returns()
func (_Ion *IonTransactor) SubmitBlock(opts *bind.TransactOpts, _id [32]byte, _rlpBlockHeader []byte, _rlpSignedBlockHeader []byte) (*types.Transaction, error) {
	return _Ion.contract.Transact(opts, "SubmitBlock", _id, _rlpBlockHeader, _rlpSignedBlockHeader)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0x52824374.
//
// Solidity: function SubmitBlock(_id bytes32, _rlpBlockHeader bytes, _rlpSignedBlockHeader bytes) returns()
func (_Ion *IonSession) SubmitBlock(_id [32]byte, _rlpBlockHeader []byte, _rlpSignedBlockHeader []byte) (*types.Transaction, error) {
	return _Ion.Contract.SubmitBlock(&_Ion.TransactOpts, _id, _rlpBlockHeader, _rlpSignedBlockHeader)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0x52824374.
//
// Solidity: function SubmitBlock(_id bytes32, _rlpBlockHeader bytes, _rlpSignedBlockHeader bytes) returns()
func (_Ion *IonTransactorSession) SubmitBlock(_id [32]byte, _rlpBlockHeader []byte, _rlpSignedBlockHeader []byte) (*types.Transaction, error) {
	return _Ion.Contract.SubmitBlock(&_Ion.TransactOpts, _id, _rlpBlockHeader, _rlpSignedBlockHeader)
}

// IonVerifiedProofIterator is returned from FilterVerifiedProof and is used to iterate over the raw logs and unpacked data for VerifiedProof events raised by the Ion contract.
type IonVerifiedProofIterator struct {
	Event *IonVerifiedProof // Event containing the contract specifics and raw log

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
func (it *IonVerifiedProofIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IonVerifiedProof)
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
		it.Event = new(IonVerifiedProof)
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
func (it *IonVerifiedProofIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IonVerifiedProofIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IonVerifiedProof represents a VerifiedProof event raised by the Ion contract.
type IonVerifiedProof struct {
	ChainId   [32]byte
	BlockHash [32]byte
	ProofType *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterVerifiedProof is a free log retrieval operation binding the contract event 0xf0bc00f5b90f382e1bbca216713ca9e2e8e298f9d7717d30847905395f287046.
//
// Solidity: e VerifiedProof(chainId bytes32, blockHash bytes32, proofType uint256)
func (_Ion *IonFilterer) FilterVerifiedProof(opts *bind.FilterOpts) (*IonVerifiedProofIterator, error) {

	logs, sub, err := _Ion.contract.FilterLogs(opts, "VerifiedProof")
	if err != nil {
		return nil, err
	}
	return &IonVerifiedProofIterator{contract: _Ion.contract, event: "VerifiedProof", logs: logs, sub: sub}, nil
}

// WatchVerifiedProof is a free log subscription operation binding the contract event 0xf0bc00f5b90f382e1bbca216713ca9e2e8e298f9d7717d30847905395f287046.
//
// Solidity: e VerifiedProof(chainId bytes32, blockHash bytes32, proofType uint256)
func (_Ion *IonFilterer) WatchVerifiedProof(opts *bind.WatchOpts, sink chan<- *IonVerifiedProof) (event.Subscription, error) {

	logs, sub, err := _Ion.contract.WatchLogs(opts, "VerifiedProof")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IonVerifiedProof)
				if err := _Ion.contract.UnpackLog(event, "VerifiedProof", log); err != nil {
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

// IonBroadcastHashIterator is returned from FilterBroadcastHash and is used to iterate over the raw logs and unpacked data for BroadcastHash events raised by the Ion contract.
type IonBroadcastHashIterator struct {
	Event *IonBroadcastHash // Event containing the contract specifics and raw log

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
func (it *IonBroadcastHashIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IonBroadcastHash)
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
		it.Event = new(IonBroadcastHash)
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
func (it *IonBroadcastHashIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IonBroadcastHashIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IonBroadcastHash represents a BroadcastHash event raised by the Ion contract.
type IonBroadcastHash struct {
	BlockHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBroadcastHash is a free log retrieval operation binding the contract event 0xcd7ee33e1a630d6301d87631aab1d4ddce7e1942593cd2689aa989f76d67cf01.
//
// Solidity: e broadcastHash(blockHash bytes32)
func (_Ion *IonFilterer) FilterBroadcastHash(opts *bind.FilterOpts) (*IonBroadcastHashIterator, error) {

	logs, sub, err := _Ion.contract.FilterLogs(opts, "broadcastHash")
	if err != nil {
		return nil, err
	}
	return &IonBroadcastHashIterator{contract: _Ion.contract, event: "broadcastHash", logs: logs, sub: sub}, nil
}

// WatchBroadcastHash is a free log subscription operation binding the contract event 0xcd7ee33e1a630d6301d87631aab1d4ddce7e1942593cd2689aa989f76d67cf01.
//
// Solidity: e broadcastHash(blockHash bytes32)
func (_Ion *IonFilterer) WatchBroadcastHash(opts *bind.WatchOpts, sink chan<- *IonBroadcastHash) (event.Subscription, error) {

	logs, sub, err := _Ion.contract.WatchLogs(opts, "broadcastHash")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IonBroadcastHash)
				if err := _Ion.contract.UnpackLog(event, "broadcastHash", log); err != nil {
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

// IonBroadcastSignatureIterator is returned from FilterBroadcastSignature and is used to iterate over the raw logs and unpacked data for BroadcastSignature events raised by the Ion contract.
type IonBroadcastSignatureIterator struct {
	Event *IonBroadcastSignature // Event containing the contract specifics and raw log

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
func (it *IonBroadcastSignatureIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IonBroadcastSignature)
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
		it.Event = new(IonBroadcastSignature)
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
func (it *IonBroadcastSignatureIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IonBroadcastSignatureIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IonBroadcastSignature represents a BroadcastSignature event raised by the Ion contract.
type IonBroadcastSignature struct {
	Signer common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBroadcastSignature is a free log retrieval operation binding the contract event 0x3fc7485a49f5212e355cd85d3fa5518c2b1e29bcc613ab74cc722a2c5b1eef43.
//
// Solidity: e broadcastSignature(signer address)
func (_Ion *IonFilterer) FilterBroadcastSignature(opts *bind.FilterOpts) (*IonBroadcastSignatureIterator, error) {

	logs, sub, err := _Ion.contract.FilterLogs(opts, "broadcastSignature")
	if err != nil {
		return nil, err
	}
	return &IonBroadcastSignatureIterator{contract: _Ion.contract, event: "broadcastSignature", logs: logs, sub: sub}, nil
}

// WatchBroadcastSignature is a free log subscription operation binding the contract event 0x3fc7485a49f5212e355cd85d3fa5518c2b1e29bcc613ab74cc722a2c5b1eef43.
//
// Solidity: e broadcastSignature(signer address)
func (_Ion *IonFilterer) WatchBroadcastSignature(opts *bind.WatchOpts, sink chan<- *IonBroadcastSignature) (event.Subscription, error) {

	logs, sub, err := _Ion.contract.WatchLogs(opts, "broadcastSignature")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IonBroadcastSignature)
				if err := _Ion.contract.UnpackLog(event, "broadcastSignature", log); err != nil {
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
