import { beforeEach, describe, expect, it } from 'vitest';
import {
  clearCutoverIntent,
  cutoverIntentKey,
  loadCutoverIntent,
  saveCutoverIntent,
  type CutoverIntentRecord,
} from '@/features/datalink/workbench-v2/state/cutoverFlow';

const record: CutoverIntentRecord = {
  workspace_id: 'ws-1',
  source_id: 'wide-1',
  kind: 'forward',
  boundary: '2099-01-01T00:00:00Z',
  review_digest: 'a'.repeat(64),
  expected_workspace_revision: 'wrev-1',
  expected_source_revision: 'rev-wide-1',
  expected_connector_revision: 'crev-1',
  request: {
    kind: 'forward',
    workspace_id: 'ws-1',
    expected_workspace_revision: 'wrev-1',
    expected_source_revision: 'rev-wide-1',
    expected_connector_revision: 'crev-1',
    boundary: '2099-01-01T00:00:00Z',
    review_digest: 'a'.repeat(64),
  },
  created_at: '2026-10-24T12:00:00Z',
};

beforeEach(() => {
  sessionStorage.clear();
});

describe('Cutover sessionStorage intent isolation and recovery', () => {
  it('saves and restores an intent scoped to its workspace', () => {
    saveCutoverIntent('ws-1', record);
    expect(loadCutoverIntent('ws-1')).toEqual(record);
    expect(loadCutoverIntent('ws-1', 'wide-1')).toEqual(record);
    // Foreign workspace scope is isolated
    expect(loadCutoverIntent('ws-2')).toBeUndefined();
  });

  it('safely handles missing, malformed or invalid stored intent data', () => {
    expect(loadCutoverIntent('ws-1')).toBeUndefined();

    // Invalid JSON
    sessionStorage.setItem(cutoverIntentKey('ws-1'), 'invalid json');
    expect(loadCutoverIntent('ws-1')).toBeUndefined();

    // Non-object JSON
    sessionStorage.setItem(cutoverIntentKey('ws-1'), '12345');
    expect(loadCutoverIntent('ws-1')).toBeUndefined();

    // Missing required fields
    sessionStorage.setItem(cutoverIntentKey('ws-1'), JSON.stringify({ workspace_id: 'ws-1' }));
    expect(loadCutoverIntent('ws-1')).toBeUndefined();
  });

  it.each(['workspace_id', 'review_digest'])('rejects a frozen request inconsistent with its stored %s', (field) => {
    const corrupt = structuredClone(record);
    if (field === 'workspace_id') corrupt.request.workspace_id = 'foreign-workspace';
    else corrupt.request.review_digest = 'b'.repeat(64);
    sessionStorage.setItem(cutoverIntentKey('ws-1'), JSON.stringify(corrupt));
    expect(loadCutoverIntent('ws-1')).toBeUndefined();
  });

  it('clears stored intent on definitive completion or rejection', () => {
    saveCutoverIntent('ws-1', record);
    expect(loadCutoverIntent('ws-1')).toBeDefined();

    clearCutoverIntent('ws-1', 'wide-1');
    expect(loadCutoverIntent('ws-1')).toBeUndefined();
    expect(loadCutoverIntent('ws-1', 'wide-1')).toBeUndefined();
  });
});
