package http

import (
	"fmt"
	"musicboxapi/database"
	"musicboxapi/models"

	"github.com/gin-gonic/gin"
)

func FetchTaskLogs(ctx *gin.Context) {
	db := database.PostgresDb{}
	defer db.CloseConnection()

	logs, err := db.GetTaskLogs(ctx.Request.Context())

	if err != nil {
		ctx.JSON(500, models.ErrorResponse(err.Error()))
	}

	ctx.JSON(200, models.OkResponse(logs, fmt.Sprintf("%d", len(logs))))
}
