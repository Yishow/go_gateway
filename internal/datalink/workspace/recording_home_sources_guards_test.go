package workspace

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestRecordingHomeSourcesReadGuards(t *testing.T) {
	valid := homeFixtures(t)[1].Before
	t.Run("nil reader", func(t *testing.T) {
		_, err := ReadRecordingHomeSources(context.Background(), "W1", nil)
		if !errors.Is(err, ErrRecordingHomeSourcesUnavailable) {
			t.Fatal(err)
		}
	})
	t.Run("nil context", func(t *testing.T) {
		//nolint:staticcheck // Deliberately verify the public nil-context rejection guard (SA1012).
		_, err := ReadRecordingHomeSources(nil, "W1", func(context.Context, string) (RecordingHomeSourceInput, error) {
			t.Fatal("must not read")
			return valid, nil
		})
		if !errors.Is(err, ErrRecordingHomeSourcesUnavailable) {
			t.Fatal(err)
		}
	})
	t.Run("empty scope", func(t *testing.T) {
		_, err := ReadRecordingHomeSources(context.Background(), "  ", nil)
		if !errors.Is(err, ErrRecordingHomeInvalidScope) {
			t.Fatal(err)
		}
	})
	t.Run("not found", func(t *testing.T) {
		_, err := ReadRecordingHomeSources(context.Background(), "W1", func(context.Context, string) (RecordingHomeSourceInput, error) {
			return valid, ErrRecordingHomeWorkspaceNotFound
		})
		if !errors.Is(err, ErrRecordingHomeWorkspaceNotFound) {
			t.Fatal(err)
		}
	})
	for _, cancelAt := range []int{0, 1, 2} {
		t.Run(string(rune('0'+cancelAt)), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			calls := 0
			if cancelAt == 0 {
				cancel()
			}
			got, err := ReadRecordingHomeSources(ctx, "W1", func(context.Context, string) (RecordingHomeSourceInput, error) {
				calls++
				if calls == cancelAt {
					cancel()
				}
				return valid, nil
			})
			if !errors.Is(err, context.Canceled) || got.RecordingHomeSourceSummary != nil || calls != cancelAt {
				t.Fatalf("got %+v err %v calls %d", got, err, calls)
			}
		})
	}
}

func TestRecordingHomeSourcesDetachAndIdentityGuards(t *testing.T) {
	for _, mutation := range []string{"revision", "duplicate-membership", "blank-device", "blank-rule", "blank-mapping", "blank-point"} {
		t.Run(mutation, func(t *testing.T) {
			input := homeFixtures(t)[1].Before
			calls := 0
			got, err := ReadRecordingHomeSources(context.Background(), "W1", func(context.Context, string) (RecordingHomeSourceInput, error) {
				calls++
				switch mutation {
				case "revision":
					if calls == 2 {
						input.Rules[0].RevisionID = "new-revision"
					}
				case "duplicate-membership":
					input.Workspace.OrderedDeviceIDs = []string{"D1", "D1"}
					input.Devices = append(input.Devices, RecordingHomeSourceDevice{ID: "D2"})
				case "blank-device":
					input.Devices[0].ID = ""
				case "blank-rule":
					input.Rules[0].ID = ""
				case "blank-mapping":
					input.Mappings[0].ID = ""
				case "blank-point":
					input.Points[0].ID = ""
				}
				return input, nil
			})
			if err != nil || got.State != "unknown" || got.RecordingHomeSourceSummary != nil || got.SourceDigest != "" {
				t.Fatalf("unsafe observation %+v err %v", got, err)
			}
			data, _ := json.Marshal(got)
			for _, forbidden := range []string{"candidates", "devices", "excluded", "source_digest"} {
				if strings.Contains(string(data), forbidden) {
					t.Fatalf("partial unknown %s", data)
				}
			}
		})
	}
}

func TestRecordingHomeSourcesDeterministicDigestAndReadOnly(t *testing.T) {
	input := homeFixtures(t)[29].Before
	before, _ := json.Marshal(input)
	read := func(v RecordingHomeSourceInput) RecordingHomeSourcesObservation {
		t.Helper()
		got, err := ReadRecordingHomeSources(context.Background(), "W1", func(context.Context, string) (RecordingHomeSourceInput, error) { return v, nil })
		if err != nil || got.State != "available" {
			t.Fatalf("%+v %v", got, err)
		}
		return got
	}
	a := read(input)
	b := read(input)
	if a.SourceDigest != b.SourceDigest {
		t.Fatal("unstable digest")
	}
	after, _ := json.Marshal(input)
	if !bytes.Equal(before, after) {
		t.Fatal("reader input mutated")
	}
	reordered := input
	reordered.Mappings = slices.Clone(input.Mappings)
	slices.Reverse(reordered.Mappings)
	reordered.Points = slices.Clone(input.Points)
	slices.Reverse(reordered.Points)
	c := read(reordered)
	if a.SourceDigest != c.SourceDigest {
		t.Fatal("source array order changed canonical digest")
	}
	if reflect.DeepEqual(a.Candidates, c.Candidates) {
		t.Fatal("mapping order not preserved")
	}
	input.Rules[0].RevisionID = "revision-2"
	d := read(input)
	if d.SourceDigest != a.SourceDigest {
		t.Fatal("non-presentable revision changed presentation digest")
	}
	if !reflect.DeepEqual(a.Candidates, d.Candidates) {
		t.Fatal("revision should not change candidates")
	}
	data, _ := json.Marshal(a)
	for _, forbidden := range []string{"revision_id", "skipped_addresses", "rule_id", "workspaceRevision", "password", "connection", "live_value"} {
		if strings.Contains(string(data), forbidden) {
			t.Fatalf("unsafe output %s", data)
		}
	}
}

func TestRecordingHomeSourcesInconsistentReader(t *testing.T) {
	for _, round := range []int{1, 2} {
		t.Run(string(rune('0'+round)), func(t *testing.T) {
			input := homeFixtures(t)[1].Before
			calls := 0
			got, err := ReadRecordingHomeSources(context.Background(), "W1", func(context.Context, string) (RecordingHomeSourceInput, error) {
				calls++
				if calls == round {
					return RecordingHomeSourceInput{}, ErrRecordingHomeSourcesInconsistent
				}
				return input, nil
			})
			if err != nil || got.State != "unknown" || got.ReasonCode != "source-inconsistent" || got.SourceDigest != "" || got.RecordingHomeSourceSummary != nil || calls != round {
				t.Fatalf("got %+v err %v calls %d", got, err, calls)
			}
		})
	}
}
