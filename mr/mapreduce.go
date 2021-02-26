package mr

import (
	"errors"
	"fmt"
	"sync"

	"github.com/x893675/gopkg/runtime"
)

const (
	defaultWorkers = 16
	minWorkers     = 1
)

var (
	// ErrCancelWithNil is an error that mapreduce was cancelled with nil.
	ErrCancelWithNil = errors.New("mapreduce cancelled with nil")
	// ErrReduceNoOutput is an error that reduce did not output a value.
	ErrReduceNoOutput = errors.New("reduce not writing value")
)

type (
	// GenerateFunc is used to let callers send elements into source.
	GenerateFunc func(source chan<- interface{})

	// MapFunc is used to do element processing and write the output to writer.
	MapFunc func(item interface{}, writer Writer)

	// MapperFunc is used to do element processing and write the output to writer,
	// use cancel func to cancel the processing.
	MapperFunc func(item interface{}, writer Writer, cancel func(error))

	// ReducerFunc is used to reduce all the mapping output and write to writer,
	// use cancel func to cancel the processing.
	ReducerFunc func(pipe <-chan interface{}, writer Writer, cancel func(error))

	// Option defines the method to customize the mapreduce.
	Option func(opts *mapReduceOptions)

	mapReduceOptions struct {
		workers int
	}

	// Writer interface wraps Write method.
	Writer interface {
		Write(v interface{})
	}
)

// MapReduce maps all elements generated from given generate func,
// and reduces the output elemenets with given reducer.
func MapReduce(generate GenerateFunc, mapper MapperFunc, reducer ReducerFunc, opts ...Option) (interface{}, error) {
	source := generator(generate)
	return MapReduceWithSource(source, mapper, reducer, opts...)
}

func generator(generate GenerateFunc) chan interface{} {
	source := make(chan interface{})
	go func() {
		defer runtime.HandleCrash()
		defer close(source)
		generate(source)
	}()
	return source
}

// MapReduceWithSource maps all elements from source, and reduce the output elements with given reducer.
func MapReduceWithSource(source <-chan interface{}, mapper MapperFunc, reducer ReducerFunc, opts ...Option) (interface{}, error) {
	options := buildOptions(opts...)
	output := make(chan interface{})
	collector := make(chan interface{}, options.workers)
	errChan := make(chan error, 1)
	defer close(errChan)
	done := make(chan struct{})
	writer := newGuardedWriter(output, done)
	var closeOnce sync.Once

	finish := func() {
		closeOnce.Do(func() {
			close(done)
			close(output)
		})
	}

	cancel := once(func(err error) {
		if err == nil {
			errChan <- ErrCancelWithNil
		} else {
			errChan <- err
		}
		drain(source)
		finish()

	})

	go func() {
		defer func() {
			if r := recover(); r != nil {
				cancel(fmt.Errorf("%v", r))
			} else {
				finish()
			}
		}()
		reducer(collector, writer, cancel)
		drain(collector)
	}()

	go executeMappers(func(item interface{}, w Writer) {
		mapper(item, w, cancel)
	}, source, collector, done, options.workers)

	value, ok := <-output
	if len(errChan) > 0 {
		return nil, <-errChan
	}
	if ok {
		return value, nil
	} else {
		return nil, ErrReduceNoOutput
	}
}

func executeMappers(mapper MapFunc, input <-chan interface{}, collector chan<- interface{},
	done <-chan struct{}, workers int) {
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		close(collector)
	}()

	pool := make(chan struct{}, workers)
	writer := newGuardedWriter(collector, done)
	for {
		select {
		case <-done:
			return
		case pool <- struct{}{}:
			item, ok := <-input
			if !ok {
				<-pool
				return
			}

			wg.Add(1)
			// better to safely run caller defined method
			go func() {
				defer runtime.HandleCrash()
				defer func() {
					wg.Done()
					<-pool
				}()
				mapper(item, writer)
			}()
		}
	}
}

// drain drains the channel.
func drain(channel <-chan interface{}) {
	// drain the channel
	for range channel {
	}
}

// WithWorkers customizes a mapreduce processing with given workers.
func WithWorkers(workers int) Option {
	return func(opts *mapReduceOptions) {
		if workers < minWorkers {
			opts.workers = minWorkers
		} else {
			opts.workers = workers
		}
	}
}

func buildOptions(opts ...Option) *mapReduceOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}

func newOptions() *mapReduceOptions {
	return &mapReduceOptions{
		workers: defaultWorkers,
	}
}

func once(fn func(error)) func(error) {
	o := new(sync.Once)
	return func(err error) {
		o.Do(func() {
			fn(err)
		})
	}
}

type guardedWriter struct {
	channel chan<- interface{}
	done    <-chan struct{}
}

func newGuardedWriter(channel chan<- interface{}, done <-chan struct{}) guardedWriter {
	return guardedWriter{
		channel: channel,
		done:    done,
	}
}

func (gw guardedWriter) Write(v interface{}) {
	select {
	case <-gw.done:
		return
	default:
		gw.channel <- v
	}
}
