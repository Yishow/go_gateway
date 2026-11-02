package collector

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionGeometryBatchPublicationRejectsWholeSetBeforeAssignment(t *testing.T) {
	for _, kind := range []string{"current", "second-reinstalled", "duplicate", "second-finalized", "busy"} {
		t.Run(kind, func(t *testing.T) {
			s := NewScheduler(DefaultSchedulerConfig(), nil)
			publications := make([]AcquisitionGeometryPublication, 0, 2)
			tokens := make([]AcquisitionGeometryMutation, 0, 2)
			for _, id := range []string{"a", "b"} {
				d := schema.Device{ID: id, Protocol: schema.ProtocolMC3E, ConnectionConfig: `{}`}
				s.AddDevice(&d)
				s.AddPoint(&schema.Point{ID: id, DeviceID: id, Address: "D0", DataType: schema.DataTypeUint16})
				token := s.BeginAcquisitionGeometryMutation(id)
				tokens = append(tokens, token)
				publications = append(publications, s.PrepareAcquisitionGeometry(token, d, AcquisitionGeometryDefinition{Revision: id}, "").Seal())
			}
			before := s.CaptureAcquisitionGeometry("a").publication
			switch kind {
			case "second-reinstalled":
				s.AddPoint(&schema.Point{ID: "b", DeviceID: "b", Address: "D0", DataType: schema.DataTypeUint16})
			case "duplicate":
				publications[1] = publications[0]
			case "second-finalized":
				s.FinishAcquisitionGeometryMutation(tokens[1])
			case "busy":
				s.mu.Lock()
			}
			receipts := make([]AcquisitionGeometryReceipt, 2)
			calls := 0
			accepted := s.PublishAcquisitionGeometries(publications, receipts, func() { calls++ })
			if kind == "busy" {
				s.mu.Unlock()
			}
			require.Equal(t, kind == "current", accepted)
			if accepted {
				require.Equal(t, 1, calls)
				require.NotNil(t, receipts[0].publication)
				require.NotNil(t, receipts[1].publication)
			} else {
				require.Zero(t, calls)
				require.Zero(t, receipts[0])
				require.Zero(t, receipts[1])
				require.Same(t, before, s.CaptureAcquisitionGeometry("a").publication)
			}
		})
	}
}

func TestAcquisitionEmptyPublicationDoesNotEraseForeignDevice(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	called := false
	require.True(t, s.PublishEmptyAcquisitionProjection(func() { called = true }))
	require.True(t, called)
	s.AddDevice(&schema.Device{ID: "foreign"})
	called = false
	require.False(t, s.PublishEmptyAcquisitionProjection(func() { called = true }))
	require.False(t, called)
}
