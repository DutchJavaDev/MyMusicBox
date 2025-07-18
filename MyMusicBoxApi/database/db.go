package database

import (
	"database/sql"
	"errors"
	"fmt"
	"musicboxapi/configuration"
	"musicboxapi/logging"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var DbInstance *sql.DB

type BaseTable struct {
	DB *sql.DB
}

func NewBaseTableInstance() BaseTable {
	return BaseTable{
		DB: DbInstance,
	}
}

func CreateDatabasConnectionPool() error {

	// Will throw an error if its missing a method implementation from interface
	// will throw a compile time error
	var _ ISongTable = (*SongTable)(nil)
	var _ IPlaylistTable = (*PlaylistTable)(nil)
	var _ IPlaylistsongTable = (*PlaylistsongTable)(nil)
	var _ ITasklogTable = (*TasklogTable)(nil)

	baseConnectionString := "user=postgres dbname=postgres password=%s %s sslmode=disable"
	password := os.Getenv("POSTGRES_PASSWORD")
	host := "host=127.0.0.1 port=5432"

	if configuration.Config.UseDevUrl {
		host = "host=127.0.0.1 port=5433"
	}

	connectionString := fmt.Sprintf(baseConnectionString, password, host)

	DB, err := sql.Open("postgres", connectionString)

	if err != nil {
		logging.Error(fmt.Sprintf("Failed to init database connection: %s", err.Error()))
		return err
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxIdleTime(1 * time.Minute)
	DB.SetConnMaxLifetime(5 * time.Minute)

	DbInstance = DB

	return nil
}

// Base methods
func (base *BaseTable) InsertWithReturningId(query string, params ...any) (lastInsertedId int, err error) {

	if !strings.Contains(query, "RETURNING") {
		logging.Error("Query does not contain RETURNING keyword")
		return -1, errors.New("Query does not contain RETURNING keyword")
	}

	transaction, err := base.DB.Begin()

	statement, err := transaction.Prepare(query)

	if err != nil {
		transaction.Rollback()
		logging.Error(fmt.Sprintf("Prepared statement error: %s", err.Error()))
		return -1, err
	}
	defer statement.Close()

	err = statement.QueryRow(params...).Scan(&lastInsertedId)

	if err != nil {
		logging.Error(fmt.Sprintf("Queryrow error: %s", err.Error()))
		transaction.Rollback()
		return -1, err
	}

	err = transaction.Commit()

	if err != nil {
		logging.Error(fmt.Sprintf("Transaction commit error: %s", err.Error()))
		transaction.Rollback()
		return -1, err
	}

	return lastInsertedId, nil
}
func (base *BaseTable) NonScalarQuery(query string, params ...any) (error error) {

	transaction, err := base.DB.Begin()

	if err != nil {
		logging.Error(fmt.Sprintf("Transaction error: %s", err.Error()))
		return err
	}

	statement, err := transaction.Prepare(query)

	if err != nil {
		logging.Error(fmt.Sprintf("Prepared statement error: %s", err.Error()))
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(params...)

	if err != nil {
		logging.Error(fmt.Sprintf("Exec error: %s", err.Error()))
		logging.Error(fmt.Sprintf("Query: %s", query))
		for index := range params {
			logging.Error(params[index])
		}
		return err
	}

	err = transaction.Commit()

	if err != nil {
		logging.Error(fmt.Sprintf("Transaction commit error: %s", err.Error()))
		return err
	}

	return nil
}
