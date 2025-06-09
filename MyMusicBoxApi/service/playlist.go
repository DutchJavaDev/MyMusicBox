package service

import (
	"context"
	"encoding/json"
	"fmt"
	"musicboxapi/database"
	"musicboxapi/logging"
	"musicboxapi/models"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

func downloadPlaylist(
	taskId int,
	storage string,
	archiveFileName string,
	idsFileName string,
	namesFileName string,
	durationFileName string,
	playlistTitleFileName string,
	playlistIdFileName string,
	imagesFolder string,
	fileExtension string,
) {
	/*
			Ignore
			[Deleted video]
		    [Deleted video]
		    [Private video]
	*/
	ids, _ := readLines(idsFileName)
	downloadIds := ids
	names, _ := readLines(namesFileName)
	durations, _ := readLines(durationFileName)
	downloadCount := len(downloadIds)
	playlistNames, _ := readLines(playlistTitleFileName)

	db := database.PostgresDb{}

	if !db.OpenConnection() {
		logging.Error(fmt.Sprintf("[downloadPlaylist] failed to open database connection: %s", db.Error.Error()))
		return
	}
	defer db.CloseConnection()

	// Check if exists, if not then create
	existingPlaylists, _ := db.FetchPlaylists(context.Background())

	playlistExists := false
	playlistId := -1
	downloadCounter := 0

	for _, playlist := range existingPlaylists {
		if playlist.Name == playlistNames[0] {
			playlistExists = true
			playlistId = playlist.Id
			break
		}
	}

	if !playlistExists {
		desc := "Custom playlist"
		_playlistId, _ := readLines(playlistIdFileName)
		playlistId, _ = db.InsertPlaylist(models.Playlist{
			Name:          playlistNames[0],
			Description:   &desc,
			ThumbnailPath: fmt.Sprintf("%s.jpg", _playlistId[0]),
			CreationDate:  time.Now(),
			IsPublic:      true,
			UpdatedAt:     time.Now(),
		})
	}

	// Update task status
	db.UpdateTaskLogStatus(taskId, int(models.Downloading))

	defaultSettings := ytdlp.New().
		ExtractAudio().
		AudioQuality("0").
		AudioFormat(fileExtension).
		PostProcessorArgs("FFmpegExtractAudio:-b:a 160k").
		DownloadArchive(archiveFileName).
		WriteThumbnail().
		ConcurrentFragments(10).
		ConvertThumbnails("jpg").
		ForceIPv4().
		//sudo apt install aria2
		Downloader("aria2c").
		DownloaderArgs("aria2c:-x 16 -s 16 -j 16").
		NoKeepVideo().
		Output(storage + "/%(id)s.%(ext)s").
		Cookies("selenium/cookies_netscape")

	var outputLogs map[string]string

	outputLogs = make(map[string]string)

	hasError := false

	for id := range downloadCount {
		name := names[id]
		if canDownload(name) {
			downloadCounter++
			ytdlpInstance := defaultSettings.Clone()

			result, err := ytdlpInstance.Run(context.Background(), fmt.Sprintf("https://www.youtube.com/watch?v=%s", ids[id]))

			outputLogs[ids[id]] = result.Stdout

			if err != nil {
				hasError = true
				logging.Error(fmt.Sprintf("Failed to download %s, error:%s", ids[id], err.Error()))
				continue
			}

			var song models.Song

			song.Name = names[id]
			song.SourceId = ids[id]
			song.Duration, _ = strconv.Atoi(durations[id])
			song.Path = fmt.Sprintf("%s/%s.%s", storage, ids[id], fileExtension)
			song.ThumbnailPath = fmt.Sprintf("%s.jpg", ids[id])

			db.InsertSong(&song)

			db.InsertPlaylistSong(playlistId, song.Id)

			oldpath := fmt.Sprintf("%s/%s", storage, song.ThumbnailPath)
			newpath := fmt.Sprintf("%s/%s", imagesFolder, song.ThumbnailPath)

			_ = os.Rename(oldpath, newpath)
		}
	}

	json, err := json.Marshal(outputLogs)

	status := models.Done

	if hasError {
		status = models.Error
	}

	err = db.EndTaskLog(taskId, int(status), json)
	if err != nil {
		logging.Error(fmt.Sprintf("Failed to update tasklog: %s", err.Error()))
	}
}

func canDownload(name string) bool {

	if strings.HasPrefix(name, "[Deleted video]") {
		return false
	}

	if strings.HasPrefix(name, "[Private video]") {
		return false
	}

	if strings.HasPrefix(name, " ") {
		return false
	}

	return true
}
