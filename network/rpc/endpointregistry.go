// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rpc

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// EndpointRegistryABI is the input ABI used to generate the binding from.
const EndpointRegistryABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"eth_address\",\"type\":\"address\"}],\"name\":\"findEndpointByAddress\",\"outputs\":[{\"name\":\"socket\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"socket\",\"type\":\"string\"}],\"name\":\"registerEndpoint\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"socket\",\"type\":\"string\"}],\"name\":\"findAddressByEndpoint\",\"outputs\":[{\"name\":\"eth_address\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"eth_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"socket\",\"type\":\"string\"}],\"name\":\"AddressRegistered\",\"type\":\"event\"}]"

// EndpointRegistryBin is the compiled bytecode used for deploying new contracts.
const EndpointRegistryBin = `0x6060604052341561000f57600080fd5b6107898061001e6000396000f3006060604052600436106100565763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663028a582e811461005b57806308b5a85a146100f1578063460123cf14610144575b600080fd5b341561006657600080fd5b61007a600160a060020a03600435166101b1565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156100b657808201518382015260200161009e565b50505050905090810190601f1680156100e35780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34156100fc57600080fd5b61014260046024813581810190830135806020601f8201819004810201604051908101604052818152929190602084018383808284375094965061027d95505050505050565b005b341561014f57600080fd5b61019560046024813581810190830135806020601f8201819004810201604051908101604052818152929190602084018383808284375094965061056495505050505050565b604051600160a060020a03909116815260200160405180910390f35b6101b96106b0565b60008083600160a060020a0316600160a060020a031681526020019081526020016000208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156102715780601f1061024657610100808354040283529160200191610271565b820191906000526020600020905b81548152906001019060200180831161025457829003601f168201915b50505050509050919050565b600081610298816020604051908101604052600081526105dc565b1515600114156102a757600080fd5b60008033600160a060020a0316600160a060020a031681526020019081526020016000209150610370828054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156103655780601f1061033a57610100808354040283529160200191610365565b820191906000526020600020905b81548152906001019060200180831161034857829003601f168201915b5050505050846105dc565b1561037a5761055f565b600060018360405180828054600181600116156101000203166002900480156103da5780601f106103b85761010080835404028352918201916103da565b820191906000526020600020905b8154815290600101906020018083116103c6575b50509283525050602001604051908190039020805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03928316179055331660009081526020819052604090208380516104359291602001906106c2565b50336001846040518082805190602001908083835b602083106104695780518252601f19909201916020918201910161044a565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051908190039020805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392831617905533167f3a62a9d7855df5303e50b0440124304fefafde7f677fc33787b784fc92cfa6188460405160208082528190810183818151815260200191508051906020019080838360005b8381101561052457808201518382015260200161050c565b50505050905090810190601f1680156105515780820380516001836020036101000a031916815260200191505b509250505060405180910390a25b505050565b60006001826040518082805190602001908083835b602083106105985780518252601f199092019160209182019101610579565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405190819003902054600160a060020a031692915050565b6000816040518082805190602001908083835b6020831061060e5780518252601f1990920191602091820191016105ef565b6001836020036101000a0380198251168184511617909252505050919091019250604091505051908190039020836040518082805190602001908083835b6020831061066b5780518252601f19909201916020918201910161064c565b6001836020036101000a038019825116818451161790925250505091909101925060409150505190819003902014156106a6575060016106aa565b5060005b92915050565b60206040519081016040526000815290565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061070357805160ff1916838001178555610730565b82800160010185558215610730579182015b82811115610730578251825591602001919060010190610715565b5061073c929150610740565b5090565b61075a91905b8082111561073c5760008155600101610746565b905600a165627a7a723058209341730c2fcadb38de8b8868113d3da9dc1a3645851e9ec4aadc8c49781b4ddc0029`

// DeployEndpointRegistry deploys a new Ethereum contract, binding an instance of EndpointRegistry to it.
func DeployEndpointRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *EndpointRegistry, error) {
	parsed, err := abi.JSON(strings.NewReader(EndpointRegistryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(EndpointRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EndpointRegistry{EndpointRegistryCaller: EndpointRegistryCaller{contract: contract}, EndpointRegistryTransactor: EndpointRegistryTransactor{contract: contract}}, nil
}

// EndpointRegistry is an auto generated Go binding around an Ethereum contract.
type EndpointRegistry struct {
	EndpointRegistryCaller     // Read-only binding to the contract
	EndpointRegistryTransactor // Write-only binding to the contract
}

// EndpointRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type EndpointRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EndpointRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EndpointRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EndpointRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EndpointRegistrySession struct {
	Contract     *EndpointRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EndpointRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EndpointRegistryCallerSession struct {
	Contract *EndpointRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// EndpointRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EndpointRegistryTransactorSession struct {
	Contract     *EndpointRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// EndpointRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type EndpointRegistryRaw struct {
	Contract *EndpointRegistry // Generic contract binding to access the raw methods on
}

// EndpointRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EndpointRegistryCallerRaw struct {
	Contract *EndpointRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// EndpointRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EndpointRegistryTransactorRaw struct {
	Contract *EndpointRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEndpointRegistry creates a new instance of EndpointRegistry, bound to a specific deployed contract.
func NewEndpointRegistry(address common.Address, backend bind.ContractBackend) (*EndpointRegistry, error) {
	contract, err := bindEndpointRegistry(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EndpointRegistry{EndpointRegistryCaller: EndpointRegistryCaller{contract: contract}, EndpointRegistryTransactor: EndpointRegistryTransactor{contract: contract}}, nil
}

// NewEndpointRegistryCaller creates a new read-only instance of EndpointRegistry, bound to a specific deployed contract.
func NewEndpointRegistryCaller(address common.Address, caller bind.ContractCaller) (*EndpointRegistryCaller, error) {
	contract, err := bindEndpointRegistry(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &EndpointRegistryCaller{contract: contract}, nil
}

// NewEndpointRegistryTransactor creates a new write-only instance of EndpointRegistry, bound to a specific deployed contract.
func NewEndpointRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*EndpointRegistryTransactor, error) {
	contract, err := bindEndpointRegistry(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &EndpointRegistryTransactor{contract: contract}, nil
}

// bindEndpointRegistry binds a generic wrapper to an already deployed contract.
func bindEndpointRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EndpointRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EndpointRegistry *EndpointRegistryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _EndpointRegistry.Contract.EndpointRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EndpointRegistry *EndpointRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EndpointRegistry.Contract.EndpointRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EndpointRegistry *EndpointRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EndpointRegistry.Contract.EndpointRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EndpointRegistry *EndpointRegistryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _EndpointRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EndpointRegistry *EndpointRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EndpointRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EndpointRegistry *EndpointRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EndpointRegistry.Contract.contract.Transact(opts, method, params...)
}

// FindAddressByEndpoint is a free data retrieval call binding the contract method 0x460123cf.
//
// Solidity: function findAddressByEndpoint(socket string) constant returns(eth_address address)
func (_EndpointRegistry *EndpointRegistryCaller) FindAddressByEndpoint(opts *bind.CallOpts, socket string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _EndpointRegistry.contract.Call(opts, out, "findAddressByEndpoint", socket)
	return *ret0, err
}

// FindAddressByEndpoint is a free data retrieval call binding the contract method 0x460123cf.
//
// Solidity: function findAddressByEndpoint(socket string) constant returns(eth_address address)
func (_EndpointRegistry *EndpointRegistrySession) FindAddressByEndpoint(socket string) (common.Address, error) {
	return _EndpointRegistry.Contract.FindAddressByEndpoint(&_EndpointRegistry.CallOpts, socket)
}

// FindAddressByEndpoint is a free data retrieval call binding the contract method 0x460123cf.
//
// Solidity: function findAddressByEndpoint(socket string) constant returns(eth_address address)
func (_EndpointRegistry *EndpointRegistryCallerSession) FindAddressByEndpoint(socket string) (common.Address, error) {
	return _EndpointRegistry.Contract.FindAddressByEndpoint(&_EndpointRegistry.CallOpts, socket)
}

// FindEndpointByAddress is a free data retrieval call binding the contract method 0x028a582e.
//
// Solidity: function findEndpointByAddress(eth_address address) constant returns(socket string)
func (_EndpointRegistry *EndpointRegistryCaller) FindEndpointByAddress(opts *bind.CallOpts, eth_address common.Address) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _EndpointRegistry.contract.Call(opts, out, "findEndpointByAddress", eth_address)
	return *ret0, err
}

// FindEndpointByAddress is a free data retrieval call binding the contract method 0x028a582e.
//
// Solidity: function findEndpointByAddress(eth_address address) constant returns(socket string)
func (_EndpointRegistry *EndpointRegistrySession) FindEndpointByAddress(eth_address common.Address) (string, error) {
	return _EndpointRegistry.Contract.FindEndpointByAddress(&_EndpointRegistry.CallOpts, eth_address)
}

// FindEndpointByAddress is a free data retrieval call binding the contract method 0x028a582e.
//
// Solidity: function findEndpointByAddress(eth_address address) constant returns(socket string)
func (_EndpointRegistry *EndpointRegistryCallerSession) FindEndpointByAddress(eth_address common.Address) (string, error) {
	return _EndpointRegistry.Contract.FindEndpointByAddress(&_EndpointRegistry.CallOpts, eth_address)
}

// RegisterEndpoint is a paid mutator transaction binding the contract method 0x08b5a85a.
//
// Solidity: function registerEndpoint(socket string) returns()
func (_EndpointRegistry *EndpointRegistryTransactor) RegisterEndpoint(opts *bind.TransactOpts, socket string) (*types.Transaction, error) {
	return _EndpointRegistry.contract.Transact(opts, "registerEndpoint", socket)
}

// RegisterEndpoint is a paid mutator transaction binding the contract method 0x08b5a85a.
//
// Solidity: function registerEndpoint(socket string) returns()
func (_EndpointRegistry *EndpointRegistrySession) RegisterEndpoint(socket string) (*types.Transaction, error) {
	return _EndpointRegistry.Contract.RegisterEndpoint(&_EndpointRegistry.TransactOpts, socket)
}

// RegisterEndpoint is a paid mutator transaction binding the contract method 0x08b5a85a.
//
// Solidity: function registerEndpoint(socket string) returns()
func (_EndpointRegistry *EndpointRegistryTransactorSession) RegisterEndpoint(socket string) (*types.Transaction, error) {
	return _EndpointRegistry.Contract.RegisterEndpoint(&_EndpointRegistry.TransactOpts, socket)
}
