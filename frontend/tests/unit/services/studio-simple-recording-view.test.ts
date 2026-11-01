import { describe, expect, it } from 'vitest';
import { parseStudioPresentation } from '@/services/studioRecordingView';
const locator = 'entry=locator&workspace_id=W&step=4&group_id=G&group_revision=R';
describe('bounded recording repair navigation', () => {
  it.each(['sources', 'destination'])('retains exact %s repair without changing the canonical locator', section => {
    expect(parseStudioPresentation(`?${locator}&view=recording&repair_section=${section}&repair_revision=W2`)).toMatchObject({
      presentation: 'recording', locatorIntent: { kind: 'locator', locator: { workspace_id: 'W', group_id: 'G', group_revision: 'R' } },
      repairIntent: { workspaceId: 'W', workspaceRevision: 'W2', groupId: 'G', groupRevision: 'R', section },
    });
  });
  it('retains original operation query identity as navigation only', () => {
    expect(parseStudioPresentation(`?${locator}&view=recording&repair_section=operation&repair_revision=W2&repair_kind=query_operation&repair_owner=start&repair_device=D&repair_circuit=C&repair_operation=O1`)).toMatchObject({
      presentation: 'recording', repairIntent: { section: 'operation', workspaceRevision: 'W2', operation: {
        kind: 'query_operation', owner: 'start', deviceId: 'D', circuitId: 'C', operationId: 'O1',
      } },
    });
  });
  it('rejects incomplete, duplicate, malformed and command-bearing recording repairs', () => {
    expect(parseStudioPresentation('?entry=locator&workspace_id=W&step=4&group_id=G&view=recording&repair_section=sources&repair_revision=W2')).toEqual({ presentation: 'advanced', locatorIntent: { kind: 'invalid' } });
    for (const suffix of ['repair_section=sources', 'repair_section=sources&repair_revision=W2&repair_section=destination',
      'repair_section=operation&repair_revision=W2&repair_owner=start&repair_device=D',
      'repair_section=sources&repair_revision=%ZZ', 'repair_section=sources&repair_revision=W2&command=start',
      'repair_section=sources&repair_revision=W2&intent=copy']) {
      expect(parseStudioPresentation(`?${locator}&view=recording&${suffix}`)).toEqual({ presentation: 'advanced', locatorIntent: { kind: 'invalid' } });
    }
  });
});
describe('bounded simple recording entry', () => {
  it('opens an explicit local new draft without fabricating a persisted identity', () => {
    expect(parseStudioPresentation('?view=recording&intent=new')).toEqual({ presentation: 'recording', locatorIntent: { kind: 'absent' } });
  });
  it('keeps the canonical identity for an opted-in saved recording', () => {
    expect(parseStudioPresentation(`?${locator}&view=recording`)).toMatchObject({ presentation: 'recording', locatorIntent: { kind: 'locator', locator: { workspace_id: 'W', group_id: 'G', group_revision: 'R' } } });
  });
  it.each(['?view=recording', '?view=recording&intent=edit', '?view=recording&intent=new&group_id=G',
    `?${locator}&view=recording&intent=new`, '?view=recording&intent=new&intent=new', '?view=recording&intent=new&return=https://evil.test'])('rejects ambiguous or forged new intent: %s', search => {
    expect(parseStudioPresentation(search)).toEqual({ presentation: 'advanced', locatorIntent: { kind: 'invalid' } });
  });
  it('makes only the empty-query default the recording overview', () => {
    expect(parseStudioPresentation('')).toEqual({ presentation: 'recordings', locatorIntent: { kind: 'absent' } });
    expect(parseStudioPresentation('?')).toEqual({ presentation: 'recordings', locatorIntent: { kind: 'absent' } });
    for (const search of ['?view=advanced', '?view=unknown', '?search=kept', '?section=unknown'])
      expect(parseStudioPresentation(search)).toEqual({ presentation: 'advanced', locatorIntent: { kind: 'absent' } });
  });
});

describe('bounded official legacy navigation survives presentation', () => {
  it.each(['?step=output&target=modbus', '?step=output&target=database&view=recordings', '?section=settings&modal=settings&view=recordings', '?modal=points&legacy=points&view=recordings'])('keeps %s as advanced rather than treating it as a locator', search => {
    expect(parseStudioPresentation(search)).toEqual({ presentation: 'advanced', locatorIntent: { kind: 'absent' } });
  });
  it.each(['?step=output&target=modbus&workspace_id=W', '?step=output&target=modbus&entry=locator', '?step=output&target=modbus&editor=G', '?step=output&target=modbus&group_revision=R', '?step=output&target=modbus&step=output', '?step=output&target=%ZZ', '?step=output&target=modbus&constructor=x', '?step=output&target=modbus&toString=x'])('never repairs mixed or malformed legacy intent: %s', search => {
    expect(parseStudioPresentation(search)).toEqual({ presentation: 'advanced', locatorIntent: { kind: 'invalid' } });
  });
});

import { buildSimpleRecordingTarget, readStudioLegacyPresentation } from '@/services/studioRecordingView';
it('builds saved presentation with real canonical IDs, never a fake draft group', () => {
  expect(buildSimpleRecordingTarget({ kind: 'new' }, '#position')).toBe('/studio/v2?view=recording&intent=new#position');
  expect(buildSimpleRecordingTarget({ kind: 'edit', workspaceId: 'W', groupId: 'G', groupRevision: 'R' }, '#position')).toBe(`/studio/v2?${locator}&view=recording#position`);
  expect(buildSimpleRecordingTarget({ kind: 'edit', workspaceId: '', groupId: 'G', groupRevision: 'R' })).toBeNull();
});
it('maps old navigation only to the original flow/settings surface', () => {
  expect(readStudioLegacyPresentation('?step=output&target=modbus')).toEqual({ view: 'flow', step: 4, target: 'modbus' });
  expect(readStudioLegacyPresentation('?section=settings&modal=settings')).toEqual({ view: 'settings', step: 1 });
  expect(readStudioLegacyPresentation('?legacy=points&modal=points')).toEqual({ view: 'flow', step: 2 });
  expect(readStudioLegacyPresentation('?step=output&target=modbus&workspace_id=W')).toBeNull();
});
it('keeps a source locator and runtime repair focus above a simple editor hint', () => {
  expect(parseStudioPresentation('?entry=locator&workspace_id=W&step=1&device_id=D&view=recording')).toMatchObject({ presentation: 'advanced', locatorIntent: { kind: 'locator', locator: { step: 1, device_id: 'D' } } });
  expect(parseStudioPresentation('?focus=readiness&step=3&issue=source-required&view=recording')).toEqual({ presentation: 'advanced', locatorIntent: { kind: 'absent' } });
});

it('opens a bounded copy intent using the canonical source locator, not an editable saved identity', () => {
  expect(parseStudioPresentation(`?${locator}&view=recording&intent=copy`)).toMatchObject({ presentation: 'recording', recordingIntent: 'copy',
    locatorIntent: { kind: 'locator', locator: { workspace_id: 'W', group_id: 'G', group_revision: 'R' } } });
  expect(buildSimpleRecordingTarget({ kind: 'copy', workspaceId: 'W', groupId: 'G', groupRevision: 'R' }, '#kept')).toBe(`/studio/v2?${locator}&view=recording&intent=copy#kept`);
});
it.each(['?view=recording&intent=copy', '?entry=locator&workspace_id=W&step=1&device_id=D&view=recording&intent=copy',
  `?${locator}&view=recording&intent=copy&intent=copy`, `?${locator}&view=recording&intent=copy&unexpected=command`,
  '?entry=locator&workspace_id=W&step=4&group_id=%ZZ&group_revision=R&view=recording&intent=copy'])('never repairs invalid copy intent: %s', search => {
  expect(parseStudioPresentation(search)).toEqual({ presentation: 'advanced', locatorIntent: { kind: 'invalid' } });
});
