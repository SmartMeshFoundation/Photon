// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rpc

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
const ChannelManagerContractABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"getChannelsParticipants\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"partner\",\"type\":\"address\"}],\"name\":\"getChannelWith\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getChannelsAddresses\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"node_address\",\"type\":\"address\"}],\"name\":\"nettingContractsByAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"channel\",\"type\":\"address\"}],\"name\":\"contractExists\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"partner\",\"type\":\"address\"},{\"name\":\"settle_timeout\",\"type\":\"uint256\"}],\"name\":\"newChannel\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"token_address\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"netting_channel\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"participant1\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"participant2\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"settle_timeout\",\"type\":\"uint256\"}],\"name\":\"ChannelNew\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"caller_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"partner\",\"type\":\"address\"}],\"name\":\"ChannelDeleted\",\"type\":\"event\"}]"

// ChannelManagerContractBin is the compiled bytecode used for deploying new contracts.
const ChannelManagerContractBin = `0x6060604052341561000f57600080fd5b6040516020806107998339810160405280805160008054600160a060020a03909216600160a060020a03199092169190911790555050610745806100546000396000f3006060604052600436106100745763ffffffff60e060020a6000350416630b74b6208114610084578063238bfba2146100ea5780636785b500146101255780636cb30fee146101385780637709bc78146101575780639d76ea581461018a578063b32c65c81461019d578063f26c6aed14610227575b341561007f57600080fd5b600080fd5b341561008f57600080fd5b610097610249565b60405160208082528190810183818151815260200191508051906020019060200280838360005b838110156100d65780820151838201526020016100be565b505050509050019250505060405180910390f35b34156100f557600080fd5b610109600160a060020a03600435166103fe565b604051600160a060020a03909116815260200160405180910390f35b341561013057600080fd5b61009761047a565b341561014357600080fd5b610097600160a060020a03600435166104e1565b341561016257600080fd5b610176600160a060020a0360043516610571565b604051901515815260200160405180910390f35b341561019557600080fd5b610109610579565b34156101a857600080fd5b6101b0610588565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156101ec5780820151838201526020016101d4565b50505050905090810190601f1680156102195780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561023257600080fd5b610109600160a060020a03600435166024356105bf565b610251610707565b60008061025c610707565b60009250828080805b6001548710156102b3576001805461029c91908990811061028257fe5b600091825260209091200154600160a060020a0316610571565b156102a8576001830192505b600190960195610265565b826002026040518059106102c45750595b9080825280602002602001820160405250945060009550600096505b6001548710156103f257600180546102fd91908990811061028257fe5b1515610308576103e7565b600180548890811061031657fe5b600091825260209091200154600160a060020a0316935083636d2381b36040518163ffffffff1660e060020a028152600401608060405180830381600087803b151561036157600080fd5b5af1151561036e57600080fd5b505050604051805190602001805190602001805190602001805190505092505091508185878151811061039d57fe5b600160a060020a0390921660209283029091019091015260019590950194808587815181106103c857fe5b600160a060020a03909216602092830290910190910152600195909501945b6001909601956102e0565b50929695505050505050565b600073__ChannelManagerLibrary.sol:ChannelMan__638a1c00e2828460405160e060020a63ffffffff85160281526004810192909252600160a060020a0316602482015260440160206040518083038186803b151561045e57600080fd5b5af4151561046b57600080fd5b50505060405180519392505050565b610482610707565b600180546020808202016040519081016040528092919081815260200182805480156104d757602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116104b9575b5050505050905090565b6104e9610707565b6000600301600083600160a060020a0316600160a060020a0316815260200190815260200160002080548060200260200160405190810160405280929190818152602001828054801561056557602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610547575b50505050509050919050565b6000903b1190565b600054600160a060020a031690565b60408051908101604052600581527f302e322e5f000000000000000000000000000000000000000000000000000000602082015281565b60008060006105cd856103fe565b9150600160a060020a03821615610625577fda8d2f351e0f7c8c368e631ce8ab15973e7582ece0c347d75a5cff49eb899eb73386604051600160a060020a039283168152911660208201526040908101905180910390a15b73__ChannelManagerLibrary.sol:ChannelMan__63941583a56000878760405160e060020a63ffffffff86160281526004810193909352600160a060020a039091166024830152604482015260640160206040518083038186803b151561068c57600080fd5b5af4151561069957600080fd5b5050506040518051905090507f7bd269696a33040df6c111efd58439c9c77909fcbe90f7511065ac277e175dac81338787604051600160a060020a039485168152928416602084015292166040808301919091526060820192909252608001905180910390a1949350505050565b602060405190810160405260008152905600a165627a7a723058202e1a52152dc54c0dcbb96c5a9883008ed9ac0740cbbc1cead3d316191a498d420029`

// DeployChannelManagerContract deploys a new Ethereum contract, binding an instance of ChannelManagerContract to it.
func DeployChannelManagerContract(auth *bind.TransactOpts, backend bind.ContractBackend, token_address common.Address) (common.Address, *types.Transaction, *ChannelManagerContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ChannelManagerContractABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ChannelManagerContractBin), backend, token_address)
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
	Caller_address common.Address
	Partner        common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterChannelDeleted is a free log retrieval operation binding the contract event 0xda8d2f351e0f7c8c368e631ce8ab15973e7582ece0c347d75a5cff49eb899eb7.
//
// Solidity: event ChannelDeleted(caller_address address, partner address)
func (_ChannelManagerContract *ChannelManagerContractFilterer) FilterChannelDeleted(opts *bind.FilterOpts) (*ChannelManagerContractChannelDeletedIterator, error) {

	logs, sub, err := _ChannelManagerContract.contract.FilterLogs(opts, "ChannelDeleted")
	if err != nil {
		return nil, err
	}
	return &ChannelManagerContractChannelDeletedIterator{contract: _ChannelManagerContract.contract, event: "ChannelDeleted", logs: logs, sub: sub}, nil
}

// WatchChannelDeleted is a free log subscription operation binding the contract event 0xda8d2f351e0f7c8c368e631ce8ab15973e7582ece0c347d75a5cff49eb899eb7.
//
// Solidity: event ChannelDeleted(caller_address address, partner address)
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
	Netting_channel common.Address
	Participant1    common.Address
	Participant2    common.Address
	Settle_timeout  *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterChannelNew is a free log retrieval operation binding the contract event 0x7bd269696a33040df6c111efd58439c9c77909fcbe90f7511065ac277e175dac.
//
// Solidity: event ChannelNew(netting_channel address, participant1 address, participant2 address, settle_timeout uint256)
func (_ChannelManagerContract *ChannelManagerContractFilterer) FilterChannelNew(opts *bind.FilterOpts) (*ChannelManagerContractChannelNewIterator, error) {

	logs, sub, err := _ChannelManagerContract.contract.FilterLogs(opts, "ChannelNew")
	if err != nil {
		return nil, err
	}
	return &ChannelManagerContractChannelNewIterator{contract: _ChannelManagerContract.contract, event: "ChannelNew", logs: logs, sub: sub}, nil
}

// WatchChannelNew is a free log subscription operation binding the contract event 0x7bd269696a33040df6c111efd58439c9c77909fcbe90f7511065ac277e175dac.
//
// Solidity: event ChannelNew(netting_channel address, participant1 address, participant2 address, settle_timeout uint256)
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
var ChannelManagerLibraryBin = `0x611407610030600b82828239805160001a6073146000811461002057610022565bfe5b5030600052607381538281f300730000000000000000000000000000000000000000301460606040526004361061006d5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416638a1c00e28114610072578063941583a5146100a5578063b32c65c8146100ca575b600080fd5b610089600435600160a060020a0360243516610149565b604051600160a060020a03909116815260200160405180910390f35b81156100b057600080fd5b610089600435600160a060020a03602435166044356101a8565b6100d261055b565b60405160208082528190810183818151815260200191508051906020019080838360005b8381101561010e5780820151838201526020016100f6565b50505050905090810190601f16801561013b5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60008060006101583385610592565b6000818152600287016020526040902054909250905080156101a057600185018054600019830190811061018857fe5b600091825260209091200154600160a060020a031692505b505092915050565b33600160a060020a0381811660009081526003860160205260408082209286168252812090928390819081908190819081906101e4908c610592565b600081815260028e0160205260409020548d549197509550600160a060020a0316338c8c610210610634565b600160a060020a0394851681529284166020840152921660408083019190915260608201929092526080019051809103906000f080151561025057600080fd5b935084156103a85760018c018054600019870190811061026c57fe5b600091825260209091200154600160a060020a0316925061028c8361062c565b1561029657600080fd5b5050600160a060020a03338116600081815260048d0160208181526040808420958f1684529481528483205491815284832093835292909252919091205460018c01805485919060001988019081106102eb57fe5b906000526020600020900160006101000a815481600160a060020a030219169083600160a060020a0316021790555083886001840381548110151561032c57fe5b906000526020600020900160006101000a815481600160a060020a030219169083600160a060020a0316021790555083876001830381548110151561036d57fe5b6000918252602090912001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a039290921691909117905561054b565b8b60010180548060010182816103be9190610644565b506000918252602090912001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03861617905587548890600181016104028382610644565b506000918252602090912001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03861617905586548790600181016104468382610644565b9160005260206000209001600086909190916101000a815481600160a060020a030219169083600160a060020a03160217905550508b600101805490508c600201600088600019166000191681526020019081526020016000208190555087805490508c600401600033600160a060020a0316600160a060020a0316815260200190815260200160002060008d600160a060020a0316600160a060020a031681526020019081526020016000208190555086805490508c60040160008d600160a060020a0316600160a060020a03168152602001908152602001600020600033600160a060020a0316600160a060020a03168152602001908152602001600020819055505b50919a9950505050505050505050565b60408051908101604052600581527f302e322e5f000000000000000000000000000000000000000000000000000000602082015281565b600081600160a060020a031683600160a060020a031610156105ec5782826040516c01000000000000000000000000600160a060020a03938416810282529190921602601482015260280160405180910390209050610626565b81836040516c01000000000000000000000000600160a060020a039384168102825291909216026014820152602801604051809103902090505b92915050565b6000903b1190565b604051610d4d8061068f83390190565b8154818355818111156106685760008381526020902061066891810190830161066d565b505050565b61068b91905b808211156106875760008155600101610673565b5090565b9056006060604052341561000f57600080fd5b604051608080610d4d833981016040528080519190602001805191906020018051919060200180519150819050600681108015906100505750622932e08111155b151561005b57600080fd5b600160a060020a03848116908416141561007457600080fd5b5060058054600160a060020a0319908116600160a060020a03958616908117909255600c805494861694821685179055600091825260136020526040808320805460ff19908116600190811790925595845290832080549095166002179094556004805496909516951694909417909255908255439055610c529081906100fb90396000f3006060604052600436106100b65763ffffffff60e060020a60003504166311da60b481146100c6578063202ac3bc146100db57806327d120fe1461017057806353af5d10146101de578063597e1fb51461020d5780635e1fc56e146102325780635f88eade146102a05780636d2381b3146102b357806373d4a13a146102fc5780637ebdc478146103515780639d76ea5814610364578063b32c65c814610377578063b6b55f2514610401578063ed531a4a1461042b575b34156100c157600080fd5b600080fd5b34156100d157600080fd5b6100d96104be565b005b34156100e657600080fd5b6100d960046024813581810190830135806020601f8201819004810201604051908101604052818152929190602084018383808284378201915050505050509190803590602001908201803590602001908080601f016020809104026020016040519081016040528181529291906020840183838082843750949650509335935061054c92505050565b341561017b57600080fd5b6100d96004803567ffffffffffffffff1690602480359160443591606435919060a49060843590810190830135806020601f820181900481020160405190810160405281815292919060208401838380828437509496506106d695505050505050565b34156101e957600080fd5b6101f1610815565b604051600160a060020a03909116815260200160405180910390f35b341561021857600080fd5b610220610824565b60405190815260200160405180910390f35b341561023d57600080fd5b6100d96004803567ffffffffffffffff1690602480359160443591606435919060a49060843590810190830135806020601f8201819004810201604051908101604052818152929190602084018383808284375094965061082a95505050505050565b34156102ab57600080fd5b610220610969565b34156102be57600080fd5b6102c661096f565b604051600160a060020a039485168152602081019390935292166040808301919091526060820192909252608001905180910390f35b341561030757600080fd5b61030f61098f565b6040519586526020860194909452604080860193909352600160a060020a03918216606086015216608084015290151560a083015260c0909101905180910390f35b341561035c57600080fd5b6102206109b7565b341561036f57600080fd5b6101f16109bd565b341561038257600080fd5b61038a6109cc565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156103c65780820151838201526020016103ae565b50505050905090810190601f1680156103f35780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561040c57600080fd5b610417600435610a03565b604051901515815260200160405180910390f35b341561043657600080fd5b6100d960046024813581810190830135806020601f8201819004810201604051908101604052818152929190602084018383808284378201915050505050509190803590602001908201803590602001908080601f016020809104026020016040519081016040528181529291906020840183838082843750949650610aec95505050505050565b73__NettingChannelLibrary.sol:NettingCha__63de394e0d600060405160e060020a63ffffffff8416028152600481019190915260240160006040518083038186803b151561050e57600080fd5b5af4151561051b57600080fd5b5050507f6713dea2491bc95585ea9be0d6993fc7790fdcd04f495a7e7592fbd80bbe00de60405160405180910390a1565b73__NettingChannelLibrary.sol:NettingCha__63c2522462600085858560405160e060020a63ffffffff871602815260048101858152606482018390526080602483019081529091604481019060840186818151815260200191508051906020019080838360005b838110156105ce5780820151838201526020016105b6565b50505050905090810190601f1680156105fb5780820380516001836020036101000a031916815260200191505b50838103825285818151815260200191508051906020019080838360005b83811015610631578082015183820152602001610619565b50505050905090810190601f16801561065e5780820380516001836020036101000a031916815260200191505b50965050505050505060006040518083038186803b151561067e57600080fd5b5af4151561068b57600080fd5b5050507fa2e2842eefea7e32abccccd9d3fae92608319362c3905ef73de44938c05925368133604051918252600160a060020a031660208201526040908101905180910390a1505050565b73__NettingChannelLibrary.sol:NettingCha__63f565eb366000878787878760405160e060020a63ffffffff89160281526004810187815267ffffffffffffffff8716602483015260448201869052606482018590526084820184905260c060a48301908152909160c40183818151815260200191508051906020019080838360005b8381101561077357808201518382015260200161075b565b50505050905090810190601f1680156107a05780820380516001836020036101000a031916815260200191505b5097505050505050505060006040518083038186803b15156107c157600080fd5b5af415156107ce57600080fd5b5050507fa0379b1bd0864245b4ff39bf6a023065e80d0e9276d2671d94b9f653b4bbcdfe33604051600160a060020a03909116815260200160405180910390a15050505050565b600354600160a060020a031690565b60025490565b73__NettingChannelLibrary.sol:NettingCha__63c800b0026000878787878760405160e060020a63ffffffff89160281526004810187815267ffffffffffffffff8716602483015260448201869052606482018590526084820184905260c060a48301908152909160c40183818151815260200191508051906020019080838360005b838110156108c75780820151838201526020016108af565b50505050905090810190601f1680156108f45780820380516001836020036101000a031916815260200191505b5097505050505050505060006040518083038186803b151561091557600080fd5b5af4151561092257600080fd5b5050507f93daadf23cc2150b386a6c3b39f6e61b9c555fc1cec423e4c93ac9d36b008fef33604051600160a060020a03909116815260200160405180910390a15050505050565b60015490565b600554600654600c54600d54600160a060020a0393841694929390911691565b600054600154600254600354600454601454600160a060020a03928316929091169060ff1686565b60005490565b600454600160a060020a031690565b60408051908101604052600581527f302e322e5f000000000000000000000000000000000000000000000000000000602082015281565b6000808073__NettingChannelLibrary.sol:NettingCha__633268a05a828660405160e060020a63ffffffff851602815260048101929092526024820152604401604080518083038186803b1515610a5b57600080fd5b5af41515610a6857600080fd5b50505060405180519060200180519193509091505060018215151415610ae5576004547f9cb02993ef7311b37acc6bdfc1a8397160be258a877d78b31f4e366caf7bfcef90600160a060020a03163383604051600160a060020a039384168152919092166020820152604080820192909252606001905180910390a15b5092915050565b73__NettingChannelLibrary.sol:NettingCha__632b6c209e600084846040518463ffffffff1660e060020a028152600401808481526020018060200180602001838103835285818151815260200191508051906020019080838360005b83811015610b63578082015183820152602001610b4b565b50505050905090810190601f168015610b905780820380516001836020036101000a031916815260200191505b50838103825284818151815260200191508051906020019080838360005b83811015610bc6578082015183820152602001610bae565b50505050905090810190601f168015610bf35780820380516001836020036101000a031916815260200191505b509550505050505060006040518083038186803b1515610c1257600080fd5b5af41515610c1f57600080fd5b50505050505600a165627a7a7230582076dd2c560f44bfad49231d36b062efa84be529abd9875348df388fe9ec19de0c0029a165627a7a72305820f52a41e895c958453724a3fcb9e25cd23612a7474fdb89e7a5a723e7afed15d20029`

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
const NettingChannelContractABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"settle\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"locked_encoded\",\"type\":\"bytes\"},{\"name\":\"merkle_proof\",\"type\":\"bytes\"},{\"name\":\"secret\",\"type\":\"bytes32\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"transferred_amount\",\"type\":\"uint256\"},{\"name\":\"locksroot\",\"type\":\"bytes32\"},{\"name\":\"extra_hash\",\"type\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"updateTransfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"closingAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"closed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"transferred_amount\",\"type\":\"uint256\"},{\"name\":\"locksroot\",\"type\":\"bytes32\"},{\"name\":\"extra_hash\",\"type\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"close\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"opened\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"addressAndBalance\",\"outputs\":[{\"name\":\"participant1\",\"type\":\"address\"},{\"name\":\"balance1\",\"type\":\"uint256\"},{\"name\":\"participant2\",\"type\":\"address\"},{\"name\":\"balance2\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"data\",\"outputs\":[{\"name\":\"settle_timeout\",\"type\":\"uint256\"},{\"name\":\"opened\",\"type\":\"uint256\"},{\"name\":\"closed\",\"type\":\"uint256\"},{\"name\":\"closing_address\",\"type\":\"address\"},{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"updated\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"settleTimeout\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"locked_encoded\",\"type\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"unwithdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"token_address\",\"type\":\"address\"},{\"name\":\"participant1\",\"type\":\"address\"},{\"name\":\"participant2\",\"type\":\"address\"},{\"name\":\"timeout\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"participant\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"balance\",\"type\":\"uint256\"}],\"name\":\"ChannelNewBalance\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"closing_address\",\"type\":\"address\"}],\"name\":\"ChannelClosed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"node_address\",\"type\":\"address\"}],\"name\":\"TransferUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"ChannelSettled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"secret\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"receiver_address\",\"type\":\"address\"}],\"name\":\"ChannelSecretRevealed\",\"type\":\"event\"}]"

// NettingChannelContractBin is the compiled bytecode used for deploying new contracts.
const NettingChannelContractBin = `0x6060604052341561000f57600080fd5b604051608080610d4d833981016040528080519190602001805191906020018051919060200180519150819050600681108015906100505750622932e08111155b151561005b57600080fd5b600160a060020a03848116908416141561007457600080fd5b5060058054600160a060020a0319908116600160a060020a03958616908117909255600c805494861694821685179055600091825260136020526040808320805460ff19908116600190811790925595845290832080549095166002179094556004805496909516951694909417909255908255439055610c529081906100fb90396000f3006060604052600436106100b65763ffffffff60e060020a60003504166311da60b481146100c6578063202ac3bc146100db57806327d120fe1461017057806353af5d10146101de578063597e1fb51461020d5780635e1fc56e146102325780635f88eade146102a05780636d2381b3146102b357806373d4a13a146102fc5780637ebdc478146103515780639d76ea5814610364578063b32c65c814610377578063b6b55f2514610401578063ed531a4a1461042b575b34156100c157600080fd5b600080fd5b34156100d157600080fd5b6100d96104be565b005b34156100e657600080fd5b6100d960046024813581810190830135806020601f8201819004810201604051908101604052818152929190602084018383808284378201915050505050509190803590602001908201803590602001908080601f016020809104026020016040519081016040528181529291906020840183838082843750949650509335935061054c92505050565b341561017b57600080fd5b6100d96004803567ffffffffffffffff1690602480359160443591606435919060a49060843590810190830135806020601f820181900481020160405190810160405281815292919060208401838380828437509496506106d695505050505050565b34156101e957600080fd5b6101f1610815565b604051600160a060020a03909116815260200160405180910390f35b341561021857600080fd5b610220610824565b60405190815260200160405180910390f35b341561023d57600080fd5b6100d96004803567ffffffffffffffff1690602480359160443591606435919060a49060843590810190830135806020601f8201819004810201604051908101604052818152929190602084018383808284375094965061082a95505050505050565b34156102ab57600080fd5b610220610969565b34156102be57600080fd5b6102c661096f565b604051600160a060020a039485168152602081019390935292166040808301919091526060820192909252608001905180910390f35b341561030757600080fd5b61030f61098f565b6040519586526020860194909452604080860193909352600160a060020a03918216606086015216608084015290151560a083015260c0909101905180910390f35b341561035c57600080fd5b6102206109b7565b341561036f57600080fd5b6101f16109bd565b341561038257600080fd5b61038a6109cc565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156103c65780820151838201526020016103ae565b50505050905090810190601f1680156103f35780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561040c57600080fd5b610417600435610a03565b604051901515815260200160405180910390f35b341561043657600080fd5b6100d960046024813581810190830135806020601f8201819004810201604051908101604052818152929190602084018383808284378201915050505050509190803590602001908201803590602001908080601f016020809104026020016040519081016040528181529291906020840183838082843750949650610aec95505050505050565b73__NettingChannelLibrary.sol:NettingCha__63de394e0d600060405160e060020a63ffffffff8416028152600481019190915260240160006040518083038186803b151561050e57600080fd5b5af4151561051b57600080fd5b5050507f6713dea2491bc95585ea9be0d6993fc7790fdcd04f495a7e7592fbd80bbe00de60405160405180910390a1565b73__NettingChannelLibrary.sol:NettingCha__63c2522462600085858560405160e060020a63ffffffff871602815260048101858152606482018390526080602483019081529091604481019060840186818151815260200191508051906020019080838360005b838110156105ce5780820151838201526020016105b6565b50505050905090810190601f1680156105fb5780820380516001836020036101000a031916815260200191505b50838103825285818151815260200191508051906020019080838360005b83811015610631578082015183820152602001610619565b50505050905090810190601f16801561065e5780820380516001836020036101000a031916815260200191505b50965050505050505060006040518083038186803b151561067e57600080fd5b5af4151561068b57600080fd5b5050507fa2e2842eefea7e32abccccd9d3fae92608319362c3905ef73de44938c05925368133604051918252600160a060020a031660208201526040908101905180910390a1505050565b73__NettingChannelLibrary.sol:NettingCha__63f565eb366000878787878760405160e060020a63ffffffff89160281526004810187815267ffffffffffffffff8716602483015260448201869052606482018590526084820184905260c060a48301908152909160c40183818151815260200191508051906020019080838360005b8381101561077357808201518382015260200161075b565b50505050905090810190601f1680156107a05780820380516001836020036101000a031916815260200191505b5097505050505050505060006040518083038186803b15156107c157600080fd5b5af415156107ce57600080fd5b5050507fa0379b1bd0864245b4ff39bf6a023065e80d0e9276d2671d94b9f653b4bbcdfe33604051600160a060020a03909116815260200160405180910390a15050505050565b600354600160a060020a031690565b60025490565b73__NettingChannelLibrary.sol:NettingCha__63c800b0026000878787878760405160e060020a63ffffffff89160281526004810187815267ffffffffffffffff8716602483015260448201869052606482018590526084820184905260c060a48301908152909160c40183818151815260200191508051906020019080838360005b838110156108c75780820151838201526020016108af565b50505050905090810190601f1680156108f45780820380516001836020036101000a031916815260200191505b5097505050505050505060006040518083038186803b151561091557600080fd5b5af4151561092257600080fd5b5050507f93daadf23cc2150b386a6c3b39f6e61b9c555fc1cec423e4c93ac9d36b008fef33604051600160a060020a03909116815260200160405180910390a15050505050565b60015490565b600554600654600c54600d54600160a060020a0393841694929390911691565b600054600154600254600354600454601454600160a060020a03928316929091169060ff1686565b60005490565b600454600160a060020a031690565b60408051908101604052600581527f302e322e5f000000000000000000000000000000000000000000000000000000602082015281565b6000808073__NettingChannelLibrary.sol:NettingCha__633268a05a828660405160e060020a63ffffffff851602815260048101929092526024820152604401604080518083038186803b1515610a5b57600080fd5b5af41515610a6857600080fd5b50505060405180519060200180519193509091505060018215151415610ae5576004547f9cb02993ef7311b37acc6bdfc1a8397160be258a877d78b31f4e366caf7bfcef90600160a060020a03163383604051600160a060020a039384168152919092166020820152604080820192909252606001905180910390a15b5092915050565b73__NettingChannelLibrary.sol:NettingCha__632b6c209e600084846040518463ffffffff1660e060020a028152600401808481526020018060200180602001838103835285818151815260200191508051906020019080838360005b83811015610b63578082015183820152602001610b4b565b50505050905090810190601f168015610b905780820380516001836020036101000a031916815260200191505b50838103825284818151815260200191508051906020019080838360005b83811015610bc6578082015183820152602001610bae565b50505050905090810190601f168015610bf35780820380516001836020036101000a031916815260200191505b509550505050505060006040518083038186803b1515610c1257600080fd5b5af41515610c1f57600080fd5b50505050505600a165627a7a7230582076dd2c560f44bfad49231d36b062efa84be529abd9875348df388fe9ec19de0c0029`

// DeployNettingChannelContract deploys a new Ethereum contract, binding an instance of NettingChannelContract to it.
func DeployNettingChannelContract(auth *bind.TransactOpts, backend bind.ContractBackend, token_address common.Address, participant1 common.Address, participant2 common.Address, timeout *big.Int) (common.Address, *types.Transaction, *NettingChannelContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NettingChannelContractABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(NettingChannelContractBin), backend, token_address, participant1, participant2, timeout)
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
// Solidity: function data() constant returns(settle_timeout uint256, opened uint256, closed uint256, closing_address address, token address, updated bool)
func (_NettingChannelContract *NettingChannelContractCaller) Data(opts *bind.CallOpts) (struct {
	Settle_timeout  *big.Int
	Opened          *big.Int
	Closed          *big.Int
	Closing_address common.Address
	Token           common.Address
	Updated         bool
}, error) {
	ret := new(struct {
		Settle_timeout  *big.Int
		Opened          *big.Int
		Closed          *big.Int
		Closing_address common.Address
		Token           common.Address
		Updated         bool
	})
	out := ret
	err := _NettingChannelContract.contract.Call(opts, out, "data")
	return *ret, err
}

// Data is a free data retrieval call binding the contract method 0x73d4a13a.
//
// Solidity: function data() constant returns(settle_timeout uint256, opened uint256, closed uint256, closing_address address, token address, updated bool)
func (_NettingChannelContract *NettingChannelContractSession) Data() (struct {
	Settle_timeout  *big.Int
	Opened          *big.Int
	Closed          *big.Int
	Closing_address common.Address
	Token           common.Address
	Updated         bool
}, error) {
	return _NettingChannelContract.Contract.Data(&_NettingChannelContract.CallOpts)
}

// Data is a free data retrieval call binding the contract method 0x73d4a13a.
//
// Solidity: function data() constant returns(settle_timeout uint256, opened uint256, closed uint256, closing_address address, token address, updated bool)
func (_NettingChannelContract *NettingChannelContractCallerSession) Data() (struct {
	Settle_timeout  *big.Int
	Opened          *big.Int
	Closed          *big.Int
	Closing_address common.Address
	Token           common.Address
	Updated         bool
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

// Unwithdraw is a paid mutator transaction binding the contract method 0xed531a4a.
//
// Solidity: function unwithdraw(locked_encoded bytes, signature bytes) returns()
func (_NettingChannelContract *NettingChannelContractTransactor) Unwithdraw(opts *bind.TransactOpts, locked_encoded []byte, signature []byte) (*types.Transaction, error) {
	return _NettingChannelContract.contract.Transact(opts, "unwithdraw", locked_encoded, signature)
}

// Unwithdraw is a paid mutator transaction binding the contract method 0xed531a4a.
//
// Solidity: function unwithdraw(locked_encoded bytes, signature bytes) returns()
func (_NettingChannelContract *NettingChannelContractSession) Unwithdraw(locked_encoded []byte, signature []byte) (*types.Transaction, error) {
	return _NettingChannelContract.Contract.Unwithdraw(&_NettingChannelContract.TransactOpts, locked_encoded, signature)
}

// Unwithdraw is a paid mutator transaction binding the contract method 0xed531a4a.
//
// Solidity: function unwithdraw(locked_encoded bytes, signature bytes) returns()
func (_NettingChannelContract *NettingChannelContractTransactorSession) Unwithdraw(locked_encoded []byte, signature []byte) (*types.Transaction, error) {
	return _NettingChannelContract.Contract.Unwithdraw(&_NettingChannelContract.TransactOpts, locked_encoded, signature)
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

// Withdraw is a paid mutator transaction binding the contract method 0x202ac3bc.
//
// Solidity: function withdraw(locked_encoded bytes, merkle_proof bytes, secret bytes32) returns()
func (_NettingChannelContract *NettingChannelContractTransactor) Withdraw(opts *bind.TransactOpts, locked_encoded []byte, merkle_proof []byte, secret [32]byte) (*types.Transaction, error) {
	return _NettingChannelContract.contract.Transact(opts, "withdraw", locked_encoded, merkle_proof, secret)
}

// Withdraw is a paid mutator transaction binding the contract method 0x202ac3bc.
//
// Solidity: function withdraw(locked_encoded bytes, merkle_proof bytes, secret bytes32) returns()
func (_NettingChannelContract *NettingChannelContractSession) Withdraw(locked_encoded []byte, merkle_proof []byte, secret [32]byte) (*types.Transaction, error) {
	return _NettingChannelContract.Contract.Withdraw(&_NettingChannelContract.TransactOpts, locked_encoded, merkle_proof, secret)
}

// Withdraw is a paid mutator transaction binding the contract method 0x202ac3bc.
//
// Solidity: function withdraw(locked_encoded bytes, merkle_proof bytes, secret bytes32) returns()
func (_NettingChannelContract *NettingChannelContractTransactorSession) Withdraw(locked_encoded []byte, merkle_proof []byte, secret [32]byte) (*types.Transaction, error) {
	return _NettingChannelContract.Contract.Withdraw(&_NettingChannelContract.TransactOpts, locked_encoded, merkle_proof, secret)
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
	Closing_address common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterChannelClosed is a free log retrieval operation binding the contract event 0x93daadf23cc2150b386a6c3b39f6e61b9c555fc1cec423e4c93ac9d36b008fef.
//
// Solidity: event ChannelClosed(closing_address address)
func (_NettingChannelContract *NettingChannelContractFilterer) FilterChannelClosed(opts *bind.FilterOpts) (*NettingChannelContractChannelClosedIterator, error) {

	logs, sub, err := _NettingChannelContract.contract.FilterLogs(opts, "ChannelClosed")
	if err != nil {
		return nil, err
	}
	return &NettingChannelContractChannelClosedIterator{contract: _NettingChannelContract.contract, event: "ChannelClosed", logs: logs, sub: sub}, nil
}

// WatchChannelClosed is a free log subscription operation binding the contract event 0x93daadf23cc2150b386a6c3b39f6e61b9c555fc1cec423e4c93ac9d36b008fef.
//
// Solidity: event ChannelClosed(closing_address address)
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
	Token_address common.Address
	Participant   common.Address
	Balance       *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterChannelNewBalance is a free log retrieval operation binding the contract event 0x9cb02993ef7311b37acc6bdfc1a8397160be258a877d78b31f4e366caf7bfcef.
//
// Solidity: event ChannelNewBalance(token_address address, participant address, balance uint256)
func (_NettingChannelContract *NettingChannelContractFilterer) FilterChannelNewBalance(opts *bind.FilterOpts) (*NettingChannelContractChannelNewBalanceIterator, error) {

	logs, sub, err := _NettingChannelContract.contract.FilterLogs(opts, "ChannelNewBalance")
	if err != nil {
		return nil, err
	}
	return &NettingChannelContractChannelNewBalanceIterator{contract: _NettingChannelContract.contract, event: "ChannelNewBalance", logs: logs, sub: sub}, nil
}

// WatchChannelNewBalance is a free log subscription operation binding the contract event 0x9cb02993ef7311b37acc6bdfc1a8397160be258a877d78b31f4e366caf7bfcef.
//
// Solidity: event ChannelNewBalance(token_address address, participant address, balance uint256)
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
	Secret           [32]byte
	Receiver_address common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterChannelSecretRevealed is a free log retrieval operation binding the contract event 0xa2e2842eefea7e32abccccd9d3fae92608319362c3905ef73de44938c0592536.
//
// Solidity: event ChannelSecretRevealed(secret bytes32, receiver_address address)
func (_NettingChannelContract *NettingChannelContractFilterer) FilterChannelSecretRevealed(opts *bind.FilterOpts) (*NettingChannelContractChannelSecretRevealedIterator, error) {

	logs, sub, err := _NettingChannelContract.contract.FilterLogs(opts, "ChannelSecretRevealed")
	if err != nil {
		return nil, err
	}
	return &NettingChannelContractChannelSecretRevealedIterator{contract: _NettingChannelContract.contract, event: "ChannelSecretRevealed", logs: logs, sub: sub}, nil
}

// WatchChannelSecretRevealed is a free log subscription operation binding the contract event 0xa2e2842eefea7e32abccccd9d3fae92608319362c3905ef73de44938c0592536.
//
// Solidity: event ChannelSecretRevealed(secret bytes32, receiver_address address)
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
	Raw types.Log // Blockchain specific contextual infos
}

// FilterChannelSettled is a free log retrieval operation binding the contract event 0x6713dea2491bc95585ea9be0d6993fc7790fdcd04f495a7e7592fbd80bbe00de.
//
// Solidity: event ChannelSettled()
func (_NettingChannelContract *NettingChannelContractFilterer) FilterChannelSettled(opts *bind.FilterOpts) (*NettingChannelContractChannelSettledIterator, error) {

	logs, sub, err := _NettingChannelContract.contract.FilterLogs(opts, "ChannelSettled")
	if err != nil {
		return nil, err
	}
	return &NettingChannelContractChannelSettledIterator{contract: _NettingChannelContract.contract, event: "ChannelSettled", logs: logs, sub: sub}, nil
}

// WatchChannelSettled is a free log subscription operation binding the contract event 0x6713dea2491bc95585ea9be0d6993fc7790fdcd04f495a7e7592fbd80bbe00de.
//
// Solidity: event ChannelSettled()
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
const NettingChannelLibraryABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"self\",\"type\":\"NettingChannelLibrary.Data storage\"},{\"name\":\"locked_encoded\",\"type\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"unwithdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"self\",\"type\":\"NettingChannelLibrary.Data storage\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"},{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"self\",\"type\":\"NettingChannelLibrary.Data storage\"},{\"name\":\"locked_encoded\",\"type\":\"bytes\"},{\"name\":\"merkle_proof\",\"type\":\"bytes\"},{\"name\":\"secret\",\"type\":\"bytes32\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"self\",\"type\":\"NettingChannelLibrary.Data storage\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"transferred_amount\",\"type\":\"uint256\"},{\"name\":\"locksroot\",\"type\":\"bytes32\"},{\"name\":\"extra_hash\",\"type\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"close\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"self\",\"type\":\"NettingChannelLibrary.Data storage\"}],\"name\":\"settle\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"self\",\"type\":\"NettingChannelLibrary.Data storage\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"transferred_amount\",\"type\":\"uint256\"},{\"name\":\"locksroot\",\"type\":\"bytes32\"},{\"name\":\"extra_hash\",\"type\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"updateTransfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// NettingChannelLibraryBin is the compiled bytecode used for deploying new contracts.
const NettingChannelLibraryBin = `0x610f0b610030600b82828239805160001a6073146000811461002057610022565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146060604052600436106100805763ffffffff60e060020a6000350416632b6c209e81146100855780633268a05a1461011f578063b32c65c814610152578063c2522462146101d1578063c800b0021461026b578063de394e0d146102dc578063f565eb36146102f2575b600080fd5b811561009057600080fd5b61011d600480359060446024803590810190830135806020601f8201819004810201604051908101604052818152929190602084018383808284378201915050505050509190803590602001908201803590602001908080601f01602080910402602001604051908101604052818152929190602084018383808284375094965061036395505050505050565b005b811561012a57600080fd5b610138600435602435610476565b604051911515825260208201526040908101905180910390f35b61015a6105f0565b60405160208082528190810183818151815260200191508051906020019080838360005b8381101561019657808201518382015260200161017e565b50505050905090810190601f1680156101c35780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b81156101dc57600080fd5b61011d600480359060446024803590810190830135806020601f8201819004810201604051908101604052818152929190602084018383808284378201915050505050509190803590602001908201803590602001908080601f016020809104026020016040519081016040528181529291906020840183838082843750949650509335935061062792505050565b811561027657600080fd5b61011d60048035906024803567ffffffffffffffff169160443591606435916084359160c49060a43590810190830135806020601f8201819004810201604051908101604052818152929190602084018383808284375094965061073995505050505050565b81156102e757600080fd5b61011d600435610815565b81156102fd57600080fd5b61011d60048035906024803567ffffffffffffffff169160443591606435916084359160c49060a43590810190830135806020601f820181900481020160405190810160405281815292919060208401838380828437509496506109f295505050505050565b600080600080600080886000816002015411151561038057600080fd5b61038a8a33610b32565b600103955060058a0160ff8716600281106103a157fe5b6007020191508160020154600019166000600102141515156103c257600080fd5b6103cb89610b64565b9098509095509350600087116103e057600080fd5b600084815260058301602052604090205460ff16151560011461040257600080fd5b600084815260068301602052604090205460ff161561042057600080fd5b60008481526006830160205260409020805460ff191660011790556104458989610b92565b8254909350600160a060020a0380851691161461046157600080fd5b50600301805490950190945550505050505050565b6000806000806000866001015411151561048f57600080fd5b60028601541561049e57600080fd5b60048601548590600160a060020a03166370a082313360405160e060020a63ffffffff8416028152600160a060020a039091166004820152602401602060405180830381600087803b15156104f257600080fd5b5af115156104ff57600080fd5b505050604051805190501015151561051657600080fd5b6105208633610b32565b91506005860160ff83166002811061053457fe5b6004880154600791909102919091019150600160a060020a03166323b872dd33308860405160e060020a63ffffffff8616028152600160a060020a0393841660048201529190921660248201526044810191909152606401602060405180830381600087803b15156105a557600080fd5b5af115156105b257600080fd5b5050506040518051945050600184151514156105df576001808201805487019081905590945092506105e7565b600093508392505b50509250929050565b60408051908101604052600581527f302e322e5f000000000000000000000000000000000000000000000000000000602082015281565b600080600080600080896000816002015411151561064457600080fd5b61064e8b33610b32565b600103955060058b0160ff87166002811061066557fe5b60070201915081600201546000191660006001021415151561068657600080fd5b61068f8a610b64565b600081815260058601602052604090205491995091965090935060ff16156106b657600080fd5b60008381526005830160205260409020805460ff191660011790554367ffffffffffffffff861610156106e857600080fd5b87604051908152602001604051908190039020831461070657600080fd5b6107108a8a610c81565b6002830154909450841461072357600080fd5b5060030180549095019094555050505050505050565b6000806000808960020154600014151561075257600080fd5b4360028b01556107628a33610b32565b60038b01805473ffffffffffffffffffffffffffffffffffffffff191633600160a060020a031617905560ff169250845160411415610809576107a88989898989610d68565b93506107b48a85610b32565b60ff169150828214156107c657600080fd5b60058a0182600281106107d557fe5b6007020160048101805467ffffffffffffffff191667ffffffffffffffff8c16179055600281018890556003810189905590505b50505050505050505050565b600080600080600080600080886000816002015411151561083557600080fd5b895460028b01548b91439101111561084c57600080fd5b60038b0154610865908c90600160a060020a0316610b32565b995060018a9003985060058b0160ff8b166002811061088057fe5b60070201935060058b0160ff8a166002811061089857fe5b60070201925082600301548460030154846001015401039650826001015484600101540197506108c88789610e6d565b94506108d5856000610e85565b9450848803955060008511156109655760048b01548354600160a060020a039182169163a9059cbb91168760405160e060020a63ffffffff8516028152600160a060020a0390921660048301526024820152604401602060405180830381600087803b151561094357600080fd5b5af1151561095057600080fd5b50505060405180519050151561096557600080fd5b60008611156109ee5760048b01548454600160a060020a039182169163a9059cbb91168860405160e060020a63ffffffff8516028152600160a060020a0390921660048301526024820152604401602060405180830381600087803b15156109cc57600080fd5b5af115156109d957600080fd5b5050506040518051905015156109ee57600080fd5b6000ff5b60008060008860008160020154111515610a0b57600080fd5b895460028b01548b914391011015610a2257600080fd5b60148b015460ff1615610a3457600080fd5b60148b01805460ff19166001179055610a4d8b33610b32565b60038c015490945033600160a060020a0390811691161415610a6e57600080fd5b610a7b8a8a8a8a8a610d68565b60038c0154909550600160a060020a03808716911614610a9a57600080fd5b600184900392508960058c0160ff851660028110610ab457fe5b6007020160040160006101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550878b6005018460ff16600281101515610af657fe5b600702016002018160001916905550888b6005018460ff16600281101515610b1a57fe5b60070201600301819055505050505050505050505050565b600160a060020a038116600090815260138301602052604081205460ff16801515610b5957fe5b600019019392505050565b60008060008351604814610b7757600080fd5b60088401519250602884015191506048840151929491935050565b60008060008060008551604114610ba857600080fd5b866040518082805190602001908083835b60208310610bd85780518252601f199092019160209182019101610bb9565b6001836020036101000a038019825116818451161790925250505091909101925060409150505180910390209350610c0f86610e9b565b9250925092506001848285856040516000815260200160405260405193845260ff9092166020808501919091526040808501929092526060840192909252608090920191516020810390808403906000865af11515610c6d57600080fd5b505060206040510351979650505050505050565b60008060008060208551811515610c9457fe5b0615610c9f57600080fd5b856040518082805190602001908083835b60208310610ccf5780518252601f199092019160209182019101610cb0565b6001836020036101000a038019825116818451161790925250505091909101925060409150505180910390209150602092505b84518311610d5f5782850151905080821015610d3857818160405191825260208201526040908101905180910390209150610d54565b8082604051918252602082015260409081019051809103902091505b602083019250610d02565b50949350505050565b60008060008060008551604114610d7e57600080fd5b898989308a60405167ffffffffffffffff95909516780100000000000000000000000000000000000000000000000002855260088501939093526028840191909152600160a060020a03166c01000000000000000000000000026048830152605c820152607c0160405180910390209350610df886610e9b565b9250925092506001848285856040516000815260200160405260405193845260ff9092166020808501919091526040808501929092526060840192909252608090920191516020810390808403906000865af11515610e5657600080fd5b5050602060405103519a9950505050505050505050565b6000818311610e7c5782610e7e565b815b9392505050565b6000818311610e945781610e7e565b5090919050565b6000806000602084015192506040840151915060ff60418501511690508060ff16601b1480610ecd57508060ff16601c145b1515610ed857600080fd5b91939092505600a165627a7a7230582016ef7d9e08bca5c566d392d07d748e024508ddc7ec90472517c195940f56d8c40029`

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
const RegistryABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"registry\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"token_address\",\"type\":\"address\"}],\"name\":\"channelManagerByToken\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokens\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"channelManagerAddresses\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenAddresses\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token_address\",\"type\":\"address\"}],\"name\":\"addToken\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"channel_manager_address\",\"type\":\"address\"}],\"name\":\"TokenAdded\",\"type\":\"event\"}]"

// RegistryBin is the compiled bytecode used for deploying new contracts.
var RegistryBin = `0x6060604052341561000f57600080fd5b610dce8061001e6000396000f3006060604052600436106100825763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663038defd7811461009257806325119b5f146100cd5780634f64b2be146100ec578063640191e214610102578063a9989b9314610168578063b32c65c81461017b578063d48bfca714610205575b341561008d57600080fd5b600080fd5b341561009d57600080fd5b6100b1600160a060020a0360043516610224565b604051600160a060020a03909116815260200160405180910390f35b34156100d857600080fd5b6100b1600160a060020a036004351661023f565b34156100f757600080fd5b6100b1600435610289565b341561010d57600080fd5b6101156102b1565b60405160208082528190810183818151815260200191508051906020019060200280838360005b8381101561015457808201518382015260200161013c565b505050509050019250505060405180910390f35b341561017357600080fd5b610115610363565b341561018657600080fd5b61018e6103cc565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156101ca5780820151838201526020016101b2565b50505050905090810190601f1680156101f75780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561021057600080fd5b6100b1600160a060020a0360043516610403565b600060208190529081526040902054600160a060020a031681565b600160a060020a038082166000908152602081905260408120549091839116151561026957600080fd5b5050600160a060020a039081166000908152602081905260409020541690565b600180548290811061029757fe5b600091825260209091200154600160a060020a0316905081565b6102b96105a0565b6000806102c46105a0565b6001546040518059106102d45750595b90808252806020026020018201604052509050600092505b60015483101561035c57600180548490811061030457fe5b6000918252602080832090910154600160a060020a03908116808452918390526040909220549093501681848151811061033a57fe5b600160a060020a039092166020928302909101909101526001909201916102ec565b9392505050565b61036b6105a0565b60018054806020026020016040519081016040528092919081815260200182805480156103c157602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116103a3575b505050505090505b90565b60408051908101604052600581527f302e322e5f000000000000000000000000000000000000000000000000000000602082015281565b600160a060020a038082166000908152602081905260408120549091829184918391161561043057600080fd5b5080600160a060020a0381166318160ddd6040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401602060405180830381600087803b151561048857600080fd5b5af1151561049557600080fd5b5050506040518051905050846104a96105b2565b600160a060020a039091168152602001604051809103906000f08015156104cf57600080fd5b600160a060020a038681166000908152602081905260409020805473ffffffffffffffffffffffffffffffffffffffff1916918316919091179055600180549194509080820161051f83826105c2565b506000918252602090912001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0387161790557fdffbd9ded1c09446f09377de547142dcce7dc541c8b0b028142b1eba7026b9e78584604051600160a060020a039283168152911660208201526040908101905180910390a150909392505050565b60206040519081016040526000815290565b6040516107998061060a83390190565b8154818355818111156105e6576000838152602090206105e69181019083016105eb565b505050565b6103c991905b8082111561060557600081556001016105f1565b509056006060604052341561000f57600080fd5b6040516020806107998339810160405280805160008054600160a060020a03909216600160a060020a03199092169190911790555050610745806100546000396000f3006060604052600436106100745763ffffffff60e060020a6000350416630b74b6208114610084578063238bfba2146100ea5780636785b500146101255780636cb30fee146101385780637709bc78146101575780639d76ea581461018a578063b32c65c81461019d578063f26c6aed14610227575b341561007f57600080fd5b600080fd5b341561008f57600080fd5b610097610249565b60405160208082528190810183818151815260200191508051906020019060200280838360005b838110156100d65780820151838201526020016100be565b505050509050019250505060405180910390f35b34156100f557600080fd5b610109600160a060020a03600435166103fe565b604051600160a060020a03909116815260200160405180910390f35b341561013057600080fd5b61009761047a565b341561014357600080fd5b610097600160a060020a03600435166104e1565b341561016257600080fd5b610176600160a060020a0360043516610571565b604051901515815260200160405180910390f35b341561019557600080fd5b610109610579565b34156101a857600080fd5b6101b0610588565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156101ec5780820151838201526020016101d4565b50505050905090810190601f1680156102195780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561023257600080fd5b610109600160a060020a03600435166024356105bf565b610251610707565b60008061025c610707565b60009250828080805b6001548710156102b3576001805461029c91908990811061028257fe5b600091825260209091200154600160a060020a0316610571565b156102a8576001830192505b600190960195610265565b826002026040518059106102c45750595b9080825280602002602001820160405250945060009550600096505b6001548710156103f257600180546102fd91908990811061028257fe5b1515610308576103e7565b600180548890811061031657fe5b600091825260209091200154600160a060020a0316935083636d2381b36040518163ffffffff1660e060020a028152600401608060405180830381600087803b151561036157600080fd5b5af1151561036e57600080fd5b505050604051805190602001805190602001805190602001805190505092505091508185878151811061039d57fe5b600160a060020a0390921660209283029091019091015260019590950194808587815181106103c857fe5b600160a060020a03909216602092830290910190910152600195909501945b6001909601956102e0565b50929695505050505050565b600073__ChannelManagerLibrary.sol:ChannelMan__638a1c00e2828460405160e060020a63ffffffff85160281526004810192909252600160a060020a0316602482015260440160206040518083038186803b151561045e57600080fd5b5af4151561046b57600080fd5b50505060405180519392505050565b610482610707565b600180546020808202016040519081016040528092919081815260200182805480156104d757602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116104b9575b5050505050905090565b6104e9610707565b6000600301600083600160a060020a0316600160a060020a0316815260200190815260200160002080548060200260200160405190810160405280929190818152602001828054801561056557602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610547575b50505050509050919050565b6000903b1190565b600054600160a060020a031690565b60408051908101604052600581527f302e322e5f000000000000000000000000000000000000000000000000000000602082015281565b60008060006105cd856103fe565b9150600160a060020a03821615610625577fda8d2f351e0f7c8c368e631ce8ab15973e7582ece0c347d75a5cff49eb899eb73386604051600160a060020a039283168152911660208201526040908101905180910390a15b73__ChannelManagerLibrary.sol:ChannelMan__63941583a56000878760405160e060020a63ffffffff86160281526004810193909352600160a060020a039091166024830152604482015260640160206040518083038186803b151561068c57600080fd5b5af4151561069957600080fd5b5050506040518051905090507f7bd269696a33040df6c111efd58439c9c77909fcbe90f7511065ac277e175dac81338787604051600160a060020a039485168152928416602084015292166040808301919091526060820192909252608001905180910390a1949350505050565b602060405190810160405260008152905600a165627a7a723058202e1a52152dc54c0dcbb96c5a9883008ed9ac0740cbbc1cead3d316191a498d420029a165627a7a723058208ec2b51612a1beb890554b544f4631855d94f619e6c33bef69b957d9a5b34c9b0029`

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
	Token_address           common.Address
	Channel_manager_address common.Address
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterTokenAdded is a free log retrieval operation binding the contract event 0xdffbd9ded1c09446f09377de547142dcce7dc541c8b0b028142b1eba7026b9e7.
//
// Solidity: event TokenAdded(token_address address, channel_manager_address address)
func (_Registry *RegistryFilterer) FilterTokenAdded(opts *bind.FilterOpts) (*RegistryTokenAddedIterator, error) {

	logs, sub, err := _Registry.contract.FilterLogs(opts, "TokenAdded")
	if err != nil {
		return nil, err
	}
	return &RegistryTokenAddedIterator{contract: _Registry.contract, event: "TokenAdded", logs: logs, sub: sub}, nil
}

// WatchTokenAdded is a free log subscription operation binding the contract event 0xdffbd9ded1c09446f09377de547142dcce7dc541c8b0b028142b1eba7026b9e7.
//
// Solidity: event TokenAdded(token_address address, channel_manager_address address)
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

// UtilsABI is the input ABI used to generate the binding from.
const UtilsABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"channel\",\"type\":\"address\"}],\"name\":\"contractExists\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// UtilsBin is the compiled bytecode used for deploying new contracts.
const UtilsBin = `0x6060604052341561000f57600080fd5b6101858061001e6000396000f30060606040526004361061004b5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416637709bc788114610050578063b32c65c814610090575b600080fd5b341561005b57600080fd5b61007c73ffffffffffffffffffffffffffffffffffffffff6004351661011a565b604051901515815260200160405180910390f35b341561009b57600080fd5b6100a3610122565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156100df5780820151838201526020016100c7565b50505050905090810190601f16801561010c5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000903b1190565b60408051908101604052600581527f302e322e5f0000000000000000000000000000000000000000000000000000006020820152815600a165627a7a72305820f7e2e40a999451e6840a61d48e3699828689744ca0113d82733de5a755f9ad410029`

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
