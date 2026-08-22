import { fireEvent, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import {
  renderApp,
  resetFoundationMocks,
} from './DatalinkWorkbenchFoundation.testHarness';

describe('DatalinkWorkbench foundation route', () => {
  beforeEach(() => {
    resetFoundationMocks();
  });

  it('redirects /datalink/workbench into /studio/v2 and renders the V2 entry', async () => {
    window.history.pushState({}, '', '/datalink/workbench');

    await renderApp();

    expect(await screen.findByTestId('workbench-v2-root')).toBeInTheDocument();
    expect(screen.getByTestId('workbench-v2-location-pathname')).toHaveTextContent('/studio/v2');
    expect(window.location.pathname).toBe('/studio/v2');
  });

  it('keeps Step 1 search, protocol filter, and create action inside one primary toolbar', async () => {
    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const toolbar = screen.getByTestId('device-primary-toolbar');

    expect(
      within(toolbar).getByRole('textbox', { name: 'workbench.device.search.label' }),
    ).toBeInTheDocument();
    expect(
      within(toolbar).getByRole('combobox', { name: 'workbench.device.filters.protocol' }),
    ).toBeInTheDocument();
    expect(
      within(toolbar).getByRole('button', { name: 'workbench.device.actions.create' }),
    ).toBeInTheDocument();
    expect(
      within(toolbar).queryByRole('button', { name: 'workbench.device.actions.clone' }),
    ).not.toBeInTheDocument();
    expect(
      within(toolbar).queryByRole('button', { name: 'workbench.device.actions.continue' }),
    ).not.toBeInTheDocument();
  });

  it('drops the decorative Step 1 hero block once the compact toolbar is available', async () => {
    await renderApp();

    expect(screen.queryByText('workbench.device.title')).not.toBeInTheDocument();
    expect(screen.queryByText('workbench.device.description')).not.toBeInTheDocument();
  });

  it('preserves deep-link query context when /datalink/workbench redirects into /studio/v2', async () => {
    window.history.pushState({}, '', '/datalink/workbench?step=output&target=database');

    await renderApp();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    });
    expect(screen.getByTestId('workbench-v2-location-pathname')).toHaveTextContent('/studio/v2');
    expect(screen.getByTestId('workbench-v2-location-search')).toHaveTextContent(
      '?step=output&target=database',
    );
    expect(window.location.pathname).toBe('/studio/v2');
    expect(window.location.search).toContain('step=output');
    expect(window.location.search).toContain('target=database');
  });

  it('does not lock step navigation after a deep-link is applied', async () => {
    window.history.pushState({}, '', '/studio?step=output&target=database');

    await renderApp();

    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /workbench\.steps\.output/ }),
      ).toHaveAttribute('aria-selected', 'true');
    });

    // User clicks the device step in the step rail — must NOT be locked back to output.
    fireEvent.click(
      screen.getByRole('tab', { name: /workbench\.steps\.device/ }),
    );

    expect(
      screen.getByRole('tab', { name: /workbench\.steps\.device/ }),
    ).toHaveAttribute('aria-selected', 'true');
    expect(
      screen.getByRole('tab', { name: /workbench\.steps\.output/ }),
    ).not.toHaveAttribute('aria-selected', 'true');
  });

  it('renders the main operator flow through /studio', async () => {
    window.history.pushState({}, '', '/studio');

    await renderApp();

    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /workbench\.steps\.device/ }),
      ).toHaveAttribute('aria-selected', 'true');
    });

    expect(screen.getByText('Mixer PLC')).toBeInTheDocument();
    expect(window.location.pathname).toBe('/studio');
  });

  it('keeps the normal operator workflow inside /studio from device through output', async () => {
    window.history.pushState({}, '', '/studio');

    await renderApp();

    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /workbench\.steps\.device/ }),
      ).toHaveAttribute('aria-selected', 'true');
    });

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    expect(window.location.pathname).toBe('/studio');

    const contextBar = screen.getByTestId('workbench-context-bar');
    fireEvent.click(within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoSource' }));
    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /workbench\.steps\.source/ }),
      ).toHaveAttribute('aria-selected', 'true');
    });
    await waitFor(() => {
      expect(screen.getByLabelText('workbench.source.planner.startAddress')).toBeInTheDocument();
    });
    expect(window.location.pathname).toBe('/studio');

    fireEvent.click(within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoTag' }));
    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /workbench\.steps\.tag/ }),
      ).toHaveAttribute('aria-selected', 'true');
    });
    expect(window.location.pathname).toBe('/studio');

    fireEvent.click(within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoOutput' }));
    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /workbench\.steps\.output/ }),
      ).toHaveAttribute('aria-selected', 'true');
    });
    expect(screen.getByTestId('active-output-target')).toBeInTheDocument();
    expect(window.location.pathname).toBe('/studio');
  });

  it('keeps source engineering tools inside /studio instead of opening another workflow route', async () => {
    window.history.pushState({}, '', '/studio');

    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      within(screen.getByTestId('workbench-context-bar')).getByRole('button', {
        name: 'workbench.contextBar.actions.gotoSource',
      }),
    );

    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /workbench\.steps\.source/ }),
      ).toHaveAttribute('aria-selected', 'true');
    });

    await waitFor(() => {
      expect(screen.getByTestId('source-toolbar-more-trigger')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByTestId('source-toolbar-more-trigger'));

    expect(screen.getByTestId('source-toolbar-more-menu')).toBeInTheDocument();
    expect(window.location.pathname).toBe('/studio');
  });
});
