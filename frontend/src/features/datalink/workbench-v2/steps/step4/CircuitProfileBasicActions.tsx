import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { BasicRecordingProgress } from './BasicRecordingProgress';
import { safeBasicNextAction, safeBasicOperationMessage } from '../../state/basicRecordingMessages';
import type { useCircuitProfileStart } from './useCircuitProfileStart';

type StartFlow = ReturnType<typeof useCircuitProfileStart>;

/** The recording-start progress of the profile panel; a prepare-only context has none. */
export const CircuitProfileStartStatus: React.FC<{ flow: StartFlow; primary?: boolean; onContinueToRuntime?: (deviceIds: string[]) => void;
  deviceNames?: Record<string, string>; groupNames?: Record<string, string> }> = ({ flow, primary, onContinueToRuntime, deviceNames, groupNames }) => {
  const { t } = useTranslation('workbench-v2');
  const { operation, status, queryFailed } = flow;
  if (!status) return null;
  return (
    <div className="space-y-2 rounded border border-slate-800 bg-slate-950/60 p-3 text-xs">
      {!operation && <p data-testid="basic-recording-status" className="font-semibold text-slate-100">{t(`step4.basic.operation.${status}`)}</p>}
      {queryFailed && <p className="text-rose-300">{t('step4.basic.query_failed')}</p>}
      {operation && (
        <BasicRecordingProgress
          operation={operation} intentMatchesCurrent={true}
          completion={flow.completion} primary={primary} onCommit={onContinueToRuntime} deviceNames={deviceNames} groupNames={groupNames}
          statusReason={safeBasicOperationMessage(operation, t)} nextAction={safeBasicNextAction(operation, t)}
        />
      )}
    </div>
  );
};

export interface CircuitProfileBasicActionsProps {
  guided?: boolean;
  flow: StartFlow; canStart: boolean; readonly: boolean; batchBusy: boolean;
  /** Target preparation only: the only action is to go back; nothing here starts, applies or stops a writer. */
  prepareOnly: boolean; onClose?: () => void;
  onContinueToRuntime?: (deviceIds: string[]) => void; onOpenAdvanced?: (groupId?: string) => void;
}

export const CircuitProfileBasicActions: React.FC<CircuitProfileBasicActionsProps> = ({ flow, canStart, readonly, batchBusy, prepareOnly, onClose, onContinueToRuntime, onOpenAdvanced, guided }) => {
  const { t } = useTranslation('workbench-v2');
  const { operationId } = flow;
  if (prepareOnly) {
    return (
      <div className="flex flex-wrap items-center gap-3">
        <button type="button" data-testid="circuit-profile-prepare-close" disabled={flow.busy || batchBusy} onClick={onClose} className="rounded border border-slate-600 px-3 py-2 text-xs text-slate-200 disabled:opacity-40">{t('step4.basic.migration.close')}</button>
      </div>
    );
  }
  return (
    <div className="flex flex-wrap items-center gap-3">
      {!guided && <button type="button" data-testid="circuit-profile-start" disabled={!canStart || flow.busy || flow.unresolved || batchBusy} onClick={() => void flow.start(false)} className="rounded bg-cyan-600 px-4 py-2 text-xs font-semibold text-white disabled:cursor-not-allowed disabled:opacity-40">
        {flow.busy ? t('step4.basic.starting') : t('step4.basic.start')}
      </button>}
      {flow.canRetrySameRequest && <button type="button" data-testid="basic-recording-retry" disabled={readonly || flow.busy || batchBusy} onClick={() => void flow.start(true)} className="rounded border border-cyan-700 px-3 py-2 text-xs text-cyan-200 disabled:opacity-40">{t('step4.basic.retry')}</button>}
      {!guided && operationId && <button type="button" data-testid="basic-recording-check" disabled={flow.busy || batchBusy} onClick={() => void flow.check()} className="rounded border border-slate-700 px-3 py-2 text-xs text-slate-200">{t('step4.basic.check')}</button>}
      {!guided && onContinueToRuntime && flow.completion.state === 'available' && flow.completion.handoffDeviceId && <button type="button" data-testid="basic-recording-continue-runtime" onClick={() => { if (flow.completion.state === 'available' && flow.completion.handoffDeviceId) onContinueToRuntime([flow.completion.handoffDeviceId]); }} className="rounded border border-emerald-700 px-3 py-2 text-xs text-emerald-200">{t('step4.basic.continue_runtime')}</button>}
      {onOpenAdvanced && <button type="button" data-testid="basic-recording-open-advanced" onClick={() => onOpenAdvanced()} className="rounded border border-slate-700 px-3 py-2 text-xs text-slate-300">{t('step4.basic.open_advanced')}</button>}
    </div>
  );
};
