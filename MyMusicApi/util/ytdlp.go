package util

import (
	"api/db"
	"api/logging"
	"api/models"
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

func StartDownloadTask(taskId int, downloadRequest models.DownloadRequest) {
	db := db.PostgresDb{}

	if db.OpenConnection() {

		defer db.CloseConnection()

		storageFolderName := "music_dev"
		archiveFileName := fmt.Sprintf("%s/video_archive", storageFolderName)
		idsFileName := fmt.Sprintf("%s/ids.%d", storageFolderName, taskId)
		namesFileName := fmt.Sprintf("%s/names.%d", storageFolderName, taskId)
		durationFileName := fmt.Sprintf("%s/durations.%d", storageFolderName, taskId)
		urlsFileName := fmt.Sprintf("%s/urls.%d", storageFolderName, taskId)
		fileExtension := "opus"

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
			PrintToFile("%(webpage_url)s", urlsFileName).
			Output(storageFolderName + "/%(id)s.%(ext)s").
			SleepInterval(5).
			MaxSleepInterval(20).
			Cookies("selenium/cookies_netscape")

			/*
				n_entries (numeric): Total number of extracted items in the playlist
				playlist_id (string): Identifier of the playlist that contains the video'
				for playlist the image extension is jpg
			*/

		// Update task status
		db.UpdateTaskLogStatus(taskId, int(models.Downloading))

		// Start download
		result, err := dlp.Run(context.Background(), downloadRequest.Url)

		if err != nil {
			query := `UPDATE TaskLog
		             SET Status = $1, OutputLog = $2, EndTime = $3
		             WHERE Id = $4`
			// Store output
			// Set Task state -> Error
			err := db.NonScalarQuery(query, int(models.Error), result.Stderr, time.Now(), taskId)

			if err != nil {
				logging.Error(fmt.Sprintf("Failed to update tasklog: %s", err.Error()))
			}

			// Delete created files?
			return
		}

		// Set task state -> Updating
		db.UpdateTaskLogStatus(taskId, int(models.Updating))

		// Read output files -> update song table
		ids, _ := ReadLines(idsFileName)
		names, _ := ReadLines(namesFileName)
		durations, _ := ReadLines(durationFileName)
		urls, _ := ReadLines(urlsFileName)

		var song models.Song

		for id := range len(ids) {
			song.Name = names[id]
			song.Duration, _ = strconv.Atoi(durations[id])
			song.SourceURL = urls[id]

			path := fmt.Sprintf("%s/%s.%s", storageFolderName, ids[id], fileExtension)
			song.Path = &path

			id, err := db.InsertSong(song)

			if err != nil {
				logging.Error(fmt.Sprintf("[StartDownloadTask] Failed to insert song (%d): %s", id, err.Error()))
			}
		}

		// Set task state -> Done
		// Update task status
		db.UpdateTaskLogStatus(taskId, int(models.Done))

		// Delete created files
		os.Remove(idsFileName)
		os.Remove(namesFileName)
		os.Remove(durationFileName)
		os.Remove(urlsFileName)
	} else {
		logging.Error(fmt.Sprintf("[StartDownloadTask] Failed to open database connection: %s", db.Error.Error()))
	}
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
