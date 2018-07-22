package main

import (
	"github.com/crxfoz/seo_metrick_parser/internal/consts"
	"github.com/crxfoz/seo_metrick_parser/internal/rq"
	"github.com/crxfoz/seo_metrick_parser/internal/storage"
	"github.com/crxfoz/seo_metrick_parser/internal/utils"
	"github.com/crxfoz/seo_metrick_parser/server"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	utils.InitApp()

	client := utils.RedisConnect()
	qpool := rq.NewQueuePool()
	err := qpool.Add(rq.NewRedisQueue(client, consts.TaskName))
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

	server.Start(os.Getenv("API_SERVER_HOST"), os.Getenv("API_SERVER_PORT"), db, qpool)
}
