package http

import (
	"fmt"
	"musicboxapi/database"
	"musicboxapi/models"

	"github.com/gin-gonic/gin"
)

func FetchSongs(ctx *gin.Context) {
	db := database.PostgresDb{}

	defer db.CloseConnection()

	if db.OpenConnection() {
		songs, err := db.FetchSongs(ctx.Request.Context())
		if err != nil {
			ctx.JSON(500, models.ErrorResponse(err.Error()))
			return
		}
		ctx.JSON(200, models.OkResponse(songs, fmt.Sprintf("Found %d songs", len(songs))))
	} else {
		ctx.JSON(500, models.ErrorResponse(db.Error))
	}
}
