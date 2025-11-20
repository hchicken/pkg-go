package response

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/gin-gonic/gin"
)

// 响应状态码常量
const (
	CodeSuccess = 0 // 成功状态码
	CodeError   = 1 // 错误状态码
)

// HTTP头常量
const (
	TraceIDHeader   = "private-trace-id"
	DefaultFileName = "download"
)

// 默认消息
const DefaultMessage = "success"

// Option 函数选项模式，用于配置响应参数
type Option func(*Options)

// Options 响应配置结构体
type Options struct {
	message     string                 // 响应消息
	code        int                    // HTTP状态码
	status      interface{}            // 业务状态码，支持int和string
	data        interface{}            // 响应数据
	async       bool                   // 是否异步处理
	customField map[string]interface{} // 自定义字段
	fileBytes   []byte                 // 文件字节数据
	filename    string                 // 文件名
}

// newDefaultOptions 创建默认配置
func newDefaultOptions() *Options {
	return &Options{
		message:     DefaultMessage,
		code:        http.StatusOK,
		status:      CodeSuccess,
		async:       false,
		customField: make(map[string]interface{}),
		filename:    DefaultFileName,
	}
}

// 配置选项函数

// Message 设置响应消息
func Message(msg string) Option {
	return func(o *Options) {
		o.message = msg
	}
}

// Error 设置错误响应
func Error(err error) Option {
	return func(o *Options) {
		if err != nil {
			o.message = err.Error()
			o.status = CodeError
			o.code = http.StatusBadRequest
		}
	}
}

// Code 设置HTTP状态码
func Code(code int) Option {
	return func(o *Options) {
		o.code = code
	}
}

// StatusInt 设置业务状态码为int类型
func StatusInt(status int) Option {
	return func(o *Options) {
		o.status = status
	}
}

// StatusString 设置业务状态码为string类型
func StatusString(status string) Option {
	return func(o *Options) {
		o.status = status
	}
}

// Data 设置响应数据
func Data(data interface{}) Option {
	return func(o *Options) {
		o.data = data
	}
}

// Async 设置异步标识
func Async(async bool) Option {
	return func(o *Options) {
		o.async = async
	}
}

// CustomField 添加自定义字段
func CustomField(key string, value interface{}) Option {
	return func(o *Options) {
		o.customField[key] = value
	}
}

// FileBytes 设置文件字节数据
func FileBytes(data []byte) Option {
	return func(o *Options) {
		o.fileBytes = data
	}
}

// FileName 设置文件名
func FileName(name string) Option {
	return func(o *Options) {
		o.filename = name
	}
}

// JSON 返回JSON格式响应
func JSON(c *gin.Context, opts ...Option) {
	options := newDefaultOptions()

	// 应用所有选项
	for _, opt := range opts {
		opt(options)
	}

	// 构建响应数据
	response := gin.H{
		"trace_id": c.GetHeader(TraceIDHeader),
		"code":     options.status,
		"message":  options.message,
		"data":     options.data,
		"async":    options.async,
	}

	// 添加自定义字段
	for key, value := range options.customField {
		response[key] = value
	}

	c.AbortWithStatusJSON(options.code, response)
}

// File 返回文件下载响应
func File(c *gin.Context, opts ...Option) {
	options := newDefaultOptions()

	// 应用所有选项
	for _, opt := range opts {
		opt(options)
	}

	if len(options.fileBytes) == 0 {
		JSON(c, Error(fmt.Errorf("文件数据为空")))
		return
	}

	// 设置文件下载头
	filename := options.filename
	if filename == "" {
		filename = DefaultFileName
	}

	disposition := fmt.Sprintf(`attachment; filename=%s`, url.PathEscape(path.Base(filename)))
	c.Header("Content-Disposition", disposition)
	c.Header("Content-Type", "application/octet-stream")
	c.Data(http.StatusOK, "application/octet-stream", options.fileBytes)
}

// Json 向后兼容的JSON响应方法
func Json(c *gin.Context, opts ...Option) {
	JSON(c, opts...)
}

// Blob 向后兼容的文件响应方法
func Blob(c *gin.Context, opts ...Option) {
	File(c, opts...)
}

// 快捷响应方法

// Success 成功响应
func Success(c *gin.Context, data interface{}, message ...string) {
	msg := DefaultMessage
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	JSON(c, Data(data), Message(msg))
}

// Fail 失败响应
func Fail(c *gin.Context, message string, httpStatus ...int) {
	status := http.StatusBadRequest
	if len(httpStatus) > 0 {
		status = httpStatus[0]
	}
	JSON(c, Message(message), StatusInt(CodeError), Code(status))
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, err error, httpStatus ...int) {
	status := http.StatusBadRequest
	if len(httpStatus) > 0 {
		status = httpStatus[0]
	}
	JSON(c, Error(err), Code(status))
}
