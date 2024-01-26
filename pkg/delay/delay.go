package delay

import (
	"sync"
	"time"
)

// Wraps a function so it can be scheduled to run after a delay. If multiple
// schedules are made before the function is run, only the last one will be
// executed.
type Function struct {
	fun            func()
	waitFor        time.Duration
	afterFuncTimer *time.Timer
	tick           chan struct{}
	mu             sync.Mutex
}

func NewFunction(fun func(), waitFor time.Duration) *Function {
	return &Function{fun: fun, waitFor: waitFor}
}

// Tick returns a channel that will be closed when the function is run.
func (f *Function) Tick() <-chan struct{} {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.tick = make(chan struct{})
	return f.tick
}

// Schedule schedules the function to run after the configured delay. If the
// function is already scheduled, the previous schedule is cancelled.
func (f *Function) Schedule() {
	if f.afterFuncTimer != nil {
		f.afterFuncTimer.Stop()
	}
	f.afterFuncTimer = time.AfterFunc(f.waitFor, func() {
		f.mu.Lock()
		defer f.mu.Unlock()
		f.fun()
		if f.tick != nil {
			close(f.tick)
		}
	})
}
