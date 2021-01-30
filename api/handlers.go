package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/samirettali/webmonitor/models"
	"github.com/samirettali/webmonitor/monitor"
)

type MonitorHandler struct {
	Monitor *monitor.Monitor
}

type Payload struct {
	Result string `json:"result"`
}

func (h *MonitorHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	checks, err := h.Monitor.GetChecks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(checks)
}

func (h *MonitorHandler) Post(w http.ResponseWriter, r *http.Request) {
	// checks, err := h.Monitor.GetChecks()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var c models.Job
	err = json.Unmarshal(body, &c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	check, err := h.Monitor.Add(&c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(check)
}

func (h *MonitorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	params := mux.Vars(r)
	id := params["id"]
	err := h.Monitor.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	p := Payload{"success"}
	json.NewEncoder(w).Encode(&p)
}
