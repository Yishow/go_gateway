import { render } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'

vi.mock('../../../src/pages/TestPage', () => ({
  default: () => <div>test-page</div>,
}))

vi.mock('../../../src/pages/TestPageShell', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}))

describe('App agentation removal', () => {
  it('renders without loading agentation', async () => {
    const { default: App } = await import('../../../src/App')

    expect(() => render(<App />)).not.toThrow()
  })
})
