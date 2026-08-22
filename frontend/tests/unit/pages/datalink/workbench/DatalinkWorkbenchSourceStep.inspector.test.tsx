import { fireEvent, screen, waitFor, within } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  getSourceStepMocks,
  openSourceRuleLayerRulesTab,
  registerSourceStepFixtures,
  renderPage,
} from './DatalinkWorkbenchSourceStep.testHarness';
const { mockSourceRules, mockUpdateSourceRuleMutation } = getSourceStepMocks();

describe('DatalinkWorkbench source step', () => {
  registerSourceStepFixtures();

  it('explains unmanaged existing points in the inspector', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByTestId('address-cell-40005'));

    expect(screen.getByTestId('source-span-inspector')).toBeInTheDocument();
    expect(screen.getByText('workbench.source.inspector.span.unmanagedNotice')).toBeInTheDocument();
  });

  it('preserves skipped addresses when editing non-geometry fields on a persisted rule', async () => {
    mockSourceRules.splice(0, mockSourceRules.length, {
      id: 'persisted-rule-1',
      device_id: 'device-1',
      start_address: '40001',
      count: 2,
      data_type: 'int16',
      naming_prefix: 'SRC',
      enabled: true,
      locked: false,
      origin: 'manual',
      template_name: '',
      skipped_addresses: ['40002'],
      created_at: '',
      updated_at: '2026-03-19T00:00:00Z',
    });

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    openSourceRuleLayerRulesTab();
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.ruleLayer.editStart' }));

    const editForm = screen.getByTestId('rule-inline-edit-form');
    fireEvent.change(within(editForm).getByLabelText('workbench.source.planner.scaleMultiplier'), {
      target: { value: '2' },
    });
    fireEvent.click(within(editForm).getByRole('button', { name: 'workbench.source.ruleLayer.editSave' }));

    await waitFor(() => {
      expect(mockUpdateSourceRuleMutation.mutateAsync).toHaveBeenCalledWith({
        id: 'persisted-rule-1',
        data: expect.objectContaining({
          start_address: '40001',
          skipped_addresses: ['40002'],
          scale_multiplier: 2,
        }),
      });
    });
  });

  it('shows rule inspector details when a source rule is selected', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(screen.getByTestId('source-rule-rule-1'));

    expect(screen.getByTestId('source-rule-inspector')).toBeInTheDocument();
    expect(screen.getByTestId('source-rule-coverage')).toHaveTextContent('40001');
  });

  it('shows span inspector details when an address cell is selected', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(screen.getByTestId('address-cell-40001'));

    expect(screen.getByTestId('source-span-inspector')).toBeInTheDocument();
    expect(screen.getByTestId('source-span-link-state')).toHaveTextContent(
      'workbench.source.link.needsPoint',
    );
  });

  it('marks the selected address cell with aria-pressed', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const addressCell = screen.getByTestId('address-cell-40001');

    expect(addressCell).toHaveAttribute('aria-pressed', 'false');

    fireEvent.click(addressCell);

    expect(addressCell).toHaveAttribute('aria-pressed', 'true');
  });

  it('retargets the inspector to the edited rule after inline save', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(screen.getByTestId('address-cell-40001'));
    expect(screen.getByTestId('source-span-inspector')).toBeInTheDocument();

    const ruleCard = screen.getByTestId('source-rule-rule-1-card');
    openSourceRuleLayerRulesTab();
    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editStart' }),
    );
    const editForm = within(ruleCard).getByTestId('rule-inline-edit-form');
    fireEvent.change(within(editForm).getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40010' },
    });
    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editSave' }),
    );

    expect(screen.queryByTestId('source-span-inspector')).not.toBeInTheDocument();
    expect(screen.getByTestId('source-rule-inspector')).toBeInTheDocument();
  });
});
