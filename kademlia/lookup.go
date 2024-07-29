package kademlia

import (
	"container/heap"
	"fmt"
	"log"
	"math/big"
	"sync"
)

type Lookup struct {
	initiator Server
	key       *big.Int
	shortlist *Shortlist
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

func (lu *Lookup) Execute() []Node {
	initialNodes := lu.initiator.routingTable.GetNearest(lu.key)
	lu.shortlist.Insert(initialNodes...)
	numSeenNodes := lu.shortlist.Len()

	for {
		lu.queryNodes()
		numNewNodes := lu.shortlist.Len() - numSeenNodes
		if numNewNodes == 0 || lu.rounds == Options.MaxIterations {
			break
		}
		numSeenNodes = lu.shortlist.Len()
		lu.rounds++
		fmt.Println("seen:", numSeenNodes)
	}
	return lu.shortlist.Closest()
}

func (lu *Lookup) queryNodes() {
	alpha := Options.Alpha
	wg := sync.WaitGroup{}
	wg.Add(alpha)
	for i := 0; i < alpha; i++ {
		go func() {
			var p Node
			for _, n := range lu.shortlist.heap.Nodes {
				if !lu.shortlist.queriedNodes.Has(n) {
					p = n
					break
				}
			}
			list, err := lu.initiator.sendFindNode(lu.key.Text(16), p)
			if err != nil {
				log.Println(err)
				lu.shortlist.Remove(p)
			}

			lu.shortlist.queriedNodes.Add(p)
			lu.shortlist.Insert(list...)
			wg.Done()
		}()
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
}

func NewShortlist(key *big.Int) *Shortlist {
	return &Shortlist{
		key:          key,
		heap:         &NodeHeap{Key: key},
		seenNodes:    &NodeSet{},
		queriedNodes: &NodeSet{},
	}
}

func (sl *Shortlist) Insert(node ...Node) {
	sl.m.Lock()
	defer sl.m.Unlock()
	for _, n := range node {
		if !sl.seenNodes.Has(n) {
			sl.seenNodes.Add(n)
			heap.Push(sl.heap, n)
		}
	}
}

func (sl *Shortlist) Remove(node Node) {
	sl.m.Lock()
	defer sl.m.Unlock()
	sl.seenNodes.Remove(node)
	sl.queriedNodes.Remove(node)
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
	return len(*sl.seenNodes)
}

func (sl *Shortlist) Pop() Node {
	sl.m.Lock()
	defer sl.m.Unlock()

	return heap.Pop(sl.heap).(Node)
}
