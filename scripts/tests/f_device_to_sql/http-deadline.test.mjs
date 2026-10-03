import assert from 'node:assert/strict';
import { createServer } from 'node:http';
import test from 'node:test';

test('all fixture HTTP helpers fail within a deadline when a server accepts but never replies', async () => {
  let requestsSeen = 0;
  const server = createServer((request) => { requestsSeen++; request.resume(); });
  await new Promise((resolve) => server.listen(0, '127.0.0.1', resolve));
  const port = String(server.address().port);
  process.env.GW_PORT = port;
  process.env.F_FIXTURE_PORT = port;
  const { control } = await import('./quality-lib.mjs');
  const { mutation } = await import('./recovery-lib.mjs');
  const { request } = await import('./capacity-cas.mjs');
  const { confirm } = await import('./test-write-lib.mjs');
  const pending = [control('/stall'), mutation('/stall', 'POST', {}), request('POST', '/stall', {}), confirm({ id: 'owned-group' }, { token: 'fixture-token', operation_id: 'op-a' })];
  let guard;
  try {
    const replies = await Promise.race([
      Promise.allSettled(pending),
      new Promise((resolve) => { guard = setTimeout(() => resolve('deadline missed'), 35000); }),
    ]);
    assert.equal(requestsSeen, 4, 'all real requests reached the stalled loopback server');
    assert.notEqual(replies, 'deadline missed', 'cleanup must regain control even if the server never responds');
    assert.equal(replies.every((reply) => reply.status === 'rejected' && reply.reason.name === 'TimeoutError'), true);
  } finally {
    clearTimeout(guard);
    server.closeAllConnections();
    await new Promise((resolve) => server.close(resolve));
    await Promise.allSettled(pending);
  }
});
