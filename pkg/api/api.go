package api

import (
	"log"
	"net/http"
	"time"
	"yandex_final_project/pkg/nextdate"
)

func NextDateHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	res.Header().Set("Content-Type", "text/plain; charset=utf8")

	if req.Method != http.MethodGet {
		errText := "Method not allowed"
		http.Error(res, errText, http.StatusMethodNotAllowed)
		return
	}

	now := req.FormValue("now")
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	nowDate, err := time.Parse(nextdate.DateFormat, now)
	if err != nil {
		errText := "invalid date format"
		logger.Printf("%s: %v", errText, err)
		http.Error(res, errText, http.StatusBadRequest)
		return
	}
	nextDate, err := nextdate.NextDate(nowDate, date, repeat)

	if err != nil {
		errText := "invalid repeat rule"
		logger.Printf("%s: %v", errText, err)
		http.Error(res, errText, http.StatusBadRequest)
		return
	}

	_, err = res.Write([]byte(nextDate))

	if err != nil {
		errText := "Respone write error"
		logger.Printf("%s: %v", errText, err)
		http.Error(res, errText, http.StatusInternalServerError)
		return
	}
}

func TaskHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	res.Header().Set("Content-Type", "application/json")

	switch req.Method {
	//other methods
	case http.MethodGet:
		getTaskHandler(res, req, logger)
	case http.MethodPost:
		addTaskHandler(res, req, logger)
	case http.MethodPut:
		putTaskHandler(res, req, logger)
	case http.MethodDelete:
		deleteTaskHandler(res, req, logger)
	default:
		errText := "Method Not Allowed"
		res.WriteHeader(http.StatusMethodNotAllowed)
		jsonError(res, errText, logger)
	}
}

func TasksHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	res.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case http.MethodGet:
		tasksHandler(res, req, logger)
	default:
		errText := "method not allowed"
		res.WriteHeader(http.StatusMethodNotAllowed)
		jsonError(res, errText, logger)
	}
}

func AddTasksHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	res.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case http.MethodGet:
		tasksHandler(res, req, logger)
	default:
		errText := "method not allowed"
		res.WriteHeader(http.StatusMethodNotAllowed)
		jsonError(res, errText, logger)
	}
}

func DoneTasksHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	res.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case http.MethodPost:
		doneTaskHandler(res, req, logger)
	default:
		errText := "method not allowed"
		res.WriteHeader(http.StatusMethodNotAllowed)
		jsonError(res, errText, logger)
	}
}
