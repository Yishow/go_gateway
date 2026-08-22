import { fireEvent, screen, within } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  getShellUiMocks,
  registerShellUiFixtures,
  renderPage,
} from './DatalinkWorkbenchShellUi.testHarness';

const { mockDevices, mockPoints } = getShellUiMocks();

describe('DatalinkWorkbench five-region shell', () => {
  registerShellUiFixtures();

  describe('WorkbenchStepRail', () => {
    it('renders all four steps in the rail', () => {
      renderPage();

      const rail = screen.getByTestId('workbench-step-rail');
      expect(within(rail).getByRole('tab', { name: /workbench\.steps\.device/ })).toBeInTheDocument();
      expect(within(rail).getByRole('tab', { name: /workbench\.steps\.source/ })).toBeInTheDocument();
      expect(within(rail).getByRole('tab', { name: /workbench\.steps\.tag/ })).toBeInTheDocument();
      expect(within(rail).getByRole('tab', { name: /workbench\.steps\.output/ })).toBeInTheDocument();
    });

    it('highlights the active step with aria-selected (MUI Tabs)', () => {
      renderPage();

      const rail = screen.getByTestId('workbench-step-rail');
      const deviceTab = within(rail).getByRole('tab', { name: /workbench\.steps\.device/ });
      expect(deviceTab).toHaveAttribute('aria-selected', 'true');

      fireEvent.click(within(rail).getByRole('tab', { name: /workbench\.steps\.source/ }));
      expect(deviceTab).toHaveAttribute('aria-selected', 'false');
      expect(within(rail).getByRole('tab', { name: /workbench\.steps\.source/ })).toHaveAttribute('aria-selected', 'true');
    });
  });

  describe('WorkbenchContextBar', () => {
    it('shows device name when a device is selected', () => {
      renderPage();

      // Select the device first by clicking it in the device step
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      expect(screen.getByTestId('context-bar-device-name')).toHaveTextContent('Mixer PLC');
    });

    it('renders a compact step summary instead of capability chips when a device is selected', () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      const summary = screen.getByTestId('context-bar-step-summary');
      expect(summary).toHaveAttribute('data-active-step', 'device');
      expect(summary).toHaveTextContent('Mixer PLC');
      expect(screen.queryByTestId('context-bar-capability-unit-id')).not.toBeInTheDocument();
      expect(screen.queryByTestId('context-bar-capability-address-base')).not.toBeInTheDocument();
      expect(screen.queryByTestId('context-bar-capability-word-order')).not.toBeInTheDocument();
      expect(screen.queryByTestId('context-bar-capability-protocol-traits')).not.toBeInTheDocument();
      expect(screen.queryByTestId('context-bar-test-status')).not.toBeInTheDocument();
    });

    it('shows no-device label when no device is selected', () => {
      mockDevices.splice(0, mockDevices.length);
      renderPage();

      expect(screen.getByTestId('context-bar-device-name')).toHaveTextContent('workbench.contextBar.noDevice');
    });

    it('shows a single primary action instead of a quick-action cluster', () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      const contextBar = screen.getByTestId('workbench-context-bar');
      expect(
        within(contextBar).getByRole('button', {
          name: 'workbench.contextBar.actions.gotoSource',
        }),
      ).toBeInTheDocument();
    });

    it('keeps the primary action aligned with the next step instead of skipping from source to output', () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(screen.getByRole('button', { name: 'workbench.contextBar.actions.gotoSource' }));

      const contextBar = screen.getByTestId('workbench-context-bar');
      expect(screen.getByTestId('context-bar-step-summary')).toHaveAttribute('data-active-step', 'source');
      expect(
        within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoTag' }),
      ).toBeInTheDocument();
    });

    it('disables the source-step primary action until source points exist', () => {
      mockPoints.splice(0, mockPoints.length);
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(screen.getByRole('button', { name: 'workbench.contextBar.actions.gotoSource' }));

      expect(screen.getByTestId('context-bar-step-summary')).toHaveAttribute('data-active-step', 'source');
      expect(
        within(screen.getByTestId('workbench-context-bar')).getByRole('button', {
          name: 'workbench.contextBar.actions.gotoTag',
        }),
      ).toBeDisabled();
    });

    it('renders the source command banner when the source step is active', () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(screen.getByRole('button', { name: 'workbench.contextBar.actions.gotoSource' }));

      expect(screen.getByTestId('source-command-banner')).toBeInTheDocument();
      expect(screen.getByTestId('source-command-banner')).toHaveTextContent(
        'workbench.source.sentryBanner.title.build',
      );
    });

    it('keeps the source command rail to the left of the source workspace', () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(screen.getByRole('button', { name: 'workbench.contextBar.actions.gotoSource' }));

      const workspaceSkeleton = screen.getByTestId('source-workspace-skeleton');
      const commandRail = screen.getByTestId('source-workspace-command-rail');
      const workspace = screen.getByTestId('sentry-source-workspace');

      expect(commandRail).toContainElement(screen.getByTestId('source-workspace-command-strip'));
      expect(commandRail).toContainElement(screen.getByTestId('source-workspace-handoff-strip'));
      expect(workspaceSkeleton.firstElementChild).toBe(commandRail);
      expect(commandRail.nextElementSibling).toBe(workspace);
    });

    it('prioritizes the build workspace over auxiliary summary panels in source mode', () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(screen.getByRole('button', { name: 'workbench.contextBar.actions.gotoSource' }));

      expect(screen.getByTestId('source-build-workspace')).toBeInTheDocument();
      expect(screen.getByTestId('source-build-rule-panel')).toBeInTheDocument();
      expect(screen.getByTestId('source-build-canvas-panel')).toBeInTheDocument();
      expect(screen.queryByTestId('source-workspace-summary-strip')).not.toBeInTheDocument();
      expect(screen.queryByTestId('source-workspace-diagnostics-strip')).not.toBeInTheDocument();
    });

    it('allocates more width to the build canvas than the rule layer on desktop', () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(screen.getByRole('button', { name: 'workbench.contextBar.actions.gotoSource' }));

      const buildWorkspace = screen.getByTestId('source-build-workspace');
      const workspaceGrid = buildWorkspace.firstElementChild as HTMLElement | null;

      expect(workspaceGrid).not.toBeNull();
      expect(workspaceGrid).toHaveClass('xl:grid-cols-[minmax(20rem,0.92fr)_minmax(0,1.08fr)]');
    });

    it('returns to tag review on step 4 when output scope has no active rule handoff', () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      const contextBar = screen.getByTestId('workbench-context-bar');
      fireEvent.click(within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoSource' }));
      fireEvent.click(within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoTag' }));
      fireEvent.click(within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoOutput' }));

      expect(screen.getByTestId('context-bar-step-summary')).toHaveAttribute('data-active-step', 'output');
      const outputCta = within(contextBar).getByTestId('context-bar-primary-action');
      expect(outputCta).toBeInTheDocument();
      expect(outputCta).toHaveTextContent('workbench.shell.actions.returnTag');
      expect(outputCta).not.toHaveTextContent('workbench.contextBar.actions.switchDevice');
    });

    it('routes back to tag review when output scope must recover the active rule handoff', () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      const contextBar = screen.getByTestId('workbench-context-bar');
      fireEvent.click(within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoSource' }));
      fireEvent.click(within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoTag' }));
      fireEvent.click(within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoOutput' }));

      const action = within(contextBar).getByTestId('context-bar-primary-action');
      expect(action).toHaveTextContent('workbench.shell.actions.returnTag');

      fireEvent.click(action);

      expect(screen.getByTestId('context-bar-step-summary')).toHaveAttribute('data-active-step', 'tag');
    });
  });
});
