package database

import (
	"database/sql"
	"fmt"
	"musicboxapi/logging"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

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

func TestQueryRow(t *testing.T) {
	// Arrange
	query := "SELECT Id, Name, Online FROM data"
	db, mock, err := CreateMockDb()

	if err != nil {
		logging.ErrorStackTrace(err)
		return
	}

	base := NewBaseTableInstance()
	base.DB = db

	mock.ExpectQuery(query).
		WillReturnRows(sqlmock.NewRows([]string{"Id", "Name", "Online"}).AddRow(1, "Test", false))

	// Act
	rows := base.QueryRow(query)

	// Assert
	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Errorf("Expectations failed: %s", err)
	}

	var id int
	var name string
	var online bool

	rows.Scan(&id, &name, &online)

	assert.Equal(t, 1, id)
	assert.Equal(t, "Test", name)
	assert.Equal(t, false, online)
}

func TestQueryRows(t *testing.T) {
	// Arrange
	query := "SELECT Id, Name, Online FROM data"
	db, mock, err := CreateMockDb()

	if err != nil {
		logging.ErrorStackTrace(err)
		return
	}

	base := NewBaseTableInstance()
	base.DB = db

	mock.ExpectQuery(query).
		WillReturnRows(sqlmock.NewRows([]string{"Id", "Name", "Online"}).AddRow(1, "Test", false).AddRow(2, "Test2", true))

	// Act
	rows, err := base.QueryRows(query)

	assert.Nil(t, err)

	defer rows.Close()

	// Assert
	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Errorf("Expectations failed: %s", err)
	}

	type Data struct {
		id     int
		name   string
		online bool
	}

	var datas []Data
	var data Data

	datas = make([]Data, 0)

	for rows.Next() {
		scanError := rows.Scan(&data.id, &data.name, &data.online)

		if scanError != nil {
			logging.Error(fmt.Sprintf("Scan error: %s", scanError.Error()))
			continue
		}

		datas = append(datas, data)
	}

	assert.Equal(t, 2, len(datas))
}

func TestQueryRowsContext(t *testing.T) {
	// Arrange
	query := "SELECT Id, Name, Online FROM data"
	db, mock, err := CreateMockDb()

	if err != nil {
		logging.ErrorStackTrace(err)
		return
	}

	base := NewBaseTableInstance()
	base.DB = db

	mock.ExpectQuery(query).
		WillReturnRows(sqlmock.NewRows([]string{"Id", "Name", "Online"}).AddRow(1, "Test", false).AddRow(2, "Test2", true))

	// Act
	rows, err := base.QueryRowsContex(t.Context(), query)

	assert.Nil(t, err)

	defer rows.Close()

	// Assert
	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Errorf("Expectations failed: %s", err)
	}

	type Data struct {
		id     int
		name   string
		online bool
	}

	var datas []Data
	var data Data

	datas = make([]Data, 0)

	for rows.Next() {
		scanError := rows.Scan(&data.id, &data.name, &data.online)

		if scanError != nil {
			logging.Error(fmt.Sprintf("Scan error: %s", scanError.Error()))
			continue
		}

		datas = append(datas, data)
	}

	assert.Equal(t, 2, len(datas))
}
