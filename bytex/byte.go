package bytex

import "fmt"

// KbSize kb单位转换
func KbSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	var unit string
	var conversion float64

	switch {
	case size < MB:
		unit, conversion = "KB", KB
	case size < GB:
		unit, conversion = "MB", MB
	case size < TB:
		unit, conversion = "GB", GB
	case size >= TB:
		unit, conversion = "TB", TB
	default:
		return fmt.Sprintf("%v", size)
	}

	return fmt.Sprintf("%.2f %s", float64(size)/conversion, unit)
}
