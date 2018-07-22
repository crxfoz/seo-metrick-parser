package rq

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	errEmptyQueue = "redis: nil"
	errRedisDown  = "EOF"
)

type QueueService interface {
	GetName() string
	Push(interface{}) error
	Pop() (string, error)
	BPop(time.Duration) (string, error)
	AddConsumer(Consumer)
	Subscribe(time.Duration) <-chan string
}

// RedisQueue implements a message RedisQueue using redis
type RedisQueue struct {
	conn *redis.Client
	name string
}

// NewRedisQueue returns a new RedisQueue
func NewRedisQueue(conn *redis.Client, name string) *RedisQueue {
	return &RedisQueue{conn: conn, name: name}
}

func (r *RedisQueue) GetName() string {
	return r.name
}

// Push RPUSH data into redis
func (r *RedisQueue) Push(value interface{}) error {
	err := r.conn.RPush(r.name, value).Err()
	if err != nil {
		return err
	}

	return nil
}

// Pop LPOP data from redis
func (r *RedisQueue) Pop() (string, error) {
	return r.conn.LPop(r.name).Result()
}

// AddConsumer adds a new consumer who subscribes on channel
// Gets data from channel until it's closed and calls Consume with this data
func (r *RedisQueue) AddConsumer(c Consumer) {
	go func() {
		jonChan := r.Subscribe(0)
		for job := range jonChan {
			c.Consume(job)
		}
	}()
}

// BPop BLPOP data from redis.
// Timeout specifies how long to wait for data
func (r *RedisQueue) BPop(timeout time.Duration) (string, error) {
	s, err := r.conn.BLPop(timeout, r.name).Result()

	if err != nil {
		// if  err.Error() != errRedisDown {
		// 	logrus.WithField("err", err).Error("got unexpected error")
		// }
		return "", err
	}

	if len(s) != 2 {
		return "", fmt.Errorf("unexpected len of result. len: %d, arr: %v", len(s), s)
	}

	if s[0] != r.name {
		return "", fmt.Errorf("unexpected list name: %s", s[0])
	}

	return s[1], nil
}

// Subscribe to the RedisQueue
// returns the channel to which it writes the received data
func (r *RedisQueue) Subscribe(timeout time.Duration) <-chan string {
	ch := make(chan string, 1)

	go func(ch chan<- string) {
		for {
			value, err := r.BPop(timeout)
			if err != nil {
				logrus.WithField("err", err).Error("error")
				close(ch)
				return
			}
			ch <- value
		}
	}(ch)

	return ch
}
