// SPDX-FileCopyrightText: 2021 The Go-SSB Authors
//
// SPDX-License-Identifier: MIT

package sbot

import (
	"net"
	"os"
	"sort"
	"time"

	"github.com/dustin/go-humanize"
	"go.cryptoscope.co/netwrap"

	"go.cryptoscope.co/ssb"
	multiserver "go.mindeco.de/ssb-multiserver"
)

// Status returns the current status of information about the bot
func (sbot *Sbot) Status() (ssb.Status, error) {
	s := ssb.Status{
		PID:   os.Getpid(),
		Root:  sbot.ReceiveLog.Seq(),
		Blobs: sbot.WantManager.AllWants(),
	}

	edps := sbot.Network.GetAllEndpoints()

	sort.Sort(byConnTime(edps))

	for _, es := range edps {
		var ms multiserver.NetAddress
		ms.Ref = es.ID
		if tcpAddr, ok := netwrap.GetAddr(es.Addr, "tcp").(*net.TCPAddr); ok {
			ms.Addr = *tcpAddr
		}
		s.Peers = append(s.Peers, ssb.PeerStatus{
			Addr:  ms.String(),
			Since: humanize.Time(time.Now().Add(-es.Since)),
		})
	}

	var idxState ssb.IndexStates
	sbot.indexStateMu.Lock()

	for n, s := range sbot.indexStates {
		idxState = append(idxState, ssb.IndexState{
			Name:  n,
			State: s,
		})
	}

	sbot.indexStateMu.Unlock()

	sort.Sort(byName(idxState))
	s.Indicies = idxState

	return s, nil
}

type byConnTime []ssb.EndpointStat

func (bct byConnTime) Len() int { return len(bct) }

func (bct byConnTime) Less(i int, j int) bool {
	return bct[i].Since < bct[j].Since
}

func (bct byConnTime) Swap(i int, j int) { bct[i], bct[j] = bct[j], bct[i] }

type byName ssb.IndexStates

func (bn byName) Len() int { return len(bn) }

func (bn byName) Less(i int, j int) bool {
	return bn[i].Name < bn[j].Name
}

func (bn byName) Swap(i int, j int) { bn[i], bn[j] = bn[j], bn[i] }
