import { render, screen, waitFor } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'

describe('DevAgentation', () => {
  it('loads and renders agentation in development mode', async () => {
    const { DevAgentation } = await import('../../../src/components/DevAgentation')
    const loadAgentation = vi.fn(async () => ({
      Agentation: () => <div data-testid="agentation-overlay" />,
    }))

    render(<DevAgentation isDev loadAgentation={loadAgentation} />)

    await waitFor(() => expect(screen.getByTestId('agentation-overlay')).toBeInTheDocument())
    expect(loadAgentation).toHaveBeenCalledTimes(1)
  })

  it('does not load agentation outside development mode', async () => {
    const { DevAgentation } = await import('../../../src/components/DevAgentation')
    const loadAgentation = vi.fn()

    const { container } = render(<DevAgentation isDev={false} loadAgentation={loadAgentation} />)

    expect(container).toBeEmptyDOMElement()
    expect(loadAgentation).not.toHaveBeenCalled()
  })
})
