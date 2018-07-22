package worker

import (
	"fmt"
	"github.com/crxfoz/webclient"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type wrapTest struct {
	fn func(*webclient.Webclient, string) (interface{}, error)
	t  *testing.T
}

var counter = 0

func TestWorker(t *testing.T) {
	wrap := &wrapTest{fn: func(client *webclient.Webclient, s string) (interface{}, error) {
		defer func() { counter++ }()
		return fmt.Sprintf("%s-%d", s, counter), nil
	}}

	wrk := NewWorker("somename", webclient.Config{}.New(), wrap.fn)
	wrk.RunBackground(1 * time.Second)
	wrap.t = t
	done := make(chan struct{}, 1)

	list := []string{
		"http://google.com",
		"http://yandex.ru",
		"http://fb.com",
	}

	go func() {
		for _, l := range list {
			wrk.AddWorkAsync(l)
		}
	}()

	go func() {
		for i := range list {
			res := wrk.GetResult()
			assert.Equal(t, fmt.Sprintf("%s-%d", list[i], i), res)
		}
		done <- struct{}{}
	}()

	for {
		select {
		case <-time.After(time.Second * 30):
			t.Errorf("should be done before")
			return
		case <-done:
			return
		}
	}
}
