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

func writeJson(res http.ResponseWriter, data any, status int) error {
	js, err := json.Marshal(data)
	if err != nil {
		errorText := fmt.Sprintf("ошибка при сериализации ответа: %v", err)
		anotherErr := jsonError(res, errorText, http.StatusInternalServerError)
		if anotherErr != nil {
			return fmt.Errorf("%v, %v", err, anotherErr)
		}
		return err
	}

	res.WriteHeader(status)
	_, err = res.Write(js)
	// если соединение оборвано - вернуть пользователю уже ничего не сможем
	if err != nil {
		return err
	}
	return nil
}

func jsonError(res http.ResponseWriter, errText string, status int) error {
	data := map[string]string{
		"error": errText,
	}

	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	res.WriteHeader(status)
	_, err = res.Write(js)
	if err != nil {
		return err
	}

	return nil
}

func addTaskHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	var newTask db.Task

	body := req.Body
	defer body.Close()

	err := json.NewDecoder(body).Decode(&newTask)
	if err != nil {
		errText := "ошибка десериализации JSON"
		logger.Printf("%s: %v", errText, err)
		err = jsonError(res, errText, http.StatusBadRequest)
		if err != nil {
			logger.Println(err)
		}
		return
	}

	err = nextdate.ValidateTask(&newTask)
	if err != nil {
		logger.Println(err)
		err = jsonError(res, err.Error(), http.StatusBadRequest)
		if err != nil {
			logger.Println(err)
		}
		return
	}

	id, err := db.AddTask(&newTask)
	if err != nil {
		errText := "db add task error"
		logger.Printf("%s: %v", errText, err)
		err = jsonError(res, errText, http.StatusBadRequest)
		if err != nil {
			logger.Println(err)
		}
		return
	}

	response := map[string]int64{
		"id": id,
	}

	err = writeJson(res, response, http.StatusCreated)
	if err != nil {
		logger.Println(err)
	}
}

func getTaskHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	taskID := req.URL.Query().Get("id")
	if taskID == "" {
		errText := "не указан идентификатор"
		logger.Println(errText)
		err := jsonError(res, errText, http.StatusBadRequest)
		if err != nil {
			logger.Println(err)
		}
		return
	}

	intID, err := strconv.Atoi(taskID)
	if err != nil {
		errText := "некорректный id задачи"
		logger.Printf("%s: %v", errText, err)
		err := jsonError(res, errText, http.StatusBadRequest)
		if err != nil {
			logger.Println(err)
		}
		return
	}

	task, err := db.GetTask(intID)
	if err != nil {
		errText := "Задача не найдена"
		logger.Printf("%s: %v", errText, err)
		err := jsonError(res, errText, http.StatusNotFound)
		if err != nil {
			logger.Println(err)
		}
		return
	}

	err = writeJson(res, task, http.StatusOK)
	if err != nil {
		logger.Println(err)
	}
}

func putTaskHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	var task db.Task

	body := req.Body
	defer body.Close()

	err := json.NewDecoder(body).Decode(&task)
	if err != nil {
		errText := "ошибка десериализации JSON"
		logger.Printf("%s: %v", errText, err)
		err := jsonError(res, errText, http.StatusBadRequest)
		if err != nil {
			logger.Println(err)
		}
		return
	}

	err = nextdate.ValidateTask(&task)

	if err != nil {
		logger.Println(err)
		err := jsonError(res, err.Error(), http.StatusBadRequest)
		if err != nil {
			logger.Println(err)
		}
		return
	}

	err = db.UpdateTask(&task)

	if err != nil {
		errText := "db update task error"
		logger.Printf("%s: %v", errText, err)
		err := jsonError(res, errText, http.StatusNotFound)
		if err != nil {
			logger.Println(err)
		}
		return
	}

	err = writeJson(res, struct{}{}, http.StatusOK)
	if err != nil {
		logger.Println(err)
	}

}

func deleteTaskHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	task := db.Task{}
	var err error
	task.ID = req.URL.Query().Get("id")
	if task.ID == "" {
		errText := "не указан идентификатор"
		logger.Println(errText)
		err := jsonError(res, errText, http.StatusBadRequest)
		if err != nil {
			logger.Println(err)
		}
		return
	}

	err = db.DeleteTask(&task)
	if err != nil {
		errText := "db delete task error"
		logger.Printf("%s: %v", errText, err)
		err := jsonError(res, errText, http.StatusInternalServerError)
		if err != nil {
			logger.Println(err)
		}
		return
	}

	err = writeJson(res, struct{}{}, http.StatusOK)
	if err != nil {
		logger.Println(err)
	}

}

func doneTaskHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	taskID := req.URL.Query().Get("id")
	if taskID == "" {
		errText := "не указан идентификатор"
		logger.Printf("%s", errText)
		err := jsonError(res, errText, http.StatusBadRequest)
		if err != nil {
			logger.Println(err)
		}
		return
	}

	intID, err := strconv.Atoi(taskID)
	if err != nil {
		errText := "некорректный id задачи"
		logger.Printf("%s: %v", errText, err)
		err := jsonError(res, errText, http.StatusBadRequest)
		if err != nil {
			logger.Println(err)
		}
		return
	}

	task, err := db.GetTask(intID)
	if err != nil {
		errText := "задача не найдена"
		logger.Printf("%s: %v", errText, err)
		err := jsonError(res, errText, http.StatusNotFound)
		if err != nil {
			logger.Println(err)
		}
		return
	}

	switch task.Repeat {
	case "":
		err = db.DeleteTask(task)
		if err != nil {
			errText := "db delete task error"
			logger.Printf("%s: %v", errText, err)
			err := jsonError(res, errText, http.StatusInternalServerError)
			if err != nil {
				logger.Println(err)
			}
			return
		}

	default:
		task.Date, err = nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			logger.Println(err)
			err := jsonError(res, err.Error(), http.StatusInternalServerError)
			if err != nil {
				logger.Println(err)
			}
			return
		}

		err = db.UpdateTask(task)
		if err != nil {
			errText := "ошибка обновления задачи в БД"
			logger.Printf("%s: %v", errText, err)
			err := jsonError(res, errText, http.StatusNotModified)
			if err != nil {
				logger.Println(err)
			}
			return
		}
	}

	err = writeJson(res, struct{}{}, http.StatusOK)
	if err != nil {
		logger.Println(err)
	}

}
