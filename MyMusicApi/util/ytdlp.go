package util

import (
	"api/db"
	"api/logging"
	"api/models"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

func StartDownloadTask(taskId int, downloadRequest models.DownloadRequestModel) {
	db := db.PostgresDb{}

	if db.OpenConnection() {

		defer db.CloseConnection()

		storageFolderName := Config.SourceFolder
		archiveFileName := fmt.Sprintf("%s/video_archive", storageFolderName)
		idsFileName := fmt.Sprintf("%s/ids.%d", storageFolderName, taskId)
		namesFileName := fmt.Sprintf("%s/names.%d", storageFolderName, taskId)
		durationFileName := fmt.Sprintf("%s/durations.%d", storageFolderName, taskId)
		playlistTitleFileName := fmt.Sprintf("%s/playlist_title.%d", storageFolderName, taskId)
		playlistIdFileName := fmt.Sprintf("%s/playlist_id.%d", storageFolderName, taskId)
		fileExtension := Config.OutputExtension

		cleanupFile := []string{
			idsFileName,
			namesFileName,
			durationFileName,
			playlistTitleFileName,
			playlistIdFileName,
		}

		isPlaylist := strings.Contains(downloadRequest.Url, "playlist?")

		dlp := ytdlp.New().
			FormatSort("bestaudio").
			ExtractAudio().
			AudioFormat(fileExtension).
			PostProcessorArgs("FFmpegExtractAudio:-b:a 160k").
			DownloadArchive(archiveFileName).
			WriteThumbnail().
			ForceIPv4().
			NoKeepVideo().
			PrintToFile("%(id)s", idsFileName).
			PrintToFile("%(title)s", namesFileName).
			PrintToFile("%(duration)s", durationFileName).
			Output(storageFolderName + "/%(id)s.%(ext)s").
			SleepInterval(5).
			MaxSleepInterval(20).
			Cookies("selenium/cookies_netscape")

		if isPlaylist {
			dlp = dlp.PrintToFile("%(playlist_title)s", playlistTitleFileName)
			dlp = dlp.PrintToFile("%(playlist_id)s", playlistIdFileName)
		}

		// Update task status
		db.UpdateTaskLogStatus(taskId, int(models.Downloading))

		// Start download
		result, err := dlp.Run(context.Background(), downloadRequest.Url)

		if err != nil {
			query := `UPDATE TaskLog
		             SET Status = $1, OutputLog = $2, EndTime = $3
		             WHERE Id = $4`

			// Set Task state -> Error
			json, err := json.Marshal(result.OutputLogs)

			err = db.NonScalarQuery(query, int(models.Error), json, time.Now(), taskId)
			if err != nil {
				logging.Error(fmt.Sprintf("Failed to update tasklog: %s", err.Error()))
			}

			// Delete created files if any
			for _, path := range cleanupFile {
				os.Remove(path)
			}
			return
		}

		// Set task state -> Updating
		db.UpdateTaskLogStatus(taskId, int(models.Updating))

		// Read output files -> update song table
		ids, _ := ReadLines(idsFileName)
		names, _ := ReadLines(namesFileName)
		durations, _ := ReadLines(durationFileName)

		playlistId := -1
		if isPlaylist {
			// create new playlist
			name, _ := ReadLines(playlistTitleFileName)

			// Check if exists, if not then create
			existingPlaylists, _ := db.FetchPlaylists(context.Background())

			playlistExists := false

			for _, playlist := range existingPlaylists {
				if playlist.Name == name[0] {
					playlistExists = true
					playlistId = playlist.Id
					break
				}
			}

			if !playlistExists {
				desc := "Custom playlist"
				_playlistId, _ := ReadLines(playlistIdFileName)
				playlistId, err = db.InsertPlaylist(models.Playlist{
					Name:          name[0],
					Description:   &desc,
					ThumbnailPath: fmt.Sprintf("%s.jpg", _playlistId[0]),
					CreationDate:  time.Now(),
					IsPublic:      true,
					UpdatedAt:     time.Now(),
				})
			}
		}

		var song models.Song

		for id := range len(ids) {
			song.Name = names[id]
			song.Duration, _ = strconv.Atoi(durations[id])
			song.SourceId = ids[id]
			song.Path = fmt.Sprintf("%s/%s.%s", storageFolderName, ids[id], fileExtension)
			song.ThumbnailPath = fmt.Sprintf("%s.webp", ids[id])
			id, err := db.InsertSong(song)

			if err != nil {
				logging.Error(fmt.Sprintf("[StartDownloadTask] Failed to insert song (%d): %s", id, err.Error()))
			}

			if isPlaylist {
				// add to playlist
				db.InsertPlaylistSong(playlistId, id)
			}
		}

		// Set task state -> Done
		// Update task status
		query := `UPDATE TaskLog
		             SET Status = $1, OutputLog = $2, EndTime = $3
		             WHERE Id = $4`

		json, err := json.Marshal(result.OutputLogs)

		err = db.NonScalarQuery(query, int(models.Done), json, time.Now(), taskId)
		if err != nil {
			logging.Error(fmt.Sprintf("Failed to update tasklog: %s", err.Error()))
		}

		// Delete created files
		for _, path := range cleanupFile {
			os.Remove(path)
		}
	} else {
		logging.Error(fmt.Sprintf("[StartDownloadTask] Failed to open database connection: %s", db.Error.Error()))
	}
}
