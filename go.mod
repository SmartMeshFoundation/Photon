module github.com/SmartMeshFoundation/Photon

require (
	github.com/DataDog/zstd v1.3.4 // indirect
	github.com/Sereal/Sereal v0.0.0-20180905114147-563b78806e28 // indirect
	github.com/StackExchange/wmi v0.0.0-20180725035823-b12b22c5341f // indirect
	github.com/ant0ine/go-json-rest v3.3.2+incompatible
	github.com/aristanetworks/goarista v0.0.0-20181101003910-5bb443fba8e0 // indirect
	github.com/asdine/storm v2.1.1+incompatible
	github.com/btcsuite/btcd v0.0.0-20180924021209-2a560b2036be // indirect
	github.com/cespare/cp v0.1.0 // indirect
	github.com/coreos/bbolt v1.3.1-coreos.6
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/edsrzf/mmap-go v0.0.0-20170320065105-0bce6a688712 // indirect
	github.com/elastic/gosigar v0.9.0 // indirect
	github.com/ethereum/go-ethereum v1.8.17
	github.com/fjl/memsize v0.0.0-20180929194037-2a09253e352a // indirect
	github.com/go-ole/go-ole v1.2.1 // indirect
	github.com/go-stack/stack v1.8.0
	github.com/hashicorp/go.net v0.0.0-20151006203346-104dcad90073 // indirect
	github.com/hashicorp/mdns v0.0.0-20170221172940-4e527d9d8081
	github.com/howeyc/gopass v0.0.0-20170109162249-bf9dde6d0d2c
	github.com/huin/goupnp v1.0.0 // indirect
	github.com/influxdata/influxdb v1.7.0 // indirect
	github.com/influxdata/platform v0.0.0-20181107003602-9b529771ebb3 // indirect
	github.com/jackpal/go-nat-pmp v1.0.1 // indirect
	github.com/karalabe/hid v0.0.0-20180420081245-2b4488a37358 // indirect
	github.com/mattn/go-colorable v0.0.9
	github.com/mattn/go-xmpp v0.0.0-20180505113305-e543ad3fcd51
	github.com/miekg/dns v1.0.13 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/onsi/gomega v1.4.2 // indirect
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pkg/errors v0.8.0
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/rs/cors v1.6.0 // indirect
	github.com/stretchr/testify v1.2.2
	github.com/syndtr/goleveldb v0.0.0-20181105012736-f9080354173f // indirect
	github.com/theckman/go-flock v0.7.0
	github.com/urfave/cli v1.20.0
	github.com/vmihailenco/msgpack v4.0.0+incompatible // indirect
	golang.org/x/sys v0.0.0-20181023152157-44b849a8bc13
	gopkg.in/karalabe/cookiejar.v2 v2.0.0-20150724131613-8dcd6a7f4951 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20180723110524-d53328019b21 // indirect
	gopkg.in/urfave/cli.v1 v1.20.0
)

replace (
	github.com/ethereum/go-ethereum v1.8.17 => github.com/nkbai/go-ethereum v1.9.1
	github.com/mattn/go-xmpp v0.0.0-20180505113305-e543ad3fcd51 => github.com/nkbai/go-xmpp v0.0.1
	golang.org/x/crypto v0.0.0-20181015023909-0c41d7ab0a0e => github.com/golang/crypto v0.0.0-20181106171534-e4dc69e5b2fd
	golang.org/x/net v0.0.0-20181023162649-9b4f9f5ad519 => github.com/golang/net v0.0.0-20181106065722-10aee1819953

	golang.org/x/sys v0.0.0-20181023152157-44b849a8bc13 => github.com/golang/sys v0.0.0-20181106135930-3a76605856fd
)
