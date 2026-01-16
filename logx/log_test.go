package logx

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger(
		LoggerPath("./test_logs"),
		LogLevel(DebugLevel),
	)
	defer func() {
		_ = logger.Close()
		_ = os.RemoveAll("./test_logs")
	}()

	if logger == nil {
		t.Fatal("NewLogger returned nil")
	}

	if logger.opts.logLevel != DebugLevel {
		t.Errorf("expected log level %v, got %v", DebugLevel, logger.opts.logLevel)
	}
}

func TestLoggerGet(t *testing.T) {
	logger := NewLogger(
		LoggerPath("./test_logs"),
		WriterType(ConsoleWriter),
	)
	defer func() {
		_ = logger.Close()
		_ = os.RemoveAll("./test_logs")
	}()

	ins1 := logger.Get("test.log")
	ins2 := logger.Get("test.log")

	if ins1 != ins2 {
		t.Error("Get should return the same instance for same file")
	}

	ins3 := logger.Get("other.log")
	if ins1 == ins3 {
		t.Error("Get should return different instances for different files")
	}
}

func TestLoggerConsole(t *testing.T) {
	logger := NewLogger(WriterType(ConsoleWriter))
	defer func() {
		_ = logger.Close()
	}()

	console := logger.Console()
	if console == nil {
		t.Fatal("Console returned nil")
	}

	// 再次调用应该返回相同实例
	console2 := logger.Console()
	if console != console2 {
		t.Error("Console should return the same instance")
	}
}

func TestLoggerSetLevel(t *testing.T) {
	logger := NewLogger(WriterType(ConsoleWriter))
	defer func() {
		_ = logger.Close()
	}()

	_ = logger.Get("test.log")

	logger.SetLevel(DebugLevel)

	if logger.GetLevel() != DebugLevel {
		t.Errorf("expected level %v, got %v", DebugLevel, logger.GetLevel())
	}
}

func TestLoggerClose(t *testing.T) {
	logger := NewLogger(WriterType(ConsoleWriter))
	_ = logger.Get("test.log")

	err := logger.Close()
	if err != nil {
		t.Errorf("Close returned error: %v", err)
	}

	// 再次关闭应该无影响
	err = logger.Close()
	if err != nil {
		t.Errorf("Second Close returned error: %v", err)
	}

	// 关闭后获取应该返回 nil
	ins := logger.Get("test.log")
	if ins != nil {
		t.Error("Get after Close should return nil")
	}
}

func TestLoggerWithFields(t *testing.T) {
	logger := NewLogger(WriterType(ConsoleWriter))
	defer func() {
		_ = logger.Close()
	}()

	ins := logger.Console()

	entry := ins.WithFields(logrus.Fields{
		"key1": "value1",
		"key2": 123,
	})

	if entry.Data["key1"] != "value1" {
		t.Errorf("expected key1 %v, got %v", "value1", entry.Data["key1"])
	}
	if entry.Data["key2"] != 123 {
		t.Errorf("expected key2 %v, got %v", 123, entry.Data["key2"])
	}
}

func TestLoggerConcurrent(t *testing.T) {
	logger := NewLogger(WriterType(ConsoleWriter))
	defer func() {
		_ = logger.Close()
	}()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			ins := logger.Get("test.log")
			ins.Info("concurrent test", n)
		}(i)
	}
	wg.Wait()
}

func TestSampledLogger(t *testing.T) {
	logger := NewLogger(
		WriterType(ConsoleWriter),
		SamplingRate(10), // 每10条记录1条
	)
	defer func() {
		_ = logger.Close()
	}()

	ins := logger.Console()
	sampled := ins.Sampled()

	// 采样率测试
	for i := 0; i < 100; i++ {
		sampled.Info("test")
	}
}

func TestGlobalLogger(t *testing.T) {
	// 测试全局函数
	SetLevel(DebugLevel)

	// 测试全局便捷方法不会 panic
	Info("test info")
	Debug("test debug")
	Warn("test warn")
}

func TestTextFormatter(t *testing.T) {
	formatter := &TextFormatter{
		TimestampFormat: time.RFC3339,
		CallerShortFile: true,
		NoColors:        true,
	}

	entry := &logrus.Entry{
		Logger:  logrus.New(),
		Level:   logrus.InfoLevel,
		Message: "test message",
		Time:    time.Now(),
		Data: logrus.Fields{
			"key": "value",
		},
	}

	output, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	if !strings.Contains(string(output), "test message") {
		t.Error("output should contain message")
	}
	if !strings.Contains(string(output), "[INFO]") {
		t.Error("output should contain level")
	}
}

func TestAsyncLogger(t *testing.T) {
	logger := NewLogger(WriterType(ConsoleWriter))
	defer func() {
		_ = logger.Close()
	}()

	ins := logger.Console()
	async := NewAsyncLogger(ins, 100)
	defer func() {
		_ = async.Close()
	}()

	async.Info("async test")
	async.Infof("async test %d", 123)
	async.WithFields(logrus.Fields{"key": "value"}).Info("with fields")

	async.Flush()
}

func TestErrorFileHook(t *testing.T) {
	tmpDir := t.TempDir()
	errFile := filepath.Join(tmpDir, "error.log")

	hook := NewErrorFileHook(errFile, nil)
	if hook == nil {
		t.Fatal("NewErrorFileHook returned nil")
	}

	levels := hook.Levels()
	if len(levels) != 3 {
		t.Errorf("expected 3 levels, got %d", len(levels))
	}

	_ = hook.Close()
}

func TestMetricsHook(t *testing.T) {
	hook := NewMetricsHook()

	entry := &logrus.Entry{
		Level: logrus.InfoLevel,
	}

	_ = hook.Fire(entry)
	_ = hook.Fire(entry)

	count := hook.GetCount(logrus.InfoLevel)
	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}

	hook.Reset()
	count = hook.GetCount(logrus.InfoLevel)
	if count != 0 {
		t.Errorf("expected count 0 after reset, got %d", count)
	}
}

func TestCallbackHook(t *testing.T) {
	called := false
	hook := NewCallbackHook(func(entry *logrus.Entry) error {
		called = true
		return nil
	}, logrus.InfoLevel)

	entry := &logrus.Entry{
		Level: logrus.InfoLevel,
	}

	_ = hook.Fire(entry)

	if !called {
		t.Error("callback was not called")
	}
}

func TestLoggerFileOutput(t *testing.T) {
	tmpDir := t.TempDir()

	logger := NewLogger(
		LoggerPath(tmpDir),
		WriterType(LumberjackWriter),
		LogLevel(InfoLevel),
	)
	defer func() {
		_ = logger.Close()
	}()

	ins := logger.Get("test.log")
	ins.Info("file output test")

	// 验证文件被创建
	time.Sleep(100 * time.Millisecond)
	files, _ := os.ReadDir(tmpDir)
	if len(files) == 0 {
		t.Error("log file was not created")
	}
}

func TestLoggerAddHook(t *testing.T) {
	logger := NewLogger(WriterType(ConsoleWriter))
	defer func() {
		_ = logger.Close()
	}()

	hook := NewMetricsHook()
	logger.AddHook(hook)

	// 先获取实例（这样 hook 才会被添加到实例上）
	ins := logger.Get("test.log")
	ins.Info("test")

	count := hook.GetCount(logrus.InfoLevel)
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}
}
