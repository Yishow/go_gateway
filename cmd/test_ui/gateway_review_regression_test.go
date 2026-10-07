package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/apphost"
	"go-gateway/internal/desktop"
	"go-gateway/internal/diagnostics"
)

type orderedTray struct {
	once             sync.Once
	mu               sync.Mutex
	states           []desktop.State
	entered, release chan struct{}
}

func (*orderedTray) Start() error { return nil }
func (*orderedTray) Close() error { return nil }
func (s *orderedTray) Update(state desktop.State) {
	if state.Process == "serving" {
		s.once.Do(func() { close(s.entered); <-s.release })
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states = append(s.states, state)
}
func TestQuitCannotBeOverwrittenByInflightServingPublication(t *testing.T) {
	tray := &orderedTray{entered: make(chan struct{}), release: make(chan struct{})}
	host := &gatewayHost{quit: make(chan struct{}), shell: tray}
	servingDone := make(chan struct{})
	go func() { host.updateState("serving"); close(servingDone) }()
	<-tray.entered
	quitDone := make(chan struct{})
	go func() { host.requestQuit(); close(quitDone) }()
	<-host.quit
	close(tray.release)
	<-servingDone
	<-quitDone
	host.updateState("serving")
	tray.mu.Lock()
	defer tray.mu.Unlock()
	last := tray.states[len(tray.states)-1]
	if last.Process != "stopping" || last.SetupURL != "" || last.LogsURL != "" {
		t.Fatalf("late serving publication: %+v", last)
	}
	if err := host.openSetup(); err == nil {
		t.Fatal("opened during shutdown")
	}
}
func TestRunCommandFileModeFailureDoesNotMirrorRawStderr(t *testing.T) {
	ln, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	t.Setenv("GATEWAY_DB_PATH", filepath.Join(root, "chosen.db"))
	t.Setenv("LOG_OUTPUT", "file")
	t.Setenv("LOG_FILE", filepath.Join(root, "logs", "runtime.jsonl"))
	t.Setenv("HOST", "127.0.0.1")
	t.Setenv("PORT", port)
	var output strings.Builder
	old := log.Writer()
	log.SetOutput(&output)
	defer log.SetOutput(old)
	if got := runCommand([]string{"--headless"}); got != 1 {
		t.Fatalf("exit=%d", got)
	}
	if output.Len() != 0 {
		t.Fatalf("file-only failure mirrored: %q", output.String())
	}
}
func TestOpenSetupRefusesStoppingCoordinator(t *testing.T) {
	host := &gatewayHost{state: desktop.State{SetupURL: "http://127.0.0.1:1/studio/v2"}}
	host.shutdown = apphost.NewShutdown(time.Second, nil, nil)
	host.shutdown.Begin()
	if err := host.shutdown.Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := host.openSetup(); err == nil {
		t.Fatal("opened after shutdown")
	}
}

func TestQuitImmediatelyClosesRequestAdmissionAndWaitsInflight(t *testing.T) {
	host := &gatewayHost{quit: make(chan struct{})}
	entered, release := make(chan struct{}), make(chan struct{})
	gate := &requestGate{stopping: host.quit, next: http.HandlerFunc(func(http.ResponseWriter, *http.Request) { close(entered); <-release })}
	requestDone := make(chan struct{})
	go func() {
		gate.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/api/v1/settings", http.NoBody))
		close(requestDone)
	}()
	<-entered
	host.requestQuit()
	rejected := httptest.NewRecorder()
	gate.ServeHTTP(rejected, httptest.NewRequest("POST", "/api/v1/settings", http.NoBody))
	if rejected.Code != 503 {
		t.Fatal("admitted after quit")
	}
	drained := make(chan struct{})
	go func() { gate.wait(); close(drained) }()
	select {
	case <-drained:
		t.Fatal("closed with live handler")
	default:
	}
	close(release)
	<-requestDone
	<-drained
}

func TestControlledNativeFailureHasSafeEventAndNonSuccess(t *testing.T) {
	broker, err := diagnostics.New(diagnostics.Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer broker.Close(t.Context())
	host := &gatewayHost{logs: broker}
	host.recordShellFailure()
	if !apphost.IsCode(host.shellFailure(), "startup.tray_failed") {
		t.Fatal("native failure discarded")
	}
	if finishGatewayResult(errLaunchHandled, host.shellFailure()) == nil {
		t.Fatal("false successful exit")
	}
	if err = broker.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	snapshot, err := broker.Snapshot(diagnostics.Query{})
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Records) != 1 || snapshot.Records[0].Code != "startup.tray_failed" {
		t.Fatalf("missing safe event: %+v", snapshot)
	}
}

func TestBrowserWaitCancellationDoesNotAdmitSecondBlockedOpener(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	host := &gatewayHost{startupCtx: ctx}
	entered, release := make(chan struct{}), make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- host.openBrowser("fixed", func(string) error { close(entered); <-release; return nil })
	}()
	<-entered
	cancel()
	if err := <-done; err == nil {
		t.Fatal("browser blocked shutdown")
	}
	if err := host.openBrowser("fixed", func(string) error { t.Error("second opener admitted"); return nil }); err == nil {
		t.Fatal("missing busy result")
	}
	close(release)
}
func TestServingResourcePublicationDoesNotWaitForBrowserOrReturn(t *testing.T) {
	if _, err := staticFiles.ReadFile("static/index.html"); err != nil {
		t.Skip("production HTTP probe test requires the built embedded SPA")
	}
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	host := newGatewayHost(ctx, apphost.Launch{Mode: apphost.Headless}, cancel)
	ln, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	host.listener = ln
	entered := make(chan struct{})
	returned := make(chan error, 1)
	go func() {
		returned <- host.serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { close(entered); <-r.Context().Done() }))
	}()
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("own HTTP readiness probe did not start")
	}
	select {
	case <-host.startupDone:
	default:
		t.Fatal("published HTTP resources still considered initializing")
	}
	select {
	case err := <-returned:
		t.Fatalf("readiness unexpectedly returned: %v", err)
	default:
	}
	host.requestQuit()
	<-returned
	if err = host.close(); err != nil {
		t.Fatal(err)
	}
}

func TestPostShutdownProducerRetainsFileOnlyPolicy(t *testing.T) {
	var output strings.Builder
	old := log.Writer()
	log.SetOutput(&output)
	defer log.SetOutput(old)
	host := &gatewayHost{}
	host.installLogging(false)
	host.restoreLogging()
	host.mirrorConsole("fixture-private-readiness")
	if output.Len() != 0 {
		t.Fatal("late producer bypassed file-only policy")
	}
}

func TestLateQuitDoesNotReplaceActualPendingShutdownPhase(t *testing.T) {
	host := &gatewayHost{quit: make(chan struct{})}
	entered, release := make(chan struct{}), make(chan struct{})
	host.shutdown = apphost.NewShutdown(time.Second, []apphost.Phase{{Name: "database", Stop: func(context.Context) error { close(entered); <-release; return nil }}}, func(state apphost.ShutdownState) { host.updateShutdownState(state.Process, state.Phase) })
	host.shutdown.Begin()
	<-entered
	host.requestQuit()
	host.stateMu.Lock()
	phase := host.state.PendingPhase
	host.stateMu.Unlock()
	close(release)
	if err := host.shutdown.Wait(t.Context()); err != nil {
		t.Fatal(err)
	}
	if phase != "database" {
		t.Fatalf("overwrote actual pending phase with %q", phase)
	}
}
