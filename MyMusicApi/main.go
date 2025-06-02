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
	util.LoadConfig()
	// If yt-dlp isn't installed yet, download and cache it for further use.
	ytdlp.MustInstall(context.TODO(), nil)

	if util.Config.UseDevUrl {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	engine.SetTrustedProxies(nil)

	apiv1Group := engine.Group(util.GetApiGroupUrlV1(util.Config.UseDevUrl))

	apiv1Group.GET("/songs", http.FetchSongs)
	apiv1Group.GET("/playlist", http.FetchPlaylists)
	apiv1Group.GET("/playlist/:playlistId", http.FetchPlaylistSongs)

	apiv1Group.POST("/playlist", http.InsertPlaylist)
	apiv1Group.POST("/playlistsong/:playlistId/:songId", http.InsertPlaylistSong)
	apiv1Group.POST("/download", http.DownloadRequest)

	apiv1Group.DELETE("/playlist/:playlistId", http.DeletePlaylist)
	apiv1Group.DELETE("playlistsong/:playlistId/:songId", http.DeletePlaylistSong)

	// Ignore
	apiv1Group.GET("/playlist/download", http.DownloadPlaylist)
	// apiv1Group.GET("/dryrun", http.DryRun)
	apiv1Group.POST("/add/song", http.AddSong)

	if util.Config.DevPort != "" {
		devPort := "127.0.0.1:" + util.Config.DevPort
		logging.Info("Running on development port")
		engine.Run(devPort)
	} else {
		engine.Run() // listen and serve on 0.0.0.0:8080
	}

}
