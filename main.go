package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"rie-kaneko/credit-cards-summary/config"
	"rie-kaneko/credit-cards-summary/internal/provider"
	"rie-kaneko/credit-cards-summary/internal/service"
)

func main() {
	a, err := provider.NewService(config.Config.AWS)
	if err != nil {
		logrus.Errorf("there was an error: %s", err.Error())
		os.Exit(0)
	}

	s := service.NewService(&config.Config, config.Config.Environment.LogLevel, a)
	s.Log.Infof("Starting...")
	if err = s.Run(); err != nil {
		logrus.Errorf("there was an error: %s", err.Error())
		os.Exit(0)
	}
	s.Log.Infof("Finished...")
}
