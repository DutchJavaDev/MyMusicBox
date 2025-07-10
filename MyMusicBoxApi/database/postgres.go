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

type PostgresDb struct {
	connection *sql.DB
	Error      error
	closed     bool
}

func (pdb *PostgresDb) OpenConnection() (created bool) {

	baseConnectionString := ""

	if configuration.Config.UseDevUrl {
		baseConnectionString = "user=postgres dbname=postgres password=%s host=127.0.0.1 port=5433 sslmode=disable"
	} else {
		baseConnectionString = "user=postgres dbname=postgres password=%s host=127.0.0.1 port=5432 sslmode=disable"
	}

	connectionString := fmt.Sprintf(baseConnectionString, os.Getenv("POSTGRES_PASSWORD"))

	pdb.connection, pdb.Error = sql.Open("postgres", connectionString)

	pdb.connection.SetMaxOpenConns(75)
	pdb.connection.SetMaxIdleConns(5)
	pdb.connection.SetConnMaxIdleTime(10 * time.Minute)
	pdb.connection.SetConnMaxLifetime(5 * time.Minute)

	if pdb.Error != nil {
		logging.Error(fmt.Sprintf("Failed to init database connection: %s", pdb.Error.Error()))
		return false
	}
	return true
}

func (pdb *PostgresDb) CloseConnection() {

	if pdb.closed {
		return
	}

	pdb.connection.Close()
	pdb.closed = true
}

func (pdb *PostgresDb) NonScalarQuery(query string, params ...any) (error error) {

	transaction, err := pdb.connection.Begin()

	if err != nil {
		logging.Error(fmt.Sprintf("[NonScalarQuery] Transaction error: %s", err.Error()))
		return err
	}

	statement, err := transaction.Prepare(query)

	if err != nil {
		logging.Error(fmt.Sprintf("[NonScalarQuery] Prepared statement error: %s", err.Error()))
		return err
	}

	_, err = statement.Exec(params...)

	if err != nil {
		logging.Error(fmt.Sprintf("[NonScalarQuery] Exec error: %s", err.Error()))
		logging.Error(fmt.Sprintf("Query: %s", query))
		for index := range params {
			logging.Error(params[index])
		}
		return err
	}

	err = transaction.Commit()

	if err != nil {
		logging.Error(fmt.Sprintf("[NonScalarQuery] Transaction commit error: %s", err.Error()))
		return err
	}

	return nil
}

func (pdb *PostgresDb) InsertWithReturningId(query string, params ...any) (lastInsertedId int, err error) {

	if !strings.Contains(query, "RETURNING") {
		logging.Error("[InsertWithReturningId] Query does not contain RETURNING keyword")
		return -1, errors.New("Query does not contain RETURNING keyword")
	}

	transaction, err := pdb.connection.Begin()

	statement, err := transaction.Prepare(query)
	defer statement.Close()

	if err != nil {
		transaction.Rollback()
		logging.Error(fmt.Sprintf("[InsertWithReturningId] Prepared statement error: %s", err.Error()))
		return -1, err
	}

	err = statement.QueryRow(params...).Scan(&lastInsertedId)

	if err != nil {
		logging.Error(fmt.Sprintf("[InsertWithReturningId] Queryrow error: %s", err.Error()))
		transaction.Rollback()
		return -1, err
	}

	err = transaction.Commit()

	if err != nil {
		logging.Error(fmt.Sprintf("[InsertWithReturningId] Transaction commit error: %s", err.Error()))
		transaction.Rollback()
		return -1, err
	}

	return lastInsertedId, nil
}
