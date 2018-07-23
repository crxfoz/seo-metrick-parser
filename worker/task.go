package worker

type Task struct {
	data []UrlConfig
}

// NewTask creates a new Task
func NewTask(data []UrlConfig) *Task {
	return &Task{data}
}

// Run runs a sheduler.
// Iterates over all urls in Task and puts it in worker
func (t *Task) Run(workers *WorkerPool) Reporter {
	result := make(Reporter)

	for _, u := range t.data {
		uurl := result.AddUrl(u.url)

		// Adds task to queue
		workers.Map(func(parserName string, info *ParserService) {
			if u.states[parserName] {
				info.Worker.AddWorkAsync(u.url)
			}
		})

		// Waiting for parser
		workers.Map(func(parserName string, info *ParserService) {
			if u.states[parserName] {
				r := info.Worker.GetResult()
				result.AddForUrl(uurl, parserName, r)
			}
		})
	}

	return result
}
