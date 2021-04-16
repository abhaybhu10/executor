package executor

import (
	"context"
	"errors"
	"time"
)

type Future interface {
	Get() (interface{}, error)
	GetWithTimeout(timeout int) (interface{}, error)
	Done() bool
	Cancel()
	IsCancelled() bool
	GetError() error
}

type MyFuture struct {
	completed bool
	cancelled bool
	result    chan interface{}
	err       error
	cancel    context.CancelFunc
}

func (f *MyFuture) Get() (interface{}, error) {
	if f.IsCancelled() {
		return nil, errors.New("Job is cancelled")
	}
	return <-f.result, nil
}

func (f *MyFuture) GetWithTimeout(timeout int) (interface{}, error) {
	deadline := time.After(time.Second)
	select {
	case data := <-f.result:
		return data, nil
	case <-deadline:
		return nil, errors.New("Timed out while waiting for result")
	}
}

func (f *MyFuture) Done() bool {
	return f.completed
}

func (f *MyFuture) Cancel() {
	f.cancelled = true
	f.cancel()
}

func (f *MyFuture) IsCancelled() bool {
	return f.cancelled
}

func (f *MyFuture) GetError() error {
	return f.err
}
