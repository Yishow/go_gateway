import { useTranslation } from 'react-i18next';
import { RecordingHomeRowActions } from './RecordingHomeRowActions';
import './recording-delete-dialog.css';
import type { HomeDeviceSummary } from './recordingHomeCollections';

/** The same menu is used in the device strip and the devices table. */
export function RecordingHomeDeviceActions({ device, available, onDelete }: {
  device: HomeDeviceSummary; available: boolean; onDelete: (device: HomeDeviceSummary) => void;
}) {
  const { t } = useTranslation('recording-home');
  return <RecordingHomeRowActions label={t('deleteDialogs.device.actions', { name: device.name })}
    items={[{ id: 'delete', label: t('deleteDialogs.device.confirm'), disabled: !available, onSelect: () => onDelete(device) }]}
    hint={available ? undefined : t('deleteDialogs.device.stale')} />;
}
