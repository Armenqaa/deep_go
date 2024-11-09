package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TNode struct {
	left, right *TNode
	key, value  int
}

type OrderedMap struct {
	head *TNode
	len  int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{
		head: nil,
		len:  0,
	}
}

func (m *OrderedMap) Insert(key, value int) {
	toInsert := TNode{key: key, value: value}
	if m.len == 0 {
		*m = OrderedMap{&toInsert, 1}
		return
	}
	m.insert(m.head, &toInsert)
}

func (m *OrderedMap) Erase(key int) {
	if m.len == 0 {
		// empty
		return
	}

	m.head = m.delete(m.head, key)
}

func (m *OrderedMap) Contains(key int) bool {
	return m.find(m.head, key) != nil
}

func (m *OrderedMap) Size() int {
	return m.len
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	m.inOrderAction(m.head, action)
}

func (m *OrderedMap) inOrderAction(curr *TNode, action func(int, int)) {
	if curr == nil {
		return
	}
	m.inOrderAction(curr.left, action)
	action(curr.key, curr.value)
	m.inOrderAction(curr.right, action)
}

func (m *OrderedMap) delete(curr *TNode, key int) *TNode {
	if curr.key < key {
		curr.right = m.delete(curr.right, key)
	} else if curr.key > key {
		curr.left = m.delete(curr.left, key)
	} else {
		// key found
		m.len--

		childrenCount := 0
		if curr.left != nil {
			childrenCount++
		}
		if curr.right != nil {
			childrenCount++
		}

		switch childrenCount {
		case 0:
			return nil
		case 1:
			if curr.left != nil {
				return curr.left
			}
			return curr.right
		case 2:
			min := m.findMin(curr.right)
			curr.key, curr.value = min.key, min.value
			curr.right = m.delete(curr.right, min.key)
		}
	}
	return curr
}

func (m *OrderedMap) find(curr *TNode, key int) *TNode {
	if curr == nil {
		return nil
	}
	if curr.key == key {
		return curr
	}

	if curr.key < key {
		return m.find(curr.right, key)
	}
	return m.find(curr.left, key)
}

func (m *OrderedMap) findMin(curr *TNode) *TNode {
	if curr.left == nil {
		return curr
	}
	return m.findMin(curr.left)
}

func (m *OrderedMap) insert(curr *TNode, toInsert *TNode) *TNode {
	if curr == nil {
		m.len++
		return toInsert
	}

	if curr.key > toInsert.key {
		curr.left = m.insert(curr.left, toInsert)
	} else if curr.key < toInsert.key {
		curr.right = m.insert(curr.right, toInsert)
	}
	return curr
}

func TestOrderedMap(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
