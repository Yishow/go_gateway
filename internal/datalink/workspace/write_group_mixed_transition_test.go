package workspace

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLegacyToMixedRequiresIndependentRecordAndFreshTarget(t *testing.T) {
	ctx := t.Context()
	h := newCircuitProfileHarness(t, ":memory:")
	legacy := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "", "", h.primary()))
	before := cloneWriteGroup(legacy)
	mixed := mixedDraft(h, WriteGroupMemberCalculationMean, 15, "new_mixed_target")
	mixed.LayoutVersion = WriteGroupLayoutManagedWideProfileV1
	converted := cloneWriteGroup(mixed)
	converted.ID = legacy.ID
	request := h.mutation(ctx, converted)
	request.ExpectedGroupRevision = legacy.Revision
	_, err := h.service.Update(ctx, legacy.ID, request)
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	// The same lineage cannot silently change its create-once canonical key.
	tx, err := h.db.BeginTx(ctx, nil)
	require.NoError(t, err)
	_, err = NewSQLWriteGroupRepository(h.db).recreatedProfileIdentityInTx(ctx, tx, legacy, converted)
	var rejected *WriteGroupTransitionError
	require.ErrorAs(t, err, &rejected)
	require.Equal(t, "canonical-identity-change", rejected.ReasonCode)
	require.NoError(t, tx.Rollback())
	created, err := h.service.Create(ctx, h.mutation(ctx, mixed))
	require.NoError(t, err)
	require.NotEqual(t, legacy.ID, created.Group.ID)
	invalidShape := cloneWriteGroup(created.Group)
	invalidShape.RecordingProfile.ResultShape = "future-shape"
	invalidRequest := h.mutation(ctx, invalidShape)
	invalidRequest.ExpectedGroupRevision = created.Group.Revision
	_, err = h.service.Update(ctx, created.Group.ID, invalidRequest)
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	untouched, err := h.service.Get(ctx, created.Group.ID)
	require.NoError(t, err)
	require.Equal(t, created.Group, untouched.Group)

	require.Equal(t, mixed.RecordID, created.Group.RecordID)
	require.Empty(t, created.Group.AppliedRevision)
	require.Empty(t, created.Group.Destination.SchemaRevision)
	require.Empty(t, created.Group.Destination.SchemaDigest)
	require.NotEqual(t, legacy.Destination.TableName, created.Group.Destination.TableName)
	after, err := h.service.Get(ctx, legacy.ID)
	require.NoError(t, err)
	require.Equal(t, before, after.Group)
	require.Equal(t, 2, h.countRows(`SELECT COUNT(*) FROM write_group_lineages`))
	again := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "", "", h.primary()))
	require.Equal(t, legacy.ID, again.ID)
	collision := mixedDraft(h, WriteGroupMemberCalculationMean, 15, legacy.Destination.TableName)
	collision.LayoutVersion = WriteGroupLayoutManagedWideProfileV1
	_, err = h.service.Create(ctx, h.mutation(ctx, collision))
	require.ErrorIs(t, err, ErrWriteGroupCircuitProfileConflict)
}
func TestMixedSchemaRejectsUnconfirmedPhysicalMeaning(t *testing.T) {
	for _, field := range []string{"provenance", "record", "dispatch", "legacy bucket", "unknown dialect"} {
		t.Run(field, func(t *testing.T) {
			g, types := mixedSchemaFixture()
			dialect := "sqlite"
			switch field {
			case "provenance":
				g.RowPolicy.ProvenanceColumn = ""
			case "record":
				g.RowPolicy.RecordKeyColumn = ""
			case "dispatch":
				g.RowPolicy.DueAtColumn = ""
			case "legacy bucket":
				g.RowPolicy.BucketStartColumn = "bucket_start"
			case "unknown dialect":
				dialect = "mysql"
			}
			_, err := managedSchemaLayout(g, types, dialect)
			require.ErrorIs(t, err, ErrWriteGroupValidation)
		})
	}
}

func TestLegacyLongToIndependentMixedKeepsExistingSharedStorageAuthority(t *testing.T) {
	ctx := t.Context()
	h := newCircuitProfileHarness(t, ":memory:")
	legacy := h.mustEnsure(h.longDraft(WriteGroupProfileInstantaneous, "main", h.primary()))
	before := cloneWriteGroup(legacy)
	mixed := mixedDraft(h, WriteGroupMemberCalculationMean, 15, "")
	mixed.LayoutVersion = WriteGroupLayoutManagedLongV1
	mixed.StorageMode = WriteGroupStorageModePlain
	mixed.Destination = legacy.Destination
	mixed.Destination.TableName = ""
	mixed.Destination.SchemaRevision = ""
	mixed.Destination.SchemaDigest = ""
	created, err := h.service.Create(ctx, h.mutation(ctx, mixed))
	require.NoError(t, err)
	require.NotEqual(t, legacy.ID, created.Group.ID)
	require.NotEqual(t, legacy.RecordID, created.Group.RecordID)
	require.Equal(t, longStorageScopeOf(legacy), longStorageScopeOf(created.Group))
	require.Equal(t, legacy.Destination.TableName, created.Group.Destination.TableName)
	require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_group_storage_scopes`))
	require.Equal(t, 2, h.countRows(`SELECT COUNT(*) FROM write_group_lineages`))
	require.Empty(t, created.Group.Destination.SchemaRevision)
	require.Empty(t, created.Group.Destination.SchemaDigest)
	require.Empty(t, created.Group.AppliedRevision)
	after, err := h.service.Get(ctx, legacy.ID)
	require.NoError(t, err)
	require.Equal(t, before, after.Group)
	// A shared physical table grants neither the old proof nor same-lineage conversion.
	converted := cloneWriteGroup(created.Group)
	converted.ID = legacy.ID
	req := h.mutation(ctx, converted)
	req.ExpectedGroupRevision = legacy.Revision
	_, err = h.service.Update(ctx, legacy.ID, req)
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	require.Zero(t, h.countRows(`SELECT COUNT(*) FROM wg_runtime_versions`))
}
