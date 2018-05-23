smoketest工具使用说明
1. 本地启动geth客户端服务
1.  初始化环境
    cd envinit
    go build
    ./envinit
    注意：默认连接的geth服务地址为：http://127.0.0.1:8545，如有不同，需在运行envinit时使用命令行参数eth-rpc-endpoint指定
          如已经运行过，可跳过该步骤
          如有现成的环境，可指定env.INI中的registry_contract_address和[ACCOUNT]下至少N0-N5共计6个账户以上，也可不运行envinit

2. 根据本地环境配置env.INI中的raidenpath指定smartraiden命令路径

3. 执行smoketest：
    cd ..
    go build
    ./smoketest
    注意:
        日志目录./log，日志文件说明:
            smoketest.log   :    测试用例日志
            before.data     :    测试运行前获取的当前节点所有数据汇总，包含节点、token、channel等数据
            after.data      :    测试运行后获取的当前节点所有数据汇总，包含节点、token、channel等数据，方便与before.data做对比
            N0-N5.log       ：   smartraiden节点日志
            killall.log     ：   smartraiden进程kill日志