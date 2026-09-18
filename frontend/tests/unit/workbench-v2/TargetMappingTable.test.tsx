import { StrictMode } from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { TargetMappingTable } from '../../../src/features/datalink/workbench-v2/steps/step4/TargetMappingTable';
import { getColumnsFor } from '../../../src/features/datalink/workbench-v2/state/dbSchemas';
import type { DbTarget, Mapping, Point } from '../../../src/features/datalink/workbench-v2/state/types';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: string | { defaultValue?: string }) => {
      if (typeof options === 'string') return options;
      return options?.defaultValue ?? key;
    },
  }),
}));

const points: Point[] = [
  {
    id: 'p-1',
    device_id: 'd-1',
    rule_id: 'r-1',
    rule_name: 'Holding Registers',
    name: 'temperature',
    address: '40001',
    data_type: 'int16',
    function: 'holding_register',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 1,
    _rule_offset: 0,
  },
  {
    id: 'p-2',
    device_id: 'd-1',
    rule_id: 'r-1',
    rule_name: 'Holding Registers',
    name: 'pressure',
    address: '40002',
    data_type: 'int16',
    function: 'holding_register',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 1,
    _rule_offset: 0,
  },
];

const mappings: Record<string, Mapping> = {
  'p-1': {
    point_id: 'p-1',
    tag_key: 'line1.temperature',
    display_name: 'Temperature',
    unit: 'C',
    target_type: 'float64',
    scale: 1,
    offset: 0,
    enabled: true,
  },
  'p-2': {
    point_id: 'p-2',
    tag_key: 'line1.pressure',
    display_name: 'Pressure',
    unit: 'kPa',
    target_type: 'float64',
    scale: 1,
    offset: 0,
    enabled: true,
  },
};

const savedTargets: Record<string, DbTarget> = {
  'p-1': {
    tag_id: 'tag.line1.temperature',
    column_name: 'temp_in_c',
    enabled: true,
    save_state: 'saved',
  },
  'p-2': {
    tag_id: 'tag.line1.pressure',
    column_name: 'pressure_main_kpa',
    enabled: true,
    save_state: 'saved',
  },
};

type TableOverrides = Partial<{
  points: Point[];
  mappings: Record<string, Mapping>;
  targets: Record<string, DbTarget>;
  onUpdateTarget: (pointId: string, patch: Partial<DbTarget>) => void;
  disabled: boolean;
}>;

function tableElement(overrides: TableOverrides = {}) {
  return (
    <TargetMappingTable
      points={overrides.points ?? points}
      mappings={overrides.mappings ?? mappings}
      targets={overrides.targets ?? savedTargets}
      columns={getColumnsFor('postgres')}
      onUpdateTarget={overrides.onUpdateTarget ?? vi.fn()}
      disabled={overrides.disabled ?? false}
    />
  );
}

function renderTable(overrides: TableOverrides = {}) {
  return render(tableElement(overrides));
}

function strictTableElement(overrides: TableOverrides = {}) {
  return <StrictMode>{tableElement(overrides)}</StrictMode>;
}

function renderStrictTable(overrides: TableOverrides = {}) {
  return render(strictTableElement(overrides));
}

describe('TargetMappingTable row editing', () => {
  it('PreserveEditingDraft keeps row B text when row A receives a refresh', () => {
    const { rerender } = renderTable();
    const rowB = screen.getByTestId('step4-target-column-p-2');

    fireEvent.change(rowB, { target: { value: 'pressure_draft' } });

    rerender(tableElement({
      targets: {
        ...savedTargets,
        'p-1': { ...savedTargets['p-1'], column_name: 'temp_remote' },
      },
    }));

    expect(screen.getByTestId('step4-target-column-p-1')).toHaveValue('temp_remote');
    expect(screen.getByTestId('step4-target-column-p-2')).toHaveValue('pressure_draft');
  });

  it('SingleEnterCommit submits one logical edit when Enter and blur both fire', () => {
    const onUpdateTarget = vi.fn();
    renderTable({ onUpdateTarget });
    const input = screen.getByTestId('step4-target-column-p-1');

    fireEvent.change(input, { target: { value: 'temperature_output' } });
    fireEvent.keyDown(input, { key: 'Enter' });
    fireEvent.blur(input);
    fireEvent.keyDown(input, { key: 'Enter' });

    expect(onUpdateTarget).toHaveBeenCalledTimes(1);
    expect(onUpdateTarget).toHaveBeenCalledWith('p-1', { column_name: 'temperature_output' });
  });

  it('CompositionDoesNotSubmit while an IME composition is active', () => {
    const onUpdateTarget = vi.fn();
    renderTable({ onUpdateTarget });
    const input = screen.getByTestId('step4-target-column-p-1');

    fireEvent.change(input, { target: { value: '溫度' } });
    fireEvent.compositionStart(input);
    fireEvent.keyDown(input, { key: 'Enter', isComposing: true });

    expect(onUpdateTarget).not.toHaveBeenCalled();
    expect(input).toHaveValue('溫度');
  });

  it('preserves a dirty draft and exposes explicit same-row conflict actions', () => {
    const onUpdateTarget = vi.fn();
    const { rerender } = renderTable({ onUpdateTarget });
    const input = screen.getByTestId('step4-target-column-p-1');

    fireEvent.change(input, { target: { value: 'temperature_draft' } });
    rerender(tableElement({
      targets: {
        ...savedTargets,
        'p-1': { ...savedTargets['p-1'], column_name: 'temperature_remote' },
      },
      onUpdateTarget,
    }));

    expect(screen.getByTestId('step4-target-column-p-1')).toHaveValue('temperature_draft');
    expect(screen.getByTestId('step4-target-remote-conflict-p-1')).toBeInTheDocument();

    fireEvent.keyDown(screen.getByTestId('step4-target-column-p-1'), { key: 'Enter' });
    expect(onUpdateTarget).not.toHaveBeenCalled();
    fireEvent.blur(screen.getByTestId('step4-target-column-p-1'));
    expect(onUpdateTarget).not.toHaveBeenCalled();

    fireEvent.click(screen.getByTestId('step4-target-reload-remote-p-1'));

    expect(screen.getByTestId('step4-target-column-p-1')).toHaveValue('temperature_remote');
    expect(onUpdateTarget).not.toHaveBeenCalled();
  });

  it('allows an operator to resubmit the kept draft after a same-row conflict', () => {
    const onUpdateTarget = vi.fn();
    const { rerender } = renderTable({ onUpdateTarget });
    const input = screen.getByTestId('step4-target-column-p-1');

    fireEvent.change(input, { target: { value: 'temperature_draft' } });
    rerender(tableElement({
      targets: {
        ...savedTargets,
        'p-1': { ...savedTargets['p-1'], column_name: 'temperature_remote' },
      },
      onUpdateTarget,
    }));

    fireEvent.click(screen.getByTestId('step4-target-resubmit-draft-p-1'));
    fireEvent.click(screen.getByTestId('step4-target-resubmit-draft-p-1'));

    expect(onUpdateTarget).toHaveBeenCalledTimes(1);
    expect(onUpdateTarget).toHaveBeenCalledWith('p-1', { column_name: 'temperature_draft' });

    rerender(tableElement({
      targets: {
        ...savedTargets,
        'p-1': { ...savedTargets['p-1'], column_name: 'temperature_draft', save_state: 'saving' },
      },
      onUpdateTarget,
    }));
    expect(screen.getByTestId('step4-target-remote-conflict-p-1')).toBeInTheDocument();

    rerender(tableElement({
      targets: {
        ...savedTargets,
        'p-1': { ...savedTargets['p-1'], column_name: 'temperature_draft', save_state: 'saved' },
      },
      onUpdateTarget,
    }));
    expect(screen.queryByTestId('step4-target-remote-conflict-p-1')).not.toBeInTheDocument();
  });

  it('renders the server save state and error while retaining the input value', () => {
    const onUpdateTarget = vi.fn();
    const { rerender } = renderTable({
      onUpdateTarget,
      targets: {
        ...savedTargets,
        'p-1': { ...savedTargets['p-1'], save_state: 'saving' },
      },
    });

    expect(screen.getByTestId('step4-target-save-state-p-1')).toHaveTextContent('step4.saveStates.saving');

    rerender(tableElement({
      targets: {
        ...savedTargets,
        'p-1': {
          ...savedTargets['p-1'],
          save_state: 'save-error',
          save_error: 'temporary failure',
        },
      },
    }));

    expect(screen.getByTestId('step4-target-save-state-p-1')).toHaveTextContent('step4.saveStates.save-error');
    expect(screen.getByTestId('step4-target-save-error-p-1')).toHaveTextContent('temporary failure');
    expect(screen.getByTestId('step4-target-column-p-1')).toHaveValue('temp_in_c');
  });

  it('allows retrying the same draft after the server reports a save error', () => {
    const onUpdateTarget = vi.fn();
    const { rerender } = renderTable({ onUpdateTarget });
    const input = screen.getByTestId('step4-target-column-p-1');

    fireEvent.change(input, { target: { value: 'temperature_output' } });
    fireEvent.keyDown(input, { key: 'Enter' });
    expect(onUpdateTarget).toHaveBeenCalledTimes(1);

    rerender(tableElement({
      targets: {
        ...savedTargets,
        'p-1': {
          ...savedTargets['p-1'],
          column_name: 'temperature_output',
          save_state: 'save-error',
          save_error: 'temporary failure',
        },
      },
      onUpdateTarget,
    }));

    fireEvent.keyDown(screen.getByTestId('step4-target-column-p-1'), { key: 'Enter' });

    expect(onUpdateTarget).toHaveBeenCalledTimes(2);
    expect(onUpdateTarget).toHaveBeenLastCalledWith('p-1', { column_name: 'temperature_output' });
  });

  it('does not submit on composition blur or after the row becomes disabled', () => {
    const onUpdateTarget = vi.fn();
    const { rerender } = renderTable({ onUpdateTarget });
    const input = screen.getByTestId('step4-target-column-p-1');

    fireEvent.change(input, { target: { value: 'temperature_draft' } });
    fireEvent.compositionStart(input);
    fireEvent.blur(input);
    expect(onUpdateTarget).not.toHaveBeenCalled();
    fireEvent.compositionEnd(input);

    rerender(tableElement({
      targets: {
        ...savedTargets,
        'p-1': { ...savedTargets['p-1'], enabled: false },
      },
      onUpdateTarget,
    }));
    fireEvent.blur(screen.getByTestId('step4-target-column-p-1'));
    expect(onUpdateTarget).not.toHaveBeenCalled();

    rerender(tableElement({ onUpdateTarget, disabled: true }));
    fireEvent.blur(screen.getByTestId('step4-target-column-p-1'));
    expect(onUpdateTarget).not.toHaveBeenCalled();
  });

  it('keeps a newer draft when its earlier save acknowledgement arrives', () => {
    const onUpdateTarget = vi.fn();
    const { rerender } = renderTable({ onUpdateTarget });
    let input = screen.getByTestId('step4-target-column-p-1');

    fireEvent.change(input, { target: { value: 'temperature_a' } });
    fireEvent.keyDown(input, { key: 'Enter' });
    rerender(tableElement({
      targets: {
        ...savedTargets,
        'p-1': { ...savedTargets['p-1'], column_name: 'temperature_a', save_state: 'saving' },
      },
      onUpdateTarget,
    }));

    input = screen.getByTestId('step4-target-column-p-1');
    fireEvent.change(input, { target: { value: 'temperature_b' } });
    rerender(tableElement({
      targets: {
        ...savedTargets,
        'p-1': { ...savedTargets['p-1'], column_name: 'temperature_a', save_state: 'saved' },
      },
      onUpdateTarget,
    }));

    expect(screen.getByTestId('step4-target-column-p-1')).toHaveValue('temperature_b');
    expect(screen.queryByTestId('step4-target-remote-conflict-p-1')).not.toBeInTheDocument();

    fireEvent.blur(screen.getByTestId('step4-target-column-p-1'));
    expect(onUpdateTarget).toHaveBeenCalledTimes(2);
    expect(onUpdateTarget).toHaveBeenLastCalledWith('p-1', { column_name: 'temperature_b' });
  });

  it('accepts an earlier value again after a newer commit is saved', () => {
    const onUpdateTarget = vi.fn();
    const { rerender } = renderTable({ onUpdateTarget });
    let input = screen.getByTestId('step4-target-column-p-1');

    fireEvent.change(input, { target: { value: 'temperature_a' } });
    fireEvent.keyDown(input, { key: 'Enter' });
    input = screen.getByTestId('step4-target-column-p-1');
    fireEvent.change(input, { target: { value: 'temperature_b' } });
    fireEvent.keyDown(input, { key: 'Enter' });
    rerender(tableElement({
      targets: { ...savedTargets, 'p-1': { ...savedTargets['p-1'], column_name: 'temperature_b', save_state: 'saving' } },
      onUpdateTarget,
    }));
    rerender(tableElement({
      targets: { ...savedTargets, 'p-1': { ...savedTargets['p-1'], column_name: 'temperature_b', save_state: 'saved' } },
      onUpdateTarget,
    }));

    input = screen.getByTestId('step4-target-column-p-1');
    fireEvent.change(input, { target: { value: 'temperature_a' } });
    fireEvent.keyDown(input, { key: 'Enter' });
    expect(onUpdateTarget).toHaveBeenCalledTimes(3);
    expect(onUpdateTarget).toHaveBeenLastCalledWith('p-1', { column_name: 'temperature_a' });
  });

  it('submits a return to an earlier value and can retry it after the newer commit fails', () => {
    const onUpdateTarget = vi.fn();
    const { rerender } = renderTable({ onUpdateTarget });
    const commit = (value: string) => {
      const input = screen.getByTestId('step4-target-column-p-1');
      fireEvent.change(input, { target: { value } });
      fireEvent.keyDown(input, { key: 'Enter' });
    };
    const withTarget = (column_name: string, save_state: DbTarget['save_state']) => rerender(tableElement({
      targets: { ...savedTargets, 'p-1': { ...savedTargets['p-1'], column_name, save_state } },
      onUpdateTarget,
    }));

    commit('temperature_a');
    commit('temperature_b');
    commit('temperature_a');
    expect(onUpdateTarget).toHaveBeenCalledTimes(3);
    expect(onUpdateTarget).toHaveBeenLastCalledWith('p-1', { column_name: 'temperature_a' });
    fireEvent.blur(screen.getByTestId('step4-target-column-p-1'));
    expect(onUpdateTarget).toHaveBeenCalledTimes(3);

    withTarget('temperature_b', 'save-error');
    withTarget('temperature_a', 'saving');
    withTarget('temperature_a', 'save-error');
    commit('temperature_a');
    expect(onUpdateTarget).toHaveBeenCalledTimes(4);
    expect(onUpdateTarget).toHaveBeenLastCalledWith('p-1', { column_name: 'temperature_a' });
  });

  it('keeps the earlier acknowledgement pure under StrictMode', () => {
    const onUpdateTarget = vi.fn();
    const { rerender } = renderStrictTable({ onUpdateTarget });
    let input = screen.getByTestId('step4-target-column-p-1');

    fireEvent.change(input, { target: { value: 'temperature_a' } });
    fireEvent.keyDown(input, { key: 'Enter' });
    rerender(strictTableElement({
      targets: {
        ...savedTargets,
        'p-1': { ...savedTargets['p-1'], column_name: 'temperature_a', save_state: 'saving' },
      },
      onUpdateTarget,
    }));

    input = screen.getByTestId('step4-target-column-p-1');
    fireEvent.change(input, { target: { value: 'temperature_b' } });
    rerender(strictTableElement({
      targets: {
        ...savedTargets,
        'p-1': { ...savedTargets['p-1'], column_name: 'temperature_a', save_state: 'saved' },
      },
      onUpdateTarget,
    }));

    expect(screen.getByTestId('step4-target-column-p-1')).toHaveValue('temperature_b');
    expect(screen.queryByTestId('step4-target-remote-conflict-p-1')).not.toBeInTheDocument();
  });

  it('keeps an unacknowledged witness when reloading another remote value', () => {
    const onUpdateTarget = vi.fn();
    const { rerender } = renderTable({ onUpdateTarget });
    let input = screen.getByTestId('step4-target-column-p-1');

    fireEvent.change(input, { target: { value: 'temperature_a' } });
    fireEvent.keyDown(input, { key: 'Enter' });
    rerender(tableElement({
      targets: {
        ...savedTargets,
        'p-1': { ...savedTargets['p-1'], column_name: 'temperature_a', save_state: 'saving' },
      },
      onUpdateTarget,
    }));

    input = screen.getByTestId('step4-target-column-p-1');
    fireEvent.change(input, { target: { value: 'temperature_b' } });
    rerender(tableElement({
      targets: {
        ...savedTargets,
        'p-1': { ...savedTargets['p-1'], column_name: 'temperature_remote', save_state: 'saved' },
      },
      onUpdateTarget,
    }));
    fireEvent.click(screen.getByTestId('step4-target-reload-remote-p-1'));
    expect(screen.getByTestId('step4-target-column-p-1')).toHaveValue('temperature_remote');

    rerender(tableElement({
      targets: {
        ...savedTargets,
        'p-1': { ...savedTargets['p-1'], column_name: 'temperature_a', save_state: 'saved' },
      },
      onUpdateTarget,
    }));

    expect(screen.getByTestId('step4-target-column-p-1')).toHaveValue('temperature_remote');
  });

  it('clears row-local editing state when a row is deleted and later re-added', () => {
    const { rerender } = renderTable();
    const rowB = screen.getByTestId('step4-target-column-p-2');

    fireEvent.change(rowB, { target: { value: 'pressure_draft' } });
    fireEvent.keyDown(rowB, { key: 'Enter' });
    rerender(tableElement({
      points: [points[0]],
      mappings: { 'p-1': mappings['p-1'] },
      targets: {
        'p-1': savedTargets['p-1'],
        'p-2': { ...savedTargets['p-2'], column_name: 'pressure_late_ack', save_state: 'saved' },
      },
    }));
    expect(screen.queryByTestId('step4-target-column-p-2')).not.toBeInTheDocument();

    rerender(tableElement());

    expect(screen.getByTestId('step4-target-column-p-2')).toHaveValue('pressure_main_kpa');
  });
});
