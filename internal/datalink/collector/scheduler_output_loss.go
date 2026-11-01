package collector

import (
	"container/list"
	"sync"
)

// OutputLossSeries describes recently dropped results for one captured identity.
// Losses covers retained detail only; evicted detail remains in the aggregate.
type OutputLossSeries struct {
	InstallationScope string
	PointRevision     string
	PointID           string
	Losses            uint64
}

// OutputLossSnapshot is a detached, read-only view of actual output-result loss.
// Total equals Unattributed + EvictedDetailLoss + the sum of Recent.Losses.
// Recent is ordered from least to most recently lost identity, at most 256 entries.
// It is not a complete history and is separate from rejected due work.
type OutputLossSnapshot struct {
	Total             uint64
	Unattributed      uint64
	EvictedDetailLoss uint64
	Recent            []OutputLossSeries
}

const outputLossRecentLimit = 256

type outputLossIdentity struct {
	scope, revision, point string
}

// outputLossCounter is zero-value usable and independent of scheduler locks.
// Only actual emitValue discards feed it; it does not observe reads or due work.
type outputLossCounter struct {
	mu                sync.Mutex
	total             uint64
	unattributed      uint64
	evictedDetailLoss uint64
	byIdentity        map[outputLossIdentity]*list.Element
	recent            list.List
}

func (c *outputLossCounter) record(value CollectedValue) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.total++
	identity := outputLossIdentity{value.InstallationScope, value.PointRevision, value.PointID}
	if identity.scope == "" || identity.revision == "" || identity.point == "" {
		c.unattributed++
		return
	}
	if element := c.byIdentity[identity]; element != nil {
		outputLossEntry(element).Losses++
		c.recent.MoveToBack(element)
		return
	}
	if c.byIdentity == nil {
		c.byIdentity = make(map[outputLossIdentity]*list.Element)
	}
	if c.recent.Len() == outputLossRecentLimit {
		oldest := c.recent.Front()
		series := outputLossEntry(oldest)
		c.evictedDetailLoss += series.Losses
		delete(c.byIdentity, outputLossIdentity{series.InstallationScope, series.PointRevision, series.PointID})
		c.recent.Remove(oldest)
	}
	series := &OutputLossSeries{
		InstallationScope: identity.scope,
		PointRevision:     identity.revision,
		PointID:           identity.point,
		Losses:            1,
	}
	c.byIdentity[identity] = c.recent.PushBack(series)
}

func (c *outputLossCounter) snapshot() OutputLossSnapshot {
	c.mu.Lock()
	defer c.mu.Unlock()
	result := OutputLossSnapshot{
		Total:             c.total,
		Unattributed:      c.unattributed,
		EvictedDetailLoss: c.evictedDetailLoss,
		Recent:            make([]OutputLossSeries, 0, c.recent.Len()),
	}
	for element := c.recent.Front(); element != nil; element = element.Next() {
		result.Recent = append(result.Recent, *outputLossEntry(element))
	}
	return result
}

// Every list entry is inserted above as this private series type. An invariant
// violation must fail explicitly, just as the original direct assertion did.
func outputLossEntry(element *list.Element) *OutputLossSeries {
	series, ok := element.Value.(*OutputLossSeries)
	if !ok {
		panic("invalid internal output loss series")
	}
	return series
}
