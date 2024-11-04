package main

import (
	// "reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type CircularQueueOpt struct {
	values             unsafe.Pointer
	start, filled, cap int
}

func NewCircularQueueOpt(size int) CircularQueueOpt {
	values := make([]int, size)
	return CircularQueue{
		values: unsafe.Pointer(&values[0]),
		cap:    size,
	}
}

func (q *CircularQueueOpt) Push(value int) bool {
	if q.Full() {
		return false
	}
	old := q.valueAt(q.start+q.filled)
	*old = value
	q.filled++
	return true

}

func (q *CircularQueueOpt) Pop() bool {
	if q.Empty() {
		return false
	}
	q.start++
	q.filled--
	return true
}

func (q *CircularQueueOpt) Front() int {
	if q.Empty() {
		return -1
	}
	return *q.valueAt(q.start)
}

func (q *CircularQueueOpt) Back() int {
	if q.Empty() {
		return -1
	}
	return *q.valueAt(q.start+q.filled-1)
}

func (q *CircularQueueOpt) Empty() bool {
	return q.filled == 0
}

func (q *CircularQueueOpt) Full() bool {
	return q.filled == q.Cap()
}

func (q *CircularQueueOpt) Cap() int {
	return q.cap
}

func (q *CircularQueueOpt) valueAt(i int) *int {
	sizeInt := (int)(unsafe.Sizeof(int(0)))
	return (*int)(unsafe.Add(q.values, sizeInt*(i%q.Cap())))
}

func TestCircularQueueOpt(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueueOpt(queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	// assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	// assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
