package storage

import (
	"context"
	"sync"

	"github.com/samirettali/webmonitor/models"
)

type Storage interface {
	Init() error
	Close() error
	SaveCheck(ctx context.Context, check *models.Check) error
	GetCheck(ctx context.Context, id string) (models.Check, error)
	GetChecks(ctx context.Context) ([]models.Check, error)
	UpdateCheck(ctx context.Context, id string, upd *models.CheckUpdate) (models.Check, error)
	DeleteCheck(ctx context.Context, id string) error
	// TODO implement get method using a filter struct
	GetChecksByInterval(ctx context.Context, interval uint64) ([]models.Check, error)
}

type MemoryStorage struct {
	checks   []*models.Check
	states map[string]string
	sync.Mutex
}
