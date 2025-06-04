package http

import (
	"musicboxapi/database"
	"musicboxapi/models"
	"musicboxapi/service"

	"github.com/gin-gonic/gin"
)

func DownloadRequest(ctx *gin.Context) {
	var request models.DownloadRequestModel
	err := ctx.ShouldBindBodyWithJSON(&request)

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
