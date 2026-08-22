import { fireEvent, screen, within } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  openSourceRuleLayerRulesTab,
  registerSourceStepFixtures,
  renderPage,
} from './DatalinkWorkbenchSourceStep.testHarness';

describe('DatalinkWorkbench source step', () => {
  registerSourceStepFixtures();

  it('supports inline editing of rule start address, count, and data type', () => {
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

    const ruleCard = screen.getByTestId('source-rule-rule-1-card');
    openSourceRuleLayerRulesTab();
    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editStart' }),
    );

    const editForm = within(ruleCard).getByTestId('rule-inline-edit-form');
    expect(editForm).toBeInTheDocument();

    const addressInput = within(editForm).getByLabelText('workbench.source.planner.startAddress');
    const countInput = within(editForm).getByLabelText('workbench.source.planner.count');

    fireEvent.change(addressInput, { target: { value: '40010' } });
    fireEvent.change(countInput, { target: { value: '3' } });

    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editSave' }),
    );

    expect(screen.getByTestId('address-cell-40010')).toHaveAttribute('data-status', 'planned');
    expect(screen.getByTestId('address-cell-40012')).toHaveAttribute('data-status', 'planned');
  });

  it('cancels inline rule editing without modifying the rule', () => {
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

    const ruleCard = screen.getByTestId('source-rule-rule-1-card');
    openSourceRuleLayerRulesTab();
    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editStart' }),
    );

    const editForm = within(ruleCard).getByTestId('rule-inline-edit-form');
    fireEvent.change(within(editForm).getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40099' },
    });

    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editCancel' }),
    );

    expect(within(ruleCard).queryByTestId('rule-inline-edit-form')).not.toBeInTheDocument();
    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-status', 'planned');
    expect(screen.queryByTestId('address-cell-40099')).not.toBeInTheDocument();
  });

  it('immediately updates the canvas when a rule is inline-edited', () => {
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

    expect(screen.getByTestId('source-summary-ready-count')).toHaveTextContent('2');

    const ruleCard = screen.getByTestId('source-rule-rule-1-card');
    openSourceRuleLayerRulesTab();
    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editStart' }),
    );

    const editForm = within(ruleCard).getByTestId('rule-inline-edit-form');
    fireEvent.change(within(editForm).getByLabelText('workbench.source.planner.count'), {
      target: { value: '4' },
    });
    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editSave' }),
    );

    expect(screen.getByTestId('source-summary-ready-count')).toHaveTextContent('4');
  });

  it('does not show selection toolbar for non-planned cells', () => {
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

    fireEvent.click(screen.getByTestId('address-cell-40005'));

    expect(screen.queryByTestId('source-selection-toolbar')).not.toBeInTheDocument();
  });
});
