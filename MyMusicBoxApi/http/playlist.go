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

	lastKnowPlaylistIdQuery := ctx.Query("lastKnowPlaylistId")

	lastKnowPlaylistId := 0

	if lastKnowPlaylistIdQuery != "" {
		lastKnowPlaylistId, _ = strconv.Atoi(lastKnowPlaylistIdQuery)
	}

	if db.OpenConnection() {
		playlists, err := db.FetchPlaylists(ctx.Request.Context(), lastKnowPlaylistId)
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

	playlistIdParameter := ctx.Param("playlistId")

	lastKnowSongPosition := 0

	lastKnowSongPositionQuery := ctx.Query("lastKnowSongPosition")

	if lastKnowSongPositionQuery != "" {
		lastKnowSongPosition, _ = strconv.Atoi(lastKnowSongPositionQuery)
	}

	if playlistIdParameter == "" {
		ctx.JSON(500, models.ErrorResponse("No playlistId inrequest"))
		return
	}

	playlistId, _ := strconv.Atoi(playlistIdParameter)

	if db.OpenConnection() {
		songs, err := db.FetchPlaylistSongs(ctx.Request.Context(), playlistId, lastKnowSongPosition)
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

	err := ctx.ShouldBindBodyWithJSON(&playlist)

	if err != nil {
		ctx.JSON(500, models.ErrorResponse(err))
		return
	}

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

	playlistIdParameter := ctx.Param("playlistId")

	if playlistIdParameter == "" {
		ctx.JSON(500, models.ErrorResponse("playlistId is empty"))
		return
	}

	id, _ := strconv.Atoi(playlistIdParameter)

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
