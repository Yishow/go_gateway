import { describe, expect, it, vi } from 'vitest';
import axios from 'axios';
import { recordingHomeSourcesAPI } from '@/services/recordingHomeSources';
const valid = { workspace_id: 'W', state: 'available', source_digest: 'a'.repeat(64), devices: [], candidates: [], excluded: [] };
function backend(value: unknown) {
  const client = axios.create();
  vi.spyOn(client, 'get').mockResolvedValue({ data: { success: true, data: value } });
  return client;
}
describe('compact home source service', () => {
  it('reads exactly one scoped summary and forwards cancellation', async () => {
    const client = backend(valid), signal = new AbortController().signal;
    expect(await recordingHomeSourcesAPI.read('W', signal, client)).toEqual(valid);
    expect(client.get).toHaveBeenCalledExactlyOnceWith('/studio-v2/recording-home/sources', { params: { workspace_id: 'W' }, signal });
  });
  it.each([null, {}, { ...valid, workspace_id: 'OTHER' }, { ...valid, source_digest: '' }, { ...valid, devices: null },
    { ...valid, state: 'other' }, { ...valid, candidates: [{}] }, { ...valid, excluded: [{ label: 'a', reason: 'made-up' }] },
    { workspace_id: 'W', state: 'unknown', reason_code: 'source-changed', devices: [] },
  ])('rejects malformed or foreign observations %j', async value => {
    await expect(recordingHomeSourcesAPI.read('W', undefined, backend(value))).rejects.toThrow('Source observation unavailable');
  });
  it.each(['source-changed', 'source-inconsistent'])('preserves safe unknown %s', async reason_code => {
    const value = { workspace_id: 'W', state: 'unknown', reason_code };
    expect(await recordingHomeSourcesAPI.read('W', undefined, backend(value))).toEqual(value);
  });
});
it('pins the captured backend base URL even when mutable client defaults change', async () => {
  const client = backend(valid);
  client.defaults.baseURL = 'https://new.example/api';
  await recordingHomeSourcesAPI.read('W', undefined, client, 'https://old.example/api');
  expect(client.get).toHaveBeenCalledExactlyOnceWith('/studio-v2/recording-home/sources', {
    params: { workspace_id: 'W' }, signal: undefined, baseURL: 'https://old.example/api',
  });
});
const device = { id: 'D', name: 'Device', protocol: 'modbus-tcp', availability: 'available' };
const candidate = { key: 'P|T', device_id: 'D', point_id: 'P', tag_id: 'T', tag_key: 'tag', label: 'Point', device_name: 'Device', address: '40001', target_type: 'float32' };
it.each([
  { ...valid, devices: [device], candidates: [candidate, { ...candidate, key: 'P|OTHER', tag_id: 'OTHER' }] },
  { ...valid, devices: [device, device] },
  { ...valid, devices: [device], candidates: [candidate, candidate] },
  ...['', ' \t\n', '\u0085', '\u2003'].map(id => ({ ...valid, devices: [{ ...device, id }] })),
  ...['', ' \t\n', '\u0085', '\u2003'].map(point_id => ({ ...valid, devices: [device], candidates: [{ ...candidate, point_id, key: `${point_id}|T` }] })),
])('rejects duplicate or blank canonical identities %j', async value => {
  await expect(recordingHomeSourcesAPI.read('W', undefined, backend(value))).rejects.toThrow('Source observation unavailable');
});
it.each(['', ' \n', '\u0085'])('rejects blank workspace identity %j', async workspace_id => {
  await expect(recordingHomeSourcesAPI.read(workspace_id, undefined, backend({ ...valid, workspace_id }))).rejects.toThrow('Source observation unavailable');
});
it.each([' D ', '\ufeff'])('preserves nonblank canonical identity verbatim %j', async id => {
  const value = { ...valid, devices: [{ ...device, id }], candidates: [{ ...candidate, device_id: id, point_id: id, key: `${id}|T` }] };
  expect(await recordingHomeSourcesAPI.read('W', undefined, backend(value))).toEqual(value);
});
