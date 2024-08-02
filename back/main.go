package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"log"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
)

const key = ``

// создание кошелька
func createWallet() (string, string, *ecdsa.PrivateKey) {
	getPrivateKey, err := crypto.GenerateKey()

	if err != nil {
		log.Println(err)
	}

	getPublicKey := crypto.FromECDSA(getPrivateKey)
	thePublicKey := hexutil.Encode(getPublicKey)

	thePublicAddress := crypto.PubkeyToAddress(getPrivateKey.PublicKey).Hex()
	return thePublicAddress, thePublicKey, getPrivateKey
}

func sendTx(client *ethclient.Client, ctx context.Context, to string, priv string) {
	RecipientAddress := common.HexToAddress(to)

	privateKey, err := crypto.HexToECDSA(priv)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Public Key Error")
	}

	SenderAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(ctx, SenderAddress)
	if err != nil {

		log.Println(err)
	}

	amount := big.NewInt(100)
	gasLimit := 3600
	gas, err := client.SuggestGasPrice(ctx)

	if err != nil {
		log.Println(err)
	}

	ChainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Println(err)
	}

	transaction := types.NewTransaction(nonce, RecipientAddress, amount, uint64(gasLimit), gas, nil)
	signedTx, err := types.SignTx(transaction, types.NewEIP155Signer(ChainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("transaction sent: %s", signedTx.Hash().Hex())
}

func main() {

	ctx := context.Background()

	pubAddress, pubKey, secretKey := createWallet()
	fmt.Println(pubAddress, pubKey, secretKey)

	conn, err := ethclient.DialContext(ctx, "https://rpc-mumbai.polygon.technology")

	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	auth, err := bind.NewTransactor(strings.NewReader(pubKey), "")
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}

	// fmt.Printf("Contract pending deploy: 0x%x\n", address)

	// opts *TransactOpts, abi abi.ABI, bytecode []byte, backend ContractBackend, params ...interface{}
	// address, tx, instance, err := bind.DeployContract(auth, LockerMetaData.ABI, "", conn)
	// if err != nil {
	// 	log.Fatalf("Failed to deploy new storage contract: %v", err)
	// }
	// fmt.Printf("Contract pending deploy: 0x%x\n", address)
	// fmt.Printf("Transaction waiting to be mined: 0x%x\n\n", tx.Hash())
	// time.Sleep(250 * time.Millisecond) // Allow it to be processed by the local node :P

	// Create an IPC based RPC connection to a remote node and an authorized transactor

	// geth: http, ws, or ipc

	// Instantiate the contract and display its name
	// NOTE update the deployment address!
	locker, err := NewLocker(common.HexToAddress("0xefb82e04318493403BCc09F97cC4a7ec8A65e8d4"), conn)
	if err != nil {
		log.Fatalf("Failed to instantiate Storage contract: %v", err)
	}

	// Call the store() function
	tx, err := locker.Withdraw(auth)
	if err != nil {
		log.Fatalf("Failed to update value: %v", err)
	}
	fmt.Println("Tx: ", tx)

}

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// LockerMetaData contains all meta data concerning the Locker contract.
var LockerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"excepted\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"NotEnoughValue\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"customer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"users\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"OrderCreated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newFeePrice\",\"type\":\"uint256\"}],\"name\":\"changeFeePrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"}],\"name\":\"claim\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"result\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"users\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"createOrderWithAmount\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"users\",\"type\":\"address[]\"}],\"name\":\"createOrderWithoutAmount\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feePrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// LockerABI is the input ABI used to generate the binding from.
// Deprecated: Use LockerMetaData.ABI instead.
var LockerABI = LockerMetaData.ABI

// Locker is an auto generated Go binding around an Ethereum contract.
type Locker struct {
	LockerCaller     // Read-only binding to the contract
	LockerTransactor // Write-only binding to the contract
	LockerFilterer   // Log filterer for contract events
}

// LockerCaller is an auto generated read-only Go binding around an Ethereum contract.
type LockerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LockerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LockerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LockerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LockerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LockerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LockerSession struct {
	Contract     *Locker           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LockerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LockerCallerSession struct {
	Contract *LockerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// LockerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LockerTransactorSession struct {
	Contract     *LockerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LockerRaw is an auto generated low-level Go binding around an Ethereum contract.
type LockerRaw struct {
	Contract *Locker // Generic contract binding to access the raw methods on
}

// LockerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LockerCallerRaw struct {
	Contract *LockerCaller // Generic read-only contract binding to access the raw methods on
}

// LockerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LockerTransactorRaw struct {
	Contract *LockerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLocker creates a new instance of Locker, bound to a specific deployed contract.
func NewLocker(address common.Address, backend bind.ContractBackend) (*Locker, error) {
	contract, err := bindLocker(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Locker{LockerCaller: LockerCaller{contract: contract}, LockerTransactor: LockerTransactor{contract: contract}, LockerFilterer: LockerFilterer{contract: contract}}, nil
}

// NewLockerCaller creates a new read-only instance of Locker, bound to a specific deployed contract.
func NewLockerCaller(address common.Address, caller bind.ContractCaller) (*LockerCaller, error) {
	contract, err := bindLocker(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LockerCaller{contract: contract}, nil
}

// NewLockerTransactor creates a new write-only instance of Locker, bound to a specific deployed contract.
func NewLockerTransactor(address common.Address, transactor bind.ContractTransactor) (*LockerTransactor, error) {
	contract, err := bindLocker(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LockerTransactor{contract: contract}, nil
}

// NewLockerFilterer creates a new log filterer instance of Locker, bound to a specific deployed contract.
func NewLockerFilterer(address common.Address, filterer bind.ContractFilterer) (*LockerFilterer, error) {
	contract, err := bindLocker(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LockerFilterer{contract: contract}, nil
}

// bindLocker binds a generic wrapper to an already deployed contract.
func bindLocker(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LockerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Locker *LockerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Locker.Contract.LockerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Locker *LockerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Locker.Contract.LockerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Locker *LockerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Locker.Contract.LockerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Locker *LockerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Locker.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Locker *LockerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Locker.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Locker *LockerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Locker.Contract.contract.Transact(opts, method, params...)
}

// FeeAmount is a free data retrieval call binding the contract method 0x69e15404.
//
// Solidity: function feeAmount() view returns(uint256)
func (_Locker *LockerCaller) FeeAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Locker.contract.Call(opts, &out, "feeAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FeeAmount is a free data retrieval call binding the contract method 0x69e15404.
//
// Solidity: function feeAmount() view returns(uint256)
func (_Locker *LockerSession) FeeAmount() (*big.Int, error) {
	return _Locker.Contract.FeeAmount(&_Locker.CallOpts)
}

// FeeAmount is a free data retrieval call binding the contract method 0x69e15404.
//
// Solidity: function feeAmount() view returns(uint256)
func (_Locker *LockerCallerSession) FeeAmount() (*big.Int, error) {
	return _Locker.Contract.FeeAmount(&_Locker.CallOpts)
}

// FeePrice is a free data retrieval call binding the contract method 0x54ad9718.
//
// Solidity: function feePrice() view returns(uint256)
func (_Locker *LockerCaller) FeePrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Locker.contract.Call(opts, &out, "feePrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FeePrice is a free data retrieval call binding the contract method 0x54ad9718.
//
// Solidity: function feePrice() view returns(uint256)
func (_Locker *LockerSession) FeePrice() (*big.Int, error) {
	return _Locker.Contract.FeePrice(&_Locker.CallOpts)
}

// FeePrice is a free data retrieval call binding the contract method 0x54ad9718.
//
// Solidity: function feePrice() view returns(uint256)
func (_Locker *LockerCallerSession) FeePrice() (*big.Int, error) {
	return _Locker.Contract.FeePrice(&_Locker.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Locker *LockerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Locker.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Locker *LockerSession) Owner() (common.Address, error) {
	return _Locker.Contract.Owner(&_Locker.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Locker *LockerCallerSession) Owner() (common.Address, error) {
	return _Locker.Contract.Owner(&_Locker.CallOpts)
}

// ChangeFeePrice is a paid mutator transaction binding the contract method 0xa957b3b4.
//
// Solidity: function changeFeePrice(uint256 newFeePrice) returns()
func (_Locker *LockerTransactor) ChangeFeePrice(opts *bind.TransactOpts, newFeePrice *big.Int) (*types.Transaction, error) {
	return _Locker.contract.Transact(opts, "changeFeePrice", newFeePrice)
}

// ChangeFeePrice is a paid mutator transaction binding the contract method 0xa957b3b4.
//
// Solidity: function changeFeePrice(uint256 newFeePrice) returns()
func (_Locker *LockerSession) ChangeFeePrice(newFeePrice *big.Int) (*types.Transaction, error) {
	return _Locker.Contract.ChangeFeePrice(&_Locker.TransactOpts, newFeePrice)
}

// ChangeFeePrice is a paid mutator transaction binding the contract method 0xa957b3b4.
//
// Solidity: function changeFeePrice(uint256 newFeePrice) returns()
func (_Locker *LockerTransactorSession) ChangeFeePrice(newFeePrice *big.Int) (*types.Transaction, error) {
	return _Locker.Contract.ChangeFeePrice(&_Locker.TransactOpts, newFeePrice)
}

// Claim is a paid mutator transaction binding the contract method 0xbd66528a.
//
// Solidity: function claim(bytes32 orderId) returns(bool result)
func (_Locker *LockerTransactor) Claim(opts *bind.TransactOpts, orderId [32]byte) (*types.Transaction, error) {
	return _Locker.contract.Transact(opts, "claim", orderId)
}

// Claim is a paid mutator transaction binding the contract method 0xbd66528a.
//
// Solidity: function claim(bytes32 orderId) returns(bool result)
func (_Locker *LockerSession) Claim(orderId [32]byte) (*types.Transaction, error) {
	return _Locker.Contract.Claim(&_Locker.TransactOpts, orderId)
}

// Claim is a paid mutator transaction binding the contract method 0xbd66528a.
//
// Solidity: function claim(bytes32 orderId) returns(bool result)
func (_Locker *LockerTransactorSession) Claim(orderId [32]byte) (*types.Transaction, error) {
	return _Locker.Contract.Claim(&_Locker.TransactOpts, orderId)
}

// CreateOrderWithAmount is a paid mutator transaction binding the contract method 0x4cd426e2.
//
// Solidity: function createOrderWithAmount(address[] users, uint256 amount) payable returns(bytes32 orderId)
func (_Locker *LockerTransactor) CreateOrderWithAmount(opts *bind.TransactOpts, users []common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Locker.contract.Transact(opts, "createOrderWithAmount", users, amount)
}

// CreateOrderWithAmount is a paid mutator transaction binding the contract method 0x4cd426e2.
//
// Solidity: function createOrderWithAmount(address[] users, uint256 amount) payable returns(bytes32 orderId)
func (_Locker *LockerSession) CreateOrderWithAmount(users []common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Locker.Contract.CreateOrderWithAmount(&_Locker.TransactOpts, users, amount)
}

// CreateOrderWithAmount is a paid mutator transaction binding the contract method 0x4cd426e2.
//
// Solidity: function createOrderWithAmount(address[] users, uint256 amount) payable returns(bytes32 orderId)
func (_Locker *LockerTransactorSession) CreateOrderWithAmount(users []common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Locker.Contract.CreateOrderWithAmount(&_Locker.TransactOpts, users, amount)
}

// CreateOrderWithoutAmount is a paid mutator transaction binding the contract method 0xbd206431.
//
// Solidity: function createOrderWithoutAmount(address[] users) payable returns(bytes32 orderId)
func (_Locker *LockerTransactor) CreateOrderWithoutAmount(opts *bind.TransactOpts, users []common.Address) (*types.Transaction, error) {
	return _Locker.contract.Transact(opts, "createOrderWithoutAmount", users)
}

// CreateOrderWithoutAmount is a paid mutator transaction binding the contract method 0xbd206431.
//
// Solidity: function createOrderWithoutAmount(address[] users) payable returns(bytes32 orderId)
func (_Locker *LockerSession) CreateOrderWithoutAmount(users []common.Address) (*types.Transaction, error) {
	return _Locker.Contract.CreateOrderWithoutAmount(&_Locker.TransactOpts, users)
}

// CreateOrderWithoutAmount is a paid mutator transaction binding the contract method 0xbd206431.
//
// Solidity: function createOrderWithoutAmount(address[] users) payable returns(bytes32 orderId)
func (_Locker *LockerTransactorSession) CreateOrderWithoutAmount(users []common.Address) (*types.Transaction, error) {
	return _Locker.Contract.CreateOrderWithoutAmount(&_Locker.TransactOpts, users)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_Locker *LockerTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Locker.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_Locker *LockerSession) Withdraw() (*types.Transaction, error) {
	return _Locker.Contract.Withdraw(&_Locker.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_Locker *LockerTransactorSession) Withdraw() (*types.Transaction, error) {
	return _Locker.Contract.Withdraw(&_Locker.TransactOpts)
}

// LockerOrderCreatedIterator is returned from FilterOrderCreated and is used to iterate over the raw logs and unpacked data for OrderCreated events raised by the Locker contract.
type LockerOrderCreatedIterator struct {
	Event *LockerOrderCreated // Event containing the contract specifics and raw log

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
func (it *LockerOrderCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockerOrderCreated)
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
		it.Event = new(LockerOrderCreated)
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
func (it *LockerOrderCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockerOrderCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockerOrderCreated represents a OrderCreated event raised by the Locker contract.
type LockerOrderCreated struct {
	OrderId  [32]byte
	Customer common.Address
	Users    []common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOrderCreated is a free log retrieval operation binding the contract event 0xd9294074e3373969bceba674da9b2253b3e133c0e18740113e1e35dcc2cf8fb7.
//
// Solidity: event OrderCreated(bytes32 indexed orderId, address indexed customer, address[] users, uint256 amount)
func (_Locker *LockerFilterer) FilterOrderCreated(opts *bind.FilterOpts, orderId [][32]byte, customer []common.Address) (*LockerOrderCreatedIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var customerRule []interface{}
	for _, customerItem := range customer {
		customerRule = append(customerRule, customerItem)
	}

	logs, sub, err := _Locker.contract.FilterLogs(opts, "OrderCreated", orderIdRule, customerRule)
	if err != nil {
		return nil, err
	}
	return &LockerOrderCreatedIterator{contract: _Locker.contract, event: "OrderCreated", logs: logs, sub: sub}, nil
}

// WatchOrderCreated is a free log subscription operation binding the contract event 0xd9294074e3373969bceba674da9b2253b3e133c0e18740113e1e35dcc2cf8fb7.
//
// Solidity: event OrderCreated(bytes32 indexed orderId, address indexed customer, address[] users, uint256 amount)
func (_Locker *LockerFilterer) WatchOrderCreated(opts *bind.WatchOpts, sink chan<- *LockerOrderCreated, orderId [][32]byte, customer []common.Address) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var customerRule []interface{}
	for _, customerItem := range customer {
		customerRule = append(customerRule, customerItem)
	}

	logs, sub, err := _Locker.contract.WatchLogs(opts, "OrderCreated", orderIdRule, customerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockerOrderCreated)
				if err := _Locker.contract.UnpackLog(event, "OrderCreated", log); err != nil {
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

// ParseOrderCreated is a log parse operation binding the contract event 0xd9294074e3373969bceba674da9b2253b3e133c0e18740113e1e35dcc2cf8fb7.
//
// Solidity: event OrderCreated(bytes32 indexed orderId, address indexed customer, address[] users, uint256 amount)
func (_Locker *LockerFilterer) ParseOrderCreated(log types.Log) (*LockerOrderCreated, error) {
	event := new(LockerOrderCreated)
	if err := _Locker.contract.UnpackLog(event, "OrderCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
