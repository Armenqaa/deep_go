import (
	// "reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type CircularQueue struct {
	values             unsafe.Pointer
	start, filled, cap int
}

func NewCircularQueue(size int) CircularQueue {
	values := make([]int, size)
	return CircularQueue{
		values: unsafe.Pointer(&values[0]),
		cap:    size,
	}
}

func (q *CircularQueue) Push(value int) bool {
	if q.Full() {
		return false
	}
	old := q.valueAt(q.start+q.filled)
	*old = value
	q.filled++
	return true

}

func (q *CircularQueue) Pop() bool {
	if q.Empty() {
		return false
	}
	q.start++
	q.filled--
	return true
}

func (q *CircularQueue) Front() int {
	if q.Empty() {
		return -1
	}
	return *q.valueAt(q.start)
}

func (q *CircularQueue) Back() int {
	if q.Empty() {
		return -1
	}
	return *q.valueAt(q.start+q.filled-1)
}

func (q *CircularQueue) Empty() bool {
	return q.filled == 0
}

func (q *CircularQueue) Full() bool {
	return q.filled == q.Cap()
}

func (q *CircularQueue) Cap() int {
	return q.cap
}

func (q *CircularQueue) valueAt(i int) *int {
	sizeInt := (int)(unsafe.Sizeof(int(0)))
	return (*int)(unsafe.Add(q.values, sizeInt*(i%q.Cap())))
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue(queueSize)

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
