package api

import (
	"log"
	"net/http"
)

func NextDateHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	res.Header().Set("Content-Type", "text/plain; charset=utf8")

	switch req.Method {
	case http.MethodGet:
		getNextDateHandler(res, req, logger)

	default:
		errText := "method not allowed"
		http.Error(res, errText, http.StatusMethodNotAllowed)
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
		err := jsonError(res, errText, http.StatusMethodNotAllowed)
		if err != nil {
			logger.Println(err)
		}
	}
}

func TasksHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	res.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case http.MethodGet:
		tasksHandler(res, req, logger)
	default:
		errText := "Method Not Allowed"
		err := jsonError(res, errText, http.StatusMethodNotAllowed)
		if err != nil {
			logger.Println(err)
		}
	}
}

func AddTasksHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	res.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case http.MethodGet:
		tasksHandler(res, req, logger)
	default:
		errText := "Method Not Allowed"
		err := jsonError(res, errText, http.StatusMethodNotAllowed)
		if err != nil {
			logger.Println(err)
		}
	}
}

func DoneTasksHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	res.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case http.MethodPost:
		doneTaskHandler(res, req, logger)
	default:
		errText := "Method Not Allowed"
		err := jsonError(res, errText, http.StatusMethodNotAllowed)
		if err != nil {
			logger.Println(err)
		}
	}
}
