package log

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/hchicken/pkg-go/stringx"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

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

// newLogger 新建一个日志记录器
func NewLogger(opts ...LoggerOption) Logger {
	// 初始化配置
	opt := LoggerOptions{
		formatter: &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyMsg:  "log_message",
				logrus.FieldKeyTime: "log_asctime",
				//logrus.FieldKeyFile:  "log_file",
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

	// 设置对应的时间
	fPath := path.Join(ls.opts.path, file)
	var newFilename string
	fileArr := strings.Split(fPath, ".")
	if len(fileArr) > 1 {
		newFilenameArr := strings.Join(fileArr[:len(fileArr)-1], ".")
		newFilename = fmt.Sprintf("%v_%v.%v", newFilenameArr, "%Y%m%d", fileArr[len(fileArr)-1])
	} else {
		newFilename = fmt.Sprintf("%v_%v", "%Y%m%d", fPath)
	}

	writer, _ := rotatelogs.New(
		newFilename,
		// 设置日志切割时间间隔(1天)(隔多久分割一次)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	l.SetLevel(logrus.InfoLevel)      // 设置打印日志等级
	l.SetFormatter(ls.opts.formatter) // 打印日志
	l.SetOutput(writer)

	return &LoggerIns{l}
}

// get ...
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
