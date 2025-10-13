package main

import (
	"log"
	"os"

	"yandex_final_project/pkg/db"
	"yandex_final_project/pkg/server"
)

func main() {
	mainLogger := log.New(
		os.Stdout,
		"server: ",
		log.LstdFlags|log.Lshortfile,
	)

	port, ok := os.LookupEnv("TODO_PORT")
	if !ok || len(port) == 0 {
		port = "7540"
	}

	dbFile, ok := os.LookupEnv("TODO_DBFILE")
	if !ok || len(dbFile) == 0 {
		dbFile = "scheduler.db"
	}

	db.Init(dbFile)
	defer db.DB.Close()

	srv := server.New(mainLogger, port)

	if err := srv.HttpServer.ListenAndServe(); err != nil {
		mainLogger.Fatal(err)
	}
}
