import { beforeEach, describe, expect, it, vi } from 'vitest';
import { readRuleProof, readRangeCompletion, matchesRangeRuleAcknowledgement } from '../../../src/pages/datalink/workbench-v2/rangeSourceReadback';
import type { RangeRuleReceipt, RangeSourceRequest, RangeMappingReceipt } from '../../../src/pages/datalink/workbench-v2/rangeSourceSubmission';
const api = vi.hoisted(() => ({ workspace: vi.fn(), rules: vi.fn(), devices: vi.fn(), candidates: vi.fn(), mappings: vi.fn(), point: vi.fn(), tag: vi.fn() }));
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2WorkspaceAPI: { get: api.workspace } }));
vi.mock('../../../src/services/studioV2Rules', () => ({ studioV2RulesAPI: { list: api.rules } }));
vi.mock('../../../src/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: api.devices } }));
vi.mock('../../../src/services/sourceRuleCandidates', () => ({ sourceRuleCandidateAPI: { get: api.candidates } }));
vi.mock('../../../src/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: api.mappings } }));
vi.mock('../../../src/services/datalink', () => ({ pointAPI: { get: api.point }, tagAPI: { get: api.tag } }));
const buffer: RangeSourceRequest['buffer'] = { name: 'Local title', start_address: '40001', count: 1, data_type: 'int16', naming_prefix: 'P', enabled: true, skipped_addresses: [], scale_multiplier: 1, scale_offset: 0, data_format: '', share_enabled: false, share_start_register: null, share_stride: null };
const request: RangeSourceRequest = { workspaceId: 'w', sessionId: 's', requestId: 'q', generation: 1, deviceId: 'd', buffer };
const fields = { ...buffer }; delete (fields as Partial<typeof buffer>).name;
const receipt: RangeRuleReceipt = { submitted: { id: 'r', device_id: 'd', ...fields }, acknowledged: { id: 'r', device_id: 'd', workspace_id: 'w', revision_id: 'rr', ...fields, locked: false, origin: 'manual', created_at: '', updated_at: '', runtime_apply_status: 'applied' } };
const mapping: RangeMappingReceipt = { submitted: { rule_id: 'r', address: '40001', tag_key: 'P1', display_name: 'P1', unit: '', target_type: 'int16', scale: 1, offset: 0, enabled: true }, acknowledged: { id: 'm', workspace_id: 'w', device_id: 'd', rule_id: 'r', address: '40001', point_id: 'p', tag_id: 't', tag_key: 'P1', display_name: 'P1', unit: '', target_type: 'int16', scale: 1, offset: 0, enabled: true, created_at: '', updated_at: '', runtime_apply_status: 'applied' } };
function snapshot(mapped = false) { return { workspace_id: 'w', workspace_revision: 'source-rev', source_rule_id: 'r', revision_id: 'rr', tags: { status: 'ready', candidates: [{ id: 'c', identity: { source_rule_id: 'r' }, address: '40001', point_id: 'p', ...(mapped ? { mapping_id: 'm', tag_id: 't' } : {}), tag_key: 'P1', display_name: 'P1', data_type: 'int16', status: 'draft' }] } }; }
beforeEach(() => {
  vi.resetAllMocks();
  api.workspace.mockResolvedValue({ id: 'w', ordered_device_ids: ['d'], modbus_share: { workspace_revision: 'source-rev' } });
  api.rules.mockResolvedValue([receipt.acknowledged]);
  api.devices.mockResolvedValue([{ id: 'd', name: 'PLC', protocol: 'modbus_tcp' }]);
  api.candidates.mockResolvedValue(snapshot());
  api.mappings.mockResolvedValue([mapping.acknowledged]);
  api.point.mockResolvedValue({ id: 'p', device_id: 'd', address: '40001', data_type: 'int16', enabled: true });
  api.tag.mockResolvedValue({ id: 't', key: 'P1', display_name: 'P1', unit: '', data_type: 'int16', status: 'active' });
});
describe('request-bound range readback', () => {
  it('reads acknowledged rule with source revision, before tags exist', async () => {
    const proof = await readRuleProof(receipt, request);
    expect(proof).toMatchObject({ record: { id: 'r' }, snapshot: { revision_id: 'rr' } });
    expect(api.candidates).toHaveBeenCalledWith('r', { workspace_id: 'w', expected_workspace_revision: 'source-rev', revision_id: 'rr' });
  });
  it.each(['ack-id', 'ack-fields', 'read-fields', 'scope', 'no-revision', 'runtime'])('does not prove %s', async kind => {
    const ack = structuredClone(receipt);
    if (kind === 'ack-id') ack.acknowledged.id = 'other';
    if (kind === 'ack-fields') ack.acknowledged.count = 2;
    if (kind === 'read-fields') api.rules.mockResolvedValue([{ ...receipt.acknowledged, count: 2 }]);
    if (kind === 'scope') api.candidates.mockResolvedValue({ ...snapshot(), workspace_id: 'other' });
    if (kind === 'no-revision') api.workspace.mockResolvedValue({ id: 'w', ordered_device_ids: ['d'] });
    if (kind === 'runtime') ack.acknowledged.runtime_apply_status = 'apply_failed';
    expect(await readRuleProof(ack, request)).toBeUndefined();
  });
  it('distinguishes rule proposal metadata from the acknowledged actual mapping and tag', async () => {
    const proof = await readRuleProof(receipt, request); expect(proof).toBeDefined();
    const candidate = snapshot(true); candidate.tags.candidates[0].tag_key = 'rule_proposed_name'; candidate.tags.candidates[0].data_type = 'float32';
    api.candidates.mockResolvedValue(candidate);
    expect(await readRangeCompletion(proof!, [mapping])).toMatchObject({ candidates: [{ tag_key: 'P1', target_type: 'int16' }] });
  });
  it('accepts an unmapped revisioned proposal only with exact live workspace mapping and point/tag receipts', async () => {
    const proof = await readRuleProof(receipt, request); expect(proof).toBeDefined();
    expect(await readRangeCompletion(proof!, [mapping])).toMatchObject({ candidates: [{ point_id: 'p', tag_id: 't' }] });
    api.mappings.mockResolvedValue([{ ...mapping.acknowledged, rule_id: 'foreign' }]);
    expect(await readRangeCompletion(proof!, [mapping])).toBeUndefined();
  });
  it('parses the HTTP empty-unit omission as empty only, while rejecting a nonempty different unit', async () => {
    const proof = await readRuleProof(receipt, request); expect(proof).toBeDefined();
    // Same serialized shape asserted by the SQLite + HTTP contract test: unit
    // is absent from Tag GET, but workspace mapping GET includes unit: "".
    const wire = '{"id":"t","key":"P1","display_name":"P1","data_type":"int16","status":"active"}';
    api.tag.mockResolvedValue(JSON.parse(wire));
    expect(await readRangeCompletion(proof!, [mapping])).toMatchObject({ candidates: [{ point_id: 'p', tag_id: 't' }] });
    api.tag.mockResolvedValue({ ...JSON.parse(wire), unit: 'C' });
    expect(await readRangeCompletion(proof!, [mapping])).toBeUndefined();
    api.tag.mockResolvedValue({ ...JSON.parse(wire), unit: null });
    expect(await readRangeCompletion(proof!, [mapping])).toBeUndefined();
  });
  it('requires all exact mapping receipts and canonical point/tag identities for completion', async () => {
    const proof = await readRuleProof(receipt, request); expect(proof).toBeDefined();
    api.candidates.mockResolvedValue(snapshot(true));
    expect(await readRangeCompletion(proof!, [mapping])).toMatchObject({ candidates: [{ key: 'p|t', point_id: 'p', tag_id: 't' }] });
    api.tag.mockResolvedValue({ id: 'foreign' });
    expect(await readRangeCompletion(proof!, [mapping])).toBeUndefined();
  });
});

describe('range proof safety boundaries', () => {
  it('uses persisted MC evidence before mapping, and rejects incompatible persisted addresses', async () => {
    const req = structuredClone(request); req.buffer.start_address = 'D100';
    const ack = structuredClone(receipt); ack.submitted.start_address = 'D100'; ack.acknowledged.start_address = 'D100';
    api.rules.mockResolvedValue([ack.acknowledged]); api.devices.mockResolvedValue([{ id: 'd', name: 'MC', protocol: 'mc_3e' }]);
    const candidate = snapshot(); candidate.tags.candidates[0].address = 'D100'; api.candidates.mockResolvedValue(candidate);
    expect(await readRuleProof(ack, req)).toMatchObject({ mcEvidence: { status: 'compatible', points: [{ point_id: 'p', address: 'D100' }] } });
    candidate.tags.candidates[0].address = 'D101';
    expect(await readRuleProof(ack, req)).toBeUndefined();
  });
  it.each(['missing', 'duplicate', 'lost-ack', 'mapping-fields', 'mapping-id', 'point-id', 'point-type', 'candidate-blocked', 'candidate-mapping-id', 'candidate-tag-id', 'live-device-id', 'live-workspace-id', 'live-point-id', 'live-tag-id', 'candidate-out-of-sync', 'candidate-error', 'runtime', 'cleanup', 'revision'])('rejects completion with %s evidence', async kind => {
    const proof = await readRuleProof(receipt, request); expect(proof).toBeDefined();
    const ack = structuredClone(mapping), candidate = snapshot(true); api.candidates.mockResolvedValue(candidate);
    let receipts = [ack];
    if (kind === 'missing') receipts = [];
    if (kind === 'duplicate') receipts = [ack, ack];
    if (kind === 'lost-ack') ack.acknowledged.id = '';
    if (kind === 'mapping-fields') api.mappings.mockResolvedValue([{ ...mapping.acknowledged, scale: 2 }]);
    if (kind === 'mapping-id') api.mappings.mockResolvedValue([{ ...mapping.acknowledged, id: 'foreign' }]);
    if (kind === 'point-id') api.point.mockResolvedValue({ id: 'foreign' });
    if (kind === 'point-type') api.point.mockResolvedValue({ id: 'p', device_id: 'd', address: '40001', enabled: true, data_type: 'float64' });
    if (kind === 'candidate-mapping-id') candidate.tags.candidates[0].mapping_id = 'foreign';
    if (kind === 'candidate-tag-id') candidate.tags.candidates[0].tag_id = 'foreign';
    for (const [test, field] of [['live-device-id', 'device_id'], ['live-workspace-id', 'workspace_id'], ['live-point-id', 'point_id'], ['live-tag-id', 'tag_id']]) {
      if (kind === test) api.mappings.mockResolvedValue([{ ...mapping.acknowledged, [field]: 'foreign' }]);
    }
    if (kind === 'candidate-out-of-sync') candidate.tags.candidates[0].status = 'out_of_sync';
    if (kind === 'candidate-error') candidate.tags.candidates[0].status = 'error';
    if (kind === 'candidate-blocked') Object.assign(candidate.tags.candidates[0], { blocking_reason: 'conflict' });
    if (kind === 'runtime') ack.acknowledged.runtime_apply_status = 'apply_failed';
    if (kind === 'cleanup') ack.acknowledged.cleanup_status = 'failed';
    if (kind === 'revision') api.workspace.mockResolvedValue({ id: 'w', ordered_device_ids: ['d'], modbus_share: { workspace_revision: 'changed' } });
    expect(await readRangeCompletion(proof!, receipts)).toBeUndefined();
  });
});

describe('offline saved identities are not selectable sources', () => {
  it.each(['explicit-disabled', 'backend-disabled'])('proves %s rule/mapping/point/tag without pretending they are enabled', async mode => {
    const req = structuredClone(request), ack = structuredClone(receipt), map = structuredClone(mapping);
    if (mode === 'explicit-disabled') { req.buffer.enabled = false; ack.submitted.enabled = false; }
    ack.acknowledged.enabled = false;
    map.submitted.enabled = false; map.acknowledged.enabled = false;
    api.rules.mockResolvedValue([ack.acknowledged]); api.mappings.mockResolvedValue([map.acknowledged]);
    api.point.mockResolvedValue({ id: 'p', device_id: 'd', address: '40001', data_type: 'int16', enabled: false });
    api.tag.mockResolvedValue({ id: 't', key: 'P1', display_name: 'P1', data_type: 'int16', status: 'draft' });
    const proof = await readRuleProof(ack, req); expect(proof).toBeDefined();
    expect(await readRangeCompletion(proof!, [map])).toEqual({
      inactive: true, savedPoints: [{ device_id: 'd', point_id: 'p', tag_id: 't' }], candidates: [], excluded: [],
    });
  });
  it('includes exact saved identities on active completion', async () => {
    const proof = await readRuleProof(receipt, request);
    expect(await readRangeCompletion(proof!, [mapping])).toMatchObject({ inactive: false, savedPoints: [{ device_id: 'd', point_id: 'p', tag_id: 't' }] });
  });
  it.each(['point', 'mapping', 'unknown-point', 'foreign-tag'])('rejects mixed inactive evidence: %s', async kind => {
    const req = structuredClone(request), ack = structuredClone(receipt), map = structuredClone(mapping);
    req.buffer.enabled = false; ack.submitted.enabled = false; ack.acknowledged.enabled = false;
    map.submitted.enabled = kind === 'mapping'; map.acknowledged.enabled = kind === 'mapping';
    api.rules.mockResolvedValue([ack.acknowledged]); api.mappings.mockResolvedValue([map.acknowledged]);
    api.point.mockResolvedValue({ id: 'p', device_id: 'd', address: '40001', data_type: 'int16', ...(kind === 'unknown-point' ? {} : { enabled: kind === 'point' }) });
    if (kind === 'foreign-tag') api.tag.mockResolvedValue({ id: 'foreign' });
    const proof = await readRuleProof(ack, req); expect(proof).toBeDefined();
    expect(await readRangeCompletion(proof!, [map])).toBeUndefined();
  });
  it('never accepts false to true activation or other field changes as a readiness downgrade', async () => {
    const req = structuredClone(request), ack = structuredClone(receipt);
    req.buffer.enabled = false; ack.submitted.enabled = false;
    expect(await readRuleProof(ack, req)).toBeUndefined();
    req.buffer.enabled = true; ack.submitted.enabled = true; ack.acknowledged.enabled = false; ack.acknowledged.count = 2;
    expect(await readRuleProof(ack, req)).toBeUndefined();
  });
});

describe('shared acknowledgement matcher', () => {
  it('accepts only the backend activation downgrade without changing the captured request', () => {
    const req = structuredClone(request), ack = structuredClone(receipt);
    ack.acknowledged.enabled = false;
    expect(matchesRangeRuleAcknowledgement(ack, req)).toBe(true);
    expect(req.buffer.enabled).toBe(true);
    expect(ack.submitted.enabled).toBe(true);
    ack.acknowledged.scale_multiplier = 2;
    expect(matchesRangeRuleAcknowledgement(ack, req)).toBe(false);
    ack.acknowledged.scale_multiplier = 1;
    req.buffer.enabled = false; ack.submitted.enabled = false; ack.acknowledged.enabled = true;
    expect(matchesRangeRuleAcknowledgement(ack, req)).toBe(false);
    expect(api.workspace).not.toHaveBeenCalled();
  });
});
