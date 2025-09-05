package http

import (
	"fmt"
	"musicboxapi/database"
	"musicboxapi/logging"
	"musicboxapi/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SongHandler struct {
	SongTable database.ISongTable
}

// @Produce json
// @Description Returns data for all songs
// @Success 200 {object} models.Song
// @Failure 500 {object} models.ApiResponseModel
// @Router /api/v1/songs [get]
func (handler *SongHandler) FetchSongs(ctx *gin.Context) {
	songs, err := handler.SongTable.FetchSongs(ctx.Request.Context())
	if err != nil {
		logging.Error(fmt.Sprintf("Failed to fetch songs: %s", err.Error()))
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to fetch songs"))
		return
	}
	ctx.JSON(http.StatusOK, models.OkResponse(songs, fmt.Sprintf("Found %d songs", len(songs))))
}
