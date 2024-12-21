package main

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var ErrPoolIsFull = errors.New("pool is full")

type token struct{}

type WorkerPool struct {
	tokens chan token
}

func NewWorkerPool(workersNumber int) *WorkerPool {
	tokens := make(chan token, workersNumber)
	for range workersNumber {
		tokens <- token{}
	}
	return &WorkerPool{tokens: tokens}
}

// Return an error if the pool is full
func (wp *WorkerPool) AddTask(task func()) error {
	select {
	case <-wp.tokens:
		go func() {
			task()
			wp.tokens <- token{}
		}()
		return nil
	default:
		return ErrPoolIsFull
	}
}

func (wp *WorkerPool) Shutdown() {
	for range cap(wp.tokens) {
		<-wp.tokens
	}
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
