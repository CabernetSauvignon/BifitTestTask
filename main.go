package main

import (
	"TestTask/internal/jobrepository"
	"TestTask/internal/jobserver"
	"TestTask/internal/jobservice"
	"log"
	"net/http"
)

func main() {
	jobRepository := jobrepository.NewJobRepository()

	jobService, jobServiceCanceller := jobservice.NewJobService(jobRepository)
	defer jobServiceCanceller()

	jobServer := server.NewJobServer(jobService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", jobServer.AddTask)
	mux.HandleFunc("PUT /jobs", jobServer.UpdateTask)
	mux.HandleFunc("GET /jobs", jobServer.GetAllTasks)

	log.Fatalln(http.ListenAndServe(":8080", mux))
}
