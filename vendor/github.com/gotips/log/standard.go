package log

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

type record struct {
	Start, End string
	Date, Time string
	Tag        string
	Level      string
	File       string
	Line       int
	Message    string
	Stack      []byte
}

// Standard 日志输出基本实现
type Standard struct {
	mu  sync.Mutex    // ensures atomic writes; protects the following fields
	out *bufio.Writer // destination for output

	format    string
	pattern   string
	colorized bool

	tpl       *template.Template
	prefixLen int
	dateFmt   string
	timeFmt   string
}

// NewStandard 返回标准实现
func NewStandard(w io.Writer, format string) *Standard {
	std := &Standard{out: bufio.NewWriter(w), colorized: true}

	// hack 如果用户不调用 SetFormat，直接用，那么也能找到主函数（main，实际是 init 函数）的所在的文件
	std.prefixLen = -5

	std.SetFormat(format)
	return std
}

// SetWriter 改变输出流
func (s *Standard) SetWriter(w io.Writer) {
	s.mu.Lock()
	s.out = bufio.NewWriter(w)
	s.mu.Unlock()
}

// Colorized 输出日志是否着色，默认着色
func (s *Standard) Colorized(c bool) {
	// 没改变
	if c == s.colorized {
		return
	}

	s.colorized = c

	s.mu.Lock()
	defer s.mu.Unlock()

	p := s.pattern
	if s.colorized {
		p = "{{.Start}}" + p + "{{.End}}"
	}
	s.tpl = template.Must(template.New("record").Parse(p))
}

// SetFormat 改变日志格式
func (s *Standard) SetFormat(format string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.format = format

	skip := 3
	if s.prefixLen == -5 {
		skip = 5
	}
	s.prefixLen = CalculatePrefixLen(format, skip)

	s.dateFmt, s.timeFmt = ExtactDateTime(format)

	p := parseFormat(format, s.prefixLen, s.dateFmt, s.timeFmt)

	s.pattern = p
	if s.colorized {
		p = "{{.Start}}" + p + "{{.End}}"
	}
	s.tpl = template.Must(template.New("record").Parse(p))
}

// Tprintf 打印日志
func (s *Standard) Tprintf(v, l Level, tag string, format string, m ...interface{}) {
	if v > l {
		return
	}

	if tag == "" {
		tag = "-"
	}
	r := record{
		Level: l.String(),
		Tag:   tag,
	}

	if s.dateFmt != "" {
		now := time.Now() // get this early.
		r.Date = now.Format(s.dateFmt)
		if s.timeFmt != "" {
			r.Time = now.Format(s.timeFmt)
		}
	}

	if s.prefixLen > -1 {
		var ok bool
		_, r.File, r.Line, ok = runtime.Caller(2) // expensive
		if ok && s.prefixLen < len(r.File) {
			r.File = r.File[s.prefixLen:]
		} else {
			r.File = "???"
		}
	}

	if format == "" {
		r.Message = fmt.Sprint(m...)
	} else {
		r.Message = fmt.Sprintf(format, m...)
	}
	r.Message = strings.TrimRight(r.Message, "\n")

	if l == StackLevel {
		r.Stack = make([]byte, 1024*1024)
		n := runtime.Stack(r.Stack, true)
		r.Stack = r.Stack[:n]
	}

	if s.colorized {
		calculateColor(l, &r)
	}

	s.mu.Lock()
	defer func() {
		s.mu.Unlock()

		if l == PanicLevel {
			panic(m)
		}

		if l == FatalLevel {
			os.Exit(-1)
		}
	}()

	s.tpl.Execute(s.out, r)
	s.out.WriteByte('\n')

	if l == StackLevel {
		s.out.Write(r.Stack)
	}

	s.out.Flush()
}

// 格式解析，把格式串替换成 token 串
func parseFormat(format string, prefixLen int, dateFmt, timeFmt string) (pattern string) {
	// 顺序最好不要变，从最长的开始匹配
	pattern = strings.Replace(format, PathToken, "{{ .File }}", -1)
	pattern = strings.Replace(pattern, PackageToken, "{{ .File }}", -1)
	pattern = strings.Replace(pattern, ProjectToken, "{{ .File }}", -1)
	pattern = strings.Replace(pattern, FileToken, "{{ .File }}", -1)
	pattern = strings.Replace(pattern, TagToken, "{{ .Tag }}", -1)
	pattern = strings.Replace(pattern, LevelToken, "{{ .Level }}", -1)
	pattern = strings.Replace(pattern, strconv.Itoa(LineToken), "{{ .Line }}", -1)
	pattern = strings.Replace(pattern, MessageToken, "{{ .Message }}", -1)

	// 提取出日期和时间的格式化模式字符串
	if dateFmt != "" {
		pattern = strings.Replace(pattern, dateFmt, "{{ .Date }}", -1)
	}
	if timeFmt != "" {
		pattern = strings.Replace(pattern, timeFmt, "{{ .Time }}", -1)
	}
	return pattern
}
//在windows命令行下，颜色处理有问题
func calculateColor(l Level, r *record) {
	return
	if l < InfoLevel {
		return
	}

	// 字背景颜色范围:40----49
	// 40:黑
	// 41:深红
	// 42:绿
	// 43:黄色
	// 44:蓝色
	// 45:紫色
	// 46:深绿
	// 47:白色
	//
	// 字颜色:30-----------39
	// 30:黑
	// 31:红
	// 32:绿
	// 33:黄
	// 34:蓝色
	// 35:紫色
	// 36:深绿
	// 37:白色
	switch l {
	case InfoLevel:
		r.Start = "\033[36;1m" // green

	case WarnLevel:
		r.Start = "\033[33;1m" // yellow

	case ErrorLevel, PanicLevel, FatalLevel:
		r.Start = "\033[31;1m" // red

	case PrintLevel, StackLevel:
		r.Start = "\033[32;1m" // green

	default:
		r.Start = "\033[0m" // default
	}
	r.End = "\033[0m"
}
