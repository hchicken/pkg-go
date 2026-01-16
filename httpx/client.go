package httpx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

// Client HTTP客户端
type Client struct {
	opts Options
}

// NewHttpClient 创建新的HTTP客户端
func NewHttpClient(opts ...Option) *Client {
	return newHttpClient(opts...)
}

func newHttpClient(opts ...Option) *Client {
	cli := new(Client)
	options := newOptions(opts...)
	cli.opts = options
	return cli
}

func (c *Client) parse() error {
	if c.opts.Result == nil {
		return nil
	}

	body := c.opts.Response.Body()

	// 检查响应体大小
	if c.opts.MaxResponseSize > 0 && int64(len(body)) > c.opts.MaxResponseSize {
		return fmt.Errorf("response body size %d exceeds limit %d", len(body), c.opts.MaxResponseSize)
	}

	// 检查响应体是否为空
	if len(body) == 0 {
		return nil
	}

	// 检查 Content-Type 是否是 JSON
	contentType := c.opts.Response.Header().Get("Content-Type")
	if !c.isJSONContentType(contentType) {
		c.opts.Logger.Printf("expected JSON response but got Content-Type: %s, body: %s", contentType, string(body))
		return fmt.Errorf("expected JSON response but got Content-Type: %s", contentType)
	}

	if err := json.Unmarshal(body, c.opts.Result); err != nil {
		c.opts.Logger.Printf("json unmarshal failed: %v, body: %s", err, string(body))
		return fmt.Errorf("json unmarshal failed: %w", err)
	}
	return nil
}

// isJSONContentType 检查是否是 JSON Content-Type
func (c *Client) isJSONContentType(contentType string) bool {
	if contentType == "" {
		// 没有 Content-Type 时尝试解析
		return true
	}
	ct := strings.ToLower(contentType)
	return strings.Contains(ct, "application/json") ||
		strings.Contains(ct, "text/json") ||
		strings.Contains(ct, "+json")
}

// isSuccessStatus 检查状态码是否成功
func (c *Client) isSuccessStatus(code int) bool {
	return code >= c.opts.SuccessStatusMin && code <= c.opts.SuccessStatusMax
}

func (c *Client) do(method string) error {
	// 执行请求拦截器
	for _, interceptor := range c.opts.RequestInterceptors {
		if err := interceptor(c.opts.Request); err != nil {
			return errors.Wrap(err, "request interceptor failed")
		}
	}

	var rsp *resty.Response
	var err error

	switch strings.ToUpper(method) {
	case "GET":
		rsp, err = c.opts.Request.Get(c.opts.URL)
	case "POST":
		rsp, err = c.opts.Request.Post(c.opts.URL)
	case "PUT":
		rsp, err = c.opts.Request.Put(c.opts.URL)
	case "PATCH":
		rsp, err = c.opts.Request.Patch(c.opts.URL)
	case "DELETE":
		rsp, err = c.opts.Request.Delete(c.opts.URL)
	case "HEAD":
		rsp, err = c.opts.Request.Head(c.opts.URL)
	case "OPTIONS":
		rsp, err = c.opts.Request.Options(c.opts.URL)
	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}

	// 先设置响应，即使有错误
	if rsp != nil {
		c.opts.Response = rsp
	}

	if err != nil {
		return err
	}

	// 执行响应拦截器
	for _, interceptor := range c.opts.ResponseInterceptors {
		if err := interceptor(rsp); err != nil {
			return errors.Wrap(err, "response interceptor failed")
		}
	}

	// 检查状态码
	if !c.opts.IgnoreStatus && !c.isSuccessStatus(rsp.StatusCode()) {
		rspBody := string(rsp.Body())
		c.opts.Logger.Printf("request failed: url=%s, status=%d, body=%s", c.opts.URL, rsp.StatusCode(), rspBody)
		return fmt.Errorf("request failed: url=%s, status=%d, body=%s", c.opts.URL, rsp.StatusCode(), rspBody)
	}

	return nil
}

// calculateBackoff 计算退避时间
func (c *Client) calculateBackoff(retryCount int) time.Duration {
	waitTime := c.opts.RetryWaitTime * time.Duration(1<<uint(retryCount))
	if waitTime > c.opts.MaxRetryWaitTime {
		waitTime = c.opts.MaxRetryWaitTime
	}
	return waitTime
}

// Do 发送请求
func (c *Client) Do(method string) error {
	return c.DoWithContext(c.opts.Ctx, method)
}

// DoWithContext 带上下文的请求
func (c *Client) DoWithContext(ctx context.Context, method string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	c.opts.Ctx = ctx
	c.opts.Request.SetContext(ctx)

	var lastErr error
	for attempt := 0; attempt <= c.opts.MaxRetryCount; attempt++ {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		lastErr = c.do(method)
		if lastErr == nil {
			return c.parse()
		}

		// 如果不需要重试或已达到最大重试次数，直接返回错误
		if !c.opts.Retry || attempt >= c.opts.MaxRetryCount {
			break
		}

		// 检查重试条件
		if c.opts.RetryCondition != nil && !c.opts.RetryCondition(c.opts.Response, lastErr) {
			break
		}

		// 计算退避时间并等待
		backoff := c.calculateBackoff(attempt)
		c.opts.Logger.Printf("request failed, retrying in %v (attempt %d/%d): %v",
			backoff, attempt+1, c.opts.MaxRetryCount, lastErr)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}

		// 重置 Request 以便重试
		c.opts.resetRequest()
	}

	return lastErr
}

// Get 发送GET请求
func (c *Client) Get() error {
	return c.Do("GET")
}

// Post 发送POST请求
func (c *Client) Post() error {
	return c.Do("POST")
}

// Put 发送PUT请求
func (c *Client) Put() error {
	return c.Do("PUT")
}

// Patch 发送PATCH请求
func (c *Client) Patch() error {
	return c.Do("PATCH")
}

// Delete 发送DELETE请求
func (c *Client) Delete() error {
	return c.Do("DELETE")
}

// GetOptions 返回Options配置
func (c *Client) GetOptions() *Options {
	return &c.opts
}

// SetHeaders 设置headers（支持链式调用）
func (c *Client) SetHeaders(headers map[string]string) *Client {
	c.opts.Request.SetHeaders(headers)
	c.opts.options = append(c.opts.options, Headers(headers))
	return c
}

// SetParam 设置param（支持链式调用）
func (c *Client) SetParam(param map[string]string) *Client {
	c.opts.Request.SetQueryParams(param)
	c.opts.options = append(c.opts.options, Param(param))
	return c
}

// SetBody 设置body（支持链式调用）
func (c *Client) SetBody(body interface{}) *Client {
	c.opts.Request.SetBody(body)
	c.opts.options = append(c.opts.options, Body(body))
	return c
}

// SetResponse 设置返回响应
func (c *Client) SetResponse(rsp *resty.Response) *Client {
	c.opts.Response = rsp
	return c
}

// SetURL 设置URL（支持链式调用）
func (c *Client) SetURL(url string) *Client {
	c.opts.URL = url
	return c
}

// SetContext 设置上下文（支持链式调用）
func (c *Client) SetContext(ctx context.Context) *Client {
	c.opts.Ctx = ctx
	c.opts.Request.SetContext(ctx)
	return c
}

// GetResponse 获取响应
func (c *Client) GetResponse() *resty.Response {
	return c.opts.Response
}

// GetStatusCode 获取响应状态码
func (c *Client) GetStatusCode() int {
	if c.opts.Response == nil {
		return 0
	}
	return c.opts.Response.StatusCode()
}

// GetBody 获取响应体
func (c *Client) GetBody() []byte {
	if c.opts.Response == nil {
		return nil
	}
	return c.opts.Response.Body()
}

// GetBodyString 获取响应体字符串
func (c *Client) GetBodyString() string {
	if c.opts.Response == nil {
		return ""
	}
	return c.opts.Response.String()
}

// SetFormData 设置表单数据（支持链式调用）
func (c *Client) SetFormData(data map[string]string) *Client {
	c.opts.Request.SetFormData(data)
	return c
}

// SetFile 设置文件上传（支持链式调用）
func (c *Client) SetFile(param, filePath string) *Client {
	c.opts.Request.SetFile(param, filePath)
	return c
}

// SetFileReader 通过Reader设置文件上传（支持链式调用）
func (c *Client) SetFileReader(param, fileName string, reader io.Reader) *Client {
	c.opts.Request.SetFileReader(param, fileName, reader)
	return c
}

// SetBasicAuth 设置Basic认证（支持链式调用）
func (c *Client) SetBasicAuth(username, password string) *Client {
	c.opts.Request.SetBasicAuth(username, password)
	return c
}

// SetAuthToken 设置Bearer Token（支持链式调用）
func (c *Client) SetAuthToken(token string) *Client {
	c.opts.Request.SetAuthToken(token)
	return c
}

// SetHeader 设置单个header（支持链式调用）
func (c *Client) SetHeader(key, value string) *Client {
	c.opts.Request.SetHeader(key, value)
	c.opts.options = append(c.opts.options, func(o *Options) {
		if o.Request != nil {
			o.Request.SetHeader(key, value)
		}
	})
	return c
}

// SetQueryParam 设置单个查询参数（支持链式调用）
func (c *Client) SetQueryParam(key, value string) *Client {
	c.opts.Request.SetQueryParam(key, value)
	c.opts.options = append(c.opts.options, func(o *Options) {
		if o.Request != nil {
			o.Request.SetQueryParam(key, value)
		}
	})
	return c
}

// GetRestyClient 获取底层resty客户端（高级用法）
func (c *Client) GetRestyClient() *resty.Client {
	return c.opts.RestyClient
}

// GetRestyRequest 获取底层resty请求（高级用法）
func (c *Client) GetRestyRequest() *resty.Request {
	return c.opts.Request
}
