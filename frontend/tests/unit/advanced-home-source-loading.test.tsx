import { Suspense, useEffect } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react'
import { afterEach, expect, it, vi } from 'vitest'
import DatalinkWorkbenchV2Page from '../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page'
import type { HomeSourceControls } from '../../src/features/datalink/workbench-v2/recording-home/RecordingHomeSourceHost'
const control = vi.hoisted(() => ({ started: false, settle: undefined as (() => void) | undefined, mounted: vi.fn() }))
vi.mock('../../src/pages/datalink/workbench-v2/AdvancedWorkbenchV2Page', async () => {
  control.started = true
  await new Promise<void>(resolve => { control.settle = resolve })
  return { default: function Source({ homeSource }: { homeSource: HomeSourceControls }) {
    useEffect(() => { control.mounted(homeSource.target.epoch) }, [homeSource.target.epoch])
    return <input aria-label="Source draft" defaultValue={`epoch ${homeSource.target.epoch}`} />
  } }
})
vi.mock('../../src/features/datalink/workbench-v2/recording-home/RecordingHomeRoute', () => ({
  RecordingHomeRoute: ({ onWorkspaceChange, onCreateDevice }: { onWorkspaceChange: (id: string) => void; onCreateDevice: (id: string) => void }) => {
    useEffect(() => { onWorkspaceChange('W') }, [onWorkspaceChange])
    return <main data-testid="reading-frame"><h1>Recordings</h1><button onClick={() => onCreateDevice('W')}>Create source</button></main>
  },
}))
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }))
afterEach(cleanup)
it('keeps home visible and source close usable during import, then mounts only the current epoch', async () => {
  const client = new QueryClient()
  render(<QueryClientProvider client={client}><Suspense fallback={<p>Outer loading</p>}><DatalinkWorkbenchV2Page presentation="recordings" /></Suspense></QueryClientProvider>)
  expect(control.started).toBe(false)
  fireEvent.click(screen.getByRole('button', { name: 'Create source' }))
  await waitFor(() => expect(control.started).toBe(true))
  expect(screen.getByTestId('reading-frame')).toBeVisible()
  expect(screen.getByTestId('advanced-load-pending')).toBeInTheDocument()
  expect(control.mounted).not.toHaveBeenCalled()
  fireEvent.click(screen.getByTestId('home-source-close'))
  expect(screen.queryByRole('dialog')).not.toBeInTheDocument()
  fireEvent.click(screen.getByRole('button', { name: 'Create source' }))
  await act(async () => { control.settle?.() })
  expect(await screen.findByRole('textbox', { name: 'Source draft' })).toHaveValue('epoch 2')
  expect(control.mounted).toHaveBeenCalledExactlyOnceWith(2)
  client.clear()
})
