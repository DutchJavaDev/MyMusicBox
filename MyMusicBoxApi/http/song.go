package http

import (
	"fmt"
	"musicboxapi/database"
	"musicboxapi/logging"
	"musicboxapi/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FetchSongs(ctx *gin.Context) {
	songTable := database.NewSongTableInstance()

	songs, err := songTable.FetchSongs(ctx.Request.Context())

	if err != nil {
		logging.Error(fmt.Sprintf("Failed to fetch songs: %s", err.Error()))
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to fetch songs"))
		return
	}

	ctx.JSON(http.StatusOK, models.OkResponse(songs, fmt.Sprintf("Found %d songs", len(songs))))
}
