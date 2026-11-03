import { useRef } from 'react';
import type { GroupCandidate } from '../state/writeGroup/candidates';
import { sourceKey } from './recordingPickerState';
import { useTranslation } from 'react-i18next';
import type { useRecordingPointNameBuffer } from './useRecordingPointNameBuffer';
export function RecordingPointNameContent({ controller, onBack, onSave, disabled, candidate, onBegin }: { controller: ReturnType<typeof useRecordingPointNameBuffer>; onBack(): void; onSave(): void; disabled: boolean; candidate?: GroupCandidate; onBegin?: () => void }) {
  const { t } = useTranslation('recording-editor');
  const input = useRef<HTMLInputElement>(null);
  const value = !candidate || controller.value && sourceKey(controller.value.candidate) === sourceKey(candidate) ? controller.value : undefined, status = value?.outcome?.status;
  if (candidate) return <div className="recording-point-name-inline" data-testid={`inline-point-name-${candidate.point_id}`}>
    <input ref={input} aria-label={`${t('name.label')} ${candidate.label}`} data-picker-focus={`name:${sourceKey(candidate)}`} value={value?.name ?? candidate.label}
      disabled={disabled} readOnly={!value || controller.pending || status !== 'ready'} aria-busy={Boolean(value && controller.pending)}
      onPointerDown={event => { if (event.currentTarget.ownerDocument.activeElement === event.currentTarget && !value && !disabled) onBegin?.(); }}
      onFocus={() => { if (!value && !disabled) onBegin?.(); }} onChange={event => controller.change(event.target.value)} />
    {value && <>
      {controller.pending && <p role="status">{t('name.pending')}</p>}
      {(value.failed || status === 'blocked') && <p role="alert">{t('name.blocked')}</p>}
      {status === 'unknown' && <><p role="alert">{t('name.unknown')}</p><button type="button" onClick={() => void controller.inspect()}>{t('name.inspect')}</button>{value.outcome?.currentName !== undefined && <p>{t('name.current', { name: value.outcome.currentName })}</p>}</>}
      <div className="recording-point-name-inline__actions"><button type="button" disabled={controller.pending} onClick={() => { input.current?.focus(); onBack(); }}>{t(status === 'unknown' ? 'name.finishInspect' : 'picker.cancel')}</button>
        <button type="button" className="gw-primary" disabled={disabled || status !== 'ready' || !value.name.trim() || !controller.dirty} onClick={() => { input.current?.focus(); onSave(); }}>{t('name.save')}</button></div>
    </>}
  </div>;
  return <><div className="gw-modal-body">
    {controller.pending && <p role="status">{t('name.pending')}</p>}
    {(value?.failed || status === 'blocked') && <p role="alert">{t('name.blocked')}</p>}
    {status === 'unknown' && <><p role="alert">{t('name.unknown')}</p><button type="button" onClick={() => void controller.inspect()}>{t('name.inspect')}</button>{value?.outcome?.currentName !== undefined && <p>{t('name.current', { name: value.outcome.currentName })}</p>}</>}
    <label>{t('name.label')}<input autoFocus aria-label={t('name.label')} value={value?.name ?? ''} disabled={disabled || status !== 'ready'} onChange={event => controller.change(event.target.value)} /></label>
  </div><footer className="gw-modal-footer"><button type="button" disabled={controller.pending} onClick={onBack}>{t(status === 'unknown' ? 'source.back' : 'picker.cancel')}</button>
    <button type="button" className="gw-primary" disabled={disabled || status !== 'ready' || !value?.name.trim() || !controller.dirty} onClick={onSave}>{t('name.save')}</button>
  </footer></>;
}
