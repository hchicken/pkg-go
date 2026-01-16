package httpx

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

// DefaultTimeout 默认超时时间
const DefaultTimeout = 30 * time.Second

// DefaultRetryWaitTime 默认重试等待时间
const DefaultRetryWaitTime = 100 * time.Millisecond

// DefaultMaxRetryWaitTime 默认最大重试等待时间
const DefaultMaxRetryWaitTime = 2 * time.Second

// Logger 日志接口，兼容 logrus.Logger、logx.LoggerIns、log.Logger 等
type Logger interface {
	Printf(format string, v ...interface{})
}

// StdLogger 标准库日志适配器
// 示例：httpx.NewStdLogger(logxIns.Warnf)
type StdLogger struct {
	LogFunc func(format string, v ...interface{})
}

// Printf 实现 Logger 接口
func (l *StdLogger) Printf(format string, v ...interface{}) {
	if l.LogFunc != nil {
		l.LogFunc(format, v...)
	}
}

// NewStdLogger 创建标准日志适配器
func NewStdLogger(logFunc func(format string, v ...interface{})) *StdLogger {
	return &StdLogger{LogFunc: logFunc}
}

// defaultLogger 默认日志实现（空实现）
type defaultLogger struct{}

func (l *defaultLogger) Printf(format string, v ...interface{}) {}

// DefaultClient 默认客户端
var DefaultClient = func(opts *Options) *resty.Client {
	cli := resty.New().SetCloseConnection(true)

	// TLS 配置
	tlsConfig := &tls.Config{InsecureSkipVerify: opts.InsecureSkipVerify}
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	// 设置超时
	if opts.Timeout > 0 {
		cli.SetTimeout(opts.Timeout)
	} else {
		cli.SetTimeout(DefaultTimeout)
	}

	cli.SetTransport(transport)
	return cli
}

// Option 选项函数类型
type Option func(*Options)

// Options http请求options
type Options struct {
	RestyClient *resty.Client
	URL         string
	Result      interface{}
	Request     *resty.Request
	Response    *resty.Response

	// 上下文
	Ctx context.Context

	// 超时配置
	Timeout time.Duration

	// TLS 配置
	InsecureSkipVerify bool

	// 重试配置
	Retry            bool
	MaxRetryCount    int
	RetryWaitTime    time.Duration
	MaxRetryWaitTime time.Duration
	RetryCondition   RetryCondition

	// 状态码配置
	IgnoreStatus     bool
	SuccessStatusMin int
	SuccessStatusMax int

	// 响应限制
	MaxResponseSize int64

	// 日志
	Logger Logger

	// 拦截器
	RequestInterceptors  []RequestInterceptor
	ResponseInterceptors []ResponseInterceptor

	// 内部使用：保存 Option 用于重试时重建 Request
	options []Option
}

// RequestInterceptor 请求拦截器
type RequestInterceptor func(*resty.Request) error

// ResponseInterceptor 响应拦截器
type ResponseInterceptor func(*resty.Response) error

// RetryCondition 重试条件函数类型
type RetryCondition func(resp *resty.Response, err error) bool

// defaultOptions 返回默认配置
func defaultOptions() Options {
	return Options{
		MaxRetryCount:    3,
		Timeout:          DefaultTimeout,
		RetryWaitTime:    DefaultRetryWaitTime,
		MaxRetryWaitTime: DefaultMaxRetryWaitTime,
		SuccessStatusMin: 200,
		SuccessStatusMax: 299,
		Logger:           &defaultLogger{},
		Ctx:              context.Background(),
	}
}

func newOptions(opts ...Option) Options {
	o := defaultOptions()
	o.options = opts

	// 应用所有选项
	o.apply(opts...)

	// 初始化 RestyClient
	if o.RestyClient == nil {
		o.RestyClient = DefaultClient(&o)
	}

	// 创建并初始化 Request
	o.initRequest()

	return o
}

// apply 应用选项
func (o *Options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// initRequest 初始化 Request
func (o *Options) initRequest() {
	o.Request = o.RestyClient.R()
	o.apply(o.options...)
	if o.Ctx != nil {
		o.Request.SetContext(o.Ctx)
	}
}

// resetRequest 重置 Request 用于重试
func (o *Options) resetRequest() {
	o.initRequest()
}

// URL 设置请求URL
func URL(uri string) Option {
	return func(o *Options) {
		o.URL = uri
	}
}

// WithContext 设置上下文
func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Ctx = ctx
	}
}

// Timeout 设置超时时间
func Timeout(d time.Duration) Option {
	return func(o *Options) {
		o.Timeout = d
	}
}

// InsecureSkipVerify 设置是否跳过TLS证书验证
func InsecureSkipVerify(skip bool) Option {
	return func(o *Options) {
		o.InsecureSkipVerify = skip
	}
}

// QueryString 请求参数(string形式)
func QueryString(s string) Option {
	return func(o *Options) {
		if o.Request != nil {
			o.Request.SetQueryString(s)
		}
	}
}

// Headers 设置 headers
func Headers(headers map[string]string) Option {
	return func(o *Options) {
		if o.Request != nil {
			o.Request.SetHeaders(headers)
		}
	}
}

// Param 请求参数(map形式)
func Param(param map[string]string) Option {
	return func(o *Options) {
		if o.Request != nil {
			o.Request.SetQueryParams(param)
		}
	}
}

// Body 请求body
func Body(body interface{}) Option {
	return func(o *Options) {
		if o.Request != nil {
			o.Request.SetBody(body)
		}
	}
}

// Result 请求结果
func Result(result interface{}) Option {
	return func(o *Options) {
		o.Result = result
	}
}

// Retry 失败重试
func Retry(b bool) Option {
	return func(o *Options) {
		o.Retry = b
	}
}

// MaxRetryCount 重试的最大次数
func MaxRetryCount(count int) Option {
	return func(o *Options) {
		o.MaxRetryCount = count
	}
}

// RetryWaitTime 重试等待时间
func RetryWaitTime(d time.Duration) Option {
	return func(o *Options) {
		o.RetryWaitTime = d
	}
}

// MaxRetryWaitTime 最大重试等待时间
func MaxRetryWaitTime(d time.Duration) Option {
	return func(o *Options) {
		o.MaxRetryWaitTime = d
	}
}

// WithRetryCondition 设置重试条件
func WithRetryCondition(condition RetryCondition) Option {
	return func(o *Options) {
		o.RetryCondition = condition
	}
}

// RetryOn5xx 只在5xx错误时重试
func RetryOn5xx() Option {
	return func(o *Options) {
		o.RetryCondition = func(resp *resty.Response, err error) bool {
			if resp == nil {
				return true
			}
			return resp.StatusCode() >= 500
		}
	}
}

// RestyClient 自定义客户端
func RestyClient(cli *resty.Client) Option {
	return func(o *Options) {
		o.RestyClient = cli
	}
}

// IgnoreStatus 忽略状态码检查
func IgnoreStatus(b bool) Option {
	return func(o *Options) {
		o.IgnoreStatus = b
	}
}

// SuccessStatusRange 设置成功状态码范围
func SuccessStatusRange(min, max int) Option {
	return func(o *Options) {
		o.SuccessStatusMin = min
		o.SuccessStatusMax = max
	}
}

// MaxResponseSize 设置最大响应体大小
func MaxResponseSize(size int64) Option {
	return func(o *Options) {
		o.MaxResponseSize = size
	}
}

// WithLogger 设置日志接口
func WithLogger(logger Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

// WithRequestInterceptor 添加请求拦截器
func WithRequestInterceptor(interceptor RequestInterceptor) Option {
	return func(o *Options) {
		o.RequestInterceptors = append(o.RequestInterceptors, interceptor)
	}
}

// WithResponseInterceptor 添加响应拦截器
func WithResponseInterceptor(interceptor ResponseInterceptor) Option {
	return func(o *Options) {
		o.ResponseInterceptors = append(o.ResponseInterceptors, interceptor)
	}
}
