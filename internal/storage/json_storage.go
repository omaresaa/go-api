package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/omaresaa/go-api/internal/models"
)

type JSONStorage struct {
	filepath string
	mu       sync.RWMutex
}

func NewJSONStorage(filePath string) *JSONStorage {
	storage := JSONStorage{
		filepath: filePath,
	}

	return &storage
}

func (s *JSONStorage) ReadTasks() ([]models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, err := os.Stat(s.filepath)
	if os.IsNotExist(err) {
		return []models.Task{}, nil
	}

	data, err := os.ReadFile(s.filepath)
	if err != nil {
		return []models.Task{}, err
	}

	if len(data) == 0 {
		return []models.Task{}, nil
	}

	var tasks []models.Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return []models.Task{}, err
	}

	return tasks, nil
}

func (s *JSONStorage) WriteTasks(tasks []models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filepath, data, 0644)
}

func (s *JSONStorage) GetTaskByID(id int) (*models.Task, error) {
	tasks, err := s.ReadTasks()
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		if task.ID == id {
			return &task, nil
		}
	}

	return nil, errors.New("task not found")
}

func (s *JSONStorage) GetNextID() (int, error) {
	tasks, err := s.ReadTasks()
	if err != nil {
		return 0, err
	}

	if len(tasks) == 0 {
		return 1, nil
	}

	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	return maxID + 1, nil
}
