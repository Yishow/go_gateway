package grouppipeline

import "testing"

func TestEmptyHandoffComponentsRemainNil(t *testing.T) {
	if got := handoffComponents(nil); got != nil {
		t.Fatalf("empty components must remain nil: %#v", got)
	}
}
