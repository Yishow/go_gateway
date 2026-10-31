package runtime

import (
	"database/sql"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
	_ "modernc.org/sqlite"
	"path/filepath"
	"testing"
)

func TestModernSQLSourceRevisionMatchesActualManualRuntimeSample(t *testing.T) {
	f := newAcquisitionInstallFixture(t, true)
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "identity.db"))
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	ws := workspace.NewService(workspace.NewSQLRepository(db))
	record, err := ws.GetOrCreate(t.Context())
	require.NoError(t, err)
	d := f.projection.Devices[0]
	d.Name = "PLC"
	require.NoError(t, device.NewSQLRepository(db).Create(t.Context(), d))
	record, err = ws.AttachDevice(t.Context(), d.ID)
	require.NoError(t, err)
	p := f.projection.Points[0]
	p.Name = "value"
	p.DataType = schema.DataTypeUint32
	p.FormatMode = schema.PointFormatInherit
	p.FormatVersion = 1
	p.RevisionID = "point-r1"
	require.NoError(t, point.NewSQLRepository(db).Create(t.Context(), p))
	r := f.projection.Rules[0]
	r.StartAddress = p.Address
	r.Count = 1
	r.DataType = p.DataType
	r.DataFormat = "DCBA"
	r.Origin = "manual"
	require.NoError(t, sourcerule.NewSQLRepository(db).Create(t.Context(), r))
	tg := &schema.Tag{ID: "T", Key: "value", DisplayName: "Value", DataType: p.DataType, Status: schema.TagStatusActive, Labels: "{}"}
	require.NoError(t, tag.NewSQLRepository(db).Create(t.Context(), tg))
	m := &schema.Mapping{ID: "M", PointID: p.ID, TagID: tg.ID, TransformPipeline: "[]", Status: schema.MappingStatusActive, Enabled: true}
	require.NoError(t, mapping.NewSQLRepository(db).Create(t.Context(), m))
	_, err = db.ExecContext(t.Context(), `INSERT INTO source_rule_links(id,rule_id,address,point_id,tag_id,mapping_id) VALUES('L',?,?,?,?,?)`, r.ID, p.Address, p.ID, tg.ID, m.ID)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `INSERT INTO database_connectors(id,name,kind,connection_config,identity_revision,status,enabled) VALUES('C','target','sqlite','{"database":"target.db","schema":"main"}','c1','ready',1)`)
	require.NoError(t, err)
	// Use persisted readback, not a separately guessed source identity.
	d, err = device.NewSQLRepository(db).GetByID(t.Context(), d.ID)
	require.NoError(t, err)
	p, err = point.NewSQLRepository(db).GetByID(t.Context(), p.ID)
	require.NoError(t, err)
	tg, err = tag.NewSQLRepository(db).GetByID(t.Context(), tg.ID)
	require.NoError(t, err)
	m, err = mapping.NewSQLRepository(db).GetByID(t.Context(), m.ID)
	require.NoError(t, err)
	r, err = sourcerule.NewSQLRepository(db).GetByID(t.Context(), r.ID)
	require.NoError(t, err)
	links, err := sourcerule.NewSQLRepository(db).ListLinks(t.Context(), r.ID)
	require.NoError(t, err)
	f.projection.WorkspaceID = record.ID
	f.projection.DeviceIDs = []string{d.ID}
	f.projection.Devices = []*schema.Device{d}
	f.projection.Points = []*schema.Point{p}
	f.projection.Rules = []*schema.SourceRule{r}
	f.projection.RuleLinks = links
	f.projection.Tags = []*schema.Tag{tg}
	f.projection.Mappings = []*schema.Mapping{m}
	record.DatabaseSetupRevision = "database-setup-distinct-from-rule"
	require.NoError(t, workspace.NewSQLRepository(db).Save(t.Context(), record))
	f.reader.record = record
	sink := &sampleSinkCapture{}
	f.service.sampleSink = sink
	f.service.writer = &mockWriter{}
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	group := &workspace.WriteGroup{WorkspaceID: record.ID, Name: "format source", Members: []workspace.WriteGroupMember{{DeviceID: d.ID, PointID: p.ID, TagID: tg.ID, TargetColumn: "value", Required: true}}, Destination: workspace.WriteGroupDestination{ConnectorID: "C", ConnectorRevision: "c1", TableName: "samples", StorageStrategy: workspace.WriteGroupStorageStrategyCustom}}
	repo := workspace.NewSQLWriteGroupRepository(db)
	require.NoError(t, repo.Create(t.Context(), group))
	values := f.service.scheduler.PollNowContext(t.Context(), []string{p.ID})
	require.Len(t, values, 1)
	require.Equal(t, schema.QualityGood, values[0].Quality)
	f.service.handleCollectedValue(t.Context(), values[0])
	require.Len(t, sink.samples, 1)
	require.Equal(t, group.Members[0].SourceRevision, sink.samples[0].SourceRevision)
	require.Equal(t, values[0].ConfigFingerprint, sink.samples[0].ConfigFingerprint)
	observation := <-f.observations
	require.Equal(t, values[0].ConfigFingerprint, observation.Binding.ConfigFingerprint)
	// Reinstall same revision: old normalized value must not be reinterpreted or
	// accepted into the newly installed sample binding.
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	f.service.handleCollectedValue(t.Context(), values[0])
	require.Len(t, sink.samples, 1)
	frozen, err := json.Marshal(sink.samples[0])
	require.NoError(t, err)
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { manager.CloseAll() })
	f.service.scheduler = collector.NewScheduler(collector.DefaultSchedulerConfig(), manager)
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	f.service.handleCollectedValue(t.Context(), values[0])
	require.Len(t, sink.samples, 1, "old raw lease cannot survive scheduler restart")
	restarted := f.service.scheduler.PollNowContext(t.Context(), []string{p.ID})
	require.Len(t, restarted, 1)
	f.service.handleCollectedValue(t.Context(), restarted[0])
	require.Len(t, sink.samples, 2)
	require.Equal(t, sink.samples[0].SourceRevision, sink.samples[1].SourceRevision)
	require.Equal(t, sink.samples[0].ConfigFingerprint, sink.samples[1].ConfigFingerprint)
	replay, err := json.Marshal(sink.samples[0])
	require.NoError(t, err)
	require.Equal(t, string(frozen), string(replay), "accepted historical envelope is never reinterpreted")
	// Rule-only default changes stale the durable source revision with unchanged
	// point identity and DataFormat; caller-supplied old revision is rejected.
	_, err = db.ExecContext(t.Context(), `UPDATE source_rules SET data_format='ABCD',revision_id='rule-r2' WHERE id=?`, r.ID)
	require.NoError(t, err)
	stale := *group
	stale.ID = ""
	stale.Name = "stale"
	stale.Members = append([]workspace.WriteGroupMember(nil), group.Members...)
	require.ErrorIs(t, repo.Create(t.Context(), &stale), workspace.ErrWriteGroupSourceRevisionConflict)
}
