package server

import (
	"TestTask/internal/jobservice"
	"encoding/json"
)

// Job представляет собой описание задачи для API сервера
type Job struct {
	jobservice.Job
}

// MarshalJSON функция переопределяет маршаллинг структуры в json
func (c Job) MarshalJSON() ([]byte, error) {
	if c.FinishedTime.IsZero() {
		return json.Marshal(struct {
			Name        string `json:"name"`
			Status      string `json:"status"`
			StartedTime string `json:"started_time"`
		}{
			Name:        c.Name,
			Status:      c.Status,
			StartedTime: c.StartedTime.Format("2006-01-02 15:04:05 -0700 MST"),
		})
	}

	return json.Marshal(struct {
		Name         string `json:"name"`
		Status       string `json:"status"`
		StartedTime  string `json:"started_time"`
		FinishedTime string `json:"finished_time"`
	}{
		Name:         c.Name,
		Status:       c.Status,
		StartedTime:  c.StartedTime.Format("2006-01-02 15:04:05 -0700 MST"),
		FinishedTime: c.FinishedTime.Format("2006-01-02 15:04:05 -0700 MST"),
	})
}
