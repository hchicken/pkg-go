package response

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// ShouldBindJSON json数据验证
func ShouldBindJSON(c *gin.Context, d interface{}) error {
	err := c.ShouldBindJSON(d)
	if err != nil {
		Json(c, Error(errors.Wrapf(err, "参数认证失败,确定请求数据是否为json")))
	}
	return err
}

// ShouldBindQuery form数据验证
func ShouldBindQuery(c *gin.Context, d interface{}) error {
	err := c.ShouldBindQuery(d)
	if err != nil {
		Json(c, Error(errors.Wrapf(err, "参数认证失败")))
	}
	return err
}
