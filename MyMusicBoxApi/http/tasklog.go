package http

import (
	"fmt"
	"musicboxapi/database"
	"musicboxapi/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskLogHandler struct {
	TasklogTable database.ITasklogTable
}

func (handler *TaskLogHandler) FetchParentTaskLogs(ctx *gin.Context) {

	logs, err := handler.TasklogTable.GetParentLogs(ctx.Request.Context())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
	}

	ctx.JSON(http.StatusOK, models.OkResponse(logs, fmt.Sprintf("%d", len(logs))))
}

func (handler *TaskLogHandler) FetchChildTaskLogs(ctx *gin.Context) {
	parentIdParameter := ctx.Param("parentId")

	parentId, err := strconv.Atoi(parentIdParameter)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	logs, err := handler.TasklogTable.GetChildLogs(ctx, parentId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.OkResponse(logs, fmt.Sprintf("%d", len(logs))))
}
