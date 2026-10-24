package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"maps"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// TimescaleChunkInterval is the chunk interval a timescale layout previews and
// freezes. It is a design default, not a throughput promise.
const TimescaleChunkInterval = 7 * 24 * time.Hour

// TimescalePartitionColumn is the only time dimension of the long layout.
const TimescalePartitionColumn = "bucket_start"

// Stable reason codes of an unusable TimescaleDB capability. They name a
// missing fact, never a value, DSN or driver text.
const (
	TimescaleReasonDialect               = "dialect-unsupported"
	TimescaleReasonNotInstalled          = "extension-not-installed"
	TimescaleReasonLibraryNotLoaded      = "library-not-loaded"
	TimescaleReasonAPIMissing            = "by-range-api-missing"
	TimescaleReasonNotExecutable         = "function-not-executable"
	TimescaleReasonSchemaCreateDenied    = "schema-create-denied"
	TimescaleReasonSchemaUsageDenied     = "schema-usage-denied"
	TimescaleReasonExtensionSchemaDenied = "extension-schema-denied"
)

// TimescaleCapability is what a read-only probe found.
type TimescaleCapability struct {
	Available bool
	// Version is informational. No capability is inferred from it.
	Version string
	// ExtensionSchema is where the installed extension's functions live.
	ExtensionSchema string
	Reasons         []string
}

// ProbeTimescale reads, and only reads, whether the selected destination can
// prepare a TimescaleDB hypertable for a table in targetSchema. It runs in a
// read-only transaction, so it cannot create a table, install an extension or
// call create_hypertable, and it decides from the catalog, not from a version:
// the extension must be installed, its own (extension-owned, not merely
// same-named) generalized by_range and create_hypertable functions must exist
// and be executable by the connected role, the role must be able to use the
// extension's schema and both use and create in the target schema (a role that
// can create there but not use it cannot resolve the table it creates), and the
// library must be loaded. Every missing fact becomes a stable reason; none is repaired here.
func ProbeTimescale(ctx context.Context, db *sql.DB, targetSchema string) (capability TimescaleCapability, err error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return TimescaleCapability{}, fmt.Errorf("begin read-only timescale probe: %w", err)
	}
	defer func() { err = errors.Join(err, rollbackReadOnly(tx)) }()

	err = tx.QueryRowContext(ctx, `
		SELECT e.extversion, n.nspname
		FROM pg_catalog.pg_extension e
		JOIN pg_catalog.pg_namespace n ON n.oid = e.extnamespace
		WHERE e.extname = 'timescaledb'`).Scan(&capability.Version, &capability.ExtensionSchema)
	if errors.Is(err, sql.ErrNoRows) {
		capability.Reasons = []string{TimescaleReasonNotInstalled}
		return capability, nil
	}
	if err != nil {
		return TimescaleCapability{}, fmt.Errorf("read timescaledb extension: %w", err)
	}

	functions, err := extensionFunctions(ctx, tx)
	if err != nil {
		return TimescaleCapability{}, err
	}
	reasons := map[string]struct{}{}
	var createSignatures, rangeSignatures []string
	for _, function := range functions {
		if function.name == "create_hypertable" {
			createSignatures = append(createSignatures, function.arguments)
		} else {
			rangeSignatures = append(rangeSignatures, function.arguments)
		}
	}
	if !hasGeneralizedAPI(createSignatures, rangeSignatures) {
		reasons[TimescaleReasonAPIMissing] = struct{}{}
	} else {
		for _, function := range functions {
			if generalizedFunction(function) && !function.executable {
				reasons[TimescaleReasonNotExecutable] = struct{}{}
			}
		}
	}

	var usable, creatable, targetUsable sql.NullBool
	if err := tx.QueryRowContext(ctx, `
		SELECT has_schema_privilege(current_user, $1, 'USAGE'),
		       CASE WHEN to_regnamespace($2) IS NULL THEN false
		            ELSE has_schema_privilege(current_user, $2, 'CREATE') END,
		       CASE WHEN to_regnamespace($2) IS NULL THEN false
		            ELSE has_schema_privilege(current_user, $2, 'USAGE') END`,
		capability.ExtensionSchema, targetSchema).Scan(&usable, &creatable, &targetUsable); err != nil {
		return TimescaleCapability{}, fmt.Errorf("read timescale schema privileges: %w", err)
	}
	if !usable.Valid || !usable.Bool {
		reasons[TimescaleReasonExtensionSchemaDenied] = struct{}{}
	}
	if !creatable.Valid || !creatable.Bool {
		reasons[TimescaleReasonSchemaCreateDenied] = struct{}{}
	}
	if !targetUsable.Valid || !targetUsable.Bool {
		reasons[TimescaleReasonSchemaUsageDenied] = struct{}{}
	}

	if len(reasons) == 0 {
		// Every catalog fact holds; the last read-only fact is that the library is
		// loaded in this session, which the information views need.
		if _, err := tx.ExecContext(ctx, `SELECT 1 FROM timescaledb_information.hypertables LIMIT 1`); err != nil {
			if isPermissionDenied(err) {
				reasons[TimescaleReasonExtensionSchemaDenied] = struct{}{}
			} else {
				reasons[TimescaleReasonLibraryNotLoaded] = struct{}{}
			}
		}
	}
	capability.Reasons = slices.Sorted(maps.Keys(reasons))
	capability.Available = len(capability.Reasons) == 0
	if capability.Available {
		capability.Reasons = nil
	}
	return capability, nil
}

type extensionFunction struct {
	name, arguments string
	executable      bool
}

// extensionFunctions lists the extension-owned create_hypertable and by_range
// overloads in the extension's own schema. Ownership comes from pg_depend, so a
// function that only shares a name is never listed.
func extensionFunctions(ctx context.Context, tx *sql.Tx) ([]extensionFunction, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT p.proname, pg_catalog.pg_get_function_identity_arguments(p.oid),
		       has_function_privilege(current_user, p.oid, 'EXECUTE')
		FROM pg_catalog.pg_proc p
		JOIN pg_catalog.pg_depend d ON d.classid = 'pg_catalog.pg_proc'::regclass AND d.objid = p.oid AND d.deptype = 'e'
		JOIN pg_catalog.pg_extension e ON d.refclassid = 'pg_catalog.pg_extension'::regclass AND e.oid = d.refobjid
		WHERE e.extname = 'timescaledb' AND p.pronamespace = e.extnamespace
		  AND p.proname IN ('create_hypertable', 'by_range')
		ORDER BY p.proname, 2`)
	if err != nil {
		return nil, fmt.Errorf("read timescaledb functions: %w", err)
	}
	defer rows.Close()
	var functions []extensionFunction
	for rows.Next() {
		var function extensionFunction
		if err := rows.Scan(&function.name, &function.arguments, &function.executable); err != nil {
			return nil, fmt.Errorf("scan timescaledb function: %w", err)
		}
		functions = append(functions, function)
	}
	return functions, rows.Err()
}

const (
	generalizedCreateArguments  = "relation regclass, dimension dimension_info, create_default_indexes boolean, if_not_exists boolean, migrate_data boolean"
	generalizedByRangeArguments = "column_name name, partition_interval anyelement, partition_func regproc"
)

var dimensionInfoQualifier = regexp.MustCompile(`[A-Za-z0-9_"]+\.dimension_info`)

// normalizeSignature drops the schema qualifier of the extension's internal
// dimension type so the comparison does not depend on where it is printed.
func normalizeSignature(arguments string) string {
	return dimensionInfoQualifier.ReplaceAllString(strings.Join(strings.Fields(arguments), " "), "dimension_info")
}

// hasGeneralizedAPI reports whether the extension offers the generalized
// create_hypertable (a dimension_info argument, created by by_range) and
// by_range. The legacy time_column_name signature alone does not count, and
// neither does a function that only shares the argument names.
func hasGeneralizedAPI(createSignatures, byRangeSignatures []string) bool {
	return slices.ContainsFunc(createSignatures, func(s string) bool { return normalizeSignature(s) == generalizedCreateArguments }) &&
		slices.ContainsFunc(byRangeSignatures, func(s string) bool { return normalizeSignature(s) == generalizedByRangeArguments })
}

func generalizedFunction(function extensionFunction) bool {
	switch function.name {
	case "create_hypertable":
		return normalizeSignature(function.arguments) == generalizedCreateArguments
	case "by_range":
		return normalizeSignature(function.arguments) == generalizedByRangeArguments
	}
	return false
}

func isPermissionDenied(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "42501"
}

// HypertableInfo is the catalog evidence of one hypertable.
type HypertableInfo struct {
	TimeColumn string
	// ChunkInterval is the catalog's chunk interval at the microsecond PostgreSQL
	// stores it. It is zero, with ChunkIntervalUnknown set, whenever the catalog
	// value is not a non-negative whole number of microseconds that a Duration
	// can hold; such a value never equals a confirmed interval.
	ChunkInterval        time.Duration
	ChunkIntervalUnknown bool
	Dimensions           int
	Chunks               int
	CompressionEnabled   bool
	// Jobs counts background jobs bound to this hypertable.
	Jobs int
}

// InspectHypertable reads the catalog evidence of one table. It returns nil
// when the table is not a hypertable, including when the extension is not
// installed (nothing can be a hypertable then). It runs in a read-only
// transaction and changes nothing.
func InspectHypertable(ctx context.Context, db *sql.DB, schemaName, tableName string) (_ *HypertableInfo, err error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, fmt.Errorf("begin read-only hypertable inspection: %w", err)
	}
	defer func() { err = errors.Join(err, rollbackReadOnly(tx)) }()

	var views bool
	if err := tx.QueryRowContext(ctx, `SELECT to_regclass('timescaledb_information.hypertables') IS NOT NULL`).Scan(&views); err != nil {
		return nil, fmt.Errorf("look up timescaledb views: %w", err)
	}
	if !views {
		return nil, nil
	}
	info := &HypertableInfo{}
	err = tx.QueryRowContext(ctx, `
		SELECT num_dimensions, num_chunks, compression_enabled FROM timescaledb_information.hypertables
		WHERE hypertable_schema = $1 AND hypertable_name = $2`, schemaName, tableName).
		Scan(&info.Dimensions, &info.Chunks, &info.CompressionEnabled)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read hypertable: %w", err)
	}
	var dimensionType string
	var micros sql.NullString
	// EXTRACT(EPOCH) is an exact numeric from PostgreSQL 14 on; an older server
	// returns a float, which cannot prove a microsecond, so it reads as unknown.
	err = tx.QueryRowContext(ctx, `
		SELECT column_name, dimension_type,
		       CASE WHEN current_setting('server_version_num')::int >= 140000
		            THEN (EXTRACT(EPOCH FROM time_interval) * 1000000)::numeric::text END
		FROM timescaledb_information.dimensions
		WHERE hypertable_schema = $1 AND hypertable_name = $2 ORDER BY dimension_number LIMIT 1`, schemaName, tableName).
		Scan(&info.TimeColumn, &dimensionType, &micros)
	if err != nil {
		return nil, fmt.Errorf("read hypertable dimension: %w", err)
	}
	var known bool
	if dimensionType == "Time" && micros.Valid {
		info.ChunkInterval, known = durationFromMicrosText(micros.String)
	}
	info.ChunkIntervalUnknown = !known
	if err := tx.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM timescaledb_information.jobs WHERE hypertable_schema = $1 AND hypertable_name = $2`,
		schemaName, tableName).Scan(&info.Jobs); err != nil {
		return nil, fmt.Errorf("read hypertable jobs: %w", err)
	}
	return info, nil
}

// rollbackReadOnly finishes a read-only transaction. Nothing was written, so
// rolling back loses nothing; a transaction that is already finished is not an
// error.
func rollbackReadOnly(tx *sql.Tx) error {
	if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		return fmt.Errorf("end read-only transaction: %w", err)
	}
	return nil
}

// durationFromMicrosText converts the catalog's exact microsecond count, a
// numeric rendered as text such as "604800000001.000000", into a duration. It
// does not go through a float or round: anything that is not a non-negative
// whole number of microseconds, or that does not fit a Duration (whose
// overflow would wrap around to an unrelated value), is not representable.
func durationFromMicrosText(text string) (time.Duration, bool) {
	whole, fraction, _ := strings.Cut(text, ".")
	if whole == "" || strings.Trim(fraction, "0") != "" || strings.IndexFunc(whole, func(r rune) bool { return r < '0' || r > '9' }) >= 0 {
		return 0, false
	}
	micros, err := strconv.ParseInt(whole, 10, 64)
	if err != nil || micros > math.MaxInt64/int64(time.Microsecond) {
		return 0, false
	}
	return time.Duration(micros) * time.Microsecond, true
}
