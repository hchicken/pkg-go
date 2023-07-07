package http

import (
	"crypto/tls"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// DefaultClient 默认客户端
var DefaultClient = func() *resty.Client {
	cli := resty.New().SetCloseConnection(true)
	cli.SetTransport(&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}})
	return cli
}

// Option ...
type Option func(*Options)

// Options http请求options
type Options struct {
	RestyClient *resty.Client
	URL         string
	QueryString string
	Param       map[string]string
	Headers     map[string]string
	Body        interface{}
	Result      interface{}
	RspBody     []byte
	Request     *resty.Request
	Response    *resty.Response
}

func newOptions(opts ...Option) Options {
	opt := Options{
		RestyClient: DefaultClient(),
		Body:        new(interface{}),
	}

	opt.Request = opt.RestyClient.R()
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// URL url
func URL(uri string) Option {
	return func(o *Options) {
		o.URL = uri
	}
}

// QueryString 请求参数(string形式)
func QueryString(s string) Option {
	return func(o *Options) {
		o.QueryString = s
		o.Request.SetQueryString(s)
	}
}

// Headers 设置 headers
func Headers(headers map[string]string) Option {
	return func(o *Options) {
		o.Headers = headers
		o.Request.SetHeaders(headers)
	}
}

// Param 请求参数(map形式)
func Param(param map[string]string) Option {
	return func(o *Options) {
		o.Param = param
		o.Request.SetQueryParams(param)
	}
}

// Body 请求body
func Body(body interface{}) Option {
	return func(o *Options) {
		o.Body = body
		o.Request.SetBody(body)
	}
}

// Result 请求结果
func Result(result interface{}) Option {
	return func(o *Options) {
		o.Result = result
	}
}

// RestyClient 客户端
func RestyClient(cli *resty.Client) Option {
	return func(o *Options) {
		o.RestyClient = cli
	}
}
