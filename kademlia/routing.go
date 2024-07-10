package kademlia

import (
	"fmt"
	"math/big"
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
		Prefix: "",
	}
}

func (rn *RTNode) StringHelper(level int) string {
	tabs := ""
	for i := 0; i < level; i++ {
		tabs += "\t"
	}
	if rn == nil {
		return tabs + "<nil>"
	}
	if rn.isLeaf() {
		return tabs + rn.Prefix + "\n" + tabs + tabs + rn.Bucket.String()
	}
	pre, left, right := "*", tabs+"nil", tabs+"nil"
	if len(rn.Prefix) > 0 {
		pre = tabs + rn.Prefix
	}
	if rn.Left != nil {
		left = tabs + rn.Left.StringHelper(level+1)
	}
	if rn.Right != nil {
		right = tabs + rn.Right.StringHelper(level+1)
	}
	return fmt.Sprintf("%s \n%s \n%s", pre, left, right)
}

func (rn *RTNode) String() string {
	return rn.StringHelper(0)
}

func (rn *RTNode) Split() {
	zeroBucket, oneBucket := NewKBucket(rn.K), NewKBucket(rn.K)
	ptr := rn.Bucket.Head
	preLen := 0
	if len(rn.Prefix) > 0 {
		preLen = len(rn.Prefix)
	}
	for ptr != nil {
		currId := ptr.Node.ID
		bit := int(currId.Bit(currId.BitLen() - preLen - 1))
		if bit == 0 {
			zeroBucket.Append(ptr.Node)
		} else {
			oneBucket.Append(ptr.Node)
		}
		ptr = ptr.Next
	}
	rn.Bucket = nil
	rn.Left = &RTNode{Bucket: zeroBucket, K: rn.K, Prefix: rn.Prefix + "0"}
	rn.Right = &RTNode{Bucket: oneBucket, K: rn.K, Prefix: rn.Prefix + "1"}
}

func (rn *RTNode) isLeaf() bool {
	return rn.Left == nil && rn.Right == nil && rn.Bucket != nil
}

func (rn *RTNode) Add(currPos int, node Node) {
	if rn.isLeaf() {
		if rn.Bucket.Size < rn.K {
			rn.Bucket.Append(node)
			return
		}
		rn.Split()
	}
	bit := int(node.ID.Bit(node.ID.BitLen() - currPos - 1))
	if bit == 0 {
		rn.Left.Add(currPos+1, node)
	} else {
		rn.Right.Add(currPos+1, node)
	}
}

type RoutingTable struct {
	K        int
	Root     *RTNode
	Size     int
	Prefixes map[string]bool
}

func (rt *RoutingTable) String() string {
	return rt.Root.String()
}

func NewRoutingTable(k int) *RoutingTable {
	return &RoutingTable{K: k, Root: NewRTNode(k)}
}

func (rt *RoutingTable) Add(node Node) {
	rt.Size += 1
	rt.Root.Add(1, node)
}

func (rt *RoutingTable) KNearestHelper(key *big.Int, rn *RTNode, nodes []Node, k int) []Node {
	if len(nodes) == k {
		return nodes
	}
	return nodes
}

func (rt *RoutingTable) KNearest(key *big.Int, k int) []Node {
	return rt.KNearestHelper(key, rt.Root, []Node{}, k)
}
