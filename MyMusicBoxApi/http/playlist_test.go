package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"musicboxapi/database"
	"musicboxapi/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockPlaylistTable struct {
	database.IPlaylistTable
	fetchPlaylists func(ctx context.Context, lastKnowPlaylistId int) (playlists []models.Playlist, error error)
	insertPlaylist func(playlist models.Playlist) (lastInsertedId int, error error)
	deletePlaylist func(playlistId int) (error error)
}

func (m *mockPlaylistTable) FetchPlaylists(ctx context.Context, lastKnowPlaylistId int) (playlists []models.Playlist, error error) {
	return m.fetchPlaylists(ctx, lastKnowPlaylistId)
}
func (m *mockPlaylistTable) InsertPlaylist(playlist models.Playlist) (lastInsertedId int, error error) {
	return m.insertPlaylist(playlist)
}
func (m *mockPlaylistTable) DeletePlaylist(playlistId int) (error error) {
	return m.deletePlaylist(playlistId)
}

func TestFetchPlaylists(t *testing.T) {
	// Arrange
	route := "/api/v1/playlist"
	router := SetupTestRouter()

	mockTable := &mockPlaylistTable{
		fetchPlaylists: func(ctx context.Context, lastKnowPlaylistId int) ([]models.Playlist, error) {
			return []models.Playlist{
				{Id: 1, Name: "Playlist_1", Description: "Best ever", ThumbnailPath: "path/image1.png", IsPublic: false},
				{Id: 2, Name: "Playlist_2", Description: "Second best", ThumbnailPath: "path/image2.png", IsPublic: true},
			}, nil
		},
	}

	playlistHandler := PlaylistHandler{
		PlaylistTable: mockTable,
	}

	router.GET(route, playlistHandler.FetchPlaylists)

	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", route, nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)

	var rawResult map[string]any

	err := json.Unmarshal(recorder.Body.Bytes(), &rawResult)

	assert.Equal(t, nil, err)

	dataBytes, err := json.Marshal(rawResult["Data"])

	var playlists []models.Playlist
	err = json.Unmarshal(dataBytes, &playlists)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(playlists))
	assert.Equal(t, "Playlist_1", playlists[0].Name)
	assert.Equal(t, "Playlist_2", playlists[1].Name)
}

func TestFetchPlaylistsLastKnowPlaylistId(t *testing.T) {
	// Arrange
	_lastKnowPlaylistId := 4
	route := "/api/v1/playlist"
	router := SetupTestRouter()

	mockTable := &mockPlaylistTable{
		fetchPlaylists: func(ctx context.Context, lastKnowPlaylistId int) ([]models.Playlist, error) {
			assert.Equal(t, _lastKnowPlaylistId, lastKnowPlaylistId)
			return []models.Playlist{
				{Id: 4, Name: "Playlist_4", Description: "Second best", ThumbnailPath: "path/image4.png", IsPublic: true},
			}, nil
		},
	}

	playlistHandler := PlaylistHandler{
		PlaylistTable: mockTable,
	}

	router.GET(route, playlistHandler.FetchPlaylists)

	recorder := httptest.NewRecorder()

	queryRoute := fmt.Sprintf("%s?lastKnowPlaylistId=%d", route, _lastKnowPlaylistId)

	req, _ := http.NewRequest("GET", queryRoute, nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)

	var rawResult map[string]any

	err := json.Unmarshal(recorder.Body.Bytes(), &rawResult)

	assert.Equal(t, nil, err)

	dataBytes, err := json.Marshal(rawResult["Data"])

	var playlists []models.Playlist
	err = json.Unmarshal(dataBytes, &playlists)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(playlists))
	assert.Equal(t, "Playlist_4", playlists[0].Name)
}

func TestInsertPlaylist(t *testing.T) {
	// Arrange
	route := "/api/v1/playlist"
	router := SetupTestRouter()

	mockTable := &mockPlaylistTable{
		insertPlaylist: func(playlist models.Playlist) (lastInsertedId int, error error) {
			return 1, nil
		},
	}

	playlistHandler := PlaylistHandler{
		PlaylistTable: mockTable,
	}

	router.POST(route, playlistHandler.InsertPlaylist)

	recorder := httptest.NewRecorder()
	recorder.Header().Set("Content-Type", "multipart/form-data")

	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	w.WriteField("playlistDescription", "Cool playlist")
	w.WriteField("playlistName", "Banger Songs")
	w.WriteField("publicPlaylist", "on")

	imgWrite, err := w.CreateFormFile("backgroundImage", "default.png")

	io.Copy(imgWrite, bytes.NewReader([]byte("0")))

	w.Close()

	req, _ := http.NewRequest("POST", route, buf)

	req.Header.Add("Content-Type", w.FormDataContentType())

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)

	var rawResult map[string]any

	err = json.Unmarshal(recorder.Body.Bytes(), &rawResult)

	assert.Equal(t, nil, err)

	dataBytes, err := json.Marshal(rawResult["Data"])

	err = json.Unmarshal(dataBytes, &rawResult)

	assert.Equal(t, 1, int(rawResult["playlistId"].(float64)))
}

func TestInsertPlaylistJsonError(t *testing.T) {
	// Arrange
	route := "/api/v1/playlist"
	router := SetupTestRouter()

	mockTable := &mockPlaylistTable{
		insertPlaylist: func(playlist models.Playlist) (lastInsertedId int, error error) {
			return 1, nil
		},
	}

	playlistHandler := PlaylistHandler{
		PlaylistTable: mockTable,
	}

	router.POST(route, playlistHandler.InsertPlaylist)

	recorder := httptest.NewRecorder()

	// Wrong type, will throw an error
	bodyBytes, _ := json.Marshal("")

	req, _ := http.NewRequest("POST", route, bytes.NewBuffer(bodyBytes))

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestDeletePlaylistPlaylistIdError(t *testing.T) {
	// Arrange
	route := "/playlist/:playlistId"
	router := SetupTestRouter()

	mockTable := &mockPlaylistTable{
		deletePlaylist: func(playlistId int) (error error) {
			return nil
		},
	}

	playlistHandler := PlaylistHandler{
		PlaylistTable: mockTable,
	}

	router.DELETE(route, playlistHandler.DeletePlaylist)

	recorder := httptest.NewRecorder()

	// Unable to parse to int, will throw error
	_route := "/playlist/sdfsd"

	req, _ := http.NewRequest("DELETE", _route, nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestDeletePlaylistPlaylistId(t *testing.T) {
	// Arrange
	route := "/playlist/:playlistId"
	router := SetupTestRouter()

	mockTable := &mockPlaylistTable{
		deletePlaylist: func(playlistId int) (error error) {
			return nil
		},
	}

	playlistHandler := PlaylistHandler{
		PlaylistTable: mockTable,
	}

	router.DELETE(route, playlistHandler.DeletePlaylist)

	recorder := httptest.NewRecorder()

	// Unable to parse to int, will throw error
	_route := "/playlist/1"

	req, _ := http.NewRequest("DELETE", _route, nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
}
