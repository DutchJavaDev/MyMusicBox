package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"musicboxapi/configuration"
	"musicboxapi/database"
	"musicboxapi/logging"
	"musicboxapi/models"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

func StartDownloadTask(taskId int, downloadRequest models.DownloadRequestModel) {

	config := configuration.Config
	storageFolderName := config.SourceFolder
	archiveFileName := fmt.Sprintf("%s/video_archive", storageFolderName)
	idsFileName := fmt.Sprintf("%s/ids.%d", storageFolderName, taskId)
	namesFileName := fmt.Sprintf("%s/names.%d", storageFolderName, taskId)
	durationFileName := fmt.Sprintf("%s/durations.%d", storageFolderName, taskId)
	playlistTitleFileName := fmt.Sprintf("%s/playlist_title.%d", storageFolderName, taskId)
	playlistIdFileName := fmt.Sprintf("%s/playlist_id.%d", storageFolderName, taskId)
	imagesFolder := fmt.Sprintf("%s/images", storageFolderName)
	fileExtension := config.OutputExtension

	if !pathExists(storageFolderName) {
		err := os.Mkdir(storageFolderName, fs.ModePerm|fs.ModeDir)
		if err != nil {
			logging.Error(err.Error())
			return
		}
	}

	if !pathExists(imagesFolder) {
		err := os.Mkdir(imagesFolder, fs.ModePerm|fs.ModeDir)
		if err != nil {
			logging.Error(err.Error())
			return
		}
	}

	isPlaylist := strings.Contains(downloadRequest.Url, "playlist?")

	cleanupFile := []string{
		idsFileName,
		namesFileName,
		durationFileName,
		playlistTitleFileName,
		playlistIdFileName,
	}

	tasklogTable := database.NewTasklogTableInstance()
	songTable := database.NewSongTableInstance()

	if isPlaylist {
		dlp := ytdlp.New().
			DownloadArchive(archiveFileName).
			ForceIPv4().
			NoKeepVideo().
			SkipDownload().
			FlatPlaylist().
			WriteThumbnail().
			PrintToFile("%(id)s", idsFileName).
			PrintToFile("%(title)s", namesFileName).
			PrintToFile("%(duration)s", durationFileName).
			PrintToFile("%(playlist_title)s", playlistTitleFileName).
			PrintToFile("%(playlist_id)s", playlistIdFileName).
			Cookies("selenium/cookies_netscape")

		// Start download (flat download)
		result, err := dlp.Run(context.Background(), downloadRequest.Url)

		if err != nil {
			// Set Task state -> Error
			json, err := json.Marshal(result.OutputLogs)

			err = tasklogTable.EndTaskLog(taskId, int(models.Error), json)
			if err != nil {
				logging.Error(fmt.Sprintf("Failed to update tasklog: %s", err.Error()))
			}

			//Delete created files if any
			for _, path := range cleanupFile {
				os.Remove(path)
			}
			return
		}

		// Check if the files have been downloaded, if not stop going further to prevent panic
		if !fileExists(idsFileName) ||
			!fileExists(namesFileName) ||
			!fileExists(durationFileName) ||
			!fileExists(playlistTitleFileName) ||
			!fileExists(playlistIdFileName) {

			errorJosn, _ := json.Marshal(models.ErrorResponse("Ytdlp files were not donwloaded, stopping task here"))

			// Set Task state -> Error
			err = tasklogTable.UpdateTaskLogError(int(models.Error), errorJosn, time.Now(), taskId)
			if err != nil {
				logging.Error(fmt.Sprintf("Failed to update tasklog: %s", err.Error()))
			}

			return
		}

		downloadPlaylist(
			taskId,
			storageFolderName,
			archiveFileName,
			idsFileName,
			namesFileName,
			durationFileName,
			playlistTitleFileName,
			playlistIdFileName,
			imagesFolder,
			fileExtension)

		// Delete created files
		for _, path := range cleanupFile {
			os.Remove(path)
		}
	} else {
		// Normal download
		dlp := ytdlp.New().
			ExtractAudio().
			AudioQuality("0").
			AudioFormat(fileExtension).
			DownloadArchive(archiveFileName).
			WriteThumbnail().
			ConcurrentFragments(10).
			ConvertThumbnails("jpg").
			ForceIPv4().
			PrintToFile("%(id)s", idsFileName).
			PrintToFile("%(title)s", namesFileName).
			PrintToFile("%(duration)s", durationFileName).
			//sudo apt install aria2
			Downloader("aria2c").
			DownloaderArgs("aria2c:-x 16 -s 16 -j 16").
			NoKeepVideo().
			Output(storageFolderName + "/%(id)s.%(ext)s").
			Cookies("selenium/cookies_netscape")

		// Update task status
		tasklogTable.UpdateTaskLogStatus(taskId, int(models.Downloading))

		// Start download
		result, err := dlp.Run(context.Background(), downloadRequest.Url)

		if err != nil {
			query := `UPDATE TaskLog
		             SET Status = $1, OutputLog = $2, EndTime = $3
		             WHERE Id = $4`

			// Set Task state -> Error
			outputlogJson, err := json.Marshal(result.OutputLogs)

			err = tasklogTable.NonScalarQuery(query, int(models.Error), outputlogJson, time.Now(), taskId)
			if err != nil {
				logging.Error(fmt.Sprintf("Failed to update tasklog: %s", err.Error()))
			}

			return
		}

		// Check if the files have been downloaded, if not stop going further to prevent panic
		if !fileExists(idsFileName) ||
			!fileExists(namesFileName) ||
			!fileExists(durationFileName) {

			errorJosn, _ := json.Marshal(models.ErrorResponse("Ytdlp files were not donwloaded, stopping task here"))

			// Set Task state -> Error
			err = tasklogTable.UpdateTaskLogError(int(models.Error), errorJosn, time.Now(), taskId)
			if err != nil {
				logging.Error(fmt.Sprintf("Failed to update tasklog: %s", err.Error()))
			}

			return
		}

		// Set task state -> Updating
		tasklogTable.UpdateTaskLogStatus(taskId, int(models.Updating))

		//Read output files -> update song table
		ids, _ := readLines(idsFileName)
		names, _ := readLines(namesFileName)
		durations, _ := readLines(durationFileName)

		var song models.Song

		indexId := 0

		song.Name = names[indexId]
		song.Duration, _ = strconv.Atoi(durations[indexId])
		song.SourceId = ids[indexId]
		song.Path = fmt.Sprintf("%s/%s.%s", storageFolderName, ids[indexId], fileExtension)
		song.ThumbnailPath = fmt.Sprintf("%s.jpg", ids[indexId])
		err = songTable.InsertSong(&song)

		if err != nil {
			// song.id might be not set..... :)
			logging.Error(fmt.Sprintf("[StartDownloadTask] Failed to insert song (%d): %s", song.Id, err.Error()))
		}

		oldpath := fmt.Sprintf("%s/%s", storageFolderName, song.ThumbnailPath)
		newpath := fmt.Sprintf("%s/%s", imagesFolder, song.ThumbnailPath)

		err = os.Rename(oldpath, newpath)

		if err != nil {
			logging.Error(fmt.Sprintf("Failed to move song image to /images folder: %s", err.Error()))
		}

		// Set task state -> Done
		// Update task status
		query := `UPDATE TaskLog
		             SET Status = $1, OutputLog = $2, EndTime = $3
		             WHERE Id = $4`
		outputJson, err := json.Marshal(result.OutputLogs)
		err = tasklogTable.NonScalarQuery(query, int(models.Done), outputJson, time.Now(), taskId)
		if err != nil {
			logging.Error(fmt.Sprintf("Failed to update tasklog: %s", err.Error()))
		}

		//Delete created files
		for _, path := range cleanupFile {
			os.Remove(path)
		}
	}
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
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

func pathExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil && (info.IsDir())
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	// No error, file exists
	return err == nil
}
