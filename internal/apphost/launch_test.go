package apphost

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLaunchModesAndExplicitOverride(t *testing.T) {
	for _, tc := range []struct {
		desktop bool
		args    []string
		browser string
		mode    Mode
		open    bool
	}{
		{true, nil, "", Desktop, true}, {true, nil, "0", Desktop, false}, {true, []string{"--headless"}, "true", Headless, false},
		{false, nil, "1", Console, true}, {false, nil, "", Console, false}, {false, []string{"--headless"}, "1", Headless, false},
	} {
		got, err := Parse(tc.args, tc.desktop)
		if err != nil || got.Mode != tc.mode || got.AutoOpen(tc.browser) != tc.open {
			t.Fatalf("%+v %v", got, err)
		}
	}
	for _, arg := range []string{"--help", "--version"} {
		got, err := Parse([]string{arg}, true)
		if err != nil || !got.Informational() {
			t.Fatalf("%+v %v", got, err)
		}
	}
	if _, err := Parse([]string{"--unknown"}, true); err == nil {
		t.Fatal("accepted unsupported option")
	}
}

func TestDataSelectionNeverCreatesOrGuesses(t *testing.T) {
	root := t.TempDir()
	exe, cwd := filepath.Join(root, "exe"), filepath.Join(root, "cwd")
	for _, dir := range []string{exe, cwd} {
		if err := os.Mkdir(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	in := DataInput{Mode: Desktop, Executable: filepath.Join(exe, "gateway.exe"), WorkingDir: cwd, Environment: map[string]string{}}
	got, err := ResolveData(in)
	if err != nil {
		t.Fatal(err)
	}
	if !got.ConfirmCreate || got.DatabasePath != filepath.Join(exe, "datalink.db") {
		t.Fatalf("%+v", got)
	}
	if _, err = os.Stat(got.DatabasePath); !os.IsNotExist(err) {
		t.Fatal("selection created storage")
	}
	if err = os.WriteFile(filepath.Join(cwd, ".env"), []byte("DB_PATH=legacy.db\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err = ResolveData(in); !IsCode(err, "startup.path_ambiguous") {
		t.Fatalf("%v", err)
	}
	in.Environment["DB_PATH"] = filepath.Join(cwd, "selected#%.db")
	if err = os.WriteFile(filepath.Join(exe, ".env"), []byte("GATEWAY_DB_PATH=wrong.db\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err = ResolveData(in)
	if err != nil || got.DatabasePath != in.Environment["DB_PATH"] || got.ConfirmCreate {
		t.Fatalf("process alias lost: %+v %v", got, err)
	}
	in.Environment["GATEWAY_DB_PATH"] = filepath.Join(cwd, "first.db")
	got, err = ResolveData(in)
	if err != nil || got.DatabasePath != in.Environment["GATEWAY_DB_PATH"] {
		t.Fatalf("alias order: %+v %v", got, err)
	}
}

func TestCLIAndDesktopConfigRootsRemainSeparate(t *testing.T) {
	root := t.TempDir()
	in := DataInput{Mode: Console, Executable: filepath.Join(root, "gateway"), WorkingDir: root, Environment: map[string]string{}}
	if err := os.WriteFile(filepath.Join(root, ".env"), []byte("BROKEN='unterminated"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ResolveData(in); err != nil {
		t.Fatal("CLI reparsed dotenv instead of loaded environment", err)
	}
	in.Mode = Desktop
	if _, err := ResolveData(in); !IsCode(err, "startup.config_invalid") {
		t.Fatal(err)
	}
	in.Environment["SQLITE_PATH"] = "chosen.db"
	got, err := ResolveData(in)
	if err != nil || got.DatabasePath != filepath.Join(root, "chosen.db") {
		t.Fatalf("%+v %v", got, err)
	}
	in.Environment["SQLITE_PATH"] = "file:foo?mode=memory"
	if _, err := ResolveData(in); !IsCode(err, "startup.path_unusable") {
		t.Fatal(err)
	}
}

func TestVerifiedListenerURLs(t *testing.T) {
	for _, tc := range []struct {
		addr, url string
		logs      bool
	}{{"0.0.0.0:8080", "http://127.0.0.1:8080", true}, {"[::]:8090", "http://[::1]:8090", true}, {"[::1]:8090", "http://[::1]:8090", true}, {"192.0.2.4:8080", "http://192.0.2.4:8080", false}, {"127.0.0.2:8080", "http://127.0.0.2:8080", false}, {"[fe80::1%Ethernet]:8080", "http://[fe80::1%25Ethernet]:8080", false}} {
		got, logs, err := ListenerURL(tc.addr)
		if err != nil || got != tc.url || logs != tc.logs {
			t.Fatalf("%s %s %v %v", tc.addr, got, logs, err)
		}
	}
	if _, _, err := ListenerURL("unknown.example:8080"); err == nil {
		t.Fatal("unverified name accepted")
	}
}
