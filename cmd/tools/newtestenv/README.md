# newTestEnv

newTestEnv is  a tool to create a new environment for test.

## pre condition
1. collect your accounts to a directory ,for example /tmp/keystore
```bash
# my keystore
ls /tmp/keystore
-rw-------   1 bai  staff   491 Sep  4 11:17 UTC--2018-09-04T03-17-31.189911766Z--65ce623b524719952093cd6cc752d48aa210173d
-rw-------   1 bai  staff   491 Sep  4 11:17 UTC--2018-09-04T03-17-38.798689713Z--82b0a532987334b2d6cfb4a85cb4ebcbcad45633
-rw-------   1 bai  staff   491 Sep  4 11:17 UTC--2018-09-04T03-17-48.188807345Z--5c61a1f89e9be46858782adb702ceaa8900a904f
-rw-------   1 bai  staff   491 Sep  4 11:17 UTC--2018-09-04T03-17-55.414338981Z--d8521dbbc38193ec8b67fd19f2fb6379a9a5e8e4
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.531722608Z--1c3f49c11e305b7b39e0bc34c76972abbfc3ec9c
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.561721508M--08b005a3f7480638d989f1d9b2e3ea37be44e796
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.561721608Z--0a2daab66a0ebf3c9f5582afd32a5556b77b3a6b
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.561722608Z--17f6d7c535033de8ba7a9185e32f93f3b228f51f
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.631722608Z--2bff1331bf17abce38f92d97591fecc1f9dcabc1
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.731722608Z--3eb5da003f498ced267ea7015abff93b777ff305
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.821722608Z--4c1f289c8893643aaccf016f652ba50c03e63204
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.831722608Z--4a2a9da65c7863f092c81bd3dc1bf421378a0327
-rw-r--r--   1 bai  staff   488 Oct 23 19:02 UTC--2018-09-21T03-37-02.841722608Z--52308e8e594516807817925945b38e714a982e61
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.861722608Z--53f8b0a72cc5a34ca52c7cbe170e26589567de28
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.871722608Z--586da0ed541e3360f6eeb79bca66763e9949ee30
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.881722608Z--68d588dd31d22544bce2aac6cd62516c651dff1e
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.891722608Z--86da83d8dd88f1cf635c4c6ad6f73e59117d5280
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.896622608Z--afeeeaf631cc8edc66fe669e234f8f0c46afd65b
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.896622618Z--c9288dd576e755b27711e3a0d3ee5653d16aaa5b
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.896622628Z--d88be7c8b316067b3dff26433e6cdb968f16f9a3
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.896622638Z--d8e9672d913031ea0812d2dca4f52f9812dc0f06
-rw-r--r--@  1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.896622648Z--e5a8386523124fb60a473b1de6674211f06c65d0
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.896622658Z--ecc5fbb97f4d93f98e5775b3eadd52ec29f094d3
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.896622668Z--f22b9449fd4d38086f1db35ea78a8a318192395f
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.896622678A--764f499acb4eb5b61d19b613e7ef01fa7f0bede2
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.896622678E--7c6998321e5a712341d14bcd3569ac6e22f57092
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.896622678Q--71585d1257f7e1cb238235e5ccc530067bb78e98
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.896622678R--7e8dc02d5d074c63c0088761c238165459c62d5a
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.896622678W--78cbcbaec517a288d61e28b1d6be5e527f9d212d
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.896622678Z--fc7d6052ea105a34c5cef272e6ba2c56b548026d
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.896722608Z--a5e61eb756111b86bf5ed4fc32c82d4dd861c4ad
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.897722608Z--91c2ea936f12deb167cea9e85680071b49e70e1b
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.898722608Z--8d75ca913df82a482a01bea35d588cf8a4f7171c
-rw-r--r--   1 bai  staff   489 Oct 23 19:02 UTC--2018-09-21T03-37-02.899722608Z--86ea659cbb3e9125feabe665a73df96a2234df21

```
2. there should be at least 7 accounts, and all these accounts should share the same password
3. all these accounts should have enough ether for  making transaction
4. there is a ether/spectrum RPC endpoint

## run
```bash
./newtestenv --keystore-path /tmp/keystore --eth-rpc-endpoint ws://127.0.0.1:8546 --password 123 --tokennum 2
```
