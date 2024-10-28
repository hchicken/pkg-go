package arraryx

import "reflect"

// Contain 判断数组里面是否有某个元素
func Contain(obj interface{}, sub interface{}) bool {
	targetValue := reflect.ValueOf(obj)
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == sub {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(sub)).IsValid() {
			return true
		}
	}
	return false
}
