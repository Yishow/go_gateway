import { QueryClient } from '@tanstack/react-query';
import { describe, expect, it } from 'vitest';
import { cachedHomeReads } from '@/features/datalink/workbench-v2/recording-home/recordingHomeCachedReads';
import { writeGroupDeliveryQueryOptions } from '@/hooks/datalink/useStudioV2WriteGroups';
import { writeGroupLifecycleQueryOptions } from '@/hooks/datalink/useStudioV2WriteGroupLifecycle';
import { delivery, group, groupScope, fixture } from './recording-home-guidance-fixtures';

describe('home observed-only cache index', () => {
  it('reads only exact scope/revision evidence without fetching any offpage group', async () => {
    const client = new QueryClient(); const item = group('G'); const scope = groupScope(item);
    const deliveryKey = writeGroupDeliveryQueryOptions('G', false, undefined, scope).queryKey;
    const lifecycleKey = writeGroupLifecycleQueryOptions('W1', 'G', false, scope).queryKey;
    const lifecycle = fixture(['D1'], true, [item]).evidence.lifecycle!.G;
    if (lifecycle.state !== 'available') throw new Error('fixture');
    client.setQueryData(deliveryKey, delivery(item)); client.setQueryData(lifecycleKey, lifecycle.value);
    expect(cachedHomeReads(client, [scope]).delivery.G.state).toBe('available');
    expect(client.isFetching()).toBe(0);
    for (const patch of [{ epoch: 999 }, { workspaceId: 'foreign' }, { workspaceRevision: 'new' }, { groupRevision: 'new' }, { appliedRevision: 'new' }]) {
      expect(cachedHomeReads(client, [{ ...scope, ...patch }]).delivery.G.state).toBe('unknown');
    }
    await client.invalidateQueries({ queryKey: deliveryKey, exact: true, refetchType: 'none' });
    expect(cachedHomeReads(client, [scope]).delivery.G.state).toBe('unknown');
    expect(cachedHomeReads(client, []).delivery).toEqual({}); client.clear();
  });
});
