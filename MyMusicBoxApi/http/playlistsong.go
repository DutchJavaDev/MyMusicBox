package http

import (
	"fmt"
	"musicboxapi/database"
	"musicboxapi/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func InsertPlaylistSong(ctx *gin.Context) {
	db := database.PostgresDb{}
	defer db.CloseConnection()

	playlistId, _ := strconv.Atoi(ctx.Param("playlistId"))
	songId, _ := strconv.Atoi(ctx.Param("songId"))

	if db.OpenConnection() {
		id, err := db.InsertPlaylistSong(playlistId, songId)

		if err != nil {
			ctx.JSON(500, models.ErrorResponse(err))
			return
		}
		ctx.JSON(200, models.OkResponse(gin.H{"playlistSongId": id}, fmt.Sprintf("Added song %d to playlist %d", songId, playlistId)))
	} else {
		ctx.JSON(500, models.ErrorResponse(db.Error))
	}
}

func DeletePlaylistSong(ctx *gin.Context) {
	db := database.PostgresDb{}
	defer db.CloseConnection()

	playlistId, _ := strconv.Atoi(ctx.Param("playlistId"))
	songId, _ := strconv.Atoi(ctx.Param("songId"))

	if db.OpenConnection() {
		err := db.DeletePlaylistSong(playlistId, songId)
		if err != nil {
			ctx.JSON(500, models.ErrorResponse(err))
			return
		}
		ctx.Status(200)
	} else {
		ctx.JSON(500, models.ErrorResponse(db.Error))
	}
}
