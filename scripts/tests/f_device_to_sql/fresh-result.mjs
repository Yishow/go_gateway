import { writeFileSync } from 'node:fs';

/** Persist the final result after cleanup has resolved its verdict. */
export function writeFreshResult(path, payload) {
  writeFileSync(path, `${JSON.stringify(payload, null, 2)}\n`);
  process.exitCode = payload.passed === true ? 0 : 1;
}
