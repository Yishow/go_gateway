import { describe, expect, it } from 'vitest';
import { recordingCollection, deviceCollection, initialCollectionState, type HomeDeviceSummary } from '@/features/datalink/workbench-v2/recording-home/recordingHomeCollections';
import { group } from './recording-home-guidance-fixtures';
const devices: HomeDeviceSummary[] = Array.from({ length: 100 }, (_, index) => ({ id: `D${index}`, name: `設備 ${index}`, protocol: 'mc_3e', availability: 'available' }));
const groups = Array.from({ length: 300 }, (_, index) => ({ ...group(`G${String(index).padStart(3, '0')}`, index % 3 === 0 ? 'draft' : 'running'),
  name: `記錄 ${String(index).padStart(3, '0')}`, members: [{ ...group('G').members[0], device_id: `D${index % 100}` }],
  destination: { ...group('G').destination, connector_id: index < 150 ? 'DB1' : 'DB2', table_name: `table_${index}` } }));
describe('recording collection is a bounded client projection', () => {
  it('renders only25 or50 rows from300 records without mutating canonical input', () => {
    const before = JSON.stringify(groups);
    expect(recordingCollection(groups, devices, initialCollectionState()).rows).toHaveLength(25);
    expect(recordingCollection(groups, devices, { ...initialCollectionState(), pageSize: 50 }).rows).toHaveLength(50);
    expect(JSON.stringify(groups)).toBe(before);
  });
  it.each(['記錄 299', 'table_299', '設備 99'])('finds labels across recording, table and device: %s', keyword => {
    expect(recordingCollection(groups, devices, { ...initialCollectionState(), keyword }).rows.map(row => row.id)).toContain('G299');
  });
  it('combines device, persisted status and destination filters, not health guesses', () => {
    const result = recordingCollection(groups, devices, { ...initialCollectionState(), deviceId: 'D99', status: 'draft', destination: 'DB2' });
    expect(result.rows.map(row => row.id)).toEqual([]);
    expect(result.empty).toBe('no-results');
    expect(recordingCollection(groups, devices, { ...initialCollectionState(), deviceId: 'D0', status: 'draft', destination: 'DB1' }).rows.map(row => row.id)).toEqual(['G000']);
  });
  it('breaks identical names by canonical id with stable ordering', () => {
    const result = recordingCollection([group('Z'), group('A')].map(value => ({ ...value, name: 'same' })), devices, initialCollectionState());
    expect(result.rows.map(row => row.id)).toEqual(['A', 'Z']);
  });
  it('clamps removed final pages while retaining valid same-context filters', () => {
    const state = { ...initialCollectionState(), page: 12, keyword: '記錄' };
    expect(recordingCollection(groups.slice(0, 28), devices, state)).toMatchObject({ page: 2, pageCount: 2, filteredCount: 28 });
    expect(recordingCollection(groups.slice(0, 28), devices, state).rows).toHaveLength(3);
  });
  it('distinguishes a genuinely empty workspace from no filter matches', () => {
    expect(recordingCollection([], devices, initialCollectionState()).empty).toBe('empty');
    expect(recordingCollection(groups, devices, { ...initialCollectionState(), keyword: 'absent' }).empty).toBe('no-results');
  });
  it('device summary has four rows while same-page management uses25/50', () => {
    expect(deviceCollection(devices, initialCollectionState()).summary).toHaveLength(4);
    expect(deviceCollection(devices, initialCollectionState()).rows).toHaveLength(25);
    expect(deviceCollection(devices, { ...initialCollectionState(), pageSize: 50, page: 2 }).rows).toHaveLength(50);
  });
});
