/** 預設每列格數（與 16-bit 列慣例一致）。 */
export const LATTICE_COLUMNS_DEFAULT = 16;

/** 使用者可設定的每列格數下限（避免欄過少無法排版）。 */
export const LATTICE_COLUMNS_MIN = 4;

/** 使用者可設定的每列格數上限。 */
export const LATTICE_COLUMNS_MAX = 32;

/** localStorage 鍵：來源畫布每列顯示幾格。 */
export const SOURCE_CANVAS_LATTICE_STORAGE_KEY =
  'datalink.workbench.sourceCanvasLatticeColumns';

/**
 * 將輸入限制在允許範圍內；無效時回傳 {@link LATTICE_COLUMNS_DEFAULT}。
 */
export function clampLatticeColumns(value: number): number {
  if (!Number.isFinite(value)) {
    return LATTICE_COLUMNS_DEFAULT;
  }
  const n = Math.floor(value);
  if (n < LATTICE_COLUMNS_MIN) {
    return LATTICE_COLUMNS_MIN;
  }
  if (n > LATTICE_COLUMNS_MAX) {
    return LATTICE_COLUMNS_MAX;
  }
  return n;
}

/**
 * 讀取本機儲存的每列格數；缺漏或無效時回傳預設 16。
 */
export function readSourceCanvasLatticeColumns(): number {
  if (typeof window === 'undefined' || typeof localStorage === 'undefined') {
    return LATTICE_COLUMNS_DEFAULT;
  }
  try {
    const raw = localStorage.getItem(SOURCE_CANVAS_LATTICE_STORAGE_KEY);
    if (raw == null || raw === '') {
      return LATTICE_COLUMNS_DEFAULT;
    }
    const parsed = parseInt(raw, 10);
    if (!Number.isFinite(parsed)) {
      return LATTICE_COLUMNS_DEFAULT;
    }
    return clampLatticeColumns(parsed);
  } catch {
    return LATTICE_COLUMNS_DEFAULT;
  }
}

/**
 * 寫入本機儲存（已 clamp）。
 */
export function writeSourceCanvasLatticeColumns(value: number): void {
  if (typeof localStorage === 'undefined') {
    return;
  }
  try {
    const next = clampLatticeColumns(value);
    if (next === LATTICE_COLUMNS_DEFAULT) {
      localStorage.removeItem(SOURCE_CANVAS_LATTICE_STORAGE_KEY);
    } else {
      localStorage.setItem(SOURCE_CANVAS_LATTICE_STORAGE_KEY, String(next));
    }
  } catch {
    /* 私密模式等略過 */
  }
}
