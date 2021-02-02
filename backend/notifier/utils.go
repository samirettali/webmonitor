package notifier

import (
	"strings"

	"github.com/samirettali/webmonitor/models"
)

func buildMessage(check *models.Check) string {
	b := strings.Builder{}
	b.WriteString("Detected difference on ")
	b.WriteString(check.URL)
	return b.String()
}
