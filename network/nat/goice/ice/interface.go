package ice

import "github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/stun"

type CandidateGetter interface {
	/*
		获取有一部分信息的candidiate.第一个是本机主要地址,最后一个是缺省 Candidate
	*/
	GetCandidates() (candidates []*Candidate, err error)
}

//treat stun and turn as the same ...
type StunTranporter interface {
	CandidateGetter
	Close()
	GetListenCandidiates() []string
	/*
		transporter is using turn?
	*/
	//IsTurn() bool
}

/*
建立连接需要的本地的 简易stun 服务器.
*/
type ServerSocker interface {
	/*
		指定 从 from 到 to 发送一个消息
		from 有可能是本地地址,也有可能是 turn server relay 的地址
	*/
	sendStunMessageSync(msg *stun.Message, fromaddr, toaddr string) (res *stun.Message, err error)
	/*
		暂时没用,先留着
	*/
	sendStunMessageWithResult(msg *stun.Message, fromaddr, toaddr string) (key stun.TransactionID, ch chan *serverSockResponse, err error)
	/*
		参数含义和 sync 是一致的,不用等待结果.
	*/
	sendStunMessageAsync(msg *stun.Message, fromaddr, toaddr string) error
	/*
		从 from 到 to 发送一个数据包,
		如果 from 是本机地址,则直接发送,
		如果是 turn server relay address, 那么需要经由 turn server 中转.
		也就是会把 data 封装到 SendIndication 或者 ChannelDataRequest中
	*/
	sendData(data []byte, fromaddr, toaddr string) error
	/*
		关闭连接
	*/
	Close()
	/*
		ice check 真正完毕以后,
		需要开启刷新权限以及 keep alive 等操作.
		todo 这里的 mode 定义并不清晰,需要梳理.
	*/
	FinishNegotiation(mode serverSockMode)
	//StartRefresh()
}
