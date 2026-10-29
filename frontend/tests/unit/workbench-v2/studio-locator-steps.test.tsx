import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { Step1Device } from '@/features/datalink/workbench-v2/steps/step1/Step1Device';
import { Step2Rule } from '@/features/datalink/workbench-v2/steps/step2/Step2Rule';
import { Step3Mapping } from '@/features/datalink/workbench-v2/steps/step3/Step3Mapping';
import type { DeviceEditorProps } from '@/features/datalink/workbench-v2/steps/step1/DeviceEditor';
import type { MappingTableProps } from '@/features/datalink/workbench-v2/steps/step3/MappingTable';
import type { Rule } from '@/features/datalink/workbench-v2/state/types';
import { studioV2LocatorHydration } from '@/pages/datalink/workbench-v2/studioV2LocatorHydration';
import { locatorRead } from './studio-locator-fixtures';
vi.mock('@/features/datalink/workbench-v2/steps/step1/DeviceEditor', () => ({ DeviceEditor: ({ device }: DeviceEditorProps) => <div data-testid="chosen-device">{device.id}</div> }));
vi.mock('@/features/datalink/workbench-v2/steps/step2/RuleEditor', () => ({ RuleEditor: ({ rule }: { rule: Rule }) => <div data-testid="chosen-rule">{rule.id}</div> }));
vi.mock('@/features/datalink/workbench-v2/steps/step3/MappingTable', () => ({ MappingTable: ({ points, selectedIdx }: MappingTableProps) =>
  <div data-testid="chosen-mapping">{selectedIdx !== null ? points[selectedIdx]?.id : 'none'}</div> }));
vi.mock('@/features/datalink/workbench-v2/steps/step3/useStep3LiveValues', () => ({ useStep3LiveValues: () => ({ rawValues: {}, connectionByDevice: {}, recoveryByDevice: {} }) }));
vi.mock('@/features/datalink/workbench-v2/steps/step3/MeasurementTemplateSelector', () => ({ MeasurementTemplateSelector: () => null }));
afterEach(cleanup);
function wrapper({ children }: { children: React.ReactNode }) { return <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}>{children}</QueryClientProvider>; }
describe('exact selected setup steps', () => {
  it('Step1 selects exact D instead of first same-named B and does not replace deleted D', () => {
    const read = locatorRead(); if (read.selection.state !== 'selected') throw new Error();
    const state = studioV2LocatorHydration(read)!; state.devices.unshift({ ...state.devices[0], id: 'B' }); const dispatch = vi.fn();
    const view = render(<Step1Device state={state} dispatch={dispatch} onContinue={vi.fn()} locatorSelection={read.selection} />, { wrapper });
    expect(screen.getByTestId('chosen-device')).toHaveTextContent(/^D$/); expect(dispatch).not.toHaveBeenCalled();
    view.rerender(<Step1Device state={{ ...state, devices: state.devices.filter(device => device.id === 'B') }} dispatch={dispatch} onContinue={vi.fn()} locatorSelection={read.selection} />);
    expect(screen.queryByTestId('chosen-device')).not.toBeInTheDocument();
  });
  it('Step2 missing exact rule never falls back to another same-named first rule', () => {
    const read = locatorRead(2); if (read.selection.state !== 'selected') throw new Error(); const state = studioV2LocatorHydration(read)!;
    const dispatch = vi.fn(); const view = render(<Step2Rule state={state} dispatch={dispatch} onContinue={vi.fn()} locatorSelection={read.selection} />, { wrapper });
    expect(screen.getByTestId('chosen-rule')).toHaveTextContent(/^R$/);
    view.rerender(<Step2Rule state={{ ...state, rules: [{ ...state.rules[0], id: 'OTHER' }] }} dispatch={dispatch} onContinue={vi.fn()} locatorSelection={read.selection} />);
    expect(screen.queryByTestId('chosen-rule')).not.toBeInTheDocument(); expect(dispatch).not.toHaveBeenCalled();
  });
  it('Step3 selects native mapping Point identity, suppresses mount initialization, and never falls back', () => {
    const read = locatorRead(3); if (read.selection.state !== 'selected') throw new Error(); const state = studioV2LocatorHydration(read)!;
    state.points.unshift({ ...state.points[0], id: 'same-address-other-point' }); const dispatch = vi.fn();
    const view = render(<Step3Mapping state={state} dispatch={dispatch} onContinue={vi.fn()} onBack={vi.fn()} locatorSelection={read.selection} />, { wrapper });
    expect(screen.getByTestId('chosen-mapping')).toHaveTextContent(/^P-native$/); expect(dispatch).not.toHaveBeenCalled();
    view.rerender(<Step3Mapping state={{ ...state, mappings: {} }} dispatch={dispatch} onContinue={vi.fn()} onBack={vi.fn()} locatorSelection={read.selection} />);
    expect(screen.getByTestId('chosen-mapping')).toHaveTextContent(/^none$/); expect(dispatch).not.toHaveBeenCalled();
  });
  it('ordinary Step3 preserves its existing mount initialization behavior', () => {
    const state = studioV2LocatorHydration(locatorRead(3))!; const dispatch = vi.fn();
    render(<Step3Mapping state={state} dispatch={dispatch} onContinue={vi.fn()} onBack={vi.fn()} />, { wrapper });
    expect(dispatch).toHaveBeenCalledWith(expect.objectContaining({ type: 'initMappingsForPoints' }));
  });
});
