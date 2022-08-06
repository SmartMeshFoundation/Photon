// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package snm

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

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// MortgageABI is the input ABI used to generate the binding from.
const MortgageABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"stopContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isRunning\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"lock_time\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"investors\",\"type\":\"address[]\"},{\"name\":\"interests\",\"type\":\"uint256[]\"}],\"name\":\"payInterest\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"getFunds\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"tryStopContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"endTimeOfFunds\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"preSubFunds\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"addFunds\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"changeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"interest\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"subFunds\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"locked\",\"outputs\":[{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"endBlock\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"mortgage\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumFunds\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"timeOfFunds\",\"type\":\"uint256\"},{\"name\":\"minFunds\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"investor\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"AddFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"investor\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"SubFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_prevOwner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"OwnerUpdate\",\"type\":\"event\"}]"

// MortgageBin is the compiled bytecode used for deploying new contracts.
const MortgageBin = `0x608060405260018054600160a060020a03191681556007805460ff1916909117905534801561002d57600080fd5b50604051604080610b1783398101604052805160209091015160008054600160a060020a03191633178155821161006357600080fd5b6000811161007057600080fd5b4391909101600555600655610a8d8061008a6000396000f3006080604052600436106100f05763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166312253a6c81146100f55780632014e5d11461010c578063480bb7c4146101355780634a426da31461015c5780634d9b3735146101dd5780635afe50db146101f25780635b8257301461020757806379ba50971461021c5780638da5cb5b146102315780639b009afa1461026f578063a26759cb14610287578063a6f9dae11461028f578063ac436bdf146102bd578063b946369a146102eb578063cbf9fe5f14610300578063d09119b414610347578063d4debf0a14610375575b600080fd5b34801561010157600080fd5b5061010a61038a565b005b34801561011857600080fd5b506101216103ca565b604080519115158252519081900360200190f35b34801561014157600080fd5b5061014a6103d3565b60408051918252519081900360200190f35b6040805160206004803580820135838102808601850190965280855261010a95369593946024949385019291829185019084908082843750506040805187358901803560208181028481018201909552818452989b9a9989019892975090820195509350839250850190849080828437509497506103d99650505050505050565b3480156101e957600080fd5b5061010a6104fc565b3480156101fe57600080fd5b5061010a610555565b34801561021357600080fd5b5061014a61058a565b34801561022857600080fd5b5061010a610590565b34801561023d57600080fd5b50610246610659565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b34801561027b57600080fd5b5061010a600435610675565b61010a610817565b34801561029b57600080fd5b5061010a73ffffffffffffffffffffffffffffffffffffffff600435166108a5565b3480156102c957600080fd5b5061014a73ffffffffffffffffffffffffffffffffffffffff60043516610938565b3480156102f757600080fd5b5061010a61094a565b34801561030c57600080fd5b5061032e73ffffffffffffffffffffffffffffffffffffffff600435166109cd565b6040805192835260208301919091528051918290030190f35b34801561035357600080fd5b5061014a73ffffffffffffffffffffffffffffffffffffffff600435166109e6565b34801561038157600080fd5b5061014a6109f8565b60005473ffffffffffffffffffffffffffffffffffffffff1633146103ae57600080fd5b60075460ff16156103be57600080fd5b6007805460ff19169055565b60075460ff1681565b619d8081565b6007546000908190819060ff1615156103f157600080fd5b83518551146103ff57600080fd5b60009250600091505b84518210156104e957838281518110151561041f57fe5b6020908102909101015190506000811161043857600080fd5b610448838263ffffffff6109fe16565b925061049c60036000878581518110151561045f57fe5b602090810290910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054829063ffffffff6109fe16565b6003600087858151811015156104ae57fe5b602090810290910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002055600190910190610408565b3483146104f557600080fd5b5050505050565b60075460009060ff161561050f57600080fd5b5033600081815260026020526040808220805490839055905190929183156108fc02918491818181858888f19350505050158015610551573d6000803e3d6000fd5b5050565b60075460ff16151561056657600080fd5b6005544311801561057957506006543031105b156100f0576007805460ff19169055565b60055481565b60015473ffffffffffffffffffffffffffffffffffffffff1633146105b457600080fd5b6000546001546040805173ffffffffffffffffffffffffffffffffffffffff938416815292909116602083015280517f343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a9281900390910190a160018054600080547fffffffffffffffffffffffff000000000000000000000000000000000000000090811673ffffffffffffffffffffffffffffffffffffffff841617909155169055565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60075460009081908190819060ff16151561068f57600080fd5b6000851161069c57600080fd5b600554431161070c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6e6f7420616c6c6f7765642066756e64696e672e000000000000000000000000604482015290519081900360640190fd5b33600090815260046020526040902080549094501561072a57600080fd5b33600090815260026020908152604080832054600390925290912054909350915061076b8261075f878663ffffffff610a1116565b9063ffffffff610a2616565b905061077d828263ffffffff610a4f16565b3360009081526003602052604090205561079d838663ffffffff610a4f16565b33600090815260026020526040902081905592506107c1858263ffffffff6109fe16565b84556107d5619d804363ffffffff6109fe16565b600185015560408051868152905133917fb697db21d22bac8617372f143c06038a4bad0b3cd9483f840e296ae42db6ac40919081900360200190a25050505050565b60075460009060ff16151561082b57600080fd5b6000341161083857600080fd5b5033600090815260026020526040902054610859813463ffffffff6109fe16565b33600081815260026020908152604091829020939093558051348152905191927ff424eeb50f7d240513b6dc4a39048768557d8465bdc7d2dd363ecc538006c2be92918290030190a250565b60005473ffffffffffffffffffffffffffffffffffffffff1633146108c957600080fd5b60005473ffffffffffffffffffffffffffffffffffffffff828116911614156108f157600080fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60036020526000908152604090205481565b600754600090819060ff16151561096057600080fd5b33600090815260046020526040812080549093501161097e57600080fd5b6001820154431161098e57600080fd5b508054600080835560018301819055604051339183156108fc02918491818181858888f193505050501580156109c8573d6000803e3d6000fd5b505050565b6004602052600090815260409020805460019091015482565b60026020526000908152604090205481565b60065481565b81810182811015610a0b57fe5b92915050565b60008183811515610a1e57fe5b049392505050565b6000821515610a3757506000610a0b565b50818102818382811515610a4757fe5b0414610a0b57fe5b600082821115610a5b57fe5b509003905600a165627a7a72305820aa40aefcd59b84d03144e856b3aadea05c9b69e8e48e7e70e4a59e6dc743e8050029`

// DeployMortgage deploys a new Ethereum contract, binding an instance of Mortgage to it.
func DeployMortgage(auth *bind.TransactOpts, backend bind.ContractBackend, timeOfFunds *big.Int, minFunds *big.Int) (common.Address, *types.Transaction, *Mortgage, error) {
	parsed, err := abi.JSON(strings.NewReader(MortgageABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MortgageBin), backend, timeOfFunds, minFunds)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Mortgage{MortgageCaller: MortgageCaller{contract: contract}, MortgageTransactor: MortgageTransactor{contract: contract}, MortgageFilterer: MortgageFilterer{contract: contract}}, nil
}

// Mortgage is an auto generated Go binding around an Ethereum contract.
type Mortgage struct {
	MortgageCaller     // Read-only binding to the contract
	MortgageTransactor // Write-only binding to the contract
	MortgageFilterer   // Log filterer for contract events
}

// MortgageCaller is an auto generated read-only Go binding around an Ethereum contract.
type MortgageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MortgageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MortgageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MortgageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MortgageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MortgageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MortgageSession struct {
	Contract     *Mortgage         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MortgageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MortgageCallerSession struct {
	Contract *MortgageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// MortgageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MortgageTransactorSession struct {
	Contract     *MortgageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// MortgageRaw is an auto generated low-level Go binding around an Ethereum contract.
type MortgageRaw struct {
	Contract *Mortgage // Generic contract binding to access the raw methods on
}

// MortgageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MortgageCallerRaw struct {
	Contract *MortgageCaller // Generic read-only contract binding to access the raw methods on
}

// MortgageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MortgageTransactorRaw struct {
	Contract *MortgageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMortgage creates a new instance of Mortgage, bound to a specific deployed contract.
func NewMortgage(address common.Address, backend bind.ContractBackend) (*Mortgage, error) {
	contract, err := bindMortgage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Mortgage{MortgageCaller: MortgageCaller{contract: contract}, MortgageTransactor: MortgageTransactor{contract: contract}, MortgageFilterer: MortgageFilterer{contract: contract}}, nil
}

// NewMortgageCaller creates a new read-only instance of Mortgage, bound to a specific deployed contract.
func NewMortgageCaller(address common.Address, caller bind.ContractCaller) (*MortgageCaller, error) {
	contract, err := bindMortgage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MortgageCaller{contract: contract}, nil
}

// NewMortgageTransactor creates a new write-only instance of Mortgage, bound to a specific deployed contract.
func NewMortgageTransactor(address common.Address, transactor bind.ContractTransactor) (*MortgageTransactor, error) {
	contract, err := bindMortgage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MortgageTransactor{contract: contract}, nil
}

// NewMortgageFilterer creates a new log filterer instance of Mortgage, bound to a specific deployed contract.
func NewMortgageFilterer(address common.Address, filterer bind.ContractFilterer) (*MortgageFilterer, error) {
	contract, err := bindMortgage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MortgageFilterer{contract: contract}, nil
}

// bindMortgage binds a generic wrapper to an already deployed contract.
func bindMortgage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MortgageABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mortgage *MortgageRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Mortgage.Contract.MortgageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mortgage *MortgageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mortgage.Contract.MortgageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mortgage *MortgageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mortgage.Contract.MortgageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mortgage *MortgageCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Mortgage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mortgage *MortgageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mortgage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mortgage *MortgageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mortgage.Contract.contract.Transact(opts, method, params...)
}

// EndTimeOfFunds is a free data retrieval call binding the contract method 0x5b825730.
//
// Solidity: function endTimeOfFunds() constant returns(uint256)
func (_Mortgage *MortgageCaller) EndTimeOfFunds(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Mortgage.contract.Call(opts, out, "endTimeOfFunds")
	return *ret0, err
}

// EndTimeOfFunds is a free data retrieval call binding the contract method 0x5b825730.
//
// Solidity: function endTimeOfFunds() constant returns(uint256)
func (_Mortgage *MortgageSession) EndTimeOfFunds() (*big.Int, error) {
	return _Mortgage.Contract.EndTimeOfFunds(&_Mortgage.CallOpts)
}

// EndTimeOfFunds is a free data retrieval call binding the contract method 0x5b825730.
//
// Solidity: function endTimeOfFunds() constant returns(uint256)
func (_Mortgage *MortgageCallerSession) EndTimeOfFunds() (*big.Int, error) {
	return _Mortgage.Contract.EndTimeOfFunds(&_Mortgage.CallOpts)
}

// Interest is a free data retrieval call binding the contract method 0xac436bdf.
//
// Solidity: function interest( address) constant returns(uint256)
func (_Mortgage *MortgageCaller) Interest(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Mortgage.contract.Call(opts, out, "interest", arg0)
	return *ret0, err
}

// Interest is a free data retrieval call binding the contract method 0xac436bdf.
//
// Solidity: function interest( address) constant returns(uint256)
func (_Mortgage *MortgageSession) Interest(arg0 common.Address) (*big.Int, error) {
	return _Mortgage.Contract.Interest(&_Mortgage.CallOpts, arg0)
}

// Interest is a free data retrieval call binding the contract method 0xac436bdf.
//
// Solidity: function interest( address) constant returns(uint256)
func (_Mortgage *MortgageCallerSession) Interest(arg0 common.Address) (*big.Int, error) {
	return _Mortgage.Contract.Interest(&_Mortgage.CallOpts, arg0)
}

// IsRunning is a free data retrieval call binding the contract method 0x2014e5d1.
//
// Solidity: function isRunning() constant returns(bool)
func (_Mortgage *MortgageCaller) IsRunning(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Mortgage.contract.Call(opts, out, "isRunning")
	return *ret0, err
}

// IsRunning is a free data retrieval call binding the contract method 0x2014e5d1.
//
// Solidity: function isRunning() constant returns(bool)
func (_Mortgage *MortgageSession) IsRunning() (bool, error) {
	return _Mortgage.Contract.IsRunning(&_Mortgage.CallOpts)
}

// IsRunning is a free data retrieval call binding the contract method 0x2014e5d1.
//
// Solidity: function isRunning() constant returns(bool)
func (_Mortgage *MortgageCallerSession) IsRunning() (bool, error) {
	return _Mortgage.Contract.IsRunning(&_Mortgage.CallOpts)
}

// LockTime is a free data retrieval call binding the contract method 0x480bb7c4.
//
// Solidity: function lock_time() constant returns(uint256)
func (_Mortgage *MortgageCaller) LockTime(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Mortgage.contract.Call(opts, out, "lock_time")
	return *ret0, err
}

// LockTime is a free data retrieval call binding the contract method 0x480bb7c4.
//
// Solidity: function lock_time() constant returns(uint256)
func (_Mortgage *MortgageSession) LockTime() (*big.Int, error) {
	return _Mortgage.Contract.LockTime(&_Mortgage.CallOpts)
}

// LockTime is a free data retrieval call binding the contract method 0x480bb7c4.
//
// Solidity: function lock_time() constant returns(uint256)
func (_Mortgage *MortgageCallerSession) LockTime() (*big.Int, error) {
	return _Mortgage.Contract.LockTime(&_Mortgage.CallOpts)
}

// Locked is a free data retrieval call binding the contract method 0xcbf9fe5f.
//
// Solidity: function locked( address) constant returns(value uint256, endBlock uint256)
func (_Mortgage *MortgageCaller) Locked(opts *bind.CallOpts, arg0 common.Address) (struct {
	Value    *big.Int
	EndBlock *big.Int
}, error) {
	ret := new(struct {
		Value    *big.Int
		EndBlock *big.Int
	})
	out := ret
	err := _Mortgage.contract.Call(opts, out, "locked", arg0)
	return *ret, err
}

// Locked is a free data retrieval call binding the contract method 0xcbf9fe5f.
//
// Solidity: function locked( address) constant returns(value uint256, endBlock uint256)
func (_Mortgage *MortgageSession) Locked(arg0 common.Address) (struct {
	Value    *big.Int
	EndBlock *big.Int
}, error) {
	return _Mortgage.Contract.Locked(&_Mortgage.CallOpts, arg0)
}

// Locked is a free data retrieval call binding the contract method 0xcbf9fe5f.
//
// Solidity: function locked( address) constant returns(value uint256, endBlock uint256)
func (_Mortgage *MortgageCallerSession) Locked(arg0 common.Address) (struct {
	Value    *big.Int
	EndBlock *big.Int
}, error) {
	return _Mortgage.Contract.Locked(&_Mortgage.CallOpts, arg0)
}

// MinimumFunds is a free data retrieval call binding the contract method 0xd4debf0a.
//
// Solidity: function minimumFunds() constant returns(uint256)
func (_Mortgage *MortgageCaller) MinimumFunds(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Mortgage.contract.Call(opts, out, "minimumFunds")
	return *ret0, err
}

// MinimumFunds is a free data retrieval call binding the contract method 0xd4debf0a.
//
// Solidity: function minimumFunds() constant returns(uint256)
func (_Mortgage *MortgageSession) MinimumFunds() (*big.Int, error) {
	return _Mortgage.Contract.MinimumFunds(&_Mortgage.CallOpts)
}

// MinimumFunds is a free data retrieval call binding the contract method 0xd4debf0a.
//
// Solidity: function minimumFunds() constant returns(uint256)
func (_Mortgage *MortgageCallerSession) MinimumFunds() (*big.Int, error) {
	return _Mortgage.Contract.MinimumFunds(&_Mortgage.CallOpts)
}

// Mortgage is a free data retrieval call binding the contract method 0xd09119b4.
//
// Solidity: function mortgage( address) constant returns(uint256)
func (_Mortgage *MortgageCaller) Mortgage(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Mortgage.contract.Call(opts, out, "mortgage", arg0)
	return *ret0, err
}

// Mortgage is a free data retrieval call binding the contract method 0xd09119b4.
//
// Solidity: function mortgage( address) constant returns(uint256)
func (_Mortgage *MortgageSession) Mortgage(arg0 common.Address) (*big.Int, error) {
	return _Mortgage.Contract.Mortgage(&_Mortgage.CallOpts, arg0)
}

// Mortgage is a free data retrieval call binding the contract method 0xd09119b4.
//
// Solidity: function mortgage( address) constant returns(uint256)
func (_Mortgage *MortgageCallerSession) Mortgage(arg0 common.Address) (*big.Int, error) {
	return _Mortgage.Contract.Mortgage(&_Mortgage.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Mortgage *MortgageCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Mortgage.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Mortgage *MortgageSession) Owner() (common.Address, error) {
	return _Mortgage.Contract.Owner(&_Mortgage.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Mortgage *MortgageCallerSession) Owner() (common.Address, error) {
	return _Mortgage.Contract.Owner(&_Mortgage.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Mortgage *MortgageTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mortgage.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Mortgage *MortgageSession) AcceptOwnership() (*types.Transaction, error) {
	return _Mortgage.Contract.AcceptOwnership(&_Mortgage.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Mortgage *MortgageTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Mortgage.Contract.AcceptOwnership(&_Mortgage.TransactOpts)
}

// AddFunds is a paid mutator transaction binding the contract method 0xa26759cb.
//
// Solidity: function addFunds() returns()
func (_Mortgage *MortgageTransactor) AddFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mortgage.contract.Transact(opts, "addFunds")
}

// AddFunds is a paid mutator transaction binding the contract method 0xa26759cb.
//
// Solidity: function addFunds() returns()
func (_Mortgage *MortgageSession) AddFunds() (*types.Transaction, error) {
	return _Mortgage.Contract.AddFunds(&_Mortgage.TransactOpts)
}

// AddFunds is a paid mutator transaction binding the contract method 0xa26759cb.
//
// Solidity: function addFunds() returns()
func (_Mortgage *MortgageTransactorSession) AddFunds() (*types.Transaction, error) {
	return _Mortgage.Contract.AddFunds(&_Mortgage.TransactOpts)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(_newOwner address) returns()
func (_Mortgage *MortgageTransactor) ChangeOwner(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _Mortgage.contract.Transact(opts, "changeOwner", _newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(_newOwner address) returns()
func (_Mortgage *MortgageSession) ChangeOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _Mortgage.Contract.ChangeOwner(&_Mortgage.TransactOpts, _newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(_newOwner address) returns()
func (_Mortgage *MortgageTransactorSession) ChangeOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _Mortgage.Contract.ChangeOwner(&_Mortgage.TransactOpts, _newOwner)
}

// GetFunds is a paid mutator transaction binding the contract method 0x4d9b3735.
//
// Solidity: function getFunds() returns()
func (_Mortgage *MortgageTransactor) GetFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mortgage.contract.Transact(opts, "getFunds")
}

// GetFunds is a paid mutator transaction binding the contract method 0x4d9b3735.
//
// Solidity: function getFunds() returns()
func (_Mortgage *MortgageSession) GetFunds() (*types.Transaction, error) {
	return _Mortgage.Contract.GetFunds(&_Mortgage.TransactOpts)
}

// GetFunds is a paid mutator transaction binding the contract method 0x4d9b3735.
//
// Solidity: function getFunds() returns()
func (_Mortgage *MortgageTransactorSession) GetFunds() (*types.Transaction, error) {
	return _Mortgage.Contract.GetFunds(&_Mortgage.TransactOpts)
}

// PayInterest is a paid mutator transaction binding the contract method 0x4a426da3.
//
// Solidity: function payInterest(investors address[], interests uint256[]) returns()
func (_Mortgage *MortgageTransactor) PayInterest(opts *bind.TransactOpts, investors []common.Address, interests []*big.Int) (*types.Transaction, error) {
	return _Mortgage.contract.Transact(opts, "payInterest", investors, interests)
}

// PayInterest is a paid mutator transaction binding the contract method 0x4a426da3.
//
// Solidity: function payInterest(investors address[], interests uint256[]) returns()
func (_Mortgage *MortgageSession) PayInterest(investors []common.Address, interests []*big.Int) (*types.Transaction, error) {
	return _Mortgage.Contract.PayInterest(&_Mortgage.TransactOpts, investors, interests)
}

// PayInterest is a paid mutator transaction binding the contract method 0x4a426da3.
//
// Solidity: function payInterest(investors address[], interests uint256[]) returns()
func (_Mortgage *MortgageTransactorSession) PayInterest(investors []common.Address, interests []*big.Int) (*types.Transaction, error) {
	return _Mortgage.Contract.PayInterest(&_Mortgage.TransactOpts, investors, interests)
}

// PreSubFunds is a paid mutator transaction binding the contract method 0x9b009afa.
//
// Solidity: function preSubFunds(value uint256) returns()
func (_Mortgage *MortgageTransactor) PreSubFunds(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _Mortgage.contract.Transact(opts, "preSubFunds", value)
}

// PreSubFunds is a paid mutator transaction binding the contract method 0x9b009afa.
//
// Solidity: function preSubFunds(value uint256) returns()
func (_Mortgage *MortgageSession) PreSubFunds(value *big.Int) (*types.Transaction, error) {
	return _Mortgage.Contract.PreSubFunds(&_Mortgage.TransactOpts, value)
}

// PreSubFunds is a paid mutator transaction binding the contract method 0x9b009afa.
//
// Solidity: function preSubFunds(value uint256) returns()
func (_Mortgage *MortgageTransactorSession) PreSubFunds(value *big.Int) (*types.Transaction, error) {
	return _Mortgage.Contract.PreSubFunds(&_Mortgage.TransactOpts, value)
}

// StopContract is a paid mutator transaction binding the contract method 0x12253a6c.
//
// Solidity: function stopContract() returns()
func (_Mortgage *MortgageTransactor) StopContract(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mortgage.contract.Transact(opts, "stopContract")
}

// StopContract is a paid mutator transaction binding the contract method 0x12253a6c.
//
// Solidity: function stopContract() returns()
func (_Mortgage *MortgageSession) StopContract() (*types.Transaction, error) {
	return _Mortgage.Contract.StopContract(&_Mortgage.TransactOpts)
}

// StopContract is a paid mutator transaction binding the contract method 0x12253a6c.
//
// Solidity: function stopContract() returns()
func (_Mortgage *MortgageTransactorSession) StopContract() (*types.Transaction, error) {
	return _Mortgage.Contract.StopContract(&_Mortgage.TransactOpts)
}

// SubFunds is a paid mutator transaction binding the contract method 0xb946369a.
//
// Solidity: function subFunds() returns()
func (_Mortgage *MortgageTransactor) SubFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mortgage.contract.Transact(opts, "subFunds")
}

// SubFunds is a paid mutator transaction binding the contract method 0xb946369a.
//
// Solidity: function subFunds() returns()
func (_Mortgage *MortgageSession) SubFunds() (*types.Transaction, error) {
	return _Mortgage.Contract.SubFunds(&_Mortgage.TransactOpts)
}

// SubFunds is a paid mutator transaction binding the contract method 0xb946369a.
//
// Solidity: function subFunds() returns()
func (_Mortgage *MortgageTransactorSession) SubFunds() (*types.Transaction, error) {
	return _Mortgage.Contract.SubFunds(&_Mortgage.TransactOpts)
}

// TryStopContract is a paid mutator transaction binding the contract method 0x5afe50db.
//
// Solidity: function tryStopContract() returns()
func (_Mortgage *MortgageTransactor) TryStopContract(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mortgage.contract.Transact(opts, "tryStopContract")
}

// TryStopContract is a paid mutator transaction binding the contract method 0x5afe50db.
//
// Solidity: function tryStopContract() returns()
func (_Mortgage *MortgageSession) TryStopContract() (*types.Transaction, error) {
	return _Mortgage.Contract.TryStopContract(&_Mortgage.TransactOpts)
}

// TryStopContract is a paid mutator transaction binding the contract method 0x5afe50db.
//
// Solidity: function tryStopContract() returns()
func (_Mortgage *MortgageTransactorSession) TryStopContract() (*types.Transaction, error) {
	return _Mortgage.Contract.TryStopContract(&_Mortgage.TransactOpts)
}

// MortgageAddFundsIterator is returned from FilterAddFunds and is used to iterate over the raw logs and unpacked data for AddFunds events raised by the Mortgage contract.
type MortgageAddFundsIterator struct {
	Event *MortgageAddFunds // Event containing the contract specifics and raw log

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
func (it *MortgageAddFundsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MortgageAddFunds)
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
		it.Event = new(MortgageAddFunds)
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
func (it *MortgageAddFundsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MortgageAddFundsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MortgageAddFunds represents a AddFunds event raised by the Mortgage contract.
type MortgageAddFunds struct {
	Investor common.Address
	Value    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAddFunds is a free log retrieval operation binding the contract event 0xf424eeb50f7d240513b6dc4a39048768557d8465bdc7d2dd363ecc538006c2be.
//
// Solidity: e AddFunds(investor indexed address, value uint256)
func (_Mortgage *MortgageFilterer) FilterAddFunds(opts *bind.FilterOpts, investor []common.Address) (*MortgageAddFundsIterator, error) {

	var investorRule []interface{}
	for _, investorItem := range investor {
		investorRule = append(investorRule, investorItem)
	}

	logs, sub, err := _Mortgage.contract.FilterLogs(opts, "AddFunds", investorRule)
	if err != nil {
		return nil, err
	}
	return &MortgageAddFundsIterator{contract: _Mortgage.contract, event: "AddFunds", logs: logs, sub: sub}, nil
}

// WatchAddFunds is a free log subscription operation binding the contract event 0xf424eeb50f7d240513b6dc4a39048768557d8465bdc7d2dd363ecc538006c2be.
//
// Solidity: e AddFunds(investor indexed address, value uint256)
func (_Mortgage *MortgageFilterer) WatchAddFunds(opts *bind.WatchOpts, sink chan<- *MortgageAddFunds, investor []common.Address) (event.Subscription, error) {

	var investorRule []interface{}
	for _, investorItem := range investor {
		investorRule = append(investorRule, investorItem)
	}

	logs, sub, err := _Mortgage.contract.WatchLogs(opts, "AddFunds", investorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MortgageAddFunds)
				if err := _Mortgage.contract.UnpackLog(event, "AddFunds", log); err != nil {
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

// MortgageOwnerUpdateIterator is returned from FilterOwnerUpdate and is used to iterate over the raw logs and unpacked data for OwnerUpdate events raised by the Mortgage contract.
type MortgageOwnerUpdateIterator struct {
	Event *MortgageOwnerUpdate // Event containing the contract specifics and raw log

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
func (it *MortgageOwnerUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MortgageOwnerUpdate)
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
		it.Event = new(MortgageOwnerUpdate)
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
func (it *MortgageOwnerUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MortgageOwnerUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MortgageOwnerUpdate represents a OwnerUpdate event raised by the Mortgage contract.
type MortgageOwnerUpdate struct {
	PrevOwner common.Address
	NewOwner  common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOwnerUpdate is a free log retrieval operation binding the contract event 0x343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a.
//
// Solidity: e OwnerUpdate(_prevOwner address, _newOwner address)
func (_Mortgage *MortgageFilterer) FilterOwnerUpdate(opts *bind.FilterOpts) (*MortgageOwnerUpdateIterator, error) {

	logs, sub, err := _Mortgage.contract.FilterLogs(opts, "OwnerUpdate")
	if err != nil {
		return nil, err
	}
	return &MortgageOwnerUpdateIterator{contract: _Mortgage.contract, event: "OwnerUpdate", logs: logs, sub: sub}, nil
}

// WatchOwnerUpdate is a free log subscription operation binding the contract event 0x343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a.
//
// Solidity: e OwnerUpdate(_prevOwner address, _newOwner address)
func (_Mortgage *MortgageFilterer) WatchOwnerUpdate(opts *bind.WatchOpts, sink chan<- *MortgageOwnerUpdate) (event.Subscription, error) {

	logs, sub, err := _Mortgage.contract.WatchLogs(opts, "OwnerUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MortgageOwnerUpdate)
				if err := _Mortgage.contract.UnpackLog(event, "OwnerUpdate", log); err != nil {
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

// MortgageSubFundsIterator is returned from FilterSubFunds and is used to iterate over the raw logs and unpacked data for SubFunds events raised by the Mortgage contract.
type MortgageSubFundsIterator struct {
	Event *MortgageSubFunds // Event containing the contract specifics and raw log

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
func (it *MortgageSubFundsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MortgageSubFunds)
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
		it.Event = new(MortgageSubFunds)
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
func (it *MortgageSubFundsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MortgageSubFundsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MortgageSubFunds represents a SubFunds event raised by the Mortgage contract.
type MortgageSubFunds struct {
	Investor common.Address
	Value    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSubFunds is a free log retrieval operation binding the contract event 0xb697db21d22bac8617372f143c06038a4bad0b3cd9483f840e296ae42db6ac40.
//
// Solidity: e SubFunds(investor indexed address, value uint256)
func (_Mortgage *MortgageFilterer) FilterSubFunds(opts *bind.FilterOpts, investor []common.Address) (*MortgageSubFundsIterator, error) {

	var investorRule []interface{}
	for _, investorItem := range investor {
		investorRule = append(investorRule, investorItem)
	}

	logs, sub, err := _Mortgage.contract.FilterLogs(opts, "SubFunds", investorRule)
	if err != nil {
		return nil, err
	}
	return &MortgageSubFundsIterator{contract: _Mortgage.contract, event: "SubFunds", logs: logs, sub: sub}, nil
}

// WatchSubFunds is a free log subscription operation binding the contract event 0xb697db21d22bac8617372f143c06038a4bad0b3cd9483f840e296ae42db6ac40.
//
// Solidity: e SubFunds(investor indexed address, value uint256)
func (_Mortgage *MortgageFilterer) WatchSubFunds(opts *bind.WatchOpts, sink chan<- *MortgageSubFunds, investor []common.Address) (event.Subscription, error) {

	var investorRule []interface{}
	for _, investorItem := range investor {
		investorRule = append(investorRule, investorItem)
	}

	logs, sub, err := _Mortgage.contract.WatchLogs(opts, "SubFunds", investorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MortgageSubFunds)
				if err := _Mortgage.contract.UnpackLog(event, "SubFunds", log); err != nil {
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

// OwnedABI is the input ABI used to generate the binding from.
const OwnedABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"changeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_prevOwner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"OwnerUpdate\",\"type\":\"event\"}]"

// OwnedBin is the compiled bytecode used for deploying new contracts.
const OwnedBin = `0x608060405260018054600160a060020a031916905534801561002057600080fd5b5060008054600160a060020a03191633179055610282806100426000396000f3006080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166379ba5097811461005b5780638da5cb5b14610072578063a6f9dae1146100b0575b600080fd5b34801561006757600080fd5b506100706100de565b005b34801561007e57600080fd5b506100876101a7565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b3480156100bc57600080fd5b5061007073ffffffffffffffffffffffffffffffffffffffff600435166101c3565b60015473ffffffffffffffffffffffffffffffffffffffff16331461010257600080fd5b6000546001546040805173ffffffffffffffffffffffffffffffffffffffff938416815292909116602083015280517f343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a9281900390910190a160018054600080547fffffffffffffffffffffffff000000000000000000000000000000000000000090811673ffffffffffffffffffffffffffffffffffffffff841617909155169055565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60005473ffffffffffffffffffffffffffffffffffffffff1633146101e757600080fd5b60005473ffffffffffffffffffffffffffffffffffffffff8281169116141561020f57600080fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff929092169190911790555600a165627a7a7230582022c24d09950b641ca65c336b41af95efaa9fccf7c4854093c4fc6a7ae73ad80d0029`

// DeployOwned deploys a new Ethereum contract, binding an instance of Owned to it.
func DeployOwned(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Owned, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnedABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OwnedBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Owned{OwnedCaller: OwnedCaller{contract: contract}, OwnedTransactor: OwnedTransactor{contract: contract}, OwnedFilterer: OwnedFilterer{contract: contract}}, nil
}

// Owned is an auto generated Go binding around an Ethereum contract.
type Owned struct {
	OwnedCaller     // Read-only binding to the contract
	OwnedTransactor // Write-only binding to the contract
	OwnedFilterer   // Log filterer for contract events
}

// OwnedCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnedCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnedTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnedFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnedSession struct {
	Contract     *Owned            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnedCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnedCallerSession struct {
	Contract *OwnedCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// OwnedTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnedTransactorSession struct {
	Contract     *OwnedTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnedRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnedRaw struct {
	Contract *Owned // Generic contract binding to access the raw methods on
}

// OwnedCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnedCallerRaw struct {
	Contract *OwnedCaller // Generic read-only contract binding to access the raw methods on
}

// OwnedTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnedTransactorRaw struct {
	Contract *OwnedTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwned creates a new instance of Owned, bound to a specific deployed contract.
func NewOwned(address common.Address, backend bind.ContractBackend) (*Owned, error) {
	contract, err := bindOwned(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Owned{OwnedCaller: OwnedCaller{contract: contract}, OwnedTransactor: OwnedTransactor{contract: contract}, OwnedFilterer: OwnedFilterer{contract: contract}}, nil
}

// NewOwnedCaller creates a new read-only instance of Owned, bound to a specific deployed contract.
func NewOwnedCaller(address common.Address, caller bind.ContractCaller) (*OwnedCaller, error) {
	contract, err := bindOwned(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedCaller{contract: contract}, nil
}

// NewOwnedTransactor creates a new write-only instance of Owned, bound to a specific deployed contract.
func NewOwnedTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnedTransactor, error) {
	contract, err := bindOwned(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedTransactor{contract: contract}, nil
}

// NewOwnedFilterer creates a new log filterer instance of Owned, bound to a specific deployed contract.
func NewOwnedFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnedFilterer, error) {
	contract, err := bindOwned(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnedFilterer{contract: contract}, nil
}

// bindOwned binds a generic wrapper to an already deployed contract.
func bindOwned(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnedABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Owned *OwnedRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Owned.Contract.OwnedCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Owned *OwnedRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.Contract.OwnedTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Owned *OwnedRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owned.Contract.OwnedTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Owned *OwnedCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Owned.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Owned *OwnedTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Owned *OwnedTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owned.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Owned *OwnedCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Owned.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Owned *OwnedSession) Owner() (common.Address, error) {
	return _Owned.Contract.Owner(&_Owned.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Owned *OwnedCallerSession) Owner() (common.Address, error) {
	return _Owned.Contract.Owner(&_Owned.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Owned *OwnedTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Owned *OwnedSession) AcceptOwnership() (*types.Transaction, error) {
	return _Owned.Contract.AcceptOwnership(&_Owned.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Owned *OwnedTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Owned.Contract.AcceptOwnership(&_Owned.TransactOpts)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(_newOwner address) returns()
func (_Owned *OwnedTransactor) ChangeOwner(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _Owned.contract.Transact(opts, "changeOwner", _newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(_newOwner address) returns()
func (_Owned *OwnedSession) ChangeOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _Owned.Contract.ChangeOwner(&_Owned.TransactOpts, _newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(_newOwner address) returns()
func (_Owned *OwnedTransactorSession) ChangeOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _Owned.Contract.ChangeOwner(&_Owned.TransactOpts, _newOwner)
}

// OwnedOwnerUpdateIterator is returned from FilterOwnerUpdate and is used to iterate over the raw logs and unpacked data for OwnerUpdate events raised by the Owned contract.
type OwnedOwnerUpdateIterator struct {
	Event *OwnedOwnerUpdate // Event containing the contract specifics and raw log

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
func (it *OwnedOwnerUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnedOwnerUpdate)
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
		it.Event = new(OwnedOwnerUpdate)
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
func (it *OwnedOwnerUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnedOwnerUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnedOwnerUpdate represents a OwnerUpdate event raised by the Owned contract.
type OwnedOwnerUpdate struct {
	PrevOwner common.Address
	NewOwner  common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOwnerUpdate is a free log retrieval operation binding the contract event 0x343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a.
//
// Solidity: e OwnerUpdate(_prevOwner address, _newOwner address)
func (_Owned *OwnedFilterer) FilterOwnerUpdate(opts *bind.FilterOpts) (*OwnedOwnerUpdateIterator, error) {

	logs, sub, err := _Owned.contract.FilterLogs(opts, "OwnerUpdate")
	if err != nil {
		return nil, err
	}
	return &OwnedOwnerUpdateIterator{contract: _Owned.contract, event: "OwnerUpdate", logs: logs, sub: sub}, nil
}

// WatchOwnerUpdate is a free log subscription operation binding the contract event 0x343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a.
//
// Solidity: e OwnerUpdate(_prevOwner address, _newOwner address)
func (_Owned *OwnedFilterer) WatchOwnerUpdate(opts *bind.WatchOpts, sink chan<- *OwnedOwnerUpdate) (event.Subscription, error) {

	logs, sub, err := _Owned.contract.WatchLogs(opts, "OwnerUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnedOwnerUpdate)
				if err := _Owned.contract.UnpackLog(event, "OwnerUpdate", log); err != nil {
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

// SafeMathABI is the input ABI used to generate the binding from.
const SafeMathABI = "[]"

// SafeMathBin is the compiled bytecode used for deploying new contracts.
const SafeMathBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a72305820436975e8a4a4cd57343cc80aef03b455ea2aa03003ee91b19d172968d204bc0b0029`

// DeploySafeMath deploys a new Ethereum contract, binding an instance of SafeMath to it.
func DeploySafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// SafeMath is an auto generated Go binding around an Ethereum contract.
type SafeMath struct {
	SafeMathCaller     // Read-only binding to the contract
	SafeMathTransactor // Write-only binding to the contract
	SafeMathFilterer   // Log filterer for contract events
}

// SafeMathCaller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMathSession struct {
	Contract     *SafeMath         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMathCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMathCallerSession struct {
	Contract *SafeMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SafeMathTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMathTransactorSession struct {
	Contract     *SafeMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SafeMathRaw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMathRaw struct {
	Contract *SafeMath // Generic contract binding to access the raw methods on
}

// SafeMathCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMathCallerRaw struct {
	Contract *SafeMathCaller // Generic read-only contract binding to access the raw methods on
}

// SafeMathTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMathTransactorRaw struct {
	Contract *SafeMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath creates a new instance of SafeMath, bound to a specific deployed contract.
func NewSafeMath(address common.Address, backend bind.ContractBackend) (*SafeMath, error) {
	contract, err := bindSafeMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// NewSafeMathCaller creates a new read-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SafeMathCaller, error) {
	contract, err := bindSafeMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathCaller{contract: contract}, nil
}

// NewSafeMathTransactor creates a new write-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeMathTransactor, error) {
	contract, err := bindSafeMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathTransactor{contract: contract}, nil
}

// NewSafeMathFilterer creates a new log filterer instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeMathFilterer, error) {
	contract, err := bindSafeMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMathFilterer{contract: contract}, nil
}

// bindSafeMath binds a generic wrapper to an already deployed contract.
func bindSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.SafeMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transact(opts, method, params...)
}
