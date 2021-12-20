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
	createContentTable = `
        CREATE TABLE IF NOT EXISTS content (
            "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
            "url" TEXT 
        );`
	insertContent = `INSERT OR IGNORE INTO content(url) VALUES (?)`
	existsContent = `SELECT EXISTS(SELECT 1 FROM content WHERE url=?);`
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

	statement, err := db.Prepare(createContentTable)
	if err != nil {
		return nil, fmt.Errorf("Couldn't prepare SQL statement:\n%s\nerr: %v", createContentTable, err)
	}
	statement.Exec()

	return &ContentDB{
		handle: db,
	}, nil
}

func (db *ContentDB) Insert(c *content.Content) error {
	statement, err := db.handle.Prepare(insertContent)
	if err != nil {
		return fmt.Errorf("Couldn't prepare SQL statement:\n%s\nerr: %v", insertContent, err)
	}
	statement.Exec(c.Url)

	return nil
}

func (db *ContentDB) Exists(c *content.Content) (bool, error) {
	statement, err := db.handle.Prepare(existsContent)
	if err != nil {
		return false, fmt.Errorf("Couldn't prepare SQL statement:\n%s\nerr: %v", existsContent, err)
	}

	var res string
	err = statement.QueryRow(c.Url).Scan(&res)
	if err != nil {
		return false, fmt.Errorf("Couldn't execute exists statement: %v", err)
	}
	if res == "0" {
		return false, nil
	}

	return true, nil
}
