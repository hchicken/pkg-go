package logx

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// AsyncLogger 异步日志实例
type AsyncLogger struct {
	*LoggerIns
	ch       chan *logEntry
	done     chan struct{}
	wg       sync.WaitGroup
	closed   bool
	closeMu  sync.Mutex
	flushMu  sync.Mutex
}

// logEntry 日志条目
type logEntry struct {
	level   logrus.Level
	message string
	fields  logrus.Fields
	time    time.Time
}

// NewAsyncLogger 创建异步日志实例
func NewAsyncLogger(ins *LoggerIns, bufferSize int) *AsyncLogger {
	if bufferSize <= 0 {
		bufferSize = 4096
	}

	a := &AsyncLogger{
		LoggerIns: ins,
		ch:        make(chan *logEntry, bufferSize),
		done:      make(chan struct{}),
	}

	a.wg.Add(1)
	go a.worker()

	return a
}

// worker 异步写入协程
func (a *AsyncLogger) worker() {
	defer a.wg.Done()

	for {
		select {
		case entry := <-a.ch:
			if entry != nil {
				a.write(entry)
			}
		case <-a.done:
			// 处理剩余日志
			for {
				select {
				case entry := <-a.ch:
					if entry != nil {
						a.write(entry)
					}
				default:
					return
				}
			}
		}
	}
}

// write 实际写入日志
func (a *AsyncLogger) write(entry *logEntry) {
	e := a.LoggerIns.WithFields(entry.fields)
	e.Time = entry.time

	switch entry.level {
	case logrus.DebugLevel:
		e.Debug(entry.message)
	case logrus.InfoLevel:
		e.Info(entry.message)
	case logrus.WarnLevel:
		e.Warn(entry.message)
	case logrus.ErrorLevel:
		e.Error(entry.message)
	case logrus.FatalLevel:
		e.Fatal(entry.message)
	case logrus.PanicLevel:
		e.Panic(entry.message)
	default:
		e.Info(entry.message)
	}
}

// send 发送日志到缓冲区
func (a *AsyncLogger) send(level logrus.Level, message string, fields logrus.Fields) {
	a.closeMu.Lock()
	if a.closed {
		a.closeMu.Unlock()
		// 已关闭，直接同步写入
		a.write(&logEntry{level: level, message: message, fields: fields, time: time.Now()})
		return
	}
	a.closeMu.Unlock()

	entry := &logEntry{
		level:   level,
		message: message,
		fields:  fields,
		time:    time.Now(),
	}

	select {
	case a.ch <- entry:
	default:
		// 缓冲区满，直接同步写入，避免丢失
		a.write(entry)
	}
}

// Debug 异步输出 debug 级别日志
func (a *AsyncLogger) Debug(args ...interface{}) {
	a.send(logrus.DebugLevel, formatArgs(args...), nil)
}

// Debugf 异步格式化输出 debug 级别日志
func (a *AsyncLogger) Debugf(format string, args ...interface{}) {
	a.send(logrus.DebugLevel, formatArgsf(format, args...), nil)
}

// Info 异步输出 info 级别日志
func (a *AsyncLogger) Info(args ...interface{}) {
	a.send(logrus.InfoLevel, formatArgs(args...), nil)
}

// Infof 异步格式化输出 info 级别日志
func (a *AsyncLogger) Infof(format string, args ...interface{}) {
	a.send(logrus.InfoLevel, formatArgsf(format, args...), nil)
}

// Warn 异步输出 warn 级别日志
func (a *AsyncLogger) Warn(args ...interface{}) {
	a.send(logrus.WarnLevel, formatArgs(args...), nil)
}

// Warnf 异步格式化输出 warn 级别日志
func (a *AsyncLogger) Warnf(format string, args ...interface{}) {
	a.send(logrus.WarnLevel, formatArgsf(format, args...), nil)
}

// Error 异步输出 error 级别日志
func (a *AsyncLogger) Error(args ...interface{}) {
	a.send(logrus.ErrorLevel, formatArgs(args...), nil)
}

// Errorf 异步格式化输出 error 级别日志
func (a *AsyncLogger) Errorf(format string, args ...interface{}) {
	a.send(logrus.ErrorLevel, formatArgsf(format, args...), nil)
}

// WithFields 返回带字段的异步日志实例
func (a *AsyncLogger) WithFields(fields logrus.Fields) *AsyncLoggerEntry {
	return &AsyncLoggerEntry{
		async:  a,
		fields: fields,
	}
}

// Flush 刷新缓冲区，等待所有日志写入完成
func (a *AsyncLogger) Flush() {
	a.flushMu.Lock()
	defer a.flushMu.Unlock()

	// 等待缓冲区清空
	for len(a.ch) > 0 {
		time.Sleep(10 * time.Millisecond)
	}
}

// Close 关闭异步日志
func (a *AsyncLogger) Close() error {
	a.closeMu.Lock()
	if a.closed {
		a.closeMu.Unlock()
		return nil
	}
	a.closed = true
	a.closeMu.Unlock()

	close(a.done)
	a.wg.Wait()
	return nil
}

// AsyncLoggerEntry 带字段的异步日志条目
type AsyncLoggerEntry struct {
	async  *AsyncLogger
	fields logrus.Fields
}

// Debug 异步输出 debug 级别日志
func (e *AsyncLoggerEntry) Debug(args ...interface{}) {
	e.async.send(logrus.DebugLevel, formatArgs(args...), e.fields)
}

// Debugf 异步格式化输出 debug 级别日志
func (e *AsyncLoggerEntry) Debugf(format string, args ...interface{}) {
	e.async.send(logrus.DebugLevel, formatArgsf(format, args...), e.fields)
}

// Info 异步输出 info 级别日志
func (e *AsyncLoggerEntry) Info(args ...interface{}) {
	e.async.send(logrus.InfoLevel, formatArgs(args...), e.fields)
}

// Infof 异步格式化输出 info 级别日志
func (e *AsyncLoggerEntry) Infof(format string, args ...interface{}) {
	e.async.send(logrus.InfoLevel, formatArgsf(format, args...), e.fields)
}

// Warn 异步输出 warn 级别日志
func (e *AsyncLoggerEntry) Warn(args ...interface{}) {
	e.async.send(logrus.WarnLevel, formatArgs(args...), e.fields)
}

// Warnf 异步格式化输出 warn 级别日志
func (e *AsyncLoggerEntry) Warnf(format string, args ...interface{}) {
	e.async.send(logrus.WarnLevel, formatArgsf(format, args...), e.fields)
}

// Error 异步输出 error 级别日志
func (e *AsyncLoggerEntry) Error(args ...interface{}) {
	e.async.send(logrus.ErrorLevel, formatArgs(args...), e.fields)
}

// Errorf 异步格式化输出 error 级别日志
func (e *AsyncLoggerEntry) Errorf(format string, args ...interface{}) {
	e.async.send(logrus.ErrorLevel, formatArgsf(format, args...), e.fields)
}

// formatArgs 格式化参数
func formatArgs(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	if len(args) == 1 {
		if s, ok := args[0].(string); ok {
			return s
		}
		return fmt.Sprint(args[0])
	}
	return fmt.Sprint(args...)
}

// formatArgsf 格式化参数
func formatArgsf(format string, args ...interface{}) string {
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}
