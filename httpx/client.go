package httpx

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

// Client 客户端
type Client struct {
	opts Options
	err  error
}

// NewHttpClient http client
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
	// 没设置返回结果直接return
	if c.opts.Result == nil {
		return nil
	}
	err := json.Unmarshal(c.opts.Response.Body(), c.opts.Result)
	if err != nil {
		return errors.Wrap(err, "parsing result")
	}
	return nil
}

func (c *Client) do(method string) error {
	var rsp *resty.Response
	var err error

	// 设置请求数据
	switch strings.ToLower(method) {
	case "get":
		rsp, err = c.opts.Request.Get(c.opts.URL)
	case "post":
		rsp, err = c.opts.Request.Post(c.opts.URL)
	case "put":
		rsp, err = c.opts.Request.Put(c.opts.URL)
	case "patch":
		rsp, err = c.opts.Request.Patch(c.opts.URL)
	case "delete":
		rsp, err = c.opts.Request.Delete(c.opts.URL)
	default:
		err = fmt.Errorf("does not support the request of the method [%v]", method)
		return err
	}
	defer c.SetResponse(rsp)
	//fmt.Println(rsp.StatusCode())
	// 设置body
	if !c.opts.IgnoreStatus && rsp.StatusCode() != 200 {
		rspBody := string(rsp.Body())
		log.Println(rspBody)
		return fmt.Errorf("invoke url [%v], code is [%v],  response body is [%v]", c.opts.URL, c.opts.Response.StatusCode(), rspBody)
	}
	return err
}

// Do 发送请求
func (c *Client) Do(method string) error {
	err := c.do(method)
	if err != nil && c.opts.Retry && c.opts.RetryCount < c.opts.MaxRetryCount {
		c.opts.RetryCount++
		return c.Do(method)
	}
	if err != nil {
		return err
	}
	return c.parse()
}

// Options 返回Options
func (c *Client) Options() *Options {
	return &c.opts
}

// SetHeaders 设置headers
func (c *Client) SetHeaders(headers map[string]string) {
	c.opts.Headers = headers
	c.opts.Request.SetHeaders(headers)
}

// SetParam 设置param
func (c *Client) SetParam(param map[string]string) {
	c.opts.Param = param
	c.opts.Request.SetQueryParams(param)
}

// SetBody 设置body
func (c *Client) SetBody(body interface{}) {
	c.opts.Body = body
	c.opts.Request.SetBody(body)
}

// SetResponse 设置返回响应body
func (c *Client) SetResponse(rsp *resty.Response) {
	c.opts.Response = rsp
}
