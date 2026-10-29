import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { evidenceContext, evidenceScope, ruleEvidenceFixture } from './runtimeRuleEvidenceFixture';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
const unavailable = { ok: false, reason: 'unavailable', action: 'refresh' };
const reply = (data: unknown, success = true) => vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success, data } });
beforeEach(() => vi.clearAllMocks());

describe('shared runtime rule evidence read service', () => {
  it('provides a scoped read method without changing ordinary context reads', async () => {
    reply(evidenceContext());
    expect(await studioV2RuntimeContextAPI.get()).not.toHaveProperty('rule_evidence');
    expect(studioV2RuntimeContextAPI.getRuleEvidence).toBeTypeOf('function');
  });
  it('keeps saved C and requested B separate from effective A before boundary and B after', async () => {
    const signal = new AbortController().signal;
    for (const revision of ['A', 'B']) {
      const evidence = ruleEvidenceFixture(evidenceScope, revision);
      reply(evidenceContext(evidence));
      const result = await studioV2RuntimeContextAPI.getRuleEvidence(evidenceScope, signal);
      expect(result).toEqual({ ok: true, value: evidence });
      if (!result.ok) throw new Error('expected verified evidence');
      const group = result.value.groups[0];
      expect([group.saved_revision, group.applied_revision, group.effective_revision]).toEqual(['C', 'B', revision]);
      expect([group.effective?.destination.table, group.effective?.interval_seconds]).toEqual(revision === 'A' ? ['meter_old', 15] : ['meter_new', 30]);
      expect(group.saved?.destination.table).toBe('meter_draft');
      expect(group.as_of).toBe(result.value.observed_at);
    }
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2);
    expect(studioV2DatalinkApi.get).toHaveBeenLastCalledWith('/studio-v2/workspace/runtime-context', { params: evidenceScope, signal });
  });
  it('degrades missing layout proof without copying saved facts into effective', async () => {
    const evidence = ruleEvidenceFixture(); evidence.groups[0].proof = 'unavailable';
    reply(evidenceContext(evidence));
    const result = await studioV2RuntimeContextAPI.getRuleEvidence(evidenceScope);
    expect(result.ok).toBe(true);
    if (!result.ok) throw new Error('expected safe degradation');
    expect(result.value.groups[0]).toMatchObject({ state: 'unavailable', proof: 'unavailable', effective: null, saved_revision: 'C' });
    expect(result.value.groups[0].saved?.destination.table).toBe('meter_draft');
  });
  it('retains ended immutable evidence without inventing intake or permissions', async () => {
    const evidence = ruleEvidenceFixture(); evidence.groups[0].effective!.ends_at = evidence.observed_at;
    evidence.groups[0].saved!.enabled = false;
    reply(evidenceContext(evidence));
    const result = await studioV2RuntimeContextAPI.getRuleEvidence(evidenceScope);
    expect(result).toEqual({ ok: true, value: evidence });
    expect(JSON.stringify(result)).not.toMatch(/"intake"|"allowed"|"start"/);
  });
  it.each(['workspace_id', 'device_id'] as const)('rejects foreign %s instead of trusting response-owned scope', async (key) => {
    const evidence = ruleEvidenceFixture(); evidence[key] = 'foreign';
    reply(evidenceContext(evidence));
    expect(await studioV2RuntimeContextAPI.getRuleEvidence(evidenceScope)).toEqual({ ok: false, reason: 'foreign-scope', action: 'open-workspace-list' });
  });
  it('rejects foreign outer context even when its nested evidence matches', async () => {
    reply(evidenceContext(ruleEvidenceFixture(), 'foreign'));
    expect(await studioV2RuntimeContextAPI.getRuleEvidence(evidenceScope)).toEqual({ ok: false, reason: 'foreign-scope', action: 'open-workspace-list' });
  });
  it.each(['workspace_revision', 'as_of', 'entity_revision'])('rejects mixed %s', async (field) => {
    const evidence = ruleEvidenceFixture();
    if (field === 'workspace_revision') evidence.workspace_revision = 'W2';
    if (field === 'as_of') evidence.groups[0].as_of = '2026-10-27T10:00:01Z';
    if (field === 'entity_revision') evidence.groups[0].effective!.binding.entity_revision = 'B';
    reply(evidenceContext(evidence));
    expect(await studioV2RuntimeContextAPI.getRuleEvidence(evidenceScope)).toEqual({ ok: false, reason: 'stale-reference', action: 'refresh' });
  });
  it.each([null, undefined, {}, { workspace_id: 'W' }])('missing evidence stays unavailable: %s', async (evidence) => {
    reply({ ...evidenceContext(), rule_evidence: evidence });
    expect(await studioV2RuntimeContextAPI.getRuleEvidence(evidenceScope)).toEqual(unavailable);
  });
  it('rejects a failed envelope even when it contains plausible evidence', async () => {
    reply(evidenceContext(), false);
    expect(await studioV2RuntimeContextAPI.getRuleEvidence(evidenceScope)).toEqual(unavailable);
  });
  it.each(['password', 'host', 'database', 'raw_error', 'params'])('does not retain unsafe %s from a server payload', async (field) => {
    const evidence = ruleEvidenceFixture(); Object.assign(evidence.groups[0].effective!.destination, { [field]: '/private/secret' });
    reply(evidenceContext(evidence));
    const result = await studioV2RuntimeContextAPI.getRuleEvidence(evidenceScope);
    expect(result).toEqual(unavailable); expect(JSON.stringify(result)).not.toContain('secret');
  });
  it.each(['db.private.internal', '10.20.30.40', 'SELECT secret FROM credentials'])('rejects unsafe syntax inside contract-approved text: %s', async (text) => {
    const evidence = ruleEvidenceFixture(); evidence.groups[0].effective!.destination.alias = text;
    reply(evidenceContext(evidence));
    expect(await studioV2RuntimeContextAPI.getRuleEvidence(evidenceScope)).toEqual(unavailable);
  });
  it.each(['stale-reference', 'foreign-scope', 'deleted', 'unavailable', 'driver: secret'])('maps server error %s without reflecting raw diagnostics', async (code) => {
    vi.mocked(studioV2DatalinkApi.get).mockRejectedValueOnce({ response: { data: { error: { code, message: '/private/secret', action: 'activate' } } } });
    const result = await studioV2RuntimeContextAPI.getRuleEvidence(evidenceScope);
    expect(result).toEqual({ ok: false, reason: code === 'foreign-scope' || code === 'stale-reference' ? code : 'unavailable', action: code === 'foreign-scope' ? 'open-workspace-list' : 'refresh' });
    expect(JSON.stringify(result)).not.toMatch(/secret|activate/);
  });
  it.each(['', '/private/secret', 'db.private.internal'])('rejects invalid request scope %s before issuing a request', async (workspace_id) => {
    expect(await studioV2RuntimeContextAPI.getRuleEvidence({ ...evidenceScope, workspace_id })).toEqual(unavailable);
    expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
  });
  it('pins the caller scope before waiting for the transport response', async () => {
    let finish!: (value: { data: { success: boolean; data: ReturnType<typeof evidenceContext> } }) => void;
    vi.mocked(studioV2DatalinkApi.get).mockReturnValueOnce(new Promise((resolve) => { finish = resolve; }));
    const scope = { ...evidenceScope };
    const pending = studioV2RuntimeContextAPI.getRuleEvidence(scope);
    scope.workspace_revision = 'W2'; scope.device_id = 'E';
    finish({ data: { success: true, data: evidenceContext() } });
    expect(await pending).toEqual({ ok: true, value: ruleEvidenceFixture() });
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/runtime-context', { params: evidenceScope, signal: undefined });
  });
  it('sends only explicit scope fields, never extra caller keys', async () => {
    reply(evidenceContext());
    const extraScope = { ...evidenceScope, secret: '/private/secret' };
    await studioV2RuntimeContextAPI.getRuleEvidence(extraScope);
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/runtime-context', { params: evidenceScope, signal: undefined });
  });
});
