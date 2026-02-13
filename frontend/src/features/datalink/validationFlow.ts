export interface StructuralValidationInput {
  selectedDeviceId: string | null;
  selectedSourceAddress: string;
  planAddressesCount: number;
  planConflictCount: number;
  hasSelectedMapping: boolean;
}

export function runStructuralValidation(input: StructuralValidationInput): {
  ok: boolean;
  error?: string;
} {
  if (!input.selectedDeviceId || !input.selectedSourceAddress) {
    return { ok: false, error: '來源段尚未完成' };
  }
  if (input.planAddressesCount === 0) {
    return { ok: false, error: '尚未完成格位規劃' };
  }
  if (input.planConflictCount > 0) {
    return { ok: false, error: `格位仍有 ${input.planConflictCount} 筆衝突` };
  }
  if (!input.hasSelectedMapping) {
    return { ok: false, error: 'Tag 段尚未完成' };
  }
  return { ok: true };
}
