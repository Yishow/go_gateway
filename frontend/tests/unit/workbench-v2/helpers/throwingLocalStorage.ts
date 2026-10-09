/**
 * 暫時以委派式 Storage 取代 window.localStorage，使指定 key 的 getItem 拋錯
 * （模擬隱私模式的 SecurityError），其餘操作與 key 仍委派原 storage。回傳還原函式。
 *
 * 不能用 vi.spyOn(localStorage, 'getItem')：jsdom 原生 Storage 會把屬性寫入
 * 當成資料項目，spy 不會生效；而 setupTests 的 shim 又不是 Storage 實例，
 * 所以 spyOn(Storage.prototype) 也無法涵蓋。替換 window.localStorage 兩種環境都適用。
 */
export function stubLocalStorageGetItemThrows(
  blockedKey: string,
  message = 'SecurityError: The operation is insecure.',
): () => void {
  const realStorage = window.localStorage;
  const originalDescriptor = Object.getOwnPropertyDescriptor(window, 'localStorage');
  const throwingStorage: Storage = {
    get length() {
      return realStorage.length;
    },
    clear: () => realStorage.clear(),
    getItem: (key: string) => {
      if (key === blockedKey) throw new Error(message);
      return realStorage.getItem(key);
    },
    key: (index: number) => realStorage.key(index),
    removeItem: (key: string) => realStorage.removeItem(key),
    setItem: (key: string, value: string) => realStorage.setItem(key, value),
  };
  Object.defineProperty(window, 'localStorage', { configurable: true, value: throwingStorage });

  return () => {
    if (originalDescriptor) {
      Object.defineProperty(window, 'localStorage', originalDescriptor);
    } else {
      delete (window as { localStorage?: Storage }).localStorage;
    }
  };
}
