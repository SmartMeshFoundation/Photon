// SPDX-FileCopyrightText: 2021 The Go-SSB Authors
//
// SPDX-License-Identifier: MIT

module go.cryptoscope.co/ssb

go 1.13

replace github.com/ant0ine/go-json-rest@v3.3.2+incompatible/rest => ../ant0ine/go-json-rest/rest

require (
	filippo.io/edwards25519 v1.0.0-rc.1
	github.com/RoaringBitmap/roaring v0.6.1
	github.com/VividCortex/gohistogram v1.0.0
	github.com/ant0ine/go-json-rest v3.3.2+incompatible
	github.com/davecgh/go-spew v1.1.1
	github.com/dgraph-io/badger/v3 v3.2011.1
	github.com/dgraph-io/sroar v0.0.0-20210524170324-9b164cbe6e02
	github.com/dustin/go-humanize v1.0.0
	github.com/ethereum/go-ethereum v1.8.17
	github.com/go-kit/kit v0.10.0
	github.com/gorilla/websocket v1.4.1
	github.com/hashicorp/go-multierror v1.1.1
	github.com/ip2location/ip2location-go/v9 v9.4.0
	github.com/karlseguin/ccache v2.0.3+incompatible // indirect
	github.com/keks/persist v0.0.0-20210520094901-9bdd97c1fad2
	github.com/keks/testops v0.1.0
	github.com/kr/pretty v0.2.0 // indirect
	github.com/kylelemons/godebug v1.1.0
	github.com/libp2p/go-reuseport v0.0.1
	github.com/machinebox/progress v0.2.0
	github.com/matryer/is v1.3.0 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/maxbrunsfeld/counterfeiter/v6 v6.2.2
	github.com/mvdan/xurls v1.1.0
	github.com/rs/cors v1.7.0
	github.com/shurcooL/go v0.0.0-20200502201357-93f07166e636 // indirect
	github.com/shurcooL/go-goon v0.0.0-20170922171312-37c2f522c041
	github.com/ssb-ngi-pointer/go-metafeed v1.1.1
	github.com/stretchr/testify v1.7.0
	github.com/ugorji/go/codec v1.2.6
	github.com/zeebo/bencode v1.0.0
	go.cryptoscope.co/luigi v0.3.6
	go.cryptoscope.co/margaret v0.4.1
	go.cryptoscope.co/muxrpc/v2 v2.0.13
	go.cryptoscope.co/netwrap v0.1.4
	go.cryptoscope.co/nocomment v0.0.0-20210520094614-fb744e81f810
	go.cryptoscope.co/secretstream v1.2.10
	go.mindeco.de v1.12.0
	go.mindeco.de/ssb-gabbygrove v0.2.1
	go.mindeco.de/ssb-multiserver v0.1.4
	go.mindeco.de/ssb-refs v0.5.1
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
	golang.org/x/exp v0.0.0-20190411193353-0480eff6dd7c // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/text v0.3.7
	gonum.org/v1/gonum v0.0.0-20190904110519-2065cbd6b42a
	gopkg.in/urfave/cli.v2 v2.0.0-20190806201727-b62605953717
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/kv v1.0.3
)
