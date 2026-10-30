import { cleanup, render, screen } from '@testing-library/react';
import { afterEach, expect, it } from 'vitest';
import { MemoryRouter } from 'react-router-dom';
import { AppRoutes } from '@/App';
import '@/i18n/config';
afterEach(cleanup);
it('exposes only the dedicated development route in an isolated document', async () => {
  render(<MemoryRouter initialEntries={['/__dev/recording-style']}><AppRoutes /></MemoryRouter>);
  const frame = await screen.findByTitle('Recording style preview — simulated');
  expect(frame).toHaveAttribute('src', '/src/dev/recording-style.html');
  expect(screen.queryByRole('navigation')).not.toBeInTheDocument();
});
