package customwg

import (
	"sync"
	"sync/atomic"
)

type CustomWaitGroup struct {
	mu      *sync.Mutex
	cond    *sync.Cond
	counter atomic.Int32
}

func NewWaitGroup() *CustomWaitGroup {
	mu := sync.Mutex{}
	return &CustomWaitGroup{
		mu:   &mu,
		cond: sync.NewCond(&mu),
	}
}

func (wg *CustomWaitGroup) Add(delta int) {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	wg.counter.Add(int32(delta))
	if wg.counter.Load() <= 0 {
		wg.cond.Broadcast()
	}
}

func (wg *CustomWaitGroup) Done() {
	wg.Add(-1)
}

func (wg *CustomWaitGroup) Wait() {
	wg.mu.Lock()
	defer wg.mu.Unlock()
	for wg.counter.Load() > 0 {
		wg.cond.Wait()
	}
}
