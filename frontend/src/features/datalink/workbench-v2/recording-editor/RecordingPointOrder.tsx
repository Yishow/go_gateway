import { ParameterHelp } from '../components/ParameterHelp';
import { useEffect, useMemo, useSyncExternalStore } from 'react';
import { useTranslation } from 'react-i18next';
import { canRestartPointOrder, type RecordingPointOrderPort } from './useRecordingPointOrder';
import { pointFormatErrorKey } from '../../../../services/pointFormatErrors';
import { wordFormats } from '../../../../services/pointFormatParser';
import type { PointWordFormat } from '../../../../types/pointFormat';
interface Props { simple?: boolean; port: RecordingPointOrderPort; workspaceId: string; deviceId: string; pointId: string; disabled: boolean }
/** Canonical source detail only: no mapping fields, local decoding or activation. */
export function RecordingPointOrder({ simple = false, port, workspaceId, deviceId, pointId, disabled }: Props) {
  const { t } = useTranslation('point-format');
  const { t: text } = useTranslation('recording-editor');
  const scope = useMemo(() => ({ workspace_id: workspaceId, device_id: deviceId, point_id: pointId }), [workspaceId, deviceId, pointId]);
  const entry = useSyncExternalStore(port.subscribe, () => port.get(scope));
  const available = port.available(scope);
  useEffect(() => { if (available) void port.load(scope); }, [port, scope, available]);
  const resolution = entry?.saved?.resolution ?? entry?.snapshot?.resolution;
  const canEdit = !disabled && entry?.status === 'ready' && entry.snapshot?.editable && resolution?.resolved;
  const value = entry?.intent?.format_mode === 'override' ? entry.intent.configured_format : entry?.intent?.format_mode ?? (resolution?.format_mode === 'override' ? resolution.configured_format : resolution?.format_mode);
  const busy = entry && ['loading', 'saving', 'inspecting', 'recovering'].includes(entry.status);
  const orderControl = resolution && <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 block"><span className={simple ? 'sr-only' : undefined}>{t('editor.order')}</span>
          <select aria-label={t('editor.order')} value={value ?? ''} disabled={!canEdit} className="ml-2 rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] px-2 py-1 text-[var(--gw-text)]"
            onChange={event => port.choose(scope, event.target.value === 'inherit' ? { format_mode: 'inherit', configured_format: '' } : { format_mode: 'override', configured_format: event.target.value as PointWordFormat })}>
            {resolution.format_mode === 'legacy' && <option value="legacy" disabled>{t('mode.legacy')}</option>}
            <option value="inherit">{t('mode.inherit')}</option>
            {wordFormats.map(format => <option key={format} value={format} disabled={!resolution.applicable}>{format} · {t('mode.override')}</option>)}
          </select>
        </label><ParameterHelp compact={simple} help="point.order" /></div>;
  const orderActions = <div className="flex flex-wrap gap-3">
          {(!simple || resolution?.format_mode === 'override') && <button type="button" disabled={!canEdit} onClick={() => port.choose(scope, { format_mode: 'inherit', configured_format: '' })}>{t('editor.reset')}</button>}
          {(!simple || (entry?.status === 'ready' && entry.intent)) && <button type="button" disabled={!canEdit || !entry?.intent} onClick={() => void port.save(scope)}>{t('editor.save')}</button>}
        </div>;
  const operationStatus = <>      {busy && <p role="status">{t('editor.working')}</p>}
      {entry?.saved && !entry.recovered && <p role="status">{t('savedPending')}</p>}
      {entry?.error && <p role="alert">{t(pointFormatErrorKey(entry.error))}</p>}
      {entry?.operation && <p className="break-all">{t('editor.operation')}: {entry.operation.request_id}</p>}
      {canRestartPointOrder(entry) && <button type="button" disabled={disabled} onClick={() => void port.restart(scope)}>{t('editor.restart')}</button>}
      {entry?.status === 'unknown' && <button type="button" disabled={disabled} onClick={() => void port.inspect(scope)}>{t('editor.inspect')}</button>}
      {entry && ['saved', 'recovery-unknown'].includes(entry.status) && <button type="button" disabled={disabled} onClick={() => void port.recover(scope)}>{t('editor.recover')}</button>}
      {entry?.recovered && <p role="status">{t(`editor.${entry.recovered.runtime_state}`)} · {t(entry.recovered.share_applied ? 'editor.shareApplied' : 'editor.shareDisabled')} {t('editor.notReading')}</p>}</>;
  if (simple) return <section className="recording-point-order-inline" data-testid={`recording-point-order-${pointId}`}>
    <details className="recording-point-order-inline__details"><summary>{t('editor.title')}</summary>
      {resolution ? <>
        <p>{t('editor.order')} · <span data-testid="point-order-effective">{resolution.applicable ? resolution.effective_format || t('editor.native') : t('editor.native')}</span></p>
        {resolution.applicable && <div className="recording-point-order-inline__controls">{orderControl}{orderActions}</div>}
        <p>{t('editor.version')}: {resolution.format_version}</p>
        {!resolution.applicable && <p>{t('editor.notApplicable')}</p>}
        <p>{t(`source.${resolution.effective_source}`)} · {t(`policy.${resolution.effective_policy}`)}</p>
        {resolution.native_compatibility && <p>{t('editor.nativeCompatibility')}</p>}
      </> : <p>{text('compact.formatUnknown')}</p>}
    </details>
    {operationStatus}
  </section>;
  return <details className="mt-1 text-sm text-[var(--gw-muted)]" data-testid={`recording-point-order-${pointId}`}>
    <summary className="cursor-pointer">{t('editor.title')}{resolution && <> · {t(`mode.${resolution.format_mode}`)} · <span data-testid="point-order-effective">{resolution.effective_format || t('editor.native')}</span></>}</summary>
    <div className="mt-2 space-y-2 rounded border border-[var(--gw-control)] p-2">
      {resolution && <>
        <p>{t(`source.${resolution.effective_source}`)} · {t(`policy.${resolution.effective_policy}`)}</p>
        <p>{t('editor.version')}: {resolution.format_version}{resolution.native_compatibility && <> · {t('editor.nativeCompatibility')}</>}</p>
        {!resolution.applicable && <p>{t('editor.notApplicable')}</p>}
        {orderControl}{orderActions}
      </>}
      {operationStatus}
    </div>
  </details>;
}
