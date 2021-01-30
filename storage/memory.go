package storage

import (
	"sync"
	"time"

	"github.com/samirettali/webmonitor/models"
)

type Storage interface {
	Init() error
	SaveJob(*models.Job) error
	GetJob(id string) (*models.Job, error)
	GetJobs() ([]*models.Job, error)
	UpdateJobState(*models.Job) error
	DeleteJob(id string) error
	Close() error
	// GetJobState(id string) (string, error)
}

// // DialFilter represents a filter used by FindDials().
// type DialFilter struct {
// 	// Filtering fields.
// 	ID         *int    `json:"id"`
// 	InviteCode *string `json:"inviteCode"`

// // Restrict to subset of range.
// Offset int `json:"offset"`
// Limit  int `json:"limit"`
// }

// type DialService interface { FindDialByID(ctx context.Context, id int) (*Dial, error)
// 	FindDials(ctx context.Context, filter DialFilter) ([]*Dial, int, error)
// 	CreateDial(ctx context.Context, dial *Dial) error
// 	UpdateDial(ctx context.Context, id int, upd DialUpdate) (*Dial, error)
// 	DeleteDial(ctx context.Context, id int) error
// }

type UpdateJob struct {
	URL      *string
	Interval *time.Duration
	State    *string
	Email    *string
}

// type Storage interface {
// 	Init() error
// 	FindJobById(ctx context.Context, id string) (*models.Job, error)
// 	FindJobs(ctx context.Context, filter JobFilter) ([]*models.Job, error) // ok double nil
// 	CreateJob(ctx contex.COntext, *models.Job) error
// 	FindJobs(ctx context.Context) ([]*models.Job, error)
// 	UpdateJob(ctx context.Context, upd UpdateJob) (*models.Job, error)
// 	UpdateJobs(ctx context.Context, ids []string *models.Job) ([]*job.Job, error)
//  DeleteJobs(ctx context.Context, ids []string)
// 	DeleteJob(ctx context.Context, id string) error
// 	Close() error
// 	// GetJobState(id string) (string, error)
// }

type MemoryStorage struct {
	jobs   []*models.Job
	states map[string]string
	sync.Mutex
}

// var ErrNotFound = errors.New("job not found")

// func NewMemoryStorage() *MemoryStorage {
// 	jobs := make([]*models.Job, 0)
// 	states := make(map[string]string)
// 	return &MemoryStorage{
// 		jobs,
// 		states,
// 		sync.Mutex{},
// 	}
// }

// func (s *MemoryStorage) SaveJob(job *models.Job) error {
// 	s.Lock()
// 	s.jobs = append(s.jobs, job)
// 	s.Unlock()
// 	return nil
// }

// func (s *MemoryStorage) GetJobs() ([]*models.Job, error) {
// 	return s.jobs, nil
// }

// func (s *MemoryStorage) SaveJobState(id string, state string) error {
// 	s.states[id] = state
// 	return nil
// }

// func (s *MemoryStorage) GetJob(id string) (*models.Job, error) {
// 	for _, job := range s.jobs {
// 		if job.ID == id {
// 			return job, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("%q: %w", id, ErrNotFound)
// }

// func (s *MemoryStorage) GetJobState(id string) (string, error) {
// 	state, ok := s.states[id]
// 	if !ok {
// 		return "", fmt.Errorf("%q: %w", id, ErrNotFound)
// 	}
// 	return state, nil
// }

// func (s *MemoryStorage) DeleteJob(id string) error {
// 	_, ok := s.states[id]
// 	if !ok {
// 		return fmt.Errorf("%q: %w", id, ErrNotFound)
// 	}
// 	for _, job :range s.jobs {
// 		if s.jobs.ID == id {

// 		}
// 	}
// 	delete(s.states, id)
// 	return nil
// }
