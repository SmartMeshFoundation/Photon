package main

type RaidenParam struct {
	//本地注释:缓存文件夹
	datadir        string
	//本地注释:API服务地址和端口
	api_address    string
	listen_address string
	//本地注释:账户地址
	address        string
	//本地注释:账户密钥地址
	keystore_path  string
	//本地注释:节点发现地址
	discovery_contract_address string
	//本地注释:合约地址
	registry_contract_address string
	//本地注释:账户密钥密码文件
	password_file             string
	//本地注释:NAT类型
	nat                       string
	//本地注释:geth服务地址
	eth_rpc_endpoint string
	//本地注释:退出事件
	conditionquit    string
	//本地注释:调试标志
	debug            bool
}

func (rp *RaidenParam) getParam() []string {
	var param []string

	param = append(param, "--datadir="+rp.datadir)
	param = append(param, "--api-address="+rp.api_address)
	param = append(param, "--listen-address="+rp.listen_address)
	param = append(param, "--address="+rp.address)
	param = append(param, "--keystore-path="+rp.keystore_path)
	param = append(param, "--discovery-contract-address="+rp.discovery_contract_address)
	param = append(param, "--registry-contract-address="+rp.registry_contract_address)
	param = append(param, "--password-file="+rp.password_file)
	param = append(param, "--nat="+rp.nat)
	param = append(param, "--eth-rpc-endpoint="+rp.eth_rpc_endpoint)
	param = append(param, "--conditionquit="+rp.conditionquit)

	if rp.debug == true {
		param = append(param, "--debug")
	}
	return param
}
