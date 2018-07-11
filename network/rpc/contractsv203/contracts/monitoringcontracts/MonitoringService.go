// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package monitoringcontracts

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
const ECVerifyBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a72305820e7d7b460f528e3740e85501ca3daf15552c792253b8738a81878262c0e1c62380029`

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

// MonitoringServiceABI is the input ABI used to generate the binding from.
const MonitoringServiceABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"beneficiary\",\"type\":\"address\"},{\"name\":\"total_deposit\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"rsb\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"contract_address\",\"type\":\"address\"}],\"name\":\"contractExists\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token_network_address\",\"type\":\"address\"},{\"name\":\"closing_participant\",\"type\":\"address\"},{\"name\":\"non_closing_participant\",\"type\":\"address\"}],\"name\":\"claimReward\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"closing_participant\",\"type\":\"address\"},{\"name\":\"non_closing_participant\",\"type\":\"address\"},{\"name\":\"balance_hash\",\"type\":\"bytes32\"},{\"name\":\"nonce\",\"type\":\"uint256\"},{\"name\":\"additional_hash\",\"type\":\"bytes32\"},{\"name\":\"closing_signature\",\"type\":\"bytes\"},{\"name\":\"non_closing_signature\",\"type\":\"bytes\"},{\"name\":\"reward_amount\",\"type\":\"uint256\"},{\"name\":\"token_network_address\",\"type\":\"address\"},{\"name\":\"reward_proof_signature\",\"type\":\"bytes\"}],\"name\":\"monitor\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_token_address\",\"type\":\"address\"},{\"name\":\"_rsb_address\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NewDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"reward_amount\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"ms_address\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"raiden_node_address\",\"type\":\"address\"}],\"name\":\"NewBalanceProofReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"ms_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"reward_identifier\",\"type\":\"bytes32\"}],\"name\":\"RewardClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawn\",\"type\":\"event\"}]"

// MonitoringServiceBin is the compiled bytecode used for deploying new contracts.
const MonitoringServiceBin = `0x608060405234801561001057600080fd5b50604051604080611157833981016040528051602090910151600160a060020a038216151561003e57600080fd5b600160a060020a038116151561005357600080fd5b6100658264010000000061014e810204565b151561007057600080fd5b6100828164010000000061014e810204565b151561008d57600080fd5b60008054600160a060020a03808516600160a060020a03199283161780845560018054868416941693909317909255604080517f18160ddd000000000000000000000000000000000000000000000000000000008152905192909116916318160ddd9160048082019260209290919082900301818787803b15801561011157600080fd5b505af1158015610125573d6000803e3d6000fd5b505050506040513d602081101561013b57600080fd5b50511161014757600080fd5b5050610156565b6000903b1190565b610ff2806101656000396000f30060806040526004361061007f5763ffffffff60e060020a60003504166327e235e381146100845780632e1a7d4d146100b757806347e7ef24146100d1578063545dcb07146100f55780637709bc7814610126578063b32c65c81461015b578063c85961b2146101e5578063d3b6c08014610212578063fc0c546a14610321575b600080fd5b34801561009057600080fd5b506100a5600160a060020a0360043516610336565b60408051918252519081900360200190f35b3480156100c357600080fd5b506100cf600435610348565b005b3480156100dd57600080fd5b506100cf600160a060020a0360043516602435610467565b34801561010157600080fd5b5061010a61059c565b60408051600160a060020a039092168252519081900360200190f35b34801561013257600080fd5b50610147600160a060020a03600435166105ab565b604080519115158252519081900360200190f35b34801561016757600080fd5b506101706105b3565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101aa578181015183820152602001610192565b50505050905090810190601f1680156101d75780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156101f157600080fd5b50610147600160a060020a03600435811690602435811690604435166105ea565b34801561021e57600080fd5b50604080516020600460a43581810135601f81018490048402850184019095528484526100cf948235600160a060020a039081169560248035909216956044359560643595608435953695929460c494909390920191819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a99988101979196509182019450925082915084018382808284375050604080516020888301358a018035601f8101839004830284018301909452838352979a89359a8a830135600160a060020a03169a919990985060609091019650919450908101925081908401838280828437509497506108d69650505050505050565b34801561032d57600080fd5b5061010a610b69565b60036020526000908152604090205481565b600160a060020a03331660009081526003602052604090205481111561036d57600080fd5b600160a060020a03338116600081815260036020908152604080832080548790039055825481517fa9059cbb000000000000000000000000000000000000000000000000000000008152600481019590955260248501879052905194169363a9059cbb93604480820194918390030190829087803b1580156103ee57600080fd5b505af1158015610402573d6000803e3d6000fd5b505050506040513d602081101561041857600080fd5b5051151561042557600080fd5b604080518281529051600160a060020a033316917f7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5919081900360200190a250565b600160a060020a038216600090815260036020526040812054821161048b57600080fd5b50600160a060020a0382166000818152600360209081526040918290208054808603908101909155825181815292519093927f2cb77763bc1e8490c1a904905c4d74b4269919aca114464f4bb4d911e60de36492908290030190a260008054604080517f23b872dd000000000000000000000000000000000000000000000000000000008152600160a060020a033381166004830152308116602483015260448201869052915191909216926323b872dd92606480820193602093909283900390910190829087803b15801561056057600080fd5b505af1158015610574573d6000803e3d6000fd5b505050506040513d602081101561058a57600080fd5b5051151561059757600080fd5b505050565b600154600160a060020a031681565b6000903b1190565b60408051808201909152600581527f302e332e5f000000000000000000000000000000000000000000000000000000602082015281565b60008060008060008088945084600160a060020a031663938bcd6789896040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182600160a060020a0316600160a060020a0316815260200192505050602060405180830381600087803b15801561066b57600080fd5b505af115801561067f573d6000803e3d6000fd5b505050506040513d602081101561069557600080fd5b50516040805160208181018490526c01000000000000000000000000600160a060020a038e160282840152825160348184030181526054909201928390528151939750909282918401908083835b602083106107025780518252601f1990920191602091820191016106e3565b5181516020939093036101000a6000190180199091169216919091179052604080519190930181900381207ff94c9e13000000000000000000000000000000000000000000000000000000008252600160a060020a038e811660048401528d811660248401529351909850928a16945063f94c9e13935060448082019360609350918290030181600087803b15801561079a57600080fd5b505af11580156107ae573d6000803e3d6000fd5b505050506040513d60608110156107c457600080fd5b5060400151915060ff8216156107d957600080fd5b50600082815260026020819052604090912090810154600160a060020a0316151561080357600080fd5b80546002820154600160a060020a03908116600090815260036020818152604080842080549690960390955585549186018054851684529285902080549092019091559054845484519081529351879491909316927fe413caa6d70a6d9b51c2af2575a2914490f614355049af8ae7cde5caab9fd201929181900390910190a350506000908152600260208190526040822082815560018101929092558101805473ffffffffffffffffffffffffffffffffffffffff199081169091556003909101805490911690555090949350505050565b600154604080517ffc7e286d000000000000000000000000000000000000000000000000000000008152600160a060020a033381811660048401529251600094859392169163fc7e286d91602480830192602092919082900301818787803b15801561094157600080fd5b505af1158015610955573d6000803e3d6000fd5b505050506040513d602081101561096b57600080fd5b50511161097757600080fd5b610986848d8d888d3389610b78565b83915081600160a060020a031663aec1dd818d8d8d8d8d8d8d6040518863ffffffff1660e060020a0281526004018088600160a060020a0316600160a060020a0316815260200187600160a060020a0316600160a060020a03168152602001866000191660001916815260200185815260200184600019166000191681526020018060200180602001838103835285818151815260200191508051906020019080838360005b83811015610a44578181015183820152602001610a2c565b50505050905090810190601f168015610a715780820380516001836020036101000a031916815260200191505b50838103825284518152845160209182019186019080838360005b83811015610aa4578181015183820152602001610a8c565b50505050905090810190601f168015610ad15780820380516001836020036101000a031916815260200191505b509950505050505050505050600060405180830381600087803b158015610af757600080fd5b505af1158015610b0b573d6000803e3d6000fd5b505050508a600160a060020a031633600160a060020a03168a7fb7c0e657d47e306f33560b0f591ec85e01ebb75f6f04a009bf02af2af6868134886040518082815260200191505060405180910390a4505050505050505050505050565b600054600160a060020a031681565b604080517f938bcd67000000000000000000000000000000000000000000000000000000008152600160a060020a0388811660048301528781166024830152915189926000928392839283929087169163938bcd679160448082019260209290919082900301818787803b158015610bef57600080fd5b505af1158015610c03573d6000803e3d6000fd5b505050506040513d6020811015610c1957600080fd5b81019080805190602001909291905050509350610ca3848a8e88600160a060020a0316633af973b16040518163ffffffff1660e060020a028152600401602060405180830381600087803b158015610c7057600080fd5b505af1158015610c84573d6000803e3d6000fd5b505050506040513d6020811015610c9a57600080fd5b50518c8b610e1d565b9250600160a060020a03808416908b1614610cbd57600080fd5b838c60405160200180836000191660001916815260200182600160a060020a0316600160a060020a03166c01000000000000000000000000028152601401925050506040516020818303038152906040526040518082805190602001908083835b60208310610d3d5780518252601f199092019160209182019101610d1e565b51815160209384036101000a600019018019909216911617905260408051929094018290039091206000818152600290925292902060018101549296509450508a119150610d8c905057600080fd5b50604080516080810182529889526020808a01988952600160a060020a039a8b168a8301908152978b1660608b0190815260009384526002918290529190922098518955965160018901559451948701805495891673ffffffffffffffffffffffffffffffffffffffff19968716179055505092516003909401805494909516939091169290921790925550505050565b6040805160208082018990528183018890526c01000000000000000000000000600160a060020a0388160260608301526074820186905260948083018690528351808403909101815260b4909201928390528151600093849392909182918401908083835b60208310610ea15780518252601f199092019160209182019101610e82565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390209050610eda8184610ee6565b98975050505050505050565b60008060008084516041141515610efc57600080fd5b50505060208201516040830151606084015160001a601b60ff82161015610f2157601b015b8060ff16601b1480610f3657508060ff16601c145b1515610f4157600080fd5b60408051600080825260208083018085528a905260ff8516838501526060830187905260808301869052925160019360a0808501949193601f19840193928390039091019190865af1158015610f9b573d6000803e3d6000fd5b5050604051601f190151945050600160a060020a0384161515610fbd57600080fd5b505050929150505600a165627a7a72305820178cfdec3b90b06d036fc5c20ed3703a8162f51c146cda00addc658efc884c560029`

// DeployMonitoringService deploys a new Ethereum contract, binding an instance of MonitoringService to it.
func DeployMonitoringService(auth *bind.TransactOpts, backend bind.ContractBackend, _token_address common.Address, _rsb_address common.Address) (common.Address, *types.Transaction, *MonitoringService, error) {
	parsed, err := abi.JSON(strings.NewReader(MonitoringServiceABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MonitoringServiceBin), backend, _token_address, _rsb_address)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MonitoringService{MonitoringServiceCaller: MonitoringServiceCaller{contract: contract}, MonitoringServiceTransactor: MonitoringServiceTransactor{contract: contract}, MonitoringServiceFilterer: MonitoringServiceFilterer{contract: contract}}, nil
}

// MonitoringService is an auto generated Go binding around an Ethereum contract.
type MonitoringService struct {
	MonitoringServiceCaller     // Read-only binding to the contract
	MonitoringServiceTransactor // Write-only binding to the contract
	MonitoringServiceFilterer   // Log filterer for contract events
}

// MonitoringServiceCaller is an auto generated read-only Go binding around an Ethereum contract.
type MonitoringServiceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MonitoringServiceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MonitoringServiceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MonitoringServiceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MonitoringServiceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MonitoringServiceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MonitoringServiceSession struct {
	Contract     *MonitoringService // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// MonitoringServiceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MonitoringServiceCallerSession struct {
	Contract *MonitoringServiceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// MonitoringServiceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MonitoringServiceTransactorSession struct {
	Contract     *MonitoringServiceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// MonitoringServiceRaw is an auto generated low-level Go binding around an Ethereum contract.
type MonitoringServiceRaw struct {
	Contract *MonitoringService // Generic contract binding to access the raw methods on
}

// MonitoringServiceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MonitoringServiceCallerRaw struct {
	Contract *MonitoringServiceCaller // Generic read-only contract binding to access the raw methods on
}

// MonitoringServiceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MonitoringServiceTransactorRaw struct {
	Contract *MonitoringServiceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMonitoringService creates a new instance of MonitoringService, bound to a specific deployed contract.
func NewMonitoringService(address common.Address, backend bind.ContractBackend) (*MonitoringService, error) {
	contract, err := bindMonitoringService(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MonitoringService{MonitoringServiceCaller: MonitoringServiceCaller{contract: contract}, MonitoringServiceTransactor: MonitoringServiceTransactor{contract: contract}, MonitoringServiceFilterer: MonitoringServiceFilterer{contract: contract}}, nil
}

// NewMonitoringServiceCaller creates a new read-only instance of MonitoringService, bound to a specific deployed contract.
func NewMonitoringServiceCaller(address common.Address, caller bind.ContractCaller) (*MonitoringServiceCaller, error) {
	contract, err := bindMonitoringService(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MonitoringServiceCaller{contract: contract}, nil
}

// NewMonitoringServiceTransactor creates a new write-only instance of MonitoringService, bound to a specific deployed contract.
func NewMonitoringServiceTransactor(address common.Address, transactor bind.ContractTransactor) (*MonitoringServiceTransactor, error) {
	contract, err := bindMonitoringService(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MonitoringServiceTransactor{contract: contract}, nil
}

// NewMonitoringServiceFilterer creates a new log filterer instance of MonitoringService, bound to a specific deployed contract.
func NewMonitoringServiceFilterer(address common.Address, filterer bind.ContractFilterer) (*MonitoringServiceFilterer, error) {
	contract, err := bindMonitoringService(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MonitoringServiceFilterer{contract: contract}, nil
}

// bindMonitoringService binds a generic wrapper to an already deployed contract.
func bindMonitoringService(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MonitoringServiceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MonitoringService *MonitoringServiceRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MonitoringService.Contract.MonitoringServiceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MonitoringService *MonitoringServiceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MonitoringService.Contract.MonitoringServiceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MonitoringService *MonitoringServiceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MonitoringService.Contract.MonitoringServiceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MonitoringService *MonitoringServiceCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MonitoringService.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MonitoringService *MonitoringServiceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MonitoringService.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MonitoringService *MonitoringServiceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MonitoringService.Contract.contract.Transact(opts, method, params...)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances( address) constant returns(uint256)
func (_MonitoringService *MonitoringServiceCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MonitoringService.contract.Call(opts, out, "balances", arg0)
	return *ret0, err
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances( address) constant returns(uint256)
func (_MonitoringService *MonitoringServiceSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _MonitoringService.Contract.Balances(&_MonitoringService.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances( address) constant returns(uint256)
func (_MonitoringService *MonitoringServiceCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _MonitoringService.Contract.Balances(&_MonitoringService.CallOpts, arg0)
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(contract_address address) constant returns(bool)
func (_MonitoringService *MonitoringServiceCaller) ContractExists(opts *bind.CallOpts, contract_address common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _MonitoringService.contract.Call(opts, out, "contractExists", contract_address)
	return *ret0, err
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(contract_address address) constant returns(bool)
func (_MonitoringService *MonitoringServiceSession) ContractExists(contract_address common.Address) (bool, error) {
	return _MonitoringService.Contract.ContractExists(&_MonitoringService.CallOpts, contract_address)
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(contract_address address) constant returns(bool)
func (_MonitoringService *MonitoringServiceCallerSession) ContractExists(contract_address common.Address) (bool, error) {
	return _MonitoringService.Contract.ContractExists(&_MonitoringService.CallOpts, contract_address)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_MonitoringService *MonitoringServiceCaller) Contract_version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _MonitoringService.contract.Call(opts, out, "contract_version")
	return *ret0, err
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_MonitoringService *MonitoringServiceSession) Contract_version() (string, error) {
	return _MonitoringService.Contract.Contract_version(&_MonitoringService.CallOpts)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_MonitoringService *MonitoringServiceCallerSession) Contract_version() (string, error) {
	return _MonitoringService.Contract.Contract_version(&_MonitoringService.CallOpts)
}

// Rsb is a free data retrieval call binding the contract method 0x545dcb07.
//
// Solidity: function rsb() constant returns(address)
func (_MonitoringService *MonitoringServiceCaller) Rsb(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _MonitoringService.contract.Call(opts, out, "rsb")
	return *ret0, err
}

// Rsb is a free data retrieval call binding the contract method 0x545dcb07.
//
// Solidity: function rsb() constant returns(address)
func (_MonitoringService *MonitoringServiceSession) Rsb() (common.Address, error) {
	return _MonitoringService.Contract.Rsb(&_MonitoringService.CallOpts)
}

// Rsb is a free data retrieval call binding the contract method 0x545dcb07.
//
// Solidity: function rsb() constant returns(address)
func (_MonitoringService *MonitoringServiceCallerSession) Rsb() (common.Address, error) {
	return _MonitoringService.Contract.Rsb(&_MonitoringService.CallOpts)
}

// TokenNetworkAddres is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_MonitoringService *MonitoringServiceCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _MonitoringService.contract.Call(opts, out, "token")
	return *ret0, err
}

// TokenNetworkAddres is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_MonitoringService *MonitoringServiceSession) Token() (common.Address, error) {
	return _MonitoringService.Contract.Token(&_MonitoringService.CallOpts)
}

// TokenNetworkAddres is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_MonitoringService *MonitoringServiceCallerSession) Token() (common.Address, error) {
	return _MonitoringService.Contract.Token(&_MonitoringService.CallOpts)
}

// ClaimReward is a paid mutator transaction binding the contract method 0xc85961b2.
//
// Solidity: function claimReward(token_network_address address, closing_participant address, non_closing_participant address) returns(bool)
func (_MonitoringService *MonitoringServiceTransactor) ClaimReward(opts *bind.TransactOpts, token_network_address common.Address, closing_participant common.Address, non_closing_participant common.Address) (*types.Transaction, error) {
	return _MonitoringService.contract.Transact(opts, "claimReward", token_network_address, closing_participant, non_closing_participant)
}

// ClaimReward is a paid mutator transaction binding the contract method 0xc85961b2.
//
// Solidity: function claimReward(token_network_address address, closing_participant address, non_closing_participant address) returns(bool)
func (_MonitoringService *MonitoringServiceSession) ClaimReward(token_network_address common.Address, closing_participant common.Address, non_closing_participant common.Address) (*types.Transaction, error) {
	return _MonitoringService.Contract.ClaimReward(&_MonitoringService.TransactOpts, token_network_address, closing_participant, non_closing_participant)
}

// ClaimReward is a paid mutator transaction binding the contract method 0xc85961b2.
//
// Solidity: function claimReward(token_network_address address, closing_participant address, non_closing_participant address) returns(bool)
func (_MonitoringService *MonitoringServiceTransactorSession) ClaimReward(token_network_address common.Address, closing_participant common.Address, non_closing_participant common.Address) (*types.Transaction, error) {
	return _MonitoringService.Contract.ClaimReward(&_MonitoringService.TransactOpts, token_network_address, closing_participant, non_closing_participant)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(beneficiary address, total_deposit uint256) returns()
func (_MonitoringService *MonitoringServiceTransactor) Deposit(opts *bind.TransactOpts, beneficiary common.Address, total_deposit *big.Int) (*types.Transaction, error) {
	return _MonitoringService.contract.Transact(opts, "deposit", beneficiary, total_deposit)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(beneficiary address, total_deposit uint256) returns()
func (_MonitoringService *MonitoringServiceSession) Deposit(beneficiary common.Address, total_deposit *big.Int) (*types.Transaction, error) {
	return _MonitoringService.Contract.Deposit(&_MonitoringService.TransactOpts, beneficiary, total_deposit)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(beneficiary address, total_deposit uint256) returns()
func (_MonitoringService *MonitoringServiceTransactorSession) Deposit(beneficiary common.Address, total_deposit *big.Int) (*types.Transaction, error) {
	return _MonitoringService.Contract.Deposit(&_MonitoringService.TransactOpts, beneficiary, total_deposit)
}

// Monitor is a paid mutator transaction binding the contract method 0xd3b6c080.
//
// Solidity: function monitor(closing_participant address, non_closing_participant address, balance_hash bytes32, nonce uint256, additional_hash bytes32, closing_signature bytes, non_closing_signature bytes, reward_amount uint256, token_network_address address, reward_proof_signature bytes) returns()
func (_MonitoringService *MonitoringServiceTransactor) Monitor(opts *bind.TransactOpts, closing_participant common.Address, non_closing_participant common.Address, balance_hash [32]byte, nonce *big.Int, additional_hash [32]byte, closing_signature []byte, non_closing_signature []byte, reward_amount *big.Int, token_network_address common.Address, reward_proof_signature []byte) (*types.Transaction, error) {
	return _MonitoringService.contract.Transact(opts, "monitor", closing_participant, non_closing_participant, balance_hash, nonce, additional_hash, closing_signature, non_closing_signature, reward_amount, token_network_address, reward_proof_signature)
}

// Monitor is a paid mutator transaction binding the contract method 0xd3b6c080.
//
// Solidity: function monitor(closing_participant address, non_closing_participant address, balance_hash bytes32, nonce uint256, additional_hash bytes32, closing_signature bytes, non_closing_signature bytes, reward_amount uint256, token_network_address address, reward_proof_signature bytes) returns()
func (_MonitoringService *MonitoringServiceSession) Monitor(closing_participant common.Address, non_closing_participant common.Address, balance_hash [32]byte, nonce *big.Int, additional_hash [32]byte, closing_signature []byte, non_closing_signature []byte, reward_amount *big.Int, token_network_address common.Address, reward_proof_signature []byte) (*types.Transaction, error) {
	return _MonitoringService.Contract.Monitor(&_MonitoringService.TransactOpts, closing_participant, non_closing_participant, balance_hash, nonce, additional_hash, closing_signature, non_closing_signature, reward_amount, token_network_address, reward_proof_signature)
}

// Monitor is a paid mutator transaction binding the contract method 0xd3b6c080.
//
// Solidity: function monitor(closing_participant address, non_closing_participant address, balance_hash bytes32, nonce uint256, additional_hash bytes32, closing_signature bytes, non_closing_signature bytes, reward_amount uint256, token_network_address address, reward_proof_signature bytes) returns()
func (_MonitoringService *MonitoringServiceTransactorSession) Monitor(closing_participant common.Address, non_closing_participant common.Address, balance_hash [32]byte, nonce *big.Int, additional_hash [32]byte, closing_signature []byte, non_closing_signature []byte, reward_amount *big.Int, token_network_address common.Address, reward_proof_signature []byte) (*types.Transaction, error) {
	return _MonitoringService.Contract.Monitor(&_MonitoringService.TransactOpts, closing_participant, non_closing_participant, balance_hash, nonce, additional_hash, closing_signature, non_closing_signature, reward_amount, token_network_address, reward_proof_signature)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(amount uint256) returns()
func (_MonitoringService *MonitoringServiceTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _MonitoringService.contract.Transact(opts, "withdraw", amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(amount uint256) returns()
func (_MonitoringService *MonitoringServiceSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _MonitoringService.Contract.Withdraw(&_MonitoringService.TransactOpts, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(amount uint256) returns()
func (_MonitoringService *MonitoringServiceTransactorSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _MonitoringService.Contract.Withdraw(&_MonitoringService.TransactOpts, amount)
}

// MonitoringServiceNewBalanceProofReceivedIterator is returned from FilterNewBalanceProofReceived and is used to iterate over the raw logs and unpacked data for NewBalanceProofReceived events raised by the MonitoringService contract.
type MonitoringServiceNewBalanceProofReceivedIterator struct {
	Event *MonitoringServiceNewBalanceProofReceived // Event containing the contract specifics and raw log

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
func (it *MonitoringServiceNewBalanceProofReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MonitoringServiceNewBalanceProofReceived)
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
		it.Event = new(MonitoringServiceNewBalanceProofReceived)
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
func (it *MonitoringServiceNewBalanceProofReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MonitoringServiceNewBalanceProofReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MonitoringServiceNewBalanceProofReceived represents a NewBalanceProofReceived event raised by the MonitoringService contract.
type MonitoringServiceNewBalanceProofReceived struct {
	Reward_amount       *big.Int
	Nonce               *big.Int
	Ms_address          common.Address
	Raiden_node_address common.Address
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterNewBalanceProofReceived is a free log retrieval operation binding the contract event 0xb7c0e657d47e306f33560b0f591ec85e01ebb75f6f04a009bf02af2af6868134.
//
// Solidity: event NewBalanceProofReceived(reward_amount uint256, nonce indexed uint256, ms_address indexed address, raiden_node_address indexed address)
func (_MonitoringService *MonitoringServiceFilterer) FilterNewBalanceProofReceived(opts *bind.FilterOpts, nonce []*big.Int, ms_address []common.Address, raiden_node_address []common.Address) (*MonitoringServiceNewBalanceProofReceivedIterator, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}
	var ms_addressRule []interface{}
	for _, ms_addressItem := range ms_address {
		ms_addressRule = append(ms_addressRule, ms_addressItem)
	}
	var raiden_node_addressRule []interface{}
	for _, raiden_node_addressItem := range raiden_node_address {
		raiden_node_addressRule = append(raiden_node_addressRule, raiden_node_addressItem)
	}

	logs, sub, err := _MonitoringService.contract.FilterLogs(opts, "NewBalanceProofReceived", nonceRule, ms_addressRule, raiden_node_addressRule)
	if err != nil {
		return nil, err
	}
	return &MonitoringServiceNewBalanceProofReceivedIterator{contract: _MonitoringService.contract, event: "NewBalanceProofReceived", logs: logs, sub: sub}, nil
}

// WatchNewBalanceProofReceived is a free log subscription operation binding the contract event 0xb7c0e657d47e306f33560b0f591ec85e01ebb75f6f04a009bf02af2af6868134.
//
// Solidity: event NewBalanceProofReceived(reward_amount uint256, nonce indexed uint256, ms_address indexed address, raiden_node_address indexed address)
func (_MonitoringService *MonitoringServiceFilterer) WatchNewBalanceProofReceived(opts *bind.WatchOpts, sink chan<- *MonitoringServiceNewBalanceProofReceived, nonce []*big.Int, ms_address []common.Address, raiden_node_address []common.Address) (event.Subscription, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}
	var ms_addressRule []interface{}
	for _, ms_addressItem := range ms_address {
		ms_addressRule = append(ms_addressRule, ms_addressItem)
	}
	var raiden_node_addressRule []interface{}
	for _, raiden_node_addressItem := range raiden_node_address {
		raiden_node_addressRule = append(raiden_node_addressRule, raiden_node_addressItem)
	}

	logs, sub, err := _MonitoringService.contract.WatchLogs(opts, "NewBalanceProofReceived", nonceRule, ms_addressRule, raiden_node_addressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MonitoringServiceNewBalanceProofReceived)
				if err := _MonitoringService.contract.UnpackLog(event, "NewBalanceProofReceived", log); err != nil {
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

// MonitoringServiceNewDepositIterator is returned from FilterNewDeposit and is used to iterate over the raw logs and unpacked data for NewDeposit events raised by the MonitoringService contract.
type MonitoringServiceNewDepositIterator struct {
	Event *MonitoringServiceNewDeposit // Event containing the contract specifics and raw log

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
func (it *MonitoringServiceNewDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MonitoringServiceNewDeposit)
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
		it.Event = new(MonitoringServiceNewDeposit)
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
func (it *MonitoringServiceNewDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MonitoringServiceNewDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MonitoringServiceNewDeposit represents a NewDeposit event raised by the MonitoringService contract.
type MonitoringServiceNewDeposit struct {
	Receiver common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNewDeposit is a free log retrieval operation binding the contract event 0x2cb77763bc1e8490c1a904905c4d74b4269919aca114464f4bb4d911e60de364.
//
// Solidity: event NewDeposit(receiver indexed address, amount uint256)
func (_MonitoringService *MonitoringServiceFilterer) FilterNewDeposit(opts *bind.FilterOpts, receiver []common.Address) (*MonitoringServiceNewDepositIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _MonitoringService.contract.FilterLogs(opts, "NewDeposit", receiverRule)
	if err != nil {
		return nil, err
	}
	return &MonitoringServiceNewDepositIterator{contract: _MonitoringService.contract, event: "NewDeposit", logs: logs, sub: sub}, nil
}

// WatchNewDeposit is a free log subscription operation binding the contract event 0x2cb77763bc1e8490c1a904905c4d74b4269919aca114464f4bb4d911e60de364.
//
// Solidity: event NewDeposit(receiver indexed address, amount uint256)
func (_MonitoringService *MonitoringServiceFilterer) WatchNewDeposit(opts *bind.WatchOpts, sink chan<- *MonitoringServiceNewDeposit, receiver []common.Address) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _MonitoringService.contract.WatchLogs(opts, "NewDeposit", receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MonitoringServiceNewDeposit)
				if err := _MonitoringService.contract.UnpackLog(event, "NewDeposit", log); err != nil {
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

// MonitoringServiceRewardClaimedIterator is returned from FilterRewardClaimed and is used to iterate over the raw logs and unpacked data for RewardClaimed events raised by the MonitoringService contract.
type MonitoringServiceRewardClaimedIterator struct {
	Event *MonitoringServiceRewardClaimed // Event containing the contract specifics and raw log

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
func (it *MonitoringServiceRewardClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MonitoringServiceRewardClaimed)
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
		it.Event = new(MonitoringServiceRewardClaimed)
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
func (it *MonitoringServiceRewardClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MonitoringServiceRewardClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MonitoringServiceRewardClaimed represents a RewardClaimed event raised by the MonitoringService contract.
type MonitoringServiceRewardClaimed struct {
	Ms_address        common.Address
	Amount            *big.Int
	Reward_identifier [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRewardClaimed is a free log retrieval operation binding the contract event 0xe413caa6d70a6d9b51c2af2575a2914490f614355049af8ae7cde5caab9fd201.
//
// Solidity: event RewardClaimed(ms_address indexed address, amount uint256, reward_identifier indexed bytes32)
func (_MonitoringService *MonitoringServiceFilterer) FilterRewardClaimed(opts *bind.FilterOpts, ms_address []common.Address, reward_identifier [][32]byte) (*MonitoringServiceRewardClaimedIterator, error) {

	var ms_addressRule []interface{}
	for _, ms_addressItem := range ms_address {
		ms_addressRule = append(ms_addressRule, ms_addressItem)
	}

	var reward_identifierRule []interface{}
	for _, reward_identifierItem := range reward_identifier {
		reward_identifierRule = append(reward_identifierRule, reward_identifierItem)
	}

	logs, sub, err := _MonitoringService.contract.FilterLogs(opts, "RewardClaimed", ms_addressRule, reward_identifierRule)
	if err != nil {
		return nil, err
	}
	return &MonitoringServiceRewardClaimedIterator{contract: _MonitoringService.contract, event: "RewardClaimed", logs: logs, sub: sub}, nil
}

// WatchRewardClaimed is a free log subscription operation binding the contract event 0xe413caa6d70a6d9b51c2af2575a2914490f614355049af8ae7cde5caab9fd201.
//
// Solidity: event RewardClaimed(ms_address indexed address, amount uint256, reward_identifier indexed bytes32)
func (_MonitoringService *MonitoringServiceFilterer) WatchRewardClaimed(opts *bind.WatchOpts, sink chan<- *MonitoringServiceRewardClaimed, ms_address []common.Address, reward_identifier [][32]byte) (event.Subscription, error) {

	var ms_addressRule []interface{}
	for _, ms_addressItem := range ms_address {
		ms_addressRule = append(ms_addressRule, ms_addressItem)
	}

	var reward_identifierRule []interface{}
	for _, reward_identifierItem := range reward_identifier {
		reward_identifierRule = append(reward_identifierRule, reward_identifierItem)
	}

	logs, sub, err := _MonitoringService.contract.WatchLogs(opts, "RewardClaimed", ms_addressRule, reward_identifierRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MonitoringServiceRewardClaimed)
				if err := _MonitoringService.contract.UnpackLog(event, "RewardClaimed", log); err != nil {
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

// MonitoringServiceWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the MonitoringService contract.
type MonitoringServiceWithdrawnIterator struct {
	Event *MonitoringServiceWithdrawn // Event containing the contract specifics and raw log

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
func (it *MonitoringServiceWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MonitoringServiceWithdrawn)
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
		it.Event = new(MonitoringServiceWithdrawn)
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
func (it *MonitoringServiceWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MonitoringServiceWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MonitoringServiceWithdrawn represents a Withdrawn event raised by the MonitoringService contract.
type MonitoringServiceWithdrawn struct {
	Account common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(account indexed address, amount uint256)
func (_MonitoringService *MonitoringServiceFilterer) FilterWithdrawn(opts *bind.FilterOpts, account []common.Address) (*MonitoringServiceWithdrawnIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _MonitoringService.contract.FilterLogs(opts, "Withdrawn", accountRule)
	if err != nil {
		return nil, err
	}
	return &MonitoringServiceWithdrawnIterator{contract: _MonitoringService.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(account indexed address, amount uint256)
func (_MonitoringService *MonitoringServiceFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *MonitoringServiceWithdrawn, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _MonitoringService.contract.WatchLogs(opts, "Withdrawn", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MonitoringServiceWithdrawn)
				if err := _MonitoringService.contract.UnpackLog(event, "Withdrawn", log); err != nil {
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

// RaidenServiceBundleABI is the input ABI used to generate the binding from.
const RaidenServiceBundleABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"contract_address\",\"type\":\"address\"}],\"name\":\"contractExists\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"deposits\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_token_address\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// RaidenServiceBundleBin is the compiled bytecode used for deploying new contracts.
const RaidenServiceBundleBin = `0x608060405234801561001057600080fd5b506040516020806104658339810160405251600160a060020a038116151561003757600080fd5b61004981640100000000610104810204565b151561005457600080fd5b60008054600160a060020a031916600160a060020a0383811691909117808355604080517f18160ddd000000000000000000000000000000000000000000000000000000008152905191909216916318160ddd91600480830192602092919082900301818787803b1580156100c857600080fd5b505af11580156100dc573d6000803e3d6000fd5b505050506040513d60208110156100f257600080fd5b5051116100fe57600080fd5b5061010c565b6000903b1190565b61034a8061011b6000396000f30060806040526004361061006c5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416637709bc788114610071578063b32c65c8146100b3578063b6b55f251461013d578063fc0c546a14610157578063fc7e286d14610195575b600080fd5b34801561007d57600080fd5b5061009f73ffffffffffffffffffffffffffffffffffffffff600435166101d5565b604080519115158252519081900360200190f35b3480156100bf57600080fd5b506100c86101dd565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101025781810151838201526020016100ea565b50505050905090810190601f16801561012f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561014957600080fd5b50610155600435610214565b005b34801561016357600080fd5b5061016c6102f0565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b3480156101a157600080fd5b506101c373ffffffffffffffffffffffffffffffffffffffff6004351661030c565b60408051918252519081900360200190f35b6000903b1190565b60408051808201909152600581527f302e332e5f000000000000000000000000000000000000000000000000000000602082015281565b6000811161022157600080fd5b73ffffffffffffffffffffffffffffffffffffffff3381166000818152600160209081526040808320805487019055825481517f23b872dd000000000000000000000000000000000000000000000000000000008152600481019590955230861660248601526044850187905290519416936323b872dd93606480820194918390030190829087803b1580156102b657600080fd5b505af11580156102ca573d6000803e3d6000fd5b505050506040513d60208110156102e057600080fd5b505115156102ed57600080fd5b50565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b600160205260009081526040902054815600a165627a7a723058204b021abc23d5add4371c9c5ea3a20dadc71c1483ba90cfe4ec4f4ff1583bbd5c0029`

// DeployRaidenServiceBundle deploys a new Ethereum contract, binding an instance of RaidenServiceBundle to it.
func DeployRaidenServiceBundle(auth *bind.TransactOpts, backend bind.ContractBackend, _token_address common.Address) (common.Address, *types.Transaction, *RaidenServiceBundle, error) {
	parsed, err := abi.JSON(strings.NewReader(RaidenServiceBundleABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(RaidenServiceBundleBin), backend, _token_address)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RaidenServiceBundle{RaidenServiceBundleCaller: RaidenServiceBundleCaller{contract: contract}, RaidenServiceBundleTransactor: RaidenServiceBundleTransactor{contract: contract}, RaidenServiceBundleFilterer: RaidenServiceBundleFilterer{contract: contract}}, nil
}

// RaidenServiceBundle is an auto generated Go binding around an Ethereum contract.
type RaidenServiceBundle struct {
	RaidenServiceBundleCaller     // Read-only binding to the contract
	RaidenServiceBundleTransactor // Write-only binding to the contract
	RaidenServiceBundleFilterer   // Log filterer for contract events
}

// RaidenServiceBundleCaller is an auto generated read-only Go binding around an Ethereum contract.
type RaidenServiceBundleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RaidenServiceBundleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RaidenServiceBundleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RaidenServiceBundleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RaidenServiceBundleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RaidenServiceBundleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RaidenServiceBundleSession struct {
	Contract     *RaidenServiceBundle // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// RaidenServiceBundleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RaidenServiceBundleCallerSession struct {
	Contract *RaidenServiceBundleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// RaidenServiceBundleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RaidenServiceBundleTransactorSession struct {
	Contract     *RaidenServiceBundleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// RaidenServiceBundleRaw is an auto generated low-level Go binding around an Ethereum contract.
type RaidenServiceBundleRaw struct {
	Contract *RaidenServiceBundle // Generic contract binding to access the raw methods on
}

// RaidenServiceBundleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RaidenServiceBundleCallerRaw struct {
	Contract *RaidenServiceBundleCaller // Generic read-only contract binding to access the raw methods on
}

// RaidenServiceBundleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RaidenServiceBundleTransactorRaw struct {
	Contract *RaidenServiceBundleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRaidenServiceBundle creates a new instance of RaidenServiceBundle, bound to a specific deployed contract.
func NewRaidenServiceBundle(address common.Address, backend bind.ContractBackend) (*RaidenServiceBundle, error) {
	contract, err := bindRaidenServiceBundle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RaidenServiceBundle{RaidenServiceBundleCaller: RaidenServiceBundleCaller{contract: contract}, RaidenServiceBundleTransactor: RaidenServiceBundleTransactor{contract: contract}, RaidenServiceBundleFilterer: RaidenServiceBundleFilterer{contract: contract}}, nil
}

// NewRaidenServiceBundleCaller creates a new read-only instance of RaidenServiceBundle, bound to a specific deployed contract.
func NewRaidenServiceBundleCaller(address common.Address, caller bind.ContractCaller) (*RaidenServiceBundleCaller, error) {
	contract, err := bindRaidenServiceBundle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RaidenServiceBundleCaller{contract: contract}, nil
}

// NewRaidenServiceBundleTransactor creates a new write-only instance of RaidenServiceBundle, bound to a specific deployed contract.
func NewRaidenServiceBundleTransactor(address common.Address, transactor bind.ContractTransactor) (*RaidenServiceBundleTransactor, error) {
	contract, err := bindRaidenServiceBundle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RaidenServiceBundleTransactor{contract: contract}, nil
}

// NewRaidenServiceBundleFilterer creates a new log filterer instance of RaidenServiceBundle, bound to a specific deployed contract.
func NewRaidenServiceBundleFilterer(address common.Address, filterer bind.ContractFilterer) (*RaidenServiceBundleFilterer, error) {
	contract, err := bindRaidenServiceBundle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RaidenServiceBundleFilterer{contract: contract}, nil
}

// bindRaidenServiceBundle binds a generic wrapper to an already deployed contract.
func bindRaidenServiceBundle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RaidenServiceBundleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RaidenServiceBundle *RaidenServiceBundleRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RaidenServiceBundle.Contract.RaidenServiceBundleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RaidenServiceBundle *RaidenServiceBundleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RaidenServiceBundle.Contract.RaidenServiceBundleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RaidenServiceBundle *RaidenServiceBundleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RaidenServiceBundle.Contract.RaidenServiceBundleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RaidenServiceBundle *RaidenServiceBundleCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RaidenServiceBundle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RaidenServiceBundle *RaidenServiceBundleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RaidenServiceBundle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RaidenServiceBundle *RaidenServiceBundleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RaidenServiceBundle.Contract.contract.Transact(opts, method, params...)
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(contract_address address) constant returns(bool)
func (_RaidenServiceBundle *RaidenServiceBundleCaller) ContractExists(opts *bind.CallOpts, contract_address common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _RaidenServiceBundle.contract.Call(opts, out, "contractExists", contract_address)
	return *ret0, err
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(contract_address address) constant returns(bool)
func (_RaidenServiceBundle *RaidenServiceBundleSession) ContractExists(contract_address common.Address) (bool, error) {
	return _RaidenServiceBundle.Contract.ContractExists(&_RaidenServiceBundle.CallOpts, contract_address)
}

// ContractExists is a free data retrieval call binding the contract method 0x7709bc78.
//
// Solidity: function contractExists(contract_address address) constant returns(bool)
func (_RaidenServiceBundle *RaidenServiceBundleCallerSession) ContractExists(contract_address common.Address) (bool, error) {
	return _RaidenServiceBundle.Contract.ContractExists(&_RaidenServiceBundle.CallOpts, contract_address)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_RaidenServiceBundle *RaidenServiceBundleCaller) Contract_version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _RaidenServiceBundle.contract.Call(opts, out, "contract_version")
	return *ret0, err
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_RaidenServiceBundle *RaidenServiceBundleSession) Contract_version() (string, error) {
	return _RaidenServiceBundle.Contract.Contract_version(&_RaidenServiceBundle.CallOpts)
}

// Contract_version is a free data retrieval call binding the contract method 0xb32c65c8.
//
// Solidity: function contract_version() constant returns(string)
func (_RaidenServiceBundle *RaidenServiceBundleCallerSession) Contract_version() (string, error) {
	return _RaidenServiceBundle.Contract.Contract_version(&_RaidenServiceBundle.CallOpts)
}

// Deposits is a free data retrieval call binding the contract method 0xfc7e286d.
//
// Solidity: function deposits( address) constant returns(uint256)
func (_RaidenServiceBundle *RaidenServiceBundleCaller) Deposits(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RaidenServiceBundle.contract.Call(opts, out, "deposits", arg0)
	return *ret0, err
}

// Deposits is a free data retrieval call binding the contract method 0xfc7e286d.
//
// Solidity: function deposits( address) constant returns(uint256)
func (_RaidenServiceBundle *RaidenServiceBundleSession) Deposits(arg0 common.Address) (*big.Int, error) {
	return _RaidenServiceBundle.Contract.Deposits(&_RaidenServiceBundle.CallOpts, arg0)
}

// Deposits is a free data retrieval call binding the contract method 0xfc7e286d.
//
// Solidity: function deposits( address) constant returns(uint256)
func (_RaidenServiceBundle *RaidenServiceBundleCallerSession) Deposits(arg0 common.Address) (*big.Int, error) {
	return _RaidenServiceBundle.Contract.Deposits(&_RaidenServiceBundle.CallOpts, arg0)
}

// TokenNetworkAddres is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_RaidenServiceBundle *RaidenServiceBundleCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _RaidenServiceBundle.contract.Call(opts, out, "token")
	return *ret0, err
}

// TokenNetworkAddres is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_RaidenServiceBundle *RaidenServiceBundleSession) Token() (common.Address, error) {
	return _RaidenServiceBundle.Contract.Token(&_RaidenServiceBundle.CallOpts)
}

// TokenNetworkAddres is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_RaidenServiceBundle *RaidenServiceBundleCallerSession) Token() (common.Address, error) {
	return _RaidenServiceBundle.Contract.Token(&_RaidenServiceBundle.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(amount uint256) returns()
func (_RaidenServiceBundle *RaidenServiceBundleTransactor) Deposit(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _RaidenServiceBundle.contract.Transact(opts, "deposit", amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(amount uint256) returns()
func (_RaidenServiceBundle *RaidenServiceBundleSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _RaidenServiceBundle.Contract.Deposit(&_RaidenServiceBundle.TransactOpts, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(amount uint256) returns()
func (_RaidenServiceBundle *RaidenServiceBundleTransactorSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _RaidenServiceBundle.Contract.Deposit(&_RaidenServiceBundle.TransactOpts, amount)
}

// SecretRegistryABI is the input ABI used to generate the binding from.
const SecretRegistryABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"secret\",\"type\":\"bytes32\"}],\"name\":\"registerSecret\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"secrethash_to_block\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"secrethash\",\"type\":\"bytes32\"}],\"name\":\"getSecretRevealBlockHeight\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"secrethash\",\"type\":\"bytes32\"}],\"name\":\"SecretRevealed\",\"type\":\"event\"}]"

// SecretRegistryBin is the compiled bytecode used for deploying new contracts.
const SecretRegistryBin = `0x608060405234801561001057600080fd5b506102cd806100206000396000f3006080604052600436106100615763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166312ad8bfc81146100665780639734030914610092578063b32c65c8146100bc578063c1f6294614610146575b600080fd5b34801561007257600080fd5b5061007e60043561015e565b604080519115158252519081900360200190f35b34801561009e57600080fd5b506100aa600435610246565b60408051918252519081900360200190f35b3480156100c857600080fd5b506100d1610258565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561010b5781810151838201526020016100f3565b50505050905090810190601f1680156101385780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561015257600080fd5b506100aa60043561028f565b604080516020808201849052825180830382018152918301928390528151600093849392909182918401908083835b602083106101ac5780518252601f19909201916020918201910161018d565b5181516020939093036101000a60001901801990911692169190911790526040519201829003909120935050841591508190506101f55750600081815260208190526040812054115b156102035760009150610240565b6000818152602081905260408082204390555182917f9b7ddc883342824bd7ddbff103e7a69f8f2e60b96c075cd1b8b8b9713ecc75a491a2600191505b50919050565b60006020819052908152604090205481565b60408051808201909152600581527f302e332e5f000000000000000000000000000000000000000000000000000000602082015281565b600090815260208190526040902054905600a165627a7a723058202e716a9ffb556c366a7bccd539d2f00eed4e6fe00a0ab33a1130d8c4807953750029`

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

// DeployToken deploys a new Ethereum contract, binding an instance of TokenNetworkAddres to it.
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

// TokenNetworkAddres is an auto generated Go binding around an Ethereum contract.
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

// NewToken creates a new instance of TokenNetworkAddres, bound to a specific deployed contract.
func NewToken(address common.Address, backend bind.ContractBackend) (*Token, error) {
	contract, err := bindToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Token{TokenCaller: TokenCaller{contract: contract}, TokenTransactor: TokenTransactor{contract: contract}, TokenFilterer: TokenFilterer{contract: contract}}, nil
}

// NewTokenCaller creates a new read-only instance of TokenNetworkAddres, bound to a specific deployed contract.
func NewTokenCaller(address common.Address, caller bind.ContractCaller) (*TokenCaller, error) {
	contract, err := bindToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCaller{contract: contract}, nil
}

// NewTokenTransactor creates a new write-only instance of TokenNetworkAddres, bound to a specific deployed contract.
func NewTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenTransactor, error) {
	contract, err := bindToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenTransactor{contract: contract}, nil
}

// NewTokenFilterer creates a new log filterer instance of TokenNetworkAddres, bound to a specific deployed contract.
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

// TokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the TokenNetworkAddres contract.
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

// TokenApproval represents a Approval event raised by the TokenNetworkAddres contract.
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

// TokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the TokenNetworkAddres contract.
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

// TokenTransfer represents a Transfer event raised by the TokenNetworkAddres contract.
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
const TokenNetworkABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"participant1\",\"type\":\"address\"},{\"name\":\"participant2\",\"type\":\"address\"},{\"name\":\"settle_timeout\",\"type\":\"uint256\"}],\"name\":\"openChannel\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"secret_registry\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"participant\",\"type\":\"address\"},{\"name\":\"partner\",\"type\":\"address\"},{\"name\":\"merkle_tree_leaves\",\"type\":\"bytes\"}],\"name\":\"unlock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"chain_id\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAX_SAFE_UINT256\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"contract_address\",\"type\":\"address\"}],\"name\":\"contractExists\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"participant1\",\"type\":\"address\"},{\"name\":\"participant1_transferred_amount\",\"type\":\"uint256\"},{\"name\":\"participant1_locked_amount\",\"type\":\"uint256\"},{\"name\":\"participant1_locksroot\",\"type\":\"bytes32\"},{\"name\":\"participant2\",\"type\":\"address\"},{\"name\":\"participant2_transferred_amount\",\"type\":\"uint256\"},{\"name\":\"participant2_locked_amount\",\"type\":\"uint256\"},{\"name\":\"participant2_locksroot\",\"type\":\"bytes32\"}],\"name\":\"settleChannel\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"channels\",\"outputs\":[{\"name\":\"settle_block_number\",\"type\":\"uint256\"},{\"name\":\"state\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"participant1_address\",\"type\":\"address\"},{\"name\":\"participant1_balance\",\"type\":\"uint256\"},{\"name\":\"participant2_address\",\"type\":\"address\"},{\"name\":\"participant2_balance\",\"type\":\"uint256\"},{\"name\":\"participant1_signature\",\"type\":\"bytes\"},{\"name\":\"participant2_signature\",\"type\":\"bytes\"}],\"name\":\"cooperativeSettle\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"participant\",\"type\":\"address\"},{\"name\":\"partner\",\"type\":\"address\"}],\"name\":\"getChannelIdentifier\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"partner\",\"type\":\"address\"},{\"name\":\"balance_hash\",\"type\":\"bytes32\"},{\"name\":\"nonce\",\"type\":\"uint256\"},{\"name\":\"additional_hash\",\"type\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"closeChannel\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"participant\",\"type\":\"address\"},{\"name\":\"total_deposit\",\"type\":\"uint256\"},{\"name\":\"partner\",\"type\":\"address\"}],\"name\":\"setTotalDeposit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"participant\",\"type\":\"address\"},{\"name\":\"partner\",\"type\":\"address\"}],\"name\":\"getChannelParticipantInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"bool\"},{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"closing_participant\",\"type\":\"address\"},{\"name\":\"non_closing_participant\",\"type\":\"address\"},{\"name\":\"balance_hash\",\"type\":\"bytes32\"},{\"name\":\"nonce\",\"type\":\"uint256\"},{\"name\":\"additional_hash\",\"type\":\"bytes32\"},{\"name\":\"closing_signature\",\"type\":\"bytes\"},{\"name\":\"non_closing_signature\",\"type\":\"bytes\"}],\"name\":\"updateNonClosingBalanceProof\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contract_version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"participant\",\"type\":\"address\"},{\"name\":\"total_withdraw\",\"type\":\"uint256\"},{\"name\":\"partner\",\"type\":\"address\"},{\"name\":\"participant_signature\",\"type\":\"bytes\"},{\"name\":\"partner_signature\",\"type\":\"bytes\"}],\"name\":\"setTotalWithdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"participant1\",\"type\":\"address\"},{\"name\":\"participant2\",\"type\":\"address\"}],\"name\":\"getChannelInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"participant1\",\"type\":\"address\"},{\"name\":\"participant2\",\"type\":\"address\"},{\"name\":\"locksroot\",\"type\":\"bytes32\"}],\"name\":\"getParticipantLockedAmount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_token_address\",\"type\":\"address\"},{\"name\":\"_secret_registry\",\"type\":\"address\"},{\"name\":\"_chain_id\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channel_identifier\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"participant1\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"participant2\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"settle_timeout\",\"type\":\"uint256\"}],\"name\":\"ChannelOpened\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channel_identifier\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"participant\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"total_deposit\",\"type\":\"uint256\"}],\"name\":\"ChannelNewDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channel_identifier\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"participant\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"total_withdraw\",\"type\":\"uint256\"}],\"name\":\"ChannelWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channel_identifier\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"closing_participant\",\"type\":\"address\"}],\"name\":\"ChannelClosed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channel_identifier\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"participant\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"unlocked_amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"returned_tokens\",\"type\":\"uint256\"}],\"name\":\"ChannelUnlocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channel_identifier\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"closing_participant\",\"type\":\"address\"}],\"name\":\"NonClosingBalanceProofUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channel_identifier\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"participant1_amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"participant2_amount\",\"type\":\"uint256\"}],\"name\":\"ChannelSettled\",\"type\":\"event\"}]"

// TokenNetworkBin is the compiled bytecode used for deploying new contracts.
const TokenNetworkBin = `0x60806040523480156200001157600080fd5b5060405160608062002929833981016040908152815160208301519190920151600160a060020a03831615156200004757600080fd5b600160a060020a03821615156200005d57600080fd5b600081116200006b57600080fd5b6200007f8364010000000062000177810204565b15156200008b57600080fd5b6200009f8264010000000062000177810204565b1515620000ab57600080fd5b60008054600160a060020a03808616600160a060020a031992831617808455600180548784169416939093179092556002849055604080517f18160ddd000000000000000000000000000000000000000000000000000000008152905192909116916318160ddd9160048082019260209290919082900301818787803b1580156200013557600080fd5b505af11580156200014a573d6000803e3d6000fd5b505050506040513d60208110156200016157600080fd5b5051116200016e57600080fd5b5050506200017f565b6000903b1190565b61279a806200018f6000396000f3006080604052600436106100ed5763ffffffff60e060020a6000350416630a798f2481146100f257806324d73a931461012e578063331d8e5d1461015f5780633af973b1146101d057806371e75992146101e55780637709bc78146101fa578063799786301461022f5780637a7ebd7b1461026b5780638568536a1461029e578063938bcd67146103555780639abe275f1461037c578063a32a6737146103f0578063ac1337091461041b578063aec1dd811461046d578063b32c65c8146104ba578063c472c7e614610544578063f94c9e1314610588578063fc0c546a146105d0578063fd5f1e03146105e5575b600080fd5b3480156100fe57600080fd5b5061011c600160a060020a036004358116906024351660443561060f565b60408051918252519081900360200190f35b34801561013a57600080fd5b506101436106d3565b60408051600160a060020a039092168252519081900360200190f35b34801561016b57600080fd5b50604080516020600460443581810135601f81018490048402850184019095528484526101ce948235600160a060020a03908116956024803590921695369594606494929301919081908401838280828437509497506106e29650505050505050565b005b3480156101dc57600080fd5b5061011c61095e565b3480156101f157600080fd5b5061011c610964565b34801561020657600080fd5b5061021b600160a060020a036004351661096a565b604080519115158252519081900360200190f35b34801561023b57600080fd5b506101ce600160a060020a036004358116906024359060443590606435906084351660a43560c43560e435610972565b34801561027757600080fd5b50610283600435610ce4565b6040805192835260ff90911660208301528051918290030190f35b3480156102aa57600080fd5b50604080516020601f6084356004818101359283018490048402850184019095528184526101ce94600160a060020a0381358116956024803596604435909316956064359536959460a49493919091019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a999881019791965091820194509250829150840183828082843750949750610d009650505050505050565b34801561036157600080fd5b5061011c600160a060020a0360043581169060243516611011565b34801561038857600080fd5b50604080516020601f6084356004818101359283018490048402850184019095528184526101ce94600160a060020a03813516946024803595604435956064359536959460a49490939101919081908401838280828437509497506111b79650505050505050565b3480156103fc57600080fd5b506101ce600160a060020a0360043581169060243590604435166112cb565b34801561042757600080fd5b50610442600160a060020a03600435811690602435166114eb565b6040805195865260208601949094529115158484015260608401526080830152519081900360a00190f35b34801561047957600080fd5b506101ce600160a060020a0360048035821691602480359091169160443591606435916084359160a43580820192908101359160c435908101910135611555565b3480156104c657600080fd5b506104cf611719565b6040805160208082528351818301528351919283929083019185019080838360005b838110156105095781810151838201526020016104f1565b50505050905090810190601f1680156105365780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561055057600080fd5b506101ce60048035600160a060020a039081169160248035926044351691606435808301929082013591608435918201910135611750565b34801561059457600080fd5b506105af600160a060020a0360043581169060243516611a05565b60408051938452602084019290925260ff1682820152519081900360600190f35b3480156105dc57600080fd5b50610143611a41565b3480156105f157600080fd5b5061011c600160a060020a0360043581169060243516604435611a50565b600080600083600681101580156106295750622932e08111155b151561063457600080fd5b61063e8787611011565b6000818152600360205260409020805491945092501561065d57600080fd5b600182015460ff161561066f57600080fd5b8482556001808301805460ff19169091179055604080518681529051600160a060020a0380891692908a169186917f448d27f1fe12f92a2070111296e68fd6ef0a01c0e05bf5819eda0dbcf267bf3d919081900360200190a4509095945050505050565b600154600160a060020a031681565b600080600080600080600087511115156106fb57600080fd5b6107058989611011565b955061071087611a84565b6000888152600460209081526040808320858452918290529091205492975090955090935090506107418484611dd8565b6000868152602083905260408120819055909450848403925084111561080a5760008054604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a038d81166004830152602482018990529151919092169263a9059cbb92604480820193602093909283900390910190829087803b1580156107d357600080fd5b505af11580156107e7573d6000803e3d6000fd5b505050506040513d60208110156107fd57600080fd5b5051151561080a57600080fd5b60008211156108bc5760008054604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a038c81166004830152602482018790529151919092169263a9059cbb92604480820193602093909283900390910190829087803b15801561088557600080fd5b505af1158015610899573d6000803e3d6000fd5b505050506040513d60208110156108af57600080fd5b505115156108bc57600080fd5b60408051858152602081018490528151600160a060020a038c169289927f6f2d495eefa4b2d91a2287258f6f88722cabdaf48f4d1410d979e38b516a258f929081900390910190a360008681526003602052604090206001015460ff161561092357600080fd5b84151561092f57600080fd5b6000831161093c57600080fd5b8183101561094957600080fd5b8383101561095357fe5b505050505050505050565b60025481565b60001981565b6000903b1190565b6000806000806109828c89611011565b60008181526003602052604090206001810154919550935060ff166002146109a957600080fd5b825443116109b657600080fd5b5050600160a060020a03808b166000908152600283016020526040808220928916825290206109e7828c8c8c611df0565b15156109f257600080fd5b6109fe81888888611df0565b1515610a0957600080fd5b610a17828c8c848b8b611ebb565b809950819d50829a50839e50505050508260020160008d600160a060020a0316600160a060020a0316815260200190815260200160002060008082016000905560018201600090556002820160006101000a81549060ff021916905560038201600090556004820160009055505082600201600089600160a060020a0316600160a060020a0316815260200190815260200160002060008082016000905560018201600090556002820160006101000a81549060ff02191690556003820160009055600482016000905550506003600085600019166000191681526020019081526020016000206000808201600090556001820160006101000a81549060ff02191690555050610b28848b8b611fb5565b610b33848787611fb5565b60008b1115610bea576000809054906101000a9004600160a060020a0316600160a060020a031663a9059cbb8d8d6040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b158015610bb357600080fd5b505af1158015610bc7573d6000803e3d6000fd5b505050506040513d6020811015610bdd57600080fd5b50511515610bea57600080fd5b6000871115610c9c5760008054604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a038c81166004830152602482018c90529151919092169263a9059cbb92604480820193602093909283900390910190829087803b158015610c6557600080fd5b505af1158015610c79573d6000803e3d6000fd5b505050506040513d6020811015610c8f57600080fd5b50511515610c9c57600080fd5b604080518c815260208101899052815186927ff94fb5c0628a82dc90648e8dc5e983f632633b0d26603d64e8cc042ca0790aa4928290030190a2505050505050505050505050565b6003602052600090815260409020805460019091015460ff1682565b600080600080600080600080610d168e8d611011565b6000818152600360205260409020600181015491995060ff90911694509250610d43888f8f8f8f8f611ff2565b9650610d53888f8f8f8f8e611ff2565b600160a060020a0380891660009081526002860160205260408082209284168252902091975092509050610d8782826120d2565b600160a060020a0380891660009081526002808701602090815260408084208481556001808201869055818501805460ff1990811690915560038084018890556004938401889055978f168752838720878155808301889055958601805482169055858801879055949091018590558e855294909152822082815590920180549092169091559095508d1115610ec5576000809054906101000a9004600160a060020a0316600160a060020a031663a9059cbb888f6040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b158015610e8e57600080fd5b505af1158015610ea2573d6000803e3d6000fd5b505050506040513d6020811015610eb857600080fd5b50511515610ec557600080fd5b60008b1115610f7c576000809054906101000a9004600160a060020a0316600160a060020a031663a9059cbb878d6040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b158015610f4557600080fd5b505af1158015610f59573d6000803e3d6000fd5b505050506040513d6020811015610f6f57600080fd5b50511515610f7c57600080fd5b60018414610f8957600080fd5b600160a060020a03878116908f1614610fa157600080fd5b600160a060020a03868116908d1614610fb957600080fd5b8c8b018514610fc757600080fd5b604080518e8152602081018d905281518a927ff94fb5c0628a82dc90648e8dc5e983f632633b0d26603d64e8cc042ca0790aa4928290030190a25050505050505050505050505050565b6000600160a060020a038316151561102857600080fd5b600160a060020a038216151561103d57600080fd5b600160a060020a03838116908316141561105657600080fd5b81600160a060020a031683600160a060020a0316101561112b5782826040516020018083600160a060020a0316600160a060020a0316606060020a02815260140182600160a060020a0316600160a060020a0316606060020a028152601401925050506040516020818303038152906040526040518082805190602001908083835b602083106110f75780518252601f1990920191602091820191016110d8565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902090506111b1565b81836040516020018083600160a060020a0316600160a060020a0316606060020a02815260140182600160a060020a0316600160a060020a0316606060020a02815260140192505050604051602081830303815290604052604051808280519060200190808383602083106110f75780518252601f1990920191602091820191016110d8565b92915050565b6000806000338860006111ca8383611011565b600081815260036020526040902060019081015491925060ff909116146111f057600080fd5b6111fa338c611011565b600081815260036020908152604080832060018082018054600260ff19918216811790925533600160a060020a03168752818401909552928520909201805490931690911790915580544301815591965090945089111561128857611262858b8b8b8b6120e8565b955061127084878b8d6121b3565b600160a060020a038b81169087161461128857600080fd5b604051600160a060020a0333169086907f16a93f97197f719d19b0258648f368a06980009827e4a4a88892dd761ba4017c90600090a35050505050505050505050565b600080600080600080888760006112e28383611011565b600081815260036020526040902060019081015491925060ff9091161461130857600080fd5b60008b1161131557600080fd5b61131f8c8b611011565b9850600360008a6000191660001916815260200190815260200160002095508560020160008d600160a060020a0316600160a060020a0316815260200190815260200160002094508560020160008b600160a060020a0316600160a060020a03168152602001908152602001600020935084600001548b039750878560000160008282540192505081905550836000015485600001540196508b600160a060020a031689600019167f0346e981e2bfa2366dc2307a8f1fa24779830a01121b1275fe565c6b98bb4d3487600001546040518082815260200191505060405180910390a360008054604080517f23b872dd000000000000000000000000000000000000000000000000000000008152600160a060020a0333811660048301523081166024830152604482018d9052915191909216926323b872dd92606480820193602093909283900390910190829087803b15801561147c57600080fd5b505af1158015611490573d6000803e3d6000fd5b505050506040513d60208110156114a657600080fd5b505115156114b357600080fd5b84548811156114c157600080fd5b84548710156114cf57600080fd5b83548710156114dd57600080fd5b505050505050505050505050565b60008060008060008060006115008989611011565b6000908152600360208181526040808420600160a060020a039d909d16845260029c8d01909152909120805460018201549b82015492820154600490920154909c60ff9093169a509098509650945050505050565b6000808080808b151561156757600080fd5b60008b1161157457600080fd5b61157e8e8e611011565b6000818152600360209081526040918290208251601f8d018390048302810183019093528b83529297509193506116049187918f918f918f918f908f90819084018382808284378201915050505050508c8c8080601f016020809104026020016040519081016040528093929190818152602001838380828437506121ef945050505050565b9350611643858d8d8d8d8d8080601f016020809104026020016040519081016040528093929190818152602001838380828437506120e8945050505050565b600160a060020a038f1660009081526002840160205260409020909350905061166e828f8d8f6121b3565b604051600160a060020a038f169086907f3558501c5224e6a6801ce5e7d81ed6f58489f1f7a88367cb14f2144b057f0b8090600090a3600182015460ff166002146116b857600080fd5b81544311156116c657600080fd5b600281015460ff1615156116d957600080fd5b600160a060020a038e8116908416146116f157600080fd5b600160a060020a038d81169085161461170957600080fd5b5050505050505050505050505050565b60408051808201909152600581527f302e332e5f000000000000000000000000000000000000000000000000000000602082015281565b60008080808080808c1161176357600080fd5b61176d8d8c611011565b955060036000876000191660001916815260200190815260200160002092508260020160008e600160a060020a0316600160a060020a0316815260200190815260200160002091508260020160008c600160a060020a0316600160a060020a0316815260200190815260200160002090508060000154826000015401945081600101548c0393508382600101600082825401925050819055506000809054906101000a9004600160a060020a0316600160a060020a031663a9059cbb8e866040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b15801561187d57600080fd5b505af1158015611891573d6000803e3d6000fd5b505050506040513d60208110156118a757600080fd5b505115156118b457600080fd5b60018084015460ff16146118c757600080fd5b611937868e8d8f8e8e8080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050508d8d8080601f016020809104026020016040519081016040528093929190818152602001838380828437506122e9945050505050565b6000841161194457600080fd5b600182015484111561195557600080fd5b60018201548c1461196557600080fd5b8060010154850382600101541115151561197e57600080fd5b815485101561198c57600080fd5b805485101561199a57600080fd5b6004820154156119a657fe5b6004810154156119b257fe5b60018201546040805191825251600160a060020a038f169188917f5860c94079516621b44c52f4423fd883c5f8f0370d5343f20f5dcb5b837738209181900360200190a350505050505050505050505050565b6000806000806000611a178787611011565b60008181526003602052604090208054600190910154919990985060ff9091169650945050505050565b600054600160a060020a031681565b6000806000611a5f8686611011565b6000908152600460209081526040808320968352959052939093205495945050505050565b805160009081908180808080606080870615611a9f57600080fd5b60608704600101604051908082528060200260200182016040528015611acf578160200160208202803883390190505b509050602095505b86861015611b1957611ae98a87612342565b9586019594509250828160608804815181101515611b0357fe5b6020908102909101015260609590950194611ad7565b6060870496505b6001871115611daf576002870615611b6d578060018803815181101515611b4357fe5b906020019060200201518188815181101515611b5b57fe5b60209081029091010152600196909601955b600095505b60018703861015611da4578086600101815181101515611b8e57fe5b602090810290910101518151829088908110611ba657fe5b602090810290910101511415611bd5578086815181101515611bc457fe5b906020019060200201519250611d7c565b8086600101815181101515611be657fe5b602090810290910101518151829088908110611bfe57fe5b602090810290910101511015611cc7578086815181101515611c1c57fe5b906020019060200201518187600101815181101515611c3757fe5b6020908102909101810151604080518084019490945283810191909152805180840382018152606090930190819052825190918291908401908083835b60208310611c935780518252601f199092019160209182019101611c74565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390209250611d7c565b8086600101815181101515611cd857fe5b906020019060200201518187815181101515611cf057fe5b6020908102909101810151604080518084019490945283810191909152805180840382018152606090930190819052825190918291908401908083835b60208310611d4c5780518252601f199092019160209182019101611d2d565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902092505b828160028804815181101515611d8e57fe5b6020908102909101015260029590950194611b72565b600286049650611b20565b806000815181101515611dbe57fe5b602090810290910101519a94995093975050505050505050565b6000818311611de75782611de9565b815b9392505050565b6003840154600090158015611e03575083155b8015611e0d575082155b8015611e17575081155b15611e2457506001611eb3565b604080516020808201879052818301869052606080830186905283518084039091018152608090920192839052815191929182918401908083835b60208310611e7e5780518252601f199092019160209182019101611e5f565b5181516020939093036101000a6000190180199091169216919091179052604051920182900390912060038901541493505050505b949350505050565b6000806000806000806000611ece612745565b611ed6612745565b8e600001548260000181815250508e600101548260200181815250508d8260400181815250508c8260600181815250508b600001548160000181815250508b600101548160200181815250508a81604001818152505089816060018181525050611f408f8d6120d2565b9250611f4c82826124b7565b9450611f588584611dd8565b94508483039350611f69858b61253f565b9a509450611f77848e61253f565b9d50935082851115611f8557fe5b82841115611f8f57fe5b8484018d018a018314611f9e57fe5b50929d919c50999a50959850949650505050505050565b6000821580611fc2575081155b15611fcc57611fec565b506000838152600460209081526040808320848452918290529091208390555b50505050565b60025460408051600160a060020a03888116606060020a908102602080850191909152603484018a90528883168202605485015260688401889052608884018c9052309092160260a883015260bc808301949094528251808303909401845260dc90910191829052825160009384939092909182918401908083835b6020831061208d5780518252601f19909201916020918201910161206e565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902090506120c68184612561565b98975050505050505050565b6001808201549083015491549254909201030390565b6002546040805160208082018890528183018790526060820186905260808201899052606060020a600160a060020a0330160260a083015260b4808301949094528251808303909401845260d490910191829052825160009384939092909182918401908083835b6020831061216f5780518252601f199092019160209182019101612150565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902090506121a88184612561565b979650505050505050565b600160a060020a03831660009081526002850160205260409020600481015483116121dd57600080fd5b60048101929092556003909101555050565b6000808686868a30600254896040516020018088600019166000191681526020018781526020018660001916600019168152602001856000191660001916815260200184600160a060020a0316600160a060020a0316606060020a02815260140183815260200182805190602001908083835b602083106122815780518252601f199092019160209182019101612262565b6001836020036101000a0380198251168184511680821785525050505050509050019750505050505050506040516020818303038152906040526040518082805190602001908083836020831061208d5780518252601f19909201916020918201910161206e565b6000806122f888888787612641565b915061230688888786612641565b9050600160a060020a038781169083161461232057600080fd5b600160a060020a038681169082161461233857600080fd5b5050505050505050565b600080600080600080600087895111151561236357955060009450856124ab565b888801805160208083015160409384015184518084018590528086018390526060808201839052865180830390910181526080909101958690528051949a509198509550929182918401908083835b602083106123d15780518252601f1990920191602091820191016123b2565b51815160209384036101000a6000190180199092169116179052604080519290940182900382206001547fc1f62946000000000000000000000000000000000000000000000000000000008452600484018a90529451909750600160a060020a03909416955063c1f62946945060248083019491935090918290030181600087803b15801561245f57600080fd5b505af1158015612473573d6000803e3d6000fd5b505050506040513d602081101561248957600080fd5b5051925082158061249a5750828511155b156124a457600093505b8084965096505b50505050509250929050565b60008060008060006124d187604001518860600151612730565b93506124e586604001518760600151612730565b9250838310156124f457600080fd5b604087015184101561250257fe5b604086015183101561251057fe5b8383039150612523828860000151612730565b905061253381886020015161253f565b50979650505050505050565b60008082841161255157600084612556565b828403835b915091509250929050565b6000806000808451604114151561257757600080fd5b50505060208201516040830151606084015160001a601b60ff8216101561259c57601b015b8060ff16601b14806125b157508060ff16601c145b15156125bc57600080fd5b60408051600080825260208083018085528a905260ff8516838501526060830187905260808301869052925160019360a0808501949193601f19840193928390039091019190865af1158015612616573d6000803e3d6000fd5b5050604051601f190151945050600160a060020a038416151561263857600080fd5b50505092915050565b600080848487306002546040516020018086600160a060020a0316600160a060020a0316606060020a028152601401858152602001846000191660001916815260200183600160a060020a0316600160a060020a0316606060020a028152601401828152602001955050505050506040516020818303038152906040526040518082805190602001908083835b602083106126ed5780518252601f1990920191602091820191016126ce565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902090506127268184612561565b9695505050505050565b600082820183811015611de957600019611eb3565b6080604051908101604052806000815260200160008152602001600081526020016000815250905600a165627a7a7230582086b23fd1b0610d7f063139f41cde9d6dd5d9b7c5aafd8c0ffa57fcbd753aa11b0029`

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

// MAX_SAFE_UINT256 is a free data retrieval call binding the contract method 0x71e75992.
//
// Solidity: function MAX_SAFE_UINT256() constant returns(uint256)
func (_TokenNetwork *TokenNetworkCaller) MAX_SAFE_UINT256(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenNetwork.contract.Call(opts, out, "MAX_SAFE_UINT256")
	return *ret0, err
}

// MAX_SAFE_UINT256 is a free data retrieval call binding the contract method 0x71e75992.
//
// Solidity: function MAX_SAFE_UINT256() constant returns(uint256)
func (_TokenNetwork *TokenNetworkSession) MAX_SAFE_UINT256() (*big.Int, error) {
	return _TokenNetwork.Contract.MAX_SAFE_UINT256(&_TokenNetwork.CallOpts)
}

// MAX_SAFE_UINT256 is a free data retrieval call binding the contract method 0x71e75992.
//
// Solidity: function MAX_SAFE_UINT256() constant returns(uint256)
func (_TokenNetwork *TokenNetworkCallerSession) MAX_SAFE_UINT256() (*big.Int, error) {
	return _TokenNetwork.Contract.MAX_SAFE_UINT256(&_TokenNetwork.CallOpts)
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
// Solidity: function channels( bytes32) constant returns(settle_block_number uint256, state uint8)
func (_TokenNetwork *TokenNetworkCaller) Channels(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Settle_block_number *big.Int
	State               uint8
}, error) {
	ret := new(struct {
		Settle_block_number *big.Int
		State               uint8
	})
	out := ret
	err := _TokenNetwork.contract.Call(opts, out, "channels", arg0)
	return *ret, err
}

// Channels is a free data retrieval call binding the contract method 0x7a7ebd7b.
//
// Solidity: function channels( bytes32) constant returns(settle_block_number uint256, state uint8)
func (_TokenNetwork *TokenNetworkSession) Channels(arg0 [32]byte) (struct {
	Settle_block_number *big.Int
	State               uint8
}, error) {
	return _TokenNetwork.Contract.Channels(&_TokenNetwork.CallOpts, arg0)
}

// Channels is a free data retrieval call binding the contract method 0x7a7ebd7b.
//
// Solidity: function channels( bytes32) constant returns(settle_block_number uint256, state uint8)
func (_TokenNetwork *TokenNetworkCallerSession) Channels(arg0 [32]byte) (struct {
	Settle_block_number *big.Int
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

// GetChannelIdentifier is a free data retrieval call binding the contract method 0x938bcd67.
//
// Solidity: function getChannelIdentifier(participant address, partner address) constant returns(bytes32)
func (_TokenNetwork *TokenNetworkCaller) GetChannelIdentifier(opts *bind.CallOpts, participant common.Address, partner common.Address) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _TokenNetwork.contract.Call(opts, out, "getChannelIdentifier", participant, partner)
	return *ret0, err
}

// GetChannelIdentifier is a free data retrieval call binding the contract method 0x938bcd67.
//
// Solidity: function getChannelIdentifier(participant address, partner address) constant returns(bytes32)
func (_TokenNetwork *TokenNetworkSession) GetChannelIdentifier(participant common.Address, partner common.Address) ([32]byte, error) {
	return _TokenNetwork.Contract.GetChannelIdentifier(&_TokenNetwork.CallOpts, participant, partner)
}

// GetChannelIdentifier is a free data retrieval call binding the contract method 0x938bcd67.
//
// Solidity: function getChannelIdentifier(participant address, partner address) constant returns(bytes32)
func (_TokenNetwork *TokenNetworkCallerSession) GetChannelIdentifier(participant common.Address, partner common.Address) ([32]byte, error) {
	return _TokenNetwork.Contract.GetChannelIdentifier(&_TokenNetwork.CallOpts, participant, partner)
}

// GetChannelInfo is a free data retrieval call binding the contract method 0xf94c9e13.
//
// Solidity: function getChannelInfo(participant1 address, participant2 address) constant returns(bytes32, uint256, uint8)
func (_TokenNetwork *TokenNetworkCaller) GetChannelInfo(opts *bind.CallOpts, participant1 common.Address, participant2 common.Address) ([32]byte, *big.Int, uint8, error) {
	var (
		ret0 = new([32]byte)
		ret1 = new(*big.Int)
		ret2 = new(uint8)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _TokenNetwork.contract.Call(opts, out, "getChannelInfo", participant1, participant2)
	return *ret0, *ret1, *ret2, err
}

// GetChannelInfo is a free data retrieval call binding the contract method 0xf94c9e13.
//
// Solidity: function getChannelInfo(participant1 address, participant2 address) constant returns(bytes32, uint256, uint8)
func (_TokenNetwork *TokenNetworkSession) GetChannelInfo(participant1 common.Address, participant2 common.Address) ([32]byte, *big.Int, uint8, error) {
	return _TokenNetwork.Contract.GetChannelInfo(&_TokenNetwork.CallOpts, participant1, participant2)
}

// GetChannelInfo is a free data retrieval call binding the contract method 0xf94c9e13.
//
// Solidity: function getChannelInfo(participant1 address, participant2 address) constant returns(bytes32, uint256, uint8)
func (_TokenNetwork *TokenNetworkCallerSession) GetChannelInfo(participant1 common.Address, participant2 common.Address) ([32]byte, *big.Int, uint8, error) {
	return _TokenNetwork.Contract.GetChannelInfo(&_TokenNetwork.CallOpts, participant1, participant2)
}

// GetChannelParticipantInfo is a free data retrieval call binding the contract method 0xac133709.
//
// Solidity: function getChannelParticipantInfo(participant address, partner address) constant returns(uint256, uint256, bool, bytes32, uint256)
func (_TokenNetwork *TokenNetworkCaller) GetChannelParticipantInfo(opts *bind.CallOpts, participant common.Address, partner common.Address) (*big.Int, *big.Int, bool, [32]byte, *big.Int, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(*big.Int)
		ret2 = new(bool)
		ret3 = new([32]byte)
		ret4 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
		ret4,
	}
	err := _TokenNetwork.contract.Call(opts, out, "getChannelParticipantInfo", participant, partner)
	return *ret0, *ret1, *ret2, *ret3, *ret4, err
}

// GetChannelParticipantInfo is a free data retrieval call binding the contract method 0xac133709.
//
// Solidity: function getChannelParticipantInfo(participant address, partner address) constant returns(uint256, uint256, bool, bytes32, uint256)
func (_TokenNetwork *TokenNetworkSession) GetChannelParticipantInfo(participant common.Address, partner common.Address) (*big.Int, *big.Int, bool, [32]byte, *big.Int, error) {
	return _TokenNetwork.Contract.GetChannelParticipantInfo(&_TokenNetwork.CallOpts, participant, partner)
}

// GetChannelParticipantInfo is a free data retrieval call binding the contract method 0xac133709.
//
// Solidity: function getChannelParticipantInfo(participant address, partner address) constant returns(uint256, uint256, bool, bytes32, uint256)
func (_TokenNetwork *TokenNetworkCallerSession) GetChannelParticipantInfo(participant common.Address, partner common.Address) (*big.Int, *big.Int, bool, [32]byte, *big.Int, error) {
	return _TokenNetwork.Contract.GetChannelParticipantInfo(&_TokenNetwork.CallOpts, participant, partner)
}

// GetParticipantLockedAmount is a free data retrieval call binding the contract method 0xfd5f1e03.
//
// Solidity: function getParticipantLockedAmount(participant1 address, participant2 address, locksroot bytes32) constant returns(uint256)
func (_TokenNetwork *TokenNetworkCaller) GetParticipantLockedAmount(opts *bind.CallOpts, participant1 common.Address, participant2 common.Address, locksroot [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenNetwork.contract.Call(opts, out, "getParticipantLockedAmount", participant1, participant2, locksroot)
	return *ret0, err
}

// GetParticipantLockedAmount is a free data retrieval call binding the contract method 0xfd5f1e03.
//
// Solidity: function getParticipantLockedAmount(participant1 address, participant2 address, locksroot bytes32) constant returns(uint256)
func (_TokenNetwork *TokenNetworkSession) GetParticipantLockedAmount(participant1 common.Address, participant2 common.Address, locksroot [32]byte) (*big.Int, error) {
	return _TokenNetwork.Contract.GetParticipantLockedAmount(&_TokenNetwork.CallOpts, participant1, participant2, locksroot)
}

// GetParticipantLockedAmount is a free data retrieval call binding the contract method 0xfd5f1e03.
//
// Solidity: function getParticipantLockedAmount(participant1 address, participant2 address, locksroot bytes32) constant returns(uint256)
func (_TokenNetwork *TokenNetworkCallerSession) GetParticipantLockedAmount(participant1 common.Address, participant2 common.Address, locksroot [32]byte) (*big.Int, error) {
	return _TokenNetwork.Contract.GetParticipantLockedAmount(&_TokenNetwork.CallOpts, participant1, participant2, locksroot)
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

// TokenNetworkAddres is a free data retrieval call binding the contract method 0xfc0c546a.
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

// TokenNetworkAddres is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_TokenNetwork *TokenNetworkSession) Token() (common.Address, error) {
	return _TokenNetwork.Contract.Token(&_TokenNetwork.CallOpts)
}

// TokenNetworkAddres is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_TokenNetwork *TokenNetworkCallerSession) Token() (common.Address, error) {
	return _TokenNetwork.Contract.Token(&_TokenNetwork.CallOpts)
}

// CloseChannel is a paid mutator transaction binding the contract method 0x9abe275f.
//
// Solidity: function closeChannel(partner address, balance_hash bytes32, nonce uint256, additional_hash bytes32, signature bytes) returns()
func (_TokenNetwork *TokenNetworkTransactor) CloseChannel(opts *bind.TransactOpts, partner common.Address, balance_hash [32]byte, nonce *big.Int, additional_hash [32]byte, signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.contract.Transact(opts, "closeChannel", partner, balance_hash, nonce, additional_hash, signature)
}

// CloseChannel is a paid mutator transaction binding the contract method 0x9abe275f.
//
// Solidity: function closeChannel(partner address, balance_hash bytes32, nonce uint256, additional_hash bytes32, signature bytes) returns()
func (_TokenNetwork *TokenNetworkSession) CloseChannel(partner common.Address, balance_hash [32]byte, nonce *big.Int, additional_hash [32]byte, signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.CloseChannel(&_TokenNetwork.TransactOpts, partner, balance_hash, nonce, additional_hash, signature)
}

// CloseChannel is a paid mutator transaction binding the contract method 0x9abe275f.
//
// Solidity: function closeChannel(partner address, balance_hash bytes32, nonce uint256, additional_hash bytes32, signature bytes) returns()
func (_TokenNetwork *TokenNetworkTransactorSession) CloseChannel(partner common.Address, balance_hash [32]byte, nonce *big.Int, additional_hash [32]byte, signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.CloseChannel(&_TokenNetwork.TransactOpts, partner, balance_hash, nonce, additional_hash, signature)
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

// OpenChannel is a paid mutator transaction binding the contract method 0x0a798f24.
//
// Solidity: function openChannel(participant1 address, participant2 address, settle_timeout uint256) returns(bytes32)
func (_TokenNetwork *TokenNetworkTransactor) OpenChannel(opts *bind.TransactOpts, participant1 common.Address, participant2 common.Address, settle_timeout *big.Int) (*types.Transaction, error) {
	return _TokenNetwork.contract.Transact(opts, "openChannel", participant1, participant2, settle_timeout)
}

// OpenChannel is a paid mutator transaction binding the contract method 0x0a798f24.
//
// Solidity: function openChannel(participant1 address, participant2 address, settle_timeout uint256) returns(bytes32)
func (_TokenNetwork *TokenNetworkSession) OpenChannel(participant1 common.Address, participant2 common.Address, settle_timeout *big.Int) (*types.Transaction, error) {
	return _TokenNetwork.Contract.OpenChannel(&_TokenNetwork.TransactOpts, participant1, participant2, settle_timeout)
}

// OpenChannel is a paid mutator transaction binding the contract method 0x0a798f24.
//
// Solidity: function openChannel(participant1 address, participant2 address, settle_timeout uint256) returns(bytes32)
func (_TokenNetwork *TokenNetworkTransactorSession) OpenChannel(participant1 common.Address, participant2 common.Address, settle_timeout *big.Int) (*types.Transaction, error) {
	return _TokenNetwork.Contract.OpenChannel(&_TokenNetwork.TransactOpts, participant1, participant2, settle_timeout)
}

// SetTotalDeposit is a paid mutator transaction binding the contract method 0xa32a6737.
//
// Solidity: function setTotalDeposit(participant address, total_deposit uint256, partner address) returns()
func (_TokenNetwork *TokenNetworkTransactor) SetTotalDeposit(opts *bind.TransactOpts, participant common.Address, total_deposit *big.Int, partner common.Address) (*types.Transaction, error) {
	return _TokenNetwork.contract.Transact(opts, "setTotalDeposit", participant, total_deposit, partner)
}

// SetTotalDeposit is a paid mutator transaction binding the contract method 0xa32a6737.
//
// Solidity: function setTotalDeposit(participant address, total_deposit uint256, partner address) returns()
func (_TokenNetwork *TokenNetworkSession) SetTotalDeposit(participant common.Address, total_deposit *big.Int, partner common.Address) (*types.Transaction, error) {
	return _TokenNetwork.Contract.SetTotalDeposit(&_TokenNetwork.TransactOpts, participant, total_deposit, partner)
}

// SetTotalDeposit is a paid mutator transaction binding the contract method 0xa32a6737.
//
// Solidity: function setTotalDeposit(participant address, total_deposit uint256, partner address) returns()
func (_TokenNetwork *TokenNetworkTransactorSession) SetTotalDeposit(participant common.Address, total_deposit *big.Int, partner common.Address) (*types.Transaction, error) {
	return _TokenNetwork.Contract.SetTotalDeposit(&_TokenNetwork.TransactOpts, participant, total_deposit, partner)
}

// SetTotalWithdraw is a paid mutator transaction binding the contract method 0xc472c7e6.
//
// Solidity: function setTotalWithdraw(participant address, total_withdraw uint256, partner address, participant_signature bytes, partner_signature bytes) returns()
func (_TokenNetwork *TokenNetworkTransactor) SetTotalWithdraw(opts *bind.TransactOpts, participant common.Address, total_withdraw *big.Int, partner common.Address, participant_signature []byte, partner_signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.contract.Transact(opts, "setTotalWithdraw", participant, total_withdraw, partner, participant_signature, partner_signature)
}

// SetTotalWithdraw is a paid mutator transaction binding the contract method 0xc472c7e6.
//
// Solidity: function setTotalWithdraw(participant address, total_withdraw uint256, partner address, participant_signature bytes, partner_signature bytes) returns()
func (_TokenNetwork *TokenNetworkSession) SetTotalWithdraw(participant common.Address, total_withdraw *big.Int, partner common.Address, participant_signature []byte, partner_signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.SetTotalWithdraw(&_TokenNetwork.TransactOpts, participant, total_withdraw, partner, participant_signature, partner_signature)
}

// SetTotalWithdraw is a paid mutator transaction binding the contract method 0xc472c7e6.
//
// Solidity: function setTotalWithdraw(participant address, total_withdraw uint256, partner address, participant_signature bytes, partner_signature bytes) returns()
func (_TokenNetwork *TokenNetworkTransactorSession) SetTotalWithdraw(participant common.Address, total_withdraw *big.Int, partner common.Address, participant_signature []byte, partner_signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.SetTotalWithdraw(&_TokenNetwork.TransactOpts, participant, total_withdraw, partner, participant_signature, partner_signature)
}

// SettleChannel is a paid mutator transaction binding the contract method 0x79978630.
//
// Solidity: function settleChannel(participant1 address, participant1_transferred_amount uint256, participant1_locked_amount uint256, participant1_locksroot bytes32, participant2 address, participant2_transferred_amount uint256, participant2_locked_amount uint256, participant2_locksroot bytes32) returns()
func (_TokenNetwork *TokenNetworkTransactor) SettleChannel(opts *bind.TransactOpts, participant1 common.Address, participant1_transferred_amount *big.Int, participant1_locked_amount *big.Int, participant1_locksroot [32]byte, participant2 common.Address, participant2_transferred_amount *big.Int, participant2_locked_amount *big.Int, participant2_locksroot [32]byte) (*types.Transaction, error) {
	return _TokenNetwork.contract.Transact(opts, "settleChannel", participant1, participant1_transferred_amount, participant1_locked_amount, participant1_locksroot, participant2, participant2_transferred_amount, participant2_locked_amount, participant2_locksroot)
}

// SettleChannel is a paid mutator transaction binding the contract method 0x79978630.
//
// Solidity: function settleChannel(participant1 address, participant1_transferred_amount uint256, participant1_locked_amount uint256, participant1_locksroot bytes32, participant2 address, participant2_transferred_amount uint256, participant2_locked_amount uint256, participant2_locksroot bytes32) returns()
func (_TokenNetwork *TokenNetworkSession) SettleChannel(participant1 common.Address, participant1_transferred_amount *big.Int, participant1_locked_amount *big.Int, participant1_locksroot [32]byte, participant2 common.Address, participant2_transferred_amount *big.Int, participant2_locked_amount *big.Int, participant2_locksroot [32]byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.SettleChannel(&_TokenNetwork.TransactOpts, participant1, participant1_transferred_amount, participant1_locked_amount, participant1_locksroot, participant2, participant2_transferred_amount, participant2_locked_amount, participant2_locksroot)
}

// SettleChannel is a paid mutator transaction binding the contract method 0x79978630.
//
// Solidity: function settleChannel(participant1 address, participant1_transferred_amount uint256, participant1_locked_amount uint256, participant1_locksroot bytes32, participant2 address, participant2_transferred_amount uint256, participant2_locked_amount uint256, participant2_locksroot bytes32) returns()
func (_TokenNetwork *TokenNetworkTransactorSession) SettleChannel(participant1 common.Address, participant1_transferred_amount *big.Int, participant1_locked_amount *big.Int, participant1_locksroot [32]byte, participant2 common.Address, participant2_transferred_amount *big.Int, participant2_locked_amount *big.Int, participant2_locksroot [32]byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.SettleChannel(&_TokenNetwork.TransactOpts, participant1, participant1_transferred_amount, participant1_locked_amount, participant1_locksroot, participant2, participant2_transferred_amount, participant2_locked_amount, participant2_locksroot)
}

// Unlock is a paid mutator transaction binding the contract method 0x331d8e5d.
//
// Solidity: function unlock(participant address, partner address, merkle_tree_leaves bytes) returns()
func (_TokenNetwork *TokenNetworkTransactor) Unlock(opts *bind.TransactOpts, participant common.Address, partner common.Address, merkle_tree_leaves []byte) (*types.Transaction, error) {
	return _TokenNetwork.contract.Transact(opts, "unlock", participant, partner, merkle_tree_leaves)
}

// Unlock is a paid mutator transaction binding the contract method 0x331d8e5d.
//
// Solidity: function unlock(participant address, partner address, merkle_tree_leaves bytes) returns()
func (_TokenNetwork *TokenNetworkSession) Unlock(participant common.Address, partner common.Address, merkle_tree_leaves []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.Unlock(&_TokenNetwork.TransactOpts, participant, partner, merkle_tree_leaves)
}

// Unlock is a paid mutator transaction binding the contract method 0x331d8e5d.
//
// Solidity: function unlock(participant address, partner address, merkle_tree_leaves bytes) returns()
func (_TokenNetwork *TokenNetworkTransactorSession) Unlock(participant common.Address, partner common.Address, merkle_tree_leaves []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.Unlock(&_TokenNetwork.TransactOpts, participant, partner, merkle_tree_leaves)
}

// UpdateNonClosingBalanceProof is a paid mutator transaction binding the contract method 0xaec1dd81.
//
// Solidity: function updateNonClosingBalanceProof(closing_participant address, non_closing_participant address, balance_hash bytes32, nonce uint256, additional_hash bytes32, closing_signature bytes, non_closing_signature bytes) returns()
func (_TokenNetwork *TokenNetworkTransactor) UpdateNonClosingBalanceProof(opts *bind.TransactOpts, closing_participant common.Address, non_closing_participant common.Address, balance_hash [32]byte, nonce *big.Int, additional_hash [32]byte, closing_signature []byte, non_closing_signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.contract.Transact(opts, "updateNonClosingBalanceProof", closing_participant, non_closing_participant, balance_hash, nonce, additional_hash, closing_signature, non_closing_signature)
}

// UpdateNonClosingBalanceProof is a paid mutator transaction binding the contract method 0xaec1dd81.
//
// Solidity: function updateNonClosingBalanceProof(closing_participant address, non_closing_participant address, balance_hash bytes32, nonce uint256, additional_hash bytes32, closing_signature bytes, non_closing_signature bytes) returns()
func (_TokenNetwork *TokenNetworkSession) UpdateNonClosingBalanceProof(closing_participant common.Address, non_closing_participant common.Address, balance_hash [32]byte, nonce *big.Int, additional_hash [32]byte, closing_signature []byte, non_closing_signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.UpdateNonClosingBalanceProof(&_TokenNetwork.TransactOpts, closing_participant, non_closing_participant, balance_hash, nonce, additional_hash, closing_signature, non_closing_signature)
}

// UpdateNonClosingBalanceProof is a paid mutator transaction binding the contract method 0xaec1dd81.
//
// Solidity: function updateNonClosingBalanceProof(closing_participant address, non_closing_participant address, balance_hash bytes32, nonce uint256, additional_hash bytes32, closing_signature bytes, non_closing_signature bytes) returns()
func (_TokenNetwork *TokenNetworkTransactorSession) UpdateNonClosingBalanceProof(closing_participant common.Address, non_closing_participant common.Address, balance_hash [32]byte, nonce *big.Int, additional_hash [32]byte, closing_signature []byte, non_closing_signature []byte) (*types.Transaction, error) {
	return _TokenNetwork.Contract.UpdateNonClosingBalanceProof(&_TokenNetwork.TransactOpts, closing_participant, non_closing_participant, balance_hash, nonce, additional_hash, closing_signature, non_closing_signature)
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
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterChannelClosed is a free log retrieval operation binding the contract event 0x16a93f97197f719d19b0258648f368a06980009827e4a4a88892dd761ba4017c.
//
// Solidity: event ChannelClosed(channel_identifier indexed bytes32, closing_participant indexed address)
func (_TokenNetwork *TokenNetworkFilterer) FilterChannelClosed(opts *bind.FilterOpts, channel_identifier [][32]byte, closing_participant []common.Address) (*TokenNetworkChannelClosedIterator, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}
	var closing_participantRule []interface{}
	for _, closing_participantItem := range closing_participant {
		closing_participantRule = append(closing_participantRule, closing_participantItem)
	}

	logs, sub, err := _TokenNetwork.contract.FilterLogs(opts, "ChannelClosed", channel_identifierRule, closing_participantRule)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkChannelClosedIterator{contract: _TokenNetwork.contract, event: "ChannelClosed", logs: logs, sub: sub}, nil
}

// WatchChannelClosed is a free log subscription operation binding the contract event 0x16a93f97197f719d19b0258648f368a06980009827e4a4a88892dd761ba4017c.
//
// Solidity: event ChannelClosed(channel_identifier indexed bytes32, closing_participant indexed address)
func (_TokenNetwork *TokenNetworkFilterer) WatchChannelClosed(opts *bind.WatchOpts, sink chan<- *TokenNetworkChannelClosed, channel_identifier [][32]byte, closing_participant []common.Address) (event.Subscription, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}
	var closing_participantRule []interface{}
	for _, closing_participantItem := range closing_participant {
		closing_participantRule = append(closing_participantRule, closing_participantItem)
	}

	logs, sub, err := _TokenNetwork.contract.WatchLogs(opts, "ChannelClosed", channel_identifierRule, closing_participantRule)
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
// Solidity: event ChannelNewDeposit(channel_identifier indexed bytes32, participant indexed address, total_deposit uint256)
func (_TokenNetwork *TokenNetworkFilterer) FilterChannelNewDeposit(opts *bind.FilterOpts, channel_identifier [][32]byte, participant []common.Address) (*TokenNetworkChannelNewDepositIterator, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}
	var participantRule []interface{}
	for _, participantItem := range participant {
		participantRule = append(participantRule, participantItem)
	}

	logs, sub, err := _TokenNetwork.contract.FilterLogs(opts, "ChannelNewDeposit", channel_identifierRule, participantRule)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkChannelNewDepositIterator{contract: _TokenNetwork.contract, event: "ChannelNewDeposit", logs: logs, sub: sub}, nil
}

// WatchChannelNewDeposit is a free log subscription operation binding the contract event 0x0346e981e2bfa2366dc2307a8f1fa24779830a01121b1275fe565c6b98bb4d34.
//
// Solidity: event ChannelNewDeposit(channel_identifier indexed bytes32, participant indexed address, total_deposit uint256)
func (_TokenNetwork *TokenNetworkFilterer) WatchChannelNewDeposit(opts *bind.WatchOpts, sink chan<- *TokenNetworkChannelNewDeposit, channel_identifier [][32]byte, participant []common.Address) (event.Subscription, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}
	var participantRule []interface{}
	for _, participantItem := range participant {
		participantRule = append(participantRule, participantItem)
	}

	logs, sub, err := _TokenNetwork.contract.WatchLogs(opts, "ChannelNewDeposit", channel_identifierRule, participantRule)
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
// Solidity: event ChannelOpened(channel_identifier indexed bytes32, participant1 indexed address, participant2 indexed address, settle_timeout uint256)
func (_TokenNetwork *TokenNetworkFilterer) FilterChannelOpened(opts *bind.FilterOpts, channel_identifier [][32]byte, participant1 []common.Address, participant2 []common.Address) (*TokenNetworkChannelOpenedIterator, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}
	var participant1Rule []interface{}
	for _, participant1Item := range participant1 {
		participant1Rule = append(participant1Rule, participant1Item)
	}
	var participant2Rule []interface{}
	for _, participant2Item := range participant2 {
		participant2Rule = append(participant2Rule, participant2Item)
	}

	logs, sub, err := _TokenNetwork.contract.FilterLogs(opts, "ChannelOpened", channel_identifierRule, participant1Rule, participant2Rule)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkChannelOpenedIterator{contract: _TokenNetwork.contract, event: "ChannelOpened", logs: logs, sub: sub}, nil
}

// WatchChannelOpened is a free log subscription operation binding the contract event 0x448d27f1fe12f92a2070111296e68fd6ef0a01c0e05bf5819eda0dbcf267bf3d.
//
// Solidity: event ChannelOpened(channel_identifier indexed bytes32, participant1 indexed address, participant2 indexed address, settle_timeout uint256)
func (_TokenNetwork *TokenNetworkFilterer) WatchChannelOpened(opts *bind.WatchOpts, sink chan<- *TokenNetworkChannelOpened, channel_identifier [][32]byte, participant1 []common.Address, participant2 []common.Address) (event.Subscription, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}
	var participant1Rule []interface{}
	for _, participant1Item := range participant1 {
		participant1Rule = append(participant1Rule, participant1Item)
	}
	var participant2Rule []interface{}
	for _, participant2Item := range participant2 {
		participant2Rule = append(participant2Rule, participant2Item)
	}

	logs, sub, err := _TokenNetwork.contract.WatchLogs(opts, "ChannelOpened", channel_identifierRule, participant1Rule, participant2Rule)
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
	Participant        common.Address
	Unlocked_amount    *big.Int
	Returned_tokens    *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterChannelUnlocked is a free log retrieval operation binding the contract event 0x6f2d495eefa4b2d91a2287258f6f88722cabdaf48f4d1410d979e38b516a258f.
//
// Solidity: event ChannelUnlocked(channel_identifier indexed bytes32, participant indexed address, unlocked_amount uint256, returned_tokens uint256)
func (_TokenNetwork *TokenNetworkFilterer) FilterChannelUnlocked(opts *bind.FilterOpts, channel_identifier [][32]byte, participant []common.Address) (*TokenNetworkChannelUnlockedIterator, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}
	var participantRule []interface{}
	for _, participantItem := range participant {
		participantRule = append(participantRule, participantItem)
	}

	logs, sub, err := _TokenNetwork.contract.FilterLogs(opts, "ChannelUnlocked", channel_identifierRule, participantRule)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkChannelUnlockedIterator{contract: _TokenNetwork.contract, event: "ChannelUnlocked", logs: logs, sub: sub}, nil
}

// WatchChannelUnlocked is a free log subscription operation binding the contract event 0x6f2d495eefa4b2d91a2287258f6f88722cabdaf48f4d1410d979e38b516a258f.
//
// Solidity: event ChannelUnlocked(channel_identifier indexed bytes32, participant indexed address, unlocked_amount uint256, returned_tokens uint256)
func (_TokenNetwork *TokenNetworkFilterer) WatchChannelUnlocked(opts *bind.WatchOpts, sink chan<- *TokenNetworkChannelUnlocked, channel_identifier [][32]byte, participant []common.Address) (event.Subscription, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}
	var participantRule []interface{}
	for _, participantItem := range participant {
		participantRule = append(participantRule, participantItem)
	}

	logs, sub, err := _TokenNetwork.contract.WatchLogs(opts, "ChannelUnlocked", channel_identifierRule, participantRule)
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
	Channel_identifier [32]byte
	Participant        common.Address
	Total_withdraw     *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterChannelWithdraw is a free log retrieval operation binding the contract event 0x5860c94079516621b44c52f4423fd883c5f8f0370d5343f20f5dcb5b83773820.
//
// Solidity: event ChannelWithdraw(channel_identifier indexed bytes32, participant indexed address, total_withdraw uint256)
func (_TokenNetwork *TokenNetworkFilterer) FilterChannelWithdraw(opts *bind.FilterOpts, channel_identifier [][32]byte, participant []common.Address) (*TokenNetworkChannelWithdrawIterator, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}
	var participantRule []interface{}
	for _, participantItem := range participant {
		participantRule = append(participantRule, participantItem)
	}

	logs, sub, err := _TokenNetwork.contract.FilterLogs(opts, "ChannelWithdraw", channel_identifierRule, participantRule)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkChannelWithdrawIterator{contract: _TokenNetwork.contract, event: "ChannelWithdraw", logs: logs, sub: sub}, nil
}

// WatchChannelWithdraw is a free log subscription operation binding the contract event 0x5860c94079516621b44c52f4423fd883c5f8f0370d5343f20f5dcb5b83773820.
//
// Solidity: event ChannelWithdraw(channel_identifier indexed bytes32, participant indexed address, total_withdraw uint256)
func (_TokenNetwork *TokenNetworkFilterer) WatchChannelWithdraw(opts *bind.WatchOpts, sink chan<- *TokenNetworkChannelWithdraw, channel_identifier [][32]byte, participant []common.Address) (event.Subscription, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}
	var participantRule []interface{}
	for _, participantItem := range participant {
		participantRule = append(participantRule, participantItem)
	}

	logs, sub, err := _TokenNetwork.contract.WatchLogs(opts, "ChannelWithdraw", channel_identifierRule, participantRule)
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

// TokenNetworkNonClosingBalanceProofUpdatedIterator is returned from FilterNonClosingBalanceProofUpdated and is used to iterate over the raw logs and unpacked data for NonClosingBalanceProofUpdated events raised by the TokenNetwork contract.
type TokenNetworkNonClosingBalanceProofUpdatedIterator struct {
	Event *TokenNetworkNonClosingBalanceProofUpdated // Event containing the contract specifics and raw log

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
func (it *TokenNetworkNonClosingBalanceProofUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenNetworkNonClosingBalanceProofUpdated)
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
		it.Event = new(TokenNetworkNonClosingBalanceProofUpdated)
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
func (it *TokenNetworkNonClosingBalanceProofUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenNetworkNonClosingBalanceProofUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenNetworkNonClosingBalanceProofUpdated represents a NonClosingBalanceProofUpdated event raised by the TokenNetwork contract.
type TokenNetworkNonClosingBalanceProofUpdated struct {
	Channel_identifier  [32]byte
	Closing_participant common.Address
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterNonClosingBalanceProofUpdated is a free log retrieval operation binding the contract event 0x3558501c5224e6a6801ce5e7d81ed6f58489f1f7a88367cb14f2144b057f0b80.
//
// Solidity: event NonClosingBalanceProofUpdated(channel_identifier indexed bytes32, closing_participant indexed address)
func (_TokenNetwork *TokenNetworkFilterer) FilterNonClosingBalanceProofUpdated(opts *bind.FilterOpts, channel_identifier [][32]byte, closing_participant []common.Address) (*TokenNetworkNonClosingBalanceProofUpdatedIterator, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}
	var closing_participantRule []interface{}
	for _, closing_participantItem := range closing_participant {
		closing_participantRule = append(closing_participantRule, closing_participantItem)
	}

	logs, sub, err := _TokenNetwork.contract.FilterLogs(opts, "NonClosingBalanceProofUpdated", channel_identifierRule, closing_participantRule)
	if err != nil {
		return nil, err
	}
	return &TokenNetworkNonClosingBalanceProofUpdatedIterator{contract: _TokenNetwork.contract, event: "NonClosingBalanceProofUpdated", logs: logs, sub: sub}, nil
}

// WatchNonClosingBalanceProofUpdated is a free log subscription operation binding the contract event 0x3558501c5224e6a6801ce5e7d81ed6f58489f1f7a88367cb14f2144b057f0b80.
//
// Solidity: event NonClosingBalanceProofUpdated(channel_identifier indexed bytes32, closing_participant indexed address)
func (_TokenNetwork *TokenNetworkFilterer) WatchNonClosingBalanceProofUpdated(opts *bind.WatchOpts, sink chan<- *TokenNetworkNonClosingBalanceProofUpdated, channel_identifier [][32]byte, closing_participant []common.Address) (event.Subscription, error) {

	var channel_identifierRule []interface{}
	for _, channel_identifierItem := range channel_identifier {
		channel_identifierRule = append(channel_identifierRule, channel_identifierItem)
	}
	var closing_participantRule []interface{}
	for _, closing_participantItem := range closing_participant {
		closing_participantRule = append(closing_participantRule, closing_participantItem)
	}

	logs, sub, err := _TokenNetwork.contract.WatchLogs(opts, "NonClosingBalanceProofUpdated", channel_identifierRule, closing_participantRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenNetworkNonClosingBalanceProofUpdated)
				if err := _TokenNetwork.contract.UnpackLog(event, "NonClosingBalanceProofUpdated", log); err != nil {
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
const UtilsBin = `0x608060405234801561001057600080fd5b50610187806100206000396000f30060806040526004361061004b5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416637709bc788114610050578063b32c65c814610092575b600080fd5b34801561005c57600080fd5b5061007e73ffffffffffffffffffffffffffffffffffffffff6004351661011c565b604080519115158252519081900360200190f35b34801561009e57600080fd5b506100a7610124565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100e15781810151838201526020016100c9565b50505050905090810190601f16801561010e5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000903b1190565b60408051808201909152600581527f302e332e5f0000000000000000000000000000000000000000000000000000006020820152815600a165627a7a72305820f85c7f4be5b7cd2b00bdfd1b7e341095deed3432ed34b676ee2fd22a83cb0ef50029`

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
