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
			api.WriteError(
				w,
				api.ErrorResponse{
					Status:  http.StatusBadRequest,
					Message: "invalid JSON body",
				},
			)
			return
		}
		result, err := h.Store.Create(parsedBody.ToDTO(nil).(models.Task))
		if err != nil {
			log.Println(err)
			api.WriteError(
				w,
				api.ErrorResponse{
					Status:  http.StatusInternalServerError,
					Message: "failed to dump the task",
				},
			)
			return
		}
		api.SerializeResponse(w, result)
	}
}

// /tasks/{id} (GET, PUT, PATCH, DELETE)
func (h *Handler) TaskItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	taskID, err := strconv.Atoi(path)
	if err != nil {
		api.WriteError(
			w,
			api.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "invalid object ID",
			},
		)
		return
	}
	existingTask, ok := h.Store.Get(taskID)
	if !ok {
		api.WriteError(
			w,
			api.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "requested object does not exist",
			},
		)
		return
	}

	switch r.Method {

	case http.MethodGet:
		api.SerializeResponse(w, existingTask)

	case http.MethodPut:
		parsedBody, ok := api.ParseBody[models.TaskReplaceBody](r)
		if !ok {
			api.WriteError(
				w,
				api.ErrorResponse{
					Status:  http.StatusBadRequest,
					Message: "invalid JSON body",
				},
			)
			return
		}
		result, err := h.Store.Update(taskID, parsedBody.ToDTO(existingTask).(models.Task))
		if err != nil {
			log.Println(err)
			api.WriteError(
				w,
				api.ErrorResponse{
					Status:  http.StatusInternalServerError,
					Message: "failed to replace the task",
				},
			)
			return
		}
		api.SerializeResponse(w, result)

	case http.MethodPatch:
		parsedBody, ok := api.ParseBody[models.TaskUpdateBody](r)
		if !ok {
			api.WriteError(
				w,
				api.ErrorResponse{
					Status:  http.StatusBadRequest,
					Message: "invalid JSON body",
				},
			)
			return
		}
		result, err := h.Store.Update(taskID, parsedBody.ToDTO(existingTask).(models.Task))
		if err != nil {
			log.Println(err)
			api.WriteError(
				w,
				api.ErrorResponse{
					Status:  http.StatusInternalServerError,
					Message: "failed to replace the task",
				},
			)
			return
		}
		api.SerializeResponse(w, result)

	case http.MethodDelete:
		err = h.Store.Delete(taskID)
		if err != nil {
			log.Println(err)
			api.WriteError(
				w,
				api.ErrorResponse{
					Status:  http.StatusInternalServerError,
					Message: "failed to delete the task",
				},
			)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
