package runtime

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"

	"github.com/stretchr/testify/require"
)

type outcomeTargetWriter struct {
	outcomes []TargetDeliveryOutcome
}

func (w *outcomeTargetWriter) WriteTagValue(context.Context, string, any, time.Time) error {
	return nil
}

func (w *outcomeTargetWriter) WriteTagValueOutcomes(context.Context, string, any, time.Time) []TargetDeliveryOutcome {
	return append([]TargetDeliveryOutcome(nil), w.outcomes...)
}

func TestHandleCollectedValue_TargetOutcomesKeepDatabaseAndModbusResultsIndependent(t *testing.T) {
	mw := &mockWriter{}
	targetWriter := &outcomeTargetWriter{outcomes: []TargetDeliveryOutcome{
		{Target: TargetDatabase, Err: errors.New("database unavailable")},
		{Target: TargetModbusShare},
	}}
	s := targetOutcomeService(mw, targetWriter)
	s.handleCollectedValue(context.Background(), collector.CollectedValue{PointID: "p1", DeviceID: "dev-1", Value: int16(7), Timestamp: time.Date(2026, 3, 16, 11, 0, 0, 0, time.UTC), Quality: schema.QualityGood})

	database := s.databaseDeliveryDiagnostics("")
	modbus := s.modbusShareDeliveryDiagnostics("")
	require.Len(t, database, 1)
	require.Equal(t, DatabaseDeliveryStatusFailed, database[0].Status)
	require.Len(t, modbus, 1)
	require.Equal(t, ModbusShareDeliveryStatusSucceeded, modbus[0].Status)
}

func TestHandleCollectedValue_TargetOutcomesRecordModbusFailureWithoutHidingDatabaseSuccess(t *testing.T) {
	mw := &mockWriter{}
	targetWriter := &outcomeTargetWriter{outcomes: []TargetDeliveryOutcome{
		{Target: TargetDatabase},
		{Target: TargetModbusShare, Err: NewDeliveryError("modbus_share", "modbus_write", errors.New("listener stopped"))},
	}}
	s := targetOutcomeService(mw, targetWriter)
	s.handleCollectedValue(context.Background(), collector.CollectedValue{PointID: "p1", DeviceID: "dev-1", Value: int16(8), Timestamp: time.Date(2026, 3, 16, 11, 5, 0, 0, time.UTC), Quality: schema.QualityGood})

	database := s.databaseDeliveryDiagnostics("")
	modbus := s.modbusShareDeliveryDiagnostics("")
	require.Len(t, database, 1)
	require.Equal(t, DatabaseDeliveryStatusSucceeded, database[0].Status)
	require.Len(t, modbus, 1)
	require.Equal(t, ModbusShareDeliveryStatusFailed, modbus[0].Status)
	require.Equal(t, "modbus_write", modbus[0].Stage)
}

func targetOutcomeService(writer storage.Writer, target TargetWriter) *Service {
	return &Service{config: Config{UpdatePointState: false}, writer: writer, target: target, mappingIndex: map[string][]mappingBinding{"p1": {{TagID: "tag-1", TagDataType: schema.DataTypeInt16}}}}
}
