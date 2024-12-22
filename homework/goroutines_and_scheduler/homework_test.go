package main

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	heap []Task
}

func NewScheduler() Scheduler {
	return Scheduler{heap: make([]Task, 0)}
}

func (s *Scheduler) AddTask(task Task) {
	s.heap = append(s.heap, task)
	s.pushTop(len(s.heap) - 1)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	heapIdx := s.findByTaskID(taskID)

	prevPriority := s.heap[heapIdx].Priority
	s.heap[heapIdx].Priority = newPriority

	if newPriority > prevPriority {
		s.pushTop(heapIdx)
	} else {
		s.pushDown(heapIdx)
	}
}

func (s *Scheduler) GetTask() Task {
	s.heap[0], s.heap[len(s.heap)-1] = s.heap[len(s.heap)-1], s.heap[0]
	task := s.heap[len(s.heap)-1]
	s.heap = s.heap[:len(s.heap)-1]
	s.pushDown(0)
	return task
}

// не придумал и не нашел способа сделать это не за O(n)
// в целом можно хранить индексы в мапе, но это тяжело поддерживать в pushTop/pushDown
func (s *Scheduler) findByTaskID(taskID int) int {
	for i, t := range s.heap {
		if t.Identifier == taskID {
			return i
		}
	}

	return -1
}

func (s *Scheduler) pushTop(i int) {
	parentIdx := (i - 1) / 2
	for i > 0 && s.heap[parentIdx].Priority < s.heap[i].Priority {
		s.heap[parentIdx], s.heap[i] = s.heap[i], s.heap[parentIdx]
		i, parentIdx = parentIdx, (parentIdx-1)/2

	}
}

func (s *Scheduler) pushDown(i int) {
	maxByIdx := func(left, right int) int {
		n := len(s.heap)
		if left >= n && right >= n {
			return -1
		}

		if right < n && s.heap[left].Priority < s.heap[right].Priority {
			return right
		}
		return left
	}

	for {
		maxIdx := maxByIdx(2*i+1, 2*i+2)
		if maxIdx == -1 || s.heap[i].Priority >= s.heap[maxIdx].Priority {
			break
		}

		s.heap[i], s.heap[maxIdx] = s.heap[maxIdx], s.heap[i]
		i = maxIdx
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
