package monitor

import (
	"context"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/samirettali/webmonitor/logger"
	"github.com/samirettali/webmonitor/models"
	"github.com/samirettali/webmonitor/notifier"
	"github.com/samirettali/webmonitor/storage"
)

const TIMEOUT = time.Second * 15
const USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; rv:68.0) Gecko/20100101 Firefox/68.0"

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

// func (m *Monitor) recoverJobs() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
// 	defer cancel()
// 	jobs, err := m.storage.GetJobs(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	for _, job := range jobs {
// 		// TODO: This does not scale, group jobs by interval
// 		go m.worker(job.ID)
// 	}
// 	return nil
// }

func (m *Monitor) Start() error {
	if err := m.storage.Init(); err != nil {
		return err
	}

	for _, interval := range INTERVALS {
		go m.worker(interval)
	}

	// if err := m.recoverJobs(); err != nil {
	// 	return err
	// }
	return nil
}

func (m *Monitor) Stop() {
	close(m.quit)
	// err := m.storage.Close()
	m.storage.Close()
	m.wg.Wait()
}

func (m *Monitor) worker(interval uint64) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	m.Logger.Debugf("worker for interval %d started\n", interval)

	// ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	// defer cancel()
	// job, err := m.storage.GetJob(ctx, id)

	// 	// In this case we want to ignore errors as a page may not exist yet
	// 	prevBody, _ := request(job.URL)

	for {
		select {
		case <-ticker.C:
			err := m.runChecks(interval)
			if err != nil {
				m.Logger.Error(err)
			}
		case <-m.quit:
			// m.Logger.Debugf("worker %d stopped\n", job.ID)
			return
		}
	}
}

func (m *Monitor) runChecks(interval uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	checks, err := m.storage.GetJobsByInterval(ctx, interval)
	if err != nil {
		werr := errors.Wrap(err, "runChecks can't get jobs")
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
		go func(check *models.Job) {
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

func (m *Monitor) runCheck(check *models.Job) error {
	body, err := request(check.URL)
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

	upd := models.JobUpdate{State: &body}
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	_, err = m.Update(ctx, check.ID, &upd)
	if err != nil {
		return errors.Wrap(err, "can't update job")
	}

	return nil
}

func (m *Monitor) Add(ctx context.Context, check *models.Job) (*models.Job, error) {
	initialState, err := request(check.URL)
	if err != nil {
		return nil, errors.Wrap(err, "can't fetch page")
	}
	check.State = initialState
	check.ID = uuid.New().String()

	err = m.storage.SaveJob(ctx, check)
	if err != nil {
		return nil, errors.Wrap(err, "can't save check")
	}

	return check, nil
}

func (m *Monitor) Delete(ctx context.Context, id string) error {
	// TODO use only predefined intervals and have workers fetch jobs every time
	return m.storage.DeleteJob(ctx, id)
}

func (m *Monitor) Update(ctx context.Context, id string, upd *models.JobUpdate) (models.Job, error) {
	return m.storage.UpdateJob(ctx, id, upd)
}

func (m *Monitor) GetChecks(ctx context.Context) ([]models.Job, error) {
	return m.storage.GetJobs(ctx)
}
func request(URL string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// log.Println(fmt.Sprintf("Requesting %s", URL))
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", USER_AGENT)

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", nil
	}

	return string(body), nil
}

// func (m *Monitor) addHandler(w http.ResponseWriter, r *http.Request) {
// 	URL := r.URL.Query().Get("url")
// 	if URL == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("Missing URL"))
// 		return
// 	}

// 	intervalString := r.URL.Query().Get("interval")
// 	if URL == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("Missing interval"))
// 		return
// 	}

// 	email := r.URL.Query().Get("email")

// 	if email == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("Missing email"))
// 		return
// 	}

// 	intervalInt, err := strconv.Atoi(intervalString)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("Invalid interval"))
// 		return
// 	}

// 	interval := time.Second * time.Duration(intervalInt)

// 	if err := m.Add(URL, interval, email); err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(err.Error()))
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Ok"))
// }

// func (m *Monitor) startHTTPServer() {
// 	m.wg.Add(1)
// 	defer m.wg.Done()
// 	addr := ":7171"
// 	mux := http.NewServeMux()
// 	server := http.Server{Addr: addr, Handler: mux}

// 	mux.HandleFunc("/add", m.addHandler)

// 	go func() {
// 		err := server.ListenAndServe()
// 		if err != nil && !errors.Is(err, http.ErrServerClosed) {
// 			// s.Logger.Error(err)
// 			log.Println(err)
// 		}
// 	}()

// 	// m.Logger.Infof("Listening on %s", addr)
// 	log.Printf("Listening on %s\n", addr)

// 	<-m.quit
// 	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()
// 	server.Shutdown(tc)
// 	// s.Logger.Debug("HTTP server shutdown")
// 	log.Println("HTTP server shutdown")
// }
