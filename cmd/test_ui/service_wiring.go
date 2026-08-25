package main

import (
	"context"
	"database/sql"
	"log"

	"go-gateway/internal/datalink/audit"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/settings"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
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
	audit        *audit.Service
	dbTarget     *dbtarget.ConnectorService
	dbMapping    *dbtarget.MappingService
	scheduler    *collector.Scheduler
	runtime      *datalinkruntime.Service
	sourceRule   *sourcerule.Service
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
	auditSvc := audit.NewService(audit.NewSQLRepository(db))
	dbTargetConnectorRepo := dbtarget.NewSQLConnectorRepository(db)
	dbTargetMappingRepo := dbtarget.NewSQLTargetMappingRepository(db)
	dbTargetConnectorSvc := dbtarget.NewConnectorService(dbTargetConnectorRepo, dbTargetMappingRepo)
	dbTargetConnectorSvc.SetTagReader(tagSvc)
	dbTargetMappingSvc := dbtarget.NewMappingService(dbTargetMappingRepo, dbTargetConnectorRepo, tagSvc)
	dbTargetWriter := dbtarget.NewWriterWithConfig(dbTargetConnectorRepo, dbTargetMappingRepo, dbtarget.WriterConfig{TagReader: tagSvc})
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), connMgr)
	runtimeWriter := storage.NewBatchWriter(storage.NewSQLiteWriter(db), storage.DefaultBatchWriterConfig())
	runtimeSvc, err := datalinkruntime.NewService(datalinkruntime.DefaultConfig(), datalinkruntime.Dependencies{
		Scheduler: scheduler, Writer: runtimeWriter, TargetWriter: newProductionTargetWriter(dbTargetWriter, modbusShareSvc),
		DeviceService: devSvc, PointService: pointSvc, MappingService: mappingSvc, TagService: tagSvc, PollingGroupService: pgSvc,
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
	return gatewayServices{device: devSvc, pollingGroup: pgSvc, point: pointSvc, tag: tagSvc, mapping: mappingSvc, settings: settingsSvc, modbusShare: modbusShareSvc, shareLoaded: shareSettingsLoaded, workspace: workspaceSvc, audit: auditSvc, dbTarget: dbTargetConnectorSvc, dbMapping: dbTargetMappingSvc, scheduler: scheduler, runtime: runtimeSvc, sourceRule: sourceRuleSvc}
}
