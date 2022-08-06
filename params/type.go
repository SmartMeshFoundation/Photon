package params

/*
ENV photon运行环境枚举
*/
type ENV string

const (
	// Dev 开发环境,即photon工作在自定义的链上
	Dev = "dev"

	// TestNet 测试网,即photon工作在Spectrum测试链上
	TestNet = "testnet"

	// MainNet 主网,即photon工作在Spectrum主链上
	MainNet = "mainnet"
)

//NetworkMode is transport status
type NetworkMode int

const (
	//NoNetwork 不与其他节点之间进行通信,仅供测试使用
	// NoNetwork : Node does not  communicates with other nodes, just for test.
	NoNetwork NetworkMode = iota + 1
	//UDPOnly 通过udp ip 端口对外暴露服务,通过使用MDNS进行节点服务发现
	// UDPOnly : expose service via udp ip,
	UDPOnly
	//XMPPOnly 通过XMPP服务器进行通信,目前暂时不使用
	// XMPPOnly : communicate via XMPP server.
	XMPPOnly
	//MixUDPXMPP 适应无网通信需要,将上面两种方式混合使用,有网时使用XMPP建立连接,无网时则使用 udp 直接暴露 ip 端口
	// MixUDPXMPP : used for Internet-free network, combining UDPOnly and XMPPOnly.
	// While Internet, it use XMPP to create connection; while Internet-free, it use udp to expose ip port.
	MixUDPXMPP
	//MixUDPMatrix Matrix and UDP at the same time
	MixUDPMatrix
)

//ConditionQuit is for test
type ConditionQuit struct {
	QuitEvent  string //name match
	IsBefore   bool   //quit before event occur
	RandomQuit bool   //random exit
}
