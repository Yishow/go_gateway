package diagnostics

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestUnknownSecretInputIsSuppressedBeforeStorage(t *testing.T) {
	b, err := New(Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = b.Close(t.Context()) }()
	secret := "postgres://private:credential@private-host/db?token=secret"
	b.Emit(Input{Code: secret, Fields: map[string]any{"message": secret, "headers": map[string]string{"Authorization": secret}}})
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	snap, err := b.Snapshot(Query{})
	if err != nil {
		t.Fatal(err)
	}
	if len(snap.Records) != 1 {
		t.Fatalf("want fixed suppression record, got %d", len(snap.Records))
	}
	data, _ := json.Marshal(snap)
	if strings.Contains(string(data), secret) || snap.Records[0].Code != "raw.suppressed" {
		t.Fatal("unsafe projection")
	}
}

func TestHardlinkedDatabaseHasSafeStartupErrorTemplate(t *testing.T) {
	const code = "startup.database_linked"
	const message = "The selected database has multiple hardlinks, so its WAL and SHM locations are ambiguous. Use one canonical database location."
	event := project(Input{Code: code, Fields: map[string]any{
		"path":   "private-database-path",
		"error":  "private-raw-error",
		"source": "private-source",
	}}, nil)
	if event.Code != code || event.Level != "error" || event.Source != "startup" || event.Message != message {
		t.Fatalf("hardlink fault was not classified safely: %#v", event)
	}
	if len(event.Fields) != 0 {
		t.Fatal("hardlink fault retained unapproved fields")
	}
	data, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "private-") {
		t.Fatal("hardlink fault leaked raw diagnostic input")
	}
}
