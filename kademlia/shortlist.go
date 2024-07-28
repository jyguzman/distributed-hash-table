package kademlia

import (
	"container/heap"
	"math/big"
	"sync"
)

type Shortlist struct {
	m         sync.Mutex
	heap      NodeHeap
	seenNodes map[Node]bool
}

func NewShortlist(key *big.Int) *Shortlist {
	heap.Init(&NodeHeap{Key: key})
	return &Shortlist{
		heap:      NodeHeap{Key: key},
		seenNodes: make(map[Node]bool),
	}
}

func (sl *Shortlist) Insert(node ...Node) {
	sl.m.Lock()
	for _, n := range node {
		_, seen := sl.seenNodes[n]
		if seen {
			continue
		}
		heap.Push(&sl.heap, n)
	}
	sl.m.Unlock()
}
