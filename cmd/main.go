package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/abhaybhu10/executor/executor"
)

type TestData struct {
	msg int
}

func main() {
	ex := executor.NewSimpleExecutor(20)
	ex.Start()
	i := 0
	futures := make([]executor.Future, 0)

	for i < 100 {
		i++
		future, _ := ex.Submit(
			func(context context.Context) (interface{}, error) {
				fmt.Printf("task %d started\n", i)
				select {
				case <-time.After(5 * time.Second):
					fmt.Printf("Task %d Finished \n", i)
					return TestData{msg: i}, nil
				case <-context.Done():
					fmt.Printf("Job %d cancelled\n", i)
					return nil, errors.New("Job cancelled")
				}
			},
		)
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
