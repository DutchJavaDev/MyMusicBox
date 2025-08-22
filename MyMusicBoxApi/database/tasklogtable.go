package database

import (
	"context"
	"musicboxapi/logging"
	"musicboxapi/models"
)

type ITasklogTable interface {
	GetParentLogs(ctx context.Context) ([]models.ParentTaskLog, error)
	GetChildLogs(ctx context.Context, parentId int) ([]models.ChildTaskLog, error)
	CreateParentTaskLog(url string) (models.ParentTaskLog, error)
	CreateChildTaskLog(parent models.ParentTaskLog) (models.ChildTaskLog, error)
	UpdateChildTaskLogStatus(child models.ChildTaskLog) error
	ChildTaskLogDone(child models.ChildTaskLog) error
	ChildTaskLogError(child models.ChildTaskLog) error
}

type TasklogTable struct {
	BaseTable
}

func NewTasklogTableInstance() *TasklogTable {
	return &TasklogTable{
		BaseTable: NewBaseTableInstance(),
	}
}
func (table *TasklogTable) GetParentLogs(ctx context.Context) ([]models.ParentTaskLog, error) {
	query := "SELECT * FROM ParentTaskLog ORDER BY AddTime desc"

	rows, err := table.QueryRowsContex(ctx, query)

	if err != nil {
		return make([]models.ParentTaskLog, 0), nil
	}

	defer rows.Close()

	var parentLog models.ParentTaskLog
	logs := make([]models.ParentTaskLog, 0)

	for rows.Next() {
		err := rows.Scan(&parentLog.Id, &parentLog.Url, &parentLog.AddTime)

		if err != nil {
			logging.ErrorStackTrace(err)
			continue
		}

		logs = append(logs, parentLog)
	}

	return logs, nil
}
func (table *TasklogTable) GetChildLogs(ctx context.Context, parentId int) ([]models.ChildTaskLog, error) {
	query := "SELECT * FROM ChildTaskLog WHERE ParentId = $1 ORDER BY StartTime desc"

	rows, err := table.QueryRowsContex(ctx, query, parentId)

	if err != nil {
		return make([]models.ChildTaskLog, 0), nil
	}

	defer rows.Close()

	var childLog models.ChildTaskLog
	logs := make([]models.ChildTaskLog, 0)

	for rows.Next() {
		err := rows.Scan(&childLog.Id, &childLog.ParentId, &childLog.StartTime, &childLog.EndTime, &childLog.Status, &childLog.OutputLog)

		if err != nil {
			logging.ErrorStackTrace(err)
			continue
		}

		logs = append(logs, childLog)
	}

	return logs, nil
}
func (table *TasklogTable) CreateParentTaskLog(url string) (models.ParentTaskLog, error) {
	query := "INSERT INTO ParentTaskLog (Url) Values($1) RETURNING Id"

	id, err := table.InsertWithReturningId(query, url)

	if err != nil {
		return models.ParentTaskLog{}, err
	}

	return models.ParentTaskLog{
		Id:  id,
		Url: url,
	}, nil
}
func (table *TasklogTable) CreateChildTaskLog(parent models.ParentTaskLog) (models.ChildTaskLog, error) {
	query := "INSERT INTO ChildTaskLog (ParentId, Status) VALUES($1,$2) RETURNING Id"

	defaultStatus := int(models.Pending)

	id, err := table.InsertWithReturningId(query, parent.Id, defaultStatus)

	if err != nil {
		return models.ChildTaskLog{}, err
	}

	return models.ChildTaskLog{
		Id:       id,
		ParentId: parent.Id,
		Status:   defaultStatus,
	}, nil
}
func (table *TasklogTable) UpdateChildTaskLogStatus(child models.ChildTaskLog) error {

	if child.Status == int(models.Downloading) {
		// set the start time to now
		query := "UPDATE ChildTaskLog SET StartTime = CURRENT_TIMESTAMP, Status = $1 WHERE Id = $2"
		return table.NonScalarQuery(query, child.Status, child.Id)
	} else {
		// just update
		query := "UPDATE ChildTaskLog SET Status = $1 WHERE Id = $2"
		return table.NonScalarQuery(query, child.Status, child.Id)
	}
}
func (table *TasklogTable) ChildTaskLogDone(child models.ChildTaskLog) error {
	query := "UPDATE ChildTaskLog SET Status = $1, OutputLog = $2, EndTime = CURRENT_TIMESTAMP WHERE Id = $3"
	return table.NonScalarQuery(query, int(models.Done), child.OutputLog, child.Id)
}
func (table *TasklogTable) ChildTaskLogError(child models.ChildTaskLog) error {
	query := "UPDATE ChildTaskLog SET Status = $1, OutputLog = $2, EndTime = CURRENT_TIMESTAMP WHERE Id = $3"
	return table.NonScalarQuery(query, int(models.Error), child.OutputLog, child.Id)
}
