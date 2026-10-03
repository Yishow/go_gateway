//go:build f_write_group_fixture

package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	datalinkruntime "go-gateway/internal/datalink/runtime"
)

const (
	fixtureDefaultStart    = "2025-12-31T23:59:50Z"
	fixtureDefaultPort     = 3355
	fixtureCaptureCapacity = 2048
	fixtureMarkerName      = ".f-write-group-fixture"
	fixtureMarkerValue     = "f_write_group_fixture"
)

type groupFixtureHooks interface {
	pipelineConfig(nodeID string, ownerID string) grouppipeline.Config
	sampleSink(*grouppipeline.Pipeline) datalinkruntime.SampleSink
	configure(context.Context, *sql.DB, *gatewayServices) (func() error, error)
}

type fixtureClock struct{ nanos atomic.Int64 }

func newFixtureClock(value time.Time) *fixtureClock {
	clock := &fixtureClock{}
	clock.nanos.Store(value.UTC().UnixNano())
	return clock
}

func (c *fixtureClock) Now() time.Time {
	return time.Unix(0, c.nanos.Load()).UTC()
}

func (c *fixtureClock) Set(value time.Time) { c.nanos.Store(value.UTC().UnixNano()) }

type capturedFixtureSample struct {
	envelope   measurement.SampleEnvelope
	status     string
	lastReason string
	attempts   int
}

type fixtureCapture struct {
	mu       sync.Mutex
	pipe     *grouppipeline.Pipeline
	capacity int
	samples  []capturedFixtureSample
}

func newFixtureCapture(capacity int) *fixtureCapture {
	return &fixtureCapture{capacity: capacity, samples: make([]capturedFixtureSample, 0, capacity)}
}

func (c *fixtureCapture) WantsSample(deviceID, pointID, tagID string) bool {
	c.mu.Lock()
	pipe := c.pipe
	c.mu.Unlock()
	return pipe != nil && pipe.WantsSample(deviceID, pointID, tagID)
}

func (c *fixtureCapture) AcceptSample(ctx context.Context, envelope measurement.SampleEnvelope) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.capacity <= 0 || len(c.samples) >= c.capacity {
		datalinkruntime.ReportFixtureSampleError(ctx)
		return errors.New("fixture capture queue full")
	}
	envelope.RawValue = cloneFixtureValue(envelope.RawValue)
	envelope.Value = cloneFixtureValue(envelope.Value)
	c.samples = append(c.samples, capturedFixtureSample{envelope: envelope, status: "captured_not_acked"})
	return nil
}

func cloneFixtureValue(value any) any {
	bytes, ok := value.([]byte)
	if !ok {
		return value
	}
	return slices.Clone(bytes)
}

func (c *fixtureCapture) length() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.samples)
}

func (c *fixtureCapture) hasCapacity() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.capacity > 0 && len(c.samples) < c.capacity
}

func (c *fixtureCapture) sample(index int) (measurement.SampleEnvelope, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if index < 0 || index >= len(c.samples) {
		return measurement.SampleEnvelope{}, fmt.Errorf("unknown capture index")
	}
	return c.samples[index].envelope, nil
}

type groupFixture struct {
	clock    *fixtureClock
	capture  *fixtureCapture
	pipe     *grouppipeline.Pipeline
	initErr  error
	port     int
	httpMu   sync.Mutex
	http     *fixtureHTTPServer
	deps     *gatewayServices
	db       *sql.DB
	faults   *fixtureFaultController
	capacity *fixtureCapacityController
	paused   atomic.Bool
	nextID   atomic.Uint64
	forcedID atomic.Value
	errors   atomic.Uint64
	opMu     sync.Mutex
}

func newGroupFixture() groupFixtureHooks {
	start, err := parseFixtureTime(os.Getenv("F_FIXTURE_START_AT"), fixtureDefaultStart)
	if err != nil {
		fallbackStart, fallbackErr := parseFixtureTime(fixtureDefaultStart, fixtureDefaultStart)
		if fallbackErr == nil {
			start = fallbackStart
		} else {
			err = errors.Join(err, fallbackErr)
		}
	}
	port, portErr := parseFixturePort(os.Getenv("F_FIXTURE_PORT"))
	if err == nil {
		err = portErr
	}
	fixture := &groupFixture{
		clock:   newFixtureClock(start),
		capture: newFixtureCapture(fixtureCaptureCapacity),
		port:    port,
		initErr: err,
	}
	fixture.forcedID.Store("")
	return fixture
}

func (f *groupFixture) pipelineConfig(nodeID, ownerID string) grouppipeline.Config {
	nodeID = fixtureWorkerNodeID(nodeID)
	return grouppipeline.Config{
		NodeID:            nodeID,
		Owner:             nodeID + "/" + ownerID,
		TickInterval:      24 * time.Hour,
		ReconcileInterval: 24 * time.Hour,
		Sender:            groupDeliverySenderConfig(),
		Now:               f.clock.Now,
		OnError:           f.recordError,
	}
}

func (f *groupFixture) sampleSink(pipe *grouppipeline.Pipeline) datalinkruntime.SampleSink {
	f.pipe = pipe
	f.capture.mu.Lock()
	f.capture.pipe = pipe
	f.capture.mu.Unlock()
	return f.capture
}

func (f *groupFixture) configure(ctx context.Context, db *sql.DB, services *gatewayServices) (func() error, error) {
	if f.initErr != nil {
		return nil, f.initErr
	}
	if err := preflightGroupFixtureDatabase(); err != nil {
		return nil, err
	}
	if db == nil || services == nil || services.scheduler == nil || services.runtime == nil || services.groupPipe == nil || services.writeGroups == nil {
		return nil, errors.New("fixture production services are incomplete")
	}
	f.db, f.deps = db, services
	f.faults = newFixtureFaultController(db, services, f.clock.Now)
	f.capacity = newFixtureCapacityController(db)
	if err := f.faults.install(); err != nil {
		return nil, err
	}
	services.scheduler.WithClock(f.clock.Now).WithAcquisitionIDFactory(f.nextAcquisitionID)
	services.writeGroups.WithClock(f.clock.Now)
	if err := f.startHTTP(ctx); err != nil {
		if closeErr := f.faults.close(ctx); closeErr != nil {
			f.recordError(closeErr)
		}
		return nil, err
	}
	return func() error {
		return errors.Join(
			f.stopHTTP(ctx),
			f.faults.close(ctx),
			f.capacity.close(ctx),
		)
	}, nil
}

func (f *groupFixture) nextAcquisitionID() string {
	forced, ok := f.forcedID.Load().(string)
	if !ok {
		forced = ""
	}
	if forced != "" {
		return forced
	}
	return fmt.Sprintf("fixture-acq-%d", f.nextID.Add(1))
}

func (f *groupFixture) setForcedAcquisitionID(value string) {
	f.forcedID.Store(value)
}

func (f *groupFixture) recordError(error) { f.errors.Add(1) }

func parseFixturePort(raw string) (int, error) {
	if strings.TrimSpace(raw) == "" {
		return fixtureDefaultPort, nil
	}
	port, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || port < 0 || port > 65535 {
		return 0, errors.New("fixture port must be between 0 and 65535")
	}
	return port, nil
}

func parseFixtureTime(raw, fallback string) (time.Time, error) {
	if strings.TrimSpace(raw) == "" {
		raw = fallback
	}
	value, err := time.Parse(time.RFC3339Nano, raw)
	if err != nil || value.Location() != time.UTC || !strings.HasSuffix(raw, "Z") {
		return time.Time{}, errors.New("fixture clock must be RFC3339Nano UTC")
	}
	return value.UTC(), nil
}

func preflightGroupFixtureDatabase() error {
	if err := preflightFixtureCapacity(); err != nil {
		return err
	}
	raw := strings.TrimSpace(os.Getenv("GATEWAY_DB_PATH"))
	if raw == "" || strings.HasPrefix(raw, "file:") || !filepath.IsAbs(raw) {
		return errors.New("f_write_group_fixture requires a disposable GATEWAY_DB_PATH")
	}
	absolute, err := filepath.Abs(raw)
	if err != nil {
		return errors.New("fixture database path is invalid")
	}
	runDir := filepath.Dir(absolute)
	runInfo, err := os.Lstat(runDir) //nolint:gosec // runDir is checked as a marker-owned temporary directory below.
	if err != nil || !runInfo.IsDir() || runInfo.Mode()&os.ModeSymlink != 0 {
		return errors.New("fixture database parent must be an owned directory")
	}
	if !strings.HasPrefix(filepath.Base(runDir), "gw-f-") {
		return errors.New("fixture database parent is not an owned fixture run")
	}
	canonicalRunDir, err := filepath.EvalSymlinks(runDir)
	if err != nil || !isFixtureTempPath(canonicalRunDir) {
		return errors.New("fixture database parent escaped the temporary directory")
	}
	markerPath := filepath.Join(runDir, fixtureMarkerName)
	markerInfo, err := os.Lstat(markerPath) //nolint:gosec // markerPath is inside the validated runDir.
	if err != nil || markerInfo.Mode()&os.ModeSymlink != 0 || !markerInfo.Mode().IsRegular() {
		return errors.New("fixture ownership marker is missing")
	}
	marker, err := os.ReadFile(markerPath) //nolint:gosec // markerPath is inside the validated runDir.
	if err != nil || string(marker) != fixtureMarkerValue {
		return errors.New("fixture ownership marker is invalid")
	}
	databaseInfo, statErr := os.Lstat(absolute) //nolint:gosec // absolute is inside the validated marker-owned runDir.
	if statErr == nil {
		if databaseInfo.Mode()&os.ModeSymlink != 0 || !databaseInfo.Mode().IsRegular() {
			return errors.New("fixture database path must be a regular non-symlink file")
		}
	} else if !errors.Is(statErr, os.ErrNotExist) {
		return errors.New("fixture database path cannot be inspected")
	}
	if err := fixtureClosureRegistrationError(); err != nil {
		return errors.New("fixture closure function registration failed")
	}
	return nil
}

func isFixtureTempPath(path string) bool {
	temporary := false
	for _, root := range []string{os.TempDir(), "/tmp", "/private/tmp"} {
		tempRoot, rootErr := filepath.Abs(root)
		if rootErr != nil {
			continue
		}
		canonicalRoot, evalErr := filepath.EvalSymlinks(tempRoot)
		if evalErr != nil {
			continue
		}
		relative, relErr := filepath.Rel(filepath.Clean(canonicalRoot), filepath.Clean(path))
		if relErr == nil && relative != "." && relative != ".." && !strings.HasPrefix(relative, ".."+string(os.PathSeparator)) {
			temporary = true
			break
		}
	}
	return temporary
}

var _ groupFixtureHooks = (*groupFixture)(nil)
var _ datalinkruntime.SampleSink = (*fixtureCapture)(nil)
var _ datalinkruntime.SampleInterest = (*fixtureCapture)(nil)
