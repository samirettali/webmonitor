package storage

import (
	"context"
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
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgreStorage) Close() error {
	return s.db.Close()
}

func (s *PostgreStorage) SaveCheck(ctx context.Context, check *models.Check) error {
	query := fmt.Sprintf("INSERT INTO %s (id, name, url, interval, state, email, active) VALUES(:id, :name, :url, :interval, :state, :email, :active)", s.Table)
	_, err := s.db.NamedExecContext(ctx, query, check)
	return err
}

func (s *PostgreStorage) GetChecks(ctx context.Context) ([]models.Check, error) {
	var checks []models.Check
	query := fmt.Sprintf("SELECT * FROM %s", s.Table)
	err := s.db.Select(&checks, query)
	if err != nil {
		return nil, err
	}

	return checks, nil
}

// TODO make this more efficient, use a query builder maybe
func (s *PostgreStorage) UpdateCheck(ctx context.Context, id string, upd *models.CheckUpdate) (models.Check, error) {
	var checks models.Check
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", s.Table)
	err := s.db.GetContext(ctx, &checks, query, id)
	if err != nil {
		return models.Check{}, err
	}

	if upd.Email != nil {
		checks.Email = *upd.Email
	}

	if upd.Interval != nil {
		checks.Interval = *upd.Interval
	}

	if upd.URL != nil {
		checks.URL = *upd.URL
	}

	if upd.State != nil {
		checks.State = *upd.State
	}

	s.Logger.Infof("Updating check %s", checks.ID)

	statement := fmt.Sprintf("UPDATE %s SET name = :name, email = :email, interval = :interval, url = :url, state = :state WHERE id = :id", s.Table)
	_, err = s.db.NamedExecContext(ctx, statement, &checks)
	if err != nil {
		return models.Check{}, err
	}

	return checks, nil
}

func (s *PostgreStorage) GetCheck(ctx context.Context, id string) (models.Check, error) {
	var check models.Check
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", s.Table)
	err := s.db.GetContext(ctx, &check, query, id)
	if err != nil {
		return models.Check{}, err
	}
	return check, nil
}

func (s *PostgreStorage) GetChecksByInterval(ctx context.Context, interval uint64) ([]models.Check, error) {
	var checks []models.Check
	query := fmt.Sprintf("SELECT * FROM %s WHERE interval=$1 and active", s.Table)
	err := s.db.SelectContext(ctx, &checks, query, interval)
	if err != nil {
		return nil, err
	}
	return checks, nil
}

func (s *PostgreStorage) DeleteCheck(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, s.Table)
	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
