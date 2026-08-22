import { fireEvent, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import { renderDesk, resetIncidentMocks, seedReadyOutputFlow } from './MuiOutputIncidentDesk.reopen.testHarness';

describe('MuiOutputIncidentDesk reopened v2 Output surface', () => {
  beforeEach(() => {
    resetIncidentMocks();
  });

  it('renders incident-desk command surfaces while preserving shared Output workboard panels', async () => {
    seedReadyOutputFlow();

    renderDesk();

    await waitFor(() => {
      expect(screen.getByTestId('output-incident-command-panel')).toBeInTheDocument();
    });
    expect(screen.getByTestId('output-incident-priority-card')).toBeInTheDocument();
    expect(screen.getByTestId('output-incident-summary-strip')).toBeInTheDocument();
    expect(screen.getByTestId('output-incident-handoff-panel')).toBeInTheDocument();
    expect(screen.getByTestId('output-primary-anchor')).toBeInTheDocument();
    expect(screen.getByTestId('output-tag-chips')).toBeInTheDocument();
    expect(screen.getByTestId('register-map-canvas')).toBeInTheDocument();
    expect(screen.getByTestId('modbus-secondary-panels')).toBeInTheDocument();
  });

  it('runs Local Modbus dry-run from the incident command panel', async () => {
    seedReadyOutputFlow();

    renderDesk();

    await waitFor(() => {
      expect(screen.getByTestId('output-incident-primary-action')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId('output-incident-primary-action'));

    await waitFor(() => {
      expect(screen.getByTestId('dry-run-results')).toBeInTheDocument();
    });
  });

  it('switches to the Database desk and opens connector editing from the command panel', async () => {
    seedReadyOutputFlow();

    renderDesk();

    await waitFor(() => {
      expect(screen.getByTestId('output-incident-target-database')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId('output-incident-target-database'));

    await waitFor(() => {
      expect(screen.getByTestId('output-incident-target-database')).toHaveAttribute(
        'aria-selected',
        'true',
      );
    });
    await waitFor(() => {
      expect(screen.getByTestId('database-selected-tag')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId('output-incident-primary-action'));

    await waitFor(() => {
      expect(
        screen.getByLabelText('workbench.output.database.connector.name'),
      ).toBeInTheDocument();
    });
  });
});
