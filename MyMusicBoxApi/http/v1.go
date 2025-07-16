package http

import (
	"musicboxapi/configuration"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func V1Endpoints(apiv1Group *gin.RouterGroup) {
	apiv1Group.GET("/songs", FetchSongs)
	apiv1Group.GET("/playlist", FetchPlaylists)
	apiv1Group.GET("/playlist/:playlistId", FetchPlaylistSongs)
	apiv1Group.GET("/play/:sourceId", Play)
	apiv1Group.GET("/tasklogs", FetchTaskLogs)

	apiv1Group.POST("/playlist", InsertPlaylist)
	apiv1Group.POST("/playlistsong/:playlistId/:songId", InsertPlaylistSong)
	apiv1Group.POST("/download", DownloadRequest)

	apiv1Group.DELETE("/playlist/:playlistId", DeletePlaylist)
	apiv1Group.DELETE("playlistsong/:playlistId/:songId", DeletePlaylistSong)

	// Serving static files
	apiv1Group.Static("/images", filepath.Join(configuration.Config.SourceFolder, "images"))
}
