import { useTranslation } from 'react-i18next';
import type { useRecordingPointNameBuffer } from './useRecordingPointNameBuffer';
export function RecordingPointNameContent({ controller, onBack, onSave, disabled }: { controller: ReturnType<typeof useRecordingPointNameBuffer>; onBack(): void; onSave(): void; disabled: boolean }) {
  const { t } = useTranslation('recording-editor'), value = controller.value, status = value?.outcome?.status;
  return <><div className="gw-modal-body">
    {controller.pending && <p role="status">{t('name.pending')}</p>}
    {(value?.failed || status === 'blocked') && <p role="alert">{t('name.blocked')}</p>}
    {status === 'unknown' && <><p role="alert">{t('name.unknown')}</p><button type="button" onClick={() => void controller.inspect()}>{t('name.inspect')}</button>{value?.outcome?.currentName !== undefined && <p>{t('name.current', { name: value.outcome.currentName })}</p>}</>}
    <label>{t('name.label')}<input autoFocus aria-label={t('name.label')} value={value?.name ?? ''} disabled={disabled || status !== 'ready'} onChange={event => controller.change(event.target.value)} /></label>
  </div><footer className="gw-modal-footer"><button type="button" disabled={controller.pending} onClick={onBack}>{t(status === 'unknown' ? 'source.back' : 'picker.cancel')}</button>
    <button type="button" className="gw-primary" disabled={disabled || status !== 'ready' || !value?.name.trim() || !controller.dirty} onClick={onSave}>{t('name.save')}</button>
  </footer></>;
}
