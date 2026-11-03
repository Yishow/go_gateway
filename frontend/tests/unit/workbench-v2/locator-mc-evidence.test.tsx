import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { RuleEditor } from '@/features/datalink/workbench-v2/steps/step2/RuleEditor';
import { useStudioV2AutosaveState } from '@/pages/datalink/workbench-v2/useStudioV2AutosaveState';
import { sourceRuleCandidateAPI } from '@/services/sourceRuleCandidates';
import type { StudioLocatorEntryRead } from '@/types/studioV2LocatorEntry';
import { studioV2RulesAPI } from '@/services/studioV2Rules';
import { locatorRead, nativeDevice, nativeRule } from './studio-locator-fixtures';

vi.mock('@/services/sourceRuleCandidates', () => ({ sourceRuleCandidateAPI: { get: vi.fn(), getTags: vi.fn() } }));
vi.mock('@/services/studioV2Rules', () => ({ studioV2RulesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));

const scope = { workspace_id: 'W', workspace_revision: 'W2' };

function entry(): StudioLocatorEntryRead {
  const read = locatorRead(2);
  read.entry.devices = [{ ...nativeDevice(), protocol: 'mc_3e', connection_config: '{"host":"fixture.invalid","port":5000}' }];
  read.entry.source_rules = [{ ...nativeRule(), start_address: 'D0', count: 2 }];
  return read;
}

function Harness({ route }: { route: StudioLocatorEntryRead }) {
  const owner = useStudioV2AutosaveState(false, route, false, scope, 'W');
  const rule = owner.state.rules.find(value => value.id === 'R');
  return rule ? <><output data-testid="rule-save-state">{rule.save_state}:{rule.revision_id}:{rule.save_error}</output>
    <output data-testid="locator-can-leave">{String(owner.locatorHydration.canLeave)}</output>
    <RuleEditor rule={rule} devices={owner.state.devices} globalShareEnabled={false} dispatch={owner.actions.dispatch} /></> : null;
}

function mount(route = entry()) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const content = (next: StudioLocatorEntryRead) => <QueryClientProvider client={client}><Harness route={next} /></QueryClientProvider>;
  const view = render(content(route));
  return { ...view, client, navigate: (next: StudioLocatorEntryRead) => view.rerender(content(next)) };
}

beforeEach(() => {
  vi.clearAllMocks();
  vi.mocked(sourceRuleCandidateAPI.getTags).mockImplementation(async (_id, query) => ({
    source_rule_id: 'R', workspace_id: 'W', workspace_revision: 'W2', revision_id: query!.revision_id,
    tags: { status: 'ready', candidates: ['D0', 'D1'].map((address, index) => ({
      id: `C${index}`, identity: { source_rule_id: 'R', candidate_type: 'tags', candidate_kind: 'tag', derived_from_rule_address: 'D0' },
      proposed_signature: `S${index}`, address, point_id: `P${index}`, tag_key: `tag${index}`, display_name: `Tag ${index}`,
      data_type: 'int16', status: 'active',
    })) },
  }));
});
afterEach(() => { cleanup(); vi.resetAllMocks(); });

it('shows the scoped saved MC points in the cold locator rule editor', async () => {
  mount();
  await waitFor(() => expect(sourceRuleCandidateAPI.getTags).toHaveBeenCalledWith('R', {
    workspace_id: 'W', expected_workspace_revision: 'W2', revision_id: 'R2',
  }));
  await waitFor(() => expect(screen.queryByTestId('rule-mc-persisted-status')).not.toBeInTheDocument());
  expect(screen.getByTestId('range-summary')).toHaveTextContent('D0–D1');
  expect(screen.getByTestId('locator-can-leave')).toHaveTextContent('true');
});

it('reads evidence for an acknowledged rule revision before saving a subsequent edit', async () => {
  let revision = 2;
  vi.mocked(studioV2RulesAPI.update).mockImplementation(async (_id, request) => ({
    ...nativeRule(), start_address: 'D0', count: 2, ...request, revision_id: `R${++revision}`, runtime_apply_status: 'aligned',
  }));
  mount();
  await waitFor(() => expect(screen.getByTestId('range-summary')).toHaveTextContent('D0–D1'));
  await waitFor(() => expect(screen.queryByTestId('rule-mc-persisted-status')).not.toBeInTheDocument());
  fireEvent.change(screen.getByTestId('rule-prefix-input'), { target: { value: 'FIRST_' } });
  await waitFor(() => expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(1));
  await waitFor(() => expect(sourceRuleCandidateAPI.getTags).toHaveBeenCalledWith('R', {
    workspace_id: 'W', expected_workspace_revision: 'W2', revision_id: 'R3',
  }));
  await waitFor(() => expect(screen.queryByTestId('rule-mc-persisted-status')).not.toBeInTheDocument());
  fireEvent.change(screen.getByTestId('rule-prefix-input'), { target: { value: 'SECOND_' } });
  await waitFor(() => expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(2));
  expect(studioV2RulesAPI.update).toHaveBeenLastCalledWith('R', expect.objectContaining({ naming_prefix: 'SECOND_' }));
  await waitFor(() => expect(screen.getByTestId('locator-can-leave')).toHaveTextContent('true'));
});

it('waits for new scoped evidence before sending an edit queued during a save', async () => {
  const getTags = vi.mocked(sourceRuleCandidateAPI.getTags).getMockImplementation()!;
  let finishProof!: () => void;
  vi.mocked(sourceRuleCandidateAPI.getTags).mockImplementation(async (id, query) => {
    const snapshot = await getTags(id, query);
    if (query?.revision_id === 'R3') await new Promise<void>(resolve => { finishProof = resolve; });
    return snapshot;
  });
  let finishFirstSave!: () => void;
  let revision = 2;
  vi.mocked(studioV2RulesAPI.update).mockImplementation(async (_id, request) => {
    const nextRevision = `R${++revision}`;
    if (nextRevision === 'R3') await new Promise<void>(resolve => { finishFirstSave = resolve; });
    return { ...nativeRule(), start_address: 'D0', count: 2, ...request, revision_id: nextRevision, runtime_apply_status: 'aligned' };
  });
  mount();
  await waitFor(() => expect(screen.getByTestId('range-summary')).toHaveTextContent('D0–D1'));
  fireEvent.change(screen.getByTestId('rule-prefix-input'), { target: { value: 'FIRST_' } });
  await waitFor(() => expect(finishFirstSave).toBeTypeOf('function'));
  fireEvent.change(screen.getByTestId('rule-prefix-input'), { target: { value: 'QUEUED_' } });
  await act(async () => finishFirstSave());
  await waitFor(() => expect(finishProof).toBeTypeOf('function'));
  expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(1);
  expect(screen.getByTestId('rule-prefix-input')).toHaveValue('QUEUED_');
  await act(async () => finishProof());
  await waitFor(() => expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(2));
  expect(studioV2RulesAPI.update).toHaveBeenLastCalledWith('R', expect.objectContaining({ naming_prefix: 'QUEUED_' }));
});

it('keeps a subsequent edit blocked when the new evidence belongs to another workspace revision', async () => {
  const getTags = vi.mocked(sourceRuleCandidateAPI.getTags).getMockImplementation()!;
  vi.mocked(sourceRuleCandidateAPI.getTags).mockImplementation(async (id, query) => {
    const snapshot = await getTags(id, query);
    return query?.revision_id === 'R3' ? { ...snapshot, workspace_revision: 'FOREIGN' } : snapshot;
  });
  vi.mocked(studioV2RulesAPI.update).mockImplementation(async (_id, request) => ({
    ...nativeRule(), start_address: 'D0', count: 2, ...request, revision_id: 'R3', runtime_apply_status: 'aligned',
  }));
  mount();
  await waitFor(() => expect(screen.getByTestId('range-summary')).toHaveTextContent('D0–D1'));
  fireEvent.change(screen.getByTestId('rule-prefix-input'), { target: { value: 'FIRST_' } });
  await waitFor(() => expect(sourceRuleCandidateAPI.getTags).toHaveBeenCalledWith('R', expect.objectContaining({ revision_id: 'R3' })));
  await waitFor(() => expect(screen.getByTestId('rule-mc-persisted-status')).toHaveAttribute('data-status', 'unavailable'));
  fireEvent.change(screen.getByTestId('rule-prefix-input'), { target: { value: 'BLOCKED_' } });
  expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(1);
  expect(screen.getByTestId('rule-save-state')).toHaveTextContent('draft-invalid:R3');
});

it('does not use a late acknowledgement to renew evidence after leaving and returning to a locator', async () => {
  let finishSave!: () => void;
  vi.mocked(studioV2RulesAPI.update).mockImplementation(async (_id, request) => {
    await new Promise<void>(resolve => { finishSave = resolve; });
    return { ...nativeRule(), start_address: 'D0', count: 2, ...request, revision_id: 'R3', runtime_apply_status: 'aligned' };
  });
  const route = entry(), view = mount(route);
  await waitFor(() => expect(screen.getByTestId('range-summary')).toHaveTextContent('D0–D1'));
  fireEvent.change(screen.getByTestId('rule-prefix-input'), { target: { value: 'FIRST_' } });
  await waitFor(() => expect(finishSave).toBeTypeOf('function'));
  view.navigate(locatorRead(1, 'B'));
  view.navigate(route);
  await act(async () => finishSave());
  expect(screen.getByTestId('rule-save-state')).toHaveTextContent('draft-invalid:R2');
  expect(vi.mocked(sourceRuleCandidateAPI.getTags).mock.calls.some(([, query]) => query?.revision_id === 'R3')).toBe(false);
});

it('does not query new evidence from an acknowledgement with different saved fields', async () => {
  vi.mocked(studioV2RulesAPI.update).mockImplementation(async (_id, request) => ({
    ...nativeRule(), start_address: 'D0', count: 2, ...request, naming_prefix: 'OTHER_', revision_id: 'R3', runtime_apply_status: 'aligned',
  }));
  mount();
  await waitFor(() => expect(screen.getByTestId('range-summary')).toHaveTextContent('D0–D1'));
  fireEvent.change(screen.getByTestId('rule-prefix-input'), { target: { value: 'FIRST_' } });
  await waitFor(() => expect(screen.getByTestId('rule-save-state')).toHaveTextContent('saved:R3'));
  expect(vi.mocked(sourceRuleCandidateAPI.getTags).mock.calls.some(([, query]) => query?.revision_id === 'R3')).toBe(false);
  expect(screen.getByTestId('rule-mc-persisted-status')).toHaveAttribute('data-status', 'unavailable');
});
