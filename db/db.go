package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"mcsweeney/content"
	"os"
)

// ContentDB is a very simple sql database abstraction
type ContentDB struct {
	handle *sql.DB
}

const (
	createClipsTable = `
        CREATE TABLE IF NOT EXISTS clips (
            "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
            "url" TEXT 
        );`
	insertClip = `INSERT OR IGNORE INTO clips(url) VALUES (?)`
	existsClip = `SELECT EXISTS(SELECT 1 FROM clips WHERE url=?);`
)

func New(filename string) (*ContentDB, error) {
	// Create new db file if one doesn't exist
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		file.Close()
		if err != nil {
			return nil, fmt.Errorf("Failed to create db %s: %v", filename, err)
		}
	}

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to load %s: %v", filename, err)
	}

	statement, err := db.Prepare(createClipsTable)
	if err != nil {
		return nil, fmt.Errorf("Couldn't prepare SQL statement:\n%s\nerr: %v", createClipsTable, err)
	}
	statement.Exec()

	return &ContentDB{
		handle: db,
	}, nil
}

func (db *ContentDB) Insert(c *content.Clip) error {
	statement, err := db.handle.Prepare(insertClip)
	if err != nil {
		return fmt.Errorf("Couldn't prepare SQL statement:\n%s\nerr: %v", insertClip, err)
	}
	statement.Exec(c.Path)

	return nil
}

func (db *ContentDB) Exists(c *content.Clip) (bool, error) {
	statement, err := db.handle.Prepare(existsClip)
	if err != nil {
		return false, fmt.Errorf("Couldn't prepare SQL statement:\n%s\nerr: %v", existsClip, err)
	}

	var res string
	err = statement.QueryRow(c.Path).Scan(&res)
	if err != nil {
		return false, fmt.Errorf("Couldn't execute exists statement: %v", err)
	}
	if res == "0" {
		return false, nil
	}

	return true, nil
}
