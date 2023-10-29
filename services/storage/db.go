package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const file = "prompt.db"
const table = "prompts"

type Prompt struct {
	ID        string
	Prompt    string
	Response  string
	CreatedAt time.Time
}

func openDB() (db *sql.DB) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func execute_sql(sql string, args ...interface{}) (rows sql.Result, err error) {
	db := openDB()
	defer db.Close()

	res, err := db.Exec(sql, args...)
	if err != nil {
		err = fmt.Errorf("execute_sql: %s .. %s", err, sql)
		log.Println("Warning: ", err)
	}
	return res, err
}

func query_db(sql string) (rows *sql.Rows) {
	db := openDB()
	defer db.Close()

	rows, err := db.Query(sql)
	if err != nil {
		err = fmt.Errorf("query_db: %s", err)
		log.Println("Warning: ", err)
	}
	return rows
}

func CreateTable() {
	createSql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY AUTOINCREMENT, prompt TEXT, response TEXT, CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP)", table)
	execute_sql(createSql)
}

func InsertPrompt(prompt string, response string) int64 {
	insertSql := fmt.Sprintf("INSERT INTO %s (prompt, response) VALUES (?, ?)", table)
	res, err := execute_sql(insertSql, prompt, response)
	if err != nil {
		log.Fatal("failed to insert prompt")
	}

	idx, err := res.LastInsertId()
	if err != nil {
		log.Println("Warnig: failed to insert prompt")
	}
	return idx
}

func ReadPrompts() (prompts []Prompt, err error) {
	db := openDB()
	defer db.Close()

	rows := query_db(fmt.Sprintf("SELECT * FROM %s ORDER BY CreatedAt DESC", table))
	defer rows.Close()

	prompts = make([]Prompt, 0)
	for rows.Next() {
		var p Prompt
		err = rows.Scan(&p.ID, &p.Prompt, &p.Response, &p.CreatedAt)
		if err != nil {
			log.Println("Warning: ", err)
			return prompts, err
		}
		prompts = append(prompts, p)
	}
	return prompts, err
}
