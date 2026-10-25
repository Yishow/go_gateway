package snapshot

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
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

// Mode selects how a group turns samples into records.
type Mode string

const (
	// ModeWide is the default: one row per entity and bucket, governed by the
	// incomplete-row policy and required members.
	ModeWide Mode = ""
	// ModeInstantaneousProfile writes one independent record per member and
	// window. A member without an eligible reading yields a record with a
	// reason instead of suppressing anyone else's, and no value is carried over.
	ModeInstantaneousProfile Mode = "instantaneous_profile"
	// ModeCumulativeProfile writes the same independent per-member record but
	// carries the latest eligible raw counter reading together with its exact
	// difference from the previous window's baseline (see Config.Counters).
	ModeCumulativeProfile Mode = "cumulative_profile"
	// ModeWideProfile writes one record per window with one value per member:
	// each member keeps its own role and reduction (a numeric mean or the newest
	// eligible reading), an optional member without a usable value is NULL with its
	// reason, and only a required member that is unusable skips the row. The window
	// is always one fixed UTC window, even when no member has a sample.
	ModeWideProfile Mode = "wide_profile"
	// ModeWideCumulative is ModeWideProfile for cumulative members: the same one
	// record per window, but every member keeps its latest eligible raw reading and
	// is closed against its own confirmed counter baseline (see Config.Counters), so
	// the typed difference rides beside the raw value. A mean is never a baseline.
	ModeWideCumulative Mode = "wide_cumulative"
)

// perMember reports whether every member is its own record in each window.
func (m Mode) perMember() bool { return m == ModeInstantaneousProfile || m == ModeCumulativeProfile }

// wide reports whether one record per window carries a value for every member.
func (m Mode) wide() bool { return m == ModeWideProfile || m == ModeWideCumulative }

// Cumulative reports whether members are closed against counter baselines.
func (m Mode) Cumulative() bool { return m == ModeCumulativeProfile || m == ModeWideCumulative }

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

// SampleReduction is how a profile member turns one window's observations into a
// value. The zero value keeps the existing rule: the newest eligible reading.
type SampleReduction string

const (
	// SampleReductionLatest keeps the newest eligible reading.
	SampleReductionLatest SampleReduction = ""
	// SampleReductionMean carries the arithmetic mean of the window's eligible
	// observations instead of one reading.
	SampleReductionMean SampleReduction = "mean"
)

// Member is one value contributing to a row.
type Member struct {
	// Key identifies the member; samples carry it as MemberKey.
	Key string
	// EntityKey partitions rows. Empty means the fixed group scope.
	EntityKey string
	Required  bool
	// SampleReduction is how this member's window value is reduced. Mean is
	// accepted only in ModeInstantaneousProfile and ModeWideProfile; the zero
	// value keeps the existing newest-eligible-reading rule.
	SampleReduction SampleReduction
	// SourceValueType is the frozen schema.DataType of the Tag this member reads.
	// It is supplied by the caller (the runtime mapping lands in 2.1); empty
	// means the frozen fact was not provided and is reported as unavailable.
	SourceValueType schema.DataType
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
	Until time.Time
	// StartsAt and EndsAt are the exact instants (half-open: observed at or
	// after StartsAt, before EndsAt) at which a forward cutover hands a source
	// to this revision or takes it away. They need not be interval-aligned: the
	// window that contains one is closed with its real bounds and marked
	// truncated. FirstBucket is the window containing StartsAt and Until the
	// window end at or after EndsAt. Zero means the ordinary aligned lifecycle.
	StartsAt, EndsAt time.Time
	IncompletePolicy IncompletePolicy
	Members          []Member
	// Mode selects the assembly rules; see ModeInstantaneousProfile.
	Mode              Mode
	MaxBucketsPerTick int
	// Counters is the confirmed counter binding of every member by member key.
	// It is required by the cumulative modes and refused otherwise.
	Counters map[string]aggregation.CounterBinding
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
	if !c.StartsAt.IsZero() {
		c.StartsAt = c.StartsAt.UTC()
		if !BucketStart(c.StartsAt, c.Interval).Equal(c.FirstBucket) {
			return invalid("starts at must fall inside the first bucket")
		}
	}
	if !c.EndsAt.IsZero() {
		c.EndsAt = c.EndsAt.UTC()
		end := BucketStart(c.EndsAt, c.Interval)
		if !end.Equal(c.EndsAt) {
			end = end.Add(c.Interval)
		}
		if !c.Until.Equal(end) {
			return invalid("until must be the window end at or after ends at")
		}
	}
	if c.Mode != ModeWide && !c.Mode.wide() && !c.Mode.perMember() {
		return invalid("unknown assembly mode")
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
		if member.SampleReduction != SampleReductionLatest && member.SampleReduction != SampleReductionMean {
			return invalid("unknown sample reduction")
		}
		if member.SampleReduction == SampleReductionMean && c.Mode != ModeInstantaneousProfile && c.Mode != ModeWideProfile {
			return invalid("mean sample reduction belongs to the instantaneous profile only")
		}
		if member.SourceValueType != "" {
			if _, ok := measurement.ExactTypeForTag(member.SourceValueType); !ok {
				return invalid("unknown frozen source value type")
			}
		}
		members[i] = member
	}
	c.Members = members
	if reason := c.countersInvalid(); reason != "" {
		return invalid(reason)
	}
	if c.Mode.wide() {
		// Optional members and required-skip are the only rules of this mode: it
		// always writes partial rows, so the storage must hold NULL and member
		// provenance, and its members are never entity-keyed.
		c.IncompletePolicy = IncompletePartial
		if !c.Storage.SupportsPartial() {
			return Config{}, fmt.Errorf("%w: %w", ErrInvalidConfig, ErrPartialUnsupported)
		}
		for _, member := range c.Members {
			if member.EntityKey != "" {
				return invalid("a wide profile member cannot be entity-keyed")
			}
		}
		return c, nil
	}
	if c.Mode.perMember() {
		// Every member is its own record; required and entity rules belong to wide rows.
		for _, member := range c.Members {
			if member.EntityKey != "" {
				return invalid("a profile member cannot be entity-keyed")
			}
		}
		return c, nil
	}
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
