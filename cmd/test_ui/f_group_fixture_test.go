//go:build f_write_group_fixture

package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
)

func fixtureOwnedDatabasePath(t *testing.T, marker string) string {
	t.Helper()
	runDir := filepath.Join(t.TempDir(), "gw-f-quality-123")
	if err := os.MkdirAll(runDir, 0o700); err != nil {
		t.Fatalf("mkdir fixture run: %v", err)
	}
	if err := os.WriteFile(filepath.Join(runDir, ".f-write-group-fixture"), []byte(marker), 0o600); err != nil {
		t.Fatalf("write fixture marker: %v", err)
	}
	databasePath := filepath.Join(runDir, "gateway.db")
	t.Setenv("GATEWAY_DB_PATH", databasePath)
	return databasePath
}

func TestFixtureClockRejectsNonUTCAndAcceptsNanoUTC(t *testing.T) {
	if _, err := parseFixtureTime("2026-01-01T00:00:00+08:00", ""); err == nil {
		t.Fatal("clock accepted a non-UTC offset")
	}
	got, err := parseFixtureTime("2026-01-01T00:00:00.123456789Z", "")
	if err != nil {
		t.Fatalf("parse UTC clock: %v", err)
	}
	if want := time.Date(2026, 1, 1, 0, 0, 0, 123456789, time.UTC); !got.Equal(want) {
		t.Fatalf("clock = %s, want %s", got, want)
	}
}

func TestFixtureInitialClockRejectsInvalidEnvironment(t *testing.T) {
	for name, value := range map[string]string{
		"malformed": "not-a-clock",
		"non-UTC":   "2026-01-01T00:00:00+08:00",
	} {
		t.Run(name, func(t *testing.T) {
			t.Setenv("F_FIXTURE_START_AT", value)
			t.Setenv("F_FIXTURE_PORT", "")
			hooks := newGroupFixture()
			fixture, ok := hooks.(*groupFixture)
			if !ok {
				t.Fatalf("fixture type = %T, want tagged controller", hooks)
			}
			if fixture.initErr == nil {
				t.Fatalf("invalid F_FIXTURE_START_AT %q was accepted", value)
			}
		})
	}

	t.Run("default", func(t *testing.T) {
		t.Setenv("F_FIXTURE_START_AT", "")
		t.Setenv("F_FIXTURE_PORT", "")
		hooks := newGroupFixture()
		fixture, ok := hooks.(*groupFixture)
		if !ok {
			t.Fatalf("fixture type = %T, want tagged controller", hooks)
		}
		if fixture.initErr != nil {
			t.Fatalf("default fixture clock was rejected: %v", fixture.initErr)
		}
	})
}

func TestFixtureCaptureIsNotAckedAndRejectsUnknownIndex(t *testing.T) {
	capture := newFixtureCapture(1)
	sample := measurement.SampleEnvelope{SampleID: "sample-1", AcquisitionID: "acq-1"}
	if err := capture.AcceptSample(t.Context(), sample); err != nil {
		t.Fatalf("capture sample: %v", err)
	}
	if got := capture.samples[0].status; got != "captured_not_acked" {
		t.Fatalf("capture status = %q, want captured_not_acked", got)
	}
	if _, err := capture.sample(1); err == nil {
		t.Fatal("unknown capture index was accepted")
	}
	if err := capture.AcceptSample(t.Context(), sample); err == nil {
		t.Fatal("bounded capture queue accepted overflow")
	}
}

func TestFixtureReleaseValidatesEveryIndexBeforeAcceptingAny(t *testing.T) {
	capture := newFixtureCapture(4)
	if err := capture.AcceptSample(t.Context(), measurement.SampleEnvelope{SampleID: "sample-1"}); err != nil {
		t.Fatalf("capture sample: %v", err)
	}
	fixture := &groupFixture{
		capture: capture,
		pipe:    grouppipeline.New(grouppipeline.Dependencies{}, grouppipeline.Config{}),
	}
	if _, err := fixture.release(t.Context(), []int{0, 99}); err == nil {
		t.Fatal("release accepted a batch containing an unknown index")
	}
	if got := capture.samples[0].attempts; got != 0 {
		t.Fatalf("valid capture attempts = %d after invalid batch, want 0", got)
	}
}

func TestFixturePointSubsetRejectsForeignAndDuplicateIDs(t *testing.T) {
	allowed := map[string]struct{}{"point-a": {}}
	if _, err := fixturePointSubset([]string{"foreign"}, allowed); err == nil {
		t.Fatal("foreign point was accepted")
	}
	if _, err := fixturePointSubset([]string{"point-a", "point-a"}, allowed); err == nil {
		t.Fatal("duplicate point was accepted")
	}
}

func TestFixturePollResolvesCanonicalWriteGroupID(t *testing.T) {
	env := newOutageEnv(t)
	group := env.createAndApply(t, "A")
	fixture, ok := env.services.fixture.(*groupFixture)
	if !ok {
		t.Fatalf("fixture type = %T, want tagged controller", env.services.fixture)
	}
	fixture.deps = &env.services
	fixture.pipe = env.services.groupPipe
	fixture.paused.Store(true)
	if _, err := env.db.ExecContext(t.Context(), `UPDATE points SET description = '' WHERE id = ?`, group.Members[0].PointID); err != nil {
		t.Fatalf("prepare canonical point fixture: %v", err)
	}
	_, err := fixture.poll(t.Context(), fixturePollRequest{
		GroupID:  group.ID,
		PointIDs: []string{group.Members[0].PointID},
	})
	if err != nil {
		t.Fatalf("canonical write-group poll: %v", err)
	}
}

func TestFixtureTickFailsWhenPipelineReportsClosureError(t *testing.T) {
	env := newOutageEnv(t)
	group := env.createAndApply(t, "A")
	fixture, ok := env.services.fixture.(*groupFixture)
	if !ok {
		t.Fatalf("fixture type = %T, want tagged controller", env.services.fixture)
	}
	fixture.pipe = env.services.groupPipe
	resolved, err := env.services.writeGroups.ResolveAppliedAt(t.Context(), group.ID, time.Now().UTC().Add(time.Minute))
	if err != nil {
		t.Fatalf("resolve applied group: %v", err)
	}
	base := resolved.EffectiveAt.Truncate(10 * time.Second).Add(10 * time.Second)
	fixture.clock.Set(base.Add(time.Second))
	if err := fixture.pipe.Reconcile(t.Context()); err != nil {
		t.Fatalf("reconcile fixture group: %v", err)
	}
	observed := base.Add(5 * time.Second)
	fixture.clock.Set(observed.Add(time.Second))
	if err := fixture.pipe.AcceptSample(t.Context(), env.envelope(group, 0, "fixture-tick-error", observed, 21.5)); err != nil {
		t.Fatalf("accept fixture sample: %v", err)
	}
	fixture.clock.Set(base.Add(11 * time.Second))
	if _, err := env.db.ExecContext(t.Context(), `CREATE TRIGGER fixture_fail_closure BEFORE INSERT ON wg_delivery_buckets BEGIN SELECT RAISE(ABORT, 'fixture closure failure'); END`); err != nil {
		t.Fatalf("install closure failure trigger: %v", err)
	}
	t.Cleanup(func() { _, _ = env.db.ExecContext(context.Background(), `DROP TRIGGER fixture_fail_closure`) })
	before := fixture.errors.Load()
	if err := fixture.tick(t.Context()); err == nil {
		t.Fatal("tick reported success after pipeline closure error")
	}
	if got := fixture.errors.Load(); got <= before {
		t.Fatalf("pipeline errors = %d, want greater than %d", got, before)
	}
}

func TestFixtureDatabasePathMustBeDisposableTempFile(t *testing.T) {
	t.Setenv("GATEWAY_DB_PATH", "/etc/gateway.db")
	if err := preflightGroupFixtureDatabase(); err == nil {
		t.Fatal("non-temporary database path was accepted")
	}
}

func TestFixtureDatabasePreflightRequiresOwnedMarker(t *testing.T) {
	fixtureOwnedDatabasePath(t, "wrong-marker")
	if err := preflightGroupFixtureDatabase(); err == nil {
		t.Fatal("fixture path without the exact ownership marker was accepted")
	}
}

func TestFixtureDatabasePreflightAcceptsOwnedUncreatedFile(t *testing.T) {
	databasePath := fixtureOwnedDatabasePath(t, "f_write_group_fixture")
	if err := preflightGroupFixtureDatabase(); err != nil {
		t.Fatalf("owned uncreated database was rejected: %v", err)
	}
	if _, err := os.Stat(databasePath); !os.IsNotExist(err) {
		t.Fatalf("preflight created or changed database file: stat err = %v", err)
	}
}

func TestFixtureDatabasePreflightRejectsDatabaseSymlink(t *testing.T) {
	databasePath := fixtureOwnedDatabasePath(t, "f_write_group_fixture")
	if err := os.WriteFile(databasePath+".real", []byte("owned"), 0o600); err != nil {
		t.Fatalf("write symlink target: %v", err)
	}
	if err := os.Symlink(databasePath+".real", databasePath); err != nil {
		t.Fatalf("create database symlink: %v", err)
	}
	if err := preflightGroupFixtureDatabase(); err == nil {
		t.Fatal("database symlink was accepted")
	}
}

func TestFixtureDatabasePreflightRejectsRunDirectoryEscape(t *testing.T) {
	targetDir := filepath.Join(t.TempDir(), "fixture-target")
	if err := os.MkdirAll(targetDir, 0o700); err != nil {
		t.Fatalf("create disposable symlink target: %v", err)
	}
	linkParent := t.TempDir()
	link := filepath.Join(linkParent, "gw-f-escape")
	if err := os.Symlink(targetDir, link); err != nil {
		t.Fatalf("create run-directory symlink: %v", err)
	}
	t.Setenv("GATEWAY_DB_PATH", filepath.Join(link, "gateway.db"))
	if err := preflightGroupFixtureDatabase(); err == nil {
		t.Fatal("run-directory symlink escape was accepted")
	}
}
