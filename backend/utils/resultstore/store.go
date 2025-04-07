package resultstore

import (
	"database/sql"
	"fmt"

	"JobQueueGo/utils/types"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB(path string) error {
	var err error
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("open db error: %w", err)
	}
	createTable := `
	CREATE TABLE IF NOT EXISTS results (
		job_id TEXT PRIMARY KEY,
		job_type TEXT,
		input TEXT,
    attempts, INTEGER,
		duration INTEGER,
		success BOOLEAN,
		result TEXT
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		return fmt.Errorf("create table error: %w", err)
	}
	return nil
}

func SaveResult(res types.Result) {
	if db == nil {
		fmt.Println("Database not initialized")
		return
	}

	stmt, err := db.Prepare(`
		INSERT OR REPLACE INTO results 
		(job_id, job_type, input, attempts, duration, success, result) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		fmt.Println("Prepare error:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		res.Job.Id,
		res.Job.Type,
		res.Job.Input,
		res.Job.Attempts,
		res.Duration,
		res.Success,
		res.Result,
	)
	if err != nil {
		fmt.Println("Exec error:", err)
	}
}

func GetResultId(id string) (types.Result, bool) {
	if db == nil {
		fmt.Println("Database not initialized")
		return types.Result{}, false
	}

	row := db.QueryRow(`
		SELECT job_id, job_type, input, attempts, duration, success, result
		FROM results
		WHERE job_id = ?
	`, id)

	var res types.Result

	err := row.Scan(
		&res.Job.Id,
		&res.Job.Type,
		&res.Job.Input,
		&res.Job.Attempts,
		&res.Duration,
		&res.Success,
		&res.Result,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Result{}, false // Inget resultat hittat
		}
		fmt.Println("Scan error:", err)
		return types.Result{}, false
	}

	return res, true
}

func GetAll() ([]types.Result, bool) {
	if db == nil {
		fmt.Println("Database not initialized")
		return nil, false
	}

	rows, err := db.Query(`
		SELECT job_id, job_type, input, attempts, duration, success, result 
		FROM results
	`)
	if err != nil {
		fmt.Println("Query error:", err)
		return nil, false
	}
	defer rows.Close()

	results := []types.Result{}

	for rows.Next() {
		var r types.Result
		err := rows.Scan(
			&r.Job.Id,
			&r.Job.Type,
			&r.Job.Input,
			&r.Job.Attempts,
			&r.Duration,
			&r.Success,
			&r.Result,
		)
		if err != nil {
			fmt.Println("Scan error:", err)
			continue // hoppa över trasiga rader
		}
		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Row iteration error:", err)
		return nil, false
	}

	if len(results) == 0 {
		return nil, false
	}
	return results, true
}

func GetMaxJobIDNumber() int64 {
	if db == nil {
		return 0
	}

	var jobID string
	err := db.QueryRow("SELECT job_id FROM results ORDER BY job_id DESC LIMIT 1").Scan(&jobID)
	if err != nil {
		return 0 // tom databas
	}

	var num int64
	_, err = fmt.Sscanf(jobID, "job-%d", &num)
	if err != nil {
		return 0
	}
	return num
}
