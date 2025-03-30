package restful

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"go.cryptoscope.co/muxrpc/v2"
	"go.cryptoscope.co/netwrap"
	"go.cryptoscope.co/secretstream"
	ssbClient "go.cryptoscope.co/ssb/client"
	"go.cryptoscope.co/ssb/restful/params"
	kitlog "go.mindeco.de/log"
	"go.mindeco.de/log/level"
	"golang.org/x/crypto/ed25519"
	"gopkg.in/urfave/cli.v2"

	/*"go.cryptoscope.co/ssb/message"
	"go.mindeco.de/ssb-refs"*/

	"bufio"
	"os"

	"math"

	"errors"

	"sync"

	"runtime"

	"net/url"

	"go.cryptoscope.co/ssb"
	"go.cryptoscope.co/ssb/dfa"
	"go.cryptoscope.co/ssb/message"
	refs "go.mindeco.de/ssb-refs"
)

var lock sync.RWMutex

var Config *params.ApiConfig

var longCtx context.Context

var channelAllSig chan interface{}

var quitSignal chan struct{}

var client *ssbClient.Client

var log kitlog.Logger

var lastAnalysisTimesnamp int64

var likeDB *PubDB

var dfax *dfa.DFA

var PhotonNodeStatus *PhotonNodeStatusStu

var pubNode *PhotonNode
var checkChannelNode *PhotonNode
var checkPayNode *PhotonNode

// 限制通道建立的次数
var ChannelDepositTimes map[string]int

const (
	SignUp                = "sign up"
	PostMessage           = "post message"
	PostComment           = "post comment"
	MintNft               = "mint a nft"
	DailyLogin            = "daily login"
	LikePost              = "like a post"
	ReceiveLike           = "receive a like"
	ReportProblematicPost = "report problematic post"
	GameEarn              = "game to earn"
	InviteEarn            = "invite to earn"
)

// Start
func Start(ctx *cli.Context) {
	Config = params.NewApiServeConfig()
	longCtx = ctx

	sclient, err := newClient(ctx)
	if err != nil {
		//level.Error(log).Log("Ssb restful api and message analysis service start err", err)
		fmt.Println(fmt.Sprintf("Ssb restful api and message analysis service start err=%s", err))
		return
	}
	client = sclient

	quitSignal := make(chan struct{})
	api := rest.NewApi()
	api.Use(rest.DefaultCommonStack...)
	/*if Config.Debug {
		api.Use(rest.DefaultDevStack...)
	} else {
		api.Use(rest.DefaultProdStack...)
	}*/
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(

		/*
			ssb pub信息
		*/
		//pub's whoami
		rest.Get("/ssb/api/pub-whoami", GetPubWhoami),

		/*
			ssb节点注册,信息查询,例如查询其绑定的钱包地址
		*/
		//get all 'about' message,e.g:'about'='eth address'
		rest.Get("/ssb/api/node-info", clientid2Profiles),
		//get the 'about' message by client id ,e.g:'about'='eth address'
		rest.Post("/ssb/api/node-info", clientid2Profile),
		//register client's eth address to it's ID
		rest.Post("/ssb/api/id2eth", UpdateEthAddr),

		/*
			受赞统计
		*/
		//likes of all client
		rest.Get("/ssb/api/likes", GetAllLikes),
		//likes of someone client
		rest.Post("/ssb/api/likes", GetSomeoneLike),

		/*
			点赞统计
		*/
		//get set like infos of all
		rest.Get("/ssb/api/set-like-info", GetAllSetLikes),
		//get set like info of someone client
		rest.Post("/ssb/api/set-like-info", GetSomeoneSetLikes),

		/*
			举报
		*/
		// tipped someone off 举报
		rest.Post("/ssb/api/tipped-who-off", TippedOff),
		//tipped off infomation 所有举报的信息汇总
		rest.Post("/ssb/api/tippedoff-info", GetTippedOffInfo),
		//tippedoff-deal pub管理员对举报的信息进行处理，认证，如属实，则对该账号进行黑名单处理
		rest.Post("/ssb/api/tippedoff-deal", DealTippedOff),

		/*
			敏感词
		*/
		//DealSensitiveWord pub管理对敏感词的处理/block or ignore
		rest.Post("/ssb/api/sensitive-word-deal", DealSensitiveWord),
		//get all sensitive-word-events from pub
		rest.Post("/ssb/api/sensitive-word-events", GetEventSensitiveWord),

		/*
			用户每日任务,数据类型：1-登录 2-发帖(Pub自动处理) 3-评论(Pub自动处理) 4-铸造NFT
		*/
		//notify pub the login infomation, pub will collect through this interface
		rest.Post("/ssb/api/notify-login", NotifyUserLogin),
		//[temporary scheme] notify the pub that user have created a NFT in metalife app
		rest.Post("/ssb/api/notify-created-nft", NotifyCreatedNFT),
		//get some user daily task infos from pub,
		//a message may appear in multiple pubs, and the client removes redundant data through messagekey and pub id
		//used by supernode to awarding or ssb-client
		rest.Post("/ssb/api/get-user-daily-task", GetUserDailyTasks),

		/*
			激励查询
		*/
		//get all or someones' reward information in PUB RULE
		rest.Post("/ssb/api/get-reward-info", GetRewardInfo),

		rest.Post("/ssb/api/get-reward-subtotals", GetRewardSubtotals),

		rest.Post("/ssb/api/get-reward-summary", GetRewardSummary),

		/*
			通过归属地获取最近的pub接入方式和邀请码
		*/
		rest.Get("/ssb/api/get-pubhost-by-ip", GetPublicIPLocation),
		/*
			用户反馈/建议
		*/
		rest.Post("/ssb/api/user/feedback", UserFeedBack),
		rest.Post("/ssb/api/user/get-feedback", GetUserFeedBack),

		/*
			通知pub强制关闭通道
		*/
		rest.Post("/ssb/api/user/notify-force-close-channel", NotifyForceCloseChannel),

		rest.Post("/ssb/api/user/get-friends-maybe", GetFriendMaybe),

		/*
			给用户直接发送邀请激励
		*/
		rest.Post("/ssb/api/internal/reward-for-some-reason", rewardForSomeReason),
	)
	if err != nil {
		level.Error(log).Log("make router err", err)
		return
	}

	api.SetApp(router)

	listen := fmt.Sprintf("%s:%d", Config.Host, Config.Port)
	server := &http.Server{Addr: listen, Handler: api.MakeHandler()}
	go server.ListenAndServe()
	fmt.Println(fmt.Sprintf(PrintTime() + "ssb restful api and message analysis service start...\nWelcome..."))

	//=======================================
	api1 := rest.NewApi()

	api1.Use(rest.DefaultCommonStack...)

	router1, err := rest.MakeRouter(
		/*
			获取游戏列表
			提交游戏凭证
			查询我的游戏提交记录
			查询游戏报酬记录
			审核游戏提交凭证，以发放激励
			上传游戏资料
		*/
		rest.Post("/ssb/api/user/upload-game-info", UploadGameInfo),
		rest.Get("/ssb/api/user/load-game-info", LoadGameInfo),
		rest.Get("/ssb/api/user/get-resource/:gamename/:resourcename", GetResource),
		//rest.Get("/ssb/api/user/get-user-photo/:eth-address/:pic-name", GetUserPhoto),
		rest.Post("/ssb/api/user/upoad-game-play", UploadGamePlay),
		rest.Post("/ssb/api/user/get-game-play", GetGamePlay),
		rest.Post("/ssb/api/user/game-earn-deal", DealGameEarn),
		rest.Post("/ssb/api/user/get-play-earn", GetPlayEarn),
	)
	if err != nil {
		level.Error(log).Log("make router err", err)
		return
	}
	api1.SetApp(router1)
	listen1 := fmt.Sprintf("0.0.0.0:%d", 10009)
	server1 := &http.Server{Addr: listen1, Handler: api1.MakeHandler()}
	go server1.ListenAndServe()
	//=======================================

	http.HandleFunc("/ssb/api/user/get-user-photo/:eth-address/:pic-name", GetUserPhoto)
	go http.ListenAndServe("0.0.0.0:10010", nil)

	channelAllSig = make(chan interface{}, 9)
	/*for {
		//v2版来调整
		select {
		case <-channelAllSig:
			SignalMessageQu()
		case <-time.After(time.Second):
		default:
			//更新所有回调的接口,异步data和并发
			Improve()
		}
	}*/
	go DoMessageTask(ctx) //todo

	//go dealBlacklist()

	//检查pub 与 所有metalife内已注册eth地址的账户的通道余额，按规定补充
	go PubCheckBalance() //todo

	//维护在线状态和通道补发激励
	go GetPeerstatusAndPubBackPay() //todo

	go func() { //不能关，没有通道无法订阅在线状态
		oneday := 30 * 24 * time.Hour
		ticker := time.NewTicker(oneday)
		for range ticker.C {
			useri, err := likeDB.SelectUserProfile("")
			if err != nil {
				fmt.Println(fmt.Sprintf(PrintTime()+"Failed to Close inactive peer's Channel,err: %s", err))
			}
			nowTime := time.Now().UnixNano() / 1e6
			for _, u := range useri {
				lat := u.LastactiveTime
				if (nowTime - lat) > time.Hour.Milliseconds()*24*90 {
					closeObj := u.EthAddress
					objAddr, err := HexToAddress(closeObj)
					if err != nil {
						continue
					} else {
						err = pubNode.Close(params.TokenAddress, objAddr.String())
						if err != nil {
							fmt.Println(fmt.Sprintf(PrintTime()+"Failed to Close inactive peer %s Channel,err: %s", objAddr.String(), err))
						} else {
							fmt.Println(fmt.Sprintf(PrintTime()+"Success to Close inactive peer %s Channel", objAddr.String()))
						}
					}
				}
			}
		}
	}()
	/*{
		go func() {

			oneday := 1 * time.Minute
			ticker := time.NewTicker(oneday)
			for range ticker.C {

				channelAll, err := checkChannelNode.GetChannels(params.TokenAddress)
				if err != nil {
					fmt.Println(fmt.Errorf(PrintTime()+"[close all channel] GetChannels err=%s", err))
					continue
				}

				//实时核对，计算量过大，改为仅查通道，不看注册
				for i := 0; i < len(channelAll); i++ {
					partnerAddr := channelAll[i].PartnerAddress
					channelState := channelAll[i].State
					if channelState == 2 {
						chanid := channelAll[i].ChannelIdentifier
						err = pubNode.Settle(chanid, params.SettleTime)
						if err != nil {
							fmt.Println(fmt.Sprintf(PrintTime()+"Failed to Settle all peer %s Channel,err: %s", partnerAddr, err))
						} else {
							fmt.Println(fmt.Sprintf(PrintTime()+"Success to Settle all peer %s Channel", partnerAddr))
						}

					}

					time.Sleep(time.Millisecond * 50)
				}
			}
		}()
	}*/

	RegisterSourceMap = make(map[string]int64)

	ChannelDepositTimes = make(map[string]int)

	<-quitSignal
	err = server.Shutdown(context.Background())
	if err != nil {
		fmt.Println(fmt.Sprintf(PrintTime()+"API restful service Shutdown err : %s", err))
	}

}

type PeerStatus struct {
	clientid   string
	clientAddr string
	status     string
}

var userChan chan interface{}
var onlineChan chan interface{}
var exitChan chan bool

func GetPeerstatusAndPubBackPay() { //通过上线来决定补发
	time.Sleep(time.Second * 15)
	for {
		useri, err := likeDB.SelectUserProfile("")
		if err != nil {
			fmt.Println(fmt.Sprintf(PrintTime()+"Failed to GetPeerstatus initPeers,err: %s", err))
		}
		userChan = make(chan interface{}, len(useri))
		onlineChan = make(chan interface{}, len(useri))
		exitChan = make(chan bool, 1) //阻塞，等待所有的routine都完成的信号

		go initPeers(useri)
		checkChannelNum := 20 //1 + len(useri)/8 //使用100个协程qu读取所有peer的状态
		for i := 0; i < checkChannelNum; i++ {
			go getStatusChan()
		}
		go func() {
			for i := 0; i < checkChannelNum; i++ {
				<-exitChan
			}
			close(onlineChan)
		}()

		//对在线的节点补发token和通道检查
		for {
			fmt.Println(fmt.Sprintf(PrintTime()+"runtime.NumGoroutine= %d", runtime.NumGoroutine()))
			ons, ok := <-onlineChan
			if !ok {
				break
			}
			onlUser := ons.(PeerStatus)
			partnerSsbid := onlUser.clientid
			partnerAddress := onlUser.clientAddr
			time.Sleep(time.Second)
			fmt.Println(fmt.Sprintf("[PubBackPay] clientid = %s, eth-addr = %s --OnLine", partnerSsbid, partnerAddress))

			//lock.Lock()
			offinfos, err := likeDB.QueryGrantFail(partnerSsbid)
			if err != nil {
				fmt.Println(fmt.Errorf(PrintTime()+"[PubBackPay] QueryGrantFail err= %s", err))
				continue
			}
			for _, info := range offinfos {
				partnerSsbid := info.ClientID
				partnerAddress := info.ClientEthAddress //可能某个用户拥有多个钱包
				//1、没有通道的重开通道;
				//2、通道被关的(closed),先结算后重开通道
				if partnerSsbid != backpackObj {
					fmt.Println(PrintTime() + " [PubBackPay] for clientid=" + partnerSsbid + ", eth-address=" + partnerAddress)
					backpackObj = partnerSsbid
				}
				/*channelX, err := checkPayNode.GetChannelWith(&PhotonNode{
					Address: partnerAddress,
				}, params.TokenAddress)
				if err != nil {
					fmt.Println(fmt.Errorf(PrintTime()+"[PubBackPay] GetChannelWith err= %s", err))
					continue
				}
				if channelX == nil {
					err = checkPayNode.OpenChannel(partnerAddress, params.TokenAddress, new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(int64(params.MinBalanceInchannel))), params.SettleTime)
					if err != nil {
						fmt.Println(fmt.Errorf(PrintTime()+"[PubBackPay] clientid= %s, eth-address= %s, OpenChannel err=%s", partnerSsbid, partnerAddress, err))
						continue
					}
					fmt.Println(fmt.Sprintf(PrintTime()+"[PubBackPay] clientid= %s, eth-address= %s, OpenChannel SUCCESS", partnerSsbid, partnerAddress))
				}*/

				//补发
				msgTime := info.MessageTime
				amount := info.GrantTokenAmount
				reason := info.RewardReason
				_, err = likeDB.UpdateRewardResult(partnerSsbid, partnerAddress, "success", msgTime, time.Now().UnixNano()/1e6)
				if err != nil {
					fmt.Println(fmt.Errorf(PrintTime()+"[PubBackPay] clientid= %s, eth-address= %s, UpdateRewardResult err= %s", partnerSsbid, partnerAddress, err))
					continue
				}
				err = checkPayNode.SendTrans(params.TokenAddress, amount, partnerAddress, true, false)
				if err != nil {
					_, err1 := likeDB.UpdateRewardResult(partnerSsbid, partnerAddress, "fail", msgTime, msgTime)
					fmt.Println(fmt.Errorf(PrintTime()+"[PubBackPay] clientid= %s, eth-address= %s, SendTrans err= %s, UpdateRewardResult= %s", partnerSsbid, partnerAddress, err, err1))
					continue
				}
				//对sign up 补发SMT激励
				if reason == SignUp {
					smtAmount := new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(int64(params.RewardOfSignupSMT)))
					//smtAmount = new(big.Int).Div(smtAmount, big.NewInt(10))
					err = checkPayNode.TransferSMT(partnerAddress, smtAmount.String())
					if err != nil {
						continue
					}
					fmt.Println(fmt.Sprintf(PrintTime()+"[PubBackPay] (SMT) clientid= %s, eth-address= %s, amount= %v, err= %v", partnerSsbid, partnerAddress, amount, err))
				}
				fmt.Println(fmt.Sprintf(PrintTime()+"[PubBackPay] clientid= %s, eth-address= %s, SUCCESS", partnerSsbid, partnerAddress))
			}
			fmt.Println(fmt.Sprintf(PrintTime()+"[PubBackPay] Pub check node's backpay info,coming in....%d", len(offinfos)))
		}
		time.Sleep(time.Second * 2)
	}

}
func getStatusChan() {
	time.Sleep(time.Second)
	for {
		uinfo, ok := <-userChan
		if !ok {
			break
		}
		user := uinfo.(*Name2ProfileReponse)
		//======
		time.Sleep(time.Millisecond * 2)
		checkPayNodeX := &PhotonNode{
			Host:       "http://" + params.PhotonHost,
			Address:    params.PhotonAddress,
			APIAddress: params.PhotonHost,
			DebugCrash: false,
		}
		channelX, err := checkPayNodeX.GetChannelWith(&PhotonNode{
			Address: user.EthAddress,
		}, params.TokenAddress)
		if err != nil {
			fmt.Println(fmt.Errorf(PrintTime()+"[getStatusChan] GetChannelWith err= %s", err))
			continue
		}
		if channelX == nil {
			_, ok := ChannelDepositTimes[user.EthAddress]
			if !ok {
				ChannelDepositTimes[user.EthAddress] = 1
			} else {
				ChannelDepositTimes[user.EthAddress] += 1
			}
			var retrytimes = ChannelDepositTimes[user.EthAddress]

			fmt.Println(fmt.Sprintf("retrytimes=%d", retrytimes))
			if retrytimes > 5 {
				continue
			}
			err = checkPayNodeX.OpenChannel(user.EthAddress, params.TokenAddress, new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(int64(params.MinBalanceInchannel))), params.SettleTime)
			if err != nil {
				fmt.Println(fmt.Errorf(PrintTime()+"[getStatusChan] clientid= %s, eth-address= %s, OpenChannel err=%s", user.ID, user.EthAddress, err))
				continue
			}
			fmt.Println(fmt.Sprintf(PrintTime()+"[getStatusChan] clientid= %s, eth-address= %s, OpenChannel SUCCESS", user.ID, user.EthAddress))
			//处理0x04B41B05346C21F4f877BA30750028077fa3cFAA建不了通道但是发钱了的问题
		}
		//=========
		nodeS, err := checkPayNodeX.GetNodeStatus(user.EthAddress)
		if err != nil {
			fmt.Println(fmt.Sprintf(PrintTime()+"[getStatusChan] eth-address= %s, GetNodeStatus err=%s", user.EthAddress, err))
			continue
		}
		if nodeS.IsOnline {
			onlinepeer := PeerStatus{
				user.ID,
				user.EthAddress,
				"online",
			}
			onlineChan <- onlinepeer
		}

	}
	exitChan <- true
}

func initPeers(useri []*Name2ProfileReponse) {
	for _, u := range useri {
		nowTime := time.Now().UnixNano() / 1e6
		lat := u.LastactiveTime

		if u.EthAddress != "" && (nowTime-lat) <= time.Hour.Milliseconds()*24*90 {
			userChan <- u
		}
	}
	close(userChan)
}

// newClient creat a client link to ssb-server
func newClient(ctx *cli.Context) (*ssbClient.Client, error) {
	sockPath := ctx.String("unixsock")
	if sockPath != "" {
		client, err := ssbClient.NewUnix(sockPath, ssbClient.WithContext(longCtx))
		if err != nil {
			level.Debug(log).Log("client", "unix-path based init failed", "err", err)
			level.Info(log).Log("client", "Now try switching to TCP working mode and init it")
			return newTCPClient(ctx)
		}
		level.Info(log).Log("client", "connected", "method", "unix sock")
		return client, nil
	}

	// Assume TCP connection
	return newTCPClient(ctx)
}

// newTCPClient create tcp client to support remote applications
func newTCPClient(ctx *cli.Context) (*ssbClient.Client, error) {
	localKey, err := ssb.LoadKeyPair(ctx.String("key"))
	if err != nil {
		return nil, err
	}

	var remotePubKey = make(ed25519.PublicKey, ed25519.PublicKeySize)
	copy(remotePubKey, localKey.ID().PubKey())
	if rk := ctx.String("remoteKey"); rk != "" {
		rk = strings.TrimSuffix(rk, ".ed25519")
		rk = strings.TrimPrefix(rk, "@")
		rpk, err := base64.StdEncoding.DecodeString(rk)
		if err != nil {
			return nil, fmt.Errorf("Init: base64 decode of --remoteKey failed: %w", err)
		}
		copy(remotePubKey, rpk)
	}

	plainAddr, err := net.ResolveTCPAddr("tcp", ctx.String("addr"))
	if err != nil {
		return nil, fmt.Errorf("Init: failed to resolve TCP address: %w", err)
	}

	shsAddr := netwrap.WrapAddr(plainAddr, secretstream.Addr{PubKey: remotePubKey})
	client, err := ssbClient.NewTCP(localKey, shsAddr,
		ssbClient.WithSHSAppKey(ctx.String("shscap")),
		ssbClient.WithContext(longCtx))
	if err != nil {
		return nil, fmt.Errorf("Init: failed to connect to %s: %w", shsAddr.String(), err)
	}

	fmt.Println(fmt.Sprintf(PrintTime()+"Client = [%s] , method = [%s] , linked pub server = [%s]", "connected", "TCP", shsAddr.String()))
	//127.0.0.1:8008|@HZnU6wM+F17J0RSLXP05x3Lag2jGv3F3LzHMjh72coE=.ed25519
	params.PubID = strings.Split(shsAddr.String(), "|")[1]
	fmt.Println(fmt.Sprintf(PrintTime()+"Init: success to work on pub [%s]", params.PubID))

	return client, nil
}

// initDb
func initDb(ctx *cli.Context) error {
	pubdatadir := ctx.String("datadir")

	likedb, err := OpenPubDB(pubdatadir)
	if err != nil {
		fmt.Println(fmt.Errorf("Failed to create database", err))
	}

	lstime, err := likedb.SelectLastScanTime()
	if err != nil {
		fmt.Println(fmt.Errorf("Failed to init database", err))
	}
	if lstime == 0 {
		_, err = likedb.UpdateLastScanTime(0)
		if err != nil {
			fmt.Println(fmt.Errorf("Failed to init database", err))
		}
	}
	lastAnalysisTimesnamp = lstime

	likeDB = likedb

	return nil
}

// DoMessageTask get message from the server copy
func DoMessageTask(ctx *cli.Context) {
	//init db
	if err := initDb(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	time.Sleep(time.Second * 1)

	//init sensitive words
	f, err := os.Open(params.SensitiveWordsFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pubNode = &PhotonNode{
		Host:       "http://" + params.PhotonHost,
		Address:    params.PhotonAddress,
		APIAddress: params.PhotonHost,
		DebugCrash: false,
	}
	checkChannelNode = &PhotonNode{
		Host:       "http://" + params.PhotonHost,
		Address:    params.PhotonAddress,
		APIAddress: params.PhotonHost,
		DebugCrash: false,
	}
	checkPayNode = &PhotonNode{
		Host:       "http://" + params.PhotonHost,
		Address:    params.PhotonAddress,
		APIAddress: params.PhotonHost,
		DebugCrash: false,
	}

	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		word := scanner.Text()
		SensitiveWords = append(SensitiveWords, word)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dfax = dfa.New()
	dfax.AddBadWords(SensitiveWords)

	time.Sleep(time.Second * 1)

	PhotonNodeStatus = &PhotonNodeStatusStu{
		nodesstatusinfo: make(map[string]*ClientOnlineStatus),
	}
	//ssb-message work
	for {
		//构建符合条件的message请求
		var ref refs.FeedRef
		if id := ctx.String("id"); id != "" {
			var err error
			ref, err = refs.ParseFeedRef(id)
			if err != nil {
				panic(err)
			}
		}
		args := message.CreateHistArgs{
			ID:     ref,
			Seq:    ctx.Int64("seq"),
			AsJSON: ctx.Bool("asJSON"),
		}
		args.Gt = message.RoundedInteger(lastAnalysisTimesnamp - 12*3600*1000)
		args.Limit = -1
		args.Seq = 0
		args.Keys = true
		args.Values = true
		args.Private = false
		src, err := client.Source(longCtx, muxrpc.TypeJSON, muxrpc.Method{"createLogStream"}, args)
		if err != nil {
			//client可能失效,则需要重建新的连接,链接资源的释放在ssb-server端
			fmt.Println(fmt.Errorf(PrintTime()+"Source stream call failed: %w ,will try other tcp connect socket...", err))
			otherClient, err := newClient(ctx)
			if err != nil {
				fmt.Println(fmt.Errorf(PrintTime()+"Try set up a ssb client tcp socket failed , will try again...", err))
				time.Sleep(time.Second * 10)
				continue
			}

			client = otherClient
			continue
		}

		//从上一次的计算点（数据库记录的毫秒时间戳）到最后一条记录的解析
		time.Sleep(time.Second)
		calcComplateTime, err := SsbMessageAnalysis(src)
		if err != nil {
			fmt.Println(fmt.Sprintf(PrintTime()+"Message pump failed: %w", err))
			time.Sleep(time.Second * 5)
			continue
		}

		var calcsumthisTurn = len(TempMsgMap)
		fmt.Println(fmt.Sprintf(PrintTime()+"A round of message data analysis has been completed ,from TimeSanmp [%v] to [%v] ,message number = [%d]", lastAnalysisTimesnamp, calcComplateTime, calcsumthisTurn))
		lastAnalysisTimesnamp = calcComplateTime

		time.Sleep(params.MsgScanInterval)
	}
}

func contactSomeone(ctx *cli.Context, dealwho string, isfollow, isblock bool) (err error) {
	if dealwho == params.PubID {
		return fmt.Errorf("Permission denied, from pub : %s", dealwho)
	}
	/*arg := map[string]interface{}{
		"contact":   dealwho,
		"type":      "contact",
		"following": isfollow,
		"blocking":  isblock,
	}
	var v string
	err = client.Async(longCtx, &v, muxrpc.TypeString, muxrpc.Method{"publish"}, arg)
	if err != nil {
		return fmt.Errorf("publish call failed: %w", err)
	}
	//newMsg, err := refs.ParseMessageRef(v)
	//if err != nil {
	//	return err
	//}
	//log.Log("event", "published", "type", "contact", "ref", newMsg.String())
	*/
	// 使用QueryEscape进行URI编码
	encodedStringBlockWho := url.QueryEscape(dealwho)
	req := &Req{
		FullURL: fmt.Sprintf("http://127.0.0.1:8888/friends/block/%s", encodedStringBlockWho),
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	_, err = req.Invoke()
	fmt.Println(fmt.Sprintf("block %s,isfollow=%s, isblock=%s, err=%s", dealwho, isfollow, isblock, err))
	return
}

func privatePublish(ctx *cli.Context, recpobj, root, branch string) (err error) {
	arg := map[string]interface{}{
		"text": ctx.Args().First(),
		"type": "post",
	}
	if r := ctx.String("root"); r != "" {
		arg["root"] = r
		if b := ctx.String("branch"); b != "" {
			arg["branch"] = b
		} else {
			arg["branch"] = r
		}
	}
	var v string
	if recps := ctx.StringSlice("recps"); len(recps) > 0 {
		err = client.Async(longCtx, &v,
			muxrpc.TypeString,
			muxrpc.Method{"private", "publish"}, arg, recps)
	} else {
		err = client.Async(longCtx, &v,
			muxrpc.TypeString,
			muxrpc.Method{"publish"}, arg)
	}
	if err != nil {
		return fmt.Errorf("publish call failed: %w", err)
	}
	return
}

func SsbMessageAnalysis(r *muxrpc.ByteSource) (int64, error) {
	var buf = &bytes.Buffer{}
	TempMsgMap = make(map[string]*TempdMessage)
	ClientID2Name = make(map[string]string)

	LikeDetail = []string{}
	UnLikeDetail = []string{}

	//不能以最后一条消息的时间作为本轮计算的时间点,后期改为从服务器上取得pub的时间,
	//计算周期越小越好,加载完本轮所有消息的时间点即为下一轮的开始时间，这样规避了在计算过程中有新消息被同步进入pub
	//注意：manyvse等客户端向服务器同步数据，延迟时间不定，如果无网状态发送过来的消息被视为空
	nowUnixTime := time.Now().UnixNano() / 1e6

	for r.Next(context.TODO()) {
		//在本轮for计算周期内如果有数据
		buf.Reset()
		err := r.Reader(func(r io.Reader) error {
			_, err := buf.ReadFrom(r)
			return err
		})
		if err != nil {
			return 0, err
		}

		var msgStruct DeserializedMessageStu
		err = json.Unmarshal(buf.Bytes(), &msgStruct)
		if err != nil {
			continue
			//fmt.Println(fmt.Errorf("Muxrpc.ByteSource Unmarshal to json err =%s", err))
			//return 0, err
		}

		//1、记录本轮所有消息ID和author的关系,保存下来,被点赞的消息基本不会在本轮被扫描到
		msgkey := fmt.Sprintf("%v", msgStruct.Key)
		msgauther := fmt.Sprintf("%v", msgStruct.Value.Author)
		var msgtime = msgStruct.Value.Timestamp
		var msgTime = int64(msgtime*math.Pow10(2)) / 100
		if IsBlackList(msgauther) {
			continue
		}

		if msgTime < lastAnalysisTimesnamp-24*3600*1000 { //专为planetary pub，因为获取不到规定内的消息，24小时外存在没获取的消息，不管了
			continue
		}

		TempMsgMap[msgkey] = &TempdMessage{
			Author: msgauther,
		}
		/*_, err = likeDB.InsertLikeDetail(msgkey, msgauther)
		if err != nil {
			fmt.Println(fmt.Errorf(PrintTime()+"Failed to InsertLikeDetail, err=%s", err))
			return 0, err
		}*/

		//2、记录like的统计结果
		contentJust := string(msgStruct.Value.Content[0])
		if contentJust == "{" {
			//time.Sleep(time.Millisecond * 100)
			//1、like的信息
			cvs := ContentVoteStru{}
			err = json.Unmarshal(msgStruct.Value.Content, &cvs)
			if err == nil {
				if string(cvs.Type) == "vote" {
					//get the Unlike tag ,先记录被like的link，再找author；由于图谱深度不一样，按照时间顺序查询存在问题，则先统一记录
					//timesp := time.Unix(int64(msgStruct.Value.Timestamp)/1e3, 0).Format("2006-01-02 15:04:05")
					if cvs.Vote == nil {
						continue
					}
					if cvs.Vote.Expression == "Unlike" {
						UnLikeDetail = append(UnLikeDetail, cvs.Vote.Link)
						//fmt.Println(PrintTime() + "unlike-time: " + timesp + "---MessageKey: " + cvs.Vote.Link)

						/*//统计我取消点赞的
						_, err = likeDB.InsertUserSetLikeInfo(msgkey, msgauther, -1, msgTime)
						if err != nil {
							fmt.Println(fmt.Errorf(PrintTime()+"%s set a unlike FAILED, err=%s", msgauther, err))
						}
						fmt.Println(fmt.Sprintf(PrintTime()+"%s set a unlike, msgkey=%s", msgauther, msgkey))*/
					} else {
						//get the Like tag ,因为like肯定在发布message后,先记录被like的link，再找author
						LikeDetail = append(LikeDetail, cvs.Vote.Link)
						//fmt.Println(PrintTime() + "like-time: " + timesp + "---MessageKey: " + cvs.Vote.Link)

						/*//统计我点赞的
						_, err = likeDB.InsertUserSetLikeInfo(msgkey, msgauther, 1, msgTime)
						if err != nil {
							fmt.Println(fmt.Errorf(PrintTime()+"%s set a like FAILED, err=%s", msgauther, err))
						}
						fmt.Println(fmt.Sprintf(PrintTime()+"%s set a like, msgkey=%s", msgauther, msgkey))*/

						{ //发送激励
							//如果点赞了，又取消了，不影响token的发放
							name2addr, err := GetNodeProfile(msgauther)
							if err != nil || len(name2addr) != 1 {
								fmt.Println(fmt.Errorf(LikePost+" Reward %s ethereum address failed, err= not found or %s", msgauther, err))
							} else {
								ehtAddr := name2addr[0].EthAddress
								PubRewardToken(ehtAddr, int64(params.RewardOfLikePost), msgauther, LikePost, msgkey, msgTime)
							}
						}
					}
				}
			} else {
				/*fmt.Println(fmt.Sprintf("Unmarshal  for vote , err %v", err))*/
				//todox 可以根据协议的扩展，记录其他的vote数据，目前没有这个需求
			}

			//3、about即修改备注名为hex-address的信息,注意:修改N次name,只需要返回最新的即可
			//此为备份方案：认定Name为ethaddr,需要同步修改API，name字段代替other1
			cau := ContentAboutStru{}
			err = json.Unmarshal(msgStruct.Value.Content, &cau)
			if err == nil {
				if cau.Type == "about" {
					ClientID2Name[fmt.Sprintf("%v", cau.About)] =
						fmt.Sprintf("%v", cau.Name)
				}
			} else {
				fmt.Println(fmt.Errorf(PrintTime()+"Unmarshal for about , err %v", err))
			}

			//4、contact触发对blakclist的处理, 通过pub关注重新进来的黑名单的消息来持续block该账户
			if msgauther == params.PubID {
				ccs := ContentContactStru{}
				err = json.Unmarshal(msgStruct.Value.Content, &ccs)
				if err == nil {
					if ccs.Type == "contact" {
						if IsBlackList(ccs.Contact) && ccs.Following && ccs.Pub {
							//block he
							err = contactSomeone(nil, ccs.Contact, false, true)
							if err != nil {
								fmt.Println(fmt.Errorf(PrintTime()+"[black-list] Unfollow and Block %s FAILED, err=%s", ccs.Contact, err))
							}
							fmt.Println(fmt.Sprintf(PrintTime()+"[black-list] Unfollow and Block %s SUCCESS", ccs.Contact))
						}
					}
				} else {
					fmt.Println(fmt.Errorf(PrintTime()+"[black-list] Unmarshal for contact, err %v", err))
				}
			}

			//5、POST Message
			cps := ContentPostStru{}
			err = json.Unmarshal(msgStruct.Value.Content, &cps)
			if err == nil {
				if cps.Type == "post" {
					postContent := cps.Text
					//5.1敏感词处理
					/*_, _, b := dfax.Check(postContent)
					if b && (msgauther != params.PubID) {
						//block he
						//err = contactSomeone(nil, msgauther, true, true)
						//if err != nil {
						//	fmt.Println(fmt.Sprintf(PrintTime()+"[sensitive-check]Unfollow and Block %s FAILED, err=%s", msgauther, err))
						//}
						//fmt.Println(fmt.Sprintf(PrintTime()+"[sensitive-check]Unfollow and Block %s SUCCESS", msgauther))
						//fix:处理违规消息由 "直接block" 转为 "提供接口人工审核处理"
						_, err = likeDB.InsertSensitiveWordRecord(params.PubID, nowUnixTime, postContent, msgkey, msgauther, "0")
						if err != nil {
							fmt.Println(fmt.Errorf(PrintTime()+"[sensitive-check]InsertSensitiveWordRecord FAILED, err=%s", err))
						}
						fmt.Println(fmt.Sprintf(PrintTime()+"[sensitive-check]InsertSensitiveWordRecord SUCCESS, author=%s, message=%s, msgkey=%s", msgauther, "see...", msgkey))
					}*/
					//5.2我发表的post
					if cps.Root == "" && PostWordCountBigThan10(postContent) { //1-登录 2-发表帖子 3-评论 4-铸造NFT
						_, err = likeDB.InsertUserTaskCollect(params.PubID, msgauther, msgkey, "2", "", msgTime, "", "", "")
						if err != nil {
							fmt.Println(fmt.Errorf(PrintTime()+"[UserTaskCollect-Post] FAILED, err=%s", err))
						}
						//fmt.Println(fmt.Sprintf(PrintTime()+"[UserTaskCollect-Post] SUCCESS, author=%s, msgkey=%s", msgauther, msgkey))

						{ //发送激励
							name2addr, err := GetNodeProfile(msgauther)
							if err != nil || len(name2addr) != 1 {
								//fmt.Println(fmt.Errorf(PostMessage+" Reward %s ethereum address failed, err= not found or %s", msgauther, err))
							} else {
								ehtAddr := name2addr[0].EthAddress
								PubRewardToken(ehtAddr, int64(params.RewardOfPostMessage), msgauther, PostMessage, msgkey, msgTime)
							}
						}
					}
					//5.3我发表的comment
					if cps.Root != "" && PostWordCountBigThan10(postContent) {
						_, err = likeDB.InsertUserTaskCollect(params.PubID, msgauther, msgkey, "3", cps.Root, msgTime, "", "", "")
						if err != nil {
							fmt.Println(fmt.Errorf(PrintTime()+"[UserTaskCollect-Comment] FAILED, err=%s", err))
						}
						//fmt.Println(fmt.Sprintf(PrintTime()+"[UserTaskCollect-Comment] SUCCESS, author=%s, msgkey=%s", msgauther, msgkey))

						{ //发送激励
							name2addr, err := GetNodeProfile(msgauther)
							if err != nil || len(name2addr) != 1 {
								fmt.Println(fmt.Errorf(PostComment+" Reward %s ethereum address failed, err= not found or %s", msgauther, err))
							} else {
								ehtAddr := name2addr[0].EthAddress
								PubRewardToken(ehtAddr, int64(params.RewardOfPostComment), msgauther, PostComment, msgkey, msgTime)
							}
						}
					}
				}
			} else {
				//fmt.Println(fmt.Errorf("json.Unmarshal(msgStruct.Value.Content err=%s", err))
			}
		}

	}

	/*//save message-result to database
	for _, likeLink := range LikeDetail { //被点赞的ID集合,标记被点赞的记录
		_, err := likeDB.UpdateLikeDetail(1, nowUnixTime, likeLink)
		if err != nil {
			fmt.Println(fmt.Errorf(PrintTime()+"Failed to UpdateLikeDetail", err))
			return 0, err
		}
	}

	for _, unLikeLink := range UnLikeDetail { //被取消点赞的ID集合
		_, err := likeDB.UpdateLikeDetail(-1, nowUnixTime, unLikeLink)
		if err != nil {
			fmt.Println(fmt.Errorf(PrintTime()+"Failed to UpdateLikeDetail", err))
			return 0, err
		}
	}*/

	_, err := likeDB.UpdateLastScanTime(nowUnixTime)
	if err != nil {
		fmt.Println(fmt.Errorf(PrintTime()+"Failed to UpdateLastScanTime", err))
		return 0, err
	}
	//更新table userethaddr
	for key := range ClientID2Name {
		_, err := likeDB.UpdateUserProfile(1, key, ClientID2Name[key], "", "", "")
		if err != nil {
			fmt.Println(fmt.Errorf(PrintTime()+"Failed to UpdateUserEthAddr", err))
			//return 0, err
		}
	}
	//fmt.Println(fmt.Sprintf(PrintTime()+"A round of message data analysis has been completed ,message number = [%v]", len(TempMsgMap)))
	/*//print for test
	for key,value := range TempMsgMap {
		fmt.Println(key, "<-this round message ID---ClientID->", value.Author)
	}
	for key := range ClientID2Name { //取map中的值err
		fmt.Println(key, "<-ClientID---Name->", ClientID2Name[key])
	}*/
	return nowUnixTime, nil
}

// NewChannelDeal
func NewChannelDeal(partnerAddress string, clientID string, messageTime int64, useInviteCode bool) (err error) {
	partnerNode := &PhotonNode{
		//:utils.APex2(rs.Config.PubAddress),
		Address:    partnerAddress,
		DebugCrash: false,
	}

	channel00, err := pubNode.GetChannelWith(partnerNode, params.TokenAddress)
	if err != nil {
		fmt.Println(fmt.Errorf(PrintTime()+SignUp+" GetChannelWith %s", err))
		return
	}
	if channel00 == nil {
		if ExceedRewardLimit(clientID, SignUp, messageTime, 0) {
			//如果一个SSB-ID连续注册地址达到2次以上，则该账号以后无法得到注册激励
			fmt.Println(fmt.Errorf(PrintTime()+SignUp+" reward %s to ethaddr=%s REJECT,reason:ExceedRewardLimit or ssbid error", clientID, partnerAddress))
			return
		}
		//create new channel with  mlt
		initRegistAmount := int64(params.MinBalanceInchannel + params.RewardOfSignup)
		err = pubNode.OpenChannel(partnerNode.Address, params.TokenAddress, new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(initRegistAmount)), params.SettleTime)
		if err != nil {
			fmt.Println(fmt.Errorf(PrintTime()+SignUp+" create channel with %s, err=%s", partnerAddress, err))
		} else {
			fmt.Println(fmt.Sprintf(PrintTime()+SignUp+" create channel success[%s], with %s", clientID, partnerAddress))
		}

		netStatus := false
		for i := 0; i < 10; i++ {
			nodeS, err := pubNode.GetNodeStatus(partnerAddress)
			if err != nil {
				fmt.Println(fmt.Errorf(PrintTime()+SignUp+" GetNodeStatus[%s], err=%s", clientID, err))
			}
			netStatus = nodeS.IsOnline
			if netStatus {
				break
			}
			time.Sleep(time.Second * 1)
		}
		if !netStatus {
			fmt.Println(fmt.Errorf(PrintTime()+SignUp+" GetNodeStatus[%s](retry 10 times) %s online=%v", clientID, partnerAddress, netStatus))
			//如果此时客户端不在线，则先记录，后续补发
			{
				//=======Record Reward Result=======
				var grantTokens int64
				if useInviteCode {
					grantTokens = int64(params.RewardOfSignup + 5)
				} else {
					grantTokens = int64(params.RewardOfSignup)
				}
				_, err = likeDB.RecordRewardResult(clientID, partnerAddress, "fail", grantTokens, SignUp, "", messageTime, messageTime)
				fmt.Println(fmt.Sprintf(PrintTime()+SignUp+" but offline ,then[RecordRewardResult] eth-address=%s for clientid=%s, reason=%s, err=%s", partnerAddress, clientID, err))
			}
			return errors.New("partner offline")
		}
		//registration award 新地址才发送注册激励
		var amount *big.Int
		if useInviteCode {
			amount = new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(int64(params.RewardOfSignup+5)))
		} else {
			amount = new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(int64(params.RewardOfSignup)))
		}
		err = pubNode.SendTrans(params.TokenAddress, amount, partnerAddress, true, false)
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf(PrintTime()+SignUp+" award[%s] to %s, amount= %v, err= %v", clientID, partnerAddress, amount, err))

		//继续发送SMT激励
		smtAmount := new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(int64(params.RewardOfSignupSMT)))
		//smtAmount = new(big.Int).Div(smtAmount, big.NewInt(10))
		err = pubNode.TransferSMT(partnerAddress, smtAmount.String())
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf(PrintTime()+SignUp+" award(SMT)[%s] to %s, amount=%v, err=%v", clientID, partnerAddress, smtAmount, err))

		{
			//=======Record Reward Result=======
			nowTime := time.Now().UnixNano() / 1e6
			var grantTokens int64
			if useInviteCode {
				grantTokens = int64(params.RewardOfSignup + 5)
			} else {
				grantTokens = int64(params.RewardOfSignup)
			}
			_, err = likeDB.RecordRewardResult(clientID, partnerAddress, "success", grantTokens, SignUp, "", messageTime, nowTime)
			fmt.Println("..........")
			fmt.Println(useInviteCode)
			fmt.Println(fmt.Sprintf(PrintTime()+SignUp+"[RecordRewardResult] reword to eth-address=%s for clientid=%s, reason=%s, err=%v", partnerAddress, clientID, SignUp, err))
		}

	} else {
		fmt.Println(fmt.Errorf(PrintTime()+"[Pub-Client-ChannelDeal-OK]channel has exist[%s], with %s", clientID, partnerAddress))
	}

	return
}

// PubRewardToken  pub paid additionally
// It is stipulated that 'the award' needs to be paid additionally by pub, and the 'min-balance-inchannel' is not used
func PubRewardToken(partnerAddress string, xamount int64, clientID, reason, messageKey string, messageTime int64) (err error) {
	_, err = HexToAddress(partnerAddress)
	if err != nil {
		err = fmt.Errorf("[PubRewardToken]verify eth-address= %s, error= %s", partnerAddress, err)
		return
	}
	time.Sleep(time.Millisecond * 50)
	//===========================
	channelX, err := pubNode.GetChannelWith(&PhotonNode{
		Address: partnerAddress,
	}, params.TokenAddress)
	if err != nil {
		fmt.Println(fmt.Errorf(PrintTime()+"[PubRewardToken] GetChannelWith err= %s", err))
		return
	}
	if channelX == nil {
		err = pubNode.OpenChannel(partnerAddress, params.TokenAddress, new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(int64(params.MinBalanceInchannel))), params.SettleTime)
		if err != nil {
			fmt.Println(fmt.Errorf(PrintTime()+"[PubRewardToken] clientid= %s, eth-address= %s, OpenChannel err=%s", clientID, partnerAddress, err))
			return
		}
		fmt.Println(fmt.Sprintf(PrintTime()+"[PubRewardToken] clientid= %s, eth-address= %s, OpenChannel SUCCESS", clientID, partnerAddress))
	}
	//===========================
	netStatus := false
	for i := 0; i < 3; i++ {
		nodeS, err := pubNode.GetNodeStatus(partnerAddress)
		if err != nil {
			fmt.Println(fmt.Sprintf(PrintTime()+" [PubRewardToken] GetNodeStatus eth-address= %s, err= %s", partnerAddress, err))
			break
		}
		netStatus = nodeS.IsOnline
		if netStatus {
			break
		}
		//time.Sleep(time.Second * 1)
	}
	if ExceedRewardLimit(clientID, reason, messageTime, xamount) {
		//fmt.Println(fmt.Errorf(PrintTime()+" [PubRewardToken] clientid= %s, ethaddr= %s, REJECT, reason: ExceedRewardLimit", clientID, partnerAddress))
		return
	}
	if err = PassRule202302(clientID, reason, messageTime, xamount); err != nil {
		//fmt.Println(fmt.Errorf(PrintTime()+" [PubRewardToken] clientid= %s, ethaddr= %s, REJECT, reason: %v", clientID, partnerAddress, err))
		return
	}
	amount := new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(xamount))
	if !netStatus {
		//如果此时客户端不在线，则先记录，后续补发
		{
			_, err = likeDB.RecordRewardResult(clientID, partnerAddress, "fail", xamount, reason, messageKey, messageTime, messageTime)
			//fmt.Println(fmt.Sprintf(PrintTime()+"[PubRewardToken] FAILED because OffLine, clientid= %s, eth-address=%s, reason= %s, messagekey= %s, RecordRewardResult= %s", clientID, partnerAddress, reason, messageKey, err))
		}
		return errors.New("partner offline")
	}
	nowTime := time.Now().UnixNano() / 1e6
	_, err = likeDB.RecordRewardResult(clientID, partnerAddress, "success", xamount, reason, messageKey, messageTime, nowTime)
	if err != nil {
		//fmt.Println(fmt.Sprintf(PrintTime()+"[PubRewardToken] clientid= %s, eth-address= %s, reason= %s, FAILED, RecordRewardResult= %s", clientID, partnerAddress, reason, err))
		return err
	} else {
		err = pubNode.SendTrans(params.TokenAddress, amount, partnerAddress, true, false)
		if err != nil {
			fmt.Println(fmt.Errorf(PrintTime()+"[PubRewardToken] clientid= %s, eth-address= %s, reason= %s, SendTrans= %s", clientID, partnerAddress, reason, err))
			//当前资金不足，等待补充后补发或再目前通道不存在
			//if err.Error() == "InsufficientBalance" {
			_, err = likeDB.UpdateRewardResult(clientID, partnerAddress, "fail", messageTime, messageTime)
			fmt.Println(fmt.Sprintf(PrintTime()+"[PubRewardToken] FAILED because %s, clientid= %s, eth-address=%s, reason= %s, messagekey= %s, UpdateRewardResult= %s", err, clientID, partnerAddress, reason, messageKey, err))
			//}
			return err
		}
	}
	fmt.Println(fmt.Sprintf(PrintTime()+"[PubRewardToken] clientid= %s, eth-address= %s, reason=%s, SUCCESS", clientID, partnerAddress, reason))
	return
}

func PubCheckBalance() {
	for {
		fmt.Println(fmt.Sprintf(PrintTime() + "Pub Check Channel Balance...."))
		time.Sleep(time.Second * 3) //数据库可能没准备好

		channelAll, err := checkChannelNode.GetChannels(params.TokenAddress)
		if err != nil {
			fmt.Println(fmt.Errorf(PrintTime()+"[PubCheckBalance] GetChannels err=%s", err))
			continue
		}

		//实时核对，计算量过大，改为仅查通道，不看注册
		for i := 0; i < len(channelAll); i++ {
			partnerAddr := channelAll[i].PartnerAddress
			channelState := channelAll[i].State
			channelToken := channelAll[i].TokenAddress
			if channelState == 1 && channelToken == params.TokenAddress {
				var minNum = new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(int64(params.MinBalanceInchannel)))
				var nowNum = channelAll[i].Balance
				var diffNum = new(big.Int).Sub(minNum, nowNum)
				if minNum.Cmp(nowNum) == 1 {
					//补充至MinBalanceInchannel
					err0 := checkChannelNode.Deposit(partnerAddr, params.TokenAddress, diffNum, 48)
					if err0 != nil {
						fmt.Println(fmt.Errorf(PrintTime()+"[PubCheckBalance] eth-address= %s, Deposit err= %s", partnerAddr, err0))
						continue
					}
					fmt.Println(fmt.Sprintf(PrintTime()+"[PubCheckBalance] eth-address= %s, num= %v, Deposit SUCCESS", partnerAddr, diffNum))
				}
			}
			if channelAll[i].State == 2 { //附属功能，对已经关闭的通道，且到达结算时间的进行结算，如果有下一次的激励行为，会主动重开通道
				if channelAll[i].BlockNumberNow >= channelAll[i].BlockNumberChannelCanSettle {
					err = checkChannelNode.Settle(channelAll[i].ChannelIdentifier, params.SettleTime)
					if err != nil {
						fmt.Println(fmt.Errorf(PrintTime()+"[PubCheckBalance]  eth-address= %s, Settle err=%s", channelAll[i].PartnerAddress, err))
						continue
					}
					fmt.Println(fmt.Errorf(PrintTime()+"[PubCheckBalance] eth-address= %s, Settle Channel SUCCESS", channelAll[i].PartnerAddress))
				}
			}
			time.Sleep(time.Millisecond * 50)
		}
		//time.Sleep(time.Second * 5)
	}
	//time.AfterFunc(params.RoundTimeOfCheckChannelBalance, checkPubChannelBalance)
}

var backpackObj = ""

/*// PubBackPay
func PubBackPay() {

	for {
		time.Sleep(time.Second * 35) //数据库可能没准备好
		fmt.Println(fmt.Sprintf(PrintTime() + "Pub Check Back Pay Node ...."))

		offinfos, err := likeDB.QueryGrantFail("")
		if err != nil {
			fmt.Println(fmt.Errorf(PrintTime()+"[PubBackPay] SelectRewardResult err= %s", err))
			continue
		}
		for _, info := range offinfos {
			partnerSsbid := info.ClientID
			partnerAddress := info.ClientEthAddress
			//1、没有通道的重开通道;
			//2、通道被关的(closed),先结算后重开通道
			if partnerSsbid != backpackObj {
				fmt.Println(PrintTime() + " [PubBackPay] for clientid=" + partnerSsbid + ", eth-address=" + partnerAddress)
				backpackObj = partnerSsbid
			}
			channelX, err := checkPayNode.GetChannelWith(&PhotonNode{
				Address: partnerAddress,
			}, params.TokenAddress)
			if err != nil {
				fmt.Println(fmt.Errorf(PrintTime()+"[PubBackPay] GetChannelWith err= %s", err))
				continue
			}
			if channelX == nil {
				err = checkPayNode.OpenChannel(partnerAddress, params.TokenAddress, new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(int64(params.MinBalanceInchannel))), params.SettleTime)
				if err != nil {
					fmt.Println(fmt.Errorf(PrintTime()+"[PubBackPay] clientid= %s, eth-address= %s, OpenChannel err=%s", partnerSsbid, partnerAddress, err))
					continue
				}
				fmt.Println(fmt.Sprintf(PrintTime()+"[PubBackPay] clientid= %s, eth-address= %s, OpenChannel SUCCESS", partnerSsbid, partnerAddress))
			}
			//------------------------
			nodeS, err := checkPayNode.GetNodeStatus(partnerAddress)
			if err != nil {
				fmt.Println(fmt.Sprintf(PrintTime()+"[PubBackPay] clientid= %s, eth-address= %s, GetNodeStatuserr=%s", partnerSsbid, partnerAddress, err))
				continue
			}

			if nodeS.IsOnline {
				fmt.Println(fmt.Sprintf(PrintTime()+"[PubBackPay] clientid= %s, eth-address= %s, online= %v", partnerSsbid, partnerAddress, nodeS.IsOnline))
				//补发
				msgTime := info.MessageTime
				amount := info.GrantTokenAmount
				reason := info.RewardReason
				_, err := likeDB.UpdateRewardResult(partnerSsbid, partnerAddress, "success", msgTime, time.Now().UnixNano()/1e6)
				if err != nil {
					fmt.Println(fmt.Errorf(PrintTime()+"[PubBackPay] clientid= %s, eth-address= %s, UpdateRewardResult err= %s", partnerSsbid, partnerAddress, err))
					continue
				}
				err = checkPayNode.SendTrans(params.TokenAddress, amount, partnerAddress, true, false)
				if err != nil {
					_, err1 := likeDB.UpdateRewardResult(partnerSsbid, partnerAddress, "fail", msgTime, msgTime)
					fmt.Println(fmt.Errorf(PrintTime()+"[PubBackPay] clientid= %s, eth-address= %s, SendTrans err= %s, UpdateRewardResult= %s", partnerSsbid, partnerAddress, err, err1))
					continue
				}
				//对sign up 补发SMT激励
				if reason == SignUp {
					smtAmount := new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(int64(params.RewardOfSignupSMT)))
					//smtAmount = new(big.Int).Div(smtAmount, big.NewInt(10))
					err = checkPayNode.TransferSMT(partnerAddress, smtAmount.String())
					if err != nil {
						continue
					}
					fmt.Println(fmt.Sprintf(PrintTime()+"[PubBackPay] (SMT) clientid= %s, eth-address= %s, amount= %v, err= %v", partnerSsbid, partnerAddress, amount, err))
				}
				fmt.Println(fmt.Sprintf(PrintTime()+"[PubBackPay] clientid= %s, eth-address= %s, SUCCESS", partnerSsbid, partnerAddress))

			}
			time.Sleep(time.Millisecond * 15)
		}
		fmt.Println(fmt.Sprintf(PrintTime()+"[PubBackPay] Pub check node's backpay info,coming in....%d", len(offinfos)))
	}
}*/

func IsBlackList(defendant string) bool {
	blacklists, err := likeDB.SelectViolationByWhere("", defendant, "", "", "1")
	if err != nil {
		fmt.Println(fmt.Errorf(PrintTime()+"selectBlacklist-Failed to get blacklist, err=%s", err))
		return false
	}
	if len(blacklists) > 0 {
		return true
	}
	return false
}

/*// dealBlacklist
func dealBlacklist() {
	for {
		time.Sleep(time.Second * 600)
		//get blacklist info
		blacklists, err := likeDB.SelectViolationByWhere("", "", "", "", "1")
		if err != nil {
			fmt.Println(fmt.Sprintf(PrintTime()+"dealBlacklist-Failed to get blacklist, err=%s", err))
		}
		for _, info := range blacklists {
			dealObj := info.Defendant
			//block him
			err = contactSomeone(nil, dealObj, false, true)
			if err != nil {
				fmt.Println(fmt.Sprintf("dealBlacklist-Unfollow and block %s failed", dealObj))
			}
			fmt.Println(fmt.Sprintf(PrintTime()+"dealBlacklist-Success to Unfollow and Block %s", dealObj))
			time.Sleep(time.Second * 3)
			//award plaintiff
			plaintiff := info.Plaintiff
			dealReward := info.Dealreward
			if strings.Index(dealReward, "-") != -1 {
				//awards have been issued
			} else {
				// No awards have been issued yet, for some reason
				name2addr, err := GetNodeProfile(plaintiff)
				if err != nil {
					fmt.Println(fmt.Sprintf("dealBlacklist-Get plaintiff's profile failed, err=%s", err))
					continue
				}
				if len(name2addr) != 1 {
					continue
				}
				addrPlaintiff := name2addr[0].EthAddress

				//另行支付
				err = sendToken(addrPlaintiff, int64(params.ReportRewarding), true, false)
				if err != nil {
					fmt.Println(fmt.Sprintf(PrintTime()+"dealBlacklist-Failed to Award to %s for ReportRewarding, err=%s", plaintiff, err))
					continue
				}
				fmt.Println(fmt.Sprintf(PrintTime()+"dealBlacklist-Success to Award to %s for ReportRewarding", plaintiff))
				_, err = likeDB.UpdateViolation(info.DealTag, info.Dealtime, string(params.ReportRewarding)+"-", plaintiff, dealObj, info.MessageKey)
				if err != nil {
					fmt.Println(fmt.Sprintf(PrintTime()+"dealBlacklist-Failed to Update ReportRewarding to %s", plaintiff))
					continue
				}

			}

		}
	}
}*/

func createStreamOfSsb(ctx *cli.Context) (msgs []*DeserializedMessageStu, err error) {
	//构建符合条件的message请求
	var ref refs.FeedRef
	if id := ctx.String("id"); id != "" {
		var err error
		ref, err = refs.ParseFeedRef(id)
		if err != nil {
			panic(err)
		}
	}
	args := message.CreateHistArgs{
		ID:     ref,
		Seq:    ctx.Int64("seq"),
		AsJSON: ctx.Bool("asJSON"),
	}
	args.Gt = message.RoundedInteger(lastAnalysisTimesnamp)
	args.Limit = -1
	args.Seq = 0
	args.Keys = true
	args.Values = true
	args.Private = false
	src, err := client.Source(longCtx, muxrpc.TypeJSON, muxrpc.Method{"createLogStream"}, args)
	if err != nil {
		//client可能失效,则需要重建新的连接,链接资源的释放在ssb-server端
		fmt.Println(fmt.Errorf(PrintTime()+"Source stream call failed: %w ", err))
	}

	var buf = &bytes.Buffer{}
	for src.Next(context.TODO()) {
		//在本轮for计算周期内如果有数据
		buf.Reset()
		err := src.Reader(func(r io.Reader) error {
			_, err := buf.ReadFrom(r)
			return err
		})
		if err != nil {
			return nil, err
		}

		var msgStruct DeserializedMessageStu
		err = json.Unmarshal(buf.Bytes(), &msgStruct)
		if err != nil {
			fmt.Println(fmt.Errorf("Muxrpc.ByteSource Unmarshal to json err =%s", err))
			return nil, err
		}
		msgs = append(msgs, &msgStruct)
	}

	return msgs, nil
}
