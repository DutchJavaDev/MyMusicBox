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
)

func downloadPlaylist(
	parentTask models.ParentTaskLog,
	storage string,
	archiveFileName string,
	idsFileName string,
	namesFileName string,
	durationFileName string,
	playlistTitleFileName string,
	playlistIdFileName string,
	imagesFolder string,
	fileExtension string,
	logfilesOutputPath string,
	logfilesOutputPathError string,
	storageFolderName string,
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

	tasklogTable := database.NewTasklogTableInstance()
	playlistTable := database.NewPlaylistTableInstance()
	playlistsongTable := database.NewPlaylistsongTableInstance()
	songTable := database.NewSongTableInstance()

	// Check if exists, if not then create
	existingPlaylists, _ := playlistTable.FetchPlaylists(context.Background(), 0)

	playlistExists := false
	playlistId := -1

	for _, playlist := range existingPlaylists {
		if playlist.Name == playlistNames[0] {
			playlistExists = true
			playlistId = playlist.Id
			break
		}
	}

	_playlistId, _ := readLines(playlistIdFileName)

	if !playlistExists {
		_newPlaylistId, err := playlistTable.InsertPlaylist(models.Playlist{
			Name:          playlistNames[0],
			Description:   "Custom playlist",
			ThumbnailPath: fmt.Sprintf("%s.jpg", _playlistId[0]),
			CreationDate:  time.Now(),
			IsPublic:      true,
			UpdatedAt:     time.Now(),
		})

		if err != nil {

			errorChildLog, err := tasklogTable.CreateChildTaskLog(parentTask)

			json, err := json.Marshal(err)

			errorChildLog.OutputLog = json

			tasklogTable.ChildTaskLogError(errorChildLog)

			logging.Error(fmt.Sprintf("[Creating custom playlist error]: %s", err.Error()))
			return
		}

		playlistId = _newPlaylistId
	}

	// Special case, thumbnail is written to root directory
	if len(playlistNames) > 0 && len(_playlistId) > 0 {
		oldpath := fmt.Sprintf("%s [%s].jpg", playlistNames[0], _playlistId[0])
		newpath := fmt.Sprintf("%s/%s.jpg", imagesFolder, _playlistId[0])
		_ = os.Rename(oldpath, newpath)
	}

	// defaultSettings := ytdlp.New().
	// 	ExtractAudio().
	// 	AudioQuality("0").
	// 	AudioFormat(fileExtension).
	// 	DownloadArchive(archiveFileName).
	// 	WriteThumbnail().
	// 	ConcurrentFragments(10).
	// 	ConvertThumbnails("jpg").
	// 	ForceIPv4().
	// 	//sudo apt install aria2
	// 	Downloader("aria2c").
	// 	DownloaderArgs("aria2c:-x 16 -s 16 -j 16").
	// 	NoKeepVideo().
	// 	Output(storage + "/%(id)s.%(ext)s").
	// 	Cookies("selenium/cookies_netscape")

	for id := range downloadCount {
		name := names[id]
		if canDownload(name) && !existInArchive(archiveFileName, ids[id]) {

			childTask, _ := tasklogTable.CreateChildTaskLog(parentTask)

			childTask.Status = int(models.Downloading)

			tasklogTable.UpdateChildTaskLogStatus(childTask)

			downloaded := FlatSingleDownload(
				archiveFileName,
				idsFileName,
				namesFileName,
				durationFileName,
				playlistTitleFileName,
				playlistIdFileName,
				fmt.Sprintf("https://www.youtube.com/watch?v=%s", ids[id]),
				logfilesOutputPath,
				logfilesOutputPathError,
				storageFolderName,
				fileExtension)

			if !downloaded {
				file, err := os.ReadFile(logfilesOutputPathError)
				json, _ := json.Marshal(file)
				childTask.OutputLog = json
				tasklogTable.ChildTaskLogError(childTask)
				logging.Error(fmt.Sprintf("Failed to download %s, error:%s", ids[id], err.Error()))
				continue
			}

			childTask.Status = int(models.Updating)

			tasklogTable.UpdateChildTaskLogStatus(childTask)

			var song models.Song

			song.Name = names[id]
			song.SourceId = ids[id]
			song.Duration, _ = strconv.Atoi(durations[id])
			song.Path = fmt.Sprintf("%s/%s.%s", storage, ids[id], fileExtension)
			song.ThumbnailPath = fmt.Sprintf("%s.jpg", ids[id])

			songTable.InsertSong(&song)

			playlistsongTable.InsertPlaylistSong(playlistId, song.Id)

			oldpath := fmt.Sprintf("%s/%s", storage, song.ThumbnailPath)
			newpath := fmt.Sprintf("%s/%s", imagesFolder, song.ThumbnailPath)

			_ = os.Rename(oldpath, newpath)

			file, err := os.ReadFile(logfilesOutputPath)

			if err != nil {
				panic(-1564654654)
			}

			json, _ := json.Marshal(file)

			childTask.OutputLog = json

			tasklogTable.ChildTaskLogDone(childTask)
		}
	}
}

func existInArchive(path string, id string) bool {
	lines, err := readLines(path)

	if err != nil {
		panic(-654654)
	}

	for line := range lines {
		if strings.Contains(lines[line], fmt.Sprintf("youtube %s", id)) {
			return true
		}
	}

	return false
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
