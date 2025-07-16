package http

import (
	"fmt"
	"musicboxapi/database"
	"musicboxapi/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FetchPlaylistSongs(ctx *gin.Context) {
	playlistIdParameter := ctx.Param("playlistId")

	lastKnowSongPosition := 0

	lastKnowSongPositionQuery := ctx.Query("lastKnowSongPosition")

	if lastKnowSongPositionQuery != "" {
		lastKnowSongPosition, _ = strconv.Atoi(lastKnowSongPositionQuery)
	}

	if playlistIdParameter == "" {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("No playlistId inrequest"))
		return
	}

	playlistId, _ := strconv.Atoi(playlistIdParameter)

	playlistsongTable := database.NewPlaylistsongTableInstance()

	songs, err := playlistsongTable.FetchPlaylistSongs(ctx.Request.Context(), playlistId, lastKnowSongPosition)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, models.OkResponse(songs, fmt.Sprintf("Found %d songs in playlist %d", len(songs), playlistId)))
}

func InsertPlaylistSong(ctx *gin.Context) {
	playlistIdParameter := ctx.Param("playlistId")
	songIdParameter := ctx.Param("songId")

	if playlistIdParameter == "" {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("playlistId is empty"))
		return
	}

	if songIdParameter == "" {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("songId is empty"))
		return
	}

	playlistId, _ := strconv.Atoi(playlistIdParameter)
	songId, _ := strconv.Atoi(songIdParameter)

	playlistsongTable := database.NewPlaylistsongTableInstance()

	id, err := playlistsongTable.InsertPlaylistSong(playlistId, songId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, models.OkResponse(gin.H{"playlistSongId": id}, fmt.Sprintf("Added song %d to playlist %d", songId, playlistId)))
}

func DeletePlaylistSong(ctx *gin.Context) {
	playlistIdParameter := ctx.Param("playlistId")
	songIdParameter := ctx.Param("songId")

	if playlistIdParameter == "" {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("playlistId is empty"))
		return
	}

	if songIdParameter == "" {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("songId is empty"))
		return
	}

	playlistId, _ := strconv.Atoi(playlistIdParameter)
	songId, _ := strconv.Atoi(songIdParameter)

	playlistsongTable := database.NewPlaylistsongTableInstance()

	err := playlistsongTable.DeletePlaylistSong(playlistId, songId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}
	ctx.Status(http.StatusOK)
}
