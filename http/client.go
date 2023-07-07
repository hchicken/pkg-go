package http

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

// Client 客户端
type Client struct {
	opts Options
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

func (c *Client) parse(rsp *resty.Response) error {
	// 设置body
	c.opts.Response = rsp
	if c.opts.Response.StatusCode() != 200 {
		fmt.Println(string(rsp.Body()))
		return fmt.Errorf(
			"the state of the request url [%v] code is %v",
			c.opts.URL, c.opts.Response.StatusCode(),
		)
	}
	// 没设置返回结果直接return
	if c.opts.Result == nil {
		return nil
	}
	err := json.Unmarshal(rsp.Body(), c.opts.Result)
	if err != nil {
		return errors.Wrap(err, "error parsing results")
	}
	return nil
}

// SetResponse 设置返回响应body
func (c *Client) SetResponse(rsp *resty.Response) {
	c.opts.Response = rsp
}

// Get 请求
func (c *Client) Get() error {
	rsp, err := c.opts.Request.Get(c.opts.URL) // 发送请求
	if err != nil {
		return err
	}
	return c.parse(rsp)
}

// Post 请求
func (c *Client) Post() error {
	// 发送请求
	rsp, err := c.opts.Request.Post(c.opts.URL) // 发送请求
	if err != nil {
		return err
	}
	return c.parse(rsp)
}

// Put 请求
func (c *Client) Put() error {
	// 发送请求
	rsp, err := c.opts.Request.Put(c.opts.URL) // 发送请求
	if err != nil {
		return err
	}
	return c.parse(rsp)
}

// Patch 请求
func (c *Client) Patch() error {
	// 发送请求
	rsp, err := c.opts.Request.Patch(c.opts.URL) // 发送请求
	if err != nil {
		return err
	}
	return c.parse(rsp)
}

// Delete 请求
func (c *Client) Delete() error {
	// 发送请求
	rsp, err := c.opts.Request.Delete(c.opts.URL) // 发送请求
	if err != nil {
		return err
	}
	return c.parse(rsp)
}

// Do 发送请求
func (c *Client) Do(method string) error {
	var err error
	switch strings.ToLower(method) {
	case "get":
		err = c.Get()
	case "post":
		err = c.Post()
	case "put":
		err = c.Put()
	case "patch":
		err = c.Patch()
	case "delete":
		err = c.Delete()
	default:
		err = fmt.Errorf("does not support the request of the method [%v]", method)
	}
	return err
}

// Options 返回Options
func (c *Client) Options() Options {
	return c.opts
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
