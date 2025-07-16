package database

import (
	"context"
	"fmt"
	"musicboxapi/logging"
	"musicboxapi/models"
)

type IPlaylistsongTable interface {
	FetchPlaylistSongs(ctx context.Context, playlistId int, lastKnowPosition int) (songs []models.Song, error error)
	InsertPlaylistSong(playlistId int, songId int) (lastInsertedId int, error error)
	DeleteAllPlaylistSongs(playlistId int) (error error)
	DeletePlaylistSong(playlistId int, songId int) (error error)
}

type PlaylistsongTable struct {
	BaseTable
}

func NewPlaylistsongTableInstance() *PlaylistsongTable {
	return &PlaylistsongTable{
		BaseTable: NewBaseTableInstance(),
	}
}

func (pt *PlaylistsongTable) FetchPlaylistSongs(ctx context.Context, playlistId int, lastKnowPosition int) (songs []models.Song, error error) {
	query := `SELECT s.Id, s.Name, s.Path, s.ThumbnailPath, s.Duration, s.SourceId, s.UpdatedAt, s.CreatedAt FROM playlistsong ps
			 INNER JOIN song s ON s.id = ps.songid
			 WHERE ps.playlistid = $1 AND ps.position >= $2`

	statement, err := pt.DB.Prepare(query)

	if err != nil {
		logging.Error(fmt.Sprintf("Prepared statement error: %s", err.Error()))
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.QueryContext(ctx, playlistId, lastKnowPosition)

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

func (pt *PlaylistsongTable) InsertPlaylistSong(playlistId int, songId int) (lastInsertedId int, error error) {
	query := `INSERT INTO PlaylistSong (SongId, PlaylistId) VALUES($1, $2) RETURNING SongId`

	lastInsertedId, err := pt.InsertWithReturningId(query,
		songId,
		playlistId,
	)

	return lastInsertedId, err
}

func (pt *PlaylistsongTable) DeleteAllPlaylistSongs(playlistId int) (error error) {
	query := `DELETE FROM PlaylistSong WHERE PlaylistId = $1`

	err := pt.NonScalarQuery(query, playlistId)

	return err
}

func (pt *PlaylistsongTable) DeletePlaylistSong(playlistId int, songId int) (error error) {
	query := `DELETE FROM PlaylistSong WHERE PlaylistId = $1 and SongId = $2`

	err := pt.NonScalarQuery(query, playlistId, songId)

	return err
}
