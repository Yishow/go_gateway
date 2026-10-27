package recordingplan

import (
	"fmt"
	"strings"
	"time"
)

// PlanMySQLNativeMonthlyExtension retains the finite prefix and MAXVALUE,
// adding at most four future months and one gap for a long missed interval.
func PlanMySQLNativeMonthlyExtension(current *MySQLPartitionEvidence, at time.Time) (*MySQLPartitionEvidence, error) {
	if at.IsZero() || validateCurrentMySQLPartitionEvidence(current) != nil {
		return nil, fmt.Errorf("%w: invalid current evidence or zero planning time", ErrInvalidGroupSchemaLayout)
	}
	planned := cloneMySQLPartitionEvidence(current)
	finite := planned.Partitions[:len(planned.Partitions)-1]
	lastBoundary, err := parseMySQLMonthBoundary(finite[len(finite)-1].LessThan)
	if err != nil {
		return nil, fmt.Errorf("%w: current final finite boundary", ErrInvalidGroupSchemaLayout)
	}
	monthStart := time.Date(at.UTC().Year(), at.UTC().Month(), 1, 0, 0, 0, 0, time.UTC)
	targetEnd := monthStart.AddDate(0, mysqlNativeMonthlyPartitionCount, 0)
	if !mysqlDateTimeYearInRange(monthStart.Year()) || !mysqlDateTimeYearInRange(targetEnd.Year()) {
		return nil, mysqlPartitionDateTimeRangeError("maintenance", targetEnd)
	}
	appendBoundary := func(boundary MySQLPartitionBoundary) {
		boundary.Ordinal = len(finite) + 1
		finite = append(finite, boundary)
	}
	if lastBoundary.Before(monthStart) {
		appendBoundary(MySQLPartitionBoundary{Name: "p_gap_" + monthStart.Format("200601"), LessThan: monthStart.Format(time.RFC3339)})
		lastBoundary = monthStart
	}
	addedMonthly := 0
	for lastBoundary.Before(targetEnd) && addedMonthly < mysqlPartitionMaintenanceBatch {
		boundary := lastBoundary.AddDate(0, 1, 0)
		if !mysqlDateTimeYearInRange(boundary.Year()) {
			return nil, mysqlPartitionDateTimeRangeError("monthly", boundary)
		}
		appendBoundary(MySQLPartitionBoundary{
			Name:     fmt.Sprintf("p_%04d%02d", lastBoundary.Year(), lastBoundary.Month()),
			LessThan: boundary.Format(time.RFC3339),
		})
		lastBoundary = boundary
		addedMonthly++
	}
	if lastBoundary.Before(targetEnd) {
		return nil, fmt.Errorf("%w: extension exceeded bounded monthly batch", ErrInvalidGroupSchemaLayout)
	}
	planned.Partitions = append(finite, MySQLPartitionBoundary{Name: mysqlPartitionFutureName, LessThan: mysqlPartitionMaxValue})
	for i := range planned.Partitions {
		planned.Partitions[i].Ordinal = i + 1
	}
	if err := validateCurrentMySQLPartitionEvidence(planned); err != nil {
		return nil, err
	}
	return planned, nil
}

func validateCurrentMySQLPartitionEvidence(evidence *MySQLPartitionEvidence) error {
	invalid := func(reason string) error {
		return fmt.Errorf("%w: mysql native monthly partition %s", ErrInvalidGroupSchemaLayout, reason)
	}
	if evidence == nil || evidence.Method != mysqlPartitionMethodRangeColumns || evidence.Expression != mysqlPartitionExpression || len(evidence.Partitions) < mysqlPartitionBoundaryCount {
		return invalid("current partition state has an invalid method, expression, or count")
	}
	parts := evidence.Partitions
	past, err := parseMySQLMonthBoundary(parts[0].LessThan)
	if err != nil || past.Day() != 1 || parts[0] != (MySQLPartitionBoundary{Name: mysqlPartitionPastName, LessThan: past.Format(time.RFC3339), Ordinal: 1}) {
		return invalid("past partition is not the canonical UTC month boundary")
	}
	seen := map[string]struct{}{strings.ToLower(parts[0].Name): {}}
	previous := past
	monthlyCount := 0
	needMonthAfterGap := false
	for index, part := range parts[1 : len(parts)-1] {
		ordinal := index + 2
		if part.Ordinal != ordinal || part.Name == "" {
			return invalid("partition names or ordinals are invalid")
		}
		name := strings.ToLower(part.Name)
		if _, exists := seen[name]; exists {
			return invalid("partition names are duplicated")
		}
		seen[name] = struct{}{}
		boundary, parseErr := parseMySQLMonthBoundary(part.LessThan)
		if parseErr != nil || boundary.Day() != 1 || boundary.Hour() != 0 || boundary.Minute() != 0 || boundary.Second() != 0 || boundary.Nanosecond() != 0 || !previous.Before(boundary) {
			return invalid("finite boundaries are not increasing UTC month starts")
		}
		if strings.HasPrefix(part.Name, "p_gap_") {
			if needMonthAfterGap || part.Name != "p_gap_"+boundary.Format("200601") {
				return invalid("gap partition name or sequence is invalid")
			}
			needMonthAfterGap = true
		} else {
			month := previous
			if needMonthAfterGap && boundary.Equal(previous.AddDate(0, 1, 0)) {
				needMonthAfterGap = false
			}
			if part.Name != fmt.Sprintf("p_%04d%02d", month.Year(), month.Month()) || !boundary.Equal(month.AddDate(0, 1, 0)) {
				return invalid("monthly partitions are not contiguous and canonical")
			}
			monthlyCount++
		}
		previous = boundary
	}
	if needMonthAfterGap || monthlyCount < mysqlNativeMonthlyPartitionCount {
		return invalid("gap is not followed by a current month or the monthly prefix is incomplete")
	}
	future := parts[len(parts)-1]
	if future != (MySQLPartitionBoundary{Name: mysqlPartitionFutureName, LessThan: mysqlPartitionMaxValue, Ordinal: len(parts)}) {
		return invalid("p_future MAXVALUE partition must remain final")
	}
	return nil
}
