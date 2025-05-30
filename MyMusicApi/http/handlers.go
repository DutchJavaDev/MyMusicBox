package http

import (
	"api/db"
	"api/models"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FetchSongs(ctx *gin.Context) {
	db := db.PostgresDb{}

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

func FetchPlaylists(ctx *gin.Context) {
	db := db.PostgresDb{}
	defer db.CloseConnection()

	if db.OpenConnection() {
		playlists, err := db.FetchPlaylists(ctx.Request.Context())
		if err != nil {
			ctx.JSON(500, models.ErrorResponse(err))
			return
		}
		ctx.JSON(200, models.OkResponse(playlists, fmt.Sprintf("Found %d playlist", len(playlists))))
	} else {
		ctx.JSON(500, models.ErrorResponse(db.Error))
	}
}

func FetchPlaylistSongs(ctx *gin.Context) {
	db := db.PostgresDb{}
	defer db.CloseConnection()

	playlistId, _ := strconv.Atoi(ctx.Param("playlistId"))

	if db.OpenConnection() {
		songs, err := db.FetchPlaylistSongs(ctx.Request.Context(), playlistId)
		if err != nil {
			ctx.JSON(500, models.ErrorResponse(err))
			return
		}
		ctx.JSON(200, models.OkResponse(songs, fmt.Sprintf("Found %d songs in playlist %d", len(songs), playlistId)))
	} else {
		ctx.JSON(500, models.ErrorResponse(db.Error))
	}
}

func InsertPlaylist(ctx *gin.Context) {

	var playlist models.Playlist

	ctx.ShouldBindBodyWithJSON(&playlist)

	db := db.PostgresDb{}
	defer db.CloseConnection()

	if db.OpenConnection() {
		id, err := db.InsertPlaylist(playlist)

		if err != nil {
			ctx.JSON(500, models.ErrorResponse(err))
			return
		}

		ctx.JSON(200, models.OkResponse(gin.H{"playlistId": id}, "Created new playlist"))
	} else {
		ctx.JSON(500, models.ErrorResponse(db.Error))
	}
}

func InsertPlaylistSong(ctx *gin.Context) {
	db := db.PostgresDb{}
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

func DeletePlaylist(ctx *gin.Context) {
	db := db.PostgresDb{}
	defer db.CloseConnection()

	id, _ := strconv.Atoi(ctx.Param("playlistId"))

	if db.OpenConnection() {
		err := db.DeletePlaylist(id)
		if err != nil {
			ctx.JSON(500, models.ErrorResponse(err))
			return
		}
		ctx.Status(200)
	} else {
		ctx.JSON(500, models.ErrorResponse(db.Error))
	}
}

func DeletePlaylistSong(ctx *gin.Context) {
	db := db.PostgresDb{}
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

func DownloadRequest(ctx *gin.Context) {
	var request models.DownloadRequest
	ctx.ShouldBindBodyWithJSON(&request)

}

// Downloads and converts playlist videos to audio only
func DownloadPlaylist(ctx *gin.Context) {
	var urlRequest models.UrlRequest

	ctx.ShouldBindBodyWithJSON(&urlRequest)

	go downloadPlaylist(urlRequest.Url)

	ctx.String(200, "Started downloading playlist...")
}

// Exports playlist data to file
// Does not download and convert video
// func DryRun(ctx *gin.Context) {
// 	var urlRequest models.UrlRequest

// 	ctx.ShouldBindBodyWithJSON(&urlRequest)

// 	go dryRun(urlRequest.Url)

// 	ctx.String(200, "Doing a dry run!")
// }

// func dryRun(link string) {
// 	dl := ytdlp.New().
// 		SkipDownload().
// 		ForceIPv4().
// 		SleepInterval(5).
// 		MaxSleepInterval(20).
// 		PrintToFile("[%(playlist)s] %(webpage_url)s %(title)s", "playlist_info").
// 		Cookies("selenium/cookies_netscape")

// 	result, err := dl.Run(context.TODO(), link)

// 	if err != nil {
// 		Error(err.Error())
// 	}

// 	Info(result)
// }

func downloadPlaylist(playlistUrl string) {
	// dl := ytdlp.New().
	// 	FormatSort("bestaudio").
	// 	ExtractAudio().
	// 	AudioFormat("opus").
	// 	PostProcessorArgs("FFmpegExtractAudio:-b:a 160k").
	// 	DownloadArchive("video_archive.db").
	// 	EmbedMetadata().
	// 	EmbedThumbnail().
	// 	WriteThumbnail().
	// 	ForceIPv4().
	// 	NoKeepVideo().
	// 	Output("music/%(playlist_title)s/%(playlist_index)02d - %(title)s.%(ext)s").
	// 	SleepInterval(8).MaxSleepInterval(20).
	// 	Cookies("selenium/cookies_netscape")

	// result, errr := dl.Run(context.TODO(), playlistUrl)

	// if errr != nil {
	// 	//logging.Error("#yt-dlp Failed")
	// 	//logging.Error(errr.Error())
	// 	return
	// }

	//logging.Info("#stdout")
	//logging.Info(result.Stdout)

	// for _, log := range result.OutputLogs {
	// 	//logging.Info(log.Line)
	// }
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
	defer db.CloseConnection()

	if db.OpenConnection() {
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
