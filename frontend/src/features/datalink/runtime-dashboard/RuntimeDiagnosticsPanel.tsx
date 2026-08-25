import { useTranslation } from 'react-i18next';
import type { TFunction } from 'i18next';
import type { RuntimeFlowDiagnostic } from '../../../types/runtimeDiagnostics';
import type {
  DatabaseDeliveryDiagnostic,
  ModbusShareDeliveryDiagnostic,
} from '../../../types/runtimeDiagnostics';

interface RuntimeDiagnosticsPanelProps {
  diagnostics?: RuntimeFlowDiagnostic[];
  databaseDelivery?: DatabaseDeliveryDiagnostic[];
  modbusShareDelivery?: ModbusShareDeliveryDiagnostic[];
  selectedDeviceId: string | null;
}

export function RuntimeDiagnosticsPanel({
  diagnostics = [],
  databaseDelivery = [],
  modbusShareDelivery = [],
  selectedDeviceId,
}: RuntimeDiagnosticsPanelProps) {
  const { t } = useTranslation('runtime-dashboard');
  const scopedDiagnostics = selectDiagnostics(
    diagnostics,
    databaseDelivery,
    modbusShareDelivery,
    selectedDeviceId,
  );

  if (scopedDiagnostics.length === 0) {
    return (
      <section
        className="rounded-3xl border border-slate-800 bg-slate-900/70 p-6"
        data-testid="runtime-dashboard-diagnostics-panel"
      >
        <h2 className="text-lg font-semibold text-slate-50">
          {t('diagnostics.title', 'Runtime diagnostics')}
        </h2>
        <p className="mt-3 text-sm text-slate-400">
          {t('diagnostics.unavailable', 'Diagnostics unavailable')}
        </p>
      </section>
    );
  }

  return (
    <section
      className="rounded-3xl border border-slate-800 bg-slate-900/70 p-6"
      data-testid="runtime-dashboard-diagnostics-panel"
    >
      <div className="mb-4">
        <h2 className="text-lg font-semibold text-slate-50">
          {t('diagnostics.title', 'Runtime diagnostics')}
        </h2>
        <p className="mt-1 text-sm text-slate-400">
          {t('diagnostics.description', 'Latest selected-scope failure context.')}
        </p>
        <p className="mt-2 text-sm text-cyan-200" data-testid="runtime-diagnostics-action">
          {resolveFailureCopy(scopedDiagnostics[0].failure_code ?? 'runtime_snapshot_unavailable', t).action}
        </p>
      </div>
      <div className="space-y-3">
        {scopedDiagnostics.map((diagnostic) => (
          <DiagnosticCard key={`${diagnostic.scope}:${diagnostic.last_failure_at ?? diagnostic.last_success_at ?? ''}`} diagnostic={diagnostic} t={t} />
        ))}
      </div>
    </section>
  );
}

function DiagnosticCard({
  diagnostic,
  t,
}: {
  diagnostic: RuntimeFlowDiagnostic;
  t: TFunction<'runtime-dashboard'>;
}) {
  const failureCode = diagnostic.failure_code ?? 'runtime_snapshot_unavailable';
  const failureCopy = resolveFailureCopy(failureCode, t);
  const rows = [
    { label: t('diagnostics.scope', 'Scope'), value: diagnostic.scope },
    { label: t('diagnostics.failureStage', 'Failure stage'), value: diagnostic.failure_stage ?? t('diagnostics.noFailureStage', 'No recent failure') },
    { label: t('diagnostics.failureReason', 'Failure reason'), value: diagnostic.failure_stage ? failureCopy.message : t('diagnostics.deliveryObserved', 'Delivery observed') },
    { label: t('diagnostics.lastFailureAt', 'Last failure at'), value: diagnostic.last_failure_at ?? diagnostic.last_success_at ?? t('diagnostics.unavailable', 'Unavailable') },
    { label: t('diagnostics.latestSuccessfulStage', 'Latest successful stage'), value: diagnostic.latest_successful_stage ?? t('diagnostics.unavailable', 'Unavailable') },
  ];
  if (diagnostic.request_id) {
    rows.push({ label: t('diagnostics.requestId', 'Request ID:'), value: ` ${diagnostic.request_id}` });
  }
  return (
    <dl className="grid gap-3 sm:grid-cols-2" data-testid="runtime-diagnostics-entry">
      {rows.map((row) => (
        <div key={row.label} className="rounded-2xl border border-slate-800 bg-slate-950/70 p-4">
          <dt className="text-xs uppercase tracking-[0.18em] text-slate-500">{row.label}</dt>
          <dd className="mt-2 text-lg font-semibold text-slate-50">{row.value}</dd>
        </div>
      ))}
    </dl>
  );
}

function resolveFailureCopy(
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

function selectDiagnostics(
  diagnostics: RuntimeFlowDiagnostic[],
  databaseDelivery: DatabaseDeliveryDiagnostic[],
  modbusShareDelivery: ModbusShareDeliveryDiagnostic[],
  selectedDeviceId: string | null,
): RuntimeFlowDiagnostic[] {
  const scoped = diagnostics.filter((diagnostic) => matchesDevice(diagnostic.device_id, diagnostic.scope, selectedDeviceId));
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

function matchesDevice(deviceId: string | undefined, scope: string, selectedDeviceId: string | null): boolean {
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
