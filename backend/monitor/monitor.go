package monitor

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/samirettali/webmonitor/logger"
	"github.com/samirettali/webmonitor/models"
	"github.com/samirettali/webmonitor/notifier"
	"github.com/samirettali/webmonitor/storage"
	"github.com/samirettali/webmonitor/utils"
)

const TIMEOUT = time.Second * 15

type Monitor struct {
	storage  storage.Storage
	notifier notifier.Notifier
	Logger   logger.Logger
	wg       *sync.WaitGroup
	quit     chan struct{}
	sem      chan struct{}
	sync.Mutex
}

func NewMonitor(storage storage.Storage, notifier notifier.Notifier, logger logger.Logger) *Monitor {
	wg := &sync.WaitGroup{}
	quit := make(chan struct{})
	sem := make(chan struct{}, 40)
	return &Monitor{
		storage,
		notifier,
		logger,
		wg,
		quit,
		sem,
		sync.Mutex{},
	}
}

var (
	INTERVALS = []uint64{1, 3, 15, 60, 300, 600}
)

func (m *Monitor) Start() error {
	if err := m.storage.Init(); err != nil {
		return err
	}

	for _, interval := range INTERVALS {
		go m.worker(interval)
	}

	return nil
}

func (m *Monitor) Stop() {
	close(m.quit)
	err := m.storage.Close()
	if err != nil {
		m.Logger.Error(err)
	}
	m.wg.Wait()
}

func (m *Monitor) worker(interval uint64) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	m.Logger.Debugf("worker for interval %d started\n", interval)
	defer m.Logger.Debugf("worker for interval %d stopped\n", interval)

	for {
		select {
		case <-ticker.C:
			err := m.runChecks(interval)
			if err != nil {
				m.Logger.Error(err)
			}
		case <-m.quit:
			return
		}
	}
}

func (m *Monitor) runChecks(interval uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	checks, err := m.storage.GetChecksByInterval(ctx, interval)
	if err != nil {
		werr := errors.Wrap(err, "runChecks can't get checks")
		return werr
		// return []error{werr}
	}

	if len(checks) == 0 {
		return nil
	}

	m.Logger.Debugf("running %d checks for interval %d", len(checks), interval)
	defer m.Logger.Debugf("ended checks for interval %d", interval)

	var wg sync.WaitGroup
	// errChan := make(chan error)
	wg.Add(len(checks))

	for _, check := range checks {
		m.sem <- struct{}{}
		go func(check *models.Check) {
			m.runCheck(check)
			err := m.runCheck(check)
			if err != nil {
				// errChan <- err
				m.Logger.Error(err)
			}
			<-m.sem
			wg.Done()
			// close(errChan)
		}(&check)
	}
	wg.Wait()
	// close(errChan)

	/* errs := make([]error, 0)
	for err := range errChan {
		errs = append(errs, err)
	}
	*/
	/* if len(errs) > 0 {
		return errs
	} */

	return nil
}

func (m *Monitor) runCheck(check *models.Check) error {
	body, err := utils.Request(check.URL)
	if err != nil {
		return err
	}

	if body == check.State {
		return nil
	}

	// check.State = body
	err = m.notifier.Notify(check)
	if err != nil {
		return errors.Wrap(err, "can't sent notification")
	}

	upd := models.CheckUpdate{State: &body}
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	_, err = m.storage.UpdateCheck(ctx, check.ID, &upd)
	if err != nil {
		return errors.Wrap(err, "can't update check")
	}

	return nil
}
