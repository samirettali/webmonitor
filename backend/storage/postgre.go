package storage

import (
	"context"
	"database/sql/driver"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/samirettali/webmonitor/logger"
	"github.com/samirettali/webmonitor/models"
)

type PostgreStorage struct {
	URI    string
	Table  string
	Logger logger.Logger

	sync.Mutex
	db *sqlx.DB
}

type Duration time.Duration

const TIMEOUT = time.Second * 15

func (s *PostgreStorage) Init() error {
	var err error
	if s.db == nil {
		s.db, err = sqlx.Open("postgres", s.URI)
		if err != nil {
			return err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	if err = s.db.PingContext(ctx); err != nil {
		return errors.Wrap(err, "cannot ping db")
	}

	err = s.initTable()
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgreStorage) initTable() error {
	s.Logger.Infof("creating table %s", s.Table)
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id TEXT PRIMARY KEY NOT NULL,
		name TEXT NOT NULL,
		url TEXT NOT NULL,
		interval BIGINT NOT NULL,
		state TEXT NOT NULL,
		email TEXT NOT NULL,
		active BOOLEAN NOT NULL
	);`, s.Table)

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	if _, err := s.db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}

func (s *PostgreStorage) Close() error {
	return s.db.Close()
}

func (s *PostgreStorage) SaveJob(ctx context.Context, job *models.Job) error {
	query := fmt.Sprintf("INSERT INTO %s (id, name, url, interval, state, email, active) VALUES(:id, :name, :url, :interval, :state, :email, :active)", s.Table)
	_, err := s.db.NamedExecContext(ctx, query, job)
	return err
}

func (s *PostgreStorage) GetJobs(ctx context.Context) ([]models.Job, error) {
	var jobs []models.Job
	query := fmt.Sprintf("SELECT * FROM %s", s.Table)
	err := s.db.Select(&jobs, query)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

// TODO make this more efficient, use a query builder maybe
func (s *PostgreStorage) UpdateJob(ctx context.Context, id string, upd *models.JobUpdate) (models.Job, error) {
	var job models.Job
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", s.Table)
	err := s.db.GetContext(ctx, &job, query, id)
	if err != nil {
		return models.Job{}, err
	}

	if upd.Email != nil {
		job.Email = *upd.Email
	}

	if upd.Interval != nil {
		job.Interval = *upd.Interval
	}

	if upd.URL != nil {
		job.URL = *upd.URL
	}

	if upd.State != nil {
		job.State = *upd.State
	}

	s.Logger.Infof("Updating job %s", job.ID)

	statement := fmt.Sprintf("UPDATE %s SET name = :name, email = :email, interval = :interval, url = :url, state = :state WHERE id = :id", s.Table)
	_, err = s.db.NamedExecContext(ctx, statement, &job)
	if err != nil {
		return models.Job{}, err
	}

	return job, nil
}

func (s *PostgreStorage) GetJob(ctx context.Context, id string) (models.Job, error) {
	var job models.Job
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", s.Table)
	err := s.db.GetContext(ctx, &job, query, id)
	if err != nil {
		return models.Job{}, err
	}
	return job, nil
}

func (s *PostgreStorage) GetJobsByInterval(ctx context.Context, interval uint64) ([]models.Job, error) {
	var jobs []models.Job
	query := fmt.Sprintf("SELECT * FROM %s WHERE interval=$1 and active", s.Table)
	err := s.db.SelectContext(ctx, &jobs, query, interval)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (s *PostgreStorage) DeleteJob(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, s.Table)
	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

// Value converts Duration to a primitive value ready to written to a database.
func (d Duration) Value() (driver.Value, error) {
	return driver.Value(int64(d)), nil
}

// Scan reads a Duration value from database driver type.
func (d *Duration) Scan(raw interface{}) error {
	switch v := raw.(type) {
	case int64:
		*d = Duration(time.Duration(v) * time.Second)
	case nil:
		*d = Duration(0)
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.Duration from: %#v", v)
	}
	return nil
}
