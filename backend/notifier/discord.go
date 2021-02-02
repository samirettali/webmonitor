package notifier

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/samirettali/webmonitor/models"
)

type Notifier interface {
	Notify(check *models.Check) error
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
