package circuit_breaker

import (
	"time"

	"github.com/thisismz/data-processor/pkg/databases/cache"
)

type circuitBreaker struct {
	Name        string
	MaxRequests uint32
	Interval    time.Duration
	Timeout     time.Duration
	IsLeader    bool
}

func New(name string, maxRequest uint32, interval time.Duration, isLeader bool) circuitBreaker {
	return circuitBreaker{
		Name:        name,
		MaxRequests: maxRequest,
		Interval:    interval,
		Timeout:     0,
		IsLeader:    isLeader,
	}
}
func (c *circuitBreaker) Run(err error) {
	cache.RedisSyncStatus = false
	timeTicker(c.Interval)
}
func timeTicker(minute time.Duration) {
	ticker := time.NewTicker(minute * time.Minute)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				_ = t
				sync()
			}
		}
	}()
}
func GetCircuitStatus() bool {
	return cache.RedisSyncStatus
}
