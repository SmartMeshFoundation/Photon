package log

import (
	"runtime"
	"strings"
	"testing"
)

func TestCalculatePrefixLen(t *testing.T) {
	format := `{"level": "info", "line": 88, "log": "message"}`
	prefixLen := CalculatePrefixLen(format, 1)
	if prefixLen != -1 {
		t.FailNow()
	}

	format = `{"level": "info", "file": "/go/src/github.com/gotips/log/examples/main.go", "line":88, "log": "message"}`
	prefixLen = CalculatePrefixLen(format, 1)
	if prefixLen != 0 {
		t.FailNow()
	}

	_, file, _, _ := runtime.Caller(0)

	format = `{"level": "info", "file": "github.com/gotips/log/examples/main.go", "line":88, "log": "message"}`
	prefixLen = CalculatePrefixLen(format, 1)
	if prefixLen != strings.Index(file, "/src/")+5 {
		t.FailNow()
	}

	println(file)
	format = `{"level": "info", "file": "examples/main.go", "line":88, "log": "message"}`
	prefixLen = CalculatePrefixLen(format, 1)
	if prefixLen != strings.LastIndex(file, "/")+1 {
		t.FailNow()
	}

	format = `{"level": "info", "file": "main.go", "line":88, "log": "message"}`
	prefixLen = CalculatePrefixLen(format, 1)
	if prefixLen != strings.LastIndex(file, "/")+1 {
		t.FailNow()
	}
}

func TestExtactDateTime(t *testing.T) {
	format := `{"level": "info", "file": "log/main.go", "line":88, "log": "message"}`
	dateFmt, timeFmt := ExtactDateTime(format)
	if dateFmt != "" && timeFmt != "" {
		t.FailNow()
	}

	format = `{"datetime": "2006-01-02 15:04:05.999999999", "level": "info", "file": "log/main.go", "line":88, "log": "message"}`
	dateFmt, timeFmt = ExtactDateTime(format)
	if dateFmt != "2006-01-02 15:04:05.999999999" && timeFmt != "" {
		t.FailNow()
	}

	format = `{"date": "2006-01-02", "time": "15:04:05.999999999", "level": "info", "file": "log/main.go", "line":88, "log": "message"}`
	dateFmt, timeFmt = ExtactDateTime(format)
	if dateFmt != "2006-01-02" && timeFmt != "15:04:05.999999999" {
		t.FailNow()
	}

	// 测试 日期模式不能重复出现在 format 中，不能判定是模式还是固定字符串
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Error("must panic, but not")
				t.FailNow()
			}
		}()

		// 有两个 2006 ，会出错
		format = `{"date": "2006-01-02", "time": "15:04:05.999999999", "Tag": "2006" "level": "info", "file": "log/main.go", "line":88, "log": "message"}`
		dateFmt, timeFmt = ExtactDateTime(format)
	}()
}
