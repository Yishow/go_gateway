import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { RuleEditor } from '../../../src/features/datalink/workbench-v2/steps/step2/RuleEditor';
import { RecordingHomeSavedRangeEditor } from '../../../src/features/datalink/workbench-v2/recording-home/RecordingHomeSavedRangeEditor';
import { buildPersistedMCRuleState, type PersistedMCRuleState } from '../../../src/features/datalink/workbench-v2/state/mcPersistedPoints';
import { hydrateStudioV2Rule } from '../../../src/features/datalink/workbench-v2/state/studioV2RuleAutosave';
import { hydrateStudioV2Device } from '../../../src/features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { useStudioV2PersistedMCPoints } from '../../../src/pages/datalink/workbench-v2/useStudioV2PersistedMCPoints';
import type { SourceRuleTagSnapshotView } from '../../../src/types/sourceRuleCandidates';
import { nativeDevice, nativeRule } from './studio-locator-fixtures';

const api = vi.hoisted(() => ({ full: vi.fn(), tags: vi.fn() }));
vi.mock('../../../src/services/sourceRuleCandidates', () => ({ sourceRuleCandidateAPI: { get: api.full, getTags: api.tags } }));

const rule = { ...nativeRule(), start_address: 'D0', count: 32, naming_prefix: 'TAG_' };
const device = { ...nativeDevice(), protocol: 'mc_3e' as const };
const scope = { workspace_id: 'W', workspace_revision: 'SOURCE-REV' };
const snapshot: SourceRuleTagSnapshotView = {
  source_rule_id: 'R', workspace_id: 'W', workspace_revision: 'SOURCE-REV', revision_id: 'R2',
  tags: { status: 'ready', candidates: Array.from({ length: 32 }, (_, index) => ({
    id: `C${index}`, identity: { source_rule_id: 'R', candidate_type: 'tags', candidate_kind: 'tag', derived_from_rule_address: 'D0' },
    proposed_signature: `sig${index}`, address: `D${index}`, point_id: `P${index}`, tag_key: `tag${index}`,
    display_name: `TAG_${index}`, data_type: 'int16', status: 'draft',
  })) },
};

it('shows the saved MC addresses when tag evidence is available even if database output validation fails', async () => {
  api.full.mockRejectedValue(new Error('database output unavailable'));
  api.tags.mockResolvedValue(snapshot);
  function View() {
    const result = useStudioV2PersistedMCPoints(true, [rule], [device], scope);
    return <RuleEditor rule={{ ...hydrateStudioV2Rule(rule), mc_persisted_points: result.evidenceByRule.R }}
      devices={[hydrateStudioV2Device(device)]} globalShareEnabled={false} dispatch={vi.fn()} simplePresentation />;
  }
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  render(<QueryClientProvider client={client}><View /></QueryClientProvider>);
  await waitFor(() => expect(screen.queryByTestId('rule-mc-persisted-status')).not.toBeInTheDocument());
  expect(screen.getByTestId('range-summary')).toHaveTextContent('D0');
  expect(screen.getByTestId('range-summary')).toHaveTextContent('D31');
  expect(api.tags).toHaveBeenCalledWith('R', { workspace_id: 'W', expected_workspace_revision: 'SOURCE-REV', revision_id: 'R2' });
  expect(api.full).not.toHaveBeenCalled();
});

it('updates an already-open MC range when matching evidence arrives without losing its draft', () => {
  const local = hydrateStudioV2Rule(rule), localDevice = hydrateStudioV2Device(device);
  const unavailable: PersistedMCRuleState = { status: 'unavailable', scope: { ...scope, revision_id: 'R2' }, points: [] };
  const compatible = buildPersistedMCRuleState(rule, { ...scope, revision_id: 'R2' }, snapshot);
  const leave = { current: undefined }, observer = { current: undefined }, dispatch = vi.fn();
  const renderEditor = (evidence: PersistedMCRuleState) => {
    const current = { ...local, mc_persisted_points: evidence };
    return <RecordingHomeSavedRangeEditor target={{ workspaceId: 'W', deviceId: 'D', ruleId: 'R', kind: 'points', epoch: 1 }}
      state={{ ...INITIAL_STATE, devices: [localDevice], rules: [current] }} rule={current} dispatch={dispatch} leave={leave} suspended={false} list={null}
      onClose={vi.fn()} headerSlot={null} onLocked={vi.fn()} onLeave={vi.fn()} ruleObserverRef={observer} onRetryShare={vi.fn()} onOpenShareSettings={vi.fn()} />;
  };
  const view = render(renderEditor(unavailable));
  expect(screen.getByTestId('rule-mc-persisted-status')).toHaveAttribute('data-status', 'unavailable');
  fireEvent.change(screen.getByTestId('rule-prefix-input'), { target: { value: 'UNSENT_' } });
  expect(screen.getByTestId('home-range-save')).toBeDisabled();
  view.rerender(renderEditor(compatible));
  expect(screen.queryByTestId('rule-mc-persisted-status')).not.toBeInTheDocument();
  expect(screen.getByTestId('range-summary')).toHaveTextContent('D31');
  expect(screen.getByTestId('rule-prefix-input')).toHaveValue('UNSENT_');
  expect(screen.getByTestId('home-range-save')).toBeEnabled();
  expect(dispatch).not.toHaveBeenCalled();
});

it.each(['conflict', 'foreign-workspace', 'old-rule-revision'] as const)('keeps saving blocked for %s evidence', kind => {
  const local = hydrateStudioV2Rule(rule), localDevice = hydrateStudioV2Device(device);
  const evidence = buildPersistedMCRuleState(rule, { ...scope, revision_id: 'R2' }, snapshot);
  if (kind === 'conflict') evidence.status = 'conflict';
  if (kind === 'foreign-workspace') evidence.scope = { ...scope, workspace_id: 'FOREIGN', revision_id: 'R2' };
  if (kind === 'old-rule-revision') evidence.scope = { ...scope, revision_id: 'OLD' };
  const current = { ...local, mc_persisted_points: evidence }, dispatch = vi.fn();
  render(<RecordingHomeSavedRangeEditor target={{ workspaceId: 'W', deviceId: 'D', ruleId: 'R', kind: 'points', epoch: 1 }}
    state={{ ...INITIAL_STATE, devices: [localDevice], rules: [current] }} rule={current} dispatch={dispatch} leave={{ current: undefined }} suspended={false} list={null}
    onClose={vi.fn()} headerSlot={null} onLocked={vi.fn()} onLeave={vi.fn()} ruleObserverRef={{ current: undefined }} onRetryShare={vi.fn()} onOpenShareSettings={vi.fn()} />);
  fireEvent.change(screen.getByTestId('rule-prefix-input'), { target: { value: 'UNSENT_' } });
  expect(screen.getByTestId('home-range-save')).toBeDisabled();
  expect(screen.getByTestId('rule-mc-persisted-status')).toBeVisible();
  expect(dispatch).not.toHaveBeenCalled();
});
