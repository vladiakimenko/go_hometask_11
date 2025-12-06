package main

import (
	"log"
	"net/http"

	"tasks-api/internal/handlers"
	"tasks-api/internal/storage"
	"tasks-api/internal/api"
)

func main() {
	var store storage.Storage = &storage.InMemoryStorage{}

	h := handlers.New(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", h.TasksCollection) // GET, POST
	mux.HandleFunc("/tasks/", h.TaskItem)       // GET, PUT, DELETE

	executionTimeMux := api.ExecutionTimeMiddleware(mux)

	log.Println("server listening on :8080")
	if err := http.ListenAndServe(":8080", executionTimeMux); err != nil {
		log.Fatal(err)
	}
}
