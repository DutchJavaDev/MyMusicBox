package http

import (
	"context"
	"encoding/json"
	"fmt"
	"musicboxapi/database"
	"musicboxapi/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockTaskLogTable struct {
	database.ITasklogTable
	getParentLogs func(ctx context.Context) ([]models.ParentTaskLog, error)
	getChildLogs  func(ctx context.Context, parentId int) ([]models.ChildTaskLog, error)

	// Rest is used internally....
	createParentTaskLog      func(url string) (models.ParentTaskLog, error)
	createChildTaskLog       func(parent models.ParentTaskLog) (models.ChildTaskLog, error)
	updateChildTaskLogStatus func(child models.ChildTaskLog) error
	childTaskLogDone         func(child models.ChildTaskLog) error
	childTaskLogError        func(child models.ChildTaskLog) error
}

func (h *mockTaskLogTable) GetParentLogs(ctx context.Context) ([]models.ParentTaskLog, error) {
	return h.getParentLogs(ctx)
}
func (h *mockTaskLogTable) GetChildLogs(ctx context.Context, parentId int) ([]models.ChildTaskLog, error) {
	return h.getChildLogs(ctx, parentId)
}

// Rest is used internally....
func (h *mockTaskLogTable) CreateParentTaskLog(url string) (models.ParentTaskLog, error) {
	return h.createParentTaskLog(url)
}
func (h *mockTaskLogTable) CreateChildTaskLog(parent models.ParentTaskLog) (models.ChildTaskLog, error) {
	return h.createChildTaskLog(parent)
}
func (h *mockTaskLogTable) UpdateChildTaskLogStatus(child models.ChildTaskLog) error {
	return h.updateChildTaskLogStatus(child)
}
func (h *mockTaskLogTable) ChildTaskLogDone(child models.ChildTaskLog) error {
	return h.childTaskLogDone(child)
}
func (h *mockTaskLogTable) ChildTaskLogError(child models.ChildTaskLog) error {
	return h.childTaskLogError(child)
}

func TestGetParentLogs(t *testing.T) {
	// Arrange
	route := "/tasklogs"
	router := SetupTestRouter()

	mockTable := &mockTaskLogTable{
		getParentLogs: func(ctx context.Context) ([]models.ParentTaskLog, error) {
			return []models.ParentTaskLog{
				{Id: 0, Url: "www.google.com", AddTime: time.Now()},
				{Id: 1, Url: "www.youtube.com", AddTime: time.Now()},
			}, nil
		},
	}

	taskLogHandler := TaskLogHandler{
		TasklogTable: mockTable,
	}

	router.GET(route, taskLogHandler.FetchParentTaskLogs)

	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", route, nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)

	var rawResult map[string]any

	err := json.Unmarshal(recorder.Body.Bytes(), &rawResult)

	assert.Equal(t, nil, err)

	dataBytes, err := json.Marshal(rawResult["Data"])

	var logs []models.ParentTaskLog
	err = json.Unmarshal(dataBytes, &logs)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(logs))
}
func TestGetChildLogs(t *testing.T) {
	// Arrange
	parentId := 1
	route := "/tasklogs/:parentId"
	router := SetupTestRouter()

	mockTable := &mockTaskLogTable{
		getChildLogs: func(ctx context.Context, parentId int) ([]models.ChildTaskLog, error) {
			childTasks := []models.ChildTaskLog{
				{Id: 0, ParentId: 1, StartTime: time.Now(), EndTime: time.Now().Add(1), Status: 4},
				{Id: 0, ParentId: 1, StartTime: time.Now(), EndTime: time.Now().Add(1), Status: 4},
			}
			return childTasks, nil
		},
	}

	taskLogHandler := TaskLogHandler{
		TasklogTable: mockTable,
	}

	router.GET(route, taskLogHandler.FetchChildTaskLogs)

	recorder := httptest.NewRecorder()

	_route := fmt.Sprintf("/tasklogs/%d", parentId)

	req, _ := http.NewRequest("GET", _route, nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)

	var rawResult map[string]any

	err := json.Unmarshal(recorder.Body.Bytes(), &rawResult)

	assert.Equal(t, nil, err)

	dataBytes, err := json.Marshal(rawResult["Data"])

	var logs []models.ChildTaskLog
	err = json.Unmarshal(dataBytes, &logs)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(logs))
	assert.Equal(t, parentId, logs[0].ParentId)
}
