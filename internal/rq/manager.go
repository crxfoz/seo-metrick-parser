package rq

// Manager is interface describing how to work with QueuePool
type Manager interface {
	Get(string) (QueueService, error)
	UnsafeGet(string) QueueService
	Add(QueueService) error
}

// var _ Manager = (*QueuePool)(nil)

// QueuePool stores RedisQueue objects.
type QueuePool struct {
	pool map[string]QueueService
}

func NewQueuePool() *QueuePool {
	return &QueuePool{make(map[string]QueueService)}
}

// Get returns a RedisQueue by RedisQueue name
// returns an error if RedisQueue doesnt exists
func (m *QueuePool) Get(name string) (QueueService, error) {
	v, ok := m.pool[name]
	if !ok {
		return nil, &ErrQueueNotExists{name}
	}
	return v, nil
}

// UnsafeGet returns a RedisQueue by name
func (m *QueuePool) UnsafeGet(name string) QueueService {
	return m.pool[name]
}

// Add a RedisQueue into pool
// returns an error if RedisQueue with name alredy exists
func (m *QueuePool) Add(q QueueService) error {
	if _, ok := m.pool[q.GetName()]; ok {
		return &ErrReservedName{q.GetName()}
	}

	m.pool[q.GetName()] = q
	return nil
}
