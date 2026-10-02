package groupdelivery

import (
	"errors"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/snapshot"
)

// GroupKey scopes every durable record to one applied group revision.
type GroupKey struct {
	WorkspaceID   string
	GroupID       string
	GroupRevision string
}

// Valid reports whether every part of the key is present.
func (k GroupKey) Valid() bool {
	return k.WorkspaceID != "" && k.GroupID != "" && k.GroupRevision != ""
}

// Restored is what a restarted boundary needs to rebuild open buckets.
type Restored struct {
	// NextClose is the checkpoint; zero when the group never closed a bucket.
	NextClose time.Time
	// Samples are accepted but not yet consumed by a closed bucket, in
	// acceptance order.
	Samples []snapshot.Sample
}

// ErrSampleConflict marks a resend of an accepted sample ID with different content.
var ErrSampleConflict = errors.New("durable sample identity conflict")

// ErrInvalidGroupKey marks a missing workspace, group or revision.
var ErrInvalidGroupKey = errors.New("durable group key invalid")

// ErrInvalidSample marks a sample without identity, member or observation time.
var ErrInvalidSample = errors.New("durable sample invalid")

// FrozenDestination is the destination and revisions an accepted row keeps for
// its whole life; later edits to the group never retarget accepted rows.
type FrozenDestination struct {
	Scope             string
	ConnectorID       string
	ConnectorRevision string
	Database          string
	TableSchema       string
	TableName         string
	// DedupeCapability states how the destination can recognize a repeated
	// effect: "receipt", "unique_key" or "none" (limited, unknown stops).
	DedupeCapability string
	// RecordKeyColumn is the destination column that holds the record ID; the
	// unique_key capability depends on it being unique.
	RecordKeyColumn string
}

// ClosedBucket is one closed entity bucket: its outcome and, for rows, the
// encoded SQL values to deliver.
type ClosedBucket struct {
	Outcome snapshot.Outcome
	Row     *dbtarget.EncodedRow
}

// Closure is everything one Tick closes, committed atomically with the
// checkpoint that moves past it.
type Closure struct {
	Destination FrozenDestination
	NextClose   time.Time
	Buckets     []ClosedBucket
}

// ErrEffectConflict marks an effect key that already exists with different content.
var ErrEffectConflict = errors.New("durable effect identity conflict")

// ErrInvalidClosure marks a closure that cannot be persisted faithfully.
var ErrInvalidClosure = errors.New("durable closure invalid")
