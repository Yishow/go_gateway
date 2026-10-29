package api

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRuntimeSourceEvidencePollingAndRecordingSchedulesStaySeparate(t *testing.T) {
	r := newRuntimeSourceRig(t)
	before := lifecycleAPIState(t, r.f.db)
	v := runtimeSourceView(t, r.get(t, r.query(t)))
	s := v.Source
	require.Equal(t, "saved-only", s.State)
	require.NotNil(t, s.Saved)
	require.NotNil(t, s.Saved.Polling)
	require.Equal(t, "polling-P", s.Saved.Polling.GroupID)
	require.EqualValues(t, 1000, s.Saved.Polling.IntervalMS)
	require.Equal(t, "polling-group", s.Saved.Polling.Source)
	require.Len(t, v.Groups, 1)
	require.NotNil(t, v.Groups[0].Effective)
	require.Equal(t, "mean", v.Groups[0].Effective.Reduction)
	require.EqualValues(t, 15, *v.Groups[0].Effective.IntervalSeconds)
	require.Equal(t, before, lifecycleAPIState(t, r.f.db), "read must not change durable configuration or proof")
}

func TestRuntimeSourceEvidenceSavedR2ScaleDoesNotBecomeEffectiveR1(t *testing.T) {
	r := newRuntimeSourceRig(t)
	s := runtimeSourceView(t, r.get(t, r.query(t))).Source
	require.Equal(t, "R2", *s.SavedRevision)
	require.Equal(t, "saved-only", s.State)
	require.Equal(t, "unavailable", s.Proof)
	require.Nil(t, s.Effective)
	require.Nil(t, s.EffectiveRevision)
	require.Equal(t, "rule-R", s.RuleID)
	require.Equal(t, "point-source", s.PointID)
	require.Equal(t, "mapping-source", *s.MappingID)
	require.Equal(t, "tag-source", *s.TagID)
	require.Equal(t, "40001", *s.Saved.Address)
	require.Equal(t, "FC03", *s.Saved.ProtocolFunction)
	require.Equal(t, "float32", *s.Saved.ReadType)
	require.Equal(t, "float32", *s.Saved.ResultType)
	require.Len(t, s.Saved.Transforms, 3)
	require.Equal(t, "byte-order", s.Saved.Transforms[0].Kind)
	require.Equal(t, "CDAB", *s.Saved.Transforms[0].Order)
	require.Equal(t, "scale", s.Saved.Transforms[1].Kind)
	require.Equal(t, 0.1, *s.Saved.Transforms[1].Multiplier)
	require.Equal(t, 2.0, *s.Saved.Transforms[1].Offset)
	require.Equal(t, "cast", s.Saved.Transforms[2].Kind)
	require.Equal(t, "float32", *s.Saved.Transforms[2].TargetType)
}

func TestRuntimeSourceEvidenceStrictSelectionAndNoImplicitFirstMapping(t *testing.T) {
	r := newRuntimeSourceRig(t)
	for _, key := range []string{"rule_id", "rule_revision", "point_id"} {
		t.Run("missing-"+key, func(t *testing.T) {
			q := r.query(t)
			q.Del(key)
			response := r.get(t, q)
			require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
			require.NotContains(t, response.Body.String(), `"rule_evidence"`)
		})
	}
	q := r.query(t)
	q.Del("mapping_id")
	s := runtimeSourceView(t, r.get(t, q)).Source
	require.Nil(t, s.MappingID)
	require.Nil(t, s.TagID)
	require.Nil(t, s.Saved.MappingRevision)
	require.Nil(t, s.Saved.ResultType)
	require.Len(t, s.Saved.Transforms, 1, "point-only decode must not silently select the linked mapping")
	for _, key := range []string{"rule_id", "rule_revision", "point_id", "mapping_id"} {
		t.Run("duplicate-"+key, func(t *testing.T) {
			q := r.query(t)
			q.Add(key, q.Get(key))
			response := r.get(t, q)
			require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
		})
	}
	q = url.Values{"rule_id": {"rule-R"}, "rule_revision": {"R2"}, "point_id": {"point-source"}}
	response := r.get(t, q)
	require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
}

func TestRuntimeSourceEvidenceOrdinaryContextProvidesObservedRevision(t *testing.T) {
	r := newRuntimeSourceRig(t)
	record, err := r.f.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	response := performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath, nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	data := decodeJSONBody(t, response)["data"].(map[string]any)
	require.Equal(t, record.DatabaseSetupRevision, data["workspace_revision"])
	require.NotContains(t, data, "rule_evidence", "ordinary context must not invent a selected source")
}

func TestRuntimeSourceEvidenceRejectsForeignAndMismatchedBinding(t *testing.T) {
	r := newRuntimeSourceRig(t)
	for _, tc := range []struct {
		name, key, value string
		status           int
	}{
		{"foreign-workspace", "workspace_id", "foreign-workspace", http.StatusForbidden},
		{"foreign-device", "device_id", "foreign-device", http.StatusForbidden},
		{"stale-rule", "rule_revision", "R1", http.StatusConflict},
		{"unlinked-point", "point_id", "point-A", http.StatusConflict},
		{"other-mapping", "mapping_id", "mapping-A", http.StatusConflict},
	} {
		t.Run(tc.name, func(t *testing.T) {
			q := r.query(t)
			q.Set(tc.key, tc.value)
			response := r.get(t, q)
			require.Equal(t, tc.status, response.Code, response.Body.String())
			require.NotContains(t, response.Body.String(), `"rule_evidence"`)
			require.Contains(t, response.Body.String(), `"action"`)
		})
	}
	_, err := r.f.db.Exec(`UPDATE mappings SET tag_id='tag-A' WHERE id='mapping-source'`)
	require.NoError(t, err)
	response := r.get(t, r.query(t))
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	require.NotContains(t, response.Body.String(), `"rule_evidence"`)
}

func TestRuntimeSourceEvidenceChangedPointPollingMappingOrTagIsStale(t *testing.T) {
	for _, tc := range []struct{ name, sql string }{
		{"point-decode", `UPDATE points SET data_format='ABCD' WHERE id='point-source'`},
		{"polling-interval", `UPDATE polling_groups SET interval_ms=2000 WHERE id='polling-P'`},
		{"mapping-scale", `UPDATE mappings SET transform_pipeline='[{"type":"scale","order":1,"params":{"scale":2,"offset":0}}]' WHERE id='mapping-source'`},
		{"tag-type", `UPDATE tags SET data_type='float64' WHERE id='tag-source'`},
		{"rule-revision", `UPDATE source_rules SET revision_id='R3' WHERE id='rule-R'`},
		{"link-binding", `UPDATE source_rule_links SET point_id='point-A' WHERE id='link-source'`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := newRuntimeSourceRig(t)
			r.withChange(func() { _, err := r.f.db.Exec(tc.sql); require.NoError(t, err) })
			response := r.get(t, r.query(t))
			require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
			require.Contains(t, response.Body.String(), `"code":"stale-reference"`)
			require.NotContains(t, response.Body.String(), `"rule_evidence"`)
		})
	}
}

func TestRuntimeSourceEvidenceTelemetryChangesDoNotInvalidateConfiguration(t *testing.T) {
	r := newRuntimeSourceRig(t)
	r.withChange(func() {
		_, err := r.f.db.Exec(`UPDATE points SET last_value='99',last_error='password=private-telemetry',last_read_at='2026-10-29T00:00:00Z',updated_at='2026-10-29T00:00:00Z' WHERE id='point-source'`)
		require.NoError(t, err)
	})
	response := r.get(t, r.query(t))
	s := runtimeSourceView(t, response).Source
	require.Equal(t, "saved-only", s.State)
	require.Equal(t, 0.1, *s.Saved.Transforms[1].Multiplier)
	require.NotContains(t, response.Body.String(), "private-telemetry")
	require.NotContains(t, response.Body.String(), "last_value")
}

func TestRuntimeSourceEvidenceByteOrderScaleCastAndUnsupportedParams(t *testing.T) {
	r := newRuntimeSourceRig(t)
	for _, order := range []string{"ABCD", "BADC", "CDAB", "DCBA"} {
		t.Run(order, func(t *testing.T) {
			_, err := r.f.db.Exec(`UPDATE points SET data_format=? WHERE id='point-source'`, order)
			require.NoError(t, err)
			s := runtimeSourceView(t, r.get(t, r.query(t))).Source
			require.Equal(t, order, *s.Saved.Transforms[0].Order)
		})
	}
	for _, pipeline := range []string{
		`[{"type":"formula","order":1,"params":{"expression":"password=fixture-secret /private/fixture.sqlite"}}]`,
		`[{"type":"scale","order":1,"params":{"scale":0.1,"offset":2,"host":"db.private.internal"}}]`,
		`[{"type":"cast","order":1,"params":{"target_type":"secret://unsafe"}}]`,
		`{ "private": "/private/fixture.sqlite" }`,
	} {
		t.Run(pipeline[:10], func(t *testing.T) {
			_, err := r.f.db.Exec(`UPDATE mappings SET transform_pipeline=? WHERE id='mapping-source'`, pipeline)
			require.NoError(t, err)
			response := r.get(t, r.query(t))
			s := runtimeSourceView(t, response).Source
			require.Equal(t, "saved-only", s.State)
			require.Nil(t, s.Effective)
			require.Contains(t, response.Body.String(), `"kind":"unsupported"`)
			for _, secret := range []string{"fixture-secret", "fixture.sqlite", "db.private.internal", "secret://unsafe"} {
				require.NotContains(t, response.Body.String(), secret)
			}
		})
	}
}

func TestRuntimeSourceEvidenceUnknownPollingAndDecodeStayUnavailable(t *testing.T) {
	r := newRuntimeSourceRig(t)
	_, err := r.f.db.Exec(`UPDATE points SET polling_group_id=NULL,data_format='' WHERE id='point-source'`)
	require.NoError(t, err)
	s := runtimeSourceView(t, r.get(t, r.query(t))).Source
	require.Nil(t, s.Saved.Polling)
	require.Nil(t, s.Effective)
	for _, step := range s.Saved.Transforms {
		require.NotEqual(t, "byte-order", step.Kind, "missing byte order must not become an invented ABCD")
	}
}

func TestRuntimeSourceEvidenceMalformedQueriesFailClosed(t *testing.T) {
	r := newRuntimeSourceRig(t)
	for _, query := range []string{
		r.query(t).Encode() + "&rule_id=%ZZ",
		r.query(t).Encode() + "&point_id=bad;query",
		"workspace_id=W&workspace_revision=R&device_id=D&rule_id=%ZZ&rule_revision=%ZZ&point_id=%ZZ",
	} {
		response := performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath+"?"+query, nil)
		require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
		require.NotContains(t, response.Body.String(), `"rule_evidence"`)
	}
	for _, key := range []string{"rule_id", "rule_revision", "point_id", "mapping_id"} {
		for _, unsafe := range []string{"", "/private/fixture.sqlite", "db.private.internal"} {
			q := r.query(t)
			q.Set(key, unsafe)
			response := r.get(t, q)
			require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
			require.NotContains(t, response.Body.String(), "private")
		}
	}
}

func TestRuntimeSourceEvidenceNeverBootstrapsOrMutates(t *testing.T) {
	r := newRuntimeSourceRig(t)
	q := r.query(t)
	var before, after int
	require.NoError(t, r.f.db.QueryRow(`SELECT total_changes()`).Scan(&before))
	runtimeSourceView(t, r.get(t, q))
	require.NoError(t, r.f.db.QueryRow(`SELECT total_changes()`).Scan(&after))
	require.Equal(t, before, after, "source read must issue no configuration or proof writes")
	_, err := r.f.db.Exec(`DELETE FROM system_settings WHERE key='studio_v2_workspace'`)
	require.NoError(t, err)
	response := r.get(t, q)
	var count int
	require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM system_settings WHERE key='studio_v2_workspace'`).Scan(&count))
	require.Zero(t, count, "read must not recreate an absent workspace")
	require.Equal(t, http.StatusServiceUnavailable, response.Code, response.Body.String())
	require.NotContains(t, response.Body.String(), `"rule_evidence"`)
}

func TestRuntimeSourceEvidenceForeignPointCannotBorrowOwnedLink(t *testing.T) {
	r := newRuntimeSourceRig(t)
	_, err := r.f.db.Exec(`INSERT INTO devices (id,name,description,protocol,status,connection_config,readiness_status) VALUES ('device-foreign','Foreign','','modbus_tcp','draft','{}','{}')`)
	require.NoError(t, err)
	_, err = r.f.db.Exec(`UPDATE points SET device_id='device-foreign' WHERE id='point-source'`)
	require.NoError(t, err)
	response := r.get(t, r.query(t))
	require.Equal(t, http.StatusForbidden, response.Code, response.Body.String())
	require.Contains(t, response.Body.String(), `"code":"foreign-scope"`)
	require.NotContains(t, response.Body.String(), `"rule_evidence"`)
}

func TestRuntimeSourceEvidenceConcurrentWorkspaceDeletionDoesNotRecreate(t *testing.T) {
	r := newRuntimeSourceRig(t)
	q := r.query(t)
	r.withChange(func() {
		_, err := r.f.db.Exec(`DELETE FROM system_settings WHERE key='studio_v2_workspace'`)
		require.NoError(t, err)
	}, 2)
	response := r.get(t, q)
	var count int
	require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM system_settings WHERE key='studio_v2_workspace'`).Scan(&count))
	require.Zero(t, count, "final source consistency read must not recreate a deleted workspace")
	require.Equal(t, http.StatusServiceUnavailable, response.Code, response.Body.String())
	require.NotContains(t, response.Body.String(), `"rule_evidence"`)
}
