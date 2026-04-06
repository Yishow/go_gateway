type TranslationTree = Record<string, unknown>;

function isPlainObject(value: unknown): value is TranslationTree {
  return Boolean(value) && typeof value === 'object' && !Array.isArray(value);
}

export function mergeTranslations<TBase extends TranslationTree>(
  base: TBase,
  ...overlays: TranslationTree[]
): TBase {
  const result: TranslationTree = {};

  for (const source of [base, ...overlays]) {
    for (const [key, value] of Object.entries(source)) {
      const currentValue = result[key];
      if (isPlainObject(currentValue) && isPlainObject(value)) {
        result[key] = mergeTranslations(currentValue, value);
        continue;
      }
      result[key] = value;
    }
  }

  return result as TBase;
}
