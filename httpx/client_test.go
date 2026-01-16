package httpx

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewHttpClient(t *testing.T) {
	client := NewHttpClient()
	assert.NotNil(t, client)
}

func TestSetHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHttpClient(URL(server.URL))
	headers := map[string]string{"Content-Type": "application/json"}
	err := client.SetHeaders(headers).Get()
	assert.Nil(t, err)
}

func TestSetParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "value1", req.URL.Query().Get("param1"))
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHttpClient(URL(server.URL))
	params := map[string]string{"param1": "value1"}
	err := client.SetParam(params).Get()
	assert.Nil(t, err)
}

func TestSetBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHttpClient(URL(server.URL))
	body := map[string]interface{}{"key": "value"}
	err := client.SetBody(body).Post()
	assert.Nil(t, err)
}

func TestSetResponse(t *testing.T) {
	client := NewHttpClient()
	response := &resty.Response{}
	client.SetResponse(response)
	assert.Equal(t, response, client.GetOptions().Response)
}

func TestDo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()

	client := NewHttpClient(URL(server.URL))
	err := client.Do("post")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, client.GetStatusCode())
	assert.Equal(t, "OK", client.GetBodyString())

	err = client.Do("unsupported")
	assert.NotNil(t, err)
}

func TestHTTPMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, method, req.Method)
				rw.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			client := NewHttpClient(URL(server.URL))
			err := client.Do(method)
			assert.Nil(t, err)
		})
	}
}

func TestConvenienceMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	t.Run("Get", func(t *testing.T) {
		err := NewHttpClient(URL(server.URL)).Get()
		assert.Nil(t, err)
	})
	t.Run("Post", func(t *testing.T) {
		err := NewHttpClient(URL(server.URL)).Post()
		assert.Nil(t, err)
	})
	t.Run("Put", func(t *testing.T) {
		err := NewHttpClient(URL(server.URL)).Put()
		assert.Nil(t, err)
	})
	t.Run("Delete", func(t *testing.T) {
		err := NewHttpClient(URL(server.URL)).Delete()
		assert.Nil(t, err)
	})
}

func TestChainedCalls(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
		assert.Equal(t, "value1", req.URL.Query().Get("param1"))
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := NewHttpClient().
		SetURL(server.URL).
		SetHeader("Content-Type", "application/json").
		SetQueryParam("param1", "value1").
		Get()
	assert.Nil(t, err)
}

func TestTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		time.Sleep(2 * time.Second)
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := NewHttpClient(URL(server.URL), Timeout(100*time.Millisecond)).Get()
	assert.NotNil(t, err)
}

func TestContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		time.Sleep(2 * time.Second)
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := NewHttpClient(URL(server.URL), WithContext(ctx), Timeout(10*time.Second)).Get()
	assert.NotNil(t, err)
}

func TestRetryWithBackoff(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if atomic.AddInt32(&attempts, 1) < 3 {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := NewHttpClient(
		URL(server.URL),
		Retry(true),
		MaxRetryCount(3),
		RetryWaitTime(10*time.Millisecond),
		MaxRetryWaitTime(50*time.Millisecond),
	).Get()
	assert.Nil(t, err)
	assert.Equal(t, int32(3), atomic.LoadInt32(&attempts))
}

func TestRetryOn5xx(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if atomic.AddInt32(&attempts, 1) < 3 {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := NewHttpClient(
		URL(server.URL),
		Retry(true),
		MaxRetryCount(3),
		RetryOn5xx(),
		RetryWaitTime(10*time.Millisecond),
	).Get()
	assert.Nil(t, err)
	assert.Equal(t, int32(3), atomic.LoadInt32(&attempts))
}

func TestRetryConditionNoRetryOn4xx(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		atomic.AddInt32(&attempts, 1)
		rw.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	err := NewHttpClient(
		URL(server.URL),
		Retry(true),
		MaxRetryCount(3),
		RetryOn5xx(),
		RetryWaitTime(10*time.Millisecond),
	).Get()
	assert.NotNil(t, err)
	assert.Equal(t, int32(1), atomic.LoadInt32(&attempts))
}

func TestSuccessStatusRange(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	// 默认 200-299 应该成功
	err := NewHttpClient(URL(server.URL)).Get()
	assert.Nil(t, err)

	// 自定义范围只允许 200
	server2 := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusCreated)
	}))
	defer server2.Close()

	err = NewHttpClient(URL(server2.URL), SuccessStatusRange(200, 200)).Get()
	assert.NotNil(t, err)
}

func TestRequestInterceptor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "intercepted", req.Header.Get("X-Custom"))
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := NewHttpClient(
		URL(server.URL),
		WithRequestInterceptor(func(req *resty.Request) error {
			req.SetHeader("X-Custom", "intercepted")
			return nil
		}),
	).Get()
	assert.Nil(t, err)
}

func TestResponseInterceptor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	var intercepted bool
	err := NewHttpClient(
		URL(server.URL),
		WithResponseInterceptor(func(rsp *resty.Response) error {
			intercepted = true
			return nil
		}),
	).Get()
	assert.Nil(t, err)
	assert.True(t, intercepted)
}

func TestJSONResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"name":"test","value":123}`))
	}))
	defer server.Close()

	var result struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	err := NewHttpClient(URL(server.URL), Result(&result)).Get()
	assert.Nil(t, err)
	assert.Equal(t, "test", result.Name)
	assert.Equal(t, 123, result.Value)
}

func TestMaxResponseSize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"data":"large response"}`))
	}))
	defer server.Close()

	var result map[string]string
	err := NewHttpClient(URL(server.URL), MaxResponseSize(10), Result(&result)).Get()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "exceeds limit")
}

func TestInsecureSkipVerify(t *testing.T) {
	client := NewHttpClient()
	assert.False(t, client.GetOptions().InsecureSkipVerify)

	client2 := NewHttpClient(InsecureSkipVerify(true))
	assert.True(t, client2.GetOptions().InsecureSkipVerify)
}

func TestCustomLogger(t *testing.T) {
	var logged bool
	logger := &testLogger{onLog: func() { logged = true }}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_ = NewHttpClient(URL(server.URL), WithLogger(logger)).Get()
	assert.True(t, logged)
}

type testLogger struct {
	onLog func()
}

func (l *testLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
	if l.onLog != nil {
		l.onLog()
	}
}

func TestAuthMethods(t *testing.T) {
	t.Run("BasicAuth", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			user, pass, ok := req.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, "user", user)
			assert.Equal(t, "pass", pass)
			rw.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		err := NewHttpClient(URL(server.URL)).SetBasicAuth("user", "pass").Get()
		assert.Nil(t, err)
	})

	t.Run("BearerToken", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			assert.Equal(t, "Bearer test-token", req.Header.Get("Authorization"))
			rw.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		err := NewHttpClient(URL(server.URL)).SetAuthToken("test-token").Get()
		assert.Nil(t, err)
	})
}

func TestFormData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		assert.Equal(t, "value1", req.FormValue("key1"))
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := NewHttpClient(URL(server.URL)).SetFormData(map[string]string{"key1": "value1"}).Post()
	assert.Nil(t, err)
}

func TestIgnoreStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// 不忽略状态码
	err := NewHttpClient(URL(server.URL)).Get()
	assert.NotNil(t, err)

	// 忽略状态码
	client := NewHttpClient(URL(server.URL), IgnoreStatus(true))
	err = client.Get()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, client.GetStatusCode())
}

func TestGetBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`test body`))
	}))
	defer server.Close()

	client := NewHttpClient(URL(server.URL))
	client.Get()
	assert.Equal(t, []byte(`test body`), client.GetBody())
	assert.Equal(t, "test body", client.GetBodyString())
}

func TestNilResponse(t *testing.T) {
	client := NewHttpClient()
	assert.Equal(t, 0, client.GetStatusCode())
	assert.Nil(t, client.GetBody())
	assert.Equal(t, "", client.GetBodyString())
	assert.Nil(t, client.GetResponse())
}

func TestGetRestyClientAndRequest(t *testing.T) {
	client := NewHttpClient()
	assert.NotNil(t, client.GetRestyClient())
	assert.NotNil(t, client.GetRestyRequest())
}
