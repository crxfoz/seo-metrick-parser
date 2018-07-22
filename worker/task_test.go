package worker

import (
	"fmt"
	"github.com/crxfoz/seo_metrick_parser/parsers"
	"github.com/crxfoz/webclient"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type wrapParser struct {
	t          *testing.T
	invokedFn1 int
	invokedFn2 int
}

func (w *wrapParser) Fn1(c *webclient.Webclient, site string) (interface{}, error) {
	w.invokedFn1++
	return fmt.Sprintf("%s-%s", site, "parser1"), nil
}

func (w *wrapParser) Fn2(c *webclient.Webclient, site string) (interface{}, error) {
	w.invokedFn2++
	return fmt.Sprintf("%s-%s", site, "parser2"), nil
}

func TestTask_Run(t *testing.T) {
	NewWorkerPool()
	cfg := []UrlConfig{
		{
			url:    "http://google.com",
			states: map[string]bool{"parserFn1": true, "parserFn2": true},
		},
		{
			url:    "http://yandex.ru",
			states: map[string]bool{"parserFn1": true, "parserFn2": false},
		},
	}

	pp := &wrapParser{
		t:          t,
		invokedFn1: 0,
		invokedFn2: 0,
	}

	list := []parsers.Parser{
		{
			Name:        "parserFn1",
			Description: "parserFn1 descr",
			Timeout:     1 * time.Second,
			Status:      true,
			ParserFn:    pp.Fn1,
		},
		{
			Name:        "parserFn2",
			Description: "parserFn2 descr",
			Timeout:     1 * time.Second,
			Status:      true,
			ParserFn:    pp.Fn2,
		},
	}

	workers := NewWorkerPool()

	for _, p := range list {
		w := NewWorker(p.Name, nil, p.ParserFn)
		w.RunBackground(p.Timeout)

		workers.Add(p, w)
	}

	newTask := NewTask(cfg)
	result := newTask.Run(workers)

	assert.Equal(t, 2, pp.invokedFn1)
	assert.Equal(t, 1, pp.invokedFn2)

	assert.Equal(t, "http://google.com-parser2", result["http://google.com"]["parserFn2"])
	assert.Equal(t, "http://yandex.ru-parser1", result["http://yandex.ru"]["parserFn1"])
}
