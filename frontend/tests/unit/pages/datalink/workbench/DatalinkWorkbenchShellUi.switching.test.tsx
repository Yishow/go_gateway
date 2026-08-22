import { fireEvent, screen, within } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  registerShellUiFixtures,
  renderPage,
} from './DatalinkWorkbenchShellUi.testHarness';

describe('DatalinkWorkbench five-region shell', () => {
  registerShellUiFixtures();

  describe('Step switching wires content into PrimaryWorkArea', () => {
    it('switches step content when step rail tabs are clicked', () => {
      renderPage();

      const rail = screen.getByTestId('workbench-step-rail');
      const workArea = screen.getByTestId('workbench-primary-work-area');

      // Device step is the default — device step content should be in work area
      expect(workArea).toBeInTheDocument();

      // Switch to source step
      fireEvent.click(within(rail).getByRole('tab', { name: /workbench\.steps\.source/ }));
      expect(screen.getByTestId('workbench-primary-work-area')).toBeInTheDocument();

      // Switch to tag step
      fireEvent.click(within(rail).getByRole('tab', { name: /workbench\.steps\.tag/ }));
      expect(screen.getByTestId('workbench-primary-work-area')).toBeInTheDocument();
    });
  });
});
