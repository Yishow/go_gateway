import { useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import type { StudioLocatorEntryRead } from '../../../../types/studioV2LocatorEntry';
import type { DeviceSourcePort } from '../../../../pages/datalink/workbench-v2/deviceSourceSubmission';
import { readStudioLocatorEditorEntry } from '../../../../services/studioV2LocatorEntry';
import { useModbusShareStatusQuery } from '../../../../hooks/datalink/useModbusShareStatus';
import { Button } from '../../../../components/ui/button';
import { SavedRecordingStartReview } from './SavedRecordingStartReview';
import { savedRecordingStartSourceIdentity } from './savedRecordingStartScope';
import { locatorSavedRecordingStartScope } from './locatorSavedRecordingStartScope';
interface Props {
  entry: StudioLocatorEntryRead; workspaceId: string; readonly: boolean; canUseScope: () => boolean;
  deviceSource?: DeviceSourcePort; onClose: () => void; onCompleted: () => void; onOperationActiveChange: (active: boolean) => void;
}
/** This host only exists after explicit review of an already hydrated scoped editor. */
export function LocatorSavedRecordingStartReview({ entry, workspaceId, readonly, canUseScope, deviceSource, onClose, onCompleted, onOperationActiveChange }: Props) {
  const { t } = useTranslation('recording-editor');
  const context = locatorSavedRecordingStartScope(entry, workspaceId);
  const share = useModbusShareStatusQuery(Boolean(context) && !readonly && canUseScope());
  const live = useRef({ entry, readonly, canUseScope }); live.current = { entry, readonly, canUseScope };
  const mounted = useRef(true), requests = useRef(new Set<AbortController>());
  useEffect(() => { mounted.current = true; const pending = requests.current; return () => { mounted.current = false; pending.forEach(request => request.abort()); }; }, []);
  const current = () => mounted.current && live.current.entry === entry && !live.current.readonly && live.current.canUseScope();
  const freshShare = !share.isError && !share.isFetching && share.data?.workspace_id === workspaceId && share.data.readiness === true &&
    share.data.readiness_token && share.data.settings_revision && share.data.workspace_revision ? share.data : undefined;
  if (!context || !freshShare) return <section role={share.isFetching ? 'status' : 'alert'} className="rounded-lg border border-[var(--gw-control)] bg-[var(--gw-surface)] p-4 text-[var(--gw-text)]">
    <p>{t('start.changed')}</p><Button className="border border-[var(--gw-control)] bg-[var(--gw-surface)] text-[var(--gw-text)] hover:bg-[var(--gw-active)] focus-visible:outline-[var(--gw-focus)]" disabled={!current() || share.isFetching} onClick={() => void share.refetch()}>{t('start.check')}</Button>
    <Button className="border border-[var(--gw-control)] bg-[var(--gw-surface)] text-[var(--gw-text)] hover:bg-[var(--gw-active)] focus-visible:outline-[var(--gw-focus)]" onClick={onClose}>{t('start.close')}</Button>
  </section>;
  return <SavedRecordingStartReview {...context} readonly={readonly || !current()} canUseScope={current} shareStatus={freshShare} deviceSource={deviceSource}
    onClose={onClose} onOperationActiveChange={onOperationActiveChange} onCompleted={() => { if (current()) onCompleted(); }}
    fetchLatest={async () => {
      if (!current()) return undefined;
      const request = new AbortController(); requests.current.add(request);
      try {
        const result = await readStudioLocatorEditorEntry(entry.locator, { workspaceRevision: context.workspaceRevision, groupRevision: context.group.revision }, request.signal);
        if (!current() || !result.ok) return undefined;
        const fresh = locatorSavedRecordingStartScope(result.value, workspaceId);
        return fresh && fresh.group.id === context.group.id && fresh.group.revision === context.group.revision &&
          fresh.workspaceRevision === context.workspaceRevision &&
          savedRecordingStartSourceIdentity(fresh.group, fresh.state) === savedRecordingStartSourceIdentity(context.group, context.state)
          ? result.value.entry.write_groups : undefined;
      } finally { requests.current.delete(request); }
    }} />;
}
