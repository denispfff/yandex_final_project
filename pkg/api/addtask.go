package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"yandex_final_project/pkg/db"
	"yandex_final_project/pkg/task"
)

func jsonError(res http.ResponseWriter, errText string, logger *log.Logger) {
	errorResponse := map[string]string{
		"error": errText,
	}

	resErr := json.NewEncoder(res).Encode(errorResponse)
	if resErr != nil {
		logger.Printf("ошибка при сериализации ошибки")
	}
}

func addTaskHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	body := req.Body
	defer body.Close()

	var newTask db.Task
	err := json.NewDecoder(body).Decode(&newTask)
	if err != nil {
		errText := "ошибка десериализации JSON"
		logger.Printf("%s: %v", errText, err)
		res.WriteHeader(http.StatusBadRequest)
		jsonError(res, errText, logger)
		return
	}

	if newTask.Title == "" {
		errText := "не указан заголовок задачи"
		logger.Printf("%s", errText)
		res.WriteHeader(http.StatusBadRequest)
		jsonError(res, errText, logger)
		return
	}

	if newTask.Date == "" {
		newTask.Date = time.Now().Format(task.DateFormat)
	}

	if newTask.Repeat != "" {
		newTask.Date, err = task.NextDate(time.Now(), newTask.Date, newTask.Repeat)
		if err != nil {
			errText := "invalid format "
			logger.Printf("%s: %v", errText, err)
			res.WriteHeader(http.StatusBadRequest)
			jsonError(res, errText, logger)
			return
		}
	} else {
		date, err := time.Parse(task.DateFormat, newTask.Date)
		if err != nil {
			errText := "invalid date format "
			logger.Printf("%s, %v", errText, err)
			res.WriteHeader(http.StatusBadRequest)
			jsonError(res, errText, logger)
			return
		}

		if date.Before(time.Now()) {
			newTask.Date = time.Now().Format(task.DateFormat)
		}
	}

	id, err := db.AddTask(&newTask)

	if err != nil {
		errText := "db add task error"
		logger.Printf("%s: %v", errText, err)
		res.WriteHeader(http.StatusBadRequest)
		jsonError(res, errText, logger)
		return
	}

	response := map[string]int64{
		"id": id,
	}

	err = json.NewEncoder(res).Encode(response)
	if err != nil {
		errText := "Respone write error"
		logger.Printf("%s: %v", errText, err)
		res.WriteHeader(http.StatusBadRequest)
		jsonError(res, errText, logger)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
}
