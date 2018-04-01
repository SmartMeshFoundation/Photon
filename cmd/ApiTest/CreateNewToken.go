package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/SmartMeshFoundation/raiden-network"
	"github.com/SmartMeshFoundation/raiden-network/abi/bind"
	"github.com/SmartMeshFoundation/raiden-network/network/rpc"
	"github.com/SmartMeshFoundation/raiden-network/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/slonzok/getpass"
	"log"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"
)

var globalPassword string = "123"

func CreateNewToken(EthRpcEndpoint string, KeystorePath string) (TokenName string, registryAddress common.Address, err error) {
	// Create an IPC based RPC connection to a remote node and an authorized transactor
	conn, err := ethclient.Dial(EthRpcEndpoint)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to connect to the Ethereum client: %v", err))
	}
	_, key := promptAccount(KeystorePath)
	//log.Println("start to deploy ...")
	registryAddress = DeployContract(key, conn)
	registry, _ := rpc.NewRegistry(registryAddress, conn)
	TokenName = createTokenAndChannels(key, conn, registry, KeystorePath, false)
	return
}
func promptAccount(keystorePath string) (addr common.Address, key *ecdsa.PrivateKey) {
	am := raiden_network.NewAccountManager(keystorePath)
	if len(am.Accounts) == 0 {
		log.Fatal(fmt.Sprintf("No Ethereum accounts found in the directory %s", keystorePath))
		os.Exit(1)
	}
	addr = am.Accounts[0].Address
	for i := 0; i < 3; i++ {
		//retries three times
		if len(globalPassword) <= 0 {
			globalPassword = getpass.Prompt("Enter the password to unlock")
		}
		//fmt.Printf("\npassword is %s\n", password)
		keybin, err := am.GetPrivateKey(addr, globalPassword)
		if err != nil && i == 3 {
			log.Fatal(fmt.Sprintf("Exhausted passphrase unlock attempts for %s. Aborting ...", addr))
			os.Exit(1)
		}
		if err != nil {
			log.Println(fmt.Sprintf("password incorrect\n Please try again or kill the process to quit.\nUsually Ctrl-c."))
			continue
		}
		key, _ = crypto.ToECDSA(keybin)
		break
	}
	return
}
func DeployContract(key *ecdsa.PrivateKey, conn *ethclient.Client) (RegistryAddress common.Address) {
	auth := bind.NewKeyedTransactor(key)
	//DeployNettingChannelLibrary
	NettingChannelLibraryAddress, tx, _, err := rpc.DeployNettingChannelLibrary(auth, conn)
	if err != nil {
		log.Fatalf("Failed to DeployNettingChannelLibrary: %v", err)
	}
	ctx := context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	//fmt.Printf("DeployNettingChannelLibrary complete...\n")
	//DeployChannelManagerLibrary link nettingchannle library before deploy
	rpc.ChannelManagerLibraryBin = strings.Replace(rpc.ChannelManagerLibraryBin, "__NettingChannelLibrary.sol:NettingCha__", NettingChannelLibraryAddress.String()[2:], -1)
	ChannelManagerLibraryAddress, tx, _, err := rpc.DeployChannelManagerLibrary(auth, conn)
	if err != nil {
		log.Fatalf("Failed to deploy new token contract: %v", err)
	}
	ctx = context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	//fmt.Printf("DeployChannelManagerLibrary complete...\n")
	//DeployRegistry link channelmanagerlibrary before deploy
	rpc.RegistryBin = strings.Replace(rpc.RegistryBin, "__ChannelManagerLibrary.sol:ChannelMan__", ChannelManagerLibraryAddress.String()[2:], -1)
	RegistryAddress, tx, _, err = rpc.DeployRegistry(auth, conn)
	if err != nil {
		log.Fatalf("Failed to deploy new token contract: %v", err)
	}
	ctx = context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	//fmt.Printf("DeployRegistry complete...\n")

	//fmt.Printf("RegistryAddress=%s\n", RegistryAddress.String())
	return
}
func createTokenAndChannels(key *ecdsa.PrivateKey, conn *ethclient.Client, registry *rpc.Registry, keystorepath string, createchannel bool) (TokenName string) {
	managerAddress, tokenAddress := NewToken(key, conn, registry)
	TokenName = tokenAddress.Hex()
	manager, _ := rpc.NewChannelManagerContract(managerAddress, conn)
	token, _ := rpc.NewToken(tokenAddress, conn)
	am := raiden_network.NewAccountManager(keystorepath)
	var accounts []common.Address
	var keys []*ecdsa.PrivateKey
	for _, account := range am.Accounts {
		accounts = append(accounts, account.Address)
		keybin, err := am.GetPrivateKey(account.Address, globalPassword)
		if err != nil {
			log.Fatalf("password error for %s", account.Address.String())
		}
		keytemp, _ := crypto.ToECDSA(keybin)
		keys = append(keys, keytemp)
	}
	fmt.Sprintf("key=%s", key)
	TransferMoneyForAccounts(key, conn, accounts[1:], token)
	if createchannel {
		CreateChannels(conn, accounts, keys, manager, token)
	}
	return
}
func NewToken(key *ecdsa.PrivateKey, conn *ethclient.Client, registry *rpc.Registry) (mgrAddress common.Address, tokenAddr common.Address) {
	auth := bind.NewKeyedTransactor(key)
	tokenAddr, tx, _, err := rpc.DeployHumanStandardToken(auth, conn, big.NewInt(50000000000), "test", 2, "test symoble")
	if err != nil {
		log.Fatalf("Failed to DeployHumanStandardToken: %v", err)
	}
	ctx := context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	//fmt.Printf("DeployHumanStandardToken complete...\n")
	tx, err = registry.AddToken(auth, tokenAddr)
	if err != nil {
		log.Fatalf("Failed to AddToken: %v", err)
	}
	ctx = context.Background()
	_, err = bind.WaitMined(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to AddToken when mining :%v", err)
	}
	mgrAddress, err = registry.ChannelManagerByToken(nil, tokenAddr)
	//fmt.Printf("DeployHumanStandardToken complete... \nNewTokenAddr=%s,mgr=%s\n", tokenAddr.String(), mgrAddress.String())
	return
}
func TransferMoneyForAccounts(key *ecdsa.PrivateKey, conn *ethclient.Client, accounts []common.Address, token *rpc.Token) {
	wg := sync.WaitGroup{}
	wg.Add(len(accounts))
	auth := bind.NewKeyedTransactor(key)
	nonce, _ := conn.PendingNonceAt(context.Background(), auth.From)
	for index, account := range accounts {
		go func(account common.Address, i int) {
			//fmt.Printf("i=%d\n", i)
			auth2 := bind.NewKeyedTransactor(key)
			auth2.Nonce = big.NewInt(int64(nonce) + int64(i))
			//fmt.Printf("transfer to %s,nonce=%s\n", account.String(), auth2.Nonce)
			tx, err := token.Transfer(auth2, account, big.NewInt(500000))
			if err != nil {
				log.Fatalf("Failed to Transfer: %v", err)
			}
			ctx := context.Background()
			_, err = bind.WaitMined(ctx, conn, tx)
			if err != nil {
				log.Fatalf("failed to Transfer when mining :%v", err)
			}
			//fmt.Printf("Transfer complete...\n")
			wg.Done()
		}(account, index)
		time.Sleep(time.Millisecond * 100)
	}
	wg.Wait()
}
func CreateChannels(conn *ethclient.Client, accounts []common.Address, keys []*ecdsa.PrivateKey, manager *rpc.ChannelManagerContract, token *rpc.Token) {
	if len(accounts) < 6 {
		panic("need 6 accounts")
	}
	AccountA := accounts[0]
	AccountB := accounts[1]
	AccountC := accounts[2]
	AccountD := accounts[3]
	AccountE := accounts[4]
	AccountF := accounts[5]
	fmt.Printf("accountA=%saccountB=%saccountC=%saccountD=%saccountE=%saccountF=%s\n", AccountA.String(), AccountB.String(), AccountC.String(), AccountD.String(), AccountE.String(), AccountF.String())
	keyA := keys[0]
	keyB := keys[1]
	keyC := keys[2]
	keyD := keys[3]
	keyE := keys[4]
	keyF := keys[5]
	/*
	   4.1 create channel A-B and save 100 both

	*/
	creatAChannelAndDeposit(AccountA, AccountB, keyA, keyB, 100, manager, token, conn)
	/*
	 4.2 create channel B-C and save 50 both
	*/
	creatAChannelAndDeposit(AccountB, AccountC, keyB, keyC, 50, manager, token, conn)
	/*
	  4.3 create channel C-E and save 100 both
	*/

	creatAChannelAndDeposit(AccountC, AccountE, keyC, keyE, 100, manager, token, conn)

	/*
	   4.4 create channel A-D and save 100 both
	*/

	creatAChannelAndDeposit(AccountA, AccountD, keyA, keyD, 100, manager, token, conn)

	/*
	  4.5 create channel B-D and save 100 both
	*/

	creatAChannelAndDeposit(AccountB, AccountD, keyB, keyD, 100, manager, token, conn)

	/*
	   4.6 create channel D-F and save 100 both
	*/

	creatAChannelAndDeposit(AccountD, AccountF, keyD, keyF, 100, manager, token, conn)

	/*
	   4.7 create channel F-E and save 100 both
	*/

	creatAChannelAndDeposit(AccountF, AccountE, keyF, keyE, 100, manager, token, conn)
	/*
	   4.8 create channel C-F and save 50 both
	*/

	creatAChannelAndDeposit(AccountC, AccountF, keyC, keyF, 50, manager, token, conn)

	/*
	   4.9     D-E 100
	*/

	creatAChannelAndDeposit(AccountD, AccountE, keyD, keyE, 100, manager, token, conn)

}
func creatAChannelAndDeposit(account1, account2 common.Address, key1, key2 *ecdsa.PrivateKey, amount int64, manager *rpc.ChannelManagerContract, token *rpc.Token, conn *ethclient.Client) {
	auth1 := bind.NewKeyedTransactor(key1)
	auth1.GasLimit = uint64(params.GAS_LIMIT)
	auth1.GasPrice = big.NewInt(params.GAS_PRICE)
	callAuth1 := &bind.CallOpts{
		Pending: false,
		From:    account1,
		Context: context.Background(),
	}
	auth2 := bind.NewKeyedTransactor(key2)
	auth2.GasLimit = uint64(params.GAS_LIMIT)
	auth2.GasPrice = big.NewInt(params.GAS_PRICE)
	tx, err := manager.NewChannel(auth1, account2, big.NewInt(600))
	if err != nil {
		log.Printf("Failed to NewChannel: %v,%s,%s", err, auth1.From.String(), account2.String())
		return
	}
	ctx := context.Background()
	_, err = bind.WaitMined(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to NewChannel when mining :%v", err)
	}
	//fmt.Printf("NewChannel complete...\n")
	//step 2 deopsit
	//step 2.1 aprove
	channelAddress, err := manager.GetChannelWith(callAuth1, account2)
	if err != nil {
		log.Fatalf("failed to get channel")
	}
	channel, _ := rpc.NewNettingChannelContract(channelAddress, conn)
	wg2 := sync.WaitGroup{}
	go func() {
		wg2.Add(1)
		tx, err := token.Approve(auth1, channelAddress, big.NewInt(amount))
		if err != nil {
			log.Fatalf("Failed to Approve: %v", err)
		}
		ctx = context.Background()
		_, err = bind.WaitMined(ctx, conn, tx)
		if err != nil {
			log.Fatalf("failed to Approve when mining :%v", err)
		}
		//fmt.Printf("Approve complete...\n")
		tx, err = channel.Deposit(auth1, big.NewInt(amount))
		if err != nil {
			log.Fatalf("Failed to Deposit: %v", err)
		}
		ctx = context.Background()
		_, err = bind.WaitMined(ctx, conn, tx)
		if err != nil {
			log.Fatalf("failed to Deposit when mining :%v", err)
		}
		//fmt.Printf("Deposit complete...\n")
		wg2.Done()
	}()
	go func() {
		wg2.Add(1)
		tx, err := token.Approve(auth2, channelAddress, big.NewInt(amount))
		if err != nil {
			log.Fatalf("Failed to Approve: %v", err)
		}
		ctx = context.Background()
		_, err = bind.WaitMined(ctx, conn, tx)
		if err != nil {
			log.Fatalf("failed to Approve when mining :%v", err)
		}
		//fmt.Printf("Approve complete...\n")
		tx, err = channel.Deposit(auth2, big.NewInt(amount))
		if err != nil {
			log.Fatalf("Failed to Deposit: %v", err)
		}
		ctx = context.Background()
		_, err = bind.WaitMined(ctx, conn, tx)
		if err != nil {
			log.Fatalf("failed to Deposit when mining :%v", err)
		}
		//fmt.Printf("Deposit complete...\n")
		wg2.Done()
	}()
	time.Sleep(time.Second)
	wg2.Wait()
}
