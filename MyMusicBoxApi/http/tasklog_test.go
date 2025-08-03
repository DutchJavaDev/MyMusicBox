package http

import (
	"context"
	"encoding/json"
	"errors"
	"musicboxapi/database"
	"musicboxapi/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTaskLogTable struct {
	database.ITasklogTable
	insertTaskLog       func() (lastInsertedId int, err error)
	updateTaskLogStatus func(taskId int, nStatus int) (err error)
	endTaskLog          func(taskId int, nStatus int, data []byte) (err error)
	updateTaskLogError  func(params ...any) (err error)
	getTaskLogs         func(ctx context.Context) ([]models.TaskLog, error)
}

func (h *mockTaskLogTable) InsertTaskLog() (lastInsertedId int, err error) {
	return h.insertTaskLog()
}
func (h *mockTaskLogTable) UpdateTaskLogStatus(taskId int, nStatus int) (err error) {
	return h.updateTaskLogStatus(taskId, nStatus)
}
func (h *mockTaskLogTable) EndTaskLog(taskId int, nStatus int, data []byte) (err error) {
	return h.endTaskLog(taskId, nStatus, data)
}
func (h *mockTaskLogTable) UpdateTaskLogError(params ...any) (err error) {
	return h.updateTaskLogError(params...)
}

func (h *mockTaskLogTable) GetTaskLogs(ctx context.Context) ([]models.TaskLog, error) {
	return h.getTaskLogs(ctx)
}

func TestGetTaskLogsError(t *testing.T) {
	// Arrange
	route := "/tasklogs"
	router := SetupTestRouter()

	mockTable := &mockTaskLogTable{
		getTaskLogs: func(ctx context.Context) (songs []models.TaskLog, error error) {
			return nil, errors.New("Woops")
		},
	}

	tasklogHandler := TaskLogHandler{
		TasklogTable: mockTable,
	}

	router.GET(route, tasklogHandler.FetchTaskLogs)

	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", route, nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestGetTaskLogs(t *testing.T) {
	// Arrange
	route := "/tasklogs"
	router := SetupTestRouter()

	mockTable := &mockTaskLogTable{
		getTaskLogs: func(ctx context.Context) (songs []models.TaskLog, error error) {
			type l struct{}

			ll := l{}

			_output, _ := json.Marshal(ll)
			return []models.TaskLog{
				{Id: 1, Status: int(models.Error), OutputLog: (*json.RawMessage)(&_output)},
			}, nil
		},
	}

	tasklogHandler := TaskLogHandler{
		TasklogTable: mockTable,
	}

	router.GET(route, tasklogHandler.FetchTaskLogs)

	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", route, nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
}
