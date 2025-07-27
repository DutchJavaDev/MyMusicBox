package database

import (
	"database/sql/driver"
	"musicboxapi/logging"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	// Arrange
	db, mock, err := CreateMockDb()

	if err != nil {
		logging.ErrorStackTrace(err)
	}

	migration := NewMigrationTableInstance()
	migration.BaseTable.DB = db

	query := "INSERT INTO Migration (filename, contents) VALUES($1, $2)"

	fileName := "999 fix db.sql"
	contents := "drop database migration"

	mock.ExpectBegin()
	mockPrepare := mock.ExpectPrepare(regexp.QuoteMeta(query))
	mockPrepare.ExpectExec().
		WithArgs(fileName, contents).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockPrepare.WillBeClosed()
	mock.ExpectCommit()

	// Act
	result := migration.Insert(fileName, contents)

	// Assert
	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Errorf("Expectations failed: %s", err)
	}

	if result != nil {
		t.Errorf("Insert failed: %s", result)
	}
}

func TestApplyMigration(t *testing.T) {
	// Arrange
	db, mock, err := CreateMockDb()

	if err != nil {
		logging.ErrorStackTrace(err)
	}

	migration := NewMigrationTableInstance()
	migration.BaseTable.DB = db

	query := "DROP DATABSE migration"

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Act
	result := migration.ApplyMigration(query)

	// Assert
	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Errorf("Expectations failed: %s", err)
	}

	if result != nil {
		t.Errorf("ApplyMigration failed: %s", result)
	}
}

func TestGetCurrentAppliedMigrationFileName(t *testing.T) {
	// Arrange
	db, mock, err := CreateMockDb()

	if err != nil {
		logging.ErrorStackTrace(err)
	}

	migration := NewMigrationTableInstance()
	migration.BaseTable.DB = db

	query := "SELECT filename FROM migration order by AppliedOn DESC LIMIT 1"

	value := []driver.Value{
		"999 fix db.sql",
	}

	row := sqlmock.NewRows([]string{"filename"}).AddRow(value...)

	mock.ExpectQuery(query).
		WillReturnRows(row)

	// Act
	filename, result := migration.GetCurrentAppliedMigrationFileName()

	// Assert
	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Errorf("Expectations failed: %s", err)
	}

	if result != nil {
		t.Errorf("GetCurrentAppliedMigrationFileName failed: %s", result)
	}

	assert.Equal(t, value[0], filename)
}
