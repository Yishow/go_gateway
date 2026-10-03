import { beforeEach, describe, expect, it, vi } from 'vitest';
import { managedSchemaPreviewResponse, managedSchemaPreviewToken } from '../../fixtures/managedSchemaPreview';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

const preview = managedSchemaPreviewToken;

const operation = {
  operation_id: 'operation-1', action: 'schema_apply', status: 'succeeded', executed_statements: 2,
  verified_digest: 'verified-digest', created_at: '2026-10-04T00:00:00Z', updated_at: '2026-10-04T00:00:01Z',
};

function envelope<T>(data: T) {
  return { data: { success: true, data } };
}

describe('canonical write-group schema service', () => {
  beforeEach(() => {
    vi.mocked(studioV2DatalinkApi.get).mockReset();
    vi.mocked(studioV2DatalinkApi.post).mockReset();
  });

  it('previews a saved group using only its revision scope and preserves the server layout', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce({ data: managedSchemaPreviewResponse } as never);

    const result = await studioV2WorkspaceWriteGroupsAPI.schemaPreview('group/1', {
      workspace_id: 'workspace-1', expected_workspace_revision: 'workspace-rev-1',
      expected_group_revision: 'group-rev-1', expected_connector_revision: 'connector-rev-1',
      token: 'client-must-not-be-sent', operation_id: 'client-must-not-be-sent',
      table_name: 'client-table-must-not-be-sent', statements: ['DROP TABLE anything'],
    } as never);

    expect(result.group_layout?.columns[0]).toEqual({ name: 'record_id', sql_type: 'TEXT', nullable: false, primary_key: true });
    expect(studioV2DatalinkApi.post).toHaveBeenCalledWith(
      '/studio-v2/workspace/write-groups/group%2F1/schema-preview',
      {
        workspace_id: 'workspace-1', expected_workspace_revision: 'workspace-rev-1',
        expected_group_revision: 'group-rev-1', expected_connector_revision: 'connector-rev-1',
      },
    );
  });

  it('rejects a write-group preview response without its complete server layout', async () => {
    const { group_layout: _groupLayout, ...legacyPreview } = preview;
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope({ ...legacyPreview, table_prefix: 'legacy_prefix' }) as never);

    await expect(studioV2WorkspaceWriteGroupsAPI.schemaPreview('group-1', {
      workspace_id: 'workspace-1', expected_workspace_revision: 'workspace-rev-1',
      expected_group_revision: 'group-rev-1', expected_connector_revision: 'connector-rev-1',
    })).rejects.toMatchObject({ name: 'RecordingPlanResponseError' });
  });

  it('confirms the same preview operation without resending table or SQL input', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope(operation) as never);

    await expect(studioV2WorkspaceWriteGroupsAPI.schemaApply('group-1', {
      workspace_id: 'workspace-1', expected_workspace_revision: 'workspace-rev-1',
      expected_group_revision: 'group-rev-1', expected_connector_revision: 'connector-rev-1',
      token: 'token-1', operation_id: 'operation-1', table_name: 'forged', statements: ['DROP TABLE forged'],
    } as never)).resolves.toEqual(operation);

    expect(studioV2DatalinkApi.post).toHaveBeenCalledWith(
      '/studio-v2/workspace/write-groups/group-1/schema-apply',
      {
        workspace_id: 'workspace-1', expected_workspace_revision: 'workspace-rev-1',
        expected_group_revision: 'group-rev-1', expected_connector_revision: 'connector-rev-1',
        token: 'token-1', operation_id: 'operation-1',
      },
    );
  });

  it('looks up the exact operation identity instead of retrying DDL', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(envelope({ ...operation, status: 'unknown' }) as never);

    await expect(studioV2WorkspaceWriteGroupsAPI.schemaOperation('operation/1')).resolves.toMatchObject({
      operation_id: 'operation-1', status: 'unknown',
    });
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/database-operations/operation%2F1');
  });
});
