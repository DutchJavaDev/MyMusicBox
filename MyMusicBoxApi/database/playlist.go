package database

import (
	"context"
	"fmt"
	"musicboxapi/logging"
	"musicboxapi/models"
)

func (pdb *PostgresDb) FetchPlaylists(ctx context.Context, lastKnowPlaylistId int) (playlists []models.Playlist, error error) {
	query := "SELECT Id, Name, ThumbnailPath, Description, CreationDate FROM Playlist WHERE Id > $1 ORDER BY Id" // order by?

	rows, err := pdb.connection.QueryContext(ctx, query, lastKnowPlaylistId)
	defer rows.Close()

	if err != nil {
		logging.Error(fmt.Sprintf("[FetchPlaylists] QueryRow error: %s", err.Error()))
		return nil, err
	}

	var playlist models.Playlist

	playlists = make([]models.Playlist, 0)

	for rows.Next() {
		scanError := rows.Scan(&playlist.Id, &playlist.Name, &playlist.ThumbnailPath, &playlist.Description, &playlist.CreationDate)

		if scanError != nil {
			logging.Error(fmt.Sprintf("[FetchPlaylists] Scan error: %s", scanError.Error()))
			continue
		}

		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func (pdb *PostgresDb) InsertPlaylist(playlist models.Playlist) (lastInsertedId int, error error) {
	query := `INSERT INTO Playlist (name, description, thumbnailPath) VALUES ($1, $2, $3) RETURNING Id`

	lastInsertedId, err := pdb.InsertWithReturningId(query,
		playlist.Name,
		playlist.Description,
		playlist.ThumbnailPath,
	)

	return lastInsertedId, err
}

func (pdb *PostgresDb) DeletePlaylist(playlistId int) (error error) {
	query := `DELETE FROM Playlist WHERE Id = $1`

	err := pdb.NonScalarQuery(query, playlistId)

	if err != nil {
		logging.Error(fmt.Sprintf("[DeletePlaylist] Failed to delete playlist: %s", err.Error()))
		return err
	}

	return pdb.deletePlaylistSongs(playlistId)
}
