package database

import (
	"context"
	"fmt"
	"musicboxapi/logging"
	"musicboxapi/models"
)

type ISongTable interface {
	InsertSong(song *models.Song) (err error)
	FetchSongs(ctx context.Context) (songs []models.Song, err error)
}

type SongTable struct {
	BaseTable
}

func NewSongTableInstance() *SongTable {
	return &SongTable{
		BaseTable: NewBaseTableInstance(),
	}
}

func (st *SongTable) InsertSong(song *models.Song) (error error) {

	query := `INSERT INTO Song (name, sourceid, path, thumbnailPath, duration) VALUES ($1, $2, $3, $4, $5) RETURNING Id`

	lastInsertedId, err := st.InsertWithReturningId(query,
		song.Name,
		song.SourceId,
		song.Path,
		song.ThumbnailPath,
		song.Duration,
	)

	// Add to main playlist
	defaultPlaylistId := 1
	playlistsongTable := NewPlaylistsongTableInstance()

	if err == nil {
		_, err = playlistsongTable.InsertPlaylistSong(defaultPlaylistId, lastInsertedId)
	}

	song.Id = lastInsertedId

	return err
}

func (st *SongTable) FetchSongs(ctx context.Context) (songs []models.Song, error error) {

	query := "SELECT Id, Name, Path, ThumbnailPath, Duration, SourceId, UpdatedAt, CreatedAt FROM Song" // order by?

	rows, err := st.DB.QueryContext(ctx, query)

	if err != nil {
		logging.Error(fmt.Sprintf("QueryRow error: %s", err.Error()))
		return nil, err
	}

	defer rows.Close()

	var song models.Song

	songs = make([]models.Song, 0)

	for rows.Next() {
		scanError := rows.Scan(&song.Id, &song.Name, &song.Path, &song.ThumbnailPath, &song.Duration, &song.SourceId, &song.UpdatedAt, &song.CreatedAt)

		if scanError != nil {
			logging.Error(fmt.Sprintf("Scan error: %s", scanError.Error()))
			continue
		}

		songs = append(songs, song)
	}

	return songs, nil
}
