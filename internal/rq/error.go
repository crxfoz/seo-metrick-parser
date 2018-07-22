package rq

import (
	"fmt"
)

type ErrQueueNotExists struct {
	name string
}

func (e *ErrQueueNotExists) Error() string {
	return fmt.Sprintf("RedisQueue not exists: %s", e.name)
}

type ErrReservedName struct {
	name string
}

func (e *ErrReservedName) Error() string {
	return fmt.Sprintf("name is reserved: %s", e.name)
}
