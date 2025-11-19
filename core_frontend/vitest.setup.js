import { vi, beforeEach } from 'vitest';

// Prevent Pinia from initializing devtools
vi.stubGlobal('__VUE_DEVTOOLS_GLOBAL_HOOK__', {
  on: vi.fn(),
  emit: vi.fn(),
});

// Provide a working Storage-like mock that behaves like DOM localStorage.
// Tests can rely on `globalThis.localStorage` and set items before dynamically
// importing modules that read storage at module-init time.
const originalLocalStorage = globalThis.localStorage;

class MockStorage {
  constructor(initial = {}) {
    this._store = new Map(Object.entries(initial));
  }
  _store = new Map();
  get length() {
    return this._store.size;
  }
  clear() {
    this._store.clear();
  }
  getItem(key) {
    return this._store.has(key) ? (this._store.get(key) ?? null) : null;
  }
  key(index) {
    return Array.from(this._store.keys())[index] ?? null;
  }
  removeItem(key) {
    this._store.delete(key);
  }
  setItem(key, value) {
    this._store.set(String(key), String(value));
  }
}

function createMockLocalStorage(initial = {}) {
  return new MockStorage(initial);
}

// Install a fresh mock before the test suite starts and reset before each test.
globalThis.localStorage = createMockLocalStorage();

beforeEach(() => {
  // reset to a fresh instance so tests start with empty storage by default
  globalThis.localStorage = createMockLocalStorage();
});

// expose helpers for tests that want to pre-install values before dynamic import
globalThis.__createMockLocalStorage = createMockLocalStorage;
globalThis.__originalLocalStorage = originalLocalStorage;
