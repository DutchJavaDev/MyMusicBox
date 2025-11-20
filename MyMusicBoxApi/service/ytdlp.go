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

	"github.com/lrstanley/go-ytdlp"
)

func StartDownloadTask(downloadRequest models.DownloadRequestModel) {

	tasklogTable := database.NewTasklogTableInstance()
	songTable := database.NewSongTableInstance()

	parentTask, _ := tasklogTable.CreateParentTaskLog(downloadRequest.Url)

	config := configuration.Config
	storageFolderName := config.SourceFolder
	archiveFileName := fmt.Sprintf("%s/video_archive", storageFolderName) // move to env / congif
	idsFileName := fmt.Sprintf("%s/ids.%d", storageFolderName, parentTask.Id)
	namesFileName := fmt.Sprintf("%s/names.%d", storageFolderName, parentTask.Id)
	durationFileName := fmt.Sprintf("%s/durations.%d", storageFolderName, parentTask.Id)
	playlistTitleFileName := fmt.Sprintf("%s/playlist_title.%d", storageFolderName, parentTask.Id)
	playlistIdFileName := fmt.Sprintf("%s/playlist_id.%d", storageFolderName, parentTask.Id)
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
			Cookies("selenium/cookies_netscape").
			IgnoreErrors()

		// Start download (flat download)
		result, err := dlp.Run(context.Background(), downloadRequest.Url)

		if err != nil {
			// Set Task state -> Error
			json, err := json.Marshal(result.OutputLogs)

			errChildTask, err := tasklogTable.CreateChildTaskLog(parentTask)

			errChildTask.OutputLog = json

			err = tasklogTable.ChildTaskLogError(errChildTask)

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

			json, _ := json.Marshal(models.ErrorResponse("Ytdlp files were not donwloaded, stopping task here"))

			errChildTask, err := tasklogTable.CreateChildTaskLog(parentTask)

			errChildTask.OutputLog = json

			err = tasklogTable.ChildTaskLogError(errChildTask)

			if err != nil {
				logging.Error(fmt.Sprintf("Failed to update tasklog: %s", err.Error()))
			}

			return
		}

		downloadPlaylist(
			parentTask,
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

		childTask, err := tasklogTable.CreateChildTaskLog(parentTask)

		if err != nil {
			return
		}

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
		childTask.Status = int(models.Downloading)

		err = tasklogTable.UpdateChildTaskLogStatus(childTask)

		if err != nil {
			return
		}

		// Start download
		result, err := dlp.Run(context.Background(), downloadRequest.Url)

		if err != nil {
			// Set Task state -> Error
			json, err := json.Marshal(result.OutputLogs)

			childTask.OutputLog = json

			err = tasklogTable.ChildTaskLogError(childTask)

			if err != nil {
				logging.Error(fmt.Sprintf("Failed to update tasklog: %s", err.Error()))
			}

			return
		}

		// Check if the files have been downloaded, if not stop going further to prevent panic
		if !fileExists(idsFileName) ||
			!fileExists(namesFileName) ||
			!fileExists(durationFileName) {

			json, _ := json.Marshal(models.ErrorResponse("Ytdlp files were not donwloaded, stopping task here"))

			childTask.OutputLog = json

			err = tasklogTable.ChildTaskLogError(childTask)

			if err != nil {
				logging.Error(fmt.Sprintf("Failed to update tasklog: %s", err.Error()))
			}

			return
		}

		// Update task status
		childTask.Status = int(models.Updating)

		err = tasklogTable.UpdateChildTaskLogStatus(childTask)

		if err != nil {
			return
		}

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
			// song.id might be not set..... :) ?
			logging.Error(fmt.Sprintf("[StartDownloadTask] Failed to insert song (%d): %s", song.Id, err.Error()))
		}

		oldpath := fmt.Sprintf("%s/%s", storageFolderName, song.ThumbnailPath)
		newpath := fmt.Sprintf("%s/%s", imagesFolder, song.ThumbnailPath)

		err = os.Rename(oldpath, newpath)

		if err != nil {
			logging.Error(fmt.Sprintf("Failed to move song image to / images folder: %s", err.Error()))
		}

		json, err := json.Marshal(result.OutputLogs)

		childTask.OutputLog = json
		childTask.Status = int(models.Done)
		err = tasklogTable.ChildTaskLogDone(childTask)

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
