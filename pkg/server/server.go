package server

import (
	"log"
	"net/http"
	"time"
)

const webDir = "./web"

type Server struct {
	logger     *log.Logger
	HttpServer *http.Server
}

func New(logger *log.Logger, port string) *Server {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(webDir)))

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
