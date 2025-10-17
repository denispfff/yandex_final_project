package api

import (
	"log"
	"net/http"
)

func AuthHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	res.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case http.MethodPost:
		authHandler(res, req, logger)

	default:
		errText := "method not allowed"
		http.Error(res, errText, http.StatusMethodNotAllowed)
		return
	}
}

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

func Init(webDir string, logger *log.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.HandleFunc("/api/signin", func(w http.ResponseWriter, r *http.Request) { AuthHandler(w, r, logger) })
	mux.HandleFunc("/api/nextdate", func(w http.ResponseWriter, r *http.Request) { NextDateHandler(w, r, logger) })
	//роуты требующие авторизации
	mux.HandleFunc("/api/task", Wrap(TaskHandler, AuthMiddleware, logger))
	mux.HandleFunc("/api/tasks", Wrap(TasksHandler, AuthMiddleware, logger))
	mux.HandleFunc("/api/task/done", Wrap(DoneTasksHandler, AuthMiddleware, logger))

	return mux
}
