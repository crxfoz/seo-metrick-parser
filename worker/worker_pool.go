package worker

import (
	"github.com/crxfoz/seo_metrick_parser/parsers"
)

// WorkerPool stores the map of *ParserService
type WorkerPool struct {
	data map[string]*ParserService
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{make(map[string]*ParserService)}
}

// Map iterates over map and calls `fn` for every key
func (l *WorkerPool) Map(fn func(index string, info *ParserService)) {
	for item := range l.data {
		fn(item, l.data[item])
	}
}

// Add adds a Parser into Reposity.
func (l *WorkerPool) Add(svc parsers.Parser, worker *Worker) {
	p := &ParserService{Parser: svc}
	l.data[svc.Name] = p
	l.data[svc.Name].Worker = worker
}
