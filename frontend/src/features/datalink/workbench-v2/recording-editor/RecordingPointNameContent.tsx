import { useEffect, useRef } from 'react';
import type { GroupCandidate } from '../state/writeGroup/candidates';
import { sourceKey } from './recordingPickerState';
import { useTranslation } from 'react-i18next';
import type { useRecordingPointNameBuffer } from './useRecordingPointNameBuffer';
export function RecordingPointNameContent({ controller, onBack, onSave, disabled, candidate, onBegin, canAutoSave }: { controller: ReturnType<typeof useRecordingPointNameBuffer>; onBack(): void; onSave(): void; disabled: boolean; candidate?: GroupCandidate; onBegin?: () => void; canAutoSave?: () => boolean }) {
  const { t } = useTranslation('recording-editor');
  const { t: mappingText } = useTranslation('mapping-errors');
  const input = useRef<HTMLInputElement>(null);
  const value = !candidate || controller.value && sourceKey(controller.value.candidate) === sourceKey(candidate) ? controller.value : undefined, status = value?.outcome?.status;
  const saveAfterPrepare = useRef<'blur' | 'enter' | false>(false);
  useEffect(() => {
    if (!value || value.preparing || !saveAfterPrepare.current || canAutoSave?.() === false) return;
    saveAfterPrepare.current = false;
    if (status !== 'ready' || value.failed) return;
    if (controller.dirty && value.name.trim()) onSave();
    else if (!controller.dirty) onBack();
  }, [value, status, controller.dirty, onSave, onBack, canAutoSave]);
  if (candidate) return <div className="recording-point-name-inline" data-testid={`inline-point-name-${candidate.point_id}`}>
    <input ref={input} aria-label={`${t('name.label')} ${candidate.label}`} data-picker-focus={`name:${sourceKey(candidate)}`} value={value?.name ?? candidate.label}
      disabled={disabled} readOnly={!value || value.failed || !(value.preparing && !value.outcome || status === 'ready' && !controller.pending)} aria-busy={Boolean(value && controller.pending)}
      aria-invalid={Boolean(value && controller.dirty && !value.name.trim())}
      onPointerDown={event => { if (event.currentTarget.ownerDocument.activeElement === event.currentTarget && !value && !disabled) onBegin?.(); }}
      onFocus={() => { if (saveAfterPrepare.current === 'blur') saveAfterPrepare.current = false; if (!value && !disabled) onBegin?.(); }} onChange={event => controller.change(event.target.value)}
      onBlur={() => { if (!value || value.failed || canAutoSave?.() === false) return; if (value.preparing && !value.outcome) { if (saveAfterPrepare.current !== 'enter') saveAfterPrepare.current = 'blur'; }
        else if (status === 'ready') { if (controller.dirty && value.name.trim()) onSave(); else if (!controller.dirty) onBack(); } }}
      onKeyDown={event => { if (event.nativeEvent.isComposing) return;
        if (event.key === 'Enter' && value && !value.failed) { event.preventDefault(); if (value.preparing && !value.outcome) saveAfterPrepare.current = 'enter'; else if (status === 'ready' && controller.dirty && value.name.trim()) onSave(); }
        if (event.key === 'Escape' && value && (status === 'ready' || value.preparing && !value.outcome)) { event.preventDefault(); event.stopPropagation(); saveAfterPrepare.current = false; onBack(); }
      }} />
    {value && <>
      {controller.pending && <p role="status" className="recording-point-name-inline__status">{t(value.preparing ? 'name.loading' : 'name.pending')}</p>}
      {controller.dirty && !value.name.trim() && <p role="alert">{t('name.required')}</p>}
      {(value.failed || status === 'blocked') && <p role="alert">{t('name.blocked')}</p>}
      {status === 'unknown' && <><p role="alert">{value?.outcome?.reason === 'source-review' ? mappingText('name_source_review') : t('name.unknown')}</p><button type="button" disabled={disabled || controller.pending} onClick={() => void controller.inspect()}>{t('name.inspect')}</button>{value?.outcome?.canDiscardRejected && <button type="button" disabled={disabled || controller.pending} onClick={() => void controller.discardRejected()}>{mappingText('cancel_rejected_name')}</button>}{value.outcome?.currentName !== undefined && <p>{t('name.current', { name: value.outcome.currentName })}</p>}</>}
      {(value.failed || status === 'blocked' || status === 'unknown' || controller.dirty && !value.name.trim()) && <div className="recording-point-name-inline__actions"><button type="button" disabled={controller.pending} onClick={() => { input.current?.focus(); onBack(); }}>{t(status === 'unknown' ? 'name.finishInspect' : 'picker.cancel')}</button></div>}
    </>}
  </div>;
  return <><div className="gw-modal-body">
    {controller.pending && <p role="status">{t('name.pending')}</p>}
    {(value?.failed || status === 'blocked') && <p role="alert">{t('name.blocked')}</p>}
    {status === 'unknown' && <><p role="alert">{value?.outcome?.reason === 'source-review' ? mappingText('name_source_review') : t('name.unknown')}</p><button type="button" disabled={disabled || controller.pending} onClick={() => void controller.inspect()}>{t('name.inspect')}</button>{value?.outcome?.canDiscardRejected && <button type="button" disabled={disabled || controller.pending} onClick={() => void controller.discardRejected()}>{mappingText('cancel_rejected_name')}</button>}{value?.outcome?.currentName !== undefined && <p>{t('name.current', { name: value.outcome.currentName })}</p>}</>}
    <label>{t('name.label')}<input autoFocus aria-label={t('name.label')} value={value?.name ?? ''} disabled={disabled || status !== 'ready'} onChange={event => controller.change(event.target.value)} /></label>
  </div><footer className="gw-modal-footer"><button type="button" disabled={controller.pending} onClick={onBack}>{t(status === 'unknown' ? 'source.back' : 'picker.cancel')}</button>
    <button type="button" className="gw-primary" disabled={disabled || status !== 'ready' || !value?.name.trim() || !controller.dirty} onClick={onSave}>{t('name.save')}</button>
  </footer></>;
}
