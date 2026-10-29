import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import type { RuleEvidenceView } from '@/types/runtimeRuleEvidence';
import { evidenceContext, evidenceScope, ruleEvidenceFixture } from './runtimeRuleEvidenceFixture';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
const selection = { rule_id: 'R', rule_revision: 'R2', point_id: 'P', mapping_id: 'M' };
const unavailable = { ok: false, reason: 'unavailable', action: 'refresh' };
const stale = { ok: false, reason: 'stale-reference', action: 'refresh' };
function fixture(): RuleEvidenceView {
  return { ...ruleEvidenceFixture(), source: {
    rule_id: 'R', point_id: 'P', mapping_id: 'M', tag_id: 'T', saved_revision: 'R2', effective_revision: null,
    state: 'saved-only', proof: 'unavailable', effective: null,
    saved: { binding: { workspace_id: 'W', workspace_revision: 'W1', entity_revision: 'R2' }, mapping_revision: 'MR2', enabled: true,
      address: '40001', protocol_function: '03', read_type: 'float32', result_type: 'float64',
      polling: { group_id: 'PG', interval_ms: 1000, source: 'polling-group' },
      transforms: [{ kind: 'byte-order', order: 'CDAB', multiplier: null, offset: null, target_type: null },
        { kind: 'scale', order: null, multiplier: 0.1, offset: 0, target_type: null }] },
  } };
}
const reply = (evidence: unknown = fixture(), success = true) => vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success, data: evidenceContext(evidence) } });
beforeEach(() => vi.mocked(studioV2DatalinkApi.get).mockReset());

describe('single-source runtime evidence read service', () => {
  it('adds explicit source selection while preserving recording-only reads', async () => {
    expect(studioV2RuntimeContextAPI.getSourceEvidence).toBeTypeOf('function');
    reply(ruleEvidenceFixture()); await studioV2RuntimeContextAPI.getRuleEvidence(evidenceScope);
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/runtime-context', { params: evidenceScope, signal: undefined });
  });
  it('keeps R2 scale 0.1 saved-only and polling 1000 ms separate from effective mean / 15 seconds', async () => {
    const signal = new AbortController().signal; reply();
    const result = await studioV2RuntimeContextAPI.getSourceEvidence(evidenceScope, selection, signal);
    expect(result).toEqual({ ok: true, value: fixture() });
    if (!result.ok) throw new Error('expected scoped source facts');
    expect(result.value.source).toMatchObject({ state: 'saved-only', proof: 'unavailable', effective: null, effective_revision: null });
    expect(result.value.source?.saved?.transforms[1].multiplier).toBe(0.1);
    expect(result.value.source?.saved?.polling).toEqual({ group_id: 'PG', interval_ms: 1000, source: 'polling-group' });
    expect(result.value.groups[0].effective).toMatchObject({ reduction: 'mean', interval_seconds: 15 });
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/runtime-context', { params: { ...evidenceScope, ...selection }, signal });
  });
  it('does not supply or accept a default mapping for point-only selection', async () => {
    const evidence = fixture(); Object.assign(evidence.source!, { mapping_id: null, tag_id: null });
    Object.assign(evidence.source!.saved!, { mapping_revision: null, result_type: null });
    evidence.source!.saved!.transforms = [evidence.source!.saved!.transforms[0]];
    const pointSelection = { rule_id: 'R', rule_revision: 'R2', point_id: 'P' };
    reply(evidence);
    expect(await studioV2RuntimeContextAPI.getSourceEvidence(evidenceScope, pointSelection)).toEqual({ ok: true, value: evidence });
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/runtime-context', { params: { ...evidenceScope, ...pointSelection }, signal: undefined });
    reply(); expect(await studioV2RuntimeContextAPI.getSourceEvidence(evidenceScope, pointSelection)).toEqual(stale);
  });
  it.each(['rule_id', 'point_id', 'mapping_id', 'saved_revision'] as const)('rejects changed returned source %s', async (key) => {
    const evidence = fixture(); evidence.source![key] = 'other';
    if (key === 'saved_revision') evidence.source!.saved!.binding.entity_revision = 'other';
    reply(evidence); expect(await studioV2RuntimeContextAPI.getSourceEvidence(evidenceScope, selection)).toEqual(stale);
  });
  it.each([null, undefined])('rejects missing selected source %s', async (source) => {
    reply({ ...fixture(), source }); expect(await studioV2RuntimeContextAPI.getSourceEvidence(evidenceScope, selection)).toEqual(unavailable);
  });
  it.each(['rule_id', 'rule_revision', 'point_id', 'mapping_id'] as const)('does not send empty or unsafe %s', async (key) => {
    for (const value of ['', '/private/secret', 'db.private.internal']) {
      expect(await studioV2RuntimeContextAPI.getSourceEvidence(evidenceScope, { ...selection, [key]: value })).toEqual(unavailable);
    }
    expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
  });
  it('pins caller scope and source selection before transport completes', async () => {
    let finish!: (value: unknown) => void;
    vi.mocked(studioV2DatalinkApi.get).mockReturnValueOnce(new Promise((resolve) => { finish = resolve; }));
    const scope = { ...evidenceScope }; const selected = { ...selection };
    const pending = studioV2RuntimeContextAPI.getSourceEvidence(scope, selected);
    scope.device_id = 'E'; selected.rule_id = 'other'; selected.rule_revision = 'R3'; selected.mapping_id = 'other';
    finish({ data: { success: true, data: evidenceContext(fixture()) } });
    expect(await pending).toEqual({ ok: true, value: fixture() });
  });
  it.each(['workspace_id', 'device_id', 'workspace_revision'] as const)('rejects returned scope %s mismatch', async (key) => {
    const evidence = fixture(); evidence[key] = 'other'; reply(evidence);
    expect(await studioV2RuntimeContextAPI.getSourceEvidence(evidenceScope, selection)).toEqual(key === 'workspace_revision' ? stale : { ok: false, reason: 'foreign-scope', action: 'open-workspace-list' });
  });
  it('rejects bad workspace scope and sends only allowlisted caller keys', async () => {
    expect(await studioV2RuntimeContextAPI.getSourceEvidence({ ...evidenceScope, workspace_revision: '' }, selection)).toEqual(unavailable);
    expect(studioV2DatalinkApi.get).not.toHaveBeenCalled(); reply();
    const scope = { ...evidenceScope, secret: '/private/secret' }; const selected = { ...selection, secret: '/private/secret' };
    await studioV2RuntimeContextAPI.getSourceEvidence(scope, selected);
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/runtime-context', { params: { ...evidenceScope, ...selection }, signal: undefined });
  });
  it.each(['stale-reference', 'foreign-scope', 'deleted', 'driver: secret'])('returns safe typed failure %s', async (code) => {
    vi.mocked(studioV2DatalinkApi.get).mockRejectedValueOnce({ response: { data: { error: { code, message: '/private/secret', action: 'activate' } } } });
    const result = await studioV2RuntimeContextAPI.getSourceEvidence(evidenceScope, selection);
    expect(result).toEqual(code === 'stale-reference' ? stale : code === 'foreign-scope' ? { ok: false, reason: code, action: 'open-workspace-list' } : unavailable);
    expect(JSON.stringify(result)).not.toMatch(/secret|activate/);
  });
  it('rejects failed envelopes and arbitrary unsafe transform parameters', async () => {
    reply(fixture(), false); expect(await studioV2RuntimeContextAPI.getSourceEvidence(evidenceScope, selection)).toEqual(unavailable);
    const evidence = fixture(); Object.assign(evidence.source!.saved!.transforms[1], { params: { password: '/private/secret' } }); reply(evidence);
    const result = await studioV2RuntimeContextAPI.getSourceEvidence(evidenceScope, selection);
    expect(result).toEqual(unavailable); expect(JSON.stringify(result)).not.toContain('secret');
  });
  it.each(['tag', 'mapping-revision', 'result-type', 'scale'])('rejects hidden mapping facts for a point-only request: %s', async (field) => {
    const evidence = fixture(); const source = evidence.source!; const saved = source.saved!;
    source.mapping_id = null; source.tag_id = null; saved.mapping_revision = null; saved.result_type = null;
    saved.transforms = [saved.transforms[0]];
    if (field === 'tag') source.tag_id = 'T';
    if (field === 'mapping-revision') saved.mapping_revision = 'MR2';
    if (field === 'result-type') saved.result_type = 'float64';
    if (field === 'scale') saved.transforms.push({ kind: 'scale', order: null, multiplier: 0.1, offset: 0, target_type: null });
    reply(evidence);
    expect(await studioV2RuntimeContextAPI.getSourceEvidence(evidenceScope, { rule_id: 'R', rule_revision: 'R2', point_id: 'P' })).toEqual(unavailable);
  });
});

describe('runtime context revision compatibility', () => {
  it('keeps the safe observed workspace revision without inventing legacy values', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success: true, data: { ...evidenceContext(), workspace_revision: 'W1' } } });
    expect(await studioV2RuntimeContextAPI.get()).toMatchObject({ workspace_id: 'W', workspace_revision: 'W1' });
    reply(); expect(await studioV2RuntimeContextAPI.get()).not.toHaveProperty('workspace_revision');
  });
  it.each(['', null, undefined, '/private/secret', 'db.private.internal', 123])('discards missing or unsafe revision %s', async (workspace_revision) => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success: true, data: { ...evidenceContext(), workspace_revision } } });
    expect(await studioV2RuntimeContextAPI.get()).not.toHaveProperty('workspace_revision');
  });
});
