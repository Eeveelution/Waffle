package database

import (
	"Waffle/helpers"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var Database *sql.DB

// Initialize initializes the MySQL Database things
func Initialize(username string, password string, location string, dbDatabase string) {
	if Database != nil {
		return
	}

	db, connErr := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, location, dbDatabase))

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if connErr != nil {
		helpers.Logger.Printf("[Database] MySQL Connection could not be established...\n")

		return
	}

	Database = db
}

func Deinitialize() {
	Database.Close()
}
