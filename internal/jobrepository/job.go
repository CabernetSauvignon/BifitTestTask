package jobrepository

import "time"

// Status представляет собой тип для перечисление возможных состояний задачи
type Status int

// Перечисление возможных состояний задачи
const (
	Unknown Status = iota
	InProgress
	Done
)

// String возвращает статус задачи в строковом формате
func (s Status) String() string {
	switch s {
	case Unknown:
		return "unknown"
	case InProgress:
		return "in progress"
	case Done:
		return "done"
	default:
		return "unknown"
	}
}

// Job представляет собой описание задачи
type Job struct {
	Name         string
	Status       Status
	StartedTime  time.Time
	FinishedTime time.Time
}
