package binding

import (
	"github.com/gin-gonic/gin"
	"github.com/hchicken/pkg-go/ginx/response"
	"github.com/pkg/errors"
)

// ShouldBindJSON json数据验证
func ShouldBindJSON(c *gin.Context, d interface{}) error {
	err := c.ShouldBindJSON(d)
	if err != nil {
		response.Json(c, response.Error(errors.Wrapf(err, "参数认证失败,确定请求数据是否为json")))
	}
	return err
}

// ShouldBindQuery form数据验证
func ShouldBindQuery(c *gin.Context, d interface{}) error {
	err := c.ShouldBindQuery(d)
	if err != nil {
		response.Json(c, response.Error(errors.Wrapf(err, "参数认证失败")))
	}
	return err
}
