import { act, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import i18n from '../../../src/i18n/config';
import { MappingRow } from '../../../src/features/datalink/workbench-v2/steps/step3/MappingRow';
import { hydrateStudioV2Mapping } from '../../../src/features/datalink/workbench-v2/state/studioV2MappingAutosave';
import { mappingFixture, mappingPoints } from './mapping-autosave-page.testHarness';
vi.mock('../../../src/services/datalink', () => ({ pointAPI: { list: vi.fn().mockResolvedValue([]) }, mappingAPI: { preview: vi.fn() } }));
vi.mock('../../../src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells', () => ({ MappingPreviewCells: () => <td /> }));
vi.mock('../../../src/features/datalink/workbench-v2/steps/step3/MappingPayloadDialog', () => ({ MappingPayloadDialog: () => <td /> }));
afterEach(async () => { await act(async () => { await i18n.changeLanguage('en'); }); });
describe('safe actionable mapping failures', () => {
  it.each(['en', 'zh-TW'].flatMap((language) => [
    { language, code: 'workspace_mapping_conflict', status: 409, action: 'review_source', en: /review.*source/i, zh: /檢查.*來源/ },
    { language, code: 'workspace_mapping_not_found', status: 404, action: 'reload', en: /no longer exists/i, zh: /已不存在/ },
    { language, code: 'workspace_mapping_save_failed', status: 500, action: 'retry', en: /could not be saved/i, zh: /無法儲存/ },
  ]))('renders safe $code copy and functional retry in $language', async ({ language, code, status, action, en, zh }) => {
    await act(async () => { await i18n.changeLanguage(language); });
    const dispatch = vi.fn();
    const point = mappingPoints[0];
    const mapping = { ...hydrateStudioV2Mapping(point, mappingFixture()), save_state: 'save-error' as const, save_error: 'SQL password https://private', save_error_detail: { code, status, action, requestId: 'req_safe_42' } };
    render(<table><tbody><MappingRow point={point} mapping={mapping} devices={[]} isSelected={false} onSelect={vi.fn()} dispatch={dispatch} /></tbody></table>);
    const alert = screen.getByTestId(`mapping-save-state-${point.id}`);
    expect(alert).toHaveTextContent(String(status));
    expect(alert).toHaveTextContent(code);
    expect(alert).toHaveTextContent('req_safe_42');
    expect(alert.textContent).toMatch(language === 'en' ? en : zh);
    expect(alert).not.toHaveTextContent(/SQL|password|https:\/\/private/);
    fireEvent.click(screen.getByTestId(`mapping-save-retry-${point.id}`));
    expect(dispatch).toHaveBeenCalledWith({ type: 'retryMappingSave', pointId: point.id });
  });
});

it('keeps failed removal actionable in the real Step 3 after its source row disappears', async () => {
  const { Step3Mapping } = await import('../../../src/features/datalink/workbench-v2/steps/step3/Step3Mapping');
  const { INITIAL_STATE } = await import('../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State');
  const { QueryClient, QueryClientProvider } = await import('@tanstack/react-query');
  const point = mappingPoints[0], dispatch = vi.fn();
  const mapping = { ...hydrateStudioV2Mapping(point, mappingFixture()), save_state: 'save-error' as const, save_error: 'failure', save_error_detail: { code: 'workspace_mapping_save_failed', status: 500, action: 'retry', operation: 'delete' as const } };
  const state = { ...INITIAL_STATE, devices: [], rules: [], points: [point], mappings: { [point.id]: mapping } };
  render(<QueryClientProvider client={new QueryClient()}><Step3Mapping state={state} dispatch={dispatch} onContinue={vi.fn()} onBack={vi.fn()} /></QueryClientProvider>);
  const alert = screen.getByTestId(`mapping-save-state-${point.id}`);
  expect(alert).toHaveTextContent(/remov/i);
  expect(alert).toHaveTextContent('500');
  fireEvent.click(screen.getByTestId(`mapping-save-retry-${point.id}`));
  expect(dispatch).toHaveBeenCalledWith({ type: 'retryMappingSave', pointId: point.id });
});

it('does not reinitialize Step 3 mapping drafts when only dispatch or save metadata changes', async () => {
  const { Step3Mapping } = await import('../../../src/features/datalink/workbench-v2/steps/step3/Step3Mapping');
  const { INITIAL_STATE } = await import('../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State');
  const { QueryClient, QueryClientProvider } = await import('@tanstack/react-query');
  const client = new QueryClient(), firstDispatch = vi.fn(), nextDispatch = vi.fn();
  const state = { ...INITIAL_STATE, devices: [], rules: [], points: [], mappings: {} };
  const props = { state, onContinue: vi.fn(), onBack: vi.fn() };
  const { rerender } = render(<QueryClientProvider client={client}><Step3Mapping {...props} dispatch={firstDispatch} /></QueryClientProvider>);
  expect(firstDispatch).toHaveBeenCalledWith({ type: 'initMappingsForPoints', points: [] });
  rerender(<QueryClientProvider client={client}><Step3Mapping {...props} dispatch={nextDispatch} /></QueryClientProvider>);
  expect(nextDispatch).not.toHaveBeenCalledWith({ type: 'initMappingsForPoints', points: [] });
});


it('reinitializes when the source point data type changes at the same rule and address', async () => {
  const { Step3Mapping } = await import('../../../src/features/datalink/workbench-v2/steps/step3/Step3Mapping');
  const { INITIAL_STATE } = await import('../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State');
  const { DEFAULT_RULE } = await import('../../../src/features/datalink/workbench-v2/state/defaults');
  const { QueryClient, QueryClientProvider } = await import('@tanstack/react-query');
  const client = new QueryClient(), dispatch = vi.fn();
  const state = { ...INITIAL_STATE, devices: INITIAL_STATE.devices, rules: [{ ...DEFAULT_RULE, count: 1, data_type: 'int16' as const }], mappings: {} };
  const props = { dispatch, onContinue: vi.fn(), onBack: vi.fn() };
  const { rerender } = render(<QueryClientProvider client={client}><Step3Mapping {...props} state={state} /></QueryClientProvider>);
  dispatch.mockClear();
  rerender(<QueryClientProvider client={client}><Step3Mapping {...props} state={{ ...state, rules: [{ ...state.rules[0], data_type: 'uint64' }] }} /></QueryClientProvider>);
  expect(dispatch).toHaveBeenCalledWith({ type: 'initMappingsForPoints', points: [expect.objectContaining({ data_type: 'uint64' })] });
});

it('shows a fixed safe cleanup warning after successful mapping deletion with no retry', async () => {
  const { Step3Mapping } = await import('../../../src/features/datalink/workbench-v2/steps/step3/Step3Mapping');
  const { INITIAL_STATE } = await import('../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State');
  const { QueryClient, QueryClientProvider } = await import('@tanstack/react-query');
  const state = { ...INITIAL_STATE, devices: [], rules: [], mappings: {}, mapping_cleanup_incomplete: true };
  render(<QueryClientProvider client={new QueryClient()}><Step3Mapping state={state} dispatch={vi.fn()} onContinue={vi.fn()} onBack={vi.fn()} /></QueryClientProvider>);
  expect(screen.getByTestId('mapping-cleanup-warning')).toHaveTextContent(/completed.*cleanup/i);
  expect(screen.queryByTestId('mapping-save-retry-rule-A-p-0')).not.toBeInTheDocument();
  expect(screen.getByTestId('mapping-cleanup-warning')).not.toHaveTextContent(/SQL|password|https:\/\/private/);
});
