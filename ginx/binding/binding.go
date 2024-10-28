package binding

import (
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
