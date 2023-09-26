package util

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
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

func id() (id string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(b)
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
		log.Fatal(err)
	}
	return res, err
}

func query_db(sql string) (rows *sql.Rows) {
	db := openDB()
	defer db.Close()

	rows, err := db.Query(sql)
	if err != nil {
		err = fmt.Errorf("query_db: %s", err)
		log.Fatal(err)
	}
	return rows
}

func CreateTable() {
	tableSql := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s'", table)
	rows := query_db(tableSql)
	defer rows.Close()

	exists := rows.Next()
	if exists {
		return
	}

	createSql := fmt.Sprintf("CREATE TABLE %s (id TEXT PRIMARY KEY, prompt TEXT, response TEXT, CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP)", table)
	execute_sql(createSql)
}

func InsertPrompt(prompt string, response string) (int64, error) {
	insertSql := fmt.Sprintf("INSERT INTO %s (id, prompt, response) VALUES (?, ?, ?)", table)
	res, err := execute_sql(insertSql, id(), prompt, response)
	if err != nil {
		log.Fatal("failed to insert prompt")
	}

	idx, err := res.LastInsertId()
	if err != nil {
		log.Fatal("failed to insert prompt")
	}
	return idx, err
}

func ReadeadPrompts() (prompts []Prompt, err error) {
	db := openDB()
	defer db.Close()

	rows := query_db(fmt.Sprintf("SELECT * FROM %s ORDER BY CreatedAt DESC", table))
	defer rows.Close()

	prompts = make([]Prompt, 0)
	for rows.Next() {
		var p Prompt
		err = rows.Scan(&p.ID, &p.Prompt, &p.Response, &p.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		prompts = append(prompts, p)
	}
	return prompts, err
}
