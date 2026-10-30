package datalink

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMixedProfileMigrationRestoreFailureDiscardsConnection(t *testing.T) {
	for _, reason := range []string{"exec", "read", "mismatch", "none"} {
		t.Run(reason, func(t *testing.T) {
			var executed []string
			discarded := 0
			failure := errors.New("restore failed")
			err := restoreSQLiteMigrationPragmas([]sqliteMigrationPragma{{"legacy_alter_table", 0}, {"foreign_keys", 1}}, func(query string) error {
				executed = append(executed, query)
				if reason == "exec" && query == "PRAGMA foreign_keys=1" {
					return failure
				}
				return nil
			}, func(query string) (int, error) {
				if query == "PRAGMA foreign_keys" {
					if reason == "read" {
						return 0, failure
					}
					if reason == "mismatch" {
						return 0, nil
					}
					return 1, nil
				}
				return 0, nil
			}, func() error { discarded++; return nil })
			require.Equal(t, []string{"PRAGMA legacy_alter_table=0", "PRAGMA foreign_keys=1"}, executed)
			if reason == "none" {
				require.NoError(t, err)
				require.Zero(t, discarded)
			} else {
				require.Error(t, err)
				require.Equal(t, 1, discarded)
			}
		})
	}
}
