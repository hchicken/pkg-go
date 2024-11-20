package arraryx

import "reflect"

func equals(a, b interface{}) bool {
	return a == b
}

func contains(slice interface{}, elem interface{}) bool {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < s.Len(); i++ {
		if equals(s.Index(i).Interface(), elem) {
			return true
		}
	}
	return false
}

// Difference Find the difference between two arrays
func Difference(a, b interface{}) interface{} {
	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)

	if aValue.Kind() != reflect.Slice || bValue.Kind() != reflect.Slice {
		return nil
	}

	resultType := reflect.SliceOf(aValue.Type().Elem())
	result := reflect.MakeSlice(resultType, 0, aValue.Len())

	for i := 0; i < aValue.Len(); i++ {
		elem := aValue.Index(i).Interface()
		if !contains(b, elem) {
			result = reflect.Append(result, reflect.ValueOf(elem))
		}
	}

	return result.Interface()
}

// Union returns the union of two arrays
func Union(a, b interface{}) interface{} {
	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)

	if aValue.Kind() != reflect.Slice || bValue.Kind() != reflect.Slice {
		return nil
	}

	resultType := reflect.SliceOf(aValue.Type().Elem())
	result := reflect.MakeSlice(resultType, 0, aValue.Len()+bValue.Len())

	// 添加第一个切片的所有元素
	for i := 0; i < aValue.Len(); i++ {
		elem := aValue.Index(i).Interface()
		if !contains(result.Interface(), elem) {
			result = reflect.Append(result, reflect.ValueOf(elem))
		}
	}

	// 添加第二个切片的所有不重复元素
	for i := 0; i < bValue.Len(); i++ {
		elem := bValue.Index(i).Interface()
		if !contains(result.Interface(), elem) {
			result = reflect.Append(result, reflect.ValueOf(elem))
		}
	}

	return result.Interface()
}

// Intersection returns the intersection of two arrays
func Intersection(a, b interface{}) interface{} {
	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)

	if aValue.Kind() != reflect.Slice || bValue.Kind() != reflect.Slice {
		return nil
	}

	resultType := reflect.SliceOf(aValue.Type().Elem())
	result := reflect.MakeSlice(resultType, 0, aValue.Len())

	// 添加在两个切片中都存在的元素
	for i := 0; i < aValue.Len(); i++ {
		elem := aValue.Index(i).Interface()
		if contains(b, elem) && !contains(result.Interface(), elem) {
			result = reflect.Append(result, reflect.ValueOf(elem))
		}
	}

	return result.Interface()
}
