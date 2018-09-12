package cases

// LocalCase : only for local test
func (cm *CaseManager) LocalCase() (err error) {
	//env, err := models.NewTestEnv("./cases/LocalCase.ENV")
	//if err != nil {
	//	return
	//}
	//defer func() {
	//	if env.Debug == false {
	//		env.KillAllRaidenNodes()
	//	}
	//}()
	//// 源数据
	//transAmount := int32(20)
	//tokenAddress := env.Tokens[0].TokenAddress.String()
	//N0, N1, N2, N3 := env.Nodes[0], env.Nodes[1], env.Nodes[2], env.Nodes[3]
	//models.Logger.Println(env.CaseName + " BEGIN ====>")
	//// 启动节点2，3
	//N0.Start(env)
	//N1.Start(env)
	//N2.StartWithConditionQuit(env, &params.ConditionQuit{
	//	QuitEvent: "EventSendRevealSecretBefore",
	//})
	//N3.Start(env)
	//go N0.SendTrans(tokenAddress, transAmount, N3.Address, false)
	//time.Sleep(180 * time.Second)
	//N2.ReStartWithoutConditionquit(env)
	//time.Sleep(1000 * time.Second)
	//models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
