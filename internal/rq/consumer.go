package rq

// Consumer interface that user must implement.
// Consume it's a function that will be called every time it receives data
type Consumer interface {
	Consume(string)
}
