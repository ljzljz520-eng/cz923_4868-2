package queue

import "sync"

type Numbering struct {
	mu   sync.Mutex
	next int
}

func New(start int) *Numbering {
	if start < 1 {
		start = 1
	}
	return &Numbering{next: start}
}
func (n *Numbering) Next() int { n.mu.Lock(); defer n.mu.Unlock(); v := n.next; n.next++; return v }
func (n *Numbering) Peek() int { n.mu.Lock(); defer n.mu.Unlock(); return n.next }
func (n *Numbering) Reset(start int) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if start < 1 {
		start = 1
	}
	n.next = start
}
func (n *Numbering) Reserve(count int) []int {
	if count < 0 {
		count = 0
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	out := make([]int, count)
	for i := range out {
		out[i] = n.next
		n.next++
	}
	return out
}
