import { describe, expect, it } from 'vitest';
import { resolveGatewayUiVersion } from '@/features/gateway/uiVersion';

describe('resolveGatewayUiVersion', () => {
  it('returns dual_quick for quick setup path', () => {
    expect(resolveGatewayUiVersion('/gateway/quick-setup')).toBe('dual_quick');
  });

  it('returns dual_expert for expert workbench path', () => {
    expect(resolveGatewayUiVersion('/gateway/expert-workbench')).toBe(
      'dual_expert'
    );
  });

  it('returns null for non-dual pages', () => {
    expect(resolveGatewayUiVersion('/datalink')).toBeNull();
    expect(resolveGatewayUiVersion('/gateway/entry')).toBeNull();
  });
});
