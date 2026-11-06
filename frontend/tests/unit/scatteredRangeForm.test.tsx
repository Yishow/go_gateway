import { describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen } from '@testing-library/react';
import '../../src/i18n/config';
import type { Device, Rule } from '../../src/features/datalink/workbench-v2/state/types';
import {
  calculateScatteredRandomUsage,
  resolveScatteredRangeSupport,
  validateScatteredItems,
} from '../../src/features/datalink/workbench-v2/state/scatteredRange';
import {
  isStudioV2RuleValid,
  toStudioV2RuleCreateRequest,
  toStudioV2RuleUpdateRequest,
} from '../../src/features/datalink/workbench-v2/state/studioV2RuleAutosave';
import {
  newRecordingRangeBuffer,
  parseRecordingRangeBuffer,
} from '../../src/features/datalink/workbench-v2/recording-editor/recordingRangeBuffer';
import { RuleEditor } from '../../src/features/datalink/workbench-v2/steps/step2/RuleEditor';

function makeDevice(
  protocol: Device['protocol'],
  config: Record<string, unknown> = {},
): Device {
  return {
    id: `dev-${protocol}`,
    name: `Device ${protocol}`,
    description: '',
    protocol,
    config,
    status: 'active',
    test: null,
    persisted: true,
    save_state: 'saved',
    save_error: null,
  };
}

function makeHomogeneousRule(deviceId: string): Rule {
  return {
    id: 'rule-1',
    device_id: deviceId,
    name: 'MIX_RANGE',
    start_address: 'D0',
    count: 1,
    data_type: 'int16',
    naming_prefix: 'MIX_',
    enabled: true,
    scale_multiplier: 1,
    scale_offset: 0,
    data_format: '',
    skipped_addresses: [],
    share_enabled: false,
    share_start_register: null,
    share_stride: null,
    persisted: true,
    save_state: 'saved',
    save_error: null,
  };
}

describe('scatteredRange pure functions & random usage calculation', () => {
  it('resolves random support and effective limit from device acquisition profile', () => {
    const modbus = makeDevice('modbus_tcp', { host: '127.0.0.1', port: 502 });
    expect(resolveScatteredRangeSupport(modbus)).toEqual({
      supported: false,
      effectiveLimit: 0,
      profile: 'undeclared',
      profileLimit: 0,
      disabledReason: 'protocol_not_supported',
    });

    const mcUndeclared = makeDevice('mc_3e', { host: '127.0.0.1', port: 5000 });
    expect(resolveScatteredRangeSupport(mcUndeclared)).toEqual({
      supported: false,
      effectiveLimit: 0,
      profile: 'undeclared',
      profileLimit: 0,
      disabledReason: 'profile_not_enabled',
    });

    const mcSafe = makeDevice('mc_3e', {
      acquisition: { version: 1, profile: 'safe-v1' },
    });
    expect(resolveScatteredRangeSupport(mcSafe)).toEqual({
      supported: false,
      effectiveLimit: 0,
      profile: 'safe-v1',
      profileLimit: 0,
      disabledReason: 'profile_not_enabled',
    });

    const mcQ = makeDevice('mc_3e', {
      acquisition: { version: 1, profile: 'mitsubishi-q-iqr-3e-binary-v1' },
    });
    expect(resolveScatteredRangeSupport(mcQ)).toEqual({
      supported: true,
      effectiveLimit: 192,
      profile: 'mitsubishi-q-iqr-3e-binary-v1',
      profileLimit: 192,
    });

    const mcCustomCap = makeDevice('mc_3e', {
      acquisition: {
        version: 1,
        profile: 'mitsubishi-fx5-cpu-3e-binary-v1',
        limits: { random_points: 120 },
      },
    });
    expect(resolveScatteredRangeSupport(mcCustomCap)).toEqual({
      supported: true,
      effectiveLimit: 120,
      profile: 'mitsubishi-fx5-cpu-3e-binary-v1',
      profileLimit: 192,
    });

    const fatek = makeDevice('fatek_fbs', {
      acquisition: { version: 1, profile: 'fatek-fbs-ascii-v1' },
    });
    expect(resolveScatteredRangeSupport(fatek)).toEqual({
      supported: true,
      effectiveLimit: 64,
      profile: 'fatek-fbs-ascii-v1',
      profileLimit: 64,
    });
  });

  it('calculates MC random usage: 16-bit and 32-bit items count 1, bit items count 1 per distinct 16-point unit', () => {
    const usage = calculateScatteredRandomUsage(
      'mc_3e',
      [
        { address: 'D39', data_type: 'int16' },
        { address: 'D100', data_type: 'int32' },
        { address: 'D200', data_type: 'float32' },
        { address: 'M0', data_type: 'bool' },
        { address: 'M5', data_type: 'bool' },
        { address: 'M15', data_type: 'bool' }, // Same 16-point unit M:0 -> 1 point total for M0/M5/M15
        { address: 'M16', data_type: 'bool' }, // Next 16-point unit M:1 -> +1 point
        { address: 'X0', data_type: 'bool' },  // Distinct area X:0 -> +1 point
      ],
      192,
    );

    expect(usage.usedPoints).toBe(6);
    expect(usage.effectiveLimit).toBe(192);
    expect(usage.isWarning).toBe(false);
    expect(usage.isOverLimit).toBe(false);
    expect(usage.status).toBe('normal');
  });

  it('calculates FATEK random usage: 16-bit counts 1 and 32-bit counts 2', () => {
    const usage = calculateScatteredRandomUsage(
      'fatek_fbs',
      [
        { address: 'R0', data_type: 'int16' },
        { address: 'D10', data_type: 'int32' },
        { address: 'D20', data_type: 'float32' },
        { address: 'M0', data_type: 'bool' },
      ],
      64,
    );

    expect(usage.usedPoints).toBe(1 + 2 + 2 + 1);
    expect(usage.effectiveLimit).toBe(64);
    expect(usage.status).toBe('normal');
  });

  it('triggers warning at 173 / 192 (>= 90%) and over-limit notice above 192 without blocking validity', () => {
    const items172 = Array.from({ length: 172 }, (_, idx) => ({
      address: `D${idx * 2}`,
      data_type: 'int16' as const,
    }));
    const usage172 = calculateScatteredRandomUsage('mc_3e', items172, 192);
    expect(usage172.usedPoints).toBe(172);
    expect(usage172.isWarning).toBe(false);
    expect(usage172.status).toBe('normal');

    const items173 = [
      ...items172,
      { address: 'D500', data_type: 'int32' as const },
    ];
    const usage173 = calculateScatteredRandomUsage('mc_3e', items173, 192);
    expect(usage173.usedPoints).toBe(173);
    expect(usage173.isWarning).toBe(true);
    expect(usage173.isOverLimit).toBe(false);
    expect(usage173.status).toBe('warning');

    const items195 = Array.from({ length: 195 }, (_, idx) => ({
      item_id: `item-${idx + 1}`,
      name: '',
      address: `D${idx * 2}`,
      register_offset: 0,
      data_type: 'int16' as const,
      semantic_kind: 'gauge' as const,
    }));
    const usage195 = calculateScatteredRandomUsage('mc_3e', items195, 192);
    expect(usage195.usedPoints).toBe(195);
    expect(usage195.isWarning).toBe(true);
    expect(usage195.isOverLimit).toBe(true);
    expect(usage195.status).toBe('over_limit');

    const overLimitRule: Rule = {
      ...makeHomogeneousRule('dev-mc_3e'),
      layout_mode: 'mixed',
      mixed_items: items195,
    };
    expect(isStudioV2RuleValid(overLimitRule)).toBe(true);
  });

  it('detects overlapping spans in scattered items (e.g., D3 uint32 overlaps D4 int16)', () => {
    const check = validateScatteredItems('mc_3e', [
      {
        item_id: 'd3',
        name: '',
        address: 'D3',
        register_offset: 0,
        data_type: 'uint32',
        semantic_kind: 'gauge',
      },
      {
        item_id: 'd4',
        name: '',
        address: 'D4',
        register_offset: 0,
        data_type: 'int16',
        semantic_kind: 'gauge',
      },
    ]);
    expect(check.valid).toBe(false);
    expect(check.overlapPairs).toEqual([['d3', 'd4']]);
  });
});

describe('RuleEditor scattered address mode UI and payload', () => {
  it('disables scattered mode and shows reason when device does not support random', () => {
    const modbusDevice = makeDevice('modbus_tcp', { host: '127.0.0.1', port: 502 });
    const rule = makeHomogeneousRule(modbusDevice.id);
    const dispatch = vi.fn();

    render(
      <RuleEditor
        rule={rule}
        devices={[modbusDevice]}
        globalShareEnabled={false}
        dispatch={dispatch}
      />,
    );

    const scatteredBtn = screen.getByTestId('range-mode-mixed');
    expect(scatteredBtn).toBeDisabled();
    expect(screen.getByTestId('scattered-mode-disabled-reason')).toBeInTheDocument();
  });

  it('renders 173 / 192 usage warning and over-limit split message in RuleEditor', () => {
    const mcDevice = makeDevice('mc_3e', {
      acquisition: { version: 1, profile: 'mitsubishi-q-iqr-3e-binary-v1' },
    });
    const items173 = Array.from({ length: 173 }, (_, idx) => ({
      item_id: `item-${idx + 1}`,
      name: '',
      address: `D${idx * 2}`,
      register_offset: 0,
      data_type: 'int16' as const,
      semantic_kind: 'gauge' as const,
    }));
    const rule173: Rule = {
      ...makeHomogeneousRule(mcDevice.id),
      layout_mode: 'mixed',
      mixed_items: items173,
    };

    const { rerender } = render(
      <RuleEditor
        rule={rule173}
        devices={[mcDevice]}
        globalShareEnabled={false}
        dispatch={vi.fn()}
      />,
    );

    expect(screen.getByTestId('scattered-usage-counter')).toHaveTextContent('173 / 192');
    expect(screen.getByTestId('scattered-usage-warning')).toBeInTheDocument();

    const items193 = [
      ...items173,
      ...Array.from({ length: 20 }, (_, idx) => ({
        item_id: `item-extra-${idx + 1}`,
        name: '',
        address: `D${400 + idx * 2}`,
        register_offset: 0,
        data_type: 'int16' as const,
        semantic_kind: 'gauge' as const,
      })),
    ];
    rerender(
      <RuleEditor
        rule={{ ...rule173, mixed_items: items193 }}
        devices={[mcDevice]}
        globalShareEnabled={false}
        dispatch={vi.fn()}
      />,
    );

    expect(screen.getByTestId('scattered-usage-counter')).toHaveTextContent('193 / 192');
    expect(screen.getByTestId('scattered-usage-over-limit')).toBeInTheDocument();
  });

  it('switches to scattered mode, edits cross-area items, and builds mixed_items save payload', () => {
    const mcDevice = makeDevice('mc_3e', {
      acquisition: { version: 1, profile: 'mitsubishi-q-iqr-3e-binary-v1' },
    });
    let currentRule: Rule = makeHomogeneousRule(mcDevice.id);
    const dispatch = vi.fn((action) => {
      if (action.type === 'updateRule') {
        currentRule = { ...currentRule, ...action.patch };
      }
    });

    const { rerender } = render(
      <RuleEditor
        rule={currentRule}
        devices={[mcDevice]}
        globalShareEnabled={false}
        dispatch={dispatch}
      />,
    );

    fireEvent.click(screen.getByTestId('range-mode-mixed'));
    expect(currentRule.layout_mode).toBe('mixed');
    expect(currentRule.mixed_items?.length).toBeGreaterThanOrEqual(1);

    rerender(
      <RuleEditor
        rule={currentRule}
        devices={[mcDevice]}
        globalShareEnabled={false}
        dispatch={dispatch}
      />,
    );

    // Edit first item to D39 int16
    fireEvent.change(screen.getByTestId('scattered-item-address-0'), {
      target: { value: 'D39' },
    });
    rerender(
      <RuleEditor
        rule={currentRule}
        devices={[mcDevice]}
        globalShareEnabled={false}
        dispatch={dispatch}
      />,
    );

    // Add second item (D100 int32)
    fireEvent.click(screen.getByTestId('scattered-add-item'));
    rerender(
      <RuleEditor
        rule={currentRule}
        devices={[mcDevice]}
        globalShareEnabled={false}
        dispatch={dispatch}
      />,
    );
    fireEvent.change(screen.getByTestId('scattered-item-address-1'), {
      target: { value: 'D100' },
    });
    rerender(
      <RuleEditor
        rule={currentRule}
        devices={[mcDevice]}
        globalShareEnabled={false}
        dispatch={dispatch}
      />,
    );
    fireEvent.change(screen.getByTestId('scattered-item-datatype-1'), {
      target: { value: 'int32' },
    });
    rerender(
      <RuleEditor
        rule={currentRule}
        devices={[mcDevice]}
        globalShareEnabled={false}
        dispatch={dispatch}
      />,
    );

    // Add third item (M5 bool)
    fireEvent.click(screen.getByTestId('scattered-add-item'));
    rerender(
      <RuleEditor
        rule={currentRule}
        devices={[mcDevice]}
        globalShareEnabled={false}
        dispatch={dispatch}
      />,
    );
    fireEvent.change(screen.getByTestId('scattered-item-address-2'), {
      target: { value: 'M5' },
    });
    rerender(
      <RuleEditor
        rule={currentRule}
        devices={[mcDevice]}
        globalShareEnabled={false}
        dispatch={dispatch}
      />,
    );
    fireEvent.change(screen.getByTestId('scattered-item-datatype-2'), {
      target: { value: 'bool' },
    });
    rerender(
      <RuleEditor
        rule={currentRule}
        devices={[mcDevice]}
        globalShareEnabled={false}
        dispatch={dispatch}
      />,
    );

    const createPayload = toStudioV2RuleCreateRequest(currentRule);
    expect(createPayload.layout_mode).toBe('mixed');
    expect(createPayload.scale_multiplier).toBeUndefined();
    expect(createPayload.scale_offset).toBeUndefined();
    expect(createPayload.mixed_items).toEqual([
      expect.objectContaining({ address: 'D39', data_type: 'int16', register_offset: 0 }),
      expect.objectContaining({ address: 'D100', data_type: 'int32', register_offset: 0 }),
      expect.objectContaining({ address: 'M5', data_type: 'bool', register_offset: 0 }),
    ]);

    const updatePayload = toStudioV2RuleUpdateRequest(currentRule);
    expect(updatePayload.layout_mode).toBe('mixed');
    expect(updatePayload.scale_multiplier).toBeUndefined();
    expect(updatePayload.scale_offset).toBeUndefined();
    expect(updatePayload.mixed_items).toHaveLength(3);
  });

  it('parses home recording range buffer in scattered mode into a mixed rule', () => {
    const buffer = {
      ...newRecordingRangeBuffer(),
      name: 'SCATTER_SRC',
      layoutMode: 'mixed' as const,
      mixedItems: [
        {
          item_id: 'item-1',
          name: 'Speed',
          address: 'd39',
          register_offset: 0,
          data_type: 'int16' as const,
          semantic_kind: 'gauge' as const,
        },
        {
          item_id: 'item-2',
          name: '',
          address: 'm5',
          register_offset: 0,
          data_type: 'bool' as const,
          semantic_kind: 'state' as const,
        },
      ],
    };
    const parsed = parseRecordingRangeBuffer(buffer, 'mc-dev', 'mc_3e');
    expect(parsed).toBeDefined();
    expect(parsed?.rule.layout_mode).toBe('mixed');
    expect(parsed?.rule.mixed_items).toEqual([
      expect.objectContaining({ item_id: 'item-1', address: 'D39', data_type: 'int16' }),
      expect.objectContaining({ item_id: 'item-2', address: 'M5', data_type: 'bool' }),
    ]);
  });

  it('summarises mixed span by data width and ends at the last word of the last item', () => {
    const buffer = {
      ...newRecordingRangeBuffer(),
      name: 'MIX',
      layoutMode: 'mixed' as const,
      mixedItems: [
        { item_id: 'item-1', name: '', address: 'D100', register_offset: 0, data_type: 'int16' as const, semantic_kind: 'gauge' as const },
        { item_id: 'item-2', name: '', address: 'D200', register_offset: 0, data_type: 'int32' as const, semantic_kind: 'gauge' as const },
        { item_id: 'item-3', name: '', address: 'D300', register_offset: 0, data_type: 'uint16' as const, semantic_kind: 'gauge' as const },
        { item_id: 'item-4', name: '', address: 'D350', register_offset: 0, data_type: 'float32' as const, semantic_kind: 'gauge' as const },
      ],
    };
    const parsed = parseRecordingRangeBuffer(buffer, 'mc-dev', 'mc_3e');
    expect(parsed).toMatchObject({ span: 6, end: 'D351' });
  });
});
