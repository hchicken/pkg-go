package logx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// 预分配的 buffer pool，减少内存分配
var bufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 256))
	},
}

// ANSI 颜色码常量
const (
	colorReset  = "\x1b[0m"
	colorRed    = "\x1b[31m"
	colorGreen  = "\x1b[32m"
	colorYellow = "\x1b[33m"
	colorBlue   = "\x1b[34m"
	colorPurple = "\x1b[35m"
	colorCyan   = "\x1b[36m"
	colorGray   = "\x1b[37m"
)

// 日志级别缩写映射，避免重复计算
var levelShortNames = map[logrus.Level]string{
	logrus.PanicLevel: "PANI",
	logrus.FatalLevel: "FATA",
	logrus.ErrorLevel: "ERRO",
	logrus.WarnLevel:  "WARN",
	logrus.InfoLevel:  "INFO",
	logrus.DebugLevel: "DEBU",
	logrus.TraceLevel: "TRAC",
}

// 日志级别颜色映射
var levelColors = map[logrus.Level]string{
	logrus.PanicLevel: colorRed,
	logrus.FatalLevel: colorRed,
	logrus.ErrorLevel: colorRed,
	logrus.WarnLevel:  colorYellow,
	logrus.InfoLevel:  colorGreen,
	logrus.DebugLevel: colorGray,
	logrus.TraceLevel: colorGray,
}

// TextFormatter logrus formatter 实现
type TextFormatter struct {
	// FieldsOrder 字段输出顺序，默认按字母排序
	FieldsOrder []string

	// TimestampFormat 时间格式，默认: "2006-01-02 15:04:05"
	TimestampFormat string

	// HideKeys 隐藏字段键名，仅显示值 [value] 而非 [key:value]
	HideKeys bool

	// NoColors 禁用颜色输出
	NoColors bool

	// NoFieldsColors 仅对日志级别着色，不对字段着色
	NoFieldsColors bool

	// NoFieldsSpace 字段之间不添加空格
	NoFieldsSpace bool

	// ShowFullLevel 显示完整级别名称 [WARNING] 而非 [WARN]
	ShowFullLevel bool

	// NoUppercaseLevel 级别名称不转大写
	NoUppercaseLevel bool

	// TrimMessages 去除消息首尾空白
	TrimMessages bool

	// CallerFirst 调用者信息放在级别之前
	CallerFirst bool

	// CallerShortFile 仅显示文件名而非完整路径
	CallerShortFile bool

	// CallerSkipFunction 不显示函数名
	CallerSkipFunction bool

	// CustomCallerFormatter 自定义调用者格式化函数
	CustomCallerFormatter func(*runtime.Frame) string
}

// Format 格式化日志条目
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 从 pool 获取 buffer
	b := bufferPool.Get().(*bytes.Buffer)
	b.Reset()
	defer bufferPool.Put(b)

	// 获取时间格式
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.DateTime
	}

	// 写入时间
	b.WriteString(entry.Time.Format(timestampFormat))

	// 调用者信息优先
	if f.CallerFirst {
		f.writeCaller(b, entry)
	}

	// 写入级别
	f.writeLevel(b, entry)

	// 写入字段
	if len(entry.Data) > 0 {
		if len(f.FieldsOrder) > 0 {
			f.writeOrderedFields(b, entry)
		} else {
			f.writeFields(b, entry)
		}
	}

	if f.NoFieldsSpace {
		b.WriteByte(' ')
	}

	// 重置颜色
	if !f.NoColors && !f.NoFieldsColors {
		b.WriteString(colorReset)
	}

	// 写入消息
	if f.TrimMessages {
		b.WriteString(strings.TrimSpace(entry.Message))
	} else {
		b.WriteString(entry.Message)
	}

	// 调用者信息在后
	if !f.CallerFirst {
		f.writeCaller(b, entry)
	}

	b.WriteByte('\n')

	// 复制结果，因为 buffer 会被回收
	result := make([]byte, b.Len())
	copy(result, b.Bytes())
	return result, nil
}

// writeLevel 写入日志级别
func (f *TextFormatter) writeLevel(b *bytes.Buffer, entry *logrus.Entry) {
	// 获取颜色
	if !f.NoColors {
		if color, ok := levelColors[entry.Level]; ok {
			b.WriteString(color)
		}
	}

	b.WriteString(" [")

	if f.ShowFullLevel {
		level := entry.Level.String()
		if !f.NoUppercaseLevel {
			level = strings.ToUpper(level)
		}
		b.WriteString(level)
	} else {
		// 使用预计算的缩写
		if shortName, ok := levelShortNames[entry.Level]; ok {
			if f.NoUppercaseLevel {
				b.WriteString(strings.ToLower(shortName))
			} else {
				b.WriteString(shortName)
			}
		} else {
			// 兜底处理
			level := entry.Level.String()
			if !f.NoUppercaseLevel {
				level = strings.ToUpper(level)
			}
			if len(level) >= 4 {
				b.WriteString(level[:4])
			} else {
				b.WriteString(level)
			}
		}
	}

	b.WriteByte(']')

	if !f.NoFieldsSpace {
		b.WriteByte(' ')
	}

	// 仅级别着色时重置颜色
	if !f.NoColors && f.NoFieldsColors {
		b.WriteString(colorReset)
	}
}

// writeCaller 写入调用者信息
func (f *TextFormatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if !entry.HasCaller() {
		return
	}

	if f.CustomCallerFormatter != nil {
		b.WriteString(f.CustomCallerFormatter(entry.Caller))
		return
	}

	b.WriteString(" (")

	// 文件路径
	file := entry.Caller.File
	if f.CallerShortFile {
		file = filepath.Base(file)
	}
	b.WriteString(file)
	b.WriteByte(':')
	b.WriteString(strconv.Itoa(entry.Caller.Line))

	// 函数名
	if !f.CallerSkipFunction {
		b.WriteByte(' ')
		// 仅保留函数名，去除包路径
		funcName := entry.Caller.Function
		if idx := strings.LastIndex(funcName, "."); idx != -1 {
			funcName = funcName[idx+1:]
		}
		b.WriteString(funcName)
	}

	b.WriteByte(')')
}

// writeFields 写入字段（按字母排序）
func (f *TextFormatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	fields := make([]string, 0, len(entry.Data))
	for field := range entry.Data {
		fields = append(fields, field)
	}
	sort.Strings(fields)

	for _, field := range fields {
		f.writeField(b, entry, field)
	}
}

// writeOrderedFields 按指定顺序写入字段
func (f *TextFormatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	written := make(map[string]struct{}, len(f.FieldsOrder))

	// 先按指定顺序写入
	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			written[field] = struct{}{}
			f.writeField(b, entry, field)
		}
	}

	// 剩余字段按字母排序写入
	remaining := len(entry.Data) - len(written)
	if remaining > 0 {
		notFound := make([]string, 0, remaining)
		for field := range entry.Data {
			if _, ok := written[field]; !ok {
				notFound = append(notFound, field)
			}
		}
		sort.Strings(notFound)

		for _, field := range notFound {
			f.writeField(b, entry, field)
		}
	}
}

// writeField 写入单个字段
func (f *TextFormatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	b.WriteByte('[')

	if !f.HideKeys {
		b.WriteString(field)
		b.WriteByte(':')
	}

	// 格式化值
	f.writeValue(b, entry.Data[field])

	b.WriteByte(']')

	if !f.NoFieldsSpace {
		b.WriteByte(' ')
	}
}

// writeValue 写入字段值
func (f *TextFormatter) writeValue(b *bytes.Buffer, value interface{}) {
	switch v := value.(type) {
	case string:
		b.WriteString(v)
	case int:
		b.WriteString(strconv.Itoa(v))
	case int64:
		b.WriteString(strconv.FormatInt(v, 10))
	case float64:
		b.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
	case bool:
		b.WriteString(strconv.FormatBool(v))
	case error:
		b.WriteString(v.Error())
	default:
		// 复杂类型使用 JSON 序列化
		if jsonData, err := json.Marshal(v); err == nil {
			b.Write(jsonData)
		} else {
			fmt.Fprintf(b, "%v", v)
		}
	}
}
