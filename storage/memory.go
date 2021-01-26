package storage

import (
	"sync"

	"github.com/samirettali/webmonitor/job"
)

type Storage interface {
	Save(*job.Job) error
	GetAll() []*job.Job
}

type MemoryStorage struct {
	jobs []*job.Job
	sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	jobs := make([]*job.Job, 0)
	return &MemoryStorage{
		jobs,
		sync.Mutex{},
	}
}

func (s *MemoryStorage) Save(job *job.Job) error {
	s.Lock()
	s.jobs = append(s.jobs, job)
	s.Unlock()
	return nil
}

func (s *MemoryStorage) GetAll() []*job.Job {
	return s.jobs
}
