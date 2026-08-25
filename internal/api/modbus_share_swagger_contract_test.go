package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestModbusShareSwaggerFailureStatusesMatchReachableHandlerPaths(t *testing.T) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate swagger contract test")
	}
	swaggerPath := filepath.Join(filepath.Dir(sourceFile), "..", "..", "docs", "swagger", "swagger.json")
	raw, err := os.ReadFile(swaggerPath)
	if err != nil {
		t.Fatalf("read generated swagger: %v", err)
	}

	var document struct {
		Paths map[string]map[string]struct {
			Responses map[string]json.RawMessage `json:"responses"`
		} `json:"paths"`
	}
	if err := json.Unmarshal(raw, &document); err != nil {
		t.Fatalf("decode generated swagger: %v", err)
	}

	expected := map[string]map[string][]string{
		"/datalink/modbus-share/start": {
			"post": {"200", "400", "409", "422", "503"},
		},
		"/datalink/modbus-share/stop": {
			"post": {"200", "400", "409", "422", "500", "503"},
		},
		"/datalink/modbus-share/mappings": {
			"get": {"200", "403", "422", "503"},
		},
		"/datalink/modbus-share/mappings/{tagId}": {
			"put":    {"200", "400", "403", "422", "503"},
			"delete": {"200", "403", "422", "503"},
		},
		"/datalink/modbus-share/reconcile": {
			"post": {"200", "400", "403", "409", "422", "503"},
		},
	}

	for path, methods := range expected {
		for method, statuses := range methods {
			operation, ok := document.Paths[path][method]
			if !ok {
				t.Fatalf("generated swagger missing %s %s", method, path)
			}
			actual := make(map[string]struct{}, len(operation.Responses))
			for status := range operation.Responses {
				actual[status] = struct{}{}
			}
			if len(actual) != len(statuses) {
				t.Fatalf("%s %s statuses = %v, want exactly %v", method, path, actual, statuses)
			}
			for _, status := range statuses {
				if _, ok := actual[status]; !ok {
					t.Fatalf("%s %s missing reachable status %s; got %v", method, path, status, actual)
				}
			}
		}
	}
}
