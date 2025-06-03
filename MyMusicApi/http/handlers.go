package http

import (
	"api/db"
	"api/logging"
	"api/models"
	"api/util"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FetchSongs(ctx *gin.Context) {
	db := db.PostgresDb{}

	defer db.CloseConnection()

	if db.OpenConnection() {
		songs, err := db.FetchSongs(ctx.Request.Context())
		if err != nil {
			ctx.JSON(500, models.ErrorResponse(err.Error()))
			return
		}
		ctx.JSON(200, models.OkResponse(songs, fmt.Sprintf("Found %d songs", len(songs))))
	} else {
		ctx.JSON(500, models.ErrorResponse(db.Error))
	}
}

func FetchPlaylists(ctx *gin.Context) {
	db := db.PostgresDb{}
	defer db.CloseConnection()

	if db.OpenConnection() {
		playlists, err := db.FetchPlaylists(ctx.Request.Context())
		if err != nil {
			ctx.JSON(500, models.ErrorResponse(err))
			return
		}
		ctx.JSON(200, models.OkResponse(playlists, fmt.Sprintf("Found %d playlist", len(playlists))))
	} else {
		ctx.JSON(500, models.ErrorResponse(db.Error))
	}
}

func FetchPlaylistSongs(ctx *gin.Context) {
	db := db.PostgresDb{}
	defer db.CloseConnection()

	playlistId, _ := strconv.Atoi(ctx.Param("playlistId"))

	if db.OpenConnection() {
		songs, err := db.FetchPlaylistSongs(ctx.Request.Context(), playlistId)
		if err != nil {
			ctx.JSON(500, models.ErrorResponse(err))
			return
		}
		ctx.JSON(200, models.OkResponse(songs, fmt.Sprintf("Found %d songs in playlist %d", len(songs), playlistId)))
	} else {
		ctx.JSON(500, models.ErrorResponse(db.Error))
	}
}

func InsertPlaylist(ctx *gin.Context) {

	var playlist models.Playlist

	ctx.ShouldBindBodyWithJSON(&playlist)

	db := db.PostgresDb{}
	defer db.CloseConnection()

	if db.OpenConnection() {
		id, err := db.InsertPlaylist(playlist)

		if err != nil {
			ctx.JSON(500, models.ErrorResponse(err))
			return
		}

		ctx.JSON(200, models.OkResponse(gin.H{"playlistId": id}, "Created new playlist"))
	} else {
		ctx.JSON(500, models.ErrorResponse(db.Error))
	}
}

func InsertPlaylistSong(ctx *gin.Context) {
	db := db.PostgresDb{}
	defer db.CloseConnection()

	playlistId, _ := strconv.Atoi(ctx.Param("playlistId"))
	songId, _ := strconv.Atoi(ctx.Param("songId"))

	if db.OpenConnection() {
		id, err := db.InsertPlaylistSong(playlistId, songId)

		if err != nil {
			ctx.JSON(500, models.ErrorResponse(err))
			return
		}
		ctx.JSON(200, models.OkResponse(gin.H{"playlistSongId": id}, fmt.Sprintf("Added song %d to playlist %d", songId, playlistId)))
	} else {
		ctx.JSON(500, models.ErrorResponse(db.Error))
	}
}

func DeletePlaylist(ctx *gin.Context) {
	db := db.PostgresDb{}
	defer db.CloseConnection()

	id, _ := strconv.Atoi(ctx.Param("playlistId"))

	if db.OpenConnection() {
		err := db.DeletePlaylist(id)
		if err != nil {
			ctx.JSON(500, models.ErrorResponse(err))
			return
		}
		ctx.Status(200)
	} else {
		ctx.JSON(500, models.ErrorResponse(db.Error))
	}
}

func DeletePlaylistSong(ctx *gin.Context) {
	db := db.PostgresDb{}
	defer db.CloseConnection()

	playlistId, _ := strconv.Atoi(ctx.Param("playlistId"))
	songId, _ := strconv.Atoi(ctx.Param("songId"))

	if db.OpenConnection() {
		err := db.DeletePlaylistSong(playlistId, songId)
		if err != nil {
			ctx.JSON(500, models.ErrorResponse(err))
			return
		}
		ctx.Status(200)
	} else {
		ctx.JSON(500, models.ErrorResponse(db.Error))
	}
}

func DownloadRequest(ctx *gin.Context) {
	var request models.DownloadRequestModel
	err := ctx.ShouldBindBodyWithJSON(&request)

	if err != nil {
		ctx.JSON(500, models.ErrorResponse(err))
	} else {
		db := db.PostgresDb{}
		defer db.CloseConnection()

		if db.OpenConnection() {
			// Insert a new task
			taskId, err := db.InsertTaskLog()

			if err != nil {
				ctx.JSON(500, models.ErrorResponse(err))
				return
			}

			go util.StartDownloadTask(taskId, request)
			ctx.JSON(200, models.OkResponse(gin.H{"taskId": taskId}, "Started task"))
		} else {
			ctx.JSON(500, models.ErrorResponse(db.Error))
		}
	}
}

func Play(ctx *gin.Context) {
	sourceId := ctx.Param("sourceId")

	if sourceId == "" {
		ctx.Status(404)
		logging.Warning(fmt.Sprintf("[Play] No sourceId in request"))
		return
	}

	path := fmt.Sprintf("%s/%s.%s", util.Config.SourceFolder, sourceId, util.Config.OutputExtension)

	file, err := os.Open(path)

	if err != nil {
		logging.Error(fmt.Sprintf("[Play] Failed to read file, %s", err.Error()))
		ctx.JSON(500, models.ErrorResponse(err))
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		logging.Error(fmt.Sprintf("[Play] Failed could not get fileinfo: %s", err.Error()))
		ctx.JSON(500, models.ErrorResponse(err))
		return
	}

	// Set proper headers for OPUS audio
	ctx.Header("Content-Type", "audio/opus")
	ctx.Header("Accept-Ranges", "bytes")

	// Use http.ServeContent which handles range requests automatically
	http.ServeContent(
		ctx.Writer,
		ctx.Request,
		fileInfo.Name(),
		fileInfo.ModTime(),
		file,
	)
}
