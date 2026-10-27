package dbtarget

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

const (
	mysqlManagedFamilyMySQL     = "mysql"
	mysqlManagedFamilyMariaDB   = "mariadb"
	mysqlManagedIndexFunctional = "functional"
	mysqlManagedIndexIgnored    = "ignored"
	mysqlManagedIndexTypeBTREE  = "BTREE"
	mysqlManagedEngineInnoDB    = "InnoDB"
)

var mysqlManagedIdentifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type mysqlManagedVisibilityRow struct {
	name    string
	visible sql.NullString
	ignored sql.NullString
}

type mysqlManagedPrimaryKeyRow struct {
	column    sql.NullString
	subPart   sql.NullInt64
	indexType string
	nonUnique int
}

// mysqlManagedCheckIsJSONValid accepts exactly one json_valid(<column>)
// expression, optionally wrapped in outer parentheses/whitespace. Anything
// combined with OR/NOT, a string literal, another column or extra arguments is
// not JSON evidence for the managed column.
func mysqlManagedCheckIsJSONValid(clause, column string) bool {
	expr := stripMySQLCheckOuterParens(strings.TrimSpace(clause))
	const functionName = "json_valid"
	if len(expr) <= len(functionName) || !strings.EqualFold(expr[:len(functionName)], functionName) {
		return false
	}
	rest := strings.TrimSpace(expr[len(functionName):])
	if len(rest) < 2 || rest[0] != '(' || rest[len(rest)-1] != ')' {
		return false
	}
	arg := strings.TrimSpace(rest[1 : len(rest)-1])
	if arg == "" || strings.ContainsAny(arg, " \t\r\n()=,+-*/%<>&|!") {
		return false
	}
	if strings.HasPrefix(arg, "`") {
		if len(arg) < 2 || arg[len(arg)-1] != '`' || strings.Contains(arg[1:len(arg)-1], "`") {
			return false
		}
		return arg[1:len(arg)-1] == column
	}
	if !mysqlManagedIdentifierPattern.MatchString(arg) {
		return false
	}
	return arg == column
}

func stripMySQLCheckOuterParens(expr string) string {
	for len(expr) >= 2 && expr[0] == '(' && expr[len(expr)-1] == ')' {
		depth, enclosed := 0, true
		for i, r := range expr {
			switch r {
			case '(':
				depth++
			case ')':
				depth--
				if depth < 0 || (depth == 0 && i != len(expr)-1) {
					enclosed = false
				}
			}
			if !enclosed {
				break
			}
		}
		if !enclosed || depth != 0 {
			break
		}
		expr = strings.TrimSpace(expr[1 : len(expr)-1])
	}
	return expr
}

func mysqlManagedPrimaryKeyColumns(rows []mysqlManagedPrimaryKeyRow) ([]string, error) {
	columns := make([]string, 0, len(rows))
	for _, row := range rows {
		if !row.column.Valid || strings.TrimSpace(row.column.String) == "" {
			return nil, mysqlManagedUnconfirmed("primary key has a functional column")
		}
		if row.subPart.Valid {
			return nil, mysqlManagedUnconfirmed("primary key uses a prefix column")
		}
		if !strings.EqualFold(row.indexType, mysqlManagedIndexTypeBTREE) {
			return nil, mysqlManagedUnconfirmed("primary key is not BTREE")
		}
		if row.nonUnique != 0 {
			return nil, mysqlManagedUnconfirmed("primary key is not unique")
		}
		columns = append(columns, row.column.String)
	}
	return columns, nil
}

func applyMySQLManagedIndexVisibility(family string, indexes map[string]*mysqlManagedIndexFacts, columns []string, rows []mysqlManagedVisibilityRow) error {
	if len(indexes) == 0 {
		return nil
	}
	keyColumn, visibleColumn, ignoredColumn := false, false, false
	for _, column := range columns {
		switch strings.ToLower(column) {
		case "key_name":
			keyColumn = true
		case "visible":
			visibleColumn = true
		case mysqlManagedIndexIgnored:
			ignoredColumn = true
		}
	}
	if !keyColumn {
		return mysqlManagedUnconfirmed("index key column is unavailable")
	}
	switch family {
	case mysqlManagedFamilyMySQL:
		if !visibleColumn {
			return mysqlManagedUnconfirmed("mysql index visibility is unavailable")
		}
	case mysqlManagedFamilyMariaDB:
		if !ignoredColumn {
			return mysqlManagedUnconfirmed("mariadb index ignored state is unavailable")
		}
	default:
		return mysqlManagedUnconfirmed("index visibility family is unknown")
	}

	seen := make(map[string]struct{}, len(indexes))
	for _, row := range rows {
		entry, expected := indexes[row.name]
		if !expected {
			continue
		}
		seen[row.name] = struct{}{}
		switch family {
		case mysqlManagedFamilyMySQL:
			if !row.visible.Valid {
				return mysqlManagedUnconfirmed("mysql index visibility is unknown")
			}
			if !strings.EqualFold(row.visible.String, "YES") {
				return mysqlManagedUnconfirmed("mysql index is invisible")
			}
			entry.visible, entry.ignored = true, false
		case mysqlManagedFamilyMariaDB:
			if !row.ignored.Valid {
				return mysqlManagedUnconfirmed("mariadb index ignored state is unknown")
			}
			if !strings.EqualFold(row.ignored.String, "NO") {
				return mysqlManagedUnconfirmed("mariadb index is ignored")
			}
			entry.visible, entry.ignored = true, false
		}
	}
	for name := range indexes {
		if _, ok := seen[name]; !ok {
			return mysqlManagedUnconfirmed("expected index row is unavailable")
		}
	}
	return nil
}

func mysqlManagedUnconfirmed(reason string) error {
	return fmt.Errorf("%w: %s", ErrMySQLManagedTableUnconfirmed, reason)
}
