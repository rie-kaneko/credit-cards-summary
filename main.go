package main

import (
	"os"
	"rie-kaneko/credit-cards-summary/config"
	"rie-kaneko/credit-cards-summary/internal/service"
)

func main() {
	s := service.NewService(&config.Config, config.Config.Environment.LogLevel)
	s.Log.Infof("Starting...")
	if err := s.Run(); err != nil {
		os.Exit(0)
	}
	s.Log.Infof("Finished...")
}
