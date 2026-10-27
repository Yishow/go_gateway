package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

const (
	mysqlCircuitCumulativeWindow = 10 * time.Minute
	mysqlCircuitFamilyMySQL      = "mysql"
	mysqlCircuitFamilyMariaDB    = "mariadb"
)

type mysqlCircuitInput struct {
	sampleID string
	offset   time.Duration
	value    uint64
}

var mysqlCircuitInputs = [...]mysqlCircuitInput{
	{sampleID: "mysql-uint64-baseline", offset: 9 * time.Minute, value: 1 << 53},
	{sampleID: "mysql-uint64-after-gap", offset: 29 * time.Minute, value: (1 << 53) + 1},
	{sampleID: "mysql-uint64-max", offset: 39 * time.Minute, value: ^uint64(0)},
	{sampleID: "mysql-uint64-reset", offset: 49 * time.Minute, value: 7},
}

type mysqlCircuitPlainFixtureOptions struct {
	Family          string
	SeedMeasurement measurement.MeasurementDefinition
	Inputs          []mysqlCircuitInput
	StorageMode     workspace.WriteGroupStorageMode
	BaseTime        time.Time
	OpenFailures    int32
	CanonicalApply  bool
}

type mysqlCircuitProfile struct {
	id       string
	group    *workspace.WriteGroup
	snapshot *workspace.WriteGroupAppliedSnapshot
	key      groupdelivery.GroupKey
}

type mysqlCircuitJournalEntry struct {
	groupID, sampleID, payload, digest string
	consumed                           int
}

type mysqlCircuitPlainFixture struct {
	t               *testing.T
	f               managedGroupFixture
	table           string
	base            time.Time
	instant         mysqlCircuitProfile
	cumulative      mysqlCircuitProfile
	definition      *measurement.MeasurementDefinition
	inputs          []mysqlCircuitInput
	store           *groupdelivery.Store
	pipe            *grouppipeline.Pipeline
	faultTargets    *mysqlCircuitFaultDestinations
	acceptedJournal []mysqlCircuitJournalEntry
	clock           atomic.Int64
}

// mysqlCircuitPlainGroupSource is the fixture-only Apply bridge. Canonical
// groups remain confirmed-but-unapplied in the database; only this GroupSource
// returns the immutable snapshots seeded after the real API and catalog proofs.
type mysqlCircuitPlainGroupSource struct {
	service   *workspace.WriteGroupService
	snapshots map[string]*workspace.WriteGroupAppliedSnapshot
}

func (s *mysqlCircuitPlainGroupSource) List(ctx context.Context) (*workspace.WriteGroupListResult, error) {
	listed, err := s.service.List(ctx)
	if err != nil {
		return nil, err
	}
	for i, group := range listed.Groups {
		snapshot, ok := s.snapshots[group.ID]
		if !ok {
			continue
		}
		applied := *group
		applied.AppliedRevision = snapshot.AppliedRevision
		applied.Status = workspace.WriteGroupStatusReady
		listed.Groups[i] = &applied
	}
	return listed, nil
}

func (s *mysqlCircuitPlainGroupSource) ResolveAppliedAt(_ context.Context, id string, at time.Time) (*workspace.WriteGroupAppliedSnapshot, error) {
	snapshot, ok := s.snapshots[id]
	if !ok || at.Before(snapshot.EffectiveAt) {
		return nil, workspace.ErrWriteGroupAppliedRevisionUnavailable
	}
	return snapshot, nil
}

// mysqlCircuitFaultDestinations delegates to the real connector service after
// a small, deterministic number of transient open failures for retry coverage.
type mysqlCircuitFaultDestinations struct {
	*dbtarget.ConnectorService
	failOpen atomic.Int32
}

func (d *mysqlCircuitFaultDestinations) OpenDestination(ctx context.Context, connectorID, expectedRevision string) (*dbtarget.OpenedDestination, error) {
	for {
		remaining := d.failOpen.Load()
		if remaining <= 0 {
			return d.ConnectorService.OpenDestination(ctx, connectorID, expectedRevision)
		}
		if d.failOpen.CompareAndSwap(remaining, remaining-1) {
			return nil, errors.New("fixture transient destination-open failure")
		}
	}
}

func newMySQLCircuitPlainFixture(t *testing.T, options ...mysqlCircuitPlainFixtureOptions) *mysqlCircuitPlainFixture {
	t.Helper()
	require.LessOrEqual(t, len(options), 1, "at most one fixture option set is supported")
	var option mysqlCircuitPlainFixtureOptions
	if len(options) == 1 {
		option = options[0]
	}
	seedMeasurement := measurement.MeasurementDefinition{
		ID: "meas-energy", DefinitionRevision: "def-5", SourceBindingRevision: "bind-5", SeriesEpoch: "epoch-3",
		Name: "Energy", Quantity: "energy", Unit: "kWh", SemanticKind: measurement.SemanticKindCounter,
		NumericEncoding: "uint64", ExactRecordingScale: option.SeedMeasurement.ExactRecordingScale,
	}
	inputs := slices.Clone(mysqlCircuitInputs[:])
	if option.Inputs != nil {
		inputs = slices.Clone(option.Inputs)
	}
	require.NotEmpty(t, inputs, "the recording fixture needs at least one sample")

	family := option.Family
	if family == "" {
		family = mysqlCircuitFamilyMySQL
	}
	require.Contains(t, []string{mysqlCircuitFamilyMySQL, mysqlCircuitFamilyMariaDB}, family,
		"fixture family must be mysql or mariadb")
	f := newMySQLSchemaRouterFixture(t, family, "uint64")
	_, err := f.db.ExecContext(t.Context(), `
		INSERT INTO measurement_definitions (
			id, workspace_id, device_id, point_id, tag_id, equipment_id, definition_revision, source_binding_revision,
			series_epoch, name, quantity, unit, semantic_kind, numeric_encoding, exact_recording_scale
		) VALUES (?, ?, 'device-A', 'point-A', 'tag-A', 'device-A', 'def-5', 'bind-5', 'epoch-3', 'Energy', 'energy', 'kWh', 'counter', ?, NULLIF(?, ''))`,
		seedMeasurement.ID, f.record.ID, seedMeasurement.NumericEncoding, seedMeasurement.ExactRecordingScale)
	require.NoError(t, err)
	definition, err := measurement.NewSQLRepository(f.db).GetByID(t.Context(), seedMeasurement.ID)
	require.NoError(t, err)
	require.Equal(t, seedMeasurement.NumericEncoding, definition.NumericEncoding)
	require.Equal(t, seedMeasurement.ExactRecordingScale, definition.ExactRecordingScale)
	require.Equal(t, [4]string{"def-5", "bind-5", "epoch-3", "uint64"}, [4]string{
		definition.DefinitionRevision, definition.SourceBindingRevision, definition.SeriesEpoch, definition.NumericEncoding,
	})

	mode := option.StorageMode
	if mode == "" {
		mode = workspace.WriteGroupStorageModePlain
	}
	instantID := ensureMySQLCircuitPlainProfile(t, f, "instantaneous", false, definition.ID, mode)
	cumulativeID := ensureMySQLCircuitPlainProfile(t, f, "cumulative", true, definition.ID, mode)
	instantGroup, err := f.groups.Get(t.Context(), instantID)
	require.NoError(t, err)
	cumulativeGroup, err := f.groups.Get(t.Context(), cumulativeID)
	require.NoError(t, err)
	table := instantGroup.Group.Destination.TableName
	require.Equal(t, table, cumulativeGroup.Group.Destination.TableName, "profiles for one circuit/source share the managed table")
	require.Equal(t, "tag-A", instantGroup.Group.Members[0].TagID)
	require.Equal(t, instantGroup.Group.Members[0].TagID, cumulativeGroup.Group.Members[0].TagID)
	instantMember, cumulativeMember := instantGroup.Group.Members[0], cumulativeGroup.Group.Members[0]
	require.Equal(t, workspace.WriteGroupValueRoleCounterRaw, instantMember.ValueRole)
	require.Equal(t, schema.DataTypeUint64, instantMember.SourceValueType)
	require.NotNil(t, instantMember.MeasurementID)
	require.Equal(t, definition.ID, *instantMember.MeasurementID)
	require.Empty(t, cumulativeMember.ValueRole, "cumulative profiles carry counter semantics without an instantaneous role")
	require.NotNil(t, cumulativeMember.MeasurementID)
	require.Equal(t, definition.ID, *cumulativeMember.MeasurementID)
	if seedMeasurement.ExactRecordingScale == "" {
		require.Equal(t, measurement.ExactUint64, instantMember.ResultValueType)
		require.Nil(t, instantMember.ExactRecordingSource)
		require.Nil(t, cumulativeMember.ExactRecordingSource)
	} else {
		for _, member := range []workspace.WriteGroupMember{instantMember, cumulativeMember} {
			proof := member.ExactRecordingSource
			require.NotNil(t, proof, "the API must derive the frozen source proof from the saved measurement")
			require.Equal(t, definition.ID, proof.MeasurementID)
			require.Equal(t, definition.DefinitionRevision, proof.DefinitionRevision)
			require.Equal(t, definition.SourceBindingRevision, proof.SourceBindingRevision)
			require.Equal(t, definition.SeriesEpoch, proof.SeriesEpoch)
			require.Equal(t, definition.ExactRecordingScale, proof.Scale)
			require.Equal(t, schema.DataTypeUint64, proof.OriginalDataType)
			require.Equal(t, schema.DataTypeUint64, member.SourceValueType)
			require.Equal(t, measurement.ExactDecimal, member.ResultValueType)
		}
	}

	var existing int
	require.NoError(t, f.target.QueryRowContext(t.Context(), `
		SELECT COUNT(*) FROM information_schema.TABLES
		WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME IN (?,?)`, table, dbtarget.EffectReceiptTable).Scan(&existing))
	require.Zero(t, existing, "the fixture must own both MySQL objects before confirmation")
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		for _, name := range []string{table, dbtarget.EffectReceiptTable} {
			require.True(t, longTablePattern.MatchString(name) || name == dbtarget.EffectReceiptTable)
			_, err := f.target.ExecContext(ctx, "DROP TABLE IF EXISTS `"+name+"`")
			require.NoError(t, err)
		}
	})

	base := time.Now().UTC().Truncate(mysqlCircuitCumulativeWindow).Add(mysqlCircuitCumulativeWindow)
	if !option.BaseTime.IsZero() {
		base = option.BaseTime.UTC()
	}
	if option.CanonicalApply {
		// Align the real Apply versions with the fixture's first window, so
		// earlier empty windows do not become part of the tested batch.
		f.groups.WithClock(func() time.Time { return base.Add(-time.Second) })
	}

	for i, id := range []string{instantID, cumulativeID} {
		token := previewLong(t, f, id)
		if i == 0 {
			require.Len(t, token.Statements, 3, "the first group confirms the table, index, and receipt table")
		} else {
			require.Empty(t, token.Statements, "the second profile confirms its own shared-table scope without inheriting the first group's confirmation")
		}
		require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, id, token).Status)
		readiness, err := f.groups.Readiness(t.Context(), id)
		require.NoError(t, err)
		require.True(t, readiness.SchemaReady, "%+v", readiness.Issues)
		if option.CanonicalApply {
			require.True(t, readiness.Ready, "%+v", readiness.Issues)
			applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
			require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
		} else {
			require.True(t, readiness.Ready, "verified operation gates require complete schema proof: %+v", readiness.Issues)
			unapplied, err := f.groups.Get(t.Context(), id)
			require.NoError(t, err)
			require.Empty(t, unapplied.Group.AppliedRevision, "the database remains unapplied")
		}
	}

	// Schema confirmation updates each group's destination lineage. Use the
	// authoritative post-confirmation revisions and digests in runtime proofs.
	instantGroup, err = f.groups.Get(t.Context(), instantID)
	require.NoError(t, err)
	cumulativeGroup, err = f.groups.Get(t.Context(), cumulativeID)
	require.NoError(t, err)

	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
	inspector := f.groups.ManagedTableInspector(targets)
	capability, err := inspector.ProbeMySQLRecordingAtRevision(t.Context(), instantGroup.Group.Destination.ConnectorID, instantGroup.Group.Destination.ConnectorRevision)
	require.NoError(t, err)
	tags := tag.NewService(tag.NewSQLRepository(f.db))
	store := groupdelivery.NewStore(f.db)
	profiles := make([]mysqlCircuitProfile, 0, 2)
	for _, saved := range []*workspace.WriteGroupSaveResult{instantGroup, cumulativeGroup} {
		group := *saved.Group
		group.AppliedRevision = group.Revision
		group.Status = workspace.WriteGroupStatusReady
		physical, err := inspector.InspectMySQLManagedTableAtRevision(t.Context(), capability.ConnectorID, capability.ConnectorRevision,
			capability.Database, capability.CapabilityRevision, table)
		require.NoError(t, err)
		receipt, err := inspector.InspectMySQLManagedTableAtRevision(t.Context(), capability.ConnectorID, capability.ConnectorRevision,
			capability.Database, capability.CapabilityRevision, dbtarget.EffectReceiptTable)
		require.NoError(t, err)
		tagTypes := make(map[string]schema.DataType, len(group.Members))
		for _, member := range group.Members {
			source, err := tags.GetByID(t.Context(), member.TagID)
			require.NoError(t, err)
			tagTypes[member.TagID] = source.DataType
		}
		frozen := &workspace.WriteGroupRuntimeLayout{
			Dialect: dbtarget.SQLDialectMySQL, Columns: physical.Inspection.Columns, TagTypes: tagTypes,
			ReceiptTableReady: true, SchemaDigest: group.Destination.SchemaDigest,
			MySQL: &dbtarget.MySQLManagedRowProof{
				ConnectorID: capability.ConnectorID, ConnectorRevision: capability.ConnectorRevision,
				Group: physical, Receipt: receipt,
			},
		}
		require.NoError(t, workspace.ValidateMySQLRuntimeLayout(&group, frozen))
		snapshot := &workspace.WriteGroupAppliedSnapshot{
			WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.Revision,
			AppliedRevision: group.Revision, EffectiveAt: base, Group: &group, RuntimeLayout: frozen,
		}
		key := groupdelivery.GroupKey{WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.Revision}
		if option.CanonicalApply {
			snapshot, err = f.groups.ResolveAppliedAt(t.Context(), group.ID, base)
			require.NoError(t, err)
			group = *snapshot.Group
		} else {
			payload, err := json.Marshal(snapshot)
			require.NoError(t, err)
			require.NoError(t, store.SaveRuntimeVersion(t.Context(), key, payload))
		}
		profiles = append(profiles, mysqlCircuitProfile{id: group.ID, group: &group, snapshot: snapshot, key: key})
	}

	groupSource := &mysqlCircuitPlainGroupSource{service: f.groups, snapshots: make(map[string]*workspace.WriteGroupAppliedSnapshot, len(profiles))}
	for _, profile := range profiles {
		groupSource.snapshots[profile.id] = profile.snapshot
	}
	var activeGroups grouppipeline.GroupSource = groupSource
	if option.CanonicalApply {
		activeGroups = f.groups
	}
	faultTargets := &mysqlCircuitFaultDestinations{ConnectorService: targets}
	faultTargets.failOpen.Store(2)
	if option.OpenFailures > 0 {
		faultTargets.failOpen.Store(option.OpenFailures)
	}
	fixture := &mysqlCircuitPlainFixture{
		t: t, f: f, table: table, base: base, definition: definition, inputs: inputs,
		instant: profiles[0], cumulative: profiles[1], store: store,
		faultTargets: faultTargets,
	}
	fixture.clock.Store(base.Add(time.Second).UnixNano())
	fixture.pipe = grouppipeline.New(grouppipeline.Dependencies{
		Groups: activeGroups, Tags: tags, Destinations: faultTargets,
		Inspector: dbtarget.NewReadOnlyTableInspector(targets), Store: store,
		Counters: grouppipeline.NewMeasurementCounterSources(measurement.NewSQLRepository(f.db)),
	}, grouppipeline.Config{
		NodeID: "mysql-circuit-plain", Owner: "mysql-circuit-plain/test",
		Now:          func() time.Time { return time.Unix(0, fixture.clock.Load()).UTC() },
		TickInterval: time.Hour, WorkerInterval: 20 * time.Millisecond, ReconcileInterval: time.Hour,
		MaxFutureSkew: time.Minute,
		Sender: groupdelivery.SenderConfig{Backoff: func(int) (string, time.Time) {
			return groupdelivery.StateRetrying, time.Now().UTC().Add(3 * time.Second)
		}},
	})
	require.NoError(t, fixture.pipe.Start(t.Context()))
	pipe := fixture.pipe
	t.Cleanup(func() { _ = pipe.Stop(time.Second) })
	for _, profile := range profiles {
		require.Equal(t, grouppipeline.StateActive, pipeStatus(fixture.pipe, profile.id).State, "%+v", fixture.pipe.Status())
	}
	return fixture
}

func ensureMySQLCircuitPlainProfile(t *testing.T, f managedGroupFixture, kind string, cumulative bool, measurementID string, mode workspace.WriteGroupStorageMode) string {
	t.Helper()
	interval, unit := any(15), "second"
	if cumulative {
		interval, unit = 10, "minute"
	}
	request := profileRequest(t, f.writeGroupRouterFixture, kind, interval, unit, "mysql-uint64-main", "")
	group := request["group"].(map[string]any)
	group["layout_version"], group["storage_mode"] = "managed_long_v1", string(mode)
	group["destination"].(map[string]any)["table_schema"] = f.namespace
	member := group["members"].([]map[string]any)[0]
	member["measurement_id"] = measurementID
	if !cumulative {
		member["value_role"] = string(workspace.WriteGroupValueRoleCounterRaw)
	}
	response := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, response))
	return saved["id"].(string)
}

func (f *mysqlCircuitPlainFixture) accept(index int) {
	f.t.Helper()
	input := f.inputs[index]
	member := f.instant.group.Members[0]
	observed := f.base.Add(input.offset)
	received := observed.Add(30 * time.Millisecond)
	f.clock.Store(observed.Add(time.Second).UnixNano())
	require.NoError(f.t, f.pipe.AcceptSample(f.t.Context(), measurement.SampleEnvelope{
		SampleID: input.sampleID, WorkspaceID: f.f.record.ID, DeviceID: member.DeviceID,
		PointID: member.PointID, TagID: member.TagID,
		SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision,
		MeasurementID: f.definition.ID, DefinitionRevision: f.definition.DefinitionRevision,
		SourceBindingRevision: f.definition.SourceBindingRevision, SeriesEpoch: f.definition.SeriesEpoch,
		AcquisitionID: input.sampleID + "-acquisition", ValueType: string(measurement.ExactUint64),
		ObservedAt: observed, ReceivedAt: received, TimeOrigin: "gateway", Value: input.value, Quality: schema.QualityGood,
	}))
	// Capture the actual durable intake before closing the window. The worker
	// may reclaim consumed journal entries after a later TickAll or restart.
	rows, err := f.f.db.QueryContext(f.t.Context(), `SELECT group_id, sample_id, payload, payload_digest, consumed
		FROM wg_delivery_samples WHERE group_id IN (?, ?) AND sample_id=?`, f.instant.id, f.cumulative.id, input.sampleID)
	require.NoError(f.t, err)
	defer rows.Close()
	var accepted []mysqlCircuitJournalEntry
	for rows.Next() {
		var entry mysqlCircuitJournalEntry
		require.NoError(f.t, rows.Scan(&entry.groupID, &entry.sampleID, &entry.payload, &entry.digest, &entry.consumed))
		require.Zero(f.t, entry.consumed)
		accepted = append(accepted, entry)
	}
	require.NoError(f.t, rows.Err())
	require.Len(f.t, accepted, 2)
	f.acceptedJournal = append(f.acceptedJournal, accepted...)
}

func (f *mysqlCircuitPlainFixture) closeAfter(offset time.Duration) {
	f.t.Helper()
	f.clock.Store(f.base.Add(offset + time.Second).UnixNano())
	f.pipe.TickAll(f.t.Context())
}

var _ grouppipeline.GroupSource = (*mysqlCircuitPlainGroupSource)(nil)
var _ grouppipeline.Destinations = (*mysqlCircuitFaultDestinations)(nil)
