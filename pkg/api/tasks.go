package api

import (
	"log"
	"net/http"
	"time"
	"yandex_final_project/pkg/db"
	"yandex_final_project/pkg/task"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	search := req.URL.Query().Get("search")
	searchDate, err := time.Parse("02.01.2006", search)
	if err == nil {
		search = searchDate.Format(task.DateFormat)
	}
	tasks, err := db.Tasks(50, search) // в параметре максимальное количество записей
	if err != nil {
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
