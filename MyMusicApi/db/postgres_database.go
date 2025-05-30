package db

import (
	"api/logging"
	"api/models"
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type PostgresDb struct {
	connection *sql.DB
	Error      error
}

func (pdb *PostgresDb) OpenConnection() (created bool) {

	baseConnectionString := "user=postgres dbname=postgres password=%s host=127.0.0.1 port=5432 sslmode=disable"

	connectionString := fmt.Sprintf(baseConnectionString, os.Getenv("POSTGRES_PASSWORD"))

	pdb.connection, pdb.Error = sql.Open("postgres", connectionString)

	if pdb.Error != nil {
		logging.Error(fmt.Sprintf("Failed to init database connection: %s", pdb.Error.Error()))
		return false
	}
	return true
}

func (pdb *PostgresDb) CloseConnection() {
	pdb.connection.Close()
}

// begin fetch
func (pdb *PostgresDb) FetchSongs(ctx context.Context) (songs []models.Song, error error) {

	query := "SELECT Id, Name, Path, Duration, SourceURL, UpdatedAt FROM Song" // order by?

	rows, err := pdb.connection.QueryContext(ctx, query)
	defer rows.Close()

	if err != nil {
		logging.Error(fmt.Sprintf("[FetchSongs] QueryRow error: %s", err.Error()))
		return nil, err
	}

	var song models.Song

	songs = make([]models.Song, 0)

	for rows.Next() {
		scanError := rows.Scan(&song.Id, &song.Name, &song.Path, &song.Duration, &song.SourceURL, &song.UpdatedAt)

		if scanError != nil {
			logging.Error(fmt.Sprintf("[FetchSongs] Scan error: %s", scanError.Error()))
			continue
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func (pdb *PostgresDb) FetchPlaylists(ctx context.Context) (playlists []models.Playlist, error error) {
	query := "SELECT Id, Name, Description FROM Playlist" // order by?

	rows, err := pdb.connection.QueryContext(ctx, query)
	defer rows.Close()

	if err != nil {
		logging.Error(fmt.Sprintf("[FetchPlaylists] QueryRow error: %s", err.Error()))
		return nil, err
	}

	var playlist models.Playlist

	playlists = make([]models.Playlist, 0)

	for rows.Next() {
		scanError := rows.Scan(&playlist.Id, &playlist.Name, &playlist.Description)

		if scanError != nil {
			logging.Error(fmt.Sprintf("[FetchPlaylists] Scan error: %s", scanError.Error()))
			continue
		}

		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func (pdb *PostgresDb) FetchPlaylistSongs(ctx context.Context, playlistId int) (songs []models.Song, error error) {

	query := `SELECT s.Id, s.Name, s.Path, s.Duration, s.SourceURL, s.UpdatedAt FROM Song s
	          INNER JOIN PlaylistSong ps ON ps.PlaylistId = $1
			  WHERE ps.SongId = s.Id
			  order by ps.Position` // order by playlist position

	statement, err := pdb.connection.Prepare(query)

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
		scanError := rows.Scan(&song.Id, &song.Name, &song.Path, &song.Duration, &song.SourceURL, &song.UpdatedAt)

		if scanError != nil {
			logging.Error(fmt.Sprintf("[FetchPlaylistSongs] Scan error: %s", scanError.Error()))
			continue
		}

		songs = append(songs, song)
	}

	return songs, nil
}

//end fetch

// begin insert
func (pdb *PostgresDb) InsertSong(song models.Song) (lastInsertedId int, error error) {

	query := `INSERT INTO Song (name, sourceurl, path, duration) VALUES ($1, $2, $3, $4) RETURNING Id`

	// could add request context?
	transaction, err := pdb.connection.Begin()

	statement, err := transaction.Prepare(query)
	defer statement.Close()

	if err != nil {
		transaction.Rollback()
		logging.Error(fmt.Sprintf("[InsertSong] Prepared statement error: %s", err.Error()))
		return -1, err
	}

	err = statement.QueryRow(song.Name, song.SourceURL, song.Path, song.Duration).Scan(&lastInsertedId)

	if err != nil {
		logging.Error(fmt.Sprintf("[InsertSong] Queryrow error: %s", err.Error()))
		transaction.Rollback()
		return -1, err
	}

	err = transaction.Commit()

	if err != nil {
		logging.Error(fmt.Sprintf("[InsertSong] Transaction commit error: %s", err.Error()))
		transaction.Rollback()
		return -1, err
	}

	return lastInsertedId, nil
}

func (pdb *PostgresDb) InsertPlaylist(playlist models.Playlist) (lastInsertedId int, error error) {
	query := `INSERT INTO Playlist (name, description) VALUES ($1, $2) RETURNING Id`

	// could add request context?
	transaction, err := pdb.connection.Begin()

	statement, err := transaction.Prepare(query)
	defer statement.Close()

	if err != nil {
		transaction.Rollback()
		logging.Error(fmt.Sprintf("[InsertPlaylist] Prepared statement error: %s", err.Error()))
		return -1, err
	}

	err = statement.QueryRow(playlist.Name, playlist.Description).Scan(&lastInsertedId)

	if err != nil {
		logging.Error(fmt.Sprintf("[InsertPlaylist] Queryrow error: %s", err.Error()))
		transaction.Rollback()
		return -1, err
	}

	err = transaction.Commit()

	if err != nil {
		logging.Error(fmt.Sprintf("[InsertPlaylist] Transaction commit error: %s", err.Error()))
		transaction.Rollback()
		return -1, err
	}

	return lastInsertedId, nil
}

func (pdb *PostgresDb) InsertPlaylistSong(playlistId int, songId int) (lastInsertedId int, error error) {
	query := `INSERT INTO PlaylistSong (SongId, PlaylistId) VALUES($1, $2) RETURNING SongId`

	// could add request context?
	transaction, err := pdb.connection.Begin()

	statement, err := transaction.Prepare(query)
	defer statement.Close()

	if err != nil {
		transaction.Rollback()
		logging.Error(fmt.Sprintf("[InsertPlaylistSong] Prepared statement error: %s", err.Error()))
		return -1, err
	}

	err = statement.QueryRow(songId, playlistId).Scan(&lastInsertedId)

	if err != nil {
		logging.Error(fmt.Sprintf("[InsertPlaylistSong] Queryrow error: %s", err.Error()))
		transaction.Rollback()
		return -1, err
	}

	err = transaction.Commit()

	if err != nil {
		logging.Error(fmt.Sprintf("[InsertPlaylistSong] Transaction commit error: %s", err.Error()))
		transaction.Rollback()
		return -1, err
	}

	return lastInsertedId, nil
}

//end insert

// begin delete
func (pdb *PostgresDb) DeletePlaylist(id int) (error error) {
	query := `DELETE FROM Playlist WHERE Id = $1`
	// Should also deleted linked playlistsong TODO
	transaction, err := pdb.connection.Begin()

	statement, err := transaction.Prepare(query)
	defer statement.Close()

	if err != nil {
		transaction.Rollback()
		logging.Error(fmt.Sprintf("[DeletePlaylistById] Prepared statement error: %s", err.Error()))
		return err
	}

	_, err = statement.Exec(id)

	if err != nil {
		logging.Error(fmt.Sprintf("[DeletePlaylistById] Execute error: %s", err.Error()))
		transaction.Rollback()
		return err
	}

	err = transaction.Commit()

	if err != nil {
		logging.Error(fmt.Sprintf("[DeletePlaylistById] Transaction commmit error: %s", err.Error()))
		transaction.Rollback()
		return err
	}

	return nil
}

func (pdb *PostgresDb) DeletePlaylistSong(playlistId int, songId int) (error error) {
	query := `DELETE FROM PlaylistSong WHERE PlaylistId = $1 and SongId = $2`

	transaction, err := pdb.connection.Begin()

	statement, err := transaction.Prepare(query)
	defer statement.Close()

	if err != nil {
		transaction.Rollback()
		logging.Error(fmt.Sprintf("[DeletePlaylistSong] Prepared statement error: %s", err.Error()))
		return err
	}

	_, err = statement.Exec(playlistId, songId)

	if err != nil {
		transaction.Rollback()
		logging.Error(fmt.Sprintf("[DeletePlaylistSong] Execute error: %s", err.Error()))
		return err
	}

	err = transaction.Commit()

	if err != nil {
		logging.Error(fmt.Sprintf("[DeletePlaylistSong] Transaction commmit error: %s", err.Error()))
		transaction.Rollback()
		return err
	}

	return nil
}

//end delete

// TODO
// func (pdb *PostgresDb) AddLog(log models.Log) {
// }
