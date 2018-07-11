#花费 gas 统计
* OpenChannel 48938
* Deposit 70228
* OpenChannelAndDeposit 91320
* CloseChannel 无证据:32920
* CloseChannel 有证据:66547
* UpdateBalanceProofDelegate:75281
* updateBalanceProof 61633
* settle channel:50975(实际),101950(估计)
* CooperativeSettle: 77088,122088
* withdraw:114035
* unlock:80591(3个锁)
* punish: 41178,56178

## 记录
INFO [07-04|12:20:44.077] OpenChannel gasLimit=48938,gasUsed=48938 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-04|12:20:53.119] Deposit complete...,gasLimit=70228,gasUsed=70228 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-04|12:21:00.152] OpenChannelWithDeposit gasLimit=91383,gasUsed=91383 fn=contracts_test.go:creatAChannelAndDeposit2:206
INFO [07-04|12:21:05.175] Deposit2 complete...,gasLimit=70228,gasUsed=70228 fn=contracts_test.go:creatAChannelAndDeposit2:229
INFO [07-04|12:21:12.207] OpenChannel gasLimit=48938,gasUsed=48938 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-04|12:21:20.254] Deposit complete...,gasLimit=70228,gasUsed=70228 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-04|12:21:24.27] CloseChannel no evidence gasLimit=32920,gasUsed=32920 fn=contracts_test.go:TestCloseChannel1:333
INFO [07-04|12:21:33.309] OpenChannel gasLimit=48938,gasUsed=48938 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-04|12:21:41.355] Deposit complete...,gasLimit=70228,gasUsed=70228 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-04|12:21:45.383] CloseChannel with evidence gasLimit=66547,gasUsed=66547 fn=contracts_test.go:TestCloseChannel2:456
INFO [07-04|12:21:53.41] OpenChannel gasLimit=48938,gasUsed=48938 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-04|12:22:01.442] Deposit complete...,gasLimit=70228,gasUsed=70228 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-04|12:22:08.511] UpdateBalanceProofDelegate gasLimit=75281,gasUsed=75281 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofDelegateAndSettle:543
INFO [07-04|12:23:30.861] SettleChannel gasLimit=101950,gasUsed=50975 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofDelegateAndSettle:597
INFO [07-04|12:23:38.884] OpenChannel gasLimit=48938,gasUsed=48938 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-04|12:23:46.925] Deposit complete...,gasLimit=70228,gasUsed=70228 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-04|12:23:58.994] UpdateBalanceProof gasLimit=61633,gasUsed=61633 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofAndSettle:654
INFO [07-04|12:25:21.343] SettleChannel gasLimit=101950,gasUsed=50975 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofAndSettle:702
INFO [07-04|12:25:29.368] OpenChannel gasLimit=48938,gasUsed=48938 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-04|12:25:37.419] Deposit complete...,gasLimit=70228,gasUsed=70228 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-04|12:25:41.466] CooperativeSettle gasLimit=122088,gasUsed=77088 fn=contracts_test.go:TestCooperateSettleChannel:770
INFO [07-04|12:25:52.528] OpenChannel gasLimit=48944,gasUsed=48944 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-04|12:26:00.583] Deposit complete...,gasLimit=70222,gasUsed=70222 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-04|12:26:09.64] close channel successful,gasused=66472,gasLimit=66472 fn=contracts_test.go:TestUnlock:860
INFO [07-04|12:26:24.72] UpdateBalanceProof successful,gasused=63696,gasLimit=63696 fn=contracts_test.go:TestUnlock:910
INFO [07-04|12:26:29.75] unlock success,gasUsed=80591,gasLimit=80591,txhash=0xacda2fa5674b0cd7a8815e077d89bfad36cb651d66219a29baa34ba824b70ff3 fn=contracts_test.go:TestUnlock:937
INFO [07-04|12:27:35.191] settle channel complete ,gasused=51992,gasLimit=103984 fn=contracts_test.go:TestUnlock:981
INFO [07-04|12:27:43.224] OpenChannel gasLimit=48938,gasUsed=48938 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-04|12:27:51.266] Deposit complete...,gasLimit=70228,gasUsed=70228 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-04|12:27:55.315] WithDraw complete.. gasLimit=114035,gasUsed=114035 fn=contracts_test.go:TestWithdraw:1118
INFO [07-04|12:28:03.347] OpenChannel gasLimit=48938,gasUsed=48938 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-04|12:28:11.39] Deposit complete...,gasLimit=70228,gasUsed=70228 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-04|12:28:15.413] close channel successful,gasused=66530,gasLimit=66530 fn=contracts_test.go:TestPunishObsoleteUnlock:1208
INFO [07-04|12:28:31.495] UpdateBalanceProofDelegate successful,gasused=77321,gasLimit=77321 fn=contracts_test.go:TestPunishObsoleteUnlock:1248
INFO [07-04|12:28:35.523] PunishObsoleteUnlock success,gasUsed=41178,gasLimit=56178,txhash=0xb0cd8fc3f92b3b64dd9c5ec23ef19e77b41ab62dbc61d7ca7e33009d4c8bacbd fn=contracts_test.go:TestPunishObsoleteUnlock:1283
