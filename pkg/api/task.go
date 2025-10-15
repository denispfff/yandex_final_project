package api

import (
	"encoding/json"
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

func addTaskHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
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

	err = nextdate.ValidateTask(&newTask)
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
		errText := "не указан идентификатор"
		logger.Println(errText)
		res.WriteHeader(http.StatusBadRequest)
		jsonError(res, errText, logger)
		return
	}

	intID, err := strconv.Atoi(taskID)
	if err != nil {
		errText := "некорректный id задачи"
		logger.Printf("%s: %v", errText, err)
		res.WriteHeader(http.StatusBadRequest)
		jsonError(res, errText, logger)
		return
	}

	task, err := db.GetTask(intID)
	if err != nil {
		errText := "Задача не найдена"
		logger.Printf("%s: %v", errText, err)
		res.WriteHeader(http.StatusNotFound)
		jsonError(res, "Задача не найдена", logger)
		return
	}

	res.WriteHeader(http.StatusOK)
	writeJson(res, task, logger)
}

func putTaskHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
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

	err = nextdate.ValidateTask(&task)

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

	res.WriteHeader(http.StatusOK)
	writeJson(res, nil, logger)

}

func deleteTaskHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	task := db.Task{}
	var err error
	task.ID = req.URL.Query().Get("id")
	if task.ID == "" {
		errText := "не указан идентификатор"
		logger.Println(errText)
		res.WriteHeader(http.StatusBadRequest)
		jsonError(res, errText, logger)
		return
	}

	err = db.DeleteTask(&task)
	if err != nil {
		errText := "db delete task error"
		logger.Printf("%s: %v", errText, err)
		res.WriteHeader(http.StatusInternalServerError)
		jsonError(res, err.Error(), logger)
		return
	}

	res.WriteHeader(http.StatusOK)
	writeJson(res, nil, logger)
}

func doneTaskHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	taskID := req.URL.Query().Get("id")
	if taskID == "" {
		errText := "не указан идентификатор"
		logger.Printf("%s", errText)
		res.WriteHeader(http.StatusBadRequest)
		jsonError(res, errText, logger)
		return
	}

	intID, err := strconv.Atoi(taskID)
	if err != nil {
		errText := "некорректный id задачи"
		logger.Printf("%s: %v", errText, err)
		res.WriteHeader(http.StatusBadRequest)
		jsonError(res, "errText", logger)
		return
	}

	task, err := db.GetTask(intID)
	if err != nil {
		errText := "задача не найдена"
		logger.Printf("%s: %v", errText, err)
		res.WriteHeader(http.StatusNotFound)
		jsonError(res, errText, logger)
		return
	}

	switch task.Repeat {
	case "":
		err = db.DeleteTask(task)
		if err != nil {
			errText := "db delete task error"
			logger.Printf("%s: %v", errText, err)
			res.WriteHeader(http.StatusInternalServerError)
			jsonError(res, errText, logger)
			return
		}

	default:
		task.Date, err = nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			logger.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			jsonError(res, err.Error(), logger)
			return
		}

		err = db.UpdateTask(task)
		if err != nil {
			errText := "ошибка обновления задачи в БД"
			logger.Printf("%s: %v", errText, err)
			res.WriteHeader(http.StatusNotModified)
			jsonError(res, errText, logger)
			return
		}
	}
	res.WriteHeader(http.StatusOK)
	writeJson(res, nil, logger)
}
