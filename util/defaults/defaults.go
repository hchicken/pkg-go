// Package defaults TODO
package defaults

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"time"
)

var (
	errInvalidType = errors.New("not a struct pointer")
)

const (
	fieldName = "default"
)

// Set initializes members in a struct referenced by a pointer.
// Maps and slices are initialized by `make` and others primitive types are set with default values.
// `ptr` should be a struct pointer
func Set(ptr interface{}) error {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return errInvalidType
	}

	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	if t.Kind() != reflect.Struct {
		return errInvalidType
	}

	for i := 0; i < t.NumField(); i++ {
		defaultVal := t.Field(i).Tag.Get(fieldName)
		if defaultVal != "-" {
			if err := setField(v.Field(i), defaultVal); err != nil {
				return err
			}
		}
	}
	callSetter(ptr)
	return nil
}

func setField(field reflect.Value, defaultVal string) error {

	// field 是否能设置
	if !field.CanSet() {
		return nil
	}

	if !shouldInitializeField(field, defaultVal) {
		return nil
	}

	isInitial := isInitialValue(field)
	if isInitial {
		err := setDefaultVal(field, defaultVal)
		if err != nil {
			return err
		}
	}

	// 指针/结构体/slice递归设置
	switch field.Kind() {
	case reflect.Ptr:
		if isInitial || field.Elem().Kind() == reflect.Struct {
			err := setField(field.Elem(), defaultVal)
			if err != nil {
				return err
			}
			callSetter(field.Interface())
		}
	case reflect.Struct:
		if err := Set(field.Addr().Interface()); err != nil {
			return err
		}
	case reflect.Slice:
		for j := 0; j < field.Len(); j++ {
			if err := setField(field.Index(j), defaultVal); err != nil {
				return err
			}
		}
	}
	return nil
}

// setDefaultVal 设置默认值
func setDefaultVal(field reflect.Value, defaultVal string) error {

	// 设置字符串
	ok, err := setString(field, defaultVal)
	if ok {
		return err
	}
	// 设置int
	ok, err = setInt(field, defaultVal)
	if ok {
		return err
	}
	ok, err = setUint(field, defaultVal)
	if ok {
		return err
	}
	ok, err = setFloat(field, defaultVal)
	if ok {
		return err
	}

	switch field.Kind() {
	case reflect.Bool:
		if val, err := strconv.ParseBool(defaultVal); err == nil {
			field.Set(reflect.ValueOf(val).Convert(field.Type()))
		}
	case reflect.Slice:
		ref := reflect.New(field.Type())
		ref.Elem().Set(reflect.MakeSlice(field.Type(), 0, 0))
		if defaultVal != "" && defaultVal != "[]" {
			err := json.Unmarshal([]byte(defaultVal), ref.Interface())
			if err != nil {
				return err
			}
		}
		field.Set(ref.Elem().Convert(field.Type()))
	case reflect.Map:
		ref := reflect.New(field.Type())
		ref.Elem().Set(reflect.MakeMap(field.Type()))
		if defaultVal != "" && defaultVal != "{}" {
			err := json.Unmarshal([]byte(defaultVal), ref.Interface())
			if err != nil {
				return err
			}
		}
		field.Set(ref.Elem().Convert(field.Type()))
	case reflect.Struct:
		if defaultVal != "" && defaultVal != "{}" {
			err := json.Unmarshal([]byte(defaultVal), field.Addr().Interface())
			if err != nil {
				return err
			}
		}
	case reflect.Ptr:
		field.Set(reflect.New(field.Type().Elem()))
	}

	return nil
}

// setString 设置int
func setString(field reflect.Value, defaultVal string) (bool, error) {
	switch field.Kind() {
	case reflect.Bool:
		if val, err := strconv.ParseBool(defaultVal); err == nil {
			field.Set(reflect.ValueOf(val).Convert(field.Type()))
		}
		return true, nil
	case reflect.String:
		field.Set(reflect.ValueOf(defaultVal).Convert(field.Type()))
		return true, nil
	}
	return false, nil
}

// setInt 设置int
func setInt(field reflect.Value, defaultVal string) (bool, error) {
	switch field.Kind() {
	case reflect.Int:
		if val, err := strconv.ParseInt(defaultVal, 0, strconv.IntSize); err == nil {
			field.Set(reflect.ValueOf(int(val)).Convert(field.Type()))
		}
		return true, nil
	case reflect.Int8:
		if val, err := strconv.ParseInt(defaultVal, 0, 8); err == nil {
			field.Set(reflect.ValueOf(int8(val)).Convert(field.Type()))
		}
		return true, nil
	case reflect.Int16:
		if val, err := strconv.ParseInt(defaultVal, 0, 16); err == nil {
			field.Set(reflect.ValueOf(int16(val)).Convert(field.Type()))
		}
		return true, nil
	case reflect.Int32:
		if val, err := strconv.ParseInt(defaultVal, 0, 32); err == nil {
			field.Set(reflect.ValueOf(int32(val)).Convert(field.Type()))
		}
		return true, nil
	case reflect.Int64:
		if val, err := time.ParseDuration(defaultVal); err == nil {
			field.Set(reflect.ValueOf(val).Convert(field.Type()))
		} else if val, err := strconv.ParseInt(defaultVal, 0, 64); err == nil {
			field.Set(reflect.ValueOf(val).Convert(field.Type()))
		}
		return true, nil
	}
	return false, nil
}

// setUint 设置uint
func setUint(field reflect.Value, defaultVal string) (bool, error) {
	switch field.Kind() {
	case reflect.Uint:
		if val, err := strconv.ParseUint(defaultVal, 0, strconv.IntSize); err == nil {
			field.Set(reflect.ValueOf(uint(val)).Convert(field.Type()))
		}
		return true, nil
	case reflect.Uint8:
		if val, err := strconv.ParseUint(defaultVal, 0, 8); err == nil {
			field.Set(reflect.ValueOf(uint8(val)).Convert(field.Type()))
		}
		return true, nil
	case reflect.Uint16:
		if val, err := strconv.ParseUint(defaultVal, 0, 16); err == nil {
			field.Set(reflect.ValueOf(uint16(val)).Convert(field.Type()))
		}
		return true, nil
	case reflect.Uint32:
		if val, err := strconv.ParseUint(defaultVal, 0, 32); err == nil {
			field.Set(reflect.ValueOf(uint32(val)).Convert(field.Type()))
		}
		return true, nil
	case reflect.Uint64:
		if val, err := strconv.ParseUint(defaultVal, 0, 64); err == nil {
			field.Set(reflect.ValueOf(val).Convert(field.Type()))
		}
		return true, nil
	case reflect.Uintptr:
		if val, err := strconv.ParseUint(defaultVal, 0, strconv.IntSize); err == nil {
			field.Set(reflect.ValueOf(uintptr(val)).Convert(field.Type()))
		}
		return true, nil
	}
	return false, nil
}

// setFloat 设置 float
func setFloat(field reflect.Value, defaultVal string) (bool, error) {
	switch field.Kind() {
	case reflect.Float32:
		if val, err := strconv.ParseFloat(defaultVal, 32); err == nil {
			field.Set(reflect.ValueOf(float32(val)).Convert(field.Type()))
		}
		return true, nil
	case reflect.Float64:
		if val, err := strconv.ParseFloat(defaultVal, 64); err == nil {
			field.Set(reflect.ValueOf(val).Convert(field.Type()))
		}
		return true, nil
	}
	return false, nil
}

// isInitialValue 判断是否值为空,为空则设置默认值
func isInitialValue(field reflect.Value) bool {
	return reflect.DeepEqual(reflect.Zero(field.Type()).Interface(), field.Interface())
}

// shouldInitializeField 验证字段
func shouldInitializeField(field reflect.Value, tag string) bool {
	switch field.Kind() {
	case reflect.Struct:
		return true
	case reflect.Ptr:
		if !field.IsNil() && field.Elem().Kind() == reflect.Struct {
			return true
		}
	case reflect.Slice:
		return field.Len() > 0 || tag != ""
	}
	return tag != ""
}
