package main

import (
	"context"
	"errors"
	"io/fs"
	"maps"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"go-gateway/internal/apphost"
	"go-gateway/internal/datalink"
	"go-gateway/internal/desktop"
)

var databaseAliasKeys = []string{"GATEWAY_DB_PATH", "DB_PATH", "SQLITE_PATH"}

// forgetEnvironment removes variables for one test and restores them afterwards,
// so dotenv loading by the code under test cannot leak into later tests.
func forgetEnvironment(t *testing.T, keys ...string) {
	t.Helper()
	for _, key := range keys {
		t.Setenv(key, "")
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
	}
}

func canonicalTempDir(t *testing.T) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func writeSelectionFixture(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

// snapshotTree records every path and file body of a disposable directory.
func snapshotTree(t *testing.T, root string) map[string]string {
	t.Helper()
	state := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if entry.IsDir() {
			state[rel] = "<dir>"
			return nil
		}
		body, err := os.ReadFile(path)
		state[rel] = string(body)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return state
}

func directoryNames(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names
}

// desktopDataFixture isolates a desktop launch: the real test-binary directory is
// the executable root and a disposable directory is the shortcut working directory.
func desktopDataFixture(t *testing.T) (executableDir, cwd string) {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	if exe, err = filepath.EvalSymlinks(exe); err != nil {
		t.Fatal(err)
	}
	executableDir = filepath.Dir(exe)
	for _, name := range []string{".env", "datalink.db"} {
		if _, err = os.Lstat(filepath.Join(executableDir, name)); err == nil {
			t.Skipf("executable directory already holds %s; the default-location branch cannot be isolated", name)
		}
	}
	forgetEnvironment(t, databaseAliasKeys...)
	cwd = canonicalTempDir(t)
	t.Chdir(cwd)
	return executableDir, cwd
}

func newSelectionHost(t *testing.T, mode apphost.Mode) *gatewayHost {
	t.Helper()
	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)
	host := newGatewayHost(ctx, apphost.Launch{Mode: mode}, cancel)
	host.newShell = func(desktop.Options) desktop.Shell {
		t.Error("native shell was created before the data source was settled")
		return &failingTray{}
	}
	return host
}

func assertNoStartupResources(t *testing.T, h *gatewayHost) {
	t.Helper()
	if h.cfg != nil || h.owner != nil || h.shell != nil || h.sink != nil || h.logs != nil || h.listener != nil || h.server != nil ||
		h.db != nil || h.connections != nil || h.runtime != nil || h.pipeline != nil || h.share != nil || h.trayReady {
		t.Fatal("startup opened a resource before the data source was settled")
	}
	if value := os.Getenv("GATEWAY_DB_PATH"); value != "" {
		t.Fatalf("an unsettled data source was published to the process: %q", value)
	}
}

func assertNoStorage(t *testing.T, database string) {
	t.Helper()
	for _, name := range []string{database, database + "-wal", database + "-shm", database + "-journal"} {
		if _, err := os.Lstat(name); !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("storage exists after a launch that must not create it: %s (%v)", name, err)
		}
	}
}

func skipWhenConfirmationIsNative(t *testing.T) {
	t.Helper()
	if runtime.GOOS == goosWindows {
		t.Skip("NOT RUN: the new-database confirmation is a native Windows dialog; cancel is verified through the non-Windows boundary that always declines")
	}
}

func TestDesktopNewDatabaseCancelLeavesStorageAndServicesUntouched(t *testing.T) {
	skipWhenConfirmationIsNative(t)
	executableDir, cwd := desktopDataFixture(t)
	executableBefore, cwdBefore := directoryNames(t, executableDir), snapshotTree(t, cwd)
	host := newSelectionHost(t, apphost.Desktop)
	err := prepareHost(host)
	if !errors.Is(err, errLaunchHandled) {
		t.Fatalf("cancel was not reported as a handled launch: %v", err)
	}
	want := filepath.Join(executableDir, "datalink.db")
	if !host.data.ConfirmCreate || host.data.Explicit || host.data.DatabasePath != want {
		t.Fatalf("the confirmation branch was not reached for the displayed location: %+v", host.data)
	}
	assertNoStartupResources(t, host)
	assertNoStorage(t, want)
	if after := directoryNames(t, executableDir); !slices.Equal(executableBefore, after) {
		t.Fatalf("cancel changed the executable directory: %v -> %v", executableBefore, after)
	}
	if after := snapshotTree(t, cwd); !maps.Equal(cwdBefore, after) {
		t.Fatalf("cancel changed the working directory: %v -> %v", cwdBefore, after)
	}
	if err = host.close(); err != nil {
		t.Fatal(err)
	}
}

func TestDesktopNewDatabaseCancelExitsSuccessfully(t *testing.T) {
	skipWhenConfirmationIsNative(t)
	executableDir, cwd := desktopDataFixture(t)
	executableBefore, cwdBefore := directoryNames(t, executableDir), snapshotTree(t, cwd)
	if err := runGateway(apphost.Launch{Mode: apphost.Desktop}); err != nil {
		t.Fatalf("declining to create a database was reported as a failure: %v", err)
	}
	assertNoStorage(t, filepath.Join(executableDir, "datalink.db"))
	if after := directoryNames(t, executableDir); !slices.Equal(executableBefore, after) {
		t.Fatalf("cancel changed the executable directory: %v -> %v", executableBefore, after)
	}
	if after := snapshotTree(t, cwd); !maps.Equal(cwdBefore, after) {
		t.Fatalf("cancel changed the working directory: %v -> %v", cwdBefore, after)
	}
}

func TestDesktopAmbiguousWorkingDirectoryStopsBeforeAnyResource(t *testing.T) {
	for name, files := range map[string]map[string]string{
		"legacy database": {"datalink.db": "legacy"},
		"legacy dotenv":   {".env": "DB_PATH=legacy.db\n"},
	} {
		t.Run(name, func(t *testing.T) {
			executableDir, cwd := desktopDataFixture(t)
			for file, body := range files {
				writeSelectionFixture(t, filepath.Join(cwd, file), body)
			}
			executableBefore, cwdBefore := directoryNames(t, executableDir), snapshotTree(t, cwd)
			host := newSelectionHost(t, apphost.Desktop)
			err := prepareHost(host)
			if !apphost.IsCode(err, "startup.path_ambiguous") {
				t.Fatalf("ambiguous launch was not stopped: %v", err)
			}
			if !strings.Contains(nativeNextStep(apphost.Code(err)), "GATEWAY_DB_PATH") {
				t.Fatal("operator is not told to configure an absolute database path")
			}
			assertNoStartupResources(t, host)
			assertNoStorage(t, filepath.Join(executableDir, "datalink.db"))
			if after := directoryNames(t, executableDir); !slices.Equal(executableBefore, after) {
				t.Fatalf("ambiguity changed the executable directory: %v -> %v", executableBefore, after)
			}
			if after := snapshotTree(t, cwd); !maps.Equal(cwdBefore, after) {
				t.Fatalf("ambiguity moved or rewrote legacy data: %v -> %v", cwdBefore, after)
			}
			if err = host.close(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDesktopExplicitPathFailureDoesNotFallBackToDefaultDatabase(t *testing.T) {
	executableDir, cwd := desktopDataFixture(t)
	configured := filepath.Join(canonicalTempDir(t), "missing-parent", "data.db")
	t.Setenv("GATEWAY_DB_PATH", configured)
	executableBefore, cwdBefore := directoryNames(t, executableDir), snapshotTree(t, cwd)
	host := newSelectionHost(t, apphost.Desktop)
	err := prepareHost(host)
	if !apphost.IsCode(err, "startup.path_unusable") {
		t.Fatalf("unusable explicit path was not reported: %v", err)
	}
	if got := os.Getenv("GATEWAY_DB_PATH"); got != configured {
		t.Fatalf("explicit path was replaced after failure: %q", got)
	}
	if host.cfg != nil || host.owner != nil || host.sink != nil || host.listener != nil || host.trayReady {
		t.Fatal("startup continued after an unusable explicit path")
	}
	assertNoStorage(t, filepath.Join(executableDir, "datalink.db"))
	assertNoStorage(t, filepath.Dir(configured))
	if after := directoryNames(t, executableDir); !slices.Equal(executableBefore, after) {
		t.Fatalf("failure changed the executable directory: %v -> %v", executableBefore, after)
	}
	if after := snapshotTree(t, cwd); !maps.Equal(cwdBefore, after) {
		t.Fatalf("failure changed the working directory: %v -> %v", cwdBefore, after)
	}
	if err = host.close(); err != nil {
		t.Fatal(err)
	}
}

func TestCLIModesSelectDatabaseFromWorkingDirectoryDotenv(t *testing.T) {
	if runtime.GOOS == goosWindows {
		t.Skip("NOT RUN: the Windows owner guard creates its lock target before bind; CLI semantics are verified on Linux/macOS")
	}
	for _, mode := range []apphost.Mode{apphost.Console, apphost.Headless} {
		t.Run(string(mode), func(t *testing.T) {
			occupied, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = occupied.Close() })
			_, port, err := net.SplitHostPort(occupied.Addr().String())
			if err != nil {
				t.Fatal(err)
			}
			base := canonicalTempDir(t)
			cwd, external := filepath.Join(base, "cwd"), filepath.Join(base, "external")
			for _, dir := range []string{cwd, external} {
				if err = os.Mkdir(dir, 0o700); err != nil {
					t.Fatal(err)
				}
			}
			selected := filepath.Join(external, "from-dotenv.db")
			forgetEnvironment(t, append([]string{"HOST", "PORT", "LOG_OUTPUT", "LOG_LEVEL", "LOG_FILE", "GIN_MODE"}, databaseAliasKeys...)...)
			writeSelectionFixture(t, filepath.Join(cwd, ".env"), "GATEWAY_DB_PATH="+selected+"\nHOST=127.0.0.1\nPORT="+port+"\nLOG_OUTPUT=console\n")
			writeSelectionFixture(t, filepath.Join(cwd, "datalink.db"), "legacy")
			t.Chdir(cwd)
			before := snapshotTree(t, base)
			host := newSelectionHost(t, mode)
			err = prepareHost(host)
			if !apphost.IsCode(err, "startup.bind_failed") {
				t.Fatalf("occupied port was not the stopping point: %v", err)
			}
			if !host.data.Explicit || host.data.ConfirmCreate || host.data.DatabasePath != selected || host.data.DataRoot != external {
				t.Fatalf("CLI ignored the working-directory dotenv database: %+v", host.data)
			}
			assertNoStorage(t, selected)
			if after := snapshotTree(t, base); !maps.Equal(before, after) {
				t.Fatalf("CLI launch changed storage before bind: %v -> %v", before, after)
			}
			if err = host.close(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDatabaseAliasOrderAgreesBetweenResolverAndHarnessDSN(t *testing.T) {
	dir := canonicalTempDir(t)
	forgetEnvironment(t, databaseAliasKeys...)
	t.Chdir(dir)
	for _, name := range []string{"first.db", "second.db", "third.db", "datalink.db"} {
		writeSelectionFixture(t, filepath.Join(dir, name), name)
	}
	for _, tc := range []struct {
		name   string
		values map[string]string
		want   string
	}{
		{"all aliases", map[string]string{"GATEWAY_DB_PATH": "first.db", "DB_PATH": "second.db", "SQLITE_PATH": "third.db"}, "first.db"},
		{"blank top alias falls through", map[string]string{"GATEWAY_DB_PATH": "  ", "DB_PATH": "second.db", "SQLITE_PATH": "third.db"}, "second.db"},
		{"SQLITE_PATH only", map[string]string{"SQLITE_PATH": "third.db"}, "third.db"},
		{"no alias keeps the embedded default", nil, "datalink.db"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for key, name := range tc.values {
				value := name
				if strings.TrimSpace(name) != "" {
					value = filepath.Join(dir, name)
				}
				t.Setenv(key, value)
			}
			selection, err := apphost.ResolveData(apphost.DataInput{Mode: apphost.Console, Executable: filepath.Join(dir, "gateway"), WorkingDir: dir, Environment: environmentValues()})
			if err != nil || selection.DatabasePath != filepath.Join(dir, tc.want) {
				t.Fatalf("resolver selected another database: %+v %v", selection, err)
			}
			if tc.values == nil && embeddedSQLiteDSN() != datalink.DefaultEmbeddedSQLiteDSN {
				t.Fatalf("harness left the embedded default: %q", embeddedSQLiteDSN())
			}
			parsed, err := url.Parse(embeddedSQLiteDSN())
			if err != nil {
				t.Fatal(err)
			}
			dsnPath := filepath.FromSlash(parsed.Path)
			if parsed.Path == "" {
				dsnPath = filepath.Join(dir, parsed.Opaque)
			} else if runtime.GOOS == goosWindows {
				dsnPath = strings.TrimPrefix(dsnPath, string(filepath.Separator))
			}
			resolved, resolvedErr := os.Stat(selection.DatabasePath)
			opened, openedErr := os.Stat(dsnPath)
			if resolvedErr != nil || openedErr != nil || !os.SameFile(resolved, opened) {
				t.Fatalf("resolver and harness DSN name different databases: %s vs %q (%v, %v)", selection.DatabasePath, embeddedSQLiteDSN(), resolvedErr, openedErr)
			}
		})
	}
}
