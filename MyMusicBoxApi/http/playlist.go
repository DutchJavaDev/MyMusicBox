package http

import (
	"fmt"
	"musicboxapi/configuration"
	"musicboxapi/database"
	"musicboxapi/logging"
	"musicboxapi/models"
	"musicboxapi/service"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PlaylistHandler struct {
	PlaylistTable database.IPlaylistTable
}

// @Produce json
// @Param lastKnowPlaylistId   path      int  false  "Last know playlist id by the client, default is 0"
// @Description Returns data for all playlist, if lastKnowPlaylistId then only the playlist after lastKnowPlaylistId
// @Success 200 {object} models.Playlist
// @Failure 500 {object} models.ApiResponseModel
// @Router /api/v1/playlist [get]
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

func (handler *PlaylistHandler) SyncPlaylist(ctx *gin.Context) {

	playlistIdParameter := ctx.Param("playlistId")

	id, err := strconv.Atoi(playlistIdParameter)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	if DefaultPlaylistId == id {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("Funky music... wrong year"))
		return
	}

	playlistTable := database.NewPlaylistTableInstance()

	playlist, err := playlistTable.FetchPlaylistsById(ctx.Request.Context(), id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	playlistId := strings.Split(playlist.ThumbnailPath, ".")[0]

	var request models.DownloadRequestModel

	request.Url = fmt.Sprintf("https://www.youtube.com/playlist?list=%s", playlistId)

	go service.StartDownloadTask(request)

	ctx.JSON(http.StatusOK, models.OkResponse(gin.H{"": ""}, "Created"))
}

func (hanlder *PlaylistHandler) InsertPlaylist(ctx *gin.Context) {

	var playlistModel models.CreatePlaylistModel

	err := ctx.ShouldBind(&playlistModel)

	if err != nil {
		logging.Error("Failed to bind model")
		logging.ErrorStackTrace(err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	var playlist models.Playlist
	var fileName string

	hasFormFile := playlistModel.Image.Size > 0

	if hasFormFile {
		fileName = fmt.Sprintf("%s.jpg", uuid.New().String())
		playlist.ThumbnailPath = fileName
	} else {
		// default_playlist_cover.jpg
		playlist.ThumbnailPath = "default_playlist_cover.jpg"
	}

	playlist.Name = playlistModel.Name
	playlist.Description = playlistModel.Description

	if strings.Contains(playlistModel.IsPublic, "on") {
		playlist.IsPublic = true
	}

	playlistId, err := hanlder.PlaylistTable.InsertPlaylist(playlist)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	if hasFormFile {
		path := filepath.Join(configuration.Config.SourceFolder, fmt.Sprintf("images/%s", fileName))

		logging.Info(path)

		err = ctx.SaveUploadedFile(playlistModel.Image, path)

		if err != nil {
			logging.ErrorStackTrace(err)
			ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
			return
		}
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

	if DefaultPlaylistId == id {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("Funky music... "))
		return
	}

	err = handler.PlaylistTable.DeletePlaylist(id)

	// TODO delete background image if its not the default image for it

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
