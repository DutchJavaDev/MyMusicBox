package database

import (
	"context"
	"musicboxapi/models"
)

type ITasklogTable interface {
	GetParentChildLogs(ctx context.Context) ([]models.ParentTaskLog, error)
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
func (table *TasklogTable) GetParentChildLogs(ctx context.Context) ([]models.ParentTaskLog, error) {
	return make([]models.ParentTaskLog, 0), nil
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
