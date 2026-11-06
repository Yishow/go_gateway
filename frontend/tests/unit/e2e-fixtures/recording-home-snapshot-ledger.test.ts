import { describe, expect, it } from 'vitest';
import { createRecordingHomeSnapshotAPI, SNAPSHOT_BASE } from '../../e2e/fixtures/recording-home-snapshot-api';
import { assertSnapshotHomeDemand, releaseSnapshotReply } from '../../e2e/helpers/recording-home-snapshot-ledger';

describe('home request ceiling assertions', () => {
  const fixture = () => {
    const api = createRecordingHomeSnapshotAPI();
    for (const group of api.snapshot().groups.slice(0, 25)) for (const action of ['lifecycle', 'delivery']) {
      api.handle({ method: 'GET', url: `${SNAPSHOT_BASE}/write-groups/${group.id}/${action}` });
    }
    return { api, ids: api.snapshot().groups.slice(0, 25).map(group => group.id) };
  };
  it('accepts exactly 25 visible groups per detail endpoint', () => {
    const { api, ids } = fixture(); expect(() => assertSnapshotHomeDemand(api.ledger, ids)).not.toThrow();
  });
  it('fails on the 26th detail, an off-page identity, a home raw read and a mutation', () => {
    const { api, ids } = fixture();
    expect(() => assertSnapshotHomeDemand([...api.ledger, api.ledger[0]], ids)).toThrow('exceeded');
    expect(() => assertSnapshotHomeDemand([{ ...api.ledger[0], groupId: 'B-0001' }], ids)).toThrow('off-page');
    expect(() => assertSnapshotHomeDemand([{ ...api.ledger[0], action: 'raw:devices' }], ids)).toThrow('raw source');
    expect(() => assertSnapshotHomeDemand([{ ...api.ledger[0], mutation: true }], ids)).toThrow('mutation');
  });
  it('rejects duplicate detail reads even below the ceiling so missing visible rows cannot be hidden', () => {
    const { api, ids } = fixture();
    expect(() => assertSnapshotHomeDemand([api.ledger[0], api.ledger[0]], ids)).toThrow('duplicate');
  });
  it('does not mislabel legitimate editor raw discovery as home transport', () => {
    const { api, ids } = fixture();
    expect(() => assertSnapshotHomeDemand([{ ...api.ledger[0], phase: 'editor', action: 'raw:devices' }], ids)).not.toThrow();
  });
});


describe('late-reply observation ordering', () => {
  it('waits for terminal delivery before asserting the current DOM, including a later stale overwrite', async () => {
    const events: string[] = [];
    let finish!: () => void, current = 'A2';
    const delivery = new Promise<void>(resolve => { finish = resolve; });
    const observed = releaseSnapshotReply(() => { events.push('release'); }, async () => { events.push('wait'); await delivery; },
      async () => { events.push('render'); }, async () => { events.push('assert'); if (current !== 'A2') throw new Error('stale DOM'); });
    const result = observed.then(() => 'passed', error => String(error));
    await Promise.resolve(); await Promise.resolve();
    expect(events).toEqual(['release', 'wait']);
    current = 'OLD-A'; finish();
    expect(await result).toContain('stale DOM');
    expect(events).toEqual(['release', 'wait', 'render', 'assert']);
  });
});
