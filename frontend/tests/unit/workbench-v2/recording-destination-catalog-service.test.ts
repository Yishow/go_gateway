import { describe, expect, it } from 'vitest';
import { parseRecordingDestinationCatalog } from '@/services/studioV2DatabaseCatalog';

const scope = { workspaceId: 'ws-1', workspaceRevision: 'setup-1', connectorId: 'conn-1', connectorRevision: 'crev-1', kind: 'postgres' as const };
const response = { workspace_id: 'ws-1', workspace_revision: 'setup-1', connector_id: 'conn-1', connector_revision: 'crev-1', kind: 'postgres', default_schema: 'public', schemas: [{ name: 'public', tables: ['readings', 'gw_group_1'] }, { name: 'empty', tables: [] }] };

describe('recording destination catalog scope', () => {
  it('accepts a complete catalog for the requested saved revisions', () => {
    expect(parseRecordingDestinationCatalog(response, scope)).toEqual(response);
  });
  it('rejects stale or malformed replies rather than offering their tables', () => {
    for (const value of [
      { ...response, workspace_revision: 'old' },
      { ...response, connector_id: 'other' },
      { ...response, connector_revision: 'old' },
      { ...response, kind: 'mysql' },
      { ...response, schemas: [{ name: 'public', tables: null }] },
      { ...response, schemas: [{ name: 'public', tables: ['readings', 'readings'] }] },
      { ...response, schemas: [{ name: 'public', tables: [] }, { name: 'public', tables: [] }] },
    ]) expect(() => parseRecordingDestinationCatalog(value, scope)).toThrow();
  });
});
