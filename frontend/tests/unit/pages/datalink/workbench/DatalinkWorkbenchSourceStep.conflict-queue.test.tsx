import { fireEvent, screen, within } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  getSourceStepMocks,
  openSourceRuleLayerRulesTab,
  registerSourceStepFixtures,
  renderPage,
} from './DatalinkWorkbenchSourceStep.testHarness';
const { mockPoints } = getSourceStepMocks();

describe('DatalinkWorkbench source step', () => {
  registerSourceStepFixtures();

  it('replaces the conflict hint with an actionable conflict queue', () => {
    mockPoints[0].address = '40002';

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const conflictQueue = screen.getByTestId('source-conflict-queue');
    expect(conflictQueue).toBeInTheDocument();
    expect(within(conflictQueue).getByText('workbench.source.conflictQueue.title')).toBeInTheDocument();
    expect(within(conflictQueue).getByText('workbench.source.conflictQueue.step3Blocked')).toBeInTheDocument();

    const conflictItem = screen.getByTestId('conflict-item-40001');
    expect(conflictItem).toBeInTheDocument();
    expect(within(conflictItem).getByText('workbench.source.conflictQueue.pointOverlap')).toBeInTheDocument();
    expect(
      within(conflictItem).getByRole('button', { name: 'workbench.source.conflictQueue.editRule' }),
    ).toBeInTheDocument();
    expect(
      within(conflictItem).getByRole('button', { name: 'workbench.source.conflictQueue.skipSpan' }),
    ).toBeInTheDocument();
  });

  it('opens visible inline edit when clicking edit rule in the conflict queue', () => {
    mockPoints[0].address = '40002';

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const conflictItem = screen.getByTestId('conflict-item-40001');
    fireEvent.click(
      within(conflictItem).getByRole('button', { name: 'workbench.source.conflictQueue.editRule' }),
    );

    expect(screen.getByTestId('source-rule-layer-tab-rules')).toHaveAttribute('aria-selected', 'true');
    expect(screen.getByTestId('source-rule-layer-tab-planner')).toHaveAttribute('aria-selected', 'false');

    const ruleCard = screen.getByTestId('source-rule-rule-1-card');
    expect(within(ruleCard).getByTestId('rule-inline-edit-form')).toBeVisible();
  });

  it('emits conflict queue item for continuation-only conflicts', () => {
    // Point at 40003, rule plans float32 at 40002 (occupies 40002+40003).
    // Conflict at continuation cell (40003) should resolve to root (40002).
    mockPoints[0].address = '40003';

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40002' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const conflictQueue = screen.getByTestId('source-conflict-queue');
    expect(within(conflictQueue).getByTestId('conflict-item-40002')).toBeInTheDocument();
  });

  it('skip-span excludes only the conflicting logical span, not the entire rule', () => {
    // Point at 40003, rule plans float32 count=2 at 40001 (40001-40002 and 40003-40004).
    // Conflict at 40003 (root of second planned span). Skip should only remove 40003-40004.
    mockPoints[0].address = '40003';

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));
    fireEvent.click(screen.getByTestId('source-desk-tab-inspect'));

    const conflictQueue = screen.getByTestId('source-conflict-queue');
    const conflictItem = within(conflictQueue).getByTestId('conflict-item-40003');

    fireEvent.click(
      within(conflictItem).getByRole('button', { name: 'workbench.source.conflictQueue.skipSpan' }),
    );

    // Conflict should be resolved — queue disappears
    expect(screen.queryByTestId('source-conflict-queue')).not.toBeInTheDocument();

    // First span (40001-40002) should still be planned
    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-status', 'planned');
    expect(screen.getByTestId('address-cell-40002')).toHaveAttribute('data-status', 'planned');

    // Rule should still be enabled (not disabled)
    expect(screen.getByTestId('source-summary-ready-count')).toHaveTextContent('1');
  });

  it('skip-span resolves correct span root when rule-overlap occurs on a continuation cell', () => {
    // Rule A: int32 at 40001 (occupies 40001-40002, mergeOffset 0-1)
    // Rule B: int16 at 40002 (occupies 40002)
    // Conflict at 40002 — continuation cell of Rule A but root cell of Rule B
    // Skip must target Rule B's span root (40002), not Rule A's (40001)
    mockPoints.splice(0, mockPoints.length);

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    // Add Rule A: int32 at 40001
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'int32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    // Add Rule B: int16 at 40002
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'int16' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40002' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    // Conflict should exist
    const conflictQueue = screen.getByTestId('source-conflict-queue');
    const conflictItem = within(conflictQueue).getByTestId('conflict-item-40001');
    expect(conflictItem).toBeInTheDocument();

    // Click skip span — targets Rule B (last ruleId)
    fireEvent.click(
      within(conflictItem).getByRole('button', { name: 'workbench.source.conflictQueue.skipSpan' }),
    );

    // Conflict should be resolved
    expect(screen.queryByTestId('source-conflict-queue')).not.toBeInTheDocument();

    // Rule A's span (40001-40002) should remain planned
    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-status', 'planned');
    expect(screen.getByTestId('address-cell-40002')).toHaveAttribute('data-status', 'planned');

    // Both rules still exist (Rule B is enabled but its only span is skipped)
    expect(screen.getByTestId('source-rule-rule-1')).toBeInTheDocument();
    expect(screen.getByTestId('source-rule-rule-2')).toBeInTheDocument();
  });

  it('clears stale skipped spans after inline geometry edits', () => {
    mockPoints[0].address = '40003';

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));
    fireEvent.click(screen.getByTestId('source-desk-tab-inspect'));

    const conflictItem = screen.getByTestId('conflict-item-40003');
    fireEvent.click(
      within(conflictItem).getByRole('button', { name: 'workbench.source.conflictQueue.skipSpan' }),
    );
    expect(screen.queryByTestId('source-conflict-queue')).not.toBeInTheDocument();
    expect(screen.getByTestId('source-summary-ready-count')).toHaveTextContent('1');

    const ruleCard = screen.getByTestId('source-rule-rule-1-card');
    openSourceRuleLayerRulesTab();

    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editStart' }),
    );
    let editForm = within(ruleCard).getByTestId('rule-inline-edit-form');
    fireEvent.change(within(editForm).getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40010' },
    });
    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editSave' }),
    );

    expect(screen.getByTestId('source-summary-ready-count')).toHaveTextContent('2');
    expect(screen.getByTestId('address-cell-40010')).toHaveAttribute('data-status', 'planned');
    expect(screen.getByTestId('address-cell-40012')).toHaveAttribute('data-status', 'planned');

    mockPoints.splice(0, mockPoints.length);

    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editStart' }),
    );
    editForm = within(ruleCard).getByTestId('rule-inline-edit-form');
    fireEvent.change(within(editForm).getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editSave' }),
    );

    expect(screen.getByTestId('source-summary-ready-count')).toHaveTextContent('2');
    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-status', 'planned');
    expect(screen.getByTestId('address-cell-40003')).toHaveAttribute('data-status', 'planned');
  });
});
