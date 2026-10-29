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
  return <section aria-label={labels.heading} className="space-y-3 rounded border border-slate-700 p-3">
    <h3 className="text-sm font-medium text-slate-100">{labels.heading}</h3>
    <p className="break-words text-sm text-slate-300">{scopeLabel} — {labels.scopes[model.configurationScope.kind]}</p>
    <p id={statusId} role="status" className="text-sm text-slate-300">{model.reasons.map((reason) => labels.reasons[reason]).join(' ')} {blockedReason}</p>
    <button type="button" data-primary="true" data-testid={testId} aria-describedby={statusId} disabled={!model.enabled || disabled}
      className="rounded bg-blue-600 px-3 py-2 text-sm text-white disabled:opacity-50"
      onClick={() => onIntent(model.intent)}>{labels.actions[model.intent.kind]}</button>
    {model.intent.kind !== 'inspect' && <button type="button"
      className="ml-2 rounded border border-slate-600 px-3 py-2 text-sm text-slate-200"
      onClick={() => onIntent({ kind: 'inspect', scope: model.intent.scope })}>{labels.actions.inspect}</button>}
    {(disabled || !model.enabled) && onRepair && <button type="button" className="ml-2 text-sm underline"
      onClick={onRepair}>{repairLabel ?? labels.actions.inspect}</button>}
  </section>;
}
