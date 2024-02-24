package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/abhaybhu10/executor/circuit-breaker/core"
)

func main() {
	settings := core.CBSetting{
		Name:                     "HTTP Circuit breaker",
		CloseToOpenThreshold:     5,
		HalfCloseToOpenThreshold: 2,
		HalfOpenCloseThreshold:   2,
		OpenDuration:             5 * time.Second,
		Timeout:                  5 * time.Second,
		IsSuccessful: func(err error) bool {
			return err == nil
		},
	}

	client := CBMiddlerWare{
		cb: core.NewCB(settings),
	}
	for i := 0; i < 100; i++ {
		url := "httpsss//google.com"
		if i > 10 && i < 20 {
			url = "https://google.com"
		}
		if _, err := client.Get(url); err != nil {
			fmt.Printf("err %v\n", err)
		}
		time.Sleep(time.Second)
	}
}

type CBMiddlerWare struct {
	cb core.CircuitBreaker
}

func (m *CBMiddlerWare) Get(url string) ([]byte, error) {
	getFunc := func() (interface{}, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	}

	body, err := m.cb.Execute(getFunc)
	if err != nil {
		return nil, err
	}
	return body.([]byte), err
}
