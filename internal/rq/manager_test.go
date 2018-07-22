package rq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	queueManager *QueuePool
	queueName    = "queue_test"
)

func init() {
	queueManager = NewQueuePool()
}

func TestQueuePool_Add(t *testing.T) {
	q := &RedisQueue{
		name: queueName,
	}

	err := queueManager.Add(q)
	assert.Nil(t, err)

	q2 := &RedisQueue{
		name: queueName,
	}
	err = queueManager.Add(q2)
	assert.NotNil(t, err)
}

func TestQueuePool_Get(t *testing.T) {
	queueManager.Add(&RedisQueue{name: "bar"})
	queueManager.Add(&RedisQueue{name: queueName})
	queueManager.Add(&RedisQueue{name: "foo"})

	ret, err := queueManager.Get(queueName)
	assert.Nil(t, err)
	assert.Equal(t, ret.GetName(), queueName)

	ret, err = queueManager.Get("foo")
	assert.Nil(t, err)
	assert.Equal(t, ret.GetName(), "foo")

}

func TestQueuePool_UnsafeGet(t *testing.T) {

	queueManager.Add(&RedisQueue{name: "bar"})
	queueManager.Add(&RedisQueue{name: queueName})
	queueManager.Add(&RedisQueue{name: "foo"})

	ret := queueManager.UnsafeGet(queueName)
	assert.Equal(t, ret.GetName(), queueName)
}
