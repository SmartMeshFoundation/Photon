#花费 gas 统计
* Approve 30418
* OpenChannel 48913
* Deposit 69956
* Deposit tokenFallback 52572
* Deposit ApproveAndCall 67702,83097
* OpenChannelAndDeposit 91618
* OpenChannelAndDeposit tokenFallback 88910
* OpenChannelAndDeposit ApproveAndCall 104098
* CloseChannel 无证据:32891
* CloseChannel 有证据:66694
* UpdateBalanceProofDelegate:75791
* updateBalanceProof 62046
* settle channel:51674(实际),103347(估计)
* CooperativeSettle: 77570,122570
* withdraw:114905
* unlock:67664 单个锁
* punish: 39064,69064

## 记录
	contracts_test.go:248: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0x7314c3E027d1AA6AB9dCb40A39b4e57659e44701 for 0xDF638ec99EeeF249Ffa68aadB4E3B8A7121B8541,gasUsed=45418,gasLimit=45418
INFO [08-03|11:47:54.886] OpenChannel gasLimit=48913,gasUsed=48913 fn=contracts_test.go:creatAChannelAndDeposit:152
INFO [08-03|11:47:58.918] Deposit complete...,gasLimit=69956,gasUsed=69956 fn=contracts_test.go:creatAChannelAndDeposit:188
	contracts_test.go:248: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0x7314c3E027d1AA6AB9dCb40A39b4e57659e44701 for 0xDF638ec99EeeF249Ffa68aadB4E3B8A7121B8541,gasUsed=30418,gasLimit=30418
INFO [08-03|11:48:02.949] OpenChannelWithDeposit gasLimit=91618,gasUsed=91618 fn=contracts_test.go:creatAChannelAndDeposit2:207
INFO [08-03|11:48:04.97] Deposit2 complete...,gasLimit=69956,gasUsed=69956 fn=contracts_test.go:creatAChannelAndDeposit2:230
	contracts_test.go:248: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0x7314c3E027d1AA6AB9dCb40A39b4e57659e44701 for 0xDF638ec99EeeF249Ffa68aadB4E3B8A7121B8541,gasUsed=30418,gasLimit=30418
INFO [08-03|11:48:09.996] OpenChannel gasLimit=48913,gasUsed=48913 fn=contracts_test.go:creatAChannelAndDeposit:152
INFO [08-03|11:48:16.025] Deposit complete...,gasLimit=69956,gasUsed=69956 fn=contracts_test.go:creatAChannelAndDeposit:188
INFO [08-03|11:48:19.039] CloseChannel no evidence gasLimit=32891,gasUsed=32891 fn=contracts_test.go:TestCloseChannel1:334
	contracts_test.go:248: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0x7314c3E027d1AA6AB9dCb40A39b4e57659e44701 for 0xDF638ec99EeeF249Ffa68aadB4E3B8A7121B8541,gasUsed=30418,gasLimit=30418
INFO [08-03|11:48:25.067] OpenChannel gasLimit=48913,gasUsed=48913 fn=contracts_test.go:creatAChannelAndDeposit:152
INFO [08-03|11:48:31.101] Deposit complete...,gasLimit=69956,gasUsed=69956 fn=contracts_test.go:creatAChannelAndDeposit:188
INFO [08-03|11:48:34.124] CloseChannel with evidence gasLimit=66694,gasUsed=66694 fn=contracts_test.go:TestCloseChannel2:448
	contracts_test.go:248: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0x7314c3E027d1AA6AB9dCb40A39b4e57659e44701 for 0xDF638ec99EeeF249Ffa68aadB4E3B8A7121B8541,gasUsed=30418,gasLimit=30418
INFO [08-03|11:48:40.146] OpenChannel gasLimit=48913,gasUsed=48913 fn=contracts_test.go:creatAChannelAndDeposit:152
INFO [08-03|11:48:46.174] Deposit complete...,gasLimit=69956,gasUsed=69956 fn=contracts_test.go:creatAChannelAndDeposit:188
INFO [08-03|11:49:06.276] UpdateBalanceProofDelegate gasLimit=75791,gasUsed=75791 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofDelegateAndSettle:561
INFO [08-03|11:49:28.367] SettleChannel gasLimit=103347,gasUsed=51674 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofDelegateAndSettle:607
	contracts_test.go:248: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0x7314c3E027d1AA6AB9dCb40A39b4e57659e44701 for 0xDF638ec99EeeF249Ffa68aadB4E3B8A7121B8541,gasUsed=30418,gasLimit=30418
INFO [08-03|11:49:33.39] OpenChannel gasLimit=48913,gasUsed=48913 fn=contracts_test.go:creatAChannelAndDeposit:152
INFO [08-03|11:49:39.419] Deposit complete...,gasLimit=69956,gasUsed=69956 fn=contracts_test.go:creatAChannelAndDeposit:188
INFO [08-03|11:49:39.422] gasLimit=21000,gasPrice=18000000000      fn=contracts_test.go:TransferTo:117
INFO [08-03|11:49:48.468] UpdateBalanceProof gasLimit=62046,gasUsed=62046 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofAndSettle:664
INFO [08-03|11:50:23.607] SettleChannel gasLimit=103347,gasUsed=51674 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofAndSettle:712
	contracts_test.go:248: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0x7314c3E027d1AA6AB9dCb40A39b4e57659e44701 for 0xDF638ec99EeeF249Ffa68aadB4E3B8A7121B8541,gasUsed=30418,gasLimit=30418
INFO [08-03|11:50:29.626] OpenChannel gasLimit=48913,gasUsed=48913 fn=contracts_test.go:creatAChannelAndDeposit:152
INFO [08-03|11:50:35.661] Deposit complete...,gasLimit=69956,gasUsed=69956 fn=contracts_test.go:creatAChannelAndDeposit:188
INFO [08-03|11:50:37.695] CooperativeSettle gasLimit=122570,gasUsed=77570 fn=contracts_test.go:TestCooperateSettleChannel:780
	contracts_test.go:248: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0x7314c3E027d1AA6AB9dCb40A39b4e57659e44701 for 0xDF638ec99EeeF249Ffa68aadB4E3B8A7121B8541,gasUsed=30418,gasLimit=30418
INFO [08-03|11:50:46.729] OpenChannel gasLimit=48913,gasUsed=48913 fn=contracts_test.go:creatAChannelAndDeposit:152
INFO [08-03|11:50:52.762] Deposit complete...,gasLimit=69956,gasUsed=69956 fn=contracts_test.go:creatAChannelAndDeposit:188
INFO [08-03|11:50:52.766] gasLimit=21000,gasPrice=18000000000      fn=contracts_test.go:TransferTo:117
INFO [08-03|11:50:57.79] close channel successful,gasused=66758,gasLimit=66758 fn=contracts_test.go:TestUnlock:870
INFO [08-03|11:51:08.846] UpdateBalanceProof successful,gasused=64103,gasLimit=64103 fn=contracts_test.go:TestUnlock:920
INFO [08-03|11:51:10.863] unlock success,gasUsed=67664,gasLimit=67664,txhash=0x4c21fbc6d9e18e1b7607e6a31b5b7c9a8b82421a38c2ad52de3bdf481ee3d5c3 fn=contracts_test.go:TestUnlock:957
INFO [08-03|11:51:34.971] settle channel complete ,gasused=52694,gasLimit=105387 fn=contracts_test.go:TestUnlock:1001
	contracts_test.go:248: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0x7314c3E027d1AA6AB9dCb40A39b4e57659e44701 for 0xDF638ec99EeeF249Ffa68aadB4E3B8A7121B8541,gasUsed=30418,gasLimit=30418
INFO [08-03|11:51:38.992] OpenChannel gasLimit=48913,gasUsed=48913 fn=contracts_test.go:creatAChannelAndDeposit:152
INFO [08-03|11:51:45.022] Deposit complete...,gasLimit=69956,gasUsed=69956 fn=contracts_test.go:creatAChannelAndDeposit:188
INFO [08-03|11:51:48.053] WithDraw complete.. gasLimit=114905,gasUsed=114905 fn=contracts_test.go:TestWithdraw:1138
	contracts_test.go:248: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0x7314c3E027d1AA6AB9dCb40A39b4e57659e44701 for 0xDF638ec99EeeF249Ffa68aadB4E3B8A7121B8541,gasUsed=30418,gasLimit=30418
INFO [08-03|11:51:54.082] OpenChannel gasLimit=48855,gasUsed=48855 fn=contracts_test.go:creatAChannelAndDeposit:152
INFO [08-03|11:52:00.112] Deposit complete...,gasLimit=69886,gasUsed=69886 fn=contracts_test.go:creatAChannelAndDeposit:188
INFO [08-03|11:52:00.116] gasLimit=21000,gasPrice=18000000000      fn=contracts_test.go:TransferTo:117
INFO [08-03|11:52:06.141] close channel successful,gasused=66683,gasLimit=66683 fn=contracts_test.go:TestPunishObsoleteUnlock:1261
INFO [08-03|11:52:18.194] UpdateBalanceProofDelegate successful,gasused=64092,gasLimit=64092 fn=contracts_test.go:TestPunishObsoleteUnlock:1302
INFO [08-03|11:52:21.219] unlockdelegate gasLimit=81159,gasUsed=81159 fn=contracts_test.go:TestPunishObsoleteUnlock:1352
INFO [08-03|11:52:24.24] PunishObsoleteUnlock success,gasUsed=39064,gasLimit=69064,txhash=0xb40b5bdf3f47db32d09bd4491ae5a7beb24570457f137901824f36ff6efc1aed fn=contracts_test.go:TestPunishObsoleteUnlock:1387
	contracts_test.go:248: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0x7314c3E027d1AA6AB9dCb40A39b4e57659e44701 for 0xDF638ec99EeeF249Ffa68aadB4E3B8A7121B8541,gasUsed=30418,gasLimit=30418
INFO [08-03|11:52:27.262] open channel and deposit by tokenFallback success,gasUsed=88910,gasLimit=88910,txhash=0x75a123813063755884f209d7932cd72a390415bc9681462895ab124e99e75cda fn=contracts_test.go:testOpenChannelAndDepositFallback:1450
INFO [08-03|11:52:30.279] deposit by tokenFallback success,gasUsed=52572,gasLimit=52572,txhash=0x91b2ab383284fc5dec5b8828df83720c2c83f23cee174e326edcac92081686ff fn=contracts_test.go:testDepositFallback:1474
INFO [08-03|11:52:33.304] open channel and deposit by ApproveAndCall success,gasUsed=104098,gasLimit=120065,txhash=0x01fbcc7caa6064e2a1bd02c542f94d157bb4dc4f98e492946a005a07d805ccf5 fn=contracts_test.go:testOpenChannelAndDepositApproveCall:1531
INFO [08-03|11:52:35.323] deposit by ApproveAndCall success,gasUsed=67702,gasLimit=83097,txhash=0xb77104c7f77106b416d60ac26b8278a42ad5db814c26aee055244671fa6fab97 fn=contracts_test.go:testDepositApproveCall:1555
