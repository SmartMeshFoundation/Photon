#花费 gas 统计
* OpenChannel 135688
* Deposit 68841
* CloseChannel 无证据:40943
* CloseChannel 有证据:99334
* UpdateNonClosingBalanceProof:96794
* settle channel:79274(实际),158547(估计)
* CooperativeSettle: 89793,179585
* withdraw:218513,308513
* unlock:69715(3个锁)
* PunishObsoleteUnlock:34384,49384

## 记录
INFO [06-29|22:17:39.343] OpenChannel complete.. gasLimit=135688,gasUsed=135688 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|22:17:47.379] Deposit2 complete.. gasLimit=68841,gasUsed=68841 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|22:17:55.41] OpenChannel complete.. gasLimit=135688,gasUsed=135688 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|22:18:03.448] Deposit2 complete.. gasLimit=68841,gasUsed=68841 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|22:18:06.461] CloseChannel no evidence complete.. gasLimit=40943,gasUsed=40943 fn=contracts_test.go:TestCloseChannel1:209
INFO [06-29|22:18:15.49] OpenChannel complete.. gasLimit=135688,gasUsed=135688 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|22:18:22.518] Deposit2 complete.. gasLimit=68841,gasUsed=68841 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|22:18:27.537] CloseChannel with evidence  complete.. gasLimit=99334,gasUsed=99334 fn=contracts_test.go:TestCloseChannel2:310
INFO [06-29|22:18:34.561] OpenChannel complete.. gasLimit=135688,gasUsed=135688 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|22:18:43.596] Deposit2 complete.. gasLimit=68841,gasUsed=68841 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|22:18:47.617] CloseChannel   complete.. gasLimit=99334,gasUsed=99334 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle:392
INFO [06-29|22:18:51.644] UpdateNonClosingBalanceProof   complete.. gasLimit=96794,gasUsed=96794 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle:408
INFO [06-29|22:19:53.89] SettleChannel  complete.. gasLimit=79274,gasUsed=158547 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle:449
INFO [06-29|22:20:00.917] OpenChannel complete.. gasLimit=135688,gasUsed=135688 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|22:20:08.951] Deposit2 complete.. gasLimit=68841,gasUsed=68841 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|22:20:12.972] CloseChannel   complete.. gasLimit=99351,gasUsed=99351 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle2:473
INFO [06-29|22:20:16.996] UpdateNonClosingBalanceProof   complete.. gasLimit=96794,gasUsed=96794 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle2:489
INFO [06-29|22:20:21.042] UpdateNonClosingBalanceProof2   complete.. gasLimit=66794,gasUsed=66794 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle2:512
INFO [06-29|22:20:29.068] OpenChannel complete.. gasLimit=135688,gasUsed=135688 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|22:20:37.111] Deposit2 complete.. gasLimit=68841,gasUsed=68841 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|22:20:41.146] CooperativeSettle   complete.. gasLimit=89793,gasUsed=179585 fn=contracts_test.go:TestCooperateSettleChannel:577
INFO [06-29|22:20:49.172] OpenChannel complete.. gasLimit=135688,gasUsed=135688 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|22:20:57.207] Deposit2 complete.. gasLimit=68841,gasUsed=68841 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|22:21:01.246] WithDraw complete.. gasLimit=218513,gasUsed=308513 fn=contracts_test.go:TestWithdraw:702
INFO [06-29|22:21:05.257] RegisterSecret success.. gasLimit=24443,gasUsed=24443 fn=contracts_test.go:TestRegisterSecret:732
INFO [06-29|22:21:13.287] OpenChannel complete.. gasLimit=135688,gasUsed=135688 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|22:21:21.326] Deposit2 complete.. gasLimit=68841,gasUsed=68841 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|22:21:24.35] close channel successful,gasused=99334,gasLimit=99334 fn=contracts_test.go:TestUnlock:783
INFO [06-29|22:21:41.422] UpdateNonClosingBalanceProof successful,gasused=113859,gasLimit=113859,locksroot=625930f371d5582539d0194b5befd5dcd409823c4be37763fd850c3884c73c1b,transferamount=10 fn=contracts_test.go:TestUnlock:828
INFO [06-29|22:21:45.441] unlock success,gasUsed=69715,gasLimit=69715,txhash=0x23aaa5a4816e14055fae0bc16a9d61149901abcebf362e0880930dc115f8401a fn=contracts_test.go:TestUnlock:846
INFO [06-29|22:22:31.636] settle channel complete ,gasused=79274,gasLimit=158547 fn=contracts_test.go:TestUnlock:894
INFO [06-29|22:22:39.662] OpenChannel complete.. gasLimit=135688,gasUsed=135688 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|22:22:47.694] Deposit2 complete.. gasLimit=68841,gasUsed=68841 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|22:22:51.715] close channel successful,gasused=99334,gasLimit=99334 fn=contracts_test.go:TestPunishObsoleteUnlock:975
INFO [06-29|22:23:07.78] UpdateNonClosingBalanceProof successful,gasused=113812,gasLimit=113812,locksroot=c6dae4c100fb84dbae2c5ea5686be48182d377f8244a2c8967d7971547de3912,transferamount=10 fn=contracts_test.go:TestPunishObsoleteUnlock:1020
INFO [06-29|22:23:11.8] PunishObsoleteUnlock success,gasUsed=34384,gasLimit=49384,txhash=0xdf111f7607cb6fbfe982bd9e56a3879bd959a1146b7091b1ba26d6a25708aa49 fn=contracts_test.go:TestPunishObsoleteUnlock:1050
