package jobservice

import "time"

// Job представляет собой описание задачи для сервиса задач
type Job struct {
	Name         string
	Status       string
	StartedTime  time.Time
	FinishedTime time.Time
}
