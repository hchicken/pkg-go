package logx

import "github.com/sirupsen/logrus"

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
