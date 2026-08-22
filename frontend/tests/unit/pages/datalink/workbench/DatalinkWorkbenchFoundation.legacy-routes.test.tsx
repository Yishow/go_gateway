import { screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import {
  renderApp,
  resetFoundationMocks,
} from './DatalinkWorkbenchFoundation.testHarness';

describe('DatalinkWorkbench foundation route', () => {
  beforeEach(() => {
    resetFoundationMocks();
  });

  it('redirects legacy local modbus entry into the new workbench output step', async () => {
    window.history.pushState({}, '', '/datalink/local-modbus?section=settings');

    await renderApp();

    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /workbench\.steps\.output/ }),
      ).toHaveAttribute('aria-selected', 'true');
    });
    expect(screen.getByTestId('active-output-target')).toHaveTextContent(
      'workbench.bottomSummary.targets.modbus',
    );
    expect(window.location.pathname).toBe('/studio');
    expect(window.location.search).toContain('step=output');
  });

  it('redirects nested legacy workbench entries into /studio while preserving query context', async () => {
    window.history.pushState({}, '', '/datalink/workbench/legacy-output?step=output&target=database');

    await renderApp();

    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /workbench\.steps\.output/ }),
      ).toHaveAttribute('aria-selected', 'true');
    });
    expect(screen.getByTestId('active-output-target')).toHaveTextContent(
      'workbench.bottomSummary.targets.database',
    );
    expect(window.location.pathname).toBe('/studio');
    expect(window.location.search).toContain('step=output');
    expect(window.location.search).toContain('target=database');
  });

  it('redirects /datalink into /studio', async () => {
    window.history.pushState({}, '', '/datalink');

    await renderApp();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    });

    expect(screen.getByTestId('workbench-v2-location-pathname')).toHaveTextContent('/studio/v2');
    expect(window.location.pathname).toBe('/studio/v2');
  });

  it('renders /test without the legacy sidebar layout shell', async () => {
    window.history.pushState({}, '', '/test');

    await renderApp();

    expect(await screen.findByTestId('test-page-mock')).toBeInTheDocument();
    expect(window.location.pathname).toBe('/test');
  });

  it('redirects legacy test utility routes into /test', async () => {
    window.history.pushState({}, '', '/templates');

    await renderApp();

    await waitFor(() => {
      expect(screen.getByTestId('test-page-mock')).toBeInTheDocument();
    });

    expect(window.location.pathname).toBe('/test');
  });
});
