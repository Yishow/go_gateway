package main

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	datalinkruntime "go-gateway/internal/datalink/runtime"

	"github.com/stretchr/testify/require"
)

type asyncDatabaseTarget struct{}

func (asyncDatabaseTarget) WriteTagValue(context.Context, string, any, time.Time) error {
	return errors.New("async acceptance must use its explicit result")
}

func (asyncDatabaseTarget) AcceptTagValue(context.Context, string, any, time.Time) (bool, error) {
	return true, nil
}

func TestProductionTargetWriter_AsyncAcceptanceDoesNotHideModbusResult(t *testing.T) {
	share := &recordingTargetWriter{}
	writer := fanoutTargetWriter{targets: []datalinkruntime.TargetWriter{asyncDatabaseTarget{}, share}}
	outcomes := writer.WriteTagValueOutcomes(t.Context(), "tag-1", 12, time.Now())
	require.Len(t, outcomes, 2)
	require.NoError(t, outcomes[0].Err)
	require.NoError(t, outcomes[1].Err)
	require.Equal(t, 1, share.calls)
	encoded, err := json.Marshal(outcomes[0])
	require.NoError(t, err)
	require.Contains(t, string(encoded), `"Accepted":true`)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "tag-1", 12, time.Now()), dbtarget.ErrSnapshotAccepted)
}
