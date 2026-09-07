import { useTranslation } from 'react-i18next';
import type { TFunction } from 'i18next';
import { GitCommit, CheckCircle2, AlertTriangle, ShieldAlert } from 'lucide-react';
import type {
  RuntimeFlowDiagnostic,
  DatabaseDeliveryDiagnostic,
  ModbusShareDeliveryDiagnostic,
} from '../../../types/runtimeDiagnostics';
import { resolveFailureCopy, selectDiagnostics } from './RuntimeDiagnosticsHelpers';

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
            全鏈路管線零異常
            <span className="sr-only" data-testid="runtime-diagnostics-action">
              {failureCopy.action}
            </span>
          </span>
        )}
      </div>

      {/* 4 階段資料管線流動視覺圖 */}
      <div className="mb-5 rounded-2xl border border-slate-800/90 bg-slate-950/70 p-4">
        <div className="mb-3 text-xs font-semibold uppercase tracking-wider text-slate-400">
          全鏈路資料管線流程 (Data Pipeline)
        </div>
        <div className="grid grid-cols-2 gap-2 sm:grid-cols-4">
          <PipelineStageNode name="1. PLC 採集" isError={failedDiagnostic?.failure_stage === 'collector'} />
          <PipelineStageNode name="2. 規則映射" isError={failedDiagnostic?.failure_stage === 'mapping'} />
          <PipelineStageNode name="3. 時序投影" isError={failedDiagnostic?.failure_stage === 'runtime_projection'} />
          <PipelineStageNode name="4. 目標交付" isError={Boolean(failedDiagnostic && (failedDiagnostic.failure_stage === 'database_delivery' || failedDiagnostic.failure_stage === 'modbus_share_delivery'))} />
        </div>
      </div>

      {/* 正常狀態 vs 異常狀態 呈現 */}
      {isHealthy ? (
        <div className="rounded-2xl border border-emerald-500/20 bg-emerald-500/5 p-4">
          <div className="flex items-center gap-3">
            <CheckCircle2 className="h-5 w-5 text-emerald-400 shrink-0" />
            <div>
              <h3 className="text-sm font-semibold text-emerald-200">
                管線資料流運作正常 (Zero Pipeline Errors)
              </h3>
              <p className="mt-0.5 text-xs text-slate-400">
                所有點位皆已持續對齊並順利交付至資料庫，無任何讀取逾時或寫入失敗。
              </p>
            </div>
          </div>
          <details className="mt-3 border-t border-emerald-500/10 pt-3">
            <summary className="cursor-pointer text-xs font-sans text-slate-400 hover:text-slate-200 transition select-none">
              檢視底層點位交付診斷清單 ({scopedDiagnostics.length}) ▾
            </summary>
            <div className="mt-2.5 space-y-1.5">
              {scopedDiagnostics.map((diagnostic) => (
                <DiagnosticCard
                  key={`${diagnostic.scope}:${diagnostic.last_failure_at ?? diagnostic.last_success_at ?? ''}`}
                  diagnostic={diagnostic}
                  t={t}
                />
              ))}
            </div>
          </details>
        </div>
      ) : (
        <div className="space-y-3">
          {scopedDiagnostics.map((diagnostic) => (
            <DiagnosticCard
              key={`${diagnostic.scope}:${diagnostic.last_failure_at ?? diagnostic.last_success_at ?? ''}`}
              diagnostic={diagnostic}
              t={t}
            />
          ))}
        </div>
      )}
    </section>
  );
}

function PipelineStageNode({ name, isError }: { name: string; isError: boolean }) {
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

  if (!isFailed) {
    return (
      <div
        className="flex flex-wrap items-center justify-between gap-2 rounded-xl border border-slate-800/80 bg-slate-900/40 px-3.5 py-2 font-mono text-xs"
        data-testid="runtime-diagnostics-entry"
      >
        <div className="flex items-center gap-2">
          <CheckCircle2 className="h-3.5 w-3.5 text-emerald-400 shrink-0" />
          <span className="font-semibold text-slate-200">{diagnostic.scope}</span>
          {diagnostic.tag_id && (
            <span className="rounded bg-slate-800 px-1.5 py-0.5 text-[11px] text-cyan-300">
              {diagnostic.tag_id}
            </span>
          )}
        </div>
        <div className="flex items-center gap-3 text-slate-400 text-[11px]">
          <span>{t('diagnostics.latestSuccessfulStage', 'Latest successful stage')}: <strong className="text-emerald-400">{diagnostic.latest_successful_stage ?? 'database_delivery'}</strong></span>
          <span>{t('diagnostics.lastFailureAt', 'Last failure at')}: {diagnostic.last_failure_at ?? diagnostic.last_success_at ?? t('diagnostics.unavailable', 'Unavailable')}</span>
        </div>
        <div className="sr-only">
          <span>{t('diagnostics.failureStage', 'Failure stage')}: {diagnostic.failure_stage ?? t('diagnostics.noFailureStage', 'No recent failure')}</span>
          <span>{t('diagnostics.failureReason', 'Failure reason')}: {t('diagnostics.deliveryObserved', 'Delivery observed')}</span>
        </div>
      </div>
    );
  }

  return (
    <div
      className="rounded-2xl border border-rose-500/30 bg-rose-500/5 p-4 transition"
      data-testid="runtime-diagnostics-entry"
    >
      <dl className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {rows.map((row) => (
          <div key={row.label} className="rounded-xl border border-rose-500/20 bg-slate-900/60 p-3">
            <dt className="text-[11px] uppercase tracking-wider text-rose-300/80">{row.label}</dt>
            <dd className="mt-1 font-mono text-xs font-semibold text-rose-100 break-all">{row.value}</dd>
          </div>
        ))}
      </dl>
    </div>
  );
}
