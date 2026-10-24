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
  assert.equal(theme('--radius-sm'), '0.125rem');
  assert.ok(declarations.some(decl => decl.parent.selector === '.shadow-sm' &&
    decl.prop === '--tw-shadow' && decl.value.includes('0 1px 2px 0')));
  assert.ok(declarations.some(decl => decl.prop === 'border-color' &&
    decl.value === '#e2e8f0' && decl.parent.selector.includes('::before')));
  assert.ok(declarations.some(decl => decl.prop === 'outline' && decl.value === '2px solid transparent'));
  assert.ok(result.css.includes('.dark'));
  assert.ok(declarations.some(decl => decl.prop === 'grid-template-columns' && decl.value === 'auto 1fr'));
  assert.ok(!declarations.some(decl => decl.prop === 'grid-template-columns' && decl.value === 'auto,1fr'));
});
