import { fireEvent, screen, waitFor, within } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  getSourceStepMocks,
  registerSourceStepFixtures,
  renderPage,
} from './DatalinkWorkbenchSourceStep.testHarness';
const { mockCreatePointMutation } = getSourceStepMocks();

describe('DatalinkWorkbench source step', () => {
  registerSourceStepFixtures();

  it('renders the source canvas as fixed 16-bit lattice rows', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '3' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const firstRow = screen.getByTestId('source-canvas-row-0');

    expect(firstRow).toHaveAttribute('data-row-start-address', '40001');
    expect(firstRow).toHaveAttribute('data-row-end-address', '40005');
    expect(firstRow).toHaveAttribute('data-lattice-columns', '16');
  });

  it('groups utility tools in the primary toolbar More menu to the left of point actions', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const primaryToolbar = screen.getByTestId('source-primary-toolbar');
    const ruleLayer = screen.getByTestId('source-rule-layer');

    expect(screen.queryByTestId('source-secondary-controls')).not.toBeInTheDocument();
    expect(
      within(primaryToolbar).getByRole('button', { name: 'workbench.source.view.plan' }),
    ).toBeInTheDocument();
    expect(
      within(primaryToolbar).getByLabelText('workbench.source.toolbar.valueFormat'),
    ).toBeInTheDocument();
    expect(
      within(primaryToolbar).queryByLabelText('workbench.source.planner.startAddress'),
    ).not.toBeInTheDocument();
    expect(
      within(primaryToolbar).queryByRole('button', {
        name: 'workbench.source.planner.addRule',
      }),
    ).not.toBeInTheDocument();
    expect(
      within(ruleLayer).getByLabelText('workbench.source.planner.startAddress'),
    ).toBeInTheDocument();
    expect(
      within(ruleLayer).getByRole('button', { name: 'workbench.source.planner.addRule' }),
    ).toBeInTheDocument();

    const moreTrigger = screen.getByTestId('source-toolbar-more-trigger');
    expect(moreTrigger).toBeInTheDocument();
    expect(
      within(primaryToolbar).getByRole('button', {
        name: 'workbench.source.actions.createSelectedPoints',
      }),
    ).toBeInTheDocument();

    fireEvent.click(moreTrigger);
    const menu = screen.getByTestId('source-toolbar-more-menu');
    expect(menu).toBeInTheDocument();
    expect(
      within(menu).getByRole('menuitem', { name: 'workbench.source.toolbar.freezeLive' }),
    ).toBeInTheDocument();
    expect(
      within(menu).getByRole('menuitem', { name: 'workbench.source.toolbar.snapshotCompare' }),
    ).toBeInTheDocument();
    expect(
      within(menu).getByRole('menuitem', { name: 'workbench.source.toolbar.showAudit' }),
    ).toBeInTheDocument();
    expect(
      within(menu).getByRole('menuitem', { name: 'workbench.source.toolbar.saveTemplate' }),
    ).toBeInTheDocument();
    expect(
      within(menu).getByRole('menuitem', { name: 'workbench.source.toolbar.loadTemplate' }),
    ).toBeInTheDocument();
  });

  it('defaults to the build workspace and still lets inspect swap the primary workspace', () => {
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

    const workspaceSkeleton = screen.getByTestId('source-workspace-skeleton');
    const handoffPanel = within(workspaceSkeleton).getByTestId('source-workspace-handoff-strip');
    const buildWorkspace = screen.getByTestId('source-build-workspace');
    const commandBanner = screen.getByTestId('source-command-banner');

    expect(buildWorkspace).toBeInTheDocument();
    expect(within(buildWorkspace).getByTestId('source-build-rule-panel')).toBeInTheDocument();
    expect(within(buildWorkspace).getByTestId('source-build-canvas-panel')).toBeInTheDocument();
    expect(screen.queryByTestId('source-workspace-summary-strip')).not.toBeInTheDocument();
    expect(screen.queryByTestId('source-workspace-diagnostics-strip')).not.toBeInTheDocument();
    expect(commandBanner).toHaveTextContent('workbench.source.sentryBanner.title.build');
    expect(
      within(handoffPanel).getByRole('button', {
        name: 'workbench.source.handoff.toTag',
      }),
    ).toBeEnabled();
    expect(within(handoffPanel).getAllByText('workbench.source.handoff.pointSummary')).toHaveLength(1);

    fireEvent.click(screen.getByTestId('source-desk-tab-inspect'));

    const inspectWorkspace = screen.getByTestId('source-inspect-workspace');
    const inspectSummary = screen.getByTestId('source-workspace-summary-strip');
    const activeRuleSummary = screen.getByTestId('source-active-rule-summary');
    expect(screen.getByTestId('source-workspace-skeleton')).toBeInTheDocument();
    expect(inspectWorkspace).toBeInTheDocument();
    expect(within(inspectWorkspace).getByTestId('source-primary-toolbar')).toBeInTheDocument();
    expect(within(inspectWorkspace).getByTestId('source-canvas-workspace')).toBeInTheDocument();
    expect(within(activeRuleSummary).getByTestId('source-active-rule-address')).toHaveTextContent(
      '40001',
    );
    expect(within(activeRuleSummary).getByTestId('source-active-rule-count')).toHaveTextContent('2');
    expect(within(inspectSummary).getByTestId('source-workspace-ready-count')).toHaveTextContent('2');
    expect(within(inspectSummary).getByTestId('source-workspace-conflict-count')).toHaveTextContent('0');
    expect(within(inspectSummary).getByTestId('source-workspace-protected-count')).toHaveTextContent('0');
    expect(
      within(activeRuleSummary).getByRole('button', {
        name: 'workbench.source.actions.createRulePoints',
      }),
    ).toBeEnabled();
    expect(commandBanner).toHaveTextContent('workbench.source.sentryBanner.title.inspect');
    expect(commandBanner).toHaveTextContent('workbench.source.sentryBanner.status.inspect');
    expect(inspectSummary.nextElementSibling).toBe(screen.getByTestId('source-workspace-diagnostics-strip'));
  });

  it('creates the selected logical span from the summary action', async () => {
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

    fireEvent.click(screen.getByTestId('address-cell-40002'));
    fireEvent.click(
      screen.getByRole('button', { name: 'workbench.source.actions.createSelectedPoints' }),
    );

    await waitFor(() => {
      expect(mockCreatePointMutation.mutateAsync).toHaveBeenCalledTimes(1);
    });

    expect(mockCreatePointMutation.mutateAsync).toHaveBeenCalledWith({
      device_id: 'device-1',
      address: '40001',
      data_type: 'float32',
      name: 'MBT_40001',
    });
  });

  it('treats the active rule layer as the primary workspace and the canvas as supporting context', () => {
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

    expect(screen.getByTestId('source-rule-layer')).toHaveAttribute(
      'data-emphasis',
      'primary',
    );
    expect(screen.getByTestId('source-canvas-workspace')).toHaveAttribute(
      'data-emphasis',
      'supporting',
    );
    expect(
      within(screen.getByTestId('source-rule-layer')).getByTestId('source-rule-rule-1'),
    ).toBeInTheDocument();
  });

  it('keeps the canvas scroll region as the dominant vertical surface in build mode', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '8' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(screen.getByTestId('source-desk-tab-build'));

    const secondaryScrollRegion = screen.getByTestId('source-coverage-overview').parentElement;
    expect(secondaryScrollRegion).not.toBeNull();
    expect(secondaryScrollRegion?.className).not.toContain('max-h-[min(46vh,26rem)]');
  });

  it('preserves merged spans and gap cells across multiple source rules', () => {
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

    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40004' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'int16' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    expect(screen.getByTestId('source-rule-rule-1')).toBeInTheDocument();
    expect(screen.getByTestId('source-rule-rule-2')).toBeInTheDocument();
    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-merge-span', '2');
    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-merge-offset', '0');
    expect(screen.getByTestId('address-cell-40002')).toHaveAttribute('data-merge-span', '2');
    expect(screen.getByTestId('address-cell-40002')).toHaveAttribute('data-merge-offset', '1');
    expect(screen.getByTestId('address-cell-40003')).toHaveAttribute('data-status', 'gap');
    expect(screen.getByTestId('address-cell-40004')).toHaveAttribute('data-status', 'planned');
  });

  it('switches plan, live, and link overlays without changing the lattice addresses', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.view.live' }));
    expect(screen.getByTestId('source-canvas')).toHaveAttribute('data-view-mode', 'live');
    expect(screen.getByTestId('address-cell-40001')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.view.link' }));
    expect(screen.getByTestId('source-canvas')).toHaveAttribute('data-view-mode', 'link');
    expect(screen.getByTestId('address-cell-40001')).toBeInTheDocument();
    expect(screen.getByText('workbench.source.link.unbound')).toBeInTheDocument();
  });
});
