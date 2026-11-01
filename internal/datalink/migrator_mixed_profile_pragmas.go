package datalink

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

type sqliteMigrationPragma struct {
	name  string
	value int
}

// restoreSQLiteMigrationPragmas attempts every restore even after one fails.
// Failed or unverifiable restoration always discards the dedicated connection.
func restoreSQLiteMigrationPragmas(settings []sqliteMigrationPragma, exec func(string) error, read func(string) (int, error), discard func() error) error {
	var result error
	for _, setting := range settings {
		if err := exec(fmt.Sprintf("PRAGMA %s=%d", setting.name, setting.value)); err != nil {
			result = errors.Join(result, err)
			continue
		}
		actual, err := read("PRAGMA " + setting.name)
		if err != nil {
			result = errors.Join(result, err)
		} else if actual != setting.value {
			result = errors.Join(result, fmt.Errorf("%s restore mismatch", setting.name))
		}
	}
	if result != nil {
		if discardErr := discard(); discardErr != nil && !errors.Is(discardErr, driver.ErrBadConn) {
			result = errors.Join(result, migrationCleanupError{cause: discardErr})
		}
	}
	return result
}
