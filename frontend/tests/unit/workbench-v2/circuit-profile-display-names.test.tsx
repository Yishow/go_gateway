import { fireEvent, render, screen, within } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { CircuitProfileConfigFields, type CircuitProfileConfigFieldsProps } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileConfigFields';
import { candidateKey, defaultCircuitProfileLocalDraft } from '@/features/datalink/workbench-v2/state/circuitProfileBasic';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { savedGroupState } from '../../fixtures/writeGroupState';
import { expandSources } from '../../fixtures/sourceReview';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string) => key }),
}));

function namedState() {
  const state = savedGroupState();
  state.mappings[state.points[0].id].tag_key = 'dev.1780862010602.842.tag.0.rd0';
  state.mappings[state.points[0].id].display_name = '主迴路有功功率';
  return state;
}

function configProps(): CircuitProfileConfigFieldsProps {
  const candidates = buildGroupCandidates(namedState()).candidates;
  return {
    workspaceId: 'ws-1', workspaceRevision: 'wrev-1', readonly: false,
    savedDevices: [{ id: 'dev-1', name: 'PLC A' }], activeDeviceId: 'dev-1',
    onSelectDevice: vi.fn(), deviceCandidates: candidates,
    localDraft: { ...defaultCircuitProfileLocalDraft(candidates, []), cumulative_keys: [candidateKey(candidates[0])] }, updateLocalDraft: vi.fn(),
    instParsed: { valid: true }, cumParsed: { valid: true },
    isPostgres: true, storageLocked: false, effectiveStorageMode: 'plain',
    cumulativeCircuitMissing: false, circuitId: '',
    unconfirmedCounterStatuses: [], onCounterConfirmed: vi.fn(),
  };
}

describe('Circuit profile source display names', () => {
  it('prefers the mapping display name without changing technical keys or persisted identities', () => {
    const state = namedState();
    const original = buildGroupCandidates(state).candidates[0];
    expect(original.label).toBe('主迴路有功功率');
    expect(original.tag_key).toBe('dev.1780862010602.842.tag.0.rd0');
    expect(original).toMatchObject({ device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0' });
    state.mappings[state.points[0].id].display_name = '進線功率';
    const renamed = buildGroupCandidates(state).candidates[0];
    expect(renamed.label).toBe('進線功率');
    expect(renamed.key).toBe(original.key);
    expect(candidateKey(renamed)).toBe(candidateKey(original));
  });

  it.each(['', '   '])('falls back to the source address for an empty display name (%j)', (displayName) => {
    const state = namedState();
    state.mappings[state.points[0].id].display_name = displayName;
    expect(buildGroupCandidates(state).candidates[0].label).toBe('40001');
  });

  it('renders the display name and address with matching accessible labels for each profile', () => {
    const props = configProps();
    // Exercise the renderer separately from the candidate builder.
    props.deviceCandidates[0].label = '主迴路有功功率';
    render(<CircuitProfileConfigFields {...props} />);
    expandSources();
    const sources = within(screen.getByTestId('circuit-profile-member-list'));
    expect(sources.getByText('主迴路有功功率')).toBeVisible();
    expect(sources.queryByText('dev.1780862010602.842.tag.0.rd0')).not.toBeInTheDocument();
    expect(sources.getByText('(40001 · float64)')).toBeVisible();
    expect(sources.getByRole('checkbox', { name: '主迴路有功功率' })).toBeChecked();
    const inst = sources.getByRole('checkbox', { name: '主迴路有功功率 step4.basic.circuit.instantaneous_label' });
    expect(sources.getByRole('checkbox', { name: '主迴路有功功率 step4.basic.circuit.cumulative_label' })).toBeChecked();
    fireEvent.click(inst);
    const updater = vi.mocked(props.updateLocalDraft).mock.calls[0][0];
    const next = updater(props.localDraft);
    expect(next.instantaneous_keys).not.toContain(candidateKey(props.deviceCandidates[0]));
    expect(next.cumulative_keys).toContain(candidateKey(props.deviceCandidates[0]));
  });
});
