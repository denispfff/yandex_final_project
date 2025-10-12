package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"yandex_final_project/pkg/db"
	"yandex_final_project/pkg/task"
)

func writeJson(res http.ResponseWriter, data any, logger *log.Logger) {
	err := json.NewEncoder(res).Encode(data)
	if err != nil {
		logger.Printf("ошибка при сериализации ответа: %v", err)

	}
}

func jsonError(res http.ResponseWriter, errText string, logger *log.Logger) {
	errorResponse := map[string]string{
		"error": errText,
	}
	res.WriteHeader(http.StatusBadRequest)
	writeJson(res, errorResponse, logger)
}

func addTaskHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	res.Header().Set("Content-Type", "application/json")

	var newTask db.Task

	body := req.Body
	defer body.Close()
	// Сразу просчитываем текущую дату, отбрасывая время для дальнейшей логики
	todayString := time.Now().Format(task.DateFormat)
	today, err := time.Parse(task.DateFormat, time.Now().Format(task.DateFormat))

	if err != nil {
		errText := "что-то с текущим временем на сервере"
		logger.Printf("%s: %v", errText, err)
		jsonError(res, errText, logger)
		return
	}

	err = json.NewDecoder(body).Decode(&newTask)
	if err != nil {
		errText := "ошибка десериализации JSON"
		logger.Printf("%s: %v", errText, err)
		jsonError(res, errText, logger)
		return
	}

	if newTask.Title == "" {
		errText := "не указан заголовок задачи"
		logger.Printf("%s", errText)
		jsonError(res, errText, logger)
		return
	}

	if newTask.Date == "" {
		newTask.Date = todayString
	}

	// Для таски без указания времени повторения -
	if newTask.Repeat == "" {
		date, err := time.Parse(task.DateFormat, newTask.Date)
		if err != nil {
			errText := "invalid date format "
			logger.Printf("%s, %v", errText, err)
			jsonError(res, errText, logger)
			return
		}

		if date.Before(today) {
			newTask.Date = todayString
		}
	} else {
		newTask.Date, err = task.NextDate(today, newTask.Date, newTask.Repeat)
		if err != nil {
			errText := "invalid format "
			logger.Printf("%s: %v", errText, err)
			jsonError(res, errText, logger)
			return
		}

		if newTask.Repeat != "" {
			newTask.Date, err = task.NextDate(today, newTask.Date, newTask.Repeat)
			if err != nil {
				errText := "invalid format "
				logger.Printf("%s: %v", errText, err)
				jsonError(res, errText, logger)
				return
			}
		} else {
			date, err := time.Parse(task.DateFormat, newTask.Date)
			if err != nil {
				errText := "invalid date format "
				logger.Printf("%s, %v", errText, err)
				jsonError(res, errText, logger)
				return
			}

			if date.Before(today) {
				newTask.Date = todayString
			}
		}

		//
	}

	id, err := db.AddTask(&newTask)

	if err != nil {
		errText := "db add task error"
		logger.Printf("%s: %v", errText, err)
		jsonError(res, errText, logger)
		return
	}

	response := map[string]int64{
		"id": id,
	}

	writeJson(res, response, logger)
	res.WriteHeader(http.StatusCreated)
}
