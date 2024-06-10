package server

import (
	"TestTask/internal/jobservice"
	"encoding/json"
	"log"
	"net/http"
)

// JobServer представляет собой API для работы с сервисом задач
type JobServer interface {
	GetAllTasks(w http.ResponseWriter, r *http.Request)
	AddTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
}

// Реализация API работы с сервисом задач
type jobServer struct {
	jobService jobservice.JobService
}

// NewJobServer создаёт новый экземпляр API сервера
func NewJobServer(service jobservice.JobService) JobServer {
	return &jobServer{service}
}

// GetAllTasks обрабатывает GET запрос на получение списка всех задач
func (s *jobServer) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	acceptHeader := r.Header.Get("Accept")
	if acceptHeader != "application/json" && acceptHeader != "*/*" {
		http.Error(w, "Expect application/json accept header", http.StatusUnsupportedMediaType)
	}

	jobs, err := s.jobService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jobsToMarshall := make([]Job, 0, len(jobs))
	for _, job := range jobs {
		jobsToMarshall = append(jobsToMarshall, Job{job})
	}

	jsonResponseBody, err := json.Marshal(jobsToMarshall)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponseBody)
	if err != nil {
		log.Println(err)
	}
}

// AddTask обрабатывает POST запрос на добавление новой задачи
func (s *jobServer) AddTask(w http.ResponseWriter, r *http.Request) {
	type CreateJobRequestBody struct {
		Name string `json:"name"`
	}

	acceptHeader := r.Header.Get("Accept")
	if acceptHeader != "application/json" && acceptHeader != "*/*" {
		http.Error(w, "Expect application/json accept header", http.StatusUnsupportedMediaType)
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Expect application/json content type", http.StatusUnsupportedMediaType)
	}

	var jsonRequestBody CreateJobRequestBody
	if err := json.NewDecoder(r.Body).Decode(&jsonRequestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.jobService.Create(jsonRequestBody.Name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateTask обрабатывает PUT запрос на обновление задачи (добавление 5 минут работы задачи)
func (s *jobServer) UpdateTask(w http.ResponseWriter, r *http.Request) {
	type UpdateJobRequestBody struct {
		Name string `json:"name"`
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Expect application/json content type", http.StatusUnsupportedMediaType)
	}

	var jsonRequestBody UpdateJobRequestBody
	if err := json.NewDecoder(r.Body).Decode(&jsonRequestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.jobService.Update(jsonRequestBody.Name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
