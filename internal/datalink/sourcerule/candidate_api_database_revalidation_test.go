package sourcerule

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_GetCandidateView_RevalidatesDatabaseScopeWhenConnectorBecomesReady(t *testing.T) {
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

	overrideTag, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "factory.db.recovered",
		DisplayName: "DB Recovered",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)
	svc.SetDatabaseTargetMappingReader(DatabaseTargetMappingListFunc(func(context.Context) ([]*schema.DatabaseTargetMapping, error) {
		return []*schema.DatabaseTargetMapping{
			{
				ID:          "db-map-recovered",
				TagID:       overrideTag.ID,
				ConnectorID: "connector-recovered",
				TableSchema: "public",
				TableName:   "measurements",
				ColumnName:  "line_a",
				WriteMode:   schema.DatabaseWriteModeInsert,
			},
		}, nil
	}))

	validation := &DatabaseTargetConnectorValidation{
		Ready: false,
		Issues: []DatabaseTargetValidationIssue{
			{
				Severity:  "error",
				MappingID: "db-map-recovered",
				Code:      "table_missing",
				Message:   "找不到資料表: public.measurements",
			},
		},
	}
	svc.SetDatabaseTargetConnectorValidator(DatabaseTargetConnectorValidatorFunc(func(context.Context, string) (*DatabaseTargetConnectorValidation, error) {
		return validation, nil
	}))

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-db-revalidate-service")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-db-revalidate-service",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	tagCandidates := tagCandidatesFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	require.Len(t, tagCandidates, 1)
	_, err = svc.UpsertTagReviewDecision(ctx, rule.ID, UpsertTagReviewDecisionRequest{
		CandidateID:   tagCandidates[0].ID,
		Action:        schema.SourceRuleTagReviewDecisionActionOverride,
		OverrideTagID: overrideTag.ID,
	})
	require.NoError(t, err)

	snapshots, err := svc.ListCandidateSnapshots(ctx, rule.ID)
	require.NoError(t, err)
	databaseCandidates, databaseSnapshot := databaseCandidatesFromSnapshots(t, snapshots)
	assert.Equal(t, schema.SourceRuleCandidateStatusBlocked, databaseSnapshot.Status)
	assert.Contains(t, databaseSnapshot.Reason, "找不到資料表")
	require.Len(t, databaseCandidates, 1)
	assert.Equal(t, "out_of_sync", databaseCandidates[0]["status"])

	validation = &DatabaseTargetConnectorValidation{
		Ready:  true,
		Issues: []DatabaseTargetValidationIssue{},
	}

	view, err := svc.GetCandidateView(ctx, rule.ID)
	require.NoError(t, err)
	assert.Equal(t, rule.ID, view.SourceRuleID)
	assert.Equal(t, rule.RevisionID, view.RevisionID)
	assert.Equal(t, schema.SourceRuleCandidateStatusReady, view.DatabaseOutputs.Status)
	assert.Empty(t, view.DatabaseOutputs.Reason)
	require.Len(t, view.DatabaseOutputs.Candidates, 1)

	candidate, ok := view.DatabaseOutputs.Candidates[0].(schema.SourceRuleDatabaseOutputCandidate)
	require.True(t, ok)
	assert.Equal(t, schema.SourceRuleOutputStatusReady, candidate.Status)
	assert.Empty(t, candidate.BlockingReason)
	assert.Equal(t, "connector-recovered", candidate.ConnectorID)
	assert.Equal(t, "public", candidate.TableSchema)
	assert.Equal(t, "measurements", candidate.TableName)
	assert.Equal(t, "line_a", candidate.ColumnName)

	refreshedSnapshots, err := svc.ListCandidateSnapshots(ctx, rule.ID)
	require.NoError(t, err)
	databaseCandidates, databaseSnapshot = databaseCandidatesFromSnapshots(t, refreshedSnapshots)
	assert.Equal(t, schema.SourceRuleCandidateStatusReady, databaseSnapshot.Status)
	assert.Empty(t, databaseSnapshot.Reason)
	require.Len(t, databaseCandidates, 1)
	assert.Equal(t, "ready", databaseCandidates[0]["status"])
	assert.Empty(t, databaseCandidates[0]["blocking_reason"])
}
