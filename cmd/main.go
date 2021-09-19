package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/abhaybhu10/executor/executor"
)

type TestTask struct {
	msg int
}

func (t *TestTask) Call(context context.Context) (interface{}, error) {
	fmt.Printf("task %d started\n", t.msg)
	select {
	case <-time.After(5 * time.Second):
		fmt.Printf("Task %d Finished \n", t.msg)
		return TestData{msg: t.msg}, nil
	case <-context.Done():
		fmt.Printf("Job %d cancelled\n", t.msg)
		return nil, errors.New("Job cancelled")
	}
}

func main() {
	ex := executor.NewSimpleExecutor(20)
	i := 0
	futures := make ([]executor.Future, 0)

	//submit 100 task
	for i < 100 {
		i++
		td := &TestTask{
			msg: i,
		}
		future, _ := ex.Submit(td)
		futures = append(futures, future)
	}

	for _, future := range futures[0:50] {
		future.Cancel()
	}

	//time.Sleep(10 * time.Second)
	for _, future := range futures {
		msg, err := future.Get()
		if err != nil {
			fmt.Printf("Error occured %v\n", err)
		} else {
			fmt.Printf("Receieved resuld %d\n", msg.(TestData).msg)
		}
	}

	time.Sleep(1000 * time.Second)
	ex.Shutdown()
}

type TestData struct {
	msg int
}
