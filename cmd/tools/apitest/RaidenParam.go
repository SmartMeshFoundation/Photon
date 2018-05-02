package main

type RaidenParam struct {
	//Caching folder
	datadir string
	//API service address and port
	api_address    string
	listen_address string
	//Account address
	address string
	//key address of the Account
	keystore_path string
	//Node discovery address
	discovery_contract_address string
	//Contract address
	registry_contract_address string
	//The key and password file of the account
	password_file string
	//NAT type
	nat string
	//Geth service address
	eth_rpc_endpoint string
	//Exiting event
	conditionquit string
	//Debug sign
	debug bool
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
	param = append(param, "--verbosity=5")
	if rp.debug == true {
		param = append(param, "--debug")
	}
	return param
}
