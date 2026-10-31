package dbtarget

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
)

// longReductionMean is the frozen reduction name of a mean member; it is the
// snapshot's own sample reduction constant, not a second spelling.
const longReductionMean = string(snapshot.SampleReductionMean)

// LongRowMember is one source of a long-layout group.
type LongRowMember struct {
	Calculation              snapshot.Calculation `json:",omitempty"`
	Interval                 time.Duration        `json:",omitempty"`
	Required                 bool                 `json:",omitempty"`
	MemberKey                string
	DeviceID, PointID, TagID string
	// CircuitID and MeasurementID are written only when known; empty is NULL.
	CircuitID, MeasurementID string
	Type                     measurement.ExactType
	// Reduction is the frozen sample reduction of this member ("mean" or empty
	// for the legacy latest rule). A mean member stores the reduction value in
	// the slot of Type, never the raw sample value.
	Reduction string
	// ValueRole and SourceValueType are the frozen recording metadata written to
	// provenance; empty means the legacy shape never provided them.
	ValueRole       string
	SourceValueType schema.DataType
	// Counter is the confirmed counter binding of a cumulative member. It
	// supplies the series epoch and unit written on every record of the member.
	Counter *aggregation.CounterBinding
}

// LongRowSpec maps closed profile records onto an inspected managed_long_v1
// table. Columns must come from real table inspection.
type LongRowSpec struct {
	ResultShape string `json:",omitempty"`
	Dialect     SQLDialect
	Columns     []ColumnInfo
	// WorkspaceID, GroupID, GroupRevision and Profile* are frozen with the
	// applied revision and written on every record.
	WorkspaceID, GroupID, GroupRevision string
	ProfileKind                         string
	ProfileInterval                     int
	ProfileUnit                         string
	Members                             []LongRowMember
}

// LongRowLayout encodes one record per member and window for the shared long
// table. It never creates a column per source.
type LongRowLayout struct {
	spec    LongRowSpec
	columns map[string]ColumnInfo
	members map[string]LongRowMember
	owner   string
}

// longMeanProvenance is the value-free evidence of a mean window. It is present
// only for a mean member; legacy latest and counter records keep the original
// payload shape.
type longMeanProvenance struct {
	Reduction       string   `json:"reduction"`
	Role            string   `json:"role,omitempty"`
	SourceType      string   `json:"source_type,omitempty"`
	ResultType      string   `json:"result_type,omitempty"`
	Version         string   `json:"version,omitempty"`
	Method          string   `json:"method,omitempty"`
	Basis           string   `json:"basis,omitempty"`
	ValidCount      int64    `json:"valid_count"`
	RejectedCount   int64    `json:"rejected_count"`
	FirstObservedAt string   `json:"first_observed_at,omitempty"`
	LastObservedAt  string   `json:"last_observed_at,omitempty"`
	Coverage        string   `json:"coverage,omitempty"`
	Reason          string   `json:"reason,omitempty"`
	Reasons         []string `json:"reasons,omitempty"`
	Sum             string   `json:"sum,omitempty"`
	Scale           int      `json:"scale"`
	Rounding        string   `json:"rounding,omitempty"`
}

type longProvenance struct {
	Type string `json:"type"`
	// SourceType/ResultType/Role are the frozen recording metadata of a latest
	// state or counter_raw member. They stay omitted for the legacy shape and
	// for mean members, whose mean block already names them.
	SourceType      string `json:"source_type,omitempty"`
	ResultType      string `json:"result_type,omitempty"`
	Role            string `json:"role,omitempty"`
	Member          string `json:"member"`
	Reason          string `json:"reason,omitempty"`
	SampleID        string `json:"sample_id,omitempty"`
	SourceRevision  string `json:"source_revision,omitempty"`
	MappingRevision string `json:"mapping_revision,omitempty"`
	Profile         struct {
		Kind     string `json:"kind"`
		Interval int    `json:"interval"`
		Unit     string `json:"unit"`
	} `json:"profile"`
	Counter *longCounterProvenance `json:"counter,omitempty"`
	Mean    *longMeanProvenance    `json:"mean,omitempty"`
	// Window is present only for a window a forward cutover cut short: which
	// side was cut and the real and nominal bounds.
	Window *longWindowProvenance `json:"window,omitempty"`
}

type longWindowProvenance struct {
	Truncated    string `json:"truncated"`
	Start        string `json:"start"`
	End          string `json:"end"`
	NominalStart string `json:"nominal_start"`
	NominalEnd   string `json:"nominal_end"`
}

// NewLongRowLayout validates every fixed column against the inspected table
// and each member's exact type against its value slot, returning all blocking
// issues at once.
func NewLongRowLayout(spec LongRowSpec) (*LongRowLayout, []LayoutIssue) {
	if issues := mixedLongIssues(spec); len(issues) > 0 {
		return nil, issues
	}
	spec = freezeMixedLong(spec)
	if spec.Dialect != SQLDialectSQLite && spec.Dialect != SQLDialectPostgres && spec.Dialect != SQLDialectMySQL && spec.Dialect != SQLDialectSQLServer {
		return nil, []LayoutIssue{{Code: issueUnsupportedDialect}}
	}
	columns := make(map[string]ColumnInfo, len(spec.Columns))
	for _, column := range spec.Columns {
		columns[strings.ToLower(column.Name)] = column
	}
	issues := validateMySQLLongIdentity(spec)
	issues = append(issues, validateSQLServerLongIdentity(spec)...)
	for _, want := range longColumnRoles {
		column, ok := columns[want.name]
		if !ok {
			issues = append(issues, LayoutIssue{Code: issueColumnMissing, Column: want.name})
			continue
		}
		if !longColumnUsable(spec.Dialect, column, want.role) {
			issues = append(issues, LayoutIssue{Code: issueUnsupportedSQLType, Column: want.name})
		}
		if !want.required && !column.Nullable {
			issues = append(issues, LayoutIssue{Code: issueColumnNotNullable, Column: want.name})
		}
	}
	owner, ownerOK := longOwnerColumn(spec.Dialect, spec.Columns)
	if !ownerOK {
		issues = append(issues, LayoutIssue{Code: issueColumnMissing, Column: "owner"})
	}
	members := make(map[string]LongRowMember, len(spec.Members))
	for _, member := range spec.Members {
		if member.MemberKey == "" || member.DeviceID == "" || member.PointID == "" || member.TagID == "" {
			issues = append(issues, LayoutIssue{Code: issueMemberIncomplete, MemberKey: member.MemberKey})
			continue
		}
		if longSlotFor(member.Type) == "" {
			issues = append(issues, LayoutIssue{Code: issueMemberIncomplete, MemberKey: member.MemberKey})
		}
		members[member.MemberKey] = member
	}
	if spec.WorkspaceID == "" || spec.GroupID == "" || spec.GroupRevision == "" || spec.ProfileKind == "" || len(members) == 0 {
		issues = append(issues, LayoutIssue{Code: issueIdentityMissing})
	}
	if len(issues) > 0 {
		return nil, issues
	}
	return &LongRowLayout{spec: spec, columns: columns, members: members, owner: owner}, nil
}

// EncodeRow encodes one member's record. A member without a usable reading is
// written with every value slot NULL, its status and its reason; the real
// observation and receive times are kept when a sample exists and are NULL
// otherwise, never invented.
func (l *LongRowLayout) EncodeRow(outcome snapshot.Outcome) (EncodedRow, error) {
	if outcome.Kind != snapshot.OutcomeRow || outcome.RecordID == "" || outcome.EffectKey == "" {
		return EncodedRow{}, &GroupRowError{Code: rowErrNotARow}
	}
	if len(outcome.Members) != 1 {
		return EncodedRow{}, &GroupRowError{Code: rowErrMemberMissing}
	}
	if l.spec.Dialect == SQLDialectSQLServer && !outcome.BucketStart.Equal(outcome.BucketStart.Truncate(time.Microsecond)) {
		return EncodedRow{}, &GroupRowError{Code: rowErrSQLValueBlock, Column: longColBucketStart}
	}
	result := outcome.Members[0]
	member, ok := l.members[result.MemberKey]
	if !ok {
		return EncodedRow{}, &GroupRowError{Code: rowErrMemberMissing}
	}
	if err := l.validateMixedOutcome(outcome, member, result); err != nil {
		return EncodedRow{}, err
	}
	if result.Counter != nil && member.Counter == nil {
		// A cumulative result must never lose its delta to a member without a binding.
		return EncodedRow{}, &GroupRowError{Code: rowErrCounterUnbound}
	}
	windowStart, windowEnd := outcome.Window()
	encoded := EncodedRow{
		DueAt: outcome.DueAt, RecordID: outcome.RecordID, EffectKey: outcome.EffectKey, EntityKey: outcome.EntityKey,
		BucketStart: windowStart.UTC(),
	}
	if l.spec.ResultShape == MixedResultShapeV1 {
		// Stable mixed identity uses the aligned bucket; SQL coverage below
		// still uses the actual (possibly partial) window.
		encoded.BucketStart = outcome.BucketStart.UTC()
	}
	if l.spec.Dialect == SQLDialectSQLServer {
		encoded.BucketStart = encoded.BucketStart.Truncate(time.Microsecond)
	}
	values := map[string]driver.Value{}
	if result.Usable() {
		value, err := windowValue(member, result)
		if err != nil {
			return EncodedRow{}, err
		}
		if value.Type() != member.Type {
			return EncodedRow{}, &GroupRowError{Code: rowErrTypeMismatch, Column: longColValueInt}
		}
		slot := longSlotFor(member.Type)
		bound, err := EncodeExactValue(l.spec.Dialect, l.columns[slot].DataType, value)
		if err != nil {
			return EncodedRow{}, &GroupRowError{Code: rowErrSQLValueBlock, Column: slot}
		}
		values[slot] = bound
	}
	text := l.text
	values[longColRecordID] = text(longColRecordID, outcome.RecordID)
	values[longColWorkspaceID] = text(longColWorkspaceID, l.spec.WorkspaceID)
	values[longColGroupID] = text(longColGroupID, l.spec.GroupID)
	values[longColGroupRev] = text(longColGroupRev, l.spec.GroupRevision)
	values[longColCircuitID] = text(longColCircuitID, member.CircuitID)
	values[longColDeviceID] = text(longColDeviceID, member.DeviceID)
	values[longColPointID] = text(longColPointID, member.PointID)
	values[longColTagID] = text(longColTagID, member.TagID)
	values[longColMeasurement] = text(longColMeasurement, member.MeasurementID)
	values[longColProfileKind] = text(longColProfileKind, l.spec.ProfileKind)
	values[longColValueType] = text(longColValueType, string(member.Type))
	values[longColStatus] = text(longColStatus, longStatus(result, outcome.Truncated != ""))
	l.counterCells(values, member, result)
	values[longColBucketStart] = l.instant(longColBucketStart, windowStart)
	values[longColBucketEnd] = l.instant(longColBucketEnd, windowEnd)
	if result.Sample != nil {
		values[longColQuality] = text(longColQuality, string(result.Sample.Quality))
		values[longColObservedAt] = l.instant(longColObservedAt, result.Sample.ObservedAt)
		values[longColReceivedAt] = l.instant(longColReceivedAt, result.Sample.ReceivedAt)
	}
	provenance, err := l.provenance(member, result, outcome)
	if err != nil {
		return EncodedRow{}, &GroupRowError{Code: rowErrSQLValueBlock, Column: longColProvenance}
	}
	values[longColProvenance] = provenance

	for _, want := range longColumnRoles {
		value := values[want.name]
		if blocked, isErr := value.(errCell); isErr {
			return EncodedRow{}, &GroupRowError{Code: rowErrSQLValueBlock, Column: blocked.column}
		}
		encoded.Cells = append(encoded.Cells, EncodedCell{Column: l.columns[want.name].Name, Value: value})
	}
	// The confirmed owner marker is written explicitly, so a compatible table
	// does not depend on an unchecked column default for every insertion.
	encoded.Cells = append(encoded.Cells, EncodedCell{Column: l.owner, Value: "gateway"})
	return encoded, nil
}

// errCell marks a value the target column cannot hold; it never reaches a row.
type errCell struct{ column string }

// instant normalizes a time to UTC microseconds. The zero time is unknown and
// stays NULL.
func (l *LongRowLayout) instant(column string, at time.Time) driver.Value {
	if at.IsZero() {
		return nil
	}
	at = at.UTC().Truncate(time.Microsecond)
	if l.spec.Dialect == SQLDialectMySQL || l.spec.Dialect == SQLDialectSQLServer || bucketColumnKind(l.spec.Dialect, l.columns[column]) == bucketTimestamp {
		return at
	}
	return at.Format(longTimeText)
}

// text binds a text cell; the empty string is NULL.
func (l *LongRowLayout) text(column, value string) driver.Value {
	if value == "" {
		return nil
	}
	bound, err := EncodeExactValue(l.spec.Dialect, l.columns[column].DataType, measurement.NewText(value))
	if err != nil {
		return errCell{column: column}
	}
	return bound
}

// longStatus names the record's state. A value that is fine but covers only
// part of its window, because a cutover cut the window short, is "partial"
// rather than "ok"; every other state (baseline, gap, missing...) keeps its
// own name and the cut shows in the provenance.
func longStatus(result snapshot.MemberResult, truncated bool) string {
	status := string(result.Status)
	switch {
	case result.Counter != nil:
		status = string(result.Counter.Status)
	case result.Usable():
		status = longStatusOK
	}
	if truncated && status == longStatusOK {
		return longStatusPartial
	}
	return status
}

// windowValue returns the value a long record stores: the reduction value for a
// mean member and the newest eligible reading for every other member. A mean
// without its reduction result, or a raw reading that would stand in for it, is
// refused instead of silently degrading to latest.
func windowValue(member LongRowMember, result snapshot.MemberResult) (measurement.ExactValue, error) {
	if member.Calculation != "" || result.Window != nil {
		return mixedMemberValue(member.Calculation, result)
	}
	if member.Reduction == longReductionMean {
		if result.Reduction == nil {
			return measurement.ExactValue{}, &GroupRowError{Code: rowErrTypeMismatch, Column: longColValueDecimal}
		}
		return result.Reduction.Value, nil
	}
	if result.Reduction != nil {
		return measurement.ExactValue{}, &GroupRowError{Code: rowErrTypeMismatch, Column: longColValueInt}
	}
	if result.Sample == nil {
		return measurement.ExactValue{}, &GroupRowError{Code: rowErrTypeMismatch, Column: longColValueInt}
	}
	return result.Sample.Value, nil
}

func (l *LongRowLayout) provenance(member LongRowMember, result snapshot.MemberResult, outcome snapshot.Outcome) (driver.Value, error) {
	if l.spec.ResultShape == MixedResultShapeV1 {
		payload, err := json.Marshal(mixedProvenance(member.mixedIntent(), result, outcome.DueAt))
		return string(payload), err
	}
	entry := longProvenance{Type: string(member.Type), Member: member.MemberKey}
	entry.Profile.Kind, entry.Profile.Interval, entry.Profile.Unit = l.spec.ProfileKind, l.spec.ProfileInterval, l.spec.ProfileUnit
	if !result.Usable() {
		entry.Reason = result.Reason
	}
	if result.Counter != nil {
		entry.Reason = result.Counter.Reason
		entry.Counter = counterProvenance(member.Counter, result.Counter)
	}
	if result.Sample != nil {
		entry.SampleID = result.Sample.SampleID
		entry.SourceRevision, entry.MappingRevision = result.Sample.SourceRevision, result.Sample.MappingRevision
	}
	if member.Reduction == longReductionMean {
		entry.Mean = meanProvenance(member.ValueRole, member.SourceValueType, member.Type, result)
	} else if member.SourceValueType != "" || member.ValueRole != "" {
		entry.SourceType, entry.ResultType, entry.Role = string(member.SourceValueType), string(member.Type), member.ValueRole
	}
	if outcome.Truncated != "" {
		start, end := outcome.Window()
		entry.Window = &longWindowProvenance{
			Truncated: outcome.Truncated, Start: start.UTC().Format(longTimeText), End: end.UTC().Format(longTimeText),
			NominalStart: outcome.BucketStart.UTC().Format(longTimeText), NominalEnd: outcome.BucketEnd.UTC().Format(longTimeText),
		}
	}
	payload, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}
	return string(payload), nil
}

// meanProvenance renders the frozen mean contract and the window evidence. It
// never fabricates an expected count or a percentage.
func meanProvenance(role string, source schema.DataType, resultType measurement.ExactType, result snapshot.MemberResult) *longMeanProvenance {
	entry := &longMeanProvenance{
		Reduction: longReductionMean, Role: role, SourceType: string(source), ResultType: string(resultType),
	}
	if evidence := result.Evidence; evidence != nil {
		entry.Version, entry.Method, entry.Basis = evidence.Version, evidence.Method, string(evidence.Basis)
		entry.ValidCount, entry.RejectedCount = evidence.ValidCount, evidence.RejectedCount
		entry.FirstObservedAt = longObservedText(evidence.FirstObservedAt)
		entry.LastObservedAt = longObservedText(evidence.LastObservedAt)
		entry.Coverage, entry.Reason = string(evidence.Coverage), evidence.Reason
		entry.Reasons = evidence.RejectionReasons
	}
	if reduction := result.Reduction; reduction != nil {
		entry.Version, entry.Method, entry.Basis = reduction.Version, reduction.Method, string(reduction.Basis)
		if digits, ok := reduction.Sum.DecimalDigits(); ok {
			entry.Sum = digits
		}
		entry.Scale, entry.Rounding = reduction.Scale, reduction.Rounding
	}
	return entry
}

func longObservedText(at time.Time) string {
	if at.IsZero() {
		return ""
	}
	return at.UTC().Format(longTimeText)
}
