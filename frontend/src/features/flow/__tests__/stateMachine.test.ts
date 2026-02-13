import { describe, expect, it } from 'vitest';
import { INITIAL_FLOW_STATE, flowReducer } from '../stateMachine';

describe('flowReducer', () => {
  it('sets source and keeps draft status', () => {
    const next = flowReducer(INITIAL_FLOW_STATE, {
      type: 'set_source',
      payload: { sourceDeviceId: 'd1', sourceAddress: 'D100', pointId: 'p1' },
    });

    expect(next.status).toBe('draft');
    expect(next.sourceDeviceId).toBe('d1');
    expect(next.sourceAddress).toBe('D100');
    expect(next.pointId).toBe('p1');
  });

  it('validates when source/address/tag are present', () => {
    const withData = flowReducer(
      flowReducer(
        flowReducer(INITIAL_FLOW_STATE, {
          type: 'set_source',
          payload: { sourceDeviceId: 'd1', sourceAddress: 'D100', pointId: 'p1' },
        }),
        { type: 'set_tag', payload: { tagId: 't1' } }
      ),
      { type: 'mark_validated' }
    );

    expect(withData.status).toBe('validated');
  });

  it('does not activate unless validated', () => {
    const next = flowReducer(INITIAL_FLOW_STATE, { type: 'mark_active' });
    expect(next.status).toBe('draft');
  });

  it('transitions validated to active', () => {
    const validated = flowReducer(
      flowReducer(
        flowReducer(INITIAL_FLOW_STATE, {
          type: 'set_source',
          payload: { sourceDeviceId: 'd1', sourceAddress: 'D100', pointId: 'p1' },
        }),
        { type: 'set_tag', payload: { tagId: 't1' } }
      ),
      { type: 'mark_validated' }
    );
    const active = flowReducer(validated, { type: 'mark_active' });
    expect(active.status).toBe('active');
  });

  it('marks segment error and updates status', () => {
    const next = flowReducer(INITIAL_FLOW_STATE, {
      type: 'mark_error',
      payload: { segment: 'sink', message: 'write failed' },
    });

    expect(next.status).toBe('error');
    expect(next.diagnostics.sink.quality).toBe('bad');
    expect(next.diagnostics.sink.error).toBe('write failed');
  });
});
