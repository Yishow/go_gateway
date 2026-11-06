package point_test

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"modernc.org/sqlite"
)

func TestDeviceIntervalPreservesConcurrentPointEdit(t *testing.T) {
	entered := make(chan struct{}, 1)
	release := make(chan struct{})
	var releaseOnce sync.Once
	t.Cleanup(func() { releaseOnce.Do(func() { close(release) }) })
	gateName := fmt.Sprintf("acquisition_gate_%d", time.Now().UnixNano())
	require.NoError(t, sqlite.RegisterScalarFunction(gateName, 0, func(_ *sqlite.FunctionContext, _ []driver.Value) (driver.Value, error) {
		select {
		case entered <- struct{}{}:
		default:
		}
		<-release
		return int64(0), nil
	}))

	db, service := acquisitionSQLiteService(t)
	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()
	shared, err := service.CreatePollingGroup(ctx, point.CreatePollingGroupRequest{Name: "shared", IntervalMs: 1000})
	require.NoError(t, err)
	first, err := service.Create(ctx, point.CreatePointRequest{DeviceID: "device-a", Name: "A-first", Address: "40001", DataType: schema.DataTypeUint16, PollingGroupID: &shared.ID})
	require.NoError(t, err)
	last, err := service.Create(ctx, point.CreatePointRequest{DeviceID: "device-a", Name: "Z-last", Address: "40002", DataType: schema.DataTypeUint16, PollingGroupID: &shared.ID})
	require.NoError(t, err)

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO points (id, device_id, name, address, data_type, mode, polling_group_id, revision_id) VALUES (?, 'device-a', ?, ?, 'uint16', 'read', ?, ?)`)
	require.NoError(t, err)
	for i := range 1500 {
		id := fmt.Sprintf("filler-%04d", i)
		_, err = stmt.ExecContext(ctx, id, fmt.Sprintf("M-%04d", i), fmt.Sprintf("D%d", i), shared.ID, id)
		require.NoError(t, err)
	}
	require.NoError(t, stmt.Close())
	require.NoError(t, tx.Commit())
	_, err = db.ExecContext(ctx, fmt.Sprintf(`CREATE TRIGGER acquisition_gate BEFORE UPDATE OF polling_group_id ON points
		WHEN OLD.id = '%s' BEGIN SELECT %s(); END`, first.ID, gateName))
	require.NoError(t, err)

	type groupResult struct {
		group *schema.PollingGroup
		err   error
	}
	intervalDone := make(chan groupResult, 1)
	go func() {
		group, setErr := service.SetDeviceAcquisitionInterval(ctx, "device-a", 500)
		intervalDone <- groupResult{group: group, err: setErr}
	}()
	select {
	case <-entered:
	case result := <-intervalDone:
		t.Fatalf("interval save finished before the SQLite assignment gate: %v", result.err)
	case <-ctx.Done():
		t.Fatal("interval save did not reach the SQLite assignment gate")
	}

	type editResult struct {
		point *schema.Point
		err   error
	}
	editDone := make(chan editResult, 1)
	editStarted := make(chan struct{})
	go func() {
		close(editStarted)
		newName := "Concurrent edit"
		for {
			edited, editErr := service.Update(ctx, last.ID, point.UpdatePointRequest{Name: &newName})
			if editErr == nil || (!errors.Is(editErr, point.ErrFormatUnavailable) && !errors.Is(editErr, point.ErrFormatConflict)) {
				editDone <- editResult{point: edited, err: editErr}
				return
			}
			select {
			case <-ctx.Done():
				editDone <- editResult{err: editErr}
				return
			case <-time.After(time.Millisecond):
			}
		}
	}()
	<-editStarted
	releaseOnce.Do(func() { close(release) })

	var edited editResult
	select {
	case edited = <-editDone:
	case <-ctx.Done():
		t.Fatal("concurrent point edit did not finish")
	}
	require.NoError(t, edited.err)
	var interval groupResult
	select {
	case interval = <-intervalDone:
	case <-ctx.Done():
		t.Fatal("interval save did not finish")
	}
	require.NoError(t, interval.err)
	got, err := service.GetByID(ctx, last.ID)
	require.NoError(t, err)
	assert.Equal(t, "Concurrent edit", got.Name)
	assert.NotEqual(t, last.RevisionID, got.RevisionID, "a saved point revision must not revert to the older snapshot")
	assert.Equal(t, &interval.group.ID, got.PollingGroupID)
}

func TestDeviceIntervalRejectsPointEditBasedOnOldAssignment(t *testing.T) {
	_, service := acquisitionSQLiteService(t)
	ctx := t.Context()
	shared, err := service.CreatePollingGroup(ctx, point.CreatePollingGroupRequest{Name: "shared", IntervalMs: 1000})
	require.NoError(t, err)
	before, err := service.Create(ctx, point.CreatePointRequest{DeviceID: "device-a", Name: "Original", Address: "40001", DataType: schema.DataTypeUint16, PollingGroupID: &shared.ID})
	require.NoError(t, err)

	group, err := service.SetDeviceAcquisitionInterval(ctx, "device-a", 500)
	require.NoError(t, err)
	newName := "Stale edit"
	_, err = service.UpdateOrdinary(ctx, before, point.UpdatePointRequest{Name: &newName}, false)
	require.ErrorIs(t, err, point.ErrFormatConflict)
	got, err := service.GetByID(ctx, before.ID)
	require.NoError(t, err)
	require.Equal(t, &group.ID, got.PollingGroupID)
	require.Equal(t, "Original", got.Name)
}
