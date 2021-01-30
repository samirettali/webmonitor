package monitor

import (
	"context"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
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
	nextID   int
	sync.Mutex
}

func NewMonitor(storage storage.Storage, notifier notifier.Notifier, logger logger.Logger) *Monitor {
	wg := &sync.WaitGroup{}
	quit := make(chan struct{})
	nextID := 1
	return &Monitor{
		storage,
		notifier,
		logger,
		wg,
		quit,
		nextID,
		sync.Mutex{},
	}
}

func (m *Monitor) recoverJobs() error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	jobs, err := m.storage.GetJobs(ctx)
	if err != nil {
		return err
	}
	for _, job := range jobs {
		// TODO: This does not scale, group jobs by interval
		go m.worker(&job)
	}
	return nil
}

func (m *Monitor) Start() error {
	if err := m.storage.Init(); err != nil {
		return err
	}

	if err := m.recoverJobs(); err != nil {
		return err
	}
	return nil
}

func (m *Monitor) Stop() {
	close(m.quit)
	// err := m.storage.Close()
	m.storage.Close()
	m.wg.Wait()
}

func (m *Monitor) Add(ctx context.Context, check *models.Job) (*models.Job, error) {
	initialState, _ := request(check.URL)
	check.State = initialState
	check.ID = uuid.New().String()

	err := m.storage.SaveJob(ctx, check)
	if err != nil {
		return nil, err
	}

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.worker(check)
	}()
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

func (m *Monitor) worker(job *models.Job) {
	m.Lock()
	id := m.nextID
	m.nextID++
	m.Unlock()

	interval := time.Duration(job.Interval)
	ticker := time.NewTicker(interval * time.Second)

	// 	// In this case we want to ignore errors as a page may not exist yet
	// 	prevBody, _ := request(job.URL)

	m.Logger.Debugf("worker %d started on %s with interval %d\n", id, job.URL, job.Interval)
	for {
		select {
		case <-ticker.C:
			body, err := request(job.URL)
			if err != nil {
				// log.Printf("Error while requesting %s\n", job.URL)
				continue
			}
			if body != job.State {
				m.Logger.Infof("Body of %s changed\n", job.URL)
				if err := m.notifier.Notify(job); err != nil {
					m.Logger.Error("can't send notification", err)
				}
				upd := models.JobUpdate{State: &body}
				ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
				m.storage.UpdateJob(ctx, job.ID, &upd)
				cancel()
			}
		case <-m.quit:
			m.Logger.Debugf("worker %d stopped\n", id)
			return
		}
	}
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
