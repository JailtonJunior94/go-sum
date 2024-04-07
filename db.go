package main

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

func NewSQLServerDatabase(driver, connString string) (*sql.DB, error) {
	sqlDB, err := sql.Open(driver, connString)
	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(200)
	return sqlDB, nil
}
