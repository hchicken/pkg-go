package util

import (
	"io/ioutil"
	"os"

	"github.com/hchicken/pkg-go/util/defaults"
	"gopkg.in/yaml.v2"
)

// FileIsExist 判断文件是否存在
func FileIsExist(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}
	return false
}

// ReadYaml 读取yaml文件
func ReadYaml(info interface{}, file string) error {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, info)
	if err != nil {
		return err
	}
	err = defaults.Set(info)
	return err
}
