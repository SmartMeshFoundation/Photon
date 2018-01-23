package log

import "io"

// 默认 debug 级别，方便调试，生产环境可以调用 SetLevel 设置 log 级别
var v Level = DebugLevel

// 默认实现，输出到 os.Std 中，可以重定向到文件中，也可以调用 SetPrinter 其他方式输出
var std Printer

// SetLevel 设置日志级别
func SetLevel(l Level) { v = l }

// Colorized 输出日志是否着色，默认着色
func Colorized(c bool) { std.Colorized(c) }

// GetLevel 返回设置的日志级别
func GetLevel() (l Level) { return v }

// SetPrinter 切换 Printer 实现
func SetPrinter(p Printer) { std = p }

// SetWriter 改变输出位置，通过这个接口，可以实现日志文件按时间或按大小滚动
func SetWriter(w io.Writer) { std.SetWriter(w) }

// SetFormat 改变日志格式
func SetFormat(format string) { std.SetFormat(format) }

// 判断各种级别的日志是否会被输出
func IsTraceEnabled() bool { return v <= TraceLevel }
func IsDebugEnabled() bool { return v <= DebugLevel }
func IsInfoEnabled() bool  { return v <= InfoLevel }
func IsWarnEnabled() bool  { return v <= WarnLevel }
func IsErrorEnabled() bool { return v <= ErrorLevel }
func IsPanicEnabled() bool { return v <= PanicLevel }
func IsFatalEnabled() bool { return v <= FatalLevel }
func IsPrintEnabled() bool { return v <= PrintLevel }
func IsStackEnabled() bool { return v <= StackLevel }

// 打印日志
func Trace(m ...interface{}) { std.Tprintf(v, TraceLevel, "", "", m...) }
func Debug(m ...interface{}) { std.Tprintf(v, DebugLevel, "", "", m...) }
func Info(m ...interface{})  { std.Tprintf(v, InfoLevel, "", "", m...) }
func Warn(m ...interface{})  { std.Tprintf(v, WarnLevel, "", "", m...) }
func Error(m ...interface{}) { std.Tprintf(v, ErrorLevel, "", "", m...) }
func Panic(m ...interface{}) { std.Tprintf(v, PanicLevel, "", "", m...) }
func Fatal(m ...interface{}) { std.Tprintf(v, FatalLevel, "", "", m...) }
func Print(m ...interface{}) { std.Tprintf(v, PrintLevel, "", "", m...) }
func Stack(m ...interface{}) { std.Tprintf(v, StackLevel, "", "", m...) }

// 按一定格式打印日志
func Tracef(format string, m ...interface{}) { std.Tprintf(v, TraceLevel, "", format, m...) }
func Debugf(format string, m ...interface{}) { std.Tprintf(v, DebugLevel, "", format, m...) }
func Infof(format string, m ...interface{})  { std.Tprintf(v, InfoLevel, "", format, m...) }
func Warnf(format string, m ...interface{})  { std.Tprintf(v, WarnLevel, "", format, m...) }
func Errorf(format string, m ...interface{}) { std.Tprintf(v, ErrorLevel, "", format, m...) }
func Panicf(format string, m ...interface{}) { std.Tprintf(v, PanicLevel, "", format, m...) }
func Fatalf(format string, m ...interface{}) { std.Tprintf(v, FatalLevel, "", format, m...) }
func Printf(format string, m ...interface{}) { std.Tprintf(v, PrintLevel, "", format, m...) }
func Stackf(format string, m ...interface{}) { std.Tprintf(v, StackLevel, "", format, m...) }

// 打印日志时带上 tag
func Ttrace(tag string, m ...interface{}) { std.Tprintf(v, TraceLevel, tag, "", m...) }
func Tdebug(tag string, m ...interface{}) { std.Tprintf(v, DebugLevel, tag, "", m...) }
func Tinfo(tag string, m ...interface{})  { std.Tprintf(v, InfoLevel, tag, "", m...) }
func Twarn(tag string, m ...interface{})  { std.Tprintf(v, WarnLevel, tag, "", m...) }
func Terror(tag string, m ...interface{}) { std.Tprintf(v, ErrorLevel, tag, "", m...) }
func Tpanic(tag string, m ...interface{}) { std.Tprintf(v, PanicLevel, tag, "", m...) }
func Tfatal(tag string, m ...interface{}) { std.Tprintf(v, FatalLevel, tag, "", m...) }
func Tprint(tag string, m ...interface{}) { std.Tprintf(v, PrintLevel, tag, "", m...) }
func Tstack(tag string, m ...interface{}) { std.Tprintf(v, StackLevel, tag, "", m...) }

// 按一定格式打印日志，并在打印日志时带上 tag
func Ttracef(tag string, format string, m ...interface{}) { std.Tprintf(v, TraceLevel, tag, format, m...)}
func Tdebugf(tag string, format string, m ...interface{}) { std.Tprintf(v, DebugLevel, tag, format, m...)}
func Tinfof(tag string, format string, m ...interface{})  { std.Tprintf(v, InfoLevel, tag, format, m...) }
func Twarnf(tag string, format string, m ...interface{})  { std.Tprintf(v, WarnLevel, tag, format, m...) }
func Terrorf(tag string, format string, m ...interface{}) { std.Tprintf(v, ErrorLevel, tag, format, m...)}
func Tpanicf(tag string, format string, m ...interface{}) { std.Tprintf(v, PanicLevel, tag, format, m...)}
func Tfatalf(tag string, format string, m ...interface{}) { std.Tprintf(v, FatalLevel, tag, format, m...)}
func Tprintf(tag string, format string, m ...interface{}) { std.Tprintf(v, PrintLevel, tag, format, m...)}
func Tstackf(tag string, format string, m ...interface{}) { std.Tprintf(v, StackLevel, tag, format, m...)}

// ======== 兼容 qiniu/log   ===============
func SetOutputLevel(l Level) { v = l }

// ======== 兼容 wothing/log ===============

// 打印日志时带上 tag
func TraceT(tag string, m ...interface{}) { std.Tprintf(v, TraceLevel, tag, "", m...) }
func DebugT(tag string, m ...interface{}) { std.Tprintf(v, DebugLevel, tag, "", m...) }
func InfoT(tag string, m ...interface{})  { std.Tprintf(v, InfoLevel, tag, "", m...) }
func WarnT(tag string, m ...interface{})  { std.Tprintf(v, WarnLevel, tag, "", m...) }
func ErrorT(tag string, m ...interface{}) { std.Tprintf(v, ErrorLevel, tag, "", m...) }
func PanicT(tag string, m ...interface{}) { std.Tprintf(v, PanicLevel, tag, "", m...) }
func FatalT(tag string, m ...interface{}) { std.Tprintf(v, FatalLevel, tag, "", m...) }
func PrintT(tag string, m ...interface{}) { std.Tprintf(v, PrintLevel, tag, "", m...) }
func StackT(tag string, m ...interface{}) { std.Tprintf(v, StackLevel, tag, "", m...) }

// 按一定格式打印日志，并在打印日志时带上 tag
func TracefT(tag string, format string, m ...interface{}) { std.Tprintf(v, TraceLevel, tag, format, m...)}
func DebugfT(tag string, format string, m ...interface{}) { std.Tprintf(v, DebugLevel, tag, format, m...)}
func InfofT(tag string, format string, m ...interface{})  { std.Tprintf(v, InfoLevel, tag, format, m...) }
func WarnfT(tag string, format string, m ...interface{})  { std.Tprintf(v, WarnLevel, tag, format, m...) }
func ErrorfT(tag string, format string, m ...interface{}) { std.Tprintf(v, ErrorLevel, tag, format, m...)}
func PanicfT(tag string, format string, m ...interface{}) { std.Tprintf(v, PanicLevel, tag, format, m...)}
func FatalfT(tag string, format string, m ...interface{}) { std.Tprintf(v, FatalLevel, tag, format, m...)}
func PrintfT(tag string, format string, m ...interface{}) {std.Tprintf(v, PrintLevel, tag, format, m...)}
func StackfT(tag string, format string, m ...interface{}) {std.Tprintf(v, StackLevel, tag, format, m...)}
