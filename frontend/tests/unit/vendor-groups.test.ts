import { expect, it } from 'vitest'
import config from '../../vite.config'
it('keeps stable dependency families without pulling application routes into vendor chunks', () => {
  const result = typeof config === 'function' ? config({ command: 'build', mode: 'production' }) : config
  const groups = result.build?.rolldownOptions?.output
  expect(groups).toBeDefined()
  // This small interface deliberately checks public config behavior, not source text.
  const output = groups as { codeSplitting: { includeDependenciesRecursively: boolean; groups: { name: (id: string) => string | null }[] } }
  expect(output.codeSplitting.includeDependenciesRecursively).toBe(false)
  const name = output.codeSplitting.groups[0].name
  for (const [dependency, expected] of [
    ['react/index.js', 'react-vendor'], ['react-dom/client.js', 'react-vendor'], ['scheduler/index.js', 'react-vendor'],
    ['recharts/es6/index.js', 'charts-vendor'], ['d3-scale/src/index.js', 'charts-vendor'],
    ['react-router/dist/index.js', 'app-vendor'], ['i18next/index.js', 'app-vendor'], ['@tanstack/react-query/build/index.js', 'app-vendor'],
    ['axios/index.js', 'misc-vendor'], ['@mui/material/index.js', 'mui-vendor'], ['@emotion/react/index.js', 'mui-vendor'],
  ]) expect(name(`/project/node_modules/${dependency}`)).toBe(expected)
  expect(name('/project/src/pages/TestPage.tsx')).toBeNull()
  expect(result.build?.manifest).toBe(true)
})
