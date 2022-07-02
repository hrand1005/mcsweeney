package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nicklaw5/helix"
)

const (
	createClipTable = `
        CREATE TABLE IF NOT EXISTS clips (
            "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
            "url" TEXT 
        );`
	insertClip = `INSERT OR IGNORE INTO clips(url) VALUES (?);`
	selectClip = `SELECT id FROM clips WHERE url=?;`
)

type clipDB struct {
	handle *sql.DB
}

// NewClipDB returns a clipDB from the given sql handle, or error upon failure
func newClipDB(db *sql.DB) (*clipDB, error) {
	_, err := db.Exec(createClipTable)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare SQL statement:\n%s\nerr: %v", createClipTable, err)
	}

	return &clipDB{
		handle: db,
	}, nil
}

// Insert adds the given clip to the database. Returns error upon failure.
func (db *clipDB) Insert(c helix.Clip) error {
	_, err := db.handle.Exec(insertClip, c.URL)
	if err != nil {
		return fmt.Errorf("encountered error inserting clip: %v", err)
	}

	return nil
}

// Exists checks if the given clip exists in the database. Returns true or false.
func (db *clipDB) Exists(c helix.Clip) bool {
	row := db.handle.QueryRow(selectClip, c.URL)
	var result string
	if row.Scan(&result) == sql.ErrNoRows {
		return false
	}
	log.Printf("clip already exists in db:\n%#v\n", c)
	return true
}

// SqliteDB creates a sql.DB handle using the sqlite3 driver and the filename.
// If the file doesn't exist, it is created
func sqliteDB(f string) (*sql.DB, error) {
	// Create new db file if one doesn't exist
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		file, err := os.Create(f)
		file.Close()
		if err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite3", f)
	if err != nil {
		return nil, err
	}

	return db, nil
}
