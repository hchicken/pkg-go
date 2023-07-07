package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CODE_SUCCESS = iota // 正确状态码
	CODE_ERROR
)

const (
	TraceId = "private-trace-id"
)

type Option func(*Options)

// Options http请求options
type Options struct {
	message string
	status  int
	code    int
	data    any
	output  any
	async   bool
}

// Message ...
func Message(m string) Option {
	return func(o *Options) {
		o.message = m
	}
}

// Error ...
func Error(err error) Option {
	return func(o *Options) {
		o.message = fmt.Sprintf("%v", err)
		o.status = CODE_ERROR
		o.code = http.StatusBadRequest
	}
}

// Data ...
func Data(data any) Option {
	return func(o *Options) {
		o.data = data
	}
}

// Async ...
func Async(ok bool) Option {
	return func(o *Options) {
		o.async = ok
	}
}

// Json ...
func Json(c *gin.Context, opts ...Option) {
	options := Options{
		message: "success",
		code:    http.StatusOK,
		status:  CODE_SUCCESS,
	}
	for _, o := range opts {
		o(&options)
	}
	c.AbortWithStatusJSON(options.code, gin.H{
		"trace_id": c.GetHeader(TraceId),
		"code":     options.status,
		"message":  options.message,
		"data":     options.data,
		"async":    options.async,
	})
}