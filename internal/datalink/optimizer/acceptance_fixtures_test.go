package optimizer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

// This is a fixture integrity check, NOT a test of an implemented planner.
// Native offsets are explicit test inputs, not an alternative address parser.
type acceptancePoint struct {
	ID            string `json:"id"`
	Address       string `json:"address"`
	Start         int    `json:"start"`
	Width         int    `json:"width"`
	Type          string `json:"type"`
	Format        string `json:"format"`
	ExpectedValue uint32 `json:"expected_value"`
	Repeat        int    `json:"repeat"`
	Stride        int    `json:"stride"`
}
type acceptanceRequest struct {
	Kind   string `json:"kind"`
	Start  int    `json:"start"`
	Length int    `json:"length"`
	Units  []int  `json:"units"`
}
type acceptanceCase struct {
	ID               string              `json:"id"`
	Protocol         string              `json:"protocol"`
	Limit            int                 `json:"limit"`
	Points           []acceptancePoint   `json:"points"`
	Requests         []acceptanceRequest `json:"requests"`
	Readable         [][2]int            `json:"readable"`
	Forbidden        [][2]int            `json:"forbidden"`
	CoverageRevision string              `json:"coverage_revision"`
	Errors           []string            `json:"errors"`
	Reason           string              `json:"reason"`
	RawWords         []uint16            `json:"raw_words"`
	RecordSelection  []string            `json:"record_selection"`
	Mode             string              `json:"mode"`
	Requested        []string            `json:"requested"`
	RandomSupported  bool                `json:"random_supported"`
	NativeBits       int                 `json:"native_bits"`
}

func TestAcquisitionAcceptanceFixtureIntegrity(t *testing.T) {
	data, err := os.ReadFile("testdata/acquisition_acceptance.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		Version   int              `json:"version"`
		Status    string           `json:"status"`
		Intervals string           `json:"intervals"`
		Cases     []acceptanceCase `json:"cases"`
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&fixture); err != nil {
		t.Fatal(err)
	}
	if fixture.Version != 1 || fixture.Status != "design-expectations-not-executed-plans" || fixture.Intervals != "half-open-native-units" {
		t.Fatal("invalid fixture header")
	}
	if len(fixture.Cases) != 22 {
		t.Fatalf("expected 22 vectors, got %d", len(fixture.Cases))
	}
	seen := map[string]bool{}
	for _, c := range fixture.Cases {
		if seen[c.ID] || c.ID == "" {
			t.Fatalf("duplicate/empty case %q", c.ID)
		}
		seen[c.ID] = true
		t.Run(c.ID, func(t *testing.T) { validateAcceptanceCase(t, c) })
	}
}

func validateAcceptanceCase(t *testing.T, c acceptanceCase) {
	t.Helper()
	if c.Limit <= 0 || (c.Protocol != "mc" && c.Protocol != "fatek" && c.Protocol != "modbus") {
		t.Fatal("invalid capability")
	}
	if len(c.Readable) > 0 && c.CoverageRevision == "" {
		t.Fatal("readable coverage lacks revision")
	}
	points := map[string]acceptancePoint{}
	active := map[int]bool{}
	errors := map[string]bool{}
	for _, id := range c.Errors {
		if errors[id] {
			t.Fatal("duplicate error")
		}
		errors[id] = true
	}
	for _, template := range c.Points {
		validateAcceptanceNativePoint(t, c, template)
		count := template.Repeat
		if count == 0 {
			count = 1
		}
		if count < 1 || template.Width < 1 || template.Start < 0 || (count > 1 && template.Stride < 1) {
			t.Fatal("invalid point span")
		}
		for i := 0; i < count; i++ {
			p := template
			if template.Repeat > 0 {
				p.ID = fmt.Sprintf("%s%d", template.ID, i)
				p.Start += i * template.Stride
			}
			if _, ok := points[p.ID]; ok || p.ID == "" {
				t.Fatal("duplicate/empty point identity")
			}
			points[p.ID] = p
			if !errors[p.ID] {
				for u := p.Start; u < p.Start+p.Width; u++ {
					active[u] = true
				}
			}
		}
	}
	for id := range errors {
		if _, ok := points[id]; !ok || c.Reason == "" {
			t.Fatal("error lacks point/reason")
		}
	}
	requestUnits := make([]map[int]bool, len(c.Requests))
	for i, r := range c.Requests {
		units := map[int]bool{}
		switch r.Kind {
		case "contiguous":
			if r.Length < 1 || r.Length > c.Limit || r.Start < 0 {
				t.Fatal("invalid contiguous limit")
			}
			for u := r.Start; u < r.Start+r.Length; u++ {
				units[u] = true
			}
		case "random":
			if !c.RandomSupported || c.Protocol == "modbus" || len(r.Units) < 1 || len(r.Units) > c.Limit {
				t.Fatal("invalid random capability/limit")
			}
			for _, u := range r.Units {
				if units[u] || u < 0 {
					t.Fatal("duplicate/invalid random unit")
				}
				units[u] = true
			}
		default:
			t.Fatal("unknown request kind")
		}
		for u := range units {
			if inAcceptanceRanges(u, c.Forbidden) || (!active[u] && !inAcceptanceRanges(u, c.Readable)) {
				t.Fatalf("unauthorized unit %d", u)
			}
		}
		requestUnits[i] = units
	}
	for id, p := range points {
		if errors[id] {
			continue
		}
		complete := 0
		for _, units := range requestUnits {
			covered := true
			for u := p.Start; u < p.Start+p.Width; u++ {
				covered = covered && units[u]
			}
			if covered {
				complete++
			}
		}
		if complete != 1 {
			t.Fatalf("point %s must fit one complete request, got %d", id, complete)
		}
		if p.Format != "" {
			if len(c.RawWords) != 2 || p.Type != "uint32" {
				t.Fatal("invalid decoder fixture")
			}
			a, b := uint32(c.RawWords[0]), uint32(c.RawWords[1])
			if p.Format == "CDAB" {
				a, b = b, a
			} else if p.Format != "ABCD" {
				t.Fatal("unsupported fixture format")
			}
			if got := a<<16 | b; got != p.ExpectedValue {
				t.Fatalf("decoder oracle mismatch %x", got)
			}
		}
	}
	for _, id := range append(append([]string{}, c.RecordSelection...), c.Requested...) {
		if _, ok := points[id]; !ok {
			t.Fatalf("unknown selection %s", id)
		}
	}
	if c.Mode == "manual" && len(c.Requested) != len(points) {
		t.Fatal("manual input must be exact requested subset")
	}
	t.Logf("design only: acquisition_points=%d expected_initial_requests=%d explicit_errors=%d", len(points), len(c.Requests), len(errors))
}

func inAcceptanceRanges(unit int, intervals [][2]int) bool {
	for _, interval := range intervals {
		if unit >= interval[0] && unit < interval[1] {
			return true
		}
	}
	return false
}
