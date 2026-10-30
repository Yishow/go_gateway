package runtime

import (
	"context"
	"slices"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

// MappingObserver observes only an executed transform result. Implementations
// must be nonblocking; diagnostic failure cannot change native mapping outcomes.
type MappingObserver func(MappingObservation)

// MappingObservationResult is a closed classification, never raw error text.
type MappingObservationResult string

const (
	MappingObservationFailed      MappingObservationResult = "failed"
	MappingObservationSucceeded   MappingObservationResult = "succeeded"
	MappingObservationUnavailable MappingObservationResult = "unavailable"
)

// MappingObservation contains execution identity, not values or transform data.
// MappingRevision is the installed canonical mapping revision, not a fingerprint.
type MappingObservation struct {
	WorkspaceID, WorkspaceRevision string
	DeviceID, PointID, TagID       string
	MappingID, MappingRevision     string
	AcquisitionID                  string
	ReceivedAt                     time.Time
	Result                         MappingObservationResult
}

type mappingObservationScope struct{ workspaceID, workspaceRevision string }

func (s *Service) mappingObservationSnapshot(ctx context.Context) *workspace.Record {
	if s == nil || s.mappingObserver == nil {
		return nil
	}
	return s.acquisitionWorkspaceSnapshot(ctx)
}

// annotateMappingObservation touches only an unpublished local map. The native
// sample WorkspaceID is deliberately not rewritten when generic RefreshMappings
// installs records without it. Scope here describes this installed input only.
func (s *Service) annotateMappingObservation(ctx context.Context, before *workspace.Record, next map[string][]mappingBinding, projection *workspace.RuntimeProjection) {
	if before == nil || before.ID == "" || before.DatabaseSetupRevision == "" {
		return
	}
	after := s.mappingObservationSnapshot(ctx)
	if after == nil || after.ID != before.ID || after.DatabaseSetupRevision != before.DatabaseSetupRevision || !slices.Equal(after.OrderedDeviceIDs, before.OrderedDeviceIDs) {
		return
	}
	if projection != nil && projection.WorkspaceID != before.ID {
		return
	}
	inputCounts := mappingObservationInputCounts(projection)
	ids, tuples := make(map[string]int), make(map[[3]string]int)
	for _, bindings := range next {
		for _, b := range bindings {
			ids[b.MappingID]++
			tuples[[3]string{b.DeviceID, b.PointID, b.TagID}]++
		}
	}
	for pointID, bindings := range next {
		for i, b := range bindings {
			if b.MappingID == "" || b.MappingRevision == "" || b.ConfigFingerprint == "" || b.DeviceID == "" || b.PointID == "" || b.TagID == "" || b.PointID != pointID || ids[b.MappingID] != 1 || tuples[[3]string{b.DeviceID, b.PointID, b.TagID}] != 1 || !slices.Contains(before.OrderedDeviceIDs, b.DeviceID) {
				continue
			}
			if b.WorkspaceID != "" && b.WorkspaceID != before.ID {
				continue
			}
			if projection != nil && !uniqueMappingObservationInput(inputCounts, b) {
				continue
			}
			bindings[i].observationScope = mappingObservationScope{before.ID, before.DatabaseSetupRevision}
		}
	}
}

// Counts exist only while preparing this local publication; no result-time map
// or persistent diagnostic authority is introduced.
func mappingObservationInputCounts(p *workspace.RuntimeProjection) map[[2]string]int {
	if p == nil {
		return nil
	}
	counts := make(map[[2]string]int)
	for _, id := range p.DeviceIDs {
		counts[[2]string{"member", id}]++
	}
	for _, v := range p.Devices {
		if v != nil {
			counts[[2]string{"device", v.ID}]++
		}
	}
	for _, v := range p.Points {
		if v != nil {
			counts[[2]string{"point", v.ID}]++
		}
	}
	for _, v := range p.Tags {
		if v != nil {
			counts[[2]string{"tag", v.ID}]++
		}
	}
	for _, v := range p.Mappings {
		if v != nil {
			counts[[2]string{"mapping", v.ID}]++
		}
	}
	return counts
}

func uniqueMappingObservationInput(counts map[[2]string]int, b mappingBinding) bool {
	return counts[[2]string{"member", b.DeviceID}] == 1 && counts[[2]string{"device", b.DeviceID}] == 1 && counts[[2]string{"point", b.PointID}] == 1 && counts[[2]string{"tag", b.TagID}] == 1 && counts[[2]string{"mapping", b.MappingID}] == 1
}

// observeMappingResult uses the exact binding already copied for execution, not
// a current mapping lookup. A successful transform is not acquisition or SQL
// success; invalid acquisition evidence cannot close an observed mapping fault.
func (s *Service) observeMappingResult(cv collector.CollectedValue, b mappingBinding, succeeded bool) {
	if s == nil || s.mappingObserver == nil {
		return
	}
	defer func() { _ = recover() }()
	o := MappingObservation{Result: MappingObservationUnavailable}
	if b.observationScope.workspaceID != "" && b.observationScope.workspaceRevision != "" && b.MappingID != "" && b.MappingRevision != "" &&
		cv.DeviceID == b.DeviceID && cv.PointID == b.PointID && cv.AcquisitionID != "" && !cv.ReceivedAt.IsZero() &&
		cv.ConfigFingerprint != "" && cv.ConfigFingerprint == b.ConfigFingerprint && cv.Quality == schema.QualityGood && cv.QualityReason == "" && cv.Error == "" {
		result := MappingObservationFailed
		if succeeded {
			result = MappingObservationSucceeded
		}
		o = MappingObservation{WorkspaceID: b.observationScope.workspaceID, WorkspaceRevision: b.observationScope.workspaceRevision, DeviceID: b.DeviceID, PointID: b.PointID, TagID: b.TagID, MappingID: b.MappingID, MappingRevision: b.MappingRevision, AcquisitionID: cv.AcquisitionID, ReceivedAt: cv.ReceivedAt, Result: result}
	}
	s.mappingObserver(o)
}
