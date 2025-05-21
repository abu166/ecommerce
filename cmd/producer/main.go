package main

import (
	"ecommerce/internal/config"
	"ecommerce/internal/producer"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := producer.Run(cfg); err != nil {
		log.Fatalf("Producer service failed: %v", err)
	}
}
