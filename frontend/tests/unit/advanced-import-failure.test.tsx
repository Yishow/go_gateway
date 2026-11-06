import { render, screen } from '@testing-library/react'
import { afterEach, expect, it, vi } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { AppRoutes } from '../../src/App'
const control = vi.hoisted(() => ({ attempts: 0 }))
vi.mock('../../src/pages/datalink/workbench-v2/AdvancedWorkbenchV2Page', () => { control.attempts++; return Promise.reject(new Error('private chunk URL and credential')) })
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }))
afterEach(() => vi.restoreAllMocks())
it('a rejected advanced import has a safe explicit recovery action, without automatic retry', async () => {
  vi.spyOn(console, 'error').mockImplementation(() => undefined)
  render(<MemoryRouter initialEntries={['/studio/v2?view=advanced#keep']}><AppRoutes /></MemoryRouter>)
  expect(await screen.findByRole('alert')).toHaveTextContent('error')
  expect(screen.getByRole('button', { name: 'retry' })).toBeInTheDocument()
  expect(screen.queryByText(/private chunk URL/)).not.toBeInTheDocument()
  expect(control.attempts).toBe(1)
})
