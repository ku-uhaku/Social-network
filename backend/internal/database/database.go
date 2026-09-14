package database

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	runMigrations(db)

	return &DB{
		Database: db,
	}
}

func runMigrations(db *sql.DB) {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal("[MIGRATION DRIVER] : ", err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite3",
		driver,
	)
	if err != nil {
		log.Fatal("[MIGRATION INIT] : ", err.Error())
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("[MIGRATION UP] : ", err.Error())
	}

	log.Println("[MIGRATION] : Database is UP!")
}
