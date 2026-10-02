package main

import (
	"context"
	"database/sql"
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
}

func wireGatewayServices(db *sql.DB, connMgr *connector.ConnectionManager) gatewayServices {
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
	shareSettingsLoaded := loadPersistedShareSettings(context.Background(), modbusShareSvc, settingsRepo)
	workspaceSvc := workspace.NewService(workspace.NewSQLRepository(db))
	writeGroupsSvc := workspace.NewWriteGroupService(workspaceSvc, workspace.NewSQLWriteGroupRepository(db)).
		WithBacklogOwnershipGuard(grouppipeline.BacklogGuard{})
	auditSvc := audit.NewService(audit.NewSQLRepository(db))
	dbTargetConnectorRepo := dbtarget.NewSQLConnectorRepository(db)
	dbTargetMappingRepo := dbtarget.NewSQLTargetMappingRepository(db)
	dbTargetConnectorSvc := dbtarget.NewConnectorService(dbTargetConnectorRepo, dbTargetMappingRepo)
	dbTargetConnectorSvc.SetTagReader(tagSvc)
	writeGroupsSvc.WithTableInspector(dbtarget.NewReadOnlyTableInspector(dbTargetConnectorSvc))
	dbTargetMappingSvc := dbtarget.NewMappingService(dbTargetMappingRepo, dbTargetConnectorRepo, tagSvc)
	dbtarget.WithLegacyWriteCoordinator(dbTargetMappingSvc, writeGroupsSvc)
	dbTargetConnectorSvc.SetDeleteGuard(writeGroupsSvc)
	groupStore, err := groupdelivery.NewStore(db).WithQuota(groupdelivery.QuotaConfig{GlobalMaxBytes: groupDeliveryQuotaBytes})
	if err != nil {
		log.Fatalf("建立群組交付容量設定失敗: %v", err)
	}
	groupPipe := grouppipeline.New(grouppipeline.Dependencies{
		Groups: writeGroupsSvc, Tags: tagSvc, Destinations: dbTargetConnectorSvc,
		Inspector: dbtarget.NewReadOnlyTableInspector(dbTargetConnectorSvc), Store: groupStore,
	}, grouppipeline.Config{NodeID: gatewayNodeID(), Owner: gatewayNodeID() + "/" + uuid.NewString(), OnError: logGroupPipelineError})
	// A canonical group that owns an output replaces the legacy writer for it, so
	// there is exactly one writer per output.
	dbTargetWriter := dbtarget.NewWriterWithConfig(dbTargetConnectorRepo, dbTargetMappingRepo, dbtarget.WriterConfig{TagReader: tagSvc, SuppressMapping: groupPipe.Owns})
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), connMgr)
	runtimeWriter := storage.NewBatchWriter(storage.NewSQLiteWriter(db), storage.DefaultBatchWriterConfig())
	runtimeSvc, err := datalinkruntime.NewService(datalinkruntime.DefaultConfig(), datalinkruntime.Dependencies{
		Scheduler: scheduler, Writer: runtimeWriter, TargetWriter: newProductionTargetWriter(dbTargetWriter, modbusShareSvc),
		DeviceService: devSvc, PointService: pointSvc, MappingService: mappingSvc, TagService: tagSvc, PollingGroupService: pgSvc,
		SampleSink: groupPipe,
	})
	if err != nil {
		log.Fatalf("建立 datalink runtime 失敗: %v", err)
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
	if err := sourceRuleSvc.SyncDerivedPointState(context.Background()); err != nil {
		log.Printf("同步來源規則衍生點位狀態失敗: %v", err)
	}
	// Explicit group test writes share the recording operation ledger and write
	// through the production row codec; they never mix with the delivery outbox.
	recordingPlanSvc := recordingplan.NewService(recordingplan.NewSQLRepository(db))
	groupTestWriteSvc := grouptestwrite.New(grouptestwrite.Dependencies{
		Groups: writeGroupsSvc, Tags: tagSvc, Destinations: dbTargetConnectorSvc,
		Inspector: dbtarget.NewReadOnlyTableInspector(dbTargetConnectorSvc), Ledger: recordingPlanSvc,
	}, grouptestwrite.Config{})
	return gatewayServices{recordingPlan: recordingPlanSvc, groupTestWrite: groupTestWriteSvc, device: devSvc, pollingGroup: pgSvc, point: pointSvc, tag: tagSvc, mapping: mappingSvc, settings: settingsSvc, modbusShare: modbusShareSvc, shareLoaded: shareSettingsLoaded, workspace: workspaceSvc, writeGroups: writeGroupsSvc, groupPipe: groupPipe, audit: auditSvc, dbTarget: dbTargetConnectorSvc, dbMapping: dbTargetMappingSvc, scheduler: scheduler, runtime: runtimeSvc, sourceRule: sourceRuleSvc}
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

func logGroupPipelineError(err error) {
	log.Printf("寫入群組 pipeline 錯誤: %v", err)
}
