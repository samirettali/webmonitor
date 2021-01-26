package notifier

import (
	"fmt"
	"log"
	"webmonitor/job"

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

func (e *EmailNotifier) Notify(job *job.Job) error {
	text := buildMessage(job)
	subject := fmt.Sprintf("WebMonitor alert: %s", job.URL)
	to := mail.NewEmail(job.Email, job.Email)
	message := mail.NewSingleEmail(e.sender, subject, to, text, "")
	// _, err := e.client.Send(message)
	// return err
	log.Printf("Send notification %s: %+v\n", job.Email, message.Sections)
	return nil
}
