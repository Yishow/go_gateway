package workspace

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

type homeParityCase struct {
	Name  string `json:"name"`
	Scope struct {
		WorkspaceID string `json:"workspaceId"`
	} `json:"scope"`
	Before    RecordingHomeSourceInput  `json:"before"`
	After     *RecordingHomeSourceInput `json:"after"`
	ReadError *struct {
		Round int `json:"round"`
	} `json:"read_error"`
	Expected struct {
		State string                      `json:"state"`
		Value *RecordingHomeSourceSummary `json:"value"`
	} `json:"expected"`
}

func homeFixtures(t *testing.T) []homeParityCase {
	t.Helper()
	b, err := os.ReadFile("../../../testdata/recording_home/source_projection.json")
	if err != nil {
		t.Fatal(err)
	}
	var f struct {
		Cases []homeParityCase `json:"cases"`
	}
	if err = json.Unmarshal(b, &f); err != nil {
		t.Fatal(err)
	}
	if len(f.Cases) != 49 {
		t.Fatalf("expected 49 independently authored shared cases, got %d", len(f.Cases))
	}
	return f.Cases
}
func TestRecordingHomeSourcesSharedParity(t *testing.T) {
	for _, f := range homeFixtures(t) {
		t.Run(f.Name, func(t *testing.T) {
			calls := 0
			got, err := ReadRecordingHomeSources(context.Background(), f.Scope.WorkspaceID, func(context.Context, string) (RecordingHomeSourceInput, error) {
				calls++
				if f.ReadError != nil && calls == f.ReadError.Round {
					return RecordingHomeSourceInput{}, errors.New("password=secret private.host")
				}
				if calls == 2 && f.After != nil {
					return *f.After, nil
				}
				return f.Before, nil
			})
			switch {
			case f.Expected.State == "error":
				if !errors.Is(err, ErrRecordingHomeSourcesUnavailable) {
					t.Fatalf("expected safe unavailable, got %v", err)
				}
			case err != nil:
				t.Fatal(err)
			case got.State != f.Expected.State:
				t.Fatalf("state %s want %s", got.State, f.Expected.State)
			}
			if !reflect.DeepEqual(got.RecordingHomeSourceSummary, f.Expected.Value) {
				t.Fatalf("projection %#v want %#v", got.RecordingHomeSourceSummary, f.Expected.Value)
			}
			expectedCalls := 2
			if f.ReadError != nil {
				expectedCalls = f.ReadError.Round
			}
			if calls != expectedCalls {
				t.Fatalf("reads %d want %d", calls, expectedCalls)
			}
			b, _ := json.Marshal(got)
			if strings.Contains(string(b), "password") || strings.Contains(string(b), "private.host") {
				t.Fatalf("unsafe response %s", b)
			}
			if got.State == "available" {
				if len(got.SourceDigest) != 64 {
					t.Fatal("missing SHA256 digest")
				}
			} else if got.SourceDigest != "" || got.RecordingHomeSourceSummary != nil {
				t.Fatal("partial unknown")
			}
		})
	}
}
