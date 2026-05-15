package main

import (
	"log"

	"feed_system/internal/router"
)

func main() {
	r := router.New()

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
