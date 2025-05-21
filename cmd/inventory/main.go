package main

import (
	"ecommerce/internal/config"
	"ecommerce/internal/inventory"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load config")
	}

	if err := inventory.Run(cfg); err != nil {
		logrus.WithError(err).Fatal("Inventory service failed")
	}

	logrus.WithField("addr", cfg.InventoryAddr).Info("Inventory service started successfully")
}
