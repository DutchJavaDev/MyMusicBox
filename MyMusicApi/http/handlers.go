package http

import (
	"api/db"
	"api/logging"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/lrstanley/go-ytdlp"
)

func AddMusic(ctx *gin.Context) {
	var musicData db.Music

	_ = ctx.ShouldBindBodyWithJSON(&musicData)

	db.InsertMusic(&musicData)

	ctx.JSON(200, gin.H{"message": "done"})
}

func GetMusic(ctx *gin.Context) {
	music := db.GetMusic()
	ctx.JSON(200, music)
}

// Downloads and converts playlist videos to audio only
func DownloadPlaylist(ctx *gin.Context) {
	var urlRequest UrlRequest

	ctx.ShouldBindBodyWithJSON(&urlRequest)

	go downloadPlaylist(urlRequest.Url)

	ctx.String(200, "Started downloading playlist...")
}

// Exports playlist data to file
// Does not download and convert video
func DryRun(ctx *gin.Context) {
	var urlRequest UrlRequest

	ctx.ShouldBindBodyWithJSON(&urlRequest)

	go dryRun(urlRequest.Url)

	ctx.String(200, "Doing a dry run!")
}

func dryRun(link string) {
	dl := ytdlp.New().
		SkipDownload().
		ForceIPv4().
		SleepInterval(5).
		MaxSleepInterval(20).
		PrintToFile("[%(playlist)s] %(webpage_url)s %(title)s", "playlist_info").
		Cookies("selenium/cookies_netscape")

	result, err := dl.Run(context.TODO(), link)

	if err != nil {
		logging.Log(err.Error())
	}

	logging.Log(result)
}

func downloadPlaylist(playlistUrl string) {
	dl := ytdlp.New().
		FormatSort("bestaudio").
		ExtractAudio().
		AudioFormat("opus").
		PostProcessorArgs("FFmpegExtractAudio:-b:a 160k").
		DownloadArchive("video_archive.db").
		EmbedMetadata().
		EmbedThumbnail().
		ForceIPv4().
		NoKeepVideo().
		Output("music/%(playlist_title)s/%(playlist_index)02d - %(title)s.%(ext)s").
		SleepInterval(8).MaxSleepInterval(20).
		Cookies("selenium/cookies_netscape")

	result, errr := dl.Run(context.TODO(), playlistUrl)

	if errr != nil {
		logging.Log("#yt-dlp Failed")
		logging.Log(errr.Error())
		return
	}

	logging.Log("#stdout")
	logging.Log(result.Stdout)

	for _, log := range result.OutputLogs {
		logging.Log(log.Line)
	}
}
