package dbtarget

import (
	"context"
	"database/sql"
	"strings"
)

const (
	mysqlRecordingCurrentUserQuery        = "SELECT CURRENT_USER()"
	mysqlRecordingCurrentRoleQuery        = "SELECT CURRENT_ROLE()"
	mysqlRecordingShowGrantsQuery         = "SHOW GRANTS"
	mysqlRecordingGlobalPrivilegesQuery   = "SELECT GRANTEE, PRIVILEGE_TYPE FROM information_schema.USER_PRIVILEGES"
	mysqlRecordingSchemaPrivilegesQuery   = "SELECT TABLE_SCHEMA, GRANTEE, PRIVILEGE_TYPE FROM information_schema.SCHEMA_PRIVILEGES"
	mysqlRecordingPermissionDeniedReason  = "operation_permissions_denied"
	mysqlRecordingPermissionUnknownReason = "operation_permissions_unconfirmed"
	mysqlRecordingPartialRevokeReason     = "operation_permissions_partial_revoke"
	mysqlRecordingRoleUnknownReason       = "operation_permissions_role_unverified"
)

// MySQLRecordingPermissions contains the read-only privilege evidence used by
// the MySQL/MariaDB operation gates. A false field is an observed missing
// privilege only when Verified is true; otherwise the probe could not prove
// the effective grants and every operation remains closed.
type MySQLRecordingPermissions struct {
	Verified     bool `json:"verified"`
	CreateTable  bool `json:"create_table"`
	AlterSchema  bool `json:"alter_schema"`
	IndexSchema  bool `json:"index_schema"`
	SelectSchema bool `json:"select_schema"`
	InsertSchema bool `json:"insert_schema"`
	DeleteSchema bool `json:"delete_schema"`
}

type mysqlRecordingPrivilege struct {
	schemaName    string
	grantee       string
	privilegeType string
}

func inspectMySQLRecordingPermissions(
	ctx context.Context,
	tx *sql.Tx,
	expectedDatabase, serverFamily string,
) (permissions MySQLRecordingPermissions, reason string) {
	if tx == nil || strings.TrimSpace(expectedDatabase) == "" {
		return permissions, mysqlRecordingPermissionUnknownReason
	}

	var currentUser, currentRole sql.NullString
	if err := tx.QueryRowContext(ctx, mysqlRecordingCurrentUserQuery).Scan(&currentUser); err != nil || !currentUser.Valid {
		return permissions, mysqlRecordingPermissionQueryReason(err)
	}
	if _, _, ok := mysqlRecordingParseAccount(currentUser.String); !ok {
		return permissions, mysqlRecordingPermissionUnknownReason
	}
	if err := tx.QueryRowContext(ctx, mysqlRecordingCurrentRoleQuery).Scan(&currentRole); err != nil {
		return permissions, mysqlRecordingPermissionQueryReason(err)
	}
	if !currentRole.Valid {
		// MariaDB 11.4 returns SQL NULL when no role is active. MySQL returns
		// the literal NONE; a NULL from another family is not evidence.
		if serverFamily != mysqlManagedFamilyMariaDB {
			return permissions, mysqlRecordingPermissionUnknownReason
		}
		currentRole.String = "NONE"
	}
	// Active roles can add or remove effective privileges. Do not infer their
	// contents from direct USER_PRIVILEGES/SCHEMA_PRIVILEGES rows; a role-aware
	// probe must establish the same effective scope before opening a gate.
	if role := strings.TrimSpace(currentRole.String); role != "" && !strings.EqualFold(role, "NONE") {
		return permissions, mysqlRecordingRoleUnknownReason
	}

	partialRevoke, err := readMySQLRecordingShowGrantsForPartialRevoke(ctx, tx)
	if err != nil {
		return permissions, mysqlRecordingPermissionQueryReason(err)
	}
	if partialRevoke {
		return permissions, mysqlRecordingPartialRevokeReason
	}

	global, err := readMySQLRecordingPrivileges(ctx, tx, mysqlRecordingGlobalPrivilegesQuery)
	if err != nil {
		return permissions, mysqlRecordingPermissionQueryReason(err)
	}
	schema, err := readMySQLRecordingSchemaPrivileges(ctx, tx)
	if err != nil {
		return permissions, mysqlRecordingPermissionQueryReason(err)
	}

	for _, grant := range global {
		if mysqlRecordingGranteeMatches(grant.grantee, currentUser.String) {
			mysqlRecordingApplyPrivilege(&permissions, grant.privilegeType)
		}
	}
	for _, grant := range schema {
		if !mysqlRecordingGranteeMatches(grant.grantee, currentUser.String) {
			continue
		}
		if !mysqlRecordingSchemaMatches(grant.schemaName, expectedDatabase) {
			if !mysqlRecordingSchemaScopeKnown(grant.schemaName) {
				return permissions, mysqlRecordingPermissionUnknownReason
			}
			continue
		}
		mysqlRecordingApplyPrivilege(&permissions, grant.privilegeType)
	}
	permissions.Verified = true
	return permissions, ""
}

func readMySQLRecordingPrivileges(
	ctx context.Context,
	tx *sql.Tx,
	query string,
	args ...any,
) ([]mysqlRecordingPrivilege, error) {
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	privileges := make([]mysqlRecordingPrivilege, 0)
	for rows.Next() {
		var grant mysqlRecordingPrivilege
		if err := rows.Scan(&grant.grantee, &grant.privilegeType); err != nil {
			return nil, err
		}
		privileges = append(privileges, grant)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return privileges, nil
}

func readMySQLRecordingSchemaPrivileges(ctx context.Context, tx *sql.Tx) ([]mysqlRecordingPrivilege, error) {
	rows, err := tx.QueryContext(ctx, mysqlRecordingSchemaPrivilegesQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	privileges := make([]mysqlRecordingPrivilege, 0)
	for rows.Next() {
		var grant mysqlRecordingPrivilege
		if err := rows.Scan(&grant.schemaName, &grant.grantee, &grant.privilegeType); err != nil {
			return nil, err
		}
		privileges = append(privileges, grant)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return privileges, nil
}

func readMySQLRecordingShowGrantsForPartialRevoke(ctx context.Context, tx *sql.Tx) (bool, error) {
	rows, err := tx.QueryContext(ctx, mysqlRecordingShowGrantsQuery)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var grant sql.NullString
		if err := rows.Scan(&grant); err != nil {
			return false, err
		}
		if !grant.Valid {
			return false, sql.ErrNoRows
		}
		grantText := strings.ToUpper(strings.TrimSpace(grant.String))
		if grantText == "REVOKE" || strings.HasPrefix(grantText, "REVOKE ") {
			return true, nil
		}
	}
	if err := rows.Err(); err != nil {
		return false, err
	}
	return false, nil
}

func mysqlRecordingPermissionQueryReason(err error) string {
	if err == nil {
		return mysqlRecordingPermissionUnknownReason
	}
	if isPermissionDeniedError(err) {
		return "operation_permissions_forbidden"
	}
	return mysqlRecordingPermissionUnknownReason
}

func mysqlRecordingGranteeMatches(grantee, currentUser string) bool {
	granteeUser, granteeHost, granteeOK := mysqlRecordingParseAccount(grantee)
	currentUserName, currentHost, currentOK := mysqlRecordingParseAccount(currentUser)
	return granteeOK && currentOK && granteeUser == currentUserName && strings.EqualFold(granteeHost, currentHost)
}

func mysqlRecordingParseAccount(value string) (user, host string, ok bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", "", false
	}
	if value[0] == '\'' || value[0] == '`' {
		var rest string
		user, rest, ok = mysqlRecordingParseQuotedAccountPart(value)
		if !ok {
			return "", "", false
		}
		rest = strings.TrimSpace(rest)
		if rest == "" || rest[0] != '@' {
			return "", "", false
		}
		rest = strings.TrimSpace(rest[1:])
		if rest == "" {
			return "", "", false
		}
		if rest[0] == '\'' || rest[0] == '`' {
			host, rest, ok = mysqlRecordingParseQuotedAccountPart(rest)
			return user, host, ok && strings.TrimSpace(rest) == ""
		}
		if strings.ContainsAny(rest, " \t\r\n") {
			return "", "", false
		}
		return user, rest, true
	}
	at := strings.LastIndexByte(value, '@')
	if at <= 0 || at == len(value)-1 || strings.TrimSpace(value[:at]) != value[:at] || strings.TrimSpace(value[at+1:]) != value[at+1:] {
		return "", "", false
	}
	return value[:at], value[at+1:], true
}

func mysqlRecordingParseQuotedAccountPart(value string) (part, rest string, ok bool) {
	quote := value[0]
	var builder strings.Builder
	for index := 1; index < len(value); index++ {
		character := value[index]
		switch character {
		case quote:
			if index+1 < len(value) && value[index+1] == quote {
				builder.WriteByte(quote)
				index++
				continue
			}
			return builder.String(), value[index+1:], true
		case '\\':
			if index+1 >= len(value) || (value[index+1] != quote && value[index+1] != '\\') {
				return "", "", false
			}
			builder.WriteByte(value[index+1])
			index++
		default:
			builder.WriteByte(character)
		}
	}
	return "", "", false
}

func mysqlRecordingSchemaMatches(observed, expected string) bool {
	observed = strings.TrimSpace(observed)
	expected = strings.TrimSpace(expected)
	if observed == expected {
		return true
	}
	decoded, ok := mysqlRecordingDecodeSchemaLiteral(observed)
	return ok && decoded == expected
}

func mysqlRecordingSchemaScopeKnown(observed string) bool {
	observed = strings.TrimSpace(observed)
	if observed == "" {
		return false
	}
	_, ok := mysqlRecordingDecodeSchemaLiteral(observed)
	return ok
}

func mysqlRecordingDecodeSchemaLiteral(value string) (string, bool) {
	var builder strings.Builder
	builder.Grow(len(value))
	for index := 0; index < len(value); index++ {
		character := value[index]
		switch character {
		case '\\':
			if index+1 >= len(value) {
				return "", false
			}
			next := value[index+1]
			if next != '\\' && next != '_' && next != '%' {
				return "", false
			}
			builder.WriteByte(next)
			index++
		case '_', '%':
			return "", false
		default:
			builder.WriteByte(character)
		}
	}
	return builder.String(), true
}

func mysqlRecordingApplyPrivilege(permissions *MySQLRecordingPermissions, privilegeType string) {
	if permissions == nil {
		return
	}
	if privilege := strings.ToUpper(strings.TrimSpace(privilegeType)); privilege == "ALL" || privilege == "ALL PRIVILEGES" {
		permissions.CreateTable = true
		permissions.AlterSchema = true
		permissions.IndexSchema = true
		permissions.SelectSchema = true
		permissions.InsertSchema = true
		permissions.DeleteSchema = true
		return
	}
	switch strings.ToUpper(strings.TrimSpace(privilegeType)) {
	case "CREATE":
		permissions.CreateTable = true
	case "ALTER":
		permissions.AlterSchema = true
	case "INDEX":
		permissions.IndexSchema = true
	case "SELECT":
		permissions.SelectSchema = true
	case "INSERT":
		permissions.InsertSchema = true
	case "DELETE":
		permissions.DeleteSchema = true
	}
}

func mysqlRecordingPermissionsComplete(permissions MySQLRecordingPermissions) bool {
	return permissions.Verified && permissions.CreateTable && permissions.AlterSchema && permissions.IndexSchema &&
		permissions.SelectSchema && permissions.InsertSchema && permissions.DeleteSchema
}

func mysqlRecordingPermissionsAllowManagedLong(permissions MySQLRecordingPermissions) bool {
	return permissions.Verified && permissions.CreateTable && permissions.IndexSchema &&
		permissions.SelectSchema && permissions.InsertSchema
}

func mysqlRecordingPermissionsAllowNativePartition(permissions MySQLRecordingPermissions) bool {
	return mysqlRecordingPermissionsAllowManagedLong(permissions) && permissions.AlterSchema
}

func mysqlRecordingPermissionsAllowGroupTestWrite(permissions MySQLRecordingPermissions) bool {
	return permissions.Verified && permissions.SelectSchema && permissions.InsertSchema && permissions.DeleteSchema
}

func mysqlRecordingPermissionsAllowMissingRepair(permissions MySQLRecordingPermissions) bool {
	return permissions.Verified && permissions.CreateTable && permissions.IndexSchema && permissions.SelectSchema
}
