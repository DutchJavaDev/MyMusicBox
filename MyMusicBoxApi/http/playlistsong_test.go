package http

import (
	"context"
	"encoding/json"
	"musicboxapi/database"
	"musicboxapi/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockPlaylistSongTable struct {
	database.IPlaylistsongTable
	fetchPlaylistSongs     func(ctx context.Context, playlistId int, lastKnowPosition int) (songs []models.Song, error error)
	insertPlaylistSong     func(playlistId int, songId int) (lastInsertedId int, error error)
	deleteAllPlaylistSongs func(playlistId int) (error error) // TODO
	deletePlaylistSong     func(playlistId int, songId int) (error error)
}

func (h *mockPlaylistSongTable) FetchPlaylistSongs(ctx context.Context, playlistId int, lastKnowPosition int) (songs []models.Song, error error) {
	return h.fetchPlaylistSongs(ctx, playlistId, lastKnowPosition)
}
func (h *mockPlaylistSongTable) InsertPlaylistSong(playlistId int, songId int) (lastInsertedId int, error error) {
	return h.insertPlaylistSong(playlistId, songId)
}

// TODO
func (h *mockPlaylistSongTable) DeleteAllPlaylistSongs(playlistId int) (error error) {
	return h.deleteAllPlaylistSongs(playlistId)
}
func (h *mockPlaylistSongTable) DeletePlaylistSong(playlistId int, songId int) (error error) {
	return h.deletePlaylistSong(playlistId, songId)
}

func TestFetchPlaylistSongsError(t *testing.T) {
	// Arrange
	route := "/playlist/:playlistId"
	router := SetupTestRouter()

	mockTable := &mockPlaylistSongTable{
		fetchPlaylistSongs: func(ctx context.Context, playlistId, lastKnowPosition int) (songs []models.Song, error error) {
			return nil, nil
		},
	}

	playlistSongHandler := PlaylistSongHandler{
		PlaylistsongTable: mockTable,
	}

	router.GET(route, playlistSongHandler.FetchPlaylistSongs)

	_route := "/playlist/sdfgsd"

	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", _route, nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestFetchPlaylistSongs(t *testing.T) {
	// Arrange
	route := "/playlist/:playlistId"
	router := SetupTestRouter()

	mockTable := &mockPlaylistSongTable{
		fetchPlaylistSongs: func(ctx context.Context, playlistId, lastKnowPosition int) (songs []models.Song, error error) {
			return []models.Song{
				{Id: 1, Name: "Chris Brown - Forever", SourceId: "fsfsdfsd", ThumbnailPath: "images/image.png", Path: "path/file.opus"},
			}, nil
		},
	}

	playlistSongHandler := PlaylistSongHandler{
		PlaylistsongTable: mockTable,
	}

	router.GET(route, playlistSongHandler.FetchPlaylistSongs)

	_route := "/playlist/1"

	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", _route, nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)

	var rawResult map[string]any

	err := json.Unmarshal(recorder.Body.Bytes(), &rawResult)

	assert.Equal(t, nil, err)

	dataBytes, err := json.Marshal(rawResult["Data"])

	var songs []models.Song
	err = json.Unmarshal(dataBytes, &songs)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(songs))
	assert.Equal(t, "Chris Brown - Forever", songs[0].Name)
}

func TestInsertPlaylistSongError(t *testing.T) {
	// Arrange
	route := "/playlistsong/:playlistId/:songId"
	router := SetupTestRouter()

	mockTable := &mockPlaylistSongTable{
		insertPlaylistSong: func(playlistId int, songId int) (lastInsertedId int, error error) {
			return -1, nil
		},
	}

	playlistSongHandler := PlaylistSongHandler{
		PlaylistsongTable: mockTable,
	}

	router.POST(route, playlistSongHandler.InsertPlaylistSong)

	_route := "/playlistsong/1/asdasd"

	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", _route, nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestInsertPlaylistSong(t *testing.T) {
	// Arrange
	route := "/playlistsong/:playlistId/:songId"
	router := SetupTestRouter()

	mockTable := &mockPlaylistSongTable{
		insertPlaylistSong: func(playlistId int, songId int) (lastInsertedId int, error error) {
			return 1, nil
		},
	}

	playlistSongHandler := PlaylistSongHandler{
		PlaylistsongTable: mockTable,
	}

	router.POST(route, playlistSongHandler.InsertPlaylistSong)

	_route := "/playlistsong/1/1"

	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", _route, nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)

	var rawResult map[string]any

	err := json.Unmarshal(recorder.Body.Bytes(), &rawResult)

	assert.Equal(t, nil, err)

	dataBytes, err := json.Marshal(rawResult["Data"])

	err = json.Unmarshal(dataBytes, &rawResult)

	assert.Equal(t, 1, int(rawResult["playlistSongId"].(float64)))
}

func TestDeleteAllPlaylistSongsError(t *testing.T) {
	// Arrange
	route := "playlistsong/:playlistId/:songId"
	router := SetupTestRouter()

	mockTable := &mockPlaylistSongTable{
		deletePlaylistSong: func(playlistId int, songId int) (error error) {
			return nil
		},
	}

	playlistSongHandler := PlaylistSongHandler{
		PlaylistsongTable: mockTable,
	}

	router.DELETE(route, playlistSongHandler.DeletePlaylistSong)

	_route := "/playlistsong/werwe/1"

	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", _route, nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestDeleteAllPlaylistSongs(t *testing.T) {
	// Arrange
	route := "playlistsong/:playlistId/:songId"
	router := SetupTestRouter()

	mockTable := &mockPlaylistSongTable{
		deletePlaylistSong: func(playlistId int, songId int) (error error) {
			return nil
		},
	}

	playlistSongHandler := PlaylistSongHandler{
		PlaylistsongTable: mockTable,
	}

	router.DELETE(route, playlistSongHandler.DeletePlaylistSong)

	_route := "/playlistsong/3/1"

	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", _route, nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
}
