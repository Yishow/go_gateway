import { screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import { mocks, renderDesk, resetIncidentMocks, seedReadyOutputFlow } from './MuiOutputIncidentDesk.reopen.testHarness';

describe('MuiOutputIncidentDesk reopened v2 Output surface', () => {
  beforeEach(() => {
    resetIncidentMocks();
  });

  it('surfaces blocker-first recovery copy above the shared Output workboard', async () => {
    seedReadyOutputFlow();
    mocks.candidateViews['rule-1'] = {
      ...mocks.candidateViews['rule-1'],
      local_modbus_outputs: {
        ...mocks.candidateViews['rule-1']!.local_modbus_outputs,
        status: 'blocked',
        reason: 'port 5020 is already in use',
      },
    };

    renderDesk();

    await waitFor(() => {
      expect(screen.getByTestId('output-incident-priority-card')).toBeInTheDocument();
    });

    expect(screen.getByTestId('output-incident-priority-card')).toHaveTextContent(
      'port 5020 is already in use',
    );
    expect(screen.getByTestId('output-incident-goto-tag')).toBeInTheDocument();
  });

  it('shows calmer summary layout with primary/secondary rows and calmer handoff copy', async () => {
    seedReadyOutputFlow();

    renderDesk();

    await waitFor(() => {
      expect(screen.getByTestId('output-incident-summary-strip')).toBeInTheDocument();
    });

    // priority card still exists as first review surface
    expect(screen.getByTestId('output-incident-priority-card')).toBeInTheDocument();

    // calmer eyebrow
    const strip = screen.getByTestId('output-incident-summary-strip');
    expect(strip).toHaveTextContent('Incident context');

    // primary row: linked + scoped metrics
    const primaryRow = screen.getByTestId('output-incident-summary-primary');
    expect(primaryRow).toBeInTheDocument();
    expect(primaryRow).toHaveTextContent('linked');
    expect(primaryRow).toHaveTextContent('scoped');

    // secondary row: attention chip with calmer wording + revision chip
    const secondaryRow = screen.getByTestId('output-incident-summary-secondary');
    expect(secondaryRow).toBeInTheDocument();
    expect(secondaryRow).toHaveTextContent('needs review');
    expect(secondaryRow).toHaveTextContent('rev rev-1');

    // revision chip must NOT appear in the primary row
    expect(primaryRow).not.toHaveTextContent('rev rev-1');

    // calmer handoff panel
    const handoff = screen.getByTestId('output-incident-handoff-panel');
    expect(handoff).toHaveTextContent('Next step');
    expect(handoff).toHaveTextContent(
      'The desk context is aligned with the shared workboard for this target.',
    );

    // repair path still intact
    expect(screen.getByTestId('output-incident-goto-tag')).toBeInTheDocument();
  });

  it('shows calmer noRule handoff copy when no source rule is available', async () => {
    // no data seeded — pass empty string so Bootstrap skips setFocusedRuleId
    // and the rule list is empty, making ruleId null → noRuleSelected = true
    renderDesk('');

    await waitFor(() => {
      expect(screen.getByTestId('output-incident-handoff-panel')).toBeInTheDocument();
    });

    const handoff = screen.getByTestId('output-incident-handoff-panel');
    expect(handoff).toHaveTextContent('Next step');
    expect(handoff).toHaveTextContent(
      'Return to Tag to restore the active source rule, then continue output review.',
    );

    // Go to Tag repair path still present
    expect(screen.getByTestId('output-incident-goto-tag')).toBeInTheDocument();
  });
});
