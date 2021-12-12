package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
    "mcsweeney/content"
	"os"
)

const (
	createTwitchTable = `
        CREATE TABLE IF NOT EXISTS twitch (
            "clipID" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
            "url" TEXT 
        );`

	insertTwitchClip = `INSERT OR IGNORE INTO twitch(url) VALUES (?)`

	existsTwitchClip = `SELECT EXISTS(SELECT 1 FROM twitch WHERE url=?);`
)

type TwitchDB struct {
	dbHandle *sql.DB
}

// TODO: Decide -- enforce new db?
func NewTwitchDB(filename string) (*TwitchDB, error) {
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

	statement, err := db.Prepare(createTwitchTable)
	if err != nil {
		return nil, fmt.Errorf("Couldn't prepare SQL statement:\n%s\nerr: %v", createTwitchTable, err)
	}
	statement.Exec()

	return &TwitchDB{
		dbHandle: db,
	}, nil
}

func (t *TwitchDB) Insert(contentObj *content.ContentObj) error {
	statement, err := t.dbHandle.Prepare(insertTwitchClip)
	if err != nil {
		return fmt.Errorf("Couldn't prepare SQL statement:\n%s\nerr: %v", insertTwitchClip, err)
	}
	statement.Exec(contentObj.Url)

	return nil
}

func (t *TwitchDB) Exists(contentObj *content.ContentObj) (bool, error) {
	statement, err := t.dbHandle.Prepare(existsTwitchClip)
	if err != nil {
		return false, fmt.Errorf("Couldn't prepare SQL statement:\n%s\nerr: %v", existsTwitchClip, err)
	}

	var res string
	err = statement.QueryRow(contentObj.Url).Scan(&res)
	if err != nil {
		return false, fmt.Errorf("Couldn't execute exists statement: %v", err)
	}
	if res == "0" {
		return false, nil
	}

	return true, nil
}
