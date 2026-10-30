package point

import (
	"errors"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestFormatMutationConcurrentServicesCommitOneRevision(t *testing.T) {
	db, first, request := formatFixture(t)
	second := NewService(NewSQLRepository(db), nil)
	second.SetFormatInvalidator(formatTestInvalidator{})
	gate := make(chan struct{})
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for index, service := range []*Service{first, second} {
		wg.Add(1)
		go func(index int, svc *Service) {
			defer wg.Done()
			<-gate
			req := request
			if index == 1 {
				req.RequestID = "request-2"
			}
			_, err := svc.MutateFormat(t.Context(), req)
			results <- err
		}(index, service)
	}
	close(gate)
	wg.Wait()
	close(results)
	successes, conflicts := 0, 0
	for err := range results {
		if err == nil {
			successes++
		} else if errors.Is(err, ErrFormatConflict) {
			conflicts++
		} else {
			require.NoError(t, err)
		}
	}
	require.Equal(t, 1, successes)
	require.Equal(t, 1, conflicts)
	var count int
	require.NoError(t, db.QueryRow(`SELECT count(*) FROM system_settings WHERE key LIKE 'point_format_receipt:%'`).Scan(&count))
	require.Equal(t, 1, count)
}

func TestFormatMutationExplicitEmptyLegacyRevisionAndReceiptScope(t *testing.T) {
	db, svc, req := formatFixture(t)
	_, err := db.Exec(`UPDATE points SET revision_id=''`)
	require.NoError(t, err)
	empty := ""
	req.ExpectedPointRevision = &empty
	result, err := svc.MutateFormat(t.Context(), req)
	require.NoError(t, err)
	require.NotEmpty(t, result.Point.RevisionID)
	// A reused ID for a different point cannot read the earlier point's receipt.
	_, err = db.Exec(`INSERT INTO points(id,device_id,name,address,data_type,mode,format_version,revision_id) VALUES('point-other',?,'other','HR10','uint32','read',0,'other-rev')`, req.DeviceID)
	require.NoError(t, err)
	other := req.FormatScope
	other.PointID = "point-other"
	_, err = svc.InspectFormatMutation(t.Context(), other)
	require.ErrorIs(t, err, ErrFormatScope)
	req.PointID = "point-other"
	otherRev := "other-rev"
	req.ExpectedPointRevision = &otherRev
	_, err = svc.MutateFormat(t.Context(), req)
	require.ErrorIs(t, err, ErrFormatConflict)
}
