package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(256) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
	);

CREATE INDEX scheduler_date ON scheduler (date);
`

var DB *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	// Насколько понял из условия - соединение оставляем открытым, это странно
	// defer db.Close()

	if install {
		_, err := db.Exec(schema)
		if err != nil {
			log.Fatal(err)
		}

		log.Print("Создана таблица scheduler")
	}

	return nil
}
