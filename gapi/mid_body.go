package gapi

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
)

func MiddleWare_BodyCaching() gin.HandlerFunc {
	//middleware to cache body
	return func(c *gin.Context) {
		// 读取请求体
		bodyBytes, err := c.GetRawData()
		if err != nil {
			c.AbortWithStatusJSON(200, httpApiError{
				Code:    1001,
				Message: "Cannot read body",
			})
			return
		}

		// 将 body 数据缓存
		c.Set(gin.BodyBytesKey, bodyBytes)

		// 将读取的 body 放回请求中，避免后续 handler 无法读取 body
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		c.Next()
		// c.ShouldBindBodyWith()
	}
}
