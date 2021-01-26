package job

import (
	"time"
)

type Job struct {
	URL      string
	Interval time.Duration
	// Notifiers []notifier.Notifier
	// Not sure about this
	Email          string
	DiscordWebhook string
}
