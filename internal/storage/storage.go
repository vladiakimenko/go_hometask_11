package storage

import (
	"fmt"
	"time"

	"tasks-api/internal/models"
)

type Storage interface {
	List() []models.Task
	Create(models.Task) (models.Task, error)
	Get(id int) (models.Task, bool)
	Update(id int, task models.Task) (models.Task, error)
	Delete(id int) error
}

// memory
type InMemoryStorage struct {
	tasks  []models.Task
	lastID int
}

func (s *InMemoryStorage) List() []models.Task {
	return s.tasks
}

func (s *InMemoryStorage) Create(task models.Task) (models.Task, error) {

	if task.ID == 0 {
		s.lastID++
		task.ID = s.lastID
	} else {
		for _, t := range s.tasks {
			if t.ID == task.ID {
				return models.Task{}, fmt.Errorf("duplicate ID")
			}
		}
	}

	if task.CreatedAt == "" {
		now := time.Now()
		task.CreatedAt = now.Format(time.RFC3339)
	}

	s.tasks = append(s.tasks, task)
	return task, nil
}

func (s *InMemoryStorage) Get(id int) (models.Task, bool) {
	for _, t := range s.tasks {
		if t.ID == id {
			return t, true
		}
	}
	return models.Task{}, false
}

func (s *InMemoryStorage) Update(id int, task models.Task) (models.Task, error) {
	for _, t := range s.tasks {
		if t.ID == task.ID && t.ID != id {
			return models.Task{}, fmt.Errorf("duplicate ID")
		}
	}
	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks[i] = task
			return task, nil
		}
	}
	return models.Task{}, fmt.Errorf("missing ID")
}

func (s *InMemoryStorage) Delete(id int) error {
	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("missing ID")
}
