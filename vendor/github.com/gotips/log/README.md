[![Build Status](https://travis-ci.org/gotips/log.svg?branch=develop)](https://travis-ci.org/gotips/log)

log
===

Golang 标准库中提供了基本的 log 模块 http://golang.org/pkg/log ，能 print/fatal/panic
日志，但唯独不能像 Java log4j 一样设置日志输出级别， debug 日志开发时输出，生产上关闭。这不能
不说是个巨大的遗憾， gopher 们只能抱怨 Golang 标准库的 log 就是个然并卵，实际项目大多不会使用
它。尽管 print 可以当成 error 使用（标准库确实把日志打印到错误输出流 os.Stdout），但是开发时
 debug/info 就没办法。

或许对标准库的设计 Golang 开发团队有自己的考虑，但是对应用开发者来说，log4j 已经成为事实上的标
准。为了向这个标准库靠近，出现了众多第三方 log 库，有在标准库基础上扩展的（也许 Golang 设计者
们也是想让开发者自己扩展标准库的 log 呢），也就另辟蹊径，也有玩各种花样的。

虽然有那么多的 log 库，但都是大同小异，我们需要的也只是个标准的可以自定义级别的 log 库而已，就
像 slf4j(Simple Logging Facade for Java) 一样，所以这个 log 库的需要完成得任务就是提供一
个标准统一的接口，同时也提供了一个基本的实现，可以自己定义模板格式，输出各种类型的日志，如
csv/json/xml，同时支持 TraceID。

使用这个 log 库打印日志，可以随时切换日志级别，可以更换不同的 logger 实现，以打印不同格式的日
志，也可以改变日志输出位置，输出到数据库、消息队列等，者所有的改变都无需修改已经写好的项目源码。


Usage
-----

安装：`go get -v -u github.com/gotips/log`

使用：
``` go
package main

import "github.com/gotips/log"

func main() {
    log.Debugf("this is a test message, %d", 1111)

	format := fmt.Sprintf("%s %s %s %s:%d %s", "2006-01-02 15:04:05.000000", log.TagToken,
		log.LevelToken, log.ProjectToken, log.LineToken, log.MessageToken)
	log.SetFormat(format)
	log.Tinfof("6ba7b814-9dad-11d1-80b4-00c04fd430c8", "this is a test message, %d", 1111)

	format = fmt.Sprintf(`{"date": "%s", "time": "%s", "level": "%s", "file": "%s", "line": %d, "log": "%s"}`,
		"2006-01-02", "15:04:05.999", log.LevelToken, log.ProjectToken, log.LineToken, log.MessageToken)
	log.SetFormat(format)
	log.Infof("this is a test message, %d", 1111)

	format = fmt.Sprintf(`<log><date>%s</date><time>%s</time><level>%s</level><file>%s</file><line>%d</line><msg>%s</msg><log>`,
		"2006-01-02", "15:04:05.000", log.LevelToken, log.ProjectToken, log.LineToken, log.MessageToken)
	log.SetFormat(format)
	log.Tinfof("6ba7b814-9dad-11d1-80b4-00c04fd430c8", "this is a test message, %d", 1111)
}
```
日志输出：
```
2016-01-16 20:28:34 debug examples/main.go:10 this is a test message, 1111
2016-01-16 20:28:34.280601 6ba7b814-9dad-11d1-80b4-00c04fd430c8 info examples/main.go:15 this is a test message, 1111
{"date": "2016-01-16", "time": "20:28:34.28", "level": "info", "file": "examples/main.go", "line": 20, "log": "this is a test message, 1111"}
<log><date>2016-01-16</date><time>20:28:34.280</time><level>info</level><file>examples/main.go</file><line>25</line><msg>this is a test message, 1111</msg><log>

```

更多用法 [examples](examples/main.go)


Go Doc and API
--------------

所有可调用的接口 API 和 文档都在 [log.go](log.go)


log/Printer/Standard
--------------------

Golang 不同于 Java，非面向对象语言（没有继承，只有组合，不能把组合实例赋给被组合的实例，即 Java
说的 子对象 赋给 父对象），为了方便使用，很多函数都是包封装的，无需创建 struct ，就可以直接调用。
（一般把裸露的方法称为函数，结构体和其他类型的方法才称为某某的方法）

log 包也一样，使用时，无需 new ，直接用。log 包有所有级别的函数可以调用，所有函数最终都调用了
print 函数。print 函数又调用了包内部变量的 std 的 Print 方法。这个 std 是一个 Printer 接
口类型，定义了打印接口。用不同的实现改变 std 就可以打印出不同格式的日志，也可以输出到不同位置。
（这个接口貌似还没有抽象好，再想想）

Printer 有个基本的实现 Standard，如果不改变，默认使用这个实现打印日志。

Standard 实现了的 Printer 接口，把日志打印到 Stdout。


性能测试
-------

环境：MacBookPro 17，4 核 8 线程 16G 内存

实际测试结果（把日志重定向到文件）：
该库平均每秒可输出 16w 行日志；
Go 语言标准库平均每秒输出 36.5w 行日志。

模板方式输出日志对性能有一定影响，其他 New Record 等也可能造成性能下降。但是实现上比标准库略微
复杂，输出格式可以灵活配置，所以整体上可以接受，后期随着对 Go 语言的学习更深入再不断优化。


TODO
----

* 测试是否支持各种格式的日期
* 处理秒和毫秒，如1:1:02.9
* 实现日志文件按一定规则自动滚动


Others
------

最近更新请移步 https://github.com/omigo/log
