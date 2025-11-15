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
	FetchPlaylistsById(ctx context.Context, playlistId int) (plst models.Playlist, error error)
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
	query := "SELECT Id, Name, ThumbnailPath, Description, CreationDate, IsPublic FROM Playlist WHERE Id > $1 ORDER BY Id" // order by?

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
		scanError := rows.Scan(&playlist.Id, &playlist.Name, &playlist.ThumbnailPath, &playlist.Description, &playlist.CreationDate, &playlist.IsPublic)

		if scanError != nil {
			logging.Error(fmt.Sprintf("Scan error: %s", scanError.Error()))
			continue
		}

		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func (table *PlaylistTable) FetchPlaylistsById(ctx context.Context, playlistId int) (plst models.Playlist, error error) {
	query := "SELECT Id, Name, ThumbnailPath, Description, CreationDate, IsPublic FROM Playlist WHERE Id = $1" //
	//rows, err := table.QueryRowsContex(ctx, query, playlistId)

	row := table.QueryRow(query, playlistId)

	var playlist models.Playlist

	scanError := row.Scan(&playlist.Id, &playlist.Name, &playlist.ThumbnailPath, &playlist.Description, &playlist.CreationDate, &playlist.IsPublic)

	if scanError != nil {
		logging.Error(fmt.Sprintf("Scan error: %s", scanError.Error()))
		return models.Playlist{}, scanError
	}

	return playlist, nil
}

func (table *PlaylistTable) InsertPlaylist(playlist models.Playlist) (lastInsertedId int, error error) {
	query := `INSERT INTO Playlist (name, description, thumbnailPath, ispublic) VALUES ($1, $2, $3, $4) RETURNING Id`

	lastInsertedId, err := table.InsertWithReturningId(query,
		playlist.Name,
		playlist.Description,
		playlist.ThumbnailPath,
		playlist.IsPublic,
	)

	return lastInsertedId, err
}

func (table *PlaylistTable) DeletePlaylist(playlistId int) (error error) {
	query := `DELETE FROM Playlist WHERE Id = $1 AND IsPublic = $2` // Prevemts private playlists (like the default one) from being deleted for real

	err := table.NonScalarQuery(query, playlistId, true)

	if err != nil {
		logging.Error(fmt.Sprintf("Failed to delete playlist: %s", err.Error()))
		return err
	}

	playlistsongTable := NewPlaylistsongTableInstance()

	return playlistsongTable.DeleteAllPlaylistSongs(playlistId)
}
