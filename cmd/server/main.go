package main

import (
	"log"

	_ "feed_system/docs"
	"feed_system/internal/bootstrap"
	"feed_system/internal/router"
)

// @title Feed System API
// @version 1.0
// @description Feed System API documentation.
// @host localhost:8080
// @BasePath /
func main() {
	app := bootstrap.New()
	r := router.New(app)

	if err := r.Run(app.Config.Server.Addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
