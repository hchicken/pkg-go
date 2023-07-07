package bytex

import "fmt"

// KbSize kb单位转换
func KbSize(size int64) string {
	switch {
	case size < 1024*1024:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(1024))
	case size < 1024*1024*1024:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(1024*1024))
	case size < 1024*1024*1024*1024:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(1024*1024*1024))
	case size < 1024*1024*1024*1024*1024:
		return fmt.Sprintf("%.2f TB", float64(size)/float64(1024*1024*1024*1024))
	default:
		return fmt.Sprintf("%v", size)
	}
}
