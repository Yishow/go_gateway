import { fireEvent, screen, waitFor } from '@testing-library/react';
import { expect } from 'vitest';
import type { WriteGroup, WriteGroupReadiness } from '@/types/studioV2WriteGroup';
import type { WriteGroupLifecycle } from '@/types/studioV2WriteGroupLifecycle';
import type { SchemaPreviewToken } from '@/types/recordingPlan';
import { managedSchemaPreviewToken } from './managedSchemaPreview';

/** Server fixtures stay bound to the caller's actual canonical snapshot, never a global ready flag. */
export function guidedReadiness(group: WriteGroup, workspaceRevision: string): WriteGroupReadiness {
  const schemaReady = group.destination.schema_revision === group.revision && Boolean(group.destination.schema_digest);
  return { workspace_id: group.workspace_id, workspace_revision: workspaceRevision, group_id: group.id, group_revision: group.revision,
    applied_revision: group.applied_revision, config_ready: true, schema_ready: schemaReady, ready: schemaReady,
    issues: schemaReady ? [] : [{ code: 'schema-confirmation-required', severity: 'blocking', step: 'Step 4', scope: '', message: '' }] };
}
export function guidedLifecycle(group: WriteGroup, workspaceRevision: string): WriteGroupLifecycle {
  return { workspace_id: group.workspace_id, workspace_revision: workspaceRevision, group_id: group.id, group_revision: group.revision,
    applied_revision: group.applied_revision, status: group.status, connector_id: group.destination.connector_id,
    connector_revision: group.destination.connector_revision,
    actions: [{ action: 'edit', allowed: true, reason_code: '', next_action: '', requires_review: false }],
    history: { canonical: true, current_group_id: group.id, generations: [] } };
}
export function guidedPreview(group: WriteGroup, workspaceRevision: string): SchemaPreviewToken {
  return { ...managedSchemaPreviewToken, token: `preview-${group.id}-${group.revision}`, operation_id: `schema-${group.id}-${group.revision}`,
    workspace_id: group.workspace_id, workspace_revision: workspaceRevision, plan_id: group.id, plan_revision: group.revision,
    connector_id: group.destination.connector_id, connector_revision: group.destination.connector_revision,
    expires_at: new Date(Date.now() + 600000).toISOString() };
}

export async function clickEnabled(testId: string) {
  const button = await screen.findByTestId(testId); await waitFor(() => expect(button).toBeEnabled()); fireEvent.click(button);
}
/** Real user commands only; this helper never supplies readiness, owner observations or tokens. */
export async function confirmPreparedStart(prefix: 'basic-recording' | 'circuit-profile') {
  for (let index = 0; index < 4; index += 1) {
    const next = await waitFor(() => {
      const start = screen.queryByTestId(`${prefix}-start`);
      if (start && !(start as HTMLButtonElement).disabled) return start;
      const preview = screen.getByTestId('group-schema-preview'); expect(preview).toBeEnabled(); return preview;
    });
    if (next.getAttribute('data-testid') === `${prefix}-start`) {
      fireEvent.click(next); expect(screen.getByRole('dialog')).toBeVisible(); await clickEnabled('recording-start-confirm'); return;
    }
    fireEvent.click(next); await screen.findByTestId('group-schema-preview-result'); await clickEnabled('group-schema-apply');
  }
  throw new Error('Expected exact saved groups to finish explicit schema confirmation');
}
export async function startGuided(prefix: 'basic-recording' | 'circuit-profile') {
  const prepare = screen.queryByTestId(`${prefix}-prepare`);
  if (prepare) { await waitFor(() => expect(prepare).toBeEnabled()); fireEvent.click(prepare); }
  await confirmPreparedStart(prefix);
}
