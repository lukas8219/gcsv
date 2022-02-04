package storage

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const databaseFile = "storage.db"

func GetConnection() *sql.Conn {
	db, err := sql.Open("sqlite3", filepath.Join(getDirPath(), databaseFile))
	if err != nil {
		log.Fatalln("Could not open connection", err)
	}

	con, err := db.Conn(context.Background())
	if err != nil {
		log.Fatal("Could not connect", err)
	}
	return con
}

func Migrate() {
	con := GetConnection()
	authenticationTable := `
	CREATE TABLE IF NOT EXISTS authentication(
		type TEXT PRIMARY KEY,
		json TEXT
	);
	`
	_, err := con.ExecContext(context.Background(), authenticationTable)
	if err != nil {
		log.Fatal("Could not execute query", err)
	}

	configTable := `
	CREATE TABLE IF NOT EXISTS config(
		property_key TEXT,
		property_value TEXT
	);
	`
	_, err = con.ExecContext(context.Background(), configTable)
	if err != nil {
		log.Fatal("Could not execute query", configTable, err)
	}

	favoritesTable := `
	CREATE TABLE IF NOT EXISTS favorites(
		id INT PRIMARY KEY,
		name TEXT,
		sheetId TEXT,
		UNIQUE(name)
	);
	`
	_, err = con.ExecContext(context.Background(), favoritesTable)
	if err != nil {
		log.Fatal("Could not execute query", favoritesTable, err)
	}
	return
}
