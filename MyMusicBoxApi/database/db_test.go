package database

import (
	"database/sql"
	"musicboxapi/logging"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
)

var Mock sqlmock.Sqlmock

func CreateMockDb() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()

	if err != nil {
		logging.ErrorStackTrace(err)
		return db, nil, nil
	}
	return db, mock, nil
}

func TestInsertWithReturningIdShouldThrowError(t *testing.T) {

	// Arrange
	query := "INSERT INTO DATA (id) VALUES ($1)"
	db, _, err := CreateMockDb()

	if err != nil {
		logging.ErrorStackTrace(err)
		return
	}

	base := NewBaseTableInstance()
	base.DB = db

	// Act
	_, err = base.InsertWithReturningId(query, 1)

	// Assert
	if err == nil {
		t.Errorf("Method should have thrown an error")
	}
}

func TestInsertWithReturningId(t *testing.T) {

	// Arrange
	query := "INSERT INTO DATA (id) VALUES($1) RETURNING Id"
	insertId := 1
	db, mock, err := CreateMockDb()

	if err != nil {
		logging.ErrorStackTrace(err)
		return
	}

	base := NewBaseTableInstance()
	base.DB = db

	mock.ExpectBegin()
	mockPrepare := mock.ExpectPrepare(regexp.QuoteMeta(query))
	mockPrepare.ExpectQuery().
		WithArgs(insertId).
		WillReturnRows(sqlmock.NewRows([]string{"Id"}).AddRow(insertId))
	mockPrepare.WillBeClosed()
	mock.ExpectCommit()

	// Act
	lastInsertedId, err := base.InsertWithReturningId(query, 1)

	// Assert
	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Errorf("Expectations failed: %s", err)
	}

	if lastInsertedId != insertId {
		t.Errorf("Expected %d, got %d instead", insertId, lastInsertedId)
	}
}

func TestNonScalarQueryShouldThrowErrt(t *testing.T) {
	// Arrange
	query := "INSERT INTO DATA (id) VALUES($1)"
	db, mock, err := CreateMockDb()
	id := 1
	if err != nil {
		logging.ErrorStackTrace(err)
		return
	}

	base := NewBaseTableInstance()
	base.DB = db

	mock.ExpectBegin()

	mockPrepare := mock.ExpectPrepare(regexp.QuoteMeta(query))
	mockPrepare.ExpectExec().
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockPrepare.WillBeClosed()
	mock.ExpectCommit()

	// Act
	result := base.NonScalarQuery(query, id)

	// Assert
	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Errorf("Expectations failed: %s", err)
	}

	if result != nil {
		t.Errorf("NonScalarQuery failed: %s", result)
	}
}
