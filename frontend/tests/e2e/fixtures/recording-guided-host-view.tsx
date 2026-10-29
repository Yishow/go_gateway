import { useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { WorkbenchV2Shell } from '../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell';
import { useWorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { resolveStudioSelection, type StudioSelectionEvidence, type StudioSelectionResult } from '../../../src/features/datalink/workbench-v2/state/studioSelection';
import { useWriteGroupsQuery } from '../../../src/hooks/datalink/useStudioV2WriteGroups';
import type { WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import { recordingGuidedHostState, type GuidedScenario } from './recording-guided-state';

/** URL-independent contract harness only, NOT the future production runtime locator receiver/provider. */
function SelectionContract({ state }: { state: WorkbenchV2State }) {
  const groups = useWriteGroupsQuery(true);
  const [result, setResult] = useState<StudioSelectionResult>();
  const inspect = (kind: 'owned' | 'historical' | 'foreign' | 'stale') => {
    const data = groups.data;
    const evidence: StudioSelectionEvidence = !data || groups.isError || groups.isFetching || data.workspace_id !== 'ws-1'
      ? { state: 'unavailable' } : kind === 'stale' ? { state: 'stale' } : {
        state: 'current', workspace: { id: 'ws-1', ordered_device_ids: state.devices.filter(device => device.persisted).map(device => device.id) },
        workspaceRevision: data.workspace_revision, devices: state.devices.filter(device => device.persisted),
        groups: data.groups.map(group => ({ value: group, currentRevision: group.revision, workspaceRevision: data.workspace_revision })),
      };
    setResult(resolveStudioSelection('ws-1', { workspaceId: kind === 'foreign' ? 'foreign-W' : 'ws-1', step: 4,
      deviceId: 'dev-1', groupId: 'G1', ...(kind === 'historical' ? { groupRevisionHint: 'HISTORICAL-R0' } : {}) }, { readCurrent: () => evidence }));
  };
  return <details data-testid="fixture-selection-contract" className="rounded border border-slate-700 p-2 text-xs">
    <summary>URL-independent selection contract harness — test evidence, not runtime URL integration</summary>
    <p>Uses the parsed canonical group query and fixed owned test-device seed; this harness does not apply the result or authorize a mutation.</p>
    <div className="flex flex-wrap gap-2">{(['owned', 'historical', 'foreign', 'stale'] as const).map(kind =>
      <button type="button" key={kind} data-testid={`fixture-selection-${kind}`} onClick={() => inspect(kind)} className="rounded border p-2">Inspect {kind} selection</button>)}</div>
    {result && <output data-testid="fixture-selection-result">{result.state === 'unavailable' ? `unavailable:${result.reason}`
      : `selected:${result.target.deviceId}:${result.target.groupId}:${result.target.groupRevision}${result.history ? `:history-${result.history.relation}` : ''}`}</output>}
  </details>;
}

/** Actual production Shell/state/Step4 with only fixture navigation destinations and test-data controls. */
export function RecordingGuidedHostView({ scenario }: { scenario: GuidedScenario }) {
  const actions = useWorkbenchV2State(recordingGuidedHostState(scenario));
  const location = useLocation(), navigate = useNavigate();
  const [mount, setMount] = useState(0);
  const deviceId = new URLSearchParams(location.search).get('device_id');
  if (location.pathname === '/studio/runtime') return <main className="min-h-screen bg-slate-950 p-4 text-slate-100" data-testid="fixture-runtime-destination">
    <h1>Runtime navigation destination — test marker, not runtime/effective evidence</h1>
    <output data-testid="fixture-runtime-scope">{actions.state.devices.some(device => device.id === deviceId) ? deviceId : 'unavailable'}</output>
    <button type="button" onClick={() => navigate(-1)}>Back to setup</button>
  </main>;
  if (new URLSearchParams(location.search).get('away') === '1') return <main data-testid="fixture-away">
    <h1>Guided fixture navigation away</h1><button type="button" onClick={() => navigate(-1)}>Back to setup</button>
  </main>;
  return <main data-workbench-v2="true" className="min-h-screen bg-slate-950 text-slate-100">
    <header className="space-y-2 px-4 py-3">
      <h1 className="text-lg font-semibold">Guided recording — real components / mock API / test data</h1>
      <p className="text-xs">No database, device, SQL or effective runtime verification. Scenario: {scenario}</p>
      <div className="flex flex-wrap gap-2 text-xs">
        <button type="button" data-testid="fixture-remount-host" onClick={() => setMount(value => value + 1)} className="rounded border p-2">Remount setup</button>
        <button type="button" data-testid="fixture-leave" onClick={() => navigate(`${location.pathname}?scenario=${scenario}&away=1`)} className="rounded border p-2">Leave fixture</button>
      </div>
      <SelectionContract state={actions.state} />
    </header>
    <WorkbenchV2Shell key={mount} state={actions.state} actions={actions} workspaceId="ws-1" navigateTo={target => navigate(target)} />
  </main>;
}
