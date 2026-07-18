package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Database *sql.DB
}

func New(path string) *DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal("[DATABASE] : ", err.Error())
	}

	if err := db.Ping(); err != nil {
		db.Close()
		log.Fatal("[DATABASE] : ", err.Error())
	}

	return &DB{
		Database: db,
	}
}
