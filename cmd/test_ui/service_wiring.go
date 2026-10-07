package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"go-gateway/internal/datalink/audit"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/grouptestwrite"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/recordingplan"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/settings"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/diagnostics"

	"github.com/google/uuid"
)

type gatewayServices struct {
	device       *device.Service
	pollingGroup *pollinggroup.Service
	point        *point.Service
	tag          *tag.Service
	mapping      *mapping.Service
	settings     *settings.Service
	modbusShare  *modbusshare.Service
	shareLoaded  bool
	workspace    *workspace.Service
	writeGroups  *workspace.WriteGroupService
	groupPipe    *grouppipeline.Pipeline
	// recordingPlan holds the operation ledger shared by schema applies and
	// explicit group test writes.
	recordingPlan  *recordingplan.Service
	groupTestWrite *grouptestwrite.Service
	audit          *audit.Service
	dbTarget       *dbtarget.ConnectorService
	dbMapping      *dbtarget.MappingService
	scheduler      *collector.Scheduler
	runtime        *datalinkruntime.Service
	sourceRule     *sourcerule.Service
	fixture        groupFixtureHooks
}

func wireGatewayServices(db *sql.DB, connMgr *connector.ConnectionManager) (gatewayServices, error) {
	return wireGatewayServicesContext(context.Background(), db, connMgr, nil)
}
func wireGatewayServicesContext(ctx context.Context, db *sql.DB, connMgr *connector.ConnectionManager, logs *diagnostics.Broker) (gatewayServices, error) {
	fixture := newGroupFixture()
	devRepo := device.NewSQLRepository(db)
	devSvc := device.NewService(devRepo, connMgr)
	pgRepo := pollinggroup.NewSQLRepository(db)
	pgSvc := pollinggroup.NewService(pgRepo)
	pointRepo := point.NewSQLRepository(db)
	pointSvc := point.NewService(pointRepo, pgRepo)
	tagRepo := tag.NewSQLRepository(db)
	tagSvc := tag.NewService(tagRepo)
	mappingRepo := mapping.NewSQLRepository(db)
	mappingSvc := mapping.NewServiceWithTagResolver(mappingRepo, tagSvc.GetByID)
	sourceRuleRepo := sourcerule.NewSQLRepository(db)
	settingsRepo := settings.NewSQLRepository(db)
	settingsSvc := settings.NewService(settingsRepo)
	modbusShareSvc := modbusshare.NewService(tagSvc, 65536)
	shareSettingsLoaded := loadPersistedShareSettings(ctx, modbusShareSvc, settingsRepo)
	workspaceSvc := workspace.NewService(workspace.NewSQLRepository(db))
	writeGroupsSvc := workspace.NewWriteGroupService(workspaceSvc, workspace.NewSQLWriteGroupRepository(db)).
		WithBacklogOwnershipGuard(grouppipeline.BacklogGuard{})
	auditSvc := audit.NewService(audit.NewSQLRepository(db))
	dbTargetConnectorRepo := dbtarget.NewSQLConnectorRepository(db)
	dbTargetMappingRepo := dbtarget.NewSQLTargetMappingRepository(db)
	dbTargetConnectorSvc := dbtarget.NewConnectorService(dbTargetConnectorRepo, dbTargetMappingRepo)
	dbTargetConnectorSvc.SetTagReader(tagSvc)
	writeGroupsSvc.WithTableInspector(dbtarget.NewReadOnlyTableInspector(dbTargetConnectorSvc))
	writeGroupsSvc.WithManagedTableInspector(dbtarget.NewManagedTableInspector(dbTargetConnectorSvc, db))
	dbTargetMappingSvc := dbtarget.NewMappingService(dbTargetMappingRepo, dbTargetConnectorRepo, tagSvc)
	dbtarget.WithLegacyWriteCoordinator(dbTargetMappingSvc, writeGroupsSvc)
	dbTargetConnectorSvc.SetDeleteGuard(writeGroupsSvc)
	groupStore, err := groupdelivery.NewStore(db).WithQuota(groupDeliveryQuotaConfig())
	if err != nil {
		return gatewayServices{}, fmt.Errorf("configure group delivery capacity: %w", err)
	}
	pipelineConfig := fixture.pipelineConfig(gatewayNodeID(), uuid.NewString())
	if logs != nil {
		previous := pipelineConfig.OnError
		pipelineConfig.OnError = func(err error) {
			if err != nil {
				logs.Emit(diagnostics.Input{Code: "runtime.group_failed"})
				if previous != nil {
					previous(err)
				}
			}
		}
	}
	groupPipe := grouppipeline.New(grouppipeline.Dependencies{
		Groups: writeGroupsSvc, Tags: tagSvc, Destinations: dbTargetConnectorSvc,
		Inspector: dbtarget.NewReadOnlyTableInspector(dbTargetConnectorSvc), Store: groupStore,
	}, pipelineConfig)
	sampleSink := fixture.sampleSink(groupPipe)
	// A canonical group that owns an output replaces the legacy writer for it, so
	// there is exactly one writer per output.
	dbTargetWriter := dbtarget.NewWriterWithConfig(dbTargetConnectorRepo, dbTargetMappingRepo, dbtarget.WriterConfig{TagReader: tagSvc, SuppressMapping: groupPipe.Owns}) //nolint:contextcheck // Persistent worker is owned by runtime shutdown.
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), connMgr)
	runtimeWriter := storage.NewBatchWriter(storage.NewSQLiteWriter(db), storage.DefaultBatchWriterConfig()) //nolint:contextcheck // Persistent timer is joined by runtime shutdown.
	runtimeSvc, err := datalinkruntime.NewService(datalinkruntime.DefaultConfig(), datalinkruntime.Dependencies{
		Scheduler: scheduler, Writer: runtimeWriter, TargetWriter: newProductionTargetWriter(dbTargetWriter, modbusShareSvc),
		DeviceService: devSvc, PointService: pointSvc, MappingService: mappingSvc, TagService: tagSvc, PollingGroupService: pgSvc,
		SampleSink: sampleSink,
	})
	if err != nil {
		targetErr := dbTargetWriter.Close(ctx)
		writerErr := runtimeWriter.CloseContext(ctx)
		waitErr := runtimeWriter.WaitClosed(context.WithoutCancel(ctx))
		return gatewayServices{}, fmt.Errorf("create datalink runtime: %w", errors.Join(err, targetErr, writerErr, waitErr))
	}
	sourceRuleSvc := sourcerule.NewService(sourceRuleRepo, devSvc, pointSvc, runtimeSvc)
	sourceRuleSvc.SetTagMappingServices(tagSvc, mappingSvc)
	sourceRuleSvc.SetDatabaseTargetMappingReader(sourcerule.DatabaseTargetMappingListFunc(func(ctx context.Context) ([]*schema.DatabaseTargetMapping, error) {
		return dbTargetMappingSvc.List(ctx, dbtarget.TargetMappingListFilter{})
	}))
	sourceRuleSvc.SetDatabaseTargetConnectorReader(dbTargetConnectorSvc)
	sourceRuleSvc.SetDatabaseTargetConnectorValidator(sourcerule.DatabaseTargetConnectorValidatorFunc(func(ctx context.Context, connectorID string) (*sourcerule.DatabaseTargetConnectorValidation, error) {
		result, err := dbTargetMappingSvc.Validate(ctx, connectorID)
		if err != nil {
			return nil, err
		}
		issues := make([]sourcerule.DatabaseTargetValidationIssue, 0, len(result.Issues))
		for _, issue := range result.Issues {
			issues = append(issues, sourcerule.DatabaseTargetValidationIssue{Severity: issue.Severity, MappingID: issue.MappingID, Code: issue.Code, Message: issue.Message})
		}
		return &sourcerule.DatabaseTargetConnectorValidation{Ready: result.Ready, Issues: issues}, nil
	}))
	workspaceSvc.WithReadinessServices(devSvc, sourceRuleSvc, dbTargetConnectorSvc, dbTargetMappingSvc)
	workspaceSvc.WithReadinessSetupReaders(tagSvc, mappingSvc)
	workspaceSvc.WithWriteGroupReadiness(writeGroupsSvc)
	workspaceSvc.WithRuntimeProjectionServices(devSvc, sourceRuleSvc, pointSvc, mappingSvc, tagSvc, dbTargetConnectorSvc, dbTargetMappingSvc, pgSvc)
	runtimeSvc.SetWorkspaceProjectionReader(workspaceSvc)
	sourceRuleSvc.SetLocalModbusMappingReader(sourcerule.LocalModbusMappingListFunc(func(context.Context) ([]sourcerule.LocalModbusMappingRecord, error) {
		mappings := modbusShareSvc.ListMappings()
		records := make([]sourcerule.LocalModbusMappingRecord, 0, len(mappings))
		for _, mappingRecord := range mappings {
			records = append(records, sourcerule.LocalModbusMappingRecord{MappingID: mappingRecord.MappingID, TagID: mappingRecord.TagID, Register: mappingRecord.Register, DataType: mappingRecord.DataType, UpdatedAt: mappingRecord.UpdatedAt})
		}
		return records, nil
	}))
	if err := sourceRuleSvc.SyncDerivedPointState(ctx); err != nil {
		log.Printf("同步來源規則衍生點位狀態失敗: %v", err)
	}
	// Explicit group test writes share the recording operation ledger and write
	// through the production row codec; they never mix with the delivery outbox.
	recordingPlanSvc := recordingplan.NewService(recordingplan.NewSQLRepository(db))
	groupTestWriteSvc := grouptestwrite.New(grouptestwrite.Dependencies{
		Groups: writeGroupsSvc, Tags: tagSvc, Destinations: dbTargetConnectorSvc,
		Inspector: dbtarget.NewReadOnlyTableInspector(dbTargetConnectorSvc), Ledger: recordingPlanSvc,
	}, grouptestwrite.Config{})
	return gatewayServices{recordingPlan: recordingPlanSvc, groupTestWrite: groupTestWriteSvc, device: devSvc, pollingGroup: pgSvc, point: pointSvc, tag: tagSvc, mapping: mappingSvc, settings: settingsSvc, modbusShare: modbusShareSvc, shareLoaded: shareSettingsLoaded, workspace: workspaceSvc, writeGroups: writeGroupsSvc, groupPipe: groupPipe, audit: auditSvc, dbTarget: dbTargetConnectorSvc, dbMapping: dbTargetMappingSvc, scheduler: scheduler, runtime: runtimeSvc, sourceRule: sourceRuleSvc, fixture: fixture}, nil
}

// groupDeliveryQuotaBytes caps accepted-but-undelivered group data, matching the
// legacy delivery queue default. It is a configured limit, not a measured
// capacity: intake is refused (never accepted data dropped) once it is reached.
const groupDeliveryQuotaBytes int64 = 500 * 1024 * 1024

// gatewayNodeID identifies this installation so a restart can tell its own
// earlier incarnation's delivery claims from another node's.
func gatewayNodeID() string {
	if host, err := os.Hostname(); err == nil && host != "" {
		return host
	}
	return "gateway"
}
