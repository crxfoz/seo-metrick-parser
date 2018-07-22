package storage

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"time"
)

// StoreService The interface describing how to work this DB
type StoreService interface {
	CreateTask() (CommitFn, int64, error)
	GetTasks() ([]Task, error)
	GetTaskByID(int64) (Task, error)
	InsertData(int64, []byte) error
	UpdateTask(int64, string, string) error
	GetDataByTaskID(int64) (Data, error)
}

type pgRepository struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *pgRepository {
	return &pgRepository{db}
}

// CommitFn use commitfn(false) for Rollback or commitfn(true) for Commit
type CommitFn func(bool) error

// CreateTask Creates a new task in DB
func (d *pgRepository) CreateTask() (CommitFn, int64, error) {
	var tx *sql.Tx
	var err error

	var fn CommitFn = func(doCommit bool) error {
		if tx == nil {
			return fmt.Errorf("tx is nil")
		}
		if doCommit {
			return tx.Commit()
		}
		return tx.Rollback()
	}

	tx, err = d.db.Begin()
	if err != nil {
		return fn, 0, err
	}

	l, err := tx.Prepare("INSERT INTO task DEFAULT VALUES RETURNING id")
	if err != nil {
		logrus.Error(err)
		return fn, 0, err
	}

	var id int64
	err = l.QueryRow().Scan(&id)
	if err != nil {
		return fn, 0, err
	}

	return fn, id, nil
}

// GetTasks returns all Task from DB
func (d *pgRepository) GetTasks() ([]Task, error) {
	var tasks []Task
	err := d.db.Select(&tasks, "SELECT * FROM task ORDER BY id")
	if err != nil {
		return tasks, err
	}

	return tasks, nil
}

// GetTaskByID returns one Task
func (d *pgRepository) GetTaskByID(id int64) (Task, error) {
	var task Task
	row := d.db.QueryRowx("SELECT * FROM task WHERE id = $1", id)
	err := row.StructScan(&task)

	if err != nil {
		return task, err
	}

	return task, nil
}

// InsertData to 'Data' table
func (d *pgRepository) InsertData(taskID int64, data []byte) error {
	_, err := d.db.Exec("INSERT INTO data (taskid, data) VALUES ($1, $2)", taskID, data)
	return err
}

// UpdateTask updates status and date params of the task with `id`.
// status - internal/consts
// paramTime - nothing or `started_at` or `finished_at`
func (d *pgRepository) UpdateTask(id int64, status string, paramTime string) error {
	var err error
	now := time.Now()

	switch paramTime {
	case "started_at":
		_, err = d.db.Exec("UPDATE task SET status = $1, started_at = $3 WHERE id = $2", status, id, now)
	case "finished_at":
		_, err = d.db.Exec("UPDATE task SET status = $1, finished_at = $3 WHERE id = $2", status, id, now)
	default:
		_, err = d.db.Exec("UPDATE task SET status = $1 WHERE id = $2", status, id)
	}

	return err
}

// GetDataByTaskID returns the data associated with task
func (d *pgRepository) GetDataByTaskID(taskID int64) (Data, error) {
	var data Data
	row := d.db.QueryRowx("SELECT * FROM data WHERE taskid = $1", taskID)
	err := row.StructScan(&data)
	return data, err
}
