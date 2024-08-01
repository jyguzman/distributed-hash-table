package kademlia

import (
	"math/big"
)

type NodeHeap struct {
	Key   *big.Int
	Nodes []Node
	isMin bool
}

func (h *NodeHeap) Len() int { return len(h.Nodes) }

func (h *NodeHeap) Less(i, j int) bool {
	iDistToKey := new(big.Int).Xor(h.Key, h.Nodes[i].Id)
	jDistToKey := new(big.Int).Xor(h.Key, h.Nodes[j].Id)
	return iDistToKey.Cmp(jDistToKey) == -1
}
func (h *NodeHeap) Swap(i, j int) { h.Nodes[i], h.Nodes[j] = h.Nodes[j], h.Nodes[i] }

func (h *NodeHeap) Push(x any) {
	h.Nodes = append(h.Nodes, x.(Node))
}

func (h *NodeHeap) Pop() any {
	old := h.Nodes
	n := len(old)
	x := old[n-1]
	h.Nodes = old[0 : n-1]
	return x
}

func (h *NodeHeap) Top() Node {
	return h.Nodes[0]
}
