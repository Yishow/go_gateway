package apphost

import (
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
)

// dataLayout is a disposable launcher layout: an executable directory, a
// different shortcut working directory and an external data directory.
type dataLayout struct{ root, exe, cwd, external string }

func newDataLayout(t *testing.T) dataLayout {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	l := dataLayout{root: root, exe: filepath.Join(root, "exe"), cwd: filepath.Join(root, "cwd"), external: filepath.Join(root, "external")}
	for _, dir := range []string{l.exe, l.cwd, l.external} {
		if err = os.Mkdir(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	return l
}

func (l dataLayout) input(mode Mode, env map[string]string) DataInput {
	return DataInput{Mode: mode, Executable: filepath.Join(l.exe, "gateway.exe"), WorkingDir: l.cwd, Environment: env}
}

func writeFixture(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

// treeState records every path and file body, so a rejected selection can prove
// that it created, moved and rewrote nothing.
func treeState(t *testing.T, root string) map[string]string {
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
		if entry.Type()&fs.ModeSymlink != 0 {
			target, linkErr := os.Readlink(path)
			state[rel] = "<link " + target + ">"
			return linkErr
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

func assertTreeUnchanged(t *testing.T, root string, before map[string]string) {
	t.Helper()
	if after := treeState(t, root); !maps.Equal(before, after) {
		t.Fatalf("selection changed storage:\nbefore=%v\nafter=%v", before, after)
	}
}

// aliasEntry turns "KEY=name" into an absolute external path; blank values stay blank.
func aliasEntry(entry, external string) (key, value string) {
	key, value, _ = strings.Cut(entry, "=")
	if strings.TrimSpace(value) != "" {
		value = filepath.Join(external, value)
	}
	return key, value
}

func TestDesktopWorkingDirectoryAmbiguityStopsBeforeStorage(t *testing.T) {
	for _, tc := range []struct {
		name          string
		cwd, exe, env map[string]string
	}{
		{name: "legacy database only", cwd: map[string]string{"datalink.db": "legacy"}},
		{name: "competing database beside an existing default", cwd: map[string]string{"datalink.db": "legacy"}, exe: map[string]string{"datalink.db": "default"}},
		{name: "competing dotenv beside an executable dotenv without a database alias", cwd: map[string]string{".env": "PORT=9000\n"}, exe: map[string]string{".env": "PORT=8000\n"}},
		{name: "relative alias is not an unambiguous choice", cwd: map[string]string{".env": "PORT=9000\n"}, env: map[string]string{"DB_PATH": "chosen.db"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			l := newDataLayout(t)
			for name, body := range tc.cwd {
				writeFixture(t, filepath.Join(l.cwd, name), body)
			}
			for name, body := range tc.exe {
				writeFixture(t, filepath.Join(l.exe, name), body)
			}
			before := treeState(t, l.root)
			got, err := ResolveData(l.input(Desktop, tc.env))
			if !IsCode(err, "startup.path_ambiguous") || got != (DataSelection{}) {
				t.Fatalf("ambiguous launch was guessed: %+v %v", got, err)
			}
			assertTreeUnchanged(t, l.root, before)
		})
	}
}

func TestDesktopWorkingDirectoryMatchingExecutableIsNotAmbiguous(t *testing.T) {
	for _, tc := range []struct {
		name     string
		alias    bool
		database bool
	}{
		{name: "same directory with existing database", database: true},
		{name: "symlink spelling with existing database", alias: true, database: true},
		{name: "symlink spelling still confirms a genuinely new database", alias: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			l := newDataLayout(t)
			writeFixture(t, filepath.Join(l.exe, ".env"), "PORT=8000\n")
			if tc.database {
				writeFixture(t, filepath.Join(l.exe, "datalink.db"), "existing")
			}
			in := l.input(Desktop, map[string]string{})
			in.WorkingDir = l.exe
			if tc.alias {
				in.WorkingDir = filepath.Join(l.root, "exe-link")
				if err := os.Symlink(l.exe, in.WorkingDir); err != nil {
					t.Skipf("directory symlinks are unavailable: %v", err)
				}
			}
			before := treeState(t, l.exe)
			got, err := ResolveData(in)
			want := DataSelection{DatabasePath: filepath.Join(l.exe, "datalink.db"), DataRoot: l.exe, ConfigRoot: l.exe, ConfirmCreate: !tc.database}
			if err != nil || got != want {
				t.Fatalf("launch from the executable directory was misjudged: %+v %v", got, err)
			}
			assertTreeUnchanged(t, l.exe, before)
		})
	}
}

func TestExplicitAbsoluteDatabaseWinsOverCompetingLocations(t *testing.T) {
	for _, tc := range []struct {
		name            string
		process, dotenv []string
		chosen          string
	}{
		{name: "process aliases keep their order", process: []string{"GATEWAY_DB_PATH=first.db", "DB_PATH=second.db", "SQLITE_PATH=third.db"}, chosen: "first.db"},
		{name: "DB_PATH precedes SQLITE_PATH", process: []string{"DB_PATH=second.db", "SQLITE_PATH=third.db"}, chosen: "second.db"},
		{name: "SQLITE_PATH alone is explicit", process: []string{"SQLITE_PATH=third.db"}, chosen: "third.db"},
		{name: "blank higher alias falls through", process: []string{"GATEWAY_DB_PATH=  ", "SQLITE_PATH=third.db"}, chosen: "third.db"},
		{name: "executable dotenv aliases keep their order", dotenv: []string{"SQLITE_PATH=third.db", "DB_PATH=second.db", "GATEWAY_DB_PATH=first.db"}, chosen: "first.db"},
		{name: "executable dotenv without the top alias", dotenv: []string{"SQLITE_PATH=third.db", "DB_PATH=second.db"}, chosen: "second.db"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			l := newDataLayout(t)
			writeFixture(t, filepath.Join(l.cwd, ".env"), "PORT=9000\n")
			writeFixture(t, filepath.Join(l.cwd, "datalink.db"), "legacy")
			writeFixture(t, filepath.Join(l.external, tc.chosen), "chosen")
			env := map[string]string{}
			for _, entry := range tc.process {
				key, value := aliasEntry(entry, l.external)
				env[key] = value
			}
			if len(tc.dotenv) > 0 {
				body := ""
				for _, entry := range tc.dotenv {
					key, value := aliasEntry(entry, l.external)
					body += key + "=" + value + "\n"
				}
				writeFixture(t, filepath.Join(l.exe, ".env"), body)
			}
			before := treeState(t, l.root)
			got, err := ResolveData(l.input(Desktop, env))
			want := DataSelection{DatabasePath: filepath.Join(l.external, tc.chosen), DataRoot: l.external, ConfigRoot: l.exe, Explicit: true}
			if err != nil || got != want {
				t.Fatalf("explicit database lost to another location: %+v %v", got, err)
			}
			assertTreeUnchanged(t, l.root, before)
		})
	}
}

func TestExplicitDatabaseUsesCanonicalIdentity(t *testing.T) {
	l := newDataLayout(t)
	database := filepath.Join(l.external, "real.db")
	writeFixture(t, database, "existing")
	fileLink, dirLink := filepath.Join(l.cwd, "alias.db"), filepath.Join(l.root, "external-link")
	if err := os.Symlink(database, fileLink); err != nil {
		t.Skipf("symlinks are unavailable: %v", err)
	}
	if err := os.Symlink(l.external, dirLink); err != nil {
		t.Skipf("directory symlinks are unavailable: %v", err)
	}
	before := treeState(t, l.root)
	got, err := ResolveData(l.input(Desktop, map[string]string{"GATEWAY_DB_PATH": fileLink}))
	if err != nil || got.DatabasePath != database || got.DataRoot != l.external {
		t.Fatalf("database symlink identity: %+v %v", got, err)
	}
	fresh := filepath.Join(l.external, "fresh.db")
	got, err = ResolveData(l.input(Desktop, map[string]string{"GATEWAY_DB_PATH": filepath.Join(dirLink, "fresh.db")}))
	if err != nil || got.DatabasePath != fresh || got.DataRoot != l.external || got.ConfirmCreate {
		t.Fatalf("directory symlink identity: %+v %v", got, err)
	}
	assertTreeUnchanged(t, l.root, before)
}

func TestPathFailuresDoNotFallBackOrCreate(t *testing.T) {
	for _, tc := range []struct {
		name  string
		setup func(t *testing.T, l dataLayout, in *DataInput)
	}{
		{"explicit parent is missing", func(_ *testing.T, l dataLayout, in *DataInput) {
			in.Environment["GATEWAY_DB_PATH"] = filepath.Join(l.external, "missing", "data.db")
		}},
		{"explicit target is a directory", func(t *testing.T, l dataLayout, in *DataInput) {
			in.Environment["DB_PATH"] = filepath.Join(l.external, "directory.db")
			if err := os.Mkdir(in.Environment["DB_PATH"], 0o700); err != nil {
				t.Fatal(err)
			}
		}},
		{"relative alias parent is missing", func(_ *testing.T, _ dataLayout, in *DataInput) {
			in.Environment["SQLITE_PATH"] = filepath.Join("missing", "data.db")
		}},
		{"executable directory is missing", func(_ *testing.T, l dataLayout, in *DataInput) {
			in.Executable = filepath.Join(l.root, "gone", "gateway.exe")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			l := newDataLayout(t)
			in := l.input(Desktop, map[string]string{})
			tc.setup(t, l, &in)
			before := treeState(t, l.root)
			got, err := ResolveData(in)
			if !IsCode(err, "startup.path_unusable") || got != (DataSelection{}) {
				t.Fatalf("unusable path produced a selection: %+v %v", got, err)
			}
			assertTreeUnchanged(t, l.root, before)
		})
	}
}

func TestUnreadableLocationsFailInsteadOfLookingEmpty(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("POSIX permission bits are required and root bypasses them; Windows ACL behavior is NOT RUN here")
	}
	for _, tc := range []struct {
		name     string
		lock     func(l dataLayout) string
		explicit bool
		codes    []string
	}{
		// The executable-side dotenv is read before the database is inspected,
		// so an unreadable executable directory is classified as invalid configuration.
		{name: "executable directory cannot be inspected", lock: func(l dataLayout) string { return l.exe }, codes: []string{"startup.config_invalid", "startup.path_unusable"}},
		{name: "working directory cannot rule out a legacy database", lock: func(l dataLayout) string { return l.cwd }, codes: []string{"startup.path_unusable"}},
		{name: "explicit database directory cannot be inspected", lock: func(l dataLayout) string { return l.external }, explicit: true, codes: []string{"startup.path_unusable"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			l := newDataLayout(t)
			in := l.input(Desktop, map[string]string{})
			if tc.explicit {
				in.Environment["GATEWAY_DB_PATH"] = filepath.Join(l.external, "data.db")
			}
			before := treeState(t, l.root)
			locked := tc.lock(l)
			if err := os.Chmod(locked, 0); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = os.Chmod(locked, 0o700) })
			got, err := ResolveData(in)
			if chmodErr := os.Chmod(locked, 0o700); chmodErr != nil {
				t.Fatal(chmodErr)
			}
			if !slices.Contains(tc.codes, Code(err)) || got != (DataSelection{}) {
				t.Fatalf("permission failure became a database choice: %+v %v", got, err)
			}
			assertTreeUnchanged(t, l.root, before)
		})
	}
}

func TestCLIModesResolveFromWorkingDirectoryWithoutConfirmation(t *testing.T) {
	for _, mode := range []Mode{Console, Headless} {
		t.Run(string(mode), func(t *testing.T) {
			l := newDataLayout(t)
			// Executable-side files belong to the desktop launcher only.
			writeFixture(t, filepath.Join(l.exe, ".env"), "GATEWAY_DB_PATH="+filepath.Join(l.external, "exe.db")+"\n")
			writeFixture(t, filepath.Join(l.exe, "datalink.db"), "desktop")
			in := l.input(mode, map[string]string{})
			before := treeState(t, l.root)
			got, err := ResolveData(in)
			want := DataSelection{DatabasePath: filepath.Join(l.cwd, "datalink.db"), DataRoot: l.cwd, ConfigRoot: l.cwd}
			if err != nil || got != want {
				t.Fatalf("CLI default moved away from the working directory: %+v %v", got, err)
			}
			in.Environment["DB_PATH"] = "chosen.db"
			got, err = ResolveData(in)
			want = DataSelection{DatabasePath: filepath.Join(l.cwd, "chosen.db"), DataRoot: l.cwd, ConfigRoot: l.cwd, Explicit: true}
			if err != nil || got != want {
				t.Fatalf("CLI relative alias left the working directory: %+v %v", got, err)
			}
			assertTreeUnchanged(t, l.root, before)
			// A legacy cwd database is the CLI's own data, never an ambiguity.
			writeFixture(t, filepath.Join(l.cwd, ".env"), "PORT=9000\n")
			writeFixture(t, filepath.Join(l.cwd, "datalink.db"), "legacy")
			delete(in.Environment, "DB_PATH")
			before = treeState(t, l.root)
			got, err = ResolveData(in)
			want = DataSelection{DatabasePath: filepath.Join(l.cwd, "datalink.db"), DataRoot: l.cwd, ConfigRoot: l.cwd}
			if err != nil || got != want {
				t.Fatalf("CLI treated its own data as ambiguous: %+v %v", got, err)
			}
			assertTreeUnchanged(t, l.root, before)
		})
	}
}
