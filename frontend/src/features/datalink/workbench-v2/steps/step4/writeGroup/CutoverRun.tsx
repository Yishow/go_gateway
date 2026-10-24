import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { CutoverScope } from '../../../state/cutoverFlow';
import type { GroupCandidate } from '../../../state/writeGroup/candidates';
import { CutoverOutcomeView } from './CutoverOutcomeView';
import { CutoverReviewView } from './CutoverReviewView';
import { useCutoverFlow } from './useCutoverFlow';

export interface CutoverRunProps {
  /** Undefined until a source, targets (or a switch to roll back) and a boundary are chosen. */
  scope?: CutoverScope;
  readonly: boolean;
  groupName: (groupId: string) => string;
  onOpenGroup: (groupId: string) => void;
  flow?: ReturnType<typeof useCutoverFlow>;
  candidates?: GroupCandidate[];
}

/** Review, confirm and settle one cutover. Nothing is sent to the server except the review until the operator confirms. */
export const CutoverRun: React.FC<CutoverRunProps> = ({ scope, readonly, groupName, onOpenGroup, flow: propFlow, candidates }) => {
  const { t } = useTranslation('workbench-v2');
  const localFlow = useCutoverFlow(scope, readonly);
  const flow = propFlow ?? localFlow;
  const rollback = scope?.kind === 'rollback';
  const isUnknown = flow.outcome.state === 'unknown';
  return (
    <div className="space-y-3">
      <button type="button" data-testid="cutover-review-button" disabled={readonly || flow.busy || !scope?.boundary || isUnknown} onClick={() => void flow.review()} className="rounded border border-slate-600 px-3 py-2 text-xs text-slate-200 disabled:opacity-40">
        {flow.busy && !flow.reviewed ? t('step4.cutover.reviewing') : t('step4.cutover.review')}
      </button>
      {flow.error && <p role="alert" data-testid="cutover-review-error" className="text-xs text-rose-300">{flow.error}</p>}
      {flow.reviewed && <CutoverReviewView review={flow.reviewed} stale={flow.stale} onOpenGroup={onOpenGroup} candidates={candidates} sourceGroup={scope?.source} />}
      {flow.reviewed && flow.outcome.state === 'idle' && (
        <button type="button" data-testid="cutover-confirm" disabled={!flow.canConfirm} onClick={() => void flow.confirm()} className="rounded bg-cyan-600 px-4 py-2 text-xs font-semibold text-white disabled:cursor-not-allowed disabled:opacity-40">
          {flow.busy ? t('step4.cutover.confirming') : t(rollback ? 'step4.cutover.confirm_rollback' : 'step4.cutover.confirm')}
        </button>
      )}
      <CutoverOutcomeView outcome={flow.outcome} groupName={groupName} busy={flow.busy} readonly={readonly} onCheck={() => void flow.check()} onResend={() => void flow.resend()} />
    </div>
  );
};
