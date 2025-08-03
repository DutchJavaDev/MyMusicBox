package http

import (
	"fmt"
	"musicboxapi/database"
	"musicboxapi/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PlaylistSongHandler struct {
	PlaylistsongTable database.IPlaylistsongTable
}

func (handler *PlaylistSongHandler) FetchPlaylistSongs(ctx *gin.Context) {
	playlistIdParameter := ctx.Param("playlistId")

	lastKnowSongPosition := 0

	// Optional
	lastKnowSongPositionQuery := ctx.Query("lastKnowSongPosition")

	if lastKnowSongPositionQuery != "" {
		lastKnowSongPosition, _ = strconv.Atoi(lastKnowSongPositionQuery)
	}

	playlistId, err := strconv.Atoi(playlistIdParameter)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	songs, err := handler.PlaylistsongTable.FetchPlaylistSongs(ctx.Request.Context(), playlistId, lastKnowSongPosition)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, models.OkResponse(songs, fmt.Sprintf("Found %d songs in playlist %d", len(songs), playlistId)))
}

func (handler *PlaylistSongHandler) InsertPlaylistSong(ctx *gin.Context) {
	playlistIdParameter := ctx.Param("playlistId")
	songIdParameter := ctx.Param("songId")

	playlistId, err := strconv.Atoi(playlistIdParameter)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	songId, err := strconv.Atoi(songIdParameter)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	id, err := handler.PlaylistsongTable.InsertPlaylistSong(playlistId, songId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, models.OkResponse(gin.H{"playlistSongId": id}, fmt.Sprintf("Added song %d to playlist %d", songId, playlistId)))
}

func (handler *PlaylistSongHandler) DeletePlaylistSong(ctx *gin.Context) {
	playlistIdParameter := ctx.Param("playlistId")
	songIdParameter := ctx.Param("songId")

	playlistId, err := strconv.Atoi(playlistIdParameter)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	songId, err := strconv.Atoi(songIdParameter)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	err = handler.PlaylistsongTable.DeletePlaylistSong(playlistId, songId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}
	ctx.Status(http.StatusOK)
}
