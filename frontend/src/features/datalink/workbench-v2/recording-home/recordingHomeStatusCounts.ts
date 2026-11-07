import type { WriteGroupStatus } from '@/types/studioV2WriteGroup';

export const recordingHomeStatuses: readonly WriteGroupStatus[] = ['draft', 'ready', 'running', 'disabled', 'deleted'];

/** Count saved configuration states over all loaded recordings. */
export function countRecordingHomeStatuses(groups: readonly { readonly status: WriteGroupStatus }[]): Record<WriteGroupStatus, number> {
  const counts: Record<WriteGroupStatus, number> = { draft: 0, ready: 0, running: 0, disabled: 0, deleted: 0 };
  for (const group of groups) counts[group.status]++;
  return counts;
}
