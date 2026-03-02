import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { MemoryRouter } from 'react-router-dom';
import GatewayExpertWorkbenchPage from '../GatewayExpertWorkbenchPage';

describe('GatewayExpertWorkbenchPage', () => {
  function renderComponent() {
    render(
      <MemoryRouter>
        <GatewayExpertWorkbenchPage />
      </MemoryRouter>,
    );
  }

  it('renders core expert workbench sections', () => {
    renderComponent();

    expect(screen.getByText('Expert Workbench')).toBeInTheDocument();
    expect(screen.getByTestId('expert-route-table')).toBeInTheDocument();
    expect(screen.getByTestId('expert-plugin-chain')).toBeInTheDocument();
    expect(screen.getByTestId('expert-raw-manifest')).toBeInTheDocument();
    expect(screen.getByTestId('expert-payload-preview')).toBeInTheDocument();
  });

  it('adds a route row when clicking add route', () => {
    renderComponent();

    const before = screen.getAllByText('刪除').length;
    fireEvent.click(screen.getByText('+ 新增路由'));
    const after = screen.getAllByText('刪除').length;

    expect(after).toBe(before + 1);
  });

  it('shows manifest parse error when raw JSON is invalid', () => {
    renderComponent();

    const textarea = screen.getByTestId('expert-raw-manifest');
    fireEvent.change(textarea, { target: { value: '{ invalid-json }' } });

    expect(screen.getByTestId('expert-manifest-error')).toBeInTheDocument();
  });
});
