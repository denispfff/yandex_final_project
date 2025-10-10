package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"yandex_final_project/pkg/db"
	"yandex_final_project/pkg/task"
)

func addTaskHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	body := req.Body
	defer body.Close()

	var newTask db.Task
	err := json.NewDecoder(body).Decode(&newTask)
	if err != nil {
		errText := "ошибка десериализации JSON"
		logger.Printf("%s: %v", errText, err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if newTask.Title == "" {
		errText := "не указан заголовок задачи"
		logger.Printf("%s", errText)
		http.Error(res, errText, http.StatusBadRequest)
	}

	if newTask.Date == "" {
		newTask.Date = time.Now().Format(task.DateFormat)
	}

	if newTask.Repeat != "" {
		newTask.Date, err = task.NextDate(time.Now(), newTask.Date, newTask.Repeat)
		if err != nil {
			errText := "invalid format "
			logger.Printf("%s: %v", errText, err)
			http.Error(res, errText, http.StatusBadRequest)
		}
	} else {
		_, err := time.Parse(task.DateFormat, newTask.Date)
		if err != nil {
			errText := "invalid date format "
			logger.Printf("%s, %v", errText, err)
			http.Error(res, errText, http.StatusBadRequest)
		}
	}

	id, err := db.AddTask(&newTask)

	if err != nil {
		errText := "db add task error"
		logger.Printf("%s: %v", errText, err)
		http.Error(res, errText, http.StatusBadRequest)
	}

	response := map[string]int64{
		"id": id,
	}

	err = json.NewEncoder(res).Encode(response)
	if err != nil {
		errText := "Respone write error"
		logger.Printf("%s: %v", errText, err)
		http.Error(res, errText, http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
}
