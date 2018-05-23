package main

// RaidenParam : smartraiden start commands
type RaidenParam struct {
	//Caching folder
	datadir string
	//API service address and port
	apiAddress    string
	listenAddress string
	//Account address
	address string
	//key address of the Account
	keystorePath string
	//Node discovery address
	discoveryContractAddress string
	//Contract address
	registryContractAddress string
	//The key and password file of the account
	passwordFile string
	//NAT type
	nat string
	//Geth service address
	ethRPCEndpoint string
	//Exiting event
	conditionquit string
	//Debug sign
	debug bool
}

func (rp *RaidenParam) getParam() []string {
	var param []string

	param = append(param, "--datadir="+rp.datadir)
	param = append(param, "--api-address="+rp.apiAddress)
	param = append(param, "--listen-address="+rp.listenAddress)
	param = append(param, "--address="+rp.address)
	param = append(param, "--keystore-path="+rp.keystorePath)
	param = append(param, "--discovery-contract-address="+rp.discoveryContractAddress)
	param = append(param, "--registry-contract-address="+rp.registryContractAddress)
	param = append(param, "--password-file="+rp.passwordFile)
	param = append(param, "--nat="+rp.nat)
	param = append(param, "--eth-rpc-endpoint="+rp.ethRPCEndpoint)
	param = append(param, "--verbosity=5")
	if rp.debug == true {
		param = append(param, "--debug")
	}
	return param
}
