package database

import (
	"fmt"
	"musicboxapi/logging"
)

type IMigrationTable interface {
	Insert(filename string, contents string) (err error)
	ApplyMigration(query string) (err error)
	GetLastAppliedMigrationFileName() (fileName string, err error)
}

type MigrationTable struct {
	BaseTable
}

func NewMigrationTableInstance() *MigrationTable {
	return &MigrationTable{
		BaseTable: NewBaseTableInstance(),
	}
}

func (mt *MigrationTable) Insert(filename string, contents string) (err error) {
	err = mt.NonScalarQuery("INSERT INTO Migration (filename, contents) VALUES($1, $2)", filename, contents)
	if err != nil {
		logging.Error(fmt.Sprintf("Failed to insert new migration: %s", err.Error()))
	}
	return err
}

func (mt *MigrationTable) ApplyMigration(query string) (err error) {
	transaction, err := mt.DB.Begin()

	if err != nil {
		logging.Error(fmt.Sprintf("Failed to begin transaction: %s", err.Error()))
	}

	_, err = transaction.Exec(query)

	if err != nil {
		logging.Error(fmt.Sprintf("Failed to execute migration, rolling back: %s", err.Error()))
		rollbackErr := transaction.Rollback()

		if rollbackErr != nil {
			logging.Error(fmt.Sprintf("Failed to roll back: %s", rollbackErr.Error()))
		}
	}

	err = transaction.Commit()

	if err != nil {
		logging.Error(fmt.Sprintf("Failed to commit migration, rolling back: %s", err.Error()))
		rollbackErr := transaction.Rollback()

		if rollbackErr != nil {
			logging.Error(fmt.Sprintf("Failed to roll back: %s", rollbackErr.Error()))
		}
	}

	return err
}

func (mt *MigrationTable) GetLastAppliedMigrationFileName() (fileName string, err error) {

	query := "SELECT filename FROM migration order by AppliedOn DESC LIMIT 1"

	row := mt.DB.QueryRow(query)

	scanError := row.Scan(&fileName)

	return fileName, scanError
}
