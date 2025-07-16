package http

import (
	"fmt"
	"musicboxapi/database"
	"musicboxapi/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FetchPlaylists(ctx *gin.Context) {
	lastKnowPlaylistIdQuery := ctx.Query("lastKnowPlaylistId")

	lastKnowPlaylistId := 0

	if lastKnowPlaylistIdQuery != "" {
		lastKnowPlaylistId, _ = strconv.Atoi(lastKnowPlaylistIdQuery)
	}

	playlistTable := database.NewPlaylistTableInstance()

	playlists, err := playlistTable.FetchPlaylists(ctx.Request.Context(), lastKnowPlaylistId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.OkResponse(playlists, fmt.Sprintf("Found %d playlist", len(playlists))))
}

func InsertPlaylist(ctx *gin.Context) {
	var playlist models.Playlist

	err := ctx.ShouldBindBodyWithJSON(&playlist)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	playlistTable := database.NewPlaylistTableInstance()

	playlistId, err := playlistTable.InsertPlaylist(playlist)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.OkResponse(gin.H{"playlistId": playlistId}, "Created new playlist"))
}

func DeletePlaylist(ctx *gin.Context) {
	playlistIdParameter := ctx.Param("playlistId")

	if playlistIdParameter == "" {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("playlistId is empty"))
		return
	}

	id, err := strconv.Atoi(playlistIdParameter)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	playlistTable := database.NewPlaylistTableInstance()

	err = playlistTable.DeletePlaylist(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
