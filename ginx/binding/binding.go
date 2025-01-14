package binding

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hchicken/pkg-go/ginx/response"
	"github.com/pkg/errors"
)

// handleError 处理错误并返回
func handleError(c *gin.Context, err error, message string) error {
	if err != nil {
		response.Json(c, response.Error(errors.Wrapf(err, message)))
	}
	return err
}

// ShouldBindJSON json数据验证
func ShouldBindJSON(c *gin.Context, d interface{}) error {
	err := c.ShouldBindJSON(d)
	return handleError(c, err, "参数认证失败")
}

// ShouldBindQuery form数据验证
func ShouldBindQuery(c *gin.Context, d interface{}) error {
	err := c.ShouldBindQuery(d)
	return handleError(c, err, "参数认证失败")
}

// ShouldBindPath 支持路径数据获取
func ShouldBindPath[T string | int](c *gin.Context, paramName string, value *T) error {
	paramValue := c.Param(paramName)
	if paramValue == "" {
		return handleError(c, fmt.Errorf("%v is null", paramName), "参数认证失败")
	}

	// 通过数据类型来判断
	switch v := any(value).(type) {
	case *int:
		intValue, err := strconv.Atoi(paramValue)
		if err != nil {
			return handleError(c, fmt.Errorf("failed to parse %v as int: %v", paramName, err), "参数认证失败")
		}
		*v = intValue
	case *string:
		*v = paramValue
	default:
		return handleError(c, fmt.Errorf("unsupported type for %v", paramName), "参数认证失败")
	}

	return nil
}
