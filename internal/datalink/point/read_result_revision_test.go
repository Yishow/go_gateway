package point

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"testing"
)

type legacyOnlyReadResultRepo struct{ Repository }

func TestReadResultRevisionCASPreservesNewConfiguration(t *testing.T) {
	for _, kind := range []string{"memory", "sqlite"} {
		t.Run(kind, func(t *testing.T) {
			var repo Repository
			deviceID := "D"
			if kind == "sqlite" {
				db := setupTestDB(t)
				defer db.Close()
				deviceID, _ = setupTestData(t, db)
				repo = NewSQLRepository(db)
			} else {
				repo = NewMemoryRepository()
			}
			p := &schema.Point{ID: "P", DeviceID: deviceID, Name: "P", Address: "D0", DataType: schema.DataTypeUint32, Mode: schema.PointModeReadOnly, RevisionID: "r1", Enabled: true}
			require.NoError(t, repo.Create(t.Context(), p))
			svc := NewService(repo, nil)
			require.ErrorIs(t, svc.UpdateReadResultAtRevision(t.Context(), p.ID, deviceID, "old", uint32(7), ""), ErrReadResultSuperseded)
			require.ErrorIs(t, svc.UpdateReadResultAtRevision(t.Context(), p.ID, "foreign", p.RevisionID, uint32(7), ""), ErrReadResultSuperseded)
			current, err := repo.GetByID(t.Context(), p.ID)
			require.NoError(t, err)
			require.Nil(t, current.LastValue)
			require.Nil(t, current.LastReadAt)
			require.NoError(t, svc.UpdateReadResultAtRevision(t.Context(), p.ID, deviceID, p.RevisionID, uint32(123), ""))
			current, err = repo.GetByID(t.Context(), p.ID)
			require.NoError(t, err)
			require.Equal(t, "123", *current.LastValue)
			require.Equal(t, "r1", current.RevisionID)
			current.RevisionID = "r2"
			require.NoError(t, repo.Update(t.Context(), current))
			before, err := repo.GetByID(t.Context(), p.ID)
			require.NoError(t, err)
			require.ErrorIs(t, svc.UpdateReadResultAtRevision(t.Context(), p.ID, deviceID, "r1", uint32(456), ""), ErrReadResultSuperseded)
			after, err := repo.GetByID(t.Context(), p.ID)
			require.NoError(t, err)
			require.Equal(t, before, after)
			require.ErrorIs(t, NewService(legacyOnlyReadResultRepo{repo}, nil).UpdateReadResultAtRevision(t.Context(), p.ID, deviceID, "r2", uint32(456), ""), ErrReadResultCASUnavailable)
		})
	}
}
