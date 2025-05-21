package main

import (
	"ecommerce/internal/config"
	"ecommerce/internal/order"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load config")
	}

	if err := order.Run(cfg); err != nil {
		logrus.WithError(err).Fatal("Order service failed")
	}

	logrus.WithField("addr", cfg.OrderAddr).Info("Order service started successfully")
}
