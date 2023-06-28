package service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"rie-kaneko/credit-cards-summary/config"
	"strings"
	"sync"
)

type Service struct {
	Config        *config.Configuration
	Log           *logrus.Logger
	CorrelationID string
}

type ServiceI interface {
	Run() error
}

func NewService(configuration *config.Configuration, logLevel string) *Service {
	return &Service{
		Config: configuration,
		Log:    config.InitLogrus(logLevel),
	}
}

func (s *Service) Run() error {
	var wg sync.WaitGroup

	batches, err := getBatches(s.Config.Environment.CsvPath)
	if err != nil {
		s.Log.Errorf("there was an error getting csv files: %s", err.Error())
		return err
	}

	if len(batches) == 0 {
		s.Log.Infof("no file to process")
		return nil
	}

	numBatches := len(batches) / maxBatchSize
	if len(batches)%maxBatchSize != 0 {
		numBatches++
	}

	for k, v := range batches {
		for i := 0; i < maxBatchSize; i++ {
			if i >= len(batches) {
				break
			}

			if i%maxBatchSize == 0 {
				wg.Add(1)
				go func(k string, v []string) {
					err = s.process(k, v)
					if err != nil {
						s.Log.Errorf("[%s] failed to process: %s", s.CorrelationID, err.Error())
					}
					wg.Done()
				}(k, v)
			}
		}
		wg.Wait()
	}

	return nil
}

func (s *Service) process(k string, v []string) error {
	s.CorrelationID = k

	credit := fmt.Sprintf("%s_%s", k, v[0])
	debit := fmt.Sprintf("%s_%s", k, v[1])

	e, err := s.readCSV(credit, debit)
	if err != nil {
		s.Log.Errorf("[%s] there was an error reading CSV file: %s", s.CorrelationID, err.Error())
		if err = moveFileToNotProcessed(s.Config.Environment.CsvPath, credit); err != nil {
			s.Log.Errorf("[%s] there was an error moving credit csv to not processed archive", s.CorrelationID)
		}
		if err = moveFileToNotProcessed(s.Config.Environment.CsvPath, debit); err != nil {
			s.Log.Errorf("[%s] there was an error moving credit csv to not processed archive", s.CorrelationID)
		}
		return err
	}

	rendered, err := s.render(e)
	if err != nil {
		s.Log.Errorf("[%s] there was an error rendering CSV file information to HTML template: %s", s.CorrelationID, err.Error())

		if err = moveFileToNotProcessed(s.Config.Environment.CsvPath, credit); err != nil {
			s.Log.Errorf("[%s] there was an error moving credit csv to not processed archive", s.CorrelationID)
		}
		if err = moveFileToNotProcessed(s.Config.Environment.CsvPath, debit); err != nil {
			s.Log.Errorf("[%s] there was an error moving debit csv to not processed archive", s.CorrelationID)
		}
		return err
	}

	_ = rendered

	//s.Log.Debugf("id %s html rendered:\n%s", k, rendered)

	_ = moveFileToProcessed(s.Config.Environment.CsvPath, credit)
	_ = moveFileToProcessed(s.Config.Environment.CsvPath, debit)

	s.Log.Debugf("[%s] files moved to processed file", s.CorrelationID)

	return nil
}

func getBatches(csvPath string) (map[string][]string, error) {
	files := make([]string, 0)

	err := filepath.Walk(fmt.Sprintf("%s/%s/", csvPath, toProcessPath), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			files = append(files, info.Name())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	batches := make(map[string][]string)
	for _, b := range files {
		split := strings.Split(b, "_")
		uid := split[0]
		cType := split[1]

		if s, ok := batches[uid]; ok {
			batches[uid] = append(s, cType)
		} else {
			s = make([]string, 0, 2)
			batches[uid] = append(s, cType)
		}
	}

	return batches, nil
}
