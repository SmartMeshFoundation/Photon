#花费 gas 统计
* OpenChannel 94386
* Deposit 71266
* CloseChannel 无证据:38162
* CloseChannel 有证据:96276
* UpdateNonClosingBalanceProof:99299
* settle channel:73789(实际),147578(估计)
* CooperativeSettle: 93575,168575
* withdraw:211953,271953
* unlock:71782(3个锁)
* PunishObsoleteUnlock:36999,66999


## 记录
INFO [06-30|13:27:59.064] OpenChannel complete.. gasLimit=94386,gasUsed=94386 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-30|13:28:07.112] Deposit2 complete.. gasLimit=71266,gasUsed=71266 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-30|13:28:15.145] OpenChannel complete.. gasLimit=94386,gasUsed=94386 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-30|13:28:22.188] Deposit2 complete.. gasLimit=71276,gasUsed=71276 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-30|13:28:27.205] CloseChannel no evidence complete.. gasLimit=38162,gasUsed=38162 fn=contracts_test.go:TestCloseChannel1:209
INFO [06-30|13:28:34.238] OpenChannel complete.. gasLimit=94386,gasUsed=94386 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-30|13:28:43.292] Deposit2 complete.. gasLimit=71266,gasUsed=71266 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-30|13:28:47.322] CloseChannel with evidence  complete.. gasLimit=96293,gasUsed=96293 fn=contracts_test.go:TestCloseChannel2:305
INFO [06-30|13:28:55.351] OpenChannel complete.. gasLimit=94386,gasUsed=94386 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-30|13:29:03.395] Deposit2 complete.. gasLimit=71266,gasUsed=71266 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-30|13:29:07.427] CloseChannel   complete.. gasLimit=96276,gasUsed=96276 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle:387
INFO [06-30|13:29:11.467] UpdateNonClosingBalanceProof   complete.. gasLimit=99299,gasUsed=99299 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle:403
INFO [06-30|13:30:12.736] SettleChannel  complete.. gasLimit=73789,gasUsed=147578 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle:444
INFO [06-30|13:30:21.77] OpenChannel complete.. gasLimit=94386,gasUsed=94386 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-30|13:30:28.816] Deposit2 complete.. gasLimit=71266,gasUsed=71266 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-30|13:30:32.846] CloseChannel   complete.. gasLimit=96212,gasUsed=96212 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle2:468
INFO [06-30|13:30:36.891] UpdateNonClosingBalanceProof   complete.. gasLimit=99299,gasUsed=99299 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle2:484
INFO [06-30|13:30:40.962] UpdateNonClosingBalanceProof2   complete.. gasLimit=69299,gasUsed=69299 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle2:507
INFO [06-30|13:30:48.989] OpenChannel complete.. gasLimit=94386,gasUsed=94386 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-30|13:30:57.038] Deposit2 complete.. gasLimit=71266,gasUsed=71266 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-30|13:31:01.087] CooperativeSettle   complete.. gasLimit=93575,gasUsed=168575 fn=contracts_test.go:TestCooperateSettleChannel:572
INFO [06-30|13:31:09.12] OpenChannel complete.. gasLimit=94322,gasUsed=94322 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-30|13:31:17.164] Deposit2 complete.. gasLimit=71202,gasUsed=71202 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-30|13:31:21.217] WithDraw complete.. gasLimit=211953,gasUsed=271953 fn=contracts_test.go:TestWithdraw:696
INFO [06-30|13:31:25.233] RegisterSecret success.. gasLimit=24443,gasUsed=24443 fn=contracts_test.go:TestRegisterSecret:726
INFO [06-30|13:31:33.27] OpenChannel complete.. gasLimit=94386,gasUsed=94386 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-30|13:31:41.325] Deposit2 complete.. gasLimit=71266,gasUsed=71266 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-30|13:31:45.358] close channel successful,gasused=96229,gasLimit=96229 fn=contracts_test.go:TestUnlock:777
INFO [06-30|13:32:01.439] UpdateNonClosingBalanceProof successful,gasused=116266,gasLimit=116266,locksroot=d13982fb2e9a11fdc6ae23134ac27c4492ffae91b6889f057177e5244ced5d55,transferamount=10 fn=contracts_test.go:TestUnlock:822
INFO [06-30|13:32:05.469] unlock success,gasUsed=71782,gasLimit=71782,txhash=0xe555cfddd84827e3f8df75128f2d12aacc0951f6c756955a11d01e63def74de5 fn=contracts_test.go:TestUnlock:840
INFO [06-30|13:32:51.633] settle channel complete ,gasused=73789,gasLimit=147578 fn=contracts_test.go:TestUnlock:888
INFO [06-30|13:32:59.668] OpenChannel complete.. gasLimit=94322,gasUsed=94322 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-30|13:33:07.713] Deposit2 complete.. gasLimit=71202,gasUsed=71202 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-30|13:33:11.754] close channel successful,gasused=96229,gasLimit=96229 fn=contracts_test.go:TestPunishObsoleteUnlock:969
INFO [06-30|13:33:27.838] UpdateNonClosingBalanceProof successful,gasused=116202,gasLimit=116202,locksroot=f3b452b400b9b5518055695d50a878e594e1a3230e4be6ca89b893095598e368,transferamount=10 fn=contracts_test.go:TestPunishObsoleteUnlock:1014
INFO [06-30|13:33:31.872] PunishObsoleteUnlock success,gasUsed=36999,gasLimit=66999,txhash=0xf8278b119a55d05a55e8217239d202ba821de632f3ba499a416eb241d7befe72 fn=contracts_test.go:TestPunishObsoleteUnlock:1044
