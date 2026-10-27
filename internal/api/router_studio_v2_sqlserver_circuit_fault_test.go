package api

import (
	"cmp"
	"context"
	"database/sql"
	"encoding/json"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/testutil/lostcommit"

	mssql "github.com/microsoft/go-mssqldb"
	"github.com/microsoft/go-mssqldb/msdsn"
	"github.com/stretchr/testify/require"
)

type sqlServerFaultDestinations struct {
	inner      *dbtarget.ConnectorService
	db         *sql.DB
	capability *dbtarget.SQLServerRecordingCapability
}

func (d sqlServerFaultDestinations) GetByID(ctx context.Context, id string) (*schema.DatabaseConnector, error) {
	return d.inner.GetByID(ctx, id)
}

func (d sqlServerFaultDestinations) OpenDestination(context.Context, string, string) (*dbtarget.OpenedDestination, error) {
	return &dbtarget.OpenedDestination{DB: d.db, Kind: "sqlserver", SQLServer: d.capability}, nil
}

// TestSQLServerCircuitPlainCommitAckLostRetriesOneEffect uses the existing
// lostcommit wrapper around the real SQL Server driver. Commit succeeds on the
// server and then reports lost response; the retry must not duplicate row or
// receipt, and the same effect must remain in the durable outbox.
func TestSQLServerCircuitPlainCommitAckLostRetriesOneEffect(t *testing.T) {
	fixture := SQLServerCircuitFixtureOrSkip(t)
	plain := newSQLServerCircuitPlainFixture(t, fixture)
	require.NoError(t, plain.pipe.Stop(time.Second))

	faultDB, faults := sqlServerCircuitFaultDB(t, fixture)
	t.Cleanup(func() { require.NoError(t, faultDB.Close()) })
	group, err := plain.groups.Get(t.Context(), plain.cumulativeID)
	require.NoError(t, err)
	capability, err := plain.targetService.ProbeSQLServerRecordingAtRevision(t.Context(),
		group.Group.Destination.ConnectorID, group.Group.Destination.ConnectorRevision)
	require.NoError(t, err)
	destinations := sqlServerFaultDestinations{inner: plain.targetService, db: faultDB, capability: capability}
	pipe := grouppipeline.New(grouppipeline.Dependencies{
		Groups: sqlServerCircuitOneGroup{plain.groups, plain.cumulativeID}, Tags: plain.tags, Destinations: destinations,
		Inspector: dbtarget.NewReadOnlyTableInspector(plain.targetService), Store: plain.store,
		Counters: grouppipeline.NewMeasurementCounterSources(measurement.NewSQLRepository(plain.db)),
	}, plain.pipeConfig)
	require.NoError(t, pipe.Start(t.Context()))
	t.Cleanup(func() { _ = pipe.Stop(time.Second) })
	plain.pipe = pipe

	capture := installSQLServerCircuitAckFault(t, faults, plain.db, plain.cumulativeID)
	t.Cleanup(func() {
		if !t.Failed() {
			return
		}
		var total, unfinished int
		_ = plain.db.QueryRowContext(context.Background(), `
			SELECT COUNT(*), COALESCE(SUM(CASE WHEN state <> 'sql_committed' THEN 1 ELSE 0 END), 0)
			FROM wg_delivery_outbox WHERE group_id = ?`, plain.cumulativeID).Scan(&total, &unfinished)
		var state, nextRetry, code, claimOwner, committedAt string
		var retries, claimEpoch int
		err := plain.db.QueryRowContext(context.Background(), `
			SELECT state, retry_count, next_retry_at, last_error_code, claim_owner, claim_epoch, committed_at
			FROM wg_delivery_outbox WHERE group_id = ?`, plain.cumulativeID).Scan(&state, &retries, &nextRetry, &code, &claimOwner, &claimEpoch, &committedAt)
		t.Logf("ACK-LOST FAILURE DIAGNOSTIC: query_err=%v total=%d unfinished=%d state=%q retries=%d next_retry=%q code=%q owner=%q epoch=%d committed_at=%q fired=%v",
			err, total, unfinished, state, retries, nextRetry, code, claimOwner, claimEpoch, committedAt, capture.fired.Load())
	})
	plain.accept(t, plain.cumulativeID, sqlServerCircuitMeasurementID, sqlServerCircuitInputs[0])
	plain.closeAt(t, 10*time.Minute)

	schemaName, table := group.Group.Destination.TableSchema, group.Group.Destination.TableName
	require.Eventually(t, func() bool {
		var total, unfinished int
		err := plain.db.QueryRowContext(t.Context(), `
			SELECT COUNT(*), SUM(CASE WHEN state <> 'sql_committed' THEN 1 ELSE 0 END)
			FROM wg_delivery_outbox WHERE group_id = ?`, plain.cumulativeID).Scan(&total, &unfinished)
		return err == nil && total == 1 && unfinished == 0
	}, 20*time.Second, 25*time.Millisecond, "ACK-lost retry must settle the one frozen effect")

	require.True(t, capture.fired.Load(), "the lost-ACK seam must actually fire")
	var recordID string
	require.NoError(t, plain.target.QueryRowContext(t.Context(), "SELECT record_id FROM "+sqlServerQualified(schemaName, table)+" WHERE workspace_id = @p1 AND group_id = @p2 AND group_revision = @p3",
		group.Group.WorkspaceID, group.Group.ID, group.Group.Revision).Scan(&recordID))

	var effectKey, payload, digest string
	require.NoError(t, plain.db.QueryRowContext(t.Context(), `
		SELECT effect_key, payload, payload_digest FROM wg_delivery_outbox WHERE group_id = ?`, plain.cumulativeID).Scan(&effectKey, &payload, &digest))
	require.NotEmpty(t, effectKey)
	require.NotEmpty(t, digest)
	require.True(t, json.Valid([]byte(payload)))
	require.NotContains(t, payload, "password")
	capture.mu.Lock()
	before, beforeErr := capture.payload, capture.err
	capture.mu.Unlock()
	require.NoError(t, beforeErr)
	require.NotEmpty(t, before, "the frozen outbox payload must be captured before retry")
	require.Equal(t, before, payload, "the ACK-lost retry must not rewrite the frozen payload")
	var receipts int
	require.NoError(t, plain.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM "+sqlServerQualified(schemaName, dbtarget.EffectReceiptTable)+" WHERE effect_key = @p1 AND payload_digest = @p2", effectKey, digest).Scan(&receipts))
	require.Equal(t, 1, receipts, "exactly the owned effect receipt must exist")
	row, err := groupdelivery.DecodeRowPayload([]byte(payload))
	require.NoError(t, err)
	cleanupSQLServerCircuitEffects(t, plain.target, schemaName, table, []string{recordID}, map[string]string{effectKey: digest}, map[string]dbtarget.EncodedRow{row.RecordID: row})
}

func sqlServerCircuitFaultDB(t *testing.T, fixture *SQLServerCircuitFixture) (*sql.DB, *lostcommit.Faults) {
	t.Helper()
	var config map[string]any
	require.NoError(t, json.Unmarshal([]byte(fixture.ConnectionConfigJSON()), &config))
	host := strings.TrimSpace(stringConfigValue(config, "host"))
	user := strings.TrimSpace(stringConfigValue(config, "user"))
	password, _ := config["password"].(string)
	port := cmp.Or(strings.TrimSpace(stringConfigValue(config, "port")), "1433")
	encrypt := "true"
	if option, present := config["encrypt"]; present {
		text := sqlServerCircuitOptionText(option)
		require.Contains(t, []string{"true", "mandatory", "yes", "1", "t", "strict"}, text)
		if text == "strict" {
			encrypt = "strict"
		}
	}
	for _, key := range []string{"trustservercertificate", "trust_server_certificate"} {
		if option, present := config[key]; present {
			trusted, err := strconv.ParseBool(sqlServerCircuitOptionText(option))
			require.True(t, err == nil && !trusted, "fault transport must preserve certificate verification")
		}
	}
	tlsMinimum := cmp.Or(stringConfigValue(config, "tlsmin"), "1.2")
	require.Contains(t, []string{"1.2", "1.3"}, tlsMinimum)
	timeout := 5 * time.Second
	if text := stringConfigValue(config, "timeout"); text != "" {
		var err error
		timeout, err = time.ParseDuration(text)
		require.True(t, err == nil && timeout > 0, "fault transport requires a bounded positive timeout")
	}
	seconds := strconv.FormatInt(int64((timeout-1)/time.Second)+1, 10)
	query := url.Values{
		"database": {fixture.Manifest.Database}, "encrypt": {encrypt},
		"TrustServerCertificate": {"false"}, "tlsmin": {tlsMinimum},
		"connection timeout": {seconds}, "dial timeout": {seconds},
	}
	if certificate, ok := config["certificate"].(string); ok && certificate != "" {
		query.Set("certificate", certificate)
	}
	if hostname, ok := config["hostnameincertificate"].(string); ok && hostname != "" {
		query.Set("hostnameincertificate", hostname)
	}
	dsn := (&url.URL{Scheme: "sqlserver", User: url.UserPassword(user, password), Host: net.JoinHostPort(strings.Trim(host, "[]"), port), RawQuery: query.Encode()}).String()
	parsed, err := msdsn.Parse(dsn)
	require.True(t, err == nil, "saved TLS connection settings cannot be parsed")
	require.NotNil(t, parsed.TLSConfig)
	require.False(t, parsed.TLSConfig.InsecureSkipVerify)
	require.False(t, parsed.TrustServerCertificate)
	require.Equal(t, fixture.Manifest.Database, parsed.Database)
	faults := &lostcommit.Faults{}
	return lostcommit.WrapFaults(&mssql.Driver{}, dsn, faults), faults
}

func sqlServerCircuitOptionText(value any) string {
	if boolean, ok := value.(bool); ok {
		return strconv.FormatBool(boolean)
	}
	text, _ := value.(string)
	return strings.ToLower(strings.TrimSpace(text))
}

// Filter only the worker list; snapshots still come from the public saved service.
type sqlServerCircuitOneGroup struct {
	*workspace.WriteGroupService
	id string
}

func (s sqlServerCircuitOneGroup) List(ctx context.Context) (*workspace.WriteGroupListResult, error) {
	listed, err := s.WriteGroupService.List(ctx)
	if err != nil {
		return nil, err
	}
	groups := make([]*workspace.WriteGroup, 0, 1)
	for _, group := range listed.Groups {
		if group.ID == s.id {
			groups = append(groups, group)
		}
	}
	listed.Groups = groups
	return listed, nil
}

type sqlServerCircuitAckCapture struct {
	fired   atomic.Bool
	mu      sync.Mutex
	payload string
	err     error
}

func installSQLServerCircuitAckFault(t *testing.T, faults *lostcommit.Faults, local *sql.DB, groupID string) *sqlServerCircuitAckCapture {
	t.Helper()
	capture := &sqlServerCircuitAckCapture{}
	faults.Lose.Store(true)
	hook := func() {
		if !capture.fired.CompareAndSwap(false, true) {
			return
		}
		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()
		var payload string
		err := local.QueryRowContext(ctx, `SELECT payload FROM wg_delivery_outbox WHERE group_id=?`, groupID).Scan(&payload)
		capture.mu.Lock()
		capture.payload, capture.err = payload, err
		capture.mu.Unlock()
	}
	faults.AfterCommit.Store(&hook)
	// Clearing Lose inside AfterCommit would cancel the fault itself. The next
	// SQL statement clears it, after Commit has already reported the lost ACK.
	nextStatement := func(string) error {
		if capture.fired.Load() {
			faults.Lose.Store(false)
		}
		return t.Context().Err()
	}
	faults.Reject.Store(&nextStatement)
	return capture
}
