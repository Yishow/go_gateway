import { useId } from 'react';
import { recordingNextAction } from '../../state/recordingWorkflowView';
import type { RecordingNextActionInput, RecordingNextIntent, RecordingNextIntentKind, RecordingNextReason } from '../../state/recordingWorkflowView';

export interface RecordingNextActionLabels {
  readonly heading: string;
  readonly scopes: Readonly<Record<'draft' | 'group' | 'managed_device', string>>;
  readonly actions: Readonly<Record<RecordingNextIntentKind, string>>;
  readonly reasons: Readonly<Record<RecordingNextReason, string>>;
}
export interface RecordingNextActionProps {
  readonly input: RecordingNextActionInput;
  readonly labels: RecordingNextActionLabels;
  readonly scopeLabel: string;
  readonly disabled?: boolean;
  readonly blockedReason?: string;
  readonly testId?: string;
  /** A disabled command still has an actionable local repair target. */
  readonly onRepair?: () => void;
  readonly repairLabel?: string;
  /** Intents distinguish navigation from explicit owner commands; neither grants mutation authority. */
  readonly onIntent: (intent: RecordingNextIntent) => void;
}

export function RecordingNextAction({ input, labels, scopeLabel, onIntent, disabled, blockedReason, testId, onRepair, repairLabel }: RecordingNextActionProps) {
  const statusId = useId();
  const model = recordingNextAction(input);
  return <section aria-label={labels.heading} className="space-y-3 rounded border border-[var(--gw-border)] p-3">
    <h3 className="text-sm font-medium text-[var(--gw-text)]">{labels.heading}</h3>
    <p className="break-words text-sm text-[var(--gw-text)]">{scopeLabel} — {labels.scopes[model.configurationScope.kind]}</p>
    <p id={statusId} role="status" className="text-sm text-[var(--gw-text)]">{model.reasons.map((reason) => labels.reasons[reason]).join(' ')} {blockedReason}</p>
    <button type="button" data-primary="true" data-testid={testId} aria-describedby={statusId} disabled={!model.enabled || disabled}
      className="rounded bg-[var(--gw-primary)] px-3 py-2 text-sm text-[var(--gw-on-primary)] disabled:text-[var(--gw-muted)] disabled:bg-[var(--gw-ground)]"
      onClick={() => onIntent(model.intent)}>{labels.actions[model.intent.kind]}</button>
    {model.intent.kind !== 'inspect' && <button type="button"
      className="ml-2 rounded border border-[var(--gw-control)] px-3 py-2 text-sm text-[var(--gw-text)]"
      onClick={() => onIntent({ kind: 'inspect', scope: model.intent.scope })}>{labels.actions.inspect}</button>}
    {(disabled || !model.enabled) && onRepair && <button type="button" className="ml-2 text-sm underline"
      onClick={onRepair}>{repairLabel ?? labels.actions.inspect}</button>}
  </section>;
}
