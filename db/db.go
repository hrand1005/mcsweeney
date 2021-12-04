package db

import (
    "database/sql"
    "os"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)


//type ContentDatabase interface {
//    Insert() error
//    Delete() error
//    Exists() error
//}


var twitchTableSQL = `CREATE TABLE twitch (
        "clipID" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
        "url" TEXT 
    );`

var insertTwitchClip = `INSERT OR IGNORE INTO twitch(url) VALUES (?)`

func CreateDatabase(name string) (*sql.DB, error) {
    dbName := name + ".sqlite"
    fmt.Println("Creating database %s", dbName)

    file, err := os.Create(dbName)
    file.Close()
    if err != nil {
        return nil, fmt.Errorf("Failed to create db file: %v", err)
    }

    db, err := sql.Open("sqlite3", dbName)
    if err != nil {
        return nil, fmt.Errorf("Failed to load db file '%s': %v", dbName, err)
    }
    fmt.Println("DB created!")

    return db, nil
}

// TODO: should be a method?
func CreateTwitchTable(db *sql.DB) error {
    fmt.Println("Creating twitch table...")
    statement, err := db.Prepare(twitchTableSQL)
    if err != nil {
        return fmt.Errorf("Couldn't prepare SQL statement:\n%s\nerr: %v", twitchTableSQL, err)
    }
    statement.Exec()
    fmt.Println("twitch table created")

    return nil
}

// TODO: should be a method?
func InsertTwitchClip(db * sql.DB, url string) error {
    fmt.Println("Inserting a clip with url: ", url)
    statement, err := db.Prepare(insertTwitchClip)
    if err != nil {
        return fmt.Errorf("Couldn't prepare SQL statement:\n%s\nerr: %v", insertTwitchClip, err)
    }
    
    statement.Exec(url)
    fmt.Println("clip inserted")
    return nil
}
