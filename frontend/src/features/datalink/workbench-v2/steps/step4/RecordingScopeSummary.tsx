import type { ReactNode } from 'react';
import { recordingWorkflowView, type RecordingWorkflowScope, type RecordingWorkflowView } from '../../state/recordingWorkflowView';
import { buildGroupCandidates, type GroupCandidate } from '../../state/writeGroup/candidates';
import { safeRuntimeAddress, safeRuntimeLabel } from '../../../../../utils/runtimeOperatorSafety';
import { RecordingRuleSummary } from './RecordingRuleSummary';
import { RecordingRuleDetails } from './RecordingRuleDetails';
import { useRecordingGuidedLabels, type RecordingGuidedLabels } from './recordingGuidedLabels';
import type { RecordingScopeSummaryProps } from './recordingScopeSummary';

const sourceKey = (source: Pick<GroupCandidate, 'device_id' | 'point_id' | 'tag_id'>) =>
  JSON.stringify([source.device_id, source.point_id, source.tag_id]);

function Field({ label, children }: { label: string; children: ReactNode }) {
  return <div className="min-w-0"><dt className="text-xs text-slate-400">{label}</dt>
    <dd className="break-words text-sm text-slate-100">{children}</dd></div>;
}

function DraftSummary({ scope, state, labels, sourceLabels }: Pick<RecordingScopeSummaryProps, 'scope' | 'state'> & {
  labels: RecordingGuidedLabels; sourceLabels: Readonly<Record<string, string>>;
}) {
  const { draft } = scope;
  const alias = safeRuntimeLabel(state.db.connector.name) || labels.destinationFallback;
  return <section aria-label={labels.unsaved} className="min-w-0 space-y-3 rounded border border-slate-700 p-3">
    <h4 className="font-medium text-slate-100">{labels.unsaved}</h4>
    <p className="text-xs text-slate-400">{labels.unsavedDescription}</p>
    <dl className="grid min-w-0 gap-3 sm:grid-cols-2">
      <Field label={labels.fields.layout}>{labels.layouts[draft.layout === 'wide' ? 'managed_wide_profile_v1' : 'managed_long_v1']}</Field>
      <Field label={labels.fields.storageMode}>{labels.storageModes[draft.storageMode]}</Field>
      <Field label={labels.fields.incomplete}>{labels.incompletePolicies[draft.incompletePolicy]}</Field>
      <Field label={labels.fields.destination}>{alias} — {labels.tablePending}</Field>
    </dl>
    {draft.profiles.length === 0 && <p role="status" className="text-sm text-slate-300">{labels.noProfiles}</p>}
    {draft.profiles.map(profile => <section key={profile.kind} aria-label={labels.profiles[profile.kind]} className="space-y-2">
      <h5 className="text-sm text-slate-200">{labels.profiles[profile.kind]}</h5>
      <dl className="grid min-w-0 gap-3 sm:grid-cols-2">
        <Field label={labels.fields.interval}>{/^[1-9][0-9]*$/.test(profile.interval) && Number.isSafeInteger(Number(profile.interval))
          ? `${profile.interval} ${labels.units[profile.unit]}` : labels.unavailable}</Field>
        <Field label={labels.fields.reduction}>{labels.reductions[profile.reduction]}</Field>
        <Field label={labels.fields.sources}>{profile.sources.length === 0 ? labels.emptySources :
          <ul className="max-h-48 space-y-1 overflow-y-auto" aria-label={labels.fields.sources}>
            {profile.sources.map((source, index) => <li key={sourceKey(source)} className="break-words">
              <span>{sourceLabels[sourceKey(source)] || `${labels.unavailableSource} ${index + 1}`}</span> — <span>{labels.sourceTypes[source.target_type] ?? labels.unavailable}</span>
            </li>)}
          </ul>}</Field>
      </dl>
    </section>)}
  </section>;
}

function SavedSummary({ view, scopeLabel, labels, sourceLabels }: {
  view: RecordingWorkflowView; scopeLabel: string; labels: RecordingGuidedLabels; sourceLabels: Readonly<Record<string, string>>;
}) {
  const details = view.state === 'available' ? {
    scope: view.scope,
    effectiveRevision: view.effective.state === 'available' ? view.effective.value.revision : null,
    asOf: view.runtime.state === 'available' ? view.runtime.value.as_of : null,
    content: <RecordingRuleDetails view={view} labels={labels} />,
  } : undefined;
  return <RecordingRuleSummary view={view} scopeLabel={scopeLabel} labels={labels} details={details} sourceLabels={sourceLabels} />;
}

/**
 * Pure view adapter. The mounted Step 4 owner provides ONE device/revision query;
 * groups and drawers never query independently. Native draft choices are display
 * only; saved and effective cards require runtime's safe version-bound facts.
 */
export function RecordingScopeSummary({ scope, state, groups, configurationState, evidence }: RecordingScopeSummaryProps) {
  const labels = useRecordingGuidedLabels();
  const device = state.devices.find(candidate => candidate.id === scope.deviceId);
  const deviceLabel = safeRuntimeLabel(device?.name) || labels.selectedDevice;
  const circuitLabel = scope.circuitId ? (safeRuntimeLabel(scope.draft.circuitName) || labels.selectedCircuit) : '';
  const scopeLabel = circuitLabel ? `${deviceLabel} / ${circuitLabel}` : deviceLabel;
  const candidates = buildGroupCandidates(state).candidates;
  const sourceLabels = Object.fromEntries(candidates.map(source => [sourceKey(source),
    safeRuntimeLabel(source.label) || safeRuntimeAddress(source.address) || labels.unavailableSource]));
  const unavailable = (message: string) => <p role="status" className="text-sm text-slate-300">{message}</p>;
  let content: ReactNode;
  const validScope = Boolean(scope.workspaceId && scope.workspaceRevision && scope.deviceId && device && scope.circuitId !== '');
  if (!validScope) content = unavailable(labels.noScope);
  else if (configurationState !== 'ready') content = unavailable(labels.reasons[configurationState]);
  else if (new Set(scope.groupIds).size !== scope.groupIds.length) content = unavailable(labels.reasons.scope_mismatch);
  else if (scope.groupIds.length === 0) {
    const candidateKeys = new Set(candidates.filter(source => source.device_id === scope.deviceId).map(sourceKey));
    const validDraftSources = scope.draft.profiles.every(profile => profile.sources.every(source =>
      source.device_id === scope.deviceId && candidateKeys.has(sourceKey(source))));
    content = validDraftSources
      ? <DraftSummary scope={scope} state={state} labels={labels} sourceLabels={sourceLabels} />
      : unavailable(labels.reasons.scope_mismatch);
  } else {
    const selectedGroups = scope.groupIds.map(id => groups.filter(group => group.id === id));
    const validGroups = selectedGroups.every(matches => matches.length === 1 && matches[0].workspace_id === scope.workspaceId && matches[0].status !== 'deleted');
    if (!validGroups) content = unavailable(labels.reasons.scope_mismatch);
    else content = selectedGroups.map(([group]) => {
      const selected: RecordingWorkflowScope = { workspaceId: scope.workspaceId, workspaceRevision: scope.workspaceRevision,
        deviceId: scope.deviceId, circuitId: scope.circuitId, groupId: group.id, groupRevision: group.revision, appliedRevision: group.applied_revision };
      const view = recordingWorkflowView({ scope: selected, savedGroup: group,
        runtimeEvidence: evidence.isError ? { state: 'error' } : evidence.isLoading || evidence.isFetching ? { state: 'loading' }
          : evidence.data?.ok ? { state: 'available', scope: selected, value: evidence.data.value } : { state: 'unavailable' } });
      return <SavedSummary key={group.id} view={view} scopeLabel={scopeLabel} labels={labels} sourceLabels={sourceLabels} />;
    });
  }
  return <section data-testid="recording-scope-summary" aria-label={labels.draftHeading} className="min-w-0 space-y-4 rounded border border-slate-700 bg-slate-950/30 p-4">
    <h3 className="text-sm font-semibold text-slate-100">{labels.draftHeading}</h3>
    <p className="break-words text-sm text-slate-300">{scopeLabel}</p>
    {content}
    {validScope && configurationState === 'ready' && scope.groupIds.length > 0 && <button type="button"
      disabled={evidence.isFetching} onClick={() => { void evidence.refetch(); }}
      className="rounded border border-slate-600 px-3 py-2 text-sm text-slate-100 disabled:opacity-50">{labels.retryEvidence}</button>}
  </section>;
}
