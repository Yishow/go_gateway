import { fireEvent, screen, within } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  openSourceRuleLayerRulesTab,
  registerSourceStepFixtures,
  renderPage,
} from './DatalinkWorkbenchSourceStep.testHarness';

describe('DatalinkWorkbench source step', () => {
  registerSourceStepFixtures();

  it('renders data type selector with grouped optgroups and all backend-aligned types enabled', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const dataTypeSelect = screen.getByLabelText('workbench.source.planner.dataType');
    const options = within(dataTypeSelect).getAllByRole('option');

    const supportedOptions = options.filter((opt) => !(opt as HTMLOptionElement).disabled);
    const disabledOptions = options.filter((opt) => (opt as HTMLOptionElement).disabled);

    expect(supportedOptions.length).toBe(10);
    expect(disabledOptions.length).toBe(0);

    const supportedValues = supportedOptions.map((opt) => (opt as HTMLOptionElement).value);
    expect(supportedValues.sort()).toEqual(
      ['bool', 'float32', 'float64', 'int16', 'int32', 'int64', 'string', 'uint16', 'uint32', 'uint64'].sort(),
    );
  });

  it('uses protect plan wording instead of lock/unlock', () => {
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
    const protectButton = within(ruleCard).getByRole('button', {
      name: 'workbench.source.ruleLayer.protectPlan',
    });
    expect(protectButton).toBeInTheDocument();

    fireEvent.click(protectButton);

    expect(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.unprotectPlan' }),
    ).toBeInTheDocument();
    expect(screen.getByTestId('rule-protect-hint-rule-1')).toHaveTextContent(
      'workbench.source.ruleLayer.protectHint',
    );
  });

  it('snaps selection to the root cell when clicking a merge continuation', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const continuationCell = screen.getByTestId('address-cell-40002');
    expect(continuationCell).toHaveAttribute('data-merge-offset', '1');

    fireEvent.click(continuationCell);

    const rootCell = screen.getByTestId('address-cell-40001');
    expect(rootCell).toHaveAttribute('aria-pressed', 'true');
    expect(continuationCell).toHaveAttribute('aria-pressed', 'true');
  });

  it('highlights the root cell of a selected logical span (continuations are visually merged)', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(screen.getByTestId('address-cell-40001'));

    const rootCell = screen.getByTestId('address-cell-40001');
    const continuationCell = screen.getByTestId('address-cell-40002');

    expect(rootCell.className).toContain('ring-2');
    // Continuation cell is sr-only (visually merged into root via gridColumn span)
    expect(continuationCell).toHaveClass('sr-only');
    expect(continuationCell).toHaveAttribute('aria-pressed', 'true');
  });

  it('aria-pressed is true for all cells of a selected logical span', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(screen.getByTestId('address-cell-40001'));

    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('aria-pressed', 'true');
    expect(screen.getByTestId('address-cell-40002')).toHaveAttribute('aria-pressed', 'true');
  });
});
