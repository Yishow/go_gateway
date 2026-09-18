package handlers

import (
	"encoding/json"
	"strings"
	"testing"

	"go-gateway/internal/datalink/schema"
)

func dsnConnector(config map[string]string) *schema.DatabaseConnector {
	raw, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	return &schema.DatabaseConnector{Kind: schema.DatabaseConnectorKindSQLite, ConnectionConfig: string(raw)}
}

func TestRecordingTargetDatabaseNeverCarriesDSNCredentials(t *testing.T) {
	cases := []struct {
		name   string
		config map[string]string
		want   string
	}{
		{
			name:   "url dsn loses userinfo",
			config: map[string]string{"dsn": "postgres://admin:s3cret@db.local:5432/metrics"},
			want:   "postgres://db.local:5432/metrics",
		},
		{
			name:   "keyword dsn loses password",
			config: map[string]string{"dsn": "host=db.local password=s3cret port=5432 dbname=metrics"},
			want:   "host=db.local port=5432 dbname=metrics",
		},
		{
			name:   "plain file path stays readable",
			config: map[string]string{"path": "/data/metrics.db"},
			want:   "/data/metrics.db",
		},
		{
			name:   "query part is dropped before display",
			config: map[string]string{"dsn": "/data/metrics.db?password=s3cret"},
			want:   "/data/metrics.db",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			connector := dsnConnector(tc.config)
			got := recordingTargetDatabase(connector, recordingDialectSQLite)
			if got != tc.want {
				t.Fatalf("recordingTargetDatabase = %q, want %q", got, tc.want)
			}
			if strings.Contains(got, "s3cret") {
				t.Fatalf("database name %q must not carry credentials", got)
			}
		})
	}
}
