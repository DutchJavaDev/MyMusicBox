package database

import (
	"context"
	"fmt"
	"musicboxapi/logging"
	"musicboxapi/models"
)

func (pdb *PostgresDb) FetchPlaylistSongs(ctx context.Context, playlistId int) (songs []models.Song, error error) {

	query := `SELECT s.Id, s.Name, s.Path, s.ThumbnailPath, s.Duration, s.SourceId, s.UpdatedAt, CreatedAt FROM Song s
	          INNER JOIN PlaylistSong ps ON ps.PlaylistId = $1
			  WHERE ps.SongId = s.Id
			  order by ps.Position` // order by playlist position

	statement, err := pdb.connection.Prepare(query)
	defer statement.Close()

	if err != nil {
		logging.Error(fmt.Sprintf("[FetchPlaylistSongs] Prepared statement error: %s", err.Error()))
		return nil, err
	}

	rows, err := statement.QueryContext(ctx, playlistId)
	defer rows.Close()

	if err != nil {
		logging.Error(fmt.Sprintf("[FetchPlaylistSongs] QueryRow error: %s", err.Error()))
		return nil, err
	}

	var song models.Song

	songs = make([]models.Song, 0)

	for rows.Next() {
		scanError := rows.Scan(&song.Id, &song.Name, &song.Path, &song.ThumbnailPath, &song.Duration, &song.SourceId, &song.UpdatedAt, &song.CreatedAt)

		if scanError != nil {
			logging.Error(fmt.Sprintf("[FetchPlaylistSongs] Scan error: %s", scanError.Error()))
			continue
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func (pdb *PostgresDb) InsertPlaylistSong(playlistId int, songId int) (lastInsertedId int, error error) {
	query := `INSERT INTO PlaylistSong (SongId, PlaylistId) VALUES($1, $2) RETURNING SongId`

	lastInsertedId, err := pdb.InsertWithReturningId(query,
		songId,
		playlistId,
	)

	return lastInsertedId, err
}

func (pdb *PostgresDb) deletePlaylistSongs(playlistId int) (error error) {
	query := `DELETE FROM PlaylistSong WHERE PlaylistId = $1`

	err := pdb.NonScalarQuery(query, playlistId)

	return err
}

func (pdb *PostgresDb) DeletePlaylistSong(playlistId int, songId int) (error error) {
	query := `DELETE FROM PlaylistSong WHERE PlaylistId = $1 and SongId = $2`

	err := pdb.NonScalarQuery(query, playlistId, songId)

	return err
}
