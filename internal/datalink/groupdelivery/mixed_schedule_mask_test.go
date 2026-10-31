package groupdelivery

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/snapshot"
)

func TestMixedClosureScheduleRejectsForgedMemberMask(t *testing.T) {
	for _, mutation := range []string{"due-as-not-due", "not-due-as-due", "not-due-bounds", "missing-member", "duplicate-member"} {
		t.Run(mutation, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			_, closure := mixedStorePlan(t, store)
			members := closure.Buckets[0].Outcome.Members
			fast, slow := 0, 1
			if members[fast].MemberKey != `["d","p","fast"]` {
				fast, slow = slow, fast
			}
			switch mutation {
			case "due-as-not-due":
				members[fast].Status = snapshot.MemberNotDue
				members[fast].Sample = nil
				members[fast].Window.ActualStart, members[fast].Window.ActualEnd = time.Time{}, time.Time{}
			case "not-due-as-due":
				members[slow].Status = members[fast].Status
				copyWindow := *members[fast].Window
				members[slow].Window = &copyWindow
			case "not-due-bounds":
				members[slow].Window.AlignedStart = members[slow].Window.AlignedStart.Add(time.Second)
			case "missing-member":
				members = members[:1]
			case "duplicate-member":
				members[slow] = members[fast]
			}
			closure.Buckets[0].Outcome.Members = members
			require.Error(t, store.CommitClosure(t.Context(), testKey, closure))
			require.Zero(t, countRows(t, db, "wg_delivery_mixed_progress"))
			require.Zero(t, countRows(t, db, "wg_delivery_mixed_outcomes"))
			restored, err := store.Restore(t.Context(), testKey)
			require.NoError(t, err)
			require.Nil(t, restored.Mixed)
			require.Len(t, restored.Samples, 2, "rejected schedule cannot consume either member")
		})
	}
}
