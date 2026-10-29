import { describe, expect, it } from 'vitest';
import { parseRuntimeStatusRecord, parseRuntimeValueRecord, parseRuntimeStreamStateRecord } from '@/utils/safeJson';
const value = { device_id: 'device-A', point_id: 'point-1', address: '40001', timestamp: '2026-10-28T00:00:00Z', raw_value: 12, transformed_value: 'Running', quality: 'good', stale: false };
const status = { device_id: 'device-A', status: 'error', points_total: 1, points_healthy: 0, points_stale: 0, points_error: 1, breaker_state: 'open', last_read_at: null };
describe('remote runtime stream safe projections', () => {
  it.each(['postgres://user:pass@db:5432', 'password=secret', '/private/tenant/data.sqlite', 'C:\\private\\data.sqlite', 'SELECT secret FROM credentials', '10.20.30.40:502', 'db.internal:5432', '~/private/db.sqlite', 'secret=fixture', 'CREATE TABLE records (value int)', 'GRANT select ON records', 'uid=operator;pwd=secret', 'file:private.sqlite', 'C:/private/fixture.sqlite', 'C:/private/fixture'])(
    'does not return unsafe scalar %s', (secret) => {
      expect(parseRuntimeValueRecord({ ...value, raw_value: secret })).toBeNull();
      expect(parseRuntimeStatusRecord({ ...status, last_error: secret })?.last_error).toBe('runtime-error');
    },
  );
  it.each([
    'path=/private/db', 'path=C:/private/db', '(/private/db)', '"/private/db"', "'/private/db'",
    '[/private/db]', '{/private/db}', '\\\\server\\share\\private', './private/db', '../private/db', '~/private/db',
    'path=./private/db', 'path=../private/db', 'path=~/private/db', 'path="C:/private/db"', 'path=\\\\server\\share\\private',
    '/', './', '../', '~/', 'path=./', 'path=../', 'path=~/', 'C:/',
    '{"path":"/private/db"}', '`/private/db`', '(../private/db)', '[~/private/db]',
    '.\\private\\db', '..\\private\\db', 'path:/private/db', 'x)/private/db', 'x]/private/db', 'x}/private/db',
  ])('rejects delimited or explicitly relative private scalar path %s', (path) => {
    expect(parseRuntimeValueRecord({ ...value, raw_value: path })).toBeNull();
    expect(parseRuntimeValueRecord({ ...value, transformed_value: path })).toBeNull();
  });
  it('preserves ordinary text/numbers and register addresses while dropping nested payloads', () => {
    expect(parseRuntimeValueRecord({ ...value, credentials: { password: 'secret' } })).toEqual(value);
    expect(parseRuntimeValueRecord({ ...value, raw_value: { password: 'secret' } })).toBeNull();
  });
  it.each(['SELECT', 'OK', 'NG', '2026-10-28', 'm/s', '40001/16', '9007199254740993', 'C temperature', 'C: temperature', 'OK/NG', 'Pump A (°C)', 'Reading [A]', 'Rate=12 m/s', '23.5 °C', '(normal reading)', '[normal]', '\\\\sensor'])('preserves ordinary scalar %s', (scalar) => {
    expect(parseRuntimeValueRecord({ ...value, raw_value: scalar })?.raw_value).toBe(scalar);
  });
  it.each(['D0100', '0xFFFF', '400001', 'R100.2', 'BFF', 'SB1A', 'W1A', 'W7FFFFFFF', 'D0000000', 'W000000001A', '00040001', '0x000000FFFF', 'D1000000', 'D16777215', 'M1000000'])('preserves validated source register %s', (address) => {
    expect(parseRuntimeValueRecord({ ...value, address })?.address).toBe(address);
  });
  it.each(['W80000000', 'WFFFFFFFF', 'W00080000000', 'SB000FFFFFFFF', 'D16777216', 'HOST1000000', 'ſ1000000', 'D0000000:502', `D${'0'.repeat(128)}`])('rejects out-of-range or unsafe register %s', (address) => {
    expect(parseRuntimeValueRecord({ ...value, address })).toBeNull();
  });
  it('retains an honest unsupported event with invalid quality and no fabricated scalar', () => {
    const unsupported = { ...value, raw_value: null, transformed_value: null, quality: 'invalid', value_state: 'unsupported' };
    expect(parseRuntimeValueRecord(unsupported)).toEqual(unsupported);
  });
  it('preserves stream freshness while replacing untrusted reason text', () => {
    const stream_state = { state: 'stale', empty: false, degraded: false, unavailable: false, stale: true, reason: 'driver SELECT secret FROM credentials' };
    expect(parseRuntimeStreamStateRecord({ stream_state })?.stream_state).toEqual({ ...stream_state, reason: 'unavailable' });
  });
  it('rejects hostile stream identity, state and timestamps instead of retaining arbitrary text', () => {
    const stream_state = { state: 'ready', empty: false, degraded: false, unavailable: false, stale: false };
    for (const payload of [{ stream_state, device_id: 'http://private-host' }, { stream_state, timestamp: '/private/db.sqlite' }, { stream_state: { ...stream_state, state: '/private/db.sqlite' } }]) {
      expect(parseRuntimeStreamStateRecord(payload)).toBeNull();
    }
  });
  it('sanitizes runtime-only envelope IDs and rejects DNS-like entity identities', () => {
    const stream_state = { state: 'ready', empty: false, degraded: false, unavailable: false, stale: false };
    expect(parseRuntimeStreamStateRecord({ stream_state, request_id: 'private.example', operation_id: 'file:private' })).toEqual({ stream_state });
    expect(parseRuntimeValueRecord({ ...value, device_id: 'private.example' })).toBeNull();
    expect(parseRuntimeStatusRecord({ ...status, device_id: 'private.example' })).toBeNull();
  });
  it('does not let nominal status fields tunnel exceptions or endpoints', () => {
    expect(parseRuntimeStatusRecord({ ...status, last_read_at: '/private/db.sqlite' })).toBeNull();
    expect(parseRuntimeStatusRecord({ ...status, breaker_state: '/private/db.sqlite' })).toBeNull();
    expect(parseRuntimeStatusRecord({ ...status, device_id: 'postgres://user:secret@db' })).toBeNull();
    expect(parseRuntimeStatusRecord({ ...status, device_id: '10.20.30.40' })).toBeNull();
    expect(parseRuntimeValueRecord({ ...value, address: 'http://private-host' })).toBeNull();
    expect(parseRuntimeValueRecord({ ...value, timestamp: '/private/db.sqlite' })).toBeNull();
  });
});
