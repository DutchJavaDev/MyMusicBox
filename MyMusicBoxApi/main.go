package main

import (
	"fmt"
	"musicboxapi/configuration"
	"musicboxapi/database"
	"musicboxapi/http"
	"musicboxapi/logging"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title MusicBoxApi API
// @version 1.0
// @BasePath  /api/v1
func main() {

	configuration.LoadConfiguration()

	err := database.CreateDatabasConnectionPool()
	defer database.DbInstance.Close()

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to create database connection: %s", err.Error())
		logging.Error(errorMessage)
		logging.ErrorStackTrace(err)
		return
	}

	database.ApplyMigrations()

	setGinMode()

	ginEngine := gin.Default()

	// Trust nginx
	ginEngine.SetTrustedProxies([]string{"127.0.0.1"})

	ginEngine.Use(corsMiddelWare())

	// V1 API
	apiv1Group := ginEngine.Group(configuration.GetApiGroupUrl("v1"))
	http.V1Endpoints(apiv1Group)

	if configuration.Config.DevPort != "" {
		devPort := "127.0.0.1:" + configuration.Config.DevPort
		logging.Info("Running on development port")
		ginEngine.Run(devPort)
	} else {
		ginEngine.Run() // listen and serve on 0.0.0.0:8080
	}
}

func corsMiddelWare() gin.HandlerFunc {
	if configuration.Config.UseDevUrl {
		return cors.Default()
	} else {
		origin := os.Getenv("CORS_ORIGIN")
		// Use Default cors
		if len(origin) == 0 {
			logging.Warning("CORS is enabled for all origins")
			return cors.Default()
		} else {
			strictCors := cors.New(cors.Config{
				AllowAllOrigins:  false,
				AllowOrigins:     []string{origin}, // move to env
				AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
				AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
				AllowCredentials: false,
				MaxAge:           12 * time.Hour,
			})
			return strictCors
		}
	}
}

func setGinMode() {
	if configuration.Config.UseDevUrl {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}
