package kademlia

import (
	"fmt"
	"math/big"
	"math/rand"
	"time"
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
	lastUsed time.Time
}

func NewKBucket(owner Node, prefix string) *KBucket {
	return &KBucket{
		Owner:    owner,
		Capacity: Options.BucketCapacity,
		Prefix:   prefix,
		lastUsed: time.Now(),
	}
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
	kb.lastUsed = time.Now()
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
		if ptr.Data.Equals(n) {
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
		if ptr.Data.Equals(n) {
			return true
		}
		ptr = ptr.Next
	}
	return false
}

func (kb *KBucket) isUnderpopulated() bool {
	return kb.Size <= Options.BucketCapacity/2
}

func (kb *KBucket) wasRecentlyUsed() bool {
	return int(time.Since(kb.lastUsed).Seconds()) <= Options.TRefresh
}

func (kb *KBucket) shouldBeRefreshed() bool {
	return !kb.wasRecentlyUsed() || kb.isUnderpopulated()
}

func (kb *KBucket) randomNum() *big.Int {
	curr := kb.Prefix
	for i := len(curr); i < 160; i++ {
		curr += []string{"0", "1"}[rand.Intn(2)]
	}
	val, ok := new(big.Int).SetString(curr, 2)
	if !ok {
		return nil
	}
	return val
}
