package main

import (
	"fmt"
	"maps"
	"slices"
	"testing"
	"unsafe"
	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Set map[uintptr]struct{}

func (s Set) WasSeen(p uintptr) bool {
	_, ok := s[p]
	return ok
}

func (s Set) Add(p uintptr) {
	s[p] = struct{}{}
}

var ptrs = make(Set)

func visit(p uintptr) {
	if p == 0x00 || ptrs.WasSeen(p) {
		return
	}
	
	ptrs.Add(p)
	nestedP := (*uintptr)(unsafe.Pointer(p))
	visit(*nestedP)

}

func Trace(stacks [][]uintptr) []uintptr {
	for i := range stacks {
		for j := range stacks[i] {
			visit(stacks[i][j])
		}
	}
	return slices.Collect(maps.Keys(ptrs))
}

func TestTrace(t *testing.T) {
	var heapObjects = []int{
		0x00, 0x00, 0x00, 0x00, 0x00,
	}

	var heapPointer1 *int = &heapObjects[1]
	var heapPointer2 *int = &heapObjects[2]
	var heapPointer3 *int = nil
	var heapPointer4 **int = &heapPointer3

	var stacks = [][]uintptr{
		{
			uintptr(unsafe.Pointer(&heapPointer1)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[0])),
			0x00, 0x00, 0x00, 0x00,
		},
		{
			uintptr(unsafe.Pointer(&heapPointer2)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[1])),
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[2])),
			uintptr(unsafe.Pointer(&heapPointer4)), 0x00, 0x00, 0x00,
		},
		{
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[3])),
		},
	}

	pointers := Trace(stacks)
	expectedPointers := []uintptr{
		uintptr(unsafe.Pointer(&heapPointer1)),
		uintptr(unsafe.Pointer(&heapObjects[0])),
		uintptr(unsafe.Pointer(&heapPointer2)),
		uintptr(unsafe.Pointer(&heapObjects[1])),
		uintptr(unsafe.Pointer(&heapObjects[2])),
		uintptr(unsafe.Pointer(&heapPointer4)),
		uintptr(unsafe.Pointer(&heapPointer3)),
		uintptr(unsafe.Pointer(&heapObjects[3])),
	}
	assert.ElementsMatch(t, expectedPointers, pointers)
}
