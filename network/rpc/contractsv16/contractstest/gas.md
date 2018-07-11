#花费 gas 统计
* OpenChannel 135638
* Deposit 68764
* CloseChannel 无证据:55875
* CloseChannel 有证据:114609
* UpdateNonClosingBalanceProof:97407
* settle channel:81722(实际),163444(估计)
* CooperativeSettle: 92228,184455
* withdraw:129081
* unlock:69590(3个锁)
* PunishObsoleteUnlock:34312,49312

## 记录
INFO [06-29|17:08:35.243] OpenChannel complete.. gasLimit=135638,gasUsed=135638 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|17:08:43.285] Deposit2 complete.. gasLimit=68764,gasUsed=68764 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|17:08:51.318] OpenChannel complete.. gasLimit=135638,gasUsed=135638 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|17:08:59.356] Deposit2 complete.. gasLimit=68764,gasUsed=68764 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|17:09:03.372] CloseChannel no evidence complete.. gasLimit=55875,gasUsed=55875 fn=contracts_test.go:TestCloseChannel1:209
INFO [06-29|17:09:11.4] OpenChannel complete.. gasLimit=135638,gasUsed=135638 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|17:09:19.438] Deposit2 complete.. gasLimit=68764,gasUsed=68764 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|17:09:23.458] CloseChannel with evidence  complete.. gasLimit=114609,gasUsed=114609 fn=contracts_test.go:TestCloseChannel2:314
INFO [06-29|17:09:31.489] OpenChannel complete.. gasLimit=135638,gasUsed=135638 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|17:09:39.53] Deposit2 complete.. gasLimit=68764,gasUsed=68764 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|17:09:42.558] CloseChannel   complete.. gasLimit=114609,gasUsed=114609 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle:398
INFO [06-29|17:09:47.59] UpdateNonClosingBalanceProof   complete.. gasLimit=97407,gasUsed=97407 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle:414
INFO [06-29|17:10:48.87] SettleChannel  complete.. gasLimit=81722,gasUsed=163444 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle:455
INFO [06-29|17:10:56.905] OpenChannel complete.. gasLimit=135638,gasUsed=135638 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|17:11:04.947] Deposit2 complete.. gasLimit=68764,gasUsed=68764 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|17:11:08.976] CloseChannel   complete.. gasLimit=114592,gasUsed=114592 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle2:479
INFO [06-29|17:11:13.016] UpdateNonClosingBalanceProof   complete.. gasLimit=97424,gasUsed=97424 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle2:495
INFO [06-29|17:11:17.075] UpdateNonClosingBalanceProof2   complete.. gasLimit=67407,gasUsed=67407 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle2:518
INFO [06-29|17:11:25.107] OpenChannel complete.. gasLimit=135638,gasUsed=135638 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|17:11:32.149] Deposit2 complete.. gasLimit=68764,gasUsed=68764 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|17:11:37.201] CooperativeSettle   complete.. gasLimit=92228,gasUsed=184455 fn=contracts_test.go:TestCooperateSettleChannel:583
INFO [06-29|17:11:45.234] OpenChannel complete.. gasLimit=135638,gasUsed=135638 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|17:11:52.279] Deposit2 complete.. gasLimit=68764,gasUsed=68764 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|17:11:57.319] WithDraw complete.. gasLimit=129081,gasUsed=129081 fn=contracts_test.go:TestWithdraw:715
INFO [06-29|17:12:01.336] RegisterSecret success.. gasLimit=24443,gasUsed=24443 fn=contracts_test.go:TestRegisterSecret:745
INFO [06-29|17:12:09.37] OpenChannel complete.. gasLimit=135638,gasUsed=135638 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|17:12:17.415] Deposit2 complete.. gasLimit=68764,gasUsed=68764 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|17:12:21.441] close channel successful,gasused=114609,gasLimit=114609 fn=contracts_test.go:TestUnlock:796
INFO [06-29|17:12:37.526] UpdateNonClosingBalanceProof successful,gasused=114438,gasLimit=114438,locksroot=9703b14733debfe573b967973e0cdc5ecfc5720bebbae832be3a64899bda0fa2,transferamount=10 fn=contracts_test.go:TestUnlock:841
INFO [06-29|17:12:41.558] unlock success,gasUsed=69590,gasLimit=69590,txhash=0xfe771d8dcdcf4bee3330e16f29f5e0c82f30070703a8dba64798bd682d342d07 fn=contracts_test.go:TestUnlock:859
INFO [06-29|17:13:26.78] settle channel complete ,gasused=81722,gasLimit=163444 fn=contracts_test.go:TestUnlock:907
INFO [06-29|17:13:35.812] OpenChannel complete.. gasLimit=135638,gasUsed=135638 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|17:13:42.858] Deposit2 complete.. gasLimit=68764,gasUsed=68764 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|17:13:46.889] close channel successful,gasused=114609,gasLimit=114609 fn=contracts_test.go:TestPunishObsoleteUnlock:988
INFO [06-29|17:14:02.971] UpdateNonClosingBalanceProof successful,gasused=114408,gasLimit=114408,locksroot=615dcb256021327781495623d283dc650ec2504989d8dfc22fa0e72066a4762f,transferamount=10 fn=contracts_test.go:TestPunishObsoleteUnlock:1033
INFO [06-29|17:14:06.994] PunishObsoleteUnlock success,gasUsed=34312,gasLimit=49312,txhash=0x45ccb1edda1c85f09045c5f66fff69232eb50f7956d60cb057586f27cea8047f fn=contracts_test.go:TestPunishObsoleteUnlock:1063
