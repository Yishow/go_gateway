package snapshot

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// IncompletePolicy decides what a bucket does when a member is not usable.
type IncompletePolicy string

const (
	// IncompleteSkipRow is the default: any unusable member skips the row.
	IncompleteSkipRow IncompletePolicy = "skip_row"
	// IncompletePartial writes SQL NULL plus a reason for unusable members and
	// is accepted only when the storage can hold NULL and member metadata.
	IncompletePartial IncompletePolicy = "partial"
)

// DefaultMaxBucketsPerTick bounds catch-up work after a long pause; the
// remaining buckets close on the following ticks, in order.
const DefaultMaxBucketsPerTick = 1000

// StorageCapabilities are facts about the destination, verified by the caller
// from real table metadata. They gate opt-in modes only: a complete good row
// needs none of them because its member provenance stays in the local envelope.
type StorageCapabilities struct {
	// NullableValues: every optional member column accepts SQL NULL.
	NullableValues bool
	// MemberQuality: the destination stores per-member quality/provenance.
	MemberQuality bool
}

// SupportsPartial reports whether partial rows can be stored honestly.
func (c StorageCapabilities) SupportsPartial() bool { return c.NullableValues && c.MemberQuality }

// ErrPartialUnsupported is the specific invalid-config cause for a partial
// policy whose storage cannot hold NULL values and member quality.
var ErrPartialUnsupported = errors.New("partial policy needs nullable storage and member quality")

// ErrInvalidConfig marks a group configuration the assembler cannot honor.
var ErrInvalidConfig = errors.New("snapshot config invalid")

// Member is one value contributing to a row.
type Member struct {
	// Key identifies the member; samples carry it as MemberKey.
	Key string
	// EntityKey partitions rows. Empty means the fixed group scope.
	EntityKey string
	Required  bool
	// MaxAge is the freshness limit measured from bucket end. Zero means the
	// bucket interval, so nothing from another bucket is ever carried over.
	MaxAge time.Duration
	// SourceRevision and MappingRevision, when set, must match each sample.
	// A changed scale or mapping is a new revision, so old pending values are
	// rejected rather than reinterpreted.
	SourceRevision  string
	MappingRevision string
}

// Config fixes one applied group revision.
type Config struct {
	// WorkspaceID, GroupID, GroupRevision and DestinationScope freeze the row
	// identity; DestinationScope comes from DestinationScope().
	WorkspaceID      string
	GroupID          string
	GroupRevision    string
	DestinationScope string

	Interval        time.Duration
	AllowedLateness time.Duration
	// MaxFutureSkew is the explicit tolerance for observed_at ahead of the
	// gateway clock. Zero tolerates none.
	MaxFutureSkew time.Duration
	// FirstBucket is the first bucket start that must close (the group's
	// effective boundary); it must be aligned to Interval in UTC.
	FirstBucket time.Time
	// Until, when set, is the first bucket start this assembler no longer
	// closes: a superseded revision stops here. It must be interval-aligned.
	Until             time.Time
	IncompletePolicy  IncompletePolicy
	Members           []Member
	MaxBucketsPerTick int
	// Storage describes verified destination capabilities; it is consulted
	// only for IncompletePartial.
	Storage StorageCapabilities
}

func (c Config) normalized() (Config, error) {
	invalid := func(reason string) (Config, error) {
		return Config{}, fmt.Errorf("%w: %s", ErrInvalidConfig, reason)
	}
	for _, field := range []struct{ name, value string }{
		{"workspace", c.WorkspaceID}, {"group", c.GroupID},
		{"group revision", c.GroupRevision}, {"destination scope", c.DestinationScope},
	} {
		if strings.TrimSpace(field.value) == "" {
			return invalid(field.name + " is required for row identity")
		}
	}
	if c.Interval < time.Second || c.Interval%time.Second != 0 {
		return invalid("interval must be a positive whole number of seconds")
	}
	if c.AllowedLateness < 0 || c.MaxFutureSkew < 0 {
		return invalid("lateness and future skew must not be negative")
	}
	if c.FirstBucket.IsZero() || !BucketStart(c.FirstBucket, c.Interval).Equal(c.FirstBucket) {
		return invalid("first bucket must be an interval-aligned UTC start")
	}
	c.FirstBucket = c.FirstBucket.UTC()
	if !c.Until.IsZero() {
		if !BucketStart(c.Until, c.Interval).Equal(c.Until) || c.Until.Before(c.FirstBucket) {
			return invalid("until must be an interval-aligned UTC start not before the first bucket")
		}
		c.Until = c.Until.UTC()
	}
	if c.IncompletePolicy == "" {
		c.IncompletePolicy = IncompleteSkipRow
	}
	if c.MaxBucketsPerTick <= 0 {
		c.MaxBucketsPerTick = DefaultMaxBucketsPerTick
	}
	if len(c.Members) == 0 {
		return invalid("at least one member is required")
	}
	members := make([]Member, len(c.Members))
	seen := make(map[string]struct{}, len(c.Members))
	for i, member := range c.Members {
		if member.Key == "" {
			return invalid("member key is required")
		}
		if _, dup := seen[member.Key]; dup {
			return invalid("duplicate member key")
		}
		seen[member.Key] = struct{}{}
		if member.MaxAge < 0 {
			return invalid("member max age must not be negative")
		}
		if member.MaxAge == 0 {
			member.MaxAge = c.Interval
		}
		members[i] = member
	}
	c.Members = members
	switch c.IncompletePolicy {
	case IncompleteSkipRow:
		for _, member := range c.Members {
			if !member.Required {
				return invalid("skip_row requires every member to be required")
			}
		}
	case IncompletePartial:
		if !c.Storage.SupportsPartial() {
			return Config{}, fmt.Errorf("%w: %w", ErrInvalidConfig, ErrPartialUnsupported)
		}
	default:
		return invalid("unknown incomplete policy")
	}
	return c, nil
}
