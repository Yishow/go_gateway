import { fireEvent, screen, waitFor, within } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  getSourceStepMocks,
  openSourceRuleLayerRulesTab,
  registerSourceStepFixtures,
  renderPage,
} from './DatalinkWorkbenchSourceStep.testHarness';
const { mockPoints, mockSourceRules, mockCreatePointMutation, mockDeletePointMutation, mockCreateSourceRuleMutation, mockUpdateSourceRuleMutation } = getSourceStepMocks();
import { SOURCE_TEMPLATE_STORAGE_KEY } from '@/features/datalink/sourceTemplateStorage';

describe('DatalinkWorkbench source step', () => {
  registerSourceStepFixtures();

  it('persists rule points through the source-rule API when applying the rule batch action', async () => {
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
    fireEvent.click(screen.getByTestId('source-desk-tab-inspect'));

    fireEvent.click(
      screen.getByRole('button', { name: 'workbench.source.actions.createRulePoints' }),
    );

    await waitFor(() => {
      expect(mockCreateSourceRuleMutation.mutateAsync).toHaveBeenCalledTimes(1);
    });

    expect(mockCreateSourceRuleMutation.mutateAsync).toHaveBeenCalledWith(
      expect.not.objectContaining({ id: expect.anything() }),
    );
    expect(mockCreateSourceRuleMutation.mutateAsync).toHaveBeenCalledWith(
      expect.objectContaining({
        device_id: 'device-1',
        start_address: '40001',
        count: 2,
        data_type: 'int16',
        naming_prefix: 'MBT',
        enabled: true,
        locked: false,
        origin: 'manual',
        skipped_addresses: [],
      }),
    );
    expect(mockCreatePointMutation.mutateAsync).not.toHaveBeenCalled();
  });

  it('reconciles persisted rules via update API when applying the rule batch action', async () => {
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
      skipped_addresses: [],
      created_at: '',
      updated_at: '2026-03-19T00:00:00Z',
    });

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByTestId('source-desk-tab-inspect'));

    fireEvent.click(
      screen.getByRole('button', { name: 'workbench.source.actions.createRulePoints' }),
    );

    await waitFor(() => {
      expect(mockUpdateSourceRuleMutation.mutateAsync).toHaveBeenCalledTimes(1);
    });

    expect(mockUpdateSourceRuleMutation.mutateAsync).toHaveBeenCalledWith({
      id: 'persisted-rule-1',
      data: { skipped_addresses: [] },
    });
    expect(mockCreateSourceRuleMutation.mutateAsync).not.toHaveBeenCalled();
    expect(mockCreatePointMutation.mutateAsync).not.toHaveBeenCalled();
  });

  it('adds and deletes rules directly from the rule layer workflow', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const ruleLayer = screen.getByTestId('source-rule-layer');

    fireEvent.change(within(ruleLayer).getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(within(ruleLayer).getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.click(
      within(ruleLayer).getByRole('button', { name: 'workbench.source.planner.addRule' }),
    );

    const ruleCard = screen.getByTestId('source-rule-rule-1-card');
    expect(ruleCard).toBeInTheDocument();

    openSourceRuleLayerRulesTab();
    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.delete' }),
    );

    expect(screen.queryByTestId('source-rule-rule-1')).not.toBeInTheDocument();
  });

  it('persists applied source rules when navigating away from and back to the source step', () => {
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

    expect(screen.getByTestId('source-rule-rule-1')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));

    expect(screen.getByTestId('source-rule-rule-1')).toBeInTheDocument();
    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-status', 'planned');
  });

  it('saves a source template locally and reapplies it to the planner inputs', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40101' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '3' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.view.live' }));

    fireEvent.click(screen.getByTestId('source-toolbar-more-trigger'));
    fireEvent.click(
      within(screen.getByTestId('source-toolbar-more-menu')).getByRole('menuitem', {
        name: 'workbench.source.toolbar.saveTemplate',
      }),
    );
    fireEvent.change(screen.getByLabelText('workbench.source.templates.name'), {
      target: { value: 'Line Float' },
    });
    fireEvent.click(
      screen.getByRole('button', { name: 'workbench.source.templates.confirmSave' }),
    );

    const storedTemplates = JSON.parse(
      localStorage.getItem(SOURCE_TEMPLATE_STORAGE_KEY) ?? '[]',
    ) as Array<Record<string, unknown>>;
    expect(storedTemplates[0]).toMatchObject({
      name: 'Line Float',
      startAddress: '40101',
      count: 3,
      dataType: 'float32',
      preferredViewMode: 'live',
    });

    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '49999' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'int16' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.view.plan' }));

    fireEvent.click(screen.getByTestId('source-toolbar-more-trigger'));
    fireEvent.click(
      within(screen.getByTestId('source-toolbar-more-menu')).getByRole('menuitem', {
        name: 'workbench.source.toolbar.loadTemplate',
      }),
    );
    fireEvent.click(screen.getByRole('button', { name: 'Line Float' }));

    expect(screen.getByLabelText('workbench.source.planner.startAddress')).toHaveValue('40101');
    expect(screen.getByLabelText('workbench.source.planner.count')).toHaveValue(3);
    expect(screen.getByLabelText('workbench.source.planner.dataType')).toHaveValue('float32');
    expect(screen.getByTestId('source-canvas')).toHaveAttribute('data-view-mode', 'live');
  });

  it('deletes orphaned points when a rule is deleted', () => {
    // Simulate a point that was created from rule-1's planned addresses
    mockPoints.splice(0, mockPoints.length, {
      id: 'point-rule-1',
      device_id: 'device-1',
      name: 'MBT_40001',
      description: '',
      data_type: 'int16',
      address: '40001',
      enabled: true,
      polling_group_id: '',
      last_value: null,
      last_read_at: '',
      last_error: '',
      error_count: 0,
      created_at: '',
      updated_at: '',
    });

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const ruleCard = screen.getByTestId('source-rule-rule-1-card');
    openSourceRuleLayerRulesTab();
    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.delete' }),
    );

    // Point at 40001 should be deleted
    expect(mockDeletePointMutation.mutate).toHaveBeenCalledWith('point-rule-1');
  });
});
