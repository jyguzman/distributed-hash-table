package kademlia

import "fmt"

type ListNode struct {
	Node Node
	Next *ListNode
	Prev *ListNode
}

type KBucket struct {
	K    int
	Head *ListNode
	Tail *ListNode
	Size int
}

func NewKBucket(k int) *KBucket {
	return &KBucket{K: k}
}

func (kb *KBucket) String() string {
	if kb.Head == nil {
		return "empty"
	}
	ptr := kb.Head
	res := ""
	for ptr != nil {
		if ptr.Next != nil {
			res += fmt.Sprintf("%s:%d -> ", ptr.Node.IP, ptr.Node.Port)
		} else {
			res += fmt.Sprintf("%s:%d ", ptr.Node.IP, ptr.Node.Port)
		}
		ptr = ptr.Next
	}
	return res
}

func (kb *KBucket) Append(n Node) {
	kb.Size += 1
	if kb.Head == nil {
		kb.Head = &ListNode{Node: n}
		kb.Tail = kb.Head
		return
	}
	newListNode := &ListNode{Node: n}
	newListNode.Next = kb.Head
	kb.Head = newListNode
}
