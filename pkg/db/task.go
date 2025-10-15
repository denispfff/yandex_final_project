package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	//Согласно требованиям к структуре и тестам ID - string, согласно требованиям к БД - ID - int
	// 1 путь - привести к соотвествию структуру и БД (к интам) НО тогда валидация на ID будет на стороне БД
	// 2 путь - оставить как есть с костылями strconv в процессе обработки post\get запросов
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

func Tasks(limit int, search string, date string) ([]*Task, error) {
	var tasks []*Task
	var rows *sql.Rows
	var err error
	// Ради производительности разбил на 3 кейса, должно быть быстрее
	// Т.К. производится поиск по столбцам - проиндексировал их
	switch {
	case date != "":
		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :search OR comment LIKE :search OR date = :date ORDER BY date LIMIT :limit`
		rows, err = DB.Query(query,
			sql.Named("search", "%"+search+"%"), // а вдруг ищем дату там в исходном форматировании
			sql.Named("date", date),
			sql.Named("limit", limit))
	case search != "":
		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE LOWER(title) LIKE LOWER(?) OR LOWER(comment) LIKE LOWER(?) ORDER BY date LIMIT ?`
		rows, err = DB.Query(query,
			"%"+search+"%",
			"%"+search+"%",
			limit)
	default:
		query := `SELECT id, date, title, comment, repeat FROM scheduler order by date LIMIT :limit`
		rows, err = DB.Query(query,
			sql.Named("limit", limit))
	}

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

func GetTask(id int) (*Task, error) {
	task := Task{}
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
	res, err := DB.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("nothing changed")
	}

	return nil
}

func DeleteTask(task *Task) error {
	query := `DELETE FROM scheduler WHERE id = :id`
	res, err := DB.Exec(query,
		sql.Named("id", task.ID))
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("nothing deleted")
	}

	return nil
}
