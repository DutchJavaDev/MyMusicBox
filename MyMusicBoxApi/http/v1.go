package http

import (
	"musicboxapi/configuration"
	"musicboxapi/database"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func V1Endpoints(apiv1Group *gin.RouterGroup) {

	songHandler := SongHandler{
		SongTable: database.NewSongTableInstance(),
	}
	playlistHandler := PlaylistHandler{
		PlaylistTable: database.NewPlaylistTableInstance(),
	}
	playlistSongHandler := PlaylistSongHandler{
		PlaylistsongTable: database.NewPlaylistsongTableInstance(),
	}
	taskLogHandler := TaskLogHandler{
		TasklogTable: database.NewTasklogTableInstance(),
	}

	apiv1Group.GET("/songs", songHandler.FetchSongs)

	apiv1Group.GET("/playlist", playlistHandler.FetchPlaylists)
	apiv1Group.GET("/playlist/:playlistId", playlistSongHandler.FetchPlaylistSongs)

	// If nginx is not configured to handle http 206
	if configuration.Config.UsePlayUrl {
		apiv1Group.GET("/play/:sourceId", Play)
	}

	apiv1Group.GET("/tasklogs", taskLogHandler.FetchParentTaskLogs)
	apiv1Group.GET("/tasklogs/:parentId", taskLogHandler.FetchChildTaskLogs)

	apiv1Group.POST("/playlist", playlistHandler.InsertPlaylist)
	apiv1Group.POST("/playlist/sync/:playlistId", playlistHandler.SyncPlaylist)
	apiv1Group.POST("/playlistsong/:playlistId/:songId", playlistSongHandler.InsertPlaylistSong)
	apiv1Group.POST("/download", DownloadRequest)

	apiv1Group.DELETE("/playlist/:playlistId", playlistHandler.DeletePlaylist)
	apiv1Group.DELETE("playlistsong/:playlistId/:songId", playlistSongHandler.DeletePlaylistSong)

	// If nginx is not configured to handle it
	if configuration.Config.UseImageUrl {
		// Serving static files
		apiv1Group.Static("/images", filepath.Join(configuration.Config.SourceFolder, "images"))
	}
}
