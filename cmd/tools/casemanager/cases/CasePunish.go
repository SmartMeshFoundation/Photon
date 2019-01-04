package cases

// CasePunish : test for punish
func (cm *CaseManager) CasePunish() (err error) {
	//env, err := models.NewTestEnv("./cases/CasePunish.ENV",cm.UseMatrix)
	//if err != nil {
	//	return
	//}
	//defer func() {
	//	if env.Debug == false {
	//		env.KillAllPhotonNodes()
	//	}
	//}()
	//// 源数据
	//transAmount := int32(20)
	//tokenAddress := env.Tokens[0].TokenAddress.String()
	//N0, N1, N2 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	//models.Logger.Println(env.CaseName + " BEGIN ====>")
	//// 启动节点,让节点0发送SendAnnounce
	//N0.StartWithConditionQuit(env, &params.ConditionQuit{
	//	QuitEvent: "EventSendAnnouncedDisposedResponseBefore",
	//})
	//N1.Start(env)
	//N2.Start(env)
	////channel := N1.GetChannelWith(N0, tokenAddress)
	//go N0.SendTrans(tokenAddress, transAmount, N2.Address, false)
	//
	////
	//models.Logger.Println("请手动调用force-unlock后输入任意键继续...")
	//var temp string
	//_, err = fmt.Scanf("%s", &temp)
	//if err != nil {
	//	panic(err)
	//}
	//N0.ReStartWithoutConditionquit(env)
	//time.Sleep(100 * time.Second)
	//models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
