# casemanager工具使用说明
1. cd 到casemanager目录,执行go build
2. 执行./casemanager list 查看当前支持的所有case名
3. 执行./casemanager --case=all按顺序执行所有case,./casemanager --case=case名，可执行指定case
4. 日志在log目录下，其中
    case名.log如CrashCaseSend01.log    为case日志，包含channel数据等，方便查阅。
    case名-Nx.log如CrashCaseSend01-N1.log 为各smartraiden节点日志。
5. 各case场景描述，参考cases目录下各case源码第一行注释
6. 如有问题，请咨询wuhan_53@163.com或联系我本人