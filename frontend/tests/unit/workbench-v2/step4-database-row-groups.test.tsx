import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { Step4Database } from '../../../src/features/datalink/workbench-v2/steps/step4/Step4Database';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

describe('Step 4 Database row groups', () => {
  it('建立 row group 時只新增群組，不自動勾選或更新所有 target', () => {
    const dispatch = vi.fn();
    const state = {
      ...INITIAL_STATE,
      points: [
        { id: 'p-1', device_id: 'd-1', rule_id: 'r-1', rule_name: 'Holding Registers', name: 't_1', address: '40001', data_type: 'int16', function: 'holding_register' as const, width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0 },
        { id: 'p-2', device_id: 'd-1', rule_id: 'r-1', rule_name: 'Holding Registers', name: 't_2', address: '40002', data_type: 'int16', function: 'holding_register' as const, width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0 },
      ],
      mappings: {
        'p-1': { point_id: 'p-1', tag_key: 'line1.t_1', display_name: 'T1', unit: 'C', target_type: 'float64' as const, scale: 1, offset: 0, enabled: true },
        'p-2': { point_id: 'p-2', tag_key: 'line1.t_2', display_name: 'T2', unit: 'C', target_type: 'float64' as const, scale: 1, offset: 0, enabled: true },
      },
      db: {
        ...INITIAL_STATE.db,
        row_groups: [{
          id: 'group-existing',
          table_schema: 'public',
          table_name: 'sensor_readings',
          member_point_ids: ['p-1', 'p-2'],
          group_key_columns: ['ts'],
        }],
        targets: {
          'p-1': { tag_id: 'tag.line1.t_1', column_name: 'temp_in_c', enabled: true, row_group_id: 'group-existing' },
          'p-2': { tag_id: 'tag.line1.t_2', column_name: 'temp_out_c', enabled: true, row_group_id: 'group-existing' },
        },
      },
    };

    render(<Step4Database state={state} dispatch={dispatch} onCommit={() => { }} />);

    expect(screen.getByTestId('step4-row-group-list')).toHaveTextContent('group-existing');
    fireEvent.click(screen.getByTestId('step4-row-group-create'));

    expect(dispatch).toHaveBeenCalledWith(expect.objectContaining({
      type: 'setDbRowGroups',
      rowGroups: expect.arrayContaining([expect.objectContaining({ id: 'row-group-2', member_point_ids: [] })]),
    }));
    expect(dispatch).not.toHaveBeenCalledWith(expect.objectContaining({ type: 'updateDbTarget' }));
    expect(screen.queryByTestId('step4-row-group-member-group-existing-p-2')).not.toBeInTheDocument();
  });

  it('在寫入欄位列選擇 row group 時才更新該 target 與群組成員', () => {
    const dispatch = vi.fn();
    const state = {
      ...INITIAL_STATE,
      points: [
        { id: 'p-1', device_id: 'd-1', rule_id: 'r-1', rule_name: 'Holding Registers', name: 't_1', address: '40001', data_type: 'int16', function: 'holding_register' as const, width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0 },
      ],
      mappings: {
        'p-1': { point_id: 'p-1', tag_key: 'line1.t_1', display_name: 'T1', unit: 'C', target_type: 'float64' as const, scale: 1, offset: 0, enabled: true },
      },
      db: {
        ...INITIAL_STATE.db,
        row_groups: [{ id: 'group-existing', table_schema: 'public', table_name: 'sensor_readings', member_point_ids: [], group_key_columns: ['ts'] }],
        targets: { 'p-1': { tag_id: 'tag.line1.t_1', column_name: 'temp_in_c', enabled: true } },
      },
    };

    render(<Step4Database state={state} dispatch={dispatch} onCommit={() => { }} />);

    fireEvent.change(screen.getByTestId('step4-target-row-group-p-1'), { target: { value: 'group-existing' } });

    expect(dispatch).toHaveBeenCalledWith(expect.objectContaining({
      type: 'setDbRowGroups',
      rowGroups: [expect.objectContaining({ id: 'group-existing', member_point_ids: ['p-1'] })],
    }));
    expect(dispatch).toHaveBeenCalledWith({
      type: 'updateDbTarget',
      pointId: 'p-1',
      patch: { row_group_id: 'group-existing' },
    });
  });

  it('變更 connector scope 時應清掉過期 row group 與 target 綁定', () => {
    const dispatch = vi.fn();
    const state = {
      ...INITIAL_STATE,
      points: [
        { id: 'p-1', device_id: 'd-1', rule_id: 'r-1', rule_name: 'Holding Registers', name: 't_1', address: '40001', data_type: 'int16', function: 'holding_register' as const, width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0 },
        { id: 'p-2', device_id: 'd-1', rule_id: 'r-1', rule_name: 'Holding Registers', name: 't_2', address: '40002', data_type: 'int16', function: 'holding_register' as const, width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0 },
      ],
      mappings: {
        'p-1': { point_id: 'p-1', tag_key: 'line1.t_1', display_name: 'T1', unit: 'C', target_type: 'float64' as const, scale: 1, offset: 0, enabled: true },
        'p-2': { point_id: 'p-2', tag_key: 'line1.t_2', display_name: 'T2', unit: 'C', target_type: 'float64' as const, scale: 1, offset: 0, enabled: true },
      },
      db: {
        ...INITIAL_STATE.db,
        connector: {
          ...INITIAL_STATE.db.connector,
          kind: 'postgres' as const,
          schema: 'public',
          table: 'sensor_readings',
        },
        row_groups: [{
          id: 'group-existing',
          table_schema: 'public',
          table_name: 'sensor_readings',
          member_point_ids: ['p-1', 'p-2'],
          group_key_columns: ['ts'],
        }],
        targets: {
          'p-1': { tag_id: 'tag.line1.t_1', column_name: 'temp_in_c', enabled: true, row_group_id: 'group-existing' },
          'p-2': { tag_id: 'tag.line1.t_2', column_name: 'temp_out_c', enabled: true, row_group_id: 'group-existing' },
        },
      },
    };

    render(<Step4Database state={state} dispatch={dispatch} onCommit={() => { }} />);

    dispatch.mockClear();
    fireEvent.change(screen.getByLabelText('step4.field_table'), { target: { value: 'sensor_values_next' } });

    expect(dispatch).toHaveBeenCalledWith({
      type: 'setDbRowGroups',
      rowGroups: [],
    });
    expect(dispatch).toHaveBeenCalledWith({
      type: 'updateDbTarget',
      pointId: 'p-1',
      patch: { row_group_id: undefined },
    });
    expect(dispatch).toHaveBeenCalledWith({
      type: 'updateDbTarget',
      pointId: 'p-2',
      patch: { row_group_id: undefined },
    });
    expect(dispatch).toHaveBeenCalledWith(expect.objectContaining({
      type: 'updateDbConnector',
      patch: expect.objectContaining({ table: 'sensor_values_next' }),
    }));
  });

  it('upsert row group 共用欄位但缺少唯一鍵時應阻擋啟用', () => {
    render(<Step4Database state={{
      ...INITIAL_STATE,
      points: [
        { id: 'p-1', device_id: 'd-1', rule_id: 'r-1', rule_name: 'R1', name: 't_1', address: '40001', data_type: 'int16', function: 'holding_register' as const, width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0 },
        { id: 'p-2', device_id: 'd-1', rule_id: 'r-2', rule_name: 'R2', name: 't_2', address: '40002', data_type: 'int16', function: 'holding_register' as const, width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0 },
      ],
      mappings: {
        'p-1': { point_id: 'p-1', tag_key: 'line.t_1', display_name: 'T1', unit: 'C', target_type: 'float64' as const, scale: 1, offset: 0, enabled: true },
        'p-2': { point_id: 'p-2', tag_key: 'line.t_2', display_name: 'T2', unit: 'C', target_type: 'float64' as const, scale: 1, offset: 0, enabled: true },
      },
      db: {
        ...INITIAL_STATE.db,
        connector: { ...INITIAL_STATE.db.connector, write_mode: 'upsert' as const },
        row_groups: [{ id: 'group-temp', table_schema: 'public', table_name: 'sensor_readings', member_point_ids: ['p-1', 'p-2'], group_key_columns: ['ts'] }],
        targets: {
          'p-1': { tag_id: 'tag.line.t_1', column_name: 'temp_in_c', enabled: true, row_group_id: 'group-temp' },
          'p-2': { tag_id: 'tag.line.t_2', column_name: 'temp_in_c', enabled: true, row_group_id: 'group-temp' },
        },
      },
    }} dispatch={vi.fn()} onCommit={() => { }} />);

    expect(screen.getByRole('button', { name: /step4\.activate_btn/ })).toBeDisabled();
  });
});
