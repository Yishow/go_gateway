import { describe, expect, it } from 'vitest';
import { parseRuntimeSourcePage } from '@/services/runtimeSourcePageParser';
const scope = { workspace_id: 'W', workspace_revision: 'W1', device_id: 'A' };
const request = { kind: 'rules' as const };
const unavailable = { ok: false, reason: 'unavailable', action: 'refresh' };
const stale = { ok: false, reason: 'stale-reference', action: 'refresh' };
const foreign = { ok: false, reason: 'foreign-scope', action: 'open-workspace-list' };
const rule = (id = 'A6') => ({ id, device_id: 'A', start_address: '40001', count: 1, data_type: 'float32', naming_prefix: '電壓', enabled: true, revision_id: 'R6' });
const mapping = () => ({ id: 'M', rule_id: 'A6', rule_revision: 'R6', device_id: 'A', point_id: 'P', address: '40001', tag_id: 'T', tag_key: 'meter.voltage', display_name: '電壓', unit: 'V', target_type: 'float32', enabled: true, status: 'active' });
const page = () => ({ contract_version: 1, ...scope, kind: 'rules', query: '', limit: 50, list_revision: 'L1', rules: [rule()], mappings: [], next_cursor: null });
describe('closed bounded source page parser', () => {
  it('returns detached safe summaries with a required parent revision', () => {
    const input = page(); const result = parseRuntimeSourcePage(input, scope, request);
    expect(result).toEqual({ ok: true, value: input });
    input.rules[0].revision_id = 'later';
    if (!result.ok) throw new Error('expected safe page');
    expect(result.value.rules[0].revision_id).toBe('R6');
  });
  it('supports exact-filter mapping rows without joining another rule page', () => {
    const input = { ...page(), kind: 'mappings', rules: [], mappings: [mapping()] };
    expect(parseRuntimeSourcePage(input, scope, { kind: 'mappings', rule_id: 'A6', rule_revision: 'R6' })).toEqual({ ok: true, value: input });
    input.mappings[0].rule_revision = 'R7';
    expect(parseRuntimeSourcePage(input, scope, { kind: 'mappings', rule_id: 'A6', rule_revision: 'R6' })).toEqual(stale);
  });
  it('allows absent optional mapping fields for an unbound point', () => {
    const input = { ...page(), kind: 'mappings', rules: [], mappings: [{ rule_id: 'A6', rule_revision: 'R6', device_id: 'A', point_id: 'P', address: '40001', enabled: false }] };
    expect(parseRuntimeSourcePage(input, scope, { kind: 'mappings' })).toEqual({ ok: true, value: input });
  });
  it.each(['workspace_id', 'device_id'] as const)('rejects foreign page %s', (key) => {
    expect(parseRuntimeSourcePage({ ...page(), [key]: 'foreign' }, scope, request)).toEqual(foreign);
  });
  it.each(['workspace_revision', 'kind', 'query', 'limit'])('rejects mismatched page %s', (key) => {
    expect(parseRuntimeSourcePage({ ...page(), [key]: key === 'limit' ? 20 : key === 'kind' ? 'mappings' : 'other' }, scope, request)).toEqual(stale);
  });
  it.each(['unknown', '', null, undefined])('rejects unknown or missing list kind %s without reclassifying it as a stale known page', (kind) => {
    expect(parseRuntimeSourcePage({ ...page(), kind }, scope, request)).toEqual(unavailable);
  });
  it('rejects foreign device rows and mixed parent revisions', () => {
    const input = page(); input.rules[0].device_id = 'B';
    expect(parseRuntimeSourcePage(input, scope, request)).toEqual(foreign);
    const map1 = mapping(); const map2 = { ...mapping(), id: 'M2', point_id: 'P2', rule_revision: 'R7' };
    expect(parseRuntimeSourcePage({ ...page(), kind: 'mappings', rules: [], mappings: [map1, map2] }, scope, { kind: 'mappings' })).toEqual(stale);
  });
  it.each(['revision_id', 'device_id', 'id'])('requires safe rule %s', (field) => {
    const row = rule(); Reflect.deleteProperty(row, field);
    expect(parseRuntimeSourcePage({ ...page(), rules: [row] }, scope, request)).toEqual(unavailable);
  });
  it('requires mapping parent revision and rejects foreign rule filters', () => {
    const row = mapping(); Reflect.deleteProperty(row, 'rule_revision');
    expect(parseRuntimeSourcePage({ ...page(), kind: 'mappings', rules: [], mappings: [row] }, scope, { kind: 'mappings' })).toEqual(unavailable);
    expect(parseRuntimeSourcePage({ ...page(), kind: 'mappings', rules: [], mappings: [mapping()] }, scope, { kind: 'mappings', rule_id: 'other', rule_revision: 'R6' })).toEqual(stale);
  });
  it.each(['params', 'password', 'rules', 'mappings'])('rejects nested unexpected rule field %s', (field) => {
    const row = { ...rule(), [field]: [{ password: '/private/secret' }] };
    const result = parseRuntimeSourcePage({ ...page(), rules: [row] }, scope, request);
    expect(result).toEqual(unavailable); expect(JSON.stringify(result)).not.toContain('secret');
  });
  it.each(['naming_prefix', 'start_address', 'revision_id'])('rejects unsafe rule %s', (field) => {
    expect(parseRuntimeSourcePage({ ...page(), rules: [{ ...rule(), [field]: '/private/secret' }] }, scope, request)).toEqual(unavailable);
  });
  it('rejects invalid types, unknown enums and non-finite counts', () => {
    for (const fields of [{ data_type: 'blob' }, { count: Infinity }, { count: 1.5 }, { enabled: 'true' }]) {
      expect(parseRuntimeSourcePage({ ...page(), rules: [{ ...rule(), ...fields }] }, scope, request)).toEqual(unavailable);
    }
  });
  it('enforces kind isolation, page size and duplicate identities', () => {
    expect(parseRuntimeSourcePage({ ...page(), mappings: [mapping()] }, scope, request)).toEqual(unavailable);
    expect(parseRuntimeSourcePage({ ...page(), rules: Array.from({ length: 51 }, (_, index) => rule(`A${index}`)) }, scope, request)).toEqual(unavailable);
    expect(parseRuntimeSourcePage({ ...page(), rules: [rule(), rule()] }, scope, request)).toEqual(unavailable);
  });
  it('accepts at most 100 rows for a clamped request, without accepting hidden remainder arrays', () => {
    const input = { ...page(), limit: 100, rules: Array.from({ length: 100 }, (_, index) => rule(`A${index}`)) };
    expect(parseRuntimeSourcePage(input, scope, { kind: 'rules', limit: 500 })).toMatchObject({ ok: true });
    input.rules.push(rule('A100'));
    expect(parseRuntimeSourcePage(input, scope, { kind: 'rules', limit: 500 })).toEqual(unavailable);
  });
  it.each([null, undefined, {}, { ...page(), extra: true }])('rejects missing or open page schema', (input) => {
    expect(parseRuntimeSourcePage(input, scope, request)).toEqual(unavailable);
  });
  it('does not evaluate accessor payloads', () => {
    const getter = () => { throw new Error('/private/secret'); };
    const input = Object.defineProperty(page(), 'rules', { get: getter });
    expect(parseRuntimeSourcePage(input, scope, request)).toEqual(unavailable);
  });
  it('validates canonical bounded cursor text and refuses a repeated cursor', () => {
    expect(parseRuntimeSourcePage({ ...page(), next_cursor: 'e30' }, scope, request)).toMatchObject({ ok: true });
    for (const next_cursor of ['', 'e30=', 'e31', 'a', 'A'.repeat(1025)]) {
      expect(parseRuntimeSourcePage({ ...page(), next_cursor }, scope, request)).toEqual(unavailable);
    }
    expect(parseRuntimeSourcePage({ ...page(), next_cursor: 'e30' }, scope, { kind: 'rules', cursor: 'e30' })).toEqual(unavailable);
  });
});
