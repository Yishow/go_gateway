// @vitest-environment node
import { describe, expect, it } from 'vitest';
import { assertDisabledSources, assertSeedOrigin, type SeedDevice, type SeedRule, type SeedPoint } from '../../e2e/helpers/recording-home-seed-contract';
function fixture() {
  const devices: SeedDevice[] = Array.from({ length: 20 }, (_, i) => ({ id: `perf-device-${i + 1}`, name: `Synthetic offline ${i + 1}`, protocol: 'modbus_tcp', status: 'draft', running: false, availability_status: 'available', last_test_at: null, last_test_success: null }));
  const rules: SeedRule[] = devices.map((d, i) => ({ id: `perf-rule-${i + 1}`, device_id: d.id, workspace_id: 'workspace', enabled: false, share_enabled: false, count: 100, revision_id: `rev-${i}`, save_state: 'saved' }));
  const points: SeedPoint[] = devices.flatMap(d => Array.from({ length: 100 }, (_, i) => ({ id: `${d.id}-point-${i}`, device_id: d.id, enabled: false, last_read_at: null, last_value: null })));
  return { devices, rules, points };
}
describe('disabled 20 by 100 fixture contract', () => {
  it('accepts the exact disabled canonical fixture', () => { const f = fixture(); expect(() => assertDisabledSources('workspace', f.devices, f.rules, f.points)).not.toThrow(); });
  it('accepts the production aligned projection ACK without treating it as collection running', () => { const f = fixture(); f.rules[0].runtime_apply_status = 'aligned'; expect(() => assertDisabledSources('workspace', f.devices, f.rules, f.points)).not.toThrow(); });
  it.each(['device-active', 'device-running', 'device-tested', 'rule-enabled', 'rule-share', 'rule-foreign', 'rule-count', 'rule-ack', 'point-enabled', 'point-read', 'duplicate', 'missing', 'point-foreign'])('rejects %s rather than measuring an unsafe or incomplete fixture', change => {
    const f = fixture();
    if (change === 'device-active') f.devices[0].status = 'active';
    if (change === 'device-running') f.devices[0].running = true;
    if (change === 'device-tested') f.devices[0].last_test_success = true;
    if (change === 'rule-enabled') f.rules[0].enabled = true;
    if (change === 'rule-share') f.rules[0].share_enabled = true;
    if (change === 'rule-foreign') f.rules[0].workspace_id = 'other';
    if (change === 'rule-count') f.rules[0].count = 99;
    if (change === 'rule-ack') f.rules[0].save_state = 'pending';
    if (change === 'point-enabled') f.points[0].enabled = true;
    if (change === 'point-read') f.points[0].last_read_at = 'now';
    if (change === 'duplicate') f.points[0].id = f.points[1].id;
    if (change === 'missing') f.points.pop();
    if (change === 'point-foreign') f.points[0].device_id = 'other';
    expect(() => assertDisabledSources('workspace', f.devices, f.rules, f.points)).toThrow();
  });
  it('accepts only the exact loopback origin returned by the owner', () => { expect(() => assertSeedOrigin('http://127.0.0.1:4567', 4567)).not.toThrow(); });
  it.each(['https://example.com:4567', 'http://localhost:4567', 'http://127.0.0.1:4568', 'http://user:password@127.0.0.1:4567', 'http://127.0.0.1:4567/path', 'http://127.0.0.1:4567?next=other'])('refuses arbitrary origin %s', origin => { expect(() => assertSeedOrigin(origin, 4567)).toThrow(); });
});
