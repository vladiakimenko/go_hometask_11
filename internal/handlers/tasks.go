package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"tasks-api/internal/api"
	"tasks-api/internal/models"
	"tasks-api/internal/storage"
)

type Handler struct{ Store storage.Storage }

func New(s storage.Storage) *Handler { return &Handler{Store: s} }

// /tasks (GET, POST)
func (h *Handler) TasksCollection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {

	case http.MethodGet:
		// no explicit status header required, default is 200
		result := h.Store.List()
		api.SerializeResponse(w, result)

	case http.MethodPost:
		w.WriteHeader(http.StatusCreated)
		parsedBody, ok := api.ParseBody[models.TaskCreateBody](r)
		if !ok {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
		result, err := h.Store.Create(parsedBody.ToDTO(nil).(models.Task))
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to dump the task", http.StatusInternalServerError)
			return
		}
		api.SerializeResponse(w, result)
	}
}

// /tasks/{id} (GET, PUT, PATCH, DELETE)
func (h *Handler) TaskItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 || parts[0] != "tasks" { // do I really need the seond part?
		http.NotFound(w, r)
		return
	}
	taskID, err := strconv.Atoi(parts[1])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	existingTask, ok := h.Store.Get(taskID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	switch r.Method {

	case http.MethodGet:
		api.SerializeResponse(w, existingTask)

	case http.MethodPut:
		parsedBody, ok := api.ParseBody[models.TaskReplaceBody](r)
		if !ok {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
		result, err := h.Store.Update(taskID, parsedBody.ToDTO(existingTask).(models.Task))
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to replace the task", http.StatusInternalServerError)
			return
		}
		api.SerializeResponse(w, result)

	case http.MethodPatch:
		parsedBody, ok := api.ParseBody[models.TaskUpdateBody](r)
		if !ok {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
		result, err := h.Store.Update(taskID, parsedBody.ToDTO(existingTask).(models.Task))
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to replace the task", http.StatusInternalServerError)
			return
		}
		api.SerializeResponse(w, result)

	case http.MethodDelete:
		err = h.Store.Delete(taskID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to delete the task", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
