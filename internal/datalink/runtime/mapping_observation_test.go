package runtime

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestMappingObservationUsesCanonicalInstalledRevision(t *testing.T) {
	f := newMappingObservationFixture(t, true)
	f.install(t)
	want, err := measurement.MappingRevision(*f.projection.Mappings[0])
	require.NoError(t, err)
	got := f.execute(t, uint16(65535))
	require.Equal(t, MappingObservationFailed, got.Result)
	require.Equal(t, "W", got.WorkspaceID)
	require.Equal(t, "WR1", got.WorkspaceRevision)
	require.Equal(t, "D", got.DeviceID)
	require.Equal(t, "P", got.PointID)
	require.Equal(t, "T", got.TagID)
	require.Equal(t, "M", got.MappingID)
	require.Equal(t, want, got.MappingRevision)
	require.NotEqual(t, f.service.mappingIndex["P"][0].ConfigFingerprint, got.MappingRevision)
	require.Equal(t, "acquisition-A", got.AcquisitionID)
	require.False(t, got.ReceivedAt.IsZero())
	require.EqualValues(t, 1, f.service.Snapshot().MappingError)
	require.Empty(t, f.writer.records)
	require.Len(t, f.sink.samples, 1)
	require.Equal(t, reasonMappingFailed, f.sink.samples[0].QualityReason)
	require.Equal(t, MappingObservationSucceeded, f.execute(t, uint16(10)).Result)
	require.Len(t, f.writer.records, 1)
	require.Equal(t, 2, f.reader.reads)
	require.Zero(t, f.reader.projections)
}

func TestMappingObservationNilObserverAddsNoReads(t *testing.T) {
	f := newMappingObservationFixture(t, false)
	f.reader.panicRead = true
	f.install(t)
	f.service.handleCollectedValue(t.Context(), f.value(t, uint16(10)))
	require.NoError(t, f.service.RefreshMappings(t.Context()))
	require.Zero(t, f.reader.reads)
	require.Len(t, f.writer.records, 1)
}

func TestMappingObservationMissingChangedOrAmbiguousScopeUnavailable(t *testing.T) {
	cases := map[string]func(mappingObservationFixture){
		"missing snapshot":  func(f mappingObservationFixture) { f.reader.record = nil },
		"empty revision":    func(f mappingObservationFixture) { f.reader.record.DatabaseSetupRevision = "" },
		"foreign workspace": func(f mappingObservationFixture) { f.reader.record.ID = "foreign" },
		"foreign device":    func(f mappingObservationFixture) { f.reader.record.OrderedDeviceIDs = []string{"foreign"} },
		"read failure":      func(f mappingObservationFixture) { f.reader.failure = errors.New("password_fixture") },
		"reader panic":      func(f mappingObservationFixture) { f.reader.panicRead = true },
		"revision changed": func(f mappingObservationFixture) {
			f.reader.beforeRead = func(n int) {
				if n == 2 {
					f.reader.record.DatabaseSetupRevision = "WR2"
				}
			}
		},
		"deleted during read": func(f mappingObservationFixture) {
			f.reader.beforeRead = func(n int) {
				if n == 2 {
					f.reader.record = nil
				}
			}
		},
		"duplicate point": func(f mappingObservationFixture) {
			f.projection.Points = append(f.projection.Points, f.projection.Points[0])
		},
		"duplicate device": func(f mappingObservationFixture) {
			f.projection.Devices = append(f.projection.Devices, f.projection.Devices[0])
		},
		"duplicate tag": func(f mappingObservationFixture) { f.projection.Tags = append(f.projection.Tags, f.projection.Tags[0]) },
		"duplicate mapping": func(f mappingObservationFixture) {
			f.projection.Mappings = append(f.projection.Mappings, f.projection.Mappings[0])
		},
		"empty mapping ID": func(f mappingObservationFixture) { f.projection.Mappings[0].ID = "" },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			f := newMappingObservationFixture(t, true)
			change(f)
			f.install(t)
			f.service.handleCollectedValue(t.Context(), f.value(t, uint16(10)))
			require.NotEmpty(t, *f.seen)
			for _, o := range *f.seen {
				require.Equal(t, MappingObservation{Result: MappingObservationUnavailable}, o)
			}
			require.NotEmpty(t, f.writer.records)
			require.Zero(t, f.reader.projections)
		})
	}
}

func TestMappingObservationFailedAcquisitionCannotRecoverFault(t *testing.T) {
	cases := map[string]func(*collector.CollectedValue){
		"foreign device":        func(cv *collector.CollectedValue) { cv.DeviceID = "other" },
		"fingerprint mismatch":  func(cv *collector.CollectedValue) { cv.ConfigFingerprint = "old-config" },
		"missing acquisition":   func(cv *collector.CollectedValue) { cv.AcquisitionID = "" },
		"missing received time": func(cv *collector.CollectedValue) { cv.ReceivedAt = time.Time{} },
		"bad quality":           func(cv *collector.CollectedValue) { cv.Quality = schema.QualityBad; cv.QualityReason = "read-failed" },
		"unknown quality":       func(cv *collector.CollectedValue) { cv.Quality = "" },
		"contradictory error":   func(cv *collector.CollectedValue) { cv.Error = "password_fixture" },
		"contradictory reason":  func(cv *collector.CollectedValue) { cv.QualityReason = "read-failed" },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			f := newMappingObservationFixture(t, true)
			f.install(t)
			cv := f.value(t, uint16(10))
			change(&cv)
			f.service.handleCollectedValue(t.Context(), cv)
			require.Equal(t, []MappingObservation{{Result: MappingObservationUnavailable}}, *f.seen)
			require.Len(t, f.writer.records, 1)
		})
	}
}

func TestMappingObservationRefreshPreservesNativeSampleIdentity(t *testing.T) {
	f := newMappingObservationFixture(t, true)
	// Refresh preserves a native legacy sample; its producer still needs a
	// current installation without changing the native workspace scope.
	f.service.scheduler.AddDevice(f.projection.Devices[0])
	require.NotZero(t, f.service.scheduler.AddPointWithObservationLease(f.projection.Points[0]))
	require.NoError(t, f.service.RefreshMappings(t.Context()))
	require.Empty(t, f.service.mappingIndex["P"][0].WorkspaceID, "diagnostics must not rewrite existing native sample scope")
	o := f.execute(t, uint16(10))
	require.Equal(t, MappingObservationSucceeded, o.Result)
	require.Equal(t, "W", o.WorkspaceID)
	require.Equal(t, "WR1", o.WorkspaceRevision)
	require.Empty(t, f.sink.samples, "native missing-workspace typed rejection remains unchanged")
	require.Len(t, f.writer.records, 1)
	require.EqualValues(t, 1, f.service.Snapshot().WriteError)
	require.Equal(t, 2, f.reader.reads)
}

func TestMappingObservationPanicDoesNotChangeNativeOutcome(t *testing.T) {
	f := newMappingObservationFixture(t, true)
	f.install(t)
	f.service.mappingObserver = func(MappingObservation) { panic("private secret") }
	require.NotPanics(t, func() {
		f.service.handleCollectedValue(t.Context(), f.value(t, uint16(65535)))
		f.service.handleCollectedValue(t.Context(), f.value(t, uint16(10)))
	})
	require.EqualValues(t, 1, f.service.Snapshot().MappingError)
	require.EqualValues(t, 1, f.service.Snapshot().WriteSuccess)
	require.Len(t, f.writer.records, 1)
	require.Len(t, f.sink.samples, 2)
}

func TestMappingObservationNoPipelineDoesNotInventExecution(t *testing.T) {
	f := newMappingObservationFixture(t, true)
	f.projection.Mappings[0].TransformPipeline = "[]"
	f.install(t)
	f.service.handleCollectedValue(t.Context(), f.value(t, uint16(10)))
	require.Empty(t, *f.seen)
	require.Len(t, f.writer.records, 1)
}

func TestMappingObservationCarrierCannotContainPayloadOrRawErrors(t *testing.T) {
	f := newMappingObservationFixture(t, true)
	f.projection.Mappings[0].TransformPipeline = `[{"type":"password_fixture","params":{"password":"secret_params","host":"private.invalid"}}]`
	f.install(t)
	o := f.execute(t, "secret_sample")
	require.Equal(t, MappingObservationFailed, o.Result)
	raw, err := json.Marshal(o)
	require.NoError(t, err)
	for _, secret := range []string{"password_fixture", "secret_params", "private.invalid", "secret_sample", "TransformPipeline", "ConfigFingerprint"} {
		require.NotContains(t, string(raw), secret)
	}
}

type mappingObservationStringer struct{ during func() }

func (v mappingObservationStringer) String() string { v.during(); return "trigger" }

func TestMappingObservationLateExecutedR1DoesNotBorrowR2(t *testing.T) {
	f := newMappingObservationFixture(t, true)
	f.projection.Mappings[0].TransformPipeline = `[{"type":"lookup","params":{"table":{"trigger":10}}},{"type":"cast","params":{"target_type":"int16"}}]`
	f.install(t)
	old := f.service.mappingIndex["P"][0]
	cv := f.value(t, mappingObservationStringer{during: func() {
		next := *f.projection.Mappings[0]
		next.TransformPipeline = `[{"type":"cast","params":{"target_type":"int16"}},{"type":"scale","params":{"multiplier":2}}]`
		f.projection.Mappings = []*schema.Mapping{&next}
		f.reader.record.DatabaseSetupRevision = "WR2"
		f.install(t)
	}})
	f.service.handleCollectedValue(t.Context(), cv)
	require.Len(t, *f.seen, 1)
	require.Equal(t, MappingObservationSucceeded, (*f.seen)[0].Result)
	require.Equal(t, "WR1", (*f.seen)[0].WorkspaceRevision)
	require.Equal(t, old.MappingRevision, (*f.seen)[0].MappingRevision)
	require.NotEqual(t, old.MappingRevision, f.service.mappingIndex["P"][0].MappingRevision)
	newResult := f.execute(t, uint16(10))
	require.Equal(t, "WR2", newResult.WorkspaceRevision)
	require.NotEqual(t, old.MappingRevision, newResult.MappingRevision)
	require.Equal(t, 4, f.reader.reads, "only two install brackets, never per-result reads")
}

func TestMappingObservationScopedInstallRetainsUnselectedProof(t *testing.T) {
	f := newMappingObservationFixture(t, true)
	p := *f.projection.Points[0]
	p.ID = "P2"
	p.Address = "40002"
	tg := *f.projection.Tags[0]
	tg.ID = "T2"
	m := *f.projection.Mappings[0]
	m.ID = "M2"
	m.PointID = p.ID
	m.TagID = tg.ID
	f.projection.Points = append(f.projection.Points, &p)
	f.projection.Tags = append(f.projection.Tags, &tg)
	f.projection.Mappings = append(f.projection.Mappings, &m)
	f.install(t)
	f.reader.record.DatabaseSetupRevision = "WR2"
	f.projection.Mappings[0].TransformPipeline = `[{"type":"scale","params":{"multiplier":2}}]`
	require.NoError(t, f.service.ApplyWorkspaceProjectionForDevicesAndPoints(t.Context(), f.projection, []string{"D"}, []string{"P"}))
	require.Equal(t, "WR2", f.execute(t, uint16(10)).WorkspaceRevision)
	cv := typedRuntimeValue(f.service.mappingIndex["P2"][0])
	cv.Value = uint16(10)
	f.service.handleCollectedValue(t.Context(), cv)
	require.Equal(t, "WR1", (*f.seen)[1].WorkspaceRevision)
	require.Equal(t, "M2", (*f.seen)[1].MappingID)
}

func TestMappingObservationRefreshMissingProofNeverBorrowsPreviousScope(t *testing.T) {
	f := newMappingObservationFixture(t, true)
	f.install(t)
	f.reader.record = nil
	require.NoError(t, f.service.RefreshMappings(t.Context()))
	require.Equal(t, MappingObservation{Result: MappingObservationUnavailable}, f.execute(t, uint16(10)))
	require.Empty(t, f.service.mappingIndex["P"][0].WorkspaceID)
}
