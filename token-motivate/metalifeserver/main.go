// SPDX-FileCopyrightText: 2021 The Go-SSB Authors
//
// SPDX-License-Identifier: MIT

// sbotcli implements a simple tool to query commands on another sbot
package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/shurcooL/go-goon"
	"go.cryptoscope.co/muxrpc/v2"
	"go.cryptoscope.co/netwrap"
	"go.cryptoscope.co/secretstream"
	kitlog "go.mindeco.de/log"
	"go.mindeco.de/log/level"
	"go.mindeco.de/log/term"
	"golang.org/x/crypto/ed25519"
	"gopkg.in/urfave/cli.v2"

	"os/signal"
	"syscall"

	"go.cryptoscope.co/ssb"
	ssbClient "go.cryptoscope.co/ssb/client"
	"go.cryptoscope.co/ssb/plugins/legacyinvites"
	"go.cryptoscope.co/ssb/restful"
	"go.cryptoscope.co/ssb/restful/params"
	refs "go.mindeco.de/ssb-refs"
)

// Version and Build are set by ldflags
var (
	// Version version of this build
	Version = "snapshot"

	// Build build time of this build
	Build = ""

	// GoVersion go version of this build
	GoVersion = ""

	// GitCommit git commit of this build
	GitCommit = ""
)

var (
	longctx      context.Context
	shutdownFunc func()

	log kitlog.Logger

	keyFileFlag  = cli.StringFlag{Name: "key,k", Value: "unset"}
	unixSockFlag = cli.StringFlag{Name: "unixsock", Usage: "if set, unix socket is used instead of tcp"}
	dataDir      = cli.StringFlag{Name: "datadir", Usage: "directory for storing pub's parsing data"}

	sensitiveWordsFlag = cli.StringFlag{Name: "sensitive-words-file", Usage: "the path of the sensitive-words file"}
)

func init() {
	u, err := user.Current()
	check(err)

	keyFileFlag.Value = filepath.Join(u.HomeDir, ".ssb-go", "secret")
	//we use tcp instead of UNIX SOCKET
	//unixSockFlag.Value = filepath.Join(u.HomeDir, ".ssb-go", "socket")
	dataDir.Value = filepath.Join(u.HomeDir, ".ssb-go", "pubdata")
	sensitiveWordsFlag.Value = filepath.Join(u.HomeDir, ".ssb-go", "sensitive.txt")
	params.Ip2LocationLiteDbPath = filepath.Join(u.HomeDir, ".ssb-go", "IP2LOCATION-LITE-DB11.IPV6.BIN")
	params.GameUserFilePath = filepath.Join(u.HomeDir, ".ssb-go", "gamefile", "userdata")
	params.GameResourcePath = filepath.Join(u.HomeDir, ".ssb-go", "gamefile", "resource")

	log = term.NewColorLogger(os.Stderr, kitlog.NewLogfmtLogger, colorFn)
}

var app = cli.App{
	Name:    os.Args[0],
	Usage:   "Metalife's application services on SSB pub",
	Version: "beta1",

	Flags: []cli.Flag{
		&cli.StringFlag{Name: "shscap", Value: "1KHLiKZvAvjbY1ziZEHMXawbCEIM6qwjCDm3VYRan/s=", Usage: "shs key"},
		&cli.StringFlag{Name: "addr", Value: params.PubTcpHostAddress, Usage: "tcp address of the sbot to connect to (or listen on)"},
		&cli.StringFlag{Name: "remoteKey", Value: "", Usage: "the remote pubkey you are connecting to (by default the local key)"},
		&dataDir,
		&cli.StringFlag{Name: "token-address", Value: "0x6601F810eaF2fa749EEa10533Fd4CC23B8C791dc", Usage: "which token is used in metalife app,if set,the default will be replaced"},
		&cli.StringFlag{Name: "photon-host", Value: "127.0.0.1:11001", Usage: "host:port link to the photon service."},
		&cli.StringFlag{Name: "pub-eth-address", Usage: "ethereum address the pub 's address is bound for reward."},
		&cli.StringFlag{Name: "anonther-serve", Usage: "host:port link to the anonther pub."},
		&cli.IntFlag{Name: "settle-timeout", Value: 40000, Usage: "set settle timeout on photon."},
		&cli.IntFlag{Name: "service-port", Value: 10008, Usage: "port' for the metalife service to listen on."},
		&cli.IntFlag{Name: "message-scan-interval", Value: 60, Usage: "the time interval at which messages are scanned and calculated (unit:second)."},
		&cli.IntFlag{Name: "min-balance-inchannel", Value: 1, Usage: "minimum balance in photon channel between this pub and ssb client (unit: 1e18 wei)."},
		&cli.IntFlag{Name: "report-rewarding", Value: 0, Usage: "pub will reward the person who provides the report (if the report is true). (unit: 1e18 wei)"},
		&cli.IntFlag{Name: "registration-rewarding-mlt", Value: 0, Usage: "pub will reward the person who provides ethereum address for his ssb client. (unit: 1e18 wei)"},
		&cli.IntFlag{Name: "registration-rewarding-smt", Value: 0, Usage: "pub will reward the person who provides ethereum address for his ssb client. (unit: 1e18 wei)"},
		&cli.IntFlag{Name: "max-daily-rewarding", Value: 0, Usage: "max daily-rewarding from 00:00:00 to 23:59:59. (unit: 1e18 wei)"},
		&cli.BoolFlag{Name: "ip-check", Value: false, Usage: "It is illegal to register with the same ip address within 1 Hour"},

		&sensitiveWordsFlag,
		&keyFileFlag,
		&unixSockFlag,
		&cli.BoolFlag{Name: "verbose,vv", Usage: "print muxrpc packets"},

		&cli.StringFlag{Name: "timeout", Value: "3600s", Usage: "pass a duration (like 3s or 5m) after which it times out, empty string to disable"},
	},

	Before: initClient,
	Commands: []*cli.Command{
		aliasCmd,
		blobsCmd,
		blockCmd,
		friendsCmd,
		getCmd,
		inviteCmds,
		logStreamCmd,
		sortedStreamCmd,
		typeStreamCmd,
		historyStreamCmd,
		partialStreamCmd,
		replicateUptoCmd,
		repliesStreamCmd,
		callCmd,
		sourceCmd,
		connectCmd,
		publishCmd,
		groupsCmd,
	},
}

// Color by error type
func colorFn(keyvals ...interface{}) term.FgBgColor {
	for i := 1; i < len(keyvals); i += 2 {
		if _, ok := keyvals[i].(error); ok {
			return term.FgBgColor{Fg: term.Red}
		}
	}
	return term.FgBgColor{}
}

func check(err error) {
	if err != nil {
		level.Error(log).Log("err", err)
		os.Exit(1)
	}
}

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s (Rev: %s, Built: %s, GitCommit: %s, GoVersion: %s)\n", c.App.Version, Version, Build, GitCommit, GoVersion)
	}

	if err := app.Run(os.Args); err != nil {
		level.Error(log).Log("run-failure", err)
	}
}

func todo(ctx *cli.Context) error {
	return fmt.Errorf("todo: %s", ctx.Command.Name)
}

func initClient(ctx *cli.Context) error {
	//init usr config
	tokenaddressStr := ctx.String("token-address")
	if tokenaddressStr == "" {
		return fmt.Errorf("Program startup parameters [token-address] must be set")
	}
	params.TokenAddress = tokenaddressStr

	apihost, apiport, err := net.SplitHostPort(ctx.String("photon-host"))
	if err != nil {
		return err
	}
	params.PhotonHost = apihost + ":" + apiport

	pubethaddressStr := ctx.String("pub-eth-address")
	if len(pubethaddressStr) != 42 || pubethaddressStr[0:2] != "0x" {
		return fmt.Errorf("Program startup parameters [pub-eth-address] must be set")
	}
	params.PhotonAddress = pubethaddressStr

	anohost, anoport, err := net.SplitHostPort(ctx.String("anonther-serve"))
	if err != nil {
		return err
	}
	params.AnotherServe = anohost + ":" + anoport

	settletimeoutInt := ctx.Int("settle-timeout")
	if settletimeoutInt <= 0 {
		return fmt.Errorf("settle timeout should > 0")
	}
	params.SettleTime = settletimeoutInt

	serveportInt := ctx.Int("service-port")
	if serveportInt <= 0 {
		return fmt.Errorf("service-port %v error", serveportInt)
	}
	params.ServePort = serveportInt

	messagescanintervalInt := ctx.Int("message-scan-interval")
	if messagescanintervalInt <= 0 {
		return fmt.Errorf("service-port %v error", messagescanintervalInt)
	}
	params.MsgScanInterval = time.Second * time.Duration(messagescanintervalInt)

	minbalance := ctx.Int("min-balance-inchannel")
	if minbalance < 0 {
		return fmt.Errorf("min-balance-inchannel %v error", minbalance)
	}
	params.MinBalanceInchannel = minbalance

	reportrewarding := ctx.Int("report-rewarding")
	if reportrewarding < 0 {
		return fmt.Errorf("report-rewarding %v error", reportrewarding)
	}
	params.RewardOfReportProblematicPost = reportrewarding

	registrationawarding := ctx.Int("registration-rewarding-mlt")
	if registrationawarding < 0 {
		return fmt.Errorf("registration-rewarding-mlt %v error", registrationawarding)
	}
	params.RewardOfSignup = registrationawarding

	registrationawardSMT := ctx.Int("registration-rewarding-smt")
	if registrationawardSMT < 0 {
		return fmt.Errorf("registration-rewarding-smt %v error", registrationawardSMT)
	}
	params.RewardOfSignupSMT = registrationawardSMT

	sensitivewordsfilepath := ctx.String("sensitive-words-file")
	if sensitivewordsfilepath == "" {
		return fmt.Errorf("Program startup parameters [sensitive-words-file] must be set")
	}
	params.SensitiveWordsFilePath = sensitivewordsfilepath

	maxdailyrewarding := ctx.Int("max-daily-rewarding")
	if maxdailyrewarding < 0 {
		return fmt.Errorf("max-daily-rewarding", maxdailyrewarding)
	}
	params.MaxDailyRewarding = maxdailyrewarding

	params.CheckIP = ctx.Bool("ip-check")
	fmt.Printf("init check-ip=%v", params.CheckIP)

	dstr := ctx.String("timeout")
	if dstr != "" {
		d, err := time.ParseDuration(dstr)
		if err != nil {
			return err
		}
		longctx, shutdownFunc = context.WithTimeout(context.Background(), d)
	} else {
		longctx, shutdownFunc = context.WithCancel(context.Background())
	}

	signalc := make(chan os.Signal)
	signal.Notify(signalc, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-signalc
		level.Warn(log).Log("event", "shutting down", "sig", s)
		shutdownFunc()
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

	fmt.Println(fmt.Sprintf("\n******os.args=%q", os.Args))

	//start pub message analysis service
	restful.Start(ctx)

	return nil
}

func newClient(ctx *cli.Context) (*ssbClient.Client, error) {
	sockPath := ctx.String("unixsock")
	if sockPath != "" {
		client, err := ssbClient.NewUnix(sockPath, ssbClient.WithContext(longctx))
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
			return nil, fmt.Errorf("init: base64 decode of --remoteKey failed: %w", err)
		}
		copy(remotePubKey, rpk)
	}

	plainAddr, err := net.ResolveTCPAddr("tcp", ctx.String("addr"))
	if err != nil {
		return nil, fmt.Errorf("int: failed to resolve TCP address: %w", err)
	}

	shsAddr := netwrap.WrapAddr(plainAddr, secretstream.Addr{PubKey: remotePubKey})
	client, err := ssbClient.NewTCP(localKey, shsAddr,
		ssbClient.WithSHSAppKey(ctx.String("shscap")),
		ssbClient.WithContext(longctx))
	if err != nil {
		return nil, fmt.Errorf("init: failed to connect to %s: %w", shsAddr.String(), err)
	}
	//level.Info(log).Log("client", "connected", "method", "tcp")
	level.Info(log).Log("client", "connected", "method", "tcp", "Working Mode", "tcp", "Pub Server", shsAddr.String())

	return client, nil
}

var callCmd = &cli.Command{
	Name:  "call",
	Usage: "make an dump* async call",
	UsageText: `SUPPORTS:
* whoami
* latestSequence
* getLatest
* get
* blobs.(has|want|rm|wants)
* gossip.(peers|add|connect)


see https://scuttlebot.io/apis/scuttlebot/ssb.html#createlogstream-source  for more

CAVEAT: only one argument...
`,
	Action: func(ctx *cli.Context) error {
		cmd := ctx.Args().Get(0)
		if cmd == "" {
			return errors.New("call: cmd can't be empty")
		}
		method := strings.Split(cmd, ".")

		args := ctx.Args().Slice()
		var sendArgs []interface{}
		if len(args) > 1 {
			sendArgs = make([]interface{}, len(args)-1)
			for i, v := range args[1:] {
				sendArgs[i] = v
			}
		}

		client, err := newClient(ctx)
		if err != nil {
			return err
		}

		var reply interface{}
		err = client.Async(longctx, &reply, muxrpc.TypeJSON, muxrpc.Method(method), sendArgs...)
		if err != nil {
			return fmt.Errorf("%s: call failed: %w", cmd, err)
		}
		level.Debug(log).Log("event", "call reply")

		jsonReply, err := json.MarshalIndent(reply, "", "  ")
		if err != nil {
			return fmt.Errorf("%s: indent failed: %w", cmd, err)
		}

		_, err = os.Stdout.Write(jsonReply)
		if err != nil {
			return fmt.Errorf("%s: result copy failed: %w", cmd, err)
		}

		return nil
	},
}

var sourceCmd = &cli.Command{
	Name:  "source",
	Usage: "make an simple source call",

	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Value: ""},
		// TODO: Slice of branches
		&cli.IntFlag{Name: "limit", Value: -1},
	},

	Action: func(ctx *cli.Context) error {
		cmd := ctx.Args().Get(0)
		if cmd == "" {
			return errors.New("call: cmd can't be empty")
		}

		v := strings.Split(cmd, ".")

		client, err := newClient(ctx)
		if err != nil {
			return err
		}

		var args = struct {
			ID    string `json:"id,omitempty"`
			Limit int    `json:"limit"`
		}{
			ID:    ctx.String("id"),
			Limit: ctx.Int("limit"),
		}

		src, err := client.Source(longctx, muxrpc.TypeJSON, muxrpc.Method(v), args)
		if err != nil {
			return fmt.Errorf("%s: call failed: %w", cmd, err)
		}
		level.Debug(log).Log("event", "call reply")

		err = jsonDrain(os.Stdout, src)
		return fmt.Errorf("%s: result copy failed: %w", cmd, err)
	},
}

var getCmd = &cli.Command{
	Name:  "get",
	Usage: "get a single message from the database by key (%...)",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "private"},
		&cli.StringFlag{Name: "format", Value: "json"},
	},
	Action: func(ctx *cli.Context) error {
		key, err := refs.ParseMessageRef(ctx.Args().First())
		if err != nil {
			return fmt.Errorf("failed to validate message ref: %w", err)
		}

		client, err := newClient(ctx)
		if err != nil {
			return err
		}

		arg := struct {
			ID      refs.MessageRef `json:"id"`
			Private bool            `json:"private"`
		}{key, ctx.Bool("private")}

		var val interface{}
		err = client.Async(longctx, &val, muxrpc.TypeJSON, muxrpc.Method{"get"}, arg)
		if err != nil {
			return err
		}
		format := strings.ToLower(ctx.String("format"))
		log.Log("event", "get reply", "format", format)
		switch format {
		case "json":
			indented, err := json.MarshalIndent(val, "", "  ")
			if err != nil {
				return err
			}
			os.Stdout.Write(indented)

		default:
			fmt.Printf("%+v\n", val)
		}
		return nil

	},
}

var connectCmd = &cli.Command{
	Name:  "connect",
	Usage: "connect to a remote peer",
	Action: func(ctx *cli.Context) error {
		to := ctx.Args().Get(0)
		if to == "" {
			return errors.New("connect: multiserv addr argument can't be empty")
		}

		client, err := newClient(ctx)
		if err != nil {
			return err
		}

		// try all three of these
		var methods = []muxrpc.Method{
			{"conn", "connect"},   // latest, npm:ssb-conn version, now also supported by go-ssb
			{"gossip", "connect"}, // previous javascript call
			{"ctrl", "connect"},   // previous go-ssb version
		}

		var val string
		for _, m := range methods {
			err = client.Async(longctx, &val, muxrpc.TypeString, m, to)
			if err == nil {
				break
			}
			level.Warn(log).Log("event", "connect command failed", "err", err, "method", m.String())

		}
		log.Log("event", "connect reply")
		goon.Dump(val)
		return nil
	},
}

var blockCmd = &cli.Command{
	Name: "block",
	Action: func(ctx *cli.Context) error {
		client, err := newClient(ctx)
		if err != nil {
			return err
		}

		var blocked = make(map[string]bool)

		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			fr, err := refs.ParseFeedRef(sc.Text())
			if err != nil {
				return err
			}
			blocked[fr.String()] = true
		}
		log.Log("blocking", len(blocked))

		var val interface{}
		err = client.Async(longctx, val, muxrpc.TypeJSON, muxrpc.Method{"ctrl", "block"}, blocked)
		if err != nil {
			return err
		}
		log.Log("event", "block reply")
		goon.Dump(val)
		return nil
	},
}

var groupsCmd = &cli.Command{
	Name:  "groups",
	Usage: "group managment (create, invite, publishTo, etc.)",
	Subcommands: []*cli.Command{
		groupsCreateCmd,
		groupsInviteCmd,
		groupsPublishToCmd,
		groupsJoinCmd,
		/* TODO:
		groupsListCmd,
		groupsMembersCmd,
		*/
	},
}

var groupsCreateCmd = &cli.Command{
	Name:  "create",
	Usage: "create a new empty group",
	Action: func(ctx *cli.Context) error {
		client, err := newClient(ctx)
		if err != nil {
			return err
		}

		name := ctx.Args().First()
		if name == "" {
			return fmt.Errorf("group name can't be empty")
		}

		var val interface{}
		err = client.Async(longctx, &val, muxrpc.TypeJSON, muxrpc.Method{"groups", "create"}, struct {
			Name string `json:"name"`
		}{name})
		if err != nil {
			return err
		}
		log.Log("event", "group created")
		goon.Dump(val)
		return nil
	},
}

var groupsInviteCmd = &cli.Command{
	Name:  "invite",
	Usage: "add people to a group",
	Action: func(ctx *cli.Context) error {
		args := ctx.Args()
		groupID, err := refs.ParseMessageRef(args.First())
		if err != nil {
			return fmt.Errorf("groupID needs to be a valid message ref: %w", err)
		}

		if groupID.Algo() != refs.RefAlgoCloakedGroup {
			return fmt.Errorf("groupID needs to be a cloaked message ref, not %s", groupID.Algo())
		}

		member, err := refs.ParseFeedRef(args.Get(1))
		if err != nil {
			return err
		}

		client, err := newClient(ctx)
		if err != nil {
			return err
		}

		var reply interface{}
		err = client.Async(longctx, &reply, muxrpc.TypeJSON, muxrpc.Method{"groups", "invite"}, groupID.String(), member.String())
		if err != nil {
			return fmt.Errorf("invite call failed: %w", err)
		}
		log.Log("event", "member added", "group", groupID.String(), "member", member.String())
		goon.Dump(reply)
		return nil
	},
}

var groupsPublishToCmd = &cli.Command{
	Name:  "publishTo",
	Usage: "publish a handcrafted JSON blob to a group",
	Action: func(ctx *cli.Context) error {
		var content interface{}
		err := json.NewDecoder(os.Stdin).Decode(&content)
		if err != nil {
			return fmt.Errorf("publish/raw: invalid json input from stdin: %w", err)
		}

		groupID, err := refs.ParseMessageRef(ctx.Args().First())
		if err != nil {
			return fmt.Errorf("groupID needs to be a valid message ref: %w", err)
		}

		if groupID.Algo() != refs.RefAlgoCloakedGroup {
			return fmt.Errorf("groupID needs to be a cloaked message ref, not %s", groupID.Algo())
		}

		client, err := newClient(ctx)
		if err != nil {
			return err
		}

		var reply interface{}
		err = client.Async(longctx, &reply, muxrpc.TypeJSON, muxrpc.Method{"groups", "publishTo"}, groupID.String(), content)
		if err != nil {
			return fmt.Errorf("publish call failed: %w", err)
		}
		log.Log("event", "publishTo", "type", "raw")
		goon.Dump(reply)
		return nil
	},
}

var groupsJoinCmd = &cli.Command{
	Name:   "join",
	Usage:  "manually join a group by adding the group key",
	Action: todo,
}

var inviteCmds = &cli.Command{
	Name: "invite",
	Subcommands: []*cli.Command{
		inviteCreateCmd,
		inviteAcceptCmd,
	},
}

var inviteCreateCmd = &cli.Command{
	Name:  "create",
	Usage: "register and return an invite for somebody else to accept",
	Flags: []cli.Flag{
		&cli.UintFlag{Name: "uses", Value: 1, Usage: "How many times an invite can be used"},
	},
	Action: func(ctx *cli.Context) error {

		client, err := newClient(ctx)
		if err != nil {
			return err
		}

		var args legacyinvites.CreateArguments
		args.Uses = ctx.Uint("uses")

		var code string
		err = client.Async(longctx, &code, muxrpc.TypeString, muxrpc.Method{"invite", "create"}, args)
		if err != nil {
			return err
		}
		fmt.Println(code)
		return nil
	},
}

var inviteAcceptCmd = &cli.Command{
	Name:   "accept",
	Usage:  "use an invite code",
	Action: todo,
}
