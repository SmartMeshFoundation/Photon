#花费 gas 统计
* OpenChannel 69099
* Deposit 73951
* CloseChannel 无证据:59502
* CloseChannel 有证据:114332
* UpdateNonClosingBalanceProof:95996
* settle channel:71334(实际),142668(估计)
* CooperativeSettle: 98614,158614
* withdraw:88926
* unlock:56879(3个锁),71879


* OpenChannel 49100
* Deposit 71005
* CloseChannel 无证据:32712
* CloseChannel 有证据:66097
* UpdateNonClosingBalanceProof:75882
* settle channel:57761(实际),115522(估计)
* CooperativeSettle: 77719,122719
* withdraw:115181
* unlock:80594(3个锁),80594
* punish: 47669,62669



## 记录
INFO [06-29|22:53:03.681] OpenChannel gasLimit=69099,gasUsed=69099 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|22:53:11.732] Deposit complete...,gasLimit=73951,gasUsed=73951 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|22:53:15.753] OpenChannel gasLimit=69099,gasUsed=69099 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|22:53:23.803] Deposit complete...,gasLimit=73951,gasUsed=73951 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|22:53:27.819] CloseChannel no evidence gasLimit=59502,gasUsed=59502 fn=contracts_test.go:TestCloseChannel1:204
INFO [06-29|22:53:30.835] OpenChannel gasLimit=69099,gasUsed=69099 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|22:53:38.878] Deposit complete...,gasLimit=73951,gasUsed=73951 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|22:53:42.906] CloseChannel with evidence gasLimit=114332,gasUsed=114332 fn=contracts_test.go:TestCloseChannel2:329
INFO [06-29|22:53:46.921] OpenChannel gasLimit=69099,gasUsed=69099 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|22:53:54.966] Deposit complete...,gasLimit=73951,gasUsed=73951 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|22:54:03.039] UpdateNonClosingBalanceProof gasLimit=95996,gasUsed=95996 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle:408
INFO [06-29|22:55:45.451] SettleChannel gasLimit=142668,gasUsed=71334 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle:455
INFO [06-29|22:55:49.465] OpenChannel gasLimit=69099,gasUsed=69099 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|22:55:57.5] Deposit complete...,gasLimit=73951,gasUsed=73951 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|22:56:01.532] CooperativeSettle gasLimit=158614,gasUsed=98614 fn=contracts_test.go:TestCooperateSettleChannel:525
INFO [06-29|22:56:05.545] OpenChannel gasLimit=69099,gasUsed=69099 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|22:56:13.581] Deposit complete...,gasLimit=73951,gasUsed=73951 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|22:56:17.632] withdraw gasLimit=88926,gasUsed=88926    fn=contracts_test.go:TestSetTotalWithdraw:602
INFO [06-29|22:56:25.658] OpenChannel gasLimit=69099,gasUsed=69099 fn=contracts_test.go:creatAChannelAndDeposit:101
INFO [06-29|22:56:33.692] Deposit complete...,gasLimit=73951,gasUsed=73951 fn=contracts_test.go:creatAChannelAndDeposit:137
INFO [06-29|22:56:37.714] close channel successful,gasused=114332,gasLimit=114332 fn=contracts_test.go:TestUnlock:687
INFO [06-29|22:56:52.773] UpdateNonClosingBalanceProof successful,gasused=95868,gasLimit=95868 fn=contracts_test.go:TestUnlock:727
INFO [06-29|22:58:23.135] settle channel complete ,gasused=82467,gasLimit=164933 fn=contracts_test.go:TestUnlock:775
INFO [06-29|22:58:27.16] unlock success,gasUsed=56879,gasLimit=71879,txhash=0x32b40a13946db192c348027cfe608a3b1a8f27fcc755db3dd216c7922389ce3d fn=contracts_test.go:TestUnlock:799
