package db

import (
	"database/sql"
	"fmt"
	"time"
	"yandex_final_project/pkg/task"
)

type Task struct {
	ID      string `json:"id" db:"id"`
	Date    string `json:"date" db:"date"`
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`

	res, err := DB.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	var tasks []*Task

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date >= :date order by date LIMIT :limit`
	rows, err := DB.Query(query,
		sql.Named("date", time.Now().Format(task.DateFormat)),
		sql.Named("limit", limit))
	if err != nil {
		return nil, fmt.Errorf("query execute error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task Task

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}
