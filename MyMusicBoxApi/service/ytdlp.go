package service

import (
	"bufio"
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
	logfilesOutputPath := fmt.Sprintf("%s/hotfix_logs/%s", storageFolderName, fmt.Sprintf("logrun_%d", parentTask.Id))
	logfilesOutputPathError := fmt.Sprintf("%s/hotfix_logs/%s", storageFolderName, fmt.Sprintf("logrunError_%d", parentTask.Id))

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
		logfilesOutputPath,
		logfilesOutputPathError,
	}

	if isPlaylist {
		downloaded := FlatPlaylistDownload(
			archiveFileName,
			idsFileName,
			namesFileName,
			durationFileName,
			playlistTitleFileName,
			playlistIdFileName,
			downloadRequest.Url,
			logfilesOutputPath,
			logfilesOutputPathError)

		// Start download (flat download)
		if !downloaded {

			fmt.Printf("Failed to download for %s check logs at %s", downloadRequest.Url, logfilesOutputPathError)

			file, err := os.ReadFile(logfilesOutputPathError)

			// Set Task state -> Error
			json, err := json.Marshal(string(file))

			errChildTask, err := tasklogTable.CreateChildTaskLog(parentTask)

			errChildTask.OutputLog = json

			err = tasklogTable.ChildTaskLogError(errChildTask)

			if err != nil {
				logging.Error(fmt.Sprintf("Failed to update tasklog: %s", err.Error()))
			}

			//Delete created files if any
			for _, path := range cleanupFile {

				if strings.Contains(path, "logrunError") {
					lines, err := readLines(path)

					if err != nil {
						panic(-654654654)
					}

					if len(lines) > 0 && len(lines[0]) > 0 {
						// skip deleting log so it can be used for debug
						continue
					}
				}
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
			fileExtension,
			logfilesOutputPath,
			logfilesOutputPathError,
			storageFolderName)

		// Delete created files
		for _, path := range cleanupFile {

			if strings.Contains(path, "logrunError") {
				lines, err := readLines(path)

				if err != nil {
					panic(-654654654)
				}

				if len(lines) > 0 && len(lines[0]) > 0 {
					// skip deleting log so it can be used for debug
					continue
				}
			}

			os.Remove(path)
		}
	} else {
		// Normal download

		childTask, err := tasklogTable.CreateChildTaskLog(parentTask)

		if err != nil {
			return
		}

		downloaded := FlatSingleDownload(
			archiveFileName,
			idsFileName,
			namesFileName,
			durationFileName,
			playlistTitleFileName,
			playlistIdFileName,
			downloadRequest.Url,
			logfilesOutputPath,
			logfilesOutputPathError,
			storageFolderName,
			fileExtension)

		childTask.Status = int(models.Downloading)

		err = tasklogTable.UpdateChildTaskLogStatus(childTask)

		if err != nil {
			return
		}

		if !downloaded {
			fmt.Printf("Failed to download for %s check logs at %s", downloadRequest.Url, logfilesOutputPathError)

			file, err := os.ReadFile(logfilesOutputPathError)

			// Set Task state -> Error
			json, err := json.Marshal(string(file))

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

		file, err := os.ReadFile(logfilesOutputPath)

		json, err := json.Marshal(string(file))

		childTask.OutputLog = json
		childTask.Status = int(models.Done)
		err = tasklogTable.ChildTaskLogDone(childTask)

		if err != nil {
			logging.Error(fmt.Sprintf("Failed to update tasklog: %s", err.Error()))
		}

		//Delete created files
		for _, path := range cleanupFile {

			if strings.Contains(path, "logrunError") {
				lines, err := readLines(path)

				if err != nil {
					panic(-654654654)
				}

				if len(lines) > 0 && len(lines[0]) > 0 {
					// skip deleting log so it can be used for debug
					continue
				}
			}
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
