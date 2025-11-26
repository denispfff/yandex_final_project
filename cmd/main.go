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

	srv := server.New(mainLogger)
	defer db.DB.Close()

	if err := srv.HttpServer.ListenAndServe(); err != nil {
		mainLogger.Fatal(err)
	}
}
