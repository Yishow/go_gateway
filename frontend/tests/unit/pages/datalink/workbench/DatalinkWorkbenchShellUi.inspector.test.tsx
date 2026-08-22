import { screen, within } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  registerShellUiFixtures,
  renderPage,
} from './DatalinkWorkbenchShellUi.testHarness';

describe('DatalinkWorkbench five-region shell', () => {
  registerShellUiFixtures();

  describe('WorkbenchInspectorPanel', () => {
    it('renders the inspector with step-dependent heading', () => {
      renderPage();

      const inspector = screen.getByTestId('workbench-inspector-panel');
      expect(inspector).toBeInTheDocument();
      expect(
        within(inspector).getByText('workbench.device.inspector.emptyTitle'),
      ).toBeInTheDocument();
    });

    it('shows device-step empty state when no item is selected', () => {
      renderPage();

      const inspector = screen.getByTestId('workbench-inspector-panel');

      expect(
        within(inspector).getByText('workbench.device.inspector.emptyTitle'),
      ).toBeInTheDocument();
      expect(
        within(inspector).getByText('workbench.device.inspector.emptyDescription'),
      ).toBeInTheDocument();
      expect(screen.queryByTestId('inspector-selection-context')).not.toBeInTheDocument();
    });
  });
});
