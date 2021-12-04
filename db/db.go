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


func CreateDatabase(name string) error {
    dbName := name + ".sqlite"
    fmt.Println("Creating database %s", dbName)

    file, err := os.Create(dbName)
    file.Close()
    if err != nil {
        return fmt.Errorf("Failed to create db file: %v", err)
    }

    dbFile, err := sql.Open("sqlite3", dbName)
    if err != nil {
        return fmt.Errorf("Failed to load db file '%s': %v", dbName, err)
    }
    defer dbFile.Close()

    err = CreateTwitchTable(dbFile)
    if err != nil {
        return fmt.Errorf("Couldn't create twitch table: %v", err)
    }

    fmt.Println("Database with twitch table created!")
    
    return nil
}


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
