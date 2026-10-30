package snapshot

import (
	"fmt"
	"go-gateway/internal/datalink/measurement"
	"maps"
	"strings"
)

func (c Config) normalizedMixed() (Config, error) {
	invalid := func(reason string) (Config, error) { return Config{}, fmt.Errorf("%w: %s", ErrInvalidConfig, reason) }
	for _, s := range []string{c.WorkspaceID, c.GroupID, c.GroupRevision, c.DestinationScope} {
		if strings.TrimSpace(s) == "" {
			return invalid("mixed scope identity is required")
		}
	}
	if c.Interval != 0 || !c.FirstBucket.IsZero() || !c.Until.IsZero() {
		return invalid("mixed has member windows, not a legacy common window")
	}
	if c.StartsAt.IsZero() {
		return invalid("mixed starts at is required")
	}
	if c.AllowedLateness < 0 || c.MaxFutureSkew < 0 {
		return invalid("negative lateness or future skew")
	}
	if _, err := CheckedWindowAdd(c.StartsAt, c.MaxFutureSkew, "future skew"); err != nil {
		return Config{}, err
	}
	if !c.EndsAt.IsZero() {
		if _, err := CheckedWindowAdd(c.EndsAt, 0, "ends at"); err != nil {
			return Config{}, err
		}
		if !c.EndsAt.After(c.StartsAt) {
			return invalid("ends at must follow starts at")
		}
	}
	if len(c.Members) == 0 {
		return invalid("at least one member is required")
	}
	if c.MaxBucketsPerTick <= 0 {
		c.MaxBucketsPerTick = DefaultMaxBucketsPerTick
	}
	c.Members = append([]Member(nil), c.Members...)
	c.Counters = maps.Clone(c.Counters)
	for key, binding := range c.Counters {
		c.Counters[key] = cloneCounterBinding(binding)
	}
	seen := map[string]bool{}
	deltaCount := 0
	for i, m := range c.Members {
		if m.Key == "" || seen[m.Key] || m.EntityKey != "" {
			return invalid("mixed member identity must be unique and not entity-keyed")
		}
		seen[m.Key] = true
		if m.SampleReduction != SampleReductionLatest {
			return invalid("mixed calculation cannot carry legacy sample reduction")
		}
		start, end, err := CheckedWindowBounds(c.StartsAt, m.Interval)
		if err != nil {
			return Config{}, err
		}
		_ = start
		if _, err = CheckedWindowAdd(end, c.AllowedLateness, "lateness"); err != nil {
			return Config{}, err
		}
		if m.MaxAge < 0 {
			return invalid("negative member max age")
		}
		if m.MaxAge == 0 {
			m.MaxAge = m.Interval
		}
		valueType, known := measurement.ExactTypeForTag(m.SourceValueType)
		if !known {
			return invalid("mixed requires a frozen source value type")
		}
		binding, bound := c.Counters[m.Key]
		switch m.Calculation {
		case CalculationLatest:
			if bound {
				return invalid("latest cannot have a calculation baseline")
			}
		case CalculationMean:
			if !numericExactType(valueType) || bound {
				return invalid("mean needs numeric gauge without a calculation baseline")
			}
		case CalculationDelta:
			deltaCount++
			if !numericExactType(valueType) || !bound || binding.Validate() != nil {
				return invalid("delta needs a confirmed numeric counter binding")
			}
			if binding.MemberKey != m.Key || binding.WorkspaceID != c.WorkspaceID || binding.GroupID != c.GroupID || binding.ProfileRevision != c.GroupRevision {
				return invalid("foreign counter binding")
			}
		default:
			return invalid("explicit mixed calculation is required")
		}
		c.Members[i] = m
	}
	if len(c.Counters) != deltaCount {
		return invalid("counter binding for non-delta member")
	}
	// Wide not-due values require nullable columns even for required members.
	needsPartial := c.Mode == ModeMixedWide
	for _, m := range c.Members {
		needsPartial = needsPartial || !m.Required
	}
	if needsPartial && !c.Storage.SupportsPartial() {
		return Config{}, fmt.Errorf("%w: %w", ErrInvalidConfig, ErrPartialUnsupported)
	}
	if c.IncompletePolicy != "" && c.IncompletePolicy != IncompletePartial {
		return invalid("mixed uses per-member required semantics")
	}
	c.IncompletePolicy = IncompletePartial
	c.StartsAt = c.StartsAt.UTC()
	c.EndsAt = c.EndsAt.UTC()
	return c, nil
}
