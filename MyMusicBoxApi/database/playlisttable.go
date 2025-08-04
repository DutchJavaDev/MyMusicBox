package database

import (
	"context"
	"fmt"
	"musicboxapi/logging"
	"musicboxapi/models"
)

type IPlaylistTable interface {
	FetchPlaylists(ctx context.Context, lastKnowPlaylistId int) (playlists []models.Playlist, error error)
	InsertPlaylist(playlist models.Playlist) (lastInsertedId int, error error)
	DeletePlaylist(playlistId int) (error error)
}

type PlaylistTable struct {
	BaseTable
}

func NewPlaylistTableInstance() IPlaylistTable {
	return &PlaylistTable{
		BaseTable: NewBaseTableInstance(),
	}
}

func (table *PlaylistTable) FetchPlaylists(ctx context.Context, lastKnowPlaylistId int) (playlists []models.Playlist, error error) {
	query := "SELECT Id, Name, ThumbnailPath, Description, CreationDate FROM Playlist WHERE Id > $1 ORDER BY Id" // order by?

	rows, err := table.QueryRowsContex(ctx, query, lastKnowPlaylistId)

	if err != nil {
		logging.Error(fmt.Sprintf("QueryRow error: %s", err.Error()))
		logging.ErrorStackTrace(err)
		return nil, err
	}

	defer rows.Close()

	var playlist models.Playlist

	playlists = make([]models.Playlist, 0)

	for rows.Next() {
		scanError := rows.Scan(&playlist.Id, &playlist.Name, &playlist.ThumbnailPath, &playlist.Description, &playlist.CreationDate)

		if scanError != nil {
			logging.Error(fmt.Sprintf("Scan error: %s", scanError.Error()))
			continue
		}

		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func (table *PlaylistTable) InsertPlaylist(playlist models.Playlist) (lastInsertedId int, error error) {
	query := `INSERT INTO Playlist (name, description, thumbnailPath) VALUES ($1, $2, $3) RETURNING Id`

	lastInsertedId, err := table.InsertWithReturningId(query,
		playlist.Name,
		playlist.Description,
		playlist.ThumbnailPath,
	)

	return lastInsertedId, err
}

func (table *PlaylistTable) DeletePlaylist(playlistId int) (error error) {
	query := `DELETE FROM Playlist WHERE Id = $1`

	err := table.NonScalarQuery(query, playlistId)

	if err != nil {
		logging.Error(fmt.Sprintf("Failed to delete playlist: %s", err.Error()))
		return err
	}

	playlistsongTable := NewPlaylistsongTableInstance()

	return playlistsongTable.DeleteAllPlaylistSongs(playlistId)
}
