package http

import (
	"fmt"
	"musicboxapi/database"
	"musicboxapi/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PlaylistHandler struct {
	PlaylistTable database.IPlaylistTable
}

func (handler *PlaylistHandler) FetchPlaylists(ctx *gin.Context) {
	lastKnowPlaylistIdQuery := ctx.Query("lastKnowPlaylistId")

	lastKnowPlaylistId := 0

	if lastKnowPlaylistIdQuery != "" {
		lastKnowPlaylistId, _ = strconv.Atoi(lastKnowPlaylistIdQuery)
	}

	playlists, err := handler.PlaylistTable.FetchPlaylists(ctx.Request.Context(), lastKnowPlaylistId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.OkResponse(playlists, fmt.Sprintf("Found %d playlist", len(playlists))))
}

func (hanlder *PlaylistHandler) InsertPlaylist(ctx *gin.Context) {
	var playlist models.Playlist

	err := ctx.ShouldBindBodyWithJSON(&playlist)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	playlistId, err := hanlder.PlaylistTable.InsertPlaylist(playlist)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.OkResponse(gin.H{"playlistId": playlistId}, "Created new playlist"))
}

func (handler *PlaylistHandler) DeletePlaylist(ctx *gin.Context) {
	playlistIdParameter := ctx.Param("playlistId")

	id, err := strconv.Atoi(playlistIdParameter)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	err = handler.PlaylistTable.DeletePlaylist(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
