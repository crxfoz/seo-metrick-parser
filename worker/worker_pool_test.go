package worker

import (
	"github.com/crxfoz/seo_metrick_parser/parsers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWorkerPool_Map(t *testing.T) {
	pool := NewWorkerPool()

	list := []parsers.Parser{
		{
			Name:        "alexa",
			Description: "alexa description",
			Timeout:     4,
			Status:      true,
		},
		{
			Name:        "tyc",
			Description: "tyc descr",
			Timeout:     7,
			Status:      true,
		},
	}

	for _, p := range list {
		w := NewWorker(p.Name, nil, p.ParserFn)
		w.RunBackground(p.Timeout)

		pool.Add(p, w)
	}

	index := 0
	pool.Map(func(s string, service *ParserService) {
		assert.Equal(t, list[index].Name, s)
		index++
	})
}
