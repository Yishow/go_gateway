import { render, screen, waitFor } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { MemoryRouter, useLocation } from 'react-router-dom'
import { AppRoutes } from '../../src/App'
import { ErrorBoundary } from '../../src/components/ErrorBoundary'
import i18n from '../../src/i18n/config'

type LazyImportControl = {
  pending: Promise<void>
  rejection: Error | null
}

const controls = vi.hoisted(() => ({
  workbenchV2: {
    pending: Promise.resolve(),
    rejection: null,
  } as LazyImportControl,
  runtime: {
    pending: Promise.resolve(),
    rejection: null,
  } as LazyImportControl,
}))

vi.mock('../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page', async () => {
  await controls.workbenchV2.pending
  if (controls.workbenchV2.rejection) {
    throw controls.workbenchV2.rejection
  }

  return {
    default: function WorkbenchV2RouteMock({
      runtimeReturnFocus,
    }: {
      runtimeReturnFocus?: { step: number; issueCode: string | null } | null
    }) {
      return (
        <div
          data-testid="workbench-v2-resolved"
          data-focus-step={runtimeReturnFocus?.step ?? ''}
        >
          Studio V2 route
        </div>
      )
    },
  }
})

vi.mock('../../src/features/datalink/runtime-dashboard/RuntimeDashboardRoute', async () => {
  await controls.runtime.pending
  if (controls.runtime.rejection) {
    throw controls.runtime.rejection
  }

  return {
    RuntimeDashboardRoute: function RuntimeDashboardRouteMock() {
      return <div data-testid="runtime-route-resolved">Runtime route</div>
    },
  }
})

function LocationProbe() {
  const location = useLocation()
  return (
    <output data-testid="location">
      {location.pathname}
      {location.search}
      {location.hash}
    </output>
  )
}

function renderApp(initialEntry: string) {
  return render(
    <ErrorBoundary>
      <MemoryRouter initialEntries={[initialEntry]}>
        <AppRoutes />
        <LocationProbe />
      </MemoryRouter>
    </ErrorBoundary>,
  )
}

function resetControls() {
  controls.workbenchV2.pending = Promise.resolve()
  controls.workbenchV2.rejection = null
  controls.runtime.pending = Promise.resolve()
  controls.runtime.rejection = null
}

describe('AppRoutes lazy-load contract', () => {
  beforeEach(async () => {
    resetControls()
    await i18n.changeLanguage('en')
    vi.spyOn(console, 'error').mockImplementation(() => undefined)
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('uses the same generic fallback for /studio and an arbitrary unknown route', async () => {
    const studioRender = renderApp('/studio?step=output#legacy')

    await waitFor(() => {
      expect(screen.getByTestId('location')).toHaveTextContent('/studio/v2?step=output#legacy')
    })
    expect(await screen.findByTestId('workbench-v2-resolved')).toBeInTheDocument()
    studioRender.unmount()

    renderApp('/unknown-route-for-retirement?step=output#legacy')

    await waitFor(() => {
      expect(screen.getByTestId('location')).toHaveTextContent('/studio/v2?step=output#legacy')
    })
    expect(await screen.findByTestId('workbench-v2-resolved')).toBeInTheDocument()
  })

  it('renders the selected route after its lazy import resolves', async () => {
    renderApp('/studio/v2')

    expect(await screen.findByTestId('workbench-v2-resolved')).toHaveTextContent(
      'Studio V2 route',
    )
    expect(screen.queryByText(i18n.t('common.loading'))).not.toBeInTheDocument()
  })

  it('surfaces a rejected lazy import through the existing error boundary', async () => {
    controls.runtime.rejection = new Error('runtime route chunk failed')

    renderApp('/studio/runtime?device_id=device-A')

    expect(await screen.findByText('發生錯誤')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '重新載入' })).toBeInTheDocument()
    expect(screen.queryByTestId('runtime-route-resolved')).not.toBeInTheDocument()
  })

  it('preserves the root redirect identity while loading the redirected route', async () => {
    renderApp('/')

    await waitFor(() => {
      expect(screen.getByTestId('location')).toHaveTextContent('/studio/v2')
    })
    expect(await screen.findByTestId('workbench-v2-resolved')).toBeInTheDocument()
  })
})
