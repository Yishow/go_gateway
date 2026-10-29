import { useQuery } from '@tanstack/react-query';
import type { EntityLocator } from '../../types/runtimeRuleEvidence';
import type { StudioLocatorEntryResult } from '../../types/studioV2LocatorEntry';
import { buildStudioEntityLocatorTarget, parseStudioEntityLocator } from '../../services/studioEntityLocator';
import { readStudioLocatorEntry } from '../../services/studioV2LocatorEntry';

/** Partial locator snapshots are isolated from all canonical setup query caches. */
export function useStudioLocatorEntry(locator: EntityLocator | null) {
  const target = locator === null ? null : buildStudioEntityLocatorTarget(locator);
  const intent = target ? parseStudioEntityLocator(target.slice('/studio/v2'.length)) : null;
  const selected = intent?.kind === 'locator' ? intent.locator : null;
  return useQuery<StudioLocatorEntryResult>({
    queryKey: ['studio-locator-entry', target],
    queryFn: ({ signal }) => selected ? readStudioLocatorEntry(selected, signal) : Promise.resolve({ ok: false, reason: 'unavailable' }),
    enabled: selected !== null, retry: false, staleTime: 0, gcTime: 0, refetchOnWindowFocus: false,
  });
}
