import type { TFunction } from 'i18next';
import type {
  RuntimeFlowDiagnostic,
  DatabaseDeliveryDiagnostic,
  ModbusShareDeliveryDiagnostic,
} from '../../../types/runtimeDiagnostics';

export function resolveFailureCopy(
  failureCode: string,
  t: TFunction<'runtime-dashboard'>,
): { message: string; action: string } {
  if (failureCode === 'runtime_snapshot_unavailable') {
    return {
      message: t('errors.runtime_snapshot_unavailable', 'Runtime snapshot is currently unavailable.'),
      action: t('errors.runtime_snapshot_unavailable_action', 'Retry snapshot'),
    };
  }
  if (failureCode === 'modbus_share_delivery') {
    return {
      message: t('errors.modbus_share_delivery', 'Modbus Share delivery is unavailable.'),
      action: t('errors.modbus_share_delivery_action', 'Retry Modbus Share delivery'),
    };
  }
  return {
    message: t('errors.generic_failure', 'Runtime diagnostics are unavailable.'),
    action: t('errors.retry', 'Retry snapshot'),
  };
}

export function selectDiagnostics(
  diagnostics: RuntimeFlowDiagnostic[],
  databaseDelivery: DatabaseDeliveryDiagnostic[],
  modbusShareDelivery: ModbusShareDeliveryDiagnostic[],
  selectedDeviceId: string | null,
): RuntimeFlowDiagnostic[] {
  const scoped = diagnostics.filter((diagnostic) =>
    matchesDevice(diagnostic.device_id, diagnostic.scope, selectedDeviceId),
  );
  const database = databaseDelivery
    .filter((item) => matchesDevice(item.device_id, `tag:${item.tag_id}`, selectedDeviceId))
    .map(toDatabaseDiagnostic);
  const share = modbusShareDelivery
    .filter((item) => matchesDevice(item.device_id, `tag:${item.tag_id}`, selectedDeviceId))
    .filter((item) => item.status !== 'disabled')
    .map(toModbusShareDiagnostic);
  return [...scoped, ...database, ...share].sort((left, right) =>
    diagnosticTimestamp(right).localeCompare(diagnosticTimestamp(left)),
  );
}

function matchesDevice(
  deviceId: string | undefined,
  scope: string,
  selectedDeviceId: string | null,
): boolean {
  return !selectedDeviceId || deviceId === selectedDeviceId || scope === `device:${selectedDeviceId}`;
}

function toDatabaseDiagnostic(item: DatabaseDeliveryDiagnostic): RuntimeFlowDiagnostic {
  const failed = item.status === 'failed';
  return {
    scope: `database_delivery:${item.tag_id}`,
    device_id: item.device_id,
    point_id: item.point_id,
    tag_id: item.tag_id,
    last_success_at: item.last_success_at,
    last_failure_at: item.last_failure_at,
    latest_successful_stage: failed ? undefined : 'database_delivery',
    failure_stage: failed ? 'database_delivery' : undefined,
    failure_code: failed ? 'database_delivery' : undefined,
    failure_reason: item.last_failure_reason ?? item.error,
    stages: item.stages.map((stage) => ({ stage, status: failed ? 'failed' : 'success' })),
  };
}

function toModbusShareDiagnostic(item: ModbusShareDeliveryDiagnostic): RuntimeFlowDiagnostic {
  const failed = item.status === 'failed' || (item.status === undefined && Boolean(item.error));
  return {
    scope: `modbus_share_delivery:${item.tag_id}`,
    device_id: item.device_id,
    point_id: item.point_id,
    tag_id: item.tag_id,
    last_success_at: failed ? undefined : item.observed_at,
    last_failure_at: failed ? item.observed_at : undefined,
    latest_successful_stage: failed ? undefined : 'modbus_share_delivery',
    failure_stage: failed ? 'modbus_share_delivery' : undefined,
    failure_code: failed ? 'modbus_share_delivery' : undefined,
    failure_reason: item.error,
    stages: [{ stage: 'modbus_share_delivery', status: failed ? 'failed' : 'success' }],
  };
}

function diagnosticTimestamp(diagnostic: RuntimeFlowDiagnostic): string {
  return diagnostic.last_failure_at ?? diagnostic.last_success_at ?? '';
}
