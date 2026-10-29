import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { RecordingScopeSummary } from '@/features/datalink/workbench-v2/steps/step4/RecordingScopeSummary';
import type { RecordingScopeSummaryProps, RecordingScopeView } from '@/features/datalink/workbench-v2/steps/step4/recordingScopeSummary';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { useRuntimeRuleEvidence } from '@/hooks/datalink/useRuntimeRuleEvidence';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import type { EvidenceScope, RuleEvidenceView } from '@/types/runtimeRuleEvidence';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { evidenceContext, ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() } }));
const clients: QueryClient[] = [];
const selected: EvidenceScope = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', device_id: 'dev-1' };
const state = savedGroupState();
state.devices.push({ ...state.devices[0], id: 'dev-2', name: 'Empty device B' });
state.mappings['rule-1-p-0'].display_name = 'Process temperature';
state.db.connector.database = '/private/customer.sqlite';
const source = buildGroupCandidates(state).candidates[0];
function selection(device = 'dev-1', groupIds: string[] = ['G']): RecordingScopeView {
  return { workspaceId: 'ws-1', workspaceRevision: 'wrev-1', deviceId: device, groupIds,
    draft: { profiles: [{ kind: 'instantaneous', interval: '60', unit: 'second', reduction: 'mean', sources: [source] }],
      layout: 'wide', storageMode: 'plain', incompletePolicy: 'partial' } };
}
function group(device = 'dev-1', id = 'G'): WriteGroup {
  return savedGroup({ id, revision: 'C', applied_revision: 'B',
    members: [{ device_id: device, point_id: 'pt-0', tag_id: 'tag-0', source_revision: 'source-C', mapping_revision: 'mapping-C',
      target_column: 'value', required: false, source_value_type: 'float64' }],
    recording_profile: { kind: 'instantaneous', interval: 60, unit: 'second', sample_reduction: 'mean' },
    row_policy: { interval_seconds: 60, allowed_lateness_seconds: 2, incomplete_policy: 'partial' } });
}
function proof(scope = selected, revision = 'A', dual = false): RuleEvidenceView {
  const value = ruleEvidenceFixture(scope, revision);
  for (const facts of [value.groups[0].saved!, value.groups[0].effective!]) {
    facts.members[0].point_id = 'pt-0'; facts.members[0].tag_id = 'tag-0';
  }
  if (dual) value.groups.push({ ...structuredClone(value.groups[0]), group_id: 'G2' });
  return value;
}
function response(value = proof()) { return { data: { success: true, data: evidenceContext(value, value.workspace_id) } }; }
function deferred() {
  let resolve!: (value: ReturnType<typeof response>) => void;
  const promise = new Promise<ReturnType<typeof response>>(done => { resolve = done; });
  return { promise, resolve };
}
type HarnessProps = Pick<RecordingScopeSummaryProps, 'scope' | 'groups' | 'configurationState'>;
function Harness(props: HarnessProps) {
  const evidence = useRuntimeRuleEvidence({ workspace_id: props.scope.workspaceId,
    workspace_revision: props.scope.workspaceRevision, device_id: props.scope.deviceId });
  return <RecordingScopeSummary {...props} state={state} evidence={evidence} />;
}
function mount(props: Partial<HarnessProps> = {}) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  const renderProps: HarnessProps = { scope: selection(), groups: [group()], configurationState: 'ready', ...props };
  const tree = (next: HarnessProps) => <QueryClientProvider client={client}><Harness {...next} /></QueryClientProvider>;
  const view = render(tree(renderProps));
  return { ...view, rerenderScope: (next: Partial<HarnessProps>) => view.rerender(tree({ ...renderProps, ...next })) };
}
beforeEach(() => { vi.clearAllMocks(); });
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });

describe('mounted recording scope summary adapter', () => {
  it('keeps saved C/requested B separate from immutable as-of A, then refreshes to B without editable GET', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response()).mockResolvedValueOnce(response(proof(selected, 'B')));
    mount();
    expect(await screen.findByText('meter_old')).toBeInTheDocument();
    expect(screen.getByText('meter_draft')).toBeInTheDocument();
    expect(screen.getByText(/^60 /)).toBeInTheDocument();
    expect(screen.getByText(/^15 /)).toBeInTheDocument();
    expect(screen.queryByText(/^30 /)).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: /retryEvidence/ }));
    expect(await screen.findByText('meter_new')).toBeInTheDocument();
    expect(screen.queryByText('meter_old')).not.toBeInTheDocument();
    expect(screen.getByText(/^60 /)).toBeInTheDocument();
    expect(screen.getByText(/^30 /)).toBeInTheDocument();
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2);
    expect(vi.mocked(studioV2DatalinkApi.get).mock.calls.every(([url]) => url === '/studio-v2/workspace/runtime-context')).toBe(true);
  });
  it('shares one request across dual groups and never refetches for source/mode/circuit views or a drawer', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response(proof(selected, 'A', true)));
    const scope = selection('dev-1', ['G', 'G2']);
    const view = mount({ scope, groups: [group(), group('dev-1', 'G2')] });
    expect(await screen.findAllByText('meter_old')).toHaveLength(2);
    const changed = { ...scope, draft: { ...scope.draft, layout: 'long' as const, profiles: [] } };
    view.rerenderScope({ scope: changed, groups: [group(), group('dev-1', 'G2')] });
    view.rerenderScope({ scope: { ...changed, circuitId: 'another-view' }, groups: [group(), group('dev-1', 'G2')] });
    view.rerenderScope({ scope, groups: [group(), group('dev-1', 'G2')] });
    fireEvent.click(screen.getAllByRole('button', { name: /^details$/ })[0]);
    expect(screen.getByRole('dialog')).toBeInTheDocument();
    fireEvent.click(within(screen.getByRole('dialog')).getByRole('button', { name: /^close$/ }));
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.put).not.toHaveBeenCalled();
  });
  it('does not let a late A response replace the selected empty B view', async () => {
    const a = deferred(), b = deferred();
    vi.mocked(studioV2DatalinkApi.get).mockReturnValueOnce(a.promise).mockReturnValueOnce(b.promise);
    const view = mount();
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1));
    const scopeB = { ...selection('dev-2', []), draft: { ...selection().draft, profiles: [] } };
    view.rerenderScope({ scope: scopeB, groups: [group()] });
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2));
    await act(async () => { b.resolve(response(proof({ ...selected, device_id: 'dev-2' }))); a.resolve(response()); });
    expect(document.body).not.toHaveTextContent(/meter_old|meter_draft|Process temperature/);
    expect(screen.getByText('noProfiles')).toBeInTheDocument();
  });
  it('shows unsaved native choices without synthetic saved/effective group facts or private destination paths', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response());
    mount({ scope: selection('dev-1', []), groups: [] });
    expect(await screen.findByText('unsaved')).toBeInTheDocument();
    expect(screen.getByText('Process temperature')).toBeInTheDocument();
    expect(screen.getByText(/^60 /)).toBeInTheDocument();
    expect(screen.getByText(/Line DB/)).toBeInTheDocument();
    expect(document.body).not.toHaveTextContent(/customer.sqlite|meter_old|meter_draft/);
    expect(screen.queryByRole('button', { name: /^details$/ })).not.toBeInTheDocument();
  });
  it('does not render invalid draft intervals or unsafe aliases as operator facts', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response());
    const scope = selection('dev-1', []);
    scope.draft.profiles = [{ ...scope.draft.profiles[0], interval: 'password=leak /private/draft' }];
    mount({ scope, groups: [] });
    expect(await screen.findByText('unsaved')).toBeInTheDocument();
    expect(document.body).not.toHaveTextContent(/password=leak|private|customer.sqlite/);
    expect(screen.getByText('unavailable')).toBeInTheDocument();
  });
  it('retains complete shared-group membership but uses human labels only for exact source identities', async () => {
    const value = proof();
    const sharedMember = { device_id: 'dev-2', point_id: 'other-point', tag_id: 'other-tag', column: 'other_column', source_value_type: 'bool' as const, circuit_id: null };
    value.groups[0].saved!.members.push(sharedMember); value.groups[0].effective!.members.push(sharedMember);
    const saved = group(); saved.members.push({ ...sharedMember, circuit_id: undefined, target_column: 'other_column', required: false, source_revision: 'S2', mapping_revision: 'M2' });
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response(value));
    mount({ groups: [saved] });
    expect(await screen.findAllByText('Process temperature')).toHaveLength(2);
    expect(document.body).not.toHaveTextContent(/other-point|other-tag|other_column/);
    fireEvent.click(screen.getByRole('button', { name: /^details$/ }));
    const drawer = screen.getByRole('dialog');
    expect(drawer).toHaveTextContent('other-point'); expect(drawer).toHaveTextContent('other-tag');
    expect(drawer).toHaveTextContent('other_column');
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
  });
  it.each(['loading', 'error', 'unknown'] as const)('does not treat %s saved configuration as an empty new draft', async configurationState => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response());
    mount({ configurationState, scope: selection('dev-1', []) });
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1));
    expect(screen.queryByText('unsaved')).not.toBeInTheDocument();
    expect(document.body).not.toHaveTextContent(/meter_old|meter_draft/);
    expect(screen.getByRole('status')).toBeInTheDocument();
  });
  it.each(['deleted', 'foreign', 'duplicate', 'missing'] as const)('rejects %s selected group rather than picking another', async kind => {
    const invalid = group();
    if (kind === 'deleted') invalid.status = 'deleted';
    if (kind === 'foreign') invalid.workspace_id = 'other-workspace';
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response());
    mount({ groups: kind === 'duplicate' ? [invalid, invalid] : kind === 'missing' ? [] : [invalid] });
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1));
    expect(document.body).not.toHaveTextContent(/meter_old|meter_draft|Process temperature/);
    expect(screen.queryByText('unsaved')).not.toBeInTheDocument();
    expect(screen.getByRole('status')).toBeInTheDocument();
  });
  it.each(['stale', 'foreign', 'malformed', 'offline'] as const)('fails closed for %s evidence without raw diagnostics or saved-as-effective substitution', async kind => {
    const value = proof();
    if (kind === 'stale') value.workspace_revision = 'older';
    if (kind === 'foreign') value.workspace_id = 'other';
    if (kind === 'malformed') value.groups[0].effective!.destination.table = '/private/raw.sql';
    if (kind === 'offline') vi.mocked(studioV2DatalinkApi.get).mockRejectedValue(new Error('password=secret /private/raw.sql'));
    else vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response(value));
    mount();
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1));
    await waitFor(() => expect(screen.getByRole('button', { name: /retryEvidence/ })).not.toBeDisabled());
    expect(document.body).not.toHaveTextContent(/meter_old|meter_draft|raw.sql|secret|customer.sqlite/);
  });
});
