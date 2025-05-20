package main

import (
	"api/db"
	"api/http"
	"api/util"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lrstanley/go-ytdlp"
)

func main() {
	config := util.GetConfig()

	db.CreateDatabase()

	// If yt-dlp isn't installed yet, download and cache it for further use.
	ytdlp.MustInstall(context.TODO(), nil)

	if config.UseDevUrl {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	engine.SetTrustedProxies(nil)

	apiGroup := engine.Group(util.GetApiGroupUrl(config.UseDevUrl))

	apiGroup.POST("/music", http.AddMusic)

	apiGroup.GET("/music", http.GetMusic)

	apiGroup.GET("/playlist/download", http.DownloadPlaylist)

	apiGroup.GET("/dryrun", http.DryRun)

	if config.DevPort != "" {
		devPort := "127.0.0.1:" + config.DevPort
		fmt.Println("Running on development port")
		engine.Run(devPort)
	} else {
		engine.Run() // listen and serve on 0.0.0.0:8080
	}

}
