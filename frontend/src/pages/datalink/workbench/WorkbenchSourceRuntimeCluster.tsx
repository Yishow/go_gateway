import { useTranslation } from 'react-i18next';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useToggleDeviceStatusMutation } from '../../../hooks/datalink/useDevices';
import { runtimeAPI } from '../../../services/datalink';
import { useWorkbench } from './WorkbenchProvider';
import { useWorkbenchSummary } from './useWorkbenchSummary';
import {
  getWorkbenchDeviceStatusChipClasses,
  getWorkbenchDeviceStatusLabelKey,
} from './workbenchDeviceFormModel';

/**
 * 將未知錯誤轉成可顯示字串。
 *
 * @param error - 擷取到的錯誤
 * @param fallback - 非 Error 時使用的後援文案
 * @returns 錯誤訊息字串
 */
function getErrorMessage(error: unknown, fallback: string): string {
  return error instanceof Error ? error.message : fallback;
}

/**
 * 依收集器狀態回傳狀態指示點的 Tailwind 類名。
 *
 * @param status - 後端 collectors[0].status
 * @returns 圓點背景／陰影類名
 */
function getRuntimeStatusDotClass(status: string | undefined): string {
  switch (status) {
    case 'running':
      return 'bg-emerald-400 shadow-[0_0_10px_rgba(52,211,153,0.45)]';
    case 'warning':
      return 'bg-amber-400 shadow-[0_0_8px_rgba(251,191,36,0.4)]';
    case 'error':
      return 'bg-rose-400 shadow-[0_0_8px_rgba(251,113,133,0.35)]';
    default:
      return 'bg-slate-500';
  }
}

/**
 * 來源步驟專用：於頂欄顯示收集狀態與開關；RUNTIME 標題與收集說明問號在來源步驟指標卡區塊。
 */
export function WorkbenchSourceRuntimeCluster() {
  const { t } = useTranslation();
  const { activeStep, selectedDeviceId, setSourceStepNotice } = useWorkbench();
  const { selectedDevice } = useWorkbenchSummary();
  const queryClient = useQueryClient();
  const toggleDeviceCollectionMutation = useToggleDeviceStatusMutation();

  const runtimeStatusQuery = useQuery({
    queryKey: ['runtime-status', selectedDeviceId],
    queryFn: () => runtimeAPI.getStatus(selectedDeviceId ?? undefined),
    enabled: Boolean(selectedDeviceId) && activeStep === 'source',
    staleTime: 10_000,
  });

  if (activeStep !== 'source' || !selectedDevice) {
    return null;
  }

  const collectorStatus = runtimeStatusQuery.data?.collectors[0]?.status;
  const runtimeStatusDotClass = getRuntimeStatusDotClass(collectorStatus);

  /**
   * 切換目前設備是否由 Runtime 輪詢收集，並將結果寫入共用的來源步驟訊息狀態。
   */
  const handleDeviceCollectionToggle = async () => {
    const wasActive = selectedDevice.status === 'active';

    try {
      await toggleDeviceCollectionMutation.mutateAsync({
        id: selectedDevice.id,
        currentStatus: selectedDevice.status,
      });
      await queryClient.invalidateQueries({ queryKey: ['runtime-status', selectedDeviceId] });
      setSourceStepNotice(
        wasActive ? t('workbench.source.collection.stopped') : t('workbench.source.collection.started'),
      );
    } catch (error) {
      setSourceStepNotice(getErrorMessage(error, t('workbench.source.collection.error')));
    }
  };

  return (
    <>
      <span className="hidden h-6 w-px shrink-0 bg-slate-800/70 sm:block" aria-hidden />
      <div
        className="flex min-w-0 max-w-full flex-wrap items-center gap-2 rounded-lg border border-slate-800/55 bg-slate-950/40 px-2 py-1 shadow-[inset_0_1px_0_0_rgba(255,255,255,0.04)] sm:flex-nowrap sm:gap-2.5 sm:py-1"
        data-testid="workbench-source-runtime-cluster"
      >
        <span
          className={[
            'inline-flex h-6 shrink-0 items-center rounded-md border px-2 text-[11px] font-semibold',
            getWorkbenchDeviceStatusChipClasses(selectedDevice.status),
          ].join(' ')}
          data-testid="context-bar-device-status-chip"
        >
          {t(getWorkbenchDeviceStatusLabelKey(selectedDevice.status))}
        </span>
        <span className="hidden h-4 w-px shrink-0 bg-slate-700/70 sm:block" aria-hidden />
        <div className="flex min-w-0 items-center gap-2">
          <span
            className={`h-2 w-2 shrink-0 rounded-full ${runtimeStatusDotClass}`}
            aria-hidden
          />
          <p
            className="m-0 min-w-0 text-xs font-semibold leading-tight text-slate-100 sm:text-sm"
            data-testid="workbench-runtime-status"
          >
            {collectorStatus === 'running'
              ? t('workbench.runtime.summary.running')
              : collectorStatus === 'warning'
                ? t('workbench.runtime.summary.warning')
                : collectorStatus === 'error'
                  ? t('workbench.runtime.summary.error')
                  : t('workbench.runtime.summary.idle')}
          </p>
        </div>
        <button
          type="button"
          data-testid="source-device-collection-toggle"
          disabled={toggleDeviceCollectionMutation.isPending}
          onClick={() => void handleDeviceCollectionToggle()}
          className={
            selectedDevice.status === 'active'
              ? 'shrink-0 rounded-lg border border-rose-500/40 bg-rose-500/10 px-2.5 py-1 text-[11px] font-semibold text-rose-100 transition hover:border-rose-400/50 hover:bg-rose-500/15 disabled:cursor-not-allowed disabled:opacity-50 sm:px-3 sm:py-1.5 sm:text-xs'
              : 'shrink-0 rounded-lg bg-cyan-500 px-2.5 py-1 text-[11px] font-semibold text-slate-950 shadow-sm shadow-cyan-900/25 transition hover:bg-cyan-400 disabled:cursor-not-allowed disabled:bg-slate-800 disabled:text-slate-500 disabled:shadow-none sm:px-3 sm:py-1.5 sm:text-xs'
          }
        >
          {selectedDevice.status === 'active'
            ? t('workbench.source.collection.stop')
            : t('workbench.source.collection.start')}
        </button>
      </div>
    </>
  );
}
