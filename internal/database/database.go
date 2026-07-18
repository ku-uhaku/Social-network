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

	// Run the migrations automatically before returning the DB instance
	runMigrations(db)

	return &DB{
		Database: db,
	}
}

func runMigrations(db *sql.DB) {
	// 1. Create a migration driver using your existing SQLite connection
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal("[MIGRATION DRIVER] : ", err.Error())
	}

	// 2. Point migrate to your migrations directory
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Path to your migrations folder
		"sqlite3",
		driver,
	)
	if err != nil {
		log.Fatal("[MIGRATION INIT] : ", err.Error())
	}

	// 3. Apply the migrations ("up")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("[MIGRATION UP] : ", err.Error())
	}

	log.Println("[MIGRATION] : Database is up to date! Users and Sessions tables verified.")
}
