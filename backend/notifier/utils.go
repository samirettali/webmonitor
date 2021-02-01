package notifier

import (
	"strings"

	"github.com/samirettali/webmonitor/models"
)

func buildMessage(job *models.Job) string {
	b := strings.Builder{}
	b.WriteString("Detected difference on ")
	b.WriteString(job.URL)
	return b.String()
}
