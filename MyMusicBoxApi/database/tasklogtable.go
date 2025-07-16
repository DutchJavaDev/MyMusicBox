package database

import (
	"context"
	"fmt"
	"musicboxapi/logging"
	"musicboxapi/models"
	"time"
)

type ITasklogTable interface {
	InsertTaskLog() (lastInsertedId int, err error)
	UpdateTaskLogStatus(taskId int, nStatus int) (err error)
	EndTaskLog(taskId int, nStatus int, data []byte) (err error)
	UpdateTaskLogError(params ...any) (err error)
	GetTaskLogs(ctx context.Context) ([]models.TaskLog, error)
}

type TasklogTable struct {
	BaseTable
}

func NewTasklogTableInstance() *TasklogTable {
	return &TasklogTable{
		BaseTable: NewBaseTableInstance(),
	}
}

func (tt *TasklogTable) InsertTaskLog() (lastInsertedId int, error error) {
	query := `INSERT INTO TaskLog (Status) VALUES($1) RETURNING Id`

	lastInsertedId, err := tt.InsertWithReturningId(query, int(models.Pending))

	return lastInsertedId, err
}

func (tt *TasklogTable) UpdateTaskLogStatus(taskId int, nStatus int) (error error) {
	query := `UPDATE TaskLog SET Status = $1 WHERE Id = $2`

	return tt.NonScalarQuery(query, nStatus, taskId)
}

func (tt *TasklogTable) EndTaskLog(taskId int, nStatus int, data []byte) error {
	query := `UPDATE TaskLog SET Status = $1, OutputLog = $2, EndTime = $3 WHERE Id = $4`

	return tt.NonScalarQuery(query, nStatus, data, time.Now(), taskId)
}

func (tt *TasklogTable) UpdateTaskLogError(params ...any) error {
	query := `UPDATE TaskLog
		             SET Status = $1, OutputLog = $2, EndTime = $3
		             WHERE Id = $4`
	return tt.NonScalarQuery(query, params...)
}

func (tt *TasklogTable) GetTaskLogs(ctx context.Context) ([]models.TaskLog, error) {
	query := `SELECT Id, StartTime, EndTime, Status, OutputLog FROM TaskLog ORDER BY Id desc` // get the latest first

	rows, err := tt.DB.QueryContext(ctx, query)

	if err != nil {
		logging.Error(fmt.Sprintf("QueryRow error: %s", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var tasklog models.TaskLog

	tasks := make([]models.TaskLog, 0)

	for rows.Next() {
		scanError := rows.Scan(&tasklog.Id, &tasklog.StartTime, &tasklog.EndTime, &tasklog.Status, &tasklog.OutputLog)

		if scanError != nil {
			logging.Error(fmt.Sprintf("Scan error: %s", scanError.Error()))
			continue
		}

		tasks = append(tasks, tasklog)
	}

	return tasks, nil
}
