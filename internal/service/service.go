package service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"rie-kaneko/credit-cards-summary/config"
)

type Service struct {
	Config *config.Configuration
	Log    *logrus.Logger
}

type ServiceI interface {
	Run() error
}

func (s *Service) Run() error {
	e, err := s.readCSV()
	if err != nil {
		s.Log.Error("there was an error reading CSV file: %s", err.Error())
		return err
	}

	rendered, err := s.render(e)
	if err != nil {
		s.Log.Error("there was an error rendering CSV file information to HTML template", err.Error())
		return err
	}

	fmt.Println(rendered)

	return nil
}
