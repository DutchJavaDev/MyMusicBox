package main

import (
	"context"
	"musicboxapi/configuration"
	"musicboxapi/http"
	"musicboxapi/logging"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lrstanley/go-ytdlp"
)

func main() {
	configuration.LoadConfig()

	// If yt-dlp isn't installed yet, download and cache it for further use.
	ytdlp.MustInstall(context.TODO(), nil)

	if configuration.Config.UseDevUrl {
		gin.SetMode(gin.DebugMode)

	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	engine.SetTrustedProxies(nil)

	if configuration.Config.UseDevUrl {
		engine.Use(cors.Default())
	}

	apiv1Group := engine.Group(configuration.GetApiGroupUrlV1(configuration.Config.UseDevUrl))

	apiv1Group.GET("/songs", http.FetchSongs)
	apiv1Group.GET("/playlist", http.FetchPlaylists)
	apiv1Group.GET("/playlist/:playlistId", http.FetchPlaylistSongs)
	apiv1Group.GET("/play/:sourceId", http.Play)

	apiv1Group.POST("/playlist", http.InsertPlaylist)
	apiv1Group.POST("/playlistsong/:playlistId/:songId", http.InsertPlaylistSong)
	apiv1Group.POST("/download", http.DownloadRequest)

	apiv1Group.DELETE("/playlist/:playlistId", http.DeletePlaylist)
	apiv1Group.DELETE("playlistsong/:playlistId/:songId", http.DeletePlaylistSong)

	if configuration.Config.DevPort != "" {
		devPort := "127.0.0.1:" + configuration.Config.DevPort
		logging.Info("Running on development port")
		engine.Run(devPort)
	} else {
		engine.Run() // listen and serve on 0.0.0.0:8080
	}

}
