package optimizer

import (
	"reflect"
	"testing"
)

func TestAcquisitionRequestSnapshotDetachedAndBounded(t *testing.T) {
	plan := BuildAcquisitionPlan(randomInput("mc", 0, 10, 20))
	before := plan.Requests()
	fingerprint := plan.Fingerprint()
	if len(before) != 1 || len(before[0].Units) != 3 {
		t.Fatal("fixture must have random raw units")
	}
	for _, index := range []int{-1, 1, 1000000} {
		got, ok := plan.Request(index)
		if ok || !reflect.DeepEqual(got, AcquisitionRequest{}) {
			t.Fatalf("invalid index %d yielded %+v", index, got)
		}
	}
	got, ok := plan.Request(0)
	if !ok || !reflect.DeepEqual(before[0], got) {
		t.Fatal("request accessor diverges")
	}
	got.Group.Identity.DeviceID = "foreign"
	got.Points[0].Point.Address = "D99"
	got.Points[0].Offset = 99
	got.Units[0].Address = 99
	after, _ := plan.Request(0)
	if !reflect.DeepEqual(before[0], after) || !reflect.DeepEqual(before, plan.Requests()) || plan.Fingerprint() != fingerprint {
		t.Fatal("request mutation leaked into immutable plan")
	}
	var empty AcquisitionPlan
	if _, ok := empty.Request(0); ok {
		t.Fatal("empty plan has a request")
	}
}
