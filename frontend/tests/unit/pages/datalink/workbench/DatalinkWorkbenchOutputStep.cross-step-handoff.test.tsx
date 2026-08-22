import { fireEvent, screen, waitFor } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { renderPage } from './DatalinkWorkbenchOutputStep.testHarness';

describe('DatalinkWorkbench output step', () => {
  describe('cross-step tag handoff (Step 3 → Step 4)', () => {
    it('pre-selects the focused tag from Step 3 when entering Step 4', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      await screen.findByTestId('output-tag-chips');

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
      const boundRow = await screen.findByTestId('tag-candidate-point-2');
      fireEvent.click(boundRow);

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      await screen.findByTestId('output-tag-chips');

      await waitFor(() => {
        expect(screen.getByTestId('output-candidate-tag-2')).toHaveAttribute('aria-pressed', 'true');
      });
    });

    it('falls back to the first candidate when focusedTagIds do not match', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      await screen.findByTestId('output-tag-chips');

      await waitFor(() => {
        expect(screen.getByTestId('output-candidate-tag-1')).toHaveAttribute('aria-pressed', 'true');
      });
    });
  });
});
