import { screen, within } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  registerShellUiFixtures,
  renderPage,
} from './DatalinkWorkbenchShellUi.testHarness';

describe('DatalinkWorkbench five-region shell', () => {
  registerShellUiFixtures();

  describe('WorkbenchFrame layout', () => {
    it('renders shell regions (step rail embedded in context bar)', () => {
      renderPage();

      expect(screen.getByTestId('workbench-frame')).toBeInTheDocument();
      expect(screen.getByTestId('workbench-context-bar')).toBeInTheDocument();
      expect(
        within(screen.getByTestId('workbench-context-bar')).getByTestId('workbench-step-rail'),
      ).toBeInTheDocument();
      expect(screen.getByTestId('workbench-primary-work-area')).toBeInTheDocument();
      expect(screen.getByTestId('workbench-inspector-panel')).toBeInTheDocument();
      expect(screen.getByTestId('workbench-bottom-summary-bar')).toBeInTheDocument();
    });

    it('locks desktop layout sizing (v2-mui uses sx-based grid)', () => {
      renderPage();

      const frame = screen.getByTestId('workbench-frame');
      const workArea = screen.getByTestId('workbench-primary-work-area');

      expect(frame).toHaveAttribute('data-variant', 'v2-mui');
      expect(frame).toBeInTheDocument();
      expect(workArea).toBeInTheDocument();
    });

    it('does not render the old ActionDock or HeaderBar', () => {
      renderPage();

      expect(screen.queryByText('workbench.actionDock.title')).not.toBeInTheDocument();
      expect(screen.queryByText('workbench.header.eyebrow')).not.toBeInTheDocument();
    });
  });
});
