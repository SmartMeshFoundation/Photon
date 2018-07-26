#CaseList
## 1. Case about channel open
    - 正确的参数
## 2. Case about channel but expect fail
    - settle_timeout = 0
    - self地址为0x0
    - self地址为""
    - self地址为0x03432
    - self地址为0x0000000000000000000000000000000000000000
    - partner地址为0x0
    - partner地址为""
    - partner地址为0x03432
    - partner地址为0x0000000000000000000000000000000000000000
    - 通道双方地址相同
    - settle_timeout = 5
    - settle_timeout = 2700001
    - 重复open
## 3. Case about opened channel state
    - 打开一个新通道并查询通道双方信息, 校验数据