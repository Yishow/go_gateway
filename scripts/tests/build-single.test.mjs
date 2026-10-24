import assert from 'node:assert/strict';
import { mkdtempSync, mkdirSync, readFileSync, readdirSync, rmSync, writeFileSync, existsSync } from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { test } from 'node:test';
import { buildSingle, prepareFrontend, resolveTarget } from '../build-single.mjs';

function fixture(t, { dependencies = true } = {}) {
  const root = mkdtempSync(path.join(os.tmpdir(), 'gateway build '));
  t.after(() => rmSync(root, { recursive: true, force: true }));
  function put(file, content) {
    mkdirSync(path.dirname(path.join(root, file)), { recursive: true });
    writeFileSync(path.join(root, file), content);
  }
  put('frontend/package.json', '{}');
  put('frontend/package-lock.json', '{}');
  if (dependencies) {
    put('frontend/node_modules/vite/bin/vite.js', '');
    put('frontend/node_modules/typescript/bin/tsc', '');
  }
  put('cmd/test_ui/static/embed-placeholder.txt', 'tracked marker\n');
  put('cmd/test_ui/static/old/stale.js', 'stale');
  const calls = [];
  function run(command, args, options) {
    calls.push({ command, args, options });
    if (command === 'npm' && args[0] === 'run') {
      put('frontend/dist/index.html', '<script src="/assets/app.js"></script>');
      put('frontend/dist/assets/app.js', 'current JS');
      put('frontend/dist/assets/nested/styles.css', 'current CSS');
      put('frontend/dist/.well-known/example', 'hidden child');
    }
    if (command === 'go') {
      // Go must see the current complete frontend before compiling.
      assert.equal(readFileSync(path.join(root, 'cmd/test_ui/static/assets/app.js'), 'utf8'), 'current JS');
      writeFileSync(args[args.indexOf('-o') + 1], 'compiled binary');
    }
  }
  return { root, put, calls, run };
}

test('six OS/architecture targets compile current embedded assets without a C toolchain', async (t) => {
  for (const osName of ['windows', 'mac', 'linux']) {
    for (const arch of ['amd64', 'arm64']) {
      await t.test(`${osName}/${arch}`, (t) => {
        const f = fixture(t);
        const output = buildSingle({ ...f, os: osName, arch });
        const goos = osName === 'mac' ? 'darwin' : osName;
        assert.equal(path.basename(output), `gateway-${goos}-${arch}${goos === 'windows' ? '.exe' : ''}`);
        assert.equal(readFileSync(output, 'utf8'), 'compiled binary');
        assert.equal(f.calls.at(-1).options.env.GOOS, goos);
        assert.equal(f.calls.at(-1).options.env.GOARCH, arch);
        assert.equal(f.calls.at(-1).options.env.CGO_ENABLED, '0');
        assert.equal(existsSync(path.join(f.root, 'cmd/test_ui/static/old')), false);
        for (const file of ['index.html', 'assets/app.js', 'assets/nested/styles.css', '.well-known/example']) {
          assert.deepEqual(readFileSync(path.join(f.root, 'frontend/dist', file)), readFileSync(path.join(f.root, 'cmd/test_ui/static', file)));
        }
        assert.equal(readFileSync(path.join(f.root, 'cmd/test_ui/static/embed-placeholder.txt'), 'utf8'), 'tracked marker\n');
      });
    }
  }
});

test('missing dependencies use npm ci before the frontend build', (t) => {
  const f = fixture(t, { dependencies: false });
  buildSingle(f);
  assert.deepEqual(f.calls.map(({ command, args }) => [command, args[0]]), [['npm', 'ci'], ['npm', 'run'], ['go', 'build']]);
});

test('TypeScript 7 compiler alias does not reinstall prepared dependencies', (t) => {
  const f = fixture(t);
  rmSync(path.join(f.root, 'frontend/node_modules/typescript/bin/tsc'));
  f.put('frontend/node_modules/@typescript/native/bin/tsc', '');
  prepareFrontend(f);
  assert.deepEqual(f.calls.map(({ command, args }) => [command, args[0]]), [['npm', 'run']]);
});

test('fresh development bootstrap emits the embedded entry without compiling a release binary', (t) => {
  const f = fixture(t);
  prepareFrontend(f);
  assert.equal(existsSync(path.join(f.root, 'cmd/test_ui/static/index.html')), true);
  assert.equal(f.calls.some(c => c.command === 'go'), false);
  assert.equal(existsSync(path.join(f.root, 'bin')), false);
});

test('failed install or frontend build cannot invoke Go or replace embedded assets', async (t) => {
  for (const step of ['ci', 'run']) {
    await t.test(step, (t) => {
      const f = fixture(t, { dependencies: step !== 'ci' });
      const run = (command, args, options) => {
        if (command === 'npm' && args[0] === step) throw new Error('fixture failure');
        f.run(command, args, options);
      };
      assert.throws(() => buildSingle({ ...f, run }), /fixture failure/);
      assert.equal(f.calls.some(c => c.command === 'go'), false);
      assert.equal(readFileSync(path.join(f.root, 'cmd/test_ui/static/old/stale.js'), 'utf8'), 'stale');
    });
  }
});

test('missing or empty new index cannot reuse stale dist or invoke Go', async (t) => {
  for (const empty of [false, true]) {
    await t.test(`empty=${empty}`, (t) => {
      const f = fixture(t);
      f.put('frontend/dist/index.html', 'stale index');
      let invokedGo = false;
      const run = (command) => {
        if (command === 'go') invokedGo = true;
        if (command === 'npm' && empty) f.put('frontend/dist/index.html', '');
      };
      assert.throws(() => buildSingle({ ...f, run }), /index.html/);
      assert.equal(invokedGo, false);
    });
  }
});

test('missing source and unusable static target stop before Go', async (t) => {
  for (const missingSource of [true, false]) {
    await t.test(`missing source=${missingSource}`, (t) => {
      const f = fixture(t);
      if (missingSource) rmSync(path.join(f.root, 'frontend'), { recursive: true });
      else {
        rmSync(path.join(f.root, 'cmd/test_ui/static'), { recursive: true });
        f.put('cmd/test_ui/static', 'cannot be a directory');
      }
      assert.throws(() => buildSingle(f));
      assert.equal(f.calls.some(c => c.command === 'go'), false);
    });
  }
});

test('failed Go build preserves the previous binary and removes pending output', (t) => {
  const f = fixture(t);
  f.put('bin/gateway-linux-amd64', 'previous binary');
  const run = (command, args, options) => {
    f.run(command, args, options);
    if (command === 'go') throw new Error('compile failed');
  };
  assert.throws(() => buildSingle({ ...f, run, os: 'linux', arch: 'amd64' }), /compile failed/);
  assert.equal(readFileSync(path.join(f.root, 'bin/gateway-linux-amd64'), 'utf8'), 'previous binary');
  assert.deepEqual(readdirSync(path.join(f.root, 'bin')), ['gateway-linux-amd64']);
});

test('invalid targets are rejected before any build command', (t) => {
  const f = fixture(t);
  assert.throws(() => buildSingle({ ...f, os: '../linux' }), /平台/);
  assert.throws(() => buildSingle({ ...f, arch: 'x86' }), /架構/);
  assert.equal(f.calls.length, 0);
  assert.deepEqual(resolveTarget('win32', 'x64'), { targetOS: 'windows', targetArch: 'amd64' });
});
