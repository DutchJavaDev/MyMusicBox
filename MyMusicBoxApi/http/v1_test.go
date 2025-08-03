package http

import "github.com/gin-gonic/gin"

func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	ginEngine := gin.Default()
	return ginEngine
}
