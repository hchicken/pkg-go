package filex

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// 最终方案-全兼容
func GetCurrentAbPath() (string, error) {
	dir, err := getCurrentAbPathByExecutable()
	if err != nil {
		return "", err
	}
	tmpDir, err := filepath.EvalSymlinks(os.TempDir())
	if err != nil {
		return "", err
	}
	if strings.Contains(dir, tmpDir) {
		return getCurrentAbPathByCaller(), nil
	}
	return dir, nil
}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(filepath.Dir(exePath))
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	_, filename, _, ok := runtime.Caller(2)
	if ok {
		return path.Dir(filename)
	}
	return ""
}
