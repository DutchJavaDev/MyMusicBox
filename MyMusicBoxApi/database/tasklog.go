package database

import (
	"context"
	"fmt"
	"musicboxapi/logging"
	"musicboxapi/models"
	"time"
)

func (pdb *PostgresDb) InsertTaskLog() (lastInsertedId int, error error) {
	query := `INSERT INTO TaskLog (Status) VALUES($1) RETURNING Id`

	lastInsertedId, err := pdb.InsertWithReturningId(query, int(models.Pending))

	return lastInsertedId, err
}

func (pdb *PostgresDb) UpdateTaskLogStatus(taskId int, nStatus int) (error error) {
	query := `UPDATE TaskLog SET Status = $1 WHERE Id = $2`

	return pdb.NonScalarQuery(query, nStatus, taskId)
}

func (pdb *PostgresDb) EndTaskLog(taskId int, nStatus int, data []byte) error {
	query := `UPDATE TaskLog SET Status = $1, OutputLog = $2, EndTime = $3 WHERE Id = $4`

	return pdb.NonScalarQuery(query, nStatus, data, time.Now(), taskId)
}

func (pdb *PostgresDb) UpdateTaskLogError(params ...any) error {
	query := `UPDATE TaskLog
		             SET Status = $1, OutputLog = $2, EndTime = $3
		             WHERE Id = $4`
	return pdb.NonScalarQuery(query, params...)
}

func (pdb *PostgresDb) GetTaskLogs(ctx context.Context) ([]models.TaskLog, error) {
	query := `SELECT Id, StartTime, EndTime, Status, OutputLog FROM TaskLog ORDER BY Id desc` // get the latest first

	rows, err := pdb.connection.QueryContext(ctx, query)
	defer rows.Close()

	if err != nil {
		logging.Error(fmt.Sprintf("[GetTaskLogs] QueryRow error: %s", err.Error()))
		return nil, err
	}

	var tasklog models.TaskLog

	tasks := make([]models.TaskLog, 0)

	for rows.Next() {
		scanError := rows.Scan(&tasklog.Id, &tasklog.StartTime, &tasklog.EndTime, &tasklog.Status, &tasklog.OutputLog)

		if scanError != nil {
			logging.Error(fmt.Sprintf("[GetTaskLogs] Scan error: %s", scanError.Error()))
			continue
		}

		tasks = append(tasks, tasklog)
	}

	return tasks, nil
}
