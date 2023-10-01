package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义状态码常量
const (
	CodeSuccess = iota // 正确状态码
	CodeError          // 错误状态码
)

// 定义TraceId常量
const TraceId = "private-trace-id"

// Option 是一个函数类型，用于修改Options
type Option func(*Options)

// Options 结构体用于http请求的选项
type Options struct {
	message     string
	status      int
	code        int
	data        interface{}
	output      interface{}
	async       bool
	customField map[string]interface{}
}

// NewOptions 创建一个新的Options实例
func NewOptions() *Options {
	return &Options{
		message:     "success",
		code:        http.StatusOK,
		status:      CodeSuccess,
		customField: make(map[string]interface{}),
	}
}

// ApplyOptions 应用选项到Options实例
func (o *Options) ApplyOptions(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// Message 设置消息选项
func Message(m string) Option {
	return func(o *Options) {
		o.message = m
	}
}

// Error 设置错误选项
func Error(err error) Option {
	return func(o *Options) {
		o.message = fmt.Sprintf("%v", err)
		o.status = CodeError
		o.code = http.StatusBadRequest
	}
}

// Data 设置数据选项
func Data(data interface{}) Option {
	return func(o *Options) {
		o.data = data
	}
}

// Async 设置异步选项
func Async(ok bool) Option {
	return func(o *Options) {
		o.async = ok
	}
}

// CustomField 设置自定义字段选项
func CustomField(key string, value interface{}) Option {
	return func(o *Options) {
		o.customField[key] = value
	}
}

// Json 生成并返回JSON响应
func Json(c *gin.Context, opts ...Option) {
	options := NewOptions()
	options.ApplyOptions(opts...)

	retData := gin.H{
		"trace_id": c.GetHeader(TraceId),
		"code":     options.status,
		"message":  options.message,
		"data":     options.data,
		"async":    options.async,
	}

	// 添加自定义返回字段
	for k, v := range options.customField {
		retData[k] = v
	}

	c.AbortWithStatusJSON(options.code, retData)
}
