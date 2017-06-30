package workergroup

import (
	"runtime"
	"sync/atomic"
	"time"
)

const waitTime = time.Nanosecond * 10

// WorkerGroup acts as a sync.WaitGroup but with a max worker count.
type WorkerGroup struct {
	maxWorkers uint32
	workers    uint32
}

// MaxWorkers sets max workers to `workers` or returns current setting if < 0
func (wg *WorkerGroup) MaxWorkers(workers uint32) uint32 {
	if workers == 0 {
		// update max workers to GOMAXPROCS if set to 0
		atomic.CompareAndSwapUint32(&wg.maxWorkers, 0, uint32(runtime.GOMAXPROCS(0)))

		return atomic.LoadUint32(&wg.maxWorkers)
	}

	atomic.StoreUint32(&wg.maxWorkers, workers)

	return workers
}

// Add delta to the worker count
func (wg *WorkerGroup) Add(delta uint32) {
	// update max workers to GOMAXPROCES if set to 0
	atomic.CompareAndSwapUint32(&wg.maxWorkers, 0, uint32(runtime.GOMAXPROCS(0)))

	var oldCount uint32
	for oldCount = atomic.LoadUint32(&wg.workers); oldCount+delta > atomic.LoadUint32(&wg.maxWorkers); oldCount = atomic.LoadUint32(&wg.workers) {
		time.Sleep(waitTime)
	}

	if atomic.CompareAndSwapUint32(&wg.workers, oldCount, oldCount+delta) {
		return
	}

	wg.Add(delta)
}

// Done decrement the worker counter
func (wg *WorkerGroup) Done() {
	atomic.AddUint32(&wg.workers, ^uint32(0))
}

// Wait until worker count is zero
func (wg *WorkerGroup) Wait() {
	for atomic.LoadUint32(&wg.workers) != 0 {
		time.Sleep(waitTime)
	}
}
