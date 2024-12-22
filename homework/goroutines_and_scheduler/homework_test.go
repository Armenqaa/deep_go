package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type KeyWithPriority struct {
	ID int
	P  int
}

type Heap struct {
	data     []KeyWithPriority
	keyToIdx map[int]int
}

func NewHeap() Heap {
	return Heap{keyToIdx: make(map[int]int), data: make([]KeyWithPriority, 0)}
}

func (h *Heap) Add(kv KeyWithPriority) {
	h.data = append(h.data, kv)
	h.keyToIdx[kv.ID] = len(h.data) - 1
	h.pushTop(len(h.data) - 1)
}

func (h *Heap) GetTop() KeyWithPriority {
	top := h.data[0]
	delete(h.keyToIdx, h.data[0].ID)

	h.data[0], h.data[len(h.data)-1] = h.data[len(h.data)-1], h.data[0]
	h.data = h.data[:len(h.data)-1]

	h.keyToIdx[h.data[0].ID] = 0
	h.pushDown(0)

	return top
}

func (h *Heap) ChangePriority(key int, newPriority int) {
	heapIdx := h.keyToIdx[key]

	prevPriority := h.data[heapIdx].P
	h.data[heapIdx].P = newPriority

	if newPriority > prevPriority {
		h.pushTop(heapIdx)
	} else {
		h.pushDown(heapIdx)
	}
}

func (h *Heap) pushTop(i int) {
	parentIdx := (i - 1) / 2
	for i > 0 && h.data[parentIdx].P < h.data[i].P {
		h.data[parentIdx], h.data[i] = h.data[i], h.data[parentIdx]
		h.keyToIdx[h.data[i].ID] = i
		h.keyToIdx[h.data[parentIdx].ID] = parentIdx
		i, parentIdx = parentIdx, (parentIdx-1)/2
	}
}

func (h *Heap) pushDown(i int) {
	maxByIdx := func(left, right int) int {
		n := len(h.data)
		if left >= n && right >= n {
			return -1
		}

		if right < n && h.data[left].P < h.data[right].P {
			return right
		}
		return left
	}

	for {
		maxIdx := maxByIdx(2*i+1, 2*i+2)
		if maxIdx == -1 || h.data[i].P >= h.data[maxIdx].P {
			break
		}

		h.data[i], h.data[maxIdx] = h.data[maxIdx], h.data[i]
		h.keyToIdx[h.data[i].ID] = i

		h.keyToIdx[h.data[maxIdx].ID] = maxIdx
		i = maxIdx
	}
}

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	heap Heap
}

func NewScheduler() Scheduler {
	return Scheduler{heap: NewHeap()}
}

func (s *Scheduler) AddTask(task Task) {
	s.heap.Add(KeyWithPriority{ID: task.Identifier, P: task.Priority})
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	s.heap.ChangePriority(taskID, newPriority)
}

func (s *Scheduler) GetTask() Task {
	kv := s.heap.GetTop()
	return Task{
		Identifier: kv.ID,
		Priority:   kv.P,
	}
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	task1.Priority = 100
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
