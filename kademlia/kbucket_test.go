package kademlia

import (
	"testing"
)

func TestKBucket_Add(t *testing.T) {
	var nodes []Node
	for i := 0; i < 10; i++ {
		nodes = append(nodes, NewNode("localhost", 8000+i, nil))
	}

	kb := KBucket{Owner: Node{}, Capacity: 11}
	for _, node := range nodes {
		kb.Add(node)
	}

	if !kb.isTail(nodes[len(nodes)-1]) {
		t.Errorf("Latest inserted node must be at the tail")
	}
	if kb.Size != len(nodes) {
		t.Errorf("Size of k bucket should be equal to number of unique nodes")
	}

	kb.Add(nodes[0])
	if !kb.isTail(nodes[0]) {
		t.Errorf("Latest inserted node must be at the tail")
	}
	if kb.Size != len(nodes) {
		t.Errorf("Inserting already added node should not increase size")
	}

	kb.Add(NewNode("localhost", nodes[len(nodes)-1].Port+1, nil))
	kb.Add(NewNode("localhost", nodes[len(nodes)-1].Port+2, nil))
	if kb.Size != 11 {
		t.Errorf("Inserting a new node should not increase size when at capacity")
	}
}

func TestKBucket_Remove(t *testing.T) {
	var nodes []Node
	numNodes := 8
	for i := 0; i < numNodes; i++ {
		nodes = append(nodes, NewNode("localhost", 8000+i, nil))
	}

	kb := KBucket{Owner: Node{}, Capacity: numNodes}
	for _, node := range nodes {
		kb.Add(node)
	}

	kb.remove(nodes[len(nodes)/2])
	if kb.contains(nodes[len(nodes)/2]) {
		t.Errorf("Node within the list should be successfully removed")
	}

	kb.remove(nodes[0])
	if !kb.isHead(nodes[1]) {
		t.Errorf("The head should be successfully removed, with the next node as the new head")
	}

	kb.remove(nodes[len(nodes)-1])
	if !kb.isTail(nodes[len(nodes)-2]) {
		t.Errorf("The tail should be successfully removed, with the previous node as the new tail")
	}

	if kb.Size != numNodes-3 {
		t.Errorf("Bucket size should have been decreased by 3")
	}
}
