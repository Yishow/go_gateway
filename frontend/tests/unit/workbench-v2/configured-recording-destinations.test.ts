import { describe, expect, it } from 'vitest';
import { configuredDestinationIdentity, configuredDestinationReason, configuredDestinationPatch, configuredDestinationStillCurrent } from '../../../src/features/datalink/workbench-v2/recording-editor/configuredRecordingDestinations';
import { batchPanelState, longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { groupToDraft } from '../../../src/features/datalink/workbench-v2/state/writeGroup/draft';
const group = longProfileGroup({ id: 'long-1', circuitId: 'c1', circuitName: 'one', interval: 15, points: [0, 1] });
const connector = batchPanelState().db.connector;
const source = { ...group, destination: { ...group.destination, database: connector.database, table_schema: connector.schema, connector_id: connector.connector_id!, connector_revision: connector.identity_revision! } };
const context = { workspaceId: source.workspace_id, connector, draft: groupToDraft(source), groupId: undefined };
describe('configured destinations use canonical configured identities, never physical inventory', () => {
 it('does not merge the same table across databases, connectors, storage modes or layouts', () => {
  const values = [source, { ...source, destination: { ...source.destination, database: 'other' } }, { ...source, destination: { ...source.destination, connector_id: 'other' } }, { ...source, storage_mode: 'timescale' as const }, { ...source, layout_version: 'managed_wide_profile_v1' as const }];
  expect(new Set(values.map(configuredDestinationIdentity)).size).toBe(5);
 });
 it('allows only same-scope canonical long intent without copying identity, proof or readiness', () => {
  expect(configuredDestinationReason(source, context)).toBeUndefined();
  expect(configuredDestinationPatch(source)).toEqual({ connector_id: source.destination.connector_id, connector_revision: source.destination.connector_revision, table_schema: source.destination.table_schema, table_name: source.destination.table_name });
 });
 it('denies another wide owner, foreign destinations and stale connector revisions with specific reasons', () => {
  expect(configuredDestinationReason({ ...source, layout_version: 'managed_wide_v1' }, context)).toBe('wideOwner');
  expect(configuredDestinationReason({ ...source, destination: { ...source.destination, database: 'elsewhere' } }, context)).toBe('destination');
  expect(configuredDestinationReason({ ...source, destination: { ...source.destination, connector_revision: 'old' } }, context)).toBe('revision');
  expect(configuredDestinationReason(source, { ...context, draft: { ...context.draft, storage_mode: 'timescale' } })).toBe('layout');
 });
 it('requires the chosen configured reference to survive a fresh pre-save read', () => {
  const list = { workspace_id: source.workspace_id, workspace_revision: 'fresh', groups: [source] };
  expect(configuredDestinationStillCurrent(source, list)).toBe(true);
  expect(configuredDestinationStillCurrent(source, { ...list, groups: [{ ...source, revision: 'changed' }] })).toBe(false);
  expect(configuredDestinationStillCurrent(source, { ...list, groups: [] })).toBe(false);
 });
});
