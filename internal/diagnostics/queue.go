package diagnostics

import (
	"maps"
	"math"
)

type record struct {
	event    Event
	sequence uint64
	size     int
}
type recordQueue struct {
	entries            []record
	head, count, bytes int
}

func newQueue(n int) recordQueue { return recordQueue{entries: make([]record, n)} }
func (q *recordQueue) push(r record) {
	q.entries[(q.head+q.count)%len(q.entries)] = r
	q.count++
	q.bytes += r.size
}
func (q *recordQueue) first() record   { return q.entries[q.head] }
func (q *recordQueue) at(i int) record { return q.entries[(q.head+i)%len(q.entries)] }
func (q *recordQueue) pop() record {
	r := q.first()
	q.entries[q.head] = record{}
	q.head = (q.head + 1) % len(q.entries)
	q.count--
	q.bytes -= r.size
	return r
}
func increment(n *uint64) {
	if *n < math.MaxUint64 {
		*n++
	}
}
func cloneEvent(e Event) Event { e.Fields = maps.Clone(e.Fields); return e }
