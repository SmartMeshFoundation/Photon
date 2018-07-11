#花费 gas 统计
* OpenChannel 48931
* Deposit 69841
* OpenChannelAndDeposit 91381
* CloseChannel 无证据:32973
* CloseChannel 有证据:66583
* UpdateBalanceProofDelegate:75629
* updateBalanceProof 61831
* settle channel:51019(实际),102037(估计)
* CooperativeSettle: 77140,122140
* withdraw:114450
* unlock:66095(一个一个解锁)
* unlockdelegate:79536
* punish: 41229,56229

## 记录
INFO [07-10|12:03:36.963] OpenChannel gasLimit=48995,gasUsed=48995 fn=contracts_test.go:creatAChannelAndDeposit:152
INFO [07-10|12:03:45.125] Deposit complete...,gasLimit=69905,gasUsed=69905 fn=contracts_test.go:creatAChannelAndDeposit:188
INFO [07-10|12:03:53.261] OpenChannelWithDeposit gasLimit=91375,gasUsed=91375 fn=contracts_test.go:creatAChannelAndDeposit2:207
INFO [07-10|12:03:57.347] Deposit2 complete...,gasLimit=69905,gasUsed=69905 fn=contracts_test.go:creatAChannelAndDeposit2:230
INFO [07-10|12:04:05.479] OpenChannel gasLimit=48995,gasUsed=48995 fn=contracts_test.go:creatAChannelAndDeposit:152
INFO [07-10|12:04:13.642] Deposit complete...,gasLimit=69905,gasUsed=69905 fn=contracts_test.go:creatAChannelAndDeposit:188
INFO [07-10|12:04:17.721] CloseChannel no evidence gasLimit=32973,gasUsed=32973 fn=contracts_test.go:TestCloseChannel1:334
INFO [07-10|12:04:24.851] OpenChannel gasLimit=48995,gasUsed=48995 fn=contracts_test.go:creatAChannelAndDeposit:152
INFO [07-10|12:04:33.025] Deposit complete...,gasLimit=69905,gasUsed=69905 fn=contracts_test.go:creatAChannelAndDeposit:188
INFO [07-10|12:04:37.131] CloseChannel with evidence gasLimit=66600,gasUsed=66600 fn=contracts_test.go:TestCloseChannel2:448
INFO [07-10|12:04:45.263] OpenChannel gasLimit=49001,gasUsed=49001 fn=contracts_test.go:creatAChannelAndDeposit:152
INFO [07-10|12:04:53.438] Deposit complete...,gasLimit=69899,gasUsed=69899 fn=contracts_test.go:creatAChannelAndDeposit:188
INFO [07-10|12:05:32.951] UpdateBalanceProofDelegate gasLimit=75588,gasUsed=75588 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofDelegateAndSettle:561
INFO [07-10|12:06:23.521] SettleChannel gasLimit=102031,gasUsed=51016 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofDelegateAndSettle:607
INFO [07-10|12:06:31.653] OpenChannel gasLimit=48995,gasUsed=48995 fn=contracts_test.go:creatAChannelAndDeposit:152
INFO [07-10|12:06:38.818] Deposit complete...,gasLimit=69905,gasUsed=69905 fn=contracts_test.go:creatAChannelAndDeposit:188
INFO [07-10|12:06:39.072] gasLimit=21000,gasPrice=18000000000      fn=contracts_test.go:TransferTo:117
INFO [07-10|12:06:51.267] UpdateBalanceProof gasLimit=61842,gasUsed=61842 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofAndSettle:664
INFO [07-10|12:08:13.2] SettleChannel gasLimit=102037,gasUsed=51019 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofAndSettle:712
INFO [07-10|12:08:21.332] OpenChannel gasLimit=48995,gasUsed=48995 fn=contracts_test.go:creatAChannelAndDeposit:152
INFO [07-10|12:08:29.495] Deposit complete...,gasLimit=69905,gasUsed=69905 fn=contracts_test.go:creatAChannelAndDeposit:188
INFO [07-10|12:08:33.653] CooperativeSettle gasLimit=122123,gasUsed=77123 fn=contracts_test.go:TestCooperateSettleChannel:780
INFO [07-10|12:08:44.862] OpenChannel gasLimit=48995,gasUsed=48995 fn=contracts_test.go:creatAChannelAndDeposit:152
INFO [07-10|12:08:53.026] Deposit complete...,gasLimit=69905,gasUsed=69905 fn=contracts_test.go:creatAChannelAndDeposit:188
INFO [07-10|12:08:53.115] gasLimit=21000,gasPrice=18000000000      fn=contracts_test.go:TransferTo:117
INFO [07-10|12:09:02.262] close channel successful,gasused=66536,gasLimit=66536 fn=contracts_test.go:TestUnlock:870
INFO [07-10|12:09:19.575] UpdateBalanceProof successful,gasused=63801,gasLimit=63801 fn=contracts_test.go:TestUnlock:920
INFO [07-10|12:09:23.662] unlock success,gasUsed=66159,gasLimit=66159,txhash=0x05f77eaff1eb760a57402e898ccc14108e72677844e494fc2f90076176516908 fn=contracts_test.go:TestUnlock:957
INFO [07-10|12:10:27.38] settle channel complete ,gasused=52039,gasLimit=104077 fn=contracts_test.go:TestUnlock:1001
INFO [07-10|12:10:35.516] OpenChannel gasLimit=48995,gasUsed=48995 fn=contracts_test.go:creatAChannelAndDeposit:152
INFO [07-10|12:10:43.683] Deposit complete...,gasLimit=69905,gasUsed=69905 fn=contracts_test.go:creatAChannelAndDeposit:188
INFO [07-10|12:10:46.798] WithDraw complete.. gasLimit=114433,gasUsed=114433 fn=contracts_test.go:TestWithdraw:1138
INFO [07-10|12:10:56.922] OpenChannel gasLimit=48995,gasUsed=48995 fn=contracts_test.go:creatAChannelAndDeposit:152
INFO [07-10|12:11:05.077] Deposit complete...,gasLimit=69905,gasUsed=69905 fn=contracts_test.go:creatAChannelAndDeposit:188
INFO [07-10|12:11:05.119] gasLimit=21000,gasPrice=18000000000      fn=contracts_test.go:TransferTo:117
INFO [07-10|12:11:13.241] close channel successful,gasused=66600,gasLimit=66600 fn=contracts_test.go:TestPunishObsoleteUnlock:1261
INFO [07-10|12:11:32.534] UpdateBalanceProofDelegate successful,gasused=63865,gasLimit=63865 fn=contracts_test.go:TestPunishObsoleteUnlock:1302
INFO [07-10|12:11:37.616] unlockdelegate gasLimit=79093,gasUsed=79093 fn=contracts_test.go:TestPunishObsoleteUnlock:1352
INFO [07-10|12:11:41.712] PunishObsoleteUnlock success,gasUsed=41357,gasLimit=56357,txhash=0x4340a85433eee059f9c3f25781e142b81d36df359e7c5e78b6d0e7c94fabb3d6 fn=contracts_test.go:TestPunishObsoleteUnlock:1387
