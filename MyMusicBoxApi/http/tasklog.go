package http

import (
	"fmt"
	"musicboxapi/database"
	"musicboxapi/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskLogHandler struct {
	TasklogTable database.ITasklogTable
}

func (handler *TaskLogHandler) FetchTaskLogs(ctx *gin.Context) {

	logs, err := handler.TasklogTable.GetParentChildLogs(ctx.Request.Context())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
	}

	ctx.JSON(http.StatusOK, models.OkResponse(logs, fmt.Sprintf("%d", len(logs))))
}
