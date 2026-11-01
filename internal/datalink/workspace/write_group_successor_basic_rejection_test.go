package workspace

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecreateBasicRejectsChangedShapeAndMeaning(t *testing.T) {
	changes := map[string]func(*testing.T, *circuitProfileHarness, *writeGroupRecreateReservationRequest){
		"partial_sources": func(_ *testing.T, _ *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.Members = r.Group.Members[:1]
		},
		"added_source": func(_ *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			s := h.seedSource(h.fixture.deviceID, "third")
			m := r.Group.Members[0]
			m.PointID, m.TagID, m.TargetColumn = s.pointID, s.tagID, "third"
			r.Group.Members = append(r.Group.Members, m)
		},
		"layout": func(_ *testing.T, _ *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.LayoutVersion = WriteGroupLayoutManagedLongV1
		},
		"profile": func(_ *testing.T, _ *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.RecordingProfile = &WriteGroupRecordingProfile{Kind: WriteGroupProfileInstantaneous}
		},
		"custom": func(_ *testing.T, _ *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.Destination.StorageStrategy = WriteGroupStorageStrategyCustom
		},
		"mode": func(_ *testing.T, _ *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.StorageMode = WriteGroupStorageModeTimescale
		},
		"interval": func(_ *testing.T, _ *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.RowPolicy.IntervalSeconds++
		},
		"lateness": func(_ *testing.T, _ *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.RowPolicy.AllowedLatenessSeconds++
		},
		"write_policy": func(_ *testing.T, _ *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.WritePolicy.Mode = "latest"
		},
		"required": func(_ *testing.T, _ *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.Members[0].Required = false
		},
		"age": func(_ *testing.T, _ *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			age := 10
			r.Group.Members[0].MaxAgeSeconds = &age
		},
		"column": func(_ *testing.T, _ *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.Members[0].TargetColumn = "other"
		},
		"type": func(_ *testing.T, _ *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.Members[0].SourceValueType = "float64"
		},
		"source_revision": func(_ *testing.T, _ *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.Members[0].SourceRevision = "changed"
		},
		"mapping_revision": func(_ *testing.T, _ *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.Members[0].MappingRevision = "changed"
		},
		"binding": func(_ *testing.T, _ *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			id := "other-measurement"
			r.Group.Members[0].MeasurementID = &id
		},
		"device": func(_ *testing.T, _ *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.Members[0].DeviceID = "other"
		},
		"table": func(_ *testing.T, _ *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.Destination.TableName = "arbitrary"
		},
		"destination": func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			recreateRepositoryNewTarget(t, h, r)
		},
		"mysql": func(t *testing.T, h *circuitProfileHarness, _ *writeGroupRecreateReservationRequest) {
			_, err := h.db.ExecContext(t.Context(), `UPDATE database_connectors SET kind='mysql'`)
			require.NoError(t, err)
		},
		"sqlserver": func(t *testing.T, h *circuitProfileHarness, _ *writeGroupRecreateReservationRequest) {
			_, err := h.db.ExecContext(t.Context(), `UPDATE database_connectors SET kind='sqlserver'`)
			require.NoError(t, err)
		},
	}
	for name, change := range changes {
		t.Run(name, func(t *testing.T) {
			h, req := recreateBasicFixture(t, "sqlite")
			change(t, h, &req)
			before := lineageReadSnapshot(t, h)
			_, err := h.service.reserveRecreatedDraft(t.Context(), req)
			require.Error(t, err)
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}

func TestRecreateBasicImmutableAppliedBasisAndExactSource(t *testing.T) {
	for _, mode := range []string{"original", "draft_interval", "missing_version", "foreign_version", "stale_source", "refreshed_source", "unknown_type", "cross_engine"} {
		t.Run(mode, func(t *testing.T) {
			h, req := recreateBasicFixture(t, "sqlite")
			seedRecreateEngineEvidence(t, h, req)
			var err error
			switch mode {
			case "original", "draft_interval":
				_, err = h.db.ExecContext(t.Context(), `UPDATE write_groups SET row_policy=json_set(row_policy,'$.interval_seconds',77),revision='draft-R2'`)
				require.NoError(t, err)
				req.ExpectedPredecessorRevision = "draft-R2"
				if mode == "draft_interval" {
					req.Group, err = h.service.repo.Get(t.Context(), req.WorkspaceID, req.PredecessorGroupID)
					require.NoError(t, err)
				}
			case "missing_version":
				_, err = h.db.ExecContext(t.Context(), `DELETE FROM write_group_versions`)
			case "foreign_version":
				_, err = h.db.ExecContext(t.Context(), `UPDATE write_group_versions SET payload=json_set(payload,'$.workspace_id','foreign')`)
			case "stale_source", "refreshed_source":
				_, err = h.db.ExecContext(t.Context(), `UPDATE tags SET data_type='float64'`)
				if mode == "refreshed_source" {
					for i := range req.Group.Members {
						req.Group.Members[i].SourceRevision = ""
					}
				}
			case "unknown_type":
				_, err = h.db.ExecContext(t.Context(), `PRAGMA ignore_check_constraints=ON; UPDATE tags SET data_type='unknown'`)
			case "cross_engine":
				_, err = h.db.ExecContext(t.Context(), `UPDATE database_connectors SET kind='postgres',connection_config='{"database":"persisted-db","schema":"persisted-schema"}'`)
			}
			require.NoError(t, err)
			before := recreateRepositorySnapshot(t, h)
			result, err := reserveRepositoryTest(t, h, req)
			if mode == "original" {
				require.NoError(t, err)
				require.NotEqual(t, 77, result.Group.RowPolicy.IntervalSeconds)
				return
			}
			require.Error(t, err)
			require.Equal(t, before, recreateRepositorySnapshot(t, h))
		})
	}
}
