package http

import (
	"musicboxapi/models"
	"musicboxapi/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Accepts json
// @Produce json
// @Param download body string true "message/rfc8259 see models.DownloadRequestModel"
// @Description Enables playback for song/file using http 206 partial content
// @Success 200 "serve song/file with range request (http 206)"
// @Success 0 {object} models.DownloadRequestModel
// @Failure 500 {object} models.ApiResponseModel
// @Router /api/v1/download [get]
func DownloadRequest(ctx *gin.Context) {
	var request models.DownloadRequestModel
	err := ctx.ShouldBindBodyWithJSON(&request)

	if err != nil {
		ctx.JSON(500, models.ErrorResponse(err))
		return
	}

	// If it contains &list= it will download but it will not update the database for all entries or create the playlist entry
	if strings.Contains(request.Url, "&list=") && strings.Contains(request.Url, "watch?") {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(gin.H{"error": "Url format is wrong, contains watch instead of playlist="}))
		return
	}

	//tasklogTable := database.NewTasklogTableInstance()
	// Insert a new task
	// parentTask, err := tasklogTable.CreateParentTaskLog(request.Url)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	go service.StartDownloadTask(request)
	ctx.JSON(http.StatusOK, models.OkResponse(gin.H{"": ""}, "Created"))
}
