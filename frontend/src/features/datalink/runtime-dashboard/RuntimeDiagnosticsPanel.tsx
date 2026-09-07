import { useTranslation } from 'react-i18next';
import type { TFunction } from 'i18next';
import { GitCommit, CheckCircle2, AlertTriangle, ShieldAlert } from 'lucide-react';
import type {
  RuntimeFlowDiagnostic,
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
        className="rounded-3xl border border-slate-800/80 bg-slate-900/70 p-6 shadow-xl backdrop-blur-sm"
        data-testid="runtime-dashboard-diagnostics-panel"
      >
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl border border-slate-700 bg-slate-800/50 text-slate-400">
            <GitCommit className="h-5 w-5" />
          </div>
          <div>
            <h2 className="text-base font-semibold text-slate-50">
              {t('diagnostics.title', 'Runtime diagnostics')}
            </h2>
            <p className="mt-1 text-xs text-slate-400">
              {t('diagnostics.unavailable', 'Diagnostics unavailable')}
            </p>
          </div>
        </div>
      </section>
    );
  }

  const failedDiagnostic = scopedDiagnostics.find((d) => d.failure_stage);
  const primaryDiagnostic = failedDiagnostic ?? scopedDiagnostics[0];
  const failureCode = primaryDiagnostic.failure_code ?? 'runtime_snapshot_unavailable';
  const failureCopy = resolveFailureCopy(failureCode, t);
  const isHealthy = !failedDiagnostic;

  return (
    <section
      className="col-span-full rounded-3xl border border-slate-800/80 bg-slate-900/70 p-6 shadow-xl backdrop-blur-sm"
      data-testid="runtime-dashboard-diagnostics-panel"
    >
      <div className="mb-5 flex flex-wrap items-center justify-between gap-3 border-b border-slate-800/80 pb-4">
        <div className="flex items-center gap-3">
          <div
            className={`flex h-10 w-10 items-center justify-center rounded-xl border ${
              isHealthy
                ? 'border-emerald-500/30 bg-emerald-500/10 text-emerald-300'
                : 'border-rose-500/30 bg-rose-500/10 text-rose-300'
            }`}
          >
            {isHealthy ? <CheckCircle2 className="h-5 w-5" /> : <AlertTriangle className="h-5 w-5" />}
          </div>
          <div>
            <h2 className="text-base font-semibold text-slate-50">
              {t('diagnostics.title', 'Runtime diagnostics')}
            </h2>
            <p className="text-xs text-slate-400">
              {t('diagnostics.description', 'Latest selected-scope failure context.')}
            </p>
          </div>
        </div>

        {failedDiagnostic ? (
          <span
            className="rounded-lg border border-rose-500/30 bg-rose-500/10 px-3 py-1 text-xs font-semibold text-rose-200"
            data-testid="runtime-diagnostics-action"
          >
            {failureCopy.action}
          </span>
        ) : (
          <span className="inline-flex items-center gap-1.5 rounded-lg border border-emerald-500/30 bg-emerald-500/10 px-3 py-1 text-xs font-semibold text-emerald-300">
            <span className="relative flex h-2 w-2">
              <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-400 opacity-75" />
              <span className="relative inline-flex h-2 w-2 rounded-full bg-emerald-500" />
            </span>
            全鏈路管線健康
            <span className="sr-only" data-testid="runtime-diagnostics-action">
              {failureCopy.action}
            </span>
          </span>
        )}
      </div>

      {/* 4 階段資料管線流動視覺圖 */}
      <div className="mb-6 rounded-2xl border border-slate-800/90 bg-slate-950/70 p-4">
        <div className="mb-3 text-xs font-semibold uppercase tracking-wider text-slate-400">
          全鏈路資料管線流程 (Data Pipeline)
        </div>
        <div className="grid grid-cols-2 gap-2 sm:grid-cols-4">
          <PipelineStageNode name="1. PLC 採集" stage="collector" isError={failedDiagnostic?.failure_stage === 'collector'} />
          <PipelineStageNode name="2. 規則映射" stage="mapping" isError={failedDiagnostic?.failure_stage === 'mapping'} />
          <PipelineStageNode name="3. 時序投影" stage="runtime_projection" isError={failedDiagnostic?.failure_stage === 'runtime_projection'} />
          <PipelineStageNode name="4. 目標交付" stage="database_delivery" isError={Boolean(failedDiagnostic && (failedDiagnostic.failure_stage === 'database_delivery' || failedDiagnostic.failure_stage === 'modbus_share_delivery'))} />
        </div>
      </div>

      {/* 診斷條目清單 */}
      <div className="space-y-3">
        {scopedDiagnostics.map((diagnostic) => (
          <DiagnosticCard
            key={`${diagnostic.scope}:${diagnostic.last_failure_at ?? diagnostic.last_success_at ?? ''}`}
            diagnostic={diagnostic}
            t={t}
          />
        ))}
      </div>
    </section>
  );
}

function PipelineStageNode({ name, isError }: { name: string; stage: string; isError: boolean }) {
  return (
    <div
      className={`flex items-center justify-between rounded-xl border p-2.5 transition ${
        isError ? 'border-rose-500/40 bg-rose-500/10 text-rose-200' : 'border-slate-800/80 bg-slate-900/50 text-slate-200'
      }`}
    >
      <div className="text-xs font-medium">{name}</div>
      {isError ? <ShieldAlert className="h-4 w-4 text-rose-400" /> : <CheckCircle2 className="h-4 w-4 text-emerald-400" />}
    </div>
  );
}

function DiagnosticCard({
  diagnostic,
  t,
}: {
  diagnostic: RuntimeFlowDiagnostic;
  t: TFunction<'runtime-dashboard'>;
}) {
  const isFailed = Boolean(diagnostic.failure_stage);
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
    <div
      className={`rounded-2xl border p-4 transition ${
        isFailed ? 'border-rose-500/30 bg-rose-500/5' : 'border-slate-800/90 bg-slate-950/70'
      }`}
      data-testid="runtime-diagnostics-entry"
    >
      <dl className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {rows.map((row) => (
          <div key={row.label} className="rounded-xl border border-slate-800/80 bg-slate-900/40 p-3">
            <dt className="text-[11px] uppercase tracking-wider text-slate-500">{row.label}</dt>
            <dd className="mt-1 font-mono text-xs font-semibold text-slate-200 break-all">{row.value}</dd>
          </div>
        ))}
      </dl>
    </div>
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
