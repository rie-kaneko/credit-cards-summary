package service

import (
	"github.com/aymerick/raymond"
	"os"
)

type Email struct {
	Name            string
	Balance         float64
	DebitAverage    float64
	CreditAverage   float64
	NumTransactions map[string]int
}

type Sender interface {
	render(e Email) (string, error)
}

func (s *Service) render(e *Email) (string, error) {
	tmpl, err := os.ReadFile(templatePath)
	if err != nil {
		return "", err
	}

	rendered, err := raymond.Render(string(tmpl), &e)
	if err != nil {
		return "", err
	}

	return rendered, nil
}
