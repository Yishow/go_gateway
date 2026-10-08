import '@testing-library/jest-dom';

// Node 22+ ships an experimental global `localStorage` that is undefined unless
// the process is started with `--localstorage-file`. During jsdom global
// population vitest skips keys already present on `globalThis`, so both bare
// `localStorage` and `window.localStorage` come up undefined and test
// setup/teardown crashes with "Cannot read properties of undefined". Install an
// in-memory Storage whenever either surface is missing, so tests exercise real
// browser storage semantics instead of depending on the Node runtime flavour.
type MemoryStorage = Pick<
  Storage,
  'length' | 'clear' | 'getItem' | 'key' | 'removeItem' | 'setItem'
>;

function createMemoryStorage(): MemoryStorage {
  const map = new Map<string, string>();
  return {
    get length() {
      return map.size;
    },
    clear: () => {
      map.clear();
    },
    getItem: (key: string) => (map.has(key) ? (map.get(key) as string) : null),
    key: (index: number) => Array.from(map.keys())[index] ?? null,
    removeItem: (key: string) => {
      map.delete(key);
    },
    setItem: (key: string, value: string) => {
      map.set(key, String(value));
    },
  };
}

function ensureStorage(target: object, instance: MemoryStorage): void {
  if (typeof (target as { localStorage?: unknown }).localStorage !== 'undefined') {
    return;
  }
  try {
    Object.defineProperty(target, 'localStorage', {
      value: instance,
      configurable: true,
      writable: true,
    });
  } catch {
    // Non-configurable native accessor: leave it alone rather than break tests.
  }
}

// One shared instance: `window.localStorage` and bare `localStorage` must be
// the same storage object, as they are in browsers.
const memoryStorage = createMemoryStorage();
ensureStorage(globalThis, memoryStorage);
if (typeof window !== 'undefined') {
  ensureStorage(window, memoryStorage);
}
