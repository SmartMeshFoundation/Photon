#花费 gas 统计
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
INFO [07-02|15:46:50.845] OpenChannel gasLimit=49100,gasUsed=49100 fn=contracts_test.go:creatAChannelAndDeposit:113
INFO [07-02|15:47:00.89] Deposit complete...,gasLimit=71005,gasUsed=71005 fn=contracts_test.go:creatAChannelAndDeposit:149
INFO [07-02|15:47:08.919] OpenChannel gasLimit=49100,gasUsed=49100 fn=contracts_test.go:creatAChannelAndDeposit:113
INFO [07-02|15:47:16.966] Deposit complete...,gasLimit=71005,gasUsed=71005 fn=contracts_test.go:creatAChannelAndDeposit:149
INFO [07-02|15:47:20.983] CloseChannel no evidence gasLimit=32776,gasUsed=32776 fn=contracts_test.go:TestCloseChannel1:224
INFO [07-02|15:47:29.009] OpenChannel gasLimit=49100,gasUsed=49100 fn=contracts_test.go:creatAChannelAndDeposit:113
INFO [07-02|15:47:39.051] Deposit complete...,gasLimit=71005,gasUsed=71005 fn=contracts_test.go:creatAChannelAndDeposit:149
INFO [07-02|15:47:43.083] CloseChannel with evidence gasLimit=66161,gasUsed=66161 fn=contracts_test.go:TestCloseChannel2:348
INFO [07-02|15:47:51.113] OpenChannel gasLimit=49100,gasUsed=49100 fn=contracts_test.go:creatAChannelAndDeposit:113
INFO [07-02|15:48:00.158] Deposit complete...,gasLimit=71005,gasUsed=71005 fn=contracts_test.go:creatAChannelAndDeposit:149
INFO [07-02|15:48:09.223] UpdateNonClosingBalanceProof gasLimit=75801,gasUsed=75801 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle:435
INFO [07-02|15:49:31.542] SettleChannel gasLimit=113474,gasUsed=56737 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle:482
INFO [07-02|15:49:39.575] OpenChannel gasLimit=49100,gasUsed=49100 fn=contracts_test.go:creatAChannelAndDeposit:113
INFO [07-02|15:49:46.618] Deposit complete...,gasLimit=71005,gasUsed=71005 fn=contracts_test.go:creatAChannelAndDeposit:149
INFO [07-02|15:49:51.655] CooperativeSettle gasLimit=122719,gasUsed=77719 fn=contracts_test.go:TestCooperateSettleChannel:550
INFO [07-02|15:50:03.702] OpenChannel gasLimit=49100,gasUsed=49100 fn=contracts_test.go:creatAChannelAndDeposit:113
INFO [07-02|15:50:10.744] Deposit complete...,gasLimit=71005,gasUsed=71005 fn=contracts_test.go:creatAChannelAndDeposit:149
INFO [07-02|15:50:14.777] close channel successful,gasused=66097,gasLimit=66097 fn=contracts_test.go:TestUnlock:634
INFO [07-02|15:50:30.86] UpdateNonClosingBalanceProof successful,gasused=77866,gasLimit=77866 fn=contracts_test.go:TestUnlock:689
INFO [07-02|15:50:34.888] unlock success,gasUsed=80604,gasLimit=80604,txhash=0x9d05916e8ef3eb60da4c3c57a4a2a5dfb5d9bd35700a51404f4c5f5cf455d96b fn=contracts_test.go:TestUnlock:716
INFO [07-02|15:51:40.174] settle channel complete ,gasused=57761,gasLimit=115522 fn=contracts_test.go:TestUnlock:762
INFO [07-02|15:51:49.208] OpenChannel gasLimit=49100,gasUsed=49100 fn=contracts_test.go:creatAChannelAndDeposit:113
INFO [07-02|15:51:57.251] Deposit complete...,gasLimit=71005,gasUsed=71005 fn=contracts_test.go:creatAChannelAndDeposit:149
INFO [07-02|15:52:00.301] WithDraw complete.. gasLimit=115181,gasUsed=115181 fn=contracts_test.go:TestWithdraw:899
INFO [07-02|15:52:09.334] OpenChannel gasLimit=49100,gasUsed=49100 fn=contracts_test.go:creatAChannelAndDeposit:113
INFO [07-02|15:52:17.376] Deposit complete...,gasLimit=71005,gasUsed=71005 fn=contracts_test.go:creatAChannelAndDeposit:149
INFO [07-02|15:52:21.408] close channel successful,gasused=66144,gasLimit=66144 fn=contracts_test.go:TestPunishObsoleteUnlock:990
INFO [07-02|15:52:37.494] UpdateNonClosingBalanceProof successful,gasused=77930,gasLimit=77930 fn=contracts_test.go:TestPunishObsoleteUnlock:1030
INFO [07-02|15:52:41.526] PunishObsoleteUnlock success,gasUsed=47669,gasLimit=62669,txhash=0xd7d82b691acf8312210b7fa32cdadeab01df2a801bd7dd1495ab7d9dea4fb996 fn=contracts_test.go:TestPunishObsoleteUnlock:1061
