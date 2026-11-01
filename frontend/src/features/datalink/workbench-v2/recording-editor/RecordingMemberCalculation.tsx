import { useTranslation } from 'react-i18next';
import type { WriteGroupMemberCalculation } from '../../../../types/studioV2WriteGroup';
import type { StudioV2LocatorReadScope } from '../../../../services/studioV2WorkspaceDatabase';
import { useStudioV2WorkspaceMeasurementsQuery } from '../../../../hooks/datalink/useStudioV2WorkspaceMeasurements';
import { canConfirmMemberGauge, confirmMemberGauge, changeMemberCalculation, memberCalculationReason, type CalculationEditorMember, type MemberCalculationContext } from './recordingMemberCalculation';
export interface RecordingMemberCalculationProps { member: CalculationEditorMember; context: MemberCalculationContext; disabled: boolean; onChange: (member: CalculationEditorMember) => void }
const FIELD = 'mt-1 block w-full rounded border border-slate-700 bg-slate-950 px-2 py-1 text-xs text-slate-100';
const methods: WriteGroupMemberCalculation[] = ['latest', 'mean', 'delta'];
export function RecordingMemberCalculation({ member, context, disabled, onChange }: RecordingMemberCalculationProps) {
  const { t } = useTranslation('recording-editor');
  const cadence = member.cadence;
  const available = Boolean(context.capabilities?.editor_supported && context.capabilities.draft_supported);
  const blocked = disabled || !available;
  const reasons = methods.map(method => memberCalculationReason(member, method, context));
  const safeReason = (reason: string) => t(`calculation.${['capability', 'dialect', 'layout', 'gauge_required', 'counter_required'].includes(reason) ? reason : 'unavailable'}`);
  return <fieldset disabled={blocked} className="mt-2 space-y-2" data-testid={`member-calculation-${member.point_id}`}>
    <legend className="text-xs font-medium">{t('calculation.title')}</legend>
    <div className="grid gap-2 sm:grid-cols-3">
      <label>{t('calculation.title')}<select className={FIELD} value={member.calculation ?? ''} onChange={event => {
        const method = event.target.value as WriteGroupMemberCalculation;
        if (!blocked && methods.includes(method)) onChange(changeMemberCalculation(member, method, context));
      }}><option value="" disabled>{t('calculation.unavailable')}</option>{methods.map((method, index) => <option key={method} value={method} disabled={Boolean(reasons[index])}>{t(`calculation.${method}`)}</option>)}</select></label>
      <label>{t('calculation.interval')}<input className={FIELD} type="text" inputMode="numeric" value={member.cadence_input ?? (Number.isFinite(cadence?.interval) ? String(cadence?.interval) : '')}
        onChange={event => { if (!blocked && cadence) onChange({ ...member, cadence_input: event.target.value, cadence: { ...cadence, interval: event.target.value.trim() ? Number(event.target.value) : NaN } }); }} /></label>
      <label>{t('calculation.unit')}<select className={FIELD} value={cadence?.unit ?? ''} onChange={event => {
        const unit = event.target.value;
        if (!blocked && cadence && (unit === 'second' || unit === 'minute') && context.capabilities?.cadence.units.includes(unit)) onChange({ ...member, cadence: { ...cadence, unit } });
      }}><option value="" disabled>{t('calculation.unavailable')}</option>{context.capabilities?.cadence.units.map(unit => <option key={unit} value={unit}>{t(`calculation.${unit}`)}</option>)}</select></label>
    </div>
    {canConfirmMemberGauge(member, context) && reasons[1] && <div className="space-y-1 text-xs text-slate-400">
      <p>{t('calculation.gaugeConfirmationMeaning')}</p>
      <button type="button" className="gw-secondary" disabled={blocked} onClick={() => { if (!blocked) onChange(confirmMemberGauge(member, context)); }}>{t('calculation.confirmGauge')}</button>
    </div>}
    {member.gauge_confirmation && !reasons[1] && <p className="text-xs text-slate-400">{t('calculation.gaugeConfirmed')}</p>}
    {reasons.map((reason, index) => reason && <p key={methods[index]} className="text-xs text-slate-400">{t(`calculation.${methods[index]}`)}: {safeReason(reason)}</p>)}
    {context.capabilities && <p className="text-xs text-slate-400">{t('calculation.bounds', { min: context.capabilities.cadence.minimum, seconds: context.capabilities.cadence.max_seconds, minutes: context.capabilities.cadence.max_minutes })}</p>}
    <p className="text-xs text-slate-400">{t('calculation.dueUnion')}</p>
  </fieldset>;
}
/** Canonical read-only query shared by React Query; this row owns no save queue or measurement mutation. */
export function RecordingMemberCalculationSource({ workspaceRevision, locatorScope, ...props }: Omit<RecordingMemberCalculationProps, 'context'> & { context: Omit<MemberCalculationContext, 'definitions' | 'definitionsAvailable'>; workspaceRevision: string; locatorScope?: StudioV2LocatorReadScope }) {
  const definitions = useStudioV2WorkspaceMeasurementsQuery(true, undefined, { workspaceId: props.context.workspaceId, workspaceRevision }, locatorScope);
  return <RecordingMemberCalculation {...props} context={{ ...props.context, definitionsAvailable: definitions.isSuccess && !definitions.isFetching && !definitions.isError, definitions: definitions.isError || definitions.isFetching ? [] : definitions.data ?? [] }} />;
}
