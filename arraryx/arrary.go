package arraryx

// RemoveDuplicates 去重
func RemoveDuplicates[T any](slice []T) []T {
	keys := make(map[interface{}]bool)
	var list []T
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// RemoveElement 删除数组中的指定元素
func RemoveElement[T comparable](slice []T, element T) []T {
	var result []T
	for _, v := range slice {
		if v != element {
			result = append(result, v)
		}
	}
	return result
}
