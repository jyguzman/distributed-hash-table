package kademlia

import (
	"fmt"
)

type ListNode struct {
	Data Node
	Next *ListNode
	Prev *ListNode
}

type KBucket struct {
	Owner    Node
	Capacity int
	Head     *ListNode
	Tail     *ListNode
	Size     int
	Prefix   string
}

func NewKBucket(owner Node, prefix string) *KBucket {
	return &KBucket{Owner: owner, Capacity: Options.BucketCapacity, Prefix: prefix}
}

func (kb *KBucket) String() string {
	if kb.Head == nil {
		return "<empty>"
	}
	ptr := kb.Head
	res := ""
	for ptr != nil {
		if ptr.Next != nil {
			res += fmt.Sprintf("%s:%d <-> ", ptr.Data.Host, ptr.Data.Port)
		} else {
			res += fmt.Sprintf("%s:%d ", ptr.Data.Host, ptr.Data.Port)
		}
		ptr = ptr.Next
	}
	return res
}

func (kb *KBucket) Add(n Node) {
	if kb.Capacity == kb.Size {
		return
	}
	if kb.Head == nil {
		kb.Head = &ListNode{Data: n}
		kb.Tail = kb.Head
		kb.Size++
		return
	}
	if kb.contains(n) {
		kb.remove(n)
	}
	newListNode := &ListNode{Data: n}
	kb.Tail.Next = newListNode
	newListNode.Prev = kb.Tail
	kb.Tail = newListNode
	kb.Size++
}

func (kb *KBucket) isTail(n Node) bool {
	return kb.Tail != nil && kb.Tail.Data.Id.Cmp(n.Id) == 0
}

func (kb *KBucket) isHead(n Node) bool {
	return kb.Head != nil && kb.Head.Data.Id.Cmp(n.Id) == 0
}

func (kb *KBucket) remove(n Node) {
	if kb.Head == nil || kb.Tail == nil {
		return
	}
	if kb.Size == 1 {
		kb.Head = nil
		kb.Tail = nil
		kb.Size = 0
		return
	}
	if kb.isTail(n) {
		kb.Tail = kb.Tail.Prev
		kb.Tail.Next = nil
		kb.Size--
		return
	}
	if kb.isHead(n) {
		kb.Head = kb.Head.Next
		kb.Head.Prev = nil
		kb.Size--
		return
	}
	ptr := kb.Head
	for ptr != nil {
		if ptr.Data.Id.Cmp(n.Id) == 0 {
			ptr.Prev.Next = ptr.Next
			ptr.Next.Prev = ptr.Prev
			kb.Size--
			return
		}
		ptr = ptr.Next
	}
}

func (kb *KBucket) contains(n Node) bool {
	ptr := kb.Head
	for ptr != nil {
		if ptr.Data.Id.Cmp(n.Id) == 0 {
			return true
		}
		ptr = ptr.Next
	}
	return false
}

func (kb *KBucket) IsUnderpopulated() bool {
	return kb.Size <= Options.BucketCapacity/2
}
