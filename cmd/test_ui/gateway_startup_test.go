package main

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"go-gateway/internal/apphost"
	"go-gateway/internal/desktop"
	"go-gateway/internal/diagnostics"
)

func TestProductionBindFailsBeforeDatabaseInitialization(t *testing.T) {
	ln, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "unused.db")
	t.Setenv("GATEWAY_DB_PATH", path)
	t.Setenv("HOST", "127.0.0.1")
	t.Setenv("PORT", port)
	t.Setenv("LOG_OUTPUT", "console")
	t.Setenv("AUTO_OPEN_BROWSER", "true")
	err = runGateway(apphost.Launch{Mode: apphost.Headless})
	if !apphost.IsCode(err, "startup.bind_failed") {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if runtime.GOOS == goosWindows {
		if err != nil || info.Size() != 0 {
			t.Fatalf("DB initialized before bind: %v", err)
		}
	} else if !os.IsNotExist(err) {
		t.Fatalf("DB created before bind: %v", err)
	}
}
func TestSelectedLiteralDatabaseHasSameSQLiteIdentity(t *testing.T) {
	for _, name := range []string{"literal#other.db", "literal%23.db", "空白 space.db"} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), name)
			db, err := sql.Open("sqlite", embeddedSQLiteDSNForPath(path))
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			if _, err = db.ExecContext(t.Context(), "CREATE TABLE selected_identity(id INTEGER)"); err != nil {
				t.Fatal(err)
			}
			var seq int
			var label, actual string
			if err = db.QueryRowContext(t.Context(), "PRAGMA database_list").Scan(&seq, &label, &actual); err != nil {
				t.Fatal(err)
			}
			a, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}
			b, err := os.Stat(actual)
			if err != nil || !os.SameFile(a, b) {
				t.Fatalf("wrong identity: %s %v", actual, err)
			}
		})
	}
}
func TestProductionManagedWritersAndFileOnlyPolicy(t *testing.T) {
	for _, mirror := range []bool{false, true} {
		t.Run(map[bool]string{false: "file", true: "console"}[mirror], func(t *testing.T) {
			broker, err := diagnostics.New(diagnostics.Options{})
			if err != nil {
				t.Fatal(err)
			}
			defer broker.Close(t.Context())
			oldWriter, oldGin, oldRecovery, flags := log.Writer(), gin.DefaultWriter, gin.DefaultErrorWriter, log.Flags()
			var output strings.Builder
			log.SetOutput(&output)
			gin.DefaultWriter = &output
			gin.DefaultErrorWriter = &output
			defer func() {
				log.SetOutput(oldWriter)
				gin.DefaultWriter = oldGin
				gin.DefaultErrorWriter = oldRecovery
				log.SetFlags(flags)
			}()
			host := &gatewayHost{logs: broker}
			host.installLogging(mirror)
			defer host.restoreLogging()
			done := make(chan struct{})
			go func() {
				log.Print("fixture-private-standard")
				slog.Error("fixture-private-slog", "password", "fixture-password")
				_, _ = io.WriteString(gin.DefaultWriter, "fixture-private-access")
				_, _ = io.WriteString(gin.DefaultErrorWriter, "fixture-private-recovery")
				close(done)
			}()
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("logging deadlock")
			}
			if mirror && !strings.Contains(output.String(), "fixture-private") {
				t.Fatal("legacy mirror lost")
			}
			if !mirror && output.Len() != 0 {
				t.Fatal("raw file-mode mirror")
			}
			if err = broker.Flush(t.Context()); err != nil {
				t.Fatal(err)
			}
			snapshot, err := broker.Snapshot(diagnostics.Query{})
			if err != nil {
				t.Fatal(err)
			}
			for _, event := range snapshot.Records {
				if strings.Contains(event.Message, "fixture-") {
					t.Fatal("secret retained")
				}
			}
			host.restoreLogging()
			if log.Flags() != flags {
				t.Fatal("logger flags not restored")
			}
		})
	}
}
func TestStoppingGateAndLateStartupCannotAdmitNewWork(t *testing.T) {
	calls := 0
	gate := &requestGate{next: http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls++ })}
	gate.stop()
	recorder := httptest.NewRecorder()
	gate.ServeHTTP(recorder, httptest.NewRequest("POST", "/api/v1/settings", http.NoBody))
	if recorder.Code != 503 || calls != 0 {
		t.Fatal("mutation admitted")
	}
	ctx, cancel := context.WithCancel(t.Context())
	host := newGatewayHost(ctx, apphost.Launch{Mode: apphost.Headless}, cancel)
	host.requestQuit()
	if err := host.serve(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls++ })); err != nil {
		t.Fatal(err)
	}
	if calls != 0 {
		t.Fatal("startup canceled but served")
	}
	_ = host.close()
}
func TestShutdownTruthAndHandledExitPreserveFailure(t *testing.T) {
	host := &gatewayHost{state: desktop.State{Acquisition: "running"}}
	host.updateShutdownState("stop-timeout", "runtime")
	if host.state.Acquisition != "stopping" {
		t.Fatal("false running/stopped")
	}
	host.updateShutdownState("stop-timeout", "pipeline")
	if host.state.Acquisition != "stopped" {
		t.Fatal("runtime completion missing")
	}
	if err := finishGatewayResult(errLaunchHandled, apphost.ErrShutdownTimeout); !errors.Is(err, apphost.ErrShutdownTimeout) {
		t.Fatal("timeout hidden")
	}
	if (&gatewayHost{}).mustObserveDiagnosticCompletion() {
		t.Fatal("pre-service fatal lacks bounded exception")
	}
	if !(&gatewayHost{trayReady: true}).mustObserveDiagnosticCompletion() || !(&gatewayHost{db: &sql.DB{}}).mustObserveDiagnosticCompletion() {
		t.Fatal("live resources need actual completion")
	}
}

type failingTray struct{ closed bool }

func (*failingTray) Start() error         { return errors.New("fixture tray failure") }
func (*failingTray) Update(desktop.State) {}
func (t *failingTray) Close() error       { t.closed = true; return nil }
func TestInitialTrayFailureCleansUpBeforeAcquisition(t *testing.T) {
	if runtime.GOOS == goosWindows {
		t.Skip("native owner requires Windows interaction acceptance; Linux injection verifies orchestration")
	}
	root := t.TempDir()
	t.Setenv("GATEWAY_DB_PATH", filepath.Join(root, "selected.db"))
	t.Setenv("LOG_FILE", filepath.Join(root, "logs", "runtime.jsonl"))
	t.Setenv("LOG_OUTPUT", "file")
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	host := newGatewayHost(ctx, apphost.Launch{Mode: apphost.Desktop}, cancel)
	tray := &failingTray{}
	host.newShell = func(desktop.Options) desktop.Shell { return tray }
	err := prepareHost(host)
	if !apphost.IsCode(err, "startup.tray_failed") {
		t.Fatal(err)
	}
	if host.db != nil || host.runtime != nil || host.listener != nil || host.trayReady {
		t.Fatal("started beyond failed tray")
	}
	if err = host.close(); err != nil {
		t.Fatal(err)
	}
	if !tray.closed {
		t.Fatal("tray resources not closed")
	}
}

type unavailableRuntime struct{}

func (unavailableRuntime) Start(context.Context) error {
	return errors.New("fixture-private-driver-error")
}
func (unavailableRuntime) IsRunning() bool { return false }
func TestRuntimeFailureProjectsDegradedWithoutRawDetails(t *testing.T) {
	broker, err := diagnostics.New(diagnostics.Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer broker.Close(t.Context())
	host := &gatewayHost{logs: broker}
	host.startAcquisition(t.Context(), unavailableRuntime{})
	if err = broker.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	snapshot, err := broker.Snapshot(diagnostics.Query{})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, event := range snapshot.Records {
		if event.Code == "runtime.started" || strings.Contains(event.Message, "fixture-private") {
			t.Fatal("false or raw runtime claim")
		}
		found = found || event.Code == "runtime.degraded"
	}
	if !found {
		t.Fatal("missing degraded code")
	}
}
func TestInfoCommandsDoNotInitializeDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "unused.db")
	t.Setenv("GATEWAY_DB_PATH", path)
	for _, arg := range []string{"--help", "--version"} {
		if runCommand([]string{"--headless", arg}) != 0 {
			t.Fatal("info command failed")
		}
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("info command initialized DB")
	}
}
