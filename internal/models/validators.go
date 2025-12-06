package models

import (
	"log"
	"time"
)

type Validator interface {
	Validate() bool
	ToDTO(existing any) any
}

// create
type TaskCreateBody struct {
	Title string `json:"title"`
}

func (v TaskCreateBody) Validate() bool {
	return validateTitle(v.Title)
}

func (v TaskCreateBody) ToDTO(existing any) any {
	return Task{Title: v.Title}
}

// replace
type TaskReplaceBody struct {
	Title     string  `json:"title"`
	Done      *bool   `json:"done,omitempty"` // using pointers for optional fields
	CreatedAt *string `json:"created_at,omitempty"`
}

func (v TaskReplaceBody) Validate() bool {
	return validateTitle(v.Title) && validateCreatedAt(*v.CreatedAt)
}

func (v TaskReplaceBody) ToDTO(existing any) any {
	existingTask, ok := existing.(Task)
	if !ok {
		panic("the argument provided to TaskReplaceBody.ToDTO is not a Task")
	}
	updated := Task{ID: existingTask.ID}
	updated.Title = v.Title
	if v.Done != nil {
		updated.Done = *v.Done
	}
	if v.CreatedAt != nil {
		updated.CreatedAt = *v.CreatedAt
	}
	return updated
}

// update
type TaskUpdateBody struct {
	Title     *string `json:"title,omitempty"`
	Done      *bool   `json:"done,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
}

func (v TaskUpdateBody) Validate() bool {
	if v.Title != nil && !validateTitle(*v.Title) {
		return false
	}
	if v.CreatedAt != nil && !validateCreatedAt(*v.CreatedAt) {
		return false
	}
	return true
}

func (v TaskUpdateBody) ToDTO(existing any) any {
	existingTask, ok := existing.(Task)
	if !ok {
		panic("the argument provided to TaskUpdateBody.ToDTO is not a Task")
	}
	updated := existingTask
	if v.Title != nil {
		updated.Title = *v.Title
	}
	if v.Done != nil {
		updated.Done = *v.Done
	}
	if v.CreatedAt != nil {
		updated.CreatedAt = *v.CreatedAt
	}
	return updated
}

// constaints
func validateTitle(value string) bool {
	if value == "" {
		log.Println("Title must not be empty")
		return false
	}
	return true
}

func validateCreatedAt(value string) bool {
	_, err := time.Parse(time.RFC3339, value)
	if err != nil {
		log.Println("CreatedAt must be a valid ISO 8601 timestamp")
		return false
	}
	return true
}
