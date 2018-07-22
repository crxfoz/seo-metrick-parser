package consumer

import (
	"encoding/json"
	"github.com/crxfoz/seo_metrick_parser/internal/consts"
	"github.com/crxfoz/seo_metrick_parser/internal/storage"
	"github.com/crxfoz/seo_metrick_parser/worker"
	"github.com/sirupsen/logrus"
)

// JobConsumer emptry struct that implements rq.Consumer interface
type JobConsumer struct {
	db      storage.StoreService
	workers worker.WorkerPoolService
}

func NewJobConsumer(db storage.StoreService, w worker.WorkerPoolService) *JobConsumer {
	return &JobConsumer{db: db, workers: w}
}

// Consume implements rq.Consumer interface
// gets raw string json from redis, creates a new parser.Task and waits until it runs out
// and then inserts results to postgres.Data and updates postgres.Task status
func (s *JobConsumer) Consume(data string) {
	logrus.Debug("consumer got a some data")

	type urlInfo struct {
		ID   int64 `json:"id"`
		Task []struct {
			Url     string   `json:"url"`
			Parsers []string `json:"parsers"`
		} `json:"task"`
	}

	var newData urlInfo

	err := json.Unmarshal([]byte(data), &newData)
	if err != nil {
		logrus.WithField("err", err).Error("could not unmarshal data from redis queue")
		return
	}

	var urlConfig []worker.UrlConfig

	for _, d := range newData.Task {
		states := make(map[string]bool)
		for _, parserName := range d.Parsers {
			states[parserName] = true
		}

		urlConfig = append(urlConfig, worker.NewUrlConfig(d.Url, states))
	}

	newTask := worker.NewTask(urlConfig)
	logrus.WithField("id", newData.ID).WithField("status", consts.TaskStatusRunned).Debug("updating task status")

	err = s.db.UpdateTask(newData.ID, consts.TaskStatusRunned, "started_at")
	if err != nil {
		logrus.WithField("err", err).WithField("status", consts.TaskStatusRunned).Error("could not update task")
		return
	}

	result := newTask.Run(s.workers)

	logrus.WithField("task_id", newData.ID).Debug("got results")

	j, err := json.Marshal(result)
	if err != nil {
		logrus.WithField("err", err).Error("could not marshal result of work")
		return
	}

	err = s.db.InsertData(newData.ID, j)
	if err != nil {
		logrus.WithField("err", err).Error("could not insert to Data")
		return
	}

	err = s.db.UpdateTask(newData.ID, consts.TaskStatusDone, "finished_at")
	if err != nil {
		return
		logrus.WithField("err", err).WithField("status", consts.TaskStatusDone).Error("could not update task")
	}
}
