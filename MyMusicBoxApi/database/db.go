package database

import (
	"database/sql"
	"errors"
	"fmt"
	"musicboxapi/configuration"
	"musicboxapi/logging"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const ReturningIdParameter = "RETURNING"
const ReturningIdParameterLower = "returning"
const DatabaseDriver = "postgres"
const MigrationFolder = "migration_scripts"
const MaxOpenConnections = 10
const MaxIdleConnections = 5
const MaxConnectionIdleTimeInMinutes = 10
const MaxConnectionLifeTimeInMinutes = 10

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
	// Should create test for these?
	var _ ISongTable = (*SongTable)(nil)
	var _ IPlaylistTable = (*PlaylistTable)(nil)
	var _ IPlaylistsongTable = (*PlaylistsongTable)(nil)
	var _ ITasklogTable = (*TasklogTable)(nil)
	var _ IMigrationTable = (*MigrationTable)(nil)

	baseConnectionString := "user=postgres dbname=postgres password=%s %s sslmode=disable"
	password := os.Getenv("POSTGRES_PASSWORD")
	host := "host=127.0.0.1 port=5432"

	if configuration.Config.UseDevUrl {
		host = "host=127.0.0.1 port=5433"
	}

	connectionString := fmt.Sprintf(baseConnectionString, password, host)

	DB, err := sql.Open(DatabaseDriver, connectionString)

	if err != nil {
		logging.Error(fmt.Sprintf("Failed to init database connection: %s", err.Error()))
		logging.ErrorStackTrace(err)
		return err
	}

	DB.SetMaxOpenConns(MaxOpenConnections)
	DB.SetMaxIdleConns(MaxIdleConnections)
	DB.SetConnMaxIdleTime(MaxConnectionIdleTimeInMinutes * time.Minute)
	DB.SetConnMaxLifetime(MaxConnectionLifeTimeInMinutes * time.Minute)

	DbInstance = DB

	return nil
}

// Base methods
func (base *BaseTable) InsertWithReturningId(query string, params ...any) (lastInsertedId int, err error) {

	if !strings.Contains(query, ReturningIdParameter) {
		return -1, errors.New("Query does not contain RETURNING keyword")
	}

	transaction, err := base.DB.Begin()

	statement, err := transaction.Prepare(query)

	if err != nil {
		logging.ErrorStackTrace(err)
		return -1, err
	}
	defer statement.Close()

	err = statement.QueryRow(params...).Scan(&lastInsertedId)

	if err != nil {
		logging.ErrorStackTrace(err)
		transaction.Rollback()
		return -1, err
	}

	err = transaction.Commit()

	if err != nil {
		logging.ErrorStackTrace(err)
		transaction.Rollback()
		return -1, err
	}

	return lastInsertedId, nil
}

func (base *BaseTable) NonScalarQuery(query string, params ...any) (error error) {

	transaction, err := base.DB.Begin()

	if err != nil {
		logging.Error(fmt.Sprintf("Transaction error: %s", err.Error()))
		logging.ErrorStackTrace(err)
		return err
	}

	statement, err := transaction.Prepare(query)

	if err != nil {
		logging.Error(fmt.Sprintf("Prepared statement error: %s", err.Error()))
		logging.ErrorStackTrace(err)
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(params...)

	if err != nil {
		logging.ErrorStackTrace(err)
		return err
	}

	err = transaction.Commit()

	if err != nil {
		logging.Error(fmt.Sprintf("Transaction commit error: %s", err.Error()))
		logging.ErrorStackTrace(err)
		return err
	}

	return nil
}

func ApplyMigrations() {
	logging.Info("Applying migrations...")
	// files will be sorted by filename
	// to make sure the migrations are executed in order
	// this naming convention must be used
	// 0 initial script.sql
	// 1 update column.sql
	// etc....
	// entries are sorted by file name
	dirs, err := os.ReadDir(MigrationFolder)

	if err != nil {
		logging.ErrorStackTrace(err)
		return
	}

	migrationTable := NewMigrationTableInstance()

	currentMigrationFileName, err := migrationTable.GetCurrentAppliedMigrationFileName()

	// start at -1 if lastMigrationFileName is empty OR migration table does not exists
	// start applying from 0
	if currentMigrationFileName == "" || err != nil {
		if strings.Contains(err.Error(), `relation "migration" does not exist`) {
			logging.Info("First time running database script migrations")
		} else {
			logging.ErrorStackTrace(err)
			return
		}

		// makes sure we start at script 0
		currentMigrationFileName = "-1 nil.sql"
	}

	currentMigrationFileId, err := strconv.Atoi(strings.Split(currentMigrationFileName, " ")[0])

	if err != nil {
		logging.ErrorStackTrace(err)
		return
	}

	for _, migrationFile := range dirs {
		filePath := filepath.Join(MigrationFolder, migrationFile.Name())

		migrationFileId, err := strconv.Atoi(strings.Split(migrationFile.Name(), " ")[0])

		if err != nil {
			logging.ErrorStackTrace(err)
			continue
		}

		if migrationFileId <= currentMigrationFileId {
			continue
		}

		migrationFileContents, err := os.ReadFile(filePath)

		if err != nil {
			logging.ErrorStackTrace(err)
			continue
		}

		err = migrationTable.ApplyMigration(string(migrationFileContents))

		if err != nil {
			logging.Error(fmt.Sprintf("Failed to apply %s", migrationFile.Name()))
			logging.ErrorStackTrace(err)
		} else {
			err = migrationTable.Insert(migrationFile.Name(), string(migrationFileContents))

			if err != nil {
				logging.Error(fmt.Sprintf("Failed to insert migration entry %s: %s", migrationFile.Name(), err.Error()))
				logging.ErrorStackTrace(err)
				return
			}

			logging.Info(fmt.Sprintf("Applied script: %s", migrationFile.Name()))
		}

	}
}
