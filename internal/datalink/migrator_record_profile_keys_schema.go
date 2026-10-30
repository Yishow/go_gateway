package datalink

import (
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema/migrations"
)

// Canonical UUID text is exactly google/uuid.UUID.String(): lower-case hex with
// four fixed hyphens; the nil UUID is not a record identity. Byte length also
// rejects a NUL suffix hidden by length().
const recordProfileUUIDCheck = `selection_scope<>'00000000-0000-0000-0000-000000000000'
 AND length(selection_scope)=36 AND length(CAST(selection_scope AS BLOB))=36
 AND substr(selection_scope,9,1)='-' AND substr(selection_scope,14,1)='-'
 AND substr(selection_scope,19,1)='-' AND substr(selection_scope,24,1)='-'
 AND length(replace(selection_scope,'-',''))=32
 AND replace(selection_scope,'-','') NOT GLOB '*[^0-9a-f]*'`

const recordProfileTupleCheck = `(profile_kind='mixed' AND canonical_role='record' AND circuit_id='' AND ` + recordProfileUUIDCheck + `)`

func recordProfileKeySchemas() (string, string, error) {
	raw, err := migrations.FS.ReadFile("030_write_group_profile_keys_sqlite.up.sql")
	if err != nil {
		return "", "", err
	}
	source := string(raw)
	start := strings.Index(source, "CREATE TABLE IF NOT EXISTS write_group_profile_keys (")
	if start < 0 {
		return "", "", fmt.Errorf("missing historical profile key declaration")
	}
	legacy := strings.TrimSuffix(strings.TrimSpace(source[start:]), ";")
	legacy = strings.Replace(legacy, "CREATE TABLE IF NOT EXISTS", "CREATE TABLE", 1)
	// SQLite ADD COLUMN inserts the new column before the first table constraint.
	legacy = strings.Replace(legacy, "    PRIMARY KEY", `    lineage_id TEXT REFERENCES write_group_lineages(lineage_id) ON DELETE RESTRICT CHECK (lineage_id IS NULL OR TRIM(lineage_id) <> ''),
    PRIMARY KEY`, 1)
	expanded := strings.Replace(legacy, "profile_kind IN ('instantaneous', 'cumulative')", "profile_kind IN ('instantaneous', 'cumulative', 'mixed')", 1)
	expanded = strings.Replace(expanded, "canonical_role IN ('circuit', 'unassigned')", "canonical_role IN ('circuit', 'unassigned', 'record')", 1)
	expanded = strings.Replace(expanded, "(canonical_role = 'circuit' AND circuit_id <> '' AND selection_scope = '')", "(canonical_role = 'circuit' AND circuit_id <> '' AND selection_scope = '' AND profile_kind IN ('instantaneous', 'cumulative'))", 1)
	expanded = strings.Replace(expanded, "OR (canonical_role = 'unassigned' AND circuit_id = '' AND selection_scope <> '' AND profile_kind = 'instantaneous')", "OR (canonical_role = 'unassigned' AND circuit_id = '' AND selection_scope <> '' AND profile_kind = 'instantaneous') OR "+recordProfileTupleCheck, 1)
	return legacy, expanded, nil
}
