package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"yandex_final_project/pkg/db"
	"yandex_final_project/pkg/nextdate"
)

func writeJson(res http.ResponseWriter, data any, logger *log.Logger) {
	err := json.NewEncoder(res).Encode(data)
	if err != nil {
		logger.Printf("ошибка при сериализации ответа: %v", err)
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
	}
}

func jsonError(res http.ResponseWriter, errText string, logger *log.Logger) {
	errorResponse := map[string]string{
		"error": errText,
	}
	writeJson(res, errorResponse, logger)
}

func processTask(newTask *db.Task) error {
	todayString := time.Now().Format(nextdate.DateFormat)
	today, err := time.Parse(nextdate.DateFormat, todayString)
	if err != nil {
		errText := "что-то с текущим временем на сервере"
		return fmt.Errorf("%s: %w", errText, err)
	}

	if newTask.Date == "" {
		newTask.Date = todayString
	}

	dateString, err := time.Parse(nextdate.DateFormat, newTask.Date)
	if err != nil {
		errText := "invalid date format "
		return fmt.Errorf("%s: %w", errText, err)
	}

	if newTask.Repeat == "" {
		if dateString.Before(today) {
			newTask.Date = todayString
		}
	} else {
		next, err := nextdate.NextDate(today, newTask.Date, newTask.Repeat)
		if err != nil {
			errText := "invalid format "
			return fmt.Errorf("%s: %w", errText, err)
		}

		if nextdate.AfterNow(today, dateString) {
			if len(newTask.Repeat) == 0 {
				newTask.Date = todayString
			} else {
				newTask.Date = next
			}
		}
	}
	return nil
}

func addTaskHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	res.Header().Set("Content-Type", "application/json")

	var newTask db.Task

	body := req.Body
	defer body.Close()

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

	err = processTask(&newTask)
	if err != nil {
		logger.Println(err)
		res.WriteHeader(http.StatusBadRequest)
		jsonError(res, err.Error(), logger)
		return
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
	res.WriteHeader(http.StatusCreated)
	writeJson(res, response, logger)
}

func getTaskHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	taskID := req.URL.Query().Get("id")
	if taskID == "" {
		res.WriteHeader(http.StatusBadRequest)
		jsonError(res, "Не указан идентификатор", logger)
		return
	}

	intID, err := strconv.Atoi(taskID)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		jsonError(res, "Invalid id", logger)
		return
	}

	task, err := db.GetTask(intID)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		jsonError(res, "Задача не найдена", logger)
		return
	}

	res.WriteHeader(http.StatusOK)
	writeJson(res, task, logger)
}

func putTaskHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	res.Header().Set("Content-Type", "application/json")

	var task db.Task

	body := req.Body
	defer body.Close()

	err := json.NewDecoder(body).Decode(&task)
	if err != nil {
		errText := "ошибка десериализации JSON"
		logger.Printf("%s: %v", errText, err)
		res.WriteHeader(http.StatusBadRequest)
		jsonError(res, errText, logger)
		return
	}

	if task.Title == "" {
		errText := "не указан заголовок задачи"
		logger.Printf("%s", errText)
		res.WriteHeader(http.StatusBadRequest)
		jsonError(res, errText, logger)
		return
	}

	err = processTask(&task)
	if err != nil {
		logger.Println(err)
		res.WriteHeader(http.StatusBadRequest)
		jsonError(res, err.Error(), logger)
		return
	}

	err = db.UpdateTask(&task)

	if err != nil {
		errText := "db update task error"
		logger.Printf("%s: %v", errText, err)
		res.WriteHeader(http.StatusNotFound)
		jsonError(res, errText, logger)
		return
	}

	response := "{}"
	res.WriteHeader(http.StatusOK)
	_, err = res.Write([]byte(response))

	if err != nil {
		errText := "Respone write error"
		logger.Printf("%s: %v", errText, err)
		res.WriteHeader(http.StatusInternalServerError)
		jsonError(res, err.Error(), logger)
		return
	}
}

func DoneTaskHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	res.Header().Set("Content-Type", "application/json")

	taskID := req.URL.Query().Get("id")
	if taskID == "" {
		res.WriteHeader(http.StatusBadRequest)
		jsonError(res, "Не указан идентификатор", logger)
		return
	}

	intID, err := strconv.Atoi(taskID)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		jsonError(res, "Invalid id", logger)
		return
	}

	task, err := db.GetTask(intID)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		jsonError(res, "Задача не найдена", logger)
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(task)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			jsonError(res, err.Error(), logger)
			return
		}
	} else {
		task.Date, err = nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			jsonError(res, err.Error(), logger)
			return
		}

		err = db.UpdateTask(task)
		if err != nil {
			res.WriteHeader(http.StatusNotModified)
			jsonError(res, err.Error(), logger)
			return
		}

		response := "{}"
		res.WriteHeader(http.StatusOK)
		_, err = res.Write([]byte(response))
		if err != nil {
			jsonError(res, err.Error(), logger)
		}
		return
	}
}
