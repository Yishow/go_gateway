import { expect, it } from 'vitest';
import { countRecordingHomeStatuses } from '@/features/datalink/workbench-v2/recording-home/recordingHomeStatusCounts';

it('counts the saved statuses of the mixed workspace from the specification', () => {
  expect(countRecordingHomeStatuses([
    { status: 'draft' }, { status: 'draft' }, { status: 'draft' }, { status: 'deleted' },
  ])).toEqual({ draft: 3, ready: 0, running: 0, disabled: 0, deleted: 1 });
  expect(countRecordingHomeStatuses([])).toEqual({ draft: 0, ready: 0, running: 0, disabled: 0, deleted: 0 });
});
