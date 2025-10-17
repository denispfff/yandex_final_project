package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"yandex_final_project/pkg/api"
	"yandex_final_project/pkg/db"
)

const webDir = "./web"

type Server struct {
	logger     *log.Logger
	HttpServer *http.Server
}

func New(logger *log.Logger) *Server {
	port, ok := os.LookupEnv("TODO_PORT")
	if !ok || len(port) == 0 {
		port = "7540"
	}

	dbFile, ok := os.LookupEnv("TODO_DBFILE")
	if !ok || len(dbFile) == 0 {
		dbFile = "scheduler.db"
	}

	db.Init(dbFile)

	mux := api.Init(webDir, logger)

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
