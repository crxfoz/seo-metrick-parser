package utils

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"os"
)

func InitApp() {
	if os.Getenv("PRODUCTION") == "1" {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func RedisConnect() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("RD_HOST"), os.Getenv("RD_PORT")),
		Password: "",
		DB:       0,
	})

	return client
}

func PostgresConnect() (*sqlx.DB, error) {
	connStr := fmt.Sprintf("port=%s host=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("PG_PORT"),
		os.Getenv("PG_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"))

	return sqlx.Connect("postgres", connStr)
}
