import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, within } from '@testing-library/react';
import { createInstance } from 'i18next';
import { I18nextProvider } from 'react-i18next';
import { describe, expect, it, vi } from 'vitest';
import type { ReactElement } from 'react';
import { MappingTable } from '@/features/datalink/workbench-v2/steps/step3/MappingTable';
import { PipelineSteps } from '@/features/datalink/workbench-v2/steps/step3/PipelineSteps';
import { Step3Mapping } from '@/features/datalink/workbench-v2/steps/step3/Step3Mapping';
import en from '@/i18n/locales/en/workbench-v2.json';
import zh from '@/i18n/locales/zh-TW/workbench-v2.json';
import { savedGroupState } from '../../fixtures/writeGroupState';

vi.mock('@/features/datalink/workbench-v2/steps/step3/MeasurementTemplateSelector', () => ({ MeasurementTemplateSelector: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step3/useStep3LiveValues', () => ({
  useStep3LiveValues: () => ({ rawValues: {}, connectionByDevice: {}, recoveryByDevice: {} }),
}));

function sourceState() {
  const state = savedGroupState();
  const mapping = state.mappings[state.points[0].id];
  mapping.tag_key = 'dev.1780862010602.842.tag.0.rd0';
  mapping.display_name = '主迴路有功功率';
  return state;
}

async function show(ui: ReactElement, lng: string) {
  const i18n = createInstance();
  await i18n.init({ lng, resources: { en: { 'workbench-v2': en }, 'zh-TW': { 'workbench-v2': zh } } });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(<I18nextProvider i18n={i18n}><QueryClientProvider client={client}>{ui}</QueryClientProvider></I18nextProvider>);
}

describe.each(['zh-TW', 'en'])('Step 3 operator source names (%s)', (lng) => {
  it('offers display-name editing without a Tag Key field in basic or expanded conversion controls', async () => {
    const state = sourceState();
    const dispatch = vi.fn();
    await show(<MappingTable points={state.points} mappings={state.mappings} selectedIdx={0}
      setSelectedIdx={vi.fn()} devices={state.devices} rawValues={{}} connectionByDevice={{}} dispatch={dispatch} />, lng);
    expect(screen.queryByTestId(`input-tag-key-${state.points[0].id}`)).not.toBeInTheDocument();
    fireEvent.click(screen.getByTestId('mapping-advanced-toggle'));
    expect(screen.queryByRole('columnheader', { name: /Tag Key/i })).not.toBeInTheDocument();
    expect(screen.queryByDisplayValue(state.mappings[state.points[0].id].tag_key)).not.toBeInTheDocument();
    fireEvent.change(screen.getByTestId(`input-display-name-${state.points[0].id}`), { target: { value: '進線功率' } });
    expect(dispatch).toHaveBeenCalledExactlyOnceWith({ type: 'updateMapping', pointId: state.points[0].id, patch: { display_name: '進線功率' } });
  });

  it('uses display name then address in the output preview while retaining the technical key in data', async () => {
    const state = sourceState();
    const point = state.points[0];
    const mapping = state.mappings[point.id];
    const view = await show(<PipelineSteps point={point} mapping={mapping} rawValue={20} />, lng);
    expect(screen.getByTestId('step-final')).toHaveTextContent(mapping.display_name);
    expect(screen.getByTestId('step-final')).not.toHaveTextContent(mapping.tag_key);
    view.unmount();
    await show(<PipelineSteps point={point} mapping={{ ...mapping, display_name: '' }} rawValue={20} />, lng);
    expect(screen.getByTestId('step-final')).toHaveTextContent(point.address);
    expect(screen.getByTestId('step-final')).not.toHaveTextContent(mapping.tag_key);
    expect(mapping.tag_key).toBe('dev.1780862010602.842.tag.0.rd0');
  });

  it('keeps technical data behind the explicit JSON action and uses a readable dialog title', async () => {
    const state = sourceState();
    const point = state.points[0];
    await show(<MappingTable points={state.points} mappings={state.mappings} selectedIdx={0}
      setSelectedIdx={vi.fn()} devices={state.devices} rawValues={{}} connectionByDevice={{}} dispatch={vi.fn()} />, lng);
    fireEvent.click(screen.getByTestId('mapping-advanced-toggle'));
    fireEvent.click(screen.getByTestId(`payload-modal-trigger-${point.id}`));
    const dialog = within(screen.getByRole('dialog'));
    expect(dialog.getByRole('heading', { name: '主迴路有功功率' })).toBeVisible();
    expect(JSON.parse(dialog.getByTestId('payload-preview').querySelector('code')!.textContent!).tag_id).toBe(state.mappings[point.id].tag_key);
    expect(state.mappings[point.id].tag_id).toBe('tag-0');
  });

  it('keeps missing source identity blocked without asking the operator to fill a Tag Key', async () => {
    const state = sourceState();
    state.mappings[state.points[0].id].tag_key = '';
    await show(<Step3Mapping state={state} dispatch={vi.fn()} onContinue={vi.fn()} onBack={vi.fn()} />, lng);
    const warning = screen.getByTestId('aside-chip-empty-tag');
    expect(warning).not.toHaveTextContent(/Tag Key/i);
    expect(warning).toHaveTextContent(lng === 'zh-TW' ? '來源' : 'source');
    expect(screen.getByTestId('btn-continue')).toBeDisabled();
  });
});
