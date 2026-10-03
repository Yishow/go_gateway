import assert from 'node:assert/strict';
import test from 'node:test';
import { responseWitness } from './capacity-cas.mjs';

test('CAS witness preserves operation facts without persisting the confirmation token', () => {
  const reply = {
    request: { method: 'POST', path: '/groups/a/test-write', body: { token: 'fixture-bearer', operation_id: 'op-a' } },
    status: 409, ok: false,
    body: { success: false, data: { token: 'fixture-bearer', operation_id: 'op-a', group_id: 'group-a' } },
  };
  const got = responseWitness(reply);
  assert.equal(JSON.stringify(got).includes('fixture-bearer'), false);
  assert.equal(got.request.body.operation_id, 'op-a');
  assert.equal(got.response.body.data.group_id, 'group-a');
  assert.equal(got.response.status, 409);
  assert.match(got.request.body.token_sha256, /^[a-f0-9]{64}$/);
  assert.equal(got.request.body.token_sha256, got.response.body.data.token_sha256);
  assert.equal(reply.request.body.token, 'fixture-bearer', 'redaction cannot change the actual HTTP request');
});
