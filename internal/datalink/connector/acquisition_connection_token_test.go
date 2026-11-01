package connector

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/wirebudget"

	"github.com/stretchr/testify/require"
)

type tokenProtocol struct {
	stubProtocol
	generation atomic.Uint64
	read       func(context.Context, PlannedReadRequest) ([]PlannedPointResult, error)
}

func (p *tokenProtocol) Connect(context.Context, string) error {
	p.connected = true
	p.generation.Add(1)
	return nil
}
func (p *tokenProtocol) AcquisitionTransportGeneration() uint64 { return p.generation.Load() }
func (p *tokenProtocol) ReadPlanned(ctx context.Context, request PlannedReadRequest) ([]PlannedPointResult, error) {
	return p.read(ctx, request)
}

func tokenRequest(t *testing.T) PlannedReadRequest {
	t.Helper()
	p := schema.Point{ID: "point", DeviceID: "device", RevisionID: "point-rev", Address: "D0", DataType: schema.DataTypeInt16}
	profile, err := ResolveReadProfile(p, decoding.FormatContext{Protocol: schema.ProtocolMC3E, ConnectionRevision: ReadConnectionRevision("{}")})
	require.NoError(t, err)
	r := profile.Resolution()
	ep := optimizer.EffectivePoint{ID: p.ID, Address: p.Address, DataType: p.DataType, NativeWidth: 1, Format: r.EffectiveFormat, FormatSource: string(r.Source), FormatRevision: p.RevisionID}
	plan := optimizer.BuildAcquisitionPlan(optimizer.AcquisitionInput{Sources: []optimizer.AcquisitionSource{{Identity: optimizer.AcquisitionIdentity{DeviceID: p.DeviceID, ConnectionRevision: ReadConnectionRevision("{}"), Protocol: "mc", Frame: "3e-binary", Cadence: time.Second}, Capability: optimizer.AcquisitionCapability{ProtocolLimit: 125}, Points: []optimizer.EffectivePoint{ep}}}})
	require.Empty(t, plan.Errors())
	return PlannedReadRequest{Plan: plan, AcquisitionID: "scan", RequestID: "request", Bindings: []PlannedPointBinding{{Point: ep, Profile: profile, PointRevision: p.RevisionID}}}
}

func TestAcquisitionTokenCapturesFinalActualAttempt(t *testing.T) {
	proto := &tokenProtocol{}
	old, exists := GetFactory(schema.ProtocolMC3E)
	Register(schema.ProtocolMC3E, func() Protocol { return proto })
	t.Cleanup(func() {
		if exists {
			Register(schema.ProtocolMC3E, old)
		} else {
			Unregister(schema.ProtocolMC3E)
		}
	})
	manager := NewConnectionManager(DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	conn, err := manager.GetOrCreate(t.Context(), "device", schema.ProtocolMC3E, "{}")
	require.NoError(t, err)
	initial, err := conn.CurrentAcquisitionToken(t.Context())
	require.NoError(t, err)
	require.True(t, initial.Valid())
	proto.read = func(ctx context.Context, req PlannedReadRequest) ([]PlannedPointResult, error) {
		require.NoError(t, wirebudget.Take(ctx))
		proto.generation.Add(1) // The adapter successfully reconnected before retry.
		require.NoError(t, wirebudget.Take(ctx))
		final := proto.generation.Load()
		proto.generation.Add(1) // A post-send update must not relabel the final wire attempt.
		return []PlannedPointResult{{Binding: req.Bindings[0], ReceivedAt: time.Now(), ReadResult: ReadResult{Value: final, Quality: schema.QualityGood}}}, nil
	}
	ctx := wirebudget.With(t.Context(), 2, 0)
	results, err := conn.ReadPlanned(ctx, tokenRequest(t))
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Equal(t, initial.ManagerGeneration, results[0].ConnectionToken.ManagerGeneration)
	require.Equal(t, initial.TransportGeneration+1, results[0].ConnectionToken.TransportGeneration)
	require.Equal(t, 2, wirebudget.From(ctx).Used())
}

func TestAcquisitionTokenNoSendCannotClaimGeneration(t *testing.T) {
	proto := &tokenProtocol{}
	proto.generation.Store(9)
	denied := errors.New("stale installation")
	proto.read = func(ctx context.Context, req PlannedReadRequest) ([]PlannedPointResult, error) {
		err := wirebudget.Take(ctx)
		return []PlannedPointResult{{Binding: req.Bindings[0]}}, err
	}
	conn := &ManagedConnection{Protocol: proto, ProtocolType: schema.ProtocolMC3E, DeviceID: "device", Config: "{}"}
	ctx := wirebudget.WithAdmission(t.Context(), func() error { return denied })
	got, err := conn.ReadPlanned(ctx, tokenRequest(t))
	require.ErrorIs(t, err, denied)
	require.Zero(t, got[0].ConnectionToken)
	current, err := conn.CurrentAcquisitionToken(t.Context())
	require.NoError(t, err)
	require.False(t, current.Valid(), "unmanaged legacy fixtures do not invent a manager incarnation")
}

func TestAcquisitionTokenGetterUsesCancelableManagedGate(t *testing.T) {
	conn := &ManagedConnection{}
	conn.Lock()
	defer conn.Unlock()
	ctx, cancel := context.WithTimeout(t.Context(), time.Millisecond)
	defer cancel()
	_, err := conn.CurrentAcquisitionToken(ctx)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

type tokenDialProtocol struct {
	tokenProtocol
	dial func(context.Context, string) error
}

func (p *tokenDialProtocol) Connect(ctx context.Context, config string) error {
	if p.dial != nil {
		if err := p.dial(ctx, config); err != nil {
			return err
		}
	}
	return p.tokenProtocol.Connect(ctx, config)
}

func TestAcquisitionTokenManagerNeverPublishesLateDial(t *testing.T) {
	typ := schema.ProtocolType(t.Name())
	entered, release := make(chan struct{}), make(chan struct{})
	registerGeneration(t, typ, func() Protocol {
		return &tokenDialProtocol{dial: func(_ context.Context, config string) error {
			if config == `{"host":"old"}` {
				close(entered)
				<-release
			}
			return nil
		}}
	})
	manager := NewConnectionManager(DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	old := acquire(t.Context(), manager, "device", typ, `{"host":"old"}`)
	<-entered
	manager.mu.RLock()
	oldEntry := manager.entries["device"]
	manager.mu.RUnlock()
	newer := result(t, acquire(t.Context(), manager, "device", typ, `{"host":"new"}`))
	require.NoError(t, newer.err)
	newToken, err := newer.conn.CurrentAcquisitionToken(t.Context())
	require.NoError(t, err)
	require.True(t, newToken.Valid())
	close(release)
	require.Error(t, result(t, old).err)
	require.NotNil(t, oldEntry.conn)
	oldToken, err := oldEntry.conn.CurrentAcquisitionToken(t.Context())
	require.NoError(t, err)
	require.Zero(t, oldToken.ManagerGeneration, "failed/late entries never publish a manager token")
	current, err := newer.conn.CurrentAcquisitionToken(t.Context())
	require.NoError(t, err)
	require.Equal(t, newToken, current)
}
