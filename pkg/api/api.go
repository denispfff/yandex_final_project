package api

import (
	"log"
	"net/http"
	"time"
	"yandex_final_project/pkg/task"
)

func NextDateHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {

	if req.Method != http.MethodGet {
		errText := "Method not allowed"
		http.Error(res, errText, http.StatusMethodNotAllowed)
		return
	}

	res.Header().Set("Content-Type", "text/plain; charset=utf8")
	now := req.FormValue("now")
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	nowDate, err := time.Parse(task.DateFormat, now)
	if err != nil {
		errText := "invalid date format"
		logger.Printf("%s: %v", errText, err)
		http.Error(res, errText, http.StatusBadRequest)
		return
	}
	nextDate, err := task.NextDate(nowDate, date, repeat)

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
	// Оказывается, Write автоматически записывает 200
	// res.WriteHeader(http.StatusOK)
}

func TaskHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	switch req.Method {
	//other methods
	case http.MethodPost:
		addTaskHandler(res, req, logger)
	default:
		errText := "Method Not Allowed"
		jsonError(res, errText, logger)
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func TasksHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	switch req.Method {
	case http.MethodGet:
		tasksHandler(res, req, logger)
	default:
		errText := "method not allowed"
		jsonError(res, errText, logger)
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}
