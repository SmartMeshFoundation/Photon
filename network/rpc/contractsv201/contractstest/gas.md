#花费 gas 统计
* OpenChannel 49199
* Deposit 70546
* OpenChannelAndDeposit 91622
* CloseChannel 无证据:33556
* CloseChannel 有证据:66919
* UpdateBalanceProofDelegate:76624
* updateBalanceProof 62955
* settle channel:56597(实际),113193(估计)
* CooperativeSettle: 77746,122746
* withdraw:115055
* unlock:81037(3个锁)
* punish: 47723,62723



## 记录
INFO [07-03|17:23:00.985] OpenChannel gasLimit=49199,gasUsed=49199 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-03|17:23:09.028] Deposit complete...,gasLimit=70546,gasUsed=70546 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-03|17:23:17.068] OpenChannelWithDeposit gasLimit=91622,gasUsed=91622 fn=contracts_test.go:creatAChannelAndDeposit2:206
INFO [07-03|17:23:21.092] Deposit2 complete...,gasLimit=70546,gasUsed=70546 fn=contracts_test.go:creatAChannelAndDeposit2:229
INFO [07-03|17:23:29.119] OpenChannel gasLimit=49199,gasUsed=49199 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-03|17:23:37.158] Deposit complete...,gasLimit=70546,gasUsed=70546 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-03|17:23:41.176] CloseChannel no evidence gasLimit=33556,gasUsed=33556 fn=contracts_test.go:TestCloseChannel1:333
INFO [07-03|17:23:49.205] OpenChannel gasLimit=49199,gasUsed=49199 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-03|17:23:57.249] Deposit complete...,gasLimit=70546,gasUsed=70546 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-03|17:24:01.28] CloseChannel with evidence gasLimit=66919,gasUsed=66919 fn=contracts_test.go:TestCloseChannel2:457
INFO [07-03|17:24:08.311] OpenChannel gasLimit=49199,gasUsed=49199 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-03|17:24:17.357] Deposit complete...,gasLimit=70546,gasUsed=70546 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-03|17:24:25.432] UpdateBalanceProofDelegate gasLimit=76624,gasUsed=76624 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofDelegateAndSettle:544
INFO [07-03|17:25:47.73] SettleChannel gasLimit=113193,gasUsed=56597 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofDelegateAndSettle:600
INFO [07-03|17:25:54.756] OpenChannel gasLimit=49199,gasUsed=49199 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-03|17:26:02.793] Deposit complete...,gasLimit=70546,gasUsed=70546 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-03|17:26:14.861] UpdateBalanceProof gasLimit=62955,gasUsed=62955 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofAndSettle:657
INFO [07-03|17:27:37.164] SettleChannel gasLimit=113193,gasUsed=56597 fn=contracts_test.go:TestCloseChannelAndUpdateBalanceProofAndSettle:707
INFO [07-03|17:27:45.197] OpenChannel gasLimit=49199,gasUsed=49199 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-03|17:27:53.241] Deposit complete...,gasLimit=70546,gasUsed=70546 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-03|17:27:56.298] CooperativeSettle gasLimit=122746,gasUsed=77746 fn=contracts_test.go:TestCooperateSettleChannel:775
INFO [07-03|17:28:09.344] OpenChannel gasLimit=49199,gasUsed=49199 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-03|17:28:17.392] Deposit complete...,gasLimit=70546,gasUsed=70546 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-03|17:28:24.437] close channel successful,gasused=66936,gasLimit=66936 fn=contracts_test.go:TestUnlock:865
INFO [07-03|17:28:40.526] UpdateBalanceProof successful,gasused=64956,gasLimit=64956 fn=contracts_test.go:TestUnlock:918
INFO [07-03|17:28:45.562] unlock success,gasUsed=81037,gasLimit=81037,txhash=0xa226592d3b4bce8c8a7652fc16872d12dcefcfc1d7bb37660df4cb00a242c7b4 fn=contracts_test.go:TestUnlock:945
INFO [07-03|17:29:50.906] settle channel complete ,gasused=57621,gasLimit=115241 fn=contracts_test.go:TestUnlock:991
INFO [07-03|17:29:58.932] OpenChannel gasLimit=49199,gasUsed=49199 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-03|17:30:06.968] Deposit complete...,gasLimit=70546,gasUsed=70546 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-03|17:30:11.019] WithDraw complete.. gasLimit=115055,gasUsed=115055 fn=contracts_test.go:TestWithdraw:1128
INFO [07-03|17:30:19.052] OpenChannel gasLimit=49199,gasUsed=49199 fn=contracts_test.go:creatAChannelAndDeposit:151
INFO [07-03|17:30:27.098] Deposit complete...,gasLimit=70546,gasUsed=70546 fn=contracts_test.go:creatAChannelAndDeposit:187
INFO [07-03|17:30:31.127] close channel successful,gasused=66855,gasLimit=66855 fn=contracts_test.go:TestPunishObsoleteUnlock:1219
INFO [07-03|17:30:47.208] UpdateBalanceProofDelegate successful,gasused=78689,gasLimit=78689 fn=contracts_test.go:TestPunishObsoleteUnlock:1259
INFO [07-03|17:30:51.239] PunishObsoleteUnlock success,gasUsed=47723,gasLimit=62723,txhash=0x67d840965e36ef139f31bf1b0499ffcf8987a4187bf03bb508619ac110b03a3b fn=contracts_test.go:TestPunishObsoleteUnlock:1298
