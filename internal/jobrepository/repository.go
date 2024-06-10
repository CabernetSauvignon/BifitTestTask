package jobrepository

import (
	"fmt"
	"sync"
)

// JobRepository представляет логику доступа к данным задач из источника данных
type JobRepository interface {
	Get(name string) (Job, error)
	GetAll() ([]Job, error)
	Add(job Job) error
	Update(job Job) error
}

// Реализация хранилища задач
type jobRepository struct {
	jobMap map[string]Job
	mutex  sync.RWMutex
}

// NewJobRepository создаёт новый экземпляр хранилища данных
func NewJobRepository() JobRepository {
	return &jobRepository{make(map[string]Job), sync.RWMutex{}}
}

// Get возвращает задачу по её имени
func (jr *jobRepository) Get(name string) (Job, error) {
	jr.mutex.RLock()
	defer jr.mutex.RUnlock()

	job, ok := jr.jobMap[name]
	if !ok {
		return job, fmt.Errorf("job not found: %s", name)
	}
	return job, nil
}

// GetAll возвращает все хранимые задачи
func (jr *jobRepository) GetAll() ([]Job, error) {
	jr.mutex.RLock()
	defer jr.mutex.RUnlock()

	var jobs []Job
	for _, job := range jr.jobMap {
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// Add добавляет новую задачу в хранилище
func (jr *jobRepository) Add(job Job) error {
	jr.mutex.Lock()
	defer jr.mutex.Unlock()

	if _, ok := jr.jobMap[job.Name]; ok {
		return fmt.Errorf("job already exists: %s", job.Name)
	}

	jr.jobMap[job.Name] = job
	return nil
}

// Update изменяет уже существующую задачу в хранилище
func (jr *jobRepository) Update(job Job) error {
	jr.mutex.Lock()
	defer jr.mutex.Unlock()

	if _, ok := jr.jobMap[job.Name]; !ok {
		return fmt.Errorf("job not found: %s", job.Name)
	}

	jr.jobMap[job.Name] = job
	return nil
}
