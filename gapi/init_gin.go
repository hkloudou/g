package gapi

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var GinEngine *gin.Engine

func init() {
	GinEngine = gin.Default()
	cfg := cors.DefaultConfig()
	cfg.AllowAllOrigins = true
	// cfg.AddAllowHeaders("token", "device_id", "*")
	cfg.AddAllowHeaders("*")
	cfg.AddExposeHeaders("*")
	GinEngine.Use(cors.New(cfg))
}
