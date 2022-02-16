package database

import (
	"co2/helpers"
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var instance *sql.DB

// Gets a new instance of the database or returns an already
// existing one if the connection hasn't died yet.
func Get() (*sql.DB, func() error) {
	// As long as the connection hasn't closed
	// return the existing instance, otherwise
	// create a new one if needed.
	if instance != nil && instance.Ping() == nil {
		return instance, instance.Close
	}

	// Try opening the database file
	db, err := sql.Open("sqlite", helpers.DatabaseFile())
	if err != nil {
		log.Fatal(err)
	}

	// Instantiate the schema
	_, err = db.Exec(schema())
	if err != nil {
		log.Printf("%q: %s\n", err, schema())
	}

	// Setup
	instance = db
	return instance, instance.Close
}

// Wrapper for neat handling of any
// errors that might occur.
//
// Nothing nice, if there's an error with the database
// layer, we're screwed so just panic.
func handle(e error) {
	if e != nil {
		panic(e)
	}
}

// Schema definition for the database.
// Simple enough, and built to not complain if it can't
// create anything because it already exists.
//
// This runs every time any of the provided commands are
// run by the user so it needs to be relatively fast.
func schema() string {
	return `
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
		env VARCHAR(64),
		created_at DATETIME default CURRENT_TIMESTAMP
	);
	`
}
