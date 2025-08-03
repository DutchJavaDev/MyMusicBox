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

type mockSongTable struct {
	database.SongTable
	fetchSongsFunc func(ctx context.Context) ([]models.Song, error)
}

func (m *mockSongTable) FetchSongs(ctx context.Context) ([]models.Song, error) {
	return m.fetchSongsFunc(ctx)
}

func TestFetchSongs(t *testing.T) {
	// Arrange
	route := "/api/v1/songs"
	router := SetupTestRouter()

	mockTable := &mockSongTable{
		fetchSongsFunc: func(ctx context.Context) ([]models.Song, error) {
			return []models.Song{
				{Id: 1, Name: "Toto - Africa", SourceId: "aqwerwe", ThumbnailPath: "image1.png", Path: "file1.opus", Duration: 360},
				{Id: 2, Name: "Hello - World", SourceId: "dfgdfgd", ThumbnailPath: "image2.png", Path: "file2.opus", Duration: 720},
			}, nil
		},
	}

	songHandler := SongHandler{
		SongTable: mockTable,
	}

	router.GET(route, songHandler.FetchSongs)

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

	var songs []models.Song
	err = json.Unmarshal(dataBytes, &songs)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(songs))
	assert.Equal(t, "Toto - Africa", songs[0].Name)
	assert.Equal(t, "Hello - World", songs[1].Name)
}
