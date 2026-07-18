package cmd

import (
	"log"
	"net/http"

	"kuu/internal/config"
	"kuu/internal/database"
	"kuu/internal/handler"
	"kuu/internal/middleware"
	"kuu/internal/repository"
	"kuu/internal/routes"
	"kuu/internal/websocket"
)

func Server() {
	cfg := config.New()

	db := database.New(cfg.DataBasePath)

	defer db.Database.Close()

	hub := websocket.New()
	go hub.Run()

	repo := repository.New(db)
	h := handler.New(repo, hub)
	m := middleware.New(repo)

	router := routes.Register(h, m)

	log.Printf("Listening on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
