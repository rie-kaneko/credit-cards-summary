package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aymerick/raymond"
	"os"
)

type Email struct {
	Name            string
	Balance         float64
	DebitAverage    float64
	CreditAverage   float64
	NumTransactions map[string]int
	Email           string
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
			CcAddresses: []*string{},
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
		Source: aws.String(sender),
	})
	if err != nil {
		return err
	}

	s.Log.Debugf("[%s] email sent, message ID: %s", s.CorrelationID, *result.MessageId)
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
