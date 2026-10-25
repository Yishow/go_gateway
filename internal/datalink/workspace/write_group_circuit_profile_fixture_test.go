package workspace

import (
	"context"
	"database/sql"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type circuitSource struct {
	deviceID, pointID, tagID string
}

type circuitProfileHarness struct {
	t           *testing.T
	db          *sql.DB
	fixture     writeGroupFixture
	workspace   *Service
	service     *WriteGroupService
	nextAddress int
}

func newCircuitProfileHarness(t *testing.T, dsn string) *circuitProfileHarness {
	t.Helper()
	db := openWorkspaceTestDB(t, dsn)
	t.Cleanup(func() { _ = db.Close() })
	fixture := seedWriteGroupFixture(t.Context(), t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	return &circuitProfileHarness{
		t: t, db: db, fixture: fixture, workspace: workspaceSvc,
		service:     NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db)),
		nextAddress: 41000,
	}
}

// newCircuitProfileHarnessOn models a process restart: fresh services over a
// reopened database that still holds the seeded fixture.
func newCircuitProfileHarnessOn(t *testing.T, db *sql.DB, previous *circuitProfileHarness) *circuitProfileHarness {
	t.Helper()
	t.Cleanup(func() { _ = db.Close() })
	workspaceSvc := NewService(NewSQLRepository(db))
	return &circuitProfileHarness{
		t: t, db: db, fixture: previous.fixture, workspace: workspaceSvc,
		service:     NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db)),
		nextAddress: previous.nextAddress,
	}
}

// primary is the source seeded by seedWriteGroupFixture.
func (h *circuitProfileHarness) primary() circuitSource {
	return circuitSource{deviceID: h.fixture.deviceID, pointID: h.fixture.pointID, tagID: h.fixture.tagID}
}

// seedSource adds one mapped point and tag to an existing device.
func (h *circuitProfileHarness) seedSource(deviceID, key string) circuitSource {
	h.t.Helper()
	ctx := h.t.Context()
	h.nextAddress++
	source := circuitSource{deviceID: deviceID, pointID: uuid.NewString(), tagID: uuid.NewString()}
	_, err := h.db.ExecContext(ctx, `
		INSERT INTO points (id, device_id, name, address, function, data_type, mode, enabled)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, source.pointID, deviceID, key, h.nextAddress, "FC03", schema.DataTypeFloat32, schema.PointModeReadOnly, true)
	require.NoError(h.t, err)
	_, err = h.db.ExecContext(ctx, `
		INSERT INTO tags (id, key, key_lower, display_name, data_type, status, labels)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, source.tagID, key, key, key, schema.DataTypeFloat32, schema.TagStatusActive, `{}`)
	require.NoError(h.t, err)
	_, err = h.db.ExecContext(ctx, `
		INSERT INTO mappings (id, point_id, tag_id, transform_pipeline, status, enabled)
		VALUES (?, ?, ?, ?, ?, ?)
	`, uuid.NewString(), source.pointID, source.tagID, `[]`, schema.MappingStatusActive, true)
	require.NoError(h.t, err)
	return source
}

// draft builds a managed circuit-profile group. An empty circuitID models the
// unassigned raw selection of one device.
func (h *circuitProfileHarness) draft(kind WriteGroupProfileKind, circuitID, circuitName string, sources ...circuitSource) *WriteGroup {
	group := newBasicManagedDraft(h.fixture)
	group.Name = circuitName + " " + string(kind)
	group.Members = nil
	for _, source := range sources {
		group.Members = append(group.Members, WriteGroupMember{
			DeviceID: source.deviceID, PointID: source.pointID, TagID: source.tagID,
			Required: true, CircuitID: circuitID, CircuitName: circuitName,
		})
	}
	profile, err := DefaultWriteGroupRecordingProfile(kind)
	require.NoError(h.t, err)
	// The shared fixture pins the legacy latest semantics explicitly so every
	// test built on it keeps protecting the latest path; the mean default has
	// its own compatibility tests.
	if kind == WriteGroupProfileInstantaneous {
		profile.SampleReduction = WriteGroupSampleReductionLatest
	}
	group.RecordingProfile = &profile
	return group
}

func (h *circuitProfileHarness) mutation(ctx context.Context, group *WriteGroup) WriteGroupMutation {
	h.t.Helper()
	current, err := h.workspace.GetOrCreate(ctx)
	require.NoError(h.t, err)
	return WriteGroupMutation{
		WorkspaceID: h.fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: group,
	}
}

func (h *circuitProfileHarness) ensure(group *WriteGroup) (*WriteGroupSaveResult, error) {
	h.t.Helper()
	return h.service.EnsureCircuitProfile(h.t.Context(), h.mutation(h.t.Context(), group))
}

func (h *circuitProfileHarness) mustEnsure(group *WriteGroup) *WriteGroup {
	h.t.Helper()
	result, err := h.ensure(group)
	require.NoError(h.t, err)
	require.NotNil(h.t, result)
	require.NotNil(h.t, result.Group)
	return result.Group
}

func (h *circuitProfileHarness) countRows(query string, args ...any) int {
	h.t.Helper()
	var count int
	require.NoError(h.t, h.db.QueryRowContext(h.t.Context(), query, args...).Scan(&count))
	return count
}
