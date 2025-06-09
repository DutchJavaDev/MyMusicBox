package database

import (
	"context"
	"fmt"
	"musicboxapi/logging"
	"musicboxapi/models"
)

func (pdb *PostgresDb) InsertSong(song *models.Song) (error error) {

	query := `INSERT INTO Song (name, sourceid, path, thumbnailPath, duration) VALUES ($1, $2, $3, $4, $5) RETURNING Id`

	lastInsertedId, err := pdb.InsertWithReturningId(query,
		song.Name,
		song.SourceId,
		song.Path,
		song.ThumbnailPath,
		song.Duration,
	)

	// Add to main playlist
	if err == nil {
		_, err = pdb.InsertPlaylistSong(1, lastInsertedId)
	}

	song.Id = lastInsertedId

	return err
}

func (pdb *PostgresDb) FetchSongs(ctx context.Context) (songs []models.Song, error error) {

	query := "SELECT Id, Name, Path, ThumbnailPath, Duration, SourceId, UpdatedAt, CreatedAt FROM Song" // order by?

	rows, err := pdb.connection.QueryContext(ctx, query)
	defer rows.Close()

	if err != nil {
		logging.Error(fmt.Sprintf("[FetchSongs] QueryRow error: %s", err.Error()))
		return nil, err
	}

	var song models.Song

	songs = make([]models.Song, 0)

	for rows.Next() {
		scanError := rows.Scan(&song.Id, &song.Name, &song.Path, &song.ThumbnailPath, &song.Duration, &song.SourceId, &song.UpdatedAt, &song.CreatedAt)

		if scanError != nil {
			logging.Error(fmt.Sprintf("[FetchSongs] Scan error: %s", scanError.Error()))
			continue
		}

		songs = append(songs, song)
	}

	return songs, nil
}
