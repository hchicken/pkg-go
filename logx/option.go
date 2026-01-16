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

// Writer 类型常量
const (
	RotateWriter     = "RotateWriter"
	LumberjackWriter = "LumberjackWriter"
	ConsoleWriter    = "ConsoleWriter" // 控制台输出，用于debug
)

// LoggerOption 配置选项函数类型
type LoggerOption func(*LoggerOptions)

// LoggerOptions 日志配置选项
type LoggerOptions struct {
	path string

	reportCaller bool
	logLevel     logrus.Level

	formatter logrus.Formatter
	table     map[string]*LoggerIns

	writerType    string
	maxSize       int
	maxBackups    int
	maxAge        int
	rotate        int
	enableConsole bool // 是否同时输出到控制台

	// 扩展配置
	hooks         []logrus.Hook // Hook 列表
	errorFilePath string        // 错误日志单独输出路径
	asyncEnabled  bool          // 是否开启异步日志
	asyncBuffer   int           // 异步缓冲区大小
	samplingRate  int           // 采样率，0表示不采样，N表示每N条记录1条
}

// newOptions 创建默认配置
func newOptions(opts ...LoggerOption) LoggerOptions {
	opt := LoggerOptions{
		reportCaller: true,
		logLevel:     InfoLevel,
		formatter: &TextFormatter{
			HideKeys:        true,
			CallerShortFile: true,
		},
		table:         make(map[string]*LoggerIns, 5),
		writerType:    RotateWriter,
		maxSize:       2048, // Max size in MB
		maxBackups:    7,    // Number of old files to retain
		maxAge:        7,    // MaxAge in days
		rotate:        24,   // 日志切割时间(小时)
		enableConsole: false,
		asyncBuffer:   4096,
		samplingRate:  0,
	}
	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// LoggerPath 设置日志路径
func LoggerPath(path string) LoggerOption {
	return func(o *LoggerOptions) {
		o.path = path
	}
}

// Formatter 设置日志格式化器
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

// LogLevel 设置日志等级
func LogLevel(level logrus.Level) LoggerOption {
	return func(o *LoggerOptions) {
		o.logLevel = level
	}
}

// WriterType 设置 Writer 类型
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

// MaxBackups 设置保留的旧文件数量
func MaxBackups(maxBackups int) LoggerOption {
	return func(o *LoggerOptions) {
		o.maxBackups = maxBackups
	}
}

// MaxAge 设置日志最大保留时间,单位:day
func MaxAge(maxAge int) LoggerOption {
	return func(o *LoggerOptions) {
		o.maxAge = maxAge
	}
}

// Rotate 设置日志切割时间,单位:h
func Rotate(rotate int) LoggerOption {
	return func(o *LoggerOptions) {
		o.rotate = rotate
	}
}

// EnableConsole 开启控制台输出，用于debug调试
func EnableConsole(enable bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.enableConsole = enable
	}
}

// WithHooks 添加 Hook
func WithHooks(hooks ...logrus.Hook) LoggerOption {
	return func(o *LoggerOptions) {
		o.hooks = append(o.hooks, hooks...)
	}
}

// ErrorFilePath 设置错误日志单独输出路径
func ErrorFilePath(path string) LoggerOption {
	return func(o *LoggerOptions) {
		o.errorFilePath = path
	}
}

// EnableAsync 开启异步日志
func EnableAsync(enable bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.asyncEnabled = enable
	}
}

// AsyncBuffer 设置异步缓冲区大小
func AsyncBuffer(size int) LoggerOption {
	return func(o *LoggerOptions) {
		o.asyncBuffer = size
	}
}

// SamplingRate 设置采样率，0表示不采样，N表示每N条记录1条
func SamplingRate(rate int) LoggerOption {
	return func(o *LoggerOptions) {
		o.samplingRate = rate
	}
}
