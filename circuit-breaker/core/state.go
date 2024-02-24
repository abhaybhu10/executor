package core

import (
	"errors"
	"fmt"
	"time"
)

type HttpReq func() (interface{}, error)

type CountInfo struct {
	ConsecutiveFailureCount int
	ConsecutiveSuccessCount int
	RequestCount            int
}

type Metadata struct {
	State               State
	LastStateChangeTime time.Time
	CBSetting
	CountInfo
}

type Open struct{}
type Closed struct{}
type HalfOpen struct{}

type State interface {
	Apply(HttpReq, *Metadata) (interface{}, error)
}

func (s Open) Apply(req HttpReq, meta *Metadata) (interface{}, error) {
	if meta.LastStateChangeTime.Add(meta.OpenDuration).Before(time.Now()) {
		changeState(&HalfOpen{}, meta)
		return meta.State.Apply(req, meta)
	}
	return nil, errors.New("circuit is open")
}

func (s HalfOpen) Apply(req HttpReq, meta *Metadata) (interface{}, error) {
	data, err := req()
	if err == nil || meta.IsSuccessful(err) {
		if meta.ConsecutiveSuccessCount > meta.HalfOpenCloseThreshold {
			changeState(&Closed{}, meta)
		}
		meta.ConsecutiveSuccessCount++
	} else {
		if meta.ConsecutiveFailureCount > meta.HalfCloseToOpenThreshold {
			changeState(&Open{}, meta)
		}
		meta.ConsecutiveFailureCount++
	}
	return data, err
}

func (s Closed) Apply(req HttpReq, meta *Metadata) (interface{}, error) {
	data, err := req()
	if err != nil && !meta.IsSuccessful(err) {
		if meta.ConsecutiveFailureCount > meta.CloseToOpenThreshold {
			changeState(&Open{}, meta)
		}
		meta.ConsecutiveFailureCount++
	}
	return data, err
}

func changeState(to State, meta *Metadata) {
	fmt.Printf("changing state to %T\n", to)
	meta.State = to
	meta.ConsecutiveFailureCount = 0
	meta.ConsecutiveFailureCount = 0
	meta.LastStateChangeTime = time.Now()
}
