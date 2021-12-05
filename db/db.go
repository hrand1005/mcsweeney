package db

import (
    "database/sql"
    "os"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

type TwitchStrategy struct {
	dbHandle *sql.DB
}

type ContentDatabase interface {
    Create() error
    Insert(string) error
    Exists(string) (bool, error)
}


var twitchTableSQL = `CREATE TABLE twitch (
        "clipID" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
        "url" TEXT 
    );`

var insertTwitchClip = `INSERT OR IGNORE INTO twitch(url) VALUES (?)`
var existsTwitchClip = `SELECT EXISTS(SELECT 1 FROM twitch WHERE url=?);`
var TWITCH = "twitch"

func ContentDB(source string) (ContentDatabase, error) {
	switch source {
	case TWITCH:
		return &TwitchStrategy{}, nil
	default:
		return nil, fmt.Errorf("Strategy %s not found", source)
	}
}

func (s *TwitchStrategy) Create() error {
    fmt.Println("Creating twitch.sqlite...")

    file, err := os.Create("twitch.sqlite")
    file.Close()
    if err != nil {
        return fmt.Errorf("Failed to create db file: %v", err)
    }

    db, err := sql.Open("sqlite3", "twitch.sqlite")
    if err != nil {
        return fmt.Errorf("Failed to load twitch.sqlite: %v", err)
    }
    fmt.Println("DB created!")
    s.dbHandle = db

    err = s.createTwitchTable()
    if err != nil {
        return fmt.Errorf("Couldn't create twitch table")
    }
    return nil
}


// TODO: should be a method?
func (s *TwitchStrategy) createTwitchTable() error {
    fmt.Println("Creating twitch table...")
    statement, err := s.dbHandle.Prepare(twitchTableSQL)
    if err != nil {
        return fmt.Errorf("Couldn't prepare SQL statement:\n%s\nerr: %v", twitchTableSQL, err)
    }
    statement.Exec()
    fmt.Println("twitch table created")

    return nil
}


// TODO: should be a method?
func (s *TwitchStrategy) Insert(url string) error {
    fmt.Println("Inserting a clip with url: ", url)
    statement, err := s.dbHandle.Prepare(insertTwitchClip)
    if err != nil {
        return fmt.Errorf("Couldn't prepare SQL statement:\n%s\nerr: %v", insertTwitchClip, err)
    }
    
    statement.Exec(url)
    fmt.Println("clip inserted")
    return nil
}


// TODO: should be a method?
func (s *TwitchStrategy) Exists(url string) (bool, error) {
    fmt.Println("Checking if clip with url exists: ", url)
    statement, err := s.dbHandle.Prepare(existsTwitchClip)
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

    fmt.Println("Clip already exists in the database.")
    fmt.Printf("Result of exists query: %s\n", res)
    return true, nil
}
    
