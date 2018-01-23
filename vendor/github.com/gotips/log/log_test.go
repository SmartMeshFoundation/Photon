package log

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

const uuid = "6ba7b814-9dad-11d1-80b4-00c04fd430c8"

func TestLogLevel(t *testing.T) {
	SetLevel(InfoLevel)
	if IsDebugEnabled() || !IsInfoEnabled() || !IsWarnEnabled() {
		t.FailNow()
	}
	SetLevel(DebugLevel) // 恢复现场，避免影响其他单元测试
}

func TestSetWriter(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 4096))
	SetWriter(buf)

	rand := time.Now().String()
	Info(rand)
	if !bytes.Contains(buf.Bytes(), ([]byte)(rand)) {
		t.FailNow()
	}
}

func TestSetFormat(t *testing.T) {
	format := fmt.Sprintf(`<log><date>%s</date><time>%s</time><level>%s</level><file>%s</file><line>%d</line><msg>%s</msg><log>`,
		"2006-01-02", "15:04:05.000", LevelToken, ProjectToken, LineToken, MessageToken)
	SetFormat(format)

	buf := bytes.NewBuffer(make([]byte, 4096))
	SetWriter(buf)

	rand := time.Now().String()
	Debug(rand)
	if bytes.HasPrefix(buf.Bytes(), ([]byte)("<log><date>")) &&
		!bytes.HasSuffix(buf.Bytes(), ([]byte)("</msg><log>")) {
		t.FailNow()
	}
}

func TestPanicLog(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fail()
		}
	}()
	Panic("test panic")
}

func TestNormalLog(t *testing.T) {
	SetLevel(AllLevel)

	Trace(AllLevel)
	Trace(TraceLevel)
	Debug(DebugLevel)
	Info(InfoLevel)
	Warn(WarnLevel)
	Error(ErrorLevel)
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Fail()
			}
		}()
		Panic(PanicLevel)
	}()
	// Fatal( FatalLevel)
	Print(PrintLevel)
	Stack(StackLevel)
}

func TestFormatLog(t *testing.T) {
	SetLevel(AllLevel)

	Tracef("%d %s", AllLevel, AllLevel)
	Tracef("%d %s", TraceLevel, TraceLevel)
	Debugf("%d %s", DebugLevel, DebugLevel)
	Infof("%d %s", InfoLevel, InfoLevel)
	Warnf("%d %s", WarnLevel, WarnLevel)
	Errorf("%d %s", ErrorLevel, ErrorLevel)
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Fail()
			}
		}()
		Panicf("%d %s", PanicLevel, PanicLevel)
	}()
	// Fatalf("%d %s", FatalLevel, FatalLevel)
	Printf("%d %s", PrintLevel, PrintLevel)
	Stackf("%d %s", StackLevel, StackLevel)
}

func TestNormalLogWithTag(t *testing.T) {
	format := "2006-01-02 15:04:05 tag info examples/main.go:88 message"
	SetFormat(format)
	SetLevel(AllLevel)

	Ttrace(uuid, AllLevel)
	Ttrace(uuid, TraceLevel)
	Tdebug(uuid, DebugLevel)
	Tinfo(uuid, InfoLevel)
	Twarn(uuid, WarnLevel)
	Terror(uuid, ErrorLevel)
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Fail()
			}
		}()
		Tpanic(uuid, PanicLevel)
	}()
	// Tfatal(uuid, FatalLevel)
	Tprint(uuid, PrintLevel)
	Tstack(uuid, StackLevel)
}

func TestFormatLogWithTag(t *testing.T) {
	format := "2006-01-02 15:04:05 tag info examples/main.go:88 message"
	SetFormat(format)
	SetLevel(AllLevel)

	Ttracef(uuid, "%d %s", AllLevel, AllLevel)
	Ttracef(uuid, "%d %s", TraceLevel, TraceLevel)
	Tdebugf(uuid, "%d %s", DebugLevel, DebugLevel)
	Tinfof(uuid, "%d %s", InfoLevel, InfoLevel)
	Twarnf(uuid, "%d %s", WarnLevel, WarnLevel)
	Terrorf(uuid, "%d %s", ErrorLevel, ErrorLevel)
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Fail()
			}
		}()
		Tpanicf(uuid, "%d %s", PanicLevel, PanicLevel)
	}()
	// Tfatalf(uuid,"%d %s", FatalLevel, FatalLevel)
	Tprintf(uuid, "%d %s", PrintLevel, PrintLevel)
	Tstackf(uuid, "%d %s", StackLevel, StackLevel)
}

func TestWothingNormalLogWithTag(t *testing.T) {
	format := "2006-01-02 15:04:05 tag info examples/main.go:88 message"
	SetFormat(format)
	SetLevel(AllLevel)

	TraceT(uuid, AllLevel)
	TraceT(uuid, TraceLevel)
	DebugT(uuid, DebugLevel)
	InfoT(uuid, InfoLevel)
	WarnT(uuid, WarnLevel)
	ErrorT(uuid, ErrorLevel)
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Fail()
			}
		}()
		PanicT(uuid, PanicLevel)
	}()
	// FatalT(uuid, FatalLevel)
	PrintT(uuid, PrintLevel)
	StackT(uuid, StackLevel)
}

func TestWothingFormatLogWithTag(t *testing.T) {
	format := "2006-01-02 15:04:05 tag info examples/main.go:88 message"
	SetFormat(format)
	SetLevel(AllLevel)

	TracefT(uuid, "%d %s", AllLevel, AllLevel)
	TracefT(uuid, "%d %s", TraceLevel, TraceLevel)
	DebugfT(uuid, "%d %s", DebugLevel, DebugLevel)
	InfofT(uuid, "%d %s", InfoLevel, InfoLevel)
	WarnfT(uuid, "%d %s", WarnLevel, WarnLevel)
	ErrorfT(uuid, "%d %s", ErrorLevel, ErrorLevel)
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Fail()
			}
		}()
		PanicfT(uuid, "%d %s", PanicLevel, PanicLevel)
	}()
	// FatalfT(uuid,"%d %s", FatalLevel, FatalLevel)
	PrintfT(uuid, "%d %s", PrintLevel, PrintLevel)
	StackfT(uuid, "%d %s", StackLevel, StackLevel)
}
