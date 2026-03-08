import { describe, expect, it } from 'vitest';
import { runStructuralValidation } from '@/features/datalink/validationFlow';

describe('validationFlow', () => {
  it('fails when source segment is incomplete', () => {
    const result = runStructuralValidation({
      selectedDeviceId: null,
      selectedSourceAddress: '',
      planAddressesCount: 3,
      planConflictCount: 0,
      hasSelectedMapping: true,
    });
    expect(result).toEqual({ ok: false, error: '來源段尚未完成' });
  });

  it('fails when structural conflicts remain', () => {
    const result = runStructuralValidation({
      selectedDeviceId: 'dev-1',
      selectedSourceAddress: '40001',
      planAddressesCount: 4,
      planConflictCount: 2,
      hasSelectedMapping: true,
    });
    expect(result).toEqual({ ok: false, error: '格位仍有 2 筆衝突' });
  });

  it('passes when structural checks are complete', () => {
    const result = runStructuralValidation({
      selectedDeviceId: 'dev-1',
      selectedSourceAddress: '40001',
      planAddressesCount: 4,
      planConflictCount: 0,
      hasSelectedMapping: true,
    });
    expect(result).toEqual({ ok: true });
  });
});
