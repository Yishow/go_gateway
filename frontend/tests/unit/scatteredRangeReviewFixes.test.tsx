import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import '../../src/i18n/config';
import type { Device, Rule } from '../../src/features/datalink/workbench-v2/state/types';
import { calculateScatteredRandomUsage } from '../../src/features/datalink/workbench-v2/state/scatteredRange';
import { RuleEditor } from '../../src/features/datalink/workbench-v2/steps/step2/RuleEditor';
import { RecordingRangeContent } from '../../src/features/datalink/workbench-v2/recording-editor/RecordingRangeContent';
import { newRecordingRangeBuffer, parseRecordingRangeBuffer } from '../../src/features/datalink/workbench-v2/recording-editor/recordingRangeBuffer';
import type { RecordingRangeController } from '../../src/features/datalink/workbench-v2/recording-editor/useRecordingRangeBuffer';

function makeMCDevice(): Device {
  return {
    id: 'dev-mc',
    name: 'MC PLC',
    description: '',
    protocol: 'mc_3e',
    config: { acquisition: { version: 1, profile: 'mitsubishi-q-iqr-3e-binary-v1' } },
    status: 'active',
    test: null,
    persisted: true,
    save_state: 'saved',
    save_error: null,
  };
}

function makeBaseRule(deviceId: string): Rule {
  return {
    id: 'rule-fix',
    device_id: deviceId,
    name: 'MIX_RULE',
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

describe('Review fixes: Warning 2 (MC 64-bit cost & unit deduplication)', () => {
  it('counts 64-bit MC items as 4 random units and deduplicates identical MC units', () => {
    const usage = calculateScatteredRandomUsage(
      'mc_3e',
      [
        { address: 'D10', data_type: 'int64' },  // 4 x 16-bit units: D10, D11, D12, D13
        { address: 'D100', data_type: 'int16' }, // 1 x 16-bit unit: D100
        { address: 'D100', data_type: 'int16' }, // duplicate of D100 16-bit -> deduplicated (+0)
        { address: 'D200', data_type: 'int32' }, // 1 x 32-bit unit: D200
        { address: 'D200', data_type: 'float32' }, // duplicate 32-bit unit at D200 -> deduplicated (+0)
      ],
      192,
    );
    expect(usage.usedPoints).toBe(4 + 1 + 1);
  });
});

describe('Review fixes: Warning 3 & Suggestion 2 (actionable errors & ARIA bindings)', () => {
  it('renders localized actionable alert for mixed_items_error without leaking raw backend text, and alerts on empty items', () => {
    const device = makeMCDevice();
    const corruptRule: Rule = {
      ...makeBaseRule(device.id),
      layout_mode: 'mixed',
      mixed_items: [],
      mixed_items_error: 'raw backend secret exception: invalid stored mixed_items',
    };

    render(
      <RuleEditor
        rule={corruptRule}
        devices={[device]}
        globalShareEnabled={false}
        dispatch={vi.fn()}
      />,
    );

    const corruptAlert = screen.getByTestId('scattered-corrupt-items-error');
    expect(corruptAlert).toHaveAttribute('role', 'alert');
    expect(corruptAlert.textContent).not.toContain('raw backend secret exception');
    expect(corruptAlert.textContent).not.toContain('invalid stored mixed_items');

    const emptyAlert = screen.getByTestId('scattered-empty-items-error');
    expect(emptyAlert).toHaveAttribute('role', 'alert');
  });

  it('associates row validation errors with inputs via aria-invalid and aria-describedby', () => {
    const device = makeMCDevice();
    const ruleWithErrors: Rule = {
      ...makeBaseRule(device.id),
      layout_mode: 'mixed',
      mixed_items: [
        { item_id: 'item-1', name: '', address: 'D3', register_offset: 0, data_type: 'uint32', semantic_kind: 'gauge' },
        { item_id: 'item-2', name: '', address: 'D4', register_offset: 0, data_type: 'int16', semantic_kind: 'gauge' },
        { item_id: 'item-3', name: '', address: 'M5', register_offset: 0, data_type: 'int16', semantic_kind: 'gauge' },
      ],
    };

    render(
      <RuleEditor
        rule={ruleWithErrors}
        devices={[device]}
        globalShareEnabled={false}
        dispatch={vi.fn()}
      />,
    );

    const addr0 = screen.getByTestId('scattered-item-address-0');
    const err0 = screen.getByTestId('scattered-item-error-0');
    expect(addr0).toHaveAttribute('aria-invalid', 'true');
    expect(err0.id).toBeTruthy();
    expect(addr0).toHaveAttribute('aria-describedby', err0.id);

    const type2 = screen.getByTestId('scattered-item-datatype-2');
    const err2 = screen.getByTestId('scattered-item-error-2');
    expect(type2).toHaveAttribute('aria-invalid', 'true');
    expect(type2).toHaveAttribute('aria-describedby', err2.id);
  });
});

describe('Review fixes: Warning 5 (RecordingRangeContent hides contiguous range.span in mixed mode)', () => {
  it('hides contiguous range.span summary when RecordingRangeContent is in scattered mode', () => {
    const mixedBuffer = {
      ...newRecordingRangeBuffer(),
      name: 'MIX_HOME',
      layoutMode: 'mixed' as const,
      mixedItems: [
        { item_id: 'item-1', name: '', address: 'D39', register_offset: 0, data_type: 'int16' as const, semantic_kind: 'gauge' as const },
        { item_id: 'item-2', name: '', address: 'D100', register_offset: 0, data_type: 'int32' as const, semantic_kind: 'gauge' as const },
      ],
    };
    const parsed = parseRecordingRangeBuffer(mixedBuffer, 'dev-mc', 'mc_3e');
    expect(parsed).toMatchObject({ span: 3, end: 'D101' });

    const fakeController: RecordingRangeController = {
      workspaceId: 'ws-1',
      buffer: mixedBuffer,
      dirty: true,
      target: { id: 'dev-mc', name: 'MC PLC' },
      device: {
        id: 'dev-mc',
        name: 'MC PLC',
        description: '',
        protocol: 'mc_3e',
        status: 'active',
        connection_config: JSON.stringify({ acquisition: { version: 1, profile: 'mitsubishi-q-iqr-3e-binary-v1' } }),
        last_test_at: null,
        last_test_success: true,
        last_test_error: '',
        created_at: '',
        updated_at: '',
      },
      outcome: undefined,
      pending: false,
      checking: false,
      preparing: false,
      blocked: false,
      activationAllowed: true,
      locked: false,
      recoverable: false,
      inactiveSources: [],
      canResume: () => false,
      canStart: () => true,
      parsed,
      change: vi.fn(),
      reset: vi.fn(),
      resume: () => false,
      start: vi.fn(),
      submit: vi.fn(),
      inspect: vi.fn(),
      prepareSaved: vi.fn(),
      submission: () => undefined,
      canInspect: false,
      canPrepareSaved: false,
      canSubmit: true,
    };

    render(
      <RecordingRangeContent
        controller={fakeController}
        available={true}
        onBack={vi.fn()}
        onSubmit={vi.fn()}
        onInspect={vi.fn()}
        onPrepare={vi.fn()}
      />,
    );

    expect(screen.queryByText(/D101/)).toBeNull();
  });
});
