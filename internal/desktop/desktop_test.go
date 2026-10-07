package desktop

import (
	"strings"
	"testing"
	"time"
)

func TestProductURLBoundary(t *testing.T) {
	for _, u := range []string{"http://127.0.0.1:8080/studio/v2", "http://[::1]:9000/studio/logs", "http://192.168.1.5:80/studio/v2", "http://[fe80::1%25Ethernet]:8080/studio/v2", "http://[fe80::1%25Ethernet%202]:8080/studio/v2"} {
		if err := validateURL(u); err != nil {
			t.Errorf("valid URL: %s: %v", u, err)
		}
	}
	for _, u := range []string{"cmd.exe", "file:///C:/test", "http://example.org:80/studio/v2", "http://user@localhost:80/studio/v2", "http://localhost:80/studio/v2?exit=1", "http://localhost:80/studio/v2#x", "http://localhost:80/other", "http://localhost/studio/v2", "http://localhost:0/studio/v2", "http://[fe80::1%25]:8080/studio/v2", "http://[fe80::1%25Ethernet%2fcommand]:8080/studio/v2", "http://[fe80::1%25Ethernet%22bad]:8080/studio/v2"} {
		if validateURL(u) == nil {
			t.Errorf("unsafe URL accepted: %s", u)
		}
	}
}
func TestNativeVersionUsesActualMetadataAndLocalPath(t *testing.T) {
	text := versionText(Options{})
	if !strings.Contains(text, "dev") || !strings.Contains(text, "unknown") || strings.Contains(text, "資料來源") {
		t.Fatal(text)
	}
	path := `C:\資料\gateway#1%20.db`
	text = versionText(Options{Version: "1.2.3", Commit: "abc123", BuildMode: "desktop", DataPath: path})
	if !strings.Contains(text, "資料來源："+path) || !strings.Contains(text, "1.2.3") {
		t.Fatal(text)
	}
}
func TestNativeErrorAdvertisesOnlySavedDiagnostic(t *testing.T) {
	info := ErrorInfo{Code: "STARTUP_FAILED", Message: "Startup failed", NextStep: "Check configuration", DiagnosticPath: "private-path"}
	if strings.Contains(errorText(info), info.DiagnosticPath) || !strings.Contains(errorText(info), "未保存") {
		t.Fatal(errorText(info))
	}
	info.DiagnosticSaved = true
	if !strings.Contains(errorText(info), info.DiagnosticPath) {
		t.Fatal(errorText(info))
	}
}
func TestMinimalMenuProjectsTruthfulState(t *testing.T) {
	items := menuItems(State{Process: "starting", Acquisition: "unknown"})
	if len(items) != 6 || items[0].Enabled || !items[1].Enabled || items[2].Enabled || items[3].Enabled {
		t.Fatalf("%+v", items)
	}
	items = menuItems(State{Process: "stopping"})
	if items[5].Enabled {
		t.Fatal("repeat quit enabled")
	}
	items = menuItems(State{Process: "stop-timeout", PendingPhase: "delivery"})
	if !items[5].Enabled || !strings.Contains(items[2].Label, "delivery") {
		t.Fatalf("%+v", items)
	}
}

func TestFaultIsRecordedBeforeQuit(t *testing.T) {
	events := make(chan string, 2)
	requestFaultStop(Options{OnFault: func() { events <- "fault" }, OnQuit: func() { events <- "quit" }})
	timeout := time.NewTimer(time.Second)
	defer timeout.Stop()
	for _, want := range []string{"fault", "quit"} {
		select {
		case got := <-events:
			if got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
		case <-timeout.C:
			t.Fatal("fault notification did not complete")
		}
	}
}
