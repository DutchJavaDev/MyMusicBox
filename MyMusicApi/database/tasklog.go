package database

import (
	"fmt"
	"musicboxapi/logging"
	"musicboxapi/models"
)

func (pdb *PostgresDb) InsertTaskLog() (lastInsertedId int, error error) {
	query := `INSERT INTO TaskLog (Status) VALUES($1) RETURNING Id`

	lastInsertedId, err := pdb.InsertWithReturningId(query, int(models.Pending))

	return lastInsertedId, err
}

func (pdb *PostgresDb) UpdateTaskLogStatus(taskId int, nStatus int) (error error) {
	query := `UPDATE TaskLog SET Status = $1 WHERE Id = $2`

	// could add request context?
	transaction, err := pdb.connection.Begin()

	statement, err := transaction.Prepare(query)
	defer statement.Close()

	if err != nil {
		transaction.Rollback()
		logging.Error(fmt.Sprintf("[UpdateTaskLogStatus] Prepared statement error: %s", err.Error()))
		return err
	}

	_, err = statement.Exec(nStatus, taskId)

	if err != nil {
		logging.Error(fmt.Sprintf("[UpdateTaskLogStatus] Queryrow error: %s", err.Error()))
		transaction.Rollback()
		return err
	}

	err = transaction.Commit()

	if err != nil {
		logging.Error(fmt.Sprintf("[UpdateTaskLogStatus] Transaction commit error: %s", err.Error()))
		transaction.Rollback()
		return err
	}

	return nil
}

// TODO
// func (pdb *PostgresDb) AddLog(log models.Log) {
// }
