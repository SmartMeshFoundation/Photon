package main

import (
	"fmt"

	"os"

	"io/ioutil"

	"encoding/hex"
	"path/filepath"

	"crypto/ecdsa"
	"math/big"

	smartraiden "github.com/SmartMeshFoundation/SmartRaiden"
	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/models"
	"github.com/SmartMeshFoundation/SmartRaiden/network"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/node"
	"github.com/slonzok/getpass"
	"github.com/urfave/cli"
)

/*
withdraw on a hash
*/
var dbpath string

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "address",
			Usage: "The ethereum address you would like raiden to use and for which a keystore file exists in your local system.",
			Value: utils.EmptyAddress.String(),
		},
		cli.StringFlag{
			Name:  "keystore-path",
			Usage: "If you have a non-standard path for the ethereum keystore directory provide it using this argument. ",
			//Value: ethutils.DirectoryString{params.DefaultKeyStoreDir()},
			Value: utils.GetHomePath() + "/privnet3/keystore",
		},
		cli.StringFlag{
			Name: "eth-rpc-endpoint",
			Usage: `"host:port" address of ethereum JSON-RPC server.\n'
	           'Also accepts a protocol prefix (ws:// or ipc channel) with optional port',`,
			Value: fmt.Sprintf("ws://%s", node.DefaultWSEndpoint()),
		},
		cli.StringFlag{
			Name:  "datadir",
			Usage: "Directory for storing raiden data.",
			Value: params.DefaultDataDir(),
		},
		cli.StringFlag{
			Name:  "password-file",
			Usage: "Text file containing password for provided account",
		},
		cli.StringFlag{
			Name:  "channel",
			Usage: "withdraw all the hashlock of this channel ,default is 0x00000",
			Value: utils.EmptyAddress.String(),
		},
		cli.StringFlag{
			Name:  "secret",
			Usage: "register this secret and withdraw ",
			Value: utils.EmptyHash.String(),
		},
	}
	app.Action = Main
	app.Name = "withdraw"
	app.Version = "0.1"
	app.Run(os.Args)
}

type WithDraw struct {
	Address                common.Address
	Conn                   *helper.SafeEthClient
	PrivateKey             *ecdsa.PrivateKey
	DbPath                 string
	bcs                    *rpc.BlockChainService
	db                     *models.ModelDB
	WithDrawChannelAddress common.Address
	Secret                 common.Hash
	ChannelAddress2Channel map[common.Address]*channel.Channel
}

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
}
func Main(ctx *cli.Context) error {
	var err error
	w := &WithDraw{
		ChannelAddress2Channel: make(map[common.Address]*channel.Channel),
	}
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlDebug, utils.MyStreamHandler(os.Stderr)))
	// Create an IPC based RPC connection to a remote node and an authorized transactor
	w.Conn, err = helper.NewSafeClient(ctx.String("eth-rpc-endpoint"))
	if err != nil {
		log.Crit(fmt.Sprintf("Failed to connect to the Ethereum client: %v", err))
	}
	w.WithDrawChannelAddress = common.HexToAddress(ctx.String("channel"))
	w.Secret = common.HexToHash(ctx.String("secret"))
	if w.WithDrawChannelAddress == utils.EmptyAddress || w.Secret == utils.EmptyHash {
		log.Crit("channel and secret muse be specified.")
	}
	address := common.HexToAddress(ctx.String("address"))
	if address == utils.EmptyAddress {
		log.Crit("must specified a valid address")
	}
	w.Address = address
	_, key := promptAccount(address, ctx.String("keystore-path"), ctx.String("password-file"))
	privateKey, err := crypto.ToECDSA(key)
	if err != nil {
		log.Crit("private key is invalid, wrong password?")
	}
	w.PrivateKey = privateKey

	//db path
	userDbPath := hex.EncodeToString(address[:])
	userDbPath = userDbPath[:8]
	w.DbPath = filepath.Join(ctx.String("datadir"), userDbPath)
	w.DbPath = filepath.Join(w.DbPath, "log.db")
	if !utils.Exists(w.DbPath) {
		log.Crit("data directory is invalid ,doesn't contain db")
	}
	w.bcs = rpc.NewBlockChainService(privateKey, utils.EmptyAddress, w.Conn)
	w.openDb()
	w.restoreChannel()
	log.Info("withdraw on channel...")
	w.WithDrawOnChannel()
	return nil
}
func promptAccount(adviceAddress common.Address, keystorePath, passwordfile string) (addr common.Address, keybin []byte) {
	am := smartraiden.NewAccountManager(keystorePath)
	if len(am.Accounts) == 0 {
		log.Error(fmt.Sprintf("No Ethereum accounts found in the directory %s", keystorePath))
		utils.SystemExit(1)
	}
	if !am.AddressInKeyStore(adviceAddress) {
		if adviceAddress != utils.EmptyAddress {
			log.Error(fmt.Sprintf("account %s could not be found on the sytstem. aborting...", adviceAddress))
			utils.SystemExit(1)
		}
		shouldPromt := true
		fmt.Println("The following accounts were found in your machine:")
		for i := 0; i < len(am.Accounts); i++ {
			fmt.Printf("%3d -  %s\n", i, am.Accounts[i].Address.String())
		}
		fmt.Println("")
		for shouldPromt {
			fmt.Printf("Select one of them by index to continue:\n")
			idx := -1
			fmt.Scanf("%d", &idx)
			if idx >= 0 && idx < len(am.Accounts) {
				shouldPromt = false
				addr = am.Accounts[idx].Address
			} else {
				fmt.Printf("Error: Provided index %d is out of bounds", idx)
			}
		}
	} else {
		addr = adviceAddress
	}
	var password string
	var err error
	if len(passwordfile) > 0 {
		data, err := ioutil.ReadFile(passwordfile)
		if err != nil {
			log.Error(fmt.Sprintf("password_file error:%s", err))
			utils.SystemExit(1)
		}
		password = string(data)
		log.Trace(fmt.Sprintf("password is %s", password))
		keybin, err = am.GetPrivateKey(addr, password)
		if err != nil {
			log.Error(fmt.Sprintf("Incorrect password for %s in file. Aborting ... %s", addr.String(), err))
			utils.SystemExit(1)
		}
	} else {
		for i := 0; i < 3; i++ {
			//retries three times
			password = getpass.Prompt("Enter the password to unlock:")
			keybin, err = am.GetPrivateKey(addr, password)
			if err != nil && i == 3 {
				log.Error(fmt.Sprintf("Exhausted passphrase unlock attempts for %s. Aborting ...", addr))
				utils.SystemExit(1)
			}
			if err != nil {
				log.Error(fmt.Sprintf("password incorrect\n Please try again or kill the process to quit.\nUsually Ctrl-c."))
				continue
			}
			break
		}
	}
	return
}

func (this *WithDraw) getChannelDetail(proxy *rpc.NettingChannelContractProxy) *network.ChannelDetails {
	addr1, b1, addr2, b2, err := proxy.AddressAndBalance()
	if err != nil {
		log.Error(fmt.Sprintf("AddressAndBalance err %s", err))
	}
	var ourAddr, partnerAddr common.Address
	var ourBalance, partnerBalance *big.Int
	if addr1 == this.Address {
		ourAddr = addr1
		partnerAddr = addr2
		ourBalance = b1
		partnerBalance = b2
	} else {
		ourAddr = addr2
		partnerAddr = addr1
		ourBalance = b2
		partnerBalance = b1
	}
	ourState := channel.NewChannelEndState(ourAddr, ourBalance, nil, transfer.EmptyMerkleTreeState)
	partenerState := channel.NewChannelEndState(partnerAddr, partnerBalance, nil, transfer.EmptyMerkleTreeState)
	channelAddress := proxy.Address
	registerChannelForHashlock := func(channel *channel.Channel, hashlock common.Hash) {

	}
	externState := channel.NewChannelExternalState(registerChannelForHashlock, proxy, channelAddress, this.bcs, this.db)
	channelDetail := &network.ChannelDetails{
		ChannelAddress:    channelAddress,
		OurState:          ourState,
		PartenerState:     partenerState,
		ExternState:       externState,
		BlockChainService: this.bcs,
		RevealTimeout:     params.DefaultRevealTimeout,
	}
	channelDetail.SettleTimeout, err = externState.NettingChannel.SettleTimeout()
	if err != nil {
		log.Error(fmt.Sprintf("SettleTimeout query err %s", err))
	}
	return channelDetail
}

func (this *WithDraw) NewChannel(channelAddress common.Address) (c *channel.Channel, err error) {
	proxy, err := this.bcs.NettingChannel(channelAddress)
	if err != nil {
		return
	}
	detail := this.getChannelDetail(proxy)
	c, err = channel.NewChannel(detail.OurState, detail.PartenerState,
		detail.ExternState, utils.EmptyAddress, channelAddress, this.bcs, detail.RevealTimeout, detail.SettleTimeout)
	return
}

func (this *WithDraw) openDb() {
	var err error
	this.db, err = models.OpenDb(this.DbPath)
	if err != nil {
		log.Crit("cannot open db")
	}
}

func (this *WithDraw) restoreChannel() error {
	var err error
	allChannels, err := this.db.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil {
		log.Crit(fmt.Sprintf("get channel list err %s", err))
		return err
	}
	for _, cs := range allChannels {
		if cs.ChannelAddress == this.WithDrawChannelAddress {
			//log.Info(fmt.Sprintf("db channel=%s", utils.StringInterface(cs, 5)))
		}
		c, err := this.NewChannel(cs.ChannelAddress)
		if err != nil {
			log.Info(fmt.Sprintf("ignore channel %s, maybe has been settled", utils.APex(cs.ChannelAddress)))
			continue
		}

		if cs.OurAddress != c.OurState.Address ||
			cs.PartnerAddress != c.PartnerState.Address {
			log.Error(fmt.Sprintf("snapshot data error, channel data error for  db=%s,contract=%s ", utils.StringInterface(cs, 3), utils.StringInterface(c, 3)))
			continue
		} else {
			c.OurState.BalanceProofState = cs.OurBalanceProof
			c.OurState.TreeState = transfer.NewMerkleTreeStateFromLeaves(cs.OurLeaves)
			c.OurState.Lock2PendingLocks = cs.OurLock2PendingLocks
			c.OurState.Lock2UnclaimedLocks = cs.OurLock2UnclaimedLocks
			c.PartnerState.BalanceProofState = cs.PartnerBalanceProof
			c.PartnerState.TreeState = transfer.NewMerkleTreeStateFromLeaves(cs.PartnerLeaves)
			c.PartnerState.Lock2PendingLocks = cs.PartnerLock2PendingLocks
			c.PartnerState.Lock2UnclaimedLocks = cs.PartnerLock2UnclaimedLocks
		}
		this.ChannelAddress2Channel[cs.ChannelAddress] = c
	}
	return nil
}
func (this *WithDraw) WithDrawOnChannel() {
	for addr, c := range this.ChannelAddress2Channel {
		if addr == this.WithDrawChannelAddress {
			err := c.RegisterSecret(this.Secret)
			if err != nil {
				log.Error(fmt.Sprint("regist secret %s on channel %s error %s", utils.HPex(this.Secret), utils.APex(this.WithDrawChannelAddress), err))
				return
			}
			err = c.ExternState.Close(c.PartnerState.BalanceProofState)
			if err != nil {
				log.Error(fmt.Sprintf("close channel %s error %s", utils.APex(c.MyAddress), err))
				break
			}
			unlockProofs2 := c.PartnerState.GetKnownUnlocks()
			err = c.ExternState.WithDraw(unlockProofs2)
			if err != nil {
				log.Error(err.Error())
			}
			break
		}
	}
}
