package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/j-dunham/openai-cli/config"
	_ "github.com/mattn/go-sqlite3"
)

const TABLE = "prompts"

type Prompt struct {
	ID        string
	Prompt    string
	Response  string
	Role      string
	CreatedAt time.Time
}

type DB struct {
	config *config.Config
	conn   *sql.DB
}

func NewDB(cfg *config.Config) *DB {
	return &DB{
		config: cfg,
		conn:   nil,
	}

}

func (db *DB) openDB() {
	conn, err := sql.Open("sqlite3", db.config.DBFile)
	if err != nil {
		log.Fatal(err)
	}
	db.conn = conn
	db.createTable()
}

func (db *DB) executeSql(sql string, args ...interface{}) (rows sql.Result, err error) {
	if db.conn == nil {
		db.openDB()
	}
	res, err := db.conn.Exec(sql, args...)
	if err != nil {
		err = fmt.Errorf("execute_sql: %s .. %s", err, sql)
		log.Println("Warning: ", err)
	}
	return res, err
}

func (db *DB) queryDb(sql string) (rows *sql.Rows) {
	if db.conn == nil {
		db.openDB()
	}
	rows, err := db.conn.Query(sql)
	if err != nil {
		err = fmt.Errorf("query_db: %s", err)
		log.Println("Warning: ", err)
	}
	return rows
}

func (db *DB) createTable() {
	createSql := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY AUTOINCREMENT, role TEXT, prompt TEXT, response TEXT, CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP)",
		TABLE,
	)
	db.executeSql(createSql)
}

func (db *DB) InsertPrompt(role, prompt, response string) int64 {
	insertSql := fmt.Sprintf("INSERT INTO %s (role, prompt, response) VALUES (?, ?, ?)", TABLE)
	res, err := db.executeSql(insertSql, role, prompt, response)
	if err != nil {
		log.Fatal("failed to insert prompt")
	}

	idx, err := res.LastInsertId()
	if err != nil {
		log.Println("Warnig: failed to insert prompt")
	}
	return idx
}

func (db *DB) ReadPrompts() (prompts []Prompt, err error) {
	rows := db.queryDb(fmt.Sprintf("SELECT * FROM %s ORDER BY CreatedAt DESC", TABLE))

	prompts = make([]Prompt, 0)
	for rows.Next() {
		var p Prompt
		err = rows.Scan(&p.ID, &p.Role, &p.Prompt, &p.Response, &p.CreatedAt)
		if err != nil {
			log.Println("Warning: ", err)
			return prompts, err
		}
		prompts = append(prompts, p)
	}
	return prompts, err
}

func (db *DB) Close() {
	if db.conn != nil {
		db.conn.Close()
	}
}
