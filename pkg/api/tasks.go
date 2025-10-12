package api

import (
	"log"
	"net/http"
	"yandex_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		// здесь вызываете функцию, которая возвращает ошибку в JSON
		// её желательно было реализовать на предыдущем шаге
		errText := "ошибка при получении записей"
		logger.Printf("%s: %v", errText, err)
		jsonError(res, errText, logger)
		return
	}

	if tasks == nil {
		tasks = []*db.Task{}
	}
	writeJson(res, TasksResp{Tasks: tasks}, logger)
}
