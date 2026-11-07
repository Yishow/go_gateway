import { useTranslation } from 'react-i18next';
import { useId } from 'react';
import { Trash2 } from 'lucide-react';
import './recording-delete-dialog.css';
import type { HomeDeviceSummary } from './recordingHomeCollections';

/** One-click delete entry for device cards; the dialog still confirms. */
export function RecordingHomeDeviceActions({ device, available, onDelete }: {
  device: HomeDeviceSummary; available: boolean; onDelete: (device: HomeDeviceSummary) => void;
}) {
  const { t } = useTranslation('recording-home');
  const hintId = useId();
  return <><button type="button" className="gw-row-icon-button gw-device-delete" disabled={!available}
    aria-label={t('deleteDialogs.device.deleteNamed', { name: device.name })}
    aria-describedby={available ? undefined : hintId}
    title={available ? t('deleteDialogs.device.confirm') : t('deleteDialogs.device.stale')}
    onClick={event => { event.currentTarget.focus(); onDelete(device); }}>
    <Trash2 aria-hidden="true" size={16} />
  </button>{!available && <p id={hintId} className="gw-row-hint">{t('deleteDialogs.device.stale')}</p>}</>;
}
