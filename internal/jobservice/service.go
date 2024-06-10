package jobservice

import (
	"TestTask/internal/jobrepository"
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	updateJobStatusInterval = time.Minute * 5
)

// JobService представляет логику работы с задачами
type JobService interface {
	GetAll() ([]Job, error)
	Create(name string) error
	Update(name string) error
}

// Реализация логики
type jobService struct {
	jobRepository jobrepository.JobRepository
	tickerMap     map[string]*time.Ticker
	mutex         sync.Mutex
	done          chan struct{}
}

// NewJobService создаёт новый экземпляр сервиса задач
func NewJobService(repository jobrepository.JobRepository) (JobService, func()) {
	newService := &jobService{
		jobRepository: repository,
		tickerMap:     make(map[string]*time.Ticker),
		mutex:         sync.Mutex{},
		done:          make(chan struct{}),
	}

	return newService, func() {
		close(newService.done)
	}
}

// GetAll возвращает все доступные задачи
func (js *jobService) GetAll() ([]Job, error) {
	tasks, err := js.jobRepository.GetAll()
	if err != nil {
		return nil, err
	}

	var res []Job
	for _, task := range tasks {
		res = append(res, Job{
			Name:         task.Name,
			Status:       task.Status.String(),
			StartedTime:  task.StartedTime,
			FinishedTime: task.FinishedTime,
		})
	}

	return res, nil
}

// Create создаёт новую задачу
func (js *jobService) Create(name string) error {
	if name == "" {
		return errors.New("name is required")
	}

	if _, ok := js.tickerMap[name]; ok {
		return fmt.Errorf("job already exists: %s", name)
	}

	newJob := jobrepository.Job{
		Name:        name,
		Status:      jobrepository.InProgress,
		StartedTime: time.Now(),
	}
	if err := js.jobRepository.Add(newJob); err != nil {
		return err
	}

	createdJob, err := js.jobRepository.Get(newJob.Name)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(updateJobStatusInterval)
	js.mutex.Lock()
	defer js.mutex.Unlock()
	js.tickerMap[newJob.Name] = ticker

	go func() {
		defer func() {
			js.mutex.Lock()
			ticker.Stop()
			delete(js.tickerMap, createdJob.Name)
			js.mutex.Unlock()
		}()

		select {
		case <-ticker.C:
			_ = js.jobRepository.Update(jobrepository.Job{
				Name:         createdJob.Name,
				Status:       jobrepository.Done,
				StartedTime:  createdJob.StartedTime,
				FinishedTime: time.Now(),
			})
		case <-js.done:
			return
		}
	}()

	return nil
}

// Update обновляет задачу, добавляя ей дополнительные 5 минут до перехода в статус Done
func (js *jobService) Update(name string) error {
	if name == "" {
		return errors.New("name is required")
	}

	job, err := js.jobRepository.Get(name)
	if err != nil {
		return err
	}

	if job.Status == jobrepository.Done {
		return fmt.Errorf("cannot update finished job: %s", job.Name)
	}

	js.mutex.Lock()
	defer js.mutex.Unlock()
	ticker, ok := js.tickerMap[job.Name]
	if !ok || ticker == nil {
		return fmt.Errorf("cannot update finish time for job: %s", job.Name)
	}
	ticker.Reset(updateJobStatusInterval)

	return nil
}
