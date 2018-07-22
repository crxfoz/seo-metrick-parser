package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var db StoreService

func init() {
	s := fmt.Sprintf("port=%s host=%s user=%s dbname=%s sslmode=disable",
		"5432",
		"0.0.0.0",
		"postgres",
		"seo_parser_test")

	if c, err := sqlx.Connect("postgres", s); err != nil {
		logrus.WithField("err", err).Panic("could not connect to DB")
		return
	} else {
		db = NewStore(c)
	}
}

func TestPgRepository_CreateTask(t *testing.T) {
	fn, id, err := db.CreateTask()
	assert.Nil(t, err)
	assert.NotZero(t, id)

	_, err = db.GetTaskByID(id)
	assert.NotNil(t, err)

	fn(false)
	_, err = db.GetTaskByID(id)
	assert.NotNil(t, err)

	fn(true)
	_, err = db.GetTaskByID(id)
	assert.NotNil(t, err)
}

func TestPgRepository_GetTasks(t *testing.T) {
	fn, _, err := db.CreateTask()
	if assert.Nil(t, err) {
		fn(true)
	}

	tasks, err := db.GetTasks()
	assert.Nil(t, err)
	assert.NotZero(t, len(tasks))
}

func TestPgRepository_GetTaskByID(t *testing.T) {
	fn, id, err := db.CreateTask()
	if assert.Nil(t, err) {
		fn(true)
	}

	task, err := db.GetTaskByID(id)
	assert.Nil(t, err)
	assert.Equal(t, "waiting", task.Status)
	assert.NotEqual(t, task.CreatedAt.String(), time.Time{}.String())
}

func TestPgRepository_InsertData(t *testing.T) {
	fn, id, err := db.CreateTask()
	if assert.Nil(t, err) {
		fn(true)
	}

	err = db.InsertData(id, []byte(`{"test": [1,2]}`))
	assert.Nil(t, err)

	data, err := db.GetDataByTaskID(id)
	assert.Nil(t, err)
	assert.Equal(t, `{"test": [1,2]}`, data.Data.String())
	assert.NotZero(t, data.ID)
	assert.Equal(t, int(id), data.TaskID)
}

func TestPgRepository_UpdateTask(t *testing.T) {
	fn, id, err := db.CreateTask()
	if assert.Nil(t, err) {
		fn(true)
	}

	err = db.UpdateTask(id, "failed", "")
	assert.Nil(t, err)

	task, err := db.GetTaskByID(id)
	assert.Nil(t, err)
	assert.Equal(t, "failed", task.Status)
	assert.NotEqual(t, time.Time{}.String(), task.CreatedAt)
	assert.Nil(t, task.StartedAt)
	assert.Nil(t, task.FinishedAt)

	err = db.UpdateTask(id, "runned", "started_at")
	assert.Nil(t, err)

	task, err = db.GetTaskByID(id)
	assert.Nil(t, err)
	assert.Equal(t, "runned", task.Status)
	assert.NotNil(t, task.StartedAt)
	assert.Nil(t, task.FinishedAt)

	err = db.UpdateTask(id, "done", "finished_at")
	assert.Nil(t, err)

	task, err = db.GetTaskByID(id)
	assert.Nil(t, err)
	assert.Equal(t, "done", task.Status)
	assert.NotNil(t, task.FinishedAt)
}
