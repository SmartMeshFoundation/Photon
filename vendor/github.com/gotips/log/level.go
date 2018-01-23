package log

import (
	"errors"
	"strconv"
)

// Level 日志级别
type Level uint8

// 所有日志级别常量，级别越高，日志越重要，级别越低，日志越详细
const (
	AllLevel Level = iota // 等同于 TraceLevel

	TraceLevel
	DebugLevel // 默认日志级别，方便开发
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel // panic 打印错误栈，但是可以 recover
	FatalLevel // fatal 表明严重错误，程序直接退出，慎用

	// 提供这个级别日志，方便在无论何种情况下，都打印必要信息，比如服务启动信息
	PrintLevel
	StackLevel // 打印堆栈信息
)

// Labels 每个级别对应的标签
var Labels = [...]string{"all", "trace", "debug", "info", "warn", "error", "panic", "fatal", "print", "stack"}

// String 返回日志标签
func (v Level) String() string {
	return Labels[v]
}

// ValueOfLevel 字符串转换成 Level, "debug" => DebugLevel
func ValueOfLevel(vstr string) (v Level, err error) {
	for i, label := range Labels {
		if vstr == label {
			return Level(i), nil
		}
	}
	return InfoLevel, errors.New("level " + vstr + " not exist")
}

// MarshalJSON 便于 JSON 解析
func (v *Level) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(v.String())), nil
}

// UnmarshalJSON 便于 JSON 解析
func (v *Level) UnmarshalJSON(b []byte) error {
	vstr, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	x, err := ValueOfLevel(vstr)
	if err != nil {
		return err
	}
	*v = x
	return nil
}

// ======== 兼容 qiniu/log   ===============
const (
	Ldebug = DebugLevel
	Linfo  = InfoLevel
	Lwarn  = WarnLevel
	Lerror = ErrorLevel
	Lpanic = PanicLevel
	Lfatal = FatalLevel
)
