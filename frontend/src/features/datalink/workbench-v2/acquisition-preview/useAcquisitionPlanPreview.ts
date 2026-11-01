import { useEffect, useMemo, useRef, useState } from 'react';
import { readAcquisitionPlanPreview } from '../../../../services/acquisitionPlanPreview';
import type { AcquisitionPlanPreviewResult, AcquisitionPreviewScope } from '../../../../types/acquisitionPlanPreview';

/** Ephemeral read receipt only: no installation, write, rebase or recording-state authority. */
export function useAcquisitionPlanPreview(scope?: AcquisitionPreviewScope) {
  const [refreshId, setRefreshId] = useState(0);
  const [receipt, setReceipt] = useState<{ owner: object; refreshId: number; result: AcquisitionPlanPreviewResult }>();
  const head = useRef<{ owner: object; revision?: string } | undefined>(undefined);
  const workspace = scope?.workspace_id, device = scope?.device_id, revision = scope?.expected_workspace_revision;
  const requestedHead = scope?.expected_preview_revision;
  const canRead = Boolean(workspace && device && typeof revision === 'string');
  const key = JSON.stringify([workspace, device, revision, requestedHead]);
  const owner = useMemo(() => ({ key }), [key]);
  useEffect(() => {
    if (!canRead || !workspace || !device || revision === undefined) { head.current = undefined; return; }
    const controller = new AbortController();
    const expectedHead = head.current?.owner === owner ? head.current.revision : requestedHead;
    const captured: AcquisitionPreviewScope = { workspace_id: workspace, device_id: device, expected_workspace_revision: revision,
      ...(expectedHead === undefined ? {} : { expected_preview_revision: expectedHead }) };
    const accept = (result: AcquisitionPlanPreviewResult) => {
      if (controller.signal.aborted) return;
      if (result.ok && (result.value.workspace_id !== workspace || result.value.device_id !== device || result.value.workspace_revision !== revision))
        result = { ok: false, reason: 'unavailable' };
      head.current = { owner, revision: result.ok && result.value.status === 'available' ? result.value.preview_revision : undefined };
      setReceipt({ owner, refreshId, result });
    };
    void readAcquisitionPlanPreview(captured, controller.signal).then(accept, () => accept({ ok: false, reason: 'unavailable' }));
    return () => controller.abort();
  }, [canRead, device, owner, refreshId, requestedHead, revision, workspace]);
  const current = canRead && receipt?.owner === owner && receipt.refreshId === refreshId ? receipt.result : undefined;
  return { result: current, loading: canRead && !current, refresh: () => setRefreshId(value => value + 1) };
}
