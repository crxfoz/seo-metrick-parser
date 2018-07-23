package worker

import (
	"github.com/crxfoz/seo_metrick_parser/parsers"
	"github.com/crxfoz/webclient"
	"github.com/sirupsen/logrus"
	"time"
)

// Worker for every parser we have
type Worker struct {
	isStopped bool
	name      string
	work      chan string
	out       chan interface{}
	stop      chan struct{}
	client    *webclient.Webclient
	timeout   time.Duration
	parserFn  parsers.ParserFn
}

// messageOutput wraps message if parser returns error
type messageOutput struct {
	Message string `json:"message"`
}

// NewWorker creates a new Worker
func NewWorker(name string, httpclient *webclient.Webclient, fn parsers.ParserFn) *Worker {
	a := new(Worker)
	a.work = make(chan string, 1)
	a.out = make(chan interface{})
	a.stop = make(chan struct{}, 1)
	a.client = httpclient
	a.name = name
	a.parserFn = fn

	return a
}

// AddWorkAsync adds job that have to be done
// puts `site` into Worker.work channel
func (a *Worker) AddWorkAsync(site string) {
	a.work <- site
}

// GetResult gets results of job.
// Gets data from Worker.out channel
func (a *Worker) GetResult() interface{} {
	return <-a.out
}

// RunBackground runs a scheduler in goroutine.
// Worker lisetening the Worker.work channel and as soon as it receives data calls Worker.parserFn
// and sends results to Worker.out channel
func (a *Worker) RunBackground(timeout time.Duration) {
	a.isStopped = false
	a.timeout = timeout
	go a.run()
}

// run runs a worker. See RunBackground for more info.
func (a *Worker) run() {
	for {
		select {
		case site := <-a.work:
			logrus.
				WithField("module", "Worker").
				WithField("url", site).
				WithField("Worker", a.name).
				Debug("Got an a task")
			result, err := a.parserFn(a.client, site)
			if err != nil {
				a.out <- messageOutput{"error"}
			} else {
				a.out <- result
			}
			time.Sleep(a.timeout)
		}
	}
}
