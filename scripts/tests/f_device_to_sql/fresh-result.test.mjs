import test from 'node:test';
import assert from 'node:assert/strict';
import { mkdtempSync, readFileSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { spawnSync } from 'node:child_process';

test('the final cleanup verdict controls both the persisted evidence and the child exit', () => {
  const work = mkdtempSync(join(tmpdir(), 'gw-f-result-'));
  try {
    for (const passed of [false, true]) {
      const path = join(work, `result-${passed}.json`);
      const payload = { passed, cleanup: { passed } };
      const code = `import { writeFreshResult } from ${JSON.stringify(new URL('./fresh-result.mjs', import.meta.url).href)};
        process.exitCode = 0;
        writeFreshResult(process.argv[1], JSON.parse(process.argv[2]));`;
      const child = spawnSync(process.execPath, ['--input-type=module', '-e', code, path, JSON.stringify(payload)]);
      assert.equal(child.signal, null);
      assert.deepEqual(JSON.parse(readFileSync(path, 'utf8')), payload);
      assert.equal(child.status, passed ? 0 : 1, 'final cleanup failure must override an earlier successful main flow');
    }
  } finally {
    rmSync(work, { recursive: true });
  }
});
