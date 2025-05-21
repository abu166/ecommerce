package main

import (
	"ecommerce/internal/config"
	"ecommerce/internal/user"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := user.Run(cfg); err != nil {
		log.Fatalf("User service failed: %v", err)
	}
}
