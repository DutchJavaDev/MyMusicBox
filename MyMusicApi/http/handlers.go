package http

import (
	"api/db"
	"api/logging"
	"api/models"
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lrstanley/go-ytdlp"
)

// begin fetch
func FetchSongs(ctx *gin.Context) {
	db := db.PostgresDb{}

	defer db.Close()

	if db.InitDatabase() {
		songs, err := db.FetchSongs(ctx.Request.Context())
		if err != nil {
			ctx.JSON(500, err.Error())
			return
		}
		ctx.JSON(200, songs)
	} else {
		ctx.JSON(500, db.Error)
	}
}

func FetchPlaylists(ctx *gin.Context) {
	db := db.PostgresDb{}
	defer db.Close()

	if db.InitDatabase() {
		songs, err := db.FetchPlaylists(ctx.Request.Context())
		if err != nil {
			ctx.JSON(500, err.Error())
			return
		}
		ctx.JSON(200, songs)
	} else {
		ctx.JSON(500, db.Error)
	}
}

func FetchPlaylistSongs(ctx *gin.Context) {
	db := db.PostgresDb{}
	defer db.Close()

	playlistId, _ := strconv.Atoi(ctx.Param("playlistId"))

	if db.InitDatabase() {
		songs, err := db.FetchPlaylistSongs(ctx.Request.Context(), playlistId)
		if err != nil {
			ctx.JSON(500, err.Error())
			return
		}
		ctx.JSON(200, songs)
	} else {
		ctx.JSON(500, db.Error)
	}
}

// end fetch

// begin insert
func InsertPlaylist(ctx *gin.Context) {

	var playlist models.Playlist

	ctx.ShouldBindBodyWithJSON(&playlist)

	db := db.PostgresDb{}
	defer db.Close()

	if db.InitDatabase() {
		id, err := db.InsertPlaylist(playlist)

		if err != nil {
			ctx.JSON(500, err.Error())
			return
		}

		ctx.JSON(200, gin.H{"playlistId": id})
	} else {
		ctx.JSON(500, db.Error)
	}
}

func InsertPlaylistSong(ctx *gin.Context) {
	db := db.PostgresDb{}
	defer db.Close()

	playlistId, _ := strconv.Atoi(ctx.Param("playlistId"))
	songId, _ := strconv.Atoi(ctx.Param("songId"))

	if db.InitDatabase() {
		id, err := db.InsertPlaylistSong(playlistId, songId)

		if err != nil {
			ctx.JSON(500, err.Error())
			return
		}
		ctx.JSON(200, gin.H{"playlistSongId": id})
	} else {
		ctx.JSON(500, db.Error)
	}
}

// end insert

// begin delete
func DeletePlaylist(ctx *gin.Context) {
	db := db.PostgresDb{}
	defer db.Close()

	id, _ := strconv.Atoi(ctx.Param("playlistId"))

	if db.InitDatabase() {
		err := db.DeletePlaylistById(id)
		if err != nil {
			ctx.JSON(500, err.Error())
			return
		}
		ctx.Status(200)
	}
}

func DeletePlaylistSong(ctx *gin.Context) {
	db := db.PostgresDb{}
	defer db.Close()

	playlistId, _ := strconv.Atoi(ctx.Param("playlistId"))
	songId, _ := strconv.Atoi(ctx.Param("songId"))

	if db.InitDatabase() {
		err := db.DeletePlaylistSong(playlistId, songId)
		if err != nil {
			ctx.JSON(500, err.Error())
			return
		}
		ctx.Status(200)
	}
}

// end delete

// Downloads and converts playlist videos to audio only
func DownloadPlaylist(ctx *gin.Context) {
	var urlRequest models.UrlRequest

	ctx.ShouldBindBodyWithJSON(&urlRequest)

	go downloadPlaylist(urlRequest.Url)

	ctx.String(200, "Started downloading playlist...")
}

// Exports playlist data to file
// Does not download and convert video
func DryRun(ctx *gin.Context) {
	var urlRequest models.UrlRequest

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
		logging.Error(err.Error())
	}

	logging.Info(result)
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
		WriteThumbnail().
		ForceIPv4().
		NoKeepVideo().
		Output("music/%(playlist_title)s/%(playlist_index)02d - %(title)s.%(ext)s").
		SleepInterval(8).MaxSleepInterval(20).
		Cookies("selenium/cookies_netscape")

	result, errr := dl.Run(context.TODO(), playlistUrl)

	if errr != nil {
		logging.Error("#yt-dlp Failed")
		logging.Error(errr.Error())
		return
	}

	logging.Info("#stdout")
	logging.Info(result.Stdout)

	for _, log := range result.OutputLogs {
		logging.Info(log.Line)
	}
}

// Ignore thise endpoint, testing purpose
func AddSong(ctx *gin.Context) {

	var song models.Song
	num := rand.Int()
	path := fmt.Sprintf("/lol/path %d", num)
	song.Name = fmt.Sprintf("Juice WRLD - Vampire %d", num)
	song.Path = &path
	song.SourceURL = fmt.Sprintf("https://www.youtube.com/watch?v=0G5a6Tm_pQQQ %d", num)

	db := db.PostgresDb{}
	defer db.Close()

	if db.InitDatabase() {
		id, err := db.InsertSong(song)

		if err != nil {
			ctx.JSON(500, err.Error())
			return
		}

		ctx.JSON(200, gin.H{"songId": id})
	} else {
		ctx.JSON(500, db.Error)
	}
}
