package api

import (
	"log"
	"net/http"
	"time"
	"yandex_final_project/pkg/nextdate"
)

func getNextDateHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
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
