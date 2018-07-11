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

// ECVerifyABI is the input ABI used to generate the binding from.
const ECVerifyABI = "[]"

// ECVerifyBin is the compiled bytecode used for deploying new contracts.
const ECVerifyBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a723058208377d8353eea6842efdc9fb1868e4fa1cf41cdd5d066638b36f4e3cdc2919d2d0029`

// DeployECVerify deploys a new Ethereum contract, binding an instance of ECVerify to it.
func DeployECVerify(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ECVerify, error) {
	parsed, err := abi.JSON(strings.NewReader(ECVerifyABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ECVerifyBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ECVerify{ECVerifyCaller: ECVerifyCaller{contract: contract}, ECVerifyTransactor: ECVerifyTransactor{contract: contract}, ECVerifyFilterer: ECVerifyFilterer{contract: contract}}, nil
}

// ECVerify is an auto generated Go binding around an Ethereum contract.
type ECVerify struct {
	ECVerifyCaller     // Read-only binding to the contract
	ECVerifyTransactor // Write-only binding to the contract
	ECVerifyFilterer   // Log filterer for contract events
}

// ECVerifyCaller is an auto generated read-only Go binding around an Ethereum contract.
type ECVerifyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ECVerifyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ECVerifyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ECVerifyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ECVerifyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ECVerifySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ECVerifySession struct {
	Contract     *ECVerify         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ECVerifyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ECVerifyCallerSession struct {
	Contract *ECVerifyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ECVerifyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ECVerifyTransactorSession struct {
	Contract     *ECVerifyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ECVerifyRaw is an auto generated low-level Go binding around an Ethereum contract.
type ECVerifyRaw struct {
	Contract *ECVerify // Generic contract binding to access the raw methods on
}

// ECVerifyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ECVerifyCallerRaw struct {
	Contract *ECVerifyCaller // Generic read-only contract binding to access the raw methods on
}

// ECVerifyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ECVerifyTransactorRaw struct {
	Contract *ECVerifyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewECVerify creates a new instance of ECVerify, bound to a specific deployed contract.
func NewECVerify(address common.Address, backend bind.ContractBackend) (*ECVerify, error) {
	contract, err := bindECVerify(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ECVerify{ECVerifyCaller: ECVerifyCaller{contract: contract}, ECVerifyTransactor: ECVerifyTransactor{contract: contract}, ECVerifyFilterer: ECVerifyFilterer{contract: contract}}, nil
}

// NewECVerifyCaller creates a new read-only instance of ECVerify, bound to a specific deployed contract.
func NewECVerifyCaller(address common.Address, caller bind.ContractCaller) (*ECVerifyCaller, error) {
	contract, err := bindECVerify(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ECVerifyCaller{contract: contract}, nil
}

// NewECVerifyTransactor creates a new write-only instance of ECVerify, bound to a specific deployed contract.
func NewECVerifyTransactor(address common.Address, transactor bind.ContractTransactor) (*ECVerifyTransactor, error) {
	contract, err := bindECVerify(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ECVerifyTransactor{contract: contract}, nil
}

// NewECVerifyFilterer creates a new log filterer instance of ECVerify, bound to a specific deployed contract.
func NewECVerifyFilterer(address common.Address, filterer bind.ContractFilterer) (*ECVerifyFilterer, error) {
	contract, err := bindECVerify(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ECVerifyFilterer{contract: contract}, nil
}

// bindECVerify binds a generic wrapper to an already deployed contract.
func bindECVerify(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ECVerifyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ECVerify *ECVerifyRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ECVerify.Contract.ECVerifyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ECVerify *ECVerifyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ECVerify.Contract.ECVerifyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ECVerify *ECVerifyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ECVerify.Contract.ECVerifyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ECVerify *ECVerifyCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ECVerify.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ECVerify *ECVerifyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ECVerify.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ECVerify *ECVerifyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ECVerify.Contract.contract.Transact(opts, method, params...)
}

// SecretRegistryABI is the input ABI used to generate the binding from.
const SecretRegistryABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"secret\",\"type\":\"bytes32\"}],\"name\":\"registerSecret\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"secrethash_to_block\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"secrethash\",\"type\":\"bytes32\"}],\"name\":\"getSecretRevealBlockHeight\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"secrethash\",\"type\":\"bytes32\"}],\"name\":\"SecretRevealed\",\"type\":\"event\"}]"

// SecretRegistryBin is the compiled bytecode used for deploying new contracts.
const SecretRegistryBin = `0x608060405234801561001057600080fd5b506102cd806100206000396000f3006080604052600436106100615763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166312ad8bfc81146100665780639734030914610092578063b32c65c8146100bc578063c1f6294614610146575b600080fd5b34801561007257600080fd5b5061007e60043561015e565b604080519115158252519081900360200190f35b34801561009e57600080fd5b506100aa600435610246565b60408051918252519081900360200190f35b3480156100c857600080fd5b506100d1610258565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561010b5781810151838201526020016100f3565b50505050905090810190601f1680156101385780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561015257600080fd5b506100aa60043561028f565b604080516020808201849052825180830382018152918301928390528151600093849392909182918401908083835b602083106101ac5780518252601f19909201916020918201910161018d565b5181516020939093036101000a60001901801990911692169190911790526040519201829003909120935050841591508190506101f55750600081815260208190526040812054115b156102035760009150610240565b6000818152602081905260408082204390555182917f9b7ddc883342824bd7ddbff103e7a69f8f2e60b96c075cd1b8b8b9713ecc75a491a2600191505b50919050565b60006020819052908152604090205481565b60408051808201909152600581527f302e332e5f000000000000000000000000000000000000000000000000000000602082015281565b600090815260208190526040902054905600a165627a7a72305820f69c8bac516d7b33d97cfff0956fa761274de6891a330e80786dd2e9e978b5e10029`

// DeploySecretRegistry deploys a new Ethereum contract, binding an instance of SecretRegistry to it.
func DeploySecretRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SecretRegistry, error) {
	parsed, err := abi.JSON(strings.NewReader(SecretRegistryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SecretRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SecretRegistry{SecretRegistryCaller: SecretRegistryCaller{contract: contract}, SecretRegistryTransactor: SecretRegistryTransactor{contract: contract}, SecretRegistryFilterer: SecretRegistryFilterer{contract: contract}}, nil
}

// SecretRegistry is an auto generated Go binding around an Ethereum contract.
type SecretRegistry struct {
	SecretRegistryCaller     // Read-only binding to the contract
	SecretRegistryTransactor // Write-only binding to the contract
	SecretRegistryFilterer   // Log filterer for contract events
}

// SecretRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type SecretRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SecretRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SecretRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SecretRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SecretRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SecretRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SecretRegistrySession struct {
	Contract     *SecretRegistry   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SecretRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SecretRegistryCallerSession struct {
	Contract *SecretRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// SecretRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SecretRegistryTransactorSession struct {
	Contract     *SecretRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// SecretRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type SecretRegistryRaw struct {
	Contract *SecretRegistry // Generic contract binding to access the raw methods on
}

// SecretRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SecretRegistryCallerRaw struct {
	Contract *SecretRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// SecretRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SecretRegistryTransactorRaw struct {
	Contract *SecretRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSecretRegistry creates a new instance of SecretRegistry, bound to a specific deployed contract.
func NewSecretRegistry(address common.Address, backend bind.ContractBackend) (*SecretRegistry, error) {
	contract, err := bindSecretRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SecretRegistry{SecretRegistryCaller: SecretRegistryCaller{contract: contract}, SecretRegistryTransactor: SecretRegistryTransactor{contract: contract}, SecretRegistryFilterer: SecretRegistryFilterer{contract: contract}}, nil
}

// NewSecretRegistryCaller creates a new read-only instance of SecretRegistry, bound to a specific deployed contract.
func NewSecretRegistryCaller(address common.Address, caller bind.ContractCaller) (*SecretRegistryCaller, error) {
	contract, err := bindSecretRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SecretRegistryCaller{contract: contract}, nil
}

// NewSecretRegistryTransactor creates a new write-only instance of SecretRegistry, bound to a specific deployed contract.
func NewSecretRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*SecretRegistryTransactor, error) {
	contract, err := bindSecretRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SecretRegistryTransactor{contract: contract}, nil
}

// NewSecretRegistryFilterer creates a new log filterer instance of SecretRegistry, bound to a specific deployed contract.
func NewSecretRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*SecretRegistryFilterer, error) {
	contract, err := bindSecretRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SecretRegistryFilterer{contract: contract}, nil
}

// bindSecretRegistry binds a generic wrapper to an already deployed contract.
func bindSecretRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SecretRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SecretRegistry *SecretRegistryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SecretRegistry.Contract.SecretRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SecretRegistry *SecretRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SecretRegistry.Contract.SecretRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SecretRegistry *SecretRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SecretRegistry.Contract.SecretRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SecretRegistry *SecretRegistryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SecretRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SecretRegistry *SecretRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SecretRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SecretRegistry *SecretRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SecretRegistry.Contract.contract.Transact(opts, method, params...)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_SecretRegistry *SecretRegistryCaller) Contract_version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _SecretRegistry.contract.Call(opts, out, "contract_version")
	return *ret0, err
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_SecretRegistry *SecretRegistrySession) Contract_version() (string, error) {
	return _SecretRegistry.Contract.Contract_version(&_SecretRegistry.CallOpts)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_SecretRegistry *SecretRegistryCallerSession) Contract_version() (string, error) {
	return _SecretRegistry.Contract.Contract_version(&_SecretRegistry.CallOpts)
}

// GetSecretRevealBlockHeight is a free data retrieval call binding the contract method 0xc1f62946.
//
// Solidity: function getSecretRevealBlockHeight(secrethash bytes32) constant returns(uint256)
func (_SecretRegistry *SecretRegistryCaller) GetSecretRevealBlockHeight(opts *bind.CallOpts, secrethash [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SecretRegistry.contract.Call(opts, out, "getSecretRevealBlockHeight", secrethash)
	return *ret0, err
}

// GetSecretRevealBlockHeight is a free data retrieval call binding the contract method 0xc1f62946.
//
// Solidity: function getSecretRevealBlockHeight(secrethash bytes32) constant returns(uint256)
func (_SecretRegistry *SecretRegistrySession) GetSecretRevealBlockHeight(secrethash [32]byte) (*big.Int, error) {
	return _SecretRegistry.Contract.GetSecretRevealBlockHeight(&_SecretRegistry.CallOpts, secrethash)
}

// GetSecretRevealBlockHeight is a free data retrieval call binding the contract method 0xc1f62946.
//
// Solidity: function getSecretRevealBlockHeight(secrethash bytes32) constant returns(uint256)
func (_SecretRegistry *SecretRegistryCallerSession) GetSecretRevealBlockHeight(secrethash [32]byte) (*big.Int, error) {
	return _SecretRegistry.Contract.GetSecretRevealBlockHeight(&_SecretRegistry.CallOpts, secrethash)
}

// Secrethash_to_block is a free data retrieval call binding the contract method 0x97340309.
//
// Solidity: function secrethash_to_block( bytes32) constant returns(uint256)
func (_SecretRegistry *SecretRegistryCaller) Secrethash_to_block(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SecretRegistry.contract.Call(opts, out, "secrethash_to_block", arg0)
	return *ret0, err
}

// Secrethash_to_block is a free data retrieval call binding the contract method 0x97340309.
//
// Solidity: function secrethash_to_block( bytes32) constant returns(uint256)
func (_SecretRegistry *SecretRegistrySession) Secrethash_to_block(arg0 [32]byte) (*big.Int, error) {
	return _SecretRegistry.Contract.Secrethash_to_block(&_SecretRegistry.CallOpts, arg0)
}

// Secrethash_to_block is a free data retrieval call binding the contract method 0x97340309.
//
// Solidity: function secrethash_to_block( bytes32) constant returns(uint256)
func (_SecretRegistry *SecretRegistryCallerSession) Secrethash_to_block(arg0 [32]byte) (*big.Int, error) {
	return _SecretRegistry.Contract.Secrethash_to_block(&_SecretRegistry.CallOpts, arg0)
}

// RegisterSecret is a paid mutator transaction binding the contract method 0x12ad8bfc.
//
// Solidity: function registerSecret(secret bytes32) returns(bool)
func (_SecretRegistry *SecretRegistryTransactor) RegisterSecret(opts *bind.TransactOpts, secret [32]byte) (*types.Transaction, error) {
	return _SecretRegistry.contract.Transact(opts, "registerSecret", secret)
}

// RegisterSecret is a paid mutator transaction binding the contract method 0x12ad8bfc.
//
// Solidity: function registerSecret(secret bytes32) returns(bool)
func (_SecretRegistry *SecretRegistrySession) RegisterSecret(secret [32]byte) (*types.Transaction, error) {
	return _SecretRegistry.Contract.RegisterSecret(&_SecretRegistry.TransactOpts, secret)
}

// RegisterSecret is a paid mutator transaction binding the contract method 0x12ad8bfc.
//
// Solidity: function registerSecret(secret bytes32) returns(bool)
func (_SecretRegistry *SecretRegistryTransactorSession) RegisterSecret(secret [32]byte) (*types.Transaction, error) {
	return _SecretRegistry.Contract.RegisterSecret(&_SecretRegistry.TransactOpts, secret)
}

// SecretRegistrySecretRevealedIterator is returned from FilterSecretRevealed and is used to iterate over the raw logs and unpacked data for SecretRevealed events raised by the SecretRegistry contract.
type SecretRegistrySecretRevealedIterator struct {
	Event *SecretRegistrySecretRevealed // Event containing the contract specifics and raw log

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
func (it *SecretRegistrySecretRevealedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SecretRegistrySecretRevealed)
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
		it.Event = new(SecretRegistrySecretRevealed)
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
func (it *SecretRegistrySecretRevealedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SecretRegistrySecretRevealedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SecretRegistrySecretRevealed represents a SecretRevealed event raised by the SecretRegistry contract.
type SecretRegistrySecretRevealed struct {
	Secrethash [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSecretRevealed is a free log retrieval operation binding the contract event 0x9b7ddc883342824bd7ddbff103e7a69f8f2e60b96c075cd1b8b8b9713ecc75a4.
//
// Solidity: event SecretRevealed(secrethash indexed bytes32)
func (_SecretRegistry *SecretRegistryFilterer) FilterSecretRevealed(opts *bind.FilterOpts, secrethash [][32]byte) (*SecretRegistrySecretRevealedIterator, error) {

	var secrethashRule []interface{}
	for _, secrethashItem := range secrethash {
		secrethashRule = append(secrethashRule, secrethashItem)
	}

	logs, sub, err := _SecretRegistry.contract.FilterLogs(opts, "SecretRevealed", secrethashRule)
	if err != nil {
		return nil, err
	}
	return &SecretRegistrySecretRevealedIterator{contract: _SecretRegistry.contract, event: "SecretRevealed", logs: logs, sub: sub}, nil
}

// WatchSecretRevealed is a free log subscription operation binding the contract event 0x9b7ddc883342824bd7ddbff103e7a69f8f2e60b96c075cd1b8b8b9713ecc75a4.
//
// Solidity: event SecretRevealed(secrethash indexed bytes32)
func (_SecretRegistry *SecretRegistryFilterer) WatchSecretRevealed(opts *bind.WatchOpts, sink chan<- *SecretRegistrySecretRevealed, secrethash [][32]byte) (event.Subscription, error) {

	var secrethashRule []interface{}
	for _, secrethashItem := range secrethash {
		secrethashRule = append(secrethashRule, secrethashItem)
	}

	logs, sub, err := _SecretRegistry.contract.WatchLogs(opts, "SecretRevealed", secrethashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SecretRegistrySecretRevealed)
				if err := _SecretRegistry.contract.UnpackLog(event, "SecretRevealed", log); err != nil {
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

// TokenNetworkABI is the input ABI used to generate the binding from.
const TokenNetworkABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"secret_registry\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"chain_id\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"participant\",\"type\":\"address\"},{\"name\":\"partner\",\"type\":\"address\"},{\"name\":\"transferered_amount\",\"type\":\"uint256\"},{\"name\":\"locksroot\",\"type\":\"bytes32\"},{\"name\":\"merkle_tree_leaves\",\"type\":\"bytes\"}],\"name\":\"unlock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"beneficiary\",\"type\":\"address\"},{\"name\":\"cheater\",\"type\":\"address\"},{\"name\":\"lockhash\",\"type\":\"bytes32\"},{\"name\":\"beneficiary_transferred_amount\",\"type\":\"uint256\"},{\"name\":\"additional_hash\",\"type\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\"},{\"name\":\"merkle_proof\",\"type\":\"bytes\"}],\"name\":\"punishObsoleteUnlock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"contract_address\",\"type\":\"address\"}],\"name\":\"contractExists\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"channels\",\"outputs\":[{\"name\":\"settle_block_number\",\"type\":\"uint64\"},{\"name\":\"open_blocknumber\",\"type\":\"uint64\"},{\"name\":\"state\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"invalid_balance_hash\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes24\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"participant1_address\",\"type\":\"address\"},{\"name\":\"participant1_balance\",\"type\":\"uint256\"},{\"name\":\"participant2_address\",\"type\":\"address\"},{\"name\":\"participant2_balance\",\"type\":\"uint256\"},{\"name\":\"participant1_signature\",\"type\":\"bytes\"},{\"name\":\"participant2_signature\",\"type\":\"bytes\"}],\"name\":\"cooperativeSettle\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"punish_block_number\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"participant1\",\"type\":\"address\"},{\"name\":\"participant1_balance\",\"type\":\"uint256\"},{\"name\":\"participant1_withdraw\",\"type\":\"uint256\"},{\"name\":\"participant2\",\"type\":\"address\"},{\"name\":\"participant2_balance\",\"type\":\"uint256\"},{\"name\":\"participant2_withdraw\",\"type\":\"uint256\"},{\"name\":\"participant1_signature\",\"type\":\"bytes\"},{\"name\":\"participant2_signature\",\"type\":\"bytes\"}],\"name\":\"withDraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"participant\",\"type\":\"address\"},{\"name\":\"transferred_amount\",\"type\":\"uint256\"},{\"name\":\"locksroot\",\"type\":\"bytes32\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"additional_hash\",\"type\":\"bytes32\"},{\"name\":\"participant_signature\",\"type\":\"bytes\"}],\"name\":\"updateBalanceProof\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"participant\",\"type\":\"address\"},{\"name\":\"partner\",\"type\":\"address\"}],\"name\":\"getChannelParticipantInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes24\"},{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"participant1\",\"type\":\"address\"},{\"name\":\"participant2\",\"type\":\"address\"},{\"name\":\"settle_timeout\",\"type\":\"uint64\"}],\"name\":\"openChannel\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"partner\",\"type\":\"address\"},{\"name\":\"transferred_amount\",\"type\":\"uint256\"},{\"name\":\"locksroot\",\"type\":\"bytes32\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"additional_hash\",\"type\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"closeChannel\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"participant\",\"type\":\"address\"},{\"name\":\"partner\",\"type\":\"address\"},{\"name\":\"total_deposit\",\"type\":\"uint256\"}],\"name\":\"setTotalDeposit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"participant1\",\"type\":\"address\"},{\"name\":\"participant1_transferred_amount\",\"type\":\"uint256\"},{\"name\":\"participant1_locksroot\",\"type\":\"bytes32\"},{\"name\":\"participant2\",\"type\":\"address\"},{\"name\":\"participant2_transferred_amount\",\"type\":\"uint256\"},{\"name\":\"participant2_locksroot\",\"type\":\"bytes32\"}],\"name\":\"settleChannel\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"participant\",\"type\":\"address\"},{\"name\":\"partner\",\"type\":\"address\"},{\"name\":\"transferred_amount\",\"type\":\"uint256\"},{\"name\":\"locksroot\",\"type\":\"bytes32\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"additional_hash\",\"type\":\"bytes32\"},{\"name\":\"participant_signature\",\"type\":\"bytes\"},{\"name\":\"partner_signature\",\"type\":\"bytes\"}],\"name\":\"updateBalanceProofDelegate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"participant1\",\"type\":\"address\"},{\"name\":\"participant2\",\"type\":\"address\"}],\"name\":\"getChannelInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"uint64\"},{\"name\":\"\",\"type\":\"uint64\"},{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"participant\",\"type\":\"address\"},{\"name\":\"partner\",\"type\":\"address\"},{\"name\":\"settle_timeout\",\"type\":\"uint64\"},{\"name\":\"deposit\",\"type\":\"uint256\"}],\"name\":\"openChannelWithDeposit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_token_address\",\"type\":\"address\"},{\"name\":\"_secret_registry\",\"type\":\"address\"},{\"name\":\"_chain_id\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channel_identifier\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"participant1\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"participant2\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"settle_timeout\",\"type\":\"uint256\"}],\"name\":\"ChannelOpened\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channel_identifier\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"participant\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"total_deposit\",\"type\":\"uint256\"}],\"name\":\"ChannelNewDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channel_identifier\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"closing_participant\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"locksroot\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"transferred_amount\",\"type\":\"uint256\"}],\"name\":\"ChannelClosed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channel_identifier\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"payer_participant\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"locskroot\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"transferred_amount\",\"type\":\"uint256\"}],\"name\":\"ChannelUnlocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channel_identifier\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"participant\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"locksroot\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"transferred_amount\",\"type\":\"uint256\"}],\"name\":\"BalanceProofUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channel_identifier\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"participant1_amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"participant2_amount\",\"type\":\"uint256\"}],\"name\":\"ChannelSettled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channel_identifier\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"participant1_amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"participant2_amount\",\"type\":\"uint256\"}],\"name\":\"ChannelCooperativeSettled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channel_identifier\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"participant1\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"participant1_balance\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"participant2\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"participant2_balance\",\"type\":\"uint256\"}],\"name\":\"ChannelWithdraw\",\"type\":\"event\"}]"

// TokenNetworkBin is the compiled bytecode used for deploying new contracts.
const TokenNetworkBin = `0x60806040523480156200001157600080fd5b506040516060806200314c833981016040908152815160208301519190920151600160a060020a03831615156200004757600080fd5b600160a060020a03821615156200005d57600080fd5b600081116200006b57600080fd5b6200007f8364010000000062000177810204565b15156200008b57600080fd5b6200009f8264010000000062000177810204565b1515620000ab57600080fd5b60008054600160a060020a03808616600160a060020a031992831617808455600180548784169416939093179092556002849055604080517f18160ddd000000000000000000000000000000000000000000000000000000008152905192909116916318160ddd9160048082019260209290919082900301818787803b1580156200013557600080fd5b505af11580156200014a573d6000803e3d6000fd5b505050506040513d60208110156200016157600080fd5b5051116200016e57600080fd5b5050506200017f565b6000903b1190565b612fbd806200018f6000396000f3006080604052600436106101035763ffffffff60e060020a60003504166324d73a9381146101085780633af973b1146101395780634e485437146101605780636b51403c146101d95780637709bc78146102945780637a7ebd7b146102c95780637ed74ad91461030f5780638568536a146103425780639375cff2146103f95780639bc6cb721461042b578063aaa3dbcc146104e8578063ac13370914610567578063aef91441146105c0578063b32c65c8146105f4578063b9eec0141461067e578063c10fd1bb146106fd578063e11cbf9914610727578063f8658b251461075c578063f94c9e1314610826578063fc0c546a14610882578063fc65697014610897575b600080fd5b34801561011457600080fd5b5061011d6108ce565b60408051600160a060020a039092168252519081900360200190f35b34801561014557600080fd5b5061014e6108dd565b60408051918252519081900360200190f35b34801561016c57600080fd5b50604080516020601f6084356004818101359283018490048402850184019095528184526101d794600160a060020a0381358116956024803590921695604435956064359536959460a494939101919081908401838280828437509497506108e39650505050505050565b005b3480156101e557600080fd5b50604080516020600460a43581810135601f81018490048402850184019095528484526101d7948235600160a060020a039081169560248035909216956044359560643595608435953695929460c494909390920191819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a999881019791965091820194509250829150840183828082843750949750610aff9650505050505050565b3480156102a057600080fd5b506102b5600160a060020a0360043516610c69565b604080519115158252519081900360200190f35b3480156102d557600080fd5b506102e1600435610c71565b6040805167ffffffffffffffff948516815292909316602083015260ff168183015290519081900360600190f35b34801561031b57600080fd5b50610324610ca5565b6040805167ffffffffffffffff199092168252519081900360200190f35b34801561034e57600080fd5b50604080516020601f6084356004818101359283018490048402850184019095528184526101d794600160a060020a0381358116956024803596604435909316956064359536959460a49493919091019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a999881019791965091820194509250829150840183828082843750949750610cc19650505050505050565b34801561040557600080fd5b5061040e610fa7565b6040805167ffffffffffffffff9092168252519081900360200190f35b34801561043757600080fd5b50604080516020601f60c4356004818101359283018490048402850184019095528184526101d794600160a060020a038135811695602480359660443596606435909416956084359560a435953695919460e49492930191819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a999881019791965091820194509250829150840183828082843750949750610fac9650505050505050565b3480156104f457600080fd5b50604080516020600460a43581810135601f81018490048402850184019095528484526101d7948235600160a060020a03169460248035956044359560643567ffffffffffffffff1695608435953695929460c494920191819084018382808284375094975061153c9650505050505050565b34801561057357600080fd5b5061058e600160a060020a03600435811690602435166116bc565b6040805193845267ffffffffffffffff19909216602084015267ffffffffffffffff1682820152519081900360600190f35b3480156105cc57600080fd5b506101d7600160a060020a036004358116906024351667ffffffffffffffff60443516611723565b34801561060057600080fd5b5061060961187f565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561064357818101518382015260200161062b565b50505050905090810190601f1680156106705780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561068a57600080fd5b50604080516020600460a43581810135601f81018490048402850184019095528484526101d7948235600160a060020a03169460248035956044359560643567ffffffffffffffff1695608435953695929460c49492019181908401838280828437509497506118b69650505050505050565b34801561070957600080fd5b506101d7600160a060020a0360043581169060243516604435611a48565b34801561073357600080fd5b506101d7600160a060020a0360043581169060243590604435906064351660843560a435611bbe565b34801561076857600080fd5b50604080516020601f60c4356004818101359283018490048402850184019095528184526101d794600160a060020a0381358116956024803590921695604435956064359567ffffffffffffffff608435169560a435953695919460e49492939091019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a999881019791965091820194509250829150840183828082843750949750611ed39650505050505050565b34801561083257600080fd5b5061084d600160a060020a03600435811690602435166120b8565b6040805194855267ffffffffffffffff9384166020860152919092168382015260ff9091166060830152519081900360800190f35b34801561088e57600080fd5b5061011d612109565b3480156108a357600080fd5b506101d7600160a060020a036004358116906024351667ffffffffffffffff60443516606435612118565b600154600160a060020a031681565b60025481565b60008060008060008060008088511115156108fd57600080fd5b6109078c8c61234d565b965060036000886000191660001916815260200190815260200160002091508160010160008d600160a060020a0316600160a060020a031681526020019081526020016000209050438260000160009054906101000a900467ffffffffffffffff1667ffffffffffffffff161015151561098057600080fd5b8154608060020a900460ff1660021461099857600080fd5b6001810154604060020a0292508815156109b157600080fd5b60018101546040805167ffffffffffffffff1986168152601881018a905267ffffffffffffffff60c060020a93849004169092026038830152805191829003810190912060008181526004602052919091205490965060ff1615610a1457600080fd5b6000868152600460205260409020805460ff19166001179055610a3688612498565b909550935060008411610a4857600080fd5b848914610a5457600080fd5b610a5e8a8a6127ec565b67ffffffffffffffff19848116911614610a7757600080fd5b98830198610a858a8a6127ec565b600182018054604060020a909204600160c060020a031990921691909117905560408051600160a060020a038e168152602081018790528082018c9052905188917f9e3b094fde58f3a83bd8b77d0a995fdb71f3169c6fa7e6d386e9f5902841e5ff919081900360600190a2505050505050505050505050565b600080600080600080610b128d8d61234d565b600081815260036020526040902080549197509350608060020a900460ff16600214610b3d57600080fd5b600160a060020a038d1660009081526001848101602052604090912090810154604060020a0267ffffffffffffffff191694509150831515610b7e57600080fd5b8254610ba09087908d90604060020a900467ffffffffffffffff168c8c612825565b600160a060020a038d8116911614610bb757600080fd5b50600160a060020a038b1660009081526001830160205260409020610bdc8b886128e7565b9450610be88a866127ec565b67ffffffffffffffff19168414610bfe57600080fd5b604080516000808252602082018190528251918290039092019020600184018054604060020a909204600160c060020a03199283161777ffffffffffffffffffffffffffffffffffffffffffffffff1690911790558154835401909255555050505050505050505050565b6000903b1190565b60036020526000908152604090205467ffffffffffffffff80821691604060020a810490911690608060020a900460ff1683565b6040805160008082526020820152815190819003909101902081565b6000806000806000806000610cd68d8c61234d565b600081815260036020526040902080549196509350608060020a900460ff16600114610d0157600080fd5b8254604060020a900467ffffffffffffffff169350610d25858e8e8e8e898f612a38565b9650600160a060020a038d811690881614610d3f57600080fd5b610d4e858e8e8e8e898e612a38565b9650600160a060020a038b811690881614610d6857600080fd5b5050600160a060020a03808c166000908152600180840160209081526040808420948e168452808420805486548688558786018790558683559482018690558986526003909352908420805470ffffffffffffffffffffffffffffffffff1916905591019650908c1115610e84576000809054906101000a9004600160a060020a0316600160a060020a031663a9059cbb8e8e6040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b158015610e4d57600080fd5b505af1158015610e61573d6000803e3d6000fd5b505050506040513d6020811015610e7757600080fd5b50511515610e8457600080fd5b60008a1115610f365760008054604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a038f81166004830152602482018f90529151919092169263a9059cbb92604480820193602093909283900390910190829087803b158015610eff57600080fd5b505af1158015610f13573d6000803e3d6000fd5b505050506040513d6020811015610f2957600080fd5b50511515610f3657600080fd5b8b8a018614610f4457600080fd5b8b861015610f5157600080fd5b89861015610f5e57600080fd5b604080518d8152602081018c9052815187927ffb2f4bc0fb2e0f1001f78d15e81a2e1981f262d31e8bd72309e26cc63bf7bb02928290030190a250505050505050505050505050565b600a81565b600080600080600080610fbf8e8c61234d565b600081815260036020526040902080549196509450608060020a900460ff16600114610fea57600080fd5b8d8d8c8c8f898960000160089054906101000a900467ffffffffffffffff166002546040516020018089600160a060020a0316600160a060020a03166c0100000000000000000000000002815260140188815260200187600160a060020a0316600160a060020a03166c0100000000000000000000000002815260140186815260200185815260200184600019166000191681526020018367ffffffffffffffff1667ffffffffffffffff1660c060020a028152600801828152602001985050505050505050506040516020818303038152906040526040518082805190602001908083835b602083106110ef5780518252601f1990920191602091820191016110d0565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902092506111288389612b2d565b600160a060020a038f811691161461113f57600080fd5b8d8d8c8c8f8d8a8a60000160089054906101000a900467ffffffffffffffff16600254604051602001808a600160a060020a0316600160a060020a03166c0100000000000000000000000002815260140189815260200188600160a060020a0316600160a060020a03166c0100000000000000000000000002815260140187815260200186815260200185815260200184600019166000191681526020018367ffffffffffffffff1667ffffffffffffffff1660c060020a02815260080182815260200199505050505050505050506040516020818303038152906040526040518082805190602001908083835b6020831061124c5780518252601f19909201916020918201910161122d565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902092506112858388612b2d565b600160a060020a038c811691161461129c57600080fd5b5050600160a060020a03808d166000908152600184016020526040808220928c168252902080548254019550858d11156112d557600080fd5b858a11156112e257600080fd5b8c8a0186146112f057600080fd5b60008c11156113a7576000809054906101000a9004600160a060020a0316600160a060020a031663a9059cbb8f8e6040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b15801561137057600080fd5b505af1158015611384573d6000803e3d6000fd5b505050506040513d602081101561139a57600080fd5b505115156113a757600080fd5b60008911156114595760008054604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a038f81166004830152602482018e90529151919092169263a9059cbb92604480820193602093909283900390910190829087803b15801561142257600080fd5b505af1158015611436573d6000803e3d6000fd5b505050506040513d602081101561144c57600080fd5b5051151561145957600080fd5b8c8c111561146657600080fd5b8989111561147357600080fd5b8b8d038260000181905550888a038160000181905550438460000160086101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555084600019167fdc5ff4ab383e66679a382f376c0e80534f51f3f3a398add646422cd81f5f815d8f8f8e8e6040518085600160a060020a0316600160a060020a0316815260200184815260200183600160a060020a0316600160a060020a0316815260200182815260200194505050505060405180910390a25050505050505050505050505050565b600080600061154b893361234d565b6000818152600360209081526040808320600160a060020a038e1684526001810190925290912081549295509093509150608060020a900460ff1660021461159257600080fd5b81544367ffffffffffffffff90911610156115ac57600080fd5b600181015467ffffffffffffffff60c060020a9091048116908716116115d157600080fd5b6115f8838989898660000160089054906101000a900467ffffffffffffffff168a8a612c0d565b600160a060020a038a811691161461160f57600080fd5b61161988886127ec565b60018201805467ffffffffffffffff891660c060020a02604060020a909304600160c060020a03199091161777ffffffffffffffffffffffffffffffffffffffffffffffff1691909117905560408051600160a060020a038b168152602081018990528082018a9052905184917f910c9237f4197a18340110a181e8fb775496506a007a94b46f9f80f2a35918f9919081900360600190a2505050505050505050565b6000806000806000806116cf888861234d565b6000908152600360209081526040808320600160a060020a039b909b16835260019a8b01909152902080549801549798604060020a89029860c060020a900467ffffffffffffffff16975095505050505050565b6000808260068167ffffffffffffffff161015801561174f5750622932e08167ffffffffffffffff1611155b151561175a57600080fd5b600160a060020a038616151561176f57600080fd5b600160a060020a038516151561178457600080fd5b600160a060020a03868116908616141561179d57600080fd5b6117a7868661234d565b600081815260036020526040902080549194509250608060020a900460ff16156117d057600080fd5b815470ff00000000000000000000000000000000194367ffffffffffffffff908116604060020a026fffffffffffffffff00000000000000001991881667ffffffffffffffff199094168417919091161716608060020a17835560408051600160a060020a03808a16825288166020820152808201929092525184917f448d27f1fe12f92a2070111296e68fd6ef0a01c0e05bf5819eda0dbcf267bf3d919081900360600190a2505050505050565b60408051808201909152600581527f302e332e5f000000000000000000000000000000000000000000000000000000602082015281565b6000806000806118c6338b61234d565b600081815260036020526040902080549195509250608060020a900460ff166001146118f157600080fd5b815467ffffffffffffffff1970ff0000000000000000000000000000000019909116700200000000000000000000000000000000179081164367ffffffffffffffff928316018216178355600090881611156119f95750600160a060020a0389166000908152600182016020526040902081546119889085908b908b908b90604060020a900467ffffffffffffffff168b8b612c0d565b9250600160a060020a038a8116908416146119a257600080fd5b6119ac89896127ec565b60018201805467ffffffffffffffff8a1660c060020a02604060020a909304600160c060020a03199091161777ffffffffffffffffffffffffffffffffffffffffffffffff169190911790555b60408051338152602081018a90528082018b9052905185917f69610baaace24c039f891a11b42c0b1df1496ab0db38b0c4ee4ed33d6d53da1a919081900360600190a250505050505050505050565b600080808080808611611a5a57600080fd5b611a64888861234d565b6000818152600360209081526040808320600160a060020a038d168452600181019092529091208054965091945092509050858410611aa257600080fd5b838603808501825560008054604080517f23b872dd000000000000000000000000000000000000000000000000000000008152336004820152306024820152604481018590529051939850600160a060020a03909116926323b872dd92606480840193602093929083900390910190829087803b158015611b2257600080fd5b505af1158015611b36573d6000803e3d6000fd5b505050506040513d6020811015611b4c57600080fd5b50511515611b5957600080fd5b8154608060020a900460ff16600114611b7157600080fd5b60408051600160a060020a038a16815260208101889052815185927f0346e981e2bfa2366dc2307a8f1fa24779830a01121b1275fe565c6b98bb4d34928290030190a25050505050505050565b600080600080600080611bd18c8a61234d565b600081815260036020526040902080549195509350608060020a900460ff16600214611bfc57600080fd5b82544367ffffffffffffffff918216600a0190911610611c1b57600080fd5b5050600160a060020a03808b166000908152600183016020526040808220928a1682529020611c4a8b8b6127ec565b6001830154604060020a0267ffffffffffffffff19908116911614611c6e57600080fd5b611c7888886127ec565b6001820154604060020a0267ffffffffffffffff19908116911614611c9c57600080fd5b805482548981018d81039850910195508b1115611cb857600095505b611cc28686612c9e565b600160a060020a03808e1660009081526001808701602090815260408084208481558301849055938e16835283832083815590910182905587825260039052908120805470ffffffffffffffffffffffffffffffffff191690558187039950909650861115611dd9576000809054906101000a9004600160a060020a0316600160a060020a031663a9059cbb8d886040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b158015611da257600080fd5b505af1158015611db6573d6000803e3d6000fd5b505050506040513d6020811015611dcc57600080fd5b50511515611dd957600080fd5b6000881115611e8b5760008054604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a038d81166004830152602482018d90529151919092169263a9059cbb92604480820193602093909283900390910190829087803b158015611e5457600080fd5b505af1158015611e68573d6000803e3d6000fd5b505050506040513d6020811015611e7e57600080fd5b50511515611e8b57600080fd5b60408051878152602081018a9052815186927ff94fb5c0628a82dc90648e8dc5e983f632633b0d26603d64e8cc042ca0790aa4928290030190a2505050505050505050505050565b6000806000611ee28b8b61234d565b925060036000846000191660001916815260200190815260200160002091508160010160008c600160a060020a0316600160a060020a0316815260200190815260200160002090508160000160109054906101000a900460ff1660ff166002141515611f4d57600080fd5b81544367ffffffffffffffff9091161015611f6757600080fd5b600181015467ffffffffffffffff60c060020a909104811690881611611f8c57600080fd5b611fb4838a8a8a8660000160089054906101000a900467ffffffffffffffff168b8b8b612cb6565b600160a060020a038b8116911614611fcb57600080fd5b611ff2838a8a8a8660000160089054906101000a900467ffffffffffffffff168b8b612c0d565b600160a060020a038c811691161461200957600080fd5b61201389896127ec565b60018201805467ffffffffffffffff8a1660c060020a02604060020a909304600160c060020a03199091161777ffffffffffffffffffffffffffffffffffffffffffffffff1691909117905560408051600160a060020a038d168152602081018a90528082018b9052905184917f910c9237f4197a18340110a181e8fb775496506a007a94b46f9f80f2a35918f9919081900360600190a25050505050505050505050565b6000806000806000806120cb888861234d565b600081815260036020526040902054909967ffffffffffffffff8083169a50604060020a8304169850608060020a90910460ff169650945050505050565b600054600160a060020a031681565b60008060008460068167ffffffffffffffff16101580156121465750622932e08167ffffffffffffffff1611155b151561215157600080fd5b600160a060020a038816151561216657600080fd5b600160a060020a038716151561217b57600080fd5b600160a060020a03888116908816141561219457600080fd5b61219e888861234d565b6000818152600360209081526040808320600160a060020a038d1684526001810190925290912081549296509094509250608060020a900460ff16156121e357600080fd5b825470ff00000000000000000000000000000000194367ffffffffffffffff908116604060020a026fffffffffffffffff000000000000000019918a1667ffffffffffffffff1990941693909317169190911716608060020a17835560008054604080517f23b872dd000000000000000000000000000000000000000000000000000000008152336004820152306024820152604481018990529051600160a060020a03909216926323b872dd926064808401936020939083900390910190829087803b1580156122b357600080fd5b505af11580156122c7573d6000803e3d6000fd5b505050506040513d60208110156122dd57600080fd5b505115156122ea57600080fd5b84825560408051600160a060020a03808b1682528916602082015267ffffffffffffffff881681830152905185917f448d27f1fe12f92a2070111296e68fd6ef0a01c0e05bf5819eda0dbcf267bf3d919081900360600190a25050505050505050565b600081600160a060020a031683600160a060020a031610156124175760408051600160a060020a038581166c0100000000000000000000000090810260208085019190915291861681026034840152300260488301528251808303603c018152605c90920192839052815191929182918401908083835b602083106123e35780518252601f1990920191602091820191016123c4565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390209050612492565b604080516c01000000000000000000000000600160a060020a03808616820260208085019190915290871682026034840152309190910260488301528251603c818403018152605c9092019283905281519192918291840190808383602083106123e35780518252601f1990920191602091820191016123c4565b92915050565b8051600090819081808080806060808706156124b357600080fd5b606087046001016040519080825280602002602001820160405280156124e3578160200160208202803883390190505b509050602095505b8686101561252d576124fd8a87612e1c565b958601959450925082816060880481518110151561251757fe5b60209081029091010152606095909501946124eb565b6060870496505b60018711156127c357600287061561258157806001880381518110151561255757fe5b90602001906020020151818881518110151561256f57fe5b60209081029091010152600196909601955b600095505b600187038610156127b85780866001018151811015156125a257fe5b6020908102909101015181518290889081106125ba57fe5b6020908102909101015114156125e95780868151811015156125d857fe5b906020019060200201519250612790565b80866001018151811015156125fa57fe5b60209081029091010151815182908890811061261257fe5b6020908102909101015110156126db57808681518110151561263057fe5b90602001906020020151818760010181518110151561264b57fe5b6020908102909101810151604080518084019490945283810191909152805180840382018152606090930190819052825190918291908401908083835b602083106126a75780518252601f199092019160209182019101612688565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390209250612790565b80866001018151811015156126ec57fe5b90602001906020020151818781518110151561270457fe5b6020908102909101810151604080518084019490945283810191909152805180840382018152606090930190819052825190918291908401908083835b602083106127605780518252601f199092019160209182019101612741565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902092505b8281600288048151811015156127a257fe5b6020908102909101015260029590950194612586565b600286049650612534565b8060008151811015156127d257fe5b602090810290910101519a94995093975050505050505050565b6000811580156127fa575082155b1561280757506000612492565b50604080519182526020820192909252815190819003909101902090565b60025460408051602080820188905281830189905260c060020a67ffffffffffffffff8816026060830152606882019390935260888082018690528251808303909101815260a890910191829052805160009384939182918401908083835b602083106128a35780518252601f199092019160209182019101612884565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902090506128dc8184612b2d565b979650505050505050565b6000806000602084518115156128f957fe5b061561290457600080fd5b602091505b83518211612a2f575082810151808510156129a357604080516020808201889052818301849052825180830384018152606090920192839052815191929182918401908083835b6020831061296f5780518252601f199092019160209182019101612950565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390209450612a24565b604080516020808201849052818301889052825180830384018152606090920192839052815191929182918401908083835b602083106129f45780518252601f1990920191602091820191016129d5565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902094505b602082019150612909565b50929392505050565b60025460408051600160a060020a038981166c01000000000000000000000000908102602080850191909152603484018b905291891602605483015260688201879052608882018b905267ffffffffffffffff861660c060020a0260a883015260b0808301949094528251808303909401845260d090910191829052825160009384939092909182918401908083835b60208310612ae75780518252601f199092019160209182019101612ac8565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390209050612b208184612b2d565b9998505050505050505050565b60008060008084516041141515612b4357600080fd5b50505060208201516040830151606084015160001a601b60ff82161015612b6857601b015b8060ff16601b1480612b7d57508060ff16601c145b1515612b8857600080fd5b60408051600080825260208083018085528a905260ff8516838501526060830187905260808301869052925160019360a0808501949193601f19840193928390039091019190865af1158015612be2573d6000803e3d6000fd5b5050604051601f190151945050600160a060020a0384161515612c0457600080fd5b50505092915050565b6002546040805160208082018a905281830189905260c060020a67ffffffffffffffff808a168202606085015260688401889052608884018d905288160260a883015260b0808301949094528251808303909401845260d0909101918290528251600093849390929091829184019080838360208310612ae75780518252601f199092019160209182019101612ac8565b6000818311612cad5782612caf565b815b9392505050565b600080888888878d8a6002548a6040516020018089815260200188600019166000191681526020018767ffffffffffffffff1667ffffffffffffffff1660c060020a028152600801866000191660001916815260200185600019166000191681526020018467ffffffffffffffff1667ffffffffffffffff1660c060020a02815260080183815260200182805190602001908083835b60208310612d6b5780518252601f199092019160209182019101612d4c565b6001836020036101000a038019825116818451168082178552505050505050905001985050505050505050506040516020818303038152906040526040518082805190602001908083835b60208310612dd55780518252601f199092019160209182019101612db6565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390209050612e0e8184612b2d565b9a9950505050505050505050565b6000806000806000806000878951111515612e3d5795506000945085612f85565b888801805160208083015160409384015184518084018590528086018390526060808201839052865180830390910181526080909101958690528051949a509198509550929182918401908083835b60208310612eab5780518252601f199092019160209182019101612e8c565b51815160209384036101000a6000190180199092169116179052604080519290940182900382206001547fc1f62946000000000000000000000000000000000000000000000000000000008452600484018a90529451909750600160a060020a03909416955063c1f62946945060248083019491935090918290030181600087803b158015612f3957600080fd5b505af1158015612f4d573d6000803e3d6000fd5b505050506040513d6020811015612f6357600080fd5b50519250821580612f745750828511155b15612f7e57600093505b8084965096505b505050505092509290505600a165627a7a72305820f3d21359b18fd00cd7455b2254197313cdf0f9c789d88ce4e8ae0bae4b0f17370029`

// DeployTokenNetwork deploys a new Ethereum contract, binding an instance of TokenNetwork to it.
func DeployTokenNetwork(auth *bind.TransactOpts, backend bind.ContractBackend, _token_address common.Address, _secret_registry common.Address, _chain_id *big.Int) (common.Address, *types.Transaction, *TokenNetwork, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenNetworkABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TokenNetworkBin), backend, _token_address, _secret_registry, _chain_id)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TokenNetwork{TokenNetworkCaller: TokenNetworkCaller{contract: contract}, TokenNetworkTransactor: TokenNetworkTransactor{contract: contract}, TokenNetworkFilterer: TokenNetworkFilterer{contract: contract}}, nil
}

// TokenNetwork is an auto generated Go binding around an Ethereum contract.
type TokenNetwork struct {
	TokenNetworkCaller     // Read-only binding to the contract
	TokenNetworkTransactor // Write-only binding to the contract
	TokenNetworkFilterer   // Log filterer for contract events
}

// TokenNetworkCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenNetworkCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenNetworkTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenNetworkTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenNetworkFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenNetworkFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenNetworkSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenNetworkSession struct {
	Contract     *TokenNetwork     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenNetworkCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenNetworkCallerSession struct {
	Contract *TokenNetworkCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// TokenNetworkTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenNetworkTransactorSession struct {
	Contract     *TokenNetworkTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// TokenNetworkRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenNetworkRaw struct {
	Contract *TokenNetwork // Generic contract binding to access the raw methods on
}

// TokenNetworkCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenNetworkCallerRaw struct {
	Contract *TokenNetworkCaller // Generic read-only contract binding to access the raw methods on
}

// TokenNetworkTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenNetworkTransactorRaw struct {
	Contract *TokenNetworkTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenNetwork creates a new instance of TokenNetwork, bound to a specific deployed contract.
func NewTokenNetwork(address common.Address, backend bind.ContractBackend) (*TokenNetwork, error) {
	contract, err := bindTokenNetwork(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenNetwork{TokenNetworkCaller: TokenNetworkCaller{contract: contract}, TokenNetworkTransactor: TokenNetworkTransactor{contract: contract}, TokenNetworkFilterer: TokenNetworkFilterer{contract: contract}}, nil
}

// NewTokenNetworkCaller creates a new read-only instance of TokenNetwork, bound to a specific deployed contract.
func NewTokenNetworkCaller(address common.Address, caller bind.ContractCaller) (*TokenNetworkCaller, error) {
	contract, err := bindTokenNetwork(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkCaller{contract: contract}, nil
}

// NewTokenNetworkTransactor creates a new write-only instance of TokenNetwork, bound to a specific deployed contract.
func NewTokenNetworkTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenNetworkTransactor, error) {
	contract, err := bindTokenNetwork(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkTransactor{contract: contract}, nil
}

// NewTokenNetworkFilterer creates a new log filterer instance of TokenNetwork, bound to a specific deployed contract.
func NewTokenNetworkFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenNetworkFilterer, error) {
	contract, err := bindTokenNetwork(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkFilterer{contract: contract}, nil
}

// bindTokenNetwork binds a generic wrapper to an already deployed contract.
func bindTokenNetwork(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenNetworkABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenNetwork *TokenNetworkRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenNetwork.Contract.TokenNetworkCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenNetwork *TokenNetworkRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenNetwork.Contract.TokenNetworkTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenNetwork *TokenNetworkRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenNetwork.Contract.TokenNetworkTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenNetwork *TokenNetworkCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenNetwork.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenNetwork *TokenNetworkTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenNetwork.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenNetwork *TokenNetworkTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenNetwork.Contract.contract.Transact(opts, method, params...)
}

// Chain_id is a free data retrieval call binding the contract method 0x3af973b1.
//
// Solidity: function chain_id() constant returns(uint256)
func (_TokenNetwork *TokenNetworkCaller) Chain_id(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenNetwork.contract.Call(opts, out, "chain_id")
	return *ret0, err
}

// Chain_id is a free data retrieval call binding the contract method 0x3af973b1.
//
// Solidity: function chain_id() constant returns(uint256)
func (_TokenNetwork *TokenNetworkSession) Chain_id() (*big.Int, error) {
	return _TokenNetwork.Contract.Chain_id(&_TokenNetwork.CallOpts)
}

// Chain_id is a free data retrieval call binding the contract method 0x3af973b1.
//
// Solidity: function chain_id() constant returns(uint256)
func (_TokenNetwork *TokenNetworkCallerSession) Chain_id() (*big.Int, error) {
	return _TokenNetwork.Contract.Chain_id(&_TokenNetwork.CallOpts)
}

// Channels is a free data retrieval call binding the contract method 0x7a7ebd7b.
//
// Solidity: function channels( bytes32) constant returns(settle_block_number uint64, open_blocknumber uint64, state uint8)
func (_TokenNetwork *TokenNetworkCaller) Channels(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Settle_block_number uint64
	Open_blocknumber    uint64
	State               uint8
}, error) {
	ret := new(struct {
		Settle_block_number uint64
		Open_blocknumber    uint64
		State               uint8
	})
	out := ret
	err := _TokenNetwork.contract.Call(opts, out, "channels", arg0)
	return *ret, err
}

// Channels is a free data retrieval call binding the contract method 0x7a7ebd7b.
//
// Solidity: function channels( bytes32) constant returns(settle_block_number uint64, open_blocknumber uint64, state uint8)
func (_TokenNetwork *TokenNetworkSession) Channels(arg0 [32]byte) (struct {
	Settle_block_number uint64
	Open_blocknumber    uint64
	State               uint8
}, error) {
	return _TokenNetwork.Contract.Channels(&_TokenNetwork.CallOpts, arg0)
}

// Channels is a free data retrieval call binding the contract method 0x7a7ebd7b.
//
// Solidity: function channels( bytes32) constant returns(settle_block_number uint64, open_blocknumber uint64, state uint8)
func (_TokenNetwork *TokenNetworkCallerSession) Channels(arg0 [32]byte) (struct {
	Settle_block_number uint64
	Open_blocknumber    uint64
	State               uint8
}, error) {
	return _TokenNetwork.Contract.Channels(&_TokenNetwork.CallOpts, arg0)
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(contract_address address) constant returns(bool)
func (_TokenNetwork *TokenNetworkCaller) ContractExists(opts *bind.CallOpts, contract_address common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TokenNetwork.contract.Call(opts, out, "contractExists", contract_address)
	return *ret0, err
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(contract_address address) constant returns(bool)
func (_TokenNetwork *TokenNetworkSession) ContractExists(contract_address common.Address) (bool, error) {
	return _TokenNetwork.Contract.ContractExists(&_TokenNetwork.CallOpts, contract_address)
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(contract_address address) constant returns(bool)
func (_TokenNetwork *TokenNetworkCallerSession) ContractExists(contract_address common.Address) (bool, error) {
	return _TokenNetwork.Contract.ContractExists(&_TokenNetwork.CallOpts, contract_address)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_TokenNetwork *TokenNetworkCaller) Contract_version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _TokenNetwork.contract.Call(opts, out, "contract_version")
	return *ret0, err
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_TokenNetwork *TokenNetworkSession) Contract_version() (string, error) {
	return _TokenNetwork.Contract.Contract_version(&_TokenNetwork.CallOpts)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_TokenNetwork *TokenNetworkCallerSession) Contract_version() (string, error) {
	return _TokenNetwork.Contract.Contract_version(&_TokenNetwork.CallOpts)
}

// GetChannelInfo is a free data retrieval call binding the contract method 0xf94c9e13.
//
// Solidity: function getChannelInfo(participant1 address, participant2 address) constant returns(bytes32, uint64, uint64, uint8)
func (_TokenNetwork *TokenNetworkCaller) GetChannelInfo(opts *bind.CallOpts, participant1 common.Address, participant2 common.Address) ([32]byte, uint64, uint64, uint8, error) {
	var (
		ret0 = new([32]byte)
		ret1 = new(uint64)
		ret2 = new(uint64)
		ret3 = new(uint8)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
	}
	err := _TokenNetwork.contract.Call(opts, out, "getChannelInfo", participant1, participant2)
	return *ret0, *ret1, *ret2, *ret3, err
}

// GetChannelInfo is a free data retrieval call binding the contract method 0xf94c9e13.
//
// Solidity: function getChannelInfo(participant1 address, participant2 address) constant returns(bytes32, uint64, uint64, uint8)
func (_TokenNetwork *TokenNetworkSession) GetChannelInfo(participant1 common.Address, participant2 common.Address) ([32]byte, uint64, uint64, uint8, error) {
	return _TokenNetwork.Contract.GetChannelInfo(&_TokenNetwork.CallOpts, participant1, participant2)
}

// GetChannelInfo is a free data retrieval call binding the contract method 0xf94c9e13.
//
// Solidity: function getChannelInfo(participant1 address, participant2 address) constant returns(bytes32, uint64, uint64, uint8)
func (_TokenNetwork *TokenNetworkCallerSession) GetChannelInfo(participant1 common.Address, participant2 common.Address) ([32]byte, uint64, uint64, uint8, error) {
	return _TokenNetwork.Contract.GetChannelInfo(&_TokenNetwork.CallOpts, participant1, participant2)
}

// GetChannelParticipantInfo is a free data retrieval call binding the contract method 0xac133709.
//
// Solidity: function getChannelParticipantInfo(participant address, partner address) constant returns(uint256, bytes24, uint64)
func (_TokenNetwork *TokenNetworkCaller) GetChannelParticipantInfo(opts *bind.CallOpts, participant common.Address, partner common.Address) (*big.Int, [24]byte, uint64, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new([24]byte)
		ret2 = new(uint64)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _TokenNetwork.contract.Call(opts, out, "getChannelParticipantInfo", participant, partner)
	return *ret0, *ret1, *ret2, err
}

// GetChannelParticipantInfo is a free data retrieval call binding the contract method 0xac133709.
//
// Solidity: function getChannelParticipantInfo(participant address, partner address) constant returns(uint256, bytes24, uint64)
func (_TokenNetwork *TokenNetworkSession) GetChannelParticipantInfo(participant common.Address, partner common.Address) (*big.Int, [24]byte, uint64, error) {
	return _TokenNetwork.Contract.GetChannelParticipantInfo(&_TokenNetwork.CallOpts, participant, partner)
}

// GetChannelParticipantInfo is a free data retrieval call binding the contract method 0xac133709.
//
// Solidity: function getChannelParticipantInfo(participant address, partner address) constant returns(uint256, bytes24, uint64)
func (_TokenNetwork *TokenNetworkCallerSession) GetChannelParticipantInfo(participant common.Address, partner common.Address) (*big.Int, [24]byte, uint64, error) {
	return _TokenNetwork.Contract.GetChannelParticipantInfo(&_TokenNetwork.CallOpts, participant, partner)
}

// Invalid_balance_hash is a free data retrieval call binding the contract method 0x7ed74ad9.
//
// Solidity: function invalid_balance_hash() constant returns(bytes24)
func (_TokenNetwork *TokenNetworkCaller) Invalid_balance_hash(opts *bind.CallOpts) ([24]byte, error) {
	var (
		ret0 = new([24]byte)
	)
	out := ret0
	err := _TokenNetwork.contract.Call(opts, out, "invalid_balance_hash")
	return *ret0, err
}

// Invalid_balance_hash is a free data retrieval call binding the contract method 0x7ed74ad9.
//
// Solidity: function invalid_balance_hash() constant returns(bytes24)
func (_TokenNetwork *TokenNetworkSession) Invalid_balance_hash() ([24]byte, error) {
	return _TokenNetwork.Contract.Invalid_balance_hash(&_TokenNetwork.CallOpts)
}

// Invalid_balance_hash is a free data retrieval call binding the contract method 0x7ed74ad9.
//
// Solidity: function invalid_balance_hash() constant returns(bytes24)
func (_TokenNetwork *TokenNetworkCallerSession) Invalid_balance_hash() ([24]byte, error) {
	return _TokenNetwork.Contract.Invalid_balance_hash(&_TokenNetwork.CallOpts)
}

// Punish_block_number is a free data retrieval call binding the contract method 0x9375cff2.
//
// Solidity: function punish_block_number() constant returns(uint64)
func (_TokenNetwork *TokenNetworkCaller) Punish_block_number(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _TokenNetwork.contract.Call(opts, out, "punish_block_number")
	return *ret0, err
}

// Punish_block_number is a free data retrieval call binding the contract method 0x9375cff2.
//
// Solidity: function punish_block_number() constant returns(uint64)
func (_TokenNetwork *TokenNetworkSession) Punish_block_number() (uint64, error) {
	return _TokenNetwork.Contract.Punish_block_number(&_TokenNetwork.CallOpts)
}

// Punish_block_number is a free data retrieval call binding the contract method 0x9375cff2.
//
// Solidity: function punish_block_number() constant returns(uint64)
func (_TokenNetwork *TokenNetworkCallerSession) Punish_block_number() (uint64, error) {
	return _TokenNetwork.Contract.Punish_block_number(&_TokenNetwork.CallOpts)
}

// Secret_registry is a free data retrieval call binding the contract method 0x24d73a93.
//
// Solidity: function secret_registry() constant returns(address)
func (_TokenNetwork *TokenNetworkCaller) Secret_registry(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenNetwork.contract.Call(opts, out, "secret_registry")
	return *ret0, err
}

// Secret_registry is a free data retrieval call binding the contract method 0x24d73a93.
//
// Solidity: function secret_registry() constant returns(address)
func (_TokenNetwork *TokenNetworkSession) Secret_registry() (common.Address, error) {
	return _TokenNetwork.Contract.Secret_registry(&_TokenNetwork.CallOpts)
}

// Secret_registry is a free data retrieval call binding the contract method 0x24d73a93.
//
// Solidity: function secret_registry() constant returns(address)
func (_TokenNetwork *TokenNetworkCallerSession) Secret_registry() (common.Address, error) {
	return _TokenNetwork.Contract.Secret_registry(&_TokenNetwork.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_TokenNetwork *TokenNetworkCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenNetwork.contract.Call(opts, out, "token")
	return *ret0, err
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_TokenNetwork *TokenNetworkSession) Token() (common.Address, error) {
	return _TokenNetwork.Contract.Token(&_TokenNetwork.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_TokenNetwork *TokenNetworkCallerSession) Token() (common.Address, error) {
	return _TokenNetwork.Contract.Token(&_TokenNetwork.CallOpts)
}

// CloseChannel is a paid mutator transaction binding the contract method 0xb9eec014.
//
// Solidity: function closeChannel(partner address, transferred_amount uint256, locksroot bytes32, nonce uint64, additional_hash bytes32, signature bytes) returns()
func (_TokenNetwork *TokenNetworkTransactor) CloseChannel(opts *bind.TransactOpts, partner common.Address, transferred_amount *big.Int, locksroot [32]byte, nonce uint64, additional_hash [32]byte, signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.contract.Transact(opts, "closeChannel", partner, transferred_amount, locksroot, nonce, additional_hash, signature)
}

// CloseChannel is a paid mutator transaction binding the contract method 0xb9eec014.
//
// Solidity: function closeChannel(partner address, transferred_amount uint256, locksroot bytes32, nonce uint64, additional_hash bytes32, signature bytes) returns()
func (_TokenNetwork *TokenNetworkSession) CloseChannel(partner common.Address, transferred_amount *big.Int, locksroot [32]byte, nonce uint64, additional_hash [32]byte, signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.CloseChannel(&_TokenNetwork.TransactOpts, partner, transferred_amount, locksroot, nonce, additional_hash, signature)
}

// CloseChannel is a paid mutator transaction binding the contract method 0xb9eec014.
//
// Solidity: function closeChannel(partner address, transferred_amount uint256, locksroot bytes32, nonce uint64, additional_hash bytes32, signature bytes) returns()
func (_TokenNetwork *TokenNetworkTransactorSession) CloseChannel(partner common.Address, transferred_amount *big.Int, locksroot [32]byte, nonce uint64, additional_hash [32]byte, signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.CloseChannel(&_TokenNetwork.TransactOpts, partner, transferred_amount, locksroot, nonce, additional_hash, signature)
}

// CooperativeSettle is a paid mutator transaction binding the contract method 0x8568536a.
//
// Solidity: function cooperativeSettle(participant1_address address, participant1_balance uint256, participant2_address address, participant2_balance uint256, participant1_signature bytes, participant2_signature bytes) returns()
func (_TokenNetwork *TokenNetworkTransactor) CooperativeSettle(opts *bind.TransactOpts, participant1_address common.Address, participant1_balance *big.Int, participant2_address common.Address, participant2_balance *big.Int, participant1_signature []byte, participant2_signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.contract.Transact(opts, "cooperativeSettle", participant1_address, participant1_balance, participant2_address, participant2_balance, participant1_signature, participant2_signature)
}

// CooperativeSettle is a paid mutator transaction binding the contract method 0x8568536a.
//
// Solidity: function cooperativeSettle(participant1_address address, participant1_balance uint256, participant2_address address, participant2_balance uint256, participant1_signature bytes, participant2_signature bytes) returns()
func (_TokenNetwork *TokenNetworkSession) CooperativeSettle(participant1_address common.Address, participant1_balance *big.Int, participant2_address common.Address, participant2_balance *big.Int, participant1_signature []byte, participant2_signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.CooperativeSettle(&_TokenNetwork.TransactOpts, participant1_address, participant1_balance, participant2_address, participant2_balance, participant1_signature, participant2_signature)
}

// CooperativeSettle is a paid mutator transaction binding the contract method 0x8568536a.
//
// Solidity: function cooperativeSettle(participant1_address address, participant1_balance uint256, participant2_address address, participant2_balance uint256, participant1_signature bytes, participant2_signature bytes) returns()
func (_TokenNetwork *TokenNetworkTransactorSession) CooperativeSettle(participant1_address common.Address, participant1_balance *big.Int, participant2_address common.Address, participant2_balance *big.Int, participant1_signature []byte, participant2_signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.CooperativeSettle(&_TokenNetwork.TransactOpts, participant1_address, participant1_balance, participant2_address, participant2_balance, participant1_signature, participant2_signature)
}

// OpenChannel is a paid mutator transaction binding the contract method 0xaef91441.
//
// Solidity: function openChannel(participant1 address, participant2 address, settle_timeout uint64) returns()
func (_TokenNetwork *TokenNetworkTransactor) OpenChannel(opts *bind.TransactOpts, participant1 common.Address, participant2 common.Address, settle_timeout uint64) (*types.Transaction, error) {
	return _TokenNetwork.contract.Transact(opts, "openChannel", participant1, participant2, settle_timeout)
}

// OpenChannel is a paid mutator transaction binding the contract method 0xaef91441.
//
// Solidity: function openChannel(participant1 address, participant2 address, settle_timeout uint64) returns()
func (_TokenNetwork *TokenNetworkSession) OpenChannel(participant1 common.Address, participant2 common.Address, settle_timeout uint64) (*types.Transaction, error) {
	return _TokenNetwork.Contract.OpenChannel(&_TokenNetwork.TransactOpts, participant1, participant2, settle_timeout)
}

// OpenChannel is a paid mutator transaction binding the contract method 0xaef91441.
//
// Solidity: function openChannel(participant1 address, participant2 address, settle_timeout uint64) returns()
func (_TokenNetwork *TokenNetworkTransactorSession) OpenChannel(participant1 common.Address, participant2 common.Address, settle_timeout uint64) (*types.Transaction, error) {
	return _TokenNetwork.Contract.OpenChannel(&_TokenNetwork.TransactOpts, participant1, participant2, settle_timeout)
}

// OpenChannelWithDeposit is a paid mutator transaction binding the contract method 0xfc656970.
//
// Solidity: function openChannelWithDeposit(participant address, partner address, settle_timeout uint64, deposit uint256) returns()
func (_TokenNetwork *TokenNetworkTransactor) OpenChannelWithDeposit(opts *bind.TransactOpts, participant common.Address, partner common.Address, settle_timeout uint64, deposit *big.Int) (*types.Transaction, error) {
	return _TokenNetwork.contract.Transact(opts, "openChannelWithDeposit", participant, partner, settle_timeout, deposit)
}

// OpenChannelWithDeposit is a paid mutator transaction binding the contract method 0xfc656970.
//
// Solidity: function openChannelWithDeposit(participant address, partner address, settle_timeout uint64, deposit uint256) returns()
func (_TokenNetwork *TokenNetworkSession) OpenChannelWithDeposit(participant common.Address, partner common.Address, settle_timeout uint64, deposit *big.Int) (*types.Transaction, error) {
	return _TokenNetwork.Contract.OpenChannelWithDeposit(&_TokenNetwork.TransactOpts, participant, partner, settle_timeout, deposit)
}

// OpenChannelWithDeposit is a paid mutator transaction binding the contract method 0xfc656970.
//
// Solidity: function openChannelWithDeposit(participant address, partner address, settle_timeout uint64, deposit uint256) returns()
func (_TokenNetwork *TokenNetworkTransactorSession) OpenChannelWithDeposit(participant common.Address, partner common.Address, settle_timeout uint64, deposit *big.Int) (*types.Transaction, error) {
	return _TokenNetwork.Contract.OpenChannelWithDeposit(&_TokenNetwork.TransactOpts, participant, partner, settle_timeout, deposit)
}

// PunishObsoleteUnlock is a paid mutator transaction binding the contract method 0x6b51403c.
//
// Solidity: function punishObsoleteUnlock(beneficiary address, cheater address, lockhash bytes32, beneficiary_transferred_amount uint256, additional_hash bytes32, signature bytes, merkle_proof bytes) returns()
func (_TokenNetwork *TokenNetworkTransactor) PunishObsoleteUnlock(opts *bind.TransactOpts, beneficiary common.Address, cheater common.Address, lockhash [32]byte, beneficiary_transferred_amount *big.Int, additional_hash [32]byte, signature []byte, merkle_proof []byte) (*types.Transaction, error) {
	return _TokenNetwork.contract.Transact(opts, "punishObsoleteUnlock", beneficiary, cheater, lockhash, beneficiary_transferred_amount, additional_hash, signature, merkle_proof)
}

// PunishObsoleteUnlock is a paid mutator transaction binding the contract method 0x6b51403c.
//
// Solidity: function punishObsoleteUnlock(beneficiary address, cheater address, lockhash bytes32, beneficiary_transferred_amount uint256, additional_hash bytes32, signature bytes, merkle_proof bytes) returns()
func (_TokenNetwork *TokenNetworkSession) PunishObsoleteUnlock(beneficiary common.Address, cheater common.Address, lockhash [32]byte, beneficiary_transferred_amount *big.Int, additional_hash [32]byte, signature []byte, merkle_proof []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.PunishObsoleteUnlock(&_TokenNetwork.TransactOpts, beneficiary, cheater, lockhash, beneficiary_transferred_amount, additional_hash, signature, merkle_proof)
}

// PunishObsoleteUnlock is a paid mutator transaction binding the contract method 0x6b51403c.
//
// Solidity: function punishObsoleteUnlock(beneficiary address, cheater address, lockhash bytes32, beneficiary_transferred_amount uint256, additional_hash bytes32, signature bytes, merkle_proof bytes) returns()
func (_TokenNetwork *TokenNetworkTransactorSession) PunishObsoleteUnlock(beneficiary common.Address, cheater common.Address, lockhash [32]byte, beneficiary_transferred_amount *big.Int, additional_hash [32]byte, signature []byte, merkle_proof []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.PunishObsoleteUnlock(&_TokenNetwork.TransactOpts, beneficiary, cheater, lockhash, beneficiary_transferred_amount, additional_hash, signature, merkle_proof)
}

// SetTotalDeposit is a paid mutator transaction binding the contract method 0xc10fd1bb.
//
// Solidity: function setTotalDeposit(participant address, partner address, total_deposit uint256) returns()
func (_TokenNetwork *TokenNetworkTransactor) SetTotalDeposit(opts *bind.TransactOpts, participant common.Address, partner common.Address, total_deposit *big.Int) (*types.Transaction, error) {
	return _TokenNetwork.contract.Transact(opts, "setTotalDeposit", participant, partner, total_deposit)
}

// SetTotalDeposit is a paid mutator transaction binding the contract method 0xc10fd1bb.
//
// Solidity: function setTotalDeposit(participant address, partner address, total_deposit uint256) returns()
func (_TokenNetwork *TokenNetworkSession) SetTotalDeposit(participant common.Address, partner common.Address, total_deposit *big.Int) (*types.Transaction, error) {
	return _TokenNetwork.Contract.SetTotalDeposit(&_TokenNetwork.TransactOpts, participant, partner, total_deposit)
}

// SetTotalDeposit is a paid mutator transaction binding the contract method 0xc10fd1bb.
//
// Solidity: function setTotalDeposit(participant address, partner address, total_deposit uint256) returns()
func (_TokenNetwork *TokenNetworkTransactorSession) SetTotalDeposit(participant common.Address, partner common.Address, total_deposit *big.Int) (*types.Transaction, error) {
	return _TokenNetwork.Contract.SetTotalDeposit(&_TokenNetwork.TransactOpts, participant, partner, total_deposit)
}

// SettleChannel is a paid mutator transaction binding the contract method 0xe11cbf99.
//
// Solidity: function settleChannel(participant1 address, participant1_transferred_amount uint256, participant1_locksroot bytes32, participant2 address, participant2_transferred_amount uint256, participant2_locksroot bytes32) returns()
func (_TokenNetwork *TokenNetworkTransactor) SettleChannel(opts *bind.TransactOpts, participant1 common.Address, participant1_transferred_amount *big.Int, participant1_locksroot [32]byte, participant2 common.Address, participant2_transferred_amount *big.Int, participant2_locksroot [32]byte) (*types.Transaction, error) {
	return _TokenNetwork.contract.Transact(opts, "settleChannel", participant1, participant1_transferred_amount, participant1_locksroot, participant2, participant2_transferred_amount, participant2_locksroot)
}

// SettleChannel is a paid mutator transaction binding the contract method 0xe11cbf99.
//
// Solidity: function settleChannel(participant1 address, participant1_transferred_amount uint256, participant1_locksroot bytes32, participant2 address, participant2_transferred_amount uint256, participant2_locksroot bytes32) returns()
func (_TokenNetwork *TokenNetworkSession) SettleChannel(participant1 common.Address, participant1_transferred_amount *big.Int, participant1_locksroot [32]byte, participant2 common.Address, participant2_transferred_amount *big.Int, participant2_locksroot [32]byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.SettleChannel(&_TokenNetwork.TransactOpts, participant1, participant1_transferred_amount, participant1_locksroot, participant2, participant2_transferred_amount, participant2_locksroot)
}

// SettleChannel is a paid mutator transaction binding the contract method 0xe11cbf99.
//
// Solidity: function settleChannel(participant1 address, participant1_transferred_amount uint256, participant1_locksroot bytes32, participant2 address, participant2_transferred_amount uint256, participant2_locksroot bytes32) returns()
func (_TokenNetwork *TokenNetworkTransactorSession) SettleChannel(participant1 common.Address, participant1_transferred_amount *big.Int, participant1_locksroot [32]byte, participant2 common.Address, participant2_transferred_amount *big.Int, participant2_locksroot [32]byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.SettleChannel(&_TokenNetwork.TransactOpts, participant1, participant1_transferred_amount, participant1_locksroot, participant2, participant2_transferred_amount, participant2_locksroot)
}

// Unlock is a paid mutator transaction binding the contract method 0x4e485437.
//
// Solidity: function unlock(participant address, partner address, transferered_amount uint256, locksroot bytes32, merkle_tree_leaves bytes) returns()
func (_TokenNetwork *TokenNetworkTransactor) Unlock(opts *bind.TransactOpts, participant common.Address, partner common.Address, transferered_amount *big.Int, locksroot [32]byte, merkle_tree_leaves []byte) (*types.Transaction, error) {
	return _TokenNetwork.contract.Transact(opts, "unlock", participant, partner, transferered_amount, locksroot, merkle_tree_leaves)
}

// Unlock is a paid mutator transaction binding the contract method 0x4e485437.
//
// Solidity: function unlock(participant address, partner address, transferered_amount uint256, locksroot bytes32, merkle_tree_leaves bytes) returns()
func (_TokenNetwork *TokenNetworkSession) Unlock(participant common.Address, partner common.Address, transferered_amount *big.Int, locksroot [32]byte, merkle_tree_leaves []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.Unlock(&_TokenNetwork.TransactOpts, participant, partner, transferered_amount, locksroot, merkle_tree_leaves)
}

// Unlock is a paid mutator transaction binding the contract method 0x4e485437.
//
// Solidity: function unlock(participant address, partner address, transferered_amount uint256, locksroot bytes32, merkle_tree_leaves bytes) returns()
func (_TokenNetwork *TokenNetworkTransactorSession) Unlock(participant common.Address, partner common.Address, transferered_amount *big.Int, locksroot [32]byte, merkle_tree_leaves []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.Unlock(&_TokenNetwork.TransactOpts, participant, partner, transferered_amount, locksroot, merkle_tree_leaves)
}

// UpdateBalanceProof is a paid mutator transaction binding the contract method 0xaaa3dbcc.
//
// Solidity: function updateBalanceProof(participant address, transferred_amount uint256, locksroot bytes32, nonce uint64, additional_hash bytes32, participant_signature bytes) returns()
func (_TokenNetwork *TokenNetworkTransactor) UpdateBalanceProof(opts *bind.TransactOpts, participant common.Address, transferred_amount *big.Int, locksroot [32]byte, nonce uint64, additional_hash [32]byte, participant_signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.contract.Transact(opts, "updateBalanceProof", participant, transferred_amount, locksroot, nonce, additional_hash, participant_signature)
}

// UpdateBalanceProof is a paid mutator transaction binding the contract method 0xaaa3dbcc.
//
// Solidity: function updateBalanceProof(participant address, transferred_amount uint256, locksroot bytes32, nonce uint64, additional_hash bytes32, participant_signature bytes) returns()
func (_TokenNetwork *TokenNetworkSession) UpdateBalanceProof(participant common.Address, transferred_amount *big.Int, locksroot [32]byte, nonce uint64, additional_hash [32]byte, participant_signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.UpdateBalanceProof(&_TokenNetwork.TransactOpts, participant, transferred_amount, locksroot, nonce, additional_hash, participant_signature)
}

// UpdateBalanceProof is a paid mutator transaction binding the contract method 0xaaa3dbcc.
//
// Solidity: function updateBalanceProof(participant address, transferred_amount uint256, locksroot bytes32, nonce uint64, additional_hash bytes32, participant_signature bytes) returns()
func (_TokenNetwork *TokenNetworkTransactorSession) UpdateBalanceProof(participant common.Address, transferred_amount *big.Int, locksroot [32]byte, nonce uint64, additional_hash [32]byte, participant_signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.UpdateBalanceProof(&_TokenNetwork.TransactOpts, participant, transferred_amount, locksroot, nonce, additional_hash, participant_signature)
}

// UpdateBalanceProofDelegate is a paid mutator transaction binding the contract method 0xf8658b25.
//
// Solidity: function updateBalanceProofDelegate(participant address, partner address, transferred_amount uint256, locksroot bytes32, nonce uint64, additional_hash bytes32, participant_signature bytes, partner_signature bytes) returns()
func (_TokenNetwork *TokenNetworkTransactor) UpdateBalanceProofDelegate(opts *bind.TransactOpts, participant common.Address, partner common.Address, transferred_amount *big.Int, locksroot [32]byte, nonce uint64, additional_hash [32]byte, participant_signature []byte, partner_signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.contract.Transact(opts, "updateBalanceProofDelegate", participant, partner, transferred_amount, locksroot, nonce, additional_hash, participant_signature, partner_signature)
}

// UpdateBalanceProofDelegate is a paid mutator transaction binding the contract method 0xf8658b25.
//
// Solidity: function updateBalanceProofDelegate(participant address, partner address, transferred_amount uint256, locksroot bytes32, nonce uint64, additional_hash bytes32, participant_signature bytes, partner_signature bytes) returns()
func (_TokenNetwork *TokenNetworkSession) UpdateBalanceProofDelegate(participant common.Address, partner common.Address, transferred_amount *big.Int, locksroot [32]byte, nonce uint64, additional_hash [32]byte, participant_signature []byte, partner_signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.UpdateBalanceProofDelegate(&_TokenNetwork.TransactOpts, participant, partner, transferred_amount, locksroot, nonce, additional_hash, participant_signature, partner_signature)
}

// UpdateBalanceProofDelegate is a paid mutator transaction binding the contract method 0xf8658b25.
//
// Solidity: function updateBalanceProofDelegate(participant address, partner address, transferred_amount uint256, locksroot bytes32, nonce uint64, additional_hash bytes32, participant_signature bytes, partner_signature bytes) returns()
func (_TokenNetwork *TokenNetworkTransactorSession) UpdateBalanceProofDelegate(participant common.Address, partner common.Address, transferred_amount *big.Int, locksroot [32]byte, nonce uint64, additional_hash [32]byte, participant_signature []byte, partner_signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.UpdateBalanceProofDelegate(&_TokenNetwork.TransactOpts, participant, partner, transferred_amount, locksroot, nonce, additional_hash, participant_signature, partner_signature)
}

// WithDraw is a paid mutator transaction binding the contract method 0x9bc6cb72.
//
// Solidity: function withDraw(participant1 address, participant1_balance uint256, participant1_withdraw uint256, participant2 address, participant2_balance uint256, participant2_withdraw uint256, participant1_signature bytes, participant2_signature bytes) returns()
func (_TokenNetwork *TokenNetworkTransactor) WithDraw(opts *bind.TransactOpts, participant1 common.Address, participant1_balance *big.Int, participant1_withdraw *big.Int, participant2 common.Address, participant2_balance *big.Int, participant2_withdraw *big.Int, participant1_signature []byte, participant2_signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.contract.Transact(opts, "withDraw", participant1, participant1_balance, participant1_withdraw, participant2, participant2_balance, participant2_withdraw, participant1_signature, participant2_signature)
}

// WithDraw is a paid mutator transaction binding the contract method 0x9bc6cb72.
//
// Solidity: function withDraw(participant1 address, participant1_balance uint256, participant1_withdraw uint256, participant2 address, participant2_balance uint256, participant2_withdraw uint256, participant1_signature bytes, participant2_signature bytes) returns()
func (_TokenNetwork *TokenNetworkSession) WithDraw(participant1 common.Address, participant1_balance *big.Int, participant1_withdraw *big.Int, participant2 common.Address, participant2_balance *big.Int, participant2_withdraw *big.Int, participant1_signature []byte, participant2_signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.WithDraw(&_TokenNetwork.TransactOpts, participant1, participant1_balance, participant1_withdraw, participant2, participant2_balance, participant2_withdraw, participant1_signature, participant2_signature)
}

// WithDraw is a paid mutator transaction binding the contract method 0x9bc6cb72.
//
// Solidity: function withDraw(participant1 address, participant1_balance uint256, participant1_withdraw uint256, participant2 address, participant2_balance uint256, participant2_withdraw uint256, participant1_signature bytes, participant2_signature bytes) returns()
func (_TokenNetwork *TokenNetworkTransactorSession) WithDraw(participant1 common.Address, participant1_balance *big.Int, participant1_withdraw *big.Int, participant2 common.Address, participant2_balance *big.Int, participant2_withdraw *big.Int, participant1_signature []byte, participant2_signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.WithDraw(&_TokenNetwork.TransactOpts, participant1, participant1_balance, participant1_withdraw, participant2, participant2_balance, participant2_withdraw, participant1_signature, participant2_signature)
}

// TokenNetworkBalanceProofUpdatedIterator is returned from FilterBalanceProofUpdated and is used to iterate over the raw logs and unpacked data for BalanceProofUpdated events raised by the TokenNetwork contract.
type TokenNetworkBalanceProofUpdatedIterator struct {
	Event *TokenNetworkBalanceProofUpdated // Event containing the contract specifics and raw log

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
func (it *TokenNetworkBalanceProofUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenNetworkBalanceProofUpdated)
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
		it.Event = new(TokenNetworkBalanceProofUpdated)
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
func (it *TokenNetworkBalanceProofUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenNetworkBalanceProofUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenNetworkBalanceProofUpdated represents a BalanceProofUpdated event raised by the TokenNetwork contract.
type TokenNetworkBalanceProofUpdated struct {
	Channel_identifier [32]byte
	Participant        common.Address
	Locksroot          [32]byte
	Transferred_amount *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterBalanceProofUpdated is a free log retrieval operation binding the contract event 0x910c9237f4197a18340110a181e8fb775496506a007a94b46f9f80f2a35918f9.
//
// Solidity: event BalanceProofUpdated(channel_identifier indexed bytes32, participant address, locksroot bytes32, transferred_amount uint256)
func (_TokenNetwork *TokenNetworkFilterer) FilterBalanceProofUpdated(opts *bind.FilterOpts, channel_identifier [][32]byte) (*TokenNetworkBalanceProofUpdatedIterator, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}

	logs, sub, err := _TokenNetwork.contract.FilterLogs(opts, "BalanceProofUpdated", channel_identifierRule)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkBalanceProofUpdatedIterator{contract: _TokenNetwork.contract, event: "BalanceProofUpdated", logs: logs, sub: sub}, nil
}

// WatchBalanceProofUpdated is a free log subscription operation binding the contract event 0x910c9237f4197a18340110a181e8fb775496506a007a94b46f9f80f2a35918f9.
//
// Solidity: event BalanceProofUpdated(channel_identifier indexed bytes32, participant address, locksroot bytes32, transferred_amount uint256)
func (_TokenNetwork *TokenNetworkFilterer) WatchBalanceProofUpdated(opts *bind.WatchOpts, sink chan<- *TokenNetworkBalanceProofUpdated, channel_identifier [][32]byte) (event.Subscription, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}

	logs, sub, err := _TokenNetwork.contract.WatchLogs(opts, "BalanceProofUpdated", channel_identifierRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenNetworkBalanceProofUpdated)
				if err := _TokenNetwork.contract.UnpackLog(event, "BalanceProofUpdated", log); err != nil {
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

// TokenNetworkChannelClosedIterator is returned from FilterChannelClosed and is used to iterate over the raw logs and unpacked data for ChannelClosed events raised by the TokenNetwork contract.
type TokenNetworkChannelClosedIterator struct {
	Event *TokenNetworkChannelClosed // Event containing the contract specifics and raw log

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
func (it *TokenNetworkChannelClosedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenNetworkChannelClosed)
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
		it.Event = new(TokenNetworkChannelClosed)
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
func (it *TokenNetworkChannelClosedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenNetworkChannelClosedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenNetworkChannelClosed represents a ChannelClosed event raised by the TokenNetwork contract.
type TokenNetworkChannelClosed struct {
	Channel_identifier  [32]byte
	Closing_participant common.Address
	Locksroot           [32]byte
	Transferred_amount  *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterChannelClosed is a free log retrieval operation binding the contract event 0x69610baaace24c039f891a11b42c0b1df1496ab0db38b0c4ee4ed33d6d53da1a.
//
// Solidity: event ChannelClosed(channel_identifier indexed bytes32, closing_participant address, locksroot bytes32, transferred_amount uint256)
func (_TokenNetwork *TokenNetworkFilterer) FilterChannelClosed(opts *bind.FilterOpts, channel_identifier [][32]byte) (*TokenNetworkChannelClosedIterator, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}

	logs, sub, err := _TokenNetwork.contract.FilterLogs(opts, "ChannelClosed", channel_identifierRule)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkChannelClosedIterator{contract: _TokenNetwork.contract, event: "ChannelClosed", logs: logs, sub: sub}, nil
}

// WatchChannelClosed is a free log subscription operation binding the contract event 0x69610baaace24c039f891a11b42c0b1df1496ab0db38b0c4ee4ed33d6d53da1a.
//
// Solidity: event ChannelClosed(channel_identifier indexed bytes32, closing_participant address, locksroot bytes32, transferred_amount uint256)
func (_TokenNetwork *TokenNetworkFilterer) WatchChannelClosed(opts *bind.WatchOpts, sink chan<- *TokenNetworkChannelClosed, channel_identifier [][32]byte) (event.Subscription, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}

	logs, sub, err := _TokenNetwork.contract.WatchLogs(opts, "ChannelClosed", channel_identifierRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenNetworkChannelClosed)
				if err := _TokenNetwork.contract.UnpackLog(event, "ChannelClosed", log); err != nil {
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

// TokenNetworkChannelCooperativeSettledIterator is returned from FilterChannelCooperativeSettled and is used to iterate over the raw logs and unpacked data for ChannelCooperativeSettled events raised by the TokenNetwork contract.
type TokenNetworkChannelCooperativeSettledIterator struct {
	Event *TokenNetworkChannelCooperativeSettled // Event containing the contract specifics and raw log

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
func (it *TokenNetworkChannelCooperativeSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenNetworkChannelCooperativeSettled)
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
		it.Event = new(TokenNetworkChannelCooperativeSettled)
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
func (it *TokenNetworkChannelCooperativeSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenNetworkChannelCooperativeSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenNetworkChannelCooperativeSettled represents a ChannelCooperativeSettled event raised by the TokenNetwork contract.
type TokenNetworkChannelCooperativeSettled struct {
	Channel_identifier  [32]byte
	Participant1_amount *big.Int
	Participant2_amount *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterChannelCooperativeSettled is a free log retrieval operation binding the contract event 0xfb2f4bc0fb2e0f1001f78d15e81a2e1981f262d31e8bd72309e26cc63bf7bb02.
//
// Solidity: event ChannelCooperativeSettled(channel_identifier indexed bytes32, participant1_amount uint256, participant2_amount uint256)
func (_TokenNetwork *TokenNetworkFilterer) FilterChannelCooperativeSettled(opts *bind.FilterOpts, channel_identifier [][32]byte) (*TokenNetworkChannelCooperativeSettledIterator, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}

	logs, sub, err := _TokenNetwork.contract.FilterLogs(opts, "ChannelCooperativeSettled", channel_identifierRule)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkChannelCooperativeSettledIterator{contract: _TokenNetwork.contract, event: "ChannelCooperativeSettled", logs: logs, sub: sub}, nil
}

// WatchChannelCooperativeSettled is a free log subscription operation binding the contract event 0xfb2f4bc0fb2e0f1001f78d15e81a2e1981f262d31e8bd72309e26cc63bf7bb02.
//
// Solidity: event ChannelCooperativeSettled(channel_identifier indexed bytes32, participant1_amount uint256, participant2_amount uint256)
func (_TokenNetwork *TokenNetworkFilterer) WatchChannelCooperativeSettled(opts *bind.WatchOpts, sink chan<- *TokenNetworkChannelCooperativeSettled, channel_identifier [][32]byte) (event.Subscription, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}

	logs, sub, err := _TokenNetwork.contract.WatchLogs(opts, "ChannelCooperativeSettled", channel_identifierRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenNetworkChannelCooperativeSettled)
				if err := _TokenNetwork.contract.UnpackLog(event, "ChannelCooperativeSettled", log); err != nil {
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

// TokenNetworkChannelNewDepositIterator is returned from FilterChannelNewDeposit and is used to iterate over the raw logs and unpacked data for ChannelNewDeposit events raised by the TokenNetwork contract.
type TokenNetworkChannelNewDepositIterator struct {
	Event *TokenNetworkChannelNewDeposit // Event containing the contract specifics and raw log

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
func (it *TokenNetworkChannelNewDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenNetworkChannelNewDeposit)
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
		it.Event = new(TokenNetworkChannelNewDeposit)
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
func (it *TokenNetworkChannelNewDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenNetworkChannelNewDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenNetworkChannelNewDeposit represents a ChannelNewDeposit event raised by the TokenNetwork contract.
type TokenNetworkChannelNewDeposit struct {
	Channel_identifier [32]byte
	Participant        common.Address
	Total_deposit      *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterChannelNewDeposit is a free log retrieval operation binding the contract event 0x0346e981e2bfa2366dc2307a8f1fa24779830a01121b1275fe565c6b98bb4d34.
//
// Solidity: event ChannelNewDeposit(channel_identifier indexed bytes32, participant address, total_deposit uint256)
func (_TokenNetwork *TokenNetworkFilterer) FilterChannelNewDeposit(opts *bind.FilterOpts, channel_identifier [][32]byte) (*TokenNetworkChannelNewDepositIterator, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}

	logs, sub, err := _TokenNetwork.contract.FilterLogs(opts, "ChannelNewDeposit", channel_identifierRule)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkChannelNewDepositIterator{contract: _TokenNetwork.contract, event: "ChannelNewDeposit", logs: logs, sub: sub}, nil
}

// WatchChannelNewDeposit is a free log subscription operation binding the contract event 0x0346e981e2bfa2366dc2307a8f1fa24779830a01121b1275fe565c6b98bb4d34.
//
// Solidity: event ChannelNewDeposit(channel_identifier indexed bytes32, participant address, total_deposit uint256)
func (_TokenNetwork *TokenNetworkFilterer) WatchChannelNewDeposit(opts *bind.WatchOpts, sink chan<- *TokenNetworkChannelNewDeposit, channel_identifier [][32]byte) (event.Subscription, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}

	logs, sub, err := _TokenNetwork.contract.WatchLogs(opts, "ChannelNewDeposit", channel_identifierRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenNetworkChannelNewDeposit)
				if err := _TokenNetwork.contract.UnpackLog(event, "ChannelNewDeposit", log); err != nil {
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

// TokenNetworkChannelOpenedIterator is returned from FilterChannelOpened and is used to iterate over the raw logs and unpacked data for ChannelOpened events raised by the TokenNetwork contract.
type TokenNetworkChannelOpenedIterator struct {
	Event *TokenNetworkChannelOpened // Event containing the contract specifics and raw log

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
func (it *TokenNetworkChannelOpenedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenNetworkChannelOpened)
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
		it.Event = new(TokenNetworkChannelOpened)
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
func (it *TokenNetworkChannelOpenedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenNetworkChannelOpenedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenNetworkChannelOpened represents a ChannelOpened event raised by the TokenNetwork contract.
type TokenNetworkChannelOpened struct {
	Channel_identifier [32]byte
	Participant1       common.Address
	Participant2       common.Address
	Settle_timeout     *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterChannelOpened is a free log retrieval operation binding the contract event 0x448d27f1fe12f92a2070111296e68fd6ef0a01c0e05bf5819eda0dbcf267bf3d.
//
// Solidity: event ChannelOpened(channel_identifier indexed bytes32, participant1 address, participant2 address, settle_timeout uint256)
func (_TokenNetwork *TokenNetworkFilterer) FilterChannelOpened(opts *bind.FilterOpts, channel_identifier [][32]byte) (*TokenNetworkChannelOpenedIterator, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}

	logs, sub, err := _TokenNetwork.contract.FilterLogs(opts, "ChannelOpened", channel_identifierRule)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkChannelOpenedIterator{contract: _TokenNetwork.contract, event: "ChannelOpened", logs: logs, sub: sub}, nil
}

// WatchChannelOpened is a free log subscription operation binding the contract event 0x448d27f1fe12f92a2070111296e68fd6ef0a01c0e05bf5819eda0dbcf267bf3d.
//
// Solidity: event ChannelOpened(channel_identifier indexed bytes32, participant1 address, participant2 address, settle_timeout uint256)
func (_TokenNetwork *TokenNetworkFilterer) WatchChannelOpened(opts *bind.WatchOpts, sink chan<- *TokenNetworkChannelOpened, channel_identifier [][32]byte) (event.Subscription, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}

	logs, sub, err := _TokenNetwork.contract.WatchLogs(opts, "ChannelOpened", channel_identifierRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenNetworkChannelOpened)
				if err := _TokenNetwork.contract.UnpackLog(event, "ChannelOpened", log); err != nil {
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

// TokenNetworkChannelSettledIterator is returned from FilterChannelSettled and is used to iterate over the raw logs and unpacked data for ChannelSettled events raised by the TokenNetwork contract.
type TokenNetworkChannelSettledIterator struct {
	Event *TokenNetworkChannelSettled // Event containing the contract specifics and raw log

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
func (it *TokenNetworkChannelSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenNetworkChannelSettled)
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
		it.Event = new(TokenNetworkChannelSettled)
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
func (it *TokenNetworkChannelSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenNetworkChannelSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenNetworkChannelSettled represents a ChannelSettled event raised by the TokenNetwork contract.
type TokenNetworkChannelSettled struct {
	Channel_identifier  [32]byte
	Participant1_amount *big.Int
	Participant2_amount *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterChannelSettled is a free log retrieval operation binding the contract event 0xf94fb5c0628a82dc90648e8dc5e983f632633b0d26603d64e8cc042ca0790aa4.
//
// Solidity: event ChannelSettled(channel_identifier indexed bytes32, participant1_amount uint256, participant2_amount uint256)
func (_TokenNetwork *TokenNetworkFilterer) FilterChannelSettled(opts *bind.FilterOpts, channel_identifier [][32]byte) (*TokenNetworkChannelSettledIterator, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}

	logs, sub, err := _TokenNetwork.contract.FilterLogs(opts, "ChannelSettled", channel_identifierRule)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkChannelSettledIterator{contract: _TokenNetwork.contract, event: "ChannelSettled", logs: logs, sub: sub}, nil
}

// WatchChannelSettled is a free log subscription operation binding the contract event 0xf94fb5c0628a82dc90648e8dc5e983f632633b0d26603d64e8cc042ca0790aa4.
//
// Solidity: event ChannelSettled(channel_identifier indexed bytes32, participant1_amount uint256, participant2_amount uint256)
func (_TokenNetwork *TokenNetworkFilterer) WatchChannelSettled(opts *bind.WatchOpts, sink chan<- *TokenNetworkChannelSettled, channel_identifier [][32]byte) (event.Subscription, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}

	logs, sub, err := _TokenNetwork.contract.WatchLogs(opts, "ChannelSettled", channel_identifierRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenNetworkChannelSettled)
				if err := _TokenNetwork.contract.UnpackLog(event, "ChannelSettled", log); err != nil {
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

// TokenNetworkChannelUnlockedIterator is returned from FilterChannelUnlocked and is used to iterate over the raw logs and unpacked data for ChannelUnlocked events raised by the TokenNetwork contract.
type TokenNetworkChannelUnlockedIterator struct {
	Event *TokenNetworkChannelUnlocked // Event containing the contract specifics and raw log

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
func (it *TokenNetworkChannelUnlockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenNetworkChannelUnlocked)
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
		it.Event = new(TokenNetworkChannelUnlocked)
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
func (it *TokenNetworkChannelUnlockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenNetworkChannelUnlockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenNetworkChannelUnlocked represents a ChannelUnlocked event raised by the TokenNetwork contract.
type TokenNetworkChannelUnlocked struct {
	Channel_identifier [32]byte
	Payer_participant  common.Address
	Locskroot          [32]byte
	Transferred_amount *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterChannelUnlocked is a free log retrieval operation binding the contract event 0x9e3b094fde58f3a83bd8b77d0a995fdb71f3169c6fa7e6d386e9f5902841e5ff.
//
// Solidity: event ChannelUnlocked(channel_identifier indexed bytes32, payer_participant address, locskroot bytes32, transferred_amount uint256)
func (_TokenNetwork *TokenNetworkFilterer) FilterChannelUnlocked(opts *bind.FilterOpts, channel_identifier [][32]byte) (*TokenNetworkChannelUnlockedIterator, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}

	logs, sub, err := _TokenNetwork.contract.FilterLogs(opts, "ChannelUnlocked", channel_identifierRule)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkChannelUnlockedIterator{contract: _TokenNetwork.contract, event: "ChannelUnlocked", logs: logs, sub: sub}, nil
}

// WatchChannelUnlocked is a free log subscription operation binding the contract event 0x9e3b094fde58f3a83bd8b77d0a995fdb71f3169c6fa7e6d386e9f5902841e5ff.
//
// Solidity: event ChannelUnlocked(channel_identifier indexed bytes32, payer_participant address, locskroot bytes32, transferred_amount uint256)
func (_TokenNetwork *TokenNetworkFilterer) WatchChannelUnlocked(opts *bind.WatchOpts, sink chan<- *TokenNetworkChannelUnlocked, channel_identifier [][32]byte) (event.Subscription, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}

	logs, sub, err := _TokenNetwork.contract.WatchLogs(opts, "ChannelUnlocked", channel_identifierRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenNetworkChannelUnlocked)
				if err := _TokenNetwork.contract.UnpackLog(event, "ChannelUnlocked", log); err != nil {
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

// TokenNetworkChannelWithdrawIterator is returned from FilterChannelWithdraw and is used to iterate over the raw logs and unpacked data for ChannelWithdraw events raised by the TokenNetwork contract.
type TokenNetworkChannelWithdrawIterator struct {
	Event *TokenNetworkChannelWithdraw // Event containing the contract specifics and raw log

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
func (it *TokenNetworkChannelWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenNetworkChannelWithdraw)
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
		it.Event = new(TokenNetworkChannelWithdraw)
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
func (it *TokenNetworkChannelWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenNetworkChannelWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenNetworkChannelWithdraw represents a ChannelWithdraw event raised by the TokenNetwork contract.
type TokenNetworkChannelWithdraw struct {
	Channel_identifier   [32]byte
	Participant1         common.Address
	Participant1_balance *big.Int
	Participant2         common.Address
	Participant2_balance *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterChannelWithdraw is a free log retrieval operation binding the contract event 0xdc5ff4ab383e66679a382f376c0e80534f51f3f3a398add646422cd81f5f815d.
//
// Solidity: event ChannelWithdraw(channel_identifier indexed bytes32, participant1 address, participant1_balance uint256, participant2 address, participant2_balance uint256)
func (_TokenNetwork *TokenNetworkFilterer) FilterChannelWithdraw(opts *bind.FilterOpts, channel_identifier [][32]byte) (*TokenNetworkChannelWithdrawIterator, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}

	logs, sub, err := _TokenNetwork.contract.FilterLogs(opts, "ChannelWithdraw", channel_identifierRule)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkChannelWithdrawIterator{contract: _TokenNetwork.contract, event: "ChannelWithdraw", logs: logs, sub: sub}, nil
}

// WatchChannelWithdraw is a free log subscription operation binding the contract event 0xdc5ff4ab383e66679a382f376c0e80534f51f3f3a398add646422cd81f5f815d.
//
// Solidity: event ChannelWithdraw(channel_identifier indexed bytes32, participant1 address, participant1_balance uint256, participant2 address, participant2_balance uint256)
func (_TokenNetwork *TokenNetworkFilterer) WatchChannelWithdraw(opts *bind.WatchOpts, sink chan<- *TokenNetworkChannelWithdraw, channel_identifier [][32]byte) (event.Subscription, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}

	logs, sub, err := _TokenNetwork.contract.WatchLogs(opts, "ChannelWithdraw", channel_identifierRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenNetworkChannelWithdraw)
				if err := _TokenNetwork.contract.UnpackLog(event, "ChannelWithdraw", log); err != nil {
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

// TokenNetworkRegistryABI is the input ABI used to generate the binding from.
const TokenNetworkRegistryABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"token_to_token_networks\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"chain_id\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_token_address\",\"type\":\"address\"}],\"name\":\"createERC20TokenNetwork\",\"outputs\":[{\"name\":\"token_network_address\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"contract_address\",\"type\":\"address\"}],\"name\":\"contractExists\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"secret_registry_address\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_secret_registry_address\",\"type\":\"address\"},{\"name\":\"_chain_id\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"token_address\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"token_network_address\",\"type\":\"address\"}],\"name\":\"TokenNetworkCreated\",\"type\":\"event\"}]"

// TokenNetworkRegistryBin is the compiled bytecode used for deploying new contracts.
const TokenNetworkRegistryBin = `0x608060405234801561001057600080fd5b5060405160408061355f8339810160405280516020909101516000811161003657600080fd5b600160a060020a038216151561004b57600080fd5b61005d82640100000000610091810204565b151561006857600080fd5b60008054600160a060020a031916600160a060020a039390931692909217909155600155610099565b6000903b1190565b6134b7806100a86000396000f3006080604052600436106100775763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630fabd9e7811461007c5780633af973b1146100b95780634cf71a04146100e05780637709bc7814610101578063b32c65c814610136578063d0ad4bec146101c0575b600080fd5b34801561008857600080fd5b5061009d600160a060020a03600435166101d5565b60408051600160a060020a039092168252519081900360200190f35b3480156100c557600080fd5b506100ce6101f0565b60408051918252519081900360200190f35b3480156100ec57600080fd5b5061009d600160a060020a03600435166101f6565b34801561010d57600080fd5b50610122600160a060020a03600435166102e1565b604080519115158252519081900360200190f35b34801561014257600080fd5b5061014b6102e9565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561018557818101518382015260200161016d565b50505050905090810190601f1680156101b25780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156101cc57600080fd5b5061009d610320565b600260205260009081526040902054600160a060020a031681565b60015481565b600160a060020a038082166000908152600260205260408120549091161561021d57600080fd5b6000546001548391600160a060020a03169061023761032f565b600160a060020a039384168152919092166020820152604080820192909252905190819003606001906000f080158015610275573d6000803e3d6000fd5b50600160a060020a03838116600081815260026020526040808220805473ffffffffffffffffffffffffffffffffffffffff1916948616948517905551939450919290917ff11a7558a113d9627989c5edf26cbd19143b7375248e621c8e30ac9e0847dc3f91a3919050565b6000903b1190565b60408051808201909152600581527f302e332e5f000000000000000000000000000000000000000000000000000000602082015281565b600054600160a060020a031681565b60405161314c8061034083390190560060806040523480156200001157600080fd5b506040516060806200314c833981016040908152815160208301519190920151600160a060020a03831615156200004757600080fd5b600160a060020a03821615156200005d57600080fd5b600081116200006b57600080fd5b6200007f8364010000000062000177810204565b15156200008b57600080fd5b6200009f8264010000000062000177810204565b1515620000ab57600080fd5b60008054600160a060020a03808616600160a060020a031992831617808455600180548784169416939093179092556002849055604080517f18160ddd000000000000000000000000000000000000000000000000000000008152905192909116916318160ddd9160048082019260209290919082900301818787803b1580156200013557600080fd5b505af11580156200014a573d6000803e3d6000fd5b505050506040513d60208110156200016157600080fd5b5051116200016e57600080fd5b5050506200017f565b6000903b1190565b612fbd806200018f6000396000f3006080604052600436106101035763ffffffff60e060020a60003504166324d73a9381146101085780633af973b1146101395780634e485437146101605780636b51403c146101d95780637709bc78146102945780637a7ebd7b146102c95780637ed74ad91461030f5780638568536a146103425780639375cff2146103f95780639bc6cb721461042b578063aaa3dbcc146104e8578063ac13370914610567578063aef91441146105c0578063b32c65c8146105f4578063b9eec0141461067e578063c10fd1bb146106fd578063e11cbf9914610727578063f8658b251461075c578063f94c9e1314610826578063fc0c546a14610882578063fc65697014610897575b600080fd5b34801561011457600080fd5b5061011d6108ce565b60408051600160a060020a039092168252519081900360200190f35b34801561014557600080fd5b5061014e6108dd565b60408051918252519081900360200190f35b34801561016c57600080fd5b50604080516020601f6084356004818101359283018490048402850184019095528184526101d794600160a060020a0381358116956024803590921695604435956064359536959460a494939101919081908401838280828437509497506108e39650505050505050565b005b3480156101e557600080fd5b50604080516020600460a43581810135601f81018490048402850184019095528484526101d7948235600160a060020a039081169560248035909216956044359560643595608435953695929460c494909390920191819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a999881019791965091820194509250829150840183828082843750949750610aff9650505050505050565b3480156102a057600080fd5b506102b5600160a060020a0360043516610c69565b604080519115158252519081900360200190f35b3480156102d557600080fd5b506102e1600435610c71565b6040805167ffffffffffffffff948516815292909316602083015260ff168183015290519081900360600190f35b34801561031b57600080fd5b50610324610ca5565b6040805167ffffffffffffffff199092168252519081900360200190f35b34801561034e57600080fd5b50604080516020601f6084356004818101359283018490048402850184019095528184526101d794600160a060020a0381358116956024803596604435909316956064359536959460a49493919091019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a999881019791965091820194509250829150840183828082843750949750610cc19650505050505050565b34801561040557600080fd5b5061040e610fa7565b6040805167ffffffffffffffff9092168252519081900360200190f35b34801561043757600080fd5b50604080516020601f60c4356004818101359283018490048402850184019095528184526101d794600160a060020a038135811695602480359660443596606435909416956084359560a435953695919460e49492930191819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a999881019791965091820194509250829150840183828082843750949750610fac9650505050505050565b3480156104f457600080fd5b50604080516020600460a43581810135601f81018490048402850184019095528484526101d7948235600160a060020a03169460248035956044359560643567ffffffffffffffff1695608435953695929460c494920191819084018382808284375094975061153c9650505050505050565b34801561057357600080fd5b5061058e600160a060020a03600435811690602435166116bc565b6040805193845267ffffffffffffffff19909216602084015267ffffffffffffffff1682820152519081900360600190f35b3480156105cc57600080fd5b506101d7600160a060020a036004358116906024351667ffffffffffffffff60443516611723565b34801561060057600080fd5b5061060961187f565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561064357818101518382015260200161062b565b50505050905090810190601f1680156106705780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561068a57600080fd5b50604080516020600460a43581810135601f81018490048402850184019095528484526101d7948235600160a060020a03169460248035956044359560643567ffffffffffffffff1695608435953695929460c49492019181908401838280828437509497506118b69650505050505050565b34801561070957600080fd5b506101d7600160a060020a0360043581169060243516604435611a48565b34801561073357600080fd5b506101d7600160a060020a0360043581169060243590604435906064351660843560a435611bbe565b34801561076857600080fd5b50604080516020601f60c4356004818101359283018490048402850184019095528184526101d794600160a060020a0381358116956024803590921695604435956064359567ffffffffffffffff608435169560a435953695919460e49492939091019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a999881019791965091820194509250829150840183828082843750949750611ed39650505050505050565b34801561083257600080fd5b5061084d600160a060020a03600435811690602435166120b8565b6040805194855267ffffffffffffffff9384166020860152919092168382015260ff9091166060830152519081900360800190f35b34801561088e57600080fd5b5061011d612109565b3480156108a357600080fd5b506101d7600160a060020a036004358116906024351667ffffffffffffffff60443516606435612118565b600154600160a060020a031681565b60025481565b60008060008060008060008088511115156108fd57600080fd5b6109078c8c61234d565b965060036000886000191660001916815260200190815260200160002091508160010160008d600160a060020a0316600160a060020a031681526020019081526020016000209050438260000160009054906101000a900467ffffffffffffffff1667ffffffffffffffff161015151561098057600080fd5b8154608060020a900460ff1660021461099857600080fd5b6001810154604060020a0292508815156109b157600080fd5b60018101546040805167ffffffffffffffff1986168152601881018a905267ffffffffffffffff60c060020a93849004169092026038830152805191829003810190912060008181526004602052919091205490965060ff1615610a1457600080fd5b6000868152600460205260409020805460ff19166001179055610a3688612498565b909550935060008411610a4857600080fd5b848914610a5457600080fd5b610a5e8a8a6127ec565b67ffffffffffffffff19848116911614610a7757600080fd5b98830198610a858a8a6127ec565b600182018054604060020a909204600160c060020a031990921691909117905560408051600160a060020a038e168152602081018790528082018c9052905188917f9e3b094fde58f3a83bd8b77d0a995fdb71f3169c6fa7e6d386e9f5902841e5ff919081900360600190a2505050505050505050505050565b600080600080600080610b128d8d61234d565b600081815260036020526040902080549197509350608060020a900460ff16600214610b3d57600080fd5b600160a060020a038d1660009081526001848101602052604090912090810154604060020a0267ffffffffffffffff191694509150831515610b7e57600080fd5b8254610ba09087908d90604060020a900467ffffffffffffffff168c8c612825565b600160a060020a038d8116911614610bb757600080fd5b50600160a060020a038b1660009081526001830160205260409020610bdc8b886128e7565b9450610be88a866127ec565b67ffffffffffffffff19168414610bfe57600080fd5b604080516000808252602082018190528251918290039092019020600184018054604060020a909204600160c060020a03199283161777ffffffffffffffffffffffffffffffffffffffffffffffff1690911790558154835401909255555050505050505050505050565b6000903b1190565b60036020526000908152604090205467ffffffffffffffff80821691604060020a810490911690608060020a900460ff1683565b6040805160008082526020820152815190819003909101902081565b6000806000806000806000610cd68d8c61234d565b600081815260036020526040902080549196509350608060020a900460ff16600114610d0157600080fd5b8254604060020a900467ffffffffffffffff169350610d25858e8e8e8e898f612a38565b9650600160a060020a038d811690881614610d3f57600080fd5b610d4e858e8e8e8e898e612a38565b9650600160a060020a038b811690881614610d6857600080fd5b5050600160a060020a03808c166000908152600180840160209081526040808420948e168452808420805486548688558786018790558683559482018690558986526003909352908420805470ffffffffffffffffffffffffffffffffff1916905591019650908c1115610e84576000809054906101000a9004600160a060020a0316600160a060020a031663a9059cbb8e8e6040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b158015610e4d57600080fd5b505af1158015610e61573d6000803e3d6000fd5b505050506040513d6020811015610e7757600080fd5b50511515610e8457600080fd5b60008a1115610f365760008054604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a038f81166004830152602482018f90529151919092169263a9059cbb92604480820193602093909283900390910190829087803b158015610eff57600080fd5b505af1158015610f13573d6000803e3d6000fd5b505050506040513d6020811015610f2957600080fd5b50511515610f3657600080fd5b8b8a018614610f4457600080fd5b8b861015610f5157600080fd5b89861015610f5e57600080fd5b604080518d8152602081018c9052815187927ffb2f4bc0fb2e0f1001f78d15e81a2e1981f262d31e8bd72309e26cc63bf7bb02928290030190a250505050505050505050505050565b600a81565b600080600080600080610fbf8e8c61234d565b600081815260036020526040902080549196509450608060020a900460ff16600114610fea57600080fd5b8d8d8c8c8f898960000160089054906101000a900467ffffffffffffffff166002546040516020018089600160a060020a0316600160a060020a03166c0100000000000000000000000002815260140188815260200187600160a060020a0316600160a060020a03166c0100000000000000000000000002815260140186815260200185815260200184600019166000191681526020018367ffffffffffffffff1667ffffffffffffffff1660c060020a028152600801828152602001985050505050505050506040516020818303038152906040526040518082805190602001908083835b602083106110ef5780518252601f1990920191602091820191016110d0565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902092506111288389612b2d565b600160a060020a038f811691161461113f57600080fd5b8d8d8c8c8f8d8a8a60000160089054906101000a900467ffffffffffffffff16600254604051602001808a600160a060020a0316600160a060020a03166c0100000000000000000000000002815260140189815260200188600160a060020a0316600160a060020a03166c0100000000000000000000000002815260140187815260200186815260200185815260200184600019166000191681526020018367ffffffffffffffff1667ffffffffffffffff1660c060020a02815260080182815260200199505050505050505050506040516020818303038152906040526040518082805190602001908083835b6020831061124c5780518252601f19909201916020918201910161122d565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902092506112858388612b2d565b600160a060020a038c811691161461129c57600080fd5b5050600160a060020a03808d166000908152600184016020526040808220928c168252902080548254019550858d11156112d557600080fd5b858a11156112e257600080fd5b8c8a0186146112f057600080fd5b60008c11156113a7576000809054906101000a9004600160a060020a0316600160a060020a031663a9059cbb8f8e6040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b15801561137057600080fd5b505af1158015611384573d6000803e3d6000fd5b505050506040513d602081101561139a57600080fd5b505115156113a757600080fd5b60008911156114595760008054604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a038f81166004830152602482018e90529151919092169263a9059cbb92604480820193602093909283900390910190829087803b15801561142257600080fd5b505af1158015611436573d6000803e3d6000fd5b505050506040513d602081101561144c57600080fd5b5051151561145957600080fd5b8c8c111561146657600080fd5b8989111561147357600080fd5b8b8d038260000181905550888a038160000181905550438460000160086101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555084600019167fdc5ff4ab383e66679a382f376c0e80534f51f3f3a398add646422cd81f5f815d8f8f8e8e6040518085600160a060020a0316600160a060020a0316815260200184815260200183600160a060020a0316600160a060020a0316815260200182815260200194505050505060405180910390a25050505050505050505050505050565b600080600061154b893361234d565b6000818152600360209081526040808320600160a060020a038e1684526001810190925290912081549295509093509150608060020a900460ff1660021461159257600080fd5b81544367ffffffffffffffff90911610156115ac57600080fd5b600181015467ffffffffffffffff60c060020a9091048116908716116115d157600080fd5b6115f8838989898660000160089054906101000a900467ffffffffffffffff168a8a612c0d565b600160a060020a038a811691161461160f57600080fd5b61161988886127ec565b60018201805467ffffffffffffffff891660c060020a02604060020a909304600160c060020a03199091161777ffffffffffffffffffffffffffffffffffffffffffffffff1691909117905560408051600160a060020a038b168152602081018990528082018a9052905184917f910c9237f4197a18340110a181e8fb775496506a007a94b46f9f80f2a35918f9919081900360600190a2505050505050505050565b6000806000806000806116cf888861234d565b6000908152600360209081526040808320600160a060020a039b909b16835260019a8b01909152902080549801549798604060020a89029860c060020a900467ffffffffffffffff16975095505050505050565b6000808260068167ffffffffffffffff161015801561174f5750622932e08167ffffffffffffffff1611155b151561175a57600080fd5b600160a060020a038616151561176f57600080fd5b600160a060020a038516151561178457600080fd5b600160a060020a03868116908616141561179d57600080fd5b6117a7868661234d565b600081815260036020526040902080549194509250608060020a900460ff16156117d057600080fd5b815470ff00000000000000000000000000000000194367ffffffffffffffff908116604060020a026fffffffffffffffff00000000000000001991881667ffffffffffffffff199094168417919091161716608060020a17835560408051600160a060020a03808a16825288166020820152808201929092525184917f448d27f1fe12f92a2070111296e68fd6ef0a01c0e05bf5819eda0dbcf267bf3d919081900360600190a2505050505050565b60408051808201909152600581527f302e332e5f000000000000000000000000000000000000000000000000000000602082015281565b6000806000806118c6338b61234d565b600081815260036020526040902080549195509250608060020a900460ff166001146118f157600080fd5b815467ffffffffffffffff1970ff0000000000000000000000000000000019909116700200000000000000000000000000000000179081164367ffffffffffffffff928316018216178355600090881611156119f95750600160a060020a0389166000908152600182016020526040902081546119889085908b908b908b90604060020a900467ffffffffffffffff168b8b612c0d565b9250600160a060020a038a8116908416146119a257600080fd5b6119ac89896127ec565b60018201805467ffffffffffffffff8a1660c060020a02604060020a909304600160c060020a03199091161777ffffffffffffffffffffffffffffffffffffffffffffffff169190911790555b60408051338152602081018a90528082018b9052905185917f69610baaace24c039f891a11b42c0b1df1496ab0db38b0c4ee4ed33d6d53da1a919081900360600190a250505050505050505050565b600080808080808611611a5a57600080fd5b611a64888861234d565b6000818152600360209081526040808320600160a060020a038d168452600181019092529091208054965091945092509050858410611aa257600080fd5b838603808501825560008054604080517f23b872dd000000000000000000000000000000000000000000000000000000008152336004820152306024820152604481018590529051939850600160a060020a03909116926323b872dd92606480840193602093929083900390910190829087803b158015611b2257600080fd5b505af1158015611b36573d6000803e3d6000fd5b505050506040513d6020811015611b4c57600080fd5b50511515611b5957600080fd5b8154608060020a900460ff16600114611b7157600080fd5b60408051600160a060020a038a16815260208101889052815185927f0346e981e2bfa2366dc2307a8f1fa24779830a01121b1275fe565c6b98bb4d34928290030190a25050505050505050565b600080600080600080611bd18c8a61234d565b600081815260036020526040902080549195509350608060020a900460ff16600214611bfc57600080fd5b82544367ffffffffffffffff918216600a0190911610611c1b57600080fd5b5050600160a060020a03808b166000908152600183016020526040808220928a1682529020611c4a8b8b6127ec565b6001830154604060020a0267ffffffffffffffff19908116911614611c6e57600080fd5b611c7888886127ec565b6001820154604060020a0267ffffffffffffffff19908116911614611c9c57600080fd5b805482548981018d81039850910195508b1115611cb857600095505b611cc28686612c9e565b600160a060020a03808e1660009081526001808701602090815260408084208481558301849055938e16835283832083815590910182905587825260039052908120805470ffffffffffffffffffffffffffffffffff191690558187039950909650861115611dd9576000809054906101000a9004600160a060020a0316600160a060020a031663a9059cbb8d886040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b158015611da257600080fd5b505af1158015611db6573d6000803e3d6000fd5b505050506040513d6020811015611dcc57600080fd5b50511515611dd957600080fd5b6000881115611e8b5760008054604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a038d81166004830152602482018d90529151919092169263a9059cbb92604480820193602093909283900390910190829087803b158015611e5457600080fd5b505af1158015611e68573d6000803e3d6000fd5b505050506040513d6020811015611e7e57600080fd5b50511515611e8b57600080fd5b60408051878152602081018a9052815186927ff94fb5c0628a82dc90648e8dc5e983f632633b0d26603d64e8cc042ca0790aa4928290030190a2505050505050505050505050565b6000806000611ee28b8b61234d565b925060036000846000191660001916815260200190815260200160002091508160010160008c600160a060020a0316600160a060020a0316815260200190815260200160002090508160000160109054906101000a900460ff1660ff166002141515611f4d57600080fd5b81544367ffffffffffffffff9091161015611f6757600080fd5b600181015467ffffffffffffffff60c060020a909104811690881611611f8c57600080fd5b611fb4838a8a8a8660000160089054906101000a900467ffffffffffffffff168b8b8b612cb6565b600160a060020a038b8116911614611fcb57600080fd5b611ff2838a8a8a8660000160089054906101000a900467ffffffffffffffff168b8b612c0d565b600160a060020a038c811691161461200957600080fd5b61201389896127ec565b60018201805467ffffffffffffffff8a1660c060020a02604060020a909304600160c060020a03199091161777ffffffffffffffffffffffffffffffffffffffffffffffff1691909117905560408051600160a060020a038d168152602081018a90528082018b9052905184917f910c9237f4197a18340110a181e8fb775496506a007a94b46f9f80f2a35918f9919081900360600190a25050505050505050505050565b6000806000806000806120cb888861234d565b600081815260036020526040902054909967ffffffffffffffff8083169a50604060020a8304169850608060020a90910460ff169650945050505050565b600054600160a060020a031681565b60008060008460068167ffffffffffffffff16101580156121465750622932e08167ffffffffffffffff1611155b151561215157600080fd5b600160a060020a038816151561216657600080fd5b600160a060020a038716151561217b57600080fd5b600160a060020a03888116908816141561219457600080fd5b61219e888861234d565b6000818152600360209081526040808320600160a060020a038d1684526001810190925290912081549296509094509250608060020a900460ff16156121e357600080fd5b825470ff00000000000000000000000000000000194367ffffffffffffffff908116604060020a026fffffffffffffffff000000000000000019918a1667ffffffffffffffff1990941693909317169190911716608060020a17835560008054604080517f23b872dd000000000000000000000000000000000000000000000000000000008152336004820152306024820152604481018990529051600160a060020a03909216926323b872dd926064808401936020939083900390910190829087803b1580156122b357600080fd5b505af11580156122c7573d6000803e3d6000fd5b505050506040513d60208110156122dd57600080fd5b505115156122ea57600080fd5b84825560408051600160a060020a03808b1682528916602082015267ffffffffffffffff881681830152905185917f448d27f1fe12f92a2070111296e68fd6ef0a01c0e05bf5819eda0dbcf267bf3d919081900360600190a25050505050505050565b600081600160a060020a031683600160a060020a031610156124175760408051600160a060020a038581166c0100000000000000000000000090810260208085019190915291861681026034840152300260488301528251808303603c018152605c90920192839052815191929182918401908083835b602083106123e35780518252601f1990920191602091820191016123c4565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390209050612492565b604080516c01000000000000000000000000600160a060020a03808616820260208085019190915290871682026034840152309190910260488301528251603c818403018152605c9092019283905281519192918291840190808383602083106123e35780518252601f1990920191602091820191016123c4565b92915050565b8051600090819081808080806060808706156124b357600080fd5b606087046001016040519080825280602002602001820160405280156124e3578160200160208202803883390190505b509050602095505b8686101561252d576124fd8a87612e1c565b958601959450925082816060880481518110151561251757fe5b60209081029091010152606095909501946124eb565b6060870496505b60018711156127c357600287061561258157806001880381518110151561255757fe5b90602001906020020151818881518110151561256f57fe5b60209081029091010152600196909601955b600095505b600187038610156127b85780866001018151811015156125a257fe5b6020908102909101015181518290889081106125ba57fe5b6020908102909101015114156125e95780868151811015156125d857fe5b906020019060200201519250612790565b80866001018151811015156125fa57fe5b60209081029091010151815182908890811061261257fe5b6020908102909101015110156126db57808681518110151561263057fe5b90602001906020020151818760010181518110151561264b57fe5b6020908102909101810151604080518084019490945283810191909152805180840382018152606090930190819052825190918291908401908083835b602083106126a75780518252601f199092019160209182019101612688565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390209250612790565b80866001018151811015156126ec57fe5b90602001906020020151818781518110151561270457fe5b6020908102909101810151604080518084019490945283810191909152805180840382018152606090930190819052825190918291908401908083835b602083106127605780518252601f199092019160209182019101612741565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902092505b8281600288048151811015156127a257fe5b6020908102909101015260029590950194612586565b600286049650612534565b8060008151811015156127d257fe5b602090810290910101519a94995093975050505050505050565b6000811580156127fa575082155b1561280757506000612492565b50604080519182526020820192909252815190819003909101902090565b60025460408051602080820188905281830189905260c060020a67ffffffffffffffff8816026060830152606882019390935260888082018690528251808303909101815260a890910191829052805160009384939182918401908083835b602083106128a35780518252601f199092019160209182019101612884565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902090506128dc8184612b2d565b979650505050505050565b6000806000602084518115156128f957fe5b061561290457600080fd5b602091505b83518211612a2f575082810151808510156129a357604080516020808201889052818301849052825180830384018152606090920192839052815191929182918401908083835b6020831061296f5780518252601f199092019160209182019101612950565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390209450612a24565b604080516020808201849052818301889052825180830384018152606090920192839052815191929182918401908083835b602083106129f45780518252601f1990920191602091820191016129d5565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902094505b602082019150612909565b50929392505050565b60025460408051600160a060020a038981166c01000000000000000000000000908102602080850191909152603484018b905291891602605483015260688201879052608882018b905267ffffffffffffffff861660c060020a0260a883015260b0808301949094528251808303909401845260d090910191829052825160009384939092909182918401908083835b60208310612ae75780518252601f199092019160209182019101612ac8565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390209050612b208184612b2d565b9998505050505050505050565b60008060008084516041141515612b4357600080fd5b50505060208201516040830151606084015160001a601b60ff82161015612b6857601b015b8060ff16601b1480612b7d57508060ff16601c145b1515612b8857600080fd5b60408051600080825260208083018085528a905260ff8516838501526060830187905260808301869052925160019360a0808501949193601f19840193928390039091019190865af1158015612be2573d6000803e3d6000fd5b5050604051601f190151945050600160a060020a0384161515612c0457600080fd5b50505092915050565b6002546040805160208082018a905281830189905260c060020a67ffffffffffffffff808a168202606085015260688401889052608884018d905288160260a883015260b0808301949094528251808303909401845260d0909101918290528251600093849390929091829184019080838360208310612ae75780518252601f199092019160209182019101612ac8565b6000818311612cad5782612caf565b815b9392505050565b600080888888878d8a6002548a6040516020018089815260200188600019166000191681526020018767ffffffffffffffff1667ffffffffffffffff1660c060020a028152600801866000191660001916815260200185600019166000191681526020018467ffffffffffffffff1667ffffffffffffffff1660c060020a02815260080183815260200182805190602001908083835b60208310612d6b5780518252601f199092019160209182019101612d4c565b6001836020036101000a038019825116818451168082178552505050505050905001985050505050505050506040516020818303038152906040526040518082805190602001908083835b60208310612dd55780518252601f199092019160209182019101612db6565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390209050612e0e8184612b2d565b9a9950505050505050505050565b6000806000806000806000878951111515612e3d5795506000945085612f85565b888801805160208083015160409384015184518084018590528086018390526060808201839052865180830390910181526080909101958690528051949a509198509550929182918401908083835b60208310612eab5780518252601f199092019160209182019101612e8c565b51815160209384036101000a6000190180199092169116179052604080519290940182900382206001547fc1f62946000000000000000000000000000000000000000000000000000000008452600484018a90529451909750600160a060020a03909416955063c1f62946945060248083019491935090918290030181600087803b158015612f3957600080fd5b505af1158015612f4d573d6000803e3d6000fd5b505050506040513d6020811015612f6357600080fd5b50519250821580612f745750828511155b15612f7e57600093505b8084965096505b505050505092509290505600a165627a7a72305820f3d21359b18fd00cd7455b2254197313cdf0f9c789d88ce4e8ae0bae4b0f17370029a165627a7a72305820bf7d8479cecb41cb2f870ffa093cf81f53a9b3b4082f844c877800aa3fd154410029`

// DeployTokenNetworkRegistry deploys a new Ethereum contract, binding an instance of TokenNetworkRegistry to it.
func DeployTokenNetworkRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, _secret_registry_address common.Address, _chain_id *big.Int) (common.Address, *types.Transaction, *TokenNetworkRegistry, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenNetworkRegistryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TokenNetworkRegistryBin), backend, _secret_registry_address, _chain_id)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TokenNetworkRegistry{TokenNetworkRegistryCaller: TokenNetworkRegistryCaller{contract: contract}, TokenNetworkRegistryTransactor: TokenNetworkRegistryTransactor{contract: contract}, TokenNetworkRegistryFilterer: TokenNetworkRegistryFilterer{contract: contract}}, nil
}

// TokenNetworkRegistry is an auto generated Go binding around an Ethereum contract.
type TokenNetworkRegistry struct {
	TokenNetworkRegistryCaller     // Read-only binding to the contract
	TokenNetworkRegistryTransactor // Write-only binding to the contract
	TokenNetworkRegistryFilterer   // Log filterer for contract events
}

// TokenNetworkRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenNetworkRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenNetworkRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenNetworkRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenNetworkRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenNetworkRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenNetworkRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenNetworkRegistrySession struct {
	Contract     *TokenNetworkRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// TokenNetworkRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenNetworkRegistryCallerSession struct {
	Contract *TokenNetworkRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// TokenNetworkRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenNetworkRegistryTransactorSession struct {
	Contract     *TokenNetworkRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// TokenNetworkRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenNetworkRegistryRaw struct {
	Contract *TokenNetworkRegistry // Generic contract binding to access the raw methods on
}

// TokenNetworkRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenNetworkRegistryCallerRaw struct {
	Contract *TokenNetworkRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// TokenNetworkRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenNetworkRegistryTransactorRaw struct {
	Contract *TokenNetworkRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenNetworkRegistry creates a new instance of TokenNetworkRegistry, bound to a specific deployed contract.
func NewTokenNetworkRegistry(address common.Address, backend bind.ContractBackend) (*TokenNetworkRegistry, error) {
	contract, err := bindTokenNetworkRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkRegistry{TokenNetworkRegistryCaller: TokenNetworkRegistryCaller{contract: contract}, TokenNetworkRegistryTransactor: TokenNetworkRegistryTransactor{contract: contract}, TokenNetworkRegistryFilterer: TokenNetworkRegistryFilterer{contract: contract}}, nil
}

// NewTokenNetworkRegistryCaller creates a new read-only instance of TokenNetworkRegistry, bound to a specific deployed contract.
func NewTokenNetworkRegistryCaller(address common.Address, caller bind.ContractCaller) (*TokenNetworkRegistryCaller, error) {
	contract, err := bindTokenNetworkRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkRegistryCaller{contract: contract}, nil
}

// NewTokenNetworkRegistryTransactor creates a new write-only instance of TokenNetworkRegistry, bound to a specific deployed contract.
func NewTokenNetworkRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenNetworkRegistryTransactor, error) {
	contract, err := bindTokenNetworkRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkRegistryTransactor{contract: contract}, nil
}

// NewTokenNetworkRegistryFilterer creates a new log filterer instance of TokenNetworkRegistry, bound to a specific deployed contract.
func NewTokenNetworkRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenNetworkRegistryFilterer, error) {
	contract, err := bindTokenNetworkRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkRegistryFilterer{contract: contract}, nil
}

// bindTokenNetworkRegistry binds a generic wrapper to an already deployed contract.
func bindTokenNetworkRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenNetworkRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenNetworkRegistry *TokenNetworkRegistryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenNetworkRegistry.Contract.TokenNetworkRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenNetworkRegistry *TokenNetworkRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenNetworkRegistry.Contract.TokenNetworkRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenNetworkRegistry *TokenNetworkRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenNetworkRegistry.Contract.TokenNetworkRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenNetworkRegistry *TokenNetworkRegistryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenNetworkRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenNetworkRegistry *TokenNetworkRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenNetworkRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenNetworkRegistry *TokenNetworkRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenNetworkRegistry.Contract.contract.Transact(opts, method, params...)
}

// Chain_id is a free data retrieval call binding the contract method 0x3af973b1.
//
// Solidity: function chain_id() constant returns(uint256)
func (_TokenNetworkRegistry *TokenNetworkRegistryCaller) Chain_id(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenNetworkRegistry.contract.Call(opts, out, "chain_id")
	return *ret0, err
}

// Chain_id is a free data retrieval call binding the contract method 0x3af973b1.
//
// Solidity: function chain_id() constant returns(uint256)
func (_TokenNetworkRegistry *TokenNetworkRegistrySession) Chain_id() (*big.Int, error) {
	return _TokenNetworkRegistry.Contract.Chain_id(&_TokenNetworkRegistry.CallOpts)
}

// Chain_id is a free data retrieval call binding the contract method 0x3af973b1.
//
// Solidity: function chain_id() constant returns(uint256)
func (_TokenNetworkRegistry *TokenNetworkRegistryCallerSession) Chain_id() (*big.Int, error) {
	return _TokenNetworkRegistry.Contract.Chain_id(&_TokenNetworkRegistry.CallOpts)
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(contract_address address) constant returns(bool)
func (_TokenNetworkRegistry *TokenNetworkRegistryCaller) ContractExists(opts *bind.CallOpts, contract_address common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TokenNetworkRegistry.contract.Call(opts, out, "contractExists", contract_address)
	return *ret0, err
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(contract_address address) constant returns(bool)
func (_TokenNetworkRegistry *TokenNetworkRegistrySession) ContractExists(contract_address common.Address) (bool, error) {
	return _TokenNetworkRegistry.Contract.ContractExists(&_TokenNetworkRegistry.CallOpts, contract_address)
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(contract_address address) constant returns(bool)
func (_TokenNetworkRegistry *TokenNetworkRegistryCallerSession) ContractExists(contract_address common.Address) (bool, error) {
	return _TokenNetworkRegistry.Contract.ContractExists(&_TokenNetworkRegistry.CallOpts, contract_address)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_TokenNetworkRegistry *TokenNetworkRegistryCaller) Contract_version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _TokenNetworkRegistry.contract.Call(opts, out, "contract_version")
	return *ret0, err
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_TokenNetworkRegistry *TokenNetworkRegistrySession) Contract_version() (string, error) {
	return _TokenNetworkRegistry.Contract.Contract_version(&_TokenNetworkRegistry.CallOpts)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_TokenNetworkRegistry *TokenNetworkRegistryCallerSession) Contract_version() (string, error) {
	return _TokenNetworkRegistry.Contract.Contract_version(&_TokenNetworkRegistry.CallOpts)
}

// Secret_registry_address is a free data retrieval call binding the contract method 0xd0ad4bec.
//
// Solidity: function secret_registry_address() constant returns(address)
func (_TokenNetworkRegistry *TokenNetworkRegistryCaller) Secret_registry_address(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenNetworkRegistry.contract.Call(opts, out, "secret_registry_address")
	return *ret0, err
}

// Secret_registry_address is a free data retrieval call binding the contract method 0xd0ad4bec.
//
// Solidity: function secret_registry_address() constant returns(address)
func (_TokenNetworkRegistry *TokenNetworkRegistrySession) Secret_registry_address() (common.Address, error) {
	return _TokenNetworkRegistry.Contract.Secret_registry_address(&_TokenNetworkRegistry.CallOpts)
}

// Secret_registry_address is a free data retrieval call binding the contract method 0xd0ad4bec.
//
// Solidity: function secret_registry_address() constant returns(address)
func (_TokenNetworkRegistry *TokenNetworkRegistryCallerSession) Secret_registry_address() (common.Address, error) {
	return _TokenNetworkRegistry.Contract.Secret_registry_address(&_TokenNetworkRegistry.CallOpts)
}

// Token_to_token_networks is a free data retrieval call binding the contract method 0x0fabd9e7.
//
// Solidity: function token_to_token_networks( address) constant returns(address)
func (_TokenNetworkRegistry *TokenNetworkRegistryCaller) Token_to_token_networks(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenNetworkRegistry.contract.Call(opts, out, "token_to_token_networks", arg0)
	return *ret0, err
}

// Token_to_token_networks is a free data retrieval call binding the contract method 0x0fabd9e7.
//
// Solidity: function token_to_token_networks( address) constant returns(address)
func (_TokenNetworkRegistry *TokenNetworkRegistrySession) Token_to_token_networks(arg0 common.Address) (common.Address, error) {
	return _TokenNetworkRegistry.Contract.Token_to_token_networks(&_TokenNetworkRegistry.CallOpts, arg0)
}

// Token_to_token_networks is a free data retrieval call binding the contract method 0x0fabd9e7.
//
// Solidity: function token_to_token_networks( address) constant returns(address)
func (_TokenNetworkRegistry *TokenNetworkRegistryCallerSession) Token_to_token_networks(arg0 common.Address) (common.Address, error) {
	return _TokenNetworkRegistry.Contract.Token_to_token_networks(&_TokenNetworkRegistry.CallOpts, arg0)
}

// CreateERC20TokenNetwork is a paid mutator transaction binding the contract method 0x4cf71a04.
//
// Solidity: function createERC20TokenNetwork(_token_address address) returns(token_network_address address)
func (_TokenNetworkRegistry *TokenNetworkRegistryTransactor) CreateERC20TokenNetwork(opts *bind.TransactOpts, _token_address common.Address) (*types.Transaction, error) {
	return _TokenNetworkRegistry.contract.Transact(opts, "createERC20TokenNetwork", _token_address)
}

// CreateERC20TokenNetwork is a paid mutator transaction binding the contract method 0x4cf71a04.
//
// Solidity: function createERC20TokenNetwork(_token_address address) returns(token_network_address address)
func (_TokenNetworkRegistry *TokenNetworkRegistrySession) CreateERC20TokenNetwork(_token_address common.Address) (*types.Transaction, error) {
	return _TokenNetworkRegistry.Contract.CreateERC20TokenNetwork(&_TokenNetworkRegistry.TransactOpts, _token_address)
}

// CreateERC20TokenNetwork is a paid mutator transaction binding the contract method 0x4cf71a04.
//
// Solidity: function createERC20TokenNetwork(_token_address address) returns(token_network_address address)
func (_TokenNetworkRegistry *TokenNetworkRegistryTransactorSession) CreateERC20TokenNetwork(_token_address common.Address) (*types.Transaction, error) {
	return _TokenNetworkRegistry.Contract.CreateERC20TokenNetwork(&_TokenNetworkRegistry.TransactOpts, _token_address)
}

// TokenNetworkRegistryTokenNetworkCreatedIterator is returned from FilterTokenNetworkCreated and is used to iterate over the raw logs and unpacked data for TokenNetworkCreated events raised by the TokenNetworkRegistry contract.
type TokenNetworkRegistryTokenNetworkCreatedIterator struct {
	Event *TokenNetworkRegistryTokenNetworkCreated // Event containing the contract specifics and raw log

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
func (it *TokenNetworkRegistryTokenNetworkCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenNetworkRegistryTokenNetworkCreated)
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
		it.Event = new(TokenNetworkRegistryTokenNetworkCreated)
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
func (it *TokenNetworkRegistryTokenNetworkCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenNetworkRegistryTokenNetworkCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenNetworkRegistryTokenNetworkCreated represents a TokenNetworkCreated event raised by the TokenNetworkRegistry contract.
type TokenNetworkRegistryTokenNetworkCreated struct {
	Token_address         common.Address
	Token_network_address common.Address
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterTokenNetworkCreated is a free log retrieval operation binding the contract event 0xf11a7558a113d9627989c5edf26cbd19143b7375248e621c8e30ac9e0847dc3f.
//
// Solidity: event TokenNetworkCreated(token_address indexed address, token_network_address indexed address)
func (_TokenNetworkRegistry *TokenNetworkRegistryFilterer) FilterTokenNetworkCreated(opts *bind.FilterOpts, token_address []common.Address, token_network_address []common.Address) (*TokenNetworkRegistryTokenNetworkCreatedIterator, error) {

	var token_addressRule []interface{}
	for _, token_addressItem := range token_address {
		token_addressRule = append(token_addressRule, token_addressItem)
	}
	var token_network_addressRule []interface{}
	for _, token_network_addressItem := range token_network_address {
		token_network_addressRule = append(token_network_addressRule, token_network_addressItem)
	}

	logs, sub, err := _TokenNetworkRegistry.contract.FilterLogs(opts, "TokenNetworkCreated", token_addressRule, token_network_addressRule)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkRegistryTokenNetworkCreatedIterator{contract: _TokenNetworkRegistry.contract, event: "TokenNetworkCreated", logs: logs, sub: sub}, nil
}

// WatchTokenNetworkCreated is a free log subscription operation binding the contract event 0xf11a7558a113d9627989c5edf26cbd19143b7375248e621c8e30ac9e0847dc3f.
//
// Solidity: event TokenNetworkCreated(token_address indexed address, token_network_address indexed address)
func (_TokenNetworkRegistry *TokenNetworkRegistryFilterer) WatchTokenNetworkCreated(opts *bind.WatchOpts, sink chan<- *TokenNetworkRegistryTokenNetworkCreated, token_address []common.Address, token_network_address []common.Address) (event.Subscription, error) {

	var token_addressRule []interface{}
	for _, token_addressItem := range token_address {
		token_addressRule = append(token_addressRule, token_addressItem)
	}
	var token_network_addressRule []interface{}
	for _, token_network_addressItem := range token_network_address {
		token_network_addressRule = append(token_network_addressRule, token_network_addressItem)
	}

	logs, sub, err := _TokenNetworkRegistry.contract.WatchLogs(opts, "TokenNetworkCreated", token_addressRule, token_network_addressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenNetworkRegistryTokenNetworkCreated)
				if err := _TokenNetworkRegistry.contract.UnpackLog(event, "TokenNetworkCreated", log); err != nil {
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
const UtilsABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"contract_address\",\"type\":\"address\"}],\"name\":\"contractExists\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// UtilsBin is the compiled bytecode used for deploying new contracts.
const UtilsBin = `0x608060405234801561001057600080fd5b50610187806100206000396000f30060806040526004361061004b5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416637709bc788114610050578063b32c65c814610092575b600080fd5b34801561005c57600080fd5b5061007e73ffffffffffffffffffffffffffffffffffffffff6004351661011c565b604080519115158252519081900360200190f35b34801561009e57600080fd5b506100a7610124565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100e15781810151838201526020016100c9565b50505050905090810190601f16801561010e5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000903b1190565b60408051808201909152600581527f302e332e5f0000000000000000000000000000000000000000000000000000006020820152815600a165627a7a7230582063bfdb817be794b9a6e6367ae4c51f7c7e599062ea99b601f169f887d4d8ca1b0029`

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
// Solidity: function contractExists(contract_address address) constant returns(bool)
func (_Utils *UtilsCaller) ContractExists(opts *bind.CallOpts, contract_address common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Utils.contract.Call(opts, out, "contractExists", contract_address)
	return *ret0, err
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(contract_address address) constant returns(bool)
func (_Utils *UtilsSession) ContractExists(contract_address common.Address) (bool, error) {
	return _Utils.Contract.ContractExists(&_Utils.CallOpts, contract_address)
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(contract_address address) constant returns(bool)
func (_Utils *UtilsCallerSession) ContractExists(contract_address common.Address) (bool, error) {
	return _Utils.Contract.ContractExists(&_Utils.CallOpts, contract_address)
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
