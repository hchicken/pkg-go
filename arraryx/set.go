package arraryx

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
