import { vi } from 'vitest';

/**
 * Studio V2 service mock registry（僅依賴 vitest，不可 import 任何 src 模組）。
 *
 * 各 Workbench V2 測試檔在 vi.mock factory 內以 dynamic import 取用此
 * registry（vi.mock 會 hoist，不可直接引用測試檔頂層 binding），確保五個
 * service 的 mock 形狀跨檔案一致；beforeEach 的 vi.clearAllMocks() 會一併
 * 重置這些 mocks。獨立成小模組是為了阻斷「factory → harness → 頁面 →
 * 被 mock 的 service → factory」的循環載入死結。
 */
export const studioV2ServiceMocks = {
  workspace: { get: vi.fn() },
  devices: {
    list: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    updateAvailability: vi.fn(),
    remove: vi.fn(),
    updateOrder: vi.fn(),
  },
  rules: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() },
  mappings: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() },
  database: {
    getConfig: vi.fn(),
    updateConfig: vi.fn(),
    listTargets: vi.fn(),
    upsertTarget: vi.fn(),
  },
};
