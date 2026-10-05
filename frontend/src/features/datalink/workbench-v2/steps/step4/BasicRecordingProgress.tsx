import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { RecordingStartOperation } from '../../../../../types/studioV2RecordingStart';

interface BasicRecordingProgressProps {
  groupNames?: Record<string, string>;
  deviceNames?: Record<string, string>;
  operation: RecordingStartOperation;
  intentMatchesCurrent: boolean;
  activatedDeviceIds: string[];
  statusReason?: string;
  nextAction?: string;
  onCommit?: (confirmedDeviceIds?: string[]) => void;
  onOpenAdvanced?: () => void;
}

export const BasicRecordingProgress: React.FC<BasicRecordingProgressProps> = ({
  groupNames = {}, deviceNames = {}, operation,
  intentMatchesCurrent,
  activatedDeviceIds,
  statusReason,
  nextAction,
  onCommit,
  onOpenAdvanced,
}) => {
  const { t } = useTranslation('workbench-v2');

  return (
    <div className="space-y-1 text-xs text-slate-300" role={operation.status === 'failed' || operation.status === 'partial' ? 'alert' : 'status'} data-testid="basic-recording-status">
      <p>{t(`step4.basic.operation.${operation.status}`)}</p>
      <p>{t(`step4.basic.stage.${operation.stage}`)}</p>
      {statusReason && <p className="text-amber-300">{statusReason}</p>}
      {nextAction && <p className="text-slate-400">{nextAction}</p>}
      {operation.groups.length > 0 && (
        <ul className="space-y-1 pt-1 text-[11px]" data-testid="basic-recording-group-progress" aria-label={t('step4.basic.progress.groups')}>
          {operation.groups.map((group) => (
            <li key={group.group_id} className="flex flex-wrap gap-x-2 gap-y-1">
              <span className="break-all text-slate-400">{t('step4.basic.progress.group')}: {groupNames[group.group_id] ?? t('step4.basic.progress.group')}</span>
              <details className="text-slate-500"><summary>{t('step4.basic.diagnostics')}</summary>{group.group_id}</details>
              <span>{group.saved ? t('step4.basic.progress.saved') : t('step4.basic.progress.not_saved')}</span>
              <span>{group.ready ? t('step4.basic.progress.readiness_verified') : t('step4.basic.progress.readiness_unverified')}</span>
              <span>{group.applied ? t('step4.basic.progress.applied') : t('step4.basic.progress.not_applied')}</span>
            </li>
          ))}
        </ul>
      )}
      {operation.devices.length > 0 && (
        <ul className="space-y-1 pt-1 text-[11px]" data-testid="basic-recording-device-progress" aria-label={t('step4.basic.progress.devices')}>
          {operation.devices.map((device) => (
            <li key={device.device_id} className="flex flex-wrap gap-x-2 gap-y-1">
              <span className="break-all text-slate-400">{t('step4.basic.progress.device')}: {deviceNames[device.device_id] ?? t('step4.basic.progress.device')}</span>
              <details className="text-slate-500"><summary>{t('step4.basic.diagnostics')}</summary>{device.device_id}</details>
              <span>{device.activated ? t('step4.basic.progress.activated') : t('step4.basic.progress.not_activated')}</span>
            </li>
          ))}
        </ul>
      )}
      {operation.status === 'succeeded' && intentMatchesCurrent
        ? <p className="text-emerald-300">{t('step4.basic.verified_activation')}</p>
        : operation.status === 'succeeded' && <p className="text-amber-300">{t('step4.basic.unconfirmed')}</p>}
      {operation.status === 'succeeded' && intentMatchesCurrent && activatedDeviceIds.length > 0 && onCommit && (
        <button type="button" onClick={() => onCommit(activatedDeviceIds)} className="mt-2 rounded border border-emerald-700 px-3 py-2 text-xs text-emerald-200" data-testid="basic-recording-continue-runtime">
          {t('step4.basic.continue_runtime')}
        </button>
      )}
      {operation.next_action === 'prepare_schema' && onOpenAdvanced && (
        <button type="button" onClick={onOpenAdvanced} className="mt-2 rounded border border-amber-700 px-3 py-2 text-xs text-amber-200" data-testid="basic-recording-open-advanced">
          {t('step4.basic.open_advanced')}
        </button>
      )}
    </div>
  );
};
