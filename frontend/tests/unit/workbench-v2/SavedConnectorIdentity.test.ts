import { describe, expect, it } from 'vitest';
import {
  buildSchemaPreviewSignature,
  createPoolConnectorPatch,
} from '../../../src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers';
import {
  hydrateStudioV2DatabaseConnector,
  isStudioV2DatabaseConnectorValid,
  toStudioV2DatabaseConfigRequest,
} from '../../../src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';

const connector = {
  ...INITIAL_STATE.db.connector,
  connector_id: 'connector-old',
  identity_revision: 'identity-old',
  kind: 'postgres' as const,
  host: 'db-old.example',
  username: 'writer',
  password: ' old fixture password ',
  database: 'metrics',
  table: 'samples',
};

const saved = {
  ...connector,
  id: 'connector-selected',
  identity_revision: 'identity-selected',
  host: 'db-selected.example',
  password: '',
  enabled: true,
  status: 'ready' as const,
  default_write_interval_seconds: 5,
};

describe('SavedConnectorIdentity', () => {
  it('keeps a saved reference without asking for the masked password', () => {
    const patch = createPoolConnectorPatch(saved);
    expect(patch).toMatchObject({
      connector_id: 'connector-selected',
      identity_revision: 'identity-selected',
      password_required: false,
    });
    expect(patch.password).toBeUndefined();
    const request = toStudioV2DatabaseConfigRequest({ ...connector, ...patch });
    expect(request).toMatchObject({
      connector_id: 'connector-selected',
      expected_connector_revision: 'identity-selected',
    });
    expect(request).not.toHaveProperty('password');
  });

  it('preserves supplied password bytes, including an all-space password', () => {
    const draft = { ...connector, password: '   ', password_required: true };
    expect(isStudioV2DatabaseConnectorValid(draft)).toBe(true);
    expect(toStudioV2DatabaseConfigRequest(draft).password).toBe('   ');
  });

  it('keeps explicit credential removal distinct from unchanged credentials', () => {
    const request = toStudioV2DatabaseConfigRequest({ ...connector, password: '', clear_password: true });
    expect(request).toHaveProperty('clear_password', true);
    expect(request).not.toHaveProperty('password');
  });

  it('does not retain an old credential while hydrating another identity', () => {
    const record = {
      ...connector,
      id: 'connector-selected',
      identity_revision: 'identity-selected',
      workspace_id: 'workspace-a',
      save_state: 'saved' as const,
      created_at: '2026-09-15T00:00:00Z',
      updated_at: '2026-09-15T00:00:00Z',
    };
    const hydrated = hydrateStudioV2DatabaseConnector(record, connector);
    expect(hydrated.password).toBeUndefined();
    expect(hydrated).toHaveProperty('identity_revision', 'identity-selected');
  });

  it('uses saved identity revisions without putting credentials or paths in preview signatures', () => {
    const first = buildSchemaPreviewSignature(connector, {});
    const next = buildSchemaPreviewSignature({ ...connector, identity_revision: 'identity-next' }, {});
    expect(first).not.toBe(next);
    expect(first).not.toContain(connector.password);
    expect(first).not.toContain(connector.host);
    expect(first).not.toContain(connector.database);
  });

  it('does not reuse a credential when a legacy response has no identity revision', () => {
    const current = { ...connector, identity_revision: undefined };
    const record = {
      ...current, id: current.connector_id, workspace_id: 'workspace-a', host: 'changed.example',
      save_state: 'saved' as const, created_at: '2026-09-15T00:00:00Z', updated_at: '2026-09-15T00:00:00Z',
    };
    expect(hydrateStudioV2DatabaseConnector(record, current).password).toBeUndefined();
  });
});
