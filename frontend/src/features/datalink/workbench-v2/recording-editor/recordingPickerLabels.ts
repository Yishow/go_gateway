import type { TFunction } from 'i18next';
import type { PickerLabels } from './RecordingPointPicker';
export function recordingPickerLabels(t: TFunction<'recording-editor'>): PickerLabels {
  return { title: t('picker.title'), search: t('picker.search'), device: t('picker.device'), pointType: t('picker.pointType'),
    allDevices: t('picker.allDevices'), allTypes: t('picker.allTypes'), types: { gauge: t('picker.gauge'), counter: t('picker.counter'), state: t('picker.state') },
    unknownType: t('picker.unknownType'), point: t('picker.point'), address: t('picker.address'), alreadyAdded: t('picker.alreadyAdded'),
    unavailable: t('picker.unavailable'), cancel: t('picker.cancel'), loading: t('picker.loading'), failed: t('picker.failed'), retry: t('picker.retry'),
    empty: t('picker.empty'), noMatches: t('picker.noMatches'), clearFilters: t('picker.clearFilters'), newSource: t('picker.newSource'),
    invalidPending: t('picker.invalidPending'), previous: t('picker.previous'), next: t('picker.next'), pageSize: t('picker.pageSize'),
    matching: (count, added) => t('picker.matching', { count, added }), pending: count => t('picker.pending', { count }),
    hidden: count => t('picker.hidden', { count }), add: count => t('picker.add', { count }), page: (current, total) => t('picker.page', { current, total }),
    removePending: identity => t('picker.removePending', { point: identity.point_id }),
  };
}
