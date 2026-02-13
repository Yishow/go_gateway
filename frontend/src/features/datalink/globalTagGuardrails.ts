export interface GlobalTagGuardrail {
  affectedMappings: number;
  requiresConfirmation: boolean;
  warningMessage: string;
}

export function buildGlobalTagGuardrail(affectedMappings: number): GlobalTagGuardrail {
  const safeCount = Number.isFinite(affectedMappings) ? Math.max(0, Math.floor(affectedMappings)) : 0;
  return {
    affectedMappings: safeCount,
    requiresConfirmation: safeCount > 0,
    warningMessage:
      safeCount > 0
        ? `請確認差異後再儲存，預估影響 ${safeCount} 個映射。`
        : '此變更不會影響既有映射，可直接儲存。',
  };
}
