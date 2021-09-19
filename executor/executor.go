package executor

import (
	"context"
	"errors"
	"math/rand"
	"sync"
)

const (
	capacity = 2000
)

type Callable interface {
	Call(context.Context) (interface{}, error)
}
type Executor interface {
	Submit(task Callable) (Future, error)
	Shutdown()
}

type SimpleExecutor struct {
	capacity    int
	concurrency int
	status      Status
	jobs        chan job
	running     chan job
	mu          *sync.Mutex
	wg          *sync.WaitGroup
	ctx         context.Context
}

func NewSimpleExecutor(concurrency int) Executor {
	ex := &SimpleExecutor{
		concurrency: concurrency,
		jobs:        make(chan job, capacity),
		running:     make(chan job, concurrency),
		capacity:    capacity,
		mu:          &sync.Mutex{},
		wg:          &sync.WaitGroup{},
		ctx:         context.Background(),
	}
	ex.start()
	return ex
}

func (SE *SimpleExecutor) Submit(task Callable) (Future, error) {

	if SE.status != Status("RUNNING") {
		return nil, errors.New("Executor is not running")
	}
	ctx, cancel := context.WithCancel(context.Background())
	job := job{
		id:     rand.Int(),
		task:   task,
		status: Status("Submitted"),
		ctx:    ctx,
		future: &MyFuture{
			cancelled: false,
			completed: false,
			result:    make(chan interface{}),
			err:       nil,
			cancel:    cancel,
		},
	}
	SE.jobs <- job
	return job.future, nil
}

func (SE *SimpleExecutor) start() {
	SE.status = Status("RUNNING")
	go func() {
		for {
			select {
			case job := <-SE.jobs:
				SE.running <- job
				SE.wg.Add(1)
				go func() {
					defer SE.wg.Done()
					SE.completeTask(job)
					<-SE.running
				}()
			case <-SE.ctx.Done():
				SE.mu.Lock()
				SE.status = Status("SHUTTING DOWN")
				SE.mu.Unlock()

				SE.wg.Wait()

				SE.mu.Lock()
				SE.status = Status("SHUTDOWN")
				SE.mu.Unlock()
			}
		}
	}()
}

func (SE *SimpleExecutor) Shutdown() {
	SE.ctx.Done()
}

func (SE *SimpleExecutor) completeTask(job job) {

	result, err := job.task.Call(job.ctx)
	if err != nil {
		job.err = err
		job.status = Status("FAILED")
	} else {
		job.status = Status("DONE")
		future := job.future.(*MyFuture)
		future.completed = true
		future.result <- result
	}
}

type task func(context.Context) (interface{}, error)

type job struct {
	id     int
	status Status
	task   Callable
	err    error
	future Future
	ctx    context.Context
}

type Status string
