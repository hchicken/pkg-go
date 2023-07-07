package structx

import (
	"encoding/json"

	jsonIterator "github.com/json-iterator/go"
)

// StructDecode 结构体转换 input ==> output(指针类型)
func StructDecode(input interface{}, output interface{}) error {
	// struct to byte
	dByte, err := json.Marshal(input)
	if err != nil {
		return err
	}
	// byte to map
	if err := json.Unmarshal(dByte, &output); err != nil {
		return err
	}
	return nil
}

// StructSpecialDecode 结构体转换 input ==> output(指针类型)
func StructSpecialDecode(input interface{}, output interface{}) error {
	// struct to byte
	var json1 = jsonIterator.ConfigCompatibleWithStandardLibrary
	dByte, err := json1.Marshal(input)
	if err != nil {
		return err
	}
	// byte to map
	if err := json1.Unmarshal(dByte, &output); err != nil {
		return err
	}
	return nil
}
