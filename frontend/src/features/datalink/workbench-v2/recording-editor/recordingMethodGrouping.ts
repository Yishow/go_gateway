import type { TFunction } from 'i18next';
import type { EditorDraft } from '../state/writeGroup/draft';
import type { CalculationEditorMember } from './recordingMemberCalculation';

/** One row of a per-point recording's method summary: members sharing a calculation and interval. */
export interface RecordingMethodGroup { key: string; method: string; interval: string; unit: string; count: number; firstPointId: string }

/** Groups members by the calculation and interval shown to the operator, in member order.
 * The summary rail and the method section render the same groups, so they cannot disagree. */
export function recordingMethodGroups(draft: Pick<EditorDraft, 'members'>, t: TFunction): RecordingMethodGroup[] {
  const groups = new Map<string, RecordingMethodGroup>();
  const members: readonly CalculationEditorMember[] = draft.members;
  for (const member of members) {
    const method = member.calculation ? t(`calculation.${member.calculation}`) : t('simple.unset');
    const interval = member.cadence_input ?? (Number.isFinite(member.cadence?.interval) ? String(member.cadence?.interval) : t('simple.unset'));
    const unit = member.cadence?.unit ? t(`calculation.${member.cadence.unit}`) : '';
    const key = JSON.stringify([method, interval, unit]);
    const current = groups.get(key);
    groups.set(key, current ? { ...current, count: current.count + 1 } : { key, method, interval, unit, count: 1, firstPointId: member.point_id });
  }
  return [...groups.values()];
}
