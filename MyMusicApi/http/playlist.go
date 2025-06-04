package http

import (
	"fmt"
	"musicboxapi/database"
	"musicboxapi/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FetchPlaylists(ctx *gin.Context) {
	db := database.PostgresDb{}
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
	db := database.PostgresDb{}
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

	db := database.PostgresDb{}
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

func DeletePlaylist(ctx *gin.Context) {
	db := database.PostgresDb{}
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
