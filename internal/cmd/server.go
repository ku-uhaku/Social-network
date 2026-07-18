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
	"kuu/internal/service" // 1. Import your brand new service package
	"kuu/internal/websocket"
)

func Server() {
	cfg := config.New()

	db := database.New(cfg.DataBasePath)
	defer db.Database.Close()

	hub := websocket.New()
	go hub.Run()

	// 2. Build Repository (Takes Database Layer)
	repo := repository.New(db)

	// 3. Build Service (Takes Repository Layer)
	svc := service.New(repo)

	// 4. Handlers and Middlewares consume the Service Layer
	h := handler.New(svc, hub)

	m := middleware.New(svc) // ◄--- CHANGE 'repo' TO 'svc' HERE

	router := routes.Register(h, m)

	log.Printf("Listening on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
