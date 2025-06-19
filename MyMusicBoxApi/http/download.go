package http

import (
	"errors"
	"musicboxapi/database"
	"musicboxapi/models"
	"musicboxapi/service"
	"strings"

	"github.com/gin-gonic/gin"
)

func DownloadRequest(ctx *gin.Context) {
	var request models.DownloadRequestModel
	err := ctx.ShouldBindBodyWithJSON(&request)

	// If it contains &list= it will download but it will not update the database for all entries or create the playlist entry
	if strings.Contains(request.Url, "&list=") {
		ctx.JSON(500, models.ErrorResponse(errors.New("Url format is wrong, containt &list= instead of playlist")))
		return
	}

	if err != nil {
		ctx.JSON(500, models.ErrorResponse(err))
	} else {
		db := database.PostgresDb{}
		defer db.CloseConnection()

		if db.OpenConnection() {
			// Insert a new task
			taskId, err := db.InsertTaskLog()

			if err != nil {
				ctx.JSON(500, models.ErrorResponse(err))
				return
			}

			go service.StartDownloadTask(taskId, request)
			ctx.JSON(200, models.OkResponse(gin.H{"taskId": taskId}, "Started task"))
		} else {
			ctx.JSON(500, models.ErrorResponse(db.Error))
		}
	}
}
