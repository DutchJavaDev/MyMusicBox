package main

import (
	"api/http"
	"api/logging"
	"api/util"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/lrstanley/go-ytdlp"
)

func main() {
	config := util.GetConfig()
	// If yt-dlp isn't installed yet, download and cache it for further use.
	ytdlp.MustInstall(context.TODO(), nil)

	if config.UseDevUrl {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	engine.SetTrustedProxies(nil)

	apiGroup := engine.Group(util.GetApiGroupUrlV1(config.UseDevUrl))

	apiGroup.POST("/add/song", http.AddSong)

	apiGroup.GET("/songs", http.GetSongs)

	apiGroup.GET("/playlist/download", http.DownloadPlaylist)

	apiGroup.GET("/dryrun", http.DryRun)

	if config.DevPort != "" {
		devPort := "127.0.0.1:" + config.DevPort
		logging.Info("Running on development port")
		engine.Run(devPort)
	} else {
		engine.Run() // listen and serve on 0.0.0.0:8080
	}

}
