package logx

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"github.com/hchicken/pkg-go/stringx"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LoggerIns 日志实例封装
type LoggerIns struct {
	*logrus.Logger
	opts    *LoggerOptions
	counter uint64 // 用于采样计数
}

// Logger 日志管理器
type Logger struct {
	sync.RWMutex
	opts   LoggerOptions
	closed bool
}

// NewLogger 新建一个日志记录器
func NewLogger(opts ...LoggerOption) *Logger {
	opt := newOptions(opts...)
	return &Logger{opts: opt}
}

// new 创建新的日志实例
func (ls *Logger) new(file string) *LoggerIns {
	l := logrus.New()

	// 仅控制台输出模式
	if ls.opts.writerType == ConsoleWriter {
		l.SetOutput(os.Stdout)
		ls.setLoggerOptions(l)
		return &LoggerIns{Logger: l, opts: &ls.opts}
	}

	// 创建目录
	if ls.opts.path != "" {
		if _, err := os.Stat(ls.opts.path); os.IsNotExist(err) {
			if err := os.MkdirAll(ls.opts.path, os.ModePerm); err != nil {
				l.SetOutput(os.Stdout)
				l.Warnf("failed to create log directory %s: %v, fallback to stdout", ls.opts.path, err)
				ls.setLoggerOptions(l)
				return &LoggerIns{Logger: l, opts: &ls.opts}
			}
		}
	}

	// 设置文件路径
	fPath := path.Join(ls.opts.path, file)

	// 获取文件 writer
	var fileWriter io.Writer

	switch ls.opts.writerType {
	case RotateWriter:
		var newFilename string
		fileArr := strings.Split(fPath, ".")
		if len(fileArr) > 1 {
			newFilenameArr := strings.Join(fileArr[:len(fileArr)-1], ".")
			newFilename = fmt.Sprintf("%v_%v.%v", newFilenameArr, "%Y%m%d", fileArr[len(fileArr)-1])
		} else {
			newFilename = fmt.Sprintf("%v_%v", "%Y%m%d", fPath)
		}
		writer, err := rotatelogs.New(
			newFilename,
			rotatelogs.WithRotationTime(time.Duration(ls.opts.rotate)*time.Hour),
			rotatelogs.WithMaxAge(time.Duration(ls.opts.maxAge)*24*time.Hour),
		)
		if err != nil {
			l.SetOutput(os.Stdout)
			l.Warnf("failed to create rotate writer: %v, fallback to stdout", err)
		} else {
			fileWriter = writer
		}
	case LumberjackWriter:
		fileWriter = &lumberjack.Logger{
			Filename:   fPath,
			MaxSize:    ls.opts.maxSize,
			MaxBackups: ls.opts.maxBackups,
			MaxAge:     ls.opts.maxAge,
			Compress:   true,
		}
	default:
		l.SetOutput(os.Stdout)
	}

	// 设置输出
	if fileWriter != nil {
		if ls.opts.enableConsole {
			l.SetOutput(io.MultiWriter(fileWriter, os.Stdout))
		} else {
			l.SetOutput(fileWriter)
		}
	}

	ls.setLoggerOptions(l)
	return &LoggerIns{Logger: l, opts: &ls.opts}
}

// setLoggerOptions 设置日志基础配置
func (ls *Logger) setLoggerOptions(l *logrus.Logger) {
	l.SetReportCaller(ls.opts.reportCaller)
	l.SetLevel(ls.opts.logLevel)
	l.SetFormatter(ls.opts.formatter)

	// 添加 hooks
	for _, hook := range ls.opts.hooks {
		l.AddHook(hook)
	}

	// 如果配置了错误日志分离，添加 ErrorHook
	if ls.opts.errorFilePath != "" {
		l.AddHook(NewErrorFileHook(ls.opts.errorFilePath, ls.opts.formatter))
	}
}

// Get 获取指定文件的日志实例
func (ls *Logger) Get(file string) *LoggerIns {
	key := stringx.StrToMd5(file)

	ls.RLock()
	if ls.closed {
		ls.RUnlock()
		return nil
	}
	if l, ok := ls.opts.table[key]; ok {
		ls.RUnlock()
		return l
	}
	ls.RUnlock()

	ls.Lock()
	defer ls.Unlock()

	if ls.closed {
		return nil
	}

	if l, ok := ls.opts.table[key]; ok {
		return l
	}

	l := ls.new(file)
	ls.opts.table[key] = l
	return l
}

// Console 获取一个仅输出到控制台的日志实例
func (ls *Logger) Console() *LoggerIns {
	const consoleKey = "__console__"

	ls.RLock()
	if l, ok := ls.opts.table[consoleKey]; ok {
		ls.RUnlock()
		return l
	}
	ls.RUnlock()

	ls.Lock()
	defer ls.Unlock()

	if l, ok := ls.opts.table[consoleKey]; ok {
		return l
	}

	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetReportCaller(ls.opts.reportCaller)
	l.SetLevel(ls.opts.logLevel)
	l.SetFormatter(ls.opts.formatter)

	ins := &LoggerIns{Logger: l, opts: &ls.opts}
	ls.opts.table[consoleKey] = ins
	return ins
}

// SetLevel 动态修改日志级别
func (ls *Logger) SetLevel(level logrus.Level) {
	ls.Lock()
	defer ls.Unlock()

	ls.opts.logLevel = level
	for _, l := range ls.opts.table {
		l.SetLevel(level)
	}
}

// GetLevel 获取当前日志级别
func (ls *Logger) GetLevel() logrus.Level {
	ls.RLock()
	defer ls.RUnlock()
	return ls.opts.logLevel
}

// AddHook 动态添加 Hook
func (ls *Logger) AddHook(hook logrus.Hook) {
	ls.Lock()
	defer ls.Unlock()

	ls.opts.hooks = append(ls.opts.hooks, hook)
	for _, l := range ls.opts.table {
		l.AddHook(hook)
	}
}

// Close 关闭日志管理器，释放资源
func (ls *Logger) Close() error {
	ls.Lock()
	defer ls.Unlock()

	if ls.closed {
		return nil
	}
	ls.closed = true

	var errs []error
	for key, l := range ls.opts.table {
		// 不关闭标准输出和标准错误
		if l.Out == os.Stdout || l.Out == os.Stderr {
			continue
		}
		if closer, ok := l.Out.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				errs = append(errs, fmt.Errorf("close %s: %w", key, err))
			}
		}
	}

	ls.opts.table = make(map[string]*LoggerIns)

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// WithField 添加单个字段
func (l *LoggerIns) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithFields 添加多个字段
func (l *LoggerIns) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

// WithError 添加错误字段
func (l *LoggerIns) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// Sampled 返回一个支持采样的日志实例
func (l *LoggerIns) Sampled() *SampledLogger {
	return &SampledLogger{
		LoggerIns: l,
		rate:      uint64(l.opts.samplingRate),
	}
}

// SampledLogger 支持采样的日志实例
type SampledLogger struct {
	*LoggerIns
	rate uint64
}

// shouldLog 判断是否应该记录日志
func (s *SampledLogger) shouldLog() bool {
	if s.rate == 0 {
		return true
	}
	return atomic.AddUint64(&s.counter, 1)%s.rate == 0
}

// Info 采样 Info 日志
func (s *SampledLogger) Info(args ...interface{}) {
	if s.shouldLog() {
		s.LoggerIns.Info(args...)
	}
}

// Infof 采样 Infof 日志
func (s *SampledLogger) Infof(format string, args ...interface{}) {
	if s.shouldLog() {
		s.LoggerIns.Infof(format, args...)
	}
}

// Debug 采样 Debug 日志
func (s *SampledLogger) Debug(args ...interface{}) {
	if s.shouldLog() {
		s.LoggerIns.Debug(args...)
	}
}

// Debugf 采样 Debugf 日志
func (s *SampledLogger) Debugf(format string, args ...interface{}) {
	if s.shouldLog() {
		s.LoggerIns.Debugf(format, args...)
	}
}

// Warn 采样 Warn 日志（警告不采样，全部记录）
func (s *SampledLogger) Warn(args ...interface{}) {
	s.LoggerIns.Warn(args...)
}

// Error 采样 Error 日志（错误不采样，全部记录）
func (s *SampledLogger) Error(args ...interface{}) {
	s.LoggerIns.Error(args...)
}
