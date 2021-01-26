package monitor

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/samirettali/webmonitor/job"
	"github.com/samirettali/webmonitor/notifier"
	"github.com/samirettali/webmonitor/storage"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; rv:68.0) Gecko/20100101 Firefox/68.0"

type Monitor struct {
	storage  storage.Storage
	notifier notifier.Notifier
	wg       *sync.WaitGroup
	quit     chan struct{}
	nextID   int
	sync.Mutex
}

func NewMonitor(storage storage.Storage, notifier notifier.Notifier) *Monitor {
	wg := &sync.WaitGroup{}
	quit := make(chan struct{})
	nextID := 1
	return &Monitor{
		storage,
		notifier,
		wg,
		quit,
		nextID,
		sync.Mutex{},
	}
}

func (m *Monitor) recoverJobs() {
	for _, job := range m.storage.GetAll() {
		go m.worker(job)
	}
}

func (m *Monitor) Start() {
	m.recoverJobs()
	go m.startHTTPServer()
}

func (m *Monitor) Stop() {
	close(m.quit)
	m.wg.Wait()
}

func (m *Monitor) Add(job *job.Job) error {
	err := m.storage.Save(job)
	if err != nil {
		return err
	}

	go m.worker(job)
	return nil
}

func (m *Monitor) worker(job *job.Job) {
	m.wg.Add(1)
	defer m.wg.Done()
	m.Lock()
	id := m.nextID
	m.nextID++
	m.Unlock()

	ticker := time.NewTicker(job.Interval)

	// In this case we want to ignore errors as a page may not exist yet
	prevBody, _ := request(job.URL)

	log.Printf("worker %d started\n", id)
	for {
		select {
		case <-ticker.C:
			body, err := request(job.URL)
			if err != nil {
				log.Printf("Error while requesting %s\n", job.URL)
				continue
			}
			if body != prevBody {
				log.Printf("Body of %s changed\n", job.URL)
				if err := m.notifier.Notify(job); err != nil {
					log.Println("can't send notification", err)
				}
				prevBody = body
			}
		case <-m.quit:
			log.Printf("worker %d stopped\n", id)
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

	req.Header.Set("User-Agent", userAgent)

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

func (m *Monitor) addHandler(w http.ResponseWriter, r *http.Request) {
	URL := r.URL.Query().Get("url")
	if URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing URL"))
		return
	}

	intervalString := r.URL.Query().Get("interval")
	if URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing interval"))
		return
	}

	email := r.URL.Query().Get("email")

	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing email"))
		return
	}

	intervalInt, err := strconv.Atoi(intervalString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid interval"))
		return
	}

	interval := time.Second * time.Duration(intervalInt)

	job := job.Job{
		URL:      URL,
		Interval: interval,
		Email:    email,
	}

	m.Add(&job)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
}

func (m *Monitor) startHTTPServer() {
	m.wg.Add(1)
	defer m.wg.Done()
	addr := ":7171"
	mux := http.NewServeMux()
	server := http.Server{Addr: addr, Handler: mux}

	mux.HandleFunc("/add", m.addHandler)

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			// s.Logger.Error(err)
			log.Println(err)
		}
	}()

	// m.Logger.Infof("Listening on %s", addr)
	log.Printf("Listening on %s\n", addr)

	<-m.quit
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(tc)
	// s.Logger.Debug("HTTP server shutdown")
	log.Println("HTTP server shutdown")
}
