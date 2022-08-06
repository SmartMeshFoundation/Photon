package params

//TestLogServer only for test, enabled if --debug flag is set
var TestLogServer = "http://transport01.smartmesh.cn:8008"

//TrustMatrixServers matrix server config
var TrustMatrixServers = map[string]string{
	"transport01.smartmesh.cn": "http://transport01.smartmesh.cn:8008",
	//"transport02.smartmesh.cn": "http://transport02.smartmesh.cn:8008",
	"transport13.smartmesh.cn": "http://transport13.smartmesh.cn:8008",
}

//ContractSignaturePrefix for EIP191 https://github.com/ethereum/EIPs/blob/master/EIPS/eip-191.md
var ContractSignaturePrefix = []byte("\x19Spectrum Signed Message:\n")

const (
	//ContractBalanceProofMessageLength balance proof  length
	ContractBalanceProofMessageLength = "176"
	//ContractBalanceProofDelegateMessageLength update balance proof delegate length
	ContractBalanceProofDelegateMessageLength = "144"
	//ContractCooperativeSettleMessageLength cooperative settle channel proof length
	ContractCooperativeSettleMessageLength = "176"
	//ContractDisposedProofMessageLength annouce disposed proof length
	ContractDisposedProofMessageLength = "136"
	//ContractWithdrawProofMessageLength withdraw proof length
	ContractWithdrawProofMessageLength = "156"
	//ContractUnlockDelegateProofMessageLength unlock delegate proof length
	ContractUnlockDelegateProofMessageLength = "188"
)
