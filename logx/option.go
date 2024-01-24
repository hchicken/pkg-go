package logx

import "github.com/sirupsen/logrus"

// LoggerOption ...
type LoggerOption func(*LoggerOptions)

// LoggerOptions ...
type LoggerOptions struct {
	path      string
	file      string
	formatter logrus.Formatter
	table     map[string]*LoggerIns
}

// newOptions ...
func newOptions(opts ...LoggerOption) LoggerOptions {
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
	return opt
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

// Formatter 格式化日志
func Formatter(formatter logrus.Formatter) LoggerOption {
	return func(o *LoggerOptions) {
		o.formatter = formatter
	}
}
