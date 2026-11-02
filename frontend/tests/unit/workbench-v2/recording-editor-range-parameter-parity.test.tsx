import { RecordingPointSelection } from '@/features/datalink/workbench-v2/recording-editor/RecordingPointSelection';
import { act, fireEvent, render, renderHook, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { newRecordingRangeBuffer, parseRecordingRangeBuffer } from '@/features/datalink/workbench-v2/recording-editor/recordingRangeBuffer';
import { useRecordingRangeBuffer } from '@/features/datalink/workbench-v2/recording-editor/useRecordingRangeBuffer';
import { RecordingRangeContent } from '@/features/datalink/workbench-v2/recording-editor/RecordingRangeContent';
import type { Rule } from '@/features/datalink/workbench-v2/state/types';
import type { RangeSourceOutcome, RangeSourcePort, RangeSourceRequest } from '@/pages/datalink/workbench-v2/rangeSourceSubmission';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
const device = { id: 'D', name: 'PLC', description: '', protocol: 'modbus_tcp' as const, status: 'draft' as const,
  connection_config: '{}', last_test_at: null, last_test_success: null, last_test_error: '', created_at: '', updated_at: '' };
const buffer = (dataFormat: Rule['data_format'] = 'CDAB', unit = 'kWh') => ({ ...newRecordingRangeBuffer(), name: 'Meter', address: '40001',
  count: '10', dataType: 'float32' as const, scale: '0.1', offset: '-2', dataFormat, unit });
function source() {
  const observations: RangeSourceOutcome[] = [];
  const submit = vi.fn<RangeSourcePort['submit']>(async request => ({ request, status: 'unknown' }));
  const port: RangeSourcePort = { prepare: async () => ({ device, activationAllowed: true }), submit,
    prepareSaved: vi.fn(), inspect: vi.fn(), observations: () => observations };
  return { port, submit, observations };
}
describe('range parameter parity without a second persistence owner', () => {
  it.each<Rule['data_format']>(['', 'ABCD', 'BADC', 'CDAB', 'DCBA'])('keeps rule format %s and leaves units outside the Rule API', dataFormat => {
    const parsed = parseRecordingRangeBuffer(buffer(dataFormat), 'D', 'modbus_tcp');
    expect(parsed?.rule.data_format).toBe(dataFormat);
    expect(parsed?.rule).not.toHaveProperty('unit');
    expect(parsed?.span).toBe(20);
    expect(parsed?.rule.scale_multiplier).toBe(0.1);
  });
  it('rejects an unsupported format rather than silently converting it to default', () => {
    expect(parseRecordingRangeBuffer(buffer('INVALID' as Rule['data_format']), 'D', 'modbus_tcp')).toBeUndefined();
  });
  it('captures format and mapping unit together and preserves newer raw input after unknown', async () => {
    const { port, submit } = source();
    let resolve!: (value: RangeSourceOutcome) => void;
    submit.mockImplementationOnce(() => new Promise(yes => { resolve = yes; }));
    const h = renderHook(() => useRecordingRangeBuffer('W', port));
    await act(async () => { await h.result.current.start(device); });
    act(() => h.result.current.change(buffer()));
    let pending!: Promise<void>;
    act(() => { pending = h.result.current.submit(); });
    const request = submit.mock.calls[0][0];
    expect(request.buffer.data_format).toBe('CDAB');
    expect(request).toHaveProperty('mappingDefaults.unit', 'kWh');
    expect(request.buffer).not.toHaveProperty('unit');
    act(() => h.result.current.change(buffer('DCBA', 'MWh')));
    await act(async () => { resolve({ request, status: 'unknown' }); await pending; });
    expect(h.result.current.buffer).toMatchObject({ dataFormat: 'DCBA', unit: 'MWh' });
    expect(h.result.current.dirty).toBe(true);
    expect(h.result.current.canSubmit).toBe(false);
    expect(submit).toHaveBeenCalledTimes(1);
  });
  it('restores the submitted format and mapping unit when resuming an unresolved request', () => {
    const { port, observations } = source();
    const rule = parseRecordingRangeBuffer(buffer(), 'D', 'modbus_tcp')!.rule;
    const request: RangeSourceRequest & { mappingDefaults: { unit: string } } = { workspaceId: 'W', deviceId: 'D',
      requestId: 'R', sessionId: 'S', generation: 5, buffer: { ...rule, data_format: 'CDAB' }, mappingDefaults: { unit: 'kWh' } };
    observations.push({ request, status: 'unknown' });
    const h = renderHook(() => useRecordingRangeBuffer('W', port));
    act(() => { expect(h.result.current.resume()).toBe(true); });
    expect(h.result.current.buffer).toMatchObject({ dataFormat: 'CDAB', unit: 'kWh' });
    expect(h.result.current.canSubmit).toBe(false);
  });
  it('shows unknown schedule evidence before persistence and labels the advanced link as global settings', async () => {
    const { port } = source();
    const h = renderHook(() => useRecordingRangeBuffer('W', port));
    await act(async () => { await h.result.current.start(device); });
    render(<RecordingRangeContent controller={h.result.current} available onBack={vi.fn()} onSubmit={vi.fn()} onInspect={vi.fn()} onPrepare={vi.fn()} />);
    expect(screen.getByTestId('recording-range-schedule')).toHaveTextContent('range.scheduleUnassigned');
    const link = screen.getByRole('link', { name: 'range.globalScheduleSettings' });
    expect(link).toHaveAttribute('href', '/studio/v2');
    expect(link).toHaveAttribute('target', '_blank');
    expect(link).toHaveAttribute('rel', 'noopener noreferrer');
    act(() => h.result.current.change(buffer()));
    fireEvent.click(link);
    expect(h.result.current.buffer).toMatchObject({ name: 'Meter', dataFormat: 'CDAB', unit: 'kWh' });
    expect(h.result.current.dirty).toBe(true);
    expect(screen.getByTestId('recording-range-schedule')).not.toHaveTextContent(/1000|100 ms/);
  });
  it('describes shared source saving rather than point selection in the range modal', async () => {
    const { port } = source();
    render(<RecordingPointSelection scope="S" workspaceId="W" candidates={[]} members={[]} onChange={vi.fn()} disabled={false}
      devices={[device]} deviceSource={{ range: port, submit: vi.fn(), inspect: vi.fn(), observations: () => [] }} />);
    fireEvent.click(screen.getByTestId('recording-points-open'));
    fireEvent.change(screen.getByLabelText('picker.device'), { target: { value: 'D' } });
    fireEvent.click(screen.getByRole('button', { name: 'picker.newSource' }));
    await screen.findByLabelText('source.fields.prefix');
    expect(screen.getByRole('dialog')).toHaveAccessibleDescription('range.description');
  });
  it('exposes editable range format and mapping unit, with no per-rule polling input', async () => {
    const { port } = source();
    const h = renderHook(() => useRecordingRangeBuffer('W', port));
    await act(async () => { await h.result.current.start(device); });
    const view = render(<RecordingRangeContent controller={h.result.current} available onBack={vi.fn()} onSubmit={vi.fn()} onInspect={vi.fn()} onPrepare={vi.fn()} />);
    fireEvent.change(screen.getByLabelText('range.dataFormat'), { target: { value: 'CDAB' } });
    view.rerender(<RecordingRangeContent controller={h.result.current} available onBack={vi.fn()} onSubmit={vi.fn()} onInspect={vi.fn()} onPrepare={vi.fn()} />);
    fireEvent.change(screen.getByLabelText('range.unit'), { target: { value: 'kWh' } });
    expect(h.result.current.buffer).toMatchObject({ dataFormat: 'CDAB', unit: 'kWh' });
    expect(screen.queryByRole('spinbutton', { name: /poll|輪詢/i })).not.toBeInTheDocument();
  });
});
