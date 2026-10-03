//go:build f_write_group_fixture

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	datalinkruntime "go-gateway/internal/datalink/runtime"

	"github.com/stretchr/testify/require"
)

func TestFixtureCapacityEndpointAppliesOwnedPageLimit(t *testing.T) {
	env := newOutageEnv(t)
	fixture := env.services.fixture.(*groupFixture)
	fixture.db = env.db
	fixture.capacity = newFixtureCapacityController(env.db)
	fixture.port = 0
	require.NoError(t, fixture.startHTTP(t.Context()))
	t.Cleanup(func() { fixture.stopHTTP(context.Background()) })
	var originalPageCount, originalMaxPageCount int64
	require.NoError(t, env.db.QueryRowContext(t.Context(), "PRAGMA page_count").Scan(&originalPageCount))
	require.NoError(t, env.db.QueryRowContext(t.Context(), "PRAGMA max_page_count").Scan(&originalMaxPageCount))

	postCapacity := func(enabled bool) fixtureCapacityState {
		t.Helper()
		body, err := json.Marshal(fixtureCapacityRequest{Kind: fixtureCapacityKindDiskFull, Enabled: enabled})
		require.NoError(t, err)
		request, err := http.NewRequestWithContext(t.Context(), http.MethodPost,
			"http://"+fixture.http.listener.Addr().String()+"/capacity", bytes.NewReader(body))
		require.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		response, err := (&http.Client{}).Do(request)
		require.NoError(t, err)
		defer response.Body.Close()
		require.Equal(t, http.StatusOK, response.StatusCode)
		var state fixtureCapacityState
		require.NoError(t, json.NewDecoder(response.Body).Decode(&state))
		return state
	}

	enabled := postCapacity(true)
	require.Equal(t, fixtureCapacityKindDiskFull, enabled.Kind)
	require.True(t, enabled.Enabled)
	require.Equal(t, originalPageCount, enabled.PageCount)
	require.Equal(t, enabled.PageCount, enabled.MaxPageCount)
	disabled := postCapacity(false)
	require.Equal(t, fixtureCapacityKindDiskFull, disabled.Kind)
	require.False(t, disabled.Enabled)
	require.Equal(t, originalMaxPageCount, disabled.MaxPageCount)
	var restoredMaxPageCount int64
	require.NoError(t, env.db.QueryRowContext(t.Context(), "PRAGMA max_page_count").Scan(&restoredMaxPageCount))
	require.Equal(t, originalMaxPageCount, restoredMaxPageCount)
}

func TestFixtureQuotaConfigAcceptsBoundedPositiveOverrides(t *testing.T) {
	t.Setenv(fixtureGroupQuotaEnv, "123")
	t.Setenv(fixtureGlobalQuotaEnv, "456")
	config, err := fixtureQuotaConfig()
	require.NoError(t, err)
	require.Equal(t, int64(123), config.GroupMaxBytes)
	require.Equal(t, int64(456), config.GlobalMaxBytes)
}

func TestFixtureQuotaPreflightRejectsInvalidOverrides(t *testing.T) {
	for name, value := range map[string]string{
		"zero":  "0",
		"minus": "-1",
		"text":  "not-a-number",
		"large": strconv.FormatInt(fixtureQuotaMaxBytes+1, 10),
	} {
		t.Run(name, func(t *testing.T) {
			fixtureOwnedDatabasePath(t, fixtureMarkerValue)
			t.Setenv(fixtureGroupQuotaEnv, value)
			if err := preflightGroupFixtureDatabase(); err == nil {
				t.Fatalf("invalid %s quota %q was accepted", name, value)
			}
		})
	}
}

func TestFixtureCapacityDiskFullRefusesNewProductionSampleAndKeepsAcceptedData(t *testing.T) {
	env := newOutageEnv(t)
	fixture := env.services.fixture.(*groupFixture)
	activeAt := time.Now().UTC()
	fixture.clock.Set(activeAt)
	env.services.writeGroups.WithClock(fixture.clock.Now)
	group := env.createAndApply(t, "A")
	fixture.pipe = env.services.groupPipe
	fixture.capacity = newFixtureCapacityController(env.db)

	resolved, err := env.services.writeGroups.ResolveAppliedAt(t.Context(), group.ID, activeAt.Add(time.Minute))
	require.NoError(t, err)
	base := resolved.EffectiveAt.Truncate(10 * time.Second).Add(10 * time.Second)
	fixture.clock.Set(base.Add(time.Second))
	require.NoError(t, fixture.pipe.Reconcile(t.Context()))
	firstObserved := base.Add(5 * time.Second)
	fixture.clock.Set(firstObserved.Add(time.Second))
	require.NoError(t, fixture.pipe.AcceptSample(t.Context(), env.envelope(group, 0, "capacity-before-full", firstObserved, 21.5)))
	var before int
	require.NoError(t, env.db.QueryRowContext(t.Context(),
		`SELECT COUNT(*) FROM wg_delivery_samples WHERE group_id = ?`, group.ID).Scan(&before))
	require.Positive(t, before)

	state, err := fixture.capacity.configure(t.Context(), true)
	require.NoError(t, err)
	require.True(t, state.Enabled)
	t.Cleanup(func() { require.NoError(t, fixture.capacity.close(context.Background())) })

	var refused error
	for i := range 512 {
		observed := base.Add(time.Duration((i+1)*10+5) * time.Second)
		fixture.clock.Set(observed.Add(time.Second))
		refused = fixture.pipe.AcceptSample(t.Context(), env.envelope(group, 0,
			fmt.Sprintf("capacity-after-full-%d", i), observed, float64(i)))
		if refused != nil {
			break
		}
	}
	var sampleError *datalinkruntime.GroupSampleError
	require.ErrorAs(t, refused, &sampleError)
	require.Equal(t, "disk-full", sampleError.Reason)
	var after int
	require.NoError(t, env.db.QueryRowContext(t.Context(),
		`SELECT COUNT(*) FROM wg_delivery_samples WHERE group_id = ?`, group.ID).Scan(&after))
	require.GreaterOrEqual(t, after, before, "disk-full refusal must not remove accepted samples")
	var stored int
	require.NoError(t, env.db.QueryRowContext(t.Context(),
		`SELECT COUNT(*) FROM wg_delivery_samples WHERE sample_id = ?`, "capacity-before-full").Scan(&stored))
	require.Equal(t, 1, stored)
}

func TestFixtureGroupQuotaOverrideRefusesNewProductionSample(t *testing.T) {
	t.Setenv(fixtureGroupQuotaEnv, "1")
	env := newOutageEnv(t)
	fixture := env.services.fixture.(*groupFixture)
	activeAt := time.Now().UTC()
	fixture.clock.Set(activeAt)
	env.services.writeGroups.WithClock(fixture.clock.Now)
	group := env.createAndApply(t, "A")
	fixture.pipe = env.services.groupPipe

	resolved, err := env.services.writeGroups.ResolveAppliedAt(t.Context(), group.ID, activeAt.Add(time.Minute))
	require.NoError(t, err)
	base := resolved.EffectiveAt.Truncate(10 * time.Second).Add(10 * time.Second)
	fixture.clock.Set(base.Add(time.Second))
	require.NoError(t, fixture.pipe.Reconcile(t.Context()))
	firstObserved := base.Add(5 * time.Second)
	fixture.clock.Set(firstObserved.Add(time.Second))
	require.NoError(t, fixture.pipe.AcceptSample(t.Context(), env.envelope(group, 0, "quota-before-limit", firstObserved, 21.5)))

	secondObserved := base.Add(15 * time.Second)
	fixture.clock.Set(secondObserved.Add(time.Second))
	err = fixture.pipe.AcceptSample(t.Context(), env.envelope(group, 0, "quota-at-limit", secondObserved, 22.5))
	var sampleError *datalinkruntime.GroupSampleError
	require.ErrorAs(t, err, &sampleError)
	require.Equal(t, "quota-hard-limit", sampleError.Reason)
}
