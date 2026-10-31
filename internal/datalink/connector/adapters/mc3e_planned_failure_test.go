package adapters

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/mcprotocol"
	"go-gateway/internal/protocol/wirebudget"
)

func TestMCPlannedMalformedResponsesNeverGoodOrFanout(t *testing.T) {
	for _, kind := range []string{"contiguous", "random", "bits"} {
		t.Run(kind, func(t *testing.T) {
			points := []schema.Point{mcPlannedPoint("a", "D0", schema.DataTypeUint32, "ABCD"), mcPlannedPoint("b", "D0", schema.DataTypeUint32, "CDAB")}
			if kind == "random" {
				points = append(points, mcPlannedPoint("c", "D10", schema.DataTypeInt16, ""))
			}
			if kind == "bits" {
				points = []schema.Point{mcPlannedPoint("a", "M0", schema.DataTypeBool, ""), mcPlannedPoint("b", "M1", schema.DataTypeBool, "")}
			}
			req := mcPlannedRequest(t, points...)
			normal := mcPlannedWords(1, 2)
			if kind == "random" {
				normal = mcPlannedWords(1, 2, 3)
			}
			if kind == "bits" {
				normal = []byte{0, 0, 0x10}
			}
			cases := [][]byte{nil, {0}, {0, 0}, normal[:len(normal)-1], append(append([]byte(nil), normal...), 0), {0x59, 0xc0}, {0x51, 0xc0}}
			if kind == "bits" {
				cases = append(cases, []byte{0, 0, 0x12}, []byte{0, 0, 0x20})
			}
			for _, response := range cases {
				tr := &mcPlannedTransport{response: response}
				ctx := wirebudget.With(t.Context(), 4, 0)
				results, err := mcPlannedConnector(tr).ReadPlanned(ctx, req)
				require.Error(t, err, "response %x", response)
				require.Len(t, results, len(points))
				require.Len(t, tr.frames, 1, "no malformed/unsupported scalar fallback")
				require.Equal(t, 1, wirebudget.From(ctx).Used())
				for _, r := range results {
					require.Equal(t, schema.QualityBad, r.Quality)
					require.Nil(t, r.Value)
					require.Empty(t, r.RawBytes)
					require.NotEmpty(t, r.Error)
					require.Equal(t, req.Plan.Fingerprint(), r.PlanFingerprint)
					require.NotEmpty(t, r.Binding.ConfigFingerprint)
				}
				if len(response) == 2 && response[0] == 0x59 {
					require.True(t, mcprotocol.IsCommandRejected(err))
				} else {
					require.False(t, mcprotocol.IsCommandRejected(err))
				}
			}
		})
	}
}
func TestMCPlannedAttemptsAndCancelCannotMintNewSends(t *testing.T) {
	req := mcPlannedRequest(t, mcPlannedPoint("a", "D0", schema.DataTypeInt16, ""), mcPlannedPoint("b", "D10", schema.DataTypeInt16, ""))
	for _, limit := range []int{1, 2, 4} {
		tr := &mcPlannedTransport{err: io.EOF}
		ctx := wirebudget.With(t.Context(), limit, 0)
		results, err := mcPlannedConnector(tr).ReadPlanned(ctx, req)
		require.ErrorIs(t, err, io.EOF)
		require.Len(t, tr.frames, limit)
		require.Equal(t, limit, wirebudget.From(ctx).Used())
		for _, r := range results {
			require.Equal(t, schema.QualityBad, r.Quality)
			require.Nil(t, r.Value)
		}
	}
	for _, before := range []bool{true, false} {
		ctx, cancel := context.WithCancel(t.Context())
		tr := &mcPlannedTransport{response: mcPlannedWords(1, 2), hook: cancel}
		if before {
			cancel()
		}
		results, err := mcPlannedConnector(tr).ReadPlanned(ctx, req)
		require.ErrorIs(t, err, context.Canceled)
		if before {
			require.Empty(t, tr.frames)
		} else {
			require.Len(t, tr.frames, 1)
		}
		for _, r := range results {
			require.Equal(t, schema.QualityBad, r.Quality)
			require.Nil(t, r.Value)
		}
	}
	tr := &mcPlannedTransport{err: io.EOF}
	ctx, cancel := context.WithTimeout(wirebudget.With(t.Context(), 4, time.Second), 20*time.Millisecond)
	defer cancel()
	_, err := mcPlannedConnector(tr).ReadPlanned(ctx, req)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Len(t, tr.frames, 1)
}
func TestMCPlannedNoCrossScanCacheAndFinalCompletion(t *testing.T) {
	req := mcPlannedRequest(t, mcPlannedPoint("a", "D0", schema.DataTypeInt16, ""))
	tr := &mcPlannedTransport{response: mcPlannedWords(1)}
	var completedAfter time.Time
	tr.hook = func() { completedAfter = time.Now() }
	c := mcPlannedConnector(tr)
	first, err := c.ReadPlanned(t.Context(), req)
	require.NoError(t, err)
	require.False(t, first[0].Timestamp.Before(completedAfter))
	require.Equal(t, first[0].Timestamp, first[0].ObservedAt)
	require.Equal(t, first[0].ObservedAt, first[0].ReceivedAt)
	tr.response = mcPlannedWords(2)
	req.AcquisitionID = "scan-2"
	req.RequestID = "request-2"
	second, err := c.ReadPlanned(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, int16(1), first[0].Value)
	require.Equal(t, int16(2), second[0].Value)
	require.Equal(t, "scan-2", second[0].AcquisitionID)
	tr.err = errors.New("failed response")
	third, err := c.ReadPlanned(t.Context(), req)
	require.Error(t, err)
	require.Nil(t, third[0].Value)
	require.Equal(t, schema.QualityBad, third[0].Quality)
	require.Len(t, tr.frames, 3)
}
func TestMCPlannedRejectsMismatchedSnapshotBeforeSend(t *testing.T) {
	for _, change := range []string{"point", "address", "format", "revision", "profile", "missing", "duplicate", "index", "request-id", "plan"} {
		t.Run(change, func(t *testing.T) {
			req := mcPlannedRequest(t, mcPlannedPoint("a", "D0", schema.DataTypeUint32, "ABCD"), mcPlannedPoint("b", "D0", schema.DataTypeUint32, "CDAB"))
			switch change {
			case "point":
				req.Bindings[0].Point.ID = "foreign"
			case "address":
				req.Bindings[0].Point.Address = "D10"
			case "format":
				req.Bindings[0].Point.Format = "DCBA"
			case "revision":
				req.Bindings[0].PointRevision = "foreign"
			case "profile":
				req.Bindings[0].Profile = req.Bindings[1].Profile
			case "missing":
				req.Bindings = req.Bindings[:1]
			case "duplicate":
				req.Bindings[1] = req.Bindings[0]
			case "index":
				req.RequestIndex = 5
			case "request-id":
				req.RequestID = ""
			case "plan":
				req.Plan = connector.PlannedReadRequest{}.Plan
			}
			tr := &mcPlannedTransport{response: mcPlannedWords(1, 2)}
			results, err := mcPlannedConnector(tr).ReadPlanned(t.Context(), req)
			require.Error(t, err)
			require.Empty(t, tr.frames)
			for _, r := range results {
				require.Equal(t, schema.QualityBad, r.Quality)
				require.Nil(t, r.Value)
			}
		})
	}
}
