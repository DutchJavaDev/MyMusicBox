package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/lrstanley/go-ytdlp"
)

func test() {}

func main() {
	// If yt-dlp isn't installed yet, download and cache it for further use.
	ytdlp.MustInstall(context.TODO(), nil)
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/randonData", func(ctx *gin.Context) {

			dl := ytdlp.New().
				FormatSort("res,ext:mp4:m4a").
				RecodeVideo("mp4").
				Output("%(extractor)s - %(title)s.%(ext)s").
				Cookies("cookies.txt")

			result, err := dl.Run(context.TODO(), "https://www.youtube.com/watch?v=dQw4w9WgXcQ")
			if err != nil {
				ctx.JSON(500, gin.H{
					"message": err.Error(),
				})
			} else {
				ctx.JSON(200, gin.H{
					"message": result,
				})
			}
		})
	}
	router.Run() // listen and serve on 0.0.0.0:8080
}
