package dbtarget

import (
	"fmt"
	"regexp"
	"strings"
)

var mysqlManagedNumericDefaultPattern = regexp.MustCompile(`^[+-]?\d+(\.\d+)?([eE][+-]?\d+)?$`)

// decodeMySQLManagedDefault turns a MariaDB information_schema COLUMN_DEFAULT
// into its actual value. MariaDB reports string defaults as SQL literals
// (for example 'gateway'), while MySQL reports the unquoted value. Only a
// complete, simple SQL string literal or a numeric literal is accepted here;
// expressions and backslash escapes stay unconfirmed instead of being guessed.
func decodeMySQLManagedDefault(family string, raw *string) (*string, error) {
	if raw == nil {
		return nil, nil
	}
	if family != mysqlManagedFamilyMariaDB {
		return raw, nil
	}
	value := strings.TrimSpace(*raw)
	if strings.EqualFold(value, "NULL") {
		return nil, nil
	}
	if value == "" {
		return nil, unsupportedMySQLManagedDefault()
	}
	if strings.HasPrefix(value, "_") {
		quote := strings.IndexByte(value, '\'')
		if quote < 2 || !isMySQLManagedCharsetName(value[1:quote]) {
			return nil, unsupportedMySQLManagedDefault()
		}
		value = value[quote:]
	}
	if strings.HasPrefix(value, "'") {
		decoded, ok := decodeMySQLManagedStringLiteral(value)
		if !ok {
			return nil, unsupportedMySQLManagedDefault()
		}
		return &decoded, nil
	}
	if mysqlManagedNumericDefaultPattern.MatchString(value) {
		return &value, nil
	}
	return nil, unsupportedMySQLManagedDefault()
}

func decodeMySQLManagedStringLiteral(value string) (string, bool) {
	if len(value) < 2 || value[0] != '\'' || value[len(value)-1] != '\'' {
		return "", false
	}
	var decoded strings.Builder
	for i := 1; i < len(value)-1; i++ {
		switch value[i] {
		case '\'':
			if i+1 < len(value)-1 && value[i+1] == '\'' {
				decoded.WriteByte('\'')
				i++
				continue
			}
			return "", false
		case '\\':
			// Backslash behavior depends on sql_mode; refusing is safer than
			// decoding the same bytes two different ways.
			return "", false
		default:
			decoded.WriteByte(value[i])
		}
	}
	return decoded.String(), true
}

func isMySQLManagedCharsetName(value string) bool {
	if value == "" {
		return false
	}
	for i, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || (i > 0 && r >= '0' && r <= '9') {
			continue
		}
		return false
	}
	return true
}

func unsupportedMySQLManagedDefault() error {
	return fmt.Errorf("%w: unsupported column default expression", ErrMySQLManagedTableUnconfirmed)
}
