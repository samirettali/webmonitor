package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/samirettali/webmonitor/logger"
	"github.com/samirettali/webmonitor/models"
	"github.com/samirettali/webmonitor/monitor"
)

type MonitorHandler struct {
	Monitor *monitor.Monitor
	Logger  logger.Logger
}

func (h *MonitorHandler) Get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	checks, err := h.Monitor.GetChecks(r.Context())
	if err != nil {
		h.Logger.Errorf("get: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&checks)
}

func (h *MonitorHandler) Post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var job models.Job
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&job)
	if err != nil {
		h.Logger.Error("decode: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	v := validator.New()
	err = v.Struct(job)

	if err != nil {
		h.Logger.Error("validate: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	check, err := h.Monitor.Add(r.Context(), &job)
	if err != nil {
		h.Logger.Errorf("add: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&check)
}

func (h *MonitorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	params := mux.Vars(r)
	id := params["id"]
	err := h.Monitor.Delete(r.Context(), id)

	if err != nil {
		h.Logger.Errorf("delete: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *MonitorHandler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var upd models.JobUpdate
	params := mux.Vars(r)
	id := params["id"]

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&upd)
	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	job, err := h.Monitor.Update(r.Context(), id, &upd)
	if err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(job)
}
