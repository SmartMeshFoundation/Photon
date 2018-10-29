# casemanager工具使用说明
1. cd 到casemanager目录,执行go build
2. 执行./casemanager list 查看当前支持的所有case名
3. 执行./casemanager --case=all按顺序执行所有case,./casemanager --case=case名，可执行指定case
4. 当执行所有case时，可添加参数--skip=true，可报错不中断跑完所有case
5. cases目录下各case配置文件中有debug参数，配置为true则case执行完毕后不杀死Photon节点，方便后续调试。
6. 各case场景描述，参考cases目录下各case源码第一行注释
7. 日志说明：
    所有日志均在log目录下
    case名.log   如   CrashCaseSend01.log                 为case日志，包含channel数据等，方便查阅。
                                                          其中channel数据命名方式为CD-节点名1-节点名2-描述，比如
                                                          CD-N1-N2-BeforeTransfer       代表N1-N2之间的通道在交易发出前的状态，即初始状态
                                                          CD-N1-N2-AfterCrash           代表N1-N2之间的通道在崩溃后的状态
                                                          CD-N1-N2Restart-AfterRestart  代表N1-N2之间的通道在崩溃节点重启后的状态，即最终状态
    case名-Nx.log如   CrashCaseSend01-N1.log              为各Photon节点日志。
    case名-Nx.log如   CrashCaseSend01-N1Restart.log       为崩溃恢复case中重启节点重启后的日志
8. 如有问题，请咨询wuhan_53@163.com或联系我本人