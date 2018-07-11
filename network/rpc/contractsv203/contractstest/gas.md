#花费 gas 统计
* OpenChannel 488897
* Deposit 70250
* OpenChannelAndDeposit 91320
* CloseChannel 无证据:33260
* CloseChannel 有证据:66474
* UpdateBalanceProofDelegate:75966
* updateBalanceProof 62493
* settle channel:56597(实际),113193(估计)
* CooperativeSettle: 76970,121970
* withdraw:113903
* unlock:80725(3个锁)
* punish: 40663,55663



## 记录
INFO [07-03|17:51:05.154] OpenChannel gasLimit=48897,gasUsed=48897 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-03|17:51:13.197] Deposit complete...,gasLimit=70250,gasUsed=70250 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-03|17:51:21.233] OpenChannelWithDeposit gasLimit=91320,gasUsed=91320 fn=contracts_test.go:creatAChannelAndDeposit2:206
INFO [07-03|17:51:24.258] Deposit2 complete...,gasLimit=70250,gasUsed=70250 fn=contracts_test.go:creatAChannelAndDeposit2:229
INFO [07-03|17:51:33.287] OpenChannel gasLimit=48903,gasUsed=48903 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-03|17:51:41.328] Deposit complete...,gasLimit=70244,gasUsed=70244 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-03|17:51:45.344] CloseChannel no evidence gasLimit=33260,gasUsed=33260 fn=contracts_test.go:TestCloseChannel1:333
INFO [07-03|17:51:53.379] OpenChannel gasLimit=48897,gasUsed=48897 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-03|17:52:01.423] Deposit complete...,gasLimit=70250,gasUsed=70250 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-03|17:52:04.45] CloseChannel with evidence gasLimit=66474,gasUsed=66474 fn=contracts_test.go:TestCloseChannel2:457
INFO [07-03|17:52:13.478] OpenChannel gasLimit=48897,gasUsed=48897 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-03|17:52:21.525] Deposit complete...,gasLimit=70250,gasUsed=70250 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-03|17:52:29.593] UpdateBalanceProofDelegate gasLimit=75966,gasUsed=75966 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofDelegateAndSettle:544
INFO [07-03|17:53:50.923] SettleChannel gasLimit=112897,gasUsed=56449 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofDelegateAndSettle:600
INFO [07-03|17:53:58.956] OpenChannel gasLimit=48833,gasUsed=48833 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-03|17:54:07] Deposit complete...,gasLimit=70186,gasUsed=70186 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-03|17:54:19.073] UpdateBalanceProof gasLimit=62493,gasUsed=62493 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofAndSettle:657
INFO [07-03|17:55:41.375] SettleChannel gasLimit=112833,gasUsed=56417 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofAndSettle:707
INFO [07-03|17:55:49.403] OpenChannel gasLimit=48897,gasUsed=48897 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-03|17:55:57.453] Deposit complete...,gasLimit=70250,gasUsed=70250 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-03|17:56:01.494] CooperativeSettle gasLimit=121970,gasUsed=76970 fn=contracts_test.go:TestCooperateSettleChannel:775
INFO [07-03|17:56:13.544] OpenChannel gasLimit=48897,gasUsed=48897 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-03|17:56:21.587] Deposit complete...,gasLimit=70250,gasUsed=70250 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-03|17:56:28.642] close channel successful,gasused=66457,gasLimit=66457 fn=contracts_test.go:TestUnlock:865
INFO [07-03|17:56:45.727] UpdateBalanceProof successful,gasused=64558,gasLimit=64558 fn=contracts_test.go:TestUnlock:918
INFO [07-03|17:56:49.753] unlock success,gasUsed=80725,gasLimit=80725,txhash=0xc14d0d2bb0fd1029aa5f3ccf5c5850850f8badf9d4438cc133e8a154b1d28c3a fn=contracts_test.go:TestUnlock:945
INFO [07-03|17:57:55.173] settle channel complete ,gasused=57473,gasLimit=114945 fn=contracts_test.go:TestUnlock:991
INFO [07-03|17:58:03.211] OpenChannel gasLimit=48897,gasUsed=48897 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-03|17:58:11.294] Deposit complete...,gasLimit=70250,gasUsed=70250 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-03|17:58:15.348] WithDraw complete.. gasLimit=113903,gasUsed=113903 fn=contracts_test.go:TestWithdraw:1128
INFO [07-03|17:58:23.385] OpenChannel gasLimit=48903,gasUsed=48903 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-03|17:58:31.433] Deposit complete...,gasLimit=70244,gasUsed=70244 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-03|17:58:35.468] close channel successful,gasused=66463,gasLimit=66463 fn=contracts_test.go:TestPunishObsoleteUnlock:1219
INFO [07-03|17:58:51.565] UpdateBalanceProofDelegate successful,gasused=78020,gasLimit=78020 fn=contracts_test.go:TestPunishObsoleteUnlock:1259
INFO [07-03|17:58:55.6] PunishObsoleteUnlock success,gasUsed=47186,gasLimit=62186,txhash=0x8fbc3d761fd2e5f25870200d2616c2d9a77174219b2da12f50f04233d7a0c9f2 fn=contracts_test.go:TestPunishObsoleteUnlock:1298
