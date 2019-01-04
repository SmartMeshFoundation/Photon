#花费 gas 统计
* Approve 30216
* Deposit(OpenChannel) 92261
* Deposit(OpenChannel) tokenFallback 88062
* Deposit(OpenChannel) ApproveAndCall 102572
* Deposit 70764 
* Deposit  tokenFallback 51495
* Deposit  ApproveAndCall 66008
* CloseChannel 无证据:33762
* CloseChannel 有证据:67929
* UpdateBalanceProofDelegate:76705 
* updateBalanceProof 63160 
* settle channel:51352 
* CooperativeSettle:78174
* withdraw:83213
* unlock:68585 单个锁
* punish: 34448 

## 记录
smoke_test.go:184: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0xE514fbb7e751CdF59C9e765C58b6daFcF7B97D49 for 0xF5DEcCfb4935eF57B500807a5214120ADDC86f74,gasUsed=30216,gasLimit=30216
INFO [12-29|11:11:12.116] Deposit complete...,gasLimit=70928,gasUsed=70764 fn=smoke_test.go:creatAChannelAndDeposit:165
    smoke_test.go:184: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0xE514fbb7e751CdF59C9e765C58b6daFcF7B97D49 for 0xF5DEcCfb4935eF57B500807a5214120ADDC86f74,gasUsed=30216,gasLimit=30216
INFO [12-29|11:11:20.463] Deposit complete...,gasLimit=70928,gasUsed=70764 fn=smoke_test.go:creatAChannelAndDeposit:165
INFO [12-29|11:11:23.581] CloseChannel no evidence gasLimit=33762,gasUsed=33762 fn=smoke_test.go:TestCloseChannel1:241
    smoke_test.go:184: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0xE514fbb7e751CdF59C9e765C58b6daFcF7B97D49 for 0xF5DEcCfb4935eF57B500807a5214120ADDC86f74,gasUsed=30216,gasLimit=30216
INFO [12-29|11:11:29.931] Deposit complete...,gasLimit=70928,gasUsed=70764 fn=smoke_test.go:creatAChannelAndDeposit:165
INFO [12-29|11:11:32.055] CloseChannel with evidence gasLimit=67929,gasUsed=67929 fn=smoke_test.go:TestCloseChannel2:367
    smoke_test.go:184: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0xE514fbb7e751CdF59C9e765C58b6daFcF7B97D49 for 0xF5DEcCfb4935eF57B500807a5214120ADDC86f74,gasUsed=30216,gasLimit=30216
INFO [12-29|11:11:40.403] Deposit complete...,gasLimit=70928,gasUsed=70764 fn=smoke_test.go:creatAChannelAndDeposit:165
args="0xE514fbb7e751CdF59C9e765C58b6daFcF7B97D49","0x292650fee408320D888e06ed89D938294Ea42f99", "0xC76F9b0aDcC0bC02a63Bda1e803E61c92fe24e98",10,"0x0000000000000000000000000000000000000000000000000000000000000000",3,"0x64e604787cbf194841e7b68d7cd28786f6c9a0a3ab9f8b0a0e87cb4387ab0107","0xfa85d549349c7bcbd0ec6ec97fcc2a679e9d3c4833f2d57f940055092a1acbcd74579613d282e8600880655997303129a97b4e1c3070ebf7bf78109f441b51ad1b","0x0f196a9e60ca90f1a46d2e7f675b2eff0362906fd3bb677cc80798f62e6f49632a828b40ae5dd513d5e05372afb8993e2472b7b8b8c092118f24f378a25e170e1c"INFO 
[12-29|11:12:00.705] UpdateBalanceProofDelegate gasLimit=76705,gasUsed=76705 fn=smoke_test.go:TestCloseChannelAndUpdateBalanceProofDelegateAndSettle:499
INFO [12-29|11:12:20.878] SettleChannel gasLimit=102704,gasUsed=51352 fn=smoke_test.go:TestCloseChannelAndUpdateBalanceProofDelegateAndSettle:546
    smoke_test.go:184: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0xE514fbb7e751CdF59C9e765C58b6daFcF7B97D49 for 0xF5DEcCfb4935eF57B500807a5214120ADDC86f74,gasUsed=30216,gasLimit=30216
INFO [12-29|11:12:27.222] Deposit complete...,gasLimit=70928,gasUsed=70764 fn=smoke_test.go:creatAChannelAndDeposit:165
INFO [12-29|11:12:35.479] UpdateBalanceProof gasLimit=63160,gasUsed=63160 fn=smoke_test.go:TestCloseChannelAndUpdateBalanceProofAndSettle:610
INFO [12-29|11:13:10.711] SettleChannel gasLimit=102704,gasUsed=51352 fn=smoke_test.go:TestCloseChannelAndUpdateBalanceProofAndSettle:659
    smoke_test.go:184: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0xE514fbb7e751CdF59C9e765C58b6daFcF7B97D49 for 0xF5DEcCfb4935eF57B500807a5214120ADDC86f74,gasUsed=30216,gasLimit=30216
INFO [12-29|11:13:17.056] Deposit complete...,gasLimit=70928,gasUsed=70764 fn=smoke_test.go:creatAChannelAndDeposit:165
INFO [12-29|11:13:19.188] CooperativeSettle gasLimit=123174,gasUsed=78174 fn=smoke_test.go:TestCooperateSettleChannel:730
    smoke_test.go:184: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0xE514fbb7e751CdF59C9e765C58b6daFcF7B97D49 for 0xF5DEcCfb4935eF57B500807a5214120ADDC86f74,gasUsed=30216,gasLimit=30216
INFO [12-29|11:13:31.667] Deposit complete...,gasLimit=70928,gasUsed=70764 fn=smoke_test.go:creatAChannelAndDeposit:165
INFO [12-29|11:13:48.389] unlock success,gasUsed=68585,gasLimit=68585,txhash=0x46923485193e9ea2140149cf05a91e1f517c3a416f6e3fe28f0f787775a4129b fn=smoke_test.go:TestUnlock:921
    smoke_test.go:184: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0xE514fbb7e751CdF59C9e765C58b6daFcF7B97D49 for 0xF5DEcCfb4935eF57B500807a5214120ADDC86f74,gasUsed=30216,gasLimit=30216
INFO [12-29|11:14:26.938] Deposit complete...,gasLimit=70928,gasUsed=70764 fn=smoke_test.go:creatAChannelAndDeposit:165
INFO [12-29|11:14:29.068] WithDraw complete.. gasLimit=83213,gasUsed=83213 fn=smoke_test.go:TestWithdraw:1062
    smoke_test.go:184: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0xE514fbb7e751CdF59C9e765C58b6daFcF7B97D49 for 0xF5DEcCfb4935eF57B500807a5214120ADDC86f74,gasUsed=30216,gasLimit=30216
INFO [12-29|11:14:38.426] Deposit complete...,gasLimit=70928,gasUsed=70764 fn=smoke_test.go:creatAChannelAndDeposit:165
INFO [12-29|11:14:55.155] unlockdelegate gasLimit=82129,gasUsed=82129 fn=smoke_test.go:TestPunishObsoleteUnlock:1282
INFO [12-29|11:14:57.279] PunishObsoleteUnlock success,gasUsed=34448,gasLimit=64448,txhash=0xa8ee7be3ccb6059edee8a18d5a8931b63f4712e15dc1aeb3ab4845f5de8cf8ec fn=smoke_test.go:TestPunishObsoleteUnlock:1318
    smoke_test.go:184: 0x292650fee408320D888e06ed89D938294Ea42f99 approve token 0xE514fbb7e751CdF59C9e765C58b6daFcF7B97D49 for 0xF5DEcCfb4935eF57B500807a5214120ADDC86f74,gasUsed=30216,gasLimit=30216
INFO [12-29|11:14:59.405] open channel and deposit by tokenFallback success,gasUsed=88062,gasLimit=88062,txhash=0x600411726da351ede36a14c54f0c1b036744e0298e495d3cf3f7d74af87208b3 fn=smoke_test.go:testOpenChannelAndDepositFallback:1369
INFO [12-29|11:15:02.532] open channel and deposit by ApproveAndCall success,gasUsed=102572,gasLimit=118605,txhash=0x14fe3957b9e4902ad8a5385ab67cdec0fa82c37e48df7dc850f684d40022439d fn=smoke_test.go:testOpenChannelAndDepositApproveCall:1412