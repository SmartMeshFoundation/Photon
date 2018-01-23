package log

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// 可以用这些串和日期、时间（包含毫秒数）任意组合，拼成各种格式的日志，如 csv/json/xml
const (
	LevelToken   string = "info"
	TagToken            = "tag"
	PathToken           = "/go/src/github.com/gotips/log/examples/main.go"
	PackageToken        = "github.com/gotips/log/examples/main.go"
	ProjectToken        = "examples/main.go"
	FileToken           = "main.go"
	LineToken    int    = 88
	MessageToken string = "message"
)

// DefaultFormat 默认日志格式
const DefaultFormat = "2006-01-02 15:04:05 info examples/main.go:88 message"

// DefaultFormatTag 默认日志格式带标签
const DefaultFormatTag = "2006-01-02 15:04:05 tag info examples/main.go:88 message"

// ExtactDateTime 抽取日期和时间格式字符串串
func ExtactDateTime(format string) (dateFmt, timeFmt string) {
	// 算法：
	// 找出两个字符串不同的部分，
	// 如果有两处不同，一个是日期模式，一个是时间模式，
	// 如果只有一个，那么只有日期或者只有时间，无关紧要，
	// 如果都相同，那么日志里没有时间，
	// 如果有三处以上不同，说明格式配置错误

	t, _ := time.ParseInLocation("2006-1-2 3:4:5.000000000", "1991-2-1 1:1:1.111111111", time.Local)
	contrast := t.Format(format)

	// println(format)
	// println(contrast)

	idxs := [10]int{}
	start := -1
	for i, l, same := 0, len(format), true; i < l; i++ {
		if start > 4 {
			panic(fmt.Sprintf("format string error at `%s`", format[i-1:]))
		}

		// fmt.Printf("%c %c %d %d\n", format[i], contrast[i], idxs, start)

		if format[i] != contrast[i] {
			if same {
				start++
				// 如果之前都是相同的，这个开始不同，那么这个就是起始位置
				idxs[start] = i
				same = false
				// println(i, diff, start, idxs[start])
				start++
			}

			idxs[start] = i + 1 // 下一个有可能是结束位置

			continue
		}

		// 如果是 空格、-、:、. ，那么它不一定是结束位置
		if format[i] == '-' || format[i] == ' ' || format[i] == ':' || format[i] == '.' {
			// 如果这些字符后面是 0（如 2006-01-02），跳过 0
			if i+1 < l && format[i+1] == '0' && contrast[i+1] == '0' {
				i++
			}
			continue
		}

		same = true
	}

	if start != -1 && start != 1 && start != 3 {
		// 正常情况是不可能到这里的，如果到这里，说明算法写错了
		panic(fmt.Sprintf("parse error %d", start))

	} else {
		if start > 0 {
			dateFmt = format[idxs[0]:idxs[1]]
			if start == 3 {
				timeFmt = format[idxs[2]:idxs[3]]
			}
		}
	}

	return dateFmt, timeFmt
}

// CalculatePrefixLen 计算包前缀，如果格式中不包含包文件路径，那么就返回 -1
func CalculatePrefixLen(format string, skip int) int {
	// 格式中不包含文件路径
	if !strings.Contains(format, "main.go") {
		return -1
	}

	_, file, _, _ := runtime.Caller(skip)

	// file with absolute path
	if strings.Contains(format, PathToken) {
		return 0
	}

	// file with package name
	if strings.Contains(format, PackageToken) {
		return strings.Index(file, "/src/") + 5
	}

	// file with project path
	if strings.Contains(format, ProjectToken) {
		// remove /<GOPATH>/src/
		prefixLen := strings.Index(file, "/src/") + 5
		file = file[prefixLen:]

		// remove github.com/
		if strings.HasPrefix(file, "github.com/") {
			prefixLen += 11
			file = file[11:]

			// remove github user or org name
			if i := strings.Index(file, "/"); i >= 0 {
				prefixLen += i + 1
				file = file[i+1:]

				// remove project name
				if i := strings.Index(file, "/"); i >= 0 {
					prefixLen += i + 1
				}
			}
		}

		return prefixLen
	}

	// file only
	if strings.Contains(format, FileToken) {
		return strings.LastIndex(file, "/") + 1
	}

	return -1
}
