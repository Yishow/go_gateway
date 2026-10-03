package dbtarget

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/snapshot"
)

// GroupRowMember binds one snapshot member to a destination column.
type GroupRowMember struct {
	MemberKey string
	// EntityKey partitions this member into the row emitted for that entity.
	// Empty means the member belongs to the fixed group scope.
	EntityKey string
	Column    string
	Type      measurement.ExactType
	Required  bool
}

// GroupRowSpec describes how snapshot rows map onto one inspected table. The
// Columns come from real table inspection, not from caller declarations.
type GroupRowSpec struct {
	Dialect SQLDialect
	Columns []ColumnInfo
	Members []GroupRowMember
	// Partial opts in to NULL members; it needs nullable optional columns and
	// a provenance column that stores member-level quality.
	Partial bool
	// EntityKeyed marks groups whose rows are partitioned by entity; the SQL
	// row must then carry the entity or the record key to stay distinguishable.
	EntityKeyed       bool
	EntityKeyColumn   string
	RecordKeyColumn   string
	BucketStartColumn string
	ProvenanceColumn  string
}

// Safe layout issue and encoding error codes. They never contain values.
const (
	issueUnsupportedDialect    = "unsupported-dialect"
	issueColumnMissing         = "column-missing"
	issueDuplicateColumn       = "duplicate-column"
	issueUnsupportedSQLType    = "unsupported-sql-type"
	issueColumnNotNullable     = "column-not-nullable"
	issuePartialNeedsProvenanc = "partial-requires-provenance"
	issueIdentityUnsupported   = "identity-column-unsupported"
	issueIdentityMissing       = "identity-column-missing"
	issueMemberIncomplete      = "member-incomplete"

	rowErrNotARow        = "not-a-row"
	rowErrMemberMissing  = "member-missing"
	rowErrTypeMismatch   = "type-mismatch"
	rowErrSQLValueBlock  = "sql-value-blocked"
	rowErrUnexpectedNull = "unexpected-null"
)

// LayoutIssue is a safe, code-only reason a layout cannot be activated.
type LayoutIssue struct {
	Code      string
	Column    string
	MemberKey string
}

// GroupRowLayout is a validated mapping from snapshot rows to table columns.
type GroupRowLayout struct {
	spec    GroupRowSpec
	columns map[string]ColumnInfo
	caps    snapshot.StorageCapabilities
}

// EncodedCell is one column value; a nil Value is SQL NULL.
type EncodedCell struct {
	Column string
	Value  driver.Value
}

// EncodedRow is a snapshot row encoded for the destination.
type EncodedRow struct {
	RecordID    string
	EffectKey   string
	EntityKey   string
	BucketStart time.Time
	Partial     bool
	Cells       []EncodedCell
}

// GroupRowError is a safe encoding failure that never carries values.
type GroupRowError struct {
	Code   string
	Column string
}

func (e *GroupRowError) Error() string { return "group row blocked: " + e.Code }

// NewGroupRowLayout validates every mapped column against the inspected table
// and returns all blocking issues at once, or a ready layout.
func NewGroupRowLayout(spec GroupRowSpec) (*GroupRowLayout, []LayoutIssue) {
	if spec.Dialect != SQLDialectSQLite && spec.Dialect != SQLDialectPostgres {
		return nil, []LayoutIssue{{Code: issueUnsupportedDialect}}
	}
	columns := make(map[string]ColumnInfo, len(spec.Columns))
	for _, column := range spec.Columns {
		columns[strings.ToLower(column.Name)] = column
	}
	var issues []LayoutIssue
	lookup := func(name string) (ColumnInfo, bool) {
		column, ok := columns[strings.ToLower(name)]
		if !ok {
			issues = append(issues, LayoutIssue{Code: issueColumnMissing, Column: name})
		}
		return column, ok
	}

	nullableOptional := true
	entityKeys := make([]string, 0, len(spec.Members))
	seenEntities := make(map[string]struct{}, len(spec.Members))
	for _, member := range spec.Members {
		if _, seen := seenEntities[member.EntityKey]; !seen {
			seenEntities[member.EntityKey] = struct{}{}
			entityKeys = append(entityKeys, member.EntityKey)
		}
		if member.MemberKey == "" || member.Column == "" || member.Type == "" {
			issues = append(issues, LayoutIssue{Code: issueMemberIncomplete, MemberKey: member.MemberKey})
			continue
		}
		column, ok := lookup(member.Column)
		if !ok {
			continue
		}
		if _, err := exactStrategyFor(spec.Dialect, column.DataType, member.Type); err != nil {
			issues = append(issues, LayoutIssue{Code: issueUnsupportedSQLType, Column: member.Column, MemberKey: member.MemberKey})
		}
		if !member.Required && !column.Nullable {
			nullableOptional = false
			if spec.Partial {
				issues = append(issues, LayoutIssue{Code: issueColumnNotNullable, Column: member.Column, MemberKey: member.MemberKey})
			}
		}
	}

	identities := []struct {
		column string
		check  func(ColumnInfo) bool
	}{
		{spec.RecordKeyColumn, func(c ColumnInfo) bool { return isTextColumn(spec.Dialect, c) }},
		{spec.EntityKeyColumn, func(c ColumnInfo) bool { return isTextColumn(spec.Dialect, c) }},
		{spec.BucketStartColumn, func(c ColumnInfo) bool { return bucketColumnKind(spec.Dialect, c) != bucketUnsupported }},
		{spec.ProvenanceColumn, func(c ColumnInfo) bool { return isProvenanceColumn(spec.Dialect, c) }},
	}
	if len(entityKeys) == 0 {
		entityKeys = append(entityKeys, "")
	}
	for _, entityKey := range entityKeys {
		used := make(map[string]struct{}, len(spec.Members)+len(identities))
		claim := func(name string) {
			if name == "" {
				return
			}
			key := strings.ToLower(name)
			if _, dup := used[key]; dup {
				issues = append(issues, LayoutIssue{Code: issueDuplicateColumn, Column: name})
			}
			used[key] = struct{}{}
		}
		for _, member := range spec.Members {
			if member.EntityKey == entityKey && member.MemberKey != "" && member.Column != "" && member.Type != "" {
				claim(member.Column)
			}
		}
		for _, identity := range identities {
			claim(identity.column)
		}
	}

	for _, identity := range identities {
		if identity.column == "" {
			continue
		}
		if column, ok := lookup(identity.column); ok && !identity.check(column) {
			issues = append(issues, LayoutIssue{Code: issueIdentityUnsupported, Column: identity.column})
		}
	}
	if spec.EntityKeyed && spec.EntityKeyColumn == "" && spec.RecordKeyColumn == "" {
		issues = append(issues, LayoutIssue{Code: issueIdentityMissing})
	}
	provenanceOK := spec.ProvenanceColumn != "" && !hasIssue(issues, issueIdentityUnsupported, spec.ProvenanceColumn)
	if spec.Partial && spec.ProvenanceColumn == "" {
		issues = append(issues, LayoutIssue{Code: issuePartialNeedsProvenanc})
	}
	if len(issues) > 0 {
		return nil, issues
	}
	return &GroupRowLayout{
		spec:    spec,
		columns: columns,
		caps:    snapshot.StorageCapabilities{NullableValues: nullableOptional, MemberQuality: provenanceOK},
	}, nil
}

func hasIssue(issues []LayoutIssue, code, column string) bool {
	for _, issue := range issues {
		if issue.Code == code && issue.Column == column {
			return true
		}
	}
	return false
}

// Capabilities reports what the inspected table can honestly store; feed it
// to the snapshot config so partial rows are only accepted when supported.
func (l *GroupRowLayout) Capabilities() snapshot.StorageCapabilities { return l.caps }

type provenanceEntry struct {
	Member     string `json:"member"`
	Status     string `json:"status"`
	Reason     string `json:"reason,omitempty"`
	SampleID   string `json:"sample_id,omitempty"`
	ObservedAt string `json:"observed_at,omitempty"`
	Quality    string `json:"quality,omitempty"`
}

func (l *GroupRowLayout) membersForEntity(entityKey string) []GroupRowMember {
	entityScoped := false
	for _, member := range l.spec.Members {
		if member.EntityKey != "" {
			entityScoped = true
			break
		}
	}
	if !entityScoped {
		return l.spec.Members
	}
	members := make([]GroupRowMember, 0, len(l.spec.Members))
	for _, member := range l.spec.Members {
		if member.EntityKey == entityKey {
			members = append(members, member)
		}
	}
	return members
}

// EncodeRow converts a snapshot row into parameterized column values. Members
// that are not usable are written as NULL only in an explicit partial layout.
func (l *GroupRowLayout) EncodeRow(outcome snapshot.Outcome) (EncodedRow, error) {
	if outcome.Kind != snapshot.OutcomeRow {
		return EncodedRow{}, &GroupRowError{Code: rowErrNotARow}
	}
	results := make(map[string]snapshot.MemberResult, len(outcome.Members))
	for _, result := range outcome.Members {
		results[result.MemberKey] = result
	}
	encoded := EncodedRow{
		RecordID: outcome.RecordID, EffectKey: outcome.EffectKey, EntityKey: outcome.EntityKey,
		BucketStart: outcome.BucketStart.UTC(), Partial: outcome.Partial,
	}
	rowMembers := l.membersForEntity(outcome.EntityKey)
	if len(rowMembers) == 0 {
		return EncodedRow{}, &GroupRowError{Code: rowErrMemberMissing}
	}
	for _, member := range rowMembers {
		result, ok := results[member.MemberKey]
		if !ok {
			return EncodedRow{}, &GroupRowError{Code: rowErrMemberMissing, Column: member.Column}
		}
		if !result.Usable() {
			if !l.spec.Partial {
				return EncodedRow{}, &GroupRowError{Code: rowErrUnexpectedNull, Column: member.Column}
			}
			encoded.Cells = append(encoded.Cells, EncodedCell{Column: member.Column})
			continue
		}
		if result.Sample == nil {
			return EncodedRow{}, &GroupRowError{Code: rowErrMemberMissing, Column: member.Column}
		}
		if result.Sample.Value.Type() != member.Type {
			return EncodedRow{}, &GroupRowError{Code: rowErrTypeMismatch, Column: member.Column}
		}
		column := l.columns[strings.ToLower(member.Column)]
		bound, err := EncodeExactValue(l.spec.Dialect, column.DataType, result.Sample.Value)
		if err != nil {
			return EncodedRow{}, &GroupRowError{Code: rowErrSQLValueBlock, Column: member.Column}
		}
		encoded.Cells = append(encoded.Cells, EncodedCell{Column: member.Column, Value: bound})
	}
	if err := l.appendIdentityCells(&encoded, outcome); err != nil {
		return EncodedRow{}, err
	}
	return encoded, nil
}

func (l *GroupRowLayout) appendIdentityCells(encoded *EncodedRow, outcome snapshot.Outcome) error {
	text := func(columnName, value string) error {
		if columnName == "" {
			return nil
		}
		bound, err := EncodeExactValue(l.spec.Dialect, l.columns[strings.ToLower(columnName)].DataType, measurement.NewText(value))
		if err != nil {
			return &GroupRowError{Code: rowErrSQLValueBlock, Column: columnName}
		}
		encoded.Cells = append(encoded.Cells, EncodedCell{Column: columnName, Value: bound})
		return nil
	}
	if err := text(l.spec.RecordKeyColumn, outcome.RecordID); err != nil {
		return err
	}
	if name := l.spec.BucketStartColumn; name != "" {
		start := outcome.BucketStart.UTC()
		if bucketColumnKind(l.spec.Dialect, l.columns[strings.ToLower(name)]) == bucketTimestamp {
			encoded.Cells = append(encoded.Cells, EncodedCell{Column: name, Value: start})
		} else if err := text(name, start.Format(time.RFC3339Nano)); err != nil {
			return err
		}
	}
	if err := text(l.spec.EntityKeyColumn, outcome.EntityKey); err != nil {
		return err
	}
	if name := l.spec.ProvenanceColumn; name != "" {
		payload, err := json.Marshal(provenanceEntries(outcome))
		if err != nil {
			return &GroupRowError{Code: rowErrSQLValueBlock, Column: name}
		}
		if isJSONColumn(l.spec.Dialect, l.columns[strings.ToLower(name)]) {
			encoded.Cells = append(encoded.Cells, EncodedCell{Column: name, Value: string(payload)})
		} else if err := text(name, string(payload)); err != nil {
			return err
		}
	}
	return nil
}

func provenanceEntries(outcome snapshot.Outcome) []provenanceEntry {
	entries := make([]provenanceEntry, 0, len(outcome.Members))
	for _, result := range outcome.Members {
		entry := provenanceEntry{Member: result.MemberKey, Status: string(result.Status), Reason: result.Reason}
		if result.Status == snapshot.MemberOK {
			entry.Reason = ""
		}
		if result.Sample != nil {
			entry.SampleID = result.Sample.SampleID
			entry.ObservedAt = result.Sample.ObservedAt.UTC().Format(time.RFC3339Nano)
			entry.Quality = string(result.Sample.Quality)
		}
		entries = append(entries, entry)
	}
	return entries
}

type bucketKind int

const (
	bucketUnsupported bucketKind = iota
	bucketText
	bucketTimestamp
)

func isTextColumn(dialect SQLDialect, column ColumnInfo) bool {
	_, err := exactStrategyFor(dialect, column.DataType, measurement.ExactText)
	return err == nil
}

func bucketColumnKind(dialect SQLDialect, column ColumnInfo) bucketKind {
	if isTextColumn(dialect, column) {
		return bucketText
	}
	parsed, err := parseDeclaredType(column.DataType)
	if err != nil {
		return bucketUnsupported
	}
	if strings.HasPrefix(parsed.base, "TIMESTAMP") || (dialect == SQLDialectSQLite && parsed.base == "DATETIME") {
		return bucketTimestamp
	}
	return bucketUnsupported
}

func isJSONColumn(dialect SQLDialect, column ColumnInfo) bool {
	if dialect != SQLDialectPostgres {
		return false
	}
	parsed, err := parseDeclaredType(column.DataType)
	return err == nil && (parsed.base == "JSON" || parsed.base == "JSONB")
}

func isProvenanceColumn(dialect SQLDialect, column ColumnInfo) bool {
	return isTextColumn(dialect, column) || isJSONColumn(dialect, column)
}
