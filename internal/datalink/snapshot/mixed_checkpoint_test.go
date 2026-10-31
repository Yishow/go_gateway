package snapshot

import (
	"encoding/json"
	"errors"
	"maps"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func checkpointAssembler(t *testing.T, cfg Config) *Assembler {
	t.Helper()
	a, err := NewAssembler(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return a
}
func checkpointOf(t *testing.T, a *Assembler) MixedCheckpoint {
	t.Helper()
	cp, err := a.MixedCheckpoint()
	if err != nil {
		t.Fatal(err)
	}
	return cp
}

func TestMixedCheckpointRoundTrip(t *testing.T) {
	for _, mode := range []Mode{ModeMixedWide, ModeMixedLong} {
		t.Run(string(mode), func(t *testing.T) {
			cfg := mixedConfig(mode)
			cfg.StartsAt = mixedTime(7)
			a := checkpointAssembler(t, cfg)
			initial := checkpointOf(t, a)
			if !initial.Progress.ClosedThrough.IsZero() || !initial.Members["a"].Start.Equal(mixedTime(0)) {
				t.Fatal(initial)
			}
			restored := checkpointAssembler(t, cfg)
			if err := restored.RestoreMixedCheckpoint(initial); err != nil {
				t.Fatal(err)
			}
			mixedOffer(t, a, mixedSample("a", "a1", 10, 22))
			mixedOffer(t, a, mixedSample("b", "b1", 10, 12))
			p := a.Plan(mixedTime(15))
			cp, err := a.MixedCheckpointForPlan(p)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(initial, checkpointOf(t, a)) {
				t.Fatal("plan projection mutated cursor")
			}
			if !cp.Members["a"].Start.Equal(mixedTime(15)) || !cp.Members["b"].Start.Equal(mixedTime(0)) {
				t.Fatal(cp)
			}
			a.Commit(p)
			if !reflect.DeepEqual(cp, checkpointOf(t, a)) {
				t.Fatal("plan and commit differ")
			}
			raw, err := json.Marshal(cp)
			if err != nil {
				t.Fatal(err)
			}
			var decoded MixedCheckpoint
			if err = json.Unmarshal(raw, &decoded); err != nil {
				t.Fatal(err)
			}
			if err = restored.RestoreMixedCheckpoint(decoded); err != nil {
				t.Fatal(err)
			}
			// Accepted slow-member samples are persisted/replayed separately, not lost
			// or silently fabricated by this progress-only projection.
			mixedOffer(t, restored, mixedSample("b", "b1", 10, 12))
			if !reflect.DeepEqual(a.Plan(mixedTime(60)).Outcomes, restored.Plan(mixedTime(60)).Outcomes) {
				t.Fatal("restart changed outputs")
			}
			if _, err = a.MixedCheckpointForPlan(p); err == nil {
				t.Fatal("accepted stale plan")
			}
			delete(cp.Members, "b")
			if len(checkpointOf(t, a).Members) != 2 {
				t.Fatal("projection aliases live state")
			}
		})
	}
}

func TestMixedCheckpointCutoffNominalAfterActual(t *testing.T) {
	cfg := mixedConfig(ModeMixedWide)
	cfg.EndsAt = mixedTime(120)
	a := checkpointAssembler(t, cfg)
	mixedOffer(t, a, mixedSample("b", "before", 10, 12))
	if err := a.SetMixedEndsAt(mixedTime(15)); err != nil {
		t.Fatal(err)
	}
	p := a.Plan(mixedTime(60))
	cp, err := a.MixedCheckpointForPlan(p)
	if err != nil {
		t.Fatal(err)
	}
	if !cp.Progress.ClosedThrough.Equal(mixedTime(60)) || !cp.EndsAt.Equal(mixedTime(15)) {
		t.Fatal(cp)
	}
	for _, m := range cp.Members {
		if !m.End.IsZero() {
			t.Fatal("ended member has next due")
		}
	}
	a.Commit(p)
	b := checkpointAssembler(t, cfg)
	if err = b.RestoreMixedCheckpoint(cp); err != nil {
		t.Fatal(err)
	}
	if !b.Done() || !b.Plan(mixedTime(120)).Empty() {
		t.Fatal("restart extended cutoff")
	}
}

func TestMixedCheckpointRejectsCorruptionAtomically(t *testing.T) {
	cfg := mixedConfig(ModeMixedWide)
	cfg.EndsAt = mixedTime(120)
	source := checkpointAssembler(t, cfg)
	cp, err := source.MixedCheckpointForPlan(source.Plan(mixedTime(15)))
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string]func(*MixedCheckpoint){
		"version":           func(c *MixedCheckpoint) { c.Version = 2 },
		"workspace":         func(c *MixedCheckpoint) { c.Progress.WorkspaceID = "foreign" },
		"group":             func(c *MixedCheckpoint) { c.Progress.GroupID = "foreign" },
		"revision":          func(c *MixedCheckpoint) { c.Progress.GroupRevision = "foreign" },
		"starts":            func(c *MixedCheckpoint) { c.StartsAt = mixedTime(1) },
		"end_extension":     func(c *MixedCheckpoint) { c.EndsAt = mixedTime(121) },
		"end_absent":        func(c *MixedCheckpoint) { c.EndsAt = time.Time{} },
		"end_before_start":  func(c *MixedCheckpoint) { c.EndsAt = mixedTime(-1) },
		"missing":           func(c *MixedCheckpoint) { delete(c.Members, "a") },
		"extra":             func(c *MixedCheckpoint) { c.Members["other"] = c.Members["a"] },
		"substitution":      func(c *MixedCheckpoint) { c.Members["other"] = c.Members["a"]; delete(c.Members, "a") },
		"phase":             func(c *MixedCheckpoint) { m := c.Members["a"]; m.Start = m.Start.Add(time.Second); c.Members["a"] = m },
		"cadence":           func(c *MixedCheckpoint) { m := c.Members["a"]; m.End = m.End.Add(time.Second); c.Members["a"] = m },
		"zero_member":       func(c *MixedCheckpoint) { c.Members["b"] = MemberProgress{} },
		"off_boundary":      func(c *MixedCheckpoint) { c.Progress.ClosedThrough = mixedTime(16) },
		"before_activation": func(c *MixedCheckpoint) { c.Progress.ClosedThrough = mixedTime(-15) },
		"overflow":          func(c *MixedCheckpoint) { c.Progress.ClosedThrough = time.Date(2300, 1, 1, 0, 0, 0, 0, time.UTC) },
	}
	for name, corrupt := range cases {
		t.Run(name, func(t *testing.T) {
			a := checkpointAssembler(t, cfg)
			before := checkpointOf(t, a)
			bad := cp
			bad.Members = maps.Clone(cp.Members)
			corrupt(&bad)
			if err := a.RestoreMixedCheckpoint(bad); err == nil {
				t.Fatal("accepted corrupt checkpoint")
			}
			if !reflect.DeepEqual(before, checkpointOf(t, a)) {
				t.Fatal("failed restore mutated live state")
			}
		})
	}
	t.Run("accepted_samples", func(t *testing.T) {
		a := checkpointAssembler(t, cfg)
		mixedOffer(t, a, mixedSample("b", "pending", 10, 1))
		before := a.Plan(mixedTime(60))
		if err := a.RestoreMixedCheckpoint(cp); err == nil {
			t.Fatal("restored over accepted samples")
		}
		after := a.Plan(mixedTime(60))
		require.NoError(t, before.Err)
		require.NoError(t, after.Err)
		require.Equal(t, before, after, "failed restore changed samples")
	})
	t.Run("committed", func(t *testing.T) {
		a := checkpointAssembler(t, cfg)
		a.Tick(mixedTime(15))
		before := checkpointOf(t, a)
		if err := a.RestoreMixedCheckpoint(cp); err == nil {
			t.Fatal("restored over committed state")
		}
		if !reflect.DeepEqual(before, checkpointOf(t, a)) {
			t.Fatal("failed restore changed progress")
		}
	})
}

func TestMixedCheckpointRejectsInvalidPlan(t *testing.T) {
	a := checkpointAssembler(t, mixedConfig(ModeMixedWide))
	for _, name := range []string{"foreign", "missing", "bounds", "error", "empty"} {
		t.Run(name, func(t *testing.T) {
			p := a.Plan(mixedTime(15))
			switch name {
			case "foreign":
				p.mixedProgress.GroupID = "other"
			case "missing":
				delete(p.mixedNext, "b")
			case "bounds":
				p.mixedNext["a"] = mixedTime(16)
			case "error":
				p.Err = errors.New("window failure")
			case "empty":
				p = Plan{}
			}
			if _, err := a.MixedCheckpointForPlan(p); err == nil {
				t.Fatal("accepted invalid plan")
			}
		})
	}
	p := a.Plan(mixedTime(1))
	cp, err := a.MixedCheckpointForPlan(p)
	if err != nil || !reflect.DeepEqual(cp, checkpointOf(t, a)) {
		t.Fatalf("empty valid plan: %+v %v", cp, err)
	}
}

func TestMixedCheckpointInitialAndBackwardCorruption(t *testing.T) {
	cfg := mixedConfig(ModeMixedWide)
	a := checkpointAssembler(t, cfg)
	cp := checkpointOf(t, a)
	cp.Members["a"] = MemberProgress{Start: mixedTime(15), End: mixedTime(30)}
	if err := a.RestoreMixedCheckpoint(cp); err == nil {
		t.Fatal("zero watermark accepted advanced member")
	}
	a.Tick(mixedTime(60))
	p := a.Plan(mixedTime(75))
	p.mixedProgress.ClosedThrough = mixedTime(15)
	p.mixedNext = map[string]time.Time{"a": mixedTime(15), "b": mixedTime(0)}
	if _, err := a.MixedCheckpointForPlan(p); err == nil {
		t.Fatal("plan checkpoint regressed live watermark")
	}
}
