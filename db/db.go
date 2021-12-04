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
var existsTwitchClip = `SELECT EXISTS(SELECT 1 FROM twitch WHERE url=?);`

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


// TODO: should be a method?
func ExistsTwitchClip(db *sql.DB, url string) (bool, error) {
    fmt.Println("Checking if clip with url exists: ", url)
    statement, err := db.Prepare(existsTwitchClip)
    if err != nil {
        return false, fmt.Errorf("Couldn't prepare SQL statement:\n%s\nerr: %v", existsTwitchClip, err)
    }

    var res string
    err = statement.QueryRow(url).Scan(&res)
    if err == sql.ErrNoRows {
        return false, nil
    }
    if err != nil {
        return false, fmt.Errorf("Couldn't execute exists statement: %v", err)
    }

    fmt.Println("Clip already exists in the database.")
    fmt.Println("Result of exists query: %s", res)
    return true, nil
}
    
