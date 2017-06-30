package workergroup

import (
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerGroupMaxWorkers(t *testing.T) {
	wg := WorkerGroup{}

	if wg.MaxWorkers(0) != uint32(runtime.GOMAXPROCS(0)) {
		t.Error("Default max workers should be GOMAXPROCS got", wg.MaxWorkers(0))
	}

	if wg.MaxWorkers(3) != 3 {
		t.Error("Failed to set max workers to 3")
	}
}

func TestWorkerGroupAdd(t *testing.T) {
	wg := WorkerGroup{}

	wg.MaxWorkers(3)

	wg.Add(1)
	wg.Add(1)
	wg.Add(1)

	waitChan := make(chan time.Time, 1)

	start := time.Now()
	go func() {
		wg.Add(1)
		waitChan <- time.Now()
	}()

	time.Sleep(time.Millisecond * 50)

	wg.Done()
	stop := <-waitChan

	if time.Duration(stop.Sub(start).Nanoseconds()) < (time.Millisecond * 50) {
		t.Error("wait channel should have waited at least 50ms, waited", stop.Sub(start))
	}
}

func TestWorkerGroupWait(t *testing.T) {
	wg := WorkerGroup{}

	wg.MaxWorkers(3)

	for i := 0; i < 32; i++ {
		go func() {
			wg.Add(1)
			time.Sleep(time.Millisecond)
			wg.Done()
		}()
	}

	// for safety
	time.Sleep(time.Millisecond)

	wg.Wait()

	time.Sleep(time.Millisecond)
	if atomic.LoadUint32(&wg.workers) != 0 {
		t.Error("Failed to wait until workers was actually 0")
	}
}
