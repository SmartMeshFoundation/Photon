package main

import (
	"fmt"

	"os"

	"encoding/hex"
	"path/filepath"

	"crypto/ecdsa"

	"bytes"

	"github.com/SmartMeshFoundation/SmartRaiden/accounts"
	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/models"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mtree"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/node"
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
	app.Action = mainctx
	app.Name = "withdraw"
	app.Version = "0.1"
	err := app.Run(os.Args)
	if err != nil {
		log.Crit(err.Error())
	}
}

type withDraw struct {
	Address                   common.Address
	Conn                      *helper.SafeEthClient
	PrivateKey                *ecdsa.PrivateKey
	DbPath                    string
	bcs                       *rpc.BlockChainService
	db                        *models.ModelDB
	WithDrawChannelIdentifier common.Hash
	Secret                    common.Hash
	ChannelIdentifier2Channel map[common.Hash]*channel.Channel
}

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
}
func mainctx(ctx *cli.Context) error {
	var err error
	w := &withDraw{
		ChannelIdentifier2Channel: make(map[common.Hash]*channel.Channel),
	}
	address := common.HexToAddress(ctx.String("address"))
	if address == utils.EmptyAddress {
		log.Crit("must specified a valid address")
	}
	w.Address = address
	//db path
	userDbPath := hex.EncodeToString(address[:])
	userDbPath = userDbPath[:8]
	w.DbPath = filepath.Join(ctx.String("datadir"), userDbPath)
	w.DbPath = filepath.Join(w.DbPath, "log.db")
	if !utils.Exists(w.DbPath) {
		log.Crit("data directory is invalid ,doesn't contain db")
	}
	w.openDb()
	if err != nil {
		return err
	}
	_, key, err := accounts.PromptAccount(address, ctx.String("keystore-path"), ctx.String("password-file"))
	if err != nil {
		log.Crit(fmt.Sprintf("unlock acccount err %s", err))
	}
	privateKey, err := crypto.ToECDSA(key)
	if err != nil {
		log.Crit("private key is invalid, wrong password?")
	}
	w.PrivateKey = privateKey
	config := &params.Config{
		EthRPCEndPoint: ctx.String("eth-rpc-endpoint"),
		PrivateKey:     privateKey,
	}
	w.bcs, err = rpc.NewBlockChainService(config, w.db)
	if err != nil {
		return err
	}
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlDebug, utils.MyStreamHandler(os.Stderr)))
	// Create an IPC based RPC connection to a remote node and an authorized transactor
	w.Conn = w.bcs.Client
	w.WithDrawChannelIdentifier = common.HexToHash(ctx.String("channel"))
	w.Secret = common.HexToHash(ctx.String("secret"))
	if w.WithDrawChannelIdentifier == utils.EmptyHash || w.Secret == utils.EmptyHash {
		log.Crit("channel and secret muse be specified.")
	}
	err = w.restoreChannel()
	if err != nil {
		log.Error(fmt.Sprintf("restore channel %s", err))
	}
	log.Info("withdraw on channel...")
	w.WithDrawOnChannel()
	return nil
}

func (w *withDraw) openDb() {
	var err error
	w.db, err = models.OpenDb(w.DbPath)
	if err != nil {
		log.Crit("cannot open db")
	}
}
func (w *withDraw) channelSerilization2Channel(c *channeltype.Serialization, tokenNetwork *rpc.TokenNetworkProxy) (ch *channel.Channel, err error) {
	OurState := channel.NewChannelEndState(c.OurAddress, c.OurContractBalance,
		c.OurBalanceProof, mtree.NewMerkleTree(c.OurLeaves))
	PartnerState := channel.NewChannelEndState(c.PartnerAddress(),
		c.PartnerContractBalance,
		c.PartnerBalanceProof, mtree.NewMerkleTree(c.PartnerLeaves))
	ExternState := channel.NewChannelExternalState(nil, tokenNetwork,
		c.ChannelIdentifier, w.PrivateKey,
		w.Conn, w.db, c.ClosedBlock,
		c.OurAddress, c.PartnerAddress())
	ch, err = channel.NewChannel(OurState, PartnerState, ExternState, c.TokenAddress(), c.ChannelIdentifier, c.RevealTimeout, c.SettleTimeout)
	if err != nil {
		return
	}

	ch.OurState.Lock2PendingLocks = c.OurLock2PendingLocks()
	ch.OurState.Lock2UnclaimedLocks = c.OurLock2UnclaimedLocks()
	ch.PartnerState.Lock2PendingLocks = c.PartnerLock2PendingLocks()
	ch.PartnerState.Lock2UnclaimedLocks = c.PartnerLock2UnclaimedLocks()
	ch.State = c.State
	ch.OurState.ContractBalance = c.OurContractBalance
	ch.PartnerState.ContractBalance = c.PartnerContractBalance
	ch.ExternState.ClosedBlock = c.ClosedBlock
	ch.ExternState.SettledBlock = c.SettledBlock
	return
}
func (w *withDraw) getTokenNetworkProxy(tokenAddress common.Address) (tokenNetwork *rpc.TokenNetworkProxy, err error) {
	registryAddress := w.db.GetRegistryAddress()
	r := w.bcs.Registry(registryAddress)
	tokenNetworkAddr, err := r.TokenNetworkByToken(tokenAddress)
	if err != nil {
		return
	}
	tokenNetwork, err = w.bcs.TokenNetwork(tokenNetworkAddr)
	return
}
func (w *withDraw) restoreChannel() error {
	var err error
	allChannels, err := w.db.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil {
		log.Crit(fmt.Sprintf("get channel list err %s", err))
		return err
	}
	for _, cs := range allChannels {
		if bytes.Compare(cs.Key, w.WithDrawChannelIdentifier[:]) != 0 {
			//continue
		}
		tn, err := w.getTokenNetworkProxy(cs.TokenAddress())
		if err != nil {
			log.Crit(fmt.Sprintf("getTokenNetworkProxy err %s", err))
			return err
		}
		c, err := w.channelSerilization2Channel(cs, tn)
		if err != nil {
			log.Info(fmt.Sprintf("ignore channel %s, maybe has been settled", utils.BPex(cs.Key)))
			continue
		}
		w.ChannelIdentifier2Channel[common.BytesToHash(cs.Key)] = c
	}
	return nil
}
func (w *withDraw) WithDrawOnChannel() {
	for addr, c := range w.ChannelIdentifier2Channel {
		if addr == w.WithDrawChannelIdentifier {
			err := c.RegisterSecret(w.Secret)
			if err != nil {
				log.Error(fmt.Sprintf("regist secret %s on channel %s error %s", utils.HPex(w.Secret), utils.HPex(w.WithDrawChannelIdentifier), err))
				return
			}
			isReg, err := w.bcs.SecretRegistryProxy.IsSecretRegistered(w.Secret)
			if err != nil {
				log.Error(fmt.Sprintf("IsSecretRegistered err %s", err))
				return
			}
			if !isReg {
				err = w.bcs.SecretRegistryProxy.RegisterSecret(w.Secret)
				if err != nil {
					log.Error(fmt.Sprintf("RegisterSecret %s", err))
				}
			}
			//todo 链上注册密码
			result := c.Close()
			err = <-result.Result
			if err != nil {
				log.Error(fmt.Sprintf("close channel %s error %s", c.ChannelIdentifier.String(), err))
				break
			}
			unlockProofs2 := c.PartnerState.GetKnownUnlocks()
			result = c.ExternState.Unlock(unlockProofs2, c.PartnerState.TransferAmount())
			err = <-result.Result
			if err != nil {
				log.Error(err.Error())
			}
			break
		}
	}
}
