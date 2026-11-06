import type { SnapshotCall } from '../fixtures/recording-home-snapshot-api';
/** Per user page intent, not an aggregate count across multiple intentional pages. */
export function assertSnapshotHomeDemand(calls: readonly SnapshotCall[], visibleIds: readonly string[], ceiling = 25) {
  const reads = calls.filter(call => call.phase === 'home');
  if (reads.some(call => call.action.startsWith('raw:'))) throw new Error('Home used editor raw source transport');
  if (reads.some(call => call.mutation)) throw new Error('Home performed a mutation');
  for (const action of ['lifecycle', 'delivery']) {
    const details = reads.filter(call => call.action === action);
    if (details.length > ceiling) throw new Error(`${action} exceeded ${ceiling} requests for one visible page`);
    if (new Set(details.map(call => call.groupId)).size !== details.length) throw new Error(`${action} made duplicate reads for one visible page`);
    if (details.some(call => !call.groupId || !visibleIds.includes(call.groupId))) throw new Error(`${action} requested an off-page group`);
  }
}

/** Release one captured reply and observe the UI at the response-completion boundary. */
export async function releaseSnapshotReply(release: () => void, waitForDelivery: () => Promise<void>,
  render: () => Promise<void>, assertCurrent: () => Promise<void>) {
  release(); await waitForDelivery(); await render(); await assertCurrent();
}
