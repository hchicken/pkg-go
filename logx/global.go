package logx

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	std     *Logger
	stdOnce sync.Once
	stdMu   sync.RWMutex
)

// 初始化默认 logger
func initStd() {
	stdOnce.Do(func() {
		std = NewLogger(
			WriterType(ConsoleWriter),
			LogLevel(InfoLevel),
		)
	})
}

// Default 获取默认 logger
func Default() *Logger {
	initStd()
	stdMu.RLock()
	defer stdMu.RUnlock()
	return std
}

// SetDefault 设置默认 logger
func SetDefault(l *Logger) {
	initStd()
	stdMu.Lock()
	defer stdMu.Unlock()
	std = l
}

// Get 使用默认 logger 获取指定文件的日志实例
func Get(file string) *LoggerIns {
	return Default().Get(file)
}

// Console 使用默认 logger 获取控制台日志实例
func Console() *LoggerIns {
	return Default().Console()
}

// SetLevel 设置默认 logger 的日志级别
func SetLevel(level logrus.Level) {
	Default().SetLevel(level)
}

// Close 关闭默认 logger
func Close() error {
	return Default().Close()
}

// 全局便捷方法 - 使用控制台输出

// Debug 输出 debug 级别日志
func Debug(args ...interface{}) {
	Default().Console().Debug(args...)
}

// Debugf 格式化输出 debug 级别日志
func Debugf(format string, args ...interface{}) {
	Default().Console().Debugf(format, args...)
}

// Info 输出 info 级别日志
func Info(args ...interface{}) {
	Default().Console().Info(args...)
}

// Infof 格式化输出 info 级别日志
func Infof(format string, args ...interface{}) {
	Default().Console().Infof(format, args...)
}

// Warn 输出 warn 级别日志
func Warn(args ...interface{}) {
	Default().Console().Warn(args...)
}

// Warnf 格式化输出 warn 级别日志
func Warnf(format string, args ...interface{}) {
	Default().Console().Warnf(format, args...)
}

// Error 输出 error 级别日志
func Error(args ...interface{}) {
	Default().Console().Error(args...)
}

// Errorf 格式化输出 error 级别日志
func Errorf(format string, args ...interface{}) {
	Default().Console().Errorf(format, args...)
}

// Fatal 输出 fatal 级别日志并退出
func Fatal(args ...interface{}) {
	Default().Console().Fatal(args...)
}

// Fatalf 格式化输出 fatal 级别日志并退出
func Fatalf(format string, args ...interface{}) {
	Default().Console().Fatalf(format, args...)
}

// Panic 输出 panic 级别日志并 panic
func Panic(args ...interface{}) {
	Default().Console().Panic(args...)
}

// Panicf 格式化输出 panic 级别日志并 panic
func Panicf(format string, args ...interface{}) {
	Default().Console().Panicf(format, args...)
}

// WithField 添加单个字段
func WithField(key string, value interface{}) *logrus.Entry {
	return Default().Console().WithField(key, value)
}

// WithFields 添加多个字段
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Default().Console().WithFields(fields)
}

// WithError 添加错误字段
func WithError(err error) *logrus.Entry {
	return Default().Console().WithError(err)
}
