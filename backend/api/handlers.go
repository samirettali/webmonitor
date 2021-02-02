package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/samirettali/webmonitor/logger"
	"github.com/samirettali/webmonitor/models"
	"github.com/samirettali/webmonitor/storage"
	"github.com/samirettali/webmonitor/utils"
)

type ChecksHandler struct {
	Storage storage.Storage
	Logger  logger.Logger
}

type Response struct {
	Error string `json:"error"`
}

func (h *ChecksHandler) Get(w http.ResponseWriter, r *http.Request) {
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

func (h *ChecksHandler) Post(w http.ResponseWriter, r *http.Request) {
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

	initialState, err := utils.Request(check.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := Response{
			Error: "The selected URL cannot be reached",
		}
		json.NewEncoder(w).Encode(&resp)
		return
	}

	check.State = initialState
	check.ID = uuid.New().String()

	err = h.Storage.SaveCheck(r.Context(), &check)
	if err != nil {
		h.Logger.Errorf("add: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&check)
}

func (h *ChecksHandler) Delete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	params := mux.Vars(r)
	id := params["id"]
	err := h.Storage.DeleteCheck(r.Context(), id)

	if err != nil {
		h.Logger.Errorf("delete: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ChecksHandler) Update(w http.ResponseWriter, r *http.Request) {
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
