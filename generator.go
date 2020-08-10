package main

import "sync"

// Generator provides a mechanism for creating a cancellable generator with the equivalent of an error promise.
type Generator struct {
	OutputCh  <-chan interface{} // Consume here.
	cancelled bool               // true once cancelled
	wg        sync.WaitGroup     // Block access to the error.
	err       error              // Any error that the generator encountered before stopping.
}

type GeneratorSink chan<- interface{}

// NewGenerator constructs an object that can be used to consume data from a generator or send it a cancel signal.
func NewGenerator(outputCh chan interface{}, function func(GeneratorSink, *Generator)) *Generator {
	generator := Generator{OutputCh: outputCh}
	generator.wg.Add(1)
	go function(outputCh, &generator)
	return &generator
}

func (e *Generator) setError(err error) {
	if err != nil && e.err == nil {
		e.err = err
	}
}

func (e *Generator) Cancel(err error) {
	e.setError(err)
	e.cancelled = true
}

func (e Generator) Cancelled() bool {
	return e.cancelled
}

func (e *Generator) Close(err error) {
	e.setError(err)
	e.wg.Done()
}

func (e *Generator) Error() error {
	e.wg.Wait()
	return e.err
}
