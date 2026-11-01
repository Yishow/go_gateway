import { useEffect, useMemo, useSyncExternalStore } from 'react';
import { useTranslation } from 'react-i18next';
import { canRestartPointOrder, type RecordingPointOrderPort } from './useRecordingPointOrder';
import { pointFormatErrorKey } from '../../../../services/pointFormatErrors';
import { wordFormats } from '../../../../services/pointFormatParser';
import type { PointWordFormat } from '../../../../types/pointFormat';
interface Props { port: RecordingPointOrderPort; workspaceId: string; deviceId: string; pointId: string; disabled: boolean }
/** Canonical source detail only: no mapping fields, local decoding or activation. */
export function RecordingPointOrder({ port, workspaceId, deviceId, pointId, disabled }: Props) {
  const { t } = useTranslation('point-format');
  const scope = useMemo(() => ({ workspace_id: workspaceId, device_id: deviceId, point_id: pointId }), [workspaceId, deviceId, pointId]);
  const entry = useSyncExternalStore(port.subscribe, () => port.get(scope));
  const available = port.available(scope);
  useEffect(() => { if (available) void port.load(scope); }, [port, scope, available]);
  const resolution = entry?.saved?.resolution ?? entry?.snapshot?.resolution;
  const canEdit = !disabled && entry?.status === 'ready' && entry.snapshot?.editable && resolution?.resolved;
  const value = entry?.intent?.format_mode === 'override' ? entry.intent.configured_format : entry?.intent?.format_mode ?? (resolution?.format_mode === 'override' ? resolution.configured_format : resolution?.format_mode);
  const busy = entry && ['loading', 'saving', 'inspecting', 'recovering'].includes(entry.status);
  return <details className="mt-1 text-xs text-slate-400" data-testid={`recording-point-order-${pointId}`}>
    <summary className="cursor-pointer">{t('editor.title')}{resolution && <> · {t(`mode.${resolution.format_mode}`)} · <span data-testid="point-order-effective">{resolution.effective_format || t('editor.native')}</span></>}</summary>
    <div className="mt-2 space-y-2 rounded border border-slate-800 p-2">
      {resolution && <>
        <p>{t(`source.${resolution.effective_source}`)} · {t(`policy.${resolution.effective_policy}`)}</p>
        <p>{t('editor.version')}: {resolution.format_version}{resolution.native_compatibility && <> · {t('editor.nativeCompatibility')}</>}</p>
        {!resolution.applicable && <p>{t('editor.notApplicable')}</p>}
        <label className="block">{t('editor.order')}
          <select aria-label={t('editor.order')} value={value ?? ''} disabled={!canEdit} className="ml-2 rounded border border-slate-700 bg-slate-950 px-2 py-1 text-slate-100"
            onChange={event => port.choose(scope, event.target.value === 'inherit' ? { format_mode: 'inherit', configured_format: '' } : { format_mode: 'override', configured_format: event.target.value as PointWordFormat })}>
            {resolution.format_mode === 'legacy' && <option value="legacy" disabled>{t('mode.legacy')}</option>}
            <option value="inherit">{t('mode.inherit')}</option>
            {wordFormats.map(format => <option key={format} value={format} disabled={!resolution.applicable}>{format} · {t('mode.override')}</option>)}
          </select>
        </label>
        <div className="flex flex-wrap gap-3">
          <button type="button" disabled={!canEdit} onClick={() => port.choose(scope, { format_mode: 'inherit', configured_format: '' })}>{t('editor.reset')}</button>
          <button type="button" disabled={!canEdit || !entry?.intent} onClick={() => void port.save(scope)}>{t('editor.save')}</button>
        </div>
      </>}
      {busy && <p role="status">{t('editor.working')}</p>}
      {entry?.saved && !entry.recovered && <p role="status">{t('savedPending')}</p>}
      {entry?.error && <p role="alert">{t(pointFormatErrorKey(entry.error))}</p>}
      {entry?.operation && <p className="break-all">{t('editor.operation')}: {entry.operation.request_id}</p>}
      {canRestartPointOrder(entry) && <button type="button" disabled={disabled} onClick={() => void port.restart(scope)}>{t('editor.restart')}</button>}
      {entry?.status === 'unknown' && <button type="button" disabled={disabled} onClick={() => void port.inspect(scope)}>{t('editor.inspect')}</button>}
      {entry && ['saved', 'recovery-unknown'].includes(entry.status) && <button type="button" disabled={disabled} onClick={() => void port.recover(scope)}>{t('editor.recover')}</button>}
      {entry?.recovered && <p role="status">{t(`editor.${entry.recovered.runtime_state}`)} · {t(entry.recovered.share_applied ? 'editor.shareApplied' : 'editor.shareDisabled')} {t('editor.notReading')}</p>}
    </div>
  </details>;
}
