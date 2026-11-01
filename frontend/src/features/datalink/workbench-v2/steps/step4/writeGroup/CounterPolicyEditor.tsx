import { ParameterHelp } from '../../../components/ParameterHelp';
import type { StudioV2LocatorReadScope } from '../../../../../../services/studioV2WorkspaceDatabase';
import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { useStudioV2WorkspaceMeasurementsQuery, useUpdateStudioV2MeasurementMutation } from '../../../../../../hooks/datalink/useStudioV2WorkspaceMeasurements';
import type { MeasurementDefinition } from '../../../../../../types/measurement';
import { getSafeErrorMessage } from '../../../../../../utils/typedErrors';
import { focusCounterConfirmation } from '../step4FocusTargets';
import {
  blockingIssues, buildCounterPolicyDefinition, counterPolicyChanged, formFromDefinition, isFloatEncoding, resolvePolicyDefinition,
  validateCounterPolicy, type CounterPolicyForm, type PolicyMember,
} from '../../../state/counterPolicyForm';

export type CounterPolicySource = PolicyMember;

export interface CounterPolicyEditorProps {
  locatorScope?: StudioV2LocatorReadScope;
  workspaceId: string;
  members: CounterPolicySource[];
  /** Display names by tag ID. */
  labels?: Record<string, string>;
  readonly: boolean;
}

type Outcome = 'saved' | 'epoch_started';
const FIELD = 'mt-1 w-full rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] px-2 py-1 text-xs text-[var(--gw-text)] disabled:opacity-50';

function SourcePolicy({ def, confirmed, label, readonly, outcome, onSaved }: {
  def: MeasurementDefinition; confirmed: boolean; label: string; readonly: boolean; outcome?: Outcome; onSaved: (outcome: Outcome) => void;
}) {
  const { t } = useTranslation('workbench-v2');
  const update = useUpdateStudioV2MeasurementMutation();
  const [form, setForm] = React.useState<CounterPolicyForm>(() => formFromDefinition(def));
  const [error, setError] = React.useState<string>();
  const inFlight = React.useRef(false);
  const tag = def.tag_id ?? '';
  const float = isFloatEncoding(def);
  const saved = formFromDefinition(def);
  // A float counter's wrap fields stay locked unless an invalid saved wrap policy must be cleared.
  const wrapLocked = float && saved.modulus === '' && saved.rate === '';
  const issues = validateCounterPolicy(def, form);
  const blocked = blockingIssues(issues, form).length > 0;
  // The fields stay locked for the whole request: later typing would be erased when the refreshed revision remounts the form.
  const locked = readonly || update.isPending;
  const canSave = !locked && !blocked && (counterPolicyChanged(def, form) || !confirmed);
  const describedBy = issues.length > 0 ? `counter-policy-issues-${tag}` : undefined;

  const save = async () => {
    if (!canSave || inFlight.current) return;
    inFlight.current = true;
    setError(undefined);
    try {
      const result = await update.mutateAsync({ id: def.id, def: buildCounterPolicyDefinition(def, form) });
      onSaved(result.epoch_transitioned ? 'epoch_started' : 'saved');
    } catch (cause) {
      setError(getSafeErrorMessage(cause, t).message);
    } finally {
      inFlight.current = false;
    }
  };

  return (
    <fieldset data-testid={`counter-policy-source-${tag}`} className="space-y-2 rounded border border-[var(--gw-control)] p-3">
      <legend className="px-1 text-xs font-medium text-[var(--gw-text)]">{label}</legend>
      <p data-testid={`counter-policy-revision-${tag}`} className="text-xs text-[var(--gw-muted)]">
        {t('step4.group.policy.unit')}: {def.unit} · {t('step4.group.policy.definition_revision')}: {def.definition_revision} · {t('step4.group.policy.series_epoch')}: {def.series_epoch}
      </p>
      <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 flex items-center gap-2 text-xs text-[var(--gw-text)]">
        <input type="checkbox" data-testid={`counter-policy-signed-${tag}`} checked={form.signed} disabled={locked}
          onChange={(event) => setForm((current) => ({ ...current, signed: event.target.checked }))} />
        <span>{t('step4.group.policy.signed')}</span>
      </label><ParameterHelp help="counterPolicy.signed" /></div>
      <div className="grid gap-2 sm:grid-cols-2">
        <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 block text-xs text-[var(--gw-muted)]">
          <span>{t('step4.group.policy.modulus')}</span>
          <input type="text" inputMode="decimal" data-testid={`counter-policy-modulus-${tag}`} value={form.modulus} disabled={locked || wrapLocked}
            aria-invalid={issues.includes('modulus_invalid')} aria-describedby={describedBy} className={FIELD}
            onChange={(event) => setForm((current) => ({ ...current, modulus: event.target.value }))} />
        </label><ParameterHelp help="counterPolicy.modulus" /></div>
        <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 block text-xs text-[var(--gw-muted)]">
          <span>{t('step4.group.policy.rate')}</span>
          <input type="text" inputMode="decimal" data-testid={`counter-policy-rate-${tag}`} value={form.rate} disabled={locked || wrapLocked}
            aria-invalid={issues.includes('rate_invalid') || issues.includes('rate_needs_modulus')} aria-describedby={describedBy} className={FIELD}
            onChange={(event) => setForm((current) => ({ ...current, rate: event.target.value }))} />
        </label><ParameterHelp help="counterPolicy.rate" /></div>
      </div>
      {!confirmed && <p role="status" data-testid={`counter-policy-unconfirmed-${tag}`} className="text-xs text-[var(--gw-warning-fg)]">{t('step4.group.policy.unconfirmed')}</p>}
      {issues.length > 0 && (
        <ul id={describedBy} data-testid={`counter-policy-issues-${tag}`} className="list-disc pl-5 text-xs text-[var(--gw-warning-fg)]">
          {issues.map((issue) => <li key={issue}>{t(`step4.group.policy.issue.${issue}`)}</li>)}
        </ul>
      )}
      {error && <p role="alert" className="text-xs text-[var(--gw-danger-fg)]">{error}</p>}
      {outcome && <p role="status" data-testid={`counter-policy-result-${tag}`} className="text-xs text-[var(--gw-success-fg)]">{t(`step4.group.policy.${outcome}`)}</p>}
      <button type="button" data-testid={`counter-policy-save-${tag}`} disabled={!canSave} onClick={() => void save()}
        className="rounded bg-[var(--gw-active)] px-3 py-1 text-xs font-semibold text-[var(--gw-text)] disabled:opacity-40">
        {update.isPending ? t('step4.group.policy.saving') : t('step4.group.policy.save')}
      </button>
    </fieldset>
  );
}

/**
 * Edits the confirmed signed and rollover policy of each cumulative source through the existing measurement
 * definition API. It never creates a definition, never guesses a unit or equipment, and sends the full saved
 * definition with only the reviewed policy changed; the server stays the authority and its revision is shown back.
 */
export const CounterPolicyEditor: React.FC<CounterPolicyEditorProps> = ({ workspaceId, members, labels = {}, readonly, locatorScope }) => {
  const { t } = useTranslation('workbench-v2');
  const query = useStudioV2WorkspaceMeasurementsQuery(Boolean(workspaceId), undefined, undefined, locatorScope);
  // Results are kept per workspace and definition, so a late answer can never be shown for another scope.
  const [outcomes, setOutcomes] = React.useState<Record<string, Outcome>>({});
  const [notOnScreen, setNotOnScreen] = React.useState<Record<string, boolean>>({});
  const definitions = Array.isArray(query.data) ? query.data : [];
  return (
    <section data-testid="counter-policy" aria-label={t('step4.group.policy.title')} className="space-y-2 rounded border border-[var(--gw-control)] p-3">
      <h4 className="text-sm font-semibold text-[var(--gw-text)]">{t('step4.group.policy.title')}</h4>
      <p className="text-xs text-[var(--gw-muted)]">{t('step4.group.policy.note')}</p>
      {query.isError && (
        <p role="alert" data-testid="counter-policy-error" className="text-xs text-[var(--gw-warning-fg)]">
          {t('step4.group.policy.load_failed')}
          <button type="button" onClick={() => void query.refetch()} className="ml-2 underline">{t('step4.group.retry')}</button>
        </p>
      )}
      {query.isLoading && <p role="status" className="text-xs text-[var(--gw-muted)]">{t('step4.group.policy.loading')}</p>}
      {!query.isLoading && !query.isError && members.map((member) => {
        const source = resolvePolicyDefinition(workspaceId, member, definitions);
        const label = labels[member.tag_id] ?? member.tag_id;
        const outcomeKey = source.kind === 'editable' ? `${workspaceId}|${source.def.id}` : '';
        return source.kind === 'editable'
          ? <SourcePolicy key={`${source.def.id}:${source.def.definition_revision}:${source.def.series_epoch}`} def={source.def} confirmed={source.confirmed} label={label}
            readonly={readonly} outcome={outcomes[outcomeKey]} onSaved={(outcome) => setOutcomes((current) => ({ ...current, [outcomeKey]: outcome }))} />
          : (
            <div key={member.tag_id} className="space-y-1 text-xs text-[var(--gw-warning-fg)]">
              <p data-testid={`counter-policy-missing-${member.tag_id}`}>{label} · {t('step4.group.policy.missing')}</p>
              <button type="button" data-testid={`counter-policy-open-basic-${member.tag_id}`} className="rounded border border-[var(--gw-warning-border)] px-2 py-1"
                onClick={() => setNotOnScreen((current) => ({ ...current, [member.tag_id]: !focusCounterConfirmation(member.tag_id) }))}>
                {t('step4.group.policy.open_basic')}
              </button>
              {notOnScreen[member.tag_id] && <p role="status" data-testid={`counter-policy-open-basic-unavailable-${member.tag_id}`}>{t('step4.group.policy.open_basic_unavailable')}</p>}
            </div>
          );
      })}
    </section>
  );
};
