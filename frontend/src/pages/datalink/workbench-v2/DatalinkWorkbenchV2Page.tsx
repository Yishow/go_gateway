import { RecordingHomeSourceHost } from '../../../features/datalink/workbench-v2/recording-home/RecordingHomeSourceHost';
import { RecordingHomeRoute } from '../../../features/datalink/workbench-v2/recording-home/RecordingHomeRoute';
import { parseStudioPresentation, type StudioPresentation } from '../../../services/studioRecordingView';
import type { StudioEntityLocatorIntent } from '../../../services/studioEntityLocator';
import type { WorkbenchV2RuntimeReturnFocus } from '../../../features/datalink/workbench-v2/shell/WorkbenchV2Shell';
import { lazy } from 'react';
import { DeferredWorkbenchBoundary } from './DeferredWorkbenchBoundary';
const AdvancedWorkbenchV2Page = lazy(() => import('./AdvancedWorkbenchV2Page'));
import '../../../features/datalink/workbench-v2/styles/workbench-v2.css';

export interface DatalinkWorkbenchV2PageProps {
  presentation?: StudioPresentation;
  presentationLocation?: { search: string; hash: string };
  navigateTo?: (target: string, options?: { replace?: boolean }) => void;
  locatorIntent?: StudioEntityLocatorIntent;
  runtimeReturnFocus?: WorkbenchV2RuntimeReturnFocus | null;
}

/** Precise locator/runtime intent always outranks an optional presentation hint. */
export default function DatalinkWorkbenchV2Page(props: DatalinkWorkbenchV2PageProps) {
  const locator = props.locatorIntent?.kind;
  if (props.presentation === 'recordings' && (!locator || locator === 'absent') && !props.runtimeReturnFocus) {
    return <RecordingHomeSourceHost renderSource={controls => <DeferredWorkbenchBoundary key={controls.target.epoch} onClose={controls.onClose}><AdvancedWorkbenchV2Page homeSource={controls} navigateTo={props.navigateTo}
      locatorIntent={controls.target.kind === 'create' ? { kind: 'absent' } : { kind: 'locator', locator: { workspace_id: controls.target.workspaceId, step: controls.target.ruleId ? 2 : 1, device_id: controls.target.deviceId,
        rule_id: controls.target.ruleId ?? null, mapping_id: null, group_id: null, group_revision: null } }} /></DeferredWorkbenchBoundary>} >
      {(open, workspaceChanged, create) => <RecordingHomeRoute location={props.presentationLocation ?? { search: '?view=recordings', hash: '' }} onNavigate={props.navigateTo}
        onOpenDevice={open} onCreateDevice={create} onWorkspaceChange={workspaceChanged} />}
    </RecordingHomeSourceHost>;
  }
  const intent = parseStudioPresentation(props.presentationLocation?.search ?? '');
  const simple = props.presentation === 'recording' && !props.runtimeReturnFocus && intent.presentation === 'recording';
  return <DeferredWorkbenchBoundary><AdvancedWorkbenchV2Page {...props} recordingIntent={simple ? intent.recordingIntent === 'copy' ? 'copy' : locator === 'locator' ? 'located' : 'new' : undefined} /></DeferredWorkbenchBoundary>;
}
