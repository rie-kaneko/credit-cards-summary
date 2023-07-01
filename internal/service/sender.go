package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aymerick/raymond"
	"os"
	"rie-kaneko/credit-cards-summary/internal/provider"
	"time"
)

type Email struct {
	ID              string
	Name            string
	Balance         string
	DebitAverage    float64
	CreditAverage   float64
	NumTransactions []NumTransactions
	Email           string
}

type NumTransactions struct {
	Month string
	Count int
}

func (s *Service) sendMail(e *Email) error {
	s.Log.Debugf("[%s] start to send mail", s.CorrelationID)

	rendered, err := render(e)
	if err != nil {
		return err
	}
	s.Log.Debugf(rendered)

	result, err := s.AwsService.SesClient.SendEmail(&ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(e.Email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(rendered),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(s.Config.Environment.EmailSender),
	})
	if err != nil {
		if result != nil {
			err = s.AwsService.PostTransaction(provider.Transaction{
				ID:         *result.MessageId,
				Sender:     s.Config.Environment.EmailSender,
				ReceiverID: e.ID,
				Receiver:   e.Email,
				Time:       time.Now(),
				Status:     provider.Failed,
			})
			if err != nil {
				s.Log.Warnf("transaction not registered: %s", err.Error())
			}
		}
		return err
	}

	s.Log.Debugf("[%s] email sent, message ID: %s", s.CorrelationID, *result.MessageId)

	err = s.AwsService.PostTransaction(provider.Transaction{
		ID:         *result.MessageId,
		Sender:     s.Config.Environment.EmailSender,
		ReceiverID: e.ID,
		Receiver:   e.Email,
		Time:       time.Now(),
		Status:     provider.Succeeded,
	})
	if err != nil {
		s.Log.Warnf("transaction not registered: %s", err.Error())
	}

	return nil
}

func render(e *Email) (string, error) {
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
