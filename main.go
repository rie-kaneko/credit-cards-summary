package main

import (
	"os"
	"rie-kaneko/credit-cards-summary/config"
	"rie-kaneko/credit-cards-summary/internal/provider"
	"rie-kaneko/credit-cards-summary/internal/service"
)

func main() {
	a, err := provider.NewService(config.Config.AWS.Region)
	if err != nil {
		os.Exit(0)
	}

	s := service.NewService(&config.Config, config.Config.Environment.LogLevel, a)
	s.Log.Infof("Starting...")
	if err = s.Run(); err != nil {
		os.Exit(0)
	}
	s.Log.Infof("Finished...")
}
