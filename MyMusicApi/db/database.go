package db

import (
	"api/logging"
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

func (pdb *PostgresDb) InitDatabase() (created bool) {

	baseConnectionString := "user=postgres dbname=postgres password=%s host=127.0.0.1 port=5432 sslmode=disable"

	connectionString := fmt.Sprintf(baseConnectionString, os.Getenv("POSTGRES_PASSWORD"))

	pdb.connection, pdb.Error = sql.Open("postgres", connectionString)

	if pdb.Error != nil {
		logging.Error(fmt.Sprintf("Failed to init database connection: %s", pdb.Error.Error()))
		return false
	}
	return true
}

func (pdb *PostgresDb) Close() {
	pdb.connection.Close()
}

func (pdb *PostgresDb) GetAllSongs(ctx context.Context) (songs []Song, error error) {

	query := "SELECT Id, Name, Path, Duration, SourceURL, UpdatedAt FROM Song" // order by?

	rows, err := pdb.connection.QueryContext(ctx, query)
	defer rows.Close()

	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	var song Song

	songs = make([]Song, 0)

	for rows.Next() {
		scanError := rows.Scan(&song.Id, &song.Name, &song.Path, &song.Duration, &song.SourceURL, &song.UpdatedAt)

		if scanError != nil {
			logging.Error(scanError.Error())
			continue
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func (pdb *PostgresDb) AddSong(song Song) (lastInsertedId int64, error error) {

	query := `INSERT INTO song (name, sourceurl, path, duration) VALUES ($1, $2, $3, $4) RETURNING Id`

	// could add request context?
	transaction, err := pdb.connection.Begin()

	statement, err := transaction.Prepare(query)
	defer statement.Close()

	if err != nil {
		transaction.Rollback()
		logging.Error(err.Error())
		return -1, err
	}

	err = statement.QueryRow(song.Name, song.SourceURL, song.Path, song.Duration).Scan(&lastInsertedId)

	if err != nil {
		logging.Error(err.Error())
		transaction.Rollback()
		return -1, err
	}

	err = transaction.Commit()

	if err != nil {
		logging.Error(err.Error())
		transaction.Rollback()
		return -1, err
	}

	return lastInsertedId, nil
}
