package service

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"math"
	"os"
	"strconv"
	"time"
)

type transaction struct {
	Id          int    `csv:"id"`
	Date        string `csv:"date"`
	Transaction string `csv:"transaction"`
}

func (s *Service) readCSV(credit string, debit string) (*Email, error) {
	s.Log.Debugf("[%s] start to read CSV files", s.CorrelationID)

	pathCredit := fmt.Sprintf("%s/%s/%s", s.Config.Environment.CsvPath, toProcessPath, credit)
	fCredit, err := os.Open(pathCredit)
	if err != nil {
		return nil, err
	}
	defer fCredit.Close()

	pathDebit := fmt.Sprintf("%s/%s/%s", s.Config.Environment.CsvPath, toProcessPath, debit)
	fDebit, err := os.Open(pathDebit)
	if err != nil {
		return nil, err
	}
	defer fDebit.Close()

	var tsCredit, tsDebit []transaction

	if err = gocsv.UnmarshalFile(fCredit, &tsCredit); err != nil {
		return nil, err
	}

	if err = gocsv.UnmarshalFile(fDebit, &tsDebit); err != nil {
		return nil, err
	}

	return s.summaryInformation(tsCredit, tsDebit)
}

func (s *Service) summaryInformation(tsCredit, tsDebit []transaction) (*Email, error) {
	numTrans := make(map[string]int)

	tCredit, avCredit, err := getSummary(tsCredit, numTrans)
	if err != nil {
		return nil, err
	}

	tDebit, avDebit, err := getSummary(tsDebit, numTrans)
	if err != nil {
		return nil, err
	}

	user, err := s.AwsService.GetUser(s.CorrelationID)

	return &Email{
		ID:              user.ID,
		Name:            user.Name,
		Balance:         fmt.Sprintf("%.2f", tCredit+tDebit),
		DebitAverage:    avDebit,
		CreditAverage:   avCredit,
		NumTransactions: numTrans,
		Email:           user.EmailAddress,
	}, err
}

func getSummary(ts []transaction, numTrans map[string]int) (float64, float64, error) {
	var total float64
	var err error
	for _, t := range ts {
		total, err = getBalance(total, t.Transaction)
		if err != nil {
			return 0, 0, err
		}

		numTrans, err = getNumTransactions(numTrans, t.Date)
		if err != nil {
			return 0, 0, err
		}
	}

	total = math.Round(total*100) / 100
	av := math.Round(total/float64(len(ts))*100) / 100

	return total, av, nil
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

func moveFileToProcessed(csvPath, fileName string) error {
	originPath := fmt.Sprintf("%s/%s/%s", csvPath, toProcessPath, fileName)
	destinationPath := fmt.Sprintf("%s/processed/%s", csvPath, fileName)

	err := os.Rename(originPath, destinationPath)
	if err != nil {
		return err
	}

	return nil
}

func moveFileToNotProcessed(csvPath, fileName string) error {
	originPath := fmt.Sprintf("%s/%s/%s", csvPath, toProcessPath, fileName)
	destinationPath := fmt.Sprintf("%s/not_processed/%s", csvPath, fileName)

	err := os.Rename(originPath, destinationPath)
	if err != nil {
		return err
	}

	return nil
}
