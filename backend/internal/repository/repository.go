package repository

import "kuu/internal/database"

type Repository struct {
	DB *database.DB
}

func New(db *database.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
