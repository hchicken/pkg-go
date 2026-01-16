package logx

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
)

// ErrorFileHook 错误日志单独输出到文件的 Hook
type ErrorFileHook struct {
	writer    io.Writer
	formatter logrus.Formatter
	levels    []logrus.Level
	mu        sync.Mutex
}

// NewErrorFileHook 创建错误日志文件 Hook
func NewErrorFileHook(filePath string, formatter logrus.Formatter) *ErrorFileHook {
	// 创建目录
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil
	}

	// 创建文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil
	}

	if formatter == nil {
		formatter = &TextFormatter{
			CallerShortFile: true,
		}
	}

	return &ErrorFileHook{
		writer:    file,
		formatter: formatter,
		levels:    []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel},
	}
}

// Levels 返回 Hook 触发的日志级别
func (h *ErrorFileHook) Levels() []logrus.Level {
	return h.levels
}

// Fire 触发 Hook
func (h *ErrorFileHook) Fire(entry *logrus.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	line, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = h.writer.Write(line)
	return err
}

// Close 关闭 Hook
func (h *ErrorFileHook) Close() error {
	if closer, ok := h.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// LevelFilterHook 日志级别过滤 Hook
type LevelFilterHook struct {
	writer    io.Writer
	formatter logrus.Formatter
	minLevel  logrus.Level
	maxLevel  logrus.Level
	mu        sync.Mutex
}

// NewLevelFilterHook 创建日志级别过滤 Hook
func NewLevelFilterHook(writer io.Writer, formatter logrus.Formatter, minLevel, maxLevel logrus.Level) *LevelFilterHook {
	if formatter == nil {
		formatter = &TextFormatter{}
	}
	return &LevelFilterHook{
		writer:    writer,
		formatter: formatter,
		minLevel:  minLevel,
		maxLevel:  maxLevel,
	}
}

// Levels 返回 Hook 触发的日志级别
func (h *LevelFilterHook) Levels() []logrus.Level {
	var levels []logrus.Level
	for l := h.maxLevel; l <= h.minLevel; l++ {
		levels = append(levels, l)
	}
	return levels
}

// Fire 触发 Hook
func (h *LevelFilterHook) Fire(entry *logrus.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	line, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = h.writer.Write(line)
	return err
}

// CallbackHook 回调 Hook，用于自定义处理
type CallbackHook struct {
	callback func(entry *logrus.Entry) error
	levels   []logrus.Level
}

// NewCallbackHook 创建回调 Hook
func NewCallbackHook(callback func(entry *logrus.Entry) error, levels ...logrus.Level) *CallbackHook {
	if len(levels) == 0 {
		levels = logrus.AllLevels
	}
	return &CallbackHook{
		callback: callback,
		levels:   levels,
	}
}

// Levels 返回 Hook 触发的日志级别
func (h *CallbackHook) Levels() []logrus.Level {
	return h.levels
}

// Fire 触发 Hook
func (h *CallbackHook) Fire(entry *logrus.Entry) error {
	if h.callback != nil {
		return h.callback(entry)
	}
	return nil
}

// MetricsHook 指标收集 Hook
type MetricsHook struct {
	counter map[logrus.Level]uint64
	mu      sync.RWMutex
}

// NewMetricsHook 创建指标收集 Hook
func NewMetricsHook() *MetricsHook {
	return &MetricsHook{
		counter: make(map[logrus.Level]uint64),
	}
}

// Levels 返回 Hook 触发的日志级别
func (h *MetricsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire 触发 Hook
func (h *MetricsHook) Fire(entry *logrus.Entry) error {
	h.mu.Lock()
	h.counter[entry.Level]++
	h.mu.Unlock()
	return nil
}

// GetCount 获取指定级别的日志计数
func (h *MetricsHook) GetCount(level logrus.Level) uint64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.counter[level]
}

// GetAllCounts 获取所有级别的日志计数
func (h *MetricsHook) GetAllCounts() map[logrus.Level]uint64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make(map[logrus.Level]uint64, len(h.counter))
	for k, v := range h.counter {
		result[k] = v
	}
	return result
}

// Reset 重置计数
func (h *MetricsHook) Reset() {
	h.mu.Lock()
	h.counter = make(map[logrus.Level]uint64)
	h.mu.Unlock()
}

// String 返回计数的字符串表示
func (h *MetricsHook) String() string {
	counts := h.GetAllCounts()
	return fmt.Sprintf("LogMetrics{panic:%d, fatal:%d, error:%d, warn:%d, info:%d, debug:%d, trace:%d}",
		counts[logrus.PanicLevel],
		counts[logrus.FatalLevel],
		counts[logrus.ErrorLevel],
		counts[logrus.WarnLevel],
		counts[logrus.InfoLevel],
		counts[logrus.DebugLevel],
		counts[logrus.TraceLevel],
	)
}
