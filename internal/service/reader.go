package service

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"math"
	"os"
	"strconv"
	"time"
)

type Reader interface {
	readCSV() (*Email, error)
}

type transaction struct {
	Id          int    `csv:"id"`
	Date        string `csv:"date"`
	Transaction string `csv:"transaction"`
}

func (s *Service) readCSV() (*Email, error) {
	path := fmt.Sprintf("%s/%s.csv", s.Config.Environment.CsvPath, "txns-credit")
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ts []transaction

	if err = gocsv.UnmarshalFile(f, &ts); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return summaryInformation(ts)
}

func summaryInformation(ts []transaction) (*Email, error) {
	var total float64
	var err error
	numTrans := make(map[string]int)

	for _, t := range ts {
		total, err = getBalance(total, t.Transaction)
		if err != nil {
			return nil, err
		}

		numTrans, err = getNumTransactions(numTrans, t.Date)
		if err != nil {
			return nil, err
		}
	}

	total = math.Round(total*100) / 100
	av := math.Round(total/float64(len(ts))*100) / 100

	return &Email{
		Balance:         total,
		DebitAverage:    0,
		CreditAverage:   av,
		NumTransactions: numTrans,
	}, err
}

func getBalance(total float64, t string) (float64, error) {
	f, err := strconv.ParseFloat(t, 32)
	if err != nil {
		return 0, err
	}
	total += f
	return total, nil
}

func getNumTransactions(numTrans map[string]int, date string) (map[string]int, error) {
	dateTime, err := time.Parse(layout, date)
	if err != nil {
		return nil, err
	}

	m := dateTime.Month().String()
	numTrans[m]++

	return numTrans, nil
}
