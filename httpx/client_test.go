package httpx

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewHttpClient(t *testing.T) {
	client := NewHttpClient()
	assert.NotNil(t, client)
}

func TestSetHeaders(t *testing.T) {
	client := NewHttpClient()
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	client.SetHeaders(headers)
	assert.Equal(t, headers, client.Options().Headers)
}

func TestSetParam(t *testing.T) {
	client := NewHttpClient()
	params := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	client.SetParam(params)
	assert.Equal(t, params, client.Options().Param)
}

func TestSetBody(t *testing.T) {
	client := NewHttpClient()
	body := map[string]interface{}{
		"key": "value",
	}
	client.SetBody(body)
	assert.Equal(t, body, client.Options().Body)
}

func TestSetResponse(t *testing.T) {
	client := NewHttpClient()
	response := &resty.Response{}
	client.SetResponse(response)
	assert.Equal(t, response, client.Options().Response)
}

func TestDo(t *testing.T) {
	// 创建一个模拟的HTTP服务器
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// 在这里，我们可以根据需要设置响应状态码、响应头和响应体
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`OK`))
	}))
	// 关闭服务器以释放资源
	defer server.Close()
	fmt.Println(server.URL)
	// 创建一个新的Client对象，并设置URL为模拟服务器的URL
	client := NewHttpClient(
		URL(server.URL),
	)

	// 测试GET请求
	err := client.Do("post")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, client.opts.Response.StatusCode())
	assert.Equal(t, "OK", string(client.opts.Response.Body()))

	// 测试不支持的请求方法
	err = client.Do("unsupported")
	assert.NotNil(t, err)
}
