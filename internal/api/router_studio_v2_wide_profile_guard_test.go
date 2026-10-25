package api

import (
	"database/sql"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

// TestWideProfileRefusesAnIncompatibleTarget pins that a wide profile group is
// never prepared on a table that cannot hold its rows: an existing table whose
// member column is NOT NULL cannot store a window where that member has no value,
// so the preview is refused with the schema-incompatible code and the table is
// left exactly as it was. Nothing is altered to make it fit.
func TestWideProfileRefusesAnIncompatibleTarget(t *testing.T) {
	for _, kind := range wideKinds() {
		t.Run(kind, func(t *testing.T) {
			f := newManagedDestinationFixture(t, kind)
			sources := []meanSource{{name: "half", tagType: "int16", role: "gauge"}}
			seedMeanSources(t, f, sources)
			created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, wideProfileRequest(t, f, sources))
			require.Equal(t, http.StatusOK, created.Code, created.Body.String())
			_, saved := groupSaveData(t, decodeJSONBody(t, created))
			groupID, table := saved["id"].(string), destinationOf(saved)["table_name"].(string)
			column := saved["members"].([]any)[0].(map[string]any)["target_column"].(string)

			// The table the preview itself would create, except that the member
			// column cannot hold NULL.
			statements := previewLong(t, f, groupID).Statements
			memberType := map[string]string{"sqlite": "TEXT", "postgres": "NUMERIC"}[kind]
			definition := `"` + column + `" ` + memberType
			require.Contains(t, statements[0], definition, "the member column is typed from the frozen decimal result")
			create := strings.Replace(statements[0], definition, definition+" NOT NULL", 1)
			if kind == "postgres" {
				create = strings.Replace(create, `CREATE TABLE "`+table+`"`, `CREATE TABLE `+qualified(kind, f.namespace, table), 1)
			}
			target := f.target
			if kind == "sqlite" {
				target = openSQLiteTarget(t, f.file)
			}
			_, err := target.ExecContext(t.Context(), create)
			require.NoError(t, err)

			response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/schema-preview", scopeBody(t, f, groupID))
			require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
			assertWriteGroupError(t, decodeJSONBody(t, response), "RECORDING_SCHEMA_INCOMPATIBLE")

			var notNull bool
			if kind == "sqlite" {
				require.NoError(t, target.QueryRowContext(t.Context(), `SELECT "notnull" = 1 FROM pragma_table_info(?) WHERE name = ?`, table, column).Scan(&notNull))
			} else {
				var nullable string
				require.NoError(t, target.QueryRowContext(t.Context(),
					`SELECT is_nullable FROM information_schema.columns WHERE table_schema = $1 AND table_name = $2 AND column_name = $3`,
					f.namespace, table, column).Scan(&nullable))
				notNull = nullable == "NO"
			}
			require.True(t, notNull, "the incompatible table is never altered to fit")
		})
	}
}

// TestWideProfileKeepsUnservedShapesClosed pins the shapes this layout does not
// serve: a cumulative wide profile whose source has no confirmed counter is saved
// at most as a draft that never prepares, never becomes ready and never applies
// (the confirmed one is served: TestWideCounterParity), and a storage mode (the
// TimescaleDB mode of the long layout) is refused when saving.
func TestWideProfileKeepsUnservedShapesClosed(t *testing.T) {
	t.Run("an unconfirmed cumulative wide profile never prepares, becomes ready or applies", func(t *testing.T) {
		f := newManagedDestinationFixture(t, "sqlite")
		seedMeanSources(t, f, []meanSource{{name: "energy", tagType: "uint64", role: "gauge"}})
		request := profileRequest(t, f.writeGroupRouterFixture, "cumulative", 10, "minute", "main", "energy")
		group := request["group"].(map[string]any)
		group["layout_version"] = wideProfileLayoutVersion
		group["destination"].(map[string]any)["table_schema"] = f.namespace
		created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
		require.Equal(t, http.StatusOK, created.Code, created.Body.String())
		_, saved := groupSaveData(t, decodeJSONBody(t, created))
		groupID := saved["id"].(string)
		preview := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/schema-preview", scopeBody(t, f, groupID))
		require.Equal(t, http.StatusUnprocessableEntity, preview.Code, preview.Body.String())
		ready, err := f.groups.Readiness(t.Context(), groupID)
		require.NoError(t, err)
		require.False(t, ready.Ready)
		applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/apply", scopeBody(t, f, groupID))
		require.Equal(t, http.StatusConflict, applied.Code, applied.Body.String())
	})

	t.Run("a storage mode is refused on the wide profile layout", func(t *testing.T) {
		f := newManagedDestinationFixture(t, "sqlite")
		sources := []meanSource{{name: "half", tagType: "int16", role: "gauge"}}
		seedMeanSources(t, f, sources)
		request := wideProfileRequest(t, f, sources)
		request["group"].(map[string]any)["storage_mode"] = "timescale"
		created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
		require.Equal(t, http.StatusUnprocessableEntity, created.Code, created.Body.String())
		assertWriteGroupError(t, decodeJSONBody(t, created), "WRITE_GROUP_INVALID")
		require.Zero(t, countProfileGroups(t, f.writeGroupRouterFixture))
	})

	t.Run("an explicit skip_row policy is refused; required members skip instead", func(t *testing.T) {
		f := newManagedDestinationFixture(t, "sqlite")
		sources := []meanSource{{name: "half", tagType: "int16", role: "gauge"}}
		seedMeanSources(t, f, sources)
		request := wideProfileRequest(t, f, sources)
		request["group"].(map[string]any)["row_policy"].(map[string]any)["incomplete_policy"] = "skip_row"
		created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
		require.Equal(t, http.StatusUnprocessableEntity, created.Code, created.Body.String())
		require.Zero(t, countProfileGroups(t, f.writeGroupRouterFixture))
	})
}

// TestWideProfileKeepsTheSingleOwnerOfItsScope pins that the layout never makes a
// second writer: a source and circuit already owned by a wide profile group
// cannot also be recorded by a long group of the same profile kind.
func TestWideProfileKeepsTheSingleOwnerOfItsScope(t *testing.T) {
	f := newManagedDestinationFixture(t, "sqlite")
	sources := []meanSource{{name: "half", tagType: "int16", role: "gauge"}}
	seedMeanSources(t, f, sources)
	created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, wideProfileRequest(t, f, sources))
	require.Equal(t, http.StatusOK, created.Code, created.Body.String())

	second := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, meanLongRequest(t, f, sources))
	require.Equal(t, http.StatusConflict, second.Code, second.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, second), "WRITE_GROUP_CIRCUIT_PROFILE_CONFLICT")
	require.Equal(t, 1, countProfileGroups(t, f.writeGroupRouterFixture), "no second writer is created from another layout")
}

// TestLegacyWideKeepsItsNoDataAndPolicy is the control for the new layout: a
// legacy managed wide group whose window saw no sample at all still writes no
// row and records no_data, and its saved row policy is not given the wide
// profile's partial default.
func TestLegacyWideKeepsItsNoDataAndPolicy(t *testing.T) {
	for _, kind := range wideKinds() {
		t.Run(kind, func(t *testing.T) {
			f := newManagedGroupFixture(t, kind)
			groupID, table := f.group["id"].(string), destinationOf(f.group)["table_name"].(string)
			_, explicit := f.group["row_policy"].(map[string]any)["incomplete_policy"]
			require.False(t, explicit, "a legacy wide group is not given a policy it never chose")

			require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, groupID, previewLong(t, f, groupID)).Status)
			ready, err := f.groups.Readiness(t.Context(), groupID)
			require.NoError(t, err)
			require.True(t, ready.Ready, "%+v", ready.Issues)
			applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/apply", scopeBody(t, f, groupID))
			require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
			snapshot, err := f.groups.ResolveAppliedAt(t.Context(), groupID, time.Now().Add(time.Minute))
			require.NoError(t, err)

			rig := &wideRig{f: f, kind: kind, groupID: groupID, table: table, snapshot: snapshot}
			rig.clock = newInjectedClock(snapshot.EffectiveAt.Add(time.Second))
			rig.pipe = newMeanPipeline(t, f, rig.clock, "node/test", true)
			require.Equal(t, grouppipeline.StateActive, pipeStatus(rig.pipe, groupID).State, "%+v", rig.pipe.Status())
			rig.closeWindow()
			logKind, reason := rig.bucketLog(t)
			require.Equal(t, [2]string{"no_data", "no-samples"}, [2]string{logKind, reason})
			target := f.target
			if kind == "sqlite" {
				target = openSQLiteTarget(t, f.file)
			}
			require.Zero(t, rig.rowCount(t, target), "a silent legacy window writes no row")
		})
	}
}

func openSQLiteTarget(t *testing.T, path string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func newInjectedClock(at time.Time) *atomic.Int64 {
	var clock atomic.Int64
	clock.Store(at.UnixNano())
	return &clock
}
