package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

const (
	createTwitchTable = `
        CREATE TABLE twitch (
            "clipID" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
            "url" TEXT 
        );`

	insertTwitchClip = `INSERT OR IGNORE INTO twitch(url) VALUES (?)`

	existsTwitchClip = `SELECT EXISTS(SELECT 1 FROM twitch WHERE url=?);`
)

type TwitchDB struct {
	dbHandle *sql.DB
}

func NewTwitchDB() (*TwitchDB, error) {
	file, err := os.Create("twitch.sqlite")
	file.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to create db file: %v", err)
	}

	db, err := sql.Open("sqlite3", "twitch.sqlite")
	if err != nil {
		return nil, fmt.Errorf("Failed to load twitch.sqlite: %v", err)
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

func (t *TwitchDB) Insert(url string) error {
	statement, err := t.dbHandle.Prepare(insertTwitchClip)
	if err != nil {
		return fmt.Errorf("Couldn't prepare SQL statement:\n%s\nerr: %v", insertTwitchClip, err)
	}
	statement.Exec(url)

	return nil
}

func (t *TwitchDB) Exists(url string) (bool, error) {
	statement, err := t.dbHandle.Prepare(existsTwitchClip)
	if err != nil {
		return false, fmt.Errorf("Couldn't prepare SQL statement:\n%s\nerr: %v", existsTwitchClip, err)
	}

	var res string
	err = statement.QueryRow(url).Scan(&res)
	if err != nil {
		return false, fmt.Errorf("Couldn't execute exists statement: %v", err)
	}
	if res == "0" {
		return false, nil
	}

	return true, nil
}
