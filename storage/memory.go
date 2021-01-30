package storage

import (
	"context"
	"sync"

	"github.com/samirettali/webmonitor/models"
)

type Storage interface {
	Init() error
	Close() error
	SaveJob(ctx context.Context, job *models.Job) error
	GetJob(ctx context.Context, id string) (models.Job, error)
	GetJobs(ctx context.Context) ([]models.Job, error)
	UpdateJob(ctx context.Context, id string, upd *models.JobUpdate) (models.Job, error)
	DeleteJob(ctx context.Context, id string) error
	// GetJobState(id string) (string, error)
}

type MemoryStorage struct {
	jobs   []*models.Job
	states map[string]string
	sync.Mutex
}
