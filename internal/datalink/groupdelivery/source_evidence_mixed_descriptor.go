package groupdelivery

import (
	"context"
	"encoding/json"
	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"time"
)

type mixedSourceDescriptor struct {
	Group struct {
		RecordID    string `json:"record_id"`
		Layout      string `json:"layout_version"`
		Destination struct {
			ConnectorID       string `json:"connector_id"`
			ConnectorRevision string `json:"connector_revision"`
		} `json:"destination"`
		RowPolicy struct {
			DueAtColumn      string `json:"due_at_column"`
			RecordKeyColumn  string `json:"record_key_column"`
			ProvenanceColumn string `json:"provenance_column"`
		} `json:"row_policy"`
		Members []struct {
			DeviceID        string                `json:"device_id"`
			PointID         string                `json:"point_id"`
			TagID           string                `json:"tag_id"`
			CircuitID       string                `json:"circuit_id"`
			Column          string                `json:"target_column"`
			MeasurementID   string                `json:"measurement_id"`
			SourceRevision  string                `json:"source_revision"`
			MappingRevision string                `json:"mapping_revision"`
			Required        bool                  `json:"required"`
			SourceType      schema.DataType       `json:"source_value_type"`
			ResultType      measurement.ExactType `json:"result_value_type"`
		} `json:"members"`
	} `json:"group"`
	Layout struct {
		Dialect dbtarget.SQLDialect   `json:"dialect"`
		Columns []dbtarget.ColumnInfo `json:"columns"`
	} `json:"runtime_layout"`
}

type mixedSourceIntent struct {
	method   snapshot.Calculation
	interval time.Duration
	required bool
}
type mixedSourceFacts struct {
	intents                         map[string]mixedSourceIntent
	member                          dbtarget.GroupRowMember
	dialect                         dbtarget.SQLDialect
	declared                        string
	sourceRevision, mappingRevision string
	counterRef                      counterproof.CounterBindingReference
}

func (s *Store) mixedSourceFacts(ctx context.Context, scope SourceScope, ref SourceRef) (mixedSourceFacts, error) {
	var facts mixedSourceFacts
	key := GroupKey{scope.WorkspaceID, scope.GroupID, scope.GroupRevision}
	d, descriptorDigest, err := readMixedDescriptor(ctx, s.db, key)
	if err != nil {
		return facts, err
	}
	if d == nil || scope.ResultShape != "member-calculation-v1" {
		return facts, mixedInvalid("source descriptor")
	}
	if _, err := d.intervals(); err != nil {
		return facts, err
	}
	var raw string
	if err = s.db.QueryRowContext(ctx, `SELECT payload FROM wg_runtime_versions WHERE workspace_id=? AND group_id=? AND group_revision=?`, scope.WorkspaceID, scope.GroupID, scope.GroupRevision).Scan(&raw); err != nil {
		return facts, err
	}
	if handoffHash(raw) != descriptorDigest {
		return facts, mixedInvalid("source descriptor changed")
	}
	var frozen mixedSourceDescriptor
	if json.Unmarshal([]byte(raw), &frozen) != nil {
		return facts, mixedInvalid("source descriptor JSON")
	}
	g := frozen.Group
	if g.RecordID == "" || g.RecordID != scope.RecordIdentity || g.Destination.ConnectorID != scope.ConnectorID || g.Destination.ConnectorRevision != scope.ConnectorRevision || (g.Layout == "managed_wide_profile_v1") != scope.Wide || (g.Layout != "managed_wide_profile_v1" && g.Layout != "managed_long_v1") {
		return facts, mixedInvalid("source scope")
	}
	if frozen.Layout.Dialect != dbtarget.SQLDialectSQLite && frozen.Layout.Dialect != dbtarget.SQLDialectPostgres {
		return facts, mixedInvalid("source dialect")
	}
	if scope.Wide && (scope.DueAtColumn != g.RowPolicy.DueAtColumn || scope.RecordKeyColumn != g.RowPolicy.RecordKeyColumn || scope.ProvenanceColumn != g.RowPolicy.ProvenanceColumn) {
		return facts, mixedInvalid("source columns")
	}
	memberKey := SourceEntityKey(ref.DeviceID, ref.PointID, ref.TagID)
	facts.intents = map[string]mixedSourceIntent{}
	for _, m := range d.Group.Members {
		interval, e := snapshot.CheckedCadenceDuration(m.Cadence.Interval, m.Cadence.Unit)
		if e != nil {
			return facts, e
		}
		facts.intents[SourceEntityKey(m.DeviceID, m.PointID, m.TagID)] = mixedSourceIntent{method: snapshot.Calculation(m.Calculation), interval: interval}
	}
	for _, m := range g.Members {
		k := SourceEntityKey(m.DeviceID, m.PointID, m.TagID)
		v, ok := facts.intents[k]
		if !ok {
			return facts, mixedInvalid("source metadata member set")
		}
		v.required = m.Required
		facts.intents[k] = v
	}
	found := false
	for _, m := range d.Group.Members {
		if SourceEntityKey(m.DeviceID, m.PointID, m.TagID) != memberKey {
			continue
		}
		interval, e := snapshot.CheckedCadenceDuration(m.Cadence.Interval, m.Cadence.Unit)
		requested, e2 := snapshot.CheckedCadenceDuration(ref.CadenceSeconds, "second")
		if e != nil || e2 != nil || interval != requested || m.Calculation != ref.Calculation || m.Cadence.Alignment != ref.Alignment {
			return facts, mixedInvalid("source intent")
		}
		facts.member = dbtarget.GroupRowMember{MemberKey: memberKey, Calculation: snapshot.Calculation(m.Calculation), Interval: interval}
		found = true
	}
	if !found {
		return facts, mixedInvalid("source member")
	}
	found = false
	for _, m := range g.Members {
		if SourceEntityKey(m.DeviceID, m.PointID, m.TagID) != memberKey {
			continue
		}
		if m.CircuitID != ref.CircuitID || m.Column != ref.Column || m.Required != ref.Required || string(m.SourceType) != ref.SourceType || string(m.ResultType) != ref.ResultType || m.MeasurementID != ref.MeasurementID {
			return facts, mixedInvalid("source frozen metadata")
		}
		facts.sourceRevision = m.SourceRevision
		facts.mappingRevision = m.MappingRevision
		facts.member.Column = m.Column
		facts.member.Required = m.Required
		facts.member.SourceValueType = m.SourceType
		facts.member.Type = m.ResultType
		found = true
	}
	if !found {
		return facts, mixedInvalid("source metadata")
	}
	facts.dialect = frozen.Layout.Dialect
	column := ref.Column
	if !scope.Wide {
		column = sourceSlotByType[ref.ResultType]
	}
	for _, c := range frozen.Layout.Columns {
		if c.Name == column {
			facts.declared = c.DataType
		}
	}
	if facts.declared == "" {
		return facts, mixedInvalid("source value column")
	}
	if ref.Calculation == "delta" {
		proofs, e := counterproof.ResolveSnapshot([]byte(raw), counterproof.SnapshotScope{WorkspaceID: scope.WorkspaceID, GroupID: scope.GroupID, AppliedRevision: scope.GroupRevision})
		if e != nil {
			return facts, e
		}
		for _, p := range proofs {
			if p.Member.SampleMemberKey() != memberKey {
				continue
			}
			mod, rate, e := p.CounterPolicy.ExactRollover()
			if e != nil {
				return facts, e
			}
			scale := ""
			if p.ScaleMode == counterproof.CounterScaleExact {
				scale = p.ExactScale
			}
			facts.counterRef, e = p.Reference()
			if e != nil {
				return facts, e
			}
			facts.member.Counter = &aggregation.CounterBinding{WorkspaceID: scope.WorkspaceID, GroupID: scope.GroupID, ProfileRevision: scope.GroupRevision, MemberKey: memberKey, MeasurementID: p.MeasurementID, DefinitionRevision: p.DefinitionRevision, SourceBindingRevision: p.SourceBindingRevision, SeriesEpoch: p.SeriesEpoch, Kind: aggregation.CounterKind(p.SemanticKind), Unit: p.Unit, ExactRecordingScale: scale, RolloverModulus: mod, MaxRatePerSecond: rate}
		}
		if facts.member.Counter == nil {
			return facts, mixedInvalid("source counter proof")
		}
	}
	return facts, nil
}
