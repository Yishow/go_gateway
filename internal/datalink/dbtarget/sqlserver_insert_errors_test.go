package dbtarget

import (
	"testing"

	mssql "github.com/microsoft/go-mssqldb"
	"github.com/stretchr/testify/require"
)

func TestSQLServerDeadlockRetry(t *testing.T) {
	for _, tc := range []struct {
		name   string
		number int32
		want   InsertErrorClass
	}{
		{"deadlock victim retries the whole frozen transaction", 1205, InsertErrorTransient},
		{"lock timeout retries", 1222, InsertErrorTransient},
		{"permission refusal does not retry forever", 229, InsertErrorTarget},
		{"missing object needs scoped repair", 208, InsertErrorTarget},
		{"bad column needs scope repair", 207, InsertErrorTarget},
		{"duplicate identity isolates the row", 2627, InsertErrorRow},
		{"truncation isolates the row", 2628, InsertErrorRow},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := &GroupInsertError{Phase: InsertPhasePreCommit, Cause: mssql.Error{Number: tc.number, Message: "fixture failure"}}
			require.Equal(t, tc.want, ClassifyInsertError(err))
		})
	}
}
