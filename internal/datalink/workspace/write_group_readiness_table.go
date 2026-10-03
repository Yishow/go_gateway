package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"os"
	"sort"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

const (
	writeGroupColumnTypeBool    = "bool"
	writeGroupColumnTypeBit     = "bit"
	writeGroupColumnTypeInteger = "integer"
	writeGroupColumnTypeFloat   = "float"
	writeGroupColumnTypeNumeric = "numeric"
	writeGroupColumnTypeString  = "string"
	writeGroupColumnTypeJSON    = "json"
	writeGroupColumnTypeUUID    = "uuid"
)

func allowWriteGroupTableInspection(snapshot *writeGroupReadinessSnapshot) (bool, ReadinessIssue) {
	if snapshot == nil || snapshot.validated == nil {
		return false, writeGroupReadinessIssue(
			"schema-scope-unverified", "destination schema scope could not be verified", groupIDOrEmpty(snapshot),
		)
	}
	if snapshot.connectorKind != writeGroupConnectorKindSQLite {
		return true, ReadinessIssue{}
	}
	targetPath, ok := sqliteWriteGroupTargetPath(snapshot.connectorConfig)
	if !ok {
		return false, writeGroupReadinessIssue(
			"schema-scope-unverified", "destination database scope could not be verified", snapshot.validated.ID,
		)
	}
	if savedDatabase := strings.TrimSpace(snapshot.validated.Destination.Database); savedDatabase != "" {
		savedPath, savedOK := normalizeSQLiteWriteGroupPath(savedDatabase)
		if !savedOK || savedPath != targetPath {
			return false, writeGroupReadinessIssue(
				"schema-scope-mismatch", "saved database scope does not match the connector file", snapshot.validated.ID,
			)
		}
	}
	info, err := os.Stat(targetPath)
	if os.IsNotExist(err) {
		return false, writeGroupReadinessIssue(
			"schema-target-missing", "destination database is not present", snapshot.validated.ID,
		)
	}
	if err != nil || !info.Mode().IsRegular() {
		return false, writeGroupReadinessIssue(
			"schema-scope-unverified", "destination database scope could not be verified", snapshot.validated.ID,
		)
	}
	return true, ReadinessIssue{}
}

func sqliteWriteGroupTargetPath(rawConfig string) (string, bool) {
	var config map[string]json.RawMessage
	if err := json.Unmarshal([]byte(rawConfig), &config); err != nil {
		return "", false
	}
	for _, key := range []string{"dsn", "path"} {
		raw, ok := config[key]
		if !ok {
			continue
		}
		var value string
		if err := json.Unmarshal(raw, &value); err != nil {
			continue
		}
		if path, ok := normalizeSQLiteWriteGroupPath(value); ok {
			return path, true
		}
	}
	return "", false
}

func normalizeSQLiteWriteGroupPath(value string) (string, bool) {
	value = strings.TrimSpace(value)
	if value == "" || strings.EqualFold(value, ":memory:") || strings.EqualFold(value, "memory") {
		return "", false
	}
	if strings.HasPrefix(strings.ToLower(value), "file:") {
		parsed, err := url.Parse(value)
		if err != nil || (parsed.Host != "" && !strings.EqualFold(parsed.Host, "localhost")) {
			return "", false
		}
		if strings.EqualFold(parsed.Query().Get("mode"), "memory") {
			return "", false
		}
		value = parsed.Path
		if value == "" {
			value = strings.TrimPrefix(value, "file:")
		}
		if decoded, err := url.PathUnescape(value); err == nil {
			value = decoded
		}
	} else if queryIndex := strings.IndexByte(value, '?'); queryIndex >= 0 {
		query := value[queryIndex+1:]
		if parsed, err := url.ParseQuery(query); err == nil &&
			strings.EqualFold(parsed.Get("mode"), "memory") {
			return "", false
		}
		value = value[:queryIndex]
	}
	value = strings.TrimSpace(value)
	if value == "" || strings.EqualFold(value, ":memory:") {
		return "", false
	}
	return value, true
}

func evaluateWriteGroupInspection(
	group *WriteGroup,
	tagTypes map[string]schema.DataType,
	inspection *dbtarget.TableInspection,
	connectorKinds ...string,
) (bool, []ReadinessIssue) {
	if inspection == nil {
		return false, []ReadinessIssue{writeGroupReadinessIssue(
			"schema-inspection-failed", "destination schema could not be verified", groupIDOrEmpty(group),
		)}
	}
	if inspection.Status != dbtarget.TableInspectionExists {
		code := "schema-table-unavailable"
		if inspection.Status == dbtarget.TableInspectionForbidden {
			code = "schema-table-forbidden"
		}
		return false, []ReadinessIssue{writeGroupReadinessIssue(
			code, "destination table could not be verified", groupIDOrEmpty(group),
		)}
	}
	issues := make([]ReadinessIssue, 0)
	dialect, dialectKnown := writeGroupReadinessDialect(connectorKinds...)
	if !dialectKnown {
		return false, []ReadinessIssue{writeGroupReadinessIssue(
			"destination-column-type-unverified", "destination SQL dialect could not be verified", group.ID,
		)}
	}
	layoutMembers := make([]dbtarget.GroupRowMember, 0, len(group.Members))
	entityKeyed := false
	if strings.TrimSpace(inspection.Schema) != strings.TrimSpace(group.Destination.TableSchema) {
		issues = append(issues, writeGroupReadinessIssue(
			"schema-scope-mismatch", "destination schema does not match the saved binding", group.ID,
		))
	}
	if strings.TrimSpace(inspection.Table) != strings.TrimSpace(group.Destination.TableName) {
		issues = append(issues, writeGroupReadinessIssue(
			"schema-scope-mismatch", "destination table does not match the saved binding", group.ID,
		))
	}
	for _, member := range group.Members {
		column, ok := findWriteGroupColumn(inspection.Columns, member.TargetColumn)
		if !ok {
			issues = append(issues, writeGroupReadinessIssue(
				"destination-column-missing", "a saved destination column is missing", group.ID,
			))
			continue
		}
		dataType := tagTypes[member.TagID]
		if !dataType.IsValid() {
			issues = append(issues, writeGroupReadinessIssue(
				"tag-type-unverified", "a source tag type could not be verified", group.ID,
			))
			continue
		}
		if _, known := normalizedWriteGroupColumnType(column.DataType); !known {
			issues = append(issues, writeGroupReadinessIssue(
				"destination-column-type-unverified", "a destination column type could not be verified", group.ID,
			))
			continue
		}
		kind, ok := measurement.ExactTypeForTag(dataType)
		if !ok {
			issues = append(issues, writeGroupReadinessIssue(
				"tag-type-unverified", "a source tag type could not be verified", group.ID,
			))
			continue
		}
		layoutMembers = append(layoutMembers, dbtarget.GroupRowMember{
			MemberKey: member.TagID, EntityKey: member.EntityKey, Column: member.TargetColumn,
			Type: kind, Required: member.Required,
		})
		entityKeyed = entityKeyed || strings.TrimSpace(member.EntityKey) != ""
	}
	_, layoutIssues := dbtarget.NewGroupRowLayout(dbtarget.GroupRowSpec{
		Dialect: dialect, Columns: inspection.Columns, Members: layoutMembers,
		EntityKeyed: entityKeyed, EntityKeyColumn: group.RowPolicy.EntityKeyColumn,
		ProvenanceColumn: group.RowPolicy.ProvenanceColumn,
	})
	for _, issue := range layoutIssues {
		code := "destination-column-type-mismatch"
		message := "a destination column is incompatible with the production writer"
		switch issue.Code {
		case "column-missing", "identity-column-missing":
			code = "destination-column-missing"
			message = "a destination column required by the production writer is missing"
		case "member-incomplete":
			code = "destination-column-type-unverified"
			message = "a production writer member type could not be verified"
		}
		issues = append(issues, writeGroupReadinessIssue(code, message, group.ID))
	}
	return len(issues) == 0, issues
}

func writeGroupReadinessDialect(connectorKinds ...string) (dbtarget.SQLDialect, bool) {
	if len(connectorKinds) == 0 || strings.EqualFold(strings.TrimSpace(connectorKinds[0]), writeGroupConnectorKindSQLite) {
		return dbtarget.SQLDialectSQLite, true
	}
	if strings.EqualFold(strings.TrimSpace(connectorKinds[0]), writeGroupConnectorKindPostgres) {
		return dbtarget.SQLDialectPostgres, true
	}
	return "", false
}

func findWriteGroupColumn(columns []dbtarget.ColumnInfo, name string) (dbtarget.ColumnInfo, bool) {
	name = strings.TrimSpace(name)
	for _, column := range columns {
		if strings.EqualFold(strings.TrimSpace(column.Name), name) {
			return column, true
		}
	}
	return dbtarget.ColumnInfo{}, false
}

func writeGroupColumnTypeCompatible(dataType schema.DataType, column dbtarget.ColumnInfo) (compatible, known bool) {
	columnType, known := normalizedWriteGroupColumnType(column.DataType)
	if !known {
		return false, false
	}
	switch dataType {
	case schema.DataTypeBool:
		return columnType == writeGroupColumnTypeBool || columnType == writeGroupColumnTypeBit || columnType == writeGroupColumnTypeInteger || columnType == writeGroupColumnTypeNumeric, true
	case schema.DataTypeInt16, schema.DataTypeUint16, schema.DataTypeInt32, schema.DataTypeUint32, schema.DataTypeInt64, schema.DataTypeUint64:
		return columnType == writeGroupColumnTypeInteger || columnType == writeGroupColumnTypeNumeric, true
	case schema.DataTypeFloat32, schema.DataTypeFloat64:
		return columnType == writeGroupColumnTypeFloat || columnType == writeGroupColumnTypeNumeric, true
	case schema.DataTypeString:
		return columnType == writeGroupColumnTypeString, true
	default:
		return false, false
	}
}

func normalizedWriteGroupColumnType(raw string) (string, bool) {
	base := strings.ToLower(strings.TrimSpace(raw))
	if base == "" || strings.Contains(base, "[]") {
		return "", false
	}
	if parenthesis := strings.IndexByte(base, '('); parenthesis >= 0 {
		base = strings.TrimSpace(base[:parenthesis])
	}
	switch strings.Join(strings.Fields(base), " ") {
	case "bool", "boolean":
		return writeGroupColumnTypeBool, true
	case "bit":
		return writeGroupColumnTypeBit, true
	case "int", "int2", "int4", "int8", "integer", "smallint", "bigint", "tinyint", "serial", "bigserial":
		return writeGroupColumnTypeInteger, true
	case "real", "float", "float4", "float8", "double", "double precision":
		return writeGroupColumnTypeFloat, true
	case "numeric", "decimal":
		return writeGroupColumnTypeNumeric, true
	case "char", "character", "varchar", "character varying", "nchar", "nvarchar", "text", "clob":
		return writeGroupColumnTypeString, true
	case "json", "jsonb":
		return writeGroupColumnTypeJSON, true
	case "uuid":
		return writeGroupColumnTypeUUID, true
	default:
		return "", false
	}
}

type writeGroupSchemaDigest struct {
	GroupID     string                    `json:"group_id"`
	ConnectorID string                    `json:"connector_id"`
	Schema      string                    `json:"schema"`
	Table       string                    `json:"table"`
	Columns     []dbtarget.ColumnInfo     `json:"columns"`
	Bindings    []writeGroupSchemaBinding `json:"bindings"`
}

type writeGroupSchemaBinding struct {
	DeviceID     string `json:"device_id"`
	PointID      string `json:"point_id"`
	TagID        string `json:"tag_id"`
	TargetColumn string `json:"target_column"`
}

func digestWriteGroupInspection(group *WriteGroup, inspection *dbtarget.TableInspection) (string, error) {
	digest := writeGroupSchemaDigest{
		GroupID: group.ID, ConnectorID: group.Destination.ConnectorID,
		Schema: group.Destination.TableSchema, Table: group.Destination.TableName,
		Columns:  append([]dbtarget.ColumnInfo(nil), inspection.Columns...),
		Bindings: make([]writeGroupSchemaBinding, 0, len(group.Members)),
	}
	sort.Slice(digest.Columns, func(i, j int) bool {
		left, right := strings.ToLower(digest.Columns[i].Name), strings.ToLower(digest.Columns[j].Name)
		if left == right {
			return digest.Columns[i].Name < digest.Columns[j].Name
		}
		return left < right
	})
	for _, member := range group.Members {
		digest.Bindings = append(digest.Bindings, writeGroupSchemaBinding{
			DeviceID: member.DeviceID, PointID: member.PointID,
			TagID: member.TagID, TargetColumn: member.TargetColumn,
		})
	}
	payload, err := json.Marshal(digest)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

func groupIDOrEmpty(value any) string {
	switch typed := value.(type) {
	case *writeGroupReadinessSnapshot:
		if typed != nil && typed.validated != nil {
			return typed.validated.ID
		}
	case *WriteGroup:
		if typed != nil {
			return typed.ID
		}
	}
	return ""
}
