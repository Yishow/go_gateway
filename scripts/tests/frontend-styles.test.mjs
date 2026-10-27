import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import { createRequire } from 'node:module';
import { fileURLToPath } from 'node:url';
import { test } from 'node:test';

const frontend = new URL('../../frontend/', import.meta.url);
const require = createRequire(new URL('package.json', frontend));
const postcss = require('postcss');
const config = require(fileURLToPath(new URL('postcss.config.cjs', frontend)));
const from = fileURLToPath(new URL('src/index.css', frontend));

test('application CSS compiles with the established theme and usable grid columns', async () => {
  const plugins = Object.entries(config.plugins).map(([name, options]) => require(name)(options));
  const result = await postcss(plugins).process(readFileSync(from, 'utf8'), { from });
  const declarations = [];
  result.root.walkDecls(decl => declarations.push(decl));
  const theme = name => declarations.filter(decl => decl.prop === name).at(-1)?.value;

  assert.ok(declarations.some(decl => decl.parent.selector === '.bg-gray-200' && decl.value === '#e2e8f0'));
  assert.equal(theme('--blur-sm'), '4px');
  assert.match(theme('--radius-sm'), /^0?\.125rem$/);
  assert.ok(declarations.some(decl => decl.parent.selector === '.shadow-sm' &&
    decl.prop === '--tw-shadow' && decl.value.includes('0 1px 2px 0')));
  assert.ok(declarations.some(decl => decl.prop === 'border-color' &&
    decl.value === '#e2e8f0' && /::?before/.test(decl.parent.selector)));
  assert.ok(declarations.some(decl => decl.prop === 'outline' &&
    ['2px solid transparent', '2px solid #0000'].includes(decl.value)));
  assert.ok(result.css.includes('.dark'));
  assert.ok(declarations.some(decl => decl.prop === 'grid-template-columns' && decl.value === 'auto 1fr'));
  assert.ok(!declarations.some(decl => decl.prop === 'grid-template-columns' && decl.value === 'auto,1fr'));
});

test('Vite loads application CSS without missing-source warnings', async t => {
  const { createLogger, createServer } = await import(require.resolve('vite'));
  const warnings = [];
  const logger = createLogger('silent');
  logger.warn = logger.warnOnce = message => warnings.push(message);
  const server = await createServer({
    root: fileURLToPath(frontend),
    customLogger: logger,
    server: { middlewareMode: true, hmr: false },
  });
  t.after(() => server.close());

  const result = await server.transformRequest('/src/index.css');
  assert.ok(result?.code.includes('@keyframes fadeIn'));
  assert.ok(result?.code.includes('@keyframes slideUp'));
  const font = await server.transformRequest('/node_modules/@fontsource/inter/400.css');
  assert.ok(font?.code.includes('/node_modules/@fontsource/inter/files/inter-latin-400-normal.woff2'));
  assert.deepEqual(warnings.filter(message => message.includes('did not pass the `from` option')), []);
});
