import { RuntimeSettingsLink } from './RuntimeSettingsLink';
import { useTranslation } from 'react-i18next';
import type { EvidenceScope, RecordingGroupFacts, RecordingGroupView, RecordingMember } from '../../../types/runtimeRuleEvidence';
import { RuntimeSetupFact } from './RuntimeSetupHelpers';
import { RuntimeRecordingDeliveryCard } from './RuntimeRecordingDeliveryCard';
import { RuntimeRecordingDeliveryDetails } from './RuntimeRecordingDeliveryDetails';
import { RuntimeRecordingTimingCard } from './RuntimeRecordingTimingCard';

/** Revision-bound configuration and independently verified delivery; neither grants permission. */
export function RuntimeRecordingGroupCard({ group, deviceId, scope, navigateTo }: { group: RecordingGroupView; deviceId: string; scope?: EvidenceScope; navigateTo?: (target: string) => void }) {
  const { t } = useTranslation('runtime-dashboard');
  const unavailable = t('recordingRules.unavailable');
  const effective = group.state === 'effective' && group.proof === 'applied-snapshot-and-layout'
    && group.effective_revision && group.effective?.binding.entity_revision === group.effective_revision ? group.effective : null;
  const locator = (facts: RecordingGroupFacts | null, revision: string | null) => scope && scope.device_id === deviceId && revision &&
    facts?.binding.workspace_id === scope.workspace_id && facts.binding.workspace_revision === scope.workspace_revision && facts.binding.entity_revision === revision &&
    facts.members.some((member) => member.device_id === deviceId)
    ? { step: 4 as const, workspace_id: scope.workspace_id, device_id: deviceId, rule_id: null, mapping_id: null, group_id: group.group_id, group_revision: revision } : null;
  const savedLocator = locator(group.saved, group.saved_revision);
  const effectiveLocator = locator(effective, group.effective_revision);
  const shared = [group.saved, effective].some((facts) => facts?.members.some((member) => member.device_id !== deviceId));
  return <article aria-label={t('recordingRules.group', { id: group.group_id })} className="rounded-xl border border-slate-700 bg-slate-950/50 p-4">
    <h3 className="font-semibold text-slate-100">{t('recordingRules.group', { id: group.group_id })}</h3>
    {(savedLocator || effectiveLocator) && <div className="mt-2 space-y-2">
      <div className="flex flex-wrap gap-2">
        <RuntimeSettingsLink locator={savedLocator} navigateTo={navigateTo}>{t('settingsLink.savedGroup')}</RuntimeSettingsLink>
        <RuntimeSettingsLink locator={effectiveLocator} navigateTo={navigateTo}>{t('settingsLink.effectiveGroup')}</RuntimeSettingsLink>
      </div>
      <p className="text-xs text-slate-400">{t('settingsLink.revisionHint')}</p>
      {shared && <p className="text-xs text-slate-400">{t('settingsLink.sharedGroup')}</p>}
    </div>}
    <dl className="mt-3 grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
      <RuntimeSetupFact label={t('recordingRules.savedRevision')} value={group.saved_revision ?? unavailable} />
      <RuntimeSetupFact label={t('recordingRules.effectiveRevision')} value={effective?.binding.entity_revision ?? unavailable} />
      <RuntimeSetupFact label={t('recordingRules.asOf')} value={group.as_of ?? unavailable} />
    </dl>
    <p className="mt-3 text-sm text-slate-400">{t('recordingRules.notIntake')}</p>
    <section aria-label={t('recordingRules.requestedTitle')} className="mt-4 rounded border border-slate-800 p-3">
      <h4 className="font-medium text-slate-200">{t('recordingRules.requestedTitle')}</h4>
      <dl className="mt-2 grid gap-3 sm:grid-cols-2">
        <RuntimeSetupFact label={t('recordingRules.requestedRevision')} value={group.applied_revision ?? unavailable} />
        <RuntimeSetupFact label={t('recordingRules.requestedBoundary')} value={group.applied_effective_at ?? unavailable} />
      </dl>
      <p className="mt-2 text-sm text-slate-400">{t('recordingRules.requestedNote')}</p>
    </section>
    <div className="mt-4 grid gap-4 xl:grid-cols-2">
      <section aria-label={t('recordingRules.savedTitle')} className="rounded border border-amber-900/60 p-3">
        <h4 className="font-medium text-amber-100">{t('recordingRules.savedTitle')}</h4>
        <p className="mt-1 text-sm text-slate-400">{t('recordingRules.savedNote')}</p>
        <GroupFacts facts={group.saved} deviceId={deviceId} />
      </section>
      <section aria-label={t('recordingRules.effectiveTitle')} className="rounded border border-slate-700 p-3">
        <h4 className="font-medium text-slate-100">{t('recordingRules.effectiveTitle')}</h4>
        {effective ? <GroupFacts facts={effective} deviceId={deviceId} /> : <p className="mt-2 text-slate-400">{t('recordingRules.effectiveUnavailable')}</p>}
      </section>
    </div>
    <RuntimeRecordingDeliveryCard group={group} />
    <RuntimeRecordingDeliveryDetails group={group} deviceId={deviceId} />
    <RuntimeRecordingTimingCard group={group} />
  </article>;
}

function GroupFacts({ facts, deviceId }: { facts: RecordingGroupFacts | null; deviceId: string }) {
  const { t } = useTranslation('runtime-dashboard');
  const unavailable = t('recordingRules.unavailable');
  if (!facts) return <p className="mt-2 text-slate-400">{unavailable}</p>;
  const show = (value: string | null) => value === null || value === 'unavailable' ? unavailable : value;
  const seconds = (value: number | null) => value === null ? unavailable : t('recordingRules.seconds', { count: value });
  const profile = facts.recording_profile;
  const selected = facts.members.filter((member) => member.device_id === deviceId);
  const others = facts.members.filter((member) => member.device_id !== deviceId);
  const fields: [string, string][] = [
    ['factRevision', facts.binding.entity_revision],
    ['connectorId', facts.destination.connector_id], ['connectorRevision', facts.destination.connector_revision],
    ['alias', show(facts.destination.alias)], ['kind', show(facts.destination.kind)],
    ['schema', show(facts.destination.schema)], ['table', facts.destination.table],
    ['layout', show(facts.layout)], ['storage', show(facts.storage_mode)],
    ['profile', profile ? `${profile.kind} · ${t(profile.unit === 'second' ? 'recordingRules.seconds' : 'recordingRules.minutes', { count: profile.interval })}` : unavailable],
    ['interval', seconds(facts.interval_seconds)], ['reduction', show(facts.reduction)],
    ['writeMode', show(facts.write_mode)], ['incompletePolicy', show(facts.incomplete_policy)],
    ['lateness', seconds(facts.allowed_lateness_seconds)],
    ['enabled', t(facts.enabled ? 'recordingRules.enabledYes' : 'recordingRules.enabledNo')],
    ['effectiveAt', facts.effective_at ?? unavailable], ['endsAt', facts.ends_at ?? unavailable],
  ];
  return <>
    <dl className="mt-3 grid gap-3 sm:grid-cols-2">{fields.map(([key, value]) => <RuntimeSetupFact key={key} label={t(`recordingRules.${key}`)} value={value} />)}</dl>
    <p className="mt-4 text-sm text-slate-400">{t('recordingRules.groupWide')}</p>
    <section aria-label={t('recordingRules.selectedMembers', { device: deviceId })} className="mt-2">
      <h5 className="text-sm font-medium text-slate-200">{t('recordingRules.selectedMembers', { device: deviceId })}</h5>
      {selected.length ? <Members members={selected} /> : <p className="mt-1 text-sm text-slate-400">{t('recordingRules.noSelectedMembers')}</p>}
    </section>
    {others.length > 0 && <details className="mt-3 rounded border border-slate-800 p-2">
      <summary className="cursor-pointer text-sm text-slate-200 focus-visible:outline focus-visible:outline-sky-400">{t('recordingRules.otherMembers', { count: others.length })}</summary>
      <Members members={others} />
    </details>}
  </>;
}

function Members({ members }: { members: RecordingMember[] }) {
  const { t } = useTranslation('runtime-dashboard');
  const unavailable = t('recordingRules.unavailable');
  return <ul className="mt-2 space-y-2">{members.map((member) => <li key={`${member.device_id}:${member.point_id}:${member.tag_id}:${member.column}`} className="rounded bg-slate-900/70 p-2">
    <dl className="grid gap-2 sm:grid-cols-2">
      <RuntimeSetupFact label={t('recordingRules.device')} value={member.device_id} />
      <RuntimeSetupFact label={t('recordingRules.point')} value={member.point_id} />
      <RuntimeSetupFact label={t('recordingRules.tag')} value={member.tag_id} />
      <RuntimeSetupFact label={t('recordingRules.column')} value={member.column} />
      <RuntimeSetupFact label={t('recordingRules.sourceType')} value={member.source_value_type ?? unavailable} />
      <RuntimeSetupFact label={t('recordingRules.circuit')} value={member.circuit_id ?? unavailable} />
    </dl>
  </li>)}</ul>;
}
