package logx

import (
	"github.com/sirupsen/logrus"
)

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel logrus.Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

const (
	RotateWriter     = "RotateWriter"
	LumberjackWriter = "LumberjackWriter"
)

// LoggerOption ...
type LoggerOption func(*LoggerOptions)

// LoggerOptions ...
type LoggerOptions struct {
	path string
	file string

	reportCaller bool
	logLevel     logrus.Level

	formatter logrus.Formatter
	table     map[string]*LoggerIns

	writerType string
	maxSize    int
	maxBackups int
	maxAge     int
	rotate     int
}

// newOptions ...
func newOptions(opts ...LoggerOption) LoggerOptions {
	// 初始化配置
	opt := LoggerOptions{
		reportCaller: true,
		logLevel:     InfoLevel,
		formatter: &TextFormatter{
			HideKeys: true,
		},
		table:      make(map[string]*LoggerIns, 5),
		writerType: RotateWriter,
		maxSize:    2048, // Max size in MB (1GB)
		maxBackups: 7,    // Number of old files to retain
		maxAge:     7,    // MaxAge
		rotate:     24,   // 日志切割时间
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

// ReportCaller 是否打印调用信息
func ReportCaller(ok bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.reportCaller = ok
	}
}

// LogLevel 日志等级
func LogLevel(level logrus.Level) LoggerOption {
	return func(o *LoggerOptions) {
		o.logLevel = level
	}
}

// WriterType Writer类型
func WriterType(writerType string) LoggerOption {
	return func(o *LoggerOptions) {
		o.writerType = writerType
	}
}

// MaxSize 设置日志最大容量,单位:MB
func MaxSize(maxSize int) LoggerOption {
	return func(o *LoggerOptions) {
		o.maxSize = maxSize
	}
}

// MaxBackups Number of old files to retain
func MaxBackups(maxBackups int) LoggerOption {
	return func(o *LoggerOptions) {
		o.maxBackups = maxBackups
	}
}

// MaxAge 日志最大保留时间,单位:day
func MaxAge(maxAge int) LoggerOption {
	return func(o *LoggerOptions) {
		o.maxAge = maxAge
	}
}

// Rotate 日志切割时间,单位:h
func Rotate(rotate int) LoggerOption {
	return func(o *LoggerOptions) {
		o.rotate = rotate
	}
}
