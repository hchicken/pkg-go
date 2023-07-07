package util

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// StrToMd5 字符串md5
func StrToMd5(str string) string {
	valueMd5 := md5.Sum([]byte(str))
	value := hex.EncodeToString(valueMd5[:])
	return value
}

// UUID 生成uuid字符串
func UUID() string {
	uuid4, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return uuid4.String()
}

// GzipEn gzip加密
func GzipEn(d interface{}) (string, error) {
	buf, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(buf); err != nil {
		panic(err)
	}
	if err := gz.Flush(); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}
	encoded := base64.StdEncoding.EncodeToString(b.Bytes())
	return encoded, nil
}

// GzipDe gzip解密
func GzipDe(in string) ([]byte, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return nil, err
	}
	reader, err := gzip.NewReader(bytes.NewReader(decodeBytes))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = reader.Close()
	}()
	return ioutil.ReadAll(reader)
}

// GzipDeValue gzip解密并赋值
func GzipDeValue(in string, v interface{}) error {
	buf, err := GzipDe(in)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, v)
	return err
}

// RandString 生成随机字符串
func RandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

// MakeUrlStr 生成url字符串
func MakeUrlStr(params map[string]string, keys ...string) string {
	keyArr := make([]string, 0)
	if len(keys) == 0 {
		// 安装字母顺序排序
		for k := range params {
			keyArr = append(keyArr, k)
		}
		sort.Strings(keyArr)
	} else {
		keyArr = keys
	}
	resultArr := make([]string, 0)
	for _, v := range keyArr {
		resultArr = append(resultArr, fmt.Sprintf("%v=%v", v, params[v]))
	}
	return strings.Join(resultArr, "&")
}
