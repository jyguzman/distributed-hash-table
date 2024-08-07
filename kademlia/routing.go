package kademlia

import (
	"container/heap"
	"math/big"
	"strings"
)

type RTNode struct {
	Bucket  *KBucket
	Left    *RTNode
	Right   *RTNode
	K       int
	Prefix  string
	RtOwner Node
}

func NewRTNode(owner Node) *RTNode {
	return &RTNode{
		Bucket:  NewKBucket(owner, ""),
		Left:    nil,
		Right:   nil,
		K:       Options.BucketCapacity,
		Prefix:  "",
		RtOwner: owner,
	}
}

func (rn *RTNode) StringHelper(level int) string {
	tabs := strings.Repeat("\t", level)
	if rn == nil {
		return tabs + "<nil>"
	}
	if rn.isLeaf() {
		return tabs + rn.Prefix + ": " + rn.Bucket.String()
	}
	prefix, left, right := "*", tabs+"nil", tabs+"nil"
	if len(rn.Prefix) > 0 {
		prefix = tabs + rn.Prefix
	}
	if rn.Left != nil {
		left = tabs + rn.Left.StringHelper(level+1)
	}
	if rn.Right != nil {
		right = tabs + rn.Right.StringHelper(level+1)
	}
	return prefix + " \n" + left + " \n" + right
}

func (rn *RTNode) String() string {
	return rn.StringHelper(0)
}

func (rn *RTNode) Split(prefixes map[string]*KBucket) {
	prfx := rn.Prefix
	zeroBucket, oneBucket := NewKBucket(rn.RtOwner, prfx+"0"), NewKBucket(rn.RtOwner, prfx+"1")
	ptr, pLen := rn.Bucket.Tail, len(prfx)
	for ptr != nil {
		currId := ptr.Data.Id
		bit := currId.Bit(159 - pLen)
		if bit == 0 {
			zeroBucket.Add(ptr.Data)
		} else {
			oneBucket.Add(ptr.Data)
		}
		ptr = ptr.Prev
	}
	rn.Bucket = nil
	delete(prefixes, prfx)
	rn.Left = &RTNode{RtOwner: rn.RtOwner, Bucket: zeroBucket, K: Options.BucketCapacity, Prefix: prfx + "0"}
	rn.Right = &RTNode{RtOwner: rn.RtOwner, Bucket: oneBucket, K: Options.BucketCapacity, Prefix: prfx + "1"}
	prefixes[rn.Left.Prefix] = rn.Left.Bucket
	prefixes[rn.Right.Prefix] = rn.Right.Bucket
}

func (rn *RTNode) isLeaf() bool {
	return rn.Left == nil && rn.Right == nil && rn.Bucket != nil
}

func (rn *RTNode) Add(currPos int, node Node, prefixes map[string]*KBucket) int {
	if rn.isLeaf() {
		if rn.Bucket.Size < rn.K {
			rn.Bucket.Add(node)
			return 1
		}
		prefix := rn.Bucket.Prefix
		if prefix == rn.RtOwner.Prefix(len(prefix)) {
			rn.Split(prefixes)
			return rn.Add(currPos, node, prefixes)
		}
	} else {
		bit := node.Id.Bit(159 - currPos)
		if bit == 0 {
			return rn.Left.Add(currPos+1, node, prefixes)
		} else {
			return rn.Right.Add(currPos+1, node, prefixes)
		}
	}
	return 0
}

type RoutingTable struct {
	Owner          Node
	K              int
	Root           *RTNode
	Size           int
	BucketPrefixes map[string]*KBucket
}

func (rt *RoutingTable) String() string {
	return rt.Root.String()
}

func NewRoutingTable(owner Node, k int) *RoutingTable {
	rt := &RoutingTable{
		Owner:          owner,
		K:              k,
		Root:           NewRTNode(owner),
		BucketPrefixes: make(map[string]*KBucket),
	}
	rt.BucketPrefixes[""] = rt.Root.Bucket
	return rt
}

func (rt *RoutingTable) Add(node Node) {
	rt.Size += rt.Root.Add(0, node, rt.BucketPrefixes)
}

func (rt *RoutingTable) GetNearest(key *big.Int) []Node {
	nodeHeap := &NodeHeap{Key: key}
	heap.Init(nodeHeap)
	for _, bucket := range rt.BucketPrefixes {
		ptr := bucket.Tail
		for ptr != nil {
			if ptr.Data.Id.Cmp(key) != 0 {
				heap.Push(nodeHeap, ptr.Data)
			}
			ptr = ptr.Prev
		}
	}
	var nodes []Node
	for i := 0; len(nodeHeap.Nodes) > 0 && i < rt.K; i++ {
		nodes = append(nodes, heap.Pop(nodeHeap).(Node))
	}
	return nodes
}
