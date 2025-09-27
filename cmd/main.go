package main

import (
	"log"
	"os"

	"yandex_final_project/server"
)

func main() {
	mainLogger := log.New(
		os.Stdout,
		"server: ",
		log.Lshortfile,
	)

	srv := server.New(mainLogger)

	if err := srv.HttpServer.ListenAndServe(); err != nil {
		mainLogger.Fatal(err)
	}
}
