package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var instance *sql.DB

func Get() (*sql.DB, func() error) {
	// As long as the connection hasn't closed
	// return the existing instance, otherwise
	// create a new one if needed.
	if instance != nil && instance.Ping() == nil {
		return instance, instance.Close
	}

	db, err := sql.Open("sqlite3", "./carbon.db")
	if err != nil {
		log.Fatal(err)
	}

	sql := `
		CREATE TABLE IF NOT EXISTS containers (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			uid VARCHAR(64), 
			name VARCHAR(64), 
			compose_file VARCHAR(64),
			created_at DATETIME default CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS stores (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			uid VARCHAR(64),
			path VARCHAR(64),
			created_at DATETIME default CURRENT_TIMESTAMP
		);
		`

	_, err = db.Exec(sql)
	if err != nil {
		log.Printf("%q: %s\n", err, sql)
	}

	instance = db

	return instance, instance.Close
}

func handle(e error) {
	if e != nil {
		panic(e)
	}
}
