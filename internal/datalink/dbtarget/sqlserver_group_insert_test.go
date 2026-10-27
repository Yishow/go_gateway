package dbtarget

import (
	"database/sql"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSQLServerGroupInsertAcceptsJSONRoundTrippedProof(t *testing.T) {
	req := sqlServerGroupInsertTestRequest()
	payload, err := json.Marshal(req.SQLServer)
	require.NoError(t, err)
	var frozen SQLServerManagedRowProof
	require.NoError(t, json.Unmarshal(payload, &frozen))
	require.False(t, frozen.Group.Inspection.IndexesInspected, "this transient flag is excluded from proof JSON")
	require.NotNil(t, frozen.Group.Indexes, "independent indexes preserve inspected-empty versus absent")
	require.NotNil(t, frozen.Receipt.Indexes)
	req.SQLServer = &frozen

	require.NoError(t, validateGroupInsert(&sql.DB{}, req))
}

func TestSQLServerGroupInsertRejectsIncompleteOrMismatchedProof(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*GroupInsertRequest)
	}{
		{name: "missing proof", mutate: func(req *GroupInsertRequest) { req.SQLServer = nil }},
		{name: "missing connector revision", mutate: func(req *GroupInsertRequest) { req.SQLServer.ConnectorRevision = " " }},
		{name: "missing group inspection", mutate: func(req *GroupInsertRequest) { req.SQLServer.Group = nil }},
		{name: "missing receipt inspection", mutate: func(req *GroupInsertRequest) { req.SQLServer.Receipt = nil }},
		{name: "schema mismatch", mutate: func(req *GroupInsertRequest) { req.SchemaName = "other" }},
		{name: "group table mismatch", mutate: func(req *GroupInsertRequest) { req.TableName = "other_table" }},
		{name: "receipt table mismatch", mutate: func(req *GroupInsertRequest) { req.SQLServer.Receipt.Inspection.Table = "foreign_receipts" }},
		{name: "unverified metadata", mutate: func(req *GroupInsertRequest) { req.SQLServer.Group.Evidence.MetadataVerified = false }},
		{name: "non-plain group", mutate: func(req *GroupInsertRequest) { req.SQLServer.Group.Evidence.Plain = false }},
		{name: "missing object identity", mutate: func(req *GroupInsertRequest) { req.SQLServer.Group.Evidence.ObjectID = 0 }},
		{name: "independent indexes absent", mutate: func(req *GroupInsertRequest) { req.SQLServer.Group.Indexes = nil }},
		{name: "unexpected group index", mutate: func(req *GroupInsertRequest) {
			req.SQLServer.Group.Indexes = append(req.SQLServer.Group.Indexes, IndexInfo{Name: "foreign_index"})
		}},
		{name: "wrong owner default", mutate: func(req *GroupInsertRequest) {
			wrong := "foreign"
			req.SQLServer.Group.Evidence.ColumnDefaults[sqlServerInsertTestOwner] = &wrong
		}},
		{name: "wrong collation", mutate: func(req *GroupInsertRequest) {
			req.SQLServer.Group.Evidence.ColumnCollations[longColRecordID] = "Latin1_General_100_CI_AS"
		}},
		{name: "non-receipt strategy", mutate: func(req *GroupInsertRequest) { req.Strategy = GroupEffectNone }},
		{name: "non-hash digest", mutate: func(req *GroupInsertRequest) { req.PayloadDigest = "digest" }},
		{name: "non-hash effect key", mutate: func(req *GroupInsertRequest) { req.Row.EffectKey = "effect" }},
		{name: "non-hash record id", mutate: func(req *GroupInsertRequest) { req.Row.RecordID = "record" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := sqlServerGroupInsertTestRequest()
			test.mutate(&req)
			require.ErrorIs(t, validateGroupInsert(&sql.DB{}, req), ErrGroupInsertUnsupported)
		})
	}
}

func TestSQLServerGroupInsertRejectsMalformedCells(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*GroupInsertRequest)
	}{
		{name: "duplicate cell", mutate: func(req *GroupInsertRequest) {
			req.Row.Cells = append(req.Row.Cells, req.Row.Cells[0])
		}},
		{name: "unknown cell", mutate: func(req *GroupInsertRequest) {
			req.Row.Cells[0].Column = "foreign_column"
		}},
		{name: "missing required cell", mutate: func(req *GroupInsertRequest) {
			req.Row.Cells = req.Row.Cells[1:]
		}},
		{name: "foreign owner", mutate: func(req *GroupInsertRequest) {
			req.Row.Cells[len(req.Row.Cells)-1].Value = "foreign"
		}},
		{name: "record identity differs from row id", mutate: func(req *GroupInsertRequest) {
			cell := sqlServerTestCell(req.Row.Cells, longColRecordID)
			cell.Value = strings.Repeat("d", 64)
		}},
		{name: "bucket start differs from row", mutate: func(req *GroupInsertRequest) {
			cell := sqlServerTestCell(req.Row.Cells, longColBucketStart)
			cell.Value = req.Row.BucketStart.Add(time.Microsecond)
		}},
		{name: "submicrosecond bucket start", mutate: func(req *GroupInsertRequest) {
			cell := sqlServerTestCell(req.Row.Cells, longColBucketStart)
			cell.Value = req.Row.BucketStart.Add(time.Nanosecond)
		}},
		{name: "decimal must remain exact text", mutate: func(req *GroupInsertRequest) {
			sqlServerTestCell(req.Row.Cells, longColValueDecimal).Value = 1.25
		}},
		{name: "invalid UTF-8 text", mutate: func(req *GroupInsertRequest) {
			sqlServerTestCell(req.Row.Cells, longColTagID).Value = string([]byte{0xff})
		}},
		{name: "empty required identity", mutate: func(req *GroupInsertRequest) {
			sqlServerTestCell(req.Row.Cells, longColTagID).Value = ""
		}},
		{name: "empty required provenance", mutate: func(req *GroupInsertRequest) {
			sqlServerTestCell(req.Row.Cells, longColProvenance).Value = ""
		}},
		{name: "non-JSON provenance", mutate: func(req *GroupInsertRequest) {
			sqlServerTestCell(req.Row.Cells, longColProvenance).Value = "not-json"
		}},
		{name: "invalid committed timestamp", mutate: func(req *GroupInsertRequest) {
			req.CommittedAt = "2026-01-02 03:04:05"
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := sqlServerGroupInsertTestRequest()
			test.mutate(&req)
			require.ErrorIs(t, validateGroupInsert(&sql.DB{}, req), ErrGroupInsertUnsupported)
		})
	}
}

func sqlServerTestCell(cells []EncodedCell, name string) *EncodedCell {
	for index := range cells {
		if cells[index].Column == name {
			return &cells[index]
		}
	}
	panic("test fixture is missing cell " + name)
}
