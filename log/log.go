package log

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/hchicken/pkg-go/stringx"
	"github.com/sirupsen/logrus"
)

// LoggerIns ...
type LoggerIns struct {
	*logrus.Logger
}

// Logger ...
type Logger struct {
	sync.Mutex
	opts LoggerOptions
}

// LoggerOption ...
type LoggerOption func(*LoggerOptions)

// LoggerOptions ...
type LoggerOptions struct {
	path      string
	file      string
	formatter logrus.Formatter
	table     map[string]*LoggerIns
}

// NewLogger 新建一个日志记录器
func NewLogger(opts ...LoggerOption) Logger {
	// 初始化配置
	opt := LoggerOptions{
		formatter: &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyMsg:   "log_message",
				logrus.FieldKeyTime:  "log_asctime",
				logrus.FieldKeyFile:  "log_file",
				logrus.FieldKeyFunc:  "log_func",
				logrus.FieldKeyLevel: "log_level",
			},
			TimestampFormat: "2006-01-02 15:04:05",
		},
		table: make(map[string]*LoggerIns, 5),
	}
	for _, o := range opts {
		o(&opt)
	}
	return Logger{opts: opt}
}

// new ...
func (ls *Logger) new(file string) *LoggerIns {
	l := logrus.New()

	// 文件生成
	dir := path.Dir(ls.opts.path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}
	fPath := fmt.Sprintf("%v%v", ls.opts.path, file)

	l.SetLevel(logrus.InfoLevel)                                                  // 设置打印日志等级
	l.SetFormatter(ls.opts.formatter)                                             // 打印日志
	l.SetReportCaller(true)                                                       // 设置显示行号
	fp, _ := os.OpenFile(fPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm) // 设置日志输入文件
	l.SetOutput(fp)
	return &LoggerIns{l}
}

// Get ...
func (ls *Logger) Get(file string) *LoggerIns {

	// 加锁
	ls.Lock()
	defer ls.Unlock()

	// 获取对应的logger
	key := stringx.StrToMd5(file)
	if l, ok := ls.opts.table[key]; ok {
		return l
	}
	l := ls.new(file)
	ls.opts.table[key] = l
	return l
}

// LoggerPath 路径
func LoggerPath(path string) LoggerOption {
	return func(o *LoggerOptions) {
		o.path = path
	}
}

// LoggerFile 文件
func LoggerFile(file string) LoggerOption {
	return func(o *LoggerOptions) {
		o.file = file
	}
}
