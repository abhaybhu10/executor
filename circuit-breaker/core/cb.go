package core

import (
	"time"
)

/*
Circuit -->
Execute method-->
failure -->
success -->

open -->
closed -->
halfOpen-->


x failure in last y seconds --> circuit is open

every z second circuit is half opened to check the respose
number of percentage of required allowed during half closed state.

after s consecutive success circuit is openend


*/

type CBSetting struct {
	Name                     string
	CloseToOpenThreshold     int
	HalfCloseToOpenThreshold int
	HalfOpenCloseThreshold   int
	OpenDuration             time.Duration
	Timeout                  time.Duration
	IsSuccessful             func(err error) bool
}

type CircuitBreaker interface {
	Execute(req HttpReq) (interface{}, error)
}

type CircuitBreakerImpl struct {
	metadata *Metadata
}

func (cb *CircuitBreakerImpl) Execute(req HttpReq) (interface{}, error) {
	return cb.metadata.State.Apply(req, cb.metadata)
}

func NewCB(setting CBSetting) CircuitBreaker {
	return &CircuitBreakerImpl{
		metadata: &Metadata{
			State:               &Closed{},
			LastStateChangeTime: time.Now(),
			CountInfo: CountInfo{
				ConsecutiveSuccessCount: 0,
				ConsecutiveFailureCount: 0,
				RequestCount:            0,
			},
			CBSetting: setting,
		},
	}
}
