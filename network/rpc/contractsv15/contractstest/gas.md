#花费 gas 统计
* OpenChannel 135682
* Deposit 68808
* CloseChannel 无证据:56048
* CloseChannel 有证据:114559
* UpdateNonClosingBalanceProof:97152
* settle channel:76331(实际),152662(估计)
* CooperativeSettle: 87241,174481
* withdraw:129333
* unlock:75094(3个锁)
* PunishObsoleteUnlock:34363,49363


## 记录
OpenChannel complete.. gasLimit=135682,gasUsed=135682 fn=contracts_test.go:creatAChannelAndDeposit:101
Deposit2 complete.. gasLimit=68808,gasUsed=68808 fn=contracts_test.go:creatAChannelAndDeposit:137
OpenChannel complete.. gasLimit=135682,gasUsed=135682 fn=contracts_test.go:creatAChannelAndDeposit:101
Deposit2 complete.. gasLimit=68808,gasUsed=68808 fn=contracts_test.go:creatAChannelAndDeposit:137
CloseChannel no evidence complete.. gasLimit=56048,gasUsed=56048 fn=contracts_test.go:TestCloseChannel1:209
OpenChannel complete.. gasLimit=135682,gasUsed=135682 fn=contracts_test.go:creatAChannelAndDeposit:101
Deposit2 complete.. gasLimit=68808,gasUsed=68808 fn=contracts_test.go:creatAChannelAndDeposit:137
CloseChannel with evidence complete.. gasLimit=114559,gasUsed=114559 fn=contracts_test.go:TestCloseChannel2:314
OpenChannel complete.. gasLimit=135682,gasUsed=135682 fn=contracts_test.go:creatAChannelAndDeposit:101
Deposit2 complete.. gasLimit=68808,gasUsed=68808 fn=contracts_test.go:creatAChannelAndDeposit:137
CloseChannel   complete.. gasLimit=114559,gasUsed=114559 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle:398
UpdateNonClosingBalanceProof   complete.. gasLimit=97152,gasUsed=97152 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle:414
SettleChannel  complete.. gasLimit=76331,gasUsed=152662 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle:455
OpenChannel complete.. gasLimit=135682,gasUsed=135682 fn=contracts_test.go:creatAChannelAndDeposit:101
Deposit2 complete.. gasLimit=68808,gasUsed=68808 fn=contracts_test.go:creatAChannelAndDeposit:137
CloseChannel   complete.. gasLimit=114559,gasUsed=114559 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle2:479
UpdateNonClosingBalanceProof   complete.. gasLimit=97105,gasUsed=97105 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle2:495
UpdateNonClosingBalanceProof2   complete.. gasLimit=67088,gasUsed=67088 fn=contracts_test.go:TestCloseChannelAndUpdateNonClosingAndSettle2:518
OpenChannel complete.. gasLimit=135682,gasUsed=135682 fn=contracts_test.go:creatAChannelAndDeposit:101
Deposit2 complete.. gasLimit=68808,gasUsed=68808 fn=contracts_test.go:creatAChannelAndDeposit:137
CooperativeSettle   complete.. gasLimit=87241,gasUsed=174481 fn=contracts_test.go:TestCooperateSettleChannel:583
OpenChannel complete.. gasLimit=135682,gasUsed=135682 fn=contracts_test.go:creatAChannelAndDeposit:101
Deposit2 complete.. gasLimit=68808,gasUsed=68808 fn=contracts_test.go:creatAChannelAndDeposit:137
WithDraw complete.. gasLimit=129333,gasUsed=129333 fn=contracts_test.go:TestWithdraw:716
OpenChannel complete.. gasLimit=135682,gasUsed=135682 fn=contracts_test.go:creatAChannelAndDeposit:101
Deposit2 complete.. gasLimit=68808,gasUsed=68808 fn=contracts_test.go:creatAChannelAndDeposit:137
close channel successful,gasused=114559,gasLimit=114559 fn=contracts_test.go:TestUnlock:796
UpdateNonClosingBalanceProof successful,gasused=114200,gasLimit=114200,locksroot=bf6298219cac23ebd9311a9cd129516a1b80700dbf9f7d7f63c7ad2b0730279e,transferamount=10 fn=contracts_test.go:TestUnlock:841
unlock success,gasUsed=75904,gasLimit=75904,txhash=0x04cdba417b2374475aaaae9a4f5beb2b0e54c21800adee7914bd155f67ec78d6 fn=contracts_test.go:TestUnlock:859
settle channel complete ,gasused=76331,gasLimit=152662 fn=contracts_test.go:TestUnlock:907
OpenChannel complete.. gasLimit=135682,gasUsed=135682 fn=contracts_test.go:creatAChannelAndDeposit:101
Deposit2 complete.. gasLimit=68808,gasUsed=68808 fn=contracts_test.go:creatAChannelAndDeposit:137
close channel successful,gasused=114512,gasLimit=114512 fn=contracts_test.go:TestPunishObsoleteUnlock:1010
UpdateNonClosingBalanceProof successful,gasused=114136,gasLimit=114136,locksroot=c9c202758167c711e24661b62721c7d23c31f9d95b2a488131d0e132141d97f3,transferamount=10 fn=contracts_test.go:TestPunishObsoleteUnlock:1055
PunishObsoleteUnlock success,gasUsed=34363,gasLimit=49363,txhash=0xfc6aea22f6af185fa3f25af6f6d07458f42ec2cd728086e2fbf66de4e180b825 fn=contracts_test.go:TestPunishObsoleteUnlock:1085
