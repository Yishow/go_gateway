import { fireEvent, screen, within } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  getShellUiMocks,
  registerShellUiFixtures,
  renderPage,
} from './DatalinkWorkbenchShellUi.testHarness';

const { mockMappings, mockPoints, mockTags } = getShellUiMocks();

describe('DatalinkWorkbench five-region shell', () => {
  registerShellUiFixtures();

  describe('WorkbenchBottomSummaryBar', () => {
    it('shows summary metrics after device is selected', () => {
      renderPage();

      // Select device to populate counters
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      expect(screen.getByTestId('summary-point-count')).toHaveTextContent('2');
      expect(screen.getByTestId('summary-tag-count')).toHaveTextContent('1');
      expect(screen.getByTestId('summary-output-count')).toHaveTextContent('1');
    });

    it('shows readiness indicators with four steps', () => {
      renderPage();

      const summaryBar = screen.getByTestId('workbench-bottom-summary-bar');
      expect(within(summaryBar).getByText('workbench.bottomSummary.readiness.device')).toBeInTheDocument();
      expect(within(summaryBar).getByText('workbench.bottomSummary.readiness.source')).toBeInTheDocument();
      expect(within(summaryBar).getByText('workbench.bottomSummary.readiness.tag')).toBeInTheDocument();
      expect(within(summaryBar).getByText('workbench.bottomSummary.readiness.output')).toBeInTheDocument();
    });

    it('shows the active output target badge', () => {
      renderPage();

      expect(screen.getByTestId('active-output-target')).toBeInTheDocument();
      expect(screen.getByText('workbench.bottomSummary.targets.modbus')).toBeInTheDocument();
    });

    it('propagates rich readiness status via data attributes', () => {
      renderPage();

      // No device selected yet → device readiness should be "draft"
      const deviceIndicator = screen.getByTestId('readiness-device');
      expect(deviceIndicator).toHaveAttribute('data-readiness', 'draft');
    });

    it('marks tag readiness as blocked when source points are missing', () => {
      mockPoints.splice(0, mockPoints.length);
      mockTags.splice(0, mockTags.length);
      mockMappings.splice(0, mockMappings.length);

      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      expect(screen.getByTestId('readiness-source')).toHaveAttribute('data-readiness', 'draft');
      expect(screen.getByTestId('readiness-source')).toHaveAttribute('data-reason', 'no-points');
      expect(screen.getByTestId('readiness-tag')).toHaveAttribute('data-readiness', 'blocked');
      expect(screen.getByTestId('readiness-tag')).toHaveAttribute('data-reason', 'no-source-points');
    });

    it('marks output readiness as blocked until tags are linked', () => {
      mockTags.splice(0, mockTags.length);
      mockMappings.splice(0, mockMappings.length);

      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      expect(screen.getByTestId('readiness-tag')).toHaveAttribute('data-readiness', 'draft');
      expect(screen.getByTestId('readiness-tag')).toHaveAttribute('data-reason', 'no-tags-linked');
      expect(screen.getByTestId('readiness-output')).toHaveAttribute('data-readiness', 'blocked');
      expect(screen.getByTestId('readiness-output')).toHaveAttribute('data-reason', 'no-tags-linked');
    });

    it('emphasizes only the active step and compacts the rest of the readiness strip', () => {
      renderPage();

      expect(screen.getByTestId('readiness-device')).toHaveAttribute('data-emphasis', 'active');
      expect(screen.getByTestId('readiness-source')).toHaveAttribute('data-emphasis', 'compact');
      expect(screen.getByTestId('readiness-tag')).toHaveAttribute('data-emphasis', 'compact');
      expect(screen.getByTestId('readiness-output')).toHaveAttribute('data-emphasis', 'compact');

      fireEvent.click(screen.getByRole('tab', { name: /workbench\.steps\.source/ }));

      expect(screen.getByTestId('readiness-device')).toHaveAttribute('data-emphasis', 'compact');
      expect(screen.getByTestId('readiness-source')).toHaveAttribute('data-emphasis', 'active');
    });
  });
});
