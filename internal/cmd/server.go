package cmd

import (
	"kuu/internal/config"
	"kuu/internal/database"
	"kuu/internal/handler"
	"kuu/internal/repository"
)

func Server() {
	cfg := config.New()

	db := database.New(cfg.DataBasePath)

	defer db.Database.Close()

	repo := repository.New(db)
	h := handler.New(repo)

	_ = h // register routes with h
}
