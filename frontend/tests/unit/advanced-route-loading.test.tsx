import { readFileSync } from 'node:fs'
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createMemoryRouter, RouterProvider } from 'react-router-dom'
import { afterEach, expect, it, vi } from 'vitest'
import { useState } from 'react'
import type { StudioEntityLocatorIntent } from '../../src/services/studioEntityLocator'
import { useRecordingDraftGuard } from '../../src/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation'
import { RecordingRouterNavigation } from '../../src/features/datalink/workbench-v2/recording-editor/RecordingRouterNavigation'

const control = vi.hoisted(() => ({ started: false, settle: undefined as (() => void) | undefined, failure: false }))
vi.mock('../../src/pages/datalink/workbench-v2/AdvancedWorkbenchV2Page', async () => {
  control.started = true
  await new Promise<void>(resolve => { control.settle = resolve })
  return { default: function Advanced({ presentationLocation, locatorIntent }: { presentationLocation?: { search: string; hash: string }; locatorIntent?: StudioEntityLocatorIntent }) {
    const [draft, setDraft] = useState('')
    useRecordingDraftGuard({ version: draft, dirty: Boolean(draft) })
    if (control.failure) throw new Error('secret endpoint / token must not leak')
    return <><output data-testid="advanced-intent">{presentationLocation?.search}{presentationLocation?.hash}</output>
      <output data-testid="advanced-locator">{JSON.stringify(locatorIntent)}</output><input aria-label="Draft" value={draft} onChange={e => setDraft(e.target.value)} /></>
  } }
})
vi.mock('../../src/hooks/datalink/useStudioV2Workspace', () => ({ useStudioV2WorkspaceQuery: () => ({ isPending: true, isFetching: true }) }))
vi.mock('../../src/pages/datalink/workbench-v2/useStudioV2AutosaveState', () => ({ useStudioV2AutosaveState: () => { throw new Error('advanced owner loaded eagerly') } }))
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }))
const clients: QueryClient[] = []
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); vi.restoreAllMocks() })

it('initial entry does not eagerly require advanced-only editor modules', () => {
  const source = readFileSync('src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx', 'utf8')
  expect(source).not.toMatch(/import .*useStudioV2AutosaveState/)
  expect(source).not.toMatch(/import \{ WorkbenchV2Shell \}/)
  expect(source).toMatch(/import\(['"]\.\/AdvancedWorkbenchV2Page['"]\)/)
})

it('keeps the reading frame before deferred imports and preserves deep-link intent and dirty Back/Forward cancellation', async () => {
  const { AppRoutes } = await import('../../src/App')
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); clients.push(client)
  const router = createMemoryRouter([{ path: '*', element: <QueryClientProvider client={client}><RecordingRouterNavigation><AppRoutes /></RecordingRouterNavigation></QueryClientProvider> }],
    { initialEntries: ['/studio/v2?view=recordings#home', '/studio/v2?entry=locator&workspace_id=W&step=2&device_id=D&rule_id=R#intent', '/studio/v2?view=recordings#next'], initialIndex: 0 })
  render(<RouterProvider router={router} />)
  expect(await screen.findByTestId('recording-home')).toHaveAttribute('aria-busy', 'true')
  expect(control.started).toBe(false)
  await act(async () => { await router.navigate(1) })
  await waitFor(() => expect(control.started).toBe(true))
  expect(screen.getByRole('status')).toBeInTheDocument()
  expect(router.state.location.hash).toBe('#intent')
  await act(async () => { control.settle?.() })
  expect(await screen.findByTestId('advanced-intent')).toHaveTextContent('?entry=locator&workspace_id=W&step=2&device_id=D&rule_id=R#intent')
  expect(screen.getByTestId('advanced-locator')).toHaveTextContent('"kind":"locator"')
  expect(screen.getByTestId('advanced-locator')).toHaveTextContent('"rule_id":"R"')
  fireEvent.change(screen.getByRole('textbox', { name: 'Draft' }), { target: { value: 'unsaved range' } })
  for (const delta of [-1, 1]) {
    await act(async () => { await router.navigate(delta) })
    expect(await screen.findByRole('alertdialog')).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: 'navigation.keep' }))
    await waitFor(() => expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument())
    expect(screen.getByRole('textbox', { name: 'Draft' })).toHaveValue('unsaved range')
    expect(router.state.location.hash).toBe('#intent')
  }
  router.dispose()
})

it('keeps a failed deferred editor actionable without exposing its exception or changing URL', async () => {
  control.failure = true
  vi.spyOn(console, 'error').mockImplementation(() => undefined)
  const { AppRoutes } = await import('../../src/App')
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); clients.push(client)
  const router = createMemoryRouter([{ path: '*', element: <QueryClientProvider client={client}><RecordingRouterNavigation><AppRoutes /></RecordingRouterNavigation></QueryClientProvider> }],
    { initialEntries: ['/studio/v2?view=advanced#retry'] })
  render(<RouterProvider router={router} />)
  expect(await screen.findByRole('alert')).toBeInTheDocument()
  expect(screen.getByRole('button', { name: 'retry' })).toBeInTheDocument()
  expect(screen.queryByText(/secret endpoint/)).not.toBeInTheDocument()
  expect(router.state.location.hash).toBe('#retry')
  router.dispose()
})
