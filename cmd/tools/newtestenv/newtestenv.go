package main

import (
	"log"
	"math"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/SmartMeshFoundation/Photon/network/helper"

	"context"

	"fmt"

	"os"

	"crypto/ecdsa"

	"math/big"

	"sync"

	"time"

	"github.com/SmartMeshFoundation/Photon/accounts"
	"github.com/SmartMeshFoundation/Photon/cmd/tools/newtestenv/createchannel"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts/test/tokens/tokenerc223"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts/test/tokens/tokenerc223approve"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts/test/tokens/tokenether"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts/test/tokens/tokenstandard"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/urfave/cli.v1"
)

var passwords = []string{"123", "111111", "123456"}

const (
	tokenERC223   = "erc223"
	tokenStandard = "standard"
	//#nosec
	tokenERC223Approve = "erc223_approve"
	tokenEther         = "ether"
)

var base = 0

func getAmount(x *big.Int) *big.Int {
	y := new(big.Int)
	y = y.Mul(x, big.NewInt(int64(math.Pow10(base))))
	return y
}
func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "address",
			Usage: "The ethereum address you would like photon to use and for which a keystore file exists in your local system.",
		},
		cli.StringFlag{
			Name:  "keystore-path",
			Usage: "If you have a non-standard path for the ethereum keystore directory provide it using this argument. ",
			//Value: ethutils.DirectoryString{params.DefaultKeyStoreDir()},
			Value: "../../../testdata/keystore",
		},
		cli.StringFlag{
			Name: "eth-rpc-endpoint",
			Usage: `"host:port" address of ethereum JSON-RPC server.\n'
	           'Also accepts a protocol prefix (ws:// or ipc channel) with optional port',`,
			Value: fmt.Sprintf("http://127.0.0.1:8545"), //, node.DefaultWSEndpoint()),
		},
		cli.BoolFlag{
			Name:  "not-create-token",
			Usage: "not-create token.",
		},
		cli.BoolFlag{
			Name:  "not-create-channel",
			Usage: "not-create channels between node for test.",
		},
		cli.IntFlag{
			Name:  "base",
			Usage: "decimal part of ERC20 Token",
			Value: 0,
		},
		cli.StringFlag{
			Name:  "password",
			Usage: "plain text password for all accounts",
			Value: "123",
		},
		cli.IntFlag{
			Name:  "tokennum",
			Usage: "how many tokens to deploy ,there are four types token to candidate. so max number is 4",
			Value: 4,
		},
	}
	app.Action = mainctx
	app.Name = "newphotonenv"
	app.Version = "0.1"
	err := app.Run(os.Args)
	if err != nil {
		log.Printf("run err %s\n", err)
	}
}

func mainctx(ctx *cli.Context) error {
	fmt.Printf("eth-rpc-endpoint:%s\n", ctx.String("eth-rpc-endpoint"))
	fmt.Printf("not-create-channel=%v\n", ctx.Bool("not-create-channel"))
	fmt.Printf("not-create-token=%v\n", ctx.Bool("not-create-token"))
	base = ctx.Int("base")
	passwords[0] = ctx.String("password")
	tokenNumber := ctx.Int("tokennum")
	//if tokenNumber <= 0 || tokenNumber > 4 {
	//	log.Fatalf("tokenum must be between 1-4")
	//}
	// Create an IPC based RPC connection to a remote node and an authorized transactor
	conn, err := helper.NewSafeClient(ctx.String("eth-rpc-endpoint"))
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to connect to the Ethereum client: %v", err))
	}

	_, key := promptAccount(ctx.String("keystore-path"))
	fmt.Println("start to deploy ...")
	registryAddress := deployContract(key, conn)
	//registryAddress := common.HexToAddress("0xDe661C5aDaF15c243475C5c6BA96634983821593")
	if ctx.Bool("not-create-token") {
		return nil
	}
	registry, err := contracts.NewTokenNetworkRegistry(registryAddress, conn)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(tokenNumber)
	lock := &sync.Mutex{}
	for i := 0; i < tokenNumber; i++ {
		switch i {
		case 0:
			go func() {
				createTokenAndChannels(key, conn, registry, ctx.String("keystore-path"), !ctx.Bool("not-create-channel"), tokenERC223Approve, lock)
				wg.Done()
			}()
		case 3:
			go func() {
				createTokenAndChannels(key, conn, registry, ctx.String("keystore-path"), !ctx.Bool("not-create-channel"), tokenERC223, lock)
				wg.Done()
			}()
		case 2:
			go func() {
				createTokenAndChannels(key, conn, registry, ctx.String("keystore-path"), !ctx.Bool("not-create-channel"), tokenStandard, lock)
				wg.Done()
			}()
		case 1:
			go func() {
				createTokenAndChannels(key, conn, registry, ctx.String("keystore-path"), !ctx.Bool("not-create-channel"), tokenERC223, lock)
				wg.Done()
			}()
		}
	}
	wg.Wait()
	return nil
}
func promptAccount(keystorePath string) (addr common.Address, key *ecdsa.PrivateKey) {
	am := accounts.NewAccountManager(keystorePath)
	if len(am.Accounts) == 0 {
		log.Fatal(fmt.Sprintf("No Ethereum accounts found in the directory %s", keystorePath))
		os.Exit(1)
	}
	addr = am.Accounts[0].Address
	log.Printf("deploy account = %s", addr.String())
	log.Printf("accounts=%q", am.Accounts)
	for i := 0; i < 3; i++ {
		//fmt.Printf("\npassword is %s\n", password)
		keybin, err := am.GetPrivateKey(addr, passwords[0])
		if err != nil && i == 3 {
			log.Fatal(fmt.Sprintf("Exhausted passphrase unlock attempts for %s. Aborting ...", addr))
			os.Exit(1)
		}
		if err != nil {
			log.Println(fmt.Sprintf("password incorrect\n Please try again or kill the process to quit.\nUsually Ctrl-c."))
			continue
		}
		key, err = crypto.ToECDSA(keybin)
		if err != nil {
			log.Println(fmt.Sprintf("private key to bytes err %s", err))
		}
		break
	}
	return
}
func deployContract(key *ecdsa.PrivateKey, conn *helper.SafeEthClient) (tokenNetworkRegistryAddress common.Address) {
	chainID, err := conn.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("failed get chain id :%s", chainID)
	}
	fmt.Printf("current chain Id=%s\n", chainID)
	auth := bind.NewKeyedTransactor(key)
	//Deploy SecretRegistry
	SecretRegistryAddress, tx, _, err := contracts.DeploySecretRegistry(auth, conn)
	if err != nil {
		log.Fatalf("Failed to Deploy SecretRegistry : %v", err)
	}
	ctx := context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	fmt.Printf("deploy SecretRegistry[%s] complete...\n", SecretRegistryAddress.String())
	//Deploy TokenNetorkRegistry
	//auth.GasLimit = 4000000 //最大gas
	//auth.GasPrice = big.NewInt(2000)
	tokenNetworkRegistryAddress, tx, _, err = contracts.DeployTokenNetworkRegistry(auth, conn, SecretRegistryAddress, chainID)
	if err != nil {
		log.Fatalf("Failed to deploy new token contract: %v", err)
	}
	fmt.Printf("tokenNetworkRegistryAddress=%s, txhash=%s\n", tokenNetworkRegistryAddress.String(), tx.Hash().String())
	ctx = context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	fmt.Printf("Deploy tokenNetworkRegistry complete...\n")

	fmt.Printf("tokenNetworkRegistryAddress=%s\n", tokenNetworkRegistryAddress.String())
	return
}
func createTokenAndChannels(key *ecdsa.PrivateKey, conn *helper.SafeEthClient, registry *contracts.TokenNetworkRegistry, keystorepath string, createchannel bool, tokenType string, lock *sync.Mutex) {
	lock.Lock()
	tokenNetworkAddress, tokenAddress := newToken(key, conn, registry, tokenType)
	//tokenAddress := common.HexToAddress("0xD29A9Cbf2Ca88981D0794ce94e68495c4bC16F28")
	//tokenNetworkAddress, _ := registry.Token_to_token_networks(nil, tokenAddress)
	token, err := contracts.NewToken(tokenAddress, conn)
	if err != nil {
		log.Fatalf("err for newtoken err %s", err)
	}
	am := accounts.NewAccountManager(keystorepath)
	var localAccounts []common.Address
	var keys []*ecdsa.PrivateKey
	for _, account := range am.Accounts {
		var keybin []byte
		for _, p := range passwords {
			keybin, err = am.GetPrivateKey(account.Address, p)
			if err != nil {
				log.Printf("password error for %s,err=%s", utils.APex2(account.Address), err)
				continue
			} else {
				break
			}
		}
		if err != nil {
			log.Printf("password error for %s,err=%s", utils.APex2(account.Address), err)
			continue
		}
		keytemp, err := crypto.ToECDSA(keybin)
		if err != nil {
			log.Fatalf("toecdsa err %s", err)
			continue
		}
		keys = append(keys, keytemp)
		localAccounts = append(localAccounts, account.Address)
	}
	lock.Unlock()
	//createerc20token合约时间较长,导致多个token同时部署的时候Tx nonce会冲突
	time.Sleep(time.Second)
	//fmt.Printf("key=%s\n", key)
	transferMoneyForAccounts(key, conn, localAccounts[1:], keys[1:], token)
	if createchannel {
		createChannels(conn.Client, localAccounts, keys, tokenNetworkAddress, token)
	}
}

func newToken(key *ecdsa.PrivateKey, conn *helper.SafeEthClient, registry *contracts.TokenNetworkRegistry, tokenType string) (tokenNetworkAddress common.Address, tokenAddr common.Address) {
	var tx *types.Transaction
	var err error
	auth := bind.NewKeyedTransactor(key)
	switch tokenType {
	case tokenERC223:
		tokenAddr, tx, _, err = tokenerc223.DeployHumanERC223Token(auth, conn, getAmount(big.NewInt(5000000000000000000)), "test erc223", uint8(base))
	case tokenStandard:
		tokenAddr, tx, _, err = tokenstandard.DeployHumanStandardToken(auth, conn, getAmount(big.NewInt(5000000000000000000)), "test standard", uint8(base))
	case tokenERC223Approve:
		tokenAddr, tx, _, err = tokenerc223approve.DeployHumanERC223Token(auth, conn, getAmount(big.NewInt(5000000000000000000)), "test erc223 approve", uint8(base))
	case tokenEther:
		auth.Value = getAmount(big.NewInt(5000000000000000000))
		tokenAddr, tx, _, err = tokenether.DeployHumanEtherToken(auth, conn, "test ether")
	}
	if err != nil {
		log.Fatalf("Failed to deploy %s: %v,account=%s", tokenType, err, auth.From.String())
	}
	fmt.Printf("token deploy tx=%s\n", tx.Hash().String())
	ctx := context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	auth.Value = nil // ether will modify auth
	fmt.Printf("Deploy %s  %s complete...\n", tokenType, tokenAddr.String())
	tx, err = registry.CreateERC20TokenNetwork(auth, tokenAddr)
	if err != nil {
		log.Fatalf("Failed to AddToken: %v", err)
	}
	log.Printf("CreateERC20TokenNetwork tx=%s", tx.Hash().String())
	ctx = context.Background()
	_, err = bind.WaitMined(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to AddToken when mining :%v", err)
	}
	tokenNetworkAddress, err = registry.TokenToTokenNetworks(nil, tokenAddr)
	fmt.Printf("Deploy %s complete... %s,tokennetwork=%s\n", tokenType, tokenAddr.String(), tokenNetworkAddress.String())
	return
}
func transferMoneyForAccounts(key *ecdsa.PrivateKey, conn *helper.SafeEthClient, accounts []common.Address, keys []*ecdsa.PrivateKey, token *contracts.Token) {
	wg := sync.WaitGroup{}
	wg.Add(len(accounts))
	//auth := bind.NewKeyedTransactor(key)
	//nonce, err := conn.PendingNonceAt(context.Background(), auth.From)
	//if err != nil {
	//	log.Fatalf("pending nonce err %s", err)
	//	return
	//}
	for index, account := range accounts {
		go func(account common.Address, i int) {
			auth2 := bind.NewKeyedTransactor(key)
			//auth2.Nonce = big.NewInt(int64(nonce) + int64(i))
			fmt.Printf("transfer to %s,nonce=%s\n", account.String(), auth2.Nonce)
			//由于生成的 Transfer 不能很好处理重载,因此需要用 approve and transfer from
			amount := getAmount(big.NewInt(500000000000))
			tx, err := token.Approve(auth2, account, amount)
			if err != nil {
				log.Fatalf("Failed to Transfer: %v", err)
			}
			ctx := context.Background()
			_, err = bind.WaitMined(ctx, conn, tx)
			if err != nil {
				log.Fatalf("failed to Transfer when mining :%v", err)
			}
			fmt.Printf("approve %s complete\n", account.String())
			auth3 := bind.NewKeyedTransactor(keys[i])
			tx, err = token.TransferFrom(auth3, auth2.From, account, amount)
			if err != nil {
				log.Fatalf("Failed to Transfer: %v", err)
			}
			fmt.Printf("transfer from %s,txhash=%s\n", account.String(), tx.Hash().String())
			ctx = context.Background()
			_, err = bind.WaitMined(ctx, conn, tx)
			if err != nil {
				log.Fatalf("failed to Transfer when mining :%v", err)
			}
			fmt.Printf("Transfer complete...\n")
			wg.Done()
		}(account, index)
		time.Sleep(time.Millisecond * 10)
	}
	wg.Wait()
	for _, account := range accounts {
		b, err := token.BalanceOf(nil, account)
		if err != nil {
			log.Fatalf("balance of err %s", err)
		}
		log.Printf("account %s has token %s\n", utils.APex(account), b)
	}
}

//path A-B-C-F-B-D-G-E
func createChannels(conn *ethclient.Client, accounts []common.Address, keys []*ecdsa.PrivateKey, tokenNetworkAddress common.Address, token *contracts.Token) {
	if len(accounts) < 6 {
		panic("need 6 accounts")
	}
	AccountA := accounts[0]
	AccountB := accounts[1]
	AccountC := accounts[2]
	AccountD := accounts[3]
	AccountE := accounts[4]
	AccountF := accounts[5]
	AccountG := accounts[6]
	fmt.Printf("accountA=%s\naccountB=%s\naccountC=%s\naccountD=%s\naccountE=%s\naccountF=%s\naccountG=%s\n",
		AccountA.String(), AccountB.String(), AccountC.String(), AccountD.String(),
		AccountE.String(), AccountF.String(), AccountG.String())
	keyA := keys[0]
	keyB := keys[1]
	keyC := keys[2]
	keyD := keys[3]
	keyE := keys[4]
	keyF := keys[5]
	keyG := keys[6]
	wg := sync.WaitGroup{}
	wg.Add(8)
	//fmt.Printf("keya=%s,keyb=%s,keyc=%s,keyd=%s,keye=%s,keyf=%s,keyg=%s\n", keyA, keyB, keyC, keyD, keyE, keyF, keyG)
	go func() {
		createchannel.CreatAChannelAndDeposit(AccountA, AccountB, keyA, keyB, getAmount(big.NewInt(100)), tokenNetworkAddress, token, conn)
		wg.Done()
	}()
	go func() {
		createchannel.CreatAChannelAndDeposit(AccountB, AccountD, keyB, keyD, getAmount(big.NewInt(90)), tokenNetworkAddress, token, conn)
		wg.Done()
	}()
	go func() {
		createchannel.CreatAChannelAndDeposit(AccountG, AccountE, keyG, keyE, getAmount(big.NewInt(80)), tokenNetworkAddress, token, conn)
		wg.Done()
	}()
	go func() {
		createchannel.CreatAChannelAndDeposit(AccountD, AccountG, keyD, keyG, getAmount(big.NewInt(190)), tokenNetworkAddress, token, conn)
		wg.Done()
	}()
	go func() {
		createchannel.CreatAChannelAndDeposit(AccountC, AccountE, keyC, keyE, getAmount(big.NewInt(10)), tokenNetworkAddress, token, conn)
		wg.Done()
	}()
	go func() {
		createchannel.CreatAChannelAndDeposit(AccountC, AccountF, keyC, keyF, getAmount(big.NewInt(60)), tokenNetworkAddress, token, conn)
		wg.Done()
	}()
	go func() {
		createchannel.CreatAChannelAndDeposit(AccountB, AccountF, keyB, keyF, getAmount(big.NewInt(70)), tokenNetworkAddress, token, conn)
		wg.Done()
	}()
	go func() {
		createchannel.CreatAChannelAndDeposit(AccountB, AccountC, keyB, keyC, getAmount(big.NewInt(50)), tokenNetworkAddress, token, conn)
		wg.Done()
	}()

	if len(accounts) >= 7 {
		wg.Add(len(accounts) - 1 - 6)
	}
	for i := 6; i < len(accounts)-1; i++ {
		go func(index int) {
			createchannel.CreatAChannelAndDeposit(accounts[index], accounts[index+1], keys[index], keys[index+1],
				getAmount(big.NewInt(100)), tokenNetworkAddress, token, conn,
			)
			wg.Done()
		}(i)

	}
	wg.Wait()
}
