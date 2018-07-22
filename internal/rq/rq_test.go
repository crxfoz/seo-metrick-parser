package rq

import (
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var queueTest QueueService
var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6379",
		Password: "",
		DB:       0,
	})
	queueTest = NewRedisQueue(client, queueName)
}

func TestQueue_Push(t *testing.T) {
	queueTest2 := NewRedisQueue(client, "test2")
	err := queueTest2.Push("dsadadasda")
	assert.Nil(t, err)
}

func TestQueue_Pop(t *testing.T) {
	queueTest2 := NewRedisQueue(client, "test3")
	err := queueTest2.Push("somedata")
	assert.Nil(t, err)

	ret, err := queueTest2.Pop()
	if assert.Nil(t, err) {
		assert.Equal(t, "somedata", ret)
	}
}

func TestQueue_BPop(t *testing.T) {
	queueTest2 := NewRedisQueue(client, "test4")
	data := "somedata1337"

	go func() {
		ret, err := queueTest2.BPop(5 * time.Second)
		if assert.Nil(t, err) {
			assert.Equal(t, data, ret)
		}
	}()

	queueTest2.Push(data)
}

func TestQueue_Subscribe(t *testing.T) {
	data := "somedata1337"
	queueTest2 := NewRedisQueue(client, "test5")
	done := make(chan struct{})

	retCh := queueTest2.Subscribe(3 * time.Second)

	counter := 0

	go func() {
		for {
			select {
			case <-time.After(5 * time.Second):
				t.Errorf("should get data before timeout")
				return
			case v, ok := <-retCh:
				if ok {
					counter++
				} else {
					assert.Equal(t, 2, counter)
					done <- struct{}{}
					return
				}
				assert.Equal(t, v, data)
			}
		}
	}()

	queueTest2.Push(data)
	queueTest2.Push(data)

	<-done

}

func TestQueue_GetName(t *testing.T) {
	assert.Equal(t, queueName, queueTest.GetName())
}

func TestQueue_Subscribe2(t *testing.T) {
	data := "somedata1337"
	queueTest2 := NewRedisQueue(client, "test5")
	done := make(chan struct{})

	retCh := queueTest2.Subscribe(3 * time.Second)

	counter := 0

	go func() {
		for {
			select {
			case <-time.After(5 * time.Second):
				t.Errorf("should get data before timeout")
				return
			case v, ok := <-retCh:
				if ok {
					counter++
				} else {
					assert.Equal(t, 0, counter, "should not get any data because timeout")
					done <- struct{}{}
					return
				}
				assert.Equal(t, v, data)
			}
		}
	}()

	time.AfterFunc(time.Second*4, func() {
		queueTest2.Push(data)
		queueTest2.Push(data)
	})

	<-done

}

type customConsumer struct {
	t       *testing.T
	invoked bool
}

func (c *customConsumer) Consume(data string) {
	c.invoked = true
	assert.Equal(c.t, "somedatafoobar", data)
}

func TestQueue_AddConsumer(t *testing.T) {
	qq := NewRedisQueue(client, "foorbarqueue")

	cons := &customConsumer{t, false}

	qq.AddConsumer(cons)
	qq.Push("somedatafoobar")

	time.Sleep(time.Second * 2)

	assert.Equal(t, true, cons.invoked, "consumer wasn't called")

}
