package gapi

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type httpApiError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func Error(c *gin.Context, err error) {
	if err == nil {
		c.AbortWithStatusJSON(200, httpApiError{Code: 1003, Message: "未知错误"})
		return
	}
	switch err := err.(type) {
	case validator.ValidationErrors:
		kv := translateErrors(err)
		//把kv的value拼接成字符串，用\n分割，strings.TrimSpace的方式不严谨
		errStr := ""
		for _, v := range kv {
			errStr += v + "\n"
		}
		errStr = strings.TrimRight(errStr, "\n")

		for k := range kv {
			if strings.ToLower(k) == "user-agent" || strings.HasPrefix(strings.ToLower(k), "x-app-") || strings.HasPrefix(strings.ToLower(k), "x-api-") {
				c.AbortWithStatusJSON(200, httpApiError{Code: 1003, Message: errStr})
				return
			}
		}
		c.AbortWithStatusJSON(200, httpApiError{Code: 1004, Message: errStr})
	case appError:
		switch err.source {
		case appErrorSourceSign:
			c.AbortWithStatusJSON(200, httpApiError{Code: 1003, Message: err.Error()})
		case appErrorSourceValidate:
			c.AbortWithStatusJSON(200, httpApiError{Code: 1004, Message: err.Error()})
		default:
			c.AbortWithStatusJSON(200, httpApiError{Code: 1005, Message: err.Error()})
		}
	default:
		c.AbortWithStatusJSON(200, httpApiError{Code: 1005, Message: err.Error()})
	}
}

func Data(c *gin.Context, T any, message ...string) {
	if len(message) == 0 {
		c.JSON(200, gin.H{"code": 1000, "data": T})
	} else {
		c.JSON(200, gin.H{"code": 1000, "data": T, "msg": message[0]})
	}
}
