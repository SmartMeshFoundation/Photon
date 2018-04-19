package ice

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
