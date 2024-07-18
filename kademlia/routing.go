package kademlia

import (
	"container/heap"
	"go-dht/bson"
	"math/big"
	"strings"
)

type RTNode struct {
	Bucket *KBucket
	Left   *RTNode
	Right  *RTNode
	K      int
	Prefix string
}

func NewRTNode(k int) *RTNode {
	return &RTNode{
		Bucket: &KBucket{},
		Left:   nil,
		Right:  nil,
		K:      k,
		Prefix: "*",
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
	zeroBucket, oneBucket := NewKBucket(rn.K), NewKBucket(rn.K)
	ptr := rn.Bucket.Head
	for ptr != nil {
		currId := ptr.Node.ID
		bit := int(currId.Bit(currId.BitLen() - len(rn.Prefix) - 1))
		if bit == 0 {
			zeroBucket.Add(ptr.Node)
		} else {
			oneBucket.Add(ptr.Node)
		}
		ptr = ptr.Next
	}
	rn.Bucket = nil
	delete(prefixes, rn.Prefix)
	rn.Left = &RTNode{Bucket: zeroBucket, K: rn.K, Prefix: rn.Prefix + "0"}
	rn.Right = &RTNode{Bucket: oneBucket, K: rn.K, Prefix: rn.Prefix + "1"}
	prefixes[rn.Left.Prefix] = rn.Left.Bucket
	prefixes[rn.Right.Prefix] = rn.Right.Bucket
}

func (rn *RTNode) isLeaf() bool {
	return rn.Left == nil && rn.Right == nil && rn.Bucket != nil
}

func (rn *RTNode) Add(currPos int, node Node, prefixes map[string]*KBucket) {
	if rn.isLeaf() {
		if rn.Bucket.Size < rn.K {
			rn.Bucket.Add(node)
			return
		}
		rn.Split(prefixes)
	}
	bit := int(node.ID.Bit(node.ID.BitLen() - currPos - 1))
	if bit == 0 {
		rn.Left.Add(currPos+1, node, prefixes)
	} else {
		rn.Right.Add(currPos+1, node, prefixes)
	}
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
		Root:           NewRTNode(k),
		BucketPrefixes: make(map[string]*KBucket),
	}
	rt.BucketPrefixes["*"] = rt.Root.Bucket
	return rt
}

func (rt *RoutingTable) Add(node Node) {
	rt.Root.Add(1, node, rt.BucketPrefixes)
	rt.Size += 1
}

func (rt *RoutingTable) GetNearest(key *big.Int) bson.A {
	nodeHeap := &NodeHeap{Key: key}
	heap.Init(nodeHeap)
	for _, bucket := range rt.BucketPrefixes {
		ptr := bucket.Head
		for ptr != nil {
			if ptr.Node.ID.Cmp(key) != 0 {
				heap.Push(nodeHeap, ptr.Node)
			}
			ptr = ptr.Next
		}
	}
	nodes := bson.A{}
	for i := 0; len(nodeHeap.Nodes) > 0 && i < rt.K; i++ {
		nodes = append(nodes, heap.Pop(nodeHeap).(Node).Tuple())
	}
	return nodes
}
