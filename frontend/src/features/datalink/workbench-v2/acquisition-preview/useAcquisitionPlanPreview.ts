import { useEffect, useMemo, useRef, useState } from 'react';
import { readAcquisitionPreviewWorkspaceHead } from '../../../../services/acquisitionPreviewWorkspaceHead';
import { readAcquisitionPlanPreview } from '../../../../services/acquisitionPlanPreview';
import type { AcquisitionPlanPreviewResult, AcquisitionPreviewScope } from '../../../../types/acquisitionPlanPreview';

/** Ephemeral read receipt only: no installation, write, rebase or recording-state authority. */
export function useAcquisitionPlanPreview(scope?: AcquisitionPreviewScope, refreshWorkspaceHead = false) {
  const [refreshId, setRefreshId] = useState(0);
  const [receipt, setReceipt] = useState<{ owner: object; refreshId: number; result: AcquisitionPlanPreviewResult }>();
  const refreshOwner = useRef<object | undefined>(undefined);
  const head = useRef<{ owner: object; workspaceRevision: string; revision?: string } | undefined>(undefined);
  const workspace = scope?.workspace_id, device = scope?.device_id, revision = scope?.expected_workspace_revision;
  const requestedHead = scope?.expected_preview_revision;
  const canRead = Boolean(workspace && device && typeof revision === 'string');
  const key = JSON.stringify([workspace, device, revision, requestedHead, refreshWorkspaceHead]);
  const owner = useMemo(() => ({ key }), [key]);
  useEffect(() => {
    if (!canRead || !workspace || !device || revision === undefined) { head.current = undefined; return; }
    const controller = new AbortController();
    const read = async () => {
      const fresh = refreshWorkspaceHead && refreshOwner.current === owner
        ? await readAcquisitionPreviewWorkspaceHead(workspace, controller.signal) : { ok: true as const, revision };
      if (controller.signal.aborted) return;
      if (!fresh.ok) { setReceipt({ owner, refreshId, result: fresh }); head.current = undefined; return; }
      const readRevision = fresh.revision;
      const expectedHead = head.current?.owner === owner && head.current.workspaceRevision === readRevision
        ? head.current.revision : readRevision === revision ? requestedHead : undefined;
      const captured: AcquisitionPreviewScope = { workspace_id: workspace, device_id: device, expected_workspace_revision: readRevision,
        ...(expectedHead === undefined ? {} : { expected_preview_revision: expectedHead }) };
      let result: AcquisitionPlanPreviewResult;
      try { result = await readAcquisitionPlanPreview(captured, controller.signal); }
      catch { result = { ok: false, reason: 'unavailable' }; }
      if (controller.signal.aborted) return;
      if (result.ok && (result.value.workspace_id !== workspace || result.value.device_id !== device || result.value.workspace_revision !== readRevision))
        result = { ok: false, reason: 'unavailable' };
      head.current = { owner, workspaceRevision: readRevision, revision: result.ok && result.value.status === 'available' ? result.value.preview_revision : undefined };
      setReceipt({ owner, refreshId, result });
    };
    void read();
    return () => controller.abort();
  }, [canRead, device, owner, refreshId, requestedHead, revision, workspace, refreshWorkspaceHead]);
  const current = canRead && receipt?.owner === owner && receipt.refreshId === refreshId ? receipt.result : undefined;
  return { result: current, loading: canRead && !current, refresh: () => { refreshOwner.current = owner; setRefreshId(value => value + 1); } };
}
