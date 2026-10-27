package dbtarget

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
)

const (
	mysqlFuturePartitionName = "p_future"
	mysqlFuturePartitionLock = 2
	mysqlFuturePartitionTTL  = 5 * time.Second
)

var (
	mysqlMonthlyPartitionPattern = regexp.MustCompile(`^p_\d{6}$`)
	mysqlGapPartitionPattern     = regexp.MustCompile(`^p_gap_\d{6}$`)
)

// ReorganizeOpenMySQLFuturePartition extends one owned native-monthly managed
// table's bounded future horizon. It only reorganizes p_future; existing
// partitions and the receipt table are never touched.
func ReorganizeOpenMySQLFuturePartition(
	ctx context.Context,
	db *sql.DB,
	expected *MySQLRecordingCapability,
	database string,
	tableProof *MySQLManagedInspection,
	after *MySQLPartitionEvidence,
) (err error) {
	if db == nil || expected == nil || tableProof == nil || after == nil {
		return fmt.Errorf("%w: mysql partition reorganize input is incomplete", ErrMySQLManagedTableUnconfirmed)
	}
	additions, err := validateMySQLFuturePartitionPlan(tableProof, after, database)
	if err != nil {
		return err
	}
	runCtx, cancel := context.WithTimeout(ctx, mysqlFuturePartitionTTL)
	defer cancel()
	actual, err := reprojectOpenMySQLCapability(runCtx, db, expected, database)
	if err != nil {
		return err
	}
	proof := tableProof.Evidence
	if proof == nil || proof.ServerFamily != actual.ServerFamily || proof.ServerVersion != actual.ServerVersion ||
		proof.CapabilityRevision != actual.CapabilityRevision || proof.Database != database {
		return fmt.Errorf("%w: mysql partition proof binding changed", ErrMySQLManagedTableUnconfirmed)
	}

	conn, err := db.Conn(runCtx)
	if err != nil {
		return fmt.Errorf("open dedicated mysql partition connection: %w", err)
	}
	originalLockWait, err := mysqlFuturePartitionLockSetting(runCtx, conn)
	if err != nil {
		return errors.Join(err, mysqlFuturePartitionDiscard(conn))
	}
	defer func() { err = errors.Join(err, mysqlFuturePartitionRelease(ctx, conn, originalLockWait)) }()

	before, err := readMySQLManagedTableFacts(runCtx, conn, actual.ServerFamily, database, tableProof.Inspection.Table)
	if err != nil {
		return fmt.Errorf("read mysql partition table before reorganize: %w", err)
	}
	if err := matchMySQLManagedInspection(tableProof, before); err != nil {
		return err
	}
	statement, err := mysqlFuturePartitionReorganizeStatement(database, tableProof.Inspection.Table, additions)
	if err != nil {
		return err
	}
	if _, err := conn.ExecContext(runCtx, statement); err != nil {
		return fmt.Errorf("execute mysql partition reorganize: %w", err)
	}
	expectedAfter := *tableProof
	expectedEvidence := *tableProof.Evidence
	expectedEvidence.Partitions = cloneMySQLPartitionEvidence(after)
	expectedAfter.Evidence = &expectedEvidence
	afterFacts, err := readMySQLManagedTableFacts(runCtx, conn, actual.ServerFamily, database, tableProof.Inspection.Table)
	if err != nil {
		return fmt.Errorf("read mysql partition table after reorganize: %w", err)
	}
	return matchMySQLManagedInspection(&expectedAfter, afterFacts)
}

func validateMySQLFuturePartitionPlan(proof *MySQLManagedInspection, after *MySQLPartitionEvidence, database string) ([]MySQLPartitionBoundary, error) {
	if proof == nil || proof.Inspection == nil || proof.Evidence == nil || proof.Evidence.Partitions == nil {
		return nil, fmt.Errorf("%w: mysql partition physical proof is incomplete", ErrMySQLManagedTableUnconfirmed)
	}
	if proof.Evidence.Partitions.Method != mysqlPartitionMethodRangeColumns ||
		proof.Evidence.Partitions.Expression != mysqlPartitionExpression {
		return nil, fmt.Errorf("%w: mysql partition proof is not RANGE COLUMNS(bucket_start)", ErrMySQLManagedTableUnconfirmed)
	}
	if !strings.HasPrefix(proof.Inspection.Table, "gw_long_") || proof.Evidence.Database != database {
		return nil, fmt.Errorf("%w: mysql partition table is outside the managed scope", ErrMySQLManagedTableUnconfirmed)
	}
	if after == nil || after.Method != mysqlPartitionMethodRangeColumns || after.Expression != mysqlPartitionExpression {
		return nil, fmt.Errorf("%w: mysql future partition evidence is incomplete", ErrMySQLManagedTableUnconfirmed)
	}
	beforeParts := proof.Evidence.Partitions.Partitions
	if len(beforeParts) < 2 || beforeParts[len(beforeParts)-1].Name != mysqlFuturePartitionName ||
		beforeParts[len(beforeParts)-1].LessThan != mysqlPartitionMaxValue {
		return nil, fmt.Errorf("%w: mysql partition proof has no final p_future MAXVALUE", ErrMySQLManagedTableUnconfirmed)
	}
	for index, partition := range beforeParts {
		if partition.Ordinal != index+1 {
			return nil, fmt.Errorf("%w: mysql partition proof ordinals are incomplete", ErrMySQLManagedTableUnconfirmed)
		}
		if index < len(beforeParts)-1 && partition.LessThan == mysqlPartitionMaxValue {
			return nil, fmt.Errorf("%w: mysql partition proof has a non-final MAXVALUE", ErrMySQLManagedTableUnconfirmed)
		}
	}
	prefix := beforeParts[:len(beforeParts)-1]
	seen := make(map[string]struct{}, len(prefix)+len(after.Partitions))
	for index, partition := range prefix {
		if partition.Ordinal != index+1 || partition.LessThan == mysqlPartitionMaxValue {
			return nil, fmt.Errorf("%w: mysql finite prefix is incomplete", ErrMySQLManagedTableUnconfirmed)
		}
		if _, err := parseMySQLPartitionTime(partition.LessThan); err != nil {
			return nil, fmt.Errorf("%w: mysql finite prefix boundary is invalid", ErrMySQLManagedTableUnconfirmed)
		}
		seen[strings.ToLower(partition.Name)] = struct{}{}
	}
	if len(after.Partitions) < len(prefix)+2 ||
		!slices.Equal(prefix, after.Partitions[:len(prefix)]) {
		return nil, fmt.Errorf("%w: mysql future partition plan changed the finite prefix", ErrMySQLManagedTableUnconfirmed)
	}
	final := after.Partitions[len(after.Partitions)-1]
	if final.Name != mysqlFuturePartitionName || final.LessThan != mysqlPartitionMaxValue || final.Ordinal != len(after.Partitions) {
		return nil, fmt.Errorf("%w: mysql future partition plan has no final p_future MAXVALUE", ErrMySQLManagedTableUnconfirmed)
	}
	additions := after.Partitions[len(prefix) : len(after.Partitions)-1]
	if len(additions) == 0 || len(additions) > 5 {
		return nil, fmt.Errorf("%w: mysql future partition additions are out of bounds", ErrMySQLManagedTableUnconfirmed)
	}
	monthlyCount, gapCount := 0, 0
	lastBoundary := time.Time{}
	if len(prefix) > 0 {
		var err error
		lastBoundary, err = parseMySQLPartitionTime(prefix[len(prefix)-1].LessThan)
		if err != nil {
			return nil, fmt.Errorf("%w: mysql finite prefix boundary is invalid", ErrMySQLManagedTableUnconfirmed)
		}
	}
	for index, addition := range additions {
		if addition.Ordinal != len(prefix)+index+1 {
			return nil, fmt.Errorf("%w: mysql future partition addition ordinal is invalid", ErrMySQLManagedTableUnconfirmed)
		}
		if _, duplicate := seen[strings.ToLower(addition.Name)]; duplicate {
			return nil, fmt.Errorf("%w: mysql future partition name is duplicated", ErrMySQLManagedTableUnconfirmed)
		}
		seen[strings.ToLower(addition.Name)] = struct{}{}
		boundary, err := parseMySQLFirstOfMonth(addition.LessThan)
		if err != nil {
			return nil, err
		}
		if !lastBoundary.IsZero() && !lastBoundary.Before(boundary) {
			return nil, fmt.Errorf("%w: mysql future partitions are not strictly increasing", ErrMySQLManagedTableUnconfirmed)
		}
		switch {
		case mysqlGapPartitionPattern.MatchString(addition.Name):
			if index != 0 || gapCount > 0 {
				return nil, fmt.Errorf("%w: mysql gap partition must be the single first addition", ErrMySQLManagedTableUnconfirmed)
			}
			if addition.Name != "p_gap_"+boundary.Format("200601") {
				return nil, fmt.Errorf("%w: mysql gap partition name does not match its upper boundary", ErrMySQLManagedTableUnconfirmed)
			}
			gapCount++
		case mysqlMonthlyPartitionPattern.MatchString(addition.Name):
			if lastBoundary.IsZero() {
				return nil, fmt.Errorf("%w: mysql monthly addition has no finite previous boundary", ErrMySQLManagedTableUnconfirmed)
			}
			if addition.Name != "p_"+lastBoundary.Format("200601") {
				return nil, fmt.Errorf("%w: mysql monthly partition name does not match the previous boundary", ErrMySQLManagedTableUnconfirmed)
			}
			if !boundary.Equal(lastBoundary.AddDate(0, 1, 0)) {
				return nil, fmt.Errorf("%w: mysql monthly partition upper boundary is not the next month", ErrMySQLManagedTableUnconfirmed)
			}
			monthlyCount++
		default:
			return nil, fmt.Errorf("%w: mysql future partition name is invalid", ErrMySQLManagedTableUnconfirmed)
		}
		lastBoundary = boundary
	}
	if monthlyCount == 0 || monthlyCount > 4 || gapCount > 1 {
		return nil, fmt.Errorf("%w: mysql future partition monthly additions are out of bounds", ErrMySQLManagedTableUnconfirmed)
	}
	return additions, nil
}

func parseMySQLFirstOfMonth(raw string) (time.Time, error) {
	parsed, err := parseMySQLPartitionTime(raw)
	if err != nil || parsed.Day() != 1 || parsed.Hour() != 0 || parsed.Minute() != 0 ||
		parsed.Second() != 0 || parsed.Nanosecond() != 0 || parsed.Format(time.RFC3339Nano) != raw {
		return time.Time{}, fmt.Errorf("%w: mysql future boundary is not a canonical UTC first-of-month", ErrMySQLManagedTableUnconfirmed)
	}
	return parsed, nil
}

func parseMySQLPartitionTime(raw string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339Nano, raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: mysql partition boundary is invalid", ErrMySQLManagedTableUnconfirmed)
	}
	parsed = parsed.UTC()
	if parsed.Year() < 1000 || parsed.Year() > 9999 {
		return time.Time{}, fmt.Errorf("%w: mysql partition boundary is outside DATETIME range", ErrMySQLManagedTableUnconfirmed)
	}
	return parsed, nil
}

func mysqlFuturePartitionReorganizeStatement(database, table string, additions []MySQLPartitionBoundary) (string, error) {
	definitions := make([]string, 0, len(additions)+1)
	for _, addition := range additions {
		boundary, err := parseMySQLFirstOfMonth(addition.LessThan)
		if err != nil {
			return "", err
		}
		definitions = append(definitions, "PARTITION "+quoteIdentifier(schema.DatabaseConnectorKindMySQL, addition.Name)+
			" VALUES LESS THAN ('"+boundary.Format("2006-01-02 15:04:05")+"')")
	}
	definitions = append(definitions, "PARTITION "+quoteIdentifier(schema.DatabaseConnectorKindMySQL, mysqlFuturePartitionName)+" VALUES LESS THAN (MAXVALUE)")
	return "ALTER TABLE " + qualifiedTableName(schema.DatabaseConnectorKindMySQL, database, table) +
		" REORGANIZE PARTITION " + quoteIdentifier(schema.DatabaseConnectorKindMySQL, mysqlFuturePartitionName) +
		" INTO (" + strings.Join(definitions, ", ") + ");", nil
}

func mysqlFuturePartitionLockSetting(ctx context.Context, conn *sql.Conn) (int, error) {
	var original int
	if err := conn.QueryRowContext(ctx, "SELECT @@SESSION.lock_wait_timeout").Scan(&original); err != nil {
		return 0, fmt.Errorf("read mysql partition lock_wait_timeout: %w", err)
	}
	if _, err := conn.ExecContext(ctx, fmt.Sprintf("SET SESSION lock_wait_timeout = %d", mysqlFuturePartitionLock)); err != nil {
		return 0, fmt.Errorf("set mysql partition lock_wait_timeout: %w", err)
	}
	return original, nil
}

func mysqlFuturePartitionRelease(ctx context.Context, conn *sql.Conn, original int) error {
	restoreCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Second)
	defer cancel()
	if _, err := conn.ExecContext(restoreCtx, fmt.Sprintf("SET SESSION lock_wait_timeout = %d", original)); err != nil {
		return errors.Join(fmt.Errorf("restore mysql partition lock timeout: %w", err), mysqlFuturePartitionDiscard(conn))
	}
	return conn.Close()
}

func mysqlFuturePartitionDiscard(conn *sql.Conn) error {
	err := conn.Raw(func(any) error { return driver.ErrBadConn })
	if errors.Is(err, driver.ErrBadConn) {
		err = nil
	}
	return errors.Join(err, conn.Close())
}
