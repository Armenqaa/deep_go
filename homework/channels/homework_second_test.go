package main_2

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var ErrPoolIsFull = errors.New("pool is full")

type task struct{}

type WorkerPool struct {
	pool chan task
	wg   sync.WaitGroup
}

func NewWorkerPool(workersNumber int) *WorkerPool {
	pool := make(chan task, workersNumber)
	return &WorkerPool{pool: pool, wg: sync.WaitGroup{}}
}

// Return an error if the pool is full
func (wp *WorkerPool) AddTask(t func()) error {
	select {
	case wp.pool <- task{}:
		wp.wg.Add(1)
		go func() {
			t()
			<-wp.pool
			wp.wg.Done()
		}()
		return nil
	default:
		return ErrPoolIsFull
	}
}

func (wp *WorkerPool) Shutdown() {
	close(wp.pool)
	wp.wg.Wait()
}

func TestWorkerPool(t *testing.T) {
	var counter atomic.Int32
	task := func() {
		time.Sleep(time.Millisecond * 500)
		counter.Add(1)
	}

	pool := NewWorkerPool(2)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)

	time.Sleep(time.Millisecond * 600)
	assert.Equal(t, int32(2), counter.Load())

	time.Sleep(time.Millisecond * 600)
	assert.Equal(t, int32(3), counter.Load())

	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	pool.Shutdown() // wait tasks

	assert.Equal(t, int32(6), counter.Load())
}
