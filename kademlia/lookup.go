package kademlia

import (
	"container/heap"
	"log"
	"math/big"
	"sync"
)

type Lookup struct {
	initiator Server
	key       *big.Int
	shortlist *Shortlist
	m         sync.Mutex
	rounds    int
	value     any
}

func NewLookup(initiator Server, key *big.Int) *Lookup {
	return &Lookup{
		initiator: initiator,
		key:       key,
		shortlist: NewShortlist(key),
	}
}

func (lu *Lookup) mark(n Node) {
	lu.m.Lock()
	lu.shortlist.queriedNodes.Add(n)
	lu.m.Unlock()
}

func (lu *Lookup) hasBeenQueried(n Node) bool {
	lu.m.Lock()
	defer lu.m.Unlock()
	return lu.shortlist.queriedNodes.Has(n)
}

func (lu *Lookup) Execute() []Node {
	initNodes := lu.initiator.routingTable.GetNearest(lu.key)
	lu.shortlist.Insert(initNodes...)
	numSeenNodes := lu.shortlist.Len()

	for {
		lu.sendRequests(lu.shortlist.GetNextAlpha())
		numNewNodes := lu.shortlist.Len() - numSeenNodes
		if numNewNodes == 0 || lu.rounds == Options.MaxIterations {
			break
		}
		numSeenNodes = lu.shortlist.Len()
		lu.rounds++
	}

	return lu.shortlist.Closest()
}

func (lu *Lookup) sendRequests(nodes []Node) {
	m, wg := sync.Mutex{}, sync.WaitGroup{}
	for _, n := range nodes {
		wg.Add(1)
		go func(n Node) {
			defer wg.Done()
			if !lu.hasBeenQueried(n) {
				list, err := lu.initiator.sendFindNode(lu.key.Text(16), n)
				if err != nil {
					log.Println(err)
					lu.shortlist.Remove(n)
					return
				}
				m.Lock()
				lu.mark(n)
				lu.shortlist.Insert(list...)
				m.Unlock()
			}
		}(n)
	}
	wg.Wait()
}

type NodeSet map[string]bool

func (ns *NodeSet) Add(n Node) {
	if !ns.Has(n) {
		(*ns)[n.Id.Text(16)] = true
	}
}

func (ns *NodeSet) Remove(n Node) {
	delete(*ns, n.Id.Text(16))
}

func (ns *NodeSet) Has(n Node) bool {
	_, ok := (*ns)[n.Id.Text(16)]
	return ok
}

type Shortlist struct {
	key          *big.Int
	m            sync.Mutex
	heap         *NodeHeap
	queriedNodes *NodeSet
	seenNodes    *NodeSet
	candidates   []Node
	index        int
}

func NewShortlist(key *big.Int) *Shortlist {
	return &Shortlist{
		key:          key,
		heap:         &NodeHeap{Key: key},
		seenNodes:    &NodeSet{},
		queriedNodes: &NodeSet{},
		candidates:   []Node{},
		index:        0,
	}
}

func (sl *Shortlist) Insert(node ...Node) {
	sl.m.Lock()
	defer sl.m.Unlock()

	for _, n := range node {
		if !sl.seenNodes.Has(n) {
			sl.seenNodes.Add(n)
			heap.Push(sl.heap, n)
			sl.candidates = append(sl.candidates, n)
		}
	}
}

func (sl *Shortlist) Remove(node Node) {
	sl.m.Lock()
	defer sl.m.Unlock()

	sl.seenNodes.Remove(node)
	sl.queriedNodes.Remove(node)
}

func (sl *Shortlist) GetNextAlpha() []Node {
	sl.m.Lock()
	defer sl.m.Unlock()

	var nodes []Node
	end := sl.index + Options.Alpha
	for i := sl.index; sl.index < len(sl.candidates) && i < end; i++ {
		if !sl.queriedNodes.Has(sl.candidates[i]) {
			nodes = append(nodes, sl.candidates[i])
			sl.index++
		}
	}

	return nodes
}

func (sl *Shortlist) Closest() []Node {
	sl.m.Lock()
	defer sl.m.Unlock()

	var closestNodes []Node
	for i := 0; len(sl.heap.Nodes) > 0 && i < Options.BucketCapacity; i++ {
		n := heap.Pop(sl.heap).(Node)
		if sl.seenNodes.Has(n) {
			closestNodes = append(closestNodes, n)
		}
	}

	return closestNodes
}

func (sl *Shortlist) Len() int {
	sl.m.Lock()
	defer sl.m.Unlock()

	return len(*sl.seenNodes)
}

func (sl *Shortlist) Pop() Node {
	sl.m.Lock()
	defer sl.m.Unlock()

	return heap.Pop(sl.heap).(Node)
}
