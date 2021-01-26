package notifier

import (
	"strings"

	"github.com/samirettali/webmonitor/job"
)

func buildMessage(job *job.Job) string {
	b := strings.Builder{}
	b.WriteString("Detected difference on ")
	b.WriteString(job.URL)
	return b.String()
}
