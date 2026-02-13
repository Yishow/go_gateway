export interface BatchNamePreviewItem {
  name: string;
  sequence: number;
  conflict: boolean;
}

export function buildBatchNamePreview(
  prefix: string,
  count: number,
  existingNames: string[]
): BatchNamePreviewItem[] {
  const normalizedPrefix = prefix.trim() || 'SRC';
  const normalizedExisting = new Set(existingNames.map((item) => item.trim().toLowerCase()));
  const size = Number.isFinite(count) ? Math.max(0, Math.floor(count)) : 0;

  return Array.from({ length: size }).map((_, index) => {
    const sequence = index + 1;
    const name = `${normalizedPrefix}-${String(sequence).padStart(3, '0')}`;
    return {
      name,
      sequence,
      conflict: normalizedExisting.has(name.toLowerCase()),
    };
  });
}
