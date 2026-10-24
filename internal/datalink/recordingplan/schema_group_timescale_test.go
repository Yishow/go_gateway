package recordingplan

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func timescaleLayout(t *testing.T) GroupSchemaLayout {
	t.Helper()
	layout, err := ManagedLongTimescaleLayout(longTestTable, longTestOwner)
	require.NoError(t, err)
	return layout
}

func timescaleScope(t *testing.T) SchemaPreviewScope {
	t.Helper()
	layout := timescaleLayout(t)
	return SchemaPreviewScope{
		WorkspaceID: "ws-ts", WorkspaceRevision: "setup-1", PlanID: "group-ts", PlanRevision: "group-rev-1",
		ConnectorID: "db-ts", ConnectorRevision: "connector-rev-1", Dialect: "postgres", Database: "gwtest", Schema: "app",
		GroupLayout: &layout, SourceDigest: "sources-1",
	}
}

// timescaleTarget answers inspections the way the production inspector does for
// a timescale layout: the group table from a state, the receipt table missing,
// and the capability and hypertable evidence on the group table only.
type timescaleTarget struct {
	layout    GroupSchemaLayout
	table     func() TargetTableInspection
	evidence  *TargetTimescaleEvidence
	inspected []string
	// receipts is true once the shared receipt table exists.
	receipts bool
}

func (s *timescaleTarget) inspect(_ context.Context, table string) (TargetTableInspection, error) {
	s.inspected = append(s.inspected, table)
	if table != s.layout.TableName {
		if s.receipts {
			state := groupReceiptCompatibleState()
			state.ColumnNullable["effect_key"] = false // PostgreSQL's primary key is NOT NULL
			return state, nil
		}
		return TargetTableInspection{Status: inspectionMissing}, nil
	}
	state := TargetTableInspection{Status: inspectionMissing}
	if s.table != nil {
		state = s.table()
	}
	state.Timescale = s.evidence
	return state, nil
}

func capableEvidence(schema string) *TargetTimescaleEvidence {
	return &TargetTimescaleEvidence{Available: true, ExtensionSchema: schema, Version: "2.31.0"}
}

// existingLong is the owned long table that already exists with its confirmed
// columns and indexes; the hypertable evidence travels separately.
func existingLong(layout GroupSchemaLayout) func() TargetTableInspection {
	return func() TargetTableInspection {
		state := groupCompatibleState(layout)
		state.IndexesInspected = true
		for _, index := range layout.Indexes {
			state.Indexes = append(state.Indexes, TargetIndexInspection{Name: index.Name, Columns: index.Columns})
		}
		return state
	}
}

func validHypertable() *TargetHypertable {
	return &TargetHypertable{TimeColumn: TimescalePartitionColumn, ChunkInterval: 7 * 24 * time.Hour, Dimensions: 1}
}

// TestTimescaleConfirmedLayout pins the timescale layout contract: the same fixed
// typed table as plain, an owned empty hypertable cut on non-null bucket_start
// with the composite key including it, the seven-day interval previewed and
// frozen into the digest, and no conversion of anything that already exists.
func TestTimescaleConfirmedLayout(t *testing.T) {
	plain, err := ManagedLongGroupLayout(dialectPostgres, longTestTable, longTestOwner)
	require.NoError(t, err)

	t.Run("the layout is the plain long table plus a frozen partition", func(t *testing.T) {
		layout := timescaleLayout(t)
		require.Equal(t, plain.Columns, layout.Columns, "typed columns, composite key and owner are the plain contract")
		require.Equal(t, plain.Indexes, layout.Indexes, "the source/time index is kept")
		require.Equal(t, &GroupSchemaPartition{Kind: PartitionKindTimescale, Column: TimescalePartitionColumn, ChunkInterval: TimescaleChunkInterval}, layout.Partition)
		var key []string
		for _, column := range layout.Columns {
			if column.PrimaryKey {
				key = append(key, column.Name)
			}
			if column.Name == TimescalePartitionColumn {
				require.False(t, column.Nullable, "the time dimension is non-null")
				require.Equal(t, "TIMESTAMPTZ", column.SQLType)
			}
		}
		require.Equal(t, []string{"bucket_start", "record_id"}, key, "the unique key includes the partition column")
		_, err := ValidateGroupSchemaLayout(layout)
		require.NoError(t, err)
	})

	t.Run("a partition that is not exactly the documented one is refused", func(t *testing.T) {
		for name, mutate := range map[string]func(*GroupSchemaPartition, *GroupSchemaLayout){
			"another kind":      func(p *GroupSchemaPartition, _ *GroupSchemaLayout) { p.Kind = "native" },
			"another interval":  func(p *GroupSchemaPartition, _ *GroupSchemaLayout) { p.ChunkInterval = "1 day" },
			"a padded interval": func(p *GroupSchemaPartition, _ *GroupSchemaLayout) { p.ChunkInterval = "7 days " },
			"another column":    func(p *GroupSchemaPartition, _ *GroupSchemaLayout) { p.Column = "bucket_end" },
			"an unknown column": func(p *GroupSchemaPartition, _ *GroupSchemaLayout) { p.Column = "missing" },
			"an unsafe schema":  func(p *GroupSchemaPartition, _ *GroupSchemaLayout) { p.ExtensionSchema = `x"; DROP TABLE y; --` },
			"a nullable time column": func(_ *GroupSchemaPartition, l *GroupSchemaLayout) {
				l.Columns[0].Nullable, l.Columns[0].PrimaryKey = true, false
			},
		} {
			layout := timescaleLayout(t)
			partition := *layout.Partition
			layout.Partition = &partition
			layout.Columns = append([]GroupSchemaColumn(nil), layout.Columns...)
			mutate(&partition, &layout)
			_, err := ValidateGroupSchemaLayout(layout)
			require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout, name)
		}
	})

	t.Run("the DDL creates the table, makes it a hypertable empty, and nothing else", func(t *testing.T) {
		layout := timescaleLayout(t)
		layout.Partition.ExtensionSchema = "tsext"
		statements, err := GenerateGroupSchemaDDL(dialectPostgres, layout)
		require.NoError(t, err)
		require.Len(t, statements, 4, strings.Join(statements, "\n"))
		require.Contains(t, statements[0], `PRIMARY KEY ("bucket_start", "record_id")`)
		require.Equal(t,
			`SELECT "tsext"."create_hypertable"('"`+longTestTable+`"'::regclass, "tsext"."by_range"('bucket_start'::name, INTERVAL '7 days'), create_default_indexes => false, if_not_exists => false, migrate_data => false);`,
			statements[1], "the generalized API with explicit, non-faking arguments")
		require.True(t, strings.HasPrefix(statements[2], "CREATE INDEX"))
		require.Contains(t, statements[3], managedEffectReceiptMarker)
		for _, statement := range statements {
			require.NotRegexp(t, `(?i)create\s+extension|retention|compress|continuous|add_job|add_policy|partition\s+by|reorder`, statement,
				"no extension, policy, job or native partition is ever created")
		}

		plainStatements, err := GenerateGroupSchemaDDL(dialectPostgres, plain)
		require.NoError(t, err)
		require.Len(t, plainStatements, 3, "a plain layout keeps its statements")

		noSchema := timescaleLayout(t)
		_, err = GenerateGroupSchemaDDL(dialectPostgres, noSchema)
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout, "the extension's schema must come from the probe")
		_, err = GenerateGroupSchemaDDL(dialectSQLite, layout)
		require.Error(t, err, "SQLite has no hypertable")
	})

	t.Run("the layout identity ignores where the extension lives but not what is partitioned", func(t *testing.T) {
		left, right := timescaleLayout(t), timescaleLayout(t)
		right.Partition.ExtensionSchema = "tsext"
		require.True(t, groupSchemaLayoutsEqual(left, right))
		right.Partition.ChunkInterval = "1 day"
		require.False(t, groupSchemaLayoutsEqual(left, right))
		require.False(t, groupSchemaLayoutsEqual(left, plain), "plain and timescale are different layouts")
	})

	t.Run("preview freezes the previewed interval and the probed schema into the digest", func(t *testing.T) {
		svc, _ := newPreviewService()
		target := &timescaleTarget{layout: timescaleLayout(t), evidence: capableEvidence("tsext")}
		token, err := svc.PrepareSchemaPreview(t.Context(), timescaleScope(t), target.inspect)
		require.NoError(t, err)
		require.Equal(t, "tsext", token.GroupLayout.Partition.ExtensionSchema)
		require.Contains(t, token.Statements[1], `"tsext"."create_hypertable"`)
		require.Equal(t, &GroupSchemaPartition{Kind: PartitionKindTimescale, Column: "bucket_start", ChunkInterval: "7 days", ExtensionSchema: "tsext"}, token.Tables[0].Partition,
			"the interval is shown in the preview itself")

		other := &timescaleTarget{layout: timescaleLayout(t), evidence: capableEvidence("public")}
		moved, err := svc.PrepareSchemaPreview(t.Context(), timescaleScope(t), other.inspect)
		require.NoError(t, err)
		require.NotEqual(t, token.Digest, moved.Digest, "the extension's schema is part of what was confirmed")
	})

	t.Run("an unavailable capability blocks the preview with its reasons and previews nothing", func(t *testing.T) {
		svc, _ := newPreviewService()
		target := &timescaleTarget{layout: timescaleLayout(t), evidence: &TargetTimescaleEvidence{Reasons: []string{"extension-not-installed"}}}
		token, err := svc.PrepareSchemaPreview(t.Context(), timescaleScope(t), target.inspect)
		require.Nil(t, token)
		require.ErrorIs(t, err, ErrTimescaleUnavailable)
		var unavailable *TimescaleUnavailableError
		require.True(t, errors.As(err, &unavailable))
		require.Equal(t, []string{"extension-not-installed"}, unavailable.Reasons)

		missing := &timescaleTarget{layout: timescaleLayout(t)}
		_, err = svc.PrepareSchemaPreview(t.Context(), timescaleScope(t), missing.inspect)
		require.ErrorIs(t, err, ErrTargetInspectionUnconfirmed, "without capability evidence nothing is assumed")
	})

	t.Run("an existing table is never converted", func(t *testing.T) {
		svc, _ := newPreviewService()
		layout := timescaleLayout(t)
		// A plain table of the same name and columns: not a hypertable, so it is not ours to convert.
		plainTable := &timescaleTarget{layout: layout, table: existingLong(layout), evidence: capableEvidence("public")}
		_, err := svc.PrepareSchemaPreview(t.Context(), timescaleScope(t), plainTable.inspect)
		require.ErrorIs(t, err, ErrIncompatibleExistingTable)

		for name, mutate := range map[string]func(*TargetHypertable){
			"another interval":    func(h *TargetHypertable) { h.ChunkInterval = 24 * time.Hour },
			"another time column": func(h *TargetHypertable) { h.TimeColumn = "bucket_end" },
			"space dimension":     func(h *TargetHypertable) { h.Dimensions = 2 },
			"compression on":      func(h *TargetHypertable) { h.CompressionEnabled = true },
			"a background job":    func(h *TargetHypertable) { h.Jobs = 1 },
		} {
			hypertable := validHypertable()
			mutate(hypertable)
			evidence := capableEvidence("public")
			evidence.Hypertable = hypertable
			target := &timescaleTarget{layout: layout, table: existingLong(layout), evidence: evidence}
			_, err := svc.PrepareSchemaPreview(t.Context(), timescaleScope(t), target.inspect)
			require.ErrorIs(t, err, ErrIncompatibleExistingTable, name)
		}

		evidence := capableEvidence("public")
		evidence.Hypertable = validHypertable()
		owned := &timescaleTarget{layout: layout, table: existingLong(layout), evidence: evidence}
		token, err := svc.PrepareSchemaPreview(t.Context(), timescaleScope(t), owned.inspect)
		require.NoError(t, err)
		require.Equal(t, SchemaTableActionUnchanged, token.Tables[0].Action, "an owned hypertable that already matches needs no change")
		for _, statement := range token.Statements {
			require.NotContains(t, statement, "create_hypertable", "and is never asked to become one again")
		}
	})

	t.Run("confirmation re-checks the capability and refuses a moved extension", func(t *testing.T) {
		svc, _ := newPreviewService()
		layout := timescaleLayout(t)
		target := &timescaleTarget{layout: layout, evidence: capableEvidence("public")}
		token, err := svc.PrepareSchemaPreview(t.Context(), timescaleScope(t), target.inspect)
		require.NoError(t, err)
		_, err = svc.ValidatePreviewTokenForApply(t.Context(), token.Token, timescaleScope(t))
		require.NoError(t, err, "the token matches the scope although the scope's layout does not carry the probed schema")

		moved := &timescaleTarget{layout: layout, evidence: capableEvidence("elsewhere")}
		executed := 0
		_, _, err = svc.ApplySchemaPreview(t.Context(), token, SchemaApplyTarget{Inspect: moved.inspect, Execute: func(context.Context, []string) (SchemaExecution, error) {
			executed++
			return SchemaExecution{}, nil
		}})
		require.ErrorIs(t, err, ErrPreviewTokenStale)
		require.Zero(t, executed, "a changed destination is never executed against")

		gone := &timescaleTarget{layout: layout, evidence: &TargetTimescaleEvidence{Reasons: []string{"function-not-executable"}}}
		_, _, err = svc.ApplySchemaPreview(t.Context(), token, SchemaApplyTarget{Inspect: gone.inspect, Execute: func(context.Context, []string) (SchemaExecution, error) {
			executed++
			return SchemaExecution{}, nil
		}})
		require.ErrorIs(t, err, ErrTimescaleUnavailable)
		require.Zero(t, executed, "a lost capability blocks execution; nothing is changed")
	})

	t.Run("a verified apply records the partition in the operation's evidence", func(t *testing.T) {
		svc, _ := newPreviewService()
		layout := timescaleLayout(t)
		state := "missing"
		target := &timescaleTarget{layout: layout, evidence: capableEvidence("public")}
		target.table = func() TargetTableInspection {
			if state == "missing" {
				return TargetTableInspection{Status: inspectionMissing}
			}
			target.evidence = &TargetTimescaleEvidence{Available: true, ExtensionSchema: "public", Hypertable: validHypertable()}
			return existingLong(layout)()
		}
		token, err := svc.PrepareSchemaPreview(t.Context(), timescaleScope(t), target.inspect)
		require.NoError(t, err)
		var executed []string
		op, outcome, err := svc.ApplySchemaPreview(t.Context(), token, SchemaApplyTarget{Inspect: target.inspect, Execute: func(_ context.Context, statements []string) (SchemaExecution, error) {
			executed = statements
			state, target.receipts = "created", true
			return SchemaExecution{Committed: len(statements)}, nil
		}})
		require.NoError(t, err)
		require.Equal(t, ClaimAcquired, outcome)
		require.Equal(t, SchemaOperationSucceeded, op.Status)
		require.Equal(t, token.Statements, executed, "only the confirmed, stored statements run")
		require.NotEmpty(t, op.VerifiedDigest)
	})
}
