// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rpc

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ChannelManagerContractABI is the input ABI used to generate the binding from.
const ChannelManagerContractABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"getChannelsParticipants\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"partner\",\"type\":\"address\"}],\"name\":\"getChannelWith\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getChannelsAddresses\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"node_address\",\"type\":\"address\"}],\"name\":\"nettingContractsByAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"channel\",\"type\":\"address\"}],\"name\":\"contractExists\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"partner\",\"type\":\"address\"},{\"name\":\"settle_timeout\",\"type\":\"uint256\"}],\"name\":\"newChannel\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"token_address\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"netting_channel\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"participant1\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"participant2\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"settle_timeout\",\"type\":\"uint256\"}],\"name\":\"ChannelNew\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"caller_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"partner\",\"type\":\"address\"}],\"name\":\"ChannelDeleted\",\"type\":\"event\"}]"

// ChannelManagerContractBin is the compiled bytecode used for deploying new contracts.
const ChannelManagerContractBin = `0x6060604052341561000f57600080fd5b6040516020806107bd8339810160405280805160008054600160a060020a03909216600160a060020a03199092169190911790555050610769806100546000396000f3006060604052600436106100745763ffffffff60e060020a6000350416630b74b6208114610084578063238bfba2146100ea5780636785b500146101255780636cb30fee146101385780637709bc78146101575780639d76ea581461018a578063b32c65c81461019d578063f26c6aed14610227575b341561007f57600080fd5b600080fd5b341561008f57600080fd5b610097610249565b60405160208082528190810183818151815260200191508051906020019060200280838360005b838110156100d65780820151838201526020016100be565b505050509050019250505060405180910390f35b34156100f557600080fd5b610109600160a060020a036004351661040a565b604051600160a060020a03909116815260200160405180910390f35b341561013057600080fd5b610097610492565b341561014357600080fd5b610097600160a060020a03600435166104f9565b341561016257600080fd5b610176600160a060020a0360043516610589565b604051901515815260200160405180910390f35b341561019557600080fd5b610109610591565b34156101a857600080fd5b6101b06105a0565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156101ec5780820151838201526020016101d4565b50505050905090810190601f1680156102195780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561023257600080fd5b610109600160a060020a03600435166024356105d7565b61025161072b565b60008061025c61072b565b60009250828080805b6001548710156102b3576001805461029c91908990811061028257fe5b600091825260209091200154600160a060020a0316610589565b156102a8576001830192505b600190960195610265565b826002026040518059106102c45750595b9080825280602002602001820160405250945060009550600096505b6001548710156103fe57600180546102fd91908990811061028257fe5b1515610308576103f3565b600180548890811061031657fe5b6000918252602082200154600160a060020a031694508490636d2381b390604051608001526040518163ffffffff1660e060020a028152600401608060405180830381600087803b151561036957600080fd5b6102c65a03f1151561037a57600080fd5b50505060405180519060200180519060200180519060200180519050509250509150818587815181106103a957fe5b600160a060020a0390921660209283029091019091015260019590950194808587815181106103d457fe5b600160a060020a03909216602092830290910190910152600195909501945b6001909601956102e0565b50929695505050505050565b600073__ChannelManagerLibrary.sol:ChannelMan__638a1c00e28284816040516020015260405160e060020a63ffffffff85160281526004810192909252600160a060020a0316602482015260440160206040518083038186803b151561047257600080fd5b6102c65a03f4151561048357600080fd5b50505060405180519392505050565b61049a61072b565b600180546020808202016040519081016040528092919081815260200182805480156104ef57602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116104d1575b5050505050905090565b61050161072b565b6000600301600083600160a060020a0316600160a060020a0316815260200190815260200160002080548060200260200160405190810160405280929190818152602001828054801561057d57602002820191906000526020600020905b8154600160a060020a0316815260019091019060200180831161055f575b50505050509050919050565b6000903b1190565b600054600160a060020a031690565b60408051908101604052600581527f302e322e5f000000000000000000000000000000000000000000000000000000602082015281565b60008060006105e58561040a565b9150600160a060020a0382161561063d577fda8d2f351e0f7c8c368e631ce8ab15973e7582ece0c347d75a5cff49eb899eb73386604051600160a060020a039283168152911660208201526040908101905180910390a15b73__ChannelManagerLibrary.sol:ChannelMan__63941583a560008787826040516020015260405160e060020a63ffffffff86160281526004810193909352600160a060020a039091166024830152604482015260640160206040518083038186803b15156106ac57600080fd5b6102c65a03f415156106bd57600080fd5b5050506040518051905090507f7bd269696a33040df6c111efd58439c9c77909fcbe90f7511065ac277e175dac81338787604051600160a060020a039485168152928416602084015292166040808301919091526060820192909252608001905180910390a1949350505050565b602060405190810160405260008152905600a165627a7a72305820522416bac5fb22953a513c007ec3c6a51edd7d61b9bd1cc408c4e121fd40c6db0029`

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
	return address, tx, &ChannelManagerContract{ChannelManagerContractCaller: ChannelManagerContractCaller{contract: contract}, ChannelManagerContractTransactor: ChannelManagerContractTransactor{contract: contract}}, nil
}

// ChannelManagerContract is an auto generated Go binding around an Ethereum contract.
type ChannelManagerContract struct {
	ChannelManagerContractCaller     // Read-only binding to the contract
	ChannelManagerContractTransactor // Write-only binding to the contract
}

// ChannelManagerContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChannelManagerContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChannelManagerContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChannelManagerContractTransactor struct {
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
	contract, err := bindChannelManagerContract(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChannelManagerContract{ChannelManagerContractCaller: ChannelManagerContractCaller{contract: contract}, ChannelManagerContractTransactor: ChannelManagerContractTransactor{contract: contract}}, nil
}

// NewChannelManagerContractCaller creates a new read-only instance of ChannelManagerContract, bound to a specific deployed contract.
func NewChannelManagerContractCaller(address common.Address, caller bind.ContractCaller) (*ChannelManagerContractCaller, error) {
	contract, err := bindChannelManagerContract(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &ChannelManagerContractCaller{contract: contract}, nil
}

// NewChannelManagerContractTransactor creates a new write-only instance of ChannelManagerContract, bound to a specific deployed contract.
func NewChannelManagerContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ChannelManagerContractTransactor, error) {
	contract, err := bindChannelManagerContract(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &ChannelManagerContractTransactor{contract: contract}, nil
}

// bindChannelManagerContract binds a generic wrapper to an already deployed contract.
func bindChannelManagerContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ChannelManagerContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
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

// NettingChannelContractABI is the input ABI used to generate the binding from.
const NettingChannelContractABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"settle\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"locked_encoded\",\"type\":\"bytes\"},{\"name\":\"merkle_proof\",\"type\":\"bytes\"},{\"name\":\"secret\",\"type\":\"bytes32\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"transferred_amount\",\"type\":\"uint256\"},{\"name\":\"locksroot\",\"type\":\"bytes32\"},{\"name\":\"extra_hash\",\"type\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"updateTransfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"closingAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"closed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"transferred_amount\",\"type\":\"uint256\"},{\"name\":\"locksroot\",\"type\":\"bytes32\"},{\"name\":\"extra_hash\",\"type\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"close\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"opened\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"addressAndBalance\",\"outputs\":[{\"name\":\"participant1\",\"type\":\"address\"},{\"name\":\"balance1\",\"type\":\"uint256\"},{\"name\":\"participant2\",\"type\":\"address\"},{\"name\":\"balance2\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"data\",\"outputs\":[{\"name\":\"settle_timeout\",\"type\":\"uint256\"},{\"name\":\"opened\",\"type\":\"uint256\"},{\"name\":\"closed\",\"type\":\"uint256\"},{\"name\":\"closing_address\",\"type\":\"address\"},{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"updated\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"settleTimeout\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"token_address\",\"type\":\"address\"},{\"name\":\"participant1\",\"type\":\"address\"},{\"name\":\"participant2\",\"type\":\"address\"},{\"name\":\"timeout\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"participant\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"balance\",\"type\":\"uint256\"}],\"name\":\"ChannelNewBalance\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"closing_address\",\"type\":\"address\"}],\"name\":\"ChannelClosed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"node_address\",\"type\":\"address\"}],\"name\":\"TransferUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"ChannelSettled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"secret\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"receiver_address\",\"type\":\"address\"}],\"name\":\"ChannelSecretRevealed\",\"type\":\"event\"}]"

// NettingChannelContractBin is the compiled bytecode used for deploying new contracts.
const NettingChannelContractBin = `0x6060604052341561000f57600080fd5b604051608080610b91833981016040528080519190602001805191906020018051919060200180519150819050600681108015906100505750622932e08111155b151561005b57600080fd5b600160a060020a03848116908416141561007457600080fd5b5060058054600160a060020a0319908116600160a060020a03958616908117909255600b805494861694821685179055600091825260116020526040808320805460ff19908116600190811790925595845290832080549095166002179094556004805496909516951694909417909255908255439055610a969081906100fb90396000f3006060604052600436106100ab5763ffffffff60e060020a60003504166311da60b481146100bb578063202ac3bc146100d057806327d120fe1461016557806353af5d10146101d3578063597e1fb5146102025780635e1fc56e146102275780635f88eade146102955780636d2381b3146102a857806373d4a13a146102f15780637ebdc478146103465780639d76ea5814610359578063b32c65c81461036c578063b6b55f25146103f6575b34156100b657600080fd5b600080fd5b34156100c657600080fd5b6100ce610420565b005b34156100db57600080fd5b6100ce60046024813581810190830135806020601f8201819004810201604051908101604052818152929190602084018383808284378201915050505050509190803590602001908201803590602001908080601f01602080910402602001604051908101604052818152929190602084018383808284375094965050933593506104b292505050565b341561017057600080fd5b6100ce6004803567ffffffffffffffff1690602480359160443591606435919060a49060843590810190830135806020601f8201819004810201604051908101604052818152929190602084018383808284375094965061064095505050505050565b34156101de57600080fd5b6101e6610783565b604051600160a060020a03909116815260200160405180910390f35b341561020d57600080fd5b610215610792565b60405190815260200160405180910390f35b341561023257600080fd5b6100ce6004803567ffffffffffffffff1690602480359160443591606435919060a49060843590810190830135806020601f8201819004810201604051908101604052818152929190602084018383808284375094965061079895505050505050565b34156102a057600080fd5b6102156108db565b34156102b357600080fd5b6102bb6108e1565b604051600160a060020a039485168152602081019390935292166040808301919091526060820192909252608001905180910390f35b34156102fc57600080fd5b610304610901565b6040519586526020860194909452604080860193909352600160a060020a03918216606086015216608084015290151560a083015260c0909101905180910390f35b341561035157600080fd5b610215610929565b341561036457600080fd5b6101e661092f565b341561037757600080fd5b61037f61093e565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156103bb5780820151838201526020016103a3565b50505050905090810190601f1680156103e85780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561040157600080fd5b61040c600435610975565b604051901515815260200160405180910390f35b73__NettingChannelLibrary.sol:NettingCha__63de394e0d600060405160e060020a63ffffffff8416028152600481019190915260240160006040518083038186803b151561047057600080fd5b6102c65a03f4151561048157600080fd5b5050507f6713dea2491bc95585ea9be0d6993fc7790fdcd04f495a7e7592fbd80bbe00de60405160405180910390a1565b73__NettingChannelLibrary.sol:NettingCha__63c2522462600085858560405160e060020a63ffffffff871602815260048101858152606482018390526080602483019081529091604481019060840186818151815260200191508051906020019080838360005b8381101561053457808201518382015260200161051c565b50505050905090810190601f1680156105615780820380516001836020036101000a031916815260200191505b50838103825285818151815260200191508051906020019080838360005b8381101561059757808201518382015260200161057f565b50505050905090810190601f1680156105c45780820380516001836020036101000a031916815260200191505b50965050505050505060006040518083038186803b15156105e457600080fd5b6102c65a03f415156105f557600080fd5b5050507fa2e2842eefea7e32abccccd9d3fae92608319362c3905ef73de44938c05925368133604051918252600160a060020a031660208201526040908101905180910390a1505050565b73__NettingChannelLibrary.sol:NettingCha__63f565eb366000878787878760405160e060020a63ffffffff89160281526004810187815267ffffffffffffffff8716602483015260448201869052606482018590526084820184905260c060a48301908152909160c40183818151815260200191508051906020019080838360005b838110156106dd5780820151838201526020016106c5565b50505050905090810190601f16801561070a5780820380516001836020036101000a031916815260200191505b5097505050505050505060006040518083038186803b151561072b57600080fd5b6102c65a03f4151561073c57600080fd5b5050507fa0379b1bd0864245b4ff39bf6a023065e80d0e9276d2671d94b9f653b4bbcdfe33604051600160a060020a03909116815260200160405180910390a15050505050565b600354600160a060020a031690565b60025490565b73__NettingChannelLibrary.sol:NettingCha__63c800b0026000878787878760405160e060020a63ffffffff89160281526004810187815267ffffffffffffffff8716602483015260448201869052606482018590526084820184905260c060a48301908152909160c40183818151815260200191508051906020019080838360005b8381101561083557808201518382015260200161081d565b50505050905090810190601f1680156108625780820380516001836020036101000a031916815260200191505b5097505050505050505060006040518083038186803b151561088357600080fd5b6102c65a03f4151561089457600080fd5b5050507f93daadf23cc2150b386a6c3b39f6e61b9c555fc1cec423e4c93ac9d36b008fef33604051600160a060020a03909116815260200160405180910390a15050505050565b60015490565b600554600654600b54600c54600160a060020a0393841694929390911691565b600054600154600254600354600454601254600160a060020a03928316929091169060ff1686565b60005490565b600454600160a060020a031690565b60408051908101604052600581527f302e322e5f000000000000000000000000000000000000000000000000000000602082015281565b6000808073__NettingChannelLibrary.sol:NettingCha__633268a05a8286816040516040015260405160e060020a63ffffffff851602815260048101929092526024820152604401604080518083038186803b15156109d557600080fd5b6102c65a03f415156109e657600080fd5b50505060405180519060200180519193509091505060018215151415610a63576004547f9cb02993ef7311b37acc6bdfc1a8397160be258a877d78b31f4e366caf7bfcef90600160a060020a03163383604051600160a060020a039384168152919092166020820152604080820192909252606001905180910390a15b50929150505600a165627a7a72305820aadeb80cd8d5788100aa7a4ed132fe819df76e8dff96190283fef48ae43433510029`

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
	return address, tx, &NettingChannelContract{NettingChannelContractCaller: NettingChannelContractCaller{contract: contract}, NettingChannelContractTransactor: NettingChannelContractTransactor{contract: contract}}, nil
}

// NettingChannelContract is an auto generated Go binding around an Ethereum contract.
type NettingChannelContract struct {
	NettingChannelContractCaller     // Read-only binding to the contract
	NettingChannelContractTransactor // Write-only binding to the contract
}

// NettingChannelContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type NettingChannelContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NettingChannelContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NettingChannelContractTransactor struct {
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
	contract, err := bindNettingChannelContract(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NettingChannelContract{NettingChannelContractCaller: NettingChannelContractCaller{contract: contract}, NettingChannelContractTransactor: NettingChannelContractTransactor{contract: contract}}, nil
}

// NewNettingChannelContractCaller creates a new read-only instance of NettingChannelContract, bound to a specific deployed contract.
func NewNettingChannelContractCaller(address common.Address, caller bind.ContractCaller) (*NettingChannelContractCaller, error) {
	contract, err := bindNettingChannelContract(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &NettingChannelContractCaller{contract: contract}, nil
}

// NewNettingChannelContractTransactor creates a new write-only instance of NettingChannelContract, bound to a specific deployed contract.
func NewNettingChannelContractTransactor(address common.Address, transactor bind.ContractTransactor) (*NettingChannelContractTransactor, error) {
	contract, err := bindNettingChannelContract(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &NettingChannelContractTransactor{contract: contract}, nil
}

// bindNettingChannelContract binds a generic wrapper to an already deployed contract.
func bindNettingChannelContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NettingChannelContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
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

// RegistryABI is the input ABI used to generate the binding from.
const RegistryABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"registry\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"token_address\",\"type\":\"address\"}],\"name\":\"channelManagerByToken\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokens\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"channelManagerAddresses\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenAddresses\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token_address\",\"type\":\"address\"}],\"name\":\"addToken\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"channel_manager_address\",\"type\":\"address\"}],\"name\":\"TokenAdded\",\"type\":\"event\"}]"

// RegistryBin is the compiled bytecode used for deploying new contracts.
const RegistryBin = `0x6060604052341561000f57600080fd5b610dff8061001e6000396000f3006060604052600436106100825763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663038defd7811461009257806325119b5f146100cd5780634f64b2be146100ec578063640191e214610102578063a9989b9314610168578063b32c65c81461017b578063d48bfca714610205575b341561008d57600080fd5b600080fd5b341561009d57600080fd5b6100b1600160a060020a0360043516610224565b604051600160a060020a03909116815260200160405180910390f35b34156100d857600080fd5b6100b1600160a060020a036004351661023f565b34156100f757600080fd5b6100b1600435610289565b341561010d57600080fd5b6101156102b1565b60405160208082528190810183818151815260200191508051906020019060200280838360005b8381101561015457808201518382015260200161013c565b505050509050019250505060405180910390f35b341561017357600080fd5b610115610363565b341561018657600080fd5b61018e6103cc565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156101ca5780820151838201526020016101b2565b50505050905090810190601f1680156101f75780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561021057600080fd5b6100b1600160a060020a0360043516610403565b600060208190529081526040902054600160a060020a031681565b600160a060020a038082166000908152602081905260408120549091839116151561026957600080fd5b5050600160a060020a039081166000908152602081905260409020541690565b600180548290811061029757fe5b600091825260209091200154600160a060020a0316905081565b6102b96105ad565b6000806102c46105ad565b6001546040518059106102d45750595b90808252806020026020018201604052509050600092505b60015483101561035c57600180548490811061030457fe5b6000918252602080832090910154600160a060020a03908116808452918390526040909220549093501681848151811061033a57fe5b600160a060020a039092166020928302909101909101526001909201916102ec565b9392505050565b61036b6105ad565b60018054806020026020016040519081016040528092919081815260200182805480156103c157602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116103a3575b505050505090505b90565b60408051908101604052600581527f302e322e5f000000000000000000000000000000000000000000000000000000602082015281565b600160a060020a038082166000908152602081905260408120549091829184918391161561043057600080fd5b5080600160a060020a0381166318160ddd6000604051602001526040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401602060405180830381600087803b151561049157600080fd5b6102c65a03f115156104a257600080fd5b5050506040518051905050846104b66105bf565b600160a060020a039091168152602001604051809103906000f08015156104dc57600080fd5b600160a060020a038681166000908152602081905260409020805473ffffffffffffffffffffffffffffffffffffffff1916918316919091179055600180549194509080820161052c83826105cf565b506000918252602090912001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0387161790557fdffbd9ded1c09446f09377de547142dcce7dc541c8b0b028142b1eba7026b9e78584604051600160a060020a039283168152911660208201526040908101905180910390a150909392505050565b60206040519081016040526000815290565b6040516107bd8061061783390190565b8154818355818115116105f3576000838152602090206105f39181019083016105f8565b505050565b6103c991905b8082111561061257600081556001016105fe565b509056006060604052341561000f57600080fd5b6040516020806107bd8339810160405280805160008054600160a060020a03909216600160a060020a03199092169190911790555050610769806100546000396000f3006060604052600436106100745763ffffffff60e060020a6000350416630b74b6208114610084578063238bfba2146100ea5780636785b500146101255780636cb30fee146101385780637709bc78146101575780639d76ea581461018a578063b32c65c81461019d578063f26c6aed14610227575b341561007f57600080fd5b600080fd5b341561008f57600080fd5b610097610249565b60405160208082528190810183818151815260200191508051906020019060200280838360005b838110156100d65780820151838201526020016100be565b505050509050019250505060405180910390f35b34156100f557600080fd5b610109600160a060020a036004351661040a565b604051600160a060020a03909116815260200160405180910390f35b341561013057600080fd5b610097610492565b341561014357600080fd5b610097600160a060020a03600435166104f9565b341561016257600080fd5b610176600160a060020a0360043516610589565b604051901515815260200160405180910390f35b341561019557600080fd5b610109610591565b34156101a857600080fd5b6101b06105a0565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156101ec5780820151838201526020016101d4565b50505050905090810190601f1680156102195780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561023257600080fd5b610109600160a060020a03600435166024356105d7565b61025161072b565b60008061025c61072b565b60009250828080805b6001548710156102b3576001805461029c91908990811061028257fe5b600091825260209091200154600160a060020a0316610589565b156102a8576001830192505b600190960195610265565b826002026040518059106102c45750595b9080825280602002602001820160405250945060009550600096505b6001548710156103fe57600180546102fd91908990811061028257fe5b1515610308576103f3565b600180548890811061031657fe5b6000918252602082200154600160a060020a031694508490636d2381b390604051608001526040518163ffffffff1660e060020a028152600401608060405180830381600087803b151561036957600080fd5b6102c65a03f1151561037a57600080fd5b50505060405180519060200180519060200180519060200180519050509250509150818587815181106103a957fe5b600160a060020a0390921660209283029091019091015260019590950194808587815181106103d457fe5b600160a060020a03909216602092830290910190910152600195909501945b6001909601956102e0565b50929695505050505050565b600073__ChannelManagerLibrary.sol:ChannelMan__638a1c00e28284816040516020015260405160e060020a63ffffffff85160281526004810192909252600160a060020a0316602482015260440160206040518083038186803b151561047257600080fd5b6102c65a03f4151561048357600080fd5b50505060405180519392505050565b61049a61072b565b600180546020808202016040519081016040528092919081815260200182805480156104ef57602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116104d1575b5050505050905090565b61050161072b565b6000600301600083600160a060020a0316600160a060020a0316815260200190815260200160002080548060200260200160405190810160405280929190818152602001828054801561057d57602002820191906000526020600020905b8154600160a060020a0316815260019091019060200180831161055f575b50505050509050919050565b6000903b1190565b600054600160a060020a031690565b60408051908101604052600581527f302e322e5f000000000000000000000000000000000000000000000000000000602082015281565b60008060006105e58561040a565b9150600160a060020a0382161561063d577fda8d2f351e0f7c8c368e631ce8ab15973e7582ece0c347d75a5cff49eb899eb73386604051600160a060020a039283168152911660208201526040908101905180910390a15b73__ChannelManagerLibrary.sol:ChannelMan__63941583a560008787826040516020015260405160e060020a63ffffffff86160281526004810193909352600160a060020a039091166024830152604482015260640160206040518083038186803b15156106ac57600080fd5b6102c65a03f415156106bd57600080fd5b5050506040518051905090507f7bd269696a33040df6c111efd58439c9c77909fcbe90f7511065ac277e175dac81338787604051600160a060020a039485168152928416602084015292166040808301919091526060820192909252608001905180910390a1949350505050565b602060405190810160405260008152905600a165627a7a72305820522416bac5fb22953a513c007ec3c6a51edd7d61b9bd1cc408c4e121fd40c6db0029a165627a7a72305820b0b1960fe79f5edcb65312056a74c7c8307e6b657689829cf345352e0c09723a0029`

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
	return address, tx, &Registry{RegistryCaller: RegistryCaller{contract: contract}, RegistryTransactor: RegistryTransactor{contract: contract}}, nil
}

// Registry is an auto generated Go binding around an Ethereum contract.
type Registry struct {
	RegistryCaller     // Read-only binding to the contract
	RegistryTransactor // Write-only binding to the contract
}

// RegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type RegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RegistryTransactor struct {
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
	contract, err := bindRegistry(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Registry{RegistryCaller: RegistryCaller{contract: contract}, RegistryTransactor: RegistryTransactor{contract: contract}}, nil
}

// NewRegistryCaller creates a new read-only instance of Registry, bound to a specific deployed contract.
func NewRegistryCaller(address common.Address, caller bind.ContractCaller) (*RegistryCaller, error) {
	contract, err := bindRegistry(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryCaller{contract: contract}, nil
}

// NewRegistryTransactor creates a new write-only instance of Registry, bound to a specific deployed contract.
func NewRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*RegistryTransactor, error) {
	contract, err := bindRegistry(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &RegistryTransactor{contract: contract}, nil
}

// bindRegistry binds a generic wrapper to an already deployed contract.
func bindRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
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
	return address, tx, &Token{TokenCaller: TokenCaller{contract: contract}, TokenTransactor: TokenTransactor{contract: contract}}, nil
}

// Token is an auto generated Go binding around an Ethereum contract.
type Token struct {
	TokenCaller     // Read-only binding to the contract
	TokenTransactor // Write-only binding to the contract
}

// TokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenTransactor struct {
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
	contract, err := bindToken(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Token{TokenCaller: TokenCaller{contract: contract}, TokenTransactor: TokenTransactor{contract: contract}}, nil
}

// NewTokenCaller creates a new read-only instance of Token, bound to a specific deployed contract.
func NewTokenCaller(address common.Address, caller bind.ContractCaller) (*TokenCaller, error) {
	contract, err := bindToken(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCaller{contract: contract}, nil
}

// NewTokenTransactor creates a new write-only instance of Token, bound to a specific deployed contract.
func NewTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenTransactor, error) {
	contract, err := bindToken(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &TokenTransactor{contract: contract}, nil
}

// bindToken binds a generic wrapper to an already deployed contract.
func bindToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
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

// UtilsABI is the input ABI used to generate the binding from.
const UtilsABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"channel\",\"type\":\"address\"}],\"name\":\"contractExists\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// UtilsBin is the compiled bytecode used for deploying new contracts.
const UtilsBin = `0x6060604052341561000f57600080fd5b6101858061001e6000396000f30060606040526004361061004b5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416637709bc788114610050578063b32c65c814610090575b600080fd5b341561005b57600080fd5b61007c73ffffffffffffffffffffffffffffffffffffffff6004351661011a565b604051901515815260200160405180910390f35b341561009b57600080fd5b6100a3610122565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156100df5780820151838201526020016100c7565b50505050905090810190601f16801561010c5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000903b1190565b60408051908101604052600581527f302e322e5f0000000000000000000000000000000000000000000000000000006020820152815600a165627a7a72305820cde168eae39ff00f1b0e4d706d39fcc56ca88fc5ea4b0791c327fa6eea0b45140029`

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
	return address, tx, &Utils{UtilsCaller: UtilsCaller{contract: contract}, UtilsTransactor: UtilsTransactor{contract: contract}}, nil
}

// Utils is an auto generated Go binding around an Ethereum contract.
type Utils struct {
	UtilsCaller     // Read-only binding to the contract
	UtilsTransactor // Write-only binding to the contract
}

// UtilsCaller is an auto generated read-only Go binding around an Ethereum contract.
type UtilsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UtilsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UtilsTransactor struct {
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
	contract, err := bindUtils(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Utils{UtilsCaller: UtilsCaller{contract: contract}, UtilsTransactor: UtilsTransactor{contract: contract}}, nil
}

// NewUtilsCaller creates a new read-only instance of Utils, bound to a specific deployed contract.
func NewUtilsCaller(address common.Address, caller bind.ContractCaller) (*UtilsCaller, error) {
	contract, err := bindUtils(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &UtilsCaller{contract: contract}, nil
}

// NewUtilsTransactor creates a new write-only instance of Utils, bound to a specific deployed contract.
func NewUtilsTransactor(address common.Address, transactor bind.ContractTransactor) (*UtilsTransactor, error) {
	contract, err := bindUtils(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &UtilsTransactor{contract: contract}, nil
}

// bindUtils binds a generic wrapper to an already deployed contract.
func bindUtils(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UtilsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
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
