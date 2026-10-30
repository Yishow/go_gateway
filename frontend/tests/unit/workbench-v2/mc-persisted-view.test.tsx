import { render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it } from 'vitest';
import { Step2Rule } from '../../../src/features/datalink/workbench-v2/steps/step2/Step2Rule';
import {
  derivePoints,
  isStep2Ready,
} from '../../../src/features/datalink/workbench-v2/state/sourceRule';
import {
  mergePersistedMCState,
  type PersistedMCRuleState,
} from '../../../src/features/datalink/workbench-v2/state/mcPersistedPoints';
import type { Point, Rule } from '../../../src/features/datalink/workbench-v2/state/types';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { ProtocolType } from '../../../src/types/datalink';

afterEach(() => {
  document.body.innerHTML = '';
});

const legacyAddresses = [
  'W0', 'W1', 'W2', 'W3', 'W4', 'W5', 'W6', 'W7', 'W8', 'W9',
  'W10', 'W11', 'W12', 'W13', 'W14', 'W15', 'W16',
];

const scope = {
  workspace_id: 'workspace-1',
  workspace_revision: 'workspace-rev-1',
  revision_id: 'rule-rev-1',
};

const identities = legacyAddresses.map((address, index) => ({
  point_id: `legacy-point-${index}`,
  address,
  tag_id: `legacy-tag-${index}`,
  mapping_id: `legacy-mapping-${index}`,
}));

function persistedEvidence(status: PersistedMCRuleState['status'], points = identities): PersistedMCRuleState {
  return { status, scope, points, ...(status === 'conflict' ? { reason: 'persisted_address_set_conflict' } : {}) };
}

function mcRule(evidence: PersistedMCRuleState = persistedEvidence('conflict'), enabled = true): Rule {
  return {
    id: 'rule-mc',
    device_id: 'dev-mc',
    workspace_id: scope.workspace_id,
    revision_id: scope.revision_id,
    name: 'Legacy MC',
    start_address: 'W0',
    count: 17,
    data_type: 'int16',
    naming_prefix: 'LEGACY_',
    enabled,
    scale_multiplier: 1,
    scale_offset: 0,
    data_format: '',
    skipped_addresses: [],
    share_enabled: false,
    share_start_register: null,
    share_stride: null,
    persisted: true,
    save_state: 'saved',
    mc_persisted_points: evidence,
  };
}

function persistedPoint(address: string, index: number): Point {
  return {
    id: `legacy-point-${index}`,
    device_id: 'dev-mc',
    rule_id: 'rule-mc',
    rule_name: 'Legacy MC',
    name: `LEGACY_${address}`,
    address,
    data_type: 'int16',
    function: 'W (Word)',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 1,
    _rule_offset: 0,
  };
}

describe('persisted MC Step 2 view', () => {
  it('uses persisted MC point identities instead of the corrected synthetic plan', () => {
    const points = derivePoints(mcRule(), 'dev-mc', new Set(), 'mc_3e');

    expect(points.map((point) => point.id)).toEqual(identities.map((point) => point.point_id));
    expect(points.map((point) => point.address)).toEqual(legacyAddresses);
  });

  it('retains last-known persisted points while a later scope is checking or unavailable', () => {
    const state = {
      ...INITIAL_STATE,
      devices: [{ ...INITIAL_STATE.devices[0], id: 'dev-mc', protocol: 'mc_3e' as ProtocolType }],
      rules: [mcRule()],
      points: legacyAddresses.map(persistedPoint),
    };
    const checking = persistedEvidence('checking', []);
    const first = mergePersistedMCState(state, { 'rule-mc': checking }, [], []);
    const unavailable = { ...persistedEvidence('unavailable', []), scope: { ...scope, workspace_revision: 'workspace-rev-2' } };
    const second = mergePersistedMCState({ ...state, rules: first.rules, points: first.points }, { 'rule-mc': unavailable }, [], []);

    expect(first.rules[0].mc_persisted_points?.points.map((point) => point.address)).toEqual(legacyAddresses);
    expect(first.points.map((point) => point.address)).toEqual(legacyAddresses);
    expect(second.rules[0].mc_persisted_points?.points.map((point) => point.address)).toEqual(legacyAddresses);
    expect(second.points.map((point) => point.address)).toEqual(legacyAddresses);
  });

  it('does not mark a persisted MC conflict as Step 2 ready', () => {
    expect(isStep2Ready([mcRule()], { 'dev-mc': 'mc_3e' })).toBe(false);
  });

  it('renders persisted legacy point IDs and avoids a corrected continuous range', () => {
    render(
      <Step2Rule
        state={{
          ...INITIAL_STATE,
          current: 2,
          devices: [{ ...INITIAL_STATE.devices[0], id: 'dev-mc', protocol: 'mc_3e' as ProtocolType }],
          rules: [mcRule()],
          selectedRuleId: 'rule-mc',
        }}
        dispatch={() => undefined}
        onContinue={() => undefined}
      />,
    );

    expect(screen.getByTestId('point-cell-legacy-point-16')).toBeInTheDocument();
    expect(screen.queryByTestId('point-cell-rule-mc-p-16')).not.toBeInTheDocument();
    expect(screen.getByTestId('range-summary')).not.toHaveTextContent('W0 ~ W10');
    expect(screen.getByTestId('continue-step3-btn')).toBeDisabled();
    expect(screen.getByTestId('readiness-issues-warning')).toHaveTextContent('Legacy MC');
    expect(screen.getByTestId('merged-point-table-container')).not.toHaveTextContent('衝突 0');
    expect(screen.queryByText('Derivation 完畢，點位將作為 Step 3 映射配置的來源。')).not.toBeInTheDocument();
  });

  it('keeps unavailable persisted points visible without claiming a clean result', () => {
    render(
      <Step2Rule
        state={{
          ...INITIAL_STATE,
          current: 2,
          devices: [{ ...INITIAL_STATE.devices[0], id: 'dev-mc', protocol: 'mc_3e' as ProtocolType }],
          rules: [mcRule(persistedEvidence('unavailable'))],
          selectedRuleId: 'rule-mc',
        }}
        dispatch={() => undefined}
        onContinue={() => undefined}
      />,
    );

    expect(screen.getByTestId('point-cell-legacy-point-16')).toBeInTheDocument();
    expect(screen.getByTestId('continue-step3-btn')).toBeDisabled();
    expect(screen.getByTestId('merged-point-table-container')).not.toHaveTextContent('衝突 0');
    expect(screen.queryByText('Derivation 完畢，點位將作為 Step 3 映射配置的來源。')).not.toBeInTheDocument();
  });

  it('marks disabled persisted MC evidence as blocked in the point summary', () => {
    render(
      <Step2Rule
        state={{
          ...INITIAL_STATE,
          current: 2,
          devices: [{ ...INITIAL_STATE.devices[0], id: 'dev-mc', protocol: 'mc_3e' as ProtocolType }],
          rules: [mcRule(persistedEvidence('conflict'), false)],
          selectedRuleId: 'rule-mc',
        }}
        dispatch={() => undefined}
        onContinue={() => undefined}
      />,
    );

    expect(screen.getByTestId('point-cell-legacy-point-16')).toBeInTheDocument();
    expect(screen.getByTestId('point-table-summary')).not.toHaveTextContent('已啟用 0 · 已略過 17 · 衝突 0');
    expect(screen.getByTestId('point-table-blocked-footer')).toBeInTheDocument();
    expect(screen.getByTestId('continue-step3-btn')).toBeDisabled();
  });
});
