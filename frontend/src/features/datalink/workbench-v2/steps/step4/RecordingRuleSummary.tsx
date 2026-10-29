import type { ReactNode } from 'react';
import { safeRuntimeLabel } from '../../../../../utils/runtimeOperatorSafety';
import * as Dialog from '@radix-ui/react-dialog';
import type { EvidenceDataType, RecordingGroupFacts } from '../../../../../types/runtimeRuleEvidence';
import type { RecordingReadonly, RecordingUnavailableReason, RecordingWorkflowScope, RecordingWorkflowView } from '../../state/recordingWorkflowView';

type SummaryField = 'sources' | 'layout' | 'storageMode' | 'profile' | 'reduction' | 'interval' | 'lateness' |
  'incomplete' | 'destination' | 'ownership' | 'workspace' | 'device' | 'circuit' | 'group' |
  'savedRevision' | 'requestedRevision' | 'effectiveRevision' | 'asOf' | 'effectiveAt' | 'endsAt';

export interface RecordingRuleSummaryLabels {
  readonly heading: string;
  readonly saved: string;
  readonly effective: string;
  readonly savedDescription: string;
  readonly effectiveDescription: string;
  readonly unavailable: string;
  readonly emptySources: string;
  readonly normalizedInterval: string;
  readonly groupScope: string;
  readonly requestedNotEffective: string;
  readonly details: string;
  readonly detailsTitle: string;
  readonly detailsDescription: string;
  readonly detailsUnavailable: string;
  readonly close: string;
  readonly fields: Readonly<Record<SummaryField, string>>;
  readonly profiles: Readonly<Record<'instantaneous' | 'cumulative', string>>;
  readonly units: Readonly<Record<'second' | 'minute', string>>;
  readonly layouts: Readonly<Record<RecordingGroupFacts['layout'], string>>;
  readonly storageModes: Readonly<Record<RecordingGroupFacts['storage_mode'], string>>;
  readonly reductions: Readonly<Record<RecordingGroupFacts['reduction'], string>>;
  readonly incompletePolicies: Readonly<Record<RecordingGroupFacts['incomplete_policy'], string>>;
  readonly sourceTypes: Readonly<Record<EvidenceDataType, string>>;
  readonly reasons: Readonly<Record<RecordingUnavailableReason, string>>;
}

/** Existing verified detail owner supplies read-only content, not raw JSON or mutation controls.
 * This slot is not a new detail/query authority. Mounting real bindings/policies/evidence stays gated.
 */
export interface RecordingRuleDetailSlot {
  readonly scope: RecordingWorkflowScope;
  readonly effectiveRevision: string | null;
  readonly asOf: string | null;
  readonly content: ReactNode;
}
export interface RecordingRuleSummaryProps {
  readonly view: RecordingWorkflowView;
  readonly scopeLabel: string;
  readonly labels: RecordingRuleSummaryLabels;
  readonly details?: RecordingRuleDetailSlot;
  /** Display-only names keyed by the complete source identity; never ownership evidence. */
  readonly sourceLabels?: Readonly<Record<string, string>>;
}

function Field({ label, children }: { label: string; children: ReactNode }) {
  return <div className="min-w-0"><dt className="text-xs text-slate-400">{label}</dt>
    <dd className="break-words text-sm text-slate-100">{children}</dd></div>;
}

function Facts({ facts, labels, sourceLabels }: { facts: RecordingReadonly<RecordingGroupFacts>; labels: RecordingRuleSummaryLabels;
  sourceLabels?: Readonly<Record<string, string>> }) {
  const profile = facts.recording_profile;
  const interval = profile?.interval ?? facts.interval_seconds;
  // Seconds is explicit in interval_seconds. It is not a guessed legacy profile
  // or a physical measurement unit; retain original profile units when proven.
  const period = interval === null ? labels.unavailable : `${interval} ${labels.units[profile?.unit ?? 'second']}`;
  const destination = facts.destination;
  return <dl className="grid min-w-0 gap-3">
    <Field label={labels.fields.sources}>{facts.members.length === 0 ? labels.emptySources :
      <ul aria-label={labels.fields.sources} className="max-h-48 space-y-1 overflow-y-auto">
        {facts.members.map((member, index) => <li key={JSON.stringify([member.device_id, member.point_id, member.tag_id])} className="break-words">
          <span>{safeRuntimeLabel(sourceLabels?.[JSON.stringify([member.device_id, member.point_id, member.tag_id])]) || `${labels.unavailable} ${index + 1}`}</span>
          {' — '}<span>{member.source_value_type === null ? labels.unavailable : labels.sourceTypes[member.source_value_type]}</span>
        </li>)}
      </ul>}</Field>
    <Field label={labels.fields.layout}>{labels.layouts[facts.layout]}</Field>
    <Field label={labels.fields.storageMode}>{labels.storageModes[facts.storage_mode]}</Field>
    <Field label={labels.fields.profile}>{profile === null ? labels.unavailable : labels.profiles[profile.kind]}</Field>
    <Field label={labels.fields.reduction}>{labels.reductions[facts.reduction]}</Field>
    <Field label={labels.fields.interval}><span>{period}</span>
      {profile === null && interval !== null && <small className="ml-2 text-slate-400">{labels.normalizedInterval}</small>}</Field>
    <Field label={labels.fields.lateness}>{facts.allowed_lateness_seconds === null ? labels.unavailable
      : `${facts.allowed_lateness_seconds} ${labels.units.second}`}</Field>
    <Field label={labels.fields.incomplete}>{labels.incompletePolicies[facts.incomplete_policy]}</Field>
    <Field label={labels.fields.destination}>{destination.alias !== null && <span>{destination.alias} — </span>}
      <span>{destination.schema === null ? destination.table : `${destination.schema}.${destination.table}`}</span></Field>
    <Field label={labels.fields.ownership}>{labels.groupScope}</Field>
  </dl>;
}

function VersionCard({ title, description, facts, reason, labels, sourceLabels }: {
  title: string; description: string; facts: RecordingReadonly<RecordingGroupFacts> | null;
  reason?: RecordingUnavailableReason; labels: RecordingRuleSummaryLabels; sourceLabels?: Readonly<Record<string, string>>;
}) {
  return <section aria-label={title} className="min-w-0 space-y-3 rounded border border-slate-700 p-3">
    <h4 className="font-medium text-slate-100">{title}</h4>
    <p className="text-xs text-slate-400">{description}</p>
    {facts === null ? <p role="status" className="text-sm text-slate-300">{labels.unavailable}{reason && ` — ${labels.reasons[reason]}`}</p>
      : <Facts facts={facts} labels={labels} sourceLabels={sourceLabels} />}
  </section>;
}

const scopeFields = ['workspaceId', 'deviceId', 'circuitId', 'groupId', 'workspaceRevision', 'groupRevision', 'appliedRevision'] as const;

/**
 * Read-only presentation. Both summaries consume only runtime's safe
 * facts, never the editable WriteGroup (which may include a private DB path).
 * A retained effective snapshot describes history, not intake or SQL success.
 * The only local state belongs to Radix's disclosure/focus management; no effects
 * query, save, apply, start, or persist selection.
 */
export function RecordingRuleSummary({ view, scopeLabel, labels, details, sourceLabels }: RecordingRuleSummaryProps) {
  if (view.state !== 'available') return <section aria-label={labels.heading} className="space-y-3">
    <h3 className="text-sm font-medium text-slate-100">{labels.heading}</h3>
    <p role="status" className="text-sm text-slate-300">{labels.reasons[view.reason]}</p>
  </section>;
  const runtime = view.runtime.state === 'available' ? view.runtime.value : null;
  const effective = view.effective.state === 'available' ? view.effective.value : null;
  const asOf = runtime?.as_of ?? null;
  const effectiveRevision = effective?.revision ?? null;
  const matchingDetails = details !== undefined && details.content !== null && details.content !== undefined &&
    details.content !== false && details.content !== '' && scopeFields.every(key => details.scope[key] === view.scope[key]) &&
    details.effectiveRevision === effectiveRevision && details.asOf === asOf;
  // Remount only the local disclosure when scope/revision changes, so an open
  // drawer never silently becomes another device's detail view or reopens on Back.
  const scopeKey = JSON.stringify(scopeFields.map(key => view.scope[key]));
  return <Dialog.Root key={scopeKey}>
    <section aria-label={labels.heading} className="min-w-0 space-y-3">
      <h3 className="text-sm font-medium text-slate-100">{labels.heading}</h3>
      <p className="break-words text-sm text-slate-300">{scopeLabel}</p>
      <div className="grid min-w-0 gap-3 sm:grid-cols-2">
        <VersionCard title={labels.saved} description={labels.savedDescription} facts={runtime?.saved ?? null}
          reason={view.runtime.state === 'unavailable' ? view.runtime.reason : undefined} labels={labels} sourceLabels={sourceLabels} />
        <VersionCard title={labels.effective} description={labels.effectiveDescription} facts={effective?.facts ?? null}
          reason={view.effective.state === 'unavailable' ? view.effective.reason : undefined} labels={labels} sourceLabels={sourceLabels} />
      </div>
      {effectiveRevision !== null && runtime?.applied_revision && runtime.applied_revision !== effectiveRevision &&
        <p role="status" className="text-sm text-slate-300">{labels.requestedNotEffective}</p>}
      <Dialog.Trigger type="button" className="rounded border border-slate-600 px-3 py-2 text-sm text-slate-100">{labels.details}</Dialog.Trigger>
    </section>
    <Dialog.Portal>
      <Dialog.Overlay className="fixed inset-0 z-50 bg-black/60" />
      <Dialog.Content className="fixed inset-y-0 right-0 z-50 w-full max-w-2xl space-y-4 overflow-y-auto border-l border-slate-700 bg-slate-950 p-4 text-slate-100">
        <Dialog.Title className="text-lg font-medium">{labels.detailsTitle}</Dialog.Title>
        <Dialog.Description className="text-sm text-slate-300">{labels.detailsDescription}</Dialog.Description>
        <Dialog.Close type="button" className="rounded border border-slate-600 px-3 py-2 text-sm">{labels.close}</Dialog.Close>
        <p className="break-words text-sm">{scopeLabel}</p>
        <dl className="grid min-w-0 gap-3 sm:grid-cols-2">
          <Field label={labels.fields.workspace}>{view.scope.workspaceId}</Field>
          <Field label={labels.fields.group}>{view.scope.groupId}</Field>
          <Field label={labels.fields.device}>{view.scope.deviceId}</Field>
          <Field label={labels.fields.circuit}>{view.scope.circuitId ?? labels.unavailable}</Field>
          <Field label={labels.fields.savedRevision}>{view.scope.groupRevision}</Field>
          <Field label={labels.fields.requestedRevision}>{runtime?.applied_revision ?? labels.unavailable}</Field>
          <Field label={labels.fields.effectiveRevision}>{effectiveRevision ?? labels.unavailable}</Field>
          <Field label={labels.fields.asOf}>{asOf ?? labels.unavailable}</Field>
          <Field label={labels.fields.effectiveAt}>{effective?.facts.effective_at ?? labels.unavailable}</Field>
          <Field label={labels.fields.endsAt}>{effective?.facts.ends_at ?? labels.unavailable}</Field>
        </dl>
        {matchingDetails ? details.content : <p role="status">{labels.detailsUnavailable}</p>}
      </Dialog.Content>
    </Dialog.Portal>
  </Dialog.Root>;
}
