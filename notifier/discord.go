package notifier

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
	"webmonitor/job"
)

type Notifier interface {
	Notify(job *job.Job) error
}

type DiscordNotifier struct {
}

func (d *DiscordNotifier) Notify(text string, destination string) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	payload := map[string]string{"content": text}

	jsonValue, _ := json.Marshal(payload)
	_, err := client.Post(destination, "application/json", bytes.NewBuffer(jsonValue))
	return err
}
