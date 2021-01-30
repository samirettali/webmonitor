package notifier

import (
	"fmt"
	"log"

	"github.com/samirettali/webmonitor/models"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailNotifier struct {
	sender *mail.Email
	client *sendgrid.Client
}

func NewEmailNotifier(sender string, apiKey string) *EmailNotifier {
	return &EmailNotifier{
		sender: mail.NewEmail("WebMonitor", sender),
		client: sendgrid.NewSendClient(apiKey),
	}
}

func (e *EmailNotifier) Notify(job *models.Job) error {
	text := buildMessage(job)
	subject := fmt.Sprintf("WebMonitor alert: %s", job.URL)
	to := mail.NewEmail(job.Email, job.Email)
	message := mail.NewSingleEmail(e.sender, subject, to, text, "")
	// _, err := e.client.Send(message)
	// return err
	log.Printf("Sent notification to %s for %+v\n", job.Email, message.Sections)
	return nil
}
