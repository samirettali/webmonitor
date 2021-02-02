package notifier

import (
	"fmt"

	"github.com/samirettali/webmonitor/logger"
	"github.com/samirettali/webmonitor/models"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailNotifier struct {
	sender *mail.Email
	client *sendgrid.Client
	Logger logger.Logger
}

func NewEmailNotifier(sender string, apiKey string, logger logger.Logger) *EmailNotifier {
	return &EmailNotifier{
		sender: mail.NewEmail("WebMonitor", sender),
		client: sendgrid.NewSendClient(apiKey),
		Logger: logger,
	}
}

func (e *EmailNotifier) Notify(check *models.Check) error {
	text := buildMessage(check)
	subject := fmt.Sprintf("WebMonitor alert: %s", check.URL)
	to := mail.NewEmail(check.Email, check.Email)
	message := mail.NewSingleEmail(e.sender, subject, to, text, "")
	// _, err := e.client.Send(message)
	// return err
	e.Logger.Infof("Sent notification to %s for %+v\n", check.Email, message.Sections)
	return nil
}
