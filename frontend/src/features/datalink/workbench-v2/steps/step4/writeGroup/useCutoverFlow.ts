import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { useApplyCutoverMutation, useCheckCutovers, useCutoverReviewMutation } from '../../../../../../hooks/datalink/useStudioV2WriteGroupCutover';
import type { WriteGroupCutover, WriteGroupCutoverRequest, WriteGroupCutoverReview } from '../../../../../../types/studioV2WriteGroupCutover';
import { applyRequestFor, buildCutoverRequest, clearCutoverIntent, cutoverFingerprint, findRecordedCutover, loadCutoverIntent, saveCutoverIntent, type CutoverScope } from '../../../state/cutoverFlow';
import { cutoverErrorMessage, isCutoverUnknown } from '../../../state/cutoverMessages';

export type CutoverOutcome =
  | { state: 'idle' }
  | { state: 'applied'; cutover: WriteGroupCutover }
  | { state: 'failed'; message: string }
  /** The answer was lost: nothing is assumed until the recorded cutovers are read. */
  | { state: 'unknown'; checked: 'never' | 'none' | 'failed' };

/**
 * One review and one confirmation of a cutover. A review belongs to exactly the selection, boundary and saved
 * revisions it was made for; confirming sends the reviewed boundary text and digest and nothing recomputed.
 */
export function useCutoverFlow(scope: CutoverScope | undefined, readonly: boolean, workspaceId?: string) {
  const { t } = useTranslation('workbench-v2');
  const reviewMutation = useCutoverReviewMutation();
  const applyMutation = useApplyCutoverMutation();
  const checkCutovers = useCheckCutovers();
  const effectiveWorkspaceId = workspaceId || scope?.workspaceId;
  const initialIntent = React.useMemo(() => {
    return effectiveWorkspaceId ? loadCutoverIntent(effectiveWorkspaceId) : undefined;
  }, [effectiveWorkspaceId]);

  const [reviewed, setReviewed] = React.useState<{ fingerprint: string; review: WriteGroupCutoverReview }>();
  const [frozen, setFrozen] = React.useState<{ sourceId: string; request: WriteGroupCutoverRequest } | undefined>(
    () => initialIntent ? { sourceId: initialIntent.source_id, request: initialIntent.request } : undefined,
  );
  const [outcome, setOutcome] = React.useState<CutoverOutcome>(
    () => initialIntent ? { state: 'unknown', checked: 'never' } : { state: 'idle' },
  );
  const [error, setError] = React.useState<string>();
  const [busy, setBusy] = React.useState(false);
  const busyRef = React.useRef(false);
  const mounted = React.useRef(true);
  React.useEffect(() => { mounted.current = true; return () => { mounted.current = false; }; }, []);

  const fingerprint = scope?.boundary ? cutoverFingerprint(scope) : '';
  const stale = Boolean(reviewed && reviewed.fingerprint !== fingerprint);
  const current = reviewed && !stale ? reviewed.review : undefined;
  const canConfirm = !readonly && !busy && Boolean(scope && current?.applicable) && outcome.state === 'idle';

  const run = async (action: () => Promise<void>) => {
    if (busyRef.current) return;
    busyRef.current = true;
    setBusy(true);
    try { await action(); } finally { busyRef.current = false; if (mounted.current) setBusy(false); }
  };

  const send = async (sourceId: string, request: WriteGroupCutoverRequest) => {
    try {
      const cutover = await applyMutation.mutateAsync({ sourceId, request });
      if (effectiveWorkspaceId) clearCutoverIntent(effectiveWorkspaceId, sourceId);
      setReviewed(undefined);
      setOutcome({ state: 'applied', cutover });
    } catch (cause) {
      if (isCutoverUnknown(cause)) {
        setOutcome({ state: 'unknown', checked: 'never' });
        return;
      }
      if (effectiveWorkspaceId) clearCutoverIntent(effectiveWorkspaceId, sourceId);
      setReviewed(undefined);
      setOutcome({ state: 'failed', message: cutoverErrorMessage(cause, t) });
    }
  };

  const review = () => run(async () => {
    if (!scope || !fingerprint || readonly || outcome.state === 'unknown') return;
    setError(undefined); setOutcome({ state: 'idle' }); setFrozen(undefined);
    try {
      const data = await reviewMutation.mutateAsync({ sourceId: scope.source.id, request: buildCutoverRequest(scope) });
      setReviewed({ fingerprint, review: data });
    } catch (cause) {
      setReviewed(undefined);
      setError(cutoverErrorMessage(cause, t));
    }
  });

  const confirm = () => run(async () => {
    if (!scope || !current?.applicable || readonly || outcome.state !== 'idle') return;
    const request = applyRequestFor(scope, current);
    setFrozen({ sourceId: scope.source.id, request });
    if (effectiveWorkspaceId) {
      saveCutoverIntent(effectiveWorkspaceId, {
        workspace_id: effectiveWorkspaceId,
        source_id: scope.source.id,
        kind: scope.kind,
        boundary: request.boundary ?? scope.boundary,
        review_digest: request.review_digest ?? '',
        expected_workspace_revision: request.expected_workspace_revision ?? '',
        expected_source_revision: request.expected_source_revision ?? '',
        expected_connector_revision: request.expected_connector_revision ?? '',
        request,
        created_at: new Date().toISOString(),
      });
    }
    await send(scope.source.id, request);
  });

  /** Reads the recorded cutovers and settles a lost answer by what was reviewed, never by recency. */
  const check = () => run(async () => {
    if (!frozen) return;
    try {
      const found = findRecordedCutover(await checkCutovers(frozen.sourceId), frozen.request, frozen.sourceId);
      if (found) {
        if (effectiveWorkspaceId) clearCutoverIntent(effectiveWorkspaceId, frozen.sourceId);
        setReviewed(undefined);
        setOutcome({ state: 'applied', cutover: found });
      } else {
        setOutcome({ state: 'unknown', checked: 'none' });
      }
    } catch {
      setOutcome({ state: 'unknown', checked: 'failed' });
    }
  });

  /** Only after a check found nothing recorded: the identical reviewed request, never a changed one. */
  const resend = () => run(async () => {
    if (!frozen || readonly || outcome.state !== 'unknown' || outcome.checked !== 'none') return;
    await send(frozen.sourceId, frozen.request);
  });

  const resetReview = React.useCallback(() => {
    setReviewed(undefined);
    setError(undefined);
  }, []);

  return { review, confirm, check, resend, resetReview, busy, stale, reviewed: reviewed?.review, canConfirm, outcome, error };
}
