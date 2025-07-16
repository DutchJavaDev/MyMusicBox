package http

import (
	"fmt"
	"musicboxapi/database"
	"musicboxapi/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FetchTaskLogs(ctx *gin.Context) {
	tasklogTable := database.NewTasklogTableInstance()

	logs, err := tasklogTable.GetTaskLogs(ctx.Request.Context())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
	}

	ctx.JSON(http.StatusOK, models.OkResponse(logs, fmt.Sprintf("%d", len(logs))))
}
