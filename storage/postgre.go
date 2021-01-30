package storage

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	_ "github.com/lib/pq"
	"github.com/samirettali/webmonitor/logger"
	"github.com/samirettali/webmonitor/models"
)

type PostgreStorage struct {
	URI   string
	Table string
	// Database string
	Logger logger.Logger

	sync.Mutex
	db *sql.DB
}

type Duration time.Duration

const TIMEOUT = time.Second * 15

var ErrPostgreDown = errors.New("PostgreSQL is down")

func (s *PostgreStorage) Init() error {
	var err error
	if s.db == nil {
		s.db, err = sql.Open("postgres", s.URI)
		if err != nil {
			return err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	if err = s.db.PingContext(ctx); err != nil {
		return errors.Wrap(err, "cannot ping db")
		// return ErrPostgreDown
	}

	// err = s.initDatabase()
	// if err != nil {
	// 	return err
	// }

	err = s.initTable()
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgreStorage) initTable() error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id TEXT PRIMARY KEY NOT NULL,
		url TEXT NOT NULL,
		interval BIGINT NOT NULL,
		state TEXT NOT NULL,
		email TEXT
	);`, s.Table)

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	if _, err := s.db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}

// func (s *PostgreStorage) dbExists() (bool, error) {
// 	statement := fmt.Sprintf(`SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = %s);`, s.Database)
// 	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
// 	defer cancel()
// 	row := s.db.QueryRowContext(ctx, statement)
// 	exists := false
// 	err := row.Scan(&exists)
// 	if err != nil {
// 		return false, err
// 	}
// 	return exists, nil
// }

// func (s *PostgreStorage) initDatabase() error {
// 	exists, err := s.dbExists()
// 	if err != nil {
// 		return errors.Wrap(err, "database exists failed")
// 	}

// 	// Cleaner IMO
// 	if exists {
// 		return nil
// 	}

// 	statement := fmt.Sprintf(`CREATE DATABASE %s;`, s.Database)
// 	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
// 	defer cancel()
// 	_, err = s.db.ExecContext(ctx, statement)
// 	return err
// }

func (s *PostgreStorage) Close() error {
	return s.db.Close()
}

func (s *PostgreStorage) SaveJob(job *models.Job) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, url, interval, state, email) VALUES($1, $2, $3, $4, $5)`, s.Table)
	_, err := s.db.Exec(query, job.ID, job.URL, job.Interval, job.State, job.Email)
	return err
}

func (s *PostgreStorage) GetJobs() ([]*models.Job, error) {
	var jobs []*models.Job
	var err error
	query := fmt.Sprintf(`SELECT * FROM %s`, s.Table)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		job := new(models.Job)
		if err := rows.Scan(&job.ID, &job.URL, &job.Interval, &job.State, &job.Email); err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, err
}

func (s *PostgreStorage) UpdateJobState(job *models.Job) error {
	statement := `UPDATE checks SET state = $2 WHERE id = $1`
	_, err := s.db.Exec(statement, job.ID, job.State)
	return err
}

func (s *PostgreStorage) GetJob(id string) (*models.Job, error) {
	var url string
	var interval uint64
	var state string
	var email string
	query := fmt.Sprintf(`SELECT url, interval, state, email FROM %s`, s.Table)
	_, err := s.db.Exec(query, url, interval, state, email)
	if err != nil {
		return nil, err
	}
	job := models.Job{
		ID:       id,
		URL:      url,
		Interval: interval,
		State:    state,
		Email:    email,
	}
	return &job, nil
	// job := new(models.Job)
	// err := s.db.Model(job).Where("id = ?").Select()
	// return job, err
	// return nil, fmt.Errorf("%q: %w", id, ErrNotFound)
}

func (s *PostgreStorage) DeleteJob(id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, s.Table)
	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

// func (s *PostgreStorage) GetJobState(id string) (string, error) {
// 	state, ok := s.states[id]
// 	if !ok {
// 		return "", fmt.Errorf("%q: %w", id, ErrNotFound)
// 	}
// 	return state, nil
// }

// func createSchema(db *pg.DB) error {
// 	models := []interface{}{
// 		(*models.Job)(nil),
// 	}
// 	for _, model := range models {
// 		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
// 			Temp: true,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

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
