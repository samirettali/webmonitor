package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/samirettali/webmonitor/logger"
	"github.com/samirettali/webmonitor/models"
	"github.com/samirettali/webmonitor/storage"
	"github.com/samirettali/webmonitor/utils"
)

type StorageHandler struct {
	Storage storage.Storage
	Logger  logger.Logger
}

type Response struct {
	Error string `json:"error"`

}

func (h *StorageHandler) GetCheck(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	params := mux.Vars(r)
	id := params["id"]
	check, err := h.Storage.GetCheck(r.Context(), id)
	if err != nil {
		h.Logger.Errorf("get: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&check)
}

func (h *StorageHandler) GetChecks(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	checks, err := h.Storage.GetChecks(r.Context())
	if err != nil {
		h.Logger.Errorf("get: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&checks)
}

func (h *StorageHandler) CreateCheck(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var check models.Check
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&check)
	if err != nil {
		h.Logger.Error("decode: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	v := validator.New()
	err = v.Struct(check)

	if err != nil {
		h.Logger.Error("validate: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := utils.Request(check.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := Response{
			Error: "The selected URL cannot be reached",
		}
		json.NewEncoder(w).Encode(&resp)
		return
	}

	check.ID = uuid.New().String()

	status := models.Status{
		ID:      uuid.New().String(),
		Content: body,
		CheckID: check.ID,
		Date:    time.Now(),
	}

	err = h.Storage.CreateCheck(r.Context(), &check)
	if err != nil {
		h.Logger.Errorf("save check: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.Storage.UpdateStatus(r.Context(), check.ID, &status)
	if err != nil {
		h.Logger.Errorf("add status: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&check)
}

func (h *StorageHandler) DeleteCheck(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	params := mux.Vars(r)
	id := params["id"]
	err := h.Storage.DeleteCheck(r.Context(), id)

	if err != nil {
		h.Logger.Errorf("delete: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *StorageHandler) UpdateCheck(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var upd models.CheckUpdate
	params := mux.Vars(r)
	id := params["id"]

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&upd)
	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	check, err := h.Storage.UpdateCheck(r.Context(), id, &upd)
	if err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(check)
}

func (h *StorageHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	params := mux.Vars(r)
	id := params["id"]

	statuses, err := h.Storage.GetHistory(r.Context(), id)
	if err != nil {
		h.Logger.Errorf("get history: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&statuses)
}