import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import { safeBasicStartError } from '../../state/basicRecordingMessages';
import {
  type BatchTargetExecutionResult,
  reconcileBatchResults,
} from '../../state/circuitProfileBatchOutcome';
import { BatchPlanError, executeBatchPlan, type ExecuteBatchPlanInput, type FrozenBatchPlan } from './circuitProfileBatchExecutor';

interface BatchRunInput {
  plan: FrozenBatchPlan | null;
  stale: boolean;
  readonly: boolean;
  connectorRevision?: string;
  workspaceRevision: string;
  savedLongGroups: WriteGroup[];
  updateWriteGroup: ExecuteBatchPlanInput['updateWriteGroup'];
  ensureCircuitProfile?: ExecuteBatchPlanInput['ensureCircuitProfile'];
  onSuccess: (updatedGroups: WriteGroup[], nextWorkspaceRevision: string) => void | Promise<void>;
  onClose: () => void;
  onVerify?: () => void | Promise<void>;
  onBusyChange?: (busy: boolean) => void;
  onOperationActiveChange?: (active: boolean) => void;
}

/**
 * Runs one reviewed plan. The lock is held until the parent has reloaded the server authority, a failed
 * reload is reported separately from the saves it follows, and a stale reply never touches a newer editor.
 */
export function useCircuitProfileBatchRun(input: BatchRunInput) {
  const { t } = useTranslation('workbench-v2');
  const { plan, savedLongGroups, workspaceRevision, onSuccess } = input;
  const [saving, setSaving] = React.useState(false);
  const [finished, setFinished] = React.useState(false);
  const [rawResults, setRawResults] = React.useState<Record<string, BatchTargetExecutionResult>>({});
  const [errorMessage, setErrorMessage] = React.useState<string>();
  const [refreshFailed, setRefreshFailed] = React.useState(false);
  const [verifying, setVerifying] = React.useState(false);
  const inFlight = React.useRef(false);
  const mounted = React.useRef(true);
  const reported = React.useRef(new Set<string>());
  // Proof found later is reported to the host that launched the run, which knows what existed before it.
  const launchHost = React.useRef(onSuccess);
  React.useEffect(() => { mounted.current = true; return () => { mounted.current = false; }; }, []);

  const reconciled = React.useMemo(() => (plan
    ? reconcileBatchResults(plan.targets, { instantaneous: plan.instInterval, cumulative: plan.cumInterval }, rawResults, savedLongGroups, plan.workspaceId)
    : { results: rawResults, proven: [] as WriteGroup[] }), [plan, rawResults, savedLongGroups]);
  const hasUnknown = Object.values(reconciled.results).some((target) => Object.values(target.profiles).some((profile) => profile?.status === 'unknown'));

  // A profile proven by the reloaded server state is reported once so the local drafts follow the server.
  React.useEffect(() => {
    const fresh = reconciled.proven.filter((group) => !reported.current.has(`${group.id}@${group.revision}`));
    if (fresh.length === 0) return;
    fresh.forEach((group) => reported.current.add(`${group.id}@${group.revision}`));
    void Promise.resolve(launchHost.current(fresh, workspaceRevision)).catch(() => setRefreshFailed(true));
  }, [reconciled.proven, workspaceRevision]);

  const reset = () => {
    setRawResults({}); setErrorMessage(undefined); setRefreshFailed(false); setFinished(false); reported.current.clear();
  };

  const confirm = async () => {
    if (inFlight.current || input.readonly || saving || finished || !plan || input.stale) return;
    inFlight.current = true;
    setSaving(true); setErrorMessage(undefined); setRefreshFailed(false);
    launchHost.current = input.onSuccess;
    input.onBusyChange?.(true);
    input.onOperationActiveChange?.(true);
    let outcome: { updatedGroups: WriteGroup[]; workspaceRevision: string; targetResults: Record<string, BatchTargetExecutionResult> };
    let failure: string | undefined;
    try {
      outcome = await executeBatchPlan({ plan, connectorRevision: input.connectorRevision, updateWriteGroup: input.updateWriteGroup, ensureCircuitProfile: input.ensureCircuitProfile });
    } catch (error: unknown) {
      if (error instanceof BatchPlanError) {
        outcome = error;
        const failed = Object.values(error.targetResults).flatMap((target) => Object.values(target.profiles)).find((profile) => profile?.status === 'failed');
        if (failed) failure = safeBasicStartError(failed.cause, t) || t('errors.internal');
      } else {
        outcome = { updatedGroups: [], workspaceRevision: plan.workspaceRevision, targetResults: {} };
        failure = t('errors.internal');
      }
    }
    let reloaded = true;
    try {
      if (outcome.updatedGroups.length > 0 && outcome.workspaceRevision) await input.onSuccess(outcome.updatedGroups, outcome.workspaceRevision);
    } catch {
      reloaded = false; // The writes are already committed; only the reload failed, so nothing is resent.
    }
    if (mounted.current) {
      setRawResults(outcome.targetResults); setErrorMessage(failure); setRefreshFailed(!reloaded); setFinished(true); setSaving(false);
    }
    inFlight.current = false;
    input.onBusyChange?.(false);
    input.onOperationActiveChange?.(false);
    if (mounted.current && !failure && reloaded && !Object.values(outcome.targetResults).some((target) => target.status === 'unknown')) input.onClose();
  };

  const verify = async () => {
    if (!input.onVerify || verifying) return;
    setVerifying(true); setRefreshFailed(false);
    try { await input.onVerify(); } catch { if (mounted.current) setRefreshFailed(true); } finally { if (mounted.current) setVerifying(false); }
  };

  return { saving, finished, results: reconciled.results, errorMessage, hasUnknown, refreshFailed, verifying, confirm, verify, reset };
}
