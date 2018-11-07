module github.com/SmartMeshFoundation/Photon

require (
	github.com/ant0ine/go-json-rest v3.3.2+incompatible
	github.com/aristanetworks/goarista v0.0.0-20181101003910-5bb443fba8e0 // indirect
	github.com/asdine/storm v2.1.1+incompatible
	github.com/coreos/bbolt v1.3.1-coreos.6
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/edsrzf/mmap-go v0.0.0-20170320065105-0bce6a688712 // indirect
	github.com/elastic/gosigar v0.9.0 // indirect
	github.com/ethereum/go-ethereum v1.8.17
	github.com/fjl/memsize v0.0.0-20180929194037-2a09253e352a // indirect
	github.com/go-stack/stack v1.8.0
	github.com/howeyc/gopass v0.0.0-20170109162249-bf9dde6d0d2c
	github.com/huin/goupnp v1.0.0 // indirect
	github.com/influxdata/influxdb v1.7.0 // indirect
	github.com/influxdata/platform v0.0.0-20181107003602-9b529771ebb3 // indirect
	github.com/jackpal/go-nat-pmp v1.0.1 // indirect
	github.com/karalabe/hid v0.0.0-20180420081245-2b4488a37358 // indirect
	github.com/mattn/go-colorable v0.0.9
	github.com/mattn/go-xmpp v0.0.1
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/nkbai/log v0.0.0-20180519141659-86998e435e8c
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pkg/errors v0.8.0
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/rs/cors v1.6.0 // indirect
	github.com/syndtr/goleveldb v0.0.0-20181105012736-f9080354173f // indirect
	github.com/theckman/go-flock v0.7.0
	go.etcd.io/bbolt v1.3.0 // indirect
	golang.org/x/crypto v0.0.1
	golang.org/x/net v0.0.1
	golang.org/x/sys v0.0.1
	gopkg.in/karalabe/cookiejar.v2 v2.0.0-20150724131613-8dcd6a7f4951 // indirect
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20180723110524-d53328019b21 // indirect
	gopkg.in/urfave/cli.v1 v1.20.0
)

replace (
	github.com/ethereum/go-ethereum v1.8.17 => github.com/nkbai/go-ethereum v1.9.1
	github.com/mattn/go-xmpp v0.0.1 => github.com/nkbai/go-xmpp v0.0.1
	golang.org/x/crypto v0.0.1 => github.com/golang/crypto v0.0.0-20181106171534-e4dc69e5b2fd
	golang.org/x/net v0.0.1 => github.com/golang/net v0.0.0-20181106065722-10aee1819953

	golang.org/x/sys v0.0.1 => github.com/golang/sys v0.0.0-20181106135930-3a76605856fd
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
	golang.org/x/tools v0.0.1 => github.com/golang/tools v0.0.0-20181106213628-e21233ffa6c3
)
