package sourcerule

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_LocalModbusConflictRefresh_PreservesDatabaseSnapshot(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-lm-db-isolation")
	require.NoError(t, err)

	ruleA, linksA := createLocalModbusConflictRule(t, ctx, svc, repo, dev.ID, "rule-lm-db-isolation-a", "40001", schema.DataTypeInt32)
	ruleB, linksB := createLocalModbusConflictRule(t, ctx, svc, repo, dev.ID, "rule-lm-db-isolation-b", "40101", schema.DataTypeInt16)

	require.NotNil(t, linksA[0].TagID)
	require.NotNil(t, linksB[0].TagID)

	svc.SetLocalModbusMappingReader(LocalModbusMappingListFunc(func(context.Context) ([]LocalModbusMappingRecord, error) {
		return []LocalModbusMappingRecord{
			{TagID: *linksA[0].TagID, Register: 10, DataType: schema.DataTypeInt32},
			{TagID: *linksB[0].TagID, Register: 11, DataType: schema.DataTypeInt16},
		}, nil
	}))

	require.NoError(t, svc.persistCandidateSnapshots(ctx, ruleA, linksA))
	require.NoError(t, svc.persistCandidateSnapshots(ctx, ruleB, linksB))

	databasePayloadA := `{"candidates":[{"id":"db-preserved-a","tag_key":"DB_PRESERVED_A","connector_id":"connector-preserved-a","table_name":"measurements","column_name":"line_a","status":"ready"}]}`
	databasePayloadB := `{"candidates":[{"id":"db-preserved-b","tag_key":"DB_PRESERVED_B","connector_id":"connector-preserved-b","table_name":"measurements","column_name":"line_b","status":"ready"}]}`
	setDatabaseSnapshotForRule(t, repo, ruleA.ID, ruleA.RevisionID, databasePayloadA, schema.SourceRuleCandidateStatusReady, "")
	setDatabaseSnapshotForRule(t, repo, ruleB.ID, ruleB.RevisionID, databasePayloadB, schema.SourceRuleCandidateStatusReady, "")

	beforeRuleA := databaseSnapshotForRule(t, repo, ruleA.ID, ruleA.RevisionID)
	beforeRuleB := databaseSnapshotForRule(t, repo, ruleB.ID, ruleB.RevisionID)

	require.NoError(t, svc.refreshLocalModbusConflictSnapshots(ctx))

	assertDatabaseSnapshotForRule(t, repo, ruleA.ID, ruleA.RevisionID, databasePayloadA, schema.SourceRuleCandidateStatusReady, "", beforeRuleA.GeneratedAt)
	assertDatabaseSnapshotForRule(t, repo, ruleB.ID, ruleB.RevisionID, databasePayloadB, schema.SourceRuleCandidateStatusReady, "", beforeRuleB.GeneratedAt)

	candidateA, snapshotA := localModbusSingleCandidateFromRule(t, ctx, svc, ruleA.ID)
	assert.Equal(t, schema.SourceRuleCandidateStatusDeferred, snapshotA.Status)
	assert.Equal(t, schema.SourceRuleLocalModbusOutputStatusBlockedConflict, candidateA.Status)
	assert.NotEmpty(t, candidateA.BlockingReason)
}

func TestService_ApplyLocalModbusOutputCandidates_PreservesDatabaseSnapshot(t *testing.T) {
	t.Parallel()

	ctx, svc, rule := createLocalModbusApplyValidationRule(t, "device-lm-apply-db-isolation", "rule-lm-apply-db-isolation", "40081")
	repo := svc.repo.(*MemoryRepository)

	databasePayload := `{"candidates":[{"id":"db-apply-preserved","tag_key":"DB_APPLY_PRESERVED","connector_id":"connector-apply-preserved","table_name":"measurements","column_name":"line_a","status":"blocked","blocking_reason":"database scope pending"}]}`
	setDatabaseSnapshotForRule(t, repo, rule.ID, rule.RevisionID, databasePayload, schema.SourceRuleCandidateStatusBlocked, "preserved database apply state")
	beforeSnapshot := databaseSnapshotForRule(t, repo, rule.ID, rule.RevisionID)

	localModbusPayload := `{"candidates":[{"id":"lm-apply-preserved","tag_id":"tag-apply-preserved","tag_key":"LM_APPLY_PRESERVED","register":9}]}`
	setLocalModbusSnapshotForRule(t, repo, rule.ID, rule.RevisionID, localModbusPayload, schema.SourceRuleCandidateStatusReady, "")

	response, err := svc.ApplyLocalModbusOutputCandidates(ctx, rule.ID, ApplyOutputCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{"lm-apply-preserved"},
	})
	require.NoError(t, err)
	require.Len(t, response.Results, 1)
	assert.Equal(t, "success", response.Results[0].Status)
	assert.Equal(t, "lm-apply-preserved", response.Results[0].CandidateID)

	assertDatabaseSnapshotForRule(t, repo, rule.ID, rule.RevisionID, databasePayload, schema.SourceRuleCandidateStatusBlocked, "preserved database apply state", beforeSnapshot.GeneratedAt)
}

func setDatabaseSnapshotForRule(
	t *testing.T,
	repo *MemoryRepository,
	ruleID string,
	revisionID string,
	payload string,
	status schema.SourceRuleCandidateStatus,
	reason string,
) {
	t.Helper()

	snapshots, err := repo.ListCandidateSnapshots(context.Background(), ruleID, revisionID)
	require.NoError(t, err)

	nextSnapshots := make([]*schema.SourceRuleCandidateSnapshot, 0, len(snapshots))
	for _, snapshot := range snapshots {
		candidateSnapshot := *snapshot
		if candidateSnapshot.CandidateType == schema.SourceRuleCandidateTypeDatabaseOutputs {
			candidateSnapshot.Payload = payload
			candidateSnapshot.Status = status
			candidateSnapshot.Reason = reason
		}
		nextSnapshots = append(nextSnapshots, &candidateSnapshot)
	}
	require.NoError(t, repo.ReplaceCandidateSnapshots(context.Background(), nextSnapshots))
}

func assertDatabaseSnapshotForRule(
	t *testing.T,
	repo *MemoryRepository,
	ruleID string,
	revisionID string,
	payload string,
	status schema.SourceRuleCandidateStatus,
	reason string,
	generatedAt time.Time,
) {
	t.Helper()

	snapshot := databaseSnapshotForRule(t, repo, ruleID, revisionID)
	assert.Equal(t, status, snapshot.Status)
	assert.Equal(t, reason, snapshot.Reason)
	assert.JSONEq(t, payload, snapshot.Payload)
	assert.True(t, snapshot.GeneratedAt.Equal(generatedAt))
}

func databaseSnapshotForRule(
	t *testing.T,
	repo *MemoryRepository,
	ruleID string,
	revisionID string,
) *schema.SourceRuleCandidateSnapshot {
	t.Helper()

	snapshots, err := repo.ListCandidateSnapshots(context.Background(), ruleID, revisionID)
	require.NoError(t, err)

	for _, snapshot := range snapshots {
		if snapshot.CandidateType != schema.SourceRuleCandidateTypeDatabaseOutputs {
			continue
		}
		cloned := *snapshot
		return &cloned
	}

	t.Fatalf("database output snapshot not found for rule %s revision %s", ruleID, revisionID)
	return nil
}
