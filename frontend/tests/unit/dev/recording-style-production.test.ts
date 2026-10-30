// @vitest-environment node
import { execFileSync } from 'node:child_process';
import { mkdtempSync, readdirSync, readFileSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import path from 'node:path';
import { expect, it } from 'vitest';

function files(dir: string): string[] {
  return readdirSync(dir, { withFileTypes: true }).flatMap(entry => entry.isDirectory() ? files(path.join(dir, entry.name)) : [path.join(dir, entry.name)]);
}
it('real production build excludes the preview route, iframe entry, fixtures and chunks', () => {
  const output = mkdtempSync(path.join(tmpdir(), 'recording-preview-production-'));
  try {
    execFileSync(process.execPath, ['node_modules/vite/bin/vite.js', 'build', '--mode', 'production', '--outDir', output], { cwd: process.cwd(), env: { ...process.env, NODE_ENV: 'production' }, stdio: 'pipe' });
    const artifacts = files(output);
    expect(artifacts.some(file => file.endsWith('index.html'))).toBe(true);
    const forbidden = /RECORDING_STYLE_PREVIEW_SIMULATED_ONLY|RecordingStylePreview|RecordingStyleFrame|previewI18n|previewTranslationsReady|recording-style(?:-entry)?\.(?:html|tsx)|\/__dev\/recording-style|SIMULATED recording · not a saved resource/;
    for (const file of artifacts) {
      expect(file).not.toMatch(forbidden);
      if (/\.(?:js|css|html|map)$/.test(file)) expect(readFileSync(file, 'utf8')).not.toMatch(forbidden);
    }
  } finally { rmSync(output, { recursive: true, force: true }); }
}, 60_000);

it('operator navigation has no dedicated preview link', () => {
  const references = files(path.join(process.cwd(), 'src')).filter(file => /\.(tsx?|json)$/.test(file))
    .filter(file => !file.includes(`${path.sep}dev${path.sep}`) && !file.endsWith(`${path.sep}App.tsx`))
    .filter(file => readFileSync(file, 'utf8').includes('/__dev/recording-style'));
  expect(references).toEqual([]);
});

// An iframe document starts its own JS realm; do not mount React through a parent portal.
it('uses a separate module entry without application/settings owners', () => {
  const html = readFileSync('src/dev/recording-style.html', 'utf8');
  expect(html).toContain('src="/src/dev/recording-style-entry.tsx"');
  const entry = readFileSync('src/dev/recording-style-entry.tsx', 'utf8');
  expect(entry).toContain('import.meta.env.DEV');
  expect(entry).toContain('createRoot(node)');
  expect(entry).toContain('<I18nextProvider i18n={previewI18n}>');
  expect(entry).not.toContain('../i18n/config');
  expect(entry).not.toMatch(/ThemeProvider|QueryClientProvider|from ['"]\.\.\/App|createPortal/);
});
