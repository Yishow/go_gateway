import { useTranslation } from 'react-i18next';
import { Archive, CircleDot, CircleHelp, CirclePause, PencilLine, ToggleLeft, ToggleRight, type LucideIcon } from 'lucide-react';
import type { WriteGroupStatus } from '../../../../types/studioV2WriteGroup';

/** Saved configuration status only; the icon never claims runtime health or SQL success. Saved
 * running is a configuration switch, not live running. */
const statusIcons: Record<WriteGroupStatus, LucideIcon> = { draft: PencilLine, ready: CircleDot, running: ToggleRight, disabled: CirclePause, deleted: Archive };
type KnownAvailability = 'available' | 'unavailable' | 'suspended';
const availabilityIcons: Record<KnownAvailability, LucideIcon> = { available: ToggleRight, unavailable: ToggleLeft, suspended: CirclePause };
/** Own keys only: untrusted strings such as `__proto__` or `constructor` stay unknown. */
const knownAvailability = (value: string): value is KnownAvailability => Object.hasOwn(availabilityIcons, value);

export function RecordingStatusBadge({ status }: { status: WriteGroupStatus }) {
  const { t } = useTranslation('recording-home');
  const Icon = Object.hasOwn(statusIcons, status) ? statusIcons[status] : CircleHelp;
  return <span className="gw-home-badge" data-status-badge={status}><Icon aria-hidden="true" size={14} />{t(`savedStatuses.${status}`)}</span>;
}
/** Persisted Studio setting availability (missing data defaults to available upstream), not a live
 * connection: labelled as a setting, switch icons only, no health colour. Anything else stays unknown. */
export function DeviceAvailabilityBadge({ availability }: { availability: string }) {
  const { t } = useTranslation('recording-home');
  const state = knownAvailability(availability) ? availability : 'unknown';
  const Icon = state === 'unknown' ? CircleHelp : availabilityIcons[state];
  return <span className="gw-home-badge" data-availability-badge={state}><Icon aria-hidden="true" size={14} />
    {t(`configAvailability.${state}`)}</span>;
}
