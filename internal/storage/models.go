package storage

import (
	"github.com/jmoiron/sqlx/types"
	"time"
)

// Data representation of postgres.Data
type Data struct {
	ID     int            `db:"id"`
	TaskID int            `db:"taskid"`
	Data   types.JSONText `db:"data"`
}

// Task representation of postgres.Task
type Task struct {
	ID         int        `db:"id"`
	Status     string     `db:"status"`
	CreatedAt  time.Time  `db:"created_at"`
	StartedAt  *time.Time `db:"started_at"`
	FinishedAt *time.Time `db:"finished_at"`
}
