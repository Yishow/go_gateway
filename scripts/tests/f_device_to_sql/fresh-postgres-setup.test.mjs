import test from 'node:test';
import assert from 'node:assert/strict';
import { createOwnedPostgresRequestGuard, inspectPostgresConfigRequest } from './fresh-postgres-setup.mjs';

const target = {
  host: '127.0.0.1', port: 55432, database: 'gwtest', username: 'postgres',
  schema: 'gw_f_owned', table: 'f_connector_placeholder', password: 'fixture-only',
};

function request(method, url, body) {
  return {
    method: () => method,
    url: () => url,
    postDataJSON: () => body,
  };
}

function ownedBody(overrides = {}) {
  return { kind: 'postgres', ...target, ...overrides };
}

test('owned guard continues GET and complete owned mutation without changing payload', () => {
  assert.deepEqual(inspectPostgresConfigRequest(request('GET', 'http://gw/api/v1/datalink/studio-v2/workspace/database-config'), null), {
    action: 'continue', blockedFields: [],
  });
  assert.deepEqual(inspectPostgresConfigRequest(request('PUT', 'http://gw/api/v1/datalink/studio-v2/workspace/database-config', ownedBody()), target), {
    action: 'continue', blockedFields: [],
  });
  assert.deepEqual(inspectPostgresConfigRequest(request('POST', 'http://gw/api/v1/datalink/studio-v2/workspace/database-config?probe=1', ownedBody()), target), {
    action: 'continue', blockedFields: [],
  });
});

test('owned guard aborts default endpoint and records only safe field names', () => {
  const guard = createOwnedPostgresRequestGuard(target);
  const decision = guard.inspect(request('PUT', 'http://gw/api/v1/datalink/studio-v2/workspace/database-config', ownedBody({
    port: 5432, database: 'gateway_metrics', schema: 'public', table: 'sensor_readings',
  })));
  assert.equal(decision.action, 'abort');
  assert.deepEqual(decision.blockedFields, ['port', 'database', 'schema', 'table']);
  assert.deepEqual(guard.snapshot(), {
    blocked_count: 1, blocked_fields: ['database', 'port', 'schema', 'table'],
  });
  assert.equal(JSON.stringify(guard.snapshot()).includes('fixture-only'), false);
});

test('guard leaves unrelated mutations alone and fails closed on malformed owned body', () => {
  assert.deepEqual(inspectPostgresConfigRequest(request('POST', 'http://gw/api/v1/datalink/studio-v2/workspace/devices', {}), target), {
    action: 'continue', blockedFields: [],
  });
  assert.deepEqual(inspectPostgresConfigRequest(request('POST', 'http://gw/api/v1/datalink/studio-v2/workspace/database-config', null), target), {
    action: 'abort', blockedFields: ['body'],
  });
});

test('guard checks username, kind, schema, and table as part of ownership', () => {
  for (const [field, value] of [['username', 'other'], ['kind', 'mysql'], ['schema', 'wrong'], ['table', 'wrong']]) {
    const decision = inspectPostgresConfigRequest(
      request('PUT', 'http://gw/api/v1/datalink/studio-v2/workspace/database-config', ownedBody({ [field]: value })), target,
    );
    assert.equal(decision.action, 'abort');
    assert.deepEqual(decision.blockedFields, [field]);
  }
});
