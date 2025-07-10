package main

import (
	"context"
	"fmt"
	"musicboxapi/configuration"
	"musicboxapi/http"
	"musicboxapi/logging"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lrstanley/go-ytdlp"
)

func main() {
	// If yt-dlp isn't installed yet, download and cache it for further use.
	go ytdlp.MustInstall(context.TODO(), nil)

	configuration.LoadConfig()

	if configuration.Config.UseDevUrl {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	engine.SetTrustedProxies([]string{"127.0.0.1"})

	if configuration.Config.UseDevUrl {
		engine.Use(cors.Default())
		logging.Warning("CORS is enabled for all origins")
	} else {

		origin := os.Getenv("CORS_ORIGIN")

		// Use Default cors
		if len(origin) == 0 {
			engine.Use(cors.Default())
			logging.Warning("CORS is enabled for all origins")
		} else {
			strictCors := cors.New(cors.Config{
				AllowAllOrigins: false,
				AllowOrigins:    []string{origin}, // move to env
			})
			engine.Use(strictCors)
		}
	}

	apiv1Group := engine.Group(configuration.GetApiGroupUrl("v1"))
	apiv1Group.GET("/songs", http.FetchSongs)
	apiv1Group.GET("/playlist", http.FetchPlaylists)
	apiv1Group.GET("/playlist/:playlistId", http.FetchPlaylistSongs)
	apiv1Group.GET("/play/:sourceId", http.Play)
	apiv1Group.GET("/tasklogs", http.FetchTaskLogs)

	apiv1Group.POST("/playlist", http.InsertPlaylist)
	apiv1Group.POST("/playlistsong/:playlistId/:songId", http.InsertPlaylistSong)
	apiv1Group.POST("/download", http.DownloadRequest)

	apiv1Group.DELETE("/playlist/:playlistId", http.DeletePlaylist)
	apiv1Group.DELETE("playlistsong/:playlistId/:songId", http.DeletePlaylistSong)

	// Serving static files
	apiv1Group.Static("/images", fmt.Sprintf("%s/%s", configuration.Config.SourceFolder, "images"))

	if configuration.Config.DevPort != "" {
		devPort := "127.0.0.1:" + configuration.Config.DevPort
		logging.Info("Running on development port")
		engine.Run(devPort)
	} else {
		engine.Run() // listen and serve on 0.0.0.0:8080
	}
}
