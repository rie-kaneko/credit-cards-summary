package main

import (
	"os"
	"rie-kaneko/credit-cards-summary/config"
	"rie-kaneko/credit-cards-summary/internal/service"
)

var s *service.Service

func init() {
	s = &service.Service{
		Config: &config.Config,
		Log:    config.InitLogrus(config.Config.Environment.LogLevel),
	}
}

func main() {
	if err := s.Run(); err != nil {
		os.Exit(0)
	}
}
