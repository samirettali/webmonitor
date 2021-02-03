package storage

import (
	"context"

	"github.com/samirettali/webmonitor/models"
)

type Storage interface {
	Init() error
	Close() error
	CreateCheck(ctx context.Context, check *models.Check) error
	GetCheck(ctx context.Context, id string) (models.Check, error)
	GetChecks(ctx context.Context) ([]models.Check, error)
	UpdateCheck(ctx context.Context, id string, upd *models.CheckUpdate) (models.Check, error)
	DeleteCheck(ctx context.Context, id string) error
	GetStatus(ctx context.Context, checkID string) (models.Status, error)
	GetHistory(ctx context.Context, checkID string) ([]models.Status, error)
	UpdateStatus(ctx context.Context, checkID string, status *models.Status) error
	// TODO implement get method using a filter struct
	GetChecksByInterval(ctx context.Context, interval uint64) ([]models.Check, error)
}
