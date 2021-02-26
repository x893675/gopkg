package mr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errDummy = errors.New("dummy")

func TestMapReduce(t *testing.T) {
	tests := []struct {
		name        string
		mapper      MapperFunc
		reducer     ReducerFunc
		expectErr   error
		expectValue interface{}
	}{
		{
			name:        "test1",
			expectErr:   nil,
			expectValue: 30,
		},
		{
			name: "test dummy error",
			mapper: func(item interface{}, writer Writer, cancel func(error)) {
				v := item.(int)
				if v%3 == 0 {
					cancel(errDummy)
				}
				writer.Write(v * v)
			},
			expectErr: errDummy,
		},
		{
			name: "test cancel with nil",
			mapper: func(item interface{}, writer Writer, cancel func(error)) {
				v := item.(int)
				if v%3 == 0 {
					cancel(nil)
				}
				writer.Write(v * v)
			},
			expectErr:   ErrCancelWithNil,
			expectValue: nil,
		},
		{
			name: "test2",
			reducer: func(pipe <-chan interface{}, writer Writer, cancel func(error)) {
				var result int
				for item := range pipe {
					result += item.(int)
					if result > 10 {
						cancel(errDummy)
					}
				}
				writer.Write(result)
			},
			expectErr: errDummy,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.mapper == nil {
				test.mapper = func(item interface{}, writer Writer, cancel func(error)) {
					v := item.(int)
					writer.Write(v * v)
				}
			}
			if test.reducer == nil {
				test.reducer = func(pipe <-chan interface{}, writer Writer, cancel func(error)) {
					var result int
					for item := range pipe {
						result += item.(int)
					}
					writer.Write(result)
				}
			}
			value, err := MapReduce(func(source chan<- interface{}) {
				for i := 1; i < 5; i++ {
					source <- i
				}
			}, test.mapper, test.reducer, WithWorkers(1))

			assert.Equal(t, test.expectErr, err)
			assert.Equal(t, test.expectValue, value)
		})
	}
}
