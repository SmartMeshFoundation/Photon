// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// ChannelManagerContractABI is the input ABI used to generate the binding from.
const ChannelManagerContractABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"getChannelsParticipants\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"partner\",\"type\":\"address\"}],\"name\":\"getChannelWith\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getChannelsAddresses\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"node_address\",\"type\":\"address\"}],\"name\":\"nettingContractsByAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"channel\",\"type\":\"address\"}],\"name\":\"contractExists\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"partner\",\"type\":\"address\"},{\"name\":\"settle_timeout\",\"type\":\"uint256\"}],\"name\":\"newChannel\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"registry_address\",\"type\":\"address\"},{\"name\":\"token_address\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"registry_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"netting_channel\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"participant1\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"participant2\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"settle_timeout\",\"type\":\"uint256\"}],\"name\":\"ChannelNew\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"registry_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"caller_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"partner\",\"type\":\"address\"}],\"name\":\"ChannelDeleted\",\"type\":\"event\"}]"

// ChannelManagerContractBin is the compiled bytecode used for deploying new contracts.
const ChannelManagerContractBin = `0x608060405234801561001057600080fd5b5060405160408061084583398101604052805160209091015160018054600160a060020a03938416600160a060020a031991821617909155600080549390921692169190911790556107de806100676000396000f30060806040526004361061008d5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630b74b620811461009f578063238bfba2146101045780636785b500146101415780636cb30fee146101565780637709bc78146101775780639d76ea58146101ac578063b32c65c8146101c1578063f26c6aed1461024b575b34801561009957600080fd5b50600080fd5b3480156100ab57600080fd5b506100b461026f565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156100f05781810151838201526020016100d8565b505050509050019250505060405180910390f35b34801561011057600080fd5b50610125600160a060020a0360043516610457565b60408051600160a060020a039092168252519081900360200190f35b34801561014d57600080fd5b506100b4610505565b34801561016257600080fd5b506100b4600160a060020a036004351661056a565b34801561018357600080fd5b50610198600160a060020a03600435166105e0565b604080519115158252519081900360200190f35b3480156101b857600080fd5b506101256105e8565b3480156101cd57600080fd5b506101d66105f7565b6040805160208082528351818301528351919283929083019185019080838360005b838110156102105781810151838201526020016101f8565b50505050905090810190601f16801561023d5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561025757600080fd5b50610125600160a060020a036004351660243561062e565b606060008082818080805b6002548710156102c557600280546102b191908990811061029757fe5b600091825260209091200154600160a060020a03166105e0565b156102ba576001015b60019096019561027a565b806002026040519080825280602002602001820160405280156102f2578160200160208202803883390190505b50945060009550600096505b60025487101561044b576002805461031b91908990811061029757fe5b151561032657610440565b600280548890811061033457fe5b9060005260206000200160009054906101000a9004600160a060020a0316915081600160a060020a0316636d2381b36040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401608060405180830381600087803b1580156103ab57600080fd5b505af11580156103bf573d6000803e3d6000fd5b505050506040513d60808110156103d557600080fd5b5080516040909101518651919550935084908690889081106103f357fe5b600160a060020a03909216602092830290910190910152845160019690960195839086908890811061042157fe5b600160a060020a03909216602092830290910190910152600195909501945b6001909601956102fe565b50929695505050505050565b604080517f8a1c00e2000000000000000000000000000000000000000000000000000000008152600060048201819052600160a060020a0384166024830152915173__ChannelManagerLibrary.sol:ChannelMan__91638a1c00e2916044808301926020929190829003018186803b1580156104d357600080fd5b505af41580156104e7573d6000803e3d6000fd5b505050506040513d60208110156104fd57600080fd5b505192915050565b6060600060020180548060200260200160405190810160405280929190818152602001828054801561056057602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610542575b5050505050905090565b600160a060020a0381166000908152600460209081526040918290208054835181840281018401909452808452606093928301828280156105d457602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116105b6575b50505050509050919050565b6000903b1190565b600054600160a060020a031690565b60408051808201909152600581527f302e322e5f000000000000000000000000000000000000000000000000000000602082015281565b600080600061063c85610457565b9150600160a060020a0382161561069c5760015460408051600160a060020a039283168152338316602082015291871682820152517f91cd7e9ad7c88602bc7b06adc62de54cb670665c75894414c01ead5f2baf09309181900360600190a15b604080517f941583a500000000000000000000000000000000000000000000000000000000815260006004820152600160a060020a038716602482015260448101869052905173__ChannelManagerLibrary.sol:ChannelMan__9163941583a5916064808301926020929190829003018186803b15801561071d57600080fd5b505af4158015610731573d6000803e3d6000fd5b505050506040513d602081101561074757600080fd5b505160015460408051600160a060020a039283168152828416602082015233831681830152918816606083015260808201879052519192507fc15f36191ba1aefd153ae0ece89afc6004a66b158bec44593e70bccff2ae7a0f919081900360a00190a19493505050505600a165627a7a723058204ca9dc2ef3777511038e010aaee35e13c0e7f876fc9f47fae0f677b65a04cf9c0029`

// DeployChannelManagerContract deploys a new Ethereum contract, binding an instance of ChannelManagerContract to it.
func DeployChannelManagerContract(auth *bind.TransactOpts, backend bind.ContractBackend, registry_address common.Address, token_address common.Address) (common.Address, *types.Transaction, *ChannelManagerContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ChannelManagerContractABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ChannelManagerContractBin), backend, registry_address, token_address)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ChannelManagerContract{ChannelManagerContractCaller: ChannelManagerContractCaller{contract: contract}, ChannelManagerContractTransactor: ChannelManagerContractTransactor{contract: contract}, ChannelManagerContractFilterer: ChannelManagerContractFilterer{contract: contract}}, nil
}

// ChannelManagerContract is an auto generated Go binding around an Ethereum contract.
type ChannelManagerContract struct {
	ChannelManagerContractCaller     // Read-only binding to the contract
	ChannelManagerContractTransactor // Write-only binding to the contract
	ChannelManagerContractFilterer   // Log filterer for contract events
}

// ChannelManagerContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChannelManagerContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChannelManagerContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChannelManagerContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChannelManagerContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChannelManagerContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChannelManagerContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChannelManagerContractSession struct {
	Contract     *ChannelManagerContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ChannelManagerContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChannelManagerContractCallerSession struct {
	Contract *ChannelManagerContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// ChannelManagerContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChannelManagerContractTransactorSession struct {
	Contract     *ChannelManagerContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// ChannelManagerContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ChannelManagerContractRaw struct {
	Contract *ChannelManagerContract // Generic contract binding to access the raw methods on
}

// ChannelManagerContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChannelManagerContractCallerRaw struct {
	Contract *ChannelManagerContractCaller // Generic read-only contract binding to access the raw methods on
}

// ChannelManagerContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChannelManagerContractTransactorRaw struct {
	Contract *ChannelManagerContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChannelManagerContract creates a new instance of ChannelManagerContract, bound to a specific deployed contract.
func NewChannelManagerContract(address common.Address, backend bind.ContractBackend) (*ChannelManagerContract, error) {
	contract, err := bindChannelManagerContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChannelManagerContract{ChannelManagerContractCaller: ChannelManagerContractCaller{contract: contract}, ChannelManagerContractTransactor: ChannelManagerContractTransactor{contract: contract}, ChannelManagerContractFilterer: ChannelManagerContractFilterer{contract: contract}}, nil
}

// NewChannelManagerContractCaller creates a new read-only instance of ChannelManagerContract, bound to a specific deployed contract.
func NewChannelManagerContractCaller(address common.Address, caller bind.ContractCaller) (*ChannelManagerContractCaller, error) {
	contract, err := bindChannelManagerContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChannelManagerContractCaller{contract: contract}, nil
}

// NewChannelManagerContractTransactor creates a new write-only instance of ChannelManagerContract, bound to a specific deployed contract.
func NewChannelManagerContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ChannelManagerContractTransactor, error) {
	contract, err := bindChannelManagerContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChannelManagerContractTransactor{contract: contract}, nil
}

// NewChannelManagerContractFilterer creates a new log filterer instance of ChannelManagerContract, bound to a specific deployed contract.
func NewChannelManagerContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ChannelManagerContractFilterer, error) {
	contract, err := bindChannelManagerContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChannelManagerContractFilterer{contract: contract}, nil
}

// bindChannelManagerContract binds a generic wrapper to an already deployed contract.
func bindChannelManagerContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ChannelManagerContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChannelManagerContract *ChannelManagerContractRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ChannelManagerContract.Contract.ChannelManagerContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChannelManagerContract *ChannelManagerContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChannelManagerContract.Contract.ChannelManagerContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChannelManagerContract *ChannelManagerContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChannelManagerContract.Contract.ChannelManagerContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChannelManagerContract *ChannelManagerContractCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ChannelManagerContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChannelManagerContract *ChannelManagerContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChannelManagerContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChannelManagerContract *ChannelManagerContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChannelManagerContract.Contract.contract.Transact(opts, method, params...)
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(channel address) constant returns(bool)
func (_ChannelManagerContract *ChannelManagerContractCaller) ContractExists(opts *bind.CallOpts, channel common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ChannelManagerContract.contract.Call(opts, out, "contractExists", channel)
	return *ret0, err
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(channel address) constant returns(bool)
func (_ChannelManagerContract *ChannelManagerContractSession) ContractExists(channel common.Address) (bool, error) {
	return _ChannelManagerContract.Contract.ContractExists(&_ChannelManagerContract.CallOpts, channel)
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(channel address) constant returns(bool)
func (_ChannelManagerContract *ChannelManagerContractCallerSession) ContractExists(channel common.Address) (bool, error) {
	return _ChannelManagerContract.Contract.ContractExists(&_ChannelManagerContract.CallOpts, channel)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_ChannelManagerContract *ChannelManagerContractCaller) Contract_version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ChannelManagerContract.contract.Call(opts, out, "contract_version")
	return *ret0, err
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_ChannelManagerContract *ChannelManagerContractSession) Contract_version() (string, error) {
	return _ChannelManagerContract.Contract.Contract_version(&_ChannelManagerContract.CallOpts)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_ChannelManagerContract *ChannelManagerContractCallerSession) Contract_version() (string, error) {
	return _ChannelManagerContract.Contract.Contract_version(&_ChannelManagerContract.CallOpts)
}

// GetChannelWith is a free data retrieval call binding the contract method 0x238bfba2.
//
// Solidity: function getChannelWith(partner address) constant returns(address)
func (_ChannelManagerContract *ChannelManagerContractCaller) GetChannelWith(opts *bind.CallOpts, partner common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ChannelManagerContract.contract.Call(opts, out, "getChannelWith", partner)
	return *ret0, err
}

// GetChannelWith is a free data retrieval call binding the contract method 0x238bfba2.
//
// Solidity: function getChannelWith(partner address) constant returns(address)
func (_ChannelManagerContract *ChannelManagerContractSession) GetChannelWith(partner common.Address) (common.Address, error) {
	return _ChannelManagerContract.Contract.GetChannelWith(&_ChannelManagerContract.CallOpts, partner)
}

// GetChannelWith is a free data retrieval call binding the contract method 0x238bfba2.
//
// Solidity: function getChannelWith(partner address) constant returns(address)
func (_ChannelManagerContract *ChannelManagerContractCallerSession) GetChannelWith(partner common.Address) (common.Address, error) {
	return _ChannelManagerContract.Contract.GetChannelWith(&_ChannelManagerContract.CallOpts, partner)
}

// GetChannelsAddresses is a free data retrieval call binding the contract method 0x6785b500.
//
// Solidity: function getChannelsAddresses() constant returns(address[])
func (_ChannelManagerContract *ChannelManagerContractCaller) GetChannelsAddresses(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _ChannelManagerContract.contract.Call(opts, out, "getChannelsAddresses")
	return *ret0, err
}

// GetChannelsAddresses is a free data retrieval call binding the contract method 0x6785b500.
//
// Solidity: function getChannelsAddresses() constant returns(address[])
func (_ChannelManagerContract *ChannelManagerContractSession) GetChannelsAddresses() ([]common.Address, error) {
	return _ChannelManagerContract.Contract.GetChannelsAddresses(&_ChannelManagerContract.CallOpts)
}

// GetChannelsAddresses is a free data retrieval call binding the contract method 0x6785b500.
//
// Solidity: function getChannelsAddresses() constant returns(address[])
func (_ChannelManagerContract *ChannelManagerContractCallerSession) GetChannelsAddresses() ([]common.Address, error) {
	return _ChannelManagerContract.Contract.GetChannelsAddresses(&_ChannelManagerContract.CallOpts)
}

// GetChannelsParticipants is a free data retrieval call binding the contract method 0x0b74b620.
//
// Solidity: function getChannelsParticipants() constant returns(address[])
func (_ChannelManagerContract *ChannelManagerContractCaller) GetChannelsParticipants(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _ChannelManagerContract.contract.Call(opts, out, "getChannelsParticipants")
	return *ret0, err
}

// GetChannelsParticipants is a free data retrieval call binding the contract method 0x0b74b620.
//
// Solidity: function getChannelsParticipants() constant returns(address[])
func (_ChannelManagerContract *ChannelManagerContractSession) GetChannelsParticipants() ([]common.Address, error) {
	return _ChannelManagerContract.Contract.GetChannelsParticipants(&_ChannelManagerContract.CallOpts)
}

// GetChannelsParticipants is a free data retrieval call binding the contract method 0x0b74b620.
//
// Solidity: function getChannelsParticipants() constant returns(address[])
func (_ChannelManagerContract *ChannelManagerContractCallerSession) GetChannelsParticipants() ([]common.Address, error) {
	return _ChannelManagerContract.Contract.GetChannelsParticipants(&_ChannelManagerContract.CallOpts)
}

// NettingContractsByAddress is a free data retrieval call binding the contract method 0x6cb30fee.
//
// Solidity: function nettingContractsByAddress(node_address address) constant returns(address[])
func (_ChannelManagerContract *ChannelManagerContractCaller) NettingContractsByAddress(opts *bind.CallOpts, node_address common.Address) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _ChannelManagerContract.contract.Call(opts, out, "nettingContractsByAddress", node_address)
	return *ret0, err
}

// NettingContractsByAddress is a free data retrieval call binding the contract method 0x6cb30fee.
//
// Solidity: function nettingContractsByAddress(node_address address) constant returns(address[])
func (_ChannelManagerContract *ChannelManagerContractSession) NettingContractsByAddress(node_address common.Address) ([]common.Address, error) {
	return _ChannelManagerContract.Contract.NettingContractsByAddress(&_ChannelManagerContract.CallOpts, node_address)
}

// NettingContractsByAddress is a free data retrieval call binding the contract method 0x6cb30fee.
//
// Solidity: function nettingContractsByAddress(node_address address) constant returns(address[])
func (_ChannelManagerContract *ChannelManagerContractCallerSession) NettingContractsByAddress(node_address common.Address) ([]common.Address, error) {
	return _ChannelManagerContract.Contract.NettingContractsByAddress(&_ChannelManagerContract.CallOpts, node_address)
}

// TokenAddress is a free data retrieval call binding the contract method 0x9d76ea58.
//
// Solidity: function tokenAddress() constant returns(address)
func (_ChannelManagerContract *ChannelManagerContractCaller) TokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ChannelManagerContract.contract.Call(opts, out, "tokenAddress")
	return *ret0, err
}

// TokenAddress is a free data retrieval call binding the contract method 0x9d76ea58.
//
// Solidity: function tokenAddress() constant returns(address)
func (_ChannelManagerContract *ChannelManagerContractSession) TokenAddress() (common.Address, error) {
	return _ChannelManagerContract.Contract.TokenAddress(&_ChannelManagerContract.CallOpts)
}

// TokenAddress is a free data retrieval call binding the contract method 0x9d76ea58.
//
// Solidity: function tokenAddress() constant returns(address)
func (_ChannelManagerContract *ChannelManagerContractCallerSession) TokenAddress() (common.Address, error) {
	return _ChannelManagerContract.Contract.TokenAddress(&_ChannelManagerContract.CallOpts)
}

// NewChannel is a paid mutator transaction binding the contract method 0xf26c6aed.
//
// Solidity: function newChannel(partner address, settle_timeout uint256) returns(address)
func (_ChannelManagerContract *ChannelManagerContractTransactor) NewChannel(opts *bind.TransactOpts, partner common.Address, settle_timeout *big.Int) (*types.Transaction, error) {
	return _ChannelManagerContract.contract.Transact(opts, "newChannel", partner, settle_timeout)
}

// NewChannel is a paid mutator transaction binding the contract method 0xf26c6aed.
//
// Solidity: function newChannel(partner address, settle_timeout uint256) returns(address)
func (_ChannelManagerContract *ChannelManagerContractSession) NewChannel(partner common.Address, settle_timeout *big.Int) (*types.Transaction, error) {
	return _ChannelManagerContract.Contract.NewChannel(&_ChannelManagerContract.TransactOpts, partner, settle_timeout)
}

// NewChannel is a paid mutator transaction binding the contract method 0xf26c6aed.
//
// Solidity: function newChannel(partner address, settle_timeout uint256) returns(address)
func (_ChannelManagerContract *ChannelManagerContractTransactorSession) NewChannel(partner common.Address, settle_timeout *big.Int) (*types.Transaction, error) {
	return _ChannelManagerContract.Contract.NewChannel(&_ChannelManagerContract.TransactOpts, partner, settle_timeout)
}

// ChannelManagerContractChannelDeletedIterator is returned from FilterChannelDeleted and is used to iterate over the raw logs and unpacked data for ChannelDeleted events raised by the ChannelManagerContract contract.
type ChannelManagerContractChannelDeletedIterator struct {
	Event *ChannelManagerContractChannelDeleted // Event containing the contract specifics and raw log

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
func (it *ChannelManagerContractChannelDeletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelManagerContractChannelDeleted)
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
		it.Event = new(ChannelManagerContractChannelDeleted)
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
func (it *ChannelManagerContractChannelDeletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChannelManagerContractChannelDeletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChannelManagerContractChannelDeleted represents a ChannelDeleted event raised by the ChannelManagerContract contract.
type ChannelManagerContractChannelDeleted struct {
	Registry_address common.Address
	Caller_address   common.Address
	Partner          common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterChannelDeleted is a free log retrieval operation binding the contract event 0x91cd7e9ad7c88602bc7b06adc62de54cb670665c75894414c01ead5f2baf0930.
//
// Solidity: event ChannelDeleted(registry_address address, caller_address address, partner address)
func (_ChannelManagerContract *ChannelManagerContractFilterer) FilterChannelDeleted(opts *bind.FilterOpts) (*ChannelManagerContractChannelDeletedIterator, error) {

	logs, sub, err := _ChannelManagerContract.contract.FilterLogs(opts, "ChannelDeleted")
	if err != nil {
		return nil, err
	}
	return &ChannelManagerContractChannelDeletedIterator{contract: _ChannelManagerContract.contract, event: "ChannelDeleted", logs: logs, sub: sub}, nil
}

// WatchChannelDeleted is a free log subscription operation binding the contract event 0x91cd7e9ad7c88602bc7b06adc62de54cb670665c75894414c01ead5f2baf0930.
//
// Solidity: event ChannelDeleted(registry_address address, caller_address address, partner address)
func (_ChannelManagerContract *ChannelManagerContractFilterer) WatchChannelDeleted(opts *bind.WatchOpts, sink chan<- *ChannelManagerContractChannelDeleted) (event.Subscription, error) {

	logs, sub, err := _ChannelManagerContract.contract.WatchLogs(opts, "ChannelDeleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChannelManagerContractChannelDeleted)
				if err := _ChannelManagerContract.contract.UnpackLog(event, "ChannelDeleted", log); err != nil {
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

// ChannelManagerContractChannelNewIterator is returned from FilterChannelNew and is used to iterate over the raw logs and unpacked data for ChannelNew events raised by the ChannelManagerContract contract.
type ChannelManagerContractChannelNewIterator struct {
	Event *ChannelManagerContractChannelNew // Event containing the contract specifics and raw log

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
func (it *ChannelManagerContractChannelNewIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelManagerContractChannelNew)
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
		it.Event = new(ChannelManagerContractChannelNew)
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
func (it *ChannelManagerContractChannelNewIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChannelManagerContractChannelNewIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChannelManagerContractChannelNew represents a ChannelNew event raised by the ChannelManagerContract contract.
type ChannelManagerContractChannelNew struct {
	Registry_address common.Address
	Netting_channel  common.Address
	Participant1     common.Address
	Participant2     common.Address
	Settle_timeout   *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterChannelNew is a free log retrieval operation binding the contract event 0xc15f36191ba1aefd153ae0ece89afc6004a66b158bec44593e70bccff2ae7a0f.
//
// Solidity: event ChannelNew(registry_address address, netting_channel address, participant1 address, participant2 address, settle_timeout uint256)
func (_ChannelManagerContract *ChannelManagerContractFilterer) FilterChannelNew(opts *bind.FilterOpts) (*ChannelManagerContractChannelNewIterator, error) {

	logs, sub, err := _ChannelManagerContract.contract.FilterLogs(opts, "ChannelNew")
	if err != nil {
		return nil, err
	}
	return &ChannelManagerContractChannelNewIterator{contract: _ChannelManagerContract.contract, event: "ChannelNew", logs: logs, sub: sub}, nil
}

// WatchChannelNew is a free log subscription operation binding the contract event 0xc15f36191ba1aefd153ae0ece89afc6004a66b158bec44593e70bccff2ae7a0f.
//
// Solidity: event ChannelNew(registry_address address, netting_channel address, participant1 address, participant2 address, settle_timeout uint256)
func (_ChannelManagerContract *ChannelManagerContractFilterer) WatchChannelNew(opts *bind.WatchOpts, sink chan<- *ChannelManagerContractChannelNew) (event.Subscription, error) {

	logs, sub, err := _ChannelManagerContract.contract.WatchLogs(opts, "ChannelNew")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChannelManagerContractChannelNew)
				if err := _ChannelManagerContract.contract.UnpackLog(event, "ChannelNew", log); err != nil {
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

// ChannelManagerLibraryABI is the input ABI used to generate the binding from.
const ChannelManagerLibraryABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"self\",\"type\":\"ChannelManagerLibrary.Data storage\"},{\"name\":\"partner\",\"type\":\"address\"}],\"name\":\"getChannelWith\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"self\",\"type\":\"ChannelManagerLibrary.Data storage\"},{\"name\":\"partner\",\"type\":\"address\"},{\"name\":\"settle_timeout\",\"type\":\"uint256\"}],\"name\":\"newChannel\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// ChannelManagerLibraryBin is the compiled bytecode used for deploying new contracts.
var ChannelManagerLibraryBin = `0x6115ff610030600b82828239805160001a6073146000811461002057610022565bfe5b5030600052607381538281f300730000000000000000000000000000000000000000301460806040526004361061006d5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416638a1c00e28114610072578063941583a5146100a5578063b32c65c8146100cc575b600080fd5b610089600435600160a060020a0360243516610149565b60408051600160a060020a039092168252519081900360200190f35b8180156100b157600080fd5b50610089600435600160a060020a03602435166044356101a8565b6100d4610579565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561010e5781810151838201526020016100f6565b50505050905090810190601f16801561013b5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b600080600061015833856105b0565b6000818152600387016020526040902054909250905080156101a057600285018054600019830190811061018857fe5b600091825260209091200154600160a060020a031692505b505092915050565b33600160a060020a0381811660009081526004860160205260408082209286168252812090928390819081908190819081906101e4908c6105b0565b600081815260038e01602052604090205460018e01548e54929850909650600160a060020a039081169116338d8d61021a61064a565b600160a060020a039586168152938516602085015291841660408085019190915293166060830152608082015290519081900360a001906000f080158015610266573d6000803e3d6000fd5b50935084156103b55760028c018054600019870190811061028357fe5b600091825260209091200154600160a060020a031692506102a383610642565b156102ad57600080fd5b5050600160a060020a03338116600081815260058d0160208181526040808420958f1684529481528483205491815284832093835292909252919091205460028c018054859190600019880190811061030257fe5b9060005260206000200160006101000a815481600160a060020a030219169083600160a060020a0316021790555083886001840381548110151561034257fe5b9060005260206000200160006101000a815481600160a060020a030219169083600160a060020a0316021790555083876001830381548110151561038257fe5b9060005260206000200160006101000a815481600160a060020a030219169083600160a060020a03160217905550610569565b8b6002018490806001815401808255809150509060018203906000526020600020016000909192909190916101000a815481600160a060020a030219169083600160a060020a0316021790555050878490806001815401808255809150509060018203906000526020600020016000909192909190916101000a815481600160a060020a030219169083600160a060020a0316021790555050868490806001815401808255809150509060018203906000526020600020016000909192909190916101000a815481600160a060020a030219169083600160a060020a03160217905550508b600201805490508c600301600088600019166000191681526020019081526020016000208190555087805490508c600501600033600160a060020a0316600160a060020a0316815260200190815260200160002060008d600160a060020a0316600160a060020a031681526020019081526020016000208190555086805490508c60050160008d600160a060020a0316600160a060020a03168152602001908152602001600020600033600160a060020a0316600160a060020a03168152602001908152602001600020819055505b50919a9950505050505050505050565b60408051808201909152600581527f302e322e5f000000000000000000000000000000000000000000000000000000602082015281565b600081600160a060020a031683600160a060020a031610156106065750604080516c01000000000000000000000000600160a060020a03808616820283528416026014820152905190819003602801902061063c565b50604080516c01000000000000000000000000600160a060020a0380851682028352851602601482015290519081900360280190205b92915050565b6000903b1190565b604051610f798061065b833901905600608060405234801561001057600080fd5b5060405160a080610f79833981016040908152815160208301519183015160608401516080909401519193909180600681108015906100525750622932e08111155b151561005d57600080fd5b600160a060020a03848116908416141561007657600080fd5b5060068054600160a060020a03948516600160a060020a03199182168117909255600c805482169486169485179055600091825260126020526040808320805460ff199081166001908117909255958452908320805460029616959095179094556004805482169786169790971790965560058054909616949093169390931790935590815543909155610e6990819061011090396000f3006080604052600436106100cf5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166311da60b481146100e157806327d120fe146100f857806353af5d101461016d57806359023f891461019e578063597e1fb5146102515780635e1fc56e146102785780635f88eade146102ed5780636d2381b31461030257806373d4a13a1461034c5780637ebdc478146103a7578063965f12e8146103bc5780639d76ea5814610463578063b32c65c814610478578063b6b55f2514610502575b3480156100db57600080fd5b50600080fd5b3480156100ed57600080fd5b506100f661052e565b005b34801561010457600080fd5b50604080516020601f6084356004818101359283018490048402850184019095528184526100f69467ffffffffffffffff813516946024803595604435956064359536959460a49490939101919081908401838280828437509497506105f19650505050505050565b34801561017957600080fd5b50610182610753565b60408051600160a060020a039092168252519081900360200190f35b3480156101aa57600080fd5b50604080516020601f6084356004818101359283018490048402850184019095528184526100f69467ffffffffffffffff813516946024803595604435956064359536959460a494909391019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a9998810197919650918201945092508291508401838280828437509497506107629650505050505050565b34801561025d57600080fd5b50610266610946565b60408051918252519081900360200190f35b34801561028457600080fd5b50604080516020601f6084356004818101359283018490048402850184019095528184526100f69467ffffffffffffffff813516946024803595604435956064359536959460a494909391019190819084018382808284375094975061094c9650505050505050565b3480156102f957600080fd5b50610266610ab9565b34801561030e57600080fd5b50610317610abf565b60408051600160a060020a03958616815260208101949094529190931682820152606082019290925290519081900360800190f35b34801561035857600080fd5b50610361610adf565b60408051978852602088019690965286860194909452600160a060020a03928316606087015290821660808601521660a0840152151560c0830152519081900360e00190f35b3480156103b357600080fd5b50610266610b0e565b3480156103c857600080fd5b5060408051602060046024803582810135601f81018590048502860185019096528585526100f6958335600160a060020a031695369560449491939091019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a9998810197919650918201945092508291508401838280828437509497505093359450610b149350505050565b34801561046f57600080fd5b50610182610cd9565b34801561048457600080fd5b5061048d610ce8565b6040805160208082528351818301528351919283929083019185019080838360005b838110156104c75781810151838201526020016104af565b50505050905090810190601f1680156104f45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561050e57600080fd5b5061051a600435610d1f565b604080519115158252519081900360200190f35b604080517fde394e0d000000000000000000000000000000000000000000000000000000008152600060048201819052915173__NettingChannelLibrary.sol:NettingCha__9263de394e0d9260248082019391829003018186803b15801561059757600080fd5b505af41580156105ab573d6000803e3d6000fd5b505060045460408051600160a060020a039092168252517f40bab4574d5dad3320b539548362899f170b309febdeaadb2e4d0367311df4e09350908190036020019150a1565b6040517ff565eb3600000000000000000000000000000000000000000000000000000000815260006004820181815267ffffffffffffffff8816602484015260448301879052606483018690526084830185905260c060a48401908152845160c4850152845173__NettingChannelLibrary.sol:NettingCha__9463f565eb3694938b938b938b938b938b939092909160e4019060208501908083838d5b838110156106a8578181015183820152602001610690565b50505050905090810190601f1680156106d55780820380516001836020036101000a031916815260200191505b5097505050505050505060006040518083038186803b1580156106f757600080fd5b505af415801561070b573d6000803e3d6000fd5b505060408051600160a060020a033316815290517fa0379b1bd0864245b4ff39bf6a023065e80d0e9276d2671d94b9f653b4bbcdfe9350908190036020019150a15050505050565b600354600160a060020a031690565b6040517ff002345600000000000000000000000000000000000000000000000000000000815260006004820181815267ffffffffffffffff8916602484015260448301889052606483018790526084830186905260e060a48401908152855160e48501528551929373__NettingChannelLibrary.sol:NettingCha__9363f00234569386938d938d938d938d938d938d9360c4820191610104019060208701908083838f5b83811015610820578181015183820152602001610808565b50505050905090810190601f16801561084d5780820380516001836020036101000a031916815260200191505b50838103825284518152845160209182019186019080838360005b83811015610880578181015183820152602001610868565b50505050905090810190601f1680156108ad5780820380516001836020036101000a031916815260200191505b50995050505050505050505060206040518083038186803b1580156108d157600080fd5b505af41580156108e5573d6000803e3d6000fd5b505050506040513d60208110156108fb57600080fd5b505160408051600160a060020a038316815290519192507fa0379b1bd0864245b4ff39bf6a023065e80d0e9276d2671d94b9f653b4bbcdfe919081900360200190a150505050505050565b60025490565b6040517fc800b00200000000000000000000000000000000000000000000000000000000815260006004820181815267ffffffffffffffff8816602484015260448301879052606483018690526084830185905260c060a48401908152845160c4850152845173__NettingChannelLibrary.sol:NettingCha__9463c800b00294938b938b938b938b938b939092909160e4019060208501908083838d5b83811015610a035781810151838201526020016109eb565b50505050905090810190601f168015610a305780820380516001836020036101000a031916815260200191505b5097505050505050505060006040518083038186803b158015610a5257600080fd5b505af4158015610a66573d6000803e3d6000fd5b505060045460408051600160a060020a03928316815233909216602083015280517f48b749184840ce6faa9f6691fd4af8e7c969cd25ee881675131e2c9358ec3118945091829003019150a15050505050565b60015490565b600654600754600c54600d54600160a060020a0393841694929390911691565b600054600154600254600354600454600554601354600160a060020a0393841693928316929091169060ff1687565b60005490565b6040517ffee4658a000000000000000000000000000000000000000000000000000000008152600060048201818152600160a060020a03871660248401526084830184905260a060448401908152865160a4850152865173__NettingChannelLibrary.sol:NettingCha__9463fee4658a94938a938a938a938a939291606482019160c4019060208801908083838d5b83811015610bbd578181015183820152602001610ba5565b50505050905090810190601f168015610bea5780820380516001836020036101000a031916815260200191505b50838103825285518152855160209182019187019080838360005b83811015610c1d578181015183820152602001610c05565b50505050905090810190601f168015610c4a5780820380516001836020036101000a031916815260200191505b5097505050505050505060006040518083038186803b158015610c6c57600080fd5b505af4158015610c80573d6000803e3d6000fd5b505060045460408051600160a060020a039283168152602081018690523390921682820152517f15b5384ec609c8a364f4f4ba5969bf39b4e2b5c6a8364773267c682afe40586f9350908190036060019150a150505050565b600554600160a060020a031690565b60408051808201909152600581527f302e322e5f000000000000000000000000000000000000000000000000000000602082015281565b60008060008073__NettingChannelLibrary.sol:NettingCha__633268a05a9091866040518363ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004018083815260200182815260200192505050604080518083038186803b158015610d9857600080fd5b505af4158015610dac573d6000803e3d6000fd5b505050506040513d6040811015610dc257600080fd5b508051602090910151909250905060018215151415610e365760045460055460408051600160a060020a039384168152918316602083015233909216818301526060810183905290517f61cf19a8c55ff5e9fda69ea0207b0350002d5073a364313b554d352ff6d7803d9181900360800190a15b50929150505600a165627a7a72305820f8fbd115c8d032d3a55554d62f83f4ad6fafdd90819d577ccd86806f06c3df820029a165627a7a7230582034c338d65f3c334081331f58e8a58e43fcc7de934e88327fb61c2f48e5f247c80029`

// DeployChannelManagerLibrary deploys a new Ethereum contract, binding an instance of ChannelManagerLibrary to it.
func DeployChannelManagerLibrary(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ChannelManagerLibrary, error) {
	parsed, err := abi.JSON(strings.NewReader(ChannelManagerLibraryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ChannelManagerLibraryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ChannelManagerLibrary{ChannelManagerLibraryCaller: ChannelManagerLibraryCaller{contract: contract}, ChannelManagerLibraryTransactor: ChannelManagerLibraryTransactor{contract: contract}, ChannelManagerLibraryFilterer: ChannelManagerLibraryFilterer{contract: contract}}, nil
}

// ChannelManagerLibrary is an auto generated Go binding around an Ethereum contract.
type ChannelManagerLibrary struct {
	ChannelManagerLibraryCaller     // Read-only binding to the contract
	ChannelManagerLibraryTransactor // Write-only binding to the contract
	ChannelManagerLibraryFilterer   // Log filterer for contract events
}

// ChannelManagerLibraryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChannelManagerLibraryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChannelManagerLibraryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChannelManagerLibraryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChannelManagerLibraryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChannelManagerLibraryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChannelManagerLibrarySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChannelManagerLibrarySession struct {
	Contract     *ChannelManagerLibrary // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ChannelManagerLibraryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChannelManagerLibraryCallerSession struct {
	Contract *ChannelManagerLibraryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// ChannelManagerLibraryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChannelManagerLibraryTransactorSession struct {
	Contract     *ChannelManagerLibraryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// ChannelManagerLibraryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ChannelManagerLibraryRaw struct {
	Contract *ChannelManagerLibrary // Generic contract binding to access the raw methods on
}

// ChannelManagerLibraryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChannelManagerLibraryCallerRaw struct {
	Contract *ChannelManagerLibraryCaller // Generic read-only contract binding to access the raw methods on
}

// ChannelManagerLibraryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChannelManagerLibraryTransactorRaw struct {
	Contract *ChannelManagerLibraryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChannelManagerLibrary creates a new instance of ChannelManagerLibrary, bound to a specific deployed contract.
func NewChannelManagerLibrary(address common.Address, backend bind.ContractBackend) (*ChannelManagerLibrary, error) {
	contract, err := bindChannelManagerLibrary(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChannelManagerLibrary{ChannelManagerLibraryCaller: ChannelManagerLibraryCaller{contract: contract}, ChannelManagerLibraryTransactor: ChannelManagerLibraryTransactor{contract: contract}, ChannelManagerLibraryFilterer: ChannelManagerLibraryFilterer{contract: contract}}, nil
}

// NewChannelManagerLibraryCaller creates a new read-only instance of ChannelManagerLibrary, bound to a specific deployed contract.
func NewChannelManagerLibraryCaller(address common.Address, caller bind.ContractCaller) (*ChannelManagerLibraryCaller, error) {
	contract, err := bindChannelManagerLibrary(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChannelManagerLibraryCaller{contract: contract}, nil
}

// NewChannelManagerLibraryTransactor creates a new write-only instance of ChannelManagerLibrary, bound to a specific deployed contract.
func NewChannelManagerLibraryTransactor(address common.Address, transactor bind.ContractTransactor) (*ChannelManagerLibraryTransactor, error) {
	contract, err := bindChannelManagerLibrary(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChannelManagerLibraryTransactor{contract: contract}, nil
}

// NewChannelManagerLibraryFilterer creates a new log filterer instance of ChannelManagerLibrary, bound to a specific deployed contract.
func NewChannelManagerLibraryFilterer(address common.Address, filterer bind.ContractFilterer) (*ChannelManagerLibraryFilterer, error) {
	contract, err := bindChannelManagerLibrary(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChannelManagerLibraryFilterer{contract: contract}, nil
}

// bindChannelManagerLibrary binds a generic wrapper to an already deployed contract.
func bindChannelManagerLibrary(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ChannelManagerLibraryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChannelManagerLibrary *ChannelManagerLibraryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ChannelManagerLibrary.Contract.ChannelManagerLibraryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChannelManagerLibrary *ChannelManagerLibraryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChannelManagerLibrary.Contract.ChannelManagerLibraryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChannelManagerLibrary *ChannelManagerLibraryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChannelManagerLibrary.Contract.ChannelManagerLibraryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChannelManagerLibrary *ChannelManagerLibraryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ChannelManagerLibrary.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChannelManagerLibrary *ChannelManagerLibraryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChannelManagerLibrary.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChannelManagerLibrary *ChannelManagerLibraryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChannelManagerLibrary.Contract.contract.Transact(opts, method, params...)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_ChannelManagerLibrary *ChannelManagerLibraryCaller) Contract_version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ChannelManagerLibrary.contract.Call(opts, out, "contract_version")
	return *ret0, err
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_ChannelManagerLibrary *ChannelManagerLibrarySession) Contract_version() (string, error) {
	return _ChannelManagerLibrary.Contract.Contract_version(&_ChannelManagerLibrary.CallOpts)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_ChannelManagerLibrary *ChannelManagerLibraryCallerSession) Contract_version() (string, error) {
	return _ChannelManagerLibrary.Contract.Contract_version(&_ChannelManagerLibrary.CallOpts)
}

// NettingChannelContractABI is the input ABI used to generate the binding from.
const NettingChannelContractABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"settle\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"transferred_amount\",\"type\":\"uint256\"},{\"name\":\"locksroot\",\"type\":\"bytes32\"},{\"name\":\"extra_hash\",\"type\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"updateTransfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"closingAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"transferred_amount\",\"type\":\"uint256\"},{\"name\":\"locksroot\",\"type\":\"bytes32\"},{\"name\":\"extra_hash\",\"type\":\"bytes32\"},{\"name\":\"closing_signature\",\"type\":\"bytes\"},{\"name\":\"non_closing_signature\",\"type\":\"bytes\"}],\"name\":\"updateTransferDelegate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"closed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"transferred_amount\",\"type\":\"uint256\"},{\"name\":\"locksroot\",\"type\":\"bytes32\"},{\"name\":\"extra_hash\",\"type\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"close\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"opened\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"addressAndBalance\",\"outputs\":[{\"name\":\"participant1\",\"type\":\"address\"},{\"name\":\"balance1\",\"type\":\"uint256\"},{\"name\":\"participant2\",\"type\":\"address\"},{\"name\":\"balance2\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"data\",\"outputs\":[{\"name\":\"settle_timeout\",\"type\":\"uint256\"},{\"name\":\"opened\",\"type\":\"uint256\"},{\"name\":\"closed\",\"type\":\"uint256\"},{\"name\":\"closing_address\",\"type\":\"address\"},{\"name\":\"registry_address\",\"type\":\"address\"},{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"updated\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"settleTimeout\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"participant\",\"type\":\"address\"},{\"name\":\"locked_encoded\",\"type\":\"bytes\"},{\"name\":\"merkle_proof\",\"type\":\"bytes\"},{\"name\":\"secret\",\"type\":\"bytes32\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"registry_address\",\"type\":\"address\"},{\"name\":\"token_address\",\"type\":\"address\"},{\"name\":\"participant1\",\"type\":\"address\"},{\"name\":\"participant2\",\"type\":\"address\"},{\"name\":\"timeout\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"registry_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"token_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"participant\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"balance\",\"type\":\"uint256\"}],\"name\":\"ChannelNewBalance\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"registry_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"closing_address\",\"type\":\"address\"}],\"name\":\"ChannelClosed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"node_address\",\"type\":\"address\"}],\"name\":\"TransferUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"registry_address\",\"type\":\"address\"}],\"name\":\"ChannelSettled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"registry_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"secret\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"receiver_address\",\"type\":\"address\"}],\"name\":\"ChannelSecretRevealed\",\"type\":\"event\"}]"

// NettingChannelContractBin is the compiled bytecode used for deploying new contracts.
const NettingChannelContractBin = `0x608060405234801561001057600080fd5b5060405160a080610f79833981016040908152815160208301519183015160608401516080909401519193909180600681108015906100525750622932e08111155b151561005d57600080fd5b600160a060020a03848116908416141561007657600080fd5b5060068054600160a060020a03948516600160a060020a03199182168117909255600c805482169486169485179055600091825260126020526040808320805460ff199081166001908117909255958452908320805460029616959095179094556004805482169786169790971790965560058054909616949093169390931790935590815543909155610e6990819061011090396000f3006080604052600436106100cf5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166311da60b481146100e157806327d120fe146100f857806353af5d101461016d57806359023f891461019e578063597e1fb5146102515780635e1fc56e146102785780635f88eade146102ed5780636d2381b31461030257806373d4a13a1461034c5780637ebdc478146103a7578063965f12e8146103bc5780639d76ea5814610463578063b32c65c814610478578063b6b55f2514610502575b3480156100db57600080fd5b50600080fd5b3480156100ed57600080fd5b506100f661052e565b005b34801561010457600080fd5b50604080516020601f6084356004818101359283018490048402850184019095528184526100f69467ffffffffffffffff813516946024803595604435956064359536959460a49490939101919081908401838280828437509497506105f19650505050505050565b34801561017957600080fd5b50610182610753565b60408051600160a060020a039092168252519081900360200190f35b3480156101aa57600080fd5b50604080516020601f6084356004818101359283018490048402850184019095528184526100f69467ffffffffffffffff813516946024803595604435956064359536959460a494909391019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a9998810197919650918201945092508291508401838280828437509497506107629650505050505050565b34801561025d57600080fd5b50610266610946565b60408051918252519081900360200190f35b34801561028457600080fd5b50604080516020601f6084356004818101359283018490048402850184019095528184526100f69467ffffffffffffffff813516946024803595604435956064359536959460a494909391019190819084018382808284375094975061094c9650505050505050565b3480156102f957600080fd5b50610266610ab9565b34801561030e57600080fd5b50610317610abf565b60408051600160a060020a03958616815260208101949094529190931682820152606082019290925290519081900360800190f35b34801561035857600080fd5b50610361610adf565b60408051978852602088019690965286860194909452600160a060020a03928316606087015290821660808601521660a0840152151560c0830152519081900360e00190f35b3480156103b357600080fd5b50610266610b0e565b3480156103c857600080fd5b5060408051602060046024803582810135601f81018590048502860185019096528585526100f6958335600160a060020a031695369560449491939091019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a9998810197919650918201945092508291508401838280828437509497505093359450610b149350505050565b34801561046f57600080fd5b50610182610cd9565b34801561048457600080fd5b5061048d610ce8565b6040805160208082528351818301528351919283929083019185019080838360005b838110156104c75781810151838201526020016104af565b50505050905090810190601f1680156104f45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561050e57600080fd5b5061051a600435610d1f565b604080519115158252519081900360200190f35b604080517fde394e0d000000000000000000000000000000000000000000000000000000008152600060048201819052915173__NettingChannelLibrary.sol:NettingCha__9263de394e0d9260248082019391829003018186803b15801561059757600080fd5b505af41580156105ab573d6000803e3d6000fd5b505060045460408051600160a060020a039092168252517f40bab4574d5dad3320b539548362899f170b309febdeaadb2e4d0367311df4e09350908190036020019150a1565b6040517ff565eb3600000000000000000000000000000000000000000000000000000000815260006004820181815267ffffffffffffffff8816602484015260448301879052606483018690526084830185905260c060a48401908152845160c4850152845173__NettingChannelLibrary.sol:NettingCha__9463f565eb3694938b938b938b938b938b939092909160e4019060208501908083838d5b838110156106a8578181015183820152602001610690565b50505050905090810190601f1680156106d55780820380516001836020036101000a031916815260200191505b5097505050505050505060006040518083038186803b1580156106f757600080fd5b505af415801561070b573d6000803e3d6000fd5b505060408051600160a060020a033316815290517fa0379b1bd0864245b4ff39bf6a023065e80d0e9276d2671d94b9f653b4bbcdfe9350908190036020019150a15050505050565b600354600160a060020a031690565b6040517ff002345600000000000000000000000000000000000000000000000000000000815260006004820181815267ffffffffffffffff8916602484015260448301889052606483018790526084830186905260e060a48401908152855160e48501528551929373__NettingChannelLibrary.sol:NettingCha__9363f00234569386938d938d938d938d938d938d9360c4820191610104019060208701908083838f5b83811015610820578181015183820152602001610808565b50505050905090810190601f16801561084d5780820380516001836020036101000a031916815260200191505b50838103825284518152845160209182019186019080838360005b83811015610880578181015183820152602001610868565b50505050905090810190601f1680156108ad5780820380516001836020036101000a031916815260200191505b50995050505050505050505060206040518083038186803b1580156108d157600080fd5b505af41580156108e5573d6000803e3d6000fd5b505050506040513d60208110156108fb57600080fd5b505160408051600160a060020a038316815290519192507fa0379b1bd0864245b4ff39bf6a023065e80d0e9276d2671d94b9f653b4bbcdfe919081900360200190a150505050505050565b60025490565b6040517fc800b00200000000000000000000000000000000000000000000000000000000815260006004820181815267ffffffffffffffff8816602484015260448301879052606483018690526084830185905260c060a48401908152845160c4850152845173__NettingChannelLibrary.sol:NettingCha__9463c800b00294938b938b938b938b938b939092909160e4019060208501908083838d5b83811015610a035781810151838201526020016109eb565b50505050905090810190601f168015610a305780820380516001836020036101000a031916815260200191505b5097505050505050505060006040518083038186803b158015610a5257600080fd5b505af4158015610a66573d6000803e3d6000fd5b505060045460408051600160a060020a03928316815233909216602083015280517f48b749184840ce6faa9f6691fd4af8e7c969cd25ee881675131e2c9358ec3118945091829003019150a15050505050565b60015490565b600654600754600c54600d54600160a060020a0393841694929390911691565b600054600154600254600354600454600554601354600160a060020a0393841693928316929091169060ff1687565b60005490565b6040517ffee4658a000000000000000000000000000000000000000000000000000000008152600060048201818152600160a060020a03871660248401526084830184905260a060448401908152865160a4850152865173__NettingChannelLibrary.sol:NettingCha__9463fee4658a94938a938a938a938a939291606482019160c4019060208801908083838d5b83811015610bbd578181015183820152602001610ba5565b50505050905090810190601f168015610bea5780820380516001836020036101000a031916815260200191505b50838103825285518152855160209182019187019080838360005b83811015610c1d578181015183820152602001610c05565b50505050905090810190601f168015610c4a5780820380516001836020036101000a031916815260200191505b5097505050505050505060006040518083038186803b158015610c6c57600080fd5b505af4158015610c80573d6000803e3d6000fd5b505060045460408051600160a060020a039283168152602081018690523390921682820152517f15b5384ec609c8a364f4f4ba5969bf39b4e2b5c6a8364773267c682afe40586f9350908190036060019150a150505050565b600554600160a060020a031690565b60408051808201909152600581527f302e322e5f000000000000000000000000000000000000000000000000000000602082015281565b60008060008073__NettingChannelLibrary.sol:NettingCha__633268a05a9091866040518363ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004018083815260200182815260200192505050604080518083038186803b158015610d9857600080fd5b505af4158015610dac573d6000803e3d6000fd5b505050506040513d6040811015610dc257600080fd5b508051602090910151909250905060018215151415610e365760045460055460408051600160a060020a039384168152918316602083015233909216818301526060810183905290517f61cf19a8c55ff5e9fda69ea0207b0350002d5073a364313b554d352ff6d7803d9181900360800190a15b50929150505600a165627a7a72305820f8fbd115c8d032d3a55554d62f83f4ad6fafdd90819d577ccd86806f06c3df820029`

// DeployNettingChannelContract deploys a new Ethereum contract, binding an instance of NettingChannelContract to it.
func DeployNettingChannelContract(auth *bind.TransactOpts, backend bind.ContractBackend, registry_address common.Address, token_address common.Address, participant1 common.Address, participant2 common.Address, timeout *big.Int) (common.Address, *types.Transaction, *NettingChannelContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NettingChannelContractABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(NettingChannelContractBin), backend, registry_address, token_address, participant1, participant2, timeout)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NettingChannelContract{NettingChannelContractCaller: NettingChannelContractCaller{contract: contract}, NettingChannelContractTransactor: NettingChannelContractTransactor{contract: contract}, NettingChannelContractFilterer: NettingChannelContractFilterer{contract: contract}}, nil
}

// NettingChannelContract is an auto generated Go binding around an Ethereum contract.
type NettingChannelContract struct {
	NettingChannelContractCaller     // Read-only binding to the contract
	NettingChannelContractTransactor // Write-only binding to the contract
	NettingChannelContractFilterer   // Log filterer for contract events
}

// NettingChannelContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type NettingChannelContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NettingChannelContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NettingChannelContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NettingChannelContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NettingChannelContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NettingChannelContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NettingChannelContractSession struct {
	Contract     *NettingChannelContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// NettingChannelContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NettingChannelContractCallerSession struct {
	Contract *NettingChannelContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// NettingChannelContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NettingChannelContractTransactorSession struct {
	Contract     *NettingChannelContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// NettingChannelContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type NettingChannelContractRaw struct {
	Contract *NettingChannelContract // Generic contract binding to access the raw methods on
}

// NettingChannelContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NettingChannelContractCallerRaw struct {
	Contract *NettingChannelContractCaller // Generic read-only contract binding to access the raw methods on
}

// NettingChannelContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NettingChannelContractTransactorRaw struct {
	Contract *NettingChannelContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNettingChannelContract creates a new instance of NettingChannelContract, bound to a specific deployed contract.
func NewNettingChannelContract(address common.Address, backend bind.ContractBackend) (*NettingChannelContract, error) {
	contract, err := bindNettingChannelContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NettingChannelContract{NettingChannelContractCaller: NettingChannelContractCaller{contract: contract}, NettingChannelContractTransactor: NettingChannelContractTransactor{contract: contract}, NettingChannelContractFilterer: NettingChannelContractFilterer{contract: contract}}, nil
}

// NewNettingChannelContractCaller creates a new read-only instance of NettingChannelContract, bound to a specific deployed contract.
func NewNettingChannelContractCaller(address common.Address, caller bind.ContractCaller) (*NettingChannelContractCaller, error) {
	contract, err := bindNettingChannelContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NettingChannelContractCaller{contract: contract}, nil
}

// NewNettingChannelContractTransactor creates a new write-only instance of NettingChannelContract, bound to a specific deployed contract.
func NewNettingChannelContractTransactor(address common.Address, transactor bind.ContractTransactor) (*NettingChannelContractTransactor, error) {
	contract, err := bindNettingChannelContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NettingChannelContractTransactor{contract: contract}, nil
}

// NewNettingChannelContractFilterer creates a new log filterer instance of NettingChannelContract, bound to a specific deployed contract.
func NewNettingChannelContractFilterer(address common.Address, filterer bind.ContractFilterer) (*NettingChannelContractFilterer, error) {
	contract, err := bindNettingChannelContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NettingChannelContractFilterer{contract: contract}, nil
}

// bindNettingChannelContract binds a generic wrapper to an already deployed contract.
func bindNettingChannelContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NettingChannelContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NettingChannelContract *NettingChannelContractRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _NettingChannelContract.Contract.NettingChannelContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NettingChannelContract *NettingChannelContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NettingChannelContract.Contract.NettingChannelContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NettingChannelContract *NettingChannelContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NettingChannelContract.Contract.NettingChannelContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NettingChannelContract *NettingChannelContractCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _NettingChannelContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NettingChannelContract *NettingChannelContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NettingChannelContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NettingChannelContract *NettingChannelContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NettingChannelContract.Contract.contract.Transact(opts, method, params...)
}

// AddressAndBalance is a free data retrieval call binding the contract method 0x6d2381b3.
//
// Solidity: function addressAndBalance() constant returns(participant1 address, balance1 uint256, participant2 address, balance2 uint256)
func (_NettingChannelContract *NettingChannelContractCaller) AddressAndBalance(opts *bind.CallOpts) (struct {
	Participant1 common.Address
	Balance1     *big.Int
	Participant2 common.Address
	Balance2     *big.Int
}, error) {
	ret := new(struct {
		Participant1 common.Address
		Balance1     *big.Int
		Participant2 common.Address
		Balance2     *big.Int
	})
	out := ret
	err := _NettingChannelContract.contract.Call(opts, out, "addressAndBalance")
	return *ret, err
}

// AddressAndBalance is a free data retrieval call binding the contract method 0x6d2381b3.
//
// Solidity: function addressAndBalance() constant returns(participant1 address, balance1 uint256, participant2 address, balance2 uint256)
func (_NettingChannelContract *NettingChannelContractSession) AddressAndBalance() (struct {
	Participant1 common.Address
	Balance1     *big.Int
	Participant2 common.Address
	Balance2     *big.Int
}, error) {
	return _NettingChannelContract.Contract.AddressAndBalance(&_NettingChannelContract.CallOpts)
}

// AddressAndBalance is a free data retrieval call binding the contract method 0x6d2381b3.
//
// Solidity: function addressAndBalance() constant returns(participant1 address, balance1 uint256, participant2 address, balance2 uint256)
func (_NettingChannelContract *NettingChannelContractCallerSession) AddressAndBalance() (struct {
	Participant1 common.Address
	Balance1     *big.Int
	Participant2 common.Address
	Balance2     *big.Int
}, error) {
	return _NettingChannelContract.Contract.AddressAndBalance(&_NettingChannelContract.CallOpts)
}

// Closed is a free data retrieval call binding the contract method 0x597e1fb5.
//
// Solidity: function closed() constant returns(uint256)
func (_NettingChannelContract *NettingChannelContractCaller) Closed(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _NettingChannelContract.contract.Call(opts, out, "closed")
	return *ret0, err
}

// Closed is a free data retrieval call binding the contract method 0x597e1fb5.
//
// Solidity: function closed() constant returns(uint256)
func (_NettingChannelContract *NettingChannelContractSession) Closed() (*big.Int, error) {
	return _NettingChannelContract.Contract.Closed(&_NettingChannelContract.CallOpts)
}

// Closed is a free data retrieval call binding the contract method 0x597e1fb5.
//
// Solidity: function closed() constant returns(uint256)
func (_NettingChannelContract *NettingChannelContractCallerSession) Closed() (*big.Int, error) {
	return _NettingChannelContract.Contract.Closed(&_NettingChannelContract.CallOpts)
}

// ClosingAddress is a free data retrieval call binding the contract method 0x53af5d10.
//
// Solidity: function closingAddress() constant returns(address)
func (_NettingChannelContract *NettingChannelContractCaller) ClosingAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _NettingChannelContract.contract.Call(opts, out, "closingAddress")
	return *ret0, err
}

// ClosingAddress is a free data retrieval call binding the contract method 0x53af5d10.
//
// Solidity: function closingAddress() constant returns(address)
func (_NettingChannelContract *NettingChannelContractSession) ClosingAddress() (common.Address, error) {
	return _NettingChannelContract.Contract.ClosingAddress(&_NettingChannelContract.CallOpts)
}

// ClosingAddress is a free data retrieval call binding the contract method 0x53af5d10.
//
// Solidity: function closingAddress() constant returns(address)
func (_NettingChannelContract *NettingChannelContractCallerSession) ClosingAddress() (common.Address, error) {
	return _NettingChannelContract.Contract.ClosingAddress(&_NettingChannelContract.CallOpts)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_NettingChannelContract *NettingChannelContractCaller) Contract_version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _NettingChannelContract.contract.Call(opts, out, "contract_version")
	return *ret0, err
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_NettingChannelContract *NettingChannelContractSession) Contract_version() (string, error) {
	return _NettingChannelContract.Contract.Contract_version(&_NettingChannelContract.CallOpts)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_NettingChannelContract *NettingChannelContractCallerSession) Contract_version() (string, error) {
	return _NettingChannelContract.Contract.Contract_version(&_NettingChannelContract.CallOpts)
}

// Data is a free data retrieval call binding the contract method 0x73d4a13a.
//
// Solidity: function data() constant returns(settle_timeout uint256, opened uint256, closed uint256, closing_address address, registry_address address, token address, updated bool)
func (_NettingChannelContract *NettingChannelContractCaller) Data(opts *bind.CallOpts) (struct {
	Settle_timeout   *big.Int
	Opened           *big.Int
	Closed           *big.Int
	Closing_address  common.Address
	Registry_address common.Address
	Token            common.Address
	Updated          bool
}, error) {
	ret := new(struct {
		Settle_timeout   *big.Int
		Opened           *big.Int
		Closed           *big.Int
		Closing_address  common.Address
		Registry_address common.Address
		Token            common.Address
		Updated          bool
	})
	out := ret
	err := _NettingChannelContract.contract.Call(opts, out, "data")
	return *ret, err
}

// Data is a free data retrieval call binding the contract method 0x73d4a13a.
//
// Solidity: function data() constant returns(settle_timeout uint256, opened uint256, closed uint256, closing_address address, registry_address address, token address, updated bool)
func (_NettingChannelContract *NettingChannelContractSession) Data() (struct {
	Settle_timeout   *big.Int
	Opened           *big.Int
	Closed           *big.Int
	Closing_address  common.Address
	Registry_address common.Address
	Token            common.Address
	Updated          bool
}, error) {
	return _NettingChannelContract.Contract.Data(&_NettingChannelContract.CallOpts)
}

// Data is a free data retrieval call binding the contract method 0x73d4a13a.
//
// Solidity: function data() constant returns(settle_timeout uint256, opened uint256, closed uint256, closing_address address, registry_address address, token address, updated bool)
func (_NettingChannelContract *NettingChannelContractCallerSession) Data() (struct {
	Settle_timeout   *big.Int
	Opened           *big.Int
	Closed           *big.Int
	Closing_address  common.Address
	Registry_address common.Address
	Token            common.Address
	Updated          bool
}, error) {
	return _NettingChannelContract.Contract.Data(&_NettingChannelContract.CallOpts)
}

// Opened is a free data retrieval call binding the contract method 0x5f88eade.
//
// Solidity: function opened() constant returns(uint256)
func (_NettingChannelContract *NettingChannelContractCaller) Opened(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _NettingChannelContract.contract.Call(opts, out, "opened")
	return *ret0, err
}

// Opened is a free data retrieval call binding the contract method 0x5f88eade.
//
// Solidity: function opened() constant returns(uint256)
func (_NettingChannelContract *NettingChannelContractSession) Opened() (*big.Int, error) {
	return _NettingChannelContract.Contract.Opened(&_NettingChannelContract.CallOpts)
}

// Opened is a free data retrieval call binding the contract method 0x5f88eade.
//
// Solidity: function opened() constant returns(uint256)
func (_NettingChannelContract *NettingChannelContractCallerSession) Opened() (*big.Int, error) {
	return _NettingChannelContract.Contract.Opened(&_NettingChannelContract.CallOpts)
}

// SettleTimeout is a free data retrieval call binding the contract method 0x7ebdc478.
//
// Solidity: function settleTimeout() constant returns(uint256)
func (_NettingChannelContract *NettingChannelContractCaller) SettleTimeout(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _NettingChannelContract.contract.Call(opts, out, "settleTimeout")
	return *ret0, err
}

// SettleTimeout is a free data retrieval call binding the contract method 0x7ebdc478.
//
// Solidity: function settleTimeout() constant returns(uint256)
func (_NettingChannelContract *NettingChannelContractSession) SettleTimeout() (*big.Int, error) {
	return _NettingChannelContract.Contract.SettleTimeout(&_NettingChannelContract.CallOpts)
}

// SettleTimeout is a free data retrieval call binding the contract method 0x7ebdc478.
//
// Solidity: function settleTimeout() constant returns(uint256)
func (_NettingChannelContract *NettingChannelContractCallerSession) SettleTimeout() (*big.Int, error) {
	return _NettingChannelContract.Contract.SettleTimeout(&_NettingChannelContract.CallOpts)
}

// TokenAddress is a free data retrieval call binding the contract method 0x9d76ea58.
//
// Solidity: function tokenAddress() constant returns(address)
func (_NettingChannelContract *NettingChannelContractCaller) TokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _NettingChannelContract.contract.Call(opts, out, "tokenAddress")
	return *ret0, err
}

// TokenAddress is a free data retrieval call binding the contract method 0x9d76ea58.
//
// Solidity: function tokenAddress() constant returns(address)
func (_NettingChannelContract *NettingChannelContractSession) TokenAddress() (common.Address, error) {
	return _NettingChannelContract.Contract.TokenAddress(&_NettingChannelContract.CallOpts)
}

// TokenAddress is a free data retrieval call binding the contract method 0x9d76ea58.
//
// Solidity: function tokenAddress() constant returns(address)
func (_NettingChannelContract *NettingChannelContractCallerSession) TokenAddress() (common.Address, error) {
	return _NettingChannelContract.Contract.TokenAddress(&_NettingChannelContract.CallOpts)
}

// Close is a paid mutator transaction binding the contract method 0x5e1fc56e.
//
// Solidity: function close(nonce uint64, transferred_amount uint256, locksroot bytes32, extra_hash bytes32, signature bytes) returns()
func (_NettingChannelContract *NettingChannelContractTransactor) Close(opts *bind.TransactOpts, nonce uint64, transferred_amount *big.Int, locksroot [32]byte, extra_hash [32]byte, signature []byte) (*types.Transaction, error) {
	return _NettingChannelContract.contract.Transact(opts, "close", nonce, transferred_amount, locksroot, extra_hash, signature)
}

// Close is a paid mutator transaction binding the contract method 0x5e1fc56e.
//
// Solidity: function close(nonce uint64, transferred_amount uint256, locksroot bytes32, extra_hash bytes32, signature bytes) returns()
func (_NettingChannelContract *NettingChannelContractSession) Close(nonce uint64, transferred_amount *big.Int, locksroot [32]byte, extra_hash [32]byte, signature []byte) (*types.Transaction, error) {
	return _NettingChannelContract.Contract.Close(&_NettingChannelContract.TransactOpts, nonce, transferred_amount, locksroot, extra_hash, signature)
}

// Close is a paid mutator transaction binding the contract method 0x5e1fc56e.
//
// Solidity: function close(nonce uint64, transferred_amount uint256, locksroot bytes32, extra_hash bytes32, signature bytes) returns()
func (_NettingChannelContract *NettingChannelContractTransactorSession) Close(nonce uint64, transferred_amount *big.Int, locksroot [32]byte, extra_hash [32]byte, signature []byte) (*types.Transaction, error) {
	return _NettingChannelContract.Contract.Close(&_NettingChannelContract.TransactOpts, nonce, transferred_amount, locksroot, extra_hash, signature)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(amount uint256) returns(bool)
func (_NettingChannelContract *NettingChannelContractTransactor) Deposit(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _NettingChannelContract.contract.Transact(opts, "deposit", amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(amount uint256) returns(bool)
func (_NettingChannelContract *NettingChannelContractSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _NettingChannelContract.Contract.Deposit(&_NettingChannelContract.TransactOpts, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(amount uint256) returns(bool)
func (_NettingChannelContract *NettingChannelContractTransactorSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _NettingChannelContract.Contract.Deposit(&_NettingChannelContract.TransactOpts, amount)
}

// Settle is a paid mutator transaction binding the contract method 0x11da60b4.
//
// Solidity: function settle() returns()
func (_NettingChannelContract *NettingChannelContractTransactor) Settle(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NettingChannelContract.contract.Transact(opts, "settle")
}

// Settle is a paid mutator transaction binding the contract method 0x11da60b4.
//
// Solidity: function settle() returns()
func (_NettingChannelContract *NettingChannelContractSession) Settle() (*types.Transaction, error) {
	return _NettingChannelContract.Contract.Settle(&_NettingChannelContract.TransactOpts)
}

// Settle is a paid mutator transaction binding the contract method 0x11da60b4.
//
// Solidity: function settle() returns()
func (_NettingChannelContract *NettingChannelContractTransactorSession) Settle() (*types.Transaction, error) {
	return _NettingChannelContract.Contract.Settle(&_NettingChannelContract.TransactOpts)
}

// UpdateTransfer is a paid mutator transaction binding the contract method 0x27d120fe.
//
// Solidity: function updateTransfer(nonce uint64, transferred_amount uint256, locksroot bytes32, extra_hash bytes32, signature bytes) returns()
func (_NettingChannelContract *NettingChannelContractTransactor) UpdateTransfer(opts *bind.TransactOpts, nonce uint64, transferred_amount *big.Int, locksroot [32]byte, extra_hash [32]byte, signature []byte) (*types.Transaction, error) {
	return _NettingChannelContract.contract.Transact(opts, "updateTransfer", nonce, transferred_amount, locksroot, extra_hash, signature)
}

// UpdateTransfer is a paid mutator transaction binding the contract method 0x27d120fe.
//
// Solidity: function updateTransfer(nonce uint64, transferred_amount uint256, locksroot bytes32, extra_hash bytes32, signature bytes) returns()
func (_NettingChannelContract *NettingChannelContractSession) UpdateTransfer(nonce uint64, transferred_amount *big.Int, locksroot [32]byte, extra_hash [32]byte, signature []byte) (*types.Transaction, error) {
	return _NettingChannelContract.Contract.UpdateTransfer(&_NettingChannelContract.TransactOpts, nonce, transferred_amount, locksroot, extra_hash, signature)
}

// UpdateTransfer is a paid mutator transaction binding the contract method 0x27d120fe.
//
// Solidity: function updateTransfer(nonce uint64, transferred_amount uint256, locksroot bytes32, extra_hash bytes32, signature bytes) returns()
func (_NettingChannelContract *NettingChannelContractTransactorSession) UpdateTransfer(nonce uint64, transferred_amount *big.Int, locksroot [32]byte, extra_hash [32]byte, signature []byte) (*types.Transaction, error) {
	return _NettingChannelContract.Contract.UpdateTransfer(&_NettingChannelContract.TransactOpts, nonce, transferred_amount, locksroot, extra_hash, signature)
}

// UpdateTransferDelegate is a paid mutator transaction binding the contract method 0x59023f89.
//
// Solidity: function updateTransferDelegate(nonce uint64, transferred_amount uint256, locksroot bytes32, extra_hash bytes32, closing_signature bytes, non_closing_signature bytes) returns()
func (_NettingChannelContract *NettingChannelContractTransactor) UpdateTransferDelegate(opts *bind.TransactOpts, nonce uint64, transferred_amount *big.Int, locksroot [32]byte, extra_hash [32]byte, closing_signature []byte, non_closing_signature []byte) (*types.Transaction, error) {
	return _NettingChannelContract.contract.Transact(opts, "updateTransferDelegate", nonce, transferred_amount, locksroot, extra_hash, closing_signature, non_closing_signature)
}

// UpdateTransferDelegate is a paid mutator transaction binding the contract method 0x59023f89.
//
// Solidity: function updateTransferDelegate(nonce uint64, transferred_amount uint256, locksroot bytes32, extra_hash bytes32, closing_signature bytes, non_closing_signature bytes) returns()
func (_NettingChannelContract *NettingChannelContractSession) UpdateTransferDelegate(nonce uint64, transferred_amount *big.Int, locksroot [32]byte, extra_hash [32]byte, closing_signature []byte, non_closing_signature []byte) (*types.Transaction, error) {
	return _NettingChannelContract.Contract.UpdateTransferDelegate(&_NettingChannelContract.TransactOpts, nonce, transferred_amount, locksroot, extra_hash, closing_signature, non_closing_signature)
}

// UpdateTransferDelegate is a paid mutator transaction binding the contract method 0x59023f89.
//
// Solidity: function updateTransferDelegate(nonce uint64, transferred_amount uint256, locksroot bytes32, extra_hash bytes32, closing_signature bytes, non_closing_signature bytes) returns()
func (_NettingChannelContract *NettingChannelContractTransactorSession) UpdateTransferDelegate(nonce uint64, transferred_amount *big.Int, locksroot [32]byte, extra_hash [32]byte, closing_signature []byte, non_closing_signature []byte) (*types.Transaction, error) {
	return _NettingChannelContract.Contract.UpdateTransferDelegate(&_NettingChannelContract.TransactOpts, nonce, transferred_amount, locksroot, extra_hash, closing_signature, non_closing_signature)
}

// Withdraw is a paid mutator transaction binding the contract method 0x965f12e8.
//
// Solidity: function withdraw(participant address, locked_encoded bytes, merkle_proof bytes, secret bytes32) returns()
func (_NettingChannelContract *NettingChannelContractTransactor) Withdraw(opts *bind.TransactOpts, participant common.Address, locked_encoded []byte, merkle_proof []byte, secret [32]byte) (*types.Transaction, error) {
	return _NettingChannelContract.contract.Transact(opts, "withdraw", participant, locked_encoded, merkle_proof, secret)
}

// Withdraw is a paid mutator transaction binding the contract method 0x965f12e8.
//
// Solidity: function withdraw(participant address, locked_encoded bytes, merkle_proof bytes, secret bytes32) returns()
func (_NettingChannelContract *NettingChannelContractSession) Withdraw(participant common.Address, locked_encoded []byte, merkle_proof []byte, secret [32]byte) (*types.Transaction, error) {
	return _NettingChannelContract.Contract.Withdraw(&_NettingChannelContract.TransactOpts, participant, locked_encoded, merkle_proof, secret)
}

// Withdraw is a paid mutator transaction binding the contract method 0x965f12e8.
//
// Solidity: function withdraw(participant address, locked_encoded bytes, merkle_proof bytes, secret bytes32) returns()
func (_NettingChannelContract *NettingChannelContractTransactorSession) Withdraw(participant common.Address, locked_encoded []byte, merkle_proof []byte, secret [32]byte) (*types.Transaction, error) {
	return _NettingChannelContract.Contract.Withdraw(&_NettingChannelContract.TransactOpts, participant, locked_encoded, merkle_proof, secret)
}

// NettingChannelContractChannelClosedIterator is returned from FilterChannelClosed and is used to iterate over the raw logs and unpacked data for ChannelClosed events raised by the NettingChannelContract contract.
type NettingChannelContractChannelClosedIterator struct {
	Event *NettingChannelContractChannelClosed // Event containing the contract specifics and raw log

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
func (it *NettingChannelContractChannelClosedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NettingChannelContractChannelClosed)
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
		it.Event = new(NettingChannelContractChannelClosed)
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
func (it *NettingChannelContractChannelClosedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NettingChannelContractChannelClosedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NettingChannelContractChannelClosed represents a ChannelClosed event raised by the NettingChannelContract contract.
type NettingChannelContractChannelClosed struct {
	Registry_address common.Address
	Closing_address  common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterChannelClosed is a free log retrieval operation binding the contract event 0x48b749184840ce6faa9f6691fd4af8e7c969cd25ee881675131e2c9358ec3118.
//
// Solidity: event ChannelClosed(registry_address address, closing_address address)
func (_NettingChannelContract *NettingChannelContractFilterer) FilterChannelClosed(opts *bind.FilterOpts) (*NettingChannelContractChannelClosedIterator, error) {

	logs, sub, err := _NettingChannelContract.contract.FilterLogs(opts, "ChannelClosed")
	if err != nil {
		return nil, err
	}
	return &NettingChannelContractChannelClosedIterator{contract: _NettingChannelContract.contract, event: "ChannelClosed", logs: logs, sub: sub}, nil
}

// WatchChannelClosed is a free log subscription operation binding the contract event 0x48b749184840ce6faa9f6691fd4af8e7c969cd25ee881675131e2c9358ec3118.
//
// Solidity: event ChannelClosed(registry_address address, closing_address address)
func (_NettingChannelContract *NettingChannelContractFilterer) WatchChannelClosed(opts *bind.WatchOpts, sink chan<- *NettingChannelContractChannelClosed) (event.Subscription, error) {

	logs, sub, err := _NettingChannelContract.contract.WatchLogs(opts, "ChannelClosed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NettingChannelContractChannelClosed)
				if err := _NettingChannelContract.contract.UnpackLog(event, "ChannelClosed", log); err != nil {
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

// NettingChannelContractChannelNewBalanceIterator is returned from FilterChannelNewBalance and is used to iterate over the raw logs and unpacked data for ChannelNewBalance events raised by the NettingChannelContract contract.
type NettingChannelContractChannelNewBalanceIterator struct {
	Event *NettingChannelContractChannelNewBalance // Event containing the contract specifics and raw log

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
func (it *NettingChannelContractChannelNewBalanceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NettingChannelContractChannelNewBalance)
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
		it.Event = new(NettingChannelContractChannelNewBalance)
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
func (it *NettingChannelContractChannelNewBalanceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NettingChannelContractChannelNewBalanceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NettingChannelContractChannelNewBalance represents a ChannelNewBalance event raised by the NettingChannelContract contract.
type NettingChannelContractChannelNewBalance struct {
	Registry_address common.Address
	Token_address    common.Address
	Participant      common.Address
	Balance          *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterChannelNewBalance is a free log retrieval operation binding the contract event 0x61cf19a8c55ff5e9fda69ea0207b0350002d5073a364313b554d352ff6d7803d.
//
// Solidity: event ChannelNewBalance(registry_address address, token_address address, participant address, balance uint256)
func (_NettingChannelContract *NettingChannelContractFilterer) FilterChannelNewBalance(opts *bind.FilterOpts) (*NettingChannelContractChannelNewBalanceIterator, error) {

	logs, sub, err := _NettingChannelContract.contract.FilterLogs(opts, "ChannelNewBalance")
	if err != nil {
		return nil, err
	}
	return &NettingChannelContractChannelNewBalanceIterator{contract: _NettingChannelContract.contract, event: "ChannelNewBalance", logs: logs, sub: sub}, nil
}

// WatchChannelNewBalance is a free log subscription operation binding the contract event 0x61cf19a8c55ff5e9fda69ea0207b0350002d5073a364313b554d352ff6d7803d.
//
// Solidity: event ChannelNewBalance(registry_address address, token_address address, participant address, balance uint256)
func (_NettingChannelContract *NettingChannelContractFilterer) WatchChannelNewBalance(opts *bind.WatchOpts, sink chan<- *NettingChannelContractChannelNewBalance) (event.Subscription, error) {

	logs, sub, err := _NettingChannelContract.contract.WatchLogs(opts, "ChannelNewBalance")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NettingChannelContractChannelNewBalance)
				if err := _NettingChannelContract.contract.UnpackLog(event, "ChannelNewBalance", log); err != nil {
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

// NettingChannelContractChannelSecretRevealedIterator is returned from FilterChannelSecretRevealed and is used to iterate over the raw logs and unpacked data for ChannelSecretRevealed events raised by the NettingChannelContract contract.
type NettingChannelContractChannelSecretRevealedIterator struct {
	Event *NettingChannelContractChannelSecretRevealed // Event containing the contract specifics and raw log

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
func (it *NettingChannelContractChannelSecretRevealedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NettingChannelContractChannelSecretRevealed)
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
		it.Event = new(NettingChannelContractChannelSecretRevealed)
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
func (it *NettingChannelContractChannelSecretRevealedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NettingChannelContractChannelSecretRevealedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NettingChannelContractChannelSecretRevealed represents a ChannelSecretRevealed event raised by the NettingChannelContract contract.
type NettingChannelContractChannelSecretRevealed struct {
	Registry_address common.Address
	Secret           [32]byte
	Receiver_address common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterChannelSecretRevealed is a free log retrieval operation binding the contract event 0x15b5384ec609c8a364f4f4ba5969bf39b4e2b5c6a8364773267c682afe40586f.
//
// Solidity: event ChannelSecretRevealed(registry_address address, secret bytes32, receiver_address address)
func (_NettingChannelContract *NettingChannelContractFilterer) FilterChannelSecretRevealed(opts *bind.FilterOpts) (*NettingChannelContractChannelSecretRevealedIterator, error) {

	logs, sub, err := _NettingChannelContract.contract.FilterLogs(opts, "ChannelSecretRevealed")
	if err != nil {
		return nil, err
	}
	return &NettingChannelContractChannelSecretRevealedIterator{contract: _NettingChannelContract.contract, event: "ChannelSecretRevealed", logs: logs, sub: sub}, nil
}

// WatchChannelSecretRevealed is a free log subscription operation binding the contract event 0x15b5384ec609c8a364f4f4ba5969bf39b4e2b5c6a8364773267c682afe40586f.
//
// Solidity: event ChannelSecretRevealed(registry_address address, secret bytes32, receiver_address address)
func (_NettingChannelContract *NettingChannelContractFilterer) WatchChannelSecretRevealed(opts *bind.WatchOpts, sink chan<- *NettingChannelContractChannelSecretRevealed) (event.Subscription, error) {

	logs, sub, err := _NettingChannelContract.contract.WatchLogs(opts, "ChannelSecretRevealed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NettingChannelContractChannelSecretRevealed)
				if err := _NettingChannelContract.contract.UnpackLog(event, "ChannelSecretRevealed", log); err != nil {
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

// NettingChannelContractChannelSettledIterator is returned from FilterChannelSettled and is used to iterate over the raw logs and unpacked data for ChannelSettled events raised by the NettingChannelContract contract.
type NettingChannelContractChannelSettledIterator struct {
	Event *NettingChannelContractChannelSettled // Event containing the contract specifics and raw log

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
func (it *NettingChannelContractChannelSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NettingChannelContractChannelSettled)
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
		it.Event = new(NettingChannelContractChannelSettled)
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
func (it *NettingChannelContractChannelSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NettingChannelContractChannelSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NettingChannelContractChannelSettled represents a ChannelSettled event raised by the NettingChannelContract contract.
type NettingChannelContractChannelSettled struct {
	Registry_address common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterChannelSettled is a free log retrieval operation binding the contract event 0x40bab4574d5dad3320b539548362899f170b309febdeaadb2e4d0367311df4e0.
//
// Solidity: event ChannelSettled(registry_address address)
func (_NettingChannelContract *NettingChannelContractFilterer) FilterChannelSettled(opts *bind.FilterOpts) (*NettingChannelContractChannelSettledIterator, error) {

	logs, sub, err := _NettingChannelContract.contract.FilterLogs(opts, "ChannelSettled")
	if err != nil {
		return nil, err
	}
	return &NettingChannelContractChannelSettledIterator{contract: _NettingChannelContract.contract, event: "ChannelSettled", logs: logs, sub: sub}, nil
}

// WatchChannelSettled is a free log subscription operation binding the contract event 0x40bab4574d5dad3320b539548362899f170b309febdeaadb2e4d0367311df4e0.
//
// Solidity: event ChannelSettled(registry_address address)
func (_NettingChannelContract *NettingChannelContractFilterer) WatchChannelSettled(opts *bind.WatchOpts, sink chan<- *NettingChannelContractChannelSettled) (event.Subscription, error) {

	logs, sub, err := _NettingChannelContract.contract.WatchLogs(opts, "ChannelSettled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NettingChannelContractChannelSettled)
				if err := _NettingChannelContract.contract.UnpackLog(event, "ChannelSettled", log); err != nil {
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

// NettingChannelContractTransferUpdatedIterator is returned from FilterTransferUpdated and is used to iterate over the raw logs and unpacked data for TransferUpdated events raised by the NettingChannelContract contract.
type NettingChannelContractTransferUpdatedIterator struct {
	Event *NettingChannelContractTransferUpdated // Event containing the contract specifics and raw log

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
func (it *NettingChannelContractTransferUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NettingChannelContractTransferUpdated)
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
		it.Event = new(NettingChannelContractTransferUpdated)
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
func (it *NettingChannelContractTransferUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NettingChannelContractTransferUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NettingChannelContractTransferUpdated represents a TransferUpdated event raised by the NettingChannelContract contract.
type NettingChannelContractTransferUpdated struct {
	Node_address common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferUpdated is a free log retrieval operation binding the contract event 0xa0379b1bd0864245b4ff39bf6a023065e80d0e9276d2671d94b9f653b4bbcdfe.
//
// Solidity: event TransferUpdated(node_address address)
func (_NettingChannelContract *NettingChannelContractFilterer) FilterTransferUpdated(opts *bind.FilterOpts) (*NettingChannelContractTransferUpdatedIterator, error) {

	logs, sub, err := _NettingChannelContract.contract.FilterLogs(opts, "TransferUpdated")
	if err != nil {
		return nil, err
	}
	return &NettingChannelContractTransferUpdatedIterator{contract: _NettingChannelContract.contract, event: "TransferUpdated", logs: logs, sub: sub}, nil
}

// WatchTransferUpdated is a free log subscription operation binding the contract event 0xa0379b1bd0864245b4ff39bf6a023065e80d0e9276d2671d94b9f653b4bbcdfe.
//
// Solidity: event TransferUpdated(node_address address)
func (_NettingChannelContract *NettingChannelContractFilterer) WatchTransferUpdated(opts *bind.WatchOpts, sink chan<- *NettingChannelContractTransferUpdated) (event.Subscription, error) {

	logs, sub, err := _NettingChannelContract.contract.WatchLogs(opts, "TransferUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NettingChannelContractTransferUpdated)
				if err := _NettingChannelContract.contract.UnpackLog(event, "TransferUpdated", log); err != nil {
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

// NettingChannelLibraryABI is the input ABI used to generate the binding from.
const NettingChannelLibraryABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"self\",\"type\":\"NettingChannelLibrary.Data storage\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"},{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"self\",\"type\":\"NettingChannelLibrary.Data storage\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"transferred_amount\",\"type\":\"uint256\"},{\"name\":\"locksroot\",\"type\":\"bytes32\"},{\"name\":\"extra_hash\",\"type\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"close\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"self\",\"type\":\"NettingChannelLibrary.Data storage\"}],\"name\":\"settle\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"self\",\"type\":\"NettingChannelLibrary.Data storage\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"transferred_amount\",\"type\":\"uint256\"},{\"name\":\"locksroot\",\"type\":\"bytes32\"},{\"name\":\"extra_hash\",\"type\":\"bytes32\"},{\"name\":\"closing_signature\",\"type\":\"bytes\"},{\"name\":\"non_closing_signature\",\"type\":\"bytes\"}],\"name\":\"updateTransferDelegate\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"self\",\"type\":\"NettingChannelLibrary.Data storage\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"transferred_amount\",\"type\":\"uint256\"},{\"name\":\"locksroot\",\"type\":\"bytes32\"},{\"name\":\"extra_hash\",\"type\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"updateTransfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"self\",\"type\":\"NettingChannelLibrary.Data storage\"},{\"name\":\"participant\",\"type\":\"address\"},{\"name\":\"locked_encoded\",\"type\":\"bytes\"},{\"name\":\"merkle_proof\",\"type\":\"bytes\"},{\"name\":\"secret\",\"type\":\"bytes32\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// NettingChannelLibraryBin is the compiled bytecode used for deploying new contracts.
const NettingChannelLibraryBin = `0x611126610030600b82828239805160001a6073146000811461002057610022565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600436106100995763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416633268a05a811461009e578063b32c65c8146100d4578063c800b00214610151578063de394e0d146101c9578063f0023456146101e1578063f565eb36146102b1578063fee4658a14610327575b600080fd5b8180156100aa57600080fd5b506100b96004356024356103d0565b60408051921515835260208301919091528051918290030190f35b6100dc61059d565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101165781810151838201526020016100fe565b50505050905090810190601f1680156101435780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b81801561015d57600080fd5b50604080516020600460a43581810135601f81018490048402850184019095528484526101c7948235946024803567ffffffffffffffff16956044359560643595608435953695929460c49492019181908401838280828437509497506105d49650505050505050565b005b8180156101d557600080fd5b506101c76004356106b3565b8180156101ed57600080fd5b50604080516020600460a43581810135601f8101849004840285018401909552848452610295948235946024803567ffffffffffffffff16956044359560643595608435953695929460c494920191819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a9998810197919650918201945092508291508401838280828437509497506108e69650505050505050565b60408051600160a060020a039092168252519081900360200190f35b8180156102bd57600080fd5b50604080516020600460a43581810135601f81018490048402850184019095528484526101c7948235946024803567ffffffffffffffff16956044359560643595608435953695929460c4949201918190840183828082843750949750610a5a9650505050505050565b81801561033357600080fd5b50604080516020600460443581810135601f81018490048402850184019095528484526101c79482359460248035600160a060020a03169536959460649492019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a9998810197919650918201945092508291508401838280828437509497505093359450610b9a9350505050565b600080600080600086600101541115156103e957600080fd5b6002860154156103f857600080fd5b6005860154604080517f70a08231000000000000000000000000000000000000000000000000000000008152600160a060020a0333811660048301529151889392909216916370a08231916024808201926020929091908290030181600087803b15801561046557600080fd5b505af1158015610479573d6000803e3d6000fd5b505050506040513d602081101561048f57600080fd5b5051101561049c57600080fd5b6104a68633610cac565b91506006860160ff8316600281106104ba57fe5b6005880154604080517f23b872dd000000000000000000000000000000000000000000000000000000008152600160a060020a0333811660048301523081166024830152604482018b9052915160069490940294909401945016916323b872dd9160648083019260209291908290030181600087803b15801561053c57600080fd5b505af1158015610550573d6000803e3d6000fd5b505050506040513d602081101561056657600080fd5b505193506001841515141561058c57600180820180548701908190559094509250610594565b600093508392505b50509250929050565b60408051808201909152600581527f302e322e5f000000000000000000000000000000000000000000000000000000602082015281565b600080600080896002015460001415156105ed57600080fd5b4360028b01556105fd8a33610cac565b60038b01805473ffffffffffffffffffffffffffffffffffffffff191633600160a060020a0316179055855160ff919091169350604114156106a7576106468989898989610cde565b93506106528a85610cac565b60ff1691508282141561066457600080fd5b60068a01826002811061067357fe5b6006020160048101805467ffffffffffffffff191667ffffffffffffffff8c16179055600281018890556003810189905590505b50505050505050505050565b60008060008060008060008088600081600201541115156106d357600080fd5b895460028b01548b9143910111156106ea57600080fd5b60038b0154610703908c90600160a060020a0316610cac565b995060018a9003985060068b0160ff8b166002811061071e57fe5b6006020193508a6006018960ff1660028110151561073857fe5b60060201925082600301548460030154846001015401039650826001015484600101540197506107688789610de6565b9450610775856000610dfe565b94508488039550600085111561082f5760058b01548354604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a039283166004820152602481018990529051919092169163a9059cbb9160448083019260209291908290030181600087803b1580156107f857600080fd5b505af115801561080c573d6000803e3d6000fd5b505050506040513d602081101561082257600080fd5b5051151561082f57600080fd5b60008611156108e25760058b01548454604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a039283166004820152602481018a90529051919092169163a9059cbb9160448083019260209291908290030181600087803b1580156108ab57600080fd5b505af11580156108bf573d6000803e3d6000fd5b505050506040513d60208110156108d557600080fd5b505115156108e257600080fd5b6000ff5b60008060008060008b6000816002015411151561090257600080fd5b8c5460028e01548e91439101101561091957600080fd5b60138e015460ff161561092b57600080fd5b60138e01805460ff191660011790556002808f01548f544391909103909102101561095557600080fd5b6109638d8d8d8d8d8d610e14565b945061096f8e86610cac565b60038f0154909450600160a060020a038681169116141561098f57600080fd5b61099c8d8d8d8d8d610cde565b60038f0154909650600160a060020a038088169116146109bb57600080fd5b600184900392508c60068f0160ff8516600281106109d557fe5b6006020160040160006101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055508a8e6006018460ff16600281101515610a1757fe5b6006020160020181600019169055508b8e6006018460ff16600281101515610a3b57fe5b6006020160030181905550849650505050505050979650505050505050565b60008060008860008160020154111515610a7357600080fd5b895460028b01548b914391011015610a8a57600080fd5b60138b015460ff1615610a9c57600080fd5b60138b01805460ff19166001179055610ab58b33610cac565b60038c015490945033600160a060020a0390811691161415610ad657600080fd5b610ae38a8a8a8a8a610cde565b60038c0154909550600160a060020a03808716911614610b0257600080fd5b600184900392508960068c0160ff851660028110610b1c57fe5b6006020160040160006101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550878b6006018460ff16600281101515610b5e57fe5b600602016002018160001916905550888b6006018460ff16600281101515610b8257fe5b60060201600301819055505050505050505050505050565b6000806000806000808a60008160020154111515610bb757600080fd5b610bc18c8c610cac565b600103955060068c0160ff871660028110610bd857fe5b600602019150816002015460001916600060010214151515610bf957600080fd5b610c028a610fb8565b600081815260058601602052604090205491995091965090935060ff1615610c2957600080fd5b60008381526005830160205260409020805460ff191660011790554367ffffffffffffffff86161015610c5b57600080fd5b6040805189815290519081900360200190208314610c7857600080fd5b610c828a8a610fe4565b60028301549094508414610c9557600080fd5b506003018054909501909455505050505050505050565b600160a060020a038116600090815260128301602052604081205460ff16801515610cd357fe5b600019019392505050565b600080600080600085516041141515610cf657600080fd5b5060408051780100000000000000000000000000000000000000000000000067ffffffffffffffff8c16028152600881018a9052602881018990526c01000000000000000000000000600160a060020a033016026048820152605c8101889052905190819003607c019020610d6a866110c4565b604080516000808252602080830180855288905260ff8516838501526060830187905260808301869052925195995093975091955060019360a0808401949293601f19830193908390039091019190865af1158015610dcd573d6000803e3d6000fd5b5050604051601f1901519b9a5050505050505050505050565b6000818311610df55782610df7565b815b9392505050565b6000818311610e0d5781610df7565b5090919050565b600080600080600086516041141515610e2c57600080fd5b8551604114610e3a57600080fd5b60405167ffffffffffffffff8c167801000000000000000000000000000000000000000000000000028152600881018b9052602881018a905230600160a060020a0381166c01000000000000000000000000026048830152605c82018a905288518d928d928d9290918d918d913391607c82019060208501908083835b60208310610ed65780518252601f199092019160209182019101610eb7565b6001836020036101000a03801982511681845116808217855250505050505090500182600160a060020a0316600160a060020a03166c0100000000000000000000000002815260140197505050505050505060405180910390209050610f3b866110c4565b604080516000808252602080830180855288905260ff8516838501526060830187905260808301869052925195995093975091955060019360a0808401949293601f19830193908390039091019190865af1158015610f9e573d6000803e3d6000fd5b5050604051601f1901519c9b505050505050505050505050565b600080600083516048141515610fcd57600080fd5b505050600881015160288201516048909201519092565b60008060008060208551811515610ff757fe5b061561100257600080fd5b856040518082805190602001908083835b602083106110325780518252601f199092019160209182019101611013565b51815160209384036101000a60001901801990921691161790526040519190930181900390209196509094505050505b845183116110bb5750838201518082101561109557604080519283526020830182905280519283900301909120906110b0565b60408051828152602081019390935280519283900301909120905b602083019250611062565b50949350505050565b60208101516040820151604183015160ff16601b8114806110e857508060ff16601c145b15156110f357600080fd5b91939092505600a165627a7a723058202b65ef8122918fe48e50d391a6a1389b8d3987ece6049623ade9d7ff6a3ce6480029`

// DeployNettingChannelLibrary deploys a new Ethereum contract, binding an instance of NettingChannelLibrary to it.
func DeployNettingChannelLibrary(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *NettingChannelLibrary, error) {
	parsed, err := abi.JSON(strings.NewReader(NettingChannelLibraryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(NettingChannelLibraryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NettingChannelLibrary{NettingChannelLibraryCaller: NettingChannelLibraryCaller{contract: contract}, NettingChannelLibraryTransactor: NettingChannelLibraryTransactor{contract: contract}, NettingChannelLibraryFilterer: NettingChannelLibraryFilterer{contract: contract}}, nil
}

// NettingChannelLibrary is an auto generated Go binding around an Ethereum contract.
type NettingChannelLibrary struct {
	NettingChannelLibraryCaller     // Read-only binding to the contract
	NettingChannelLibraryTransactor // Write-only binding to the contract
	NettingChannelLibraryFilterer   // Log filterer for contract events
}

// NettingChannelLibraryCaller is an auto generated read-only Go binding around an Ethereum contract.
type NettingChannelLibraryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NettingChannelLibraryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NettingChannelLibraryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NettingChannelLibraryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NettingChannelLibraryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NettingChannelLibrarySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NettingChannelLibrarySession struct {
	Contract     *NettingChannelLibrary // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// NettingChannelLibraryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NettingChannelLibraryCallerSession struct {
	Contract *NettingChannelLibraryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// NettingChannelLibraryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NettingChannelLibraryTransactorSession struct {
	Contract     *NettingChannelLibraryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// NettingChannelLibraryRaw is an auto generated low-level Go binding around an Ethereum contract.
type NettingChannelLibraryRaw struct {
	Contract *NettingChannelLibrary // Generic contract binding to access the raw methods on
}

// NettingChannelLibraryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NettingChannelLibraryCallerRaw struct {
	Contract *NettingChannelLibraryCaller // Generic read-only contract binding to access the raw methods on
}

// NettingChannelLibraryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NettingChannelLibraryTransactorRaw struct {
	Contract *NettingChannelLibraryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNettingChannelLibrary creates a new instance of NettingChannelLibrary, bound to a specific deployed contract.
func NewNettingChannelLibrary(address common.Address, backend bind.ContractBackend) (*NettingChannelLibrary, error) {
	contract, err := bindNettingChannelLibrary(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NettingChannelLibrary{NettingChannelLibraryCaller: NettingChannelLibraryCaller{contract: contract}, NettingChannelLibraryTransactor: NettingChannelLibraryTransactor{contract: contract}, NettingChannelLibraryFilterer: NettingChannelLibraryFilterer{contract: contract}}, nil
}

// NewNettingChannelLibraryCaller creates a new read-only instance of NettingChannelLibrary, bound to a specific deployed contract.
func NewNettingChannelLibraryCaller(address common.Address, caller bind.ContractCaller) (*NettingChannelLibraryCaller, error) {
	contract, err := bindNettingChannelLibrary(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NettingChannelLibraryCaller{contract: contract}, nil
}

// NewNettingChannelLibraryTransactor creates a new write-only instance of NettingChannelLibrary, bound to a specific deployed contract.
func NewNettingChannelLibraryTransactor(address common.Address, transactor bind.ContractTransactor) (*NettingChannelLibraryTransactor, error) {
	contract, err := bindNettingChannelLibrary(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NettingChannelLibraryTransactor{contract: contract}, nil
}

// NewNettingChannelLibraryFilterer creates a new log filterer instance of NettingChannelLibrary, bound to a specific deployed contract.
func NewNettingChannelLibraryFilterer(address common.Address, filterer bind.ContractFilterer) (*NettingChannelLibraryFilterer, error) {
	contract, err := bindNettingChannelLibrary(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NettingChannelLibraryFilterer{contract: contract}, nil
}

// bindNettingChannelLibrary binds a generic wrapper to an already deployed contract.
func bindNettingChannelLibrary(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NettingChannelLibraryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NettingChannelLibrary *NettingChannelLibraryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _NettingChannelLibrary.Contract.NettingChannelLibraryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NettingChannelLibrary *NettingChannelLibraryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NettingChannelLibrary.Contract.NettingChannelLibraryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NettingChannelLibrary *NettingChannelLibraryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NettingChannelLibrary.Contract.NettingChannelLibraryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NettingChannelLibrary *NettingChannelLibraryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _NettingChannelLibrary.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NettingChannelLibrary *NettingChannelLibraryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NettingChannelLibrary.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NettingChannelLibrary *NettingChannelLibraryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NettingChannelLibrary.Contract.contract.Transact(opts, method, params...)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_NettingChannelLibrary *NettingChannelLibraryCaller) Contract_version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _NettingChannelLibrary.contract.Call(opts, out, "contract_version")
	return *ret0, err
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_NettingChannelLibrary *NettingChannelLibrarySession) Contract_version() (string, error) {
	return _NettingChannelLibrary.Contract.Contract_version(&_NettingChannelLibrary.CallOpts)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_NettingChannelLibrary *NettingChannelLibraryCallerSession) Contract_version() (string, error) {
	return _NettingChannelLibrary.Contract.Contract_version(&_NettingChannelLibrary.CallOpts)
}

// RegistryABI is the input ABI used to generate the binding from.
const RegistryABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"registry\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"token_address\",\"type\":\"address\"}],\"name\":\"channelManagerByToken\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokens\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"channelManagerAddresses\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenAddresses\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token_address\",\"type\":\"address\"}],\"name\":\"addToken\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"registry_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"token_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"channel_manager_address\",\"type\":\"address\"}],\"name\":\"TokenAdded\",\"type\":\"event\"}]"

// RegistryBin is the compiled bytecode used for deploying new contracts.
var RegistryBin = `0x608060405234801561001057600080fd5b50610e43806100206000396000f3006080604052600436106100825763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663038defd7811461009457806325119b5f146100d15780634f64b2be146100f2578063640191e21461010a578063a9989b931461016f578063b32c65c814610184578063d48bfca71461020e575b34801561008e57600080fd5b50600080fd5b3480156100a057600080fd5b506100b5600160a060020a036004351661022f565b60408051600160a060020a039092168252519081900360200190f35b3480156100dd57600080fd5b506100b5600160a060020a036004351661024a565b3480156100fe57600080fd5b506100b5600435610294565b34801561011657600080fd5b5061011f6102bc565b60408051602080825283518183015283519192839290830191858101910280838360005b8381101561015b578181015183820152602001610143565b505050509050019250505060405180910390f35b34801561017b57600080fd5b5061011f610376565b34801561019057600080fd5b506101996103d8565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101d35781810151838201526020016101bb565b50505050905090810190601f1680156102005780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561021a57600080fd5b506100b5600160a060020a036004351661040f565b600060208190529081526040902054600160a060020a031681565b600160a060020a038082166000908152602081905260408120549091839116151561027457600080fd5b5050600160a060020a039081166000908152602081905260409020541690565b60018054829081106102a257fe5b600091825260209091200154600160a060020a0316905081565b606060008060606001805490506040519080825280602002602001820160405280156102f2578160200160208202803883390190505b509050600092505b60015483101561036f57600180548490811061031257fe5b6000918252602080832090910154600160a060020a039081168084529183905260409092205483519194509091169082908590811061034d57fe5b600160a060020a039092166020928302909101909101526001909201916102fa565b9392505050565b606060018054806020026020016040519081016040528092919081815260200182805480156103ce57602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116103b0575b5050505050905090565b60408051808201909152600581527f302e322e5f000000000000000000000000000000000000000000000000000000602082015281565b600160a060020a038082166000908152602081905260408120549091829184918391161561043c57600080fd5b81905080600160a060020a03166318160ddd6040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401602060405180830381600087803b15801561049657600080fd5b505af11580156104aa573d6000803e3d6000fd5b505050506040513d60208110156104c057600080fd5b50309050856104cd6105c2565b600160a060020a03928316815291166020820152604080519182900301906000f080158015610500573d6000803e3d6000fd5b50600160a060020a03808716600081815260208181526040808320805486881673ffffffffffffffffffffffffffffffffffffffff19918216811790925560018054808201825595527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf69094018054909416851790935580513090951685529084019290925282820152519194507ffc43233c964efa713b168e2361b2c57eafddc32aa7f7d0f85c92e66e113aa28a919081900360600190a150909392505050565b604051610845806105d3833901905600608060405234801561001057600080fd5b5060405160408061084583398101604052805160209091015160018054600160a060020a03938416600160a060020a031991821617909155600080549390921692169190911790556107de806100676000396000f30060806040526004361061008d5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630b74b620811461009f578063238bfba2146101045780636785b500146101415780636cb30fee146101565780637709bc78146101775780639d76ea58146101ac578063b32c65c8146101c1578063f26c6aed1461024b575b34801561009957600080fd5b50600080fd5b3480156100ab57600080fd5b506100b461026f565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156100f05781810151838201526020016100d8565b505050509050019250505060405180910390f35b34801561011057600080fd5b50610125600160a060020a0360043516610457565b60408051600160a060020a039092168252519081900360200190f35b34801561014d57600080fd5b506100b4610505565b34801561016257600080fd5b506100b4600160a060020a036004351661056a565b34801561018357600080fd5b50610198600160a060020a03600435166105e0565b604080519115158252519081900360200190f35b3480156101b857600080fd5b506101256105e8565b3480156101cd57600080fd5b506101d66105f7565b6040805160208082528351818301528351919283929083019185019080838360005b838110156102105781810151838201526020016101f8565b50505050905090810190601f16801561023d5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561025757600080fd5b50610125600160a060020a036004351660243561062e565b606060008082818080805b6002548710156102c557600280546102b191908990811061029757fe5b600091825260209091200154600160a060020a03166105e0565b156102ba576001015b60019096019561027a565b806002026040519080825280602002602001820160405280156102f2578160200160208202803883390190505b50945060009550600096505b60025487101561044b576002805461031b91908990811061029757fe5b151561032657610440565b600280548890811061033457fe5b9060005260206000200160009054906101000a9004600160a060020a0316915081600160a060020a0316636d2381b36040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401608060405180830381600087803b1580156103ab57600080fd5b505af11580156103bf573d6000803e3d6000fd5b505050506040513d60808110156103d557600080fd5b5080516040909101518651919550935084908690889081106103f357fe5b600160a060020a03909216602092830290910190910152845160019690960195839086908890811061042157fe5b600160a060020a03909216602092830290910190910152600195909501945b6001909601956102fe565b50929695505050505050565b604080517f8a1c00e2000000000000000000000000000000000000000000000000000000008152600060048201819052600160a060020a0384166024830152915173__ChannelManagerLibrary.sol:ChannelMan__91638a1c00e2916044808301926020929190829003018186803b1580156104d357600080fd5b505af41580156104e7573d6000803e3d6000fd5b505050506040513d60208110156104fd57600080fd5b505192915050565b6060600060020180548060200260200160405190810160405280929190818152602001828054801561056057602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610542575b5050505050905090565b600160a060020a0381166000908152600460209081526040918290208054835181840281018401909452808452606093928301828280156105d457602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116105b6575b50505050509050919050565b6000903b1190565b600054600160a060020a031690565b60408051808201909152600581527f302e322e5f000000000000000000000000000000000000000000000000000000602082015281565b600080600061063c85610457565b9150600160a060020a0382161561069c5760015460408051600160a060020a039283168152338316602082015291871682820152517f91cd7e9ad7c88602bc7b06adc62de54cb670665c75894414c01ead5f2baf09309181900360600190a15b604080517f941583a500000000000000000000000000000000000000000000000000000000815260006004820152600160a060020a038716602482015260448101869052905173__ChannelManagerLibrary.sol:ChannelMan__9163941583a5916064808301926020929190829003018186803b15801561071d57600080fd5b505af4158015610731573d6000803e3d6000fd5b505050506040513d602081101561074757600080fd5b505160015460408051600160a060020a039283168152828416602082015233831681830152918816606083015260808201879052519192507fc15f36191ba1aefd153ae0ece89afc6004a66b158bec44593e70bccff2ae7a0f919081900360a00190a19493505050505600a165627a7a723058204ca9dc2ef3777511038e010aaee35e13c0e7f876fc9f47fae0f677b65a04cf9c0029a165627a7a72305820178176ebe89d9837f67e187a349fa2e152c7fed608861b58262bc5329f5582260029`

// DeployRegistry deploys a new Ethereum contract, binding an instance of Registry to it.
func DeployRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Registry, error) {
	parsed, err := abi.JSON(strings.NewReader(RegistryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(RegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Registry{RegistryCaller: RegistryCaller{contract: contract}, RegistryTransactor: RegistryTransactor{contract: contract}, RegistryFilterer: RegistryFilterer{contract: contract}}, nil
}

// Registry is an auto generated Go binding around an Ethereum contract.
type Registry struct {
	RegistryCaller     // Read-only binding to the contract
	RegistryTransactor // Write-only binding to the contract
	RegistryFilterer   // Log filterer for contract events
}

// RegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type RegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RegistrySession struct {
	Contract     *Registry         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RegistryCallerSession struct {
	Contract *RegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// RegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RegistryTransactorSession struct {
	Contract     *RegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// RegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type RegistryRaw struct {
	Contract *Registry // Generic contract binding to access the raw methods on
}

// RegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RegistryCallerRaw struct {
	Contract *RegistryCaller // Generic read-only contract binding to access the raw methods on
}

// RegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RegistryTransactorRaw struct {
	Contract *RegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRegistry creates a new instance of Registry, bound to a specific deployed contract.
func NewRegistry(address common.Address, backend bind.ContractBackend) (*Registry, error) {
	contract, err := bindRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Registry{RegistryCaller: RegistryCaller{contract: contract}, RegistryTransactor: RegistryTransactor{contract: contract}, RegistryFilterer: RegistryFilterer{contract: contract}}, nil
}

// NewRegistryCaller creates a new read-only instance of Registry, bound to a specific deployed contract.
func NewRegistryCaller(address common.Address, caller bind.ContractCaller) (*RegistryCaller, error) {
	contract, err := bindRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryCaller{contract: contract}, nil
}

// NewRegistryTransactor creates a new write-only instance of Registry, bound to a specific deployed contract.
func NewRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*RegistryTransactor, error) {
	contract, err := bindRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryTransactor{contract: contract}, nil
}

// NewRegistryFilterer creates a new log filterer instance of Registry, bound to a specific deployed contract.
func NewRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*RegistryFilterer, error) {
	contract, err := bindRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RegistryFilterer{contract: contract}, nil
}

// bindRegistry binds a generic wrapper to an already deployed contract.
func bindRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.RegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transact(opts, method, params...)
}

// ChannelManagerAddresses is a free data retrieval call binding the contract method 0x640191e2.
//
// Solidity: function channelManagerAddresses() constant returns(address[])
func (_Registry *RegistryCaller) ChannelManagerAddresses(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "channelManagerAddresses")
	return *ret0, err
}

// ChannelManagerAddresses is a free data retrieval call binding the contract method 0x640191e2.
//
// Solidity: function channelManagerAddresses() constant returns(address[])
func (_Registry *RegistrySession) ChannelManagerAddresses() ([]common.Address, error) {
	return _Registry.Contract.ChannelManagerAddresses(&_Registry.CallOpts)
}

// ChannelManagerAddresses is a free data retrieval call binding the contract method 0x640191e2.
//
// Solidity: function channelManagerAddresses() constant returns(address[])
func (_Registry *RegistryCallerSession) ChannelManagerAddresses() ([]common.Address, error) {
	return _Registry.Contract.ChannelManagerAddresses(&_Registry.CallOpts)
}

// ChannelManagerByToken is a free data retrieval call binding the contract method 0x25119b5f.
//
// Solidity: function channelManagerByToken(token_address address) constant returns(address)
func (_Registry *RegistryCaller) ChannelManagerByToken(opts *bind.CallOpts, token_address common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "channelManagerByToken", token_address)
	return *ret0, err
}

// ChannelManagerByToken is a free data retrieval call binding the contract method 0x25119b5f.
//
// Solidity: function channelManagerByToken(token_address address) constant returns(address)
func (_Registry *RegistrySession) ChannelManagerByToken(token_address common.Address) (common.Address, error) {
	return _Registry.Contract.ChannelManagerByToken(&_Registry.CallOpts, token_address)
}

// ChannelManagerByToken is a free data retrieval call binding the contract method 0x25119b5f.
//
// Solidity: function channelManagerByToken(token_address address) constant returns(address)
func (_Registry *RegistryCallerSession) ChannelManagerByToken(token_address common.Address) (common.Address, error) {
	return _Registry.Contract.ChannelManagerByToken(&_Registry.CallOpts, token_address)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_Registry *RegistryCaller) Contract_version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "contract_version")
	return *ret0, err
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_Registry *RegistrySession) Contract_version() (string, error) {
	return _Registry.Contract.Contract_version(&_Registry.CallOpts)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_Registry *RegistryCallerSession) Contract_version() (string, error) {
	return _Registry.Contract.Contract_version(&_Registry.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x038defd7.
//
// Solidity: function registry( address) constant returns(address)
func (_Registry *RegistryCaller) Registry(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "registry", arg0)
	return *ret0, err
}

// Registry is a free data retrieval call binding the contract method 0x038defd7.
//
// Solidity: function registry( address) constant returns(address)
func (_Registry *RegistrySession) Registry(arg0 common.Address) (common.Address, error) {
	return _Registry.Contract.Registry(&_Registry.CallOpts, arg0)
}

// Registry is a free data retrieval call binding the contract method 0x038defd7.
//
// Solidity: function registry( address) constant returns(address)
func (_Registry *RegistryCallerSession) Registry(arg0 common.Address) (common.Address, error) {
	return _Registry.Contract.Registry(&_Registry.CallOpts, arg0)
}

// TokenAddresses is a free data retrieval call binding the contract method 0xa9989b93.
//
// Solidity: function tokenAddresses() constant returns(address[])
func (_Registry *RegistryCaller) TokenAddresses(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "tokenAddresses")
	return *ret0, err
}

// TokenAddresses is a free data retrieval call binding the contract method 0xa9989b93.
//
// Solidity: function tokenAddresses() constant returns(address[])
func (_Registry *RegistrySession) TokenAddresses() ([]common.Address, error) {
	return _Registry.Contract.TokenAddresses(&_Registry.CallOpts)
}

// TokenAddresses is a free data retrieval call binding the contract method 0xa9989b93.
//
// Solidity: function tokenAddresses() constant returns(address[])
func (_Registry *RegistryCallerSession) TokenAddresses() ([]common.Address, error) {
	return _Registry.Contract.TokenAddresses(&_Registry.CallOpts)
}

// Tokens is a free data retrieval call binding the contract method 0x4f64b2be.
//
// Solidity: function tokens( uint256) constant returns(address)
func (_Registry *RegistryCaller) Tokens(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "tokens", arg0)
	return *ret0, err
}

// Tokens is a free data retrieval call binding the contract method 0x4f64b2be.
//
// Solidity: function tokens( uint256) constant returns(address)
func (_Registry *RegistrySession) Tokens(arg0 *big.Int) (common.Address, error) {
	return _Registry.Contract.Tokens(&_Registry.CallOpts, arg0)
}

// Tokens is a free data retrieval call binding the contract method 0x4f64b2be.
//
// Solidity: function tokens( uint256) constant returns(address)
func (_Registry *RegistryCallerSession) Tokens(arg0 *big.Int) (common.Address, error) {
	return _Registry.Contract.Tokens(&_Registry.CallOpts, arg0)
}

// AddToken is a paid mutator transaction binding the contract method 0xd48bfca7.
//
// Solidity: function addToken(token_address address) returns(address)
func (_Registry *RegistryTransactor) AddToken(opts *bind.TransactOpts, token_address common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "addToken", token_address)
}

// AddToken is a paid mutator transaction binding the contract method 0xd48bfca7.
//
// Solidity: function addToken(token_address address) returns(address)
func (_Registry *RegistrySession) AddToken(token_address common.Address) (*types.Transaction, error) {
	return _Registry.Contract.AddToken(&_Registry.TransactOpts, token_address)
}

// AddToken is a paid mutator transaction binding the contract method 0xd48bfca7.
//
// Solidity: function addToken(token_address address) returns(address)
func (_Registry *RegistryTransactorSession) AddToken(token_address common.Address) (*types.Transaction, error) {
	return _Registry.Contract.AddToken(&_Registry.TransactOpts, token_address)
}

// RegistryTokenAddedIterator is returned from FilterTokenAdded and is used to iterate over the raw logs and unpacked data for TokenAdded events raised by the Registry contract.
type RegistryTokenAddedIterator struct {
	Event *RegistryTokenAdded // Event containing the contract specifics and raw log

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
func (it *RegistryTokenAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryTokenAdded)
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
		it.Event = new(RegistryTokenAdded)
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
func (it *RegistryTokenAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryTokenAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryTokenAdded represents a TokenAdded event raised by the Registry contract.
type RegistryTokenAdded struct {
	Registry_address        common.Address
	Token_address           common.Address
	Channel_manager_address common.Address
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterTokenAdded is a free log retrieval operation binding the contract event 0xfc43233c964efa713b168e2361b2c57eafddc32aa7f7d0f85c92e66e113aa28a.
//
// Solidity: event TokenAdded(registry_address address, token_address address, channel_manager_address address)
func (_Registry *RegistryFilterer) FilterTokenAdded(opts *bind.FilterOpts) (*RegistryTokenAddedIterator, error) {

	logs, sub, err := _Registry.contract.FilterLogs(opts, "TokenAdded")
	if err != nil {
		return nil, err
	}
	return &RegistryTokenAddedIterator{contract: _Registry.contract, event: "TokenAdded", logs: logs, sub: sub}, nil
}

// WatchTokenAdded is a free log subscription operation binding the contract event 0xfc43233c964efa713b168e2361b2c57eafddc32aa7f7d0f85c92e66e113aa28a.
//
// Solidity: event TokenAdded(registry_address address, token_address address, channel_manager_address address)
func (_Registry *RegistryFilterer) WatchTokenAdded(opts *bind.WatchOpts, sink chan<- *RegistryTokenAdded) (event.Subscription, error) {

	logs, sub, err := _Registry.contract.WatchLogs(opts, "TokenAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryTokenAdded)
				if err := _Registry.contract.UnpackLog(event, "TokenAdded", log); err != nil {
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

// TokenABI is the input ABI used to generate the binding from.
const TokenABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"supply\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// TokenBin is the compiled bytecode used for deploying new contracts.
const TokenBin = `0x`

// DeployToken deploys a new Ethereum contract, binding an instance of Token to it.
func DeployToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Token, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Token{TokenCaller: TokenCaller{contract: contract}, TokenTransactor: TokenTransactor{contract: contract}, TokenFilterer: TokenFilterer{contract: contract}}, nil
}

// Token is an auto generated Go binding around an Ethereum contract.
type Token struct {
	TokenCaller     // Read-only binding to the contract
	TokenTransactor // Write-only binding to the contract
	TokenFilterer   // Log filterer for contract events
}

// TokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenSession struct {
	Contract     *Token            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenCallerSession struct {
	Contract *TokenCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// TokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenTransactorSession struct {
	Contract     *TokenTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenRaw struct {
	Contract *Token // Generic contract binding to access the raw methods on
}

// TokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenCallerRaw struct {
	Contract *TokenCaller // Generic read-only contract binding to access the raw methods on
}

// TokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenTransactorRaw struct {
	Contract *TokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewToken creates a new instance of Token, bound to a specific deployed contract.
func NewToken(address common.Address, backend bind.ContractBackend) (*Token, error) {
	contract, err := bindToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Token{TokenCaller: TokenCaller{contract: contract}, TokenTransactor: TokenTransactor{contract: contract}, TokenFilterer: TokenFilterer{contract: contract}}, nil
}

// NewTokenCaller creates a new read-only instance of Token, bound to a specific deployed contract.
func NewTokenCaller(address common.Address, caller bind.ContractCaller) (*TokenCaller, error) {
	contract, err := bindToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCaller{contract: contract}, nil
}

// NewTokenTransactor creates a new write-only instance of Token, bound to a specific deployed contract.
func NewTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenTransactor, error) {
	contract, err := bindToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenTransactor{contract: contract}, nil
}

// NewTokenFilterer creates a new log filterer instance of Token, bound to a specific deployed contract.
func NewTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenFilterer, error) {
	contract, err := bindToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenFilterer{contract: contract}, nil
}

// bindToken binds a generic wrapper to an already deployed contract.
func bindToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Token *TokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Token.Contract.TokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Token *TokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.Contract.TokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Token *TokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Token.Contract.TokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Token *TokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Token.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Token *TokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Token *TokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Token.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_Token *TokenCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "allowance", _owner, _spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_Token *TokenSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _Token.Contract.Allowance(&_Token.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_Token *TokenCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _Token.Contract.Allowance(&_Token.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_Token *TokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_Token *TokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Token.Contract.BalanceOf(&_Token.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_Token *TokenCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Token.Contract.BalanceOf(&_Token.CallOpts, _owner)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(supply uint256)
func (_Token *TokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(supply uint256)
func (_Token *TokenSession) TotalSupply() (*big.Int, error) {
	return _Token.Contract.TotalSupply(&_Token.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(supply uint256)
func (_Token *TokenCallerSession) TotalSupply() (*big.Int, error) {
	return _Token.Contract.TotalSupply(&_Token.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_Token *TokenTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_Token *TokenSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Approve(&_Token.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_Token *TokenTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Approve(&_Token.TransactOpts, _spender, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_Token *TokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_Token *TokenSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Transfer(&_Token.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_Token *TokenTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Transfer(&_Token.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_Token *TokenTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_Token *TokenSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.TransferFrom(&_Token.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_Token *TokenTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.TransferFrom(&_Token.TransactOpts, _from, _to, _value)
}

// TokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Token contract.
type TokenApprovalIterator struct {
	Event *TokenApproval // Event containing the contract specifics and raw log

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
func (it *TokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenApproval)
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
		it.Event = new(TokenApproval)
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
func (it *TokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenApproval represents a Approval event raised by the Token contract.
type TokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(_owner indexed address, _spender indexed address, _value uint256)
func (_Token *TokenFilterer) FilterApproval(opts *bind.FilterOpts, _owner []common.Address, _spender []common.Address) (*TokenApprovalIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _Token.contract.FilterLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return &TokenApprovalIterator{contract: _Token.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(_owner indexed address, _spender indexed address, _value uint256)
func (_Token *TokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TokenApproval, _owner []common.Address, _spender []common.Address) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _Token.contract.WatchLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenApproval)
				if err := _Token.contract.UnpackLog(event, "Approval", log); err != nil {
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

// TokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Token contract.
type TokenTransferIterator struct {
	Event *TokenTransfer // Event containing the contract specifics and raw log

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
func (it *TokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenTransfer)
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
		it.Event = new(TokenTransfer)
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
func (it *TokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenTransfer represents a Transfer event raised by the Token contract.
type TokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(_from indexed address, _to indexed address, _value uint256)
func (_Token *TokenFilterer) FilterTransfer(opts *bind.FilterOpts, _from []common.Address, _to []common.Address) (*TokenTransferIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _Token.contract.FilterLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &TokenTransferIterator{contract: _Token.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(_from indexed address, _to indexed address, _value uint256)
func (_Token *TokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TokenTransfer, _from []common.Address, _to []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _Token.contract.WatchLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenTransfer)
				if err := _Token.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// UtilsABI is the input ABI used to generate the binding from.
const UtilsABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"channel\",\"type\":\"address\"}],\"name\":\"contractExists\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// UtilsBin is the compiled bytecode used for deploying new contracts.
const UtilsBin = `0x608060405234801561001057600080fd5b50610187806100206000396000f30060806040526004361061004b5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416637709bc788114610050578063b32c65c814610092575b600080fd5b34801561005c57600080fd5b5061007e73ffffffffffffffffffffffffffffffffffffffff6004351661011c565b604080519115158252519081900360200190f35b34801561009e57600080fd5b506100a7610124565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100e15781810151838201526020016100c9565b50505050905090810190601f16801561010e5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000903b1190565b60408051808201909152600581527f302e322e5f0000000000000000000000000000000000000000000000000000006020820152815600a165627a7a72305820d67873f03c27cd4bec55bf24ab7fa7bb85740108eea23308d3160fe2dc320b4f0029`

// DeployUtils deploys a new Ethereum contract, binding an instance of Utils to it.
func DeployUtils(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Utils, error) {
	parsed, err := abi.JSON(strings.NewReader(UtilsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(UtilsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Utils{UtilsCaller: UtilsCaller{contract: contract}, UtilsTransactor: UtilsTransactor{contract: contract}, UtilsFilterer: UtilsFilterer{contract: contract}}, nil
}

// Utils is an auto generated Go binding around an Ethereum contract.
type Utils struct {
	UtilsCaller     // Read-only binding to the contract
	UtilsTransactor // Write-only binding to the contract
	UtilsFilterer   // Log filterer for contract events
}

// UtilsCaller is an auto generated read-only Go binding around an Ethereum contract.
type UtilsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UtilsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UtilsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UtilsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UtilsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UtilsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UtilsSession struct {
	Contract     *Utils            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UtilsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UtilsCallerSession struct {
	Contract *UtilsCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// UtilsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UtilsTransactorSession struct {
	Contract     *UtilsTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UtilsRaw is an auto generated low-level Go binding around an Ethereum contract.
type UtilsRaw struct {
	Contract *Utils // Generic contract binding to access the raw methods on
}

// UtilsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UtilsCallerRaw struct {
	Contract *UtilsCaller // Generic read-only contract binding to access the raw methods on
}

// UtilsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UtilsTransactorRaw struct {
	Contract *UtilsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUtils creates a new instance of Utils, bound to a specific deployed contract.
func NewUtils(address common.Address, backend bind.ContractBackend) (*Utils, error) {
	contract, err := bindUtils(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Utils{UtilsCaller: UtilsCaller{contract: contract}, UtilsTransactor: UtilsTransactor{contract: contract}, UtilsFilterer: UtilsFilterer{contract: contract}}, nil
}

// NewUtilsCaller creates a new read-only instance of Utils, bound to a specific deployed contract.
func NewUtilsCaller(address common.Address, caller bind.ContractCaller) (*UtilsCaller, error) {
	contract, err := bindUtils(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UtilsCaller{contract: contract}, nil
}

// NewUtilsTransactor creates a new write-only instance of Utils, bound to a specific deployed contract.
func NewUtilsTransactor(address common.Address, transactor bind.ContractTransactor) (*UtilsTransactor, error) {
	contract, err := bindUtils(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UtilsTransactor{contract: contract}, nil
}

// NewUtilsFilterer creates a new log filterer instance of Utils, bound to a specific deployed contract.
func NewUtilsFilterer(address common.Address, filterer bind.ContractFilterer) (*UtilsFilterer, error) {
	contract, err := bindUtils(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UtilsFilterer{contract: contract}, nil
}

// bindUtils binds a generic wrapper to an already deployed contract.
func bindUtils(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UtilsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Utils *UtilsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Utils.Contract.UtilsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Utils *UtilsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Utils.Contract.UtilsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Utils *UtilsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Utils.Contract.UtilsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Utils *UtilsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Utils.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Utils *UtilsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Utils.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Utils *UtilsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Utils.Contract.contract.Transact(opts, method, params...)
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(channel address) constant returns(bool)
func (_Utils *UtilsCaller) ContractExists(opts *bind.CallOpts, channel common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Utils.contract.Call(opts, out, "contractExists", channel)
	return *ret0, err
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(channel address) constant returns(bool)
func (_Utils *UtilsSession) ContractExists(channel common.Address) (bool, error) {
	return _Utils.Contract.ContractExists(&_Utils.CallOpts, channel)
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(channel address) constant returns(bool)
func (_Utils *UtilsCallerSession) ContractExists(channel common.Address) (bool, error) {
	return _Utils.Contract.ContractExists(&_Utils.CallOpts, channel)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_Utils *UtilsCaller) Contract_version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Utils.contract.Call(opts, out, "contract_version")
	return *ret0, err
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_Utils *UtilsSession) Contract_version() (string, error) {
	return _Utils.Contract.Contract_version(&_Utils.CallOpts)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_Utils *UtilsCallerSession) Contract_version() (string, error) {
	return _Utils.Contract.Contract_version(&_Utils.CallOpts)
}
