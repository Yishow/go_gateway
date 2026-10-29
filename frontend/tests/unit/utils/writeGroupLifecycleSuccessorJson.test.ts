import { describe, expect, it } from 'vitest';
import { parseWriteGroupLifecycle } from '@/utils/studioV2WriteGroupLifecycleJson';
import { parseWriteGroupSuccessorEnvelope, parseWriteGroupSuccessorActivationEnvelope, parseWriteGroupSuccessorActivationReadiness, parseWriteGroupSuccessorOperationStatus } from '@/utils/studioV2WriteGroupSuccessorJson';
import { parseStudioV2WriteGroup } from '@/utils/studioV2WriteGroupJson';
import { activation, candidate, envelope, history, lifecycle, operationStatus } from '../../fixtures/writeGroupSuccessor';

describe('strict lifecycle and successor contracts', () => {
  it('keeps samples separate from outbox states and accepts the sanitized advisory seed', () => {
    expect(parseWriteGroupLifecycle(lifecycle())).toEqual(lifecycle());
    expect(parseStudioV2WriteGroup(candidate())).toBeNull();
  });
  it('parses lowercase handoff wire identities without adding them to row counts', () => {
    const value = lifecycle();
    const key = { workspace_id: 'ws-1', group_id: 'wide-1', group_revision: 'applied-wide' };
    value.history.generations[0].revisions[0].pending.handoffs = [{ origin: key, source: key, target: { ...key, group_id: 'new-2' }, state: 'accepted', samples: 2 }];
    expect(parseWriteGroupLifecycle(value)).toEqual(value);
  });
  it('preserves the complete original operation envelope including blank candidate identities', () => {
    expect(parseWriteGroupSuccessorEnvelope(envelope())).toEqual(envelope());
  });
  it('binds activation v2 to its kind and permits only kindless v1 cutover', () => {
    expect(parseWriteGroupSuccessorActivationEnvelope(activation())).toEqual(activation());
    const { activation_kind: _kind, ...legacy } = activation();
    expect(parseWriteGroupSuccessorActivationEnvelope({ ...legacy, version: 'successor-activation-review-v1' })).not.toBeNull();
    expect(parseWriteGroupSuccessorActivationEnvelope({ ...legacy, version: 'successor-activation-review-v2' })).toBeNull();
    expect(parseWriteGroupSuccessorActivationEnvelope({ ...activation(), version: 'successor-activation-review-v1' })).toBeNull();
  });
  it('retains a blocked readiness without inventing a readiness digest', () => {
    const value = { review: { envelope: activation(), review_digest: 'a'.repeat(64), history: history(), blockers: [] },
      readiness: { workspace_id: 'ws-1', workspace_revision: 'ws-r3', group_id: 'new-2', group_revision: 'g-r2', applied_revision: '', config_ready: false, schema_ready: false, ready: false, issues: [] }, readiness_digest: '' };
    expect(parseWriteGroupSuccessorActivationReadiness(value)).toEqual(value);
    expect(parseWriteGroupSuccessorActivationReadiness({ ...value, readiness_digest: 'not-a-digest' })).toBeNull();
  });
  it('accepts only recorded operation outcomes, never a latest result or a bare success', () => {
    expect(parseWriteGroupSuccessorOperationStatus(operationStatus())).toEqual(operationStatus());
    expect(parseWriteGroupSuccessorOperationStatus({ ...operationStatus(), reserve: { status: 'succeeded' } })).toBeNull();
  });
  it('fails closed for malformed numeric units, unknown actions, foreign seed and malformed nested evidence', () => {
    const bad = lifecycle(); bad.history.generations[0].revisions[0].pending.samples = -1;
    expect(parseWriteGroupLifecycle(bad)).toBeNull();
    expect(parseWriteGroupLifecycle({ ...lifecycle(), actions: [{ ...lifecycle().actions[0], action: 'start_anything' }] })).toBeNull();
    const foreign = lifecycle(); if (foreign.successor_seed) foreign.successor_seed.group.workspace_id = 'foreign';
    expect(parseWriteGroupLifecycle(foreign)).toBeNull();
    const badEnvelope = envelope(); badEnvelope.request.group.members[0].required = 'yes' as unknown as boolean;
    expect(parseWriteGroupSuccessorEnvelope(badEnvelope)).toBeNull();
  });
});
