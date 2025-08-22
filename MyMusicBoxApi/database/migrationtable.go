package database

import (
	"fmt"
	"musicboxapi/logging"
)

type IMigrationTable interface {
	Insert(filename string, contents string) (err error)
	ApplyMigration(query string) (err error)
	GetCurrentAppliedMigrationFileName() (fileName string, err error)
}

type MigrationTable struct {
	BaseTable
}

func NewMigrationTableInstance() *MigrationTable {
	return &MigrationTable{
		BaseTable: NewBaseTableInstance(),
	}
}

func (table *MigrationTable) Insert(filename string, contents string) (err error) {
	err = table.NonScalarQuery("INSERT INTO Migration (filename, contents) VALUES($1, $2)", filename, contents)
	if err != nil {
		logging.Error(fmt.Sprintf("Failed to insert new migration: %s", err.Error()))
	}
	return err
}

func (table *MigrationTable) ApplyMigration(query string) (err error) {

	transaction, err := table.DB.Begin()

	if err != nil {
		logging.Error(fmt.Sprintf("Failed to creare transaction %s", err.Error()))
		return
	}

	_, err = transaction.Exec(query)

	if err != nil {
		logging.Error(fmt.Sprintf("Failed to execute query %s", err.Error()))
		return
	}

	return transaction.Commit()
}

func (table *MigrationTable) GetCurrentAppliedMigrationFileName() (fileName string, err error) {
	row := table.QueryRow("SELECT filename FROM migration order by AppliedOn DESC LIMIT 1")
	scanError := row.Scan(&fileName)
	return fileName, scanError
}
