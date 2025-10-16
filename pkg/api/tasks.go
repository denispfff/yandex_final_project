package api

import (
	"log"
	"net/http"
	"time"
	"yandex_final_project/pkg/db"
	"yandex_final_project/pkg/nextdate"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	var limit int = 50
	var searchDate string

	search := req.URL.Query().Get("search")
	if search != "" {
		date, err := time.Parse("02.01.2006", search)
		if err == nil {
			searchDate = date.Format(nextdate.DateFormat)
		}
	}
	tasks, err := db.Tasks(limit, search, searchDate)
	if err != nil {
		errText := "ошибка при получении записей"
		logger.Printf("%s: %v", errText, err)
		err := jsonError(res, errText, http.StatusInternalServerError)
		if err != nil {
			logger.Println(err)
		}
		return
	}

	if tasks == nil {
		tasks = []*db.Task{}
	}

	err = writeJson(res, TasksResp{Tasks: tasks}, http.StatusOK)
	if err != nil {
		logger.Println(err)
	}
}
