// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package debug

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"

	"github.com/SmartMeshFoundation/Photon/utils"

	//need pprof
	_ "net/http/pprof"
	"os"
	"runtime"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/log/term"
	"github.com/mattn/go-colorable"
	"gopkg.in/urfave/cli.v1"
)

var (
	verbosityFlag = cli.IntFlag{
		Name:  "verbosity",
		Usage: "Logging verbosity: 0=silent, 1=error, 2=warn, 3=info, 4=debug, 5=trace",
		Value: 3,
	}
	vmoduleFlag = cli.StringFlag{
		Name:  "vmodule",
		Usage: "Per-module verbosity: comma-separated list of <pattern>=<level> (e.g. eth/*=5,p2p=4)",
		Value: "",
	}
	backtraceAtFlag = cli.StringFlag{
		Name:  "backtrace",
		Usage: "Request a stack trace at a specific logging statement (e.g. \"block.go:271\")",
		Value: "",
	}
	debugFlag = cli.BoolFlag{
		Name:  "debug",
		Usage: "Prepends log messages with call-site location (file and line number)",
	}
	pprofFlag = cli.BoolFlag{
		Name:  "pprof",
		Usage: "Enable the pprof HTTP server",
	}
	pprofPortFlag = cli.IntFlag{
		Name:  "pprofport",
		Usage: "pprof HTTP server listening port",
		Value: 6060,
	}
	pprofAddrFlag = cli.StringFlag{
		Name:  "pprofaddr",
		Usage: "pprof HTTP server listening interface",
		Value: "127.0.0.1",
	}
	memprofilerateFlag = cli.IntFlag{
		Name:  "memprofilerate",
		Usage: "Turn on memory profiling with the given rate",
		Value: runtime.MemProfileRate,
	}
	blockprofilerateFlag = cli.IntFlag{
		Name:  "blockprofilerate",
		Usage: "Turn on block profiling with the given rate",
	}
	cpuprofileFlag = cli.StringFlag{
		Name:  "cpuprofile",
		Usage: "Write CPU profile to the given file",
	}
	traceFlag = cli.StringFlag{
		Name:  "trace",
		Usage: "Write execution trace to the given file",
	}
	logFileFlag = cli.StringFlag{
		Name:  "logfile",
		Usage: "redirect log to this the given file",
	}
)

// Flags holds all command-line flags required for debugging.
var Flags = []cli.Flag{
	verbosityFlag, vmoduleFlag, backtraceAtFlag, debugFlag,
	pprofFlag, pprofAddrFlag, pprofPortFlag,
	memprofilerateFlag, blockprofilerateFlag, cpuprofileFlag, traceFlag, logFileFlag,
}

var glogger *log.GlogHandler

func init() {
}

//获取本机mac地址作为id,如果有多个mac就拼在一起,长度不超过32,如果没有mac地址,就返回一个随机字符串
func mac() string {
	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Error : " + err.Error())
	}
	sort.Slice(interfaces, func(i, j int) bool {
		return bytes.Compare(interfaces[i].HardwareAddr, interfaces[j].HardwareAddr) < 0
	})
	var r string
	for _, s := range interfaces {
		r += hex.EncodeToString(s.HardwareAddr)
	}
	if r == "" {
		r = utils.RandomString(20)
	}
	if len(r) > 32 {
		r = r[:32]
	}
	return r
}

// Setup initializes profiling and logging based on the CLI flags.
// It should be called as early as possible in the program.
func Setup(ctx *cli.Context) (err error) {

	var fileHandler log.Handler

	// file handler
	if len(ctx.String(logFileFlag.Name)) > 0 {
		fmt.Printf("log will be write to %s\n", ctx.String(logFileFlag.Name))
		fileHandler, err = log.FileHandler(ctx.String(logFileFlag.Name), log.TerminalFormat(false))
		if err != nil {
			return
		}
	}
	// console handler
	usecolor := term.IsTty(os.Stderr.Fd()) && os.Getenv("TERM") != "dumb"
	output := io.Writer(os.Stderr)
	if usecolor {
		output = colorable.NewColorableStderr()
	}
	consoleHandler := log.StreamHandler(output, log.TerminalFormat(usecolor))
	glogger = log.NewGlogHandler(log.TeeHandler(consoleHandler, fileHandler, nil))

	// logging
	log.PrintOrigins(ctx.GlobalBool(debugFlag.Name))
	glogger.Verbosity(log.Lvl(ctx.GlobalInt(verbosityFlag.Name)))
	err = glogger.Vmodule(ctx.GlobalString(vmoduleFlag.Name))
	if err != nil {
		return err
	}
	err = glogger.BacktraceAt(ctx.GlobalString(backtraceAtFlag.Name))
	if err != nil {
		//todo fixit ,return error when backtraceAtFlag is empty
	}
	log.Root().SetHandler(glogger)

	// profiling, tracing
	runtime.MemProfileRate = ctx.GlobalInt(memprofilerateFlag.Name)
	Handler.SetBlockProfileRate(ctx.GlobalInt(blockprofilerateFlag.Name))
	if traceFile := ctx.GlobalString(traceFlag.Name); traceFile != "" {
		if err := Handler.StartGoTrace(traceFile); err != nil {
			return err
		}
	}
	if cpuFile := ctx.GlobalString(cpuprofileFlag.Name); cpuFile != "" {
		if err := Handler.StartCPUProfile(cpuFile); err != nil {
			return err
		}
	}

	// pprof server
	if ctx.GlobalBool(pprofFlag.Name) {

		address := fmt.Sprintf("%s:%d", ctx.GlobalString(pprofAddrFlag.Name), ctx.GlobalInt(pprofPortFlag.Name))
		go func() {
			log.Info("Starting pprof server", "addr", fmt.Sprintf("http://%s/debug/pprof", address))
			if err := http.ListenAndServe(address, nil); err != nil {
				log.Error("Failure in running pprof server", "err", err)
			}
		}()
	}
	return nil
}

// Exit stops all running profiles, flushing their output to the
// respective file.
func Exit() {
	err := Handler.StopCPUProfile()
	err = Handler.StopGoTrace()
	if err != nil {
		fmt.Printf("StopGoTrace err %s ", err)
	}
}
