package main

import (
	"api/http"
	"api/logging"
	"api/util"
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lrstanley/go-ytdlp"
)

func main() {
	config := util.GetConfig()

	// If yt-dlp isn't installed yet, download and cache it for further use.
	ytdlp.MustInstall(context.TODO(), nil)

	if config.UseDevUrl {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()
	engine.SetTrustedProxies(nil)

	apiGroup := engine.Group(util.GetApiGroupUrl(config.UseDevUrl))

	apiGroup.GET("/test", func(ctx *gin.Context) {

		cookiesPath := "selenium/cookies_netscape"

		contents, err := os.ReadFile(cookiesPath)

		if err != nil {
			logging.Log(err.Error())
			ctx.Writer.WriteHeader(500)
			ctx.Writer.WriteString(err.Error())
			ctx.Writer.Flush()
			return
		}

		ctx.Writer.WriteHeader(200)
		ctx.Writer.WriteString(string(contents))
		ctx.Writer.Flush()
	})

	apiGroup.POST("/data", func(ctx *gin.Context) {
		var url http.UrlRequest

		err := ctx.ShouldBindBodyWithJSON(&url)

		if err != nil {
			message := fmt.Sprintf("Error parsing json, %s", err.Error())
			logging.Log(message)
			response := http.ErrorResponse{
				Message: message,
			}
			ctx.JSON(500, response)
			return
		}

		if url.Url == "" {
			message := "Empty request url"
			response := http.ErrorResponse{
				Message: message,
			}
			logging.Log(message)
			ctx.JSON(500, response)
			return
		}

		ctx.JSON(200, url)
	})

	apiGroup.GET("/randonData", func(ctx *gin.Context) {

		cookiesPath := "selenium/cookies_netscape"

		dl := ytdlp.New().
			FormatSort("res,ext:mp4:m4a").
			ExtractAudio().
			AudioQuality("0").
			AudioFormat("opus").
			NoKeepVideo().
			// Output("%(extractor)s - %(title)s.%(ext)s").
			Output("%(title)s.%(ext)s").
			Cookies(string(cookiesPath))

		result, err := dl.Run(context.TODO(), "https://www.youtube.com/watch?v=dQw4w9WgXcQ")
		if err != nil {
			logging.Log(err.Error())
			ctx.JSON(500, gin.H{
				"message": err.Error(),
			})
		} else {
			ctx.JSON(200, gin.H{
				"message": result,
			})
		}
	})

	if config.DevPort != "" {
		devPort := "127.0.0.1:" + config.DevPort
		fmt.Println("Running on development port")
		engine.Run(devPort)
	} else {
		engine.Run() // listen and serve on 0.0.0.0:8080
	}

}
