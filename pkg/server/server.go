package server

import (
	"log"
	"net/http"
	"time"

	"yandex_final_project/pkg/api"
)

const webDir = "./web"

type Server struct {
	logger     *log.Logger
	HttpServer *http.Server
}

func New(logger *log.Logger, port string) *Server {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.HandleFunc("/api/nextdate", func(w http.ResponseWriter, r *http.Request) { api.NextDateHandler(w, r, logger) })
	mux.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) { api.TaskHandler(w, r, logger) })
	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) { api.TasksHandler(w, r, logger) })
	mux.HandleFunc("/api/task/done", func(w http.ResponseWriter, r *http.Request) { api.DoneTasksHandler(w, r, logger) })

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{logger: logger, HttpServer: server}
}
