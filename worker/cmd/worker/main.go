package main

import (
	"github.com/crxfoz/seo_metrick_parser/internal/consts"
	"github.com/crxfoz/seo_metrick_parser/internal/rq"
	"github.com/crxfoz/seo_metrick_parser/internal/storage"
	"github.com/crxfoz/seo_metrick_parser/internal/utils"
	"github.com/crxfoz/seo_metrick_parser/parsers"
	"github.com/crxfoz/seo_metrick_parser/worker"
	"github.com/crxfoz/seo_metrick_parser/worker/consumer"
	"github.com/crxfoz/webclient"
	"github.com/sirupsen/logrus"
)

func main() {
	forever := make(chan struct{})

	utils.InitApp()

	// Creates a new WorkerPool and puts parsers into it
	workers := worker.NewWorkerPool()

	for _, p := range parsers.ParsersList {
		httpclient := webclient.Config{Timeout: 10, FollowRedirect: false, UseKeepAlive: true}.New()

		w := worker.NewWorker(p.Name, httpclient, p.ParserFn)
		w.RunBackground(p.Timeout)

		workers.Add(p, w)
	}

	// connect to redis, init new QueuePool, create a new queue and put it into pool
	client := utils.RedisConnect()
	qpool := rq.NewQueuePool()
	newQueue := rq.NewRedisQueue(client, consts.TaskName)
	err := qpool.Add(newQueue)
	if err != nil {
		logrus.WithField("err", err).Fatal("could not initialize a new queque")
		return
	}

	connect, err := utils.PostgresConnect()
	if err != nil {
		logrus.WithField("err", err).Fatal("could not connect to postgres DB")
		return
	}

	db := storage.NewStore(connect)

	logrus.Info("Parser successfully started")
	cons := consumer.NewJobConsumer(db, workers)
	newQueue.AddConsumer(cons)

	<-forever
}
